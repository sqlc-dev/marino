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

// Column definitions and table constraints: ColumnDef, ColumnOption(List),
// DefaultValueExpr and its NowSym/BuiltinFunction/NextValue helpers,
// ConstraintElem and the index-clause productions it shares with indexes
// (IndexNameAndTypeOpt, IndexPartSpecification, IndexOption, ...), and
// ReferenceDef.

import (
	"github.com/sqlc-dev/marino/ast"
	"github.com/sqlc-dev/marino/mysql"
	"github.com/sqlc-dev/marino/types"
)

// parseColumnDef implements ColumnDef.
func (r *rdParser) parseColumnDef() *ast.ColumnDef {
	name := r.parseColumnName()
	if r.tok() == serial {
		// ColumnDef: ColumnName "SERIAL" ColumnOptionListOpt
		r.advance()
		// TODO: check flen 0
		tp := types.NewFieldType(mysql.TypeLonglong)
		options := []*ast.ColumnOption{{Tp: ast.ColumnOptionNotNull}, {Tp: ast.ColumnOptionAutoIncrement}, {Tp: ast.ColumnOptionUniqKey}}
		options = append(options, r.parseColumnOptionListOpt().Options...)
		tp.AddFlag(mysql.UnsignedFlag)
		colDef := &ast.ColumnDef{Name: name, Tp: tp, Options: options}
		if err := colDef.Validate(); err != nil {
			r.actionError(err)
		}
		return colDef
	}
	// ColumnDef: ColumnName Type ColumnOptionListOpt
	tp := r.parseType()
	options := r.parseColumnOptionListOpt().Options
	colDef := &ast.ColumnDef{Name: name, Tp: tp, Options: options}
	if err := colDef.Validate(); err != nil {
		r.actionError(err)
	}
	return colDef
}

// isColumnOptionStart reports whether the current token can begin a
// ColumnOption.
func (r *rdParser) isColumnOptionStart() bool {
	switch r.tok() {
	case not, not2, null, autoIncrement, primary, key, unique, defaultKwd,
		serial, on, comment, check, constraint, generated, as, references,
		collate, columnFormat, storage, autoRandom, secondaryEngineAttribute,
		srid, visible, invisible, engine_attribute:
		return true
	}
	return false
}

// parseColumnOptionListOpt implements ColumnOptionListOpt and
// ColumnOptionList, including the multiple-COLLATE check.
func (r *rdParser) parseColumnOptionListOpt() ast.ColumnOptionList {
	if !r.isColumnOptionStart() {
		return ast.ColumnOptionList{}
	}
	// ColumnOptionList: ColumnOption
	var columnOptionList ast.ColumnOptionList
	first := r.parseColumnOption()
	if columnOption, ok := first.(*ast.ColumnOption); ok {
		hasCollateOption := false
		if columnOption.Tp == ast.ColumnOptionCollate {
			hasCollateOption = true
		}
		columnOptionList = ast.ColumnOptionList{
			HasCollateOption: hasCollateOption,
			Options:          []*ast.ColumnOption{columnOption},
		}
	} else {
		columnOptionList = first.(ast.ColumnOptionList)
	}
	// ColumnOptionList: ColumnOptionList ColumnOption
	for r.isColumnOptionStart() {
		next := r.parseColumnOption()
		if columnOption, ok := next.(*ast.ColumnOption); ok {
			if columnOption.Tp == ast.ColumnOptionCollate && columnOptionList.HasCollateOption {
				r.actionError(ErrParse.GenWithStackByArgs("Multiple COLLATE clauses", r.actionErrorf("").Error()))
			}
			columnOptionList.Options = append(columnOptionList.Options, columnOption)
		} else {
			if columnOptionList.HasCollateOption && next.(ast.ColumnOptionList).HasCollateOption {
				r.actionError(ErrParse.GenWithStackByArgs("Multiple COLLATE clauses", r.actionErrorf("").Error()))
			}
			columnOptionList.Options = append(columnOptionList.Options, next.(ast.ColumnOptionList).Options...)
		}
	}
	return columnOptionList
}

