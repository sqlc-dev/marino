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

// The DROP statement family: DropTableStmt, DropViewStmt,
// DropIndexStmt (with IndexLockAndAlgorithmOpt), DropDatabaseStmt,
// DropSequenceStmt, DropUserStmt, DropRoleStmt, DropStatsStmt,
// DropStatisticsStmt, DropPolicyStmt, DropResourceGroupStmt, and the
// "DROP" "PREPARE" spelling of DeallocateStmt.

import (
	"fmt"
	"strings"

	"github.com/sqlc-dev/marino/ast"
	"github.com/sqlc-dev/marino/auth"
)

func init() {
	rdRegister(drop, (*rdParser).parseDropStmtFamily)
}

// parseDropStmtFamily dispatches the DROP statement family on the
// token after "DROP".
func (r *rdParser) parseDropStmtFamily() ast.StmtNode {
	switch r.la(1) {
	case tableKwd, tables, temporary:
		return r.parseDropTableStmt()
	case global:
		if r.la(2) == temporary {
			return r.parseDropTableStmt()
		}
		// DropBindingStmt: "DROP" GlobalScope "BINDING" ...
		return r.parseDropBindingStmt()
	case view:
		return r.parseDropViewStmt()
	case index, hypo:
		return r.parseDropIndexStmt()
	case database:
		return r.parseDropDatabaseStmt()
	case sequence:
		return r.parseDropSequenceStmt()
	case user:
		return r.parseDropUserStmt()
	case role:
		return r.parseDropRoleStmt()
	case stats:
		return r.parseDropStatsStmt()
	case statistics:
		// DropStatisticsStmt: "DROP" "STATISTICS" Identifier
		r.advance()
		r.advance()
		return &ast.DropStatisticsStmt{StatsName: r.parseIdentifier()}
	case placement:
		// DropPolicyStmt: "DROP" "PLACEMENT" "POLICY" IfExists PolicyName
		r.advance()
		r.advance()
		r.expect(policy)
		ifExists := r.parseIfExists()
		return &ast.DropPlacementPolicyStmt{
			IfExists:   ifExists,
			PolicyName: ast.NewCIStr(r.parseIdentifier()),
		}
	case resource:
		// DropResourceGroupStmt: "DROP" "RESOURCE" "GROUP" IfExists
		// ResourceGroupName
		r.advance()
		r.advance()
		r.expect(group)
		ifExists := r.parseIfExists()
		return &ast.DropResourceGroupStmt{
			IfExists:          ifExists,
			ResourceGroupName: ast.NewCIStr(r.parseResourceGroupName()),
		}
	case prepare:
		// DeallocateStmt: DeallocateSym "PREPARE" Identifier, for the
		// "DROP" spelling of DeallocateSym (parse_misc.go owns the
		// "DEALLOCATE" spelling and builds the same AST).
		r.advance()
		r.advance()
		return &ast.DeallocateStmt{Name: r.parseIdentifier()}
	case binding, session:
		// DropBindingStmt: "DROP" GlobalScope "BINDING" ...
		return r.parseDropBindingStmt()
	case procedure:
		return r.parseDropProcedureStmt()
	default:
		r.unsupported(fmt.Sprintf("DROP %s", r.at(r.i+1).lit))
	}
	return nil
}

// parseIfExists implements IfExists: [ "IF" "EXISTS" ].
func (r *rdParser) parseIfExists() bool {
	if r.tok() != ifKwd {
		return false
	}
	r.advance()
	r.expect(exists)
	return true
}

// parseRestrictOrCascadeOpt consumes RestrictOrCascadeOpt:
// empty | "RESTRICT" | "CASCADE" (value ignored).
func (r *rdParser) parseRestrictOrCascadeOpt() {
	if r.tok() == restrict || r.tok() == cascade {
		r.advance()
	}
}

// parseDropTableStmt implements:
//
//	DropTableStmt: "DROP" OptTemporary TableOrTables IfExists
//	    TableNameList RestrictOrCascadeOpt
func (r *rdParser) parseDropTableStmt() ast.StmtNode {
	r.expect(drop)
	// OptTemporary
	temporaryKeyword := ast.TemporaryNone
	switch r.tok() {
	case temporary:
		r.advance()
		temporaryKeyword = ast.TemporaryLocal
	case global:
		r.advance()
		r.expect(temporary)
		temporaryKeyword = ast.TemporaryGlobal
	}
	// TableOrTables
	if r.tok() != tableKwd && r.tok() != tables {
		r.syntaxError()
	}
	r.advance()
	ifExists := r.parseIfExists()
	tableNames := r.parseTableNameList()
	r.parseRestrictOrCascadeOpt()
	return &ast.DropTableStmt{IfExists: ifExists, Tables: tableNames, IsView: false, TemporaryKeyword: temporaryKeyword}
}

