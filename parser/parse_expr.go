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

// Expressions. The grammar's precedence layering is preserved:
//
//	Expression > BoolPri > PredicateExpr > BitExpr > SimpleExpr
//
// with the %left/%right ladder of parser.y transcribed into the loops
// below, and every LALR conflict resolution replicated as explicit
// lookahead with a comment naming the pseudo-token or rule it replaces.
//
// Mirroring the generated parser's reduction epilogue, every expression
// node records the byte offset of the production's first token via
// SetOriginTextPosition (unless position recording is disabled).

import (
	"strings"

	"github.com/sqlc-dev/marino/ast"
	"github.com/sqlc-dev/marino/charset"
	"github.com/sqlc-dev/marino/internal/panics"
	"github.com/sqlc-dev/marino/mysql"
	"github.com/sqlc-dev/marino/opcode"
	"github.com/sqlc-dev/marino/types"
)

// setOrigin reproduces yySetOffset: after each reduction to an <expr>
// symbol, the generated parser stamps the production's start offset.
func (r *rdParser) setOrigin(n ast.ExprNode, offset int) ast.ExprNode {
	if !r.sc.skipPositionRecording {
		n.SetOriginTextPosition(offset)
	}
	return n
}

// actionError reproduces a `yylex.AppendError(...); return 1` grammar
// action: the error is recorded and the parse aborts.
func (r *rdParser) actionError(err error) {
	r.sc.AppendError(err)
	panic(rdActionAbort{})
}

// rdActionAbort unwinds an aborting grammar action; the recorded errors
// make parseRD fall back to goyacc (which reproduces them) while it is
// in-tree.
type rdActionAbort struct{}

// try runs f speculatively: on rdSyntaxError the position rewinds and
// try returns false. Unimplemented-grammar and action aborts propagate.
func (r *rdParser) try(f func()) (ok bool) {
	m := r.mark()
	defer panics.Handle(func(e interface{}) {
		if _, isSyntax := e.(rdSyntaxError); isSyntax {
			r.rewind(m)
			return
		}
		r.unmark()
		panic(e)
	})
	f()
	r.unmark()
	return true
}

// parseExpression implements Expression: the OR/XOR/AND ladder over
// parseExprUnit operands (rules `Expression logOr Expression %prec pipes`,
// `Expression "XOR" Expression %prec xor`,
// `Expression logAnd Expression %prec andand`).
func (r *rdParser) parseExpression() ast.ExprNode {
	return r.parseExpressionPrec(0)
}

const (
	exprPrecOr  = 1
	exprPrecXor = 2
	exprPrecAnd = 3
)

func exprOpPrec(tok int) (opcode.Op, int) {
	switch tok {
	case or, pipesAsOr:
		return opcode.LogicOr, exprPrecOr
	case xor:
		return opcode.LogicXor, exprPrecXor
	case and, andand:
		return opcode.LogicAnd, exprPrecAnd
	}
	return 0, 0
}

func (r *rdParser) parseExpressionPrec(minPrec int) ast.ExprNode {
	start := r.cur().offset
	v := r.parseExprUnit()
	for {
		op, prec := exprOpPrec(r.tok())
		if prec == 0 || prec < minPrec {
			return v
		}
		r.advance()
		rhs := r.parseExpressionPrec(prec + 1)
		v = r.setOrigin(&ast.BinaryOperationExpr{Op: op, L: v, R: rhs}, start)
	}
}