// parseColumnOption implements ColumnOption. Most alternatives produce a
// *ast.ColumnOption; "SERIAL" "DEFAULT" "VALUE" and the CHECK-with-NOT-NULL
// workaround produce an ast.ColumnOptionList, exactly as the actions do.
func (r *rdParser) parseColumnOption() interface{} {
	switch r.tok() {
	case not, not2:
		// ColumnOption: NotSym "NULL"
		r.advance()
		r.expect(null)
		return &ast.ColumnOption{Tp: ast.ColumnOptionNotNull}
	case null:
		r.advance()
		return &ast.ColumnOption{Tp: ast.ColumnOptionNull}
	case autoIncrement:
		r.advance()
		return &ast.ColumnOption{Tp: ast.ColumnOptionAutoIncrement}
	case primary, key:
		// ColumnOption: PrimaryOpt "KEY" GlobalOrLocalOpt
		//             | PrimaryOpt "KEY" WithClustered GlobalOrLocalOpt
		// KEY is normally a synonym for INDEX. The key attribute PRIMARY KEY
		// can also be specified as just KEY when given in a column definition.
		// See http://dev.mysql.com/doc/refman/5.7/en/create-table.html
		if r.tok() == primary {
			r.advance()
		}
		r.expect(key)
		switch r.tok() {
		case clustered:
			r.advance()
			return &ast.ColumnOption{Tp: ast.ColumnOptionPrimaryKey, PrimaryKeyTp: ast.PrimaryKeyTypeClustered, StrValue: r.parseGlobalOrLocalOpt()}
		case nonclustered:
			r.advance()
			return &ast.ColumnOption{Tp: ast.ColumnOptionPrimaryKey, PrimaryKeyTp: ast.PrimaryKeyTypeNonClustered, StrValue: r.parseGlobalOrLocalOpt()}
		}
		return &ast.ColumnOption{Tp: ast.ColumnOptionPrimaryKey, StrValue: r.parseGlobalOrLocalOpt()}
	case unique:
		// ColumnOption: "UNIQUE" "GLOBAL" | "UNIQUE" "LOCAL"
		//             | "UNIQUE" %prec lowerThanKey | "UNIQUE" "KEY" GlobalOrLocalOpt
		r.advance()
		switch r.tok() {
		case global:
			r.advance()
			return &ast.ColumnOption{Tp: ast.ColumnOptionUniqKey, StrValue: "Global"}
		case local:
			r.advance()
			return &ast.ColumnOption{Tp: ast.ColumnOptionUniqKey}
		case key:
			r.advance()
			return &ast.ColumnOption{Tp: ast.ColumnOptionUniqKey, StrValue: r.parseGlobalOrLocalOpt()}
		}
		return &ast.ColumnOption{Tp: ast.ColumnOptionUniqKey}
	case defaultKwd:
		// ColumnOption: "DEFAULT" DefaultValueExpr
		r.advance()
		return &ast.ColumnOption{Tp: ast.ColumnOptionDefaultValue, Expr: r.parseDefaultValueExpr()}
	case serial:
		// ColumnOption: "SERIAL" "DEFAULT" "VALUE"
		r.advance()
		r.expect(defaultKwd)
		r.expect(value)
		return ast.ColumnOptionList{
			Options: []*ast.ColumnOption{{Tp: ast.ColumnOptionNotNull}, {Tp: ast.ColumnOptionAutoIncrement}, {Tp: ast.ColumnOptionUniqKey}},
		}
	case on:
		// ColumnOption: "ON" "UPDATE" NowSymOptionFraction
		r.advance()
		r.expect(update)
		return &ast.ColumnOption{Tp: ast.ColumnOptionOnUpdate, Expr: r.parseNowSymOptionFraction()}
	case comment:
		r.advance()
		s := r.expect(stringLit).lit
		return &ast.ColumnOption{Tp: ast.ColumnOptionComment, Expr: ast.NewValueExpr(s, "", "")}
	case check, constraint:
		// ColumnOption: ConstraintKeywordOpt "CHECK" '(' Expression ')' EnforcedOrNotOrNotNullOpt
		// See https://dev.mysql.com/doc/refman/5.7/en/create-table.html
		// The CHECK clause is parsed but ignored by all storage engines.
		var constraintName interface{}
		if r.tok() == constraint {
			r.advance()
			if isIdentifierTok(r.tok()) {
				// ConstraintKeywordOpt: "CONSTRAINT" Symbol
				constraintName = r.parseIdentifier()
			}
		}
		r.expect(check)
		r.expect(int('('))
		expr := r.parseExpression()
		r.expect(int(')'))
		enf := r.parseEnforcedOrNotOrNotNullOpt()
		optionCheck := &ast.ColumnOption{
			Tp:       ast.ColumnOptionCheck,
			Expr:     expr,
			Enforced: true,
		}
		// Keep the column type check constraint name.
		if constraintName != nil {
			optionCheck.ConstraintName = constraintName.(string)
		}
		switch enf {
		case 0:
			return ast.ColumnOptionList{
				Options: []*ast.ColumnOption{optionCheck, {Tp: ast.ColumnOptionNotNull}},
			}
		case 1:
			optionCheck.Enforced = true
			return optionCheck
		case 2:
			optionCheck.Enforced = false
			return optionCheck
		}
		return optionCheck
	case generated, as:
		// ColumnOption: GeneratedAlways "AS" '(' Expression ')' VirtualOrStored
		if r.tok() == generated {
			r.advance()
			r.expect(always)
		}
		r.expect(as)
		r.expect(int('('))
		startOffset := r.cur().offset
		expr := r.parseExpression()
		rp := r.expect(int(')'))
		stored := r.parseVirtualOrStored()
		r.p.setNodeText(expr, r.src[startOffset:r.endOffsetAt(rp.offset)])
		return &ast.ColumnOption{
			Tp:     ast.ColumnOptionGenerated,
			Expr:   expr,
			Stored: stored,
		}
	case references:
		// ColumnOption: ReferDef
		return &ast.ColumnOption{
			Tp:    ast.ColumnOptionReference,
			Refer: r.parseReferDef(),
		}
	case collate:
		r.advance()
		return &ast.ColumnOption{Tp: ast.ColumnOptionCollate, StrValue: r.parseCollationName()}
	case columnFormat:
		r.advance()
		return &ast.ColumnOption{Tp: ast.ColumnOptionColumnFormat, StrValue: r.parseColumnFormat()}
	case storage:
		// ColumnOption: "STORAGE" StorageMedia
		r.advance()
		opt := &ast.ColumnOption{Tp: ast.ColumnOptionStorage, StrValue: r.parseStorageMedia()}
		r.appendWarnf("The STORAGE clause is parsed but ignored by all storage engines.")
		return opt
	case autoRandom:
		// ColumnOption: "AUTO_RANDOM" AutoRandomOpt
		r.advance()
		return &ast.ColumnOption{Tp: ast.ColumnOptionAutoRandom, AutoRandOpt: r.parseAutoRandomOpt()}
	case secondaryEngineAttribute:
		// ColumnOption: "SECONDARY_ENGINE_ATTRIBUTE" EqOpt stringLit
		r.advance()
		r.parseEqOpt()
		return &ast.ColumnOption{
			Tp:       ast.ColumnOptionSecondaryEngineAttribute,
			StrValue: r.expect(stringLit).lit,
		}
	case engine_attribute:
		// ColumnOption: "ENGINE_ATTRIBUTE" EqOpt stringLit
		r.advance()
		r.parseEqOpt()
		return &ast.ColumnOption{
			Tp:       ast.ColumnOptionEngineAttribute,
			StrValue: r.expect(stringLit).lit,
		}
	case visible:
		// ColumnOption: "VISIBLE"
		r.advance()
		return &ast.ColumnOption{Tp: ast.ColumnOptionVisible}
	case invisible:
		// ColumnOption: "INVISIBLE"
		r.advance()
		return &ast.ColumnOption{Tp: ast.ColumnOptionInvisible}
	case srid:
		// ColumnOption: "SRID" LengthNum — the spatial column attribute
		// (MySQL 26.7 §13.1.20.10); postdates the goyacc grammar.
		r.advance()
		return &ast.ColumnOption{
			Tp:        ast.ColumnOptionSrid,
			UintValue: r.parseLengthNum(),
		}
	}
	r.syntaxError()
	return nil
}

