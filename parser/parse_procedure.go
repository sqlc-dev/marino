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

// Stored procedures: CreateProcedureStmt with the procedure parameter
// list (OptSpPdparams/SpPdparams/SpPdparam/SpOptInout) and the whole
// procedure body grammar (ProcedureProcStmt and friends),
// DropProcedureStmt, and CallStmt/ProcedureCall.

import (
	"strings"

	"github.com/sqlc-dev/marino/ast"
)

func init() {
	rdRegister(call, (*rdParser).parseCallStmt)
}

// parseCallStmt implements CallStmt: "CALL" ProcedureCall.
func (r *rdParser) parseCallStmt() ast.StmtNode {
	r.expect(call)
	return &ast.CallStmt{
		Procedure: r.parseProcedureCall(),
	}
}

// parseProcedureCall implements ProcedureCall:
//
//	identifier | Identifier '.' Identifier
//	| identifier '(' ExpressionListOpt ')'
//	| Identifier '.' Identifier '(' ExpressionListOpt ')'
//
// An unqualified procedure name must be a plain identifier token; only
// the schema-qualified forms accept unreserved keywords (Identifier).
func (r *rdParser) parseProcedureCall() *ast.FuncCallExpr {
	start := r.cur().offset
	var fc *ast.FuncCallExpr
	if isIdentifierTok(r.tok()) && r.la(1) == int('.') {
		schema := r.parseIdentifier()
		r.expect(int('.'))
		fc = &ast.FuncCallExpr{
			Tp:     ast.FuncCallExprTypeGeneric,
			Schema: ast.NewCIStr(schema),
			FnName: ast.NewCIStr(r.parseIdentifier()),
			Args:   []ast.ExprNode{},
		}
	} else {
		fc = &ast.FuncCallExpr{
			Tp:     ast.FuncCallExprTypeGeneric,
			FnName: ast.NewCIStr(r.expect(identifier).lit),
			Args:   []ast.ExprNode{},
		}
	}
	if r.tok() == int('(') {
		r.advance()
		fc.Args = r.parseExpressionListOpt()
		r.expect(int(')'))
	}
	return r.setOrigin(fc, start).(*ast.FuncCallExpr)
}

// parseCreateProcedureStmt implements:
//
//	CreateProcedureStmt: "CREATE" DefinerOpt "PROCEDURE" IfNotExists
//	    TableName '(' OptSpPdparams ')' ProcedureProcStmt
//
// The action records the body statement's source text — from the body's
// first token through the reduce-time lookahead (parser.yylval.offset) —
// and the parameter list's source text between the parentheses.
func (r *rdParser) parseCreateProcedureStmt() ast.StmtNode {
	r.expect(create)
	definerV := r.parseDefinerOpt()
	r.expect(procedure)
	ifNotExists := r.parseIfNotExists()
	procName := r.parseTableName()
	lparen := r.expect(int('('))
	params := r.parseOptSpPdparams()
	rparen := r.expect(int(')'))
	bodyStart := r.cur().offset
	body := r.parseProcedureProcStmt()
	x := &ast.ProcedureInfo{
		IfNotExists:    ifNotExists,
		ProcedureName:  procName,
		ProcedureParam: params,
		ProcedureBody:  body,
		Definer:        definerV,
	}
	r.p.setNodeText(body, strings.TrimSpace(r.src[bodyStart:r.cur().offset]))
	startOffset := lparen.offset
	if r.src[startOffset] == '(' {
		startOffset++
	}
	x.ProcedureParamStr = strings.TrimSpace(r.src[startOffset:rparen.offset])
	return x
}

// parseDropProcedureStmt implements DropProcedureStmt:
// "DROP" "PROCEDURE" IfExists TableName.
func (r *rdParser) parseDropProcedureStmt() ast.StmtNode {
	r.expect(drop)
	r.expect(procedure)
	ifExists := r.parseIfExists()
	return &ast.DropProcedureStmt{
		IfExists:      ifExists,
		ProcedureName: r.parseTableName(),
	}
}