// parseExprUnit parses one operand of the OR/XOR/AND ladder: the
// assignment, NOT, and MATCH alternatives of Expression, then BoolPri
// with the Expression-level `BoolPri IsOrNotOp TRUE/FALSE/UNKNOWN`
// postfix (which, unlike IS NULL, terminates the BoolPri chain).
func (r *rdParser) parseExprUnit() ast.ExprNode {
	start := r.cur().offset
	switch r.tok() {
	case singleAtIdentifier:
		// Expression: singleAtIdentifier assignmentEq Expression %prec assignmentEq
		if r.la(1) == assignmentEq {
			name := strings.TrimPrefix(r.cur().lit, "@")
			r.advance()
			r.advance()
			v := r.parseExpression()
			return r.setOrigin(&ast.VariableExpr{Name: name, Value: v}, start)
		}
	case not:
		// Expression: "NOT" Expression %prec not
		r.advance()
		v := r.parseExprUnit()
		if x, ok := v.(*ast.ExistsSubqueryExpr); ok {
			x.Not = !x.Not
			return r.setOrigin(x, start)
		}
		return r.setOrigin(&ast.UnaryOperationExpr{Op: opcode.Not, V: v}, start)
	case match:
		// Expression: "MATCH" '(' ColumnNameList ')' "AGAINST" '(' BitExpr FulltextSearchModifierOpt ')'
		r.advance()
		r.expect(int('('))
		cols := r.parseColumnNameList()
		r.expect(int(')'))
		r.expect(against)
		r.expect(int('('))
		expr := r.parseBitExpr(0)
		modifier := r.parseFulltextSearchModifierOpt()
		r.expect(int(')'))
		return r.setOrigin(&ast.MatchAgainst{
			ColumnNames: cols,
			Against:     expr,
			Modifier:    ast.FulltextSearchModifier(modifier),
		}, start)
	}
	v := r.parseBoolPri()
	// Expression: BoolPri IsOrNotOp trueKwd/falseKwd/"UNKNOWN" %prec is.
	// parseBoolPri leaves `IS` unconsumed unless it is IS [NOT] NULL.
	if r.tok() == is {
		notFlag := false
		k := 1
		if r.la(1) == not || r.la(1) == not2 {
			notFlag = true
			k = 2
		}
		switch r.la(k) {
		case trueKwd:
			r.skip(k + 1)
			return r.setOrigin(&ast.IsTruthExpr{Expr: v, Not: notFlag, True: int64(1)}, start)
		case falseKwd:
			r.skip(k + 1)
			return r.setOrigin(&ast.IsTruthExpr{Expr: v, Not: notFlag, True: int64(0)}, start)
		case unknown:
			r.skip(k + 1)
			return r.setOrigin(&ast.IsNullExpr{Expr: v, Not: notFlag}, start)
		}
	}
	return v
}

func (r *rdParser) skip(n int) {
	for range n {
		r.advance()
	}
}

func compareOp(tok int) (opcode.Op, bool) {
	switch tok {
	case ge:
		return opcode.GE, true
	case int('>'):
		return opcode.GT, true
	case le:
		return opcode.LE, true
	case int('<'):
		return opcode.LT, true
	case neq, neqSynonym:
		return opcode.NE, true
	case eq:
		return opcode.EQ, true
	case nulleq:
		return opcode.NullEQ, true
	}
	return 0, false
}