// parseDropViewStmt implements:
//
//	DropViewStmt: "DROP" "VIEW" ["IF" "EXISTS"] TableNameList
//	    RestrictOrCascadeOpt
func (r *rdParser) parseDropViewStmt() ast.StmtNode {
	r.expect(drop)
	r.expect(view)
	ifExists := r.parseIfExists()
	tableNames := r.parseTableNameList()
	r.parseRestrictOrCascadeOpt()
	return &ast.DropTableStmt{IfExists: ifExists, Tables: tableNames, IsView: true}
}

// parseDropIndexStmt implements:
//
//	DropIndexStmt: "DROP" "INDEX" IfExists Identifier "ON" TableName
//	    IndexLockAndAlgorithmOpt
//	  | "DROP" "HYPO" "INDEX" IfExists Identifier "ON" TableName
func (r *rdParser) parseDropIndexStmt() ast.StmtNode {
	r.expect(drop)
	if r.accept(hypo) {
		r.expect(index)
		ifExists := r.parseIfExists()
		indexName := r.parseIdentifier()
		r.expect(on)
		return &ast.DropIndexStmt{IfExists: ifExists, IndexName: indexName, Table: r.parseTableName(), IsHypo: true}
	}
	r.expect(index)
	ifExists := r.parseIfExists()
	indexName := r.parseIdentifier()
	r.expect(on)
	table := r.parseTableName()
	lockAlgOpt := r.parseDropIndexLockAndAlgorithm()
	var indexLockAndAlgorithm *ast.IndexLockAndAlgorithm
	if lockAlgOpt != nil {
		indexLockAndAlgorithm = lockAlgOpt
		if indexLockAndAlgorithm.LockTp == ast.LockTypeDefault && indexLockAndAlgorithm.AlgorithmTp == ast.AlgorithmTypeDefault {
			indexLockAndAlgorithm = nil
		}
	}
	return &ast.DropIndexStmt{IfExists: ifExists, IndexName: indexName, Table: table, LockAlg: indexLockAndAlgorithm}
}

// parseDropIndexLockAndAlgorithm implements IndexLockAndAlgorithmOpt
// (nil for the empty alternative; a missing clause is
// LockTypeDefault/AlgorithmTypeDefault as in the one-clause
// alternatives). Named for its DROP INDEX home to stay clear of the
// ALTER TABLE family's clause helpers; CREATE INDEX reuses it.
func (r *rdParser) parseDropIndexLockAndAlgorithm() *ast.IndexLockAndAlgorithm {
	switch r.tok() {
	case lock:
		lockTp := r.parseDropIndexLockClause()
		if r.tok() == algorithm {
			return &ast.IndexLockAndAlgorithm{
				LockTp:      lockTp,
				AlgorithmTp: r.parseDropIndexAlgorithmClause(),
			}
		}
		return &ast.IndexLockAndAlgorithm{
			LockTp:      lockTp,
			AlgorithmTp: ast.AlgorithmTypeDefault,
		}
	case algorithm:
		algorithmTp := r.parseDropIndexAlgorithmClause()
		if r.tok() == lock {
			return &ast.IndexLockAndAlgorithm{
				LockTp:      r.parseDropIndexLockClause(),
				AlgorithmTp: algorithmTp,
			}
		}
		return &ast.IndexLockAndAlgorithm{
			LockTp:      ast.LockTypeDefault,
			AlgorithmTp: algorithmTp,
		}
	}
	return nil
}

// parseDropIndexLockClause implements LockClause.
func (r *rdParser) parseDropIndexLockClause() ast.LockType {
	r.expect(lock)
	r.parseEqOpt()
	if r.accept(defaultKwd) {
		// "LOCK" EqOpt "DEFAULT"
		return ast.LockTypeDefault
	}
	// "LOCK" EqOpt Identifier
	name := r.parseIdentifier()
	id := strings.ToUpper(name)
	if id == "NONE" {
		return ast.LockTypeNone
	} else if id == "SHARED" {
		return ast.LockTypeShared
	} else if id == "EXCLUSIVE" {
		return ast.LockTypeExclusive
	}
	r.actionError(ErrUnknownAlterLock.GenWithStackByArgs(name))
	return ast.LockTypeDefault
}

