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

// The SHOW statement family: ShowStmt with all its helper productions
// (ShowTargetFilterable, ShowLikeOrWhereOpt, ShowDatabaseNameOpt, ...).

import (
	"github.com/sqlc-dev/marino/ast"
	"github.com/sqlc-dev/marino/mysql"
)

func init() {
	rdRegister(show, (*rdParser).parseShowStmt)
}

// applyShowLikeOrWhereOpt implements ShowLikeOrWhereOpt together with the
// "SHOW" ShowTargetFilterable ShowLikeOrWhereOpt action: a LIKE pattern
// goes to stmt.Pattern, a WHERE expression to stmt.Where.
func (r *rdParser) applyShowLikeOrWhereOpt(stmt *ast.ShowStmt) ast.StmtNode {
	switch r.tok() {
	case like:
		// "LIKE" SimpleExpr
		r.advance()
		stmt.Pattern = &ast.PatternLikeOrIlikeExpr{
			Pattern:        r.parseSimpleExpr(),
			Escape:         '\\',
			EscapeExplicit: false,
			IsLike:         true,
		}
	case where:
		// "WHERE" Expression
		r.advance()
		stmt.Where = r.parseExpression()
	}
	return stmt
}

// parseShowDatabaseNameOpt implements ShowDatabaseNameOpt: [FromOrIn DBName].
func (r *rdParser) parseShowDatabaseNameOpt() string {
	if r.tok() == from || r.tok() == in {
		r.advance()
		return r.parseIdentifier()
	}
	return ""
}

// parseShowTableAliasOpt implements ShowTableAliasOpt: FromOrIn TableName.
func (r *rdParser) parseShowTableAliasOpt() *ast.TableName {
	if r.tok() != from && r.tok() != in {
		r.syntaxError()
	}
	r.advance()
	return r.parseTableName()
}

// parseIfNotExists implements IfNotExists: [ "IF" NotSym "EXISTS" ].
func (r *rdParser) parseIfNotExists() bool {
	if r.tok() != ifKwd {
		return false
	}
	r.advance()
	// NotSym: not | not2
	if r.tok() != not && r.tok() != not2 {
		r.syntaxError()
	}
	r.advance()
	r.expect(exists)
	return true
}