// parseBoolPri implements BoolPri: a left-associative chain of
// comparisons and IS [NOT] NULL over PredicateExpr operands.
func (r *rdParser) parseBoolPri() ast.ExprNode {
	start := r.cur().offset
	v := r.parsePredicate()
	for {
		if r.tok() == is {
			// BoolPri: BoolPri IsOrNotOp "NULL" %prec is
			notFlag := false
			k := 1
			if r.la(1) == not || r.la(1) == not2 {
				notFlag = true
				k = 2
			}
			if r.la(k) == null {
				r.skip(k + 1)
				v = r.setOrigin(&ast.IsNullExpr{Expr: v, Not: notFlag}, start)
				continue
			}
			// IS TRUE/FALSE/UNKNOWN belongs to parseExprUnit.
			return v
		}
		op, ok := compareOp(r.tok())
		if !ok {
			return v
		}
		r.advance()
		switch {
		case (r.tok() == any || r.tok() == some || r.tok() == all) && r.la(1) == int('('):
			// BoolPri: BoolPri CompareOp AnyOrAll SubSelect %prec eq.
			// ANY/SOME can otherwise be identifiers; '(' commits (an
			// unqualified generic call needs a raw identifier token).
			allFlag := r.tok() == all
			r.advance()
			sq := r.parseSubSelect()
			sq.MultiRows = true
			v = r.setOrigin(&ast.CompareSubqueryExpr{Op: op, L: v, R: sq, All: allFlag}, start)
		case r.tok() == singleAtIdentifier && r.la(1) == assignmentEq:
			// BoolPri: BoolPri CompareOp singleAtIdentifier assignmentEq PredicateExpr %prec assignmentEq
			vstart := r.cur().offset
			name := strings.TrimPrefix(r.cur().lit, "@")
			r.advance()
			r.advance()
			val := r.parsePredicate()
			variable := r.setOrigin(&ast.VariableExpr{Name: name, Value: val}, vstart)
			v = r.setOrigin(&ast.BinaryOperationExpr{Op: op, L: v, R: variable}, start)
		default:
			rhs := r.parsePredicate()
			v = r.setOrigin(&ast.BinaryOperationExpr{Op: op, L: v, R: rhs}, start)
		}
	}
}

// parsePredicate implements PredicateExpr: BitExpr optionally followed by
// one [NOT] IN/BETWEEN/LIKE/ILIKE/REGEXP or MEMBER OF predicate.
func (r *rdParser) parsePredicate() ast.ExprNode {
	start := r.cur().offset
	v := r.parseBitExpr(0)
	notFlag := false
	if (r.tok() == not || r.tok() == not2) && predicateFollows(r.la(1)) {
		// NotSym in InOrNotOp/BetweenOrNotOp/LikeOrNotOp/...
		notFlag = true
		r.advance()
	}
	switch r.tok() {
	case in:
		// PredicateExpr: BitExpr InOrNotOp '(' ExpressionList ')'
		//              | BitExpr InOrNotOp SubSelect
		// A '(' directly followed by SELECT/WITH is the subquery form;
		// otherwise the list form is preferred (`in ((select 1))` is a
		// one-element list holding a scalar subquery) with the subquery
		// derivation as the fallback (`in ((select 1) union (select 2))`).
		r.advance()
		if r.tok() == int('(') && (r.la(1) == selectKwd || r.la(1) == with) {
			sq := r.parseSubSelect()
			sq.MultiRows = true
			return r.setOrigin(&ast.PatternInExpr{Expr: v, Not: notFlag, Sel: sq}, start)
		}
		var list []ast.ExprNode
		if r.try(func() {
			r.expect(int('('))
			list = r.parseExpressionList()
			r.expect(int(')'))
		}) {
			return r.setOrigin(&ast.PatternInExpr{Expr: v, Not: notFlag, List: list}, start)
		}
		sq := r.parseSubSelect()
		sq.MultiRows = true
		return r.setOrigin(&ast.PatternInExpr{Expr: v, Not: notFlag, Sel: sq}, start)
	case between:
		// PredicateExpr: BitExpr BetweenOrNotOp BitExpr "AND" PredicateExpr
		r.advance()
		left := r.parseBitExpr(0)
		r.expect(and)
		right := r.parsePredicate()
		return r.setOrigin(&ast.BetweenExpr{Expr: v, Left: left, Right: right, Not: notFlag}, start)
	case like, ilike:
		// PredicateExpr: BitExpr LikeOrNotOp SimpleExpr LikeOrIlikeEscapeOpt
		//              | BitExpr IlikeOrNotOp SimpleExpr LikeOrIlikeEscapeOpt
		isLike := r.tok() == like
		r.advance()
		pattern := r.parseSimpleExpr()
		escapeStr, explicit := "\\", false
		if r.tok() == escape {
			r.advance()
			escapeStr, explicit = r.expect(stringLit).lit, true
		}
		if len(escapeStr) > 1 {
			r.actionError(ErrWrongArguments.GenWithStackByArgs("ESCAPE"))
		}
		var escapeChar byte
		if len(escapeStr) > 0 {
			escapeChar = escapeStr[0]
		}
		return r.setOrigin(&ast.PatternLikeOrIlikeExpr{
			Expr:           v,
			Pattern:        pattern,
			Not:            notFlag,
			Escape:         escapeChar,
			EscapeExplicit: explicit,
			IsLike:         isLike,
		}, start)
	case regexpKwd, rlike:
		// PredicateExpr: BitExpr RegexpOrNotOp SimpleExpr
		r.advance()
		pattern := r.parseSimpleExpr()
		return r.setOrigin(&ast.PatternRegexpExpr{Expr: v, Pattern: pattern, Not: notFlag}, start)
	case memberof:
		// PredicateExpr: BitExpr memberof '(' SimpleExpr ')'
		r.advance()
		r.expect(int('('))
		arg := r.parseSimpleExpr()
		r.expect(int(')'))
		return r.setOrigin(&ast.FuncCallExpr{FnName: ast.NewCIStr(ast.JSONMemberOf), Args: []ast.ExprNode{v, arg}}, start)
	default:
		if notFlag {
			r.syntaxError()
		}
		return v
	}
}

