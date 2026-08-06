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

// SELECT statements: SelectStmt, SetOprStmt, SelectStmtWithClause,
// SubSelect, CTEs, table references and joins, and their many optional
// clauses.

import (
	"strings"

	"github.com/sqlc-dev/marino/ast"
	"github.com/sqlc-dev/marino/mysql"
)

// parseSelectCore parses the select family shared by statements and
// subqueries:
//
//	SelectStmt | SetOprStmt | SelectStmtWithClause | SubSelect
//
// It returns either stmt (a *SelectStmt or *SetOprStmt, when the input
// reduced to one of the Stmt forms) or sub (when the input was exactly
// one parenthesized SubSelect with no WITH, set operator, or trailing
// clause — the caller applies its context's action to it).
func (r *rdParser) parseSelectCore() (ast.StmtNode, *ast.SubqueryExpr) {
	var withClause *ast.WithClause
	if r.tok() == with {
		withClause = r.parseWithClause()
	}

	// First SetOprClause.
	var firstSel *ast.SelectStmt
	var firstSub *ast.SubqueryExpr
	if r.tok() == int('(') {
		firstSub = r.parseSubSelect()
	} else {
		firstSel = r.parseSelectStmt()
	}

	if !isSetOprTok(r.tok()) {
		if firstSel != nil {
			// SelectStmt, or SelectStmtWithClause: WithClause SelectStmt.
			if withClause != nil {
				firstSel.With = withClause
			}
			return firstSel, nil
		}
		// A single parenthesized query.
		if r.tok() == order || r.tok() == limit || r.tok() == fetch {
			// SetOprStmtWithLimitOrderBy: SubSelect OrderBy / SubSelect
			// SelectStmtLimit / SubSelect OrderBy SelectStmtLimit.
			var setOprList []ast.Node
			switch x := firstSub.Query.(type) {
			case *ast.SelectStmt:
				setOprList = []ast.Node{&ast.SetOprSelectList{Selects: []ast.Node{x}, With: x.With}}
			case *ast.SetOprStmt:
				setOprList = []ast.Node{&ast.SetOprSelectList{Selects: x.SelectList.Selects, With: x.With, Limit: x.Limit, OrderBy: x.OrderBy}}
			}
			setOpr := &ast.SetOprStmt{SelectList: &ast.SetOprSelectList{Selects: setOprList}}
			if r.tok() == order {
				setOpr.OrderBy = r.parseOrderBy()
			}
			if r.tok() == limit || r.tok() == fetch {
				setOpr.Limit = r.parseSelectStmtLimit()
			}
			// SetOprStmt: WithClause SetOprStmtWithLimitOrderBy.
			setOpr.With = withClause
			return setOpr, nil
		}
		if withClause != nil {
			// SelectStmtWithClause: WithClause SubSelect.
			switch x := firstSub.Query.(type) {
			case *ast.SelectStmt:
				x.IsInBraces = true
				x.WithBeforeBraces = true
				x.With = withClause
				return x, nil
			case *ast.SetOprStmt:
				x.IsInBraces = true
				x.With = withClause
				return x, nil
			}
		}
		return nil, firstSub
	}

	// Set operations: SetOprClauseList (SetOpr SetOprClause)* with the
	// last clause folded per SetOprStmtWoutLimitOrderBy /
	// SetOprStmtWithLimitOrderBy.
	var clauses []ast.Node
	if firstSel != nil {
		// SetOprClause: SelectStmt
		clauses = []ast.Node{firstSel}
	} else {
		// SetOprClause: SubSelect
		clauses = []ast.Node{subSelectToSetOprList(firstSub)}
	}
	for {
		opOffset := r.cur().offset
		op := r.parseSetOpr()
		// The SetOprClauseList/SetOprStmt* actions trim the previous
		// non-braced select's last field text at the set operator.
		if sel, isSelect := clauses[len(clauses)-1].(*ast.SelectStmt); isSelect && !sel.IsInBraces {
			r.p.setLastSelectFieldText(sel, r.endOffsetAt(opOffset))
		}
		if r.tok() == int('(') {
			sub := r.parseSubSelect()
			if isSetOprTok(r.tok()) {
				// Intermediate clause: SetOprClauseList SetOpr SetOprClause.
				lst := subSelectToSetOprList(sub)
				lst.AfterSetOperator = op
				clauses = append(clauses, lst)
				continue
			}
			// Final clause: SetOprStmtWoutLimitOrderBy: ... SetOpr SubSelect
			// (optionally followed by OrderBy / SelectStmtLimit).
			var setOprList2 []ast.Node
			var with2 *ast.WithClause
			var limit2 *ast.Limit
			var orderBy2 *ast.OrderByClause
			switch x := sub.Query.(type) {
			case *ast.SelectStmt:
				setOprList2 = []ast.Node{x}
				with2 = x.With
			case *ast.SetOprStmt:
				setOprList2 = x.SelectList.Selects
				with2 = x.With
				limit2 = x.Limit
				orderBy2 = x.OrderBy
			}
			nextSetOprList := &ast.SetOprSelectList{Selects: setOprList2, With: with2, Limit: limit2, OrderBy: orderBy2}
			nextSetOprList.AfterSetOperator = op
			setOpr := &ast.SetOprStmt{SelectList: &ast.SetOprSelectList{Selects: append(clauses, nextSetOprList)}}
			if r.tok() == order {
				setOpr.OrderBy = r.parseOrderBy()
			}
			if r.tok() == limit || r.tok() == fetch {
				setOpr.Limit = r.parseSelectStmtLimit()
			}
			setOpr.With = withClause
			return setOpr, nil
		}
		sel := r.parseSelectStmt()
		if isSetOprTok(r.tok()) {
			// Intermediate clause: SetOprClause: SelectStmt.
			sel.AfterSetOperator = op
			clauses = append(clauses, sel)
			continue
		}
		// Final clause: SetOprStmtWoutLimitOrderBy: SetOprClauseList SetOpr
		// SelectStmt — the select's own ORDER BY/LIMIT hoist to the set
		// operation.
		setOpr := &ast.SetOprStmt{SelectList: &ast.SetOprSelectList{Selects: clauses}}
		setOpr.Limit = sel.Limit
		setOpr.OrderBy = sel.OrderBy
		sel.Limit = nil
		sel.OrderBy = nil
		sel.AfterSetOperator = op
		setOpr.SelectList.Selects = append(setOpr.SelectList.Selects, sel)
		setOpr.With = withClause
		return setOpr, nil
	}
}