// parseOptSpPdparams implements OptSpPdparams: empty | SpPdparams,
// and SpPdparams: SpPdparam (',' SpPdparam)*.
func (r *rdParser) parseOptSpPdparams() []*ast.StoreParameter {
	params := []*ast.StoreParameter{}
	if r.tok() == int(')') {
		return params
	}
	params = append(params, r.parseSpPdparam())
	for r.accept(int(',')) {
		params = append(params, r.parseSpPdparam())
	}
	return params
}

// parseSpPdparam implements SpPdparam: SpOptInout Identifier Type,
// with SpOptInout: empty | "IN" | "OUT" | "INOUT".
func (r *rdParser) parseSpPdparam() *ast.StoreParameter {
	paramStatus := ast.MODE_IN
	switch r.tok() {
	case in:
		r.advance()
	case out:
		r.advance()
		paramStatus = ast.MODE_OUT
	case inout:
		r.advance()
		paramStatus = ast.MODE_INOUT
	}
	name := r.parseIdentifier()
	return &ast.StoreParameter{
		Paramstatus: paramStatus,
		ParamType:   r.parseType(),
		ParamName:   name,
	}
}

// parseProcedureProcStmt implements ProcedureProcStmt, dispatching on the
// leading token:
//
//	ProcedureStatementStmt | ProcedureUnlabeledBlock | ProcedureIfstmt
//	| ProcedureCaseStmt | ProcedureUnlabelLoopBlock | ProcedureOpenCur
//	| ProcedureCloseCur | ProcedureFetchInto | ProcedureLabeledBlock
//	| ProcedurelabeledLoopStmt | ProcedureIterate | ProcedureLeave
//
// Labels (and cursor names) are plain identifier tokens in the grammar,
// so the keyword-led alternatives never conflict with them.
func (r *rdParser) parseProcedureProcStmt() ast.StmtNode {
	switch r.tok() {
	case begin:
		// ProcedureUnlabeledBlock: ProcedureBlockContent
		return r.parseProcedureBlockContent()
	case ifKwd:
		return r.parseProcedureIfstmt()
	case caseKwd:
		return r.parseProcedureCaseStmt()
	case while, repeat, loop:
		// ProcedureUnlabelLoopBlock: ProcedureUnlabelLoopStmt
		return r.parseProcedureUnlabelLoopStmt()
	case open:
		// ProcedureOpenCur: "OPEN" identifier
		r.advance()
		return &ast.ProcedureOpenCur{CurName: strings.ToLower(r.expect(identifier).lit)}
	case close:
		// ProcedureCloseCur: "CLOSE" identifier
		r.advance()
		return &ast.ProcedureCloseCur{CurName: strings.ToLower(r.expect(identifier).lit)}
	case fetch:
		return r.parseProcedureFetchInto()
	case iterate:
		// ProcedureIterate: "ITERATE" identifier
		r.advance()
		return &ast.ProcedureJump{Name: r.expect(identifier).lit, IsLeave: false}
	case leave:
		// ProcedureLeave: "LEAVE" identifier
		r.advance()
		return &ast.ProcedureJump{Name: r.expect(identifier).lit, IsLeave: true}
	case returnKwd:
		// ReturnStmt: "RETURN" Expression — the routine_body form of
		// stored functions; MySQL rejects it elsewhere at a later stage.
		r.advance()
		return &ast.ReturnStmt{Expr: r.parseExpression()}
	case identifier:
		if r.la(1) == int(':') {
			return r.parseProcedureLabeled()
		}
	}
	return r.parseProcedureStatementStmt()
}