func predicateFollows(tok int) bool {
	switch tok {
	case in, between, like, ilike, regexpKwd, rlike:
		return true
	}
	return false
}

// Binding powers of the BitExpr operators, from the %left ladder:
// '|' < '&' < shifts < '+'/'-' < '*'/'/'/'%'/DIV/MOD < '^'.
func bitOpPrec(tok int) (opcode.Op, int) {
	switch tok {
	case int('|'):
		return opcode.Or, 1
	case int('&'):
		return opcode.And, 2
	case lsh:
		return opcode.LeftShift, 3
	case rsh:
		return opcode.RightShift, 3
	case int('+'):
		return opcode.Plus, 4
	case int('-'):
		return opcode.Minus, 4
	case int('*'):
		return opcode.Mul, 5
	case int('/'):
		return opcode.Div, 5
	case int('%'):
		return opcode.Mod, 5
	case div:
		return opcode.IntDiv, 5
	case mod:
		return opcode.Mod, 5
	case int('^'):
		return opcode.Xor, 6
	}
	return 0, 0
}

// parseBitExpr implements BitExpr as precedence climbing, including the
// three INTERVAL date-arithmetic productions.
func (r *rdParser) parseBitExpr(minPrec int) ast.ExprNode {
	start := r.cur().offset
	var v ast.ExprNode
	if r.tok() == interval {
		// BitExpr: "INTERVAL" Expression TimeUnit '+' BitExpr %prec '+'.
		// INTERVAL( is also FunctionNameConflict; speculate the
		// date-arithmetic form the way the LALR states keep both alive.
		var expr ast.ExprNode
		var unit ast.TimeUnitType
		if r.try(func() {
			r.advance()
			expr = r.parseExpression()
			unit = r.parseTimeUnit()
			r.expect(int('+'))
		}) {
			rhs := r.parseBitExpr(5) // '+' operands bind at the mul level
			v = r.setOrigin(&ast.FuncCallExpr{
				FnName: ast.NewCIStr("DATE_ADD"),
				Args:   []ast.ExprNode{rhs, expr, &ast.TimeUnitExpr{Unit: unit}},
			}, start)
		}
	}
	if v == nil {
		v = r.parseSimpleExpr()
	}
	for {
		op, prec := bitOpPrec(r.tok())
		if prec == 0 || prec < minPrec {
			return v
		}
		opTok := r.tok()
		r.advance()
		if (opTok == int('+') || opTok == int('-')) && r.tok() == interval {
			// BitExpr: BitExpr '+'/'-' "INTERVAL" Expression TimeUnit %prec '+'
			var expr ast.ExprNode
			var unit ast.TimeUnitType
			if r.try(func() {
				r.advance()
				expr = r.parseExpression()
				unit = r.parseTimeUnit()
			}) {
				name := "DATE_ADD"
				if opTok == int('-') {
					name = "DATE_SUB"
				}
				v = r.setOrigin(&ast.FuncCallExpr{
					FnName: ast.NewCIStr(name),
					Args:   []ast.ExprNode{v, expr, &ast.TimeUnitExpr{Unit: unit}},
				}, start)
				continue
			}
		}
		rhs := r.parseBitExpr(prec + 1)
		v = r.setOrigin(&ast.BinaryOperationExpr{Op: op, L: v, R: rhs}, start)
	}
}

