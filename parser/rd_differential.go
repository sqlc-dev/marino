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

// The differential harness of the parser rewrite (PLAN.md): while the
// goyacc parser is in-tree it is the oracle, and under `go test` every
// input the recursive-descent parser handles is re-parsed by goyacc and
// the two results compared field-for-field (internal/dump renderings,
// error strings, warning strings). A divergence panics, so the entire
// existing test suite doubles as a differential corpus. rdDifferential is
// switched on by an init in the package's internal test file.

import (
	"fmt"
	"strings"

	"github.com/sqlc-dev/marino/ast"
	astdump "github.com/sqlc-dev/marino/internal/dump"
)

var rdDifferential bool

// rdDifferentialMaxLen bounds the inputs the in-process check re-parses:
// the suite asserts allocation ceilings on multi-kilobyte statements, and
// a second parse plus two tree dumps would triple-count them. Inputs this
// size are exercised differentially out of band by cmd/difftest instead.
const rdDifferentialMaxLen = 2500

// rdDifferentialCheck re-parses sql with the goyacc oracle on a fresh
// Parser carrying the same configuration and panics on any divergence
// from the RD result.
func (parser *Parser) rdDifferentialCheck(sql string, params []ParseParam, stmts []ast.StmtNode, warns []error, rdErr error) {
	oracle := New()
	oracle.lexer.sqlMode = parser.lexer.sqlMode
	oracle.lexer.supportWindowFunc = parser.lexer.supportWindowFunc
	oracle.lexer.skipPositionRecording = parser.lexer.skipPositionRecording
	oracle.lexer.keepHint = parser.lexer.keepHint
	oracle.enableMariaDB = parser.enableMariaDB
	oracle.strictDoubleFieldType = parser.strictDoubleFieldType
	oStmts, oWarns, oErr := oracle.yaccParseSQL(sql, params...)

	fail := func(aspect, rd, yacc string) {
		panic(fmt.Sprintf("rd/yacc parser divergence on %s\nSQL: %s\n%s", aspect, sql, firstDiff(rd, yacc)))
	}
	if (rdErr == nil) != (oErr == nil) || (rdErr != nil && rdErr.Error() != oErr.Error()) {
		fail("error", fmt.Sprintf("%v", rdErr), fmt.Sprintf("%v", oErr))
	}
	if len(warns) != len(oWarns) {
		fail("warning count", fmt.Sprintf("%v", warns), fmt.Sprintf("%v", oWarns))
	}
	for i := range warns {
		if warns[i].Error() != oWarns[i].Error() {
			fail("warning", warns[i].Error(), oWarns[i].Error())
		}
	}
	rdDump, oDump := astdump.Statements(stmts), astdump.Statements(oStmts)
	if rdDump != oDump {
		fail("AST", rdDump, oDump)
	}
}

// firstDiff renders the first run of differing lines between two dumps,
// with a little context.
func firstDiff(a, b string) string {
	al, bl := strings.Split(a, "\n"), strings.Split(b, "\n")
	n := len(al)
	if len(bl) < n {
		n = len(bl)
	}
	i := 0
	for i < n && al[i] == bl[i] {
		i++
	}
	var sb strings.Builder
	lo := i - 3
	if lo < 0 {
		lo = 0
	}
	for k := lo; k < i; k++ {
		fmt.Fprintf(&sb, "  %s\n", al[k])
	}
	for k := i; k < i+8; k++ {
		if k < len(al) {
			fmt.Fprintf(&sb, "rd  > %s\n", al[k])
		}
	}
	for k := i; k < i+8; k++ {
		if k < len(bl) {
			fmt.Fprintf(&sb, "yacc> %s\n", bl[k])
		}
	}
	return sb.String()
}