// parseShowStmt implements ShowStmt, dispatching on the token after SHOW.
func (r *rdParser) parseShowStmt() ast.StmtNode {
	r.expect(show)
	switch r.tok() {
	case backup:
		// BRIEStmt: "SHOW" "BACKUP" "LOGS" "STATUS"
		//         | "SHOW" "BACKUP" "LOGS" "METADATA" "FROM" stringLit
		//         | "SHOW" "BACKUP" "METADATA" "FROM" stringLit
		r.advance()
		stmt := &ast.BRIEStmt{}
		if r.accept(logs) {
			if r.accept(status) {
				stmt.Kind = ast.BRIEKindStreamStatus
				return stmt
			}
			r.expect(metadata)
			r.expect(from)
			stmt.Kind = ast.BRIEKindStreamMetaData
			stmt.Storage = r.expect(stringLit).lit
			return stmt
		}
		r.expect(metadata)
		r.expect(from)
		stmt.Kind = ast.BRIEKindShowBackupMeta
		stmt.Storage = r.expect(stringLit).lit
		return stmt
	case br:
		// BRIEStmt: "SHOW" "BR" "JOB" Int64Num
		//         | "SHOW" "BR" "JOB" "QUERY" Int64Num
		r.advance()
		r.expect(job)
		stmt := &ast.BRIEStmt{}
		if r.accept(query) {
			stmt.Kind = ast.BRIEKindShowQuery
		} else {
			stmt.Kind = ast.BRIEKindShowJob
		}
		stmt.JobID = r.parseInt64Num()
		return stmt
	case traffic:
		// TrafficStmt: "SHOW" "TRAFFIC" "JOBS"
		r.advance()
		r.expect(jobs)
		return &ast.TrafficStmt{OpType: ast.TrafficOpShow}
	case create:
		return r.parseShowCreate()
	case masking:
		// "SHOW" "MASKING" "POLICIES" "FOR" TableName WhereClauseOptional
		r.advance()
		r.expect(policies)
		r.expect(forKwd)
		stmt := &ast.ShowStmt{
			Tp:    ast.ShowMaskingPolicies,
			Table: r.parseTableName(),
		}
		if r.accept(where) {
			stmt.Where = r.parseExpression()
		}
		return stmt
	case tableKwd:
		return r.parseShowTable()
	case grants:
		// "SHOW" "GRANTS" ["FOR" Username UsingRoles]
		r.advance()
		if r.accept(forKwd) {
			user := r.parseUsername()
			if r.accept(using) {
				// UsingRoles: "USING" RolenameList
				return &ast.ShowStmt{
					Tp:    ast.ShowGrants,
					User:  user,
					Roles: r.parseRolenameList(),
				}
			}
			return &ast.ShowStmt{
				Tp:    ast.ShowGrants,
				User:  user,
				Roles: nil,
			}
		}
		return &ast.ShowStmt{Tp: ast.ShowGrants}
	case master:
		// "SHOW" "MASTER" "STATUS"
		r.advance()
		r.expect(status)
		return &ast.ShowStmt{
			Tp: ast.ShowMasterStatus,
		}
	case binaryType:
		// "SHOW" "BINARY" "LOG" "STATUS" | "SHOW" "BINARY" "LOGS"
		r.advance()
		if r.accept(logs) {
			return &ast.ShowStmt{
				Tp: ast.ShowBinaryLogs,
			}
		}
		r.expect(log)
		r.expect(status)
		return &ast.ShowStmt{
			Tp: ast.ShowBinlogStatus,
		}
	case binlog:
		// "SHOW" "BINLOG" "EVENTS" ShowLogEventsTail
		r.advance()
		r.expect(events)
		return r.parseShowLogEventsTail(ast.ShowBinlogEvents)
	case relaylog:
		// "SHOW" "RELAYLOG" "EVENTS" ShowLogEventsTail
		r.advance()
		r.expect(events)
		return r.parseShowLogEventsTail(ast.ShowRelaylogEvents)
	case replicas:
		// "SHOW" "REPLICAS"
		r.advance()
		return &ast.ShowStmt{
			Tp: ast.ShowReplicas,
		}
	case engine:
		// "SHOW" "ENGINE" (Identifier | stringLit) ("STATUS" | "MUTEX");
		// the statement postdates the goyacc grammar, so the shape follows
		// the MySQL 26.7 reference manual.
		r.advance()
		stmt := &ast.ShowStmt{}
		if r.tok() == stringLit {
			stmt.Engine = r.cur().lit
			r.advance()
		} else {
			stmt.Engine = r.parseIdentifier()
		}
		switch r.tok() {
		case status:
			r.advance()
			stmt.Tp = ast.ShowEngineStatus
		case mutex:
			r.advance()
			stmt.Tp = ast.ShowEngineMutex
		default:
			r.syntaxError()
		}
		return stmt
	case parseTree:
		// "SHOW" "PARSE_TREE" select statement; the statement postdates
		// the goyacc grammar, and MySQL 26.7 limits its argument to a
		// select statement.
		r.advance()
		stmt := &ast.ShowStmt{Tp: ast.ShowParseTree}
		switch r.tok() {
		case selectKwd, tableKwd, values, int('('):
			stmt.ParseTreeStmt = r.finishSelectFamily(nil)
		case with:
			stmt.ParseTreeStmt = r.finishSelectFamily(r.parseWithClause())
		default:
			r.syntaxError()
		}
		return stmt
	case replica, slave:
		// "SHOW" Replica "STATUS" ForChannelOpt
		r.advance()
		r.expect(status)
		return &ast.ShowStmt{
			Tp:          ast.ShowReplicaStatus,
			ChannelName: r.parseForChannelOpt(),
		}
	case processlist:
		// "SHOW" OptFull "PROCESSLIST" (empty OptFull)
		r.advance()
		return &ast.ShowStmt{
			Tp:   ast.ShowProcessList,
			Full: false,
		}
	case full:
		r.advance()
		switch r.tok() {
		case processlist:
			// "SHOW" OptFull "PROCESSLIST"
			r.advance()
			return &ast.ShowStmt{
				Tp:   ast.ShowProcessList,
				Full: true,
			}
		case tables:
			// OptFull "TABLES" ShowDatabaseNameOpt
			r.advance()
			return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{
				Tp:     ast.ShowTables,
				DBName: r.parseShowDatabaseNameOpt(),
				Full:   true,
			})
		case fields, columns:
			// OptFull FieldsOrColumns ShowTableAliasOpt ShowDatabaseNameOpt
			r.advance()
			stmt := &ast.ShowStmt{
				Tp:   ast.ShowColumns,
				Full: true,
			}
			stmt.Table = r.parseShowTableAliasOpt()
			stmt.DBName = r.parseShowDatabaseNameOpt()
			return r.applyShowLikeOrWhereOpt(stmt)
		}
		r.syntaxError()
	case profiles:
		// "SHOW" "PROFILES"
		r.advance()
		return &ast.ShowStmt{
			Tp: ast.ShowProfiles,
		}
	case profile:
		// "SHOW" "PROFILE" ShowProfileTypesOpt ShowProfileArgsOpt SelectStmtLimitOpt
		r.advance()
		v := &ast.ShowStmt{
			Tp: ast.ShowProfile,
		}
		if types := r.parseShowProfileTypesOpt(); types != nil {
			v.ShowProfileTypes = types
		}
		if r.tok() == forKwd {
			// ShowProfileArgsOpt: "FOR" "QUERY" Int64Num
			r.advance()
			r.expect(query)
			n := r.parseInt64Num()
			v.ShowProfileArgs = &n
		}
		if r.tok() == limit || r.tok() == fetch {
			v.ShowProfileLimit = r.parseSelectStmtLimit()
		}
		return v
	case privileges:
		r.advance()
		return &ast.ShowStmt{
			Tp: ast.ShowPrivileges,
		}
	case builtins:
		r.advance()
		return &ast.ShowStmt{
			Tp: ast.ShowBuiltins,
		}
	case placement:
		if r.la(1) == forKwd {
			// "SHOW" "PLACEMENT" "FOR" ShowPlacementTarget
			r.advance()
			r.advance()
			return r.parseShowPlacementTarget()
		}
		if r.la(1) == labels {
			// ShowTargetFilterable: "PLACEMENT" "LABELS"
			r.advance()
			r.advance()
			return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{Tp: ast.ShowPlacementLabels})
		}
		// ShowTargetFilterable: "PLACEMENT"
		r.advance()
		return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{Tp: ast.ShowPlacement})
	case raw:
		// ShowImportJob(s)Target: "RAW" "IMPORT" ("JOB" Int64Num | "JOBS")
		r.advance()
		r.expect(importKwd)
		if r.accept(job) {
			v := r.parseInt64Num()
			return &ast.ShowStmt{
				Tp:           ast.ShowImportJobs,
				ImportJobID:  &v,
				ImportJobRaw: true,
			}
		}
		r.expect(jobs)
		return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{Tp: ast.ShowImportJobs, ImportJobRaw: true})
	case importKwd:
		r.advance()
		switch r.tok() {
		case job:
			// "SHOW" ShowImportJobTarget Int64Num
			r.advance()
			v := r.parseInt64Num()
			return &ast.ShowStmt{
				Tp:           ast.ShowImportJobs,
				ImportJobID:  &v,
				ImportJobRaw: false,
			}
		case jobs:
			// ShowImportJobsTarget: "IMPORT" "JOBS"
			r.advance()
			return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{Tp: ast.ShowImportJobs, ImportJobRaw: false})
		case groups:
			// "IMPORT" "GROUPS"
			r.advance()
			return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{Tp: ast.ShowImportGroups})
		case group:
			// "IMPORT" "GROUP" stringLit
			r.advance()
			return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{Tp: ast.ShowImportGroups, ShowGroupKey: r.expect(stringLit).lit})
		}
		r.syntaxError()
	case distribution:
		r.advance()
		if r.accept(job) {
			// "SHOW" "DISTRIBUTION" "JOB" Int64Num
			v := r.parseInt64Num()
			return &ast.ShowStmt{Tp: ast.ShowDistributionJobs, DistributionJobID: &v}
		}
		// ShowTargetFilterable: "DISTRIBUTION" "JOBS"
		r.expect(jobs)
		return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{Tp: ast.ShowDistributionJobs})
	}
	return r.parseShowTargetFilterable()
}