// parseSimpleExpr implements the SimpleExpr level: units joined by ||
// under PIPES_AS_CONCAT (rule `SimpleExpr pipes SimpleExpr`; without the
// mode the lexer produces pipesAsOr, handled at the Expression level).
func (r *rdParser) parseSimpleExpr() ast.ExprNode {
	start := r.cur().offset
	v := r.parseSimpleExprUnit()
	for r.tok() == pipes {
		r.advance()
		rhs := r.parseSimpleExprUnit()
		v = r.setOrigin(&ast.FuncCallExpr{FnName: ast.NewCIStr(ast.Concat), Args: []ast.ExprNode{v, rhs}}, start)
	}
	return v
}

// parseSimpleExprUnit parses unary prefixes, an atom, and the COLLATE
// postfix (`SimpleExpr "COLLATE" CollationName`, %right collate — which
// binds tighter than the unary %prec neg rules, hence inside the
// recursion, while pipes binds looser, hence outside).
func (r *rdParser) parseSimpleExprUnit() ast.ExprNode {
	start := r.cur().offset
	var v ast.ExprNode
	switch r.tok() {
	case int('!'):
		// SimpleExpr: '!' SimpleExpr %prec neg
		r.advance()
		v = r.setOrigin(&ast.UnaryOperationExpr{Op: opcode.Not2, V: r.parseSimpleExprUnit()}, start)
	case int('~'):
		r.advance()
		v = r.setOrigin(&ast.UnaryOperationExpr{Op: opcode.BitNeg, V: r.parseSimpleExprUnit()}, start)
	case int('-'):
		r.advance()
		v = r.setOrigin(&ast.UnaryOperationExpr{Op: opcode.Minus, V: r.parseSimpleExprUnit()}, start)
	case int('+'):
		r.advance()
		v = r.setOrigin(&ast.UnaryOperationExpr{Op: opcode.Plus, V: r.parseSimpleExprUnit()}, start)
	case not2:
		// SimpleExpr: not2 SimpleExpr %prec neg (HIGH_NOT_PRECEDENCE)
		r.advance()
		v = r.setOrigin(&ast.UnaryOperationExpr{Op: opcode.Not2, V: r.parseSimpleExprUnit()}, start)
	case binaryType:
		// SimpleExpr: "BINARY" SimpleExpr %prec neg
		r.advance()
		tp := types.NewFieldType(mysql.TypeString)
		tp.SetCharset(charset.CharsetBin)
		tp.SetCollate(charset.CharsetBin)
		tp.AddFlag(mysql.BinaryFlag)
		v = r.setOrigin(&ast.FuncCastExpr{
			Expr:         r.parseSimpleExprUnit(),
			Tp:           tp,
			FunctionType: ast.CastBinaryOperator,
		}, start)
	default:
		v = r.parseSimpleExprAtom()
	}
	for r.tok() == collate {
		r.advance()
		name := r.parseCollationName()
		v = r.setOrigin(&ast.SetCollationExpr{Expr: v, Collate: name}, start)
	}
	return v
}

// parseIdentifier implements Identifier: identifier | UnReservedKeyword |
// NotKeywordToken | TiDBKeyword, returning the source spelling.
func (r *rdParser) parseIdentifier() string {
	if !isIdentifierTok(r.tok()) {
		r.syntaxError()
	}
	lit := r.cur().lit
	r.advance()
	return lit
}