func isSetOprTok(tok int) bool {
	return tok == union || tok == except || tok == intersect
}

// parseSetOpr implements SetOpr: UNION/EXCEPT/INTERSECT with
// DefaultTrueDistinctOpt.
func (r *rdParser) parseSetOpr() *ast.SetOprType {
	var distinctTp, allTp ast.SetOprType
	switch r.tok() {
	case union:
		distinctTp, allTp = ast.Union, ast.UnionAll
	case except:
		distinctTp, allTp = ast.Except, ast.ExceptAll
	case intersect:
		distinctTp, allTp = ast.Intersect, ast.IntersectAll
	default:
		r.syntaxError()
	}
	r.advance()
	tp := distinctTp
	switch r.tok() {
	case all:
		r.advance()
		tp = allTp
	case distinct, distinctRow:
		r.advance()
	}
	return &tp
}

// subSelectToSetOprList implements the SetOprClause: SubSelect action.
func subSelectToSetOprList(sub *ast.SubqueryExpr) *ast.SetOprSelectList {
	switch x := sub.Query.(type) {
	case *ast.SelectStmt:
		return &ast.SetOprSelectList{Selects: []ast.Node{x}}
	case *ast.SetOprStmt:
		return &ast.SetOprSelectList{Selects: x.SelectList.Selects, With: x.With, Limit: x.Limit, OrderBy: x.OrderBy}
	}
	return nil
}

// parseSubSelect implements SubSelect: a parenthesized query with the
// statement-text bookkeeping of its grammar actions.
func (r *rdParser) parseSubSelect() *ast.SubqueryExpr {
	start := r.cur().offset
	r.expect(int('('))
	innerStart := r.cur().offset
	stmt, sub := r.parseSelectCore()
	rparen := r.expect(int(')'))
	if sub != nil {
		// SubSelect: '(' SubSelect ')' — strip the redundant brackets.
		stmt = sub.Query.(ast.StmtNode)
	}
	var result *ast.SubqueryExpr
	switch rs := stmt.(type) {
	case *ast.SelectStmt:
		r.p.setLastSelectFieldText(rs, r.endOffsetAt(rparen.offset))
		r.p.setNodeText(rs, r.src[innerStart:rparen.offset])
		result = &ast.SubqueryExpr{Query: rs}
	case *ast.SetOprStmt:
		r.p.setNodeText(rs, r.src[innerStart:rparen.offset])
		result = &ast.SubqueryExpr{Query: rs}
	default:
		r.syntaxError()
	}
	// SubSelect is an <expr> nonterminal: its reduction stamps the '('.
	r.setOrigin(result, start)
	return result
}

