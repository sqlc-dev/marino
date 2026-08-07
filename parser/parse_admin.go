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

// The ADMIN statement family (AdminStmt and its helper productions
// AdminShowSlow, AdminStmtLimitOpt, AlterJobOptionList, HandleRangeList,
// NumList, BDRRole, StatementScope).

import (
	"strings"

	"github.com/sqlc-dev/marino/ast"
)

func init() {
	rdRegister(admin, (*rdParser).parseAdminStmt)
}

// parseAdminStmt implements AdminStmt (every alternative).
func (r *rdParser) parseAdminStmt() ast.StmtNode {
	r.expect(admin)
	switch r.tok() {
	case show:
		return r.parseAdminShow()
	case check:
		r.advance()
		switch r.tok() {
		case tableKwd:
			// "ADMIN" "CHECK" "TABLE" TableNameList
			r.advance()
			return &ast.AdminStmt{
				Tp:     ast.AdminCheckTable,
				Tables: r.parseTableNameList(),
			}
		case index:
			// "ADMIN" "CHECK" "INDEX" TableName Identifier [HandleRangeList]
			r.advance()
			table := r.parseTableName()
			indexName := r.parseIdentifier()
			if r.tok() == int('(') {
				return &ast.AdminStmt{
					Tp:           ast.AdminCheckIndexRange,
					Tables:       []*ast.TableName{table},
					Index:        indexName,
					HandleRanges: r.parseHandleRangeList(),
				}
			}
			return &ast.AdminStmt{
				Tp:     ast.AdminCheckIndex,
				Tables: []*ast.TableName{table},
				Index:  indexName,
			}
		}
	case checksum:
		// "ADMIN" "CHECKSUM" "TABLE" TableNameList
		r.advance()
		r.expect(tableKwd)
		return &ast.AdminStmt{
			Tp:     ast.AdminChecksumTable,
			Tables: r.parseTableNameList(),
		}
	case recover:
		// "ADMIN" "RECOVER" "INDEX" TableName Identifier
		r.advance()
		r.expect(index)
		return &ast.AdminStmt{
			Tp:     ast.AdminRecoverIndex,
			Tables: []*ast.TableName{r.parseTableName()},
			Index:  r.parseIdentifier(),
		}
	case create:
		// "ADMIN" "CREATE" "WORKLOAD" "SNAPSHOT"
		r.advance()
		r.expect(workload)
		r.expect(snapshot)
		return &ast.AdminStmt{Tp: ast.AdminWorkloadRepoCreate}
	case cleanup:
		r.advance()
		switch r.tok() {
		case index:
			// "ADMIN" "CLEANUP" "INDEX" TableName Identifier
			r.advance()
			return &ast.AdminStmt{
				Tp:     ast.AdminCleanupIndex,
				Tables: []*ast.TableName{r.parseTableName()},
				Index:  r.parseIdentifier(),
			}
		case tableKwd:
			// "ADMIN" "CLEANUP" "TABLE" "LOCK" TableNameList
			r.advance()
			r.expect(lock)
			return &ast.CleanupTableLockStmt{
				Tables: r.parseTableNameList(),
			}
		}
	case cancel:
		// "ADMIN" "CANCEL" "DDL" "JOBS" NumList
		r.advance()
		r.expect(ddl)
		r.expect(jobs)
		return &ast.AdminStmt{
			Tp:     ast.AdminCancelDDLJobs,
			JobIDs: r.parseNumList(),
		}
	case pause:
		// "ADMIN" "PAUSE" "DDL" "JOBS" NumList
		r.advance()
		r.expect(ddl)
		r.expect(jobs)
		return &ast.AdminStmt{
			Tp:     ast.AdminPauseDDLJobs,
			JobIDs: r.parseNumList(),
		}
	case resume:
		// "ADMIN" "RESUME" "DDL" "JOBS" NumList
		r.advance()
		r.expect(ddl)
		r.expect(jobs)
		return &ast.AdminStmt{
			Tp:     ast.AdminResumeDDLJobs,
			JobIDs: r.parseNumList(),
		}
	case alter:
		// "ADMIN" "ALTER" "DDL" "JOBS" Int64Num AlterJobOptionList
		r.advance()
		r.expect(ddl)
		r.expect(jobs)
		return &ast.AdminStmt{
			Tp:              ast.AdminAlterDDLJob,
			JobNumber:       r.parseInt64Num(),
			AlterJobOptions: r.parseAlterJobOptionList(),
		}
	case reload:
		r.advance()
		switch r.tok() {
		case exprPushdownBlacklist:
			// "ADMIN" "RELOAD" "EXPR_PUSHDOWN_BLACKLIST"
			r.advance()
			return &ast.AdminStmt{Tp: ast.AdminReloadExprPushdownBlacklist}
		case optRuleBlacklist:
			// "ADMIN" "RELOAD" "OPT_RULE_BLACKLIST"
			r.advance()
			return &ast.AdminStmt{Tp: ast.AdminReloadOptRuleBlacklist}
		case bindings:
			// "ADMIN" "RELOAD" "BINDINGS"
			r.advance()
			return &ast.AdminStmt{Tp: ast.AdminReloadBindings}
		case cluster:
			// "ADMIN" "RELOAD" "CLUSTER" "BINDINGS"
			r.advance()
			r.expect(bindings)
			return &ast.AdminStmt{Tp: ast.AdminReloadClusterBindings}
		case statsExtended:
			// "ADMIN" "RELOAD" "STATS_EXTENDED"
			r.advance()
			return &ast.AdminStmt{Tp: ast.AdminReloadStatistics}
		case statistics:
			// "ADMIN" "RELOAD" "STATISTICS"
			r.advance()
			return &ast.AdminStmt{Tp: ast.AdminReloadStatistics}
		}
	case plugins:
		r.advance()
		switch r.tok() {
		case enable:
			// "ADMIN" "PLUGINS" "ENABLE" PluginNameList
			r.advance()
			return &ast.AdminStmt{
				Tp:      ast.AdminPluginEnable,
				Plugins: r.parsePluginNameList(),
			}
		case disable:
			// "ADMIN" "PLUGINS" "DISABLE" PluginNameList
			r.advance()
			return &ast.AdminStmt{
				Tp:      ast.AdminPluginDisable,
				Plugins: r.parsePluginNameList(),
			}
		}
	case repair:
		// "ADMIN" "REPAIR" "TABLE" TableName CreateTableStmt
		r.advance()
		r.expect(tableKwd)
		table := r.parseTableName()
		return &ast.RepairTableStmt{
			Table:      table,
			CreateStmt: r.parseCreateTableStmt().(*ast.CreateTableStmt),
		}
	case flush:
		r.advance()
		if r.tok() == bindings {
			// "ADMIN" "FLUSH" "BINDINGS"
			r.advance()
			return &ast.AdminStmt{Tp: ast.AdminFlushBindings}
		}
		// "ADMIN" "FLUSH" StatementScope "PLAN_CACHE"
		scope := r.parseStatementScope()
		r.expect(planCache)
		return &ast.AdminStmt{
			Tp:             ast.AdminFlushPlanCache,
			StatementScope: scope,
		}
	case capture:
		// "ADMIN" "CAPTURE" "BINDINGS"
		r.advance()
		r.expect(bindings)
		return &ast.AdminStmt{Tp: ast.AdminCaptureBindings}
	case evolve:
		// "ADMIN" "EVOLVE" "BINDINGS"
		r.advance()
		r.expect(bindings)
		return &ast.AdminStmt{Tp: ast.AdminEvolveBindings}
	case set:
		// "ADMIN" "SET" "BDR" "ROLE" BDRRole
		r.advance()
		r.expect(bdr)
		r.expect(role)
		return &ast.AdminStmt{
			Tp:      ast.AdminSetBDRRole,
			BDRRole: r.parseBDRRole(),
		}
	case unset:
		// "ADMIN" "UNSET" "BDR" "ROLE"
		r.advance()
		r.expect(bdr)
		r.expect(role)
		return &ast.AdminStmt{Tp: ast.AdminUnsetBDRRole}
	}
	r.syntaxError()
	return nil
}

