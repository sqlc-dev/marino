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

// DML statements: INSERT/REPLACE, UPDATE, DELETE, LOAD DATA, IMPORT
// INTO, and the BATCH non-transactional wrapper.

import (
	"strings"

	"github.com/sqlc-dev/marino/ast"
	"github.com/sqlc-dev/marino/mysql"
	"github.com/sqlc-dev/marino/opcode"
)

// parseTableOptimizerHintsOpt implements TableOptimizerHintsOpt.
func (r *rdParser) parseTableOptimizerHintsOpt() []*ast.TableOptimizerHint {
	if r.tok() != hintComment {
		return nil
	}
	return r.parseTableOptimizerHints()
}

// parsePriorityOpt implements PriorityOpt.
func (r *rdParser) parsePriorityOpt() mysql.PriorityEnum {
	switch r.tok() {
	case lowPriority:
		r.advance()
		return mysql.LowPriority
	case highPriority:
		r.advance()
		return mysql.HighPriority
	case delayed:
		r.advance()
		return mysql.DelayedPriority
	}
	return mysql.NoPriority
}

// parseInsertIntoStmt implements InsertIntoStmt.
func (r *rdParser) parseInsertIntoStmt() ast.StmtNode {
	// "INSERT" TableOptimizerHintsOpt PriorityOpt IgnoreOptional IntoOpt
	// TableName PartitionNameListOpt InsertValues OnDuplicateKeyUpdate
	r.expect(insert)
	hints := r.parseTableOptimizerHintsOpt()
	priority := r.parsePriorityOpt()
	ignoreErr := r.accept(ignore)
	r.accept(into)
	tn := r.parseTableName()
	partitions := r.parsePartitionNameListOpt()
	x := r.parseInsertValues()
	x.Priority = priority
	x.IgnoreErr = ignoreErr
	ts := &ast.TableSource{Source: tn}
	x.Table = &ast.TableRefsClause{TableRefs: &ast.Join{Left: ts}}
	if r.tok() == on {
		// OnDuplicateKeyUpdate: "ON" "DUPLICATE" "KEY" "UPDATE" AssignmentList
		r.advance()
		r.expect(duplicate)
		r.expect(key)
		r.expect(update)
		x.OnDuplicate = r.parseAssignmentList()
	}
	if hints != nil {
		x.TableHints = hints
	}
	x.PartitionNames = partitions
	return x
}

// parseReplaceIntoStmt implements ReplaceIntoStmt.
func (r *rdParser) parseReplaceIntoStmt() ast.StmtNode {
	// "REPLACE" TableOptimizerHintsOpt PriorityOpt IntoOpt TableName
	// PartitionNameListOpt InsertValues
	r.expect(replace)
	hints := r.parseTableOptimizerHintsOpt()
	priority := r.parsePriorityOpt()
	r.accept(into)
	tn := r.parseTableName()
	partitions := r.parsePartitionNameListOpt()
	x := r.parseInsertValues()
	if hints != nil {
		x.TableHints = hints
	}
	x.IsReplace = true
	x.Priority = priority
	ts := &ast.TableSource{Source: tn}
	x.Table = &ast.TableRefsClause{TableRefs: &ast.Join{Left: ts}}
	x.PartitionNames = partitions
	return x
}

// parseInsertValues implements InsertValues.
func (r *rdParser) parseInsertValues() *ast.InsertStmt {
	switch {
	case r.tok() == int('(') && !r.subSelectFollows():
		// '(' ColumnNameListOpt ')' followed by values or a query.
		r.advance()
		var cols []*ast.ColumnName
		if r.tok() != int(')') {
			cols = r.parseColumnNameList()
		} else {
			cols = []*ast.ColumnName{}
		}
		r.expect(int(')'))
		if r.tok() == value || (r.tok() == values && r.la(1) != row) {
			r.advance()
			is := &ast.InsertStmt{Columns: cols, Lists: r.parseValuesList()}
			is.RowAlias, is.ColumnAliases = r.parseRowAliasOpt()
			return is
		}
		// '(' ColumnNameListOpt ')' SetOprStmt/SelectStmt/SelectStmtWithClause/SubSelect
		return &ast.InsertStmt{Columns: cols, Select: r.parseInsertSelectValue()}
	case r.tok() == value || (r.tok() == values && r.la(1) != row):
		// ValueSym ValuesList RowAliasOpt %prec insertValues
		r.advance()
		is := &ast.InsertStmt{Lists: r.parseValuesList()}
		is.RowAlias, is.ColumnAliases = r.parseRowAliasOpt()
		return is
	case r.tok() == set:
		// "SET" ColumnSetValueList RowAliasOpt
		r.advance()
		is := &ast.InsertStmt{Setlist: true, Lists: [][]ast.ExprNode{{}}}
		for {
			col := r.parseColumnName()
			if !r.accept(eq) && !r.accept(assignmentEq) {
				r.syntaxError()
			}
			expr := r.parseExprOrDefault()
			is.Columns = append(is.Columns, col)
			is.Lists[0] = append(is.Lists[0], expr)
			if !r.accept(int(',')) {
				break
			}
		}
		is.RowAlias, is.ColumnAliases = r.parseRowAliasOpt()
		return is
	default:
		return &ast.InsertStmt{Select: r.parseInsertSelectValue()}
	}
}