// parseSelectStmt implements the SelectStmt production: an
// unparenthesized SELECT/TABLE/VALUES with its clause tail.
func (r *rdParser) parseSelectStmt() *ast.SelectStmt {
	switch r.tok() {
	case selectKwd:
		st := r.parseSelectStmtBasic()
		switch {
		case r.tok() == from && st.Having == nil && r.la(1) == dual:
			// SelectStmtFromDualTable: SelectStmtBasic FromDual
			// WhereClauseOptional, then SelectStmtGroup OrderByOptional ...
			fromTok := *r.cur()
			r.advance()
			r.expect(dual)
			lastField := st.Fields.Fields[len(st.Fields.Fields)-1]
			if lastField.Expr != nil && lastField.AsName.O == "" {
				r.p.setNodeText(lastField, r.src[lastField.Offset:fromTok.offset-1])
			}
			if r.accept(where) {
				st.Where = r.parseExpression()
			}
			if r.tok() == group {
				st.GroupBy = r.parseGroupByClause()
			}
			r.parseSelectTail(st)
		case r.tok() == from && st.Having == nil:
			// SelectStmtFromTable: SelectStmtBasic "FROM" TableRefsClause
			// WhereClauseOptional SelectStmtGroup HavingClause
			// WindowClauseOptional, then OrderByOptional ...
			fromTok := *r.cur()
			r.advance()
			st.From = &ast.TableRefsClause{TableRefs: r.parseTableRefs()}
			lastField := st.Fields.Fields[len(st.Fields.Fields)-1]
			if lastField.Expr != nil && lastField.AsName.O == "" {
				r.p.setNodeText(lastField, r.src[lastField.Offset:r.endOffsetAt(fromTok.offset)])
			}
			if r.accept(where) {
				st.Where = r.parseExpression()
			}
			if r.tok() == group {
				st.GroupBy = r.parseGroupByClause()
			}
			if r.accept(having) {
				st.Having = &ast.HavingClause{Expr: r.parseExpression()}
			}
			if r.accept(window) {
				st.WindowSpecs = r.parseWindowDefinitionList()
			}
			r.parseSelectTail(st)
		default:
			// SelectStmt: SelectStmtBasic WhereClauseOptional
			// SelectStmtGroup OrderByOptional SelectStmtLimitOpt ...
			if r.accept(where) {
				st.Where = r.parseExpression()
			}
			if r.tok() == group {
				st.GroupBy = r.parseGroupByClause()
			}
			r.parseSelectTail(st)
		}
		return st
	case tableKwd:
		// SelectStmt: "TABLE" TableName OrderByOptional ...
		r.advance()
		st := &ast.SelectStmt{
			Kind:   ast.SelectStmtKindTable,
			Fields: &ast.FieldList{Fields: []*ast.SelectField{{WildCard: &ast.WildCardField{}}}},
		}
		ts := &ast.TableSource{Source: r.parseTableName()}
		st.From = &ast.TableRefsClause{TableRefs: &ast.Join{Left: ts}}
		r.parseSelectTail(st)
		return st
	case values:
		// SelectStmt: "VALUES" ValuesStmtList OrderByOptional ...
		r.advance()
		st := &ast.SelectStmt{
			Kind:   ast.SelectStmtKindValues,
			Fields: &ast.FieldList{Fields: []*ast.SelectField{{WildCard: &ast.WildCardField{}}}},
		}
		st.Lists = []*ast.RowExpr{r.parseRowStmt()}
		for r.accept(int(',')) {
			st.Lists = append(st.Lists, r.parseRowStmt())
		}
		r.parseSelectTail(st)
		return st
	}
	r.syntaxError()
	return nil
}

// parseSelectTail parses OrderByOptional SelectStmtLimitOpt SelectLockOpt
// SelectStmtIntoOption, shared by the SelectStmt alternatives.
func (r *rdParser) parseSelectTail(st *ast.SelectStmt) {
	if r.tok() == order {
		st.OrderBy = r.parseOrderBy()
	}
	if r.tok() == limit || r.tok() == fetch {
		st.Limit = r.parseSelectStmtLimit()
	}
	if lock := r.parseSelectLockOpt(); lock != nil {
		st.LockInfo = lock
	}
	if opt := r.parseSelectStmtIntoOption(); opt != nil {
		st.SelectIntoOpt = opt
	}
}

// parseSelectStmtBasic implements SelectStmtBasic:
// "SELECT" SelectStmtOpts SelectStmtFieldList HavingClause.
func (r *rdParser) parseSelectStmtBasic() *ast.SelectStmt {
	r.expect(selectKwd)
	opts := r.parseSelectStmtOpts()
	fields := r.parseFieldList()
	st := &ast.SelectStmt{
		SelectStmtOpts: opts,
		Distinct:       opts.Distinct,
		Fields:         &ast.FieldList{Fields: fields},
		Kind:           ast.SelectStmtKindSelect,
	}
	if st.SelectStmtOpts.TableHints != nil {
		st.TableHints = st.SelectStmtOpts.TableHints
	}
	if r.accept(having) {
		st.Having = &ast.HavingClause{Expr: r.parseExpression()}
	}
	return st
}

// parseSelectStmtOpts implements SelectStmtOpts/SelectStmtOptsList with
// the merge semantics of the grammar action.
func (r *rdParser) parseSelectStmtOpts() *ast.SelectStmtOpts {
	opts := &ast.SelectStmtOpts{}
	opts.SQLCache = true
	for {
		switch r.tok() {
		case hintComment:
			// TableOptimizerHints; always use the first hint.
			hints := r.parseTableOptimizerHints()
			if hints != nil && opts.TableHints == nil {
				opts.TableHints = hints
			}
		case distinct, distinctRow:
			r.advance()
			opts.Distinct = true
		case all:
			r.advance()
			opts.ExplicitAll = true
		case lowPriority:
			r.advance()
			opts.Priority = mysql.LowPriority
		case highPriority:
			r.advance()
			opts.Priority = mysql.HighPriority
		case delayed:
			r.advance()
			opts.Priority = mysql.DelayedPriority
		case sqlSmallResult:
			r.advance()
			opts.SQLSmallResult = true
		case sqlBigResult:
			r.advance()
			opts.SQLBigResult = true
		case sqlBufferResult:
			r.advance()
			opts.SQLBufferResult = true
		case sqlCache:
			r.advance()
		case sqlNoCache:
			r.advance()
			opts.SQLCache = false
		case sqlCalcFoundRows:
			r.advance()
			opts.CalcFoundRows = true
		case straightJoin:
			r.advance()
			opts.StraightJoin = true
		default:
			if opts.Distinct && opts.ExplicitAll {
				r.actionError(ErrWrongUsage.GenWithStackByArgs("ALL", "DISTINCT"))
			}
			return opts
		}
		if opts.Distinct && opts.ExplicitAll {
			r.actionError(ErrWrongUsage.GenWithStackByArgs("ALL", "DISTINCT"))
		}
	}
}

