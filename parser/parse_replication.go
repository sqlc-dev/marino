// Copyright 2026 The sqlc Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package parser

// The MySQL replication statements (MySQL 26.7 §15.4): CHANGE
// REPLICATION SOURCE TO, CHANGE REPLICATION FILTER, PURGE BINARY LOGS,
// RESET BINARY LOGS AND GTIDS, RESET REPLICA, START/STOP REPLICA, and
// START/STOP GROUP_REPLICATION. All of these postdate the goyacc
// grammar; the productions are written from the MySQL 26.7 reference
// manual in the same style as parser.y. The removed pre-8.4 spellings
// (CHANGE MASTER TO, PURGE MASTER LOGS, RESET MASTER, RESET/START/STOP
// SLAVE) are not parsed.
//
// This file owns the START, STOP, and PURGE statement heads and routes
// their non-replication alternatives onward: START to
// BeginTransactionStmt (parse_txn.go), STOP and PURGE to the BRIE
// statement family (parse_brie.go). The RESET head stays with
// parse_mysql_admin.go, which routes RESET REPLICA and RESET BINARY
// LOGS AND GTIDS here.

import (
	"strings"

	"github.com/sqlc-dev/marino/ast"
)

func init() {
	rdRegister(change, (*rdParser).parseChangeStmtFamily)
	rdRegister(start, (*rdParser).parseStartStmtFamily)
	rdRegister(stop, (*rdParser).parseStopStmtFamily)
	rdRegister(purge, (*rdParser).parsePurgeStmtFamily)
}

// parseChangeStmtFamily dispatches the CHANGE statement family:
// ChangeReplicationSourceStmt | ChangeReplicationFilterStmt.
func (r *rdParser) parseChangeStmtFamily() ast.StmtNode {
	if r.la(2) == filter {
		return r.parseChangeReplicationFilterStmt()
	}
	return r.parseChangeReplicationSourceStmt()
}

// parseStartStmtFamily dispatches the START statement family on the
// token after "START": StartReplicaStmt, StartGroupReplicationStmt, or
// the "START TRANSACTION" alternative of BeginTransactionStmt.
func (r *rdParser) parseStartStmtFamily() ast.StmtNode {
	switch r.la(1) {
	case replica:
		return r.parseStartReplicaStmt()
	case groupReplication:
		return r.parseStartGroupReplicationStmt()
	}
	return r.parseBeginTransactionStmt()
}

// parseStopStmtFamily dispatches the STOP statement family on the token
// after "STOP": StopReplicaStmt, StopGroupReplicationStmt, or the
// "STOP BACKUP LOGS" BRIEStmt.
func (r *rdParser) parseStopStmtFamily() ast.StmtNode {
	switch r.la(1) {
	case replica:
		return r.parseStopReplicaStmt()
	case groupReplication:
		// StopGroupReplicationStmt: "STOP" "GROUP_REPLICATION"
		r.advance()
		r.advance()
		return &ast.StopGroupReplicationStmt{}
	}
	return r.parseStopBackupStmt()
}

// parsePurgeStmtFamily dispatches the PURGE statement family on the
// token after "PURGE": PurgeBinaryLogsStmt or the "PURGE BACKUP LOGS"
// BRIEStmt.
func (r *rdParser) parsePurgeStmtFamily() ast.StmtNode {
	if r.la(1) == binaryType {
		return r.parsePurgeBinaryLogsStmt()
	}
	return r.parsePurgeBackupStmt()
}