// parseSimpleIdent implements SimpleIdent: a 1-3 part column reference.
func (r *rdParser) parseSimpleIdent() *ast.ColumnNameExpr {
	name := &ast.ColumnName{}
	first := r.parseIdentifier()
	if r.tok() == int('.') && isIdentifierTok(r.la(1)) {
		r.advance()
		second := r.parseIdentifier()
		if r.tok() == int('.') && isIdentifierTok(r.la(1)) {
			r.advance()
			name.Schema = ast.NewCIStr(first)
			name.Table = ast.NewCIStr(second)
			name.Name = ast.NewCIStr(r.parseIdentifier())
		} else {
			name.Table = ast.NewCIStr(first)
			name.Name = ast.NewCIStr(second)
		}
	} else {
		name.Name = ast.NewCIStr(first)
	}
	return &ast.ColumnNameExpr{Name: name}
}

// parseColumnName implements ColumnName (same shape as SimpleIdent).
func (r *rdParser) parseColumnName() *ast.ColumnName {
	return r.parseSimpleIdent().Name
}

// parseColumnNameList implements ColumnNameList.
func (r *rdParser) parseColumnNameList() []*ast.ColumnName {
	cols := []*ast.ColumnName{r.parseColumnName()}
	for r.accept(int(',')) {
		cols = append(cols, r.parseColumnName())
	}
	return cols
}

// parseExpressionList implements ExpressionList.
func (r *rdParser) parseExpressionList() []ast.ExprNode {
	list := []ast.ExprNode{r.parseExpression()}
	for r.accept(int(',')) {
		list = append(list, r.parseExpression())
	}
	return list
}

// parseExpressionListOpt implements ExpressionListOpt (empty before ')').
func (r *rdParser) parseExpressionListOpt() []ast.ExprNode {
	if r.tok() == int(')') {
		return []ast.ExprNode{}
	}
	return r.parseExpressionList()
}

// parseFulltextSearchModifierOpt implements FulltextSearchModifierOpt.
func (r *rdParser) parseFulltextSearchModifierOpt() int {
	switch r.tok() {
	case in:
		if r.la(1) == natural {
			r.expect(in)
			r.expect(natural)
			r.expect(language)
			r.expect(mode)
			if r.tok() == with {
				r.expect(with)
				r.expect(query)
				r.expect(expansion)
				return ast.FulltextSearchModifierNaturalLanguageMode | ast.FulltextSearchModifierWithQueryExpansion
			}
			return ast.FulltextSearchModifierNaturalLanguageMode
		}
		r.expect(in)
		r.expect(booleanType)
		r.expect(mode)
		return ast.FulltextSearchModifierBooleanMode
	case with:
		r.expect(with)
		r.expect(query)
		r.expect(expansion)
		return ast.FulltextSearchModifierWithQueryExpansion
	default:
		return ast.FulltextSearchModifierNaturalLanguageMode
	}
}

// parseTimeUnit implements TimeUnit.
func (r *rdParser) parseTimeUnit() ast.TimeUnitType {
	u, ok := timeUnitFor(r.tok())
	if !ok {
		r.syntaxError()
	}
	r.advance()
	return u
}

func timeUnitFor(tok int) (ast.TimeUnitType, bool) {
	switch tok {
	case microsecond:
		return ast.TimeUnitMicrosecond, true
	case second:
		return ast.TimeUnitSecond, true
	case minute:
		return ast.TimeUnitMinute, true
	case hour:
		return ast.TimeUnitHour, true
	case day:
		return ast.TimeUnitDay, true
	case week:
		return ast.TimeUnitWeek, true
	case month:
		return ast.TimeUnitMonth, true
	case quarter:
		return ast.TimeUnitQuarter, true
	case yearType:
		return ast.TimeUnitYear, true
	case secondMicrosecond:
		return ast.TimeUnitSecondMicrosecond, true
	case minuteMicrosecond:
		return ast.TimeUnitMinuteMicrosecond, true
	case minuteSecond:
		return ast.TimeUnitMinuteSecond, true
	case hourMicrosecond:
		return ast.TimeUnitHourMicrosecond, true
	case hourSecond:
		return ast.TimeUnitHourSecond, true
	case hourMinute:
		return ast.TimeUnitHourMinute, true
	case dayMicrosecond:
		return ast.TimeUnitDayMicrosecond, true
	case daySecond:
		return ast.TimeUnitDaySecond, true
	case dayMinute:
		return ast.TimeUnitDayMinute, true
	case dayHour:
		return ast.TimeUnitDayHour, true
	case yearMonth:
		return ast.TimeUnitYearMonth, true
	}
	return 0, false
}

