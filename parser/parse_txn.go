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

// Transaction control statements: BEGIN/START TRANSACTION, COMMIT,
// ROLLBACK, SAVEPOINT, and RELEASE SAVEPOINT.

import (
	"github.com/sqlc-dev/marino/ast"
)

func init() {
	rdRegister(begin, (*rdParser).parseBeginTransactionStmt)
	rdRegister(start, (*rdParser).parseBeginTransactionStmt)
	rdRegister(commit, (*rdParser).parseCommitStmt)
	rdRegister(rollback, (*rdParser).parseRollbackStmt)
	rdRegister(savepoint, (*rdParser).parseSavepointStmt)
	rdRegister(release, (*rdParser).parseReleaseSavepointStmt)
}

// parseAsOfClause implements AsOfClause: asof "TIMESTAMP" Expression.
func (r *rdParser) parseAsOfClause() *ast.AsOfClause {
	r.expect(asof)
	r.expect(timestampType)
	return &ast.AsOfClause{TsExpr: r.parseExpression()}
}

// parseBeginTransactionStmt implements BeginTransactionStmt (both the
// "BEGIN" and "START TRANSACTION" alternatives).
func (r *rdParser) parseBeginTransactionStmt() ast.StmtNode {
	if r.accept(begin) {
		switch r.tok() {
		case pessimistic:
			// "BEGIN" "PESSIMISTIC"
			r.advance()
			return &ast.BeginStmt{Mode: ast.Pessimistic}
		case optimistic:
			// "BEGIN" "OPTIMISTIC"
			r.advance()
			return &ast.BeginStmt{Mode: ast.Optimistic}
		}
		// "BEGIN"
		return &ast.BeginStmt{}
	}
	r.expect(start)
	r.expect(transaction)
	switch r.tok() {
	case read:
		r.advance()
		if r.accept(write) {
			// "START" "TRANSACTION" "READ" "WRITE"
			return &ast.BeginStmt{}
		}
		r.expect(only)
		if r.tok() == asof {
			// "START" "TRANSACTION" "READ" "ONLY" AsOfClause
			return &ast.BeginStmt{
				ReadOnly: true,
				AsOf:     r.parseAsOfClause(),
			}
		}
		// "START" "TRANSACTION" "READ" "ONLY"
		return &ast.BeginStmt{
			ReadOnly: true,
		}
	case with:
		r.advance()
		if r.accept(consistent) {
			// "START" "TRANSACTION" "WITH" "CONSISTENT" "SNAPSHOT"
			r.expect(snapshot)
			return &ast.BeginStmt{}
		}
		// "START" "TRANSACTION" "WITH" "CAUSAL" "CONSISTENCY" "ONLY"
		r.expect(causal)
		r.expect(consistency)
		r.expect(only)
		return &ast.BeginStmt{
			CausalConsistencyOnly: true,
		}
	}
	// "START" "TRANSACTION"
	return &ast.BeginStmt{}
}

// parseCompletionType implements CompletionTypeWithinTransaction.
func (r *rdParser) parseCompletionType() ast.CompletionType {
	switch r.tok() {
	case and:
		r.advance()
		if r.accept(no) {
			r.expect(chain)
			if r.accept(no) {
				// "AND" "NO" "CHAIN" "NO" "RELEASE"
				r.expect(release)
				return ast.CompletionTypeDefault
			}
			if r.accept(release) {
				// "AND" "NO" "CHAIN" "RELEASE"
				return ast.CompletionTypeRelease
			}
			// "AND" "NO" "CHAIN"
			return ast.CompletionTypeDefault
		}
		r.expect(chain)
		if r.accept(no) {
			// "AND" "CHAIN" "NO" "RELEASE"
			r.expect(release)
		}
		// "AND" "CHAIN"
		return ast.CompletionTypeChain
	case release:
		r.advance()
		return ast.CompletionTypeRelease
	case no:
		// "NO" "RELEASE"
		r.advance()
		r.expect(release)
		return ast.CompletionTypeDefault
	}
	r.syntaxError()
	return ast.CompletionTypeDefault
}

// parseCommitStmt implements CommitStmt.
func (r *rdParser) parseCommitStmt() ast.StmtNode {
	r.expect(commit)
	switch r.tok() {
	case and, release, no:
		// "COMMIT" CompletionTypeWithinTransaction
		return &ast.CommitStmt{CompletionType: r.parseCompletionType()}
	}
	return &ast.CommitStmt{}
}

// parseRollbackStmt implements RollbackStmt.
func (r *rdParser) parseRollbackStmt() ast.StmtNode {
	r.expect(rollback)
	switch r.tok() {
	case to:
		r.advance()
		if r.tok() == savepoint && isIdentifierTok(r.la(1)) {
			// "ROLLBACK" "TO" "SAVEPOINT" Identifier
			r.advance()
			return &ast.RollbackStmt{SavepointName: r.parseIdentifier()}
		}
		// "ROLLBACK" "TO" Identifier
		return &ast.RollbackStmt{SavepointName: r.parseIdentifier()}
	case and, release, no:
		// "ROLLBACK" CompletionTypeWithinTransaction
		return &ast.RollbackStmt{CompletionType: r.parseCompletionType()}
	}
	return &ast.RollbackStmt{}
}

// parseSavepointStmt implements SavepointStmt: "SAVEPOINT" Identifier.
func (r *rdParser) parseSavepointStmt() ast.StmtNode {
	r.expect(savepoint)
	return &ast.SavepointStmt{Name: r.parseIdentifier()}
}

// parseReleaseSavepointStmt implements ReleaseSavepointStmt:
// "RELEASE" "SAVEPOINT" Identifier.
func (r *rdParser) parseReleaseSavepointStmt() ast.StmtNode {
	r.expect(release)
	r.expect(savepoint)
	return &ast.ReleaseSavepointStmt{Name: r.parseIdentifier()}
}