// parseShowCreate implements the "SHOW" "CREATE" ... alternatives.
func (r *rdParser) parseShowCreate() ast.StmtNode {
	r.expect(create)
	switch r.tok() {
	case tableKwd:
		// "SHOW" "CREATE" "TABLE" TableName
		r.advance()
		return &ast.ShowStmt{
			Tp:    ast.ShowCreateTable,
			Table: r.parseTableName(),
		}
	case view:
		// "SHOW" "CREATE" "VIEW" TableName
		r.advance()
		return &ast.ShowStmt{
			Tp:    ast.ShowCreateView,
			Table: r.parseTableName(),
		}
	case database:
		// "SHOW" "CREATE" "DATABASE" IfNotExists DBName
		r.advance()
		ifNotExists := r.parseIfNotExists()
		return &ast.ShowStmt{
			Tp:          ast.ShowCreateDatabase,
			IfNotExists: ifNotExists,
			DBName:      r.parseIdentifier(),
		}
	case sequence:
		// "SHOW" "CREATE" "SEQUENCE" TableName
		r.advance()
		return &ast.ShowStmt{
			Tp:    ast.ShowCreateSequence,
			Table: r.parseTableName(),
		}
	case placement:
		// "SHOW" "CREATE" "PLACEMENT" "POLICY" PolicyName
		r.advance()
		r.expect(policy)
		return &ast.ShowStmt{
			Tp:     ast.ShowCreatePlacementPolicy,
			DBName: r.parseIdentifier(),
		}
	case resource:
		// "SHOW" "CREATE" "RESOURCE" "GROUP" ResourceGroupName
		r.advance()
		r.expect(group)
		return &ast.ShowStmt{
			Tp:                ast.ShowCreateResourceGroup,
			ResourceGroupName: r.parseResourceGroupName(),
		}
	case user:
		// "SHOW" "CREATE" "USER" Username
		// See https://dev.mysql.com/doc/refman/5.7/en/show-create-user.html
		r.advance()
		return &ast.ShowStmt{
			Tp:   ast.ShowCreateUser,
			User: r.parseUsername(),
		}
	case procedure:
		// "SHOW" "CREATE" "PROCEDURE" TableName
		r.advance()
		return &ast.ShowStmt{
			Tp:        ast.ShowCreateProcedure,
			Procedure: r.parseTableName(),
		}
	case function:
		// "SHOW" "CREATE" "FUNCTION" TableName
		r.advance()
		return &ast.ShowStmt{
			Tp:        ast.ShowCreateFunction,
			Procedure: r.parseTableName(),
		}
	case event:
		// "SHOW" "CREATE" "EVENT" TableName
		r.advance()
		return &ast.ShowStmt{
			Tp:    ast.ShowCreateEvent,
			Table: r.parseTableName(),
		}
	case trigger:
		// "SHOW" "CREATE" "TRIGGER" TableName
		r.advance()
		return &ast.ShowStmt{
			Tp:    ast.ShowCreateTrigger,
			Table: r.parseTableName(),
		}
	case library:
		// "SHOW" "CREATE" "LIBRARY" TableName
		r.advance()
		return &ast.ShowStmt{
			Tp:    ast.ShowCreateLibrary,
			Table: r.parseTableName(),
		}
	case masking:
		// "SHOW" "CREATE" "MASKING" "POLICY" PolicyName
		r.advance()
		r.expect(policy)
		return &ast.ShowStmt{
			Tp:     ast.ShowCreateMaskingPolicy,
			DBName: r.parseIdentifier(),
		}
	}
	r.syntaxError()
	return nil
}

