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

import (
	"strings"

	"github.com/sqlc-dev/marino/ast"
)

func init() {
	rdRegister(change, (*rdParser).parseChangeReplicationSourceStmt)
}

// parseChangeReplicationSourceStmt implements ChangeReplicationSourceStmt.
// The statement postdates the goyacc grammar (parser.y has no production
// for it); the shape follows the MySQL 26.7 reference manual:
//
//	"CHANGE" "REPLICATION" "SOURCE" "TO" ReplicationSourceOption
//	    ("," ReplicationSourceOption)* ("FOR" "CHANNEL" stringLit)?
//
// Option names are not validated against the server's option list, so the
// MySQL 26.7 Change Stream Applier options (APPLIER_VERSION,
// APPLIER_WORKER_COUNT, APPLIER_EVENT_MEMORY_LIMIT) parse like any other.
// The removed CHANGE MASTER TO spelling is not parsed, matching MySQL 8.4+.
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
	if r.accept(forKwd) {
		r.expect(channel)
		stmt.Channel = r.expect(stringLit).lit
	}
	return stmt
}

// parseReplicationSourceOption implements ReplicationSourceOption:
// Identifier eq (stringLit | intLit | decLit | floatLit).
func (r *rdParser) parseReplicationSourceOption() *ast.ReplicationSourceOption {
	name := strings.ToUpper(r.parseIdentifier())
	r.expect(eq)
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
	return &ast.ReplicationSourceOption{Name: name, Value: value}
}
