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

// Miscellaneous statements: USE, KILL, TRUNCATE, FLUSH, HELP, SHUTDOWN,
// RESTART, LOCK/UNLOCK TABLES, BINLOG, PREPARE/EXECUTE/DEALLOCATE,
// EXPLAIN/DESC/DESCRIBE, and TRACE.

import (
	"strings"

	"github.com/sqlc-dev/marino/ast"
)

func init() {
	rdRegister(use, (*rdParser).parseUseStmt)
	rdRegister(kill, (*rdParser).parseKillStmt)
	rdRegister(truncate, (*rdParser).parseTruncateTableStmt)
	rdRegister(flush, (*rdParser).parseFlushStmt)
	rdRegister(help, (*rdParser).parseHelpStmt)
	rdRegister(shutdown, (*rdParser).parseShutdownStmt)
	rdRegister(restart, (*rdParser).parseRestartStmt)
	rdRegister(lock, (*rdParser).parseLockStmtFamily)
	rdRegister(unlock, (*rdParser).parseUnlockStmtFamily)
	rdRegister(binlog, (*rdParser).parseBinlogStmt)
	rdRegister(prepare, (*rdParser).parsePreparedStmt)
	rdRegister(execute, (*rdParser).parseExecuteStmt)
	rdRegister(deallocate, (*rdParser).parseDeallocateStmt)
	rdRegister(explain, (*rdParser).parseExplainStmt)
	rdRegister(describe, (*rdParser).parseExplainStmt)
	rdRegister(desc, (*rdParser).parseExplainStmt)
	rdRegister(trace, (*rdParser).parseTraceStmt)
}

// parseUseStmt implements UseStmt: "USE" DBName.
func (r *rdParser) parseUseStmt() ast.StmtNode {
	r.expect(use)
	return &ast.UseStmt{DBName: r.parseIdentifier()}
}

// parseKillStmt implements KillStmt and KillOrKillTiDB.
func (r *rdParser) parseKillStmt() ast.StmtNode {
	r.expect(kill)
	// KillOrKillTiDB: "KILL" | "KILL" "TIDB"
	tidbExtension := r.accept(tidb)
	switch r.tok() {
	case connection:
		// KillOrKillTiDB "CONNECTION" NUM
		r.advance()
		return &ast.KillStmt{
			ConnectionID:  getUint64FromNUM(r.expect(intLit).item),
			TiDBExtension: tidbExtension,
		}
	case query:
		// KillOrKillTiDB "QUERY" NUM
		r.advance()
		return &ast.KillStmt{
			ConnectionID:  getUint64FromNUM(r.expect(intLit).item),
			Query:         true,
			TiDBExtension: tidbExtension,
		}
	case intLit:
		// KillOrKillTiDB NUM
		v := r.cur().item
		r.advance()
		return &ast.KillStmt{
			ConnectionID:  getUint64FromNUM(v),
			TiDBExtension: tidbExtension,
		}
	}
	// KillOrKillTiDB BuiltinFunction
	return &ast.KillStmt{
		TiDBExtension: tidbExtension,
		Expr:          r.parseBuiltinFunction(),
	}
}

// parseBuiltinFunction implements BuiltinFunction:
//
//	'(' BuiltinFunction ')' | identifier '(' ')' |
//	identifier '(' ExpressionList ')' | "REPLACE" '(' ExpressionList ')'
func (r *rdParser) parseBuiltinFunction() ast.ExprNode {
	start := r.cur().offset
	if r.accept(int('(')) {
		f := r.parseBuiltinFunction()
		r.expect(int(')'))
		return r.setOrigin(f, start)
	}
	var name string
	requireArgs := false
	switch r.tok() {
	case identifier:
		name = r.cur().lit
	case replace:
		name = r.cur().lit
		requireArgs = true
	default:
		r.syntaxError()
	}
	r.advance()
	r.expect(int('('))
	if !requireArgs && r.accept(int(')')) {
		return r.setOrigin(&ast.FuncCallExpr{
			FnName: ast.NewCIStr(name),
		}, start)
	}
	args := r.parseExpressionList()
	r.expect(int(')'))
	return r.setOrigin(&ast.FuncCallExpr{
		FnName: ast.NewCIStr(name),
		Args:   args,
	}, start)
}