// parseTableOptimizerHints implements TableOptimizerHints: hintComment.
func (r *rdParser) parseTableOptimizerHints() []*ast.TableOptimizerHint {
	t := r.expect(hintComment)
	hints, warns := r.parseHintComment(t)
	for _, w := range warns {
		r.sc.AppendError(w)
		r.sc.lastErrorAsWarn()
	}
	return hints
}

// parseHintComment routes a hintComment token through the Parser's hint
// parser using the token's recorded hint position.
func (r *rdParser) parseHintComment(t rdToken) ([]*ast.TableOptimizerHint, []error) {
	if r.p.hintParser == nil {
		r.p.hintParser = newHintParser()
	}
	return r.p.hintParser.parse(t.lit, r.sc.GetSQLMode(), t.hintPos)
}

// parseFieldList implements FieldList with the per-field Offset and text
// bookkeeping (text runs to the lookahead token, trimmed).
func (r *rdParser) parseFieldList() []*ast.SelectField {
	var fields []*ast.SelectField
	for {
		start := r.cur().offset
		field := r.parseField()
		field.Offset = start
		if field.Expr != nil {
			endOffset := r.cur().offset
			r.p.setNodeText(field, strings.TrimSpace(r.src[field.Offset:endOffset]))
		}
		fields = append(fields, field)
		if !r.accept(int(',')) {
			return fields
		}
	}
}

// parseField implements Field: wildcards or Expression FieldAsNameOpt.
func (r *rdParser) parseField() *ast.SelectField {
	if r.tok() == int('*') {
		// Field: '*' %prec '*'
		r.advance()
		return &ast.SelectField{WildCard: &ast.WildCardField{}}
	}
	if isIdentifierTok(r.tok()) && r.la(1) == int('.') {
		if r.la(2) == int('*') {
			// Field: Identifier '.' '*' %prec '*'
			table := r.parseIdentifier()
			r.advance()
			r.advance()
			return &ast.SelectField{WildCard: &ast.WildCardField{Table: ast.NewCIStr(table)}}
		}
		if isIdentifierTok(r.la(2)) && r.la(3) == int('.') && r.la(4) == int('*') {
			// Field: Identifier '.' Identifier '.' '*' %prec '*'
			schema := r.parseIdentifier()
			r.advance()
			table := r.parseIdentifier()
			r.advance()
			r.advance()
			return &ast.SelectField{WildCard: &ast.WildCardField{Schema: ast.NewCIStr(schema), Table: ast.NewCIStr(table)}}
		}
	}
	// Field: Expression FieldAsNameOpt
	expr := r.parseExpression()
	asName := r.parseFieldAsNameOpt()
	return &ast.SelectField{Expr: expr, AsName: ast.NewCIStr(asName)}
}

// parseFieldAsNameOpt implements FieldAsNameOpt/FieldAsName.
func (r *rdParser) parseFieldAsNameOpt() string {
	switch {
	case r.tok() == as:
		r.advance()
		if r.tok() == stringLit {
			lit := r.cur().lit
			r.advance()
			return lit
		}
		return r.parseIdentifier()
	case r.tok() == stringLit:
		lit := r.cur().lit
		r.advance()
		return lit
	case isIdentifierTok(r.tok()):
		return r.parseIdentifier()
	}
	return ""
}

// parseGroupByClause implements GroupByClause with WithRollupClause.
func (r *rdParser) parseGroupByClause() *ast.GroupByClause {
	r.expect(group)
	r.expect(by)
	items := r.parseByList()
	rollupFlag := false
	if r.tok() == with && r.la(1) == rollup {
		r.advance()
		r.advance()
		rollupFlag = true
	}
	return &ast.GroupByClause{Items: items, Rollup: rollupFlag}
}

// parseWindowDefinitionList implements WindowDefinitionList.
func (r *rdParser) parseWindowDefinitionList() []ast.WindowSpec {
	specs := []ast.WindowSpec{r.parseWindowDefinition()}
	for r.accept(int(',')) {
		specs = append(specs, r.parseWindowDefinition())
	}
	return specs
}

// parseWindowDefinition implements WindowDefinition: WindowName "AS" WindowSpec.
func (r *rdParser) parseWindowDefinition() ast.WindowSpec {
	name := ast.NewCIStr(r.parseIdentifier())
	r.expect(as)
	spec := r.parseWindowSpec()
	spec.Name = name
	return spec
}