// parseEnforcedOrNotOrNotNullOpt implements EnforcedOrNotOrNotNullOpt
// (0 = NOT NULL, 1 = enforced, 2 = not enforced), the workaround for the
// 2-token lookahead of { [NOT] NULL | CHECK(...) [NOT] ENFORCED }.
func (r *rdParser) parseEnforcedOrNotOrNotNullOpt() int {
	switch r.tok() {
	case not, not2:
		r.advance()
		switch r.tok() {
		case null:
			r.advance()
			return 0
		case enforced:
			// EnforcedOrNot: NotSym "ENFORCED"
			r.advance()
			return 2
		}
		r.syntaxError()
	case enforced:
		r.advance()
		return 1
	}
	// EnforcedOrNotOpt: empty %prec lowerThanNot
	return 1
}

// parseEnforcedOrNotOpt implements EnforcedOrNotOpt.
func (r *rdParser) parseEnforcedOrNotOpt() bool {
	switch r.tok() {
	case enforced:
		r.advance()
		return true
	case not, not2:
		r.advance()
		r.expect(enforced)
		return false
	}
	return true
}

// parseGlobalOrLocalOpt implements GlobalOrLocalOpt.
func (r *rdParser) parseGlobalOrLocalOpt() string {
	switch r.tok() {
	case local:
		r.advance()
		return ""
	case global:
		r.advance()
		return "Global"
	}
	return ""
}

// parseColumnFormat implements ColumnFormat.
func (r *rdParser) parseColumnFormat() string {
	switch r.tok() {
	case defaultKwd:
		r.advance()
		return "DEFAULT"
	case fixed:
		r.advance()
		return "FIXED"
	case dynamic:
		r.advance()
		return "DYNAMIC"
	}
	r.syntaxError()
	return ""
}

// parseStorageMedia implements StorageMedia ($$ is the token's spelling).
func (r *rdParser) parseStorageMedia() string {
	switch r.tok() {
	case defaultKwd, disk, memory:
		lit := r.cur().lit
		r.advance()
		return lit
	}
	r.syntaxError()
	return ""
}

// parseVirtualOrStored implements VirtualOrStored.
func (r *rdParser) parseVirtualOrStored() bool {
	switch r.tok() {
	case virtual:
		r.advance()
		return false
	case stored:
		r.advance()
		return true
	}
	return false
}

// parseAutoRandomOpt implements AutoRandomOpt.
func (r *rdParser) parseAutoRandomOpt() ast.AutoRandomOption {
	if r.tok() != int('(') {
		return ast.AutoRandomOption{ShardBits: types.UnspecifiedLength, RangeBits: types.UnspecifiedLength}
	}
	r.advance()
	shardBits := int(r.parseLengthNum())
	if r.accept(int(',')) {
		rangeBits := int(r.parseLengthNum())
		r.expect(int(')'))
		return ast.AutoRandomOption{ShardBits: shardBits, RangeBits: rangeBits}
	}
	r.expect(int(')'))
	return ast.AutoRandomOption{ShardBits: shardBits, RangeBits: types.UnspecifiedLength}
}

// parseEqOpt consumes EqOpt.
func (r *rdParser) parseEqOpt() {
	r.accept(eq)
}

// parseDefaultValueExpr implements DefaultValueExpr:
//
//	NowSymOptionFractionParentheses | SignedLiteral
//	| NextValueForSequenceParentheses | BuiltinFunction
//	| '(' Identifier ')' | '(' SignedLiteral ')'
func (r *rdParser) parseDefaultValueExpr() ast.ExprNode {
	start := r.cur().offset
	var v ast.ExprNode
	switch r.tok() {
	case currentTs, localTime, localTs, builtinNow, builtinCurDate, currentDate:
		v = r.parseNowSymOptionFraction()
	case next, nextval:
		v = r.parseNextValueForSequence()
	case identifier:
		// BuiltinFunction: identifier '(' [ExpressionList] ')' — the
		// identifier is always shifted; a missing '(' errors after it.
		v = r.parseBuiltinFunction()
	case replace:
		// BuiltinFunction: "REPLACE" '(' ExpressionList ')'
		v = r.parseBuiltinFunction()
	case int('('):
		v = r.parseDefaultValueExprParen()
	default:
		v = r.parseSignedLiteral()
	}
	return r.setOrigin(v, start)
}