// parseShowLogEventsTail implements the shared tail of "SHOW" "BINLOG"
// "EVENTS" and "SHOW" "RELAYLOG" "EVENTS"; the statements postdate the
// goyacc grammar, so the shape follows the MySQL 26.7 reference manual:
//
//	["IN" stringLit] ["FROM" Int64Num] [SelectStmtLimit]
//	    ["FOR" "CHANNEL" stringLit]
func (r *rdParser) parseShowLogEventsTail(tp ast.ShowStmtType) ast.StmtNode {
	stmt := &ast.ShowStmt{Tp: tp}
	if r.accept(in) {
		stmt.LogName = r.expect(stringLit).lit
	}
	if r.accept(from) {
		v := r.parseInt64Num()
		stmt.Position = &v
	}
	if r.tok() == limit {
		stmt.Limit = r.parseSelectStmtLimit()
	}
	if r.accept(forKwd) {
		r.expect(channel)
		stmt.ChannelName = r.expect(stringLit).lit
	}
	return stmt
}

// parseShowTable distinguishes the "SHOW" "TABLE" TableName ... targets
// from ShowTargetFilterable's "TABLE" "STATUS" (where STATUS could also
// be a table name; the target forms are tried first and rewound like the
// LALR states keep both alive).
func (r *rdParser) parseShowTable() ast.StmtNode {
	var result ast.StmtNode
	if r.try(func() { result = r.parseShowTableTarget() }) {
		return result
	}
	// ShowTargetFilterable: "TABLE" "STATUS" ShowDatabaseNameOpt
	r.expect(tableKwd)
	r.expect(status)
	return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{
		Tp:     ast.ShowTableStatus,
		DBName: r.parseShowDatabaseNameOpt(),
	})
}