// parseChangeReplicationSourceStmt implements ChangeReplicationSourceStmt:
//
//	"CHANGE" "REPLICATION" "SOURCE" "TO" ReplicationSourceOption
//	    ("," ReplicationSourceOption)* ForChannelOpt
//
// Option names are not validated against the server's option list, so the
// MySQL 26.7 Change Stream Applier options (APPLIER_VERSION,
// APPLIER_WORKER_COUNT, APPLIER_EVENT_MEMORY_LIMIT) parse like any other.
func (r *rdParser) parseChangeReplicationSourceStmt() ast.StmtNode {
	r.expect(change)
	r.expect(replication)
	r.expect(source)
	r.expect(to)
	stmt := &ast.ChangeReplicationSourceStmt{
		Options: []*ast.ReplicationSourceOption{r.parseReplicationSourceOption()},
	}
	for r.accept(int(',')) {
		stmt.Options = append(stmt.Options, r.parseReplicationSourceOption())
	}
	stmt.Channel = r.parseForChannelOpt()
	return stmt
}

// parseReplicationSourceOption implements ReplicationSourceOption:
// Identifier eq (stringLit | intLit | decLit | floatLit), plus the
// non-literal value forms selected by the option name: a parenthesized
// server-id list (IGNORE_SERVER_IDS), an account name or NULL
// (PRIVILEGE_CHECKS_USER), and bare keyword values
// (REQUIRE_TABLE_PRIMARY_KEY_CHECK = STREAM|GENERATE|ON|OFF,
// ASSIGN_GTIDS_TO_ANONYMOUS_TRANSACTIONS = OFF|LOCAL|uuid).
func (r *rdParser) parseReplicationSourceOption() *ast.ReplicationSourceOption {
	name := strings.ToUpper(r.parseIdentifier())
	r.expect(eq)
	opt := &ast.ReplicationSourceOption{Name: name}
	switch name {
	case "IGNORE_SERVER_IDS":
		r.expect(int('('))
		opt.ServerIDs = []ast.ValueExpr{}
		if r.tok() != int(')') {
			opt.ServerIDs = append(opt.ServerIDs, r.parseReplicationOptionValue())
			for r.accept(int(',')) {
				opt.ServerIDs = append(opt.ServerIDs, r.parseReplicationOptionValue())
			}
		}
		r.expect(int(')'))
	case "PRIVILEGE_CHECKS_USER":
		if r.accept(null) {
			opt.KeywordValue = "NULL"
		} else {
			opt.User = r.parseGrantUsername()
		}
	case "REQUIRE_TABLE_PRIMARY_KEY_CHECK", "ASSIGN_GTIDS_TO_ANONYMOUS_TRANSACTIONS":
		if r.accept(on) {
			opt.KeywordValue = "ON"
		} else if r.tok() == stringLit {
			// ASSIGN_GTIDS_TO_ANONYMOUS_TRANSACTIONS = 'uuid'
			opt.Value = r.parseReplicationOptionValue()
		} else {
			opt.KeywordValue = strings.ToUpper(r.parseIdentifier())
		}
	default:
		opt.Value = r.parseReplicationOptionValue()
	}
	return opt
}

// parseReplicationOptionValue implements the literal value of a
// replication option: stringLit | intLit | decLit | floatLit.
func (r *rdParser) parseReplicationOptionValue() ast.ValueExpr {
	var value ast.ValueExpr
	switch r.tok() {
	case stringLit:
		value = ast.NewValueExpr(r.cur().lit, "", "")
		r.advance()
	case intLit, decLit, floatLit:
		value = ast.NewValueExpr(r.cur().item, "", "")
		r.advance()
	default:
		r.syntaxError()
	}
	return value
}

// parseForChannelOpt implements ForChannelOpt:
// empty | "FOR" "CHANNEL" stringLit. It returns "" for the empty
// alternative.
func (r *rdParser) parseForChannelOpt() string {
	if !r.accept(forKwd) {
		return ""
	}
	r.expect(channel)
	return r.expect(stringLit).lit
}