// parseDefaultValueExprParen parses the '('-led DefaultValueExpr
// alternatives: '(' Identifier ')' and '(' SignedLiteral ')' (one level
// only), plus the parenthesized recursions of BuiltinFunction,
// NowSymOptionFractionParentheses, and NextValueForSequenceParentheses.
// When none of those alternatives spans the parentheses, the content
// reparses as a full '(' Expression ')' — the MySQL DEFAULT (expr) form
// with operators.
func (r *rdParser) parseDefaultValueExprParen() ast.ExprNode {
	start := r.cur().offset
	var v ast.ExprNode
	if ok := r.try(func() {
		r.expect(int('('))
		switch r.tok() {
		case currentTs, localTime, localTs, builtinNow, builtinCurDate, currentDate:
			v = r.parseNowSymOptionFraction()
		case next:
			// NEXT can also be an Identifier; "VALUE" commits to
			// NextValueForSequence the way the LALR states do.
			if r.la(1) == value {
				v = r.parseNextValueForSequence()
			} else {
				v = r.parseParenIdentifierExpr()
			}
		case nextval:
			if r.la(1) == int('(') {
				v = r.parseNextValueForSequence()
			} else {
				v = r.parseParenIdentifierExpr()
			}
		case identifier:
			if r.la(1) == int('(') {
				v = r.parseBuiltinFunction()
			} else {
				v = r.parseParenIdentifierExpr()
			}
		case replace:
			v = r.parseBuiltinFunction()
		case int('('):
			v = r.parseDefaultNestedParen()
		default:
			if isIdentifierTok(r.tok()) {
				v = r.parseParenIdentifierExpr()
			} else {
				v = r.parseSignedLiteral()
			}
		}
		r.expect(int(')'))
	}); !ok {
		// '(' Expression ')'
		r.expect(int('('))
		v = r.parseExpression()
		r.expect(int(')'))
	}
	return r.setOrigin(v, start)
}

// parseParenIdentifierExpr builds the '(' Identifier ')' node (the
// parentheses are consumed by the caller).
func (r *rdParser) parseParenIdentifierExpr() ast.ExprNode {
	return &ast.ColumnNameExpr{Name: &ast.ColumnName{
		Name: ast.NewCIStr(r.parseIdentifier()),
	}}
}

// parseDefaultNestedParen parses parentheses nested two or more deep in a
// DefaultValueExpr, where only BuiltinFunction,
// NowSymOptionFractionParentheses, and NextValueForSequenceParentheses
// productions recurse — '(' Identifier ')' and '(' SignedLiteral ')' do not.
func (r *rdParser) parseDefaultNestedParen() ast.ExprNode {
	start := r.cur().offset
	r.expect(int('('))
	var v ast.ExprNode
	switch r.tok() {
	case currentTs, localTime, localTs, builtinNow, builtinCurDate, currentDate:
		v = r.parseNowSymOptionFraction()
	case next, nextval:
		v = r.parseNextValueForSequence()
	case identifier:
		if r.la(1) != int('(') {
			r.syntaxError()
		}
		v = r.parseBuiltinFunction()
	case replace:
		v = r.parseBuiltinFunction()
	case int('('):
		v = r.parseDefaultNestedParen()
	default:
		r.syntaxError()
	}
	r.expect(int(')'))
	return r.setOrigin(v, start)
}

// parseNowSymOptionFraction implements NowSymOptionFraction (NowSym,
// NowSymFunc, and CurdateSym forms).
func (r *rdParser) parseNowSymOptionFraction() ast.ExprNode {
	start := r.cur().offset
	var v ast.ExprNode
	switch r.tok() {
	case currentTs, localTime, localTs:
		// NowSym | NowSymFunc '(' ')' | NowSymFunc '(' NUM ')'
		r.advance()
		if r.tok() == int('(') {
			r.advance()
			if r.tok() == intLit {
				num := r.cur().item
				r.advance()
				r.expect(int(')'))
				v = &ast.FuncCallExpr{FnName: ast.NewCIStr("CURRENT_TIMESTAMP"), Args: []ast.ExprNode{ast.NewValueExpr(num, r.p.charset, r.p.collation)}}
			} else {
				r.expect(int(')'))
				v = &ast.FuncCallExpr{FnName: ast.NewCIStr("CURRENT_TIMESTAMP")}
			}
		} else {
			v = &ast.FuncCallExpr{FnName: ast.NewCIStr("CURRENT_TIMESTAMP")}
		}
	case builtinNow:
		// NowSymFunc: builtinNow — always the call form.
		r.advance()
		r.expect(int('('))
		if r.tok() == intLit {
			num := r.cur().item
			r.advance()
			r.expect(int(')'))
			v = &ast.FuncCallExpr{FnName: ast.NewCIStr("CURRENT_TIMESTAMP"), Args: []ast.ExprNode{ast.NewValueExpr(num, r.p.charset, r.p.collation)}}
		} else {
			r.expect(int(')'))
			v = &ast.FuncCallExpr{FnName: ast.NewCIStr("CURRENT_TIMESTAMP")}
		}
	case builtinCurDate:
		// CurdateSym '(' ')'
		r.advance()
		r.expect(int('('))
		r.expect(int(')'))
		v = &ast.FuncCallExpr{FnName: ast.NewCIStr("CURRENT_DATE")}
	case currentDate:
		// "CURRENT_DATE" | CurdateSym '(' ')'
		r.advance()
		if r.accept(int('(')) {
			r.expect(int(')'))
		}
		v = &ast.FuncCallExpr{FnName: ast.NewCIStr("CURRENT_DATE")}
	default:
		r.syntaxError()
	}
	return r.setOrigin(v, start)
}