// parseShowTableTarget implements the "SHOW" "TABLE" TableName
// alternatives (REGIONS, NEXT_ROW_ID, INDEX ... REGIONS, DISTRIBUTIONS).
func (r *rdParser) parseShowTableTarget() ast.StmtNode {
	r.expect(tableKwd)
	tn := r.parseTableName()
	// PartitionNameListOpt: PARTITION is always shifted here, so a
	// missing '(' errors after it.
	partitions := []ast.CIStr{}
	if r.accept(partition) {
		r.expect(int('('))
		partitions = append(partitions, ast.NewCIStr(r.parseIdentifier()))
		for r.accept(int(',')) {
			partitions = append(partitions, ast.NewCIStr(r.parseIdentifier()))
		}
		r.expect(int(')'))
	}
	switch r.tok() {
	case regions:
		// "SHOW" "TABLE" TableName PartitionNameListOpt "REGIONS" WhereClauseOptional
		r.advance()
		stmt := &ast.ShowStmt{
			Tp:    ast.ShowRegions,
			Table: tn,
		}
		stmt.Table.PartitionNames = partitions
		if r.accept(where) {
			stmt.Where = r.parseExpression()
		}
		return stmt
	case next_row_id:
		// "SHOW" "TABLE" TableName "NEXT_ROW_ID" (no partition list)
		if len(partitions) != 0 {
			r.syntaxError()
		}
		r.advance()
		return &ast.ShowStmt{
			Tp:    ast.ShowTableNextRowId,
			Table: tn,
		}
	case index:
		// "SHOW" "TABLE" TableName PartitionNameListOpt "INDEX" Identifier
		// "REGIONS" WhereClauseOptional
		r.advance()
		indexName := r.parseIdentifier()
		r.expect(regions)
		stmt := &ast.ShowStmt{
			Tp:        ast.ShowRegions,
			Table:     tn,
			IndexName: ast.NewCIStr(indexName),
		}
		stmt.Table.PartitionNames = partitions
		if r.accept(where) {
			stmt.Where = r.parseExpression()
		}
		return stmt
	case distributions:
		// "SHOW" "TABLE" TableName PartitionNameListOpt "DISTRIBUTIONS" WhereClauseOptional
		r.advance()
		stmt := &ast.ShowStmt{
			Tp:    ast.ShowDistributions,
			Table: tn,
		}
		stmt.Table.PartitionNames = partitions
		if r.accept(where) {
			stmt.Where = r.parseExpression()
		}
		return stmt
	}
	r.syntaxError()
	return nil
}