// parseChangeReplicationFilterStmt implements ChangeReplicationFilterStmt:
//
//	"CHANGE" "REPLICATION" "FILTER" ReplicationFilter
//	    ("," ReplicationFilter)* ForChannelOpt
func (r *rdParser) parseChangeReplicationFilterStmt() ast.StmtNode {
	r.expect(change)
	r.expect(replication)
	r.expect(filter)
	stmt := &ast.ChangeReplicationFilterStmt{
		Filters: []*ast.ReplicationFilter{r.parseReplicationFilter()},
	}
	for r.accept(int(',')) {
		stmt.Filters = append(stmt.Filters, r.parseReplicationFilter())
	}
	stmt.Channel = r.parseForChannelOpt()
	return stmt
}

// parseReplicationFilter implements ReplicationFilter. The filter names
// are not keywords (non-reserved in MySQL) and are matched
// case-insensitively in identifier position; the name selects the value
// form:
//
//	{"REPLICATE_DO_DB" | "REPLICATE_IGNORE_DB"} eq '(' DBNameListOpt ')'
//	| {"REPLICATE_DO_TABLE" | "REPLICATE_IGNORE_TABLE"} eq
//	    '(' TableNameListOpt ')'
//	| {"REPLICATE_WILD_DO_TABLE" | "REPLICATE_WILD_IGNORE_TABLE"} eq
//	    '(' StringListOpt ')'
//	| "REPLICATE_REWRITE_DB" eq '(' DBPairListOpt ')'
//
// An empty '(' ')' value clears the filter.
func (r *rdParser) parseReplicationFilter() *ast.ReplicationFilter {
	if !isIdentifierTok(r.tok()) {
		r.syntaxError()
	}
	var tp ast.ReplicationFilterType
	switch strings.ToUpper(r.cur().lit) {
	case "REPLICATE_DO_DB":
		tp = ast.ReplicationFilterDoDB
	case "REPLICATE_IGNORE_DB":
		tp = ast.ReplicationFilterIgnoreDB
	case "REPLICATE_DO_TABLE":
		tp = ast.ReplicationFilterDoTable
	case "REPLICATE_IGNORE_TABLE":
		tp = ast.ReplicationFilterIgnoreTable
	case "REPLICATE_WILD_DO_TABLE":
		tp = ast.ReplicationFilterWildDoTable
	case "REPLICATE_WILD_IGNORE_TABLE":
		tp = ast.ReplicationFilterWildIgnoreTable
	case "REPLICATE_REWRITE_DB":
		tp = ast.ReplicationFilterRewriteDB
	default:
		r.syntaxError()
	}
	r.advance()
	r.expect(eq)
	r.expect(int('('))
	f := &ast.ReplicationFilter{Tp: tp}
	if r.tok() != int(')') {
		switch tp {
		case ast.ReplicationFilterDoDB, ast.ReplicationFilterIgnoreDB:
			f.DBNames = []ast.CIStr{ast.NewCIStr(r.parseIdentifier())}
			for r.accept(int(',')) {
				f.DBNames = append(f.DBNames, ast.NewCIStr(r.parseIdentifier()))
			}
		case ast.ReplicationFilterDoTable, ast.ReplicationFilterIgnoreTable:
			f.Tables = []*ast.TableName{r.parseTableName()}
			for r.accept(int(',')) {
				f.Tables = append(f.Tables, r.parseTableName())
			}
		case ast.ReplicationFilterWildDoTable, ast.ReplicationFilterWildIgnoreTable:
			f.Patterns = []string{r.expect(stringLit).lit}
			for r.accept(int(',')) {
				f.Patterns = append(f.Patterns, r.expect(stringLit).lit)
			}
		case ast.ReplicationFilterRewriteDB:
			f.Rewrites = []*ast.ReplicationRewriteDB{r.parseReplicationRewriteDB()}
			for r.accept(int(',')) {
				f.Rewrites = append(f.Rewrites, r.parseReplicationRewriteDB())
			}
		}
	}
	r.expect(int(')'))
	return f
}