// parseProcedureStatementStmt implements ProcedureStatementStmt:
//
//	SelectStmt | SelectStmtWithClause | SubSelect | SetStmt | UpdateStmt
//	| UseStmt | InsertIntoStmt | ReplaceIntoStmt | CommitStmt
//	| RollbackStmt | ExplainStmt | SetOprStmt | DeleteFromStmt
//	| AnalyzeTableStmt | TruncateTableStmt | SignalStmt | ResignalStmt
//	| GetDiagnosticsStmt
//
// The SubSelect alternative carries the usual IsInBraces action, which
// finishSelectFamily reproduces.
func (r *rdParser) parseProcedureStatementStmt() ast.StmtNode {
	switch r.tok() {
	case selectKwd, tableKwd, values, int('('):
		return r.finishSelectFamily(nil)
	case with:
		// SelectStmtWithClause, or the WithClause forms of
		// UpdateStmt/DeleteFromStmt.
		withClause := r.parseWithClause()
		switch r.tok() {
		case update:
			return r.parseUpdateStmt(withClause)
		case deleteKwd:
			return r.parseDeleteStmt(withClause)
		}
		return r.finishSelectFamily(withClause)
	case set:
		return r.parseSetStmt()
	case update:
		return r.parseUpdateStmt(nil)
	case use:
		return r.parseUseStmt()
	case insert:
		return r.parseInsertIntoStmt()
	case replace:
		return r.parseReplaceIntoStmt()
	case commit:
		return r.parseCommitStmt()
	case rollback:
		return r.parseRollbackStmt()
	case explain, describe, desc:
		return r.parseExplainStmt()
	case deleteKwd:
		return r.parseDeleteStmt(nil)
	case analyze:
		// ProcedureStatementStmt: AnalyzeTableStmt
		return r.parseAnalyzeTableStmt()
	case truncate:
		return r.parseTruncateTableStmt()
	case signal:
		return r.parseSignalStmt()
	case resignal:
		return r.parseResignalStmt()
	case get:
		return r.parseGetDiagnosticsStmt()
	}
	r.syntaxError()
	return nil
}

// parseProcedureBlockContent implements ProcedureBlockContent:
// "BEGIN" ProcedureDeclsOpt ProcedureProcStmts "END", with
// ProcedureDeclsOpt/ProcedureDecls: (ProcedureDecl ';')* and
// ProcedureProcStmts: (ProcedureProcStmt ';')*.
func (r *rdParser) parseProcedureBlockContent() *ast.ProcedureBlock {
	r.expect(begin)
	decls := []ast.DeclNode{}
	for r.tok() == declare {
		decls = append(decls, r.parseProcedureDecl())
		r.expect(int(';'))
	}
	stmts := []ast.StmtNode{}
	for r.tok() != end && r.tok() != 0 {
		stmts = append(stmts, r.parseProcedureProcStmt())
		r.expect(int(';'))
	}
	r.expect(end)
	return &ast.ProcedureBlock{
		ProcedureVars:      decls,
		ProcedureProcStmts: stmts,
	}
}

// parseProcedureDecl implements ProcedureDecl:
//
//	"DECLARE" ProcedureDeclIdents Type ProcedureOptDefault
//	| "DECLARE" identifier "CURSOR" "FOR" ProcedureCursorSelectStmt
//	| "DECLARE" ProcedureHandlerType "HANDLER" "FOR" ProcedureHcondList
//	  ProcedureProcStmt
//
// CONTINUE/EXIT are reserved words, so the handler form is committed by
// its first token; the cursor form needs a plain identifier followed by
// the reserved word CURSOR.
func (r *rdParser) parseProcedureDecl() ast.DeclNode {
	r.expect(declare)
	switch {
	case r.tok() == continueKwd || r.tok() == exit:
		// ProcedureHandlerType: "CONTINUE" | "EXIT"
		handlerType := ast.PROCEDUR_CONTINUE
		if r.tok() == exit {
			handlerType = ast.PROCEDUR_EXIT
		}
		r.advance()
		r.expect(handler)
		r.expect(forKwd)
		conds := r.parseProcedureHcondList()
		return &ast.ProcedureErrorControl{
			ControlHandle: handlerType,
			ErrorCon:      conds,
			Operate:       r.parseProcedureProcStmt(),
		}
	case r.tok() == identifier && r.la(1) == cursor:
		name := strings.ToLower(r.expect(identifier).lit)
		r.expect(cursor)
		r.expect(forKwd)
		return &ast.ProcedureCursor{
			CurName:      name,
			Selectstring: r.parseProcedureCursorSelectStmt(),
		}
	case r.tok() == identifier && r.la(1) == condition:
		// "DECLARE" identifier "CONDITION" "FOR" ProcedurceCond
		name := strings.ToLower(r.expect(identifier).lit)
		r.expect(condition)
		r.expect(forKwd)
		decl := &ast.ProcedureConditionDecl{Name: name}
		if r.tok() == intLit {
			decl.Value = &ast.ProcedureErrorVal{
				ErrorNum: getUint64FromNUM(r.expect(intLit).item),
			}
		} else {
			r.expect(sqlstate)
			// optValue: empty | "VALUE"
			r.accept(value)
			decl.Value = &ast.ProcedureErrorState{
				CodeStatus: r.expect(stringLit).lit,
			}
		}
		return decl
	}
	// ProcedureDeclIdents: Identifier (',' Identifier)*
	names := []string{strings.ToLower(r.parseIdentifier())}
	for r.accept(int(',')) {
		names = append(names, strings.ToLower(r.parseIdentifier()))
	}
	x := &ast.ProcedureDecl{
		DeclNames: names,
		DeclType:  r.parseType(),
	}
	// ProcedureOptDefault: empty | "DEFAULT" Expression
	if r.accept(defaultKwd) {
		x.DeclDefault = r.parseExpression()
	}
	return x
}