// parseTimestampUnit implements TimestampUnit (the single-unit subset).
func (r *rdParser) parseTimestampUnit() ast.TimeUnitType {
	switch r.tok() {
	case microsecond, second, minute, hour, day, week, month, quarter, yearType:
		return r.parseTimeUnit()
	}
	r.syntaxError()
	return 0
}

// parseStringName implements StringName: stringLit | Identifier.
func (r *rdParser) parseStringName() string {
	if r.tok() == stringLit {
		lit := r.cur().lit
		r.advance()
		return lit
	}
	return r.parseIdentifier()
}

// parseCharsetName implements CharsetName, validating like MySQL.
func (r *rdParser) parseCharsetName() string {
	if r.tok() == binaryType {
		r.advance()
		return charset.CharsetBin
	}
	name := r.parseStringName()
	cs, err := charset.GetCharsetInfo(name)
	if err != nil {
		r.actionError(ErrUnknownCharacterSet.GenWithStackByArgs(name))
	}
	return cs.Name
}

// parseCollationName implements CollationName.
func (r *rdParser) parseCollationName() string {
	if r.tok() == binaryType {
		r.advance()
		return charset.CollationBin
	}
	name := r.parseStringName()
	info, err := charset.GetCollationByName(name)
	if err != nil {
		r.actionError(err)
	}
	return info.Name
}

// parseTableName implements TableName.
func (r *rdParser) parseTableName() *ast.TableName {
	if r.tok() == int('*') && r.la(1) == int('.') {
		// TableName: '*' '.' Identifier
		r.advance()
		r.advance()
		return &ast.TableName{Schema: ast.NewCIStr("*"), Name: ast.NewCIStr(r.parseIdentifier())}
	}
	first := r.parseIdentifier()
	if r.tok() == int('.') && isIdentifierTok(r.la(1)) {
		r.advance()
		if isInCorrectIdentifierName(first) {
			r.actionError(ErrWrongDBName.GenWithStackByArgs(first))
		}
		return &ast.TableName{Schema: ast.NewCIStr(first), Name: ast.NewCIStr(r.parseIdentifier())}
	}
	return &ast.TableName{Name: ast.NewCIStr(first)}
}

// parseFieldLen implements FieldLen: '(' LengthNum ')'.
func (r *rdParser) parseFieldLen() int {
	r.expect(int('('))
	n := r.parseLengthNum()
	r.expect(int(')'))
	return int(n)
}

// parseLengthNum implements LengthNum.
func (r *rdParser) parseLengthNum() uint64 {
	t := r.expect(intLit)
	return getUint64FromNUM(t.item)
}

// parseOptFieldLen implements OptFieldLen.
func (r *rdParser) parseOptFieldLen() int {
	if r.tok() == int('(') {
		return r.parseFieldLen()
	}
	return types.UnspecifiedLength
}

// subSelectFollows reports whether the current '(' opens a subquery, the
// decision the LALR states make when '(' can start either SubSelect or a
// parenthesized expression: scan past nested '(' to the first meaningful
// token.
func (r *rdParser) subSelectFollows() bool {
	k := 0
	for r.la(k) == int('(') {
		k++
	}
	switch r.la(k) {
	case selectKwd, with:
		return true
	}
	return false
}