// parseAdminShow implements the "ADMIN" "SHOW" alternatives of AdminStmt.
func (r *rdParser) parseAdminShow() ast.StmtNode {
	r.expect(show)
	switch r.tok() {
	case ddl:
		r.advance()
		switch r.tok() {
		case jobs:
			// "ADMIN" "SHOW" "DDL" "JOBS" [Int64Num] WhereClauseOptional
			r.advance()
			stmt := &ast.AdminStmt{Tp: ast.AdminShowDDLJobs}
			if r.tok() == intLit {
				stmt.JobNumber = r.parseInt64Num()
			}
			if where := r.parseWhereClauseOptional(); where != nil {
				stmt.Where = where
			}
			return stmt
		case job:
			// "ADMIN" "SHOW" "DDL" "JOB" "QUERIES" (NumList|AdminStmtLimitOpt)
			r.advance()
			r.expect(queries)
			if r.tok() == limit {
				ret := &ast.AdminStmt{
					Tp: ast.AdminShowDDLJobQueriesWithRange,
				}
				limitSimple := r.parseAdminStmtLimit()
				ret.LimitSimple.Count = limitSimple.Count
				ret.LimitSimple.Offset = limitSimple.Offset
				return ret
			}
			return &ast.AdminStmt{
				Tp:     ast.AdminShowDDLJobQueries,
				JobIDs: r.parseNumList(),
			}
		}
		// "ADMIN" "SHOW" "DDL"
		return &ast.AdminStmt{Tp: ast.AdminShowDDL}
	case slow:
		// "ADMIN" "SHOW" "SLOW" AdminShowSlow
		r.advance()
		return &ast.AdminStmt{
			Tp:       ast.AdminShowSlow,
			ShowSlow: r.parseAdminShowSlow(),
		}
	case bdr:
		// "ADMIN" "SHOW" "BDR" "ROLE"
		r.advance()
		r.expect(role)
		return &ast.AdminStmt{Tp: ast.AdminShowBDRRole}
	}
	// "ADMIN" "SHOW" TableName "NEXT_ROW_ID"
	stmt := &ast.AdminStmt{
		Tp:     ast.AdminShowNextRowID,
		Tables: []*ast.TableName{r.parseTableName()},
	}
	r.expect(next_row_id)
	return stmt
}