// parseSelectStmtLimit implements SelectStmtLimit.
func (r *rdParser) parseSelectStmtLimit() *ast.Limit {
	if r.accept(fetch) {
		// "FETCH" FirstOrNext FetchFirstOpt RowOrRows "ONLY"
		if r.tok() != first && r.tok() != next {
			r.syntaxError()
		}
		r.advance()
		var count ast.ExprNode
		if r.tok() == row || r.tok() == rows {
			count = ast.NewValueExpr(uint64(1), r.p.charset, r.p.collation)
		} else {
			count = r.parseLimitOption()
		}
		if r.tok() != row && r.tok() != rows {
			r.syntaxError()
		}
		r.advance()
		r.expect(only)
		return &ast.Limit{Count: count}
	}
	r.expect(limit)
	cnt := r.parseLimitOption()
	switch {
	case r.accept(int(',')):
		return &ast.Limit{Offset: cnt, Count: r.parseLimitOption()}
	case r.accept(offset):
		return &ast.Limit{Offset: r.parseLimitOption(), Count: cnt}
	}
	return &ast.Limit{Count: cnt}
}

// parseLimitOption implements LimitOption: LengthNum | paramMarker.
func (r *rdParser) parseLimitOption() ast.ExprNode {
	if r.tok() == paramMarker {
		offset := r.cur().offset
		r.advance()
		return ast.NewParamMarkerExpr(offset)
	}
	return ast.NewValueExpr(r.parseLengthNum(), r.p.charset, r.p.collation)
}

// parseLimitClause implements LimitClause: optional "LIMIT" LimitOption.
func (r *rdParser) parseLimitClause() *ast.Limit {
	if !r.accept(limit) {
		return nil
	}
	return &ast.Limit{Count: r.parseLimitOption().(ast.ValueExpr)}
}

// parseSelectLockOpt implements SelectLockOpt.
func (r *rdParser) parseSelectLockOpt() *ast.SelectLockInfo {
	switch r.tok() {
	case forKwd:
		r.advance()
		var forUpdate bool
		switch r.tok() {
		case update:
			forUpdate = true
		case share:
			forUpdate = false
		default:
			r.syntaxError()
		}
		r.advance()
		tables := []*ast.TableName{}
		if r.accept(of) {
			tables = []*ast.TableName{r.parseTableName()}
			for r.accept(int(',')) {
				tables = append(tables, r.parseTableName())
			}
		}
		info := &ast.SelectLockInfo{Tables: tables}
		switch {
		case r.accept(nowait):
			if forUpdate {
				info.LockType = ast.SelectLockForUpdateNoWait
			} else {
				info.LockType = ast.SelectLockForShareNoWait
			}
		case forUpdate && r.tok() == wait:
			r.advance()
			info.LockType = ast.SelectLockForUpdateWaitN
			info.WaitSec = getUint64FromNUM(r.expect(intLit).item)
		case r.tok() == skip && r.la(1) == locked:
			r.advance()
			r.advance()
			if forUpdate {
				info.LockType = ast.SelectLockForUpdateSkipLocked
			} else {
				info.LockType = ast.SelectLockForShareSkipLocked
			}
		default:
			if forUpdate {
				info.LockType = ast.SelectLockForUpdate
			} else {
				info.LockType = ast.SelectLockForShare
			}
		}
		return info
	case lock:
		// "LOCK" "IN" "SHARE" "MODE"
		r.advance()
		r.expect(in)
		r.expect(share)
		r.expect(mode)
		return &ast.SelectLockInfo{LockType: ast.SelectLockForShare, Tables: []*ast.TableName{}}
	}
	return nil
}

// parseSelectStmtIntoOption implements SelectStmtIntoOption.
func (r *rdParser) parseSelectStmtIntoOption() *ast.SelectIntoOption {
	if r.tok() != into {
		return nil
	}
	r.advance()
	r.expect(outfile)
	x := &ast.SelectIntoOption{Tp: ast.SelectIntoOutfile, FileName: r.expect(stringLit).lit}
	if fields := r.parseFieldsClause(); fields != nil {
		x.FieldsInfo = fields
	}
	if lines := r.parseLinesClause(); lines != nil {
		x.LinesInfo = lines
	}
	return x
}

// parseFieldsClause implements Fields: FieldsOrColumns FieldItemList.
func (r *rdParser) parseFieldsClause() *ast.FieldsClause {
	if r.tok() != fields && r.tok() != columns {
		return nil
	}
	r.advance()
	fieldsClause := &ast.FieldsClause{}
	for {
		item := r.parseFieldItem()
		if item == nil {
			return fieldsClause
		}
		switch item.Type {
		case ast.Terminated:
			fieldsClause.Terminated = &item.Value
		case ast.Enclosed:
			fieldsClause.Enclosed = &item.Value
			fieldsClause.OptEnclosed = item.OptEnclosed
		case ast.Escaped:
			fieldsClause.Escaped = &item.Value
		case ast.DefinedNullBy:
			fieldsClause.DefinedNullBy = &item.Value
			fieldsClause.NullValueOptEnclosed = item.OptEnclosed
		}
	}
}