// parseTruncateTableStmt implements TruncateTableStmt: "TRUNCATE" OptTable TableName.
func (r *rdParser) parseTruncateTableStmt() ast.StmtNode {
	r.expect(truncate)
	r.accept(tableKwd)
	return &ast.TruncateTableStmt{Table: r.parseTableName()}
}

// parseFlushStmt implements FlushStmt: "FLUSH" NoWriteToBinLogAliasOpt FlushOption.
func (r *rdParser) parseFlushStmt() ast.StmtNode {
	r.expect(flush)
	// NoWriteToBinLogAliasOpt: empty | "NO_WRITE_TO_BINLOG" | "LOCAL"
	noWrite := false
	if r.tok() == noWriteToBinLog || r.tok() == local {
		noWrite = true
		r.advance()
	}
	st := r.parseFlushOption()
	st.NoWriteToBinLog = noWrite
	return st
}

// parseFlushOption implements FlushOption.
func (r *rdParser) parseFlushOption() *ast.FlushStmt {
	switch r.tok() {
	case privileges:
		r.advance()
		return &ast.FlushStmt{Tp: ast.FlushPrivileges}
	case status:
		r.advance()
		return &ast.FlushStmt{Tp: ast.FlushStatus}
	case tidb:
		// "TIDB" "PLUGINS" PluginNameList
		r.advance()
		r.expect(plugins)
		return &ast.FlushStmt{
			Tp:      ast.FlushTiDBPlugin,
			Plugins: r.parsePluginNameList(),
		}
	case hosts:
		r.advance()
		return &ast.FlushStmt{Tp: ast.FlushHosts}
	case logs:
		// LogTypeOpt "LOGS" with the empty LogTypeOpt.
		r.advance()
		return &ast.FlushStmt{
			Tp:      ast.FlushLogs,
			LogType: ast.LogTypeDefault,
		}
	case binaryType, engine, errorKwd, general, slow:
		// LogTypeOpt "LOGS"
		var logType ast.LogType
		switch r.tok() {
		case binaryType:
			logType = ast.LogTypeBinary
		case engine:
			logType = ast.LogTypeEngine
		case errorKwd:
			logType = ast.LogTypeError
		case general:
			logType = ast.LogTypeGeneral
		case slow:
			logType = ast.LogTypeSlow
		}
		r.advance()
		r.expect(logs)
		return &ast.FlushStmt{
			Tp:      ast.FlushLogs,
			LogType: logType,
		}
	case tableKwd, tables:
		// TableOrTables TableNameListOpt WithReadLockOpt
		r.advance()
		st := &ast.FlushStmt{
			Tp:     ast.FlushTables,
			Tables: []*ast.TableName{},
		}
		if isIdentifierTok(r.tok()) || r.tok() == int('*') {
			st.Tables = r.parseTableNameList()
		}
		if r.accept(with) {
			r.expect(read)
			r.expect(lock)
			st.ReadLock = true
		}
		return st
	case clientErrorsSummary:
		r.advance()
		return &ast.FlushStmt{Tp: ast.FlushClientErrorsSummary}
	case statsDelta:
		// "STATS_DELTA" StatsObjectList ClusterOpt
		r.advance()
		objects := r.parseStatsObjectList()
		return &ast.FlushStmt{
			Tp:           ast.FlushStatsDelta,
			FlushObjects: objects,
			IsCluster:    r.accept(cluster),
		}
	}
	r.syntaxError()
	return nil
}

// parsePluginNameList implements PluginNameList.
func (r *rdParser) parsePluginNameList() []string {
	list := []string{r.parseIdentifier()}
	for r.accept(int(',')) {
		list = append(list, r.parseIdentifier())
	}
	return list
}