// parseDropIndexAlgorithmClause implements AlgorithmClause. The
// unknown-algorithm error reports $1 — the "ALGORITHM" keyword's
// spelling — exactly as the goyacc action does.
func (r *rdParser) parseDropIndexAlgorithmClause() ast.AlgorithmType {
	algorithmTok := r.expect(algorithm)
	r.parseEqOpt()
	switch r.tok() {
	case defaultKwd:
		r.advance()
		return ast.AlgorithmTypeDefault
	case copyKwd:
		r.advance()
		return ast.AlgorithmTypeCopy
	case inplace:
		r.advance()
		return ast.AlgorithmTypeInplace
	case instant:
		r.advance()
		return ast.AlgorithmTypeInstant
	case identifier:
		r.advance()
		r.actionError(ErrUnknownAlterAlgorithm.GenWithStackByArgs(algorithmTok.lit))
	}
	r.syntaxError()
	return ast.AlgorithmTypeDefault
}

// parseDropDatabaseStmt implements:
//
//	DropDatabaseStmt: "DROP" DatabaseSym IfExists DBName
func (r *rdParser) parseDropDatabaseStmt() ast.StmtNode {
	r.expect(drop)
	r.expect(database)
	ifExists := r.parseIfExists()
	return &ast.DropDatabaseStmt{IfExists: ifExists, Name: ast.NewCIStr(r.parseIdentifier())}
}

// parseDropSequenceStmt implements:
//
//	DropSequenceStmt: "DROP" "SEQUENCE" IfExists TableNameList
func (r *rdParser) parseDropSequenceStmt() ast.StmtNode {
	r.expect(drop)
	r.expect(sequence)
	ifExists := r.parseIfExists()
	return &ast.DropSequenceStmt{
		IfExists:  ifExists,
		Sequences: r.parseTableNameList(),
	}
}

// parseDropUserStmt implements:
//
//	DropUserStmt: "DROP" "USER" ["IF" "EXISTS"] UsernameList
func (r *rdParser) parseDropUserStmt() ast.StmtNode {
	r.expect(drop)
	r.expect(user)
	ifExists := r.parseIfExists()
	return &ast.DropUserStmt{IsDropRole: false, IfExists: ifExists, UserList: r.parseUsernameList()}
}

// parseDropRoleStmt implements:
//
//	DropRoleStmt: "DROP" "ROLE" ["IF" "EXISTS"] RolenameList
func (r *rdParser) parseDropRoleStmt() ast.StmtNode {
	r.expect(drop)
	r.expect(role)
	ifExists := r.parseIfExists()
	tmp := make([]*auth.UserIdentity, 0, 10)
	roleList := r.parseRolenameList()
	for _, v := range roleList {
		tmp = append(tmp, &auth.UserIdentity{Username: v.Username, Hostname: v.Hostname})
	}
	return &ast.DropUserStmt{IsDropRole: true, IfExists: ifExists, UserList: tmp}
}

// parseDropStatsStmt implements:
//
//	DropStatsStmt: "DROP" "STATS" TableNameList
//	  | "DROP" "STATS" TableName "PARTITION" PartitionNameList
//	  | "DROP" "STATS" TableName "GLOBAL"
func (r *rdParser) parseDropStatsStmt() ast.StmtNode {
	r.expect(drop)
	r.expect(stats)
	tableNames := r.parseTableNameList()
	switch r.tok() {
	case partition:
		if len(tableNames) != 1 {
			r.syntaxError()
		}
		r.advance()
		// PartitionNameList: Identifier | PartitionNameList ',' Identifier
		partitionNames := []ast.CIStr{ast.NewCIStr(r.parseIdentifier())}
		for r.accept(int(',')) {
			partitionNames = append(partitionNames, ast.NewCIStr(r.parseIdentifier()))
		}
		r.cur() // warning position: the reduce-time lookahead
		r.sc.AppendError(ErrWarnDeprecatedSyntaxNoReplacement.FastGenByArgs("'DROP STATS ... PARTITION ...'", ""))
		r.sc.lastErrorAsWarn()
		return &ast.DropStatsStmt{
			Tables:         tableNames,
			PartitionNames: partitionNames,
		}
	case global:
		if len(tableNames) != 1 {
			r.syntaxError()
		}
		r.advance()
		r.cur() // warning position: the reduce-time lookahead
		r.sc.AppendError(ErrWarnDeprecatedSyntax.FastGenByArgs("DROP STATS ... GLOBAL", "DROP STATS ..."))
		r.sc.lastErrorAsWarn()
		return &ast.DropStatsStmt{
			Tables:        tableNames,
			IsGlobalStats: true,
		}
	}
	return &ast.DropStatsStmt{Tables: tableNames}
}