// parseFieldItem implements FieldItem; nil means the list ended.
func (r *rdParser) parseFieldItem() *ast.FieldItem {
	switch r.tok() {
	case terminated:
		r.advance()
		r.expect(by)
		return &ast.FieldItem{Type: ast.Terminated, Value: r.parseFieldTerminator()}
	case optionallyEnclosedBy:
		r.advance()
		str := r.parseFieldTerminator()
		if str != "\\" && len(str) > 1 {
			r.actionError(ErrWrongFieldTerminators.GenWithStackByArgs())
		}
		return &ast.FieldItem{Type: ast.Enclosed, Value: str, OptEnclosed: true}
	case enclosed:
		r.advance()
		r.expect(by)
		str := r.parseFieldTerminator()
		if str != "\\" && len(str) > 1 {
			r.actionError(ErrWrongFieldTerminators.GenWithStackByArgs())
		}
		return &ast.FieldItem{Type: ast.Enclosed, Value: str}
	case escaped:
		r.advance()
		r.expect(by)
		str := r.parseFieldTerminator()
		if str != "\\" && len(str) > 1 {
			r.actionError(ErrWrongFieldTerminators.GenWithStackByArgs())
		}
		return &ast.FieldItem{Type: ast.Escaped, Value: str}
	case defined:
		r.advance()
		r.expect(null)
		r.expect(by)
		ts := r.parseTextString()
		item := &ast.FieldItem{Type: ast.DefinedNullBy, Value: ts.Value}
		if r.tok() == optionally && r.la(1) == enclosed {
			r.advance()
			r.advance()
			item.OptEnclosed = true
		}
		return item
	}
	return nil
}

// parseFieldTerminator implements FieldTerminator.
func (r *rdParser) parseFieldTerminator() string {
	switch r.tok() {
	case stringLit:
		lit := r.cur().lit
		r.advance()
		return lit
	case hexLit, bitLit:
		v := r.cur().item.(ast.BinaryLiteral).ToString()
		r.advance()
		return v
	}
	r.syntaxError()
	return ""
}

// parseTextString implements TextString.
func (r *rdParser) parseTextString() *ast.TextString {
	switch r.tok() {
	case stringLit:
		lit := r.cur().lit
		r.advance()
		return &ast.TextString{Value: lit}
	case hexLit:
		v := r.cur().item.(ast.BinaryLiteral)
		r.advance()
		return &ast.TextString{Value: v.ToString(), IsBinaryLiteral: true}
	case bitLit:
		v := r.cur().item.(ast.BinaryLiteral)
		r.advance()
		return &ast.TextString{Value: v.ToString(), IsBinaryLiteral: true}
	}
	r.syntaxError()
	return nil
}

// parseLinesClause implements Lines: "LINES" Starting LinesTerminated.
func (r *rdParser) parseLinesClause() *ast.LinesClause {
	if !r.accept(lines) {
		return nil
	}
	clause := &ast.LinesClause{}
	if r.accept(starting) {
		r.expect(by)
		s := r.parseFieldTerminator()
		clause.Starting = &s
	}
	if r.accept(terminated) {
		r.expect(by)
		s := r.parseFieldTerminator()
		clause.Terminated = &s
	}
	return clause
}

// parseRowStmt implements RowStmt: "ROW" RowValue.
func (r *rdParser) parseRowStmt() *ast.RowExpr {
	r.expect(row)
	return &ast.RowExpr{Values: r.parseRowValue()}
}

// parseRowValue implements RowValue: '(' ValuesOpt ')'.
func (r *rdParser) parseRowValue() []ast.ExprNode {
	r.expect(int('('))
	values := []ast.ExprNode{}
	if r.tok() != int(')') {
		values = append(values, r.parseExprOrDefault())
		for r.accept(int(',')) {
			values = append(values, r.parseExprOrDefault())
		}
	}
	r.expect(int(')'))
	return values
}

// parseExprOrDefault implements ExprOrDefault: Expression | "DEFAULT".
func (r *rdParser) parseExprOrDefault() ast.ExprNode {
	if r.tok() == defaultKwd && r.la(1) != int('(') {
		start := r.cur().offset
		r.advance()
		return r.setOrigin(&ast.DefaultExpr{}, start)
	}
	return r.parseExpression()
}

// parseWithClause implements WithClause.
func (r *rdParser) parseWithClause() *ast.WithClause {
	r.expect(with)
	ws := &ast.WithClause{}
	if r.accept(recursive) {
		ws.IsRecursive = true
	}
	ws.CTEs = make([]*ast.CommonTableExpression, 0, 4)
	for {
		cte := &ast.CommonTableExpression{}
		cte.Name = ast.NewCIStr(r.parseIdentifier())
		cte.ColNameList = r.parseIdentListWithParenOpt()
		r.expect(as)
		cte.Query = r.parseSubSelect()
		if ws.IsRecursive {
			cte.IsRecursive = true
		}
		ws.CTEs = append(ws.CTEs, cte)
		if !r.accept(int(',')) {
			return ws
		}
	}
}

// parseIdentListWithParenOpt implements IdentListWithParenOpt.
func (r *rdParser) parseIdentListWithParenOpt() []ast.CIStr {
	if r.tok() != int('(') {
		return []ast.CIStr{}
	}
	r.advance()
	list := []ast.CIStr{ast.NewCIStr(r.parseIdentifier())}
	for r.accept(int(',')) {
		list = append(list, ast.NewCIStr(r.parseIdentifier()))
	}
	r.expect(int(')'))
	return list
}