// parseShowPlacementTarget implements ShowPlacementTarget.
func (r *rdParser) parseShowPlacementTarget() ast.StmtNode {
	switch r.tok() {
	case database:
		// DatabaseSym DBName
		r.advance()
		return &ast.ShowStmt{
			Tp:     ast.ShowPlacementForDatabase,
			DBName: r.parseIdentifier(),
		}
	case tableKwd:
		r.advance()
		tn := r.parseTableName()
		if r.accept(partition) {
			// "TABLE" TableName "PARTITION" Identifier
			return &ast.ShowStmt{
				Tp:        ast.ShowPlacementForPartition,
				Table:     tn,
				Partition: ast.NewCIStr(r.parseIdentifier()),
			}
		}
		// "TABLE" TableName
		return &ast.ShowStmt{
			Tp:    ast.ShowPlacementForTable,
			Table: tn,
		}
	}
	r.syntaxError()
	return nil
}

// parseShowProfileTypesOpt implements ShowProfileTypesOpt/ShowProfileTypes,
// returning nil for the empty alternative.
func (r *rdParser) parseShowProfileTypesOpt() []int {
	first, ok := r.parseShowProfileType()
	if !ok {
		return nil
	}
	l := []int{first}
	for r.accept(int(',')) {
		t, ok := r.parseShowProfileType()
		if !ok {
			r.syntaxError()
		}
		l = append(l, t)
	}
	return l
}

// parseShowProfileType implements ShowProfileType; ok is false when the
// current token cannot start one.
func (r *rdParser) parseShowProfileType() (int, bool) {
	switch r.tok() {
	case cpu:
		r.advance()
		return ast.ProfileTypeCPU, true
	case memory:
		r.advance()
		return ast.ProfileTypeMemory, true
	case block:
		// "BLOCK" "IO"
		r.advance()
		r.expect(io)
		return ast.ProfileTypeBlockIo, true
	case context:
		// "CONTEXT" "SWITCHES"
		r.advance()
		r.expect(switchesSym)
		return ast.ProfileTypeContextSwitch, true
	case pageSym:
		// "PAGE" "FAULTS"
		r.advance()
		r.expect(faultsSym)
		return ast.ProfileTypePageFaults, true
	case ipc:
		r.advance()
		return ast.ProfileTypeIpc, true
	case swaps:
		r.advance()
		return ast.ProfileTypeSwaps, true
	case source:
		r.advance()
		return ast.ProfileTypeSource, true
	case all:
		r.advance()
		return ast.ProfileTypeAll, true
	}
	return 0, false
}