// parseProcedureCursorSelectStmt implements ProcedureCursorSelectStmt:
// SelectStmt | SelectStmtWithClause | SubSelect | SetOprStmt (with the
// SubSelect IsInBraces action).
func (r *rdParser) parseProcedureCursorSelectStmt() ast.StmtNode {
	switch r.tok() {
	case selectKwd, tableKwd, values, int('('):
		return r.finishSelectFamily(nil)
	case with:
		return r.finishSelectFamily(r.parseWithClause())
	}
	r.syntaxError()
	return nil
}

// parseProcedureHcondList implements ProcedureHcondList:
// ProcedureHcond (',' ProcedureHcond)*.
func (r *rdParser) parseProcedureHcondList() []ast.ErrNode {
	conds := []ast.ErrNode{r.parseProcedureHcond()}
	for r.accept(int(',')) {
		conds = append(conds, r.parseProcedureHcond())
	}
	return conds
}

// parseProcedureHcond implements ProcedureHcond:
//
//	ProcedurceCond | "SQLWARNING" | "NOT" "FOUND" | "SQLEXCEPTION"
//
// and ProcedurceCond: NUM | "SQLSTATE" optValue stringLit.
func (r *rdParser) parseProcedureHcond() ast.ErrNode {
	switch r.tok() {
	case intLit:
		return &ast.ProcedureErrorVal{
			ErrorNum: getUint64FromNUM(r.expect(intLit).item),
		}
	case sqlstate:
		r.advance()
		// optValue: empty | "VALUE"
		r.accept(value)
		return &ast.ProcedureErrorState{
			CodeStatus: r.expect(stringLit).lit,
		}
	case sqlwarning:
		r.advance()
		return &ast.ProcedureErrorCon{
			ErrorCon: ast.PROCEDUR_SQLWARNING,
		}
	case not:
		r.advance()
		r.expect(found)
		return &ast.ProcedureErrorCon{
			ErrorCon: ast.PROCEDUR_NOT_FOUND,
		}
	case sqlexception:
		r.advance()
		return &ast.ProcedureErrorCon{
			ErrorCon: ast.PROCEDUR_SQLEXCEPTION,
		}
	case identifier:
		// A condition name declared with DECLARE ... CONDITION.
		return &ast.ProcedureErrorName{
			Name: strings.ToLower(r.expect(identifier).lit),
		}
	}
	r.syntaxError()
	return nil
}

// parseProcedureFetchInto implements ProcedureFetchInto:
// "FETCH" ProcedureOptFetchNo identifier "INTO" ProcedureFetchList, with
// ProcedureOptFetchNo: empty | "NEXT" "FROM" | "FROM" and
// ProcedureFetchList: identifier (',' identifier)*.
func (r *rdParser) parseProcedureFetchInto() ast.StmtNode {
	r.expect(fetch)
	switch r.tok() {
	case next:
		r.advance()
		r.expect(from)
	case from:
		r.advance()
	}
	name := strings.ToLower(r.expect(identifier).lit)
	r.expect(into)
	vars := []string{strings.ToLower(r.expect(identifier).lit)}
	for r.accept(int(',')) {
		vars = append(vars, strings.ToLower(r.expect(identifier).lit))
	}
	return &ast.ProcedureFetchInto{
		CurName:   name,
		Variables: vars,
	}
}