// parseInsertSelectValue parses the query-form InsertValues alternatives
// (SetOprStmt | SelectStmt | SelectStmtWithClause | SubSelect).
func (r *rdParser) parseInsertSelectValue() ast.ResultSetNode {
	stmt, sub := r.parseSelectCore()
	if sub != nil {
		switch x := sub.Query.(type) {
		case *ast.SelectStmt:
			x.IsInBraces = true
			return x
		case *ast.SetOprStmt:
			x.IsInBraces = true
			return x
		}
	}
	switch x := stmt.(type) {
	case *ast.SelectStmt:
		return x
	case *ast.SetOprStmt:
		return x
	}
	r.syntaxError()
	return nil
}

// parseValuesList implements ValuesList.
func (r *rdParser) parseValuesList() [][]ast.ExprNode {
	lists := [][]ast.ExprNode{r.parseRowValue()}
	for r.accept(int(',')) {
		lists = append(lists, r.parseRowValue())
	}
	return lists
}

// parseRowAliasOpt implements RowAliasOpt.
func (r *rdParser) parseRowAliasOpt() (*ast.CIStr, []*ast.CIStr) {
	if r.tok() != as {
		return nil, nil
	}
	r.advance()
	name := ast.NewCIStr(r.parseIdentifier())
	if r.tok() != int('(') {
		return &name, nil
	}
	r.advance()
	idents := []ast.CIStr{ast.NewCIStr(r.parseIdentifier())}
	for r.accept(int(',')) {
		idents = append(idents, ast.NewCIStr(r.parseIdentifier()))
	}
	r.expect(int(')'))
	cols := make([]*ast.CIStr, 0, len(idents))
	for i := range idents {
		c := idents[i]
		cols = append(cols, &c)
	}
	return &name, cols
}

// parseAssignmentList implements AssignmentList.
func (r *rdParser) parseAssignmentList() []*ast.Assignment {
	list := []*ast.Assignment{r.parseAssignment()}
	for r.accept(int(',')) {
		list = append(list, r.parseAssignment())
	}
	return list
}

// parseAssignment implements Assignment: ColumnName EqOrAssignmentEq ExprOrDefault.
func (r *rdParser) parseAssignment() *ast.Assignment {
	col := r.parseColumnName()
	if !r.accept(eq) && !r.accept(assignmentEq) {
		r.syntaxError()
	}
	return &ast.Assignment{Column: col, Expr: r.parseExprOrDefault()}
}

// parseUpdateStmt implements UpdateStmt/UpdateStmtNoWith. The WithClause,
// if any, was consumed by the statement dispatcher.
func (r *rdParser) parseUpdateStmt(withClause *ast.WithClause) ast.StmtNode {
	r.expect(update)
	hints := r.parseTableOptimizerHintsOpt()
	priority := r.parsePriorityOpt()
	ignoreErr := r.accept(ignore)
	single := r.parseTableRef()
	multi := false
	refs, isJoin := single.(*ast.Join)
	if !isJoin {
		refs = &ast.Join{Left: single}
	}
	for r.accept(int(',')) {
		// The TableRefs alternative (no ORDER BY/LIMIT tail).
		multi = true
		right := r.parseEscapedTableRef()
		refs = &ast.Join{Left: refs, Right: right, Tp: ast.CrossJoin}
	}
	r.expect(set)
	st := &ast.UpdateStmt{
		Priority:  priority,
		TableRefs: &ast.TableRefsClause{TableRefs: refs},
		List:      r.parseAssignmentList(),
		IgnoreErr: ignoreErr,
	}
	if hints != nil {
		st.TableHints = hints
	}
	if r.accept(where) {
		st.Where = r.parseExpression()
	}
	if !multi {
		if r.tok() == order {
			st.Order = r.parseOrderBy()
		}
		st.Limit = r.parseLimitClause()
	}
	st.With = withClause
	return st
}