// parseTableRefs implements TableRefs: comma-separated EscapedTableRefs
// folded into cross joins.
func (r *rdParser) parseTableRefs() *ast.Join {
	first := r.parseEscapedTableRef()
	var res *ast.Join
	if j, ok := first.(*ast.Join); ok {
		res = j
	} else {
		res = &ast.Join{Left: first, Right: nil}
	}
	for r.accept(int(',')) {
		right := r.parseEscapedTableRef()
		res = &ast.Join{Left: res, Right: right, Tp: ast.CrossJoin}
	}
	return res
}

// parseEscapedTableRef implements EscapedTableRef, including the ODBC
// '{' OJ TableRef '}' escape.
func (r *rdParser) parseEscapedTableRef() ast.ResultSetNode {
	if r.tok() == int('{') {
		r.advance()
		r.parseIdentifier() // conventionally OJ; the grammar takes any Identifier
		ref := r.parseTableRef()
		r.expect(int('}'))
		return ref
	}
	return r.parseTableRef()
}

// parseTableRef implements TableRef/JoinTable. The plain CrossOpt form
// parses its right side greedily (the %prec tableRefPriority shift) and
// rebalances with ast.NewCrossJoin; NATURAL and STRAIGHT_JOIN reduce
// left-associatively, so their right side is a TableFactor.
func (r *rdParser) parseTableRef() ast.ResultSetNode {
	v := r.parseTableFactor()
	for {
		switch {
		case r.tok() == join || (r.tok() == cross && r.la(1) == join) || (r.tok() == inner && r.la(1) == join):
			// JoinTable: TableRef CrossOpt TableRef [ON Expression | USING (...)]
			if r.tok() != join {
				r.advance()
			}
			r.advance()
			right := r.parseTableRef()
			switch {
			case r.accept(on):
				v = &ast.Join{Left: v, Right: right, Tp: ast.CrossJoin, On: &ast.OnCondition{Expr: r.parseExpression()}}
			case r.accept(using):
				r.expect(int('('))
				cols := r.parseColumnNameList()
				r.expect(int(')'))
				v = &ast.Join{Left: v, Right: right, Tp: ast.CrossJoin, Using: cols}
			default:
				v = ast.NewCrossJoin(v, right)
			}
		case r.tok() == left || r.tok() == right:
			// JoinTable: TableRef JoinType OuterOpt "JOIN" TableRef ON/USING
			tp := ast.LeftJoin
			if r.tok() == right {
				tp = ast.RightJoin
			}
			r.advance()
			r.accept(outer)
			r.expect(join)
			rhs := r.parseTableRef()
			switch {
			case r.accept(on):
				v = &ast.Join{Left: v, Right: rhs, Tp: tp, On: &ast.OnCondition{Expr: r.parseExpression()}}
			case r.accept(using):
				r.expect(int('('))
				cols := r.parseColumnNameList()
				r.expect(int(')'))
				v = &ast.Join{Left: v, Right: rhs, Tp: tp, Using: cols}
			default:
				r.syntaxError()
			}
		case r.tok() == natural:
			// JoinTable: TableRef "NATURAL" ["LEFT"/"RIGHT" [OUTER]] "JOIN" TableRef
			r.advance()
			switch r.tok() {
			case join:
				r.advance()
				v = &ast.Join{Left: v, Right: r.parseTableFactor(), NaturalJoin: true}
			case left, right:
				tp := ast.LeftJoin
				if r.tok() == right {
					tp = ast.RightJoin
				}
				r.advance()
				r.accept(outer)
				r.expect(join)
				v = &ast.Join{Left: v, Right: r.parseTableFactor(), Tp: tp, NaturalJoin: true}
			default:
				r.syntaxError()
			}
		case r.tok() == straightJoin:
			// JoinTable: TableRef "STRAIGHT_JOIN" TableRef [ON/USING]
			r.advance()
			rhs := r.parseTableFactor()
			switch {
			case r.accept(on):
				v = &ast.Join{Left: v, Right: rhs, StraightJoin: true, On: &ast.OnCondition{Expr: r.parseExpression()}}
			case r.accept(using):
				r.expect(int('('))
				cols := r.parseColumnNameList()
				r.expect(int(')'))
				v = &ast.Join{Left: v, Right: rhs, StraightJoin: true, Using: cols}
			default:
				v = &ast.Join{Left: v, Right: rhs, StraightJoin: true}
			}
		default:
			return v
		}
	}
}

