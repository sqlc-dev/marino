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

// The HANDLER statements (MySQL 26.7 §15.2.5). These postdate the
// goyacc grammar; the productions below are written from the MySQL 26.7
// reference manual in the same style as parser.y.

import (
	"github.com/sqlc-dev/marino/ast"
)

func init() {
	rdRegister(handler, (*rdParser).parseHandlerStmt)
}

// parseHandlerStmt implements the HANDLER statement family:
//
//	"HANDLER" TableName "OPEN" [["AS"] Identifier]
//	| "HANDLER" TableName "READ" HandlerReadTail
//	| "HANDLER" TableName "CLOSE"
func (r *rdParser) parseHandlerStmt() ast.StmtNode {
	r.expect(handler)
	table := r.parseTableName()
	switch r.tok() {
	case open:
		r.advance()
		stmt := &ast.HandlerOpenStmt{Table: table}
		if r.accept(as) {
			stmt.Alias = ast.NewCIStr(r.parseIdentifier())
		} else if isIdentifierTok(r.tok()) {
			stmt.Alias = ast.NewCIStr(r.parseIdentifier())
		}
		return stmt
	case read:
		r.advance()
		return r.parseHandlerReadTail(table)
	case close:
		r.advance()
		return &ast.HandlerCloseStmt{Table: table}
	}
	r.syntaxError()
	return nil
}

// parseHandlerReadTail implements the tail of HANDLER ... READ:
//
//	Identifier HandlerCompareOp '(' ExpressionListOpt ')' HandlerReadTail2
//	| Identifier ("FIRST" | "NEXT" | "PREV" | "LAST") HandlerReadTail2
//	| ("FIRST" | "NEXT") HandlerReadTail2
//
// with HandlerReadTail2: ["WHERE" Expression] [LimitClause]. An index
// name may also be "PRIMARY", naming the primary key.
func (r *rdParser) parseHandlerReadTail(table *ast.TableName) ast.StmtNode {
	stmt := &ast.HandlerReadStmt{Table: table}
	switch r.tok() {
	case first:
		r.advance()
		stmt.Direction = ast.HandlerReadFirst
	case next:
		r.advance()
		stmt.Direction = ast.HandlerReadNext
	default:
		stmt.IndexName = r.parseCacheIndexName()
		switch r.tok() {
		case first:
			r.advance()
			stmt.Direction = ast.HandlerReadFirst
		case next:
			r.advance()
			stmt.Direction = ast.HandlerReadNext
		case prev:
			r.advance()
			stmt.Direction = ast.HandlerReadPrev
		case last:
			r.advance()
			stmt.Direction = ast.HandlerReadLast
		case eq:
			r.advance()
			stmt.Op = ast.HandlerOpEQ
		case le:
			r.advance()
			stmt.Op = ast.HandlerOpLE
		case ge:
			r.advance()
			stmt.Op = ast.HandlerOpGE
		case int('<'):
			r.advance()
			stmt.Op = ast.HandlerOpLT
		case int('>'):
			r.advance()
			stmt.Op = ast.HandlerOpGT
		default:
			r.syntaxError()
		}
		if stmt.Op != ast.HandlerOpNone {
			r.expect(int('('))
			stmt.Values = r.parseExpressionListOpt()
			r.expect(int(')'))
		}
	}
	if r.accept(where) {
		stmt.Where = r.parseExpression()
	}
	if r.tok() == limit {
		stmt.Limit = r.parseSelectStmtLimit()
	}
	return stmt
}