// parseDeleteStmt implements DeleteFromStmt (all three table forms).
func (r *rdParser) parseDeleteStmt(withClause *ast.WithClause) ast.StmtNode {
	r.expect(deleteKwd)
	hints := r.parseTableOptimizerHintsOpt()
	priority := r.parsePriorityOpt()
	quickFlag := r.accept(quick)
	ignoreErr := r.accept(ignore)

	x := &ast.DeleteStmt{
		Priority:  priority,
		Quick:     quickFlag,
		IgnoreErr: ignoreErr,
	}
	if hints != nil {
		x.TableHints = hints
	}

	if !r.accept(from) {
		// DeleteWithoutUsingStmt (multi-table, tables before FROM):
		// "DELETE" ... TableAliasRefList "FROM" TableRefs WhereClauseOptional
		x.IsMultiTable = true
		x.BeforeFrom = true
		x.Tables = &ast.DeleteTableList{Tables: r.parseTableAliasRefList()}
		r.expect(from)
		x.TableRefs = &ast.TableRefsClause{TableRefs: r.parseTableRefs()}
		if r.accept(where) {
			x.Where = r.parseExpression()
		}
		x.With = withClause
		return x
	}

	// After FROM: either the single-table form or the USING form. Parse
	// TableNameOptWild first; a wildcard suffix, a comma, or USING commits
	// to the USING form.
	first, hadWild := r.parseTableNameOptWild()
	if hadWild || r.tok() == int(',') || r.tok() == using {
		// DeleteWithUsingStmt:
		// "DELETE" ... "FROM" TableAliasRefList "USING" TableRefs WhereClauseOptional
		tables := []*ast.TableName{first}
		for r.accept(int(',')) {
			tn, _ := r.parseTableNameOptWild()
			tables = append(tables, tn)
		}
		r.expect(using)
		x.IsMultiTable = true
		x.Tables = &ast.DeleteTableList{Tables: tables}
		x.TableRefs = &ast.TableRefsClause{TableRefs: r.parseTableRefs()}
		if r.accept(where) {
			x.Where = r.parseExpression()
		}
		x.With = withClause
		return x
	}

	// DeleteWithoutUsingStmt (single table):
	// "DELETE" ... "FROM" TableName PartitionNameListOpt TableAsNameOpt
	// IndexHintListOpt WhereClauseOptional OrderByOptional LimitClause
	tn := first
	tn.PartitionNames = r.parsePartitionNameListOpt()
	asName := r.parseTableAsNameOpt()
	tn.IndexHints = r.parseIndexHintListOpt()
	join := &ast.Join{Left: &ast.TableSource{Source: tn, AsName: asName}, Right: nil}
	x.TableRefs = &ast.TableRefsClause{TableRefs: join}
	if r.accept(where) {
		x.Where = r.parseExpression()
	}
	if r.tok() == order {
		x.Order = r.parseOrderBy()
	}
	x.Limit = r.parseLimitClause()
	x.With = withClause
	return x
}

// parseTableNameOptWild implements TableNameOptWild, reporting whether
// the '.*' suffix was present.
func (r *rdParser) parseTableNameOptWild() (*ast.TableName, bool) {
	first := r.parseIdentifier()
	tn := &ast.TableName{Name: ast.NewCIStr(first)}
	if r.tok() == int('.') && isIdentifierTok(r.la(1)) {
		r.advance()
		tn.Schema = ast.NewCIStr(first)
		tn.Name = ast.NewCIStr(r.parseIdentifier())
	}
	if r.tok() == int('.') && r.la(1) == int('*') {
		r.advance()
		r.advance()
		return tn, true
	}
	return tn, false
}

// parseTableAliasRefList implements TableAliasRefList.
func (r *rdParser) parseTableAliasRefList() []*ast.TableName {
	tn, _ := r.parseTableNameOptWild()
	tables := []*ast.TableName{tn}
	for r.accept(int(',')) {
		tn, _ := r.parseTableNameOptWild()
		tables = append(tables, tn)
	}
	return tables
}

