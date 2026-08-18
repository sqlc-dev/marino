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

// The MySQL XA transaction statements (MySQL 26.7 §15.3.8). These
// postdate the goyacc grammar; the productions below are written from
// the MySQL 26.7 reference manual in the same style as parser.y. (The
// LOCK INSTANCE FOR BACKUP / UNLOCK INSTANCE statements of §15.3.5 are
// alternatives of the LOCK/UNLOCK statement families in parse_misc.go.)

import (
	"github.com/sqlc-dev/marino/ast"
)

func init() {
	rdRegister(xa, (*rdParser).parseXAStmt)
}

// parseXAStmt implements XAStmt:
//
//	"XA" ("START" | "BEGIN") XID ["JOIN" | "RESUME"]
//	| "XA" "END" XID ["SUSPEND" ["FOR" "MIGRATE"]]
//	| "XA" "PREPARE" XID
//	| "XA" "COMMIT" XID ["ONE" "PHASE"]
//	| "XA" "ROLLBACK" XID
//	| "XA" "RECOVER" ["CONVERT" "XID"]
//
// The XA BEGIN spelling parses to XAOpStart, so Restore() canonicalizes
// it to XA START.
func (r *rdParser) parseXAStmt() ast.StmtNode {
	r.expect(xa)
	switch r.tok() {
	case start, begin:
		r.advance()
		stmt := &ast.XAStmt{Op: ast.XAOpStart, Xid: r.parseXID()}
		switch r.tok() {
		case join:
			r.advance()
			stmt.StartOption = ast.XAStartJoin
		case resume:
			r.advance()
			stmt.StartOption = ast.XAStartResume
		}
		return stmt
	case end:
		r.advance()
		stmt := &ast.XAStmt{Op: ast.XAOpEnd, Xid: r.parseXID()}
		if r.accept(suspend) {
			stmt.Suspend = true
			if r.accept(forKwd) {
				r.expect(migrate)
				stmt.ForMigrate = true
			}
		}
		return stmt
	case prepare:
		r.advance()
		return &ast.XAStmt{Op: ast.XAOpPrepare, Xid: r.parseXID()}
	case commit:
		r.advance()
		stmt := &ast.XAStmt{Op: ast.XAOpCommit, Xid: r.parseXID()}
		if r.accept(one) {
			r.expect(phase)
			stmt.OnePhase = true
		}
		return stmt
	case rollback:
		r.advance()
		return &ast.XAStmt{Op: ast.XAOpRollback, Xid: r.parseXID()}
	case recover:
		r.advance()
		stmt := &ast.XAStmt{Op: ast.XAOpRecover}
		if r.accept(convert) {
			r.expect(xid)
			stmt.ConvertXID = true
		}
		return stmt
	}
	r.syntaxError()
	return nil
}

// parseXID implements the xid production:
// stringLit [',' stringLit [',' NUM]].
func (r *rdParser) parseXID() *ast.XID {
	x := &ast.XID{Gtrid: r.expect(stringLit).lit}
	if r.accept(int(',')) {
		x.Bqual = r.expect(stringLit).lit
		x.HasBqual = true
		if r.accept(int(',')) {
			x.FormatID = getUint64FromNUM(r.expect(intLit).item)
			x.HasFormatID = true
		}
	}
	return x
}