// parseReplicationRewriteDB implements the (from_db, to_db) pair of a
// REPLICATE_REWRITE_DB filter.
func (r *rdParser) parseReplicationRewriteDB() *ast.ReplicationRewriteDB {
	r.expect(int('('))
	pair := &ast.ReplicationRewriteDB{FromDB: ast.NewCIStr(r.parseIdentifier())}
	r.expect(int(','))
	pair.ToDB = ast.NewCIStr(r.parseIdentifier())
	r.expect(int(')'))
	return pair
}

// parsePurgeBinaryLogsStmt implements PurgeBinaryLogsStmt:
// "PURGE" "BINARY" "LOGS" ("TO" stringLit | "BEFORE" Expression).
func (r *rdParser) parsePurgeBinaryLogsStmt() ast.StmtNode {
	r.expect(purge)
	r.expect(binaryType)
	r.expect(logs)
	switch r.tok() {
	case to:
		r.advance()
		return &ast.PurgeBinaryLogsStmt{To: r.expect(stringLit).lit}
	case before:
		r.advance()
		return &ast.PurgeBinaryLogsStmt{Before: r.parseExpression()}
	}
	r.syntaxError()
	return nil
}

// parseResetReplicaStmt implements ResetReplicaStmt:
// "RESET" "REPLICA" ["ALL"] ForChannelOpt.
func (r *rdParser) parseResetReplicaStmt() ast.StmtNode {
	r.expect(reset)
	r.expect(replica)
	stmt := &ast.ResetReplicaStmt{All: r.accept(all)}
	stmt.Channel = r.parseForChannelOpt()
	return stmt
}

// parseResetBinaryLogsAndGtidsStmt implements
// ResetBinaryLogsAndGtidsStmt:
// "RESET" "BINARY" "LOGS" "AND" "GTIDS" ["TO" Int64Num].
func (r *rdParser) parseResetBinaryLogsAndGtidsStmt() ast.StmtNode {
	r.expect(reset)
	r.expect(binaryType)
	r.expect(logs)
	r.expect(and)
	r.expect(gtids)
	stmt := &ast.ResetBinaryLogsAndGtidsStmt{}
	if r.accept(to) {
		stmt.To = r.parseInt64Num()
	}
	return stmt
}

// parseReplicaThreadTypes implements the thread_types list of
// START/STOP REPLICA: empty | ReplicaThreadType (',' ReplicaThreadType)*
// with ReplicaThreadType: "IO_THREAD" | "SQL_THREAD".
func (r *rdParser) parseReplicaThreadTypes() []ast.ReplicaThreadType {
	if r.tok() != ioThread && r.tok() != sqlThread {
		return nil
	}
	types := []ast.ReplicaThreadType{r.parseReplicaThreadType()}
	for r.accept(int(',')) {
		types = append(types, r.parseReplicaThreadType())
	}
	return types
}

func (r *rdParser) parseReplicaThreadType() ast.ReplicaThreadType {
	switch r.tok() {
	case ioThread:
		r.advance()
		return ast.ReplicaThreadIO
	case sqlThread:
		r.advance()
		return ast.ReplicaThreadSQL
	}
	r.syntaxError()
	return ast.ReplicaThreadIO
}

// parseStartReplicaStmt implements StartReplicaStmt:
//
//	"START" "REPLICA" ReplicaThreadTypes ["UNTIL" ReplicaUntilOption
//	    ("," ReplicaUntilOption)*] ReplicaConnectionOptions ForChannelOpt
//
// A ReplicaUntilOption is a name eq value pair (SQL_BEFORE_GTIDS,
// SOURCE_LOG_FILE, RELAY_LOG_POS, ...) or the bare SQL_AFTER_MTS_GAPS;
// like the CHANGE REPLICATION SOURCE TO options the names are matched
// in identifier position and not validated against the server's list.
func (r *rdParser) parseStartReplicaStmt() ast.StmtNode {
	r.expect(start)
	r.expect(replica)
	stmt := &ast.StartReplicaStmt{ThreadTypes: r.parseReplicaThreadTypes()}
	if r.accept(until) {
		stmt.Until = []*ast.ReplicationSourceOption{r.parseReplicaUntilOption()}
		for r.accept(int(',')) {
			stmt.Until = append(stmt.Until, r.parseReplicaUntilOption())
		}
	}
	stmt.ConnectionOptions = r.parseReplicaConnectionOptions()
	stmt.Channel = r.parseForChannelOpt()
	return stmt
}