// parseLoadDataStmt implements LoadDataStmt.
func (r *rdParser) parseLoadDataStmt() ast.StmtNode {
	r.expect(load)
	r.expect(data)
	x := &ast.LoadDataStmt{FileLocRef: ast.FileLocServerOrRemote}
	x.LowPriority = r.accept(lowPriority)
	isLocal := r.accept(local)
	r.expect(infile)
	x.Path = r.expect(stringLit).lit
	if r.accept(format) {
		s := r.expect(stringLit).lit
		x.Format = &s
	}
	x.OnDuplicate = ast.OnDuplicateKeyHandlingError
	switch r.tok() {
	case ignore:
		r.advance()
		x.OnDuplicate = ast.OnDuplicateKeyHandlingIgnore
	case replace:
		r.advance()
		x.OnDuplicate = ast.OnDuplicateKeyHandlingReplace
	}
	r.expect(into)
	r.expect(tableKwd)
	x.Table = r.parseTableName()
	if (r.tok() == character || r.tok() == charType) && r.la(1) == set {
		r.advance()
		r.advance()
		cs := r.parseCharsetName()
		x.Charset = &cs
	}
	x.FieldsInfo = r.parseFieldsClause()
	x.LinesInfo = r.parseLinesClause()
	if r.tok() == ignore {
		// IgnoreLines: "IGNORE" NUM "LINES" — always shifted here.
		r.advance()
		v := getUint64FromNUM(r.expect(intLit).item)
		r.expect(lines)
		x.IgnoreLines = &v
	}
	x.ColumnsAndUserVars = r.parseColumnNameOrUserVarListOptWithBrackets()
	if r.accept(set) {
		// LoadDataSetSpecOpt: "SET" LoadDataSetList
		x.ColumnAssignments = r.parseLoadDataSetList()
	}
	x.Options = r.parseLoadDataOptionListOpt()
	if isLocal {
		x.FileLocRef = ast.FileLocClient
		if x.OnDuplicate == ast.OnDuplicateKeyHandlingError {
			x.OnDuplicate = ast.OnDuplicateKeyHandlingIgnore
		}
	}
	columns := []*ast.ColumnName{}
	for _, v := range x.ColumnsAndUserVars {
		if v.ColumnName != nil {
			columns = append(columns, v.ColumnName)
		}
	}
	x.Columns = columns
	return x
}

// parseColumnNameOrUserVarListOptWithBrackets implements the production
// of the same name.
func (r *rdParser) parseColumnNameOrUserVarListOptWithBrackets() []*ast.ColumnNameOrUserVar {
	if r.tok() != int('(') {
		return []*ast.ColumnNameOrUserVar{}
	}
	r.advance()
	list := []*ast.ColumnNameOrUserVar{}
	for r.tok() != int(')') {
		if len(list) > 0 {
			r.expect(int(','))
		}
		list = append(list, r.parseColumnNameOrUserVariable())
	}
	r.expect(int(')'))
	return list
}

// parseColumnNameOrUserVariable implements ColumnNameOrUserVariable.
func (r *rdParser) parseColumnNameOrUserVariable() *ast.ColumnNameOrUserVar {
	if r.tok() == singleAtIdentifier {
		start := r.cur().offset
		name := strings.TrimPrefix(r.cur().lit, "@")
		r.advance()
		v := r.setOrigin(&ast.VariableExpr{Name: name}, start).(*ast.VariableExpr)
		return &ast.ColumnNameOrUserVar{UserVar: v}
	}
	return &ast.ColumnNameOrUserVar{ColumnName: r.parseColumnName()}
}

// parseLoadDataSetList implements LoadDataSetList.
func (r *rdParser) parseLoadDataSetList() []*ast.Assignment {
	var list []*ast.Assignment
	for {
		// LoadDataSetItem: SimpleIdent "=" ExprOrDefault
		col := r.parseSimpleIdent()
		r.expect(eq)
		list = append(list, &ast.Assignment{Column: col.Name, Expr: r.parseExprOrDefault()})
		if !r.accept(int(',')) {
			return list
		}
	}
}

// parseLoadDataOptionListOpt implements LoadDataOptionListOpt.
func (r *rdParser) parseLoadDataOptionListOpt() []*ast.LoadDataOpt {
	if r.tok() != with {
		return []*ast.LoadDataOpt{}
	}
	r.advance()
	var opts []*ast.LoadDataOpt
	for {
		name := strings.ToLower(r.expect(identifier).lit)
		opt := &ast.LoadDataOpt{Name: name}
		if r.accept(eq) {
			opt.Value = r.parseSignedLiteral()
		}
		opts = append(opts, opt)
		if !r.accept(int(',')) {
			return opts
		}
	}
}

// parseSignedLiteral implements SignedLiteral.
func (r *rdParser) parseSignedLiteral() ast.ExprNode {
	switch r.tok() {
	case int('+'):
		r.advance()
		v := r.parseNumLiteralValue()
		return &ast.UnaryOperationExpr{Op: opcode.Plus, V: ast.NewValueExpr(v, r.p.charset, r.p.collation)}
	case int('-'):
		r.advance()
		v := r.parseNumLiteralValue()
		return &ast.UnaryOperationExpr{Op: opcode.Minus, V: ast.NewValueExpr(v, r.p.charset, r.p.collation)}
	default:
		start := r.cur().offset
		lit := r.parseLiteralExpr(start)
		return ast.NewValueExpr(lit, r.p.charset, r.p.collation)
	}
}