// parseNextValueForSequence implements NextValueForSequence.
func (r *rdParser) parseNextValueForSequence() ast.ExprNode {
	start := r.cur().offset
	var tn *ast.TableName
	switch r.tok() {
	case next:
		// "NEXT" "VALUE" forKwd TableName
		r.advance()
		r.expect(value)
		r.expect(forKwd)
		tn = r.parseTableName()
	case nextval:
		// "NEXTVAL" '(' TableName ')'
		r.advance()
		r.expect(int('('))
		tn = r.parseTableName()
		r.expect(int(')'))
	default:
		r.syntaxError()
	}
	objNameExpr := &ast.TableNameExpr{
		Name: tn,
	}
	return r.setOrigin(&ast.FuncCallExpr{
		FnName: ast.NewCIStr(ast.NextVal),
		Args:   []ast.ExprNode{objNameExpr},
	}, start)
}

// parseReferDef implements ReferDef.
func (r *rdParser) parseReferDef() *ast.ReferenceDef {
	// "REFERENCES" TableName IndexPartSpecificationListOpt MatchOpt OnDeleteUpdateOpt
	r.expect(references)
	tn := r.parseTableName()
	// IndexPartSpecificationListOpt
	var parts []*ast.IndexPartSpecification
	if r.tok() == int('(') {
		r.advance()
		parts = r.parseIndexPartSpecificationList()
		r.expect(int(')'))
	}
	match := r.parseMatchOpt()
	onDelete, onUpdate := r.parseOnDeleteUpdateOpt()
	return &ast.ReferenceDef{
		Table:                   tn,
		IndexPartSpecifications: parts,
		OnDelete:                onDelete,
		OnUpdate:                onUpdate,
		Match:                   match,
	}
}

// parseMatchOpt implements MatchOpt and Match.
func (r *rdParser) parseMatchOpt() ast.MatchType {
	if r.tok() != match {
		return ast.MatchNone
	}
	r.advance()
	var m ast.MatchType
	switch r.tok() {
	case full:
		m = ast.MatchFull
	case partial:
		m = ast.MatchPartial
	case simple:
		m = ast.MatchSimple
	default:
		r.syntaxError()
	}
	r.advance()
	r.appendWarnf("The MATCH clause is parsed but ignored by all storage engines.")
	return m
}

// parseOnDeleteUpdateOpt implements OnDeleteUpdateOpt, OnDelete, and
// OnUpdate. An "ON" here always commits to ON DELETE/ON UPDATE ReferOpt
// (the `on` token outranks the %prec lowerThanOn empty alternative).
func (r *rdParser) parseOnDeleteUpdateOpt() (*ast.OnDeleteOpt, *ast.OnUpdateOpt) {
	onDelete := &ast.OnDeleteOpt{}
	onUpdate := &ast.OnUpdateOpt{}
	if r.tok() != on {
		return onDelete, onUpdate
	}
	r.advance()
	if r.accept(deleteKwd) {
		// OnDelete [OnUpdate]
		onDelete = &ast.OnDeleteOpt{ReferOpt: r.parseReferOpt()}
		if r.tok() == on {
			r.advance()
			r.expect(update)
			onUpdate = &ast.OnUpdateOpt{ReferOpt: r.parseReferOpt()}
		}
		return onDelete, onUpdate
	}
	// OnUpdate [OnDelete]
	r.expect(update)
	onUpdate = &ast.OnUpdateOpt{ReferOpt: r.parseReferOpt()}
	if r.tok() == on {
		r.advance()
		r.expect(deleteKwd)
		onDelete = &ast.OnDeleteOpt{ReferOpt: r.parseReferOpt()}
	}
	return onDelete, onUpdate
}

// parseReferOpt implements ReferOpt.
func (r *rdParser) parseReferOpt() ast.ReferOptionType {
	switch r.tok() {
	case restrict:
		r.advance()
		return ast.ReferOptionRestrict
	case cascade:
		r.advance()
		return ast.ReferOptionCascade
	case set:
		r.advance()
		if r.accept(null) {
			return ast.ReferOptionSetNull
		}
		r.expect(defaultKwd)
		r.appendWarnf("The SET DEFAULT clause is parsed but ignored by all storage engines.")
		return ast.ReferOptionSetDefault
	case no:
		r.advance()
		r.expect(action)
		return ast.ReferOptionNoAction
	}
	r.syntaxError()
	return 0
}

// parseConstraint implements Constraint: ConstraintKeywordOpt ConstraintElem.
func (r *rdParser) parseConstraint() *ast.Constraint {
	var name interface{}
	if r.tok() == constraint {
		r.advance()
		if isIdentifierTok(r.tok()) {
			// ConstraintKeywordOpt: "CONSTRAINT" Symbol
			name = r.parseIdentifier()
		}
	}
	cst := r.parseConstraintElem()
	if name != nil {
		cst.Name = name.(string)
		cst.IsEmptyIndex = len(cst.Name) == 0
	}
	return cst
}