// parseTableFactor implements TableFactor.
func (r *rdParser) parseTableFactor() ast.ResultSetNode {
	switch {
	case r.tok() == int('('):
		if r.subSelectFollows() {
			// TableFactor: SubSelect TableAsNameOpt
			sub := r.parseSubSelect()
			return &ast.TableSource{Source: sub.Query.(ast.ResultSetNode), AsName: r.parseTableAsNameOpt()}
		}
		// TableFactor: '(' TableRefs ')'
		r.advance()
		j := r.parseTableRefs()
		r.expect(int(')'))
		j.ExplicitParens = true
		return j
	case r.tok() == lateral:
		// TableFactor: "LATERAL" SubSelect TableAsName IdentListWithParenOpt
		r.advance()
		sub := r.parseSubSelect()
		ts := &ast.TableSource{Source: sub.Query.(ast.ResultSetNode), AsName: r.parseTableAsName()}
		ts.Lateral = true
		ts.ColumnNames = r.parseIdentListWithParenOpt()
		return ts
	default:
		// TableFactor: TableName PartitionNameListOpt TableAsNameOpt
		// AsOfClauseOpt IndexHintListOpt TableSampleOpt
		tn := r.parseTableName()
		tn.PartitionNames = r.parsePartitionNameListOpt()
		asName := r.parseTableAsNameOpt()
		if r.tok() == asof {
			r.advance()
			r.expect(timestampType)
			tn.AsOf = &ast.AsOfClause{TsExpr: r.parseExpression()}
		}
		tn.IndexHints = r.parseIndexHintListOpt()
		if r.tok() == tableSample {
			tn.TableSample = r.parseTableSample()
		}
		return &ast.TableSource{Source: tn, AsName: asName}
	}
}

// parsePartitionNameListOpt implements PartitionNameListOpt.
func (r *rdParser) parsePartitionNameListOpt() []ast.CIStr {
	if r.tok() != partition || r.la(1) != int('(') {
		return []ast.CIStr{}
	}
	r.advance()
	r.advance()
	list := []ast.CIStr{ast.NewCIStr(r.parseIdentifier())}
	for r.accept(int(',')) {
		list = append(list, ast.NewCIStr(r.parseIdentifier()))
	}
	r.expect(int(')'))
	return list
}

// parseTableAsNameOpt implements TableAsNameOpt.
func (r *rdParser) parseTableAsNameOpt() ast.CIStr {
	if r.tok() == as || isIdentifierTok(r.tok()) {
		return r.parseTableAsName()
	}
	return ast.CIStr{}
}

// parseTableAsName implements TableAsName.
func (r *rdParser) parseTableAsName() ast.CIStr {
	r.accept(as)
	return ast.NewCIStr(r.parseIdentifier())
}

// parseIndexHintListOpt implements IndexHintListOpt.
func (r *rdParser) parseIndexHintListOpt() []*ast.IndexHint {
	hints := []*ast.IndexHint{}
	for {
		var tp ast.IndexHintType
		switch r.tok() {
		case use:
			tp = ast.HintUse
		case ignore:
			tp = ast.HintIgnore
		case force:
			tp = ast.HintForce
		default:
			return hints
		}
		if r.la(1) != key && r.la(1) != index {
			return hints
		}
		r.advance()
		r.advance()
		scope := ast.HintForScan
		if r.tok() == forKwd {
			r.advance()
			switch r.tok() {
			case join:
				r.advance()
				scope = ast.HintForJoin
			case order:
				r.advance()
				r.expect(by)
				scope = ast.HintForOrderBy
			case group:
				r.advance()
				r.expect(by)
				scope = ast.HintForGroupBy
			default:
				r.syntaxError()
			}
		}
		r.expect(int('('))
		var names []ast.CIStr
		for r.tok() != int(')') {
			if len(names) > 0 {
				r.expect(int(','))
			}
			if r.tok() == primary {
				names = append(names, ast.NewCIStr(r.cur().lit))
				r.advance()
			} else {
				names = append(names, ast.NewCIStr(r.parseIdentifier()))
			}
		}
		r.expect(int(')'))
		hints = append(hints, &ast.IndexHint{IndexNames: names, HintType: tp, HintScope: scope})
	}
}

// parseTableSample implements the TABLESAMPLE alternatives of
// TableSampleOpt.
func (r *rdParser) parseTableSample() *ast.TableSample {
	r.expect(tableSample)
	method := ast.SampleMethodTypeNone
	switch r.tok() {
	case system:
		method = ast.SampleMethodTypeSystem
		r.advance()
	case bernoulli:
		method = ast.SampleMethodTypeBernoulli
		r.advance()
	case regions:
		method = ast.SampleMethodTypeTiDBRegion
		r.advance()
	}
	r.expect(int('('))
	ts := &ast.TableSample{SampleMethod: method}
	if r.tok() != int(')') {
		ts.Expr = ast.NewValueExpr(r.parseExpressionRaw(), r.p.charset, r.p.collation)
		switch r.tok() {
		case rows:
			ts.SampleClauseUnit = ast.SampleClauseUnitTypeRow
			r.advance()
		case percent:
			ts.SampleClauseUnit = ast.SampleClauseUnitTypePercent
			r.advance()
		}
	}
	r.expect(int(')'))
	if r.tok() == repeatable {
		r.advance()
		r.expect(int('('))
		seed := r.parseExpressionRaw()
		r.expect(int(')'))
		ts.RepeatableSeed = ast.NewValueExpr(seed, r.p.charset, r.p.collation)
	}
	return ts
}

// parseExpressionRaw parses an Expression and returns it as a plain
// value for the TableSample actions, which wrap the expression node in
// NewValueExpr — the grammar passes the *node* itself through.
func (r *rdParser) parseExpressionRaw() interface{} {
	return r.parseExpression()
}