// parseShowTargetFilterable implements the remaining ShowTargetFilterable
// alternatives and applies ShowLikeOrWhereOpt.
func (r *rdParser) parseShowTargetFilterable() ast.StmtNode {
	switch r.tok() {
	case engines:
		r.advance()
		return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{Tp: ast.ShowEngines})
	case storage:
		// "STORAGE" "ENGINES"; Restore() canonicalizes to SHOW ENGINES.
		r.advance()
		r.expect(engines)
		return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{Tp: ast.ShowEngines})
	case databases:
		r.advance()
		return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{Tp: ast.ShowDatabases})
	case config:
		r.advance()
		return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{Tp: ast.ShowConfig})
	case tables:
		// OptFull "TABLES" ShowDatabaseNameOpt (empty OptFull)
		r.advance()
		return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{
			Tp:     ast.ShowTables,
			DBName: r.parseShowDatabaseNameOpt(),
			Full:   false,
		})
	case open:
		// "OPEN" "TABLES" ShowDatabaseNameOpt
		r.advance()
		r.expect(tables)
		return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{
			Tp:     ast.ShowOpenTables,
			DBName: r.parseShowDatabaseNameOpt(),
		})
	case index, keys, indexes:
		// ShowIndexKwd FromOrIn TableName
		// ShowIndexKwd FromOrIn Identifier FromOrIn Identifier
		r.advance()
		if r.tok() != from && r.tok() != in {
			r.syntaxError()
		}
		r.advance()
		tn := r.parseTableName()
		if (r.tok() == from || r.tok() == in) && tn.Schema.O == "" {
			r.advance()
			tn = &ast.TableName{Name: tn.Name, Schema: ast.NewCIStr(r.parseIdentifier())}
		}
		return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{
			Tp:    ast.ShowIndex,
			Table: tn,
		})
	case fields, columns:
		// OptFull FieldsOrColumns ShowTableAliasOpt ShowDatabaseNameOpt
		r.advance()
		stmt := &ast.ShowStmt{
			Tp:   ast.ShowColumns,
			Full: false,
		}
		stmt.Table = r.parseShowTableAliasOpt()
		stmt.DBName = r.parseShowDatabaseNameOpt()
		return r.applyShowLikeOrWhereOpt(stmt)
	case extended:
		// "EXTENDED" ShowIndexKwd FromOrIn TableName
		// | "EXTENDED" OptFull FieldsOrColumns ShowTableAliasOpt ShowDatabaseNameOpt
		r.advance()
		if r.tok() == index || r.tok() == keys || r.tok() == indexes {
			r.advance()
			if r.tok() != from && r.tok() != in {
				r.syntaxError()
			}
			r.advance()
			tn := r.parseTableName()
			if (r.tok() == from || r.tok() == in) && tn.Schema.O == "" {
				r.advance()
				tn = &ast.TableName{Name: tn.Name, Schema: ast.NewCIStr(r.parseIdentifier())}
			}
			return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{
				Tp:       ast.ShowIndex,
				Table:    tn,
				Extended: true,
			})
		}
		fullFlag := r.accept(full)
		if r.tok() != fields && r.tok() != columns {
			r.syntaxError()
		}
		r.advance()
		stmt := &ast.ShowStmt{
			Tp:       ast.ShowColumns,
			Full:     fullFlag,
			Extended: true,
		}
		stmt.Table = r.parseShowTableAliasOpt()
		stmt.DBName = r.parseShowDatabaseNameOpt()
		return r.applyShowLikeOrWhereOpt(stmt)
	case builtinCount:
		// builtinCount '(' '*' ')' ("WARNINGS" | "ERRORS")
		r.advance()
		r.expect(int('('))
		r.expect(int('*'))
		r.expect(int(')'))
		switch r.tok() {
		case warnings:
			r.advance()
			return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{Tp: ast.ShowWarnings, CountWarningsOrErrors: true})
		case identSQLErrors:
			r.advance()
			return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{Tp: ast.ShowErrors, CountWarningsOrErrors: true})
		}
		r.syntaxError()
	case warnings:
		// "WARNINGS" SelectStmtLimitOpt
		r.advance()
		stmt := &ast.ShowStmt{Tp: ast.ShowWarnings}
		if r.tok() == limit || r.tok() == fetch {
			stmt.Limit = r.parseSelectStmtLimit()
		}
		return r.applyShowLikeOrWhereOpt(stmt)
	case identSQLErrors:
		// "ERRORS" SelectStmtLimitOpt
		r.advance()
		stmt := &ast.ShowStmt{Tp: ast.ShowErrors}
		if r.tok() == limit || r.tok() == fetch {
			stmt.Limit = r.parseSelectStmtLimit()
		}
		return r.applyShowLikeOrWhereOpt(stmt)
	case global, session, variables, status, bindings:
		// GlobalScope ("VARIABLES" | "STATUS" | "BINDINGS")
		globalScope := false
		if r.tok() == global || r.tok() == session {
			globalScope = r.tok() == global
			r.advance()
		}
		var tp ast.ShowStmtType
		switch r.tok() {
		case variables:
			tp = ast.ShowVariables
		case status:
			tp = ast.ShowStatus
		case bindings:
			tp = ast.ShowBindings
		default:
			r.syntaxError()
		}
		r.advance()
		return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{
			Tp:          tp,
			GlobalScope: globalScope,
		})
	case collation:
		r.advance()
		return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{
			Tp: ast.ShowCollation,
		})
	case triggers:
		// "TRIGGERS" ShowDatabaseNameOpt
		r.advance()
		return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{
			Tp:     ast.ShowTriggers,
			DBName: r.parseShowDatabaseNameOpt(),
		})
	case bindingCache:
		// "BINDING_CACHE" "STATUS"
		r.advance()
		r.expect(status)
		return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{
			Tp: ast.ShowBindingCacheStatus,
		})
	case procedure:
		// "PROCEDURE" "STATUS" | "PROCEDURE" "CODE" TableName
		r.advance()
		if r.accept(code) {
			return &ast.ShowStmt{
				Tp:        ast.ShowProcedureCode,
				Procedure: r.parseTableName(),
			}
		}
		r.expect(status)
		return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{
			Tp: ast.ShowProcedureStatus,
		})
	case function:
		// "FUNCTION" "STATUS" | "FUNCTION" "CODE" TableName
		r.advance()
		if r.accept(code) {
			return &ast.ShowStmt{
				Tp:        ast.ShowFunctionCode,
				Procedure: r.parseTableName(),
			}
		}
		r.expect(status)
		return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{
			Tp: ast.ShowFunctionStatus,
		})
	case library:
		// "LIBRARY" "STATUS"
		r.advance()
		r.expect(status)
		return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{
			Tp: ast.ShowLibraryStatus,
		})
	case events:
		// "EVENTS" ShowDatabaseNameOpt
		r.advance()
		return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{
			Tp:     ast.ShowEvents,
			DBName: r.parseShowDatabaseNameOpt(),
		})
	case plugins:
		r.advance()
		return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{
			Tp: ast.ShowPlugins,
		})
	case sessionStates:
		r.advance()
		return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{Tp: ast.ShowSessionStates})
	case statsExtended:
		r.advance()
		return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{Tp: ast.ShowStatsExtended})
	case statsMeta:
		r.advance()
		return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{Tp: ast.ShowStatsMeta, Table: &ast.TableName{Name: ast.NewCIStr("STATS_META"), Schema: ast.NewCIStr(mysql.SystemDB)}})
	case statsHistograms:
		r.advance()
		return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{Tp: ast.ShowStatsHistograms, Table: &ast.TableName{Name: ast.NewCIStr("STATS_HISTOGRAMS"), Schema: ast.NewCIStr(mysql.SystemDB)}})
	case statsTopN:
		r.advance()
		return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{Tp: ast.ShowStatsTopN})
	case statsBuckets:
		r.advance()
		return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{Tp: ast.ShowStatsBuckets, Table: &ast.TableName{Name: ast.NewCIStr("STATS_BUCKETS"), Schema: ast.NewCIStr(mysql.SystemDB)}})
	case statsHealthy:
		r.advance()
		return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{Tp: ast.ShowStatsHealthy})
	case statsLocked:
		r.advance()
		return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{Tp: ast.ShowStatsLocked, Table: &ast.TableName{Name: ast.NewCIStr("STATS_TABLE_LOCKED"), Schema: ast.NewCIStr(mysql.SystemDB)}})
	case histogramsInFlight:
		r.advance()
		return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{Tp: ast.ShowHistogramsInFlight})
	case columnStatsUsage:
		r.advance()
		return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{Tp: ast.ShowColumnStatsUsage})
	case affinity:
		r.advance()
		return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{Tp: ast.ShowAffinity})
	case analyze:
		// "ANALYZE" "STATUS"
		r.advance()
		r.expect(status)
		return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{Tp: ast.ShowAnalyzeStatus})
	case backups:
		r.advance()
		return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{Tp: ast.ShowBackups})
	case restores:
		r.advance()
		return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{Tp: ast.ShowRestores})
	}
	if r.isCharsetKw() {
		// CharsetKw
		r.parseCharsetKw()
		return r.applyShowLikeOrWhereOpt(&ast.ShowStmt{Tp: ast.ShowCharset})
	}
	r.unsupported("SHOW statement target")
	return nil
}