// parseConstraintElem implements ConstraintElem.
func (r *rdParser) parseConstraintElem() *ast.Constraint {
	switch r.tok() {
	case primary:
		// "PRIMARY" "KEY" IndexNameAndTypeOpt '(' IndexPartSpecificationList ')' IndexOptionList
		r.advance()
		r.expect(key)
		name, indexType := r.parseIndexNameAndTypeOpt()
		r.expect(int('('))
		keys := r.parseIndexPartSpecificationList()
		r.expect(int(')'))
		option := r.parseIndexOptionList()
		c := &ast.Constraint{
			Tp:           ast.ConstraintPrimaryKey,
			Keys:         keys,
			Name:         name.String,
			IsEmptyIndex: name.Empty,
		}
		if option != nil {
			c.Option = option
		}
		if indexType != nil {
			if c.Option == nil {
				c.Option = &ast.IndexOption{}
			}
			c.Option.Tp = indexType.(ast.IndexType)
		}
		return c
	case fulltext:
		// "FULLTEXT" KeyOrIndexOpt IndexName '(' IndexPartSpecificationList ')' IndexOptionList
		r.advance()
		if r.tok() == key || r.tok() == index {
			r.advance()
		}
		name := r.parseIndexName()
		r.expect(int('('))
		keys := r.parseIndexPartSpecificationList()
		r.expect(int(')'))
		option := r.parseIndexOptionList()
		c := &ast.Constraint{
			Tp:           ast.ConstraintFulltext,
			Keys:         keys,
			Name:         name.String,
			IsEmptyIndex: name.Empty,
		}
		if option != nil {
			c.Option = option
		} else {
			c.Option = &ast.IndexOption{}
		}
		return c
	case spatial:
		// "SPATIAL" KeyOrIndexOpt IndexName '(' IndexPartSpecificationList ')' IndexOptionList
		// — the MySQL 26.7 §15.1.20 table-constraint form, which postdates
		// the goyacc grammar (it only had CREATE SPATIAL INDEX). Mirrors
		// the FULLTEXT alternative.
		r.advance()
		if r.tok() == key || r.tok() == index {
			r.advance()
		}
		name := r.parseIndexName()
		r.expect(int('('))
		keys := r.parseIndexPartSpecificationList()
		r.expect(int(')'))
		option := r.parseIndexOptionList()
		c := &ast.Constraint{
			Tp:           ast.ConstraintSpatial,
			Keys:         keys,
			Name:         name.String,
			IsEmptyIndex: name.Empty,
		}
		if option != nil {
			c.Option = option
		} else {
			c.Option = &ast.IndexOption{}
		}
		return c
	case key, index:
		// KeyOrIndex IfNotExists IndexNameAndTypeOpt '(' IndexPartSpecificationList ')' IndexOptionList
		r.advance()
		ifNotExists := r.parseIfNotExists()
		name, indexType := r.parseIndexNameAndTypeOpt()
		r.expect(int('('))
		keys := r.parseIndexPartSpecificationList()
		r.expect(int(')'))
		option := r.parseIndexOptionList()
		c := &ast.Constraint{
			IfNotExists:  ifNotExists,
			Tp:           ast.ConstraintIndex,
			Keys:         keys,
			Name:         name.String,
			IsEmptyIndex: name.Empty,
		}
		if option != nil {
			c.Option = option
		}
		if indexType != nil {
			if c.Option == nil {
				c.Option = &ast.IndexOption{}
			}
			c.Option.Tp = indexType.(ast.IndexType)
		}
		return c
	case unique:
		// "UNIQUE" KeyOrIndexOpt IndexNameAndTypeOpt '(' IndexPartSpecificationList ')' IndexOptionList
		r.advance()
		if r.tok() == key || r.tok() == index {
			r.advance()
		}
		name, indexType := r.parseIndexNameAndTypeOpt()
		r.expect(int('('))
		keys := r.parseIndexPartSpecificationList()
		r.expect(int(')'))
		option := r.parseIndexOptionList()
		c := &ast.Constraint{
			Tp:           ast.ConstraintUniq,
			Keys:         keys,
			Name:         name.String,
			IsEmptyIndex: name.Empty,
		}
		if option != nil {
			c.Option = option
		}
		if indexType != nil {
			if c.Option == nil {
				c.Option = &ast.IndexOption{}
			}
			c.Option.Tp = indexType.(ast.IndexType)
		}
		return c
	case foreign:
		// "FOREIGN" "KEY" IfNotExists IndexName '(' IndexPartSpecificationList ')' ReferDef
		r.advance()
		r.expect(key)
		ifNotExists := r.parseIfNotExists()
		name := r.parseIndexName()
		r.expect(int('('))
		keys := r.parseIndexPartSpecificationList()
		r.expect(int(')'))
		refer := r.parseReferDef()
		return &ast.Constraint{
			IfNotExists:  ifNotExists,
			Tp:           ast.ConstraintForeignKey,
			Keys:         keys,
			Name:         name.String,
			Refer:        refer,
			IsEmptyIndex: name.Empty,
		}
	case check:
		// "CHECK" '(' Expression ')' EnforcedOrNotOpt
		r.advance()
		r.expect(int('('))
		expr := r.parseExpression()
		r.expect(int(')'))
		return &ast.Constraint{
			Tp:       ast.ConstraintCheck,
			Expr:     expr,
			Enforced: r.parseEnforcedOrNotOpt(),
		}
	}
	r.syntaxError()
	return nil
}