// parseTableNameList implements TableNameList.
func (r *rdParser) parseTableNameList() []*ast.TableName {
	list := []*ast.TableName{r.parseTableName()}
	for r.accept(int(',')) {
		list = append(list, r.parseTableName())
	}
	return list
}

// parseStatsObjectList implements StatsObjectList.
func (r *rdParser) parseStatsObjectList() []*ast.StatsObject {
	list := []*ast.StatsObject{r.parseStatsObject()}
	for r.accept(int(',')) {
		list = append(list, r.parseStatsObject())
	}
	return list
}

// parseStatsObject implements StatsObject.
func (r *rdParser) parseStatsObject() *ast.StatsObject {
	if r.tok() == int('*') {
		// '*' '.' '*'
		r.advance()
		r.expect(int('.'))
		r.expect(int('*'))
		return &ast.StatsObject{
			StatsObjectScope: ast.StatsObjectScopeGlobal,
		}
	}
	first := r.parseIdentifier()
	if r.accept(int('.')) {
		if r.accept(int('*')) {
			// Identifier '.' '*'
			return &ast.StatsObject{
				StatsObjectScope: ast.StatsObjectScopeDatabase,
				DBName:           ast.NewCIStr(first),
			}
		}
		// Identifier '.' Identifier
		return &ast.StatsObject{
			StatsObjectScope: ast.StatsObjectScopeTable,
			DBName:           ast.NewCIStr(first),
			TableName:        ast.NewCIStr(r.parseIdentifier()),
		}
	}
	// Identifier
	return &ast.StatsObject{
		StatsObjectScope: ast.StatsObjectScopeTable,
		TableName:        ast.NewCIStr(first),
	}
}

// parseHelpStmt implements HelpStmt: "HELP" stringLit.
func (r *rdParser) parseHelpStmt() ast.StmtNode {
	r.expect(help)
	return &ast.HelpStmt{Topic: r.expect(stringLit).lit}
}

// parseShutdownStmt implements ShutdownStmt: "SHUTDOWN".
func (r *rdParser) parseShutdownStmt() ast.StmtNode {
	r.expect(shutdown)
	return &ast.ShutdownStmt{}
}

// parseRestartStmt implements RestartStmt: "RESTART".
func (r *rdParser) parseRestartStmt() ast.StmtNode {
	r.expect(restart)
	return &ast.RestartStmt{}
}

// parseLockStmtFamily dispatches LOCK-leading statements: LockTablesStmt
// is ported here; LockStatsStmt falls back.
func (r *rdParser) parseLockStmtFamily() ast.StmtNode {
	if r.la(1) == stats {
		r.unsupported("LOCK STATS statement")
	}
	// LockTablesStmt: "LOCK" TablesTerminalSym TableLockList
	r.expect(lock)
	if !r.accept(tables) && !r.accept(tableKwd) {
		r.syntaxError()
	}
	locks := []ast.TableLock{r.parseTableLock()}
	for r.accept(int(',')) {
		locks = append(locks, r.parseTableLock())
	}
	return &ast.LockTablesStmt{
		TableLocks: locks,
	}
}

// parseTableLock implements TableLock and LockType.
func (r *rdParser) parseTableLock() ast.TableLock {
	tn := r.parseTableName()
	var lockType ast.TableLockType
	switch r.tok() {
	case read:
		r.advance()
		lockType = ast.TableLockRead
		if r.accept(local) {
			lockType = ast.TableLockReadLocal
		}
	case write:
		r.advance()
		lockType = ast.TableLockWrite
		if r.accept(local) {
			lockType = ast.TableLockWriteLocal
		}
	default:
		r.syntaxError()
	}
	return ast.TableLock{
		Table: tn,
		Type:  lockType,
	}
}

