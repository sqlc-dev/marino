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

// SQL plan bindings: CreateBindingStmt, DropBindingStmt, SetBindingStmt,
// with BindableStmt and StringLitOrUserVariable(List). The actions
// record each bindable statement's source text: a statement followed by
// "USING" is sliced up to the USING token's offset, the last one runs to
// the end of the source, both trimmed with strings.TrimSpace.

import (
	"strings"

	"github.com/sqlc-dev/marino/ast"
)

// parseGlobalScope implements GlobalScope:
// empty (false) | "GLOBAL" (true) | "SESSION" (false).
func (r *rdParser) parseGlobalScope() bool {
	switch r.tok() {
	case global:
		r.advance()
		return true
	case session:
		r.advance()
	}
	return false
}

// parseBindableStmt implements BindableStmt:
//
//	SetOprStmt | SelectStmt | SelectStmtWithClause | SubSelect
//	| UpdateStmt | DeleteWithoutUsingStmt | InsertIntoStmt
//	| ReplaceIntoStmt
//
// The SubSelect alternative carries the usual IsInBraces action, which
// finishSelectFamily reproduces. UpdateStmt includes a WithClause form;
// DeleteWithoutUsingStmt has none, and the USING form of DELETE is not
// bindable at all.
func (r *rdParser) parseBindableStmt() ast.StmtNode {
	switch r.tok() {
	case selectKwd, tableKwd, values, int('('):
		return r.finishSelectFamily(nil)
	case with:
		withClause := r.parseWithClause()
		if r.tok() == update {
			return r.parseUpdateStmt(withClause)
		}
		return r.finishSelectFamily(withClause)
	case update:
		return r.parseUpdateStmt(nil)
	case deleteKwd:
		stmt := r.parseDeleteStmt(nil)
		if d := stmt.(*ast.DeleteStmt); d.IsMultiTable && !d.BeforeFrom {
			// DeleteWithUsingStmt is not a BindableStmt alternative.
			r.syntaxError()
		}
		return stmt
	case insert:
		return r.parseInsertIntoStmt()
	case replace:
		return r.parseReplaceIntoStmt()
	}
	r.syntaxError()
	return nil
}

// parseCreateBindingStmt implements CreateBindingStmt:
//
//	"CREATE" GlobalScope "BINDING" "FOR" BindableStmt "USING" BindableStmt
//	| "CREATE" GlobalScope "BINDING" "USING" BindableStmt
//	| "CREATE" GlobalScope "BINDING" "FROM" "HISTORY" "USING" "PLAN"
//	  "DIGEST" StringLitOrUserVariableList
func (r *rdParser) parseCreateBindingStmt() ast.StmtNode {
	r.expect(create)
	globalScope := r.parseGlobalScope()
	r.expect(binding)
	switch r.tok() {
	case forKwd:
		r.advance()
		startOffset := r.cur().offset
		originStmt := r.parseBindableStmt()
		usingTok := r.expect(using)
		r.p.setNodeText(originStmt, strings.TrimSpace(r.src[startOffset:usingTok.offset]))

		startOffset = r.cur().offset
		hintedStmt := r.parseBindableStmt()
		r.p.setNodeText(hintedStmt, strings.TrimSpace(r.src[startOffset:]))

		return &ast.CreateBindingStmt{
			OriginNode:  originStmt,
			HintedNode:  hintedStmt,
			GlobalScope: globalScope,
		}
	case using:
		r.advance()
		startOffset := r.cur().offset
		hintedStmt := r.parseBindableStmt()
		r.p.setNodeText(hintedStmt, strings.TrimSpace(r.src[startOffset:]))

		return &ast.CreateBindingStmt{
			OriginNode:  hintedStmt,
			HintedNode:  hintedStmt,
			GlobalScope: globalScope,
		}
	case from:
		r.advance()
		r.expect(history)
		r.expect(using)
		r.expect(plan)
		r.expect(digest)
		return &ast.CreateBindingStmt{
			GlobalScope: globalScope,
			PlanDigests: r.parseStringLitOrUserVariableList(),
		}
	}
	r.syntaxError()
	return nil
}