// parseConstraintVectorOrColumnarIndex implements ConstraintVectorIndex
// and ConstraintColumnarIndex (which share their shape).
func (r *rdParser) parseConstraintVectorOrColumnarIndex() *ast.Constraint {
	var tp ast.ConstraintType
	switch r.tok() {
	case vectorType:
		// "VECTOR" "INDEX" IfNotExists IndexNameAndTypeOpt '(' IndexPartSpecificationList ')' IndexOptionList
		tp = ast.ConstraintVector
	case columnar:
		// "COLUMNAR" "INDEX" IfNotExists IndexNameAndTypeOpt '(' IndexPartSpecificationList ')' IndexOptionList
		tp = ast.ConstraintColumnar
	default:
		r.syntaxError()
	}
	r.advance()
	r.expect(index)
	ifNotExists := r.parseIfNotExists()
	name, indexType := r.parseIndexNameAndTypeOpt()
	r.expect(int('('))
	keys := r.parseIndexPartSpecificationList()
	r.expect(int(')'))
	option := r.parseIndexOptionList()
	c := &ast.Constraint{
		IfNotExists:  ifNotExists,
		Tp:           tp,
		Keys:         keys,
		Name:         name.String,
		IsEmptyIndex: name.Empty,
	}
	if option != nil {
		c.Option = option
	} else {
		c.Option = &ast.IndexOption{}
	}
	if indexType != nil {
		c.Option.Tp = indexType.(ast.IndexType)
	}
	return c
}

// parseIndexName implements IndexName.
func (r *rdParser) parseIndexName() *ast.NullString {
	if !isIdentifierTok(r.tok()) {
		return &ast.NullString{
			String: "",
			Empty:  false,
		}
	}
	id := r.parseIdentifier()
	return &ast.NullString{
		String: id,
		Empty:  len(id) == 0,
	}
}

// parseIndexNameAndTypeOpt implements IndexNameAndTypeOpt. The TYPE
// spelling of the index-type clause is accepted only after an explicit
// index name (see the comment in parser.y: `INDEX TYPE <index_type>` is
// ambiguous, so TYPE requires a name).
func (r *rdParser) parseIndexNameAndTypeOpt() (*ast.NullString, interface{}) {
	if r.tok() == using {
		// IndexName(empty) "USING" IndexTypeName
		r.advance()
		return &ast.NullString{String: "", Empty: false}, r.parseIndexTypeName()
	}
	if !isIdentifierTok(r.tok()) {
		return &ast.NullString{String: "", Empty: false}, nil
	}
	id := r.parseIdentifier()
	name := &ast.NullString{String: id, Empty: len(id) == 0}
	switch r.tok() {
	case using:
		// IndexName "USING" IndexTypeName
		r.advance()
		return name, r.parseIndexTypeName()
	case tp:
		// Identifier "TYPE" IndexTypeName
		r.advance()
		return name, r.parseIndexTypeName()
	}
	return name, nil
}

// parseIndexTypeName implements IndexTypeName.
func (r *rdParser) parseIndexTypeName() ast.IndexType {
	var t ast.IndexType
	switch r.tok() {
	case btree:
		t = ast.IndexTypeBtree
	case hash:
		t = ast.IndexTypeHash
	case rtree:
		t = ast.IndexTypeRtree
	case hypo:
		t = ast.IndexTypeHypo
	case hnsw:
		t = ast.IndexTypeHNSW
	case inverted:
		t = ast.IndexTypeInverted
	default:
		r.syntaxError()
	}
	r.advance()
	return t
}

// parseIndexPartSpecificationList implements IndexPartSpecificationList.
func (r *rdParser) parseIndexPartSpecificationList() []*ast.IndexPartSpecification {
	list := []*ast.IndexPartSpecification{r.parseIndexPartSpecification()}
	for r.accept(int(',')) {
		list = append(list, r.parseIndexPartSpecification())
	}
	return list
}

// parseIndexPartSpecification implements IndexPartSpecification.
func (r *rdParser) parseIndexPartSpecification() *ast.IndexPartSpecification {
	if r.tok() == int('(') {
		// '(' Expression ')' OptOrder
		r.advance()
		expr := r.parseExpression()
		r.expect(int(')'))
		return &ast.IndexPartSpecification{Expr: expr, Desc: r.parseOptOrder()}
	}
	// ColumnName OptFieldLen OptOrder
	col := r.parseColumnName()
	length := r.parseOptFieldLen()
	return &ast.IndexPartSpecification{Column: col, Length: length, Desc: r.parseOptOrder()}
}

// parseOptOrder implements OptOrder.
func (r *rdParser) parseOptOrder() bool {
	switch r.tok() {
	case asc:
		r.advance()
		return false
	case desc:
		r.advance()
		return true
	}
	return false // ASC by default
}

// isIndexOptionStart reports whether the current token can begin an
// IndexOption.
func (r *rdParser) isIndexOptionStart() bool {
	switch r.tok() {
	case keyBlockSize, addColumnarReplicaOnDemand, using, tp, with, comment,
		visible, invisible, clustered, nonclustered, global, local,
		preSplitRegions, engine_attribute, secondaryEngineAttribute, where:
		return true
	}
	return false
}