// parseWhereClauseOptional implements WhereClauseOptional (and
// WhereClause); nil when absent.
func (r *rdParser) parseWhereClauseOptional() ast.ExprNode {
	if !r.accept(where) {
		return nil
	}
	return r.parseExpression()
}

// parseAdminStmtLimit implements AdminStmtLimitOpt.
func (r *rdParser) parseAdminStmtLimit() *ast.LimitSimple {
	r.expect(limit)
	n := r.parseLengthNum()
	switch {
	case r.accept(int(',')):
		// "LIMIT" LengthNum ',' LengthNum
		return &ast.LimitSimple{Offset: n, Count: r.parseLengthNum()}
	case r.accept(offset):
		// "LIMIT" LengthNum "OFFSET" LengthNum
		return &ast.LimitSimple{Offset: r.parseLengthNum(), Count: n}
	}
	// "LIMIT" LengthNum
	return &ast.LimitSimple{Offset: 0, Count: n}
}

// parseBDRRole implements BDRRole: "PRIMARY" | "SECONDARY".
func (r *rdParser) parseBDRRole() ast.BDRRole {
	switch r.tok() {
	case primary:
		r.advance()
		return ast.BDRRolePrimary
	case secondary:
		r.advance()
		return ast.BDRRoleSecondary
	}
	r.syntaxError()
	return ""
}

// parseStatementScope implements StatementScope:
// empty | "GLOBAL" | "INSTANCE" | "SESSION".
func (r *rdParser) parseStatementScope() ast.StatementScope {
	switch r.tok() {
	case global:
		r.advance()
		return ast.StatementScopeGlobal
	case instance:
		r.advance()
		return ast.StatementScopeInstance
	case session:
		r.advance()
		return ast.StatementScopeSession
	}
	return ast.StatementScopeSession
}

// parseAdminShowSlow implements AdminShowSlow.
func (r *rdParser) parseAdminShowSlow() *ast.ShowSlow {
	switch r.tok() {
	case recent:
		// "RECENT" NUM
		r.advance()
		return &ast.ShowSlow{
			Tp:    ast.ShowSlowRecent,
			Count: getUint64FromNUM(r.expect(intLit).item),
		}
	case top:
		r.advance()
		switch r.tok() {
		case internal:
			// "TOP" "INTERNAL" NUM
			r.advance()
			return &ast.ShowSlow{
				Tp:    ast.ShowSlowTop,
				Kind:  ast.ShowSlowKindInternal,
				Count: getUint64FromNUM(r.expect(intLit).item),
			}
		case all:
			// "TOP" "ALL" NUM
			r.advance()
			return &ast.ShowSlow{
				Tp:    ast.ShowSlowTop,
				Kind:  ast.ShowSlowKindAll,
				Count: getUint64FromNUM(r.expect(intLit).item),
			}
		}
		// "TOP" NUM
		return &ast.ShowSlow{
			Tp:    ast.ShowSlowTop,
			Kind:  ast.ShowSlowKindDefault,
			Count: getUint64FromNUM(r.expect(intLit).item),
		}
	}
	r.syntaxError()
	return nil
}

// parseHandleRangeList implements HandleRangeList and HandleRange.
func (r *rdParser) parseHandleRangeList() []ast.HandleRange {
	list := []ast.HandleRange{r.parseHandleRange()}
	for r.accept(int(',')) {
		list = append(list, r.parseHandleRange())
	}
	return list
}

// parseHandleRange implements HandleRange: '(' Int64Num ',' Int64Num ')'.
func (r *rdParser) parseHandleRange() ast.HandleRange {
	r.expect(int('('))
	begin := r.parseInt64Num()
	r.expect(int(','))
	end := r.parseInt64Num()
	r.expect(int(')'))
	return ast.HandleRange{Begin: begin, End: end}
}

// parseNumList implements NumList: Int64Num | NumList ',' Int64Num.
func (r *rdParser) parseNumList() []int64 {
	list := []int64{r.parseInt64Num()}
	for r.accept(int(',')) {
		list = append(list, r.parseInt64Num())
	}
	return list
}

// parseAlterJobOptionList implements AlterJobOptionList.
func (r *rdParser) parseAlterJobOptionList() []*ast.AlterJobOption {
	list := []*ast.AlterJobOption{r.parseAlterJobOption()}
	for r.accept(int(',')) {
		list = append(list, r.parseAlterJobOption())
	}
	return list
}

// parseAlterJobOption implements AlterJobOption:
// identifier "=" SignedLiteral (the bare identifier token, not the
// Identifier production).
func (r *rdParser) parseAlterJobOption() *ast.AlterJobOption {
	name := r.expect(identifier).lit
	r.expect(eq)
	return &ast.AlterJobOption{
		Name:  strings.ToLower(name),
		Value: r.parseSignedLiteral(),
	}
}