// parseUnlockStmtFamily dispatches UNLOCK-leading statements:
// UnlockTablesStmt is ported here; UnlockStatsStmt falls back.
func (r *rdParser) parseUnlockStmtFamily() ast.StmtNode {
	if r.la(1) == stats {
		r.unsupported("UNLOCK STATS statement")
	}
	// UnlockTablesStmt: "UNLOCK" TablesTerminalSym
	r.expect(unlock)
	if !r.accept(tables) && !r.accept(tableKwd) {
		r.syntaxError()
	}
	return &ast.UnlockTablesStmt{}
}

// parseBinlogStmt implements BinlogStmt: "BINLOG" stringLit.
func (r *rdParser) parseBinlogStmt() ast.StmtNode {
	r.expect(binlog)
	return &ast.BinlogStmt{Str: r.expect(stringLit).lit}
}

// parsePreparedStmt implements PreparedStmt: "PREPARE" Identifier "FROM"
// PrepareSQL, with PrepareSQL: stringLit | UserVariable.
func (r *rdParser) parsePreparedStmt() ast.StmtNode {
	r.expect(prepare)
	name := r.parseIdentifier()
	r.expect(from)
	var sqlText string
	var sqlVar *ast.VariableExpr
	if r.tok() == stringLit {
		sqlText = r.cur().lit
		r.advance()
	} else {
		sqlVar = r.parseUserVariable()
	}
	return &ast.PrepareStmt{
		Name:    name,
		SQLText: sqlText,
		SQLVar:  sqlVar,
	}
}

// parseUserVariable implements UserVariable: singleAtIdentifier.
func (r *rdParser) parseUserVariable() *ast.VariableExpr {
	start := r.cur().offset
	t := r.expect(singleAtIdentifier)
	v := strings.TrimPrefix(t.lit, "@")
	return r.setOrigin(&ast.VariableExpr{Name: v, IsGlobal: false, IsSystem: false}, start).(*ast.VariableExpr)
}

// parseExecuteStmt implements ExecuteStmt and UserVariableList.
func (r *rdParser) parseExecuteStmt() ast.StmtNode {
	r.expect(execute)
	name := r.parseIdentifier()
	if r.accept(using) {
		// "EXECUTE" Identifier "USING" UserVariableList
		vars := []ast.ExprNode{r.parseUserVariable()}
		for r.accept(int(',')) {
			vars = append(vars, r.parseUserVariable())
		}
		return &ast.ExecuteStmt{
			Name:      name,
			UsingVars: vars,
		}
	}
	return &ast.ExecuteStmt{Name: name}
}

// parseDeallocateStmt implements DeallocateStmt for the "DEALLOCATE"
// spelling of DeallocateSym ("DROP" "PREPARE" belongs to the DDL family's
// DROP token).
func (r *rdParser) parseDeallocateStmt() ast.StmtNode {
	r.expect(deallocate)
	r.expect(prepare)
	return &ast.DeallocateStmt{Name: r.parseIdentifier()}
}