// parseProcedureIfstmt implements ProcedureIfstmt:
// "IF" ProcedureIf "END" "IF".
func (r *rdParser) parseProcedureIfstmt() ast.StmtNode {
	r.expect(ifKwd)
	body := r.parseProcedureIf()
	r.expect(end)
	r.expect(ifKwd)
	return &ast.ProcedureIfInfo{
		IfBody: body,
	}
}

// parseProcedureIf implements ProcedureIf:
// Expression "THEN" ProcedureProcStmt1s procedurceElseIfs, with
// procedurceElseIfs: empty | "ELSEIF" ProcedureIf
// | "ELSE" ProcedureProcStmt1s.
func (r *rdParser) parseProcedureIf() *ast.ProcedureIfBlock {
	expr := r.parseExpression()
	r.expect(then)
	ifBlock := &ast.ProcedureIfBlock{
		IfExpr:           expr,
		ProcedureIfStmts: r.parseProcedureProcStmt1s(),
	}
	switch r.tok() {
	case elseIfKwd:
		r.advance()
		ifBlock.ProcedureElseStmt = &ast.ProcedureElseIfBlock{
			ProcedureIfStmt: r.parseProcedureIf(),
		}
	case elseKwd:
		r.advance()
		ifBlock.ProcedureElseStmt = &ast.ProcedureElseBlock{
			ProcedureIfStmts: r.parseProcedureProcStmt1s(),
		}
	}
	return ifBlock
}

// parseProcedureCaseStmt implements ProcedureCaseStmt:
// ProcedureSimpleCase | ProcedureSearchedCase, told apart by whether
// WHEN follows CASE directly:
//
//	ProcedureSimpleCase: "CASE" Expression SimpleWhenThenList ElseCaseOpt
//	    "END" "CASE"
//	ProcedureSearchedCase: "CASE" SearchedWhenThenList ElseCaseOpt "END"
//	    "CASE"
func (r *rdParser) parseProcedureCaseStmt() ast.StmtNode {
	r.expect(caseKwd)
	if r.tok() == when {
		// SearchedWhenThenList: SearchWhenThen+
		whens := []*ast.SearchWhenThenStmt{r.parseSearchWhenThen()}
		for r.tok() == when {
			whens = append(whens, r.parseSearchWhenThen())
		}
		caseStmt := &ast.SearchCaseStmt{
			WhenCases: whens,
		}
		// ElseCaseOpt: empty | "ELSE" ProcedureProcStmt1s
		if r.accept(elseKwd) {
			caseStmt.ElseCases = r.parseProcedureProcStmt1s()
		}
		r.expect(end)
		r.expect(caseKwd)
		return caseStmt
	}
	cond := r.parseExpression()
	// SimpleWhenThenList: SimpleWhenThen+
	whens := []*ast.SimpleWhenThenStmt{r.parseSimpleWhenThen()}
	for r.tok() == when {
		whens = append(whens, r.parseSimpleWhenThen())
	}
	caseStmt := &ast.SimpleCaseStmt{
		Condition: cond,
		WhenCases: whens,
	}
	if r.accept(elseKwd) {
		caseStmt.ElseCases = r.parseProcedureProcStmt1s()
	}
	r.expect(end)
	r.expect(caseKwd)
	return caseStmt
}

// parseSimpleWhenThen implements SimpleWhenThen:
// "WHEN" Expression "THEN" ProcedureProcStmt1s.
func (r *rdParser) parseSimpleWhenThen() *ast.SimpleWhenThenStmt {
	r.expect(when)
	expr := r.parseExpression()
	r.expect(then)
	return &ast.SimpleWhenThenStmt{
		Expr:           expr,
		ProcedureStmts: r.parseProcedureProcStmt1s(),
	}
}

