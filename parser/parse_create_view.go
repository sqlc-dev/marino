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

// CREATE VIEW: CreateViewStmt with OrReplace, ViewAlgorithm,
// ViewDefiner, ViewSQLSecurity, ViewName, ViewFieldList,
// CreateViewSelectOpt, and ViewCheckOption.

import (
	"strings"

	"github.com/sqlc-dev/marino/ast"
	"github.com/sqlc-dev/marino/auth"
)

// parseCreateViewStmt implements:
//
//	CreateViewStmt: "CREATE" OrReplace ViewAlgorithm ViewDefiner
//	    ViewSQLSecurity "VIEW" ViewName ViewFieldList "AS"
//	    CreateViewSelectOpt ViewCheckOption
//
// The action records the select statement's source text: startOffset is
// the first token of CreateViewSelectOpt; endOffset is the reduce-time
// lookahead (parser.yylval.offset), moved back to the start of the
// ViewCheckOption when one is present — in both cases the offset of the
// token that follows the select.
func (r *rdParser) parseCreateViewStmt() ast.StmtNode {
	r.expect(create)
	// OrReplace
	orReplace := false
	if r.tok() == or {
		r.advance()
		r.expect(replace)
		orReplace = true
	}
	// ViewAlgorithm
	algorithmV := ast.AlgorithmUndefined
	if r.tok() == algorithm {
		r.advance()
		r.expect(eq)
		switch r.tok() {
		case undefined:
			algorithmV = ast.AlgorithmUndefined
		case merge:
			algorithmV = ast.AlgorithmMerge
		case temptable:
			algorithmV = ast.AlgorithmTemptable
		default:
			r.syntaxError()
		}
		r.advance()
	}
	// ViewDefiner
	definerV := &auth.UserIdentity{CurrentUser: true}
	if r.tok() == definer {
		r.advance()
		r.expect(eq)
		definerV = r.parseUsername()
	}
	// ViewSQLSecurity
	securityV := ast.SecurityDefiner
	if r.tok() == sql {
		r.advance()
		r.expect(security)
		switch r.tok() {
		case definer:
			securityV = ast.SecurityDefiner
		case invoker:
			securityV = ast.SecurityInvoker
		default:
			r.syntaxError()
		}
		r.advance()
	}
	r.expect(view)
	// ViewName: TableName
	viewName := r.parseTableName()
	// ViewFieldList: empty | '(' ColumnList ')'
	var cols []ast.CIStr
	if r.tok() == int('(') {
		r.advance()
		cols = []ast.CIStr{ast.NewCIStr(r.parseIdentifier())}
		for r.accept(int(',')) {
			cols = append(cols, ast.NewCIStr(r.parseIdentifier()))
		}
		r.expect(int(')'))
	}
	r.expect(as)
	startOffset := r.cur().offset
	selStmt := r.parseCreateViewSelectOpt()
	endOffset := r.cur().offset
	x := &ast.CreateViewStmt{
		OrReplace: orReplace,
		ViewName:  viewName,
		Select:    selStmt,
		Algorithm: algorithmV,
		Definer:   definerV,
		Security:  securityV,
	}
	if cols != nil {
		x.Cols = cols
	}
	// ViewCheckOption: empty | "WITH" ("CASCADED"|"LOCAL") "CHECK" "OPTION"
	if r.tok() == with {
		r.advance()
		switch r.tok() {
		case cascaded:
			x.CheckOption = ast.CheckOptionCascaded
		case local:
			x.CheckOption = ast.CheckOptionLocal
		default:
			r.syntaxError()
		}
		r.advance()
		r.expect(check)
		r.expect(option)
	} else {
		x.CheckOption = ast.CheckOptionCascaded
	}
	r.p.setNodeText(selStmt, strings.TrimSpace(r.src[startOffset:endOffset]))
	return x
}

// parseCreateViewSelectOpt implements CreateViewSelectOpt:
// SetOprStmt | SelectStmt | SelectStmtWithClause | SubSelect,
// with the SubSelect alternative's IsInBraces action.
func (r *rdParser) parseCreateViewSelectOpt() ast.StmtNode {
	stmt, sub := r.parseSelectCore()
	if sub != nil {
		var sel ast.StmtNode
		switch x := sub.Query.(type) {
		case *ast.SelectStmt:
			x.IsInBraces = true
			sel = x
		case *ast.SetOprStmt:
			x.IsInBraces = true
			sel = x
		}
		return sel
	}
	return stmt
}