// parseExplainStmt implements ExplainStmt (and ExplainSym: "EXPLAIN" |
// "DESCRIBE" | "DESC").
func (r *rdParser) parseExplainStmt() ast.StmtNode {
	switch r.tok() {
	case explain, describe, desc:
		r.advance()
	default:
		r.syntaxError()
	}
	if r.tok() == explore &&
		(r.la(1) == selectKwd || r.la(1) == stringLit || r.la(1) == analyze) {
		// ExplainSym "EXPLORE" ["ANALYZE"] (SelectStmt | stringLit)
		r.advance()
		analyzeFlag := r.accept(analyze)
		if r.tok() == stringLit {
			digest := r.cur().lit
			r.advance()
			return &ast.ExplainStmt{
				SQLDigest: digest,
				Explore:   true,
				Analyze:   analyzeFlag,
			}
		}
		startOffset := r.cur().offset
		stmt := r.parseSelectStmt()
		r.p.setNodeText(stmt, strings.TrimSpace(r.src[startOffset:]))
		return &ast.ExplainStmt{
			Stmt:    stmt,
			Explore: true,
			Analyze: analyzeFlag,
		}
	}
	if r.tok() == forKwd {
		// ExplainSym "FOR" "CONNECTION" NUM
		r.advance()
		r.expect(connection)
		return &ast.ExplainForStmt{
			Format:       "row",
			ConnectionID: getUint64FromNUM(r.expect(intLit).item),
		}
	}
	if r.tok() == analyze {
		r.advance()
		if r.tok() == format {
			// ExplainSym "ANALYZE" "FORMAT" "=" (ExplainFormatType|stringLit)
			// (ExplainableStmt|stringLit)
			r.advance()
			r.expect(eq)
			formatStr := r.parseExplainFormat()
			if r.tok() == stringLit {
				digest := r.cur().lit
				r.advance()
				return &ast.ExplainStmt{
					PlanDigest: digest,
					Format:     formatStr,
					Analyze:    true,
				}
			}
			return &ast.ExplainStmt{
				Stmt:    r.parseExplainableStmt(),
				Format:  formatStr,
				Analyze: true,
			}
		}
		if r.tok() == stringLit {
			// ExplainSym "ANALYZE" stringLit
			digest := r.cur().lit
			r.advance()
			return &ast.ExplainStmt{
				PlanDigest: digest,
				Format:     "row",
				Analyze:    true,
			}
		}
		// ExplainSym "ANALYZE" ExplainableStmt
		return &ast.ExplainStmt{
			Stmt:    r.parseExplainableStmt(),
			Format:  "row",
			Analyze: true,
		}
	}
	if r.tok() == format && r.la(1) == eq {
		// ExplainSym "FORMAT" "=" (ExplainFormatType|stringLit) ...
		r.advance()
		r.advance()
		formatStr := r.parseExplainFormat()
		switch r.tok() {
		case forKwd:
			// ... "FOR" "CONNECTION" NUM
			r.advance()
			r.expect(connection)
			return &ast.ExplainForStmt{
				Format:       formatStr,
				ConnectionID: getUint64FromNUM(r.expect(intLit).item),
			}
		case stringLit:
			// ... stringLit (a plan digest)
			digest := r.cur().lit
			r.advance()
			return &ast.ExplainStmt{
				PlanDigest: digest,
				Format:     formatStr,
			}
		}
		// ... ExplainableStmt
		return &ast.ExplainStmt{
			Stmt:   r.parseExplainableStmt(),
			Format: formatStr,
		}
	}
	if r.tok() == stringLit {
		// ExplainSym stringLit
		digest := r.cur().lit
		r.advance()
		return &ast.ExplainStmt{
			PlanDigest: digest,
			Format:     "row",
		}
	}
	switch r.tok() {
	case deleteKwd, update, insert, replace, selectKwd, tableKwd, values, int('('), with, alter, importKwd:
		// ExplainSym ExplainableStmt
		return &ast.ExplainStmt{
			Stmt:   r.parseExplainableStmt(),
			Format: "row",
		}
	}
	// ExplainSym TableName [ColumnName]
	showStmt := &ast.ShowStmt{
		Tp:    ast.ShowColumns,
		Table: r.parseTableName(),
	}
	if isIdentifierTok(r.tok()) {
		showStmt.Column = r.parseColumnName()
	}
	return &ast.ExplainStmt{
		Stmt: showStmt,
	}
}

// parseExplainFormat parses the symbol after "FORMAT" "=": either a
// stringLit or an ExplainFormatType keyword; both yield their spelling.
func (r *rdParser) parseExplainFormat() string {
	switch r.tok() {
	case stringLit, traditional, jsonType, row, dotType, briefType, verboseType, trueCardCost, tidbJson:
		lit := r.cur().lit
		r.advance()
		return lit
	}
	r.syntaxError()
	return ""
}