// parseSearchWhenThen implements SearchWhenThen:
// "WHEN" Expression "THEN" ProcedureProcStmt1s.
func (r *rdParser) parseSearchWhenThen() *ast.SearchWhenThenStmt {
	r.expect(when)
	expr := r.parseExpression()
	r.expect(then)
	return &ast.SearchWhenThenStmt{
		Expr:           expr,
		ProcedureStmts: r.parseProcedureProcStmt1s(),
	}
}

// parseProcedureUnlabelLoopStmt implements ProcedureUnlabelLoopStmt:
//
//	"WHILE" Expression "DO" ProcedureProcStmt1s "END" "WHILE"
//	| "REPEAT" ProcedureProcStmt1s "UNTIL" Expression "END" "REPEAT"
func (r *rdParser) parseProcedureUnlabelLoopStmt() ast.StmtNode {
	if r.tok() == while {
		r.advance()
		cond := r.parseExpression()
		r.expect(do)
		body := r.parseProcedureProcStmt1s()
		r.expect(end)
		r.expect(while)
		return &ast.ProcedureWhileStmt{
			Condition: cond,
			Body:      body,
		}
	}
	if r.tok() == loop {
		// "LOOP" ProcedureProcStmt1s "END" "LOOP"
		r.advance()
		body := r.parseProcedureProcStmt1s()
		r.expect(end)
		r.expect(loop)
		return &ast.ProcedureLoopStmt{
			Body: body,
		}
	}
	r.expect(repeat)
	body := r.parseProcedureProcStmt1s()
	r.expect(until)
	cond := r.parseExpression()
	r.expect(end)
	r.expect(repeat)
	return &ast.ProcedureRepeatStmt{
		Body:      body,
		Condition: cond,
	}
}

// parseProcedureLabeled implements ProcedureLabeledBlock and
// ProcedurelabeledLoopStmt:
//
//	identifier ':' ProcedureBlockContent ProcedurceLabelOpt
//	identifier ':' ProcedureUnlabelLoopStmt ProcedurceLabelOpt
func (r *rdParser) parseProcedureLabeled() ast.StmtNode {
	label := r.expect(identifier).lit
	r.expect(int(':'))
	switch r.tok() {
	case begin:
		labelBlock := &ast.ProcedureLabelBlock{
			LabelName: label,
			Block:     r.parseProcedureBlockContent(),
		}
		if endLabel := r.parseProcedurceLabelOpt(); endLabel != "" && label != endLabel {
			labelBlock.LabelError = true
			labelBlock.LabelEnd = endLabel
		}
		return labelBlock
	case while, repeat, loop:
		labelLoop := &ast.ProcedureLabelLoop{
			LabelName: label,
			Block:     r.parseProcedureUnlabelLoopStmt(),
		}
		if endLabel := r.parseProcedurceLabelOpt(); endLabel != "" && label != endLabel {
			labelLoop.LabelError = true
			labelLoop.LabelEnd = endLabel
		}
		return labelLoop
	}
	r.syntaxError()
	return nil
}

// parseProcedurceLabelOpt implements ProcedurceLabelOpt:
// empty | identifier.
func (r *rdParser) parseProcedurceLabelOpt() string {
	if r.tok() == identifier {
		return r.expect(identifier).lit
	}
	return ""
}

// parseProcedureProcStmt1s implements ProcedureProcStmt1s:
// (ProcedureProcStmt ';')+. The list ends at the tokens that can follow
// it in its uses (END, ELSEIF, ELSE, UNTIL, WHEN), none of which can
// begin a ProcedureProcStmt.
func (r *rdParser) parseProcedureProcStmt1s() []ast.StmtNode {
	stmts := []ast.StmtNode{r.parseProcedureProcStmt()}
	r.expect(int(';'))
	for {
		switch r.tok() {
		case end, elseIfKwd, elseKwd, until, when, 0:
			return stmts
		}
		stmts = append(stmts, r.parseProcedureProcStmt())
		r.expect(int(';'))
	}
}