// parseDropBindingStmt implements DropBindingStmt:
//
//	"DROP" GlobalScope "BINDING" "FOR" BindableStmt
//	| "DROP" GlobalScope "BINDING" "FOR" BindableStmt "USING" BindableStmt
//	| "DROP" GlobalScope "BINDING" "FOR" "SQL" "DIGEST"
//	  StringLitOrUserVariableList
func (r *rdParser) parseDropBindingStmt() ast.StmtNode {
	r.expect(drop)
	globalScope := r.parseGlobalScope()
	r.expect(binding)
	r.expect(forKwd)
	if r.tok() == sql {
		// "SQL" is a reserved word, so no BindableStmt starts with it.
		r.advance()
		r.expect(digest)
		return &ast.DropBindingStmt{
			GlobalScope: globalScope,
			SQLDigests:  r.parseStringLitOrUserVariableList(),
		}
	}
	startOffset := r.cur().offset
	originStmt := r.parseBindableStmt()
	if r.tok() == using {
		usingTok := r.expect(using)
		r.p.setNodeText(originStmt, strings.TrimSpace(r.src[startOffset:usingTok.offset]))

		startOffset = r.cur().offset
		hintedStmt := r.parseBindableStmt()
		r.p.setNodeText(hintedStmt, strings.TrimSpace(r.src[startOffset:]))

		return &ast.DropBindingStmt{
			OriginNode:  originStmt,
			HintedNode:  hintedStmt,
			GlobalScope: globalScope,
		}
	}
	r.p.setNodeText(originStmt, strings.TrimSpace(r.src[startOffset:]))

	return &ast.DropBindingStmt{
		OriginNode:  originStmt,
		GlobalScope: globalScope,
	}
}

// parseSetBindingStmt implements SetBindingStmt:
//
//	"SET" "BINDING" BindingStatusType "FOR" BindableStmt
//	| "SET" "BINDING" BindingStatusType "FOR" BindableStmt "USING"
//	  BindableStmt
//	| "SET" "BINDING" BindingStatusType "FOR" "SQL" "DIGEST" stringLit
func (r *rdParser) parseSetBindingStmt() ast.StmtNode {
	r.expect(set)
	r.expect(binding)
	// BindingStatusType: "ENABLED" | "DISABLED"
	var statusType ast.BindingStatusType
	switch r.tok() {
	case enabled:
		statusType = ast.BindingStatusTypeEnabled
	case disabled:
		statusType = ast.BindingStatusTypeDisabled
	default:
		r.syntaxError()
	}
	r.advance()
	r.expect(forKwd)
	if r.tok() == sql {
		r.advance()
		r.expect(digest)
		return &ast.SetBindingStmt{
			BindingStatusType: statusType,
			SQLDigest:         r.expect(stringLit).lit,
		}
	}
	startOffset := r.cur().offset
	originStmt := r.parseBindableStmt()
	if r.tok() == using {
		usingTok := r.expect(using)
		r.p.setNodeText(originStmt, strings.TrimSpace(r.src[startOffset:usingTok.offset]))

		startOffset = r.cur().offset
		hintedStmt := r.parseBindableStmt()
		r.p.setNodeText(hintedStmt, strings.TrimSpace(r.src[startOffset:]))

		return &ast.SetBindingStmt{
			BindingStatusType: statusType,
			OriginNode:        originStmt,
			HintedNode:        hintedStmt,
		}
	}
	r.p.setNodeText(originStmt, strings.TrimSpace(r.src[startOffset:]))

	return &ast.SetBindingStmt{
		BindingStatusType: statusType,
		OriginNode:        originStmt,
	}
}

// parseStringLitOrUserVariableList implements StringLitOrUserVariableList:
// StringLitOrUserVariable (',' StringLitOrUserVariable)*.
func (r *rdParser) parseStringLitOrUserVariableList() []*ast.StringOrUserVar {
	list := []*ast.StringOrUserVar{r.parseStringLitOrUserVariable()}
	for r.accept(int(',')) {
		list = append(list, r.parseStringLitOrUserVariable())
	}
	return list
}

// parseStringLitOrUserVariable implements StringLitOrUserVariable:
// stringLit | UserVariable.
func (r *rdParser) parseStringLitOrUserVariable() *ast.StringOrUserVar {
	if r.tok() == stringLit {
		return &ast.StringOrUserVar{StringLit: r.expect(stringLit).lit}
	}
	return &ast.StringOrUserVar{UserVar: r.parseUserVariable()}
}