// parseExplainableStmt implements ExplainableStmt.
func (r *rdParser) parseExplainableStmt() ast.StmtNode {
	switch r.tok() {
	case deleteKwd:
		return r.parseDeleteStmt(nil)
	case update:
		return r.parseUpdateStmt(nil)
	case insert:
		return r.parseInsertIntoStmt()
	case replace:
		return r.parseReplaceIntoStmt()
	case selectKwd, tableKwd, values, int('('):
		// SetOprStmt | SelectStmt | SubSelect (the SubSelect action sets
		// IsInBraces, as in finishSelectFamily).
		return r.finishSelectFamily(nil)
	case with:
		// SelectStmtWithClause, or the WithClause forms of
		// UpdateStmt/DeleteFromStmt.
		withClause := r.parseWithClause()
		switch r.tok() {
		case update:
			return r.parseUpdateStmt(withClause)
		case deleteKwd:
			return r.parseDeleteStmt(withClause)
		}
		return r.finishSelectFamily(withClause)
	case alter:
		r.unsupported("EXPLAIN ALTER TABLE statement")
	case importKwd:
		return r.parseImportIntoStmt()
	}
	r.syntaxError()
	return nil
}

// parseTraceStmt implements TraceStmt.
func (r *rdParser) parseTraceStmt() ast.StmtNode {
	r.expect(trace)
	switch r.tok() {
	case format:
		// "TRACE" "FORMAT" "=" stringLit TraceableStmt
		r.advance()
		r.expect(eq)
		formatStr := r.expect(stringLit).lit
		startOffset := r.cur().offset
		stmt := r.parseTraceableStmt()
		st := &ast.TraceStmt{
			Stmt:      stmt,
			Format:    formatStr,
			TracePlan: false,
		}
		r.p.setNodeText(stmt, string(r.src[startOffset:]))
		return st
	case plan:
		r.advance()
		if r.accept(target) {
			// "TRACE" "PLAN" "TARGET" "=" stringLit TraceableStmt
			r.expect(eq)
			targetStr := r.expect(stringLit).lit
			startOffset := r.cur().offset
			stmt := r.parseTraceableStmt()
			st := &ast.TraceStmt{
				Stmt:            stmt,
				TracePlan:       true,
				TracePlanTarget: targetStr,
			}
			r.p.setNodeText(stmt, string(r.src[startOffset:]))
			return st
		}
		// "TRACE" "PLAN" TraceableStmt
		startOffset := r.cur().offset
		stmt := r.parseTraceableStmt()
		st := &ast.TraceStmt{
			Stmt:      stmt,
			TracePlan: true,
		}
		r.p.setNodeText(stmt, string(r.src[startOffset:]))
		return st
	}
	// "TRACE" TraceableStmt
	startOffset := r.cur().offset
	stmt := r.parseTraceableStmt()
	st := &ast.TraceStmt{
		Stmt:      stmt,
		Format:    "row",
		TracePlan: false,
	}
	r.p.setNodeText(stmt, string(r.src[startOffset:]))
	return st
}

// parseTraceableStmt implements TraceableStmt.
func (r *rdParser) parseTraceableStmt() ast.StmtNode {
	switch r.tok() {
	case deleteKwd:
		return r.parseDeleteStmt(nil)
	case update:
		return r.parseUpdateStmt(nil)
	case insert:
		return r.parseInsertIntoStmt()
	case replace:
		return r.parseReplaceIntoStmt()
	case selectKwd, tableKwd, values, int('('):
		return r.finishSelectFamily(nil)
	case with:
		withClause := r.parseWithClause()
		switch r.tok() {
		case update:
			return r.parseUpdateStmt(withClause)
		case deleteKwd:
			return r.parseDeleteStmt(withClause)
		}
		return r.finishSelectFamily(withClause)
	case load:
		return r.parseLoadDataStmt()
	case begin, start:
		return r.parseBeginTransactionStmt()
	case commit:
		return r.parseCommitStmt()
	case savepoint:
		return r.parseSavepointStmt()
	case release:
		return r.parseReleaseSavepointStmt()
	case rollback:
		return r.parseRollbackStmt()
	case set:
		return r.parseSetStmt()
	case analyze:
		r.unsupported("TRACE ANALYZE TABLE statement")
	}
	r.syntaxError()
	return nil
}