// parseIndexOptionList implements IndexOptionList, including its
// field-by-field merge of successive options.
func (r *rdParser) parseIndexOptionList() *ast.IndexOption {
	var opt1 *ast.IndexOption
	for r.isIndexOptionStart() {
		opt2 := r.parseIndexOption()
		// Merge the options
		if opt1 == nil {
			opt1 = opt2
			continue
		}
		if len(opt2.Comment) > 0 {
			opt1.Comment = opt2.Comment
		} else if opt2.Tp != 0 {
			opt1.Tp = opt2.Tp
		} else if opt2.KeyBlockSize > 0 {
			opt1.KeyBlockSize = opt2.KeyBlockSize
		} else if len(opt2.ParserName.O) > 0 {
			opt1.ParserName = opt2.ParserName
		} else if opt2.Visibility != ast.IndexVisibilityDefault {
			opt1.Visibility = opt2.Visibility
		} else if opt2.PrimaryKeyTp != ast.PrimaryKeyTypeDefault {
			opt1.PrimaryKeyTp = opt2.PrimaryKeyTp
		} else if opt2.AddColumnarReplicaOnDemand > 0 {
			opt1.AddColumnarReplicaOnDemand = opt2.AddColumnarReplicaOnDemand
		} else if opt2.Global {
			opt1.Global = true
		} else if opt2.SplitOpt != nil {
			opt1.SplitOpt = opt2.SplitOpt
		} else if len(opt2.EngineAttr) > 0 {
			opt1.EngineAttr = opt2.EngineAttr
		} else if len(opt2.SecondaryEngineAttr) > 0 {
			opt1.SecondaryEngineAttr = opt2.SecondaryEngineAttr
		} else if opt2.Condition != nil {
			opt1.Condition = opt2.Condition
		}
	}
	return opt1
}

// parseIndexOption implements IndexOption.
func (r *rdParser) parseIndexOption() *ast.IndexOption {
	switch r.tok() {
	case keyBlockSize:
		// "KEY_BLOCK_SIZE" EqOpt LengthNum
		r.advance()
		r.parseEqOpt()
		return &ast.IndexOption{
			KeyBlockSize: r.parseLengthNum(),
		}
	case addColumnarReplicaOnDemand:
		r.advance()
		return &ast.IndexOption{
			AddColumnarReplicaOnDemand: 1,
		}
	case using, tp:
		// IndexType: "USING" IndexTypeName | "TYPE" IndexTypeName
		r.advance()
		return &ast.IndexOption{
			Tp: r.parseIndexTypeName(),
		}
	case with:
		// "WITH" "PARSER" Identifier
		r.advance()
		r.expect(parser)
		return &ast.IndexOption{
			ParserName: ast.NewCIStr(r.parseIdentifier()),
		}
	case comment:
		r.advance()
		return &ast.IndexOption{
			Comment: r.expect(stringLit).lit,
		}
	case visible:
		// IndexInvisible: "VISIBLE"
		r.advance()
		return &ast.IndexOption{
			Visibility: ast.IndexVisibilityVisible,
		}
	case invisible:
		r.advance()
		return &ast.IndexOption{
			Visibility: ast.IndexVisibilityInvisible,
		}
	case clustered:
		// WithClustered: "CLUSTERED"
		r.advance()
		return &ast.IndexOption{
			PrimaryKeyTp: ast.PrimaryKeyTypeClustered,
		}
	case nonclustered:
		r.advance()
		return &ast.IndexOption{
			PrimaryKeyTp: ast.PrimaryKeyTypeNonClustered,
		}
	case global:
		r.advance()
		return &ast.IndexOption{
			Global: true,
		}
	case local:
		r.advance()
		return &ast.IndexOption{
			Global: false,
		}
	case preSplitRegions:
		// "PRE_SPLIT_REGIONS" EqOpt ('(' SplitOption ')' | Int64Num)
		r.advance()
		r.parseEqOpt()
		if r.tok() == int('(') {
			r.advance()
			so := r.parseSplitOption()
			r.expect(int(')'))
			return &ast.IndexOption{
				SplitOpt: so,
			}
		}
		return &ast.IndexOption{
			SplitOpt: &ast.SplitOption{
				Num: r.parseInt64Num(),
			},
		}
	case engine_attribute:
		// "ENGINE_ATTRIBUTE" EqOpt stringLit
		r.advance()
		r.parseEqOpt()
		return &ast.IndexOption{EngineAttr: r.expect(stringLit).lit}
	case secondaryEngineAttribute:
		// "SECONDARY_ENGINE_ATTRIBUTE" EqOpt stringLit
		r.advance()
		r.parseEqOpt()
		return &ast.IndexOption{SecondaryEngineAttr: r.expect(stringLit).lit}
	case where:
		// "WHERE" Expression
		r.advance()
		return &ast.IndexOption{
			Condition: r.parseExpression(),
		}
	}
	r.syntaxError()
	return nil
}

// parseSplitOption implements SplitOption.
func (r *rdParser) parseSplitOption() *ast.SplitOption {
	if r.accept(by) {
		// "BY" ValuesList
		return &ast.SplitOption{
			ValueLists: r.parseValuesList(),
		}
	}
	return r.parseSplitOptionBetween()
}

// parseSplitOptionBetween implements SplitOptionBetween.
func (r *rdParser) parseSplitOptionBetween() *ast.SplitOption {
	// "BETWEEN" RowValue "AND" RowValue "REGIONS" Int64Num
	r.expect(between)
	lower := r.parseRowValue()
	r.expect(and)
	upper := r.parseRowValue()
	r.expect(regions)
	return &ast.SplitOption{
		Lower: lower,
		Upper: upper,
		Num:   r.parseInt64Num(),
	}
}
