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

// FLASHBACK statements (FlashbackTableStmt, FlashbackDatabaseStmt,
// FlashbackToTimestampStmt) and RecoverTableStmt. The "TO TIMESTAMP" /
// "TO TSO" spellings arrive as the single tokens toTimestamp / toTSO:
// the lexer only fuses them when the right literal follows.

import (
	"github.com/sqlc-dev/marino/ast"
)

func init() {
	rdRegister(flashback, (*rdParser).parseFlashbackStmtFamily)
	rdRegister(recover, (*rdParser).parseRecoverTableStmt)
}

// parseFlashbackStmtFamily implements:
//
//	FlashbackToTimestampStmt:
//	    "FLASHBACK" "CLUSTER" toTimestamp stringLit
//	  | "FLASHBACK" "TABLE" TableNameList toTimestamp stringLit
//	  | "FLASHBACK" DatabaseSym DBName toTimestamp stringLit
//	  | "FLASHBACK" "CLUSTER" toTSO LengthNum
//	  | "FLASHBACK" "TABLE" TableNameList toTSO LengthNum
//	  | "FLASHBACK" DatabaseSym DBName toTSO LengthNum
//	FlashbackTableStmt: "FLASHBACK" "TABLE" TableName FlashbackToNewName
//	FlashbackDatabaseStmt: "FLASHBACK" DatabaseSym DBName FlashbackToNewName
func (r *rdParser) parseFlashbackStmtFamily() ast.StmtNode {
	r.expect(flashback)
	switch r.tok() {
	case cluster:
		r.advance()
		switch r.tok() {
		case toTimestamp:
			r.advance()
			return &ast.FlashBackToTimestampStmt{
				FlashbackTS:  ast.NewValueExpr(r.expect(stringLit).lit, "", ""),
				FlashbackTSO: 0,
			}
		case toTSO:
			r.advance()
			tsoValue := r.parseLengthNum()
			if tsoValue > 0 {
				return &ast.FlashBackToTimestampStmt{
					FlashbackTSO: tsoValue,
				}
			}
			r.actionError(r.actionErrorf("Invalid TSO value provided: %d", tsoValue))
		}
		r.syntaxError()
	case tableKwd:
		r.advance()
		tables := r.parseTableNameList()
		switch r.tok() {
		case toTimestamp:
			r.advance()
			return &ast.FlashBackToTimestampStmt{
				Tables:       tables,
				FlashbackTS:  ast.NewValueExpr(r.expect(stringLit).lit, "", ""),
				FlashbackTSO: 0,
			}
		case toTSO:
			r.advance()
			tsoValue := r.parseLengthNum()
			if tsoValue > 0 {
				return &ast.FlashBackToTimestampStmt{
					Tables:       tables,
					FlashbackTSO: tsoValue,
				}
			}
			r.actionError(r.actionErrorf("Invalid TSO value provided: %d", tsoValue))
		}
		// FlashbackTableStmt takes a single TableName.
		if len(tables) != 1 {
			r.syntaxError()
		}
		return &ast.FlashBackTableStmt{
			Table:   tables[0],
			NewName: r.parseFlashbackToNewName(),
		}
	case database:
		r.advance()
		dbName := r.parseIdentifier()
		switch r.tok() {
		case toTimestamp:
			r.advance()
			return &ast.FlashBackToTimestampStmt{
				DBName:       ast.NewCIStr(dbName),
				FlashbackTS:  ast.NewValueExpr(r.expect(stringLit).lit, "", ""),
				FlashbackTSO: 0,
			}
		case toTSO:
			r.advance()
			tsoValue := r.parseLengthNum()
			if tsoValue > 0 {
				return &ast.FlashBackToTimestampStmt{
					DBName:       ast.NewCIStr(dbName),
					FlashbackTSO: tsoValue,
				}
			}
			r.actionError(r.actionErrorf("Invalid TSO value provided: %d", tsoValue))
		}
		return &ast.FlashBackDatabaseStmt{
			DBName:  ast.NewCIStr(dbName),
			NewName: r.parseFlashbackToNewName(),
		}
	}
	r.syntaxError()
	return nil
}

// parseFlashbackToNewName implements FlashbackToNewName:
// empty | "TO" Identifier.
func (r *rdParser) parseFlashbackToNewName() string {
	if r.accept(to) {
		return r.parseIdentifier()
	}
	return ""
}

// parseRecoverTableStmt implements:
//
//	RecoverTableStmt: "RECOVER" "TABLE" "BY" "JOB" Int64Num
//	  | "RECOVER" "TABLE" TableName
//	  | "RECOVER" "TABLE" TableName Int64Num
func (r *rdParser) parseRecoverTableStmt() ast.StmtNode {
	r.expect(recover)
	r.expect(tableKwd)
	if r.tok() == by {
		r.advance()
		r.expect(job)
		return &ast.RecoverTableStmt{
			JobID: r.parseInt64Num(),
		}
	}
	table := r.parseTableName()
	if r.tok() == intLit {
		return &ast.RecoverTableStmt{
			Table:  table,
			JobNum: r.parseInt64Num(),
		}
	}
	return &ast.RecoverTableStmt{
		Table: table,
	}
}