// parseNumLiteralValue consumes NumLiteral and returns its value.
func (r *rdParser) parseNumLiteralValue() interface{} {
	switch r.tok() {
	case intLit, floatLit, decLit:
		v := r.cur().item
		r.advance()
		return v
	}
	r.syntaxError()
	return nil
}

// parseLiteralExpr parses the Literal production (used where a bare
// literal — not a full expression — is required).
func (r *rdParser) parseLiteralExpr(start int) ast.ExprNode {
	switch r.tok() {
	case trueKwd:
		r.advance()
		return r.setOrigin(ast.NewValueExpr(true, r.p.charset, r.p.collation), start)
	case falseKwd:
		r.advance()
		return r.setOrigin(ast.NewValueExpr(false, r.p.charset, r.p.collation), start)
	case null:
		r.advance()
		return r.setOrigin(ast.NewValueExpr(nil, r.p.charset, r.p.collation), start)
	case intLit, decLit, floatLit, hexLit, bitLit:
		v := r.cur().item
		r.advance()
		return r.setOrigin(ast.NewValueExpr(v, r.p.charset, r.p.collation), start)
	case stringLit:
		return r.parseStringLiteral(start)
	case underscoreCS:
		return r.parseUnderscoreCharsetLiteral(start)
	}
	r.syntaxError()
	return nil
}

// parseImportIntoStmt implements ImportIntoStmt.
func (r *rdParser) parseImportIntoStmt() ast.StmtNode {
	r.expect(importKwd)
	r.expect(into)
	st := &ast.ImportIntoStmt{}
	st.Table = r.parseTableName()
	st.ColumnsAndUserVars = r.parseColumnNameOrUserVarListOptWithBrackets()
	if r.accept(set) {
		st.ColumnAssignments = r.parseLoadDataSetList()
	}
	r.expect(from)
	if r.tok() == stringLit {
		st.Path = r.cur().lit
		r.advance()
		if r.accept(format) {
			s := r.expect(stringLit).lit
			st.Format = &s
		}
		st.Options = r.parseLoadDataOptionListOpt()
		return st
	}
	// "FROM" ImportFromSelectStmt
	sel, sub := r.parseSelectCore()
	if sub != nil {
		switch x := sub.Query.(type) {
		case *ast.SelectStmt:
			x.IsInBraces = true
			sel = x
		case *ast.SetOprStmt:
			x.IsInBraces = true
			sel = x
		}
	}
	st.Select = sel.(ast.ResultSetNode)
	for _, cu := range st.ColumnsAndUserVars {
		if cu.ColumnName == nil {
			r.actionError(r.actionErrorf("Cannot use user variable(%s) in IMPORT INTO FROM SELECT statement.", cu.UserVar.Name))
		}
	}
	if st.ColumnAssignments != nil {
		r.actionError(r.actionErrorf("Cannot use SET clause in IMPORT INTO FROM SELECT statement."))
	}
	st.Options = r.parseLoadDataOptionListOpt()
	return st
}

// parseNonTransactionalDMLStmt implements NonTransactionalDMLStmt:
// "BATCH" OptionalShardColumn "LIMIT" NUM DryRunOptions ShardableStmt.
func (r *rdParser) parseNonTransactionalDMLStmt() ast.StmtNode {
	r.expect(batch)
	st := &ast.NonTransactionalDMLStmt{DryRun: ast.NoDryRun}
	if r.accept(on) {
		st.ShardColumn = r.parseColumnName()
	}
	r.expect(limit)
	st.Limit = getUint64FromNUM(r.expect(intLit).item)
	if r.accept(dry) {
		r.expect(run)
		if r.accept(query) {
			st.DryRun = ast.DryRunQuery
		} else {
			st.DryRun = ast.DryRunSplitDml
		}
	}
	var withClause *ast.WithClause
	if r.tok() == with {
		withClause = r.parseWithClause()
	}
	var dml ast.StmtNode
	switch r.tok() {
	case deleteKwd:
		dml = r.parseDeleteStmt(withClause)
	case update:
		dml = r.parseUpdateStmt(withClause)
	case insert:
		if withClause != nil {
			r.syntaxError()
		}
		dml = r.parseInsertIntoStmt()
	case replace:
		if withClause != nil {
			r.syntaxError()
		}
		dml = r.parseReplaceIntoStmt()
	default:
		r.syntaxError()
	}
	st.DMLStmt = dml.(ast.ShardableDMLStmt)
	return st
}