// parseReplicaUntilOption implements ReplicaUntilOption:
// Identifier [eq (stringLit | intLit | decLit | floatLit)].
func (r *rdParser) parseReplicaUntilOption() *ast.ReplicationSourceOption {
	name := strings.ToUpper(r.parseIdentifier())
	opt := &ast.ReplicationSourceOption{Name: name}
	if r.accept(eq) {
		opt.Value = r.parseReplicationOptionValue()
	}
	return opt
}

// parseReplicaConnectionOptions implements the space-separated
// connection options of START REPLICA: ("USER" | "PASSWORD" |
// "DEFAULT_AUTH" | "PLUGIN_DIR") eq stringLit, in any order.
// DEFAULT_AUTH and PLUGIN_DIR are matched in identifier position.
func (r *rdParser) parseReplicaConnectionOptions() []*ast.ReplicationSourceOption {
	var opts []*ast.ReplicationSourceOption
	for {
		var name string
		switch {
		case r.tok() == user || r.tok() == password:
			name = strings.ToUpper(r.cur().lit)
		case isIdentifierTok(r.tok()):
			name = strings.ToUpper(r.cur().lit)
			if name != "DEFAULT_AUTH" && name != "PLUGIN_DIR" {
				return opts
			}
		default:
			return opts
		}
		r.advance()
		r.expect(eq)
		opts = append(opts, &ast.ReplicationSourceOption{
			Name:  name,
			Value: ast.NewValueExpr(r.expect(stringLit).lit, "", ""),
		})
	}
}

// parseStopReplicaStmt implements StopReplicaStmt:
// "STOP" "REPLICA" ReplicaThreadTypes ForChannelOpt.
func (r *rdParser) parseStopReplicaStmt() ast.StmtNode {
	r.expect(stop)
	r.expect(replica)
	stmt := &ast.StopReplicaStmt{ThreadTypes: r.parseReplicaThreadTypes()}
	stmt.Channel = r.parseForChannelOpt()
	return stmt
}

// parseStartGroupReplicationStmt implements StartGroupReplicationStmt:
//
//	"START" "GROUP_REPLICATION" [GroupReplicationOption
//	    ("," GroupReplicationOption)*]
//
// with GroupReplicationOption: ("USER" | "PASSWORD" | "DEFAULT_AUTH")
// eq stringLit; DEFAULT_AUTH is matched in identifier position.
func (r *rdParser) parseStartGroupReplicationStmt() ast.StmtNode {
	r.expect(start)
	r.expect(groupReplication)
	stmt := &ast.StartGroupReplicationStmt{}
	if r.tok() != user && r.tok() != password &&
		!(isIdentifierTok(r.tok()) && strings.EqualFold(r.cur().lit, "DEFAULT_AUTH")) {
		return stmt
	}
	stmt.Options = []*ast.ReplicationSourceOption{r.parseGroupReplicationOption()}
	for r.accept(int(',')) {
		stmt.Options = append(stmt.Options, r.parseGroupReplicationOption())
	}
	return stmt
}

func (r *rdParser) parseGroupReplicationOption() *ast.ReplicationSourceOption {
	switch {
	case r.tok() == user || r.tok() == password:
	case isIdentifierTok(r.tok()) && strings.EqualFold(r.cur().lit, "DEFAULT_AUTH"):
	default:
		r.syntaxError()
	}
	name := strings.ToUpper(r.cur().lit)
	r.advance()
	r.expect(eq)
	return &ast.ReplicationSourceOption{
		Name:  name,
		Value: ast.NewValueExpr(r.expect(stringLit).lit, "", ""),
	}
}
