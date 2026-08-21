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

// SimpleExpr atoms: literals, identifiers, variables, and the
// FunctionCallKeyword / FunctionCallNonKeyword / FunctionCallGeneric /
// SumExpr / WindowFuncCall productions.

import (
	"strings"
	"unicode"

	"github.com/sqlc-dev/marino/ast"
	"github.com/sqlc-dev/marino/charset"
	"github.com/sqlc-dev/marino/mysql"
	"github.com/sqlc-dev/marino/opcode"
	"github.com/sqlc-dev/marino/types"
)

// functionNameConflictTok is the FunctionNameConflict token set: keywords
// that act as generic function names when followed by '('.
var functionNameConflictTok = map[int]struct{}{
	ascii: {}, charsetKwd: {}, coalesce: {}, collation: {}, dateType: {},
	database: {}, day: {}, hour: {}, ifKwd: {}, interval: {}, log: {},
	format: {}, left: {}, microsecond: {}, minute: {}, month: {},
	builtinNow: {}, point: {}, quarter: {}, repeat: {}, replace: {},
	reverse: {}, right: {}, rowCount: {}, second: {}, timeType: {},
	timestampType: {}, truncate: {}, user: {}, week: {}, yearType: {},
}

// parseSimpleExprAtom parses the atom alternatives of SimpleExpr.
func (r *rdParser) parseSimpleExprAtom() ast.ExprNode {
	start := r.cur().offset
	tok := r.tok()
	switch tok {
	case identifier:
		// FunctionCallGeneric: identifier '(' ExpressionListOpt ')'.
		// The unqualified form requires a raw identifier token.
		if r.la(1) == int('(') {
			name := r.cur().lit
			if strings.EqualFold(name, "JSON_VALUE") {
				// JSON_VALUE with a RETURNING or ON EMPTY/ON ERROR
				// clause parses as a JSONValueExpr; the plain call falls
				// through to the generic form.
				var jv ast.ExprNode
				if r.try(func() { jv = r.parseJSONValueExpr() }) {
					return r.setOrigin(jv, start)
				}
			}
			r.advance()
			r.advance()
			args := r.parseExpressionListOpt()
			r.expect(int(')'))
			return r.setOrigin(&ast.FuncCallExpr{FnName: ast.NewCIStr(name), Args: args}, start)
		}
		return r.parseSimpleIdentAtom(start)

	case stringLit:
		return r.parseStringLiteral(start)
	case intLit, decLit, floatLit, hexLit, bitLit:
		v := r.cur().item
		r.advance()
		return r.setOrigin(ast.NewValueExpr(v, r.p.charset, r.p.collation), start)
	case trueKwd:
		r.advance()
		return r.setOrigin(ast.NewValueExpr(true, r.p.charset, r.p.collation), start)
	case falseKwd:
		r.advance()
		return r.setOrigin(ast.NewValueExpr(false, r.p.charset, r.p.collation), start)
	case null:
		r.advance()
		return r.setOrigin(ast.NewValueExpr(nil, r.p.charset, r.p.collation), start)
	case underscoreCS:
		return r.parseUnderscoreCharsetLiteral(start)
	case paramMarker:
		r.advance()
		return r.setOrigin(ast.NewParamMarkerExpr(start), start)

	case singleAtIdentifier:
		// UserVariable: singleAtIdentifier
		name := strings.TrimPrefix(r.cur().lit, "@")
		r.advance()
		return r.setOrigin(&ast.VariableExpr{Name: name}, start)
	case doubleAtIdentifier:
		return r.parseSystemVariable(start)

	case int('('):
		// SimpleExpr: SubSelect | '(' Expression ')' | '(' ExpressionList ',' Expression ')'.
		// A '(' directly followed by SELECT/WITH is a subquery; otherwise
		// the parenthesized-expression route is preferred (`((select 1))`
		// is ParenthesesExpr around a scalar subquery), and only when
		// that route fails — `((select 1) union (select 2))` — is the
		// SubSelect derivation taken, mirroring the LALR resolution.
		if querySubSelectHead(r.la(1), r.la(2)) {
			return r.setOrigin(r.parseSubSelect(), start)
		}
		var result ast.ExprNode
		if r.try(func() { result = r.parseParenExpr(start) }) {
			return result
		}
		return r.setOrigin(r.parseSubSelect(), start)
	case jsonDualityObject:
		// SimpleExpr: "JSON_DUALITY_OBJECT" '(' [stringLit ':' Expression
		// (',' stringLit ':' Expression)*] ')' — the select-list
		// constructor of JSON duality views (MySQL 26.7 §15.1.14).
		r.advance()
		r.expect(int('('))
		x := &ast.JSONDualityObjectExpr{}
		if r.tok() != int(')') {
			for {
				key := r.expect(stringLit).lit
				r.expect(int(':'))
				x.Items = append(x.Items, &ast.JSONDualityObjectItem{Key: key, Value: r.parseExpression()})
				if !r.accept(int(',')) {
					break
				}
			}
		}
		r.expect(int(')'))
		return r.setOrigin(x, start)
	case row:
		// SimpleExpr: "ROW" '(' ExpressionList ',' Expression ')'
		r.advance()
		r.expect(int('('))
		values := []ast.ExprNode{r.parseExpression()}
		r.expect(int(','))
		values = append(values, r.parseExpression())
		for r.accept(int(',')) {
			values = append(values, r.parseExpression())
		}
		r.expect(int(')'))
		return r.setOrigin(&ast.RowExpr{Values: values}, start)
	case exists:
		// SimpleExpr: "EXISTS" SubSelect
		r.advance()
		sq := r.parseSubSelect()
		sq.Exists = true
		return r.setOrigin(&ast.ExistsSubqueryExpr{Sel: sq}, start)
	case int('{'):
		return r.parseODBCExpr(start)

	case caseKwd:
		// SimpleExpr: "CASE" ExpressionOpt WhenClauseList ElseOpt "END"
		r.advance()
		x := &ast.CaseExpr{}
		if r.tok() != when {
			x.Value = r.parseExpression()
		}
		for r.tok() == when {
			r.advance()
			cond := r.parseExpression()
			r.expect(then)
			result := r.parseExpression()
			x.WhenClauses = append(x.WhenClauses, &ast.WhenClause{Expr: cond, Result: result})
		}
		if len(x.WhenClauses) == 0 {
			r.syntaxError()
		}
		if r.accept(elseKwd) {
			x.ElseClause = r.parseExpression()
		}
		r.expect(end)
		return r.setOrigin(x, start)

	case builtinCast:
		// SimpleExpr: builtinCast '(' Expression ["AT" "TIME" "ZONE"
		// stringLit] "AS" CastType ArrayKwdOpt ')'
		r.advance()
		r.expect(int('('))
		expr := r.parseExpression()
		var atTimeZone string
		if r.tok() == at {
			r.advance()
			r.expect(timeType)
			r.expect(zone)
			atTimeZone = r.expect(stringLit).lit
		}
		r.expect(as)
		tp := r.parseCastType()
		isArray := r.accept(array)
		r.expect(int(')'))
		defaultFlen, defaultDecimal := mysql.GetDefaultFieldLengthAndDecimalForCast(tp.GetType())
		if tp.GetFlen() == types.UnspecifiedLength {
			tp.SetFlen(defaultFlen)
		}
		if tp.GetDecimal() == types.UnspecifiedLength {
			tp.SetDecimal(defaultDecimal)
		}
		tp.SetArray(isArray)
		explicitCharset := r.p.explicitCharset
		if isArray && !explicitCharset && tp.GetCharset() != charset.CharsetBin {
			tp.SetCharset(charset.CharsetUTF8MB4)
			tp.SetCollate(charset.CollationUTF8MB4)
		}
		r.p.explicitCharset = false
		return r.setOrigin(&ast.FuncCastExpr{
			Expr:            expr,
			Tp:              tp,
			FunctionType:    ast.CastFunction,
			ExplicitCharSet: explicitCharset,
			AtTimeZone:      atTimeZone,
		}, start)
	case jsonSumCrc32:
		// SimpleExpr: jsonSumCrc32 '(' Expression "AS" CastType "ARRAY" ')'
		r.advance()
		r.expect(int('('))
		expr := r.parseExpression()
		r.expect(as)
		tp := r.parseCastType()
		r.expect(array)
		r.expect(int(')'))
		defaultFlen, defaultDecimal := mysql.GetDefaultFieldLengthAndDecimalForCast(tp.GetType())
		if tp.GetFlen() == types.UnspecifiedLength {
			tp.SetFlen(defaultFlen)
		}
		if tp.GetDecimal() == types.UnspecifiedLength {
			tp.SetDecimal(defaultDecimal)
		}
		tp.SetArray(true)
		explicitCharset := r.p.explicitCharset
		if !explicitCharset && tp.GetCharset() != charset.CharsetBin {
			tp.SetCharset(charset.CharsetUTF8MB4)
			tp.SetCollate(charset.CollationUTF8MB4)
		}
		r.p.explicitCharset = false
		return r.setOrigin(&ast.JSONSumCrc32Expr{
			Expr:            expr,
			Tp:              tp,
			ExplicitCharSet: explicitCharset,
		}, start)
	case convert:
		// FunctionCallKeyword: "CONVERT" '(' Expression ',' CastType ')'
		//                    | "CONVERT" '(' Expression "USING" CharsetName ')'
		name := r.cur().lit
		r.advance()
		r.expect(int('('))
		expr := r.parseExpression()
		if r.accept(using) {
			cs := r.parseCharsetName()
			r.expect(int(')'))
			return r.setOrigin(&ast.FuncCallExpr{
				FnName: ast.NewCIStr(name),
				Args:   []ast.ExprNode{expr, ast.NewValueExpr(cs, "", "")},
			}, start)
		}
		r.expect(int(','))
		tp := r.parseCastType()
		r.expect(int(')'))
		defaultFlen, defaultDecimal := mysql.GetDefaultFieldLengthAndDecimalForCast(tp.GetType())
		if tp.GetFlen() == types.UnspecifiedLength {
			tp.SetFlen(defaultFlen)
		}
		if tp.GetDecimal() == types.UnspecifiedLength {
			tp.SetDecimal(defaultDecimal)
		}
		explicitCharset := r.p.explicitCharset
		r.p.explicitCharset = false
		return r.setOrigin(&ast.FuncCastExpr{
			Expr:            expr,
			Tp:              tp,
			FunctionType:    ast.CastConvertFunction,
			ExplicitCharSet: explicitCharset,
		}, start)
	case defaultKwd:
		// SimpleExpr: "DEFAULT" '(' SimpleIdent ')'
		r.advance()
		r.expect(int('('))
		col := r.parseSimpleIdent()
		r.expect(int(')'))
		return r.setOrigin(&ast.DefaultExpr{Name: col.Name}, start)
	case values:
		// SimpleExpr: "VALUES" '(' SimpleIdent ')' %prec lowerThanInsertValues
		r.advance()
		r.expect(int('('))
		colStart := r.cur().offset
		col := r.parseSimpleIdent()
		r.setOrigin(col, colStart)
		r.expect(int(')'))
		return r.setOrigin(&ast.ValuesExpr{Column: col}, start)

	case insert:
		// FunctionCallKeyword: "INSERT" '(' ExpressionListOpt ')'
		r.advance()
		r.expect(int('('))
		args := r.parseExpressionListOpt()
		r.expect(int(')'))
		return r.setOrigin(&ast.FuncCallExpr{FnName: ast.NewCIStr(ast.InsertFunc), Args: args}, start)
	case mod:
		// FunctionCallKeyword: "MOD" '(' Expression ',' Expression ')'
		r.advance()
		r.expect(int('('))
		l := r.parseExpression()
		r.expect(int(','))
		rr := r.parseExpression()
		r.expect(int(')'))
		return r.setOrigin(&ast.BinaryOperationExpr{Op: opcode.Mod, L: l, R: rr}, start)
	case password:
		if r.la(1) != int('(') {
			return r.parseSimpleIdentAtom(start)
		}
		r.advance()
		r.expect(int('('))
		args := r.parseExpressionListOpt()
		r.expect(int(')'))
		return r.setOrigin(&ast.FuncCallExpr{FnName: ast.NewCIStr(ast.PasswordFunc), Args: args}, start)
	case charType:
		// FunctionCallKeyword: "CHAR" '(' ExpressionList ')'
		//                    | "CHAR" '(' ExpressionList "USING" CharsetName ')'
		r.advance()
		r.expect(int('('))
		args := r.parseExpressionList()
		if r.accept(using) {
			cs := r.parseCharsetName()
			r.expect(int(')'))
			return r.setOrigin(&ast.FuncCallExpr{
				FnName: ast.NewCIStr(ast.CharFunc),
				Args:   append(args, ast.NewValueExpr(cs, "", "")),
			}, start)
		}
		r.expect(int(')'))
		return r.setOrigin(&ast.FuncCallExpr{
			FnName: ast.NewCIStr(ast.CharFunc),
			Args:   append(args, ast.NewValueExpr(nil, r.p.charset, r.p.collation)),
		}, start)

	case currentUser, currentDate, currentRole, utcDate, tidbCurrentTSO:
		// FunctionCallKeyword: FunctionNameOptionalBraces OptionalBraces
		name := r.cur().lit
		r.advance()
		if r.tok() == int('(') && r.la(1) == int(')') {
			r.advance()
			r.advance()
		}
		return r.setOrigin(&ast.FuncCallExpr{FnName: ast.NewCIStr(name)}, start)
	case currentTime, currentTs, localTime, localTs, utcTime, utcTimestamp:
		// FunctionCallKeyword: FunctionNameDatetimePrecision FuncDatetimePrec
		name := r.cur().lit
		r.advance()
		args := []ast.ExprNode{}
		if r.tok() == int('(') {
			r.advance()
			if r.tok() == intLit {
				args = append(args, ast.NewValueExpr(r.cur().item, r.p.charset, r.p.collation))
				r.advance()
			}
			r.expect(int(')'))
		}
		return r.setOrigin(&ast.FuncCallExpr{FnName: ast.NewCIStr(name), Args: args}, start)
	case builtinUser:
		name := r.cur().lit
		r.advance()
		r.expect(int('('))
		args := r.parseExpressionListOpt()
		r.expect(int(')'))
		return r.setOrigin(&ast.FuncCallExpr{FnName: ast.NewCIStr(name), Args: args}, start)
	case builtinCurDate:
		// FunctionCallKeyword: builtinCurDate '(' ')'
		name := r.cur().lit
		r.advance()
		r.expect(int('('))
		r.expect(int(')'))
		return r.setOrigin(&ast.FuncCallExpr{FnName: ast.NewCIStr(name)}, start)
	case builtinCurTime, builtinSysDate:
		// FunctionCallNonKeyword: builtinCurTime/builtinSysDate '(' FuncDatetimePrecListOpt ')'
		name := r.cur().lit
		r.advance()
		r.expect(int('('))
		args := []ast.ExprNode{}
		if r.tok() == intLit {
			args = append(args, ast.NewValueExpr(r.cur().item, r.p.charset, r.p.collation))
			r.advance()
		}
		r.expect(int(')'))
		return r.setOrigin(&ast.FuncCallExpr{FnName: ast.NewCIStr(name), Args: args}, start)

	case addDate, subDate:
		// FunctionCallNonKeyword: FunctionNameDateArithMultiForms
		name := r.cur().lit
		r.advance()
		r.expect(int('('))
		arg1 := r.parseExpression()
		r.expect(int(','))
		if r.accept(interval) {
			arg2 := r.parseExpression()
			unit := r.parseTimeUnit()
			r.expect(int(')'))
			return r.setOrigin(&ast.FuncCallExpr{
				FnName: ast.NewCIStr(name),
				Args:   []ast.ExprNode{arg1, arg2, &ast.TimeUnitExpr{Unit: unit}},
			}, start)
		}
		arg2 := r.parseExpression()
		r.expect(int(')'))
		return r.setOrigin(&ast.FuncCallExpr{
			FnName: ast.NewCIStr(name),
			Args:   []ast.ExprNode{arg1, arg2, &ast.TimeUnitExpr{Unit: ast.TimeUnitDay}},
		}, start)
	case builtinDateAdd, builtinDateSub:
		// FunctionCallNonKeyword: FunctionNameDateArith '(' Expression ',' "INTERVAL" Expression TimeUnit ')'
		name := r.cur().lit
		r.advance()
		r.expect(int('('))
		arg1 := r.parseExpression()
		r.expect(int(','))
		r.expect(interval)
		arg2 := r.parseExpression()
		unit := r.parseTimeUnit()
		r.expect(int(')'))
		return r.setOrigin(&ast.FuncCallExpr{
			FnName: ast.NewCIStr(name),
			Args:   []ast.ExprNode{arg1, arg2, &ast.TimeUnitExpr{Unit: unit}},
		}, start)
	case builtinExtract:
		// FunctionCallNonKeyword: builtinExtract '(' TimeUnit "FROM" Expression ')'
		name := r.cur().lit
		r.advance()
		r.expect(int('('))
		unit := r.parseTimeUnit()
		r.expect(from)
		expr := r.parseExpression()
		r.expect(int(')'))
		return r.setOrigin(&ast.FuncCallExpr{
			FnName: ast.NewCIStr(name),
			Args:   []ast.ExprNode{&ast.TimeUnitExpr{Unit: unit}, expr},
		}, start)
	case getFormat:
		// FunctionCallNonKeyword: "GET_FORMAT" '(' GetFormatSelector ',' Expression ')'
		name := r.cur().lit
		r.advance()
		r.expect(int('('))
		var sel ast.GetFormatSelectorType
		switch r.tok() {
		case dateType:
			sel = ast.GetFormatSelectorDate
		case datetimeType, timestampType:
			sel = ast.GetFormatSelectorDatetime
		case timeType:
			sel = ast.GetFormatSelectorTime
		default:
			r.syntaxError()
		}
		r.advance()
		r.expect(int(','))
		expr := r.parseExpression()
		r.expect(int(')'))
		return r.setOrigin(&ast.FuncCallExpr{
			FnName: ast.NewCIStr(name),
			Args:   []ast.ExprNode{&ast.GetFormatSelectorExpr{Selector: sel}, expr},
		}, start)
	case builtinPosition:
		// FunctionCallNonKeyword: builtinPosition '(' BitExpr "IN" Expression ')'
		name := r.cur().lit
		r.advance()
		r.expect(int('('))
		arg1 := r.parseBitExpr(0)
		r.expect(in)
		arg2 := r.parseExpression()
		r.expect(int(')'))
		return r.setOrigin(&ast.FuncCallExpr{FnName: ast.NewCIStr(name), Args: []ast.ExprNode{arg1, arg2}}, start)
	case builtinSubstring:
		// FunctionCallNonKeyword: builtinSubstring, the ','/FROM/FOR forms
		name := r.cur().lit
		r.advance()
		r.expect(int('('))
		arg1 := r.parseExpression()
		args := []ast.ExprNode{arg1}
		if r.accept(from) {
			args = append(args, r.parseExpression())
			if r.accept(forKwd) {
				args = append(args, r.parseExpression())
			}
		} else {
			r.expect(int(','))
			args = append(args, r.parseExpression())
			if r.accept(int(',')) {
				args = append(args, r.parseExpression())
			}
		}
		r.expect(int(')'))
		return r.setOrigin(&ast.FuncCallExpr{FnName: ast.NewCIStr(name), Args: args}, start)
	case timestampAdd, timestampDiff:
		// FunctionCallNonKeyword: "TIMESTAMPADD"/"TIMESTAMPDIFF" '(' TimestampUnit ',' Expression ',' Expression ')'
		// Both are identifier-class keywords when not followed by '('.
		if r.la(1) != int('(') {
			return r.parseSimpleIdentAtom(start)
		}
		name := r.cur().lit
		r.advance()
		r.expect(int('('))
		unit := r.parseTimestampUnit()
		r.expect(int(','))
		arg1 := r.parseExpression()
		r.expect(int(','))
		arg2 := r.parseExpression()
		r.expect(int(')'))
		return r.setOrigin(&ast.FuncCallExpr{
			FnName: ast.NewCIStr(name),
			Args:   []ast.ExprNode{&ast.TimeUnitExpr{Unit: unit}, arg1, arg2},
		}, start)
	case builtinTrim:
		return r.parseTrim(start)
	case weightString:
		return r.parseWeightString(start)
	case builtinTranslate:
		// FunctionCallNonKeyword: builtinTranslate '(' Expression ',' Expression ',' Expression ')'
		name := r.cur().lit
		r.advance()
		r.expect(int('('))
		a1 := r.parseExpression()
		r.expect(int(','))
		a2 := r.parseExpression()
		r.expect(int(','))
		a3 := r.parseExpression()
		r.expect(int(')'))
		return r.setOrigin(&ast.FuncCallExpr{FnName: ast.NewCIStr(name), Args: []ast.ExprNode{a1, a2, a3}}, start)
	case compress:
		if r.la(1) != int('(') {
			return r.parseSimpleIdentAtom(start)
		}
		// FunctionCallNonKeyword: "COMPRESS" '(' ExpressionListOpt ')'
		name := r.cur().lit
		r.advance()
		r.expect(int('('))
		args := r.parseExpressionListOpt()
		r.expect(int(')'))
		return r.setOrigin(&ast.FuncCallExpr{FnName: ast.NewCIStr(name), Args: args}, start)

	case lastval:
		// FunctionNameSequence: "LASTVAL" '(' TableName ')'
		if r.la(1) != int('(') {
			return r.parseSimpleIdentAtom(start)
		}
		r.advance()
		r.expect(int('('))
		tn := r.parseTableName()
		r.expect(int(')'))
		return r.setOrigin(&ast.FuncCallExpr{
			FnName: ast.NewCIStr(ast.LastVal),
			Args:   []ast.ExprNode{&ast.TableNameExpr{Name: tn}},
		}, start)
	case setval:
		// FunctionNameSequence: "SETVAL" '(' TableName ',' SignedNum ')'
		if r.la(1) != int('(') {
			return r.parseSimpleIdentAtom(start)
		}
		r.advance()
		r.expect(int('('))
		tn := r.parseTableName()
		r.expect(int(','))
		num := r.parseSignedNum()
		r.expect(int(')'))
		return r.setOrigin(&ast.FuncCallExpr{
			FnName: ast.NewCIStr(ast.SetVal),
			Args:   []ast.ExprNode{&ast.TableNameExpr{Name: tn}, ast.NewValueExpr(num, r.p.charset, r.p.collation)},
		}, start)
	case next:
		// NextValueForSequence: "NEXT" "VALUE" forKwd TableName
		if r.la(1) != value || r.la(2) != forKwd {
			return r.parseSimpleIdentAtom(start)
		}
		r.advance()
		r.expect(value)
		r.expect(forKwd)
		tn := r.parseTableName()
		return r.setOrigin(&ast.FuncCallExpr{
			FnName: ast.NewCIStr(ast.NextVal),
			Args:   []ast.ExprNode{&ast.TableNameExpr{Name: tn}},
		}, start)
	case nextval:
		// NextValueForSequence: "NEXTVAL" '(' TableName ')'
		if r.la(1) != int('(') {
			return r.parseSimpleIdentAtom(start)
		}
		r.advance()
		r.expect(int('('))
		tn := r.parseTableName()
		r.expect(int(')'))
		return r.setOrigin(&ast.FuncCallExpr{
			FnName: ast.NewCIStr(ast.NextVal),
			Args:   []ast.ExprNode{&ast.TableNameExpr{Name: tn}},
		}, start)

	case builtinApproxCountDistinct, builtinApproxPercentile, builtinBitAnd,
		builtinBitOr, builtinBitXor, builtinCount, builtinGroupConcat,
		builtinMax, builtinMin, builtinSum, builtinStddevPop, builtinStddevSamp,
		builtinVarPop, builtinVarSamp:
		return r.parseSumExprBuiltin(start)
	case avg:
		if r.la(1) != int('(') {
			return r.parseSimpleIdentAtom(start)
		}
		return r.parseSumExprBuiltin(start)
	case jsonArrayagg, jsonObjectAgg:
		if r.la(1) != int('(') {
			return r.parseSimpleIdentAtom(start)
		}
		return r.parseSumExprBuiltin(start)

	case rowNumber, rank, denseRank, cumeDist, percentRank, ntile, lead, lag,
		firstValue, lastValue, nthValue:
		return r.parseWindowFuncCall(start)

	default:
		if _, ok := functionNameConflictTok[tok]; ok && r.la(1) == int('(') {
			// FunctionCallKeyword: FunctionNameConflict '(' ExpressionListOpt ')'
			name := r.cur().lit
			r.advance()
			r.advance()
			args := r.parseExpressionListOpt()
			r.expect(int(')'))
			return r.setOrigin(&ast.FuncCallExpr{FnName: ast.NewCIStr(name), Args: args}, start)
		}
		switch tok {
		case dateType, timeType, timestampType:
			// FunctionCallKeyword: "DATE"/"TIME"/"TIMESTAMP" stringLit
			if r.la(1) == stringLit {
				var fn string
				switch tok {
				case dateType:
					fn = ast.DateLiteral
				case timeType:
					fn = ast.TimeLiteral
				default:
					fn = ast.TimestampLiteral
				}
				r.advance()
				lit := r.cur().lit
				r.advance()
				return r.setOrigin(&ast.FuncCallExpr{
					FnName: ast.NewCIStr(fn),
					Args:   []ast.ExprNode{ast.NewValueExpr(lit, "", "")},
				}, start)
			}
		}
		if isIdentifierTok(tok) {
			return r.parseSimpleIdentAtom(start)
		}
		r.syntaxError()
		return nil
	}
}

// parseParenExpr parses '(' Expression ')' (ParenthesesExpr, with the
// grammar action's node-text span) or '(' ExpressionList ',' Expression
// ')' (RowExpr).
func (r *rdParser) parseParenExpr(start int) ast.ExprNode {
	r.expect(int('('))
	exprStart := r.cur().offset
	expr := r.parseExpression()
	if r.tok() == int(',') {
		values := []ast.ExprNode{expr}
		for r.accept(int(',')) {
			values = append(values, r.parseExpression())
		}
		r.expect(int(')'))
		return r.setOrigin(&ast.RowExpr{Values: values}, start)
	}
	closing := r.expect(int(')'))
	r.setNodeTextSpan(expr, exprStart, r.endOffsetAt(closing.offset))
	return r.setOrigin(&ast.ParenthesesExpr{Expr: expr}, start)
}

// parseSimpleIdentAtom parses SimpleIdent, the qualified
// FunctionCallGeneric form, and the jss/juss JSON-extract postfixes
// (whose left side is restricted to SimpleIdent).
func (r *rdParser) parseSimpleIdentAtom(start int) ast.ExprNode {
	first := r.parseIdentifier()
	if r.tok() == int('.') && isIdentifierTok(r.la(1)) && r.la(2) == int('(') {
		// FunctionCallGeneric: Identifier '.' Identifier '(' ExpressionListOpt ')'
		r.advance()
		fnName := r.cur().lit
		r.advance()
		r.advance()
		args := r.parseExpressionListOpt()
		r.expect(int(')'))
		tp := ast.FuncCallExprTypeGeneric
		if isInTokenMap(fnName) {
			tp = ast.FuncCallExprTypeKeyword
		}
		return r.setOrigin(&ast.FuncCallExpr{
			Tp:     tp,
			Schema: ast.NewCIStr(first),
			FnName: ast.NewCIStr(fnName),
			Args:   args,
		}, start)
	}
	// The expr and its ColumnName are allocated as one block: column
	// references are the most-allocated node in typical queries, and the
	// two objects always live and die together.
	block := &struct {
		expr ast.ColumnNameExpr
		name ast.ColumnName
	}{}
	name := &block.name
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
	col := &block.expr
	col.Name = name
	r.setOrigin(col, start)
	switch r.tok() {
	case jss:
		// SimpleExpr: SimpleIdent jss stringLit
		r.advance()
		path := ast.NewValueExpr(r.expect(stringLit).lit, r.p.charset, r.p.collation)
		return r.setOrigin(&ast.FuncCallExpr{FnName: ast.NewCIStr(ast.JSONExtract), Args: []ast.ExprNode{col, path}}, start)
	case juss:
		// SimpleExpr: SimpleIdent juss stringLit
		r.advance()
		path := ast.NewValueExpr(r.expect(stringLit).lit, r.p.charset, r.p.collation)
		extract := &ast.FuncCallExpr{FnName: ast.NewCIStr(ast.JSONExtract), Args: []ast.ExprNode{col, path}}
		return r.setOrigin(&ast.FuncCallExpr{FnName: ast.NewCIStr(ast.JSONUnquote), Args: []ast.ExprNode{extract}}, start)
	}
	return col
}

// parseStringLiteral implements StringLiteral: adjacent string literals
// concatenate, remembering the first fragment's length as the projection
// offset.
func (r *rdParser) parseStringLiteral(start int) ast.ExprNode {
	expr := ast.NewValueExpr(r.expect(stringLit).lit, r.p.charset, r.p.collation)
	for r.tok() == stringLit {
		valExpr := expr
		strLit := valExpr.GetString()
		next := ast.NewValueExpr(strLit+r.cur().lit, r.p.charset, r.p.collation)
		if valExpr.GetProjectionOffset() >= 0 {
			next.SetProjectionOffset(valExpr.GetProjectionOffset())
		} else {
			next.SetProjectionOffset(len(strLit))
		}
		expr = next
		r.advance()
	}
	return r.setOrigin(expr, start)
}

// parseUnderscoreCharsetLiteral implements the "UNDERSCORE_CHARSET"
// stringLit/hexLit/bitLit alternatives of Literal.
func (r *rdParser) parseUnderscoreCharsetLiteral(start int) ast.ExprNode {
	cs := r.cur().lit
	r.advance()
	var val interface{}
	switch r.tok() {
	case stringLit:
		val = r.cur().lit
	case hexLit, bitLit:
		val = r.cur().item
	default:
		r.syntaxError()
	}
	r.advance()
	co, err := charset.GetDefaultCollationLegacy(cs)
	if err != nil {
		r.actionError(ast.ErrUnknownCharacterSet.GenWithStack("Unsupported character introducer: '%-.64s'", cs))
	}
	expr := ast.NewValueExpr(val, cs, co)
	tp := expr.GetType()
	tp.SetCharset(cs)
	tp.SetCollate(co)
	tp.AddFlag(mysql.UnderScoreCharsetFlag)
	if tp.GetCollate() == charset.CollationBin {
		tp.AddFlag(mysql.BinaryFlag)
	}
	return r.setOrigin(expr, start)
}

// parseSystemVariable implements SystemVariable.
func (r *rdParser) parseSystemVariable(start int) ast.ExprNode {
	v := strings.ToLower(r.cur().lit)
	r.advance()
	var isGlobal, isInstance bool
	explicitScope := true
	switch {
	case strings.HasPrefix(v, "@@global."):
		isGlobal = true
		v = strings.TrimPrefix(v, "@@global.")
	case strings.HasPrefix(v, "@@instance."):
		isInstance = true
		v = strings.TrimPrefix(v, "@@instance.")
	case strings.HasPrefix(v, "@@session."):
		v = strings.TrimPrefix(v, "@@session.")
	case strings.HasPrefix(v, "@@local."):
		v = strings.TrimPrefix(v, "@@local.")
	case strings.HasPrefix(v, "@@"):
		v, explicitScope = strings.TrimPrefix(v, "@@"), false
	}
	return r.setOrigin(&ast.VariableExpr{
		Name: v, IsGlobal: isGlobal, IsInstance: isInstance,
		IsSystem: true, ExplicitScope: explicitScope,
	}, start)
}

// parseODBCExpr implements SimpleExpr: '{' Identifier Expression '}'.
func (r *rdParser) parseODBCExpr(start int) ast.ExprNode {
	r.expect(int('{'))
	ident := r.parseIdentifier()
	expr := r.parseExpression()
	r.expect(int('}'))
	tp := expr.GetType()
	var fn string
	switch ident {
	case "d":
		fn = ast.DateLiteral
	case "t":
		fn = ast.TimeLiteral
	case "ts":
		fn = ast.TimestampLiteral
	default:
		return expr
	}
	tp.SetCharset("")
	tp.SetCollate("")
	return r.setOrigin(&ast.FuncCallExpr{FnName: ast.NewCIStr(fn), Args: []ast.ExprNode{expr}}, start)
}

// parseTrim implements the four builtinTrim alternatives.
func (r *rdParser) parseTrim(start int) ast.ExprNode {
	name := r.cur().lit
	r.advance()
	r.expect(int('('))
	var direction ast.TrimDirectionType
	switch r.tok() {
	case both:
		direction = ast.TrimBoth
	case leading:
		direction = ast.TrimLeading
	case trailing:
		direction = ast.TrimTrailing
	}
	if direction != ast.TrimBothDefault {
		r.advance()
		dirExpr := &ast.TrimDirectionExpr{Direction: direction}
		if r.accept(from) {
			// builtinTrim '(' TrimDirection "FROM" Expression ')'
			str := r.parseExpression()
			r.expect(int(')'))
			spaceVal := ast.NewValueExpr(" ", r.p.charset, r.p.collation)
			return r.setOrigin(&ast.FuncCallExpr{
				FnName: ast.NewCIStr(name),
				Args:   []ast.ExprNode{str, spaceVal, dirExpr},
			}, start)
		}
		// builtinTrim '(' TrimDirection Expression "FROM" Expression ')'
		pat := r.parseExpression()
		r.expect(from)
		str := r.parseExpression()
		r.expect(int(')'))
		return r.setOrigin(&ast.FuncCallExpr{
			FnName: ast.NewCIStr(name),
			Args:   []ast.ExprNode{str, pat, dirExpr},
		}, start)
	}
	expr := r.parseExpression()
	if r.accept(from) {
		// builtinTrim '(' Expression "FROM" Expression ')'
		str := r.parseExpression()
		r.expect(int(')'))
		return r.setOrigin(&ast.FuncCallExpr{FnName: ast.NewCIStr(name), Args: []ast.ExprNode{str, expr}}, start)
	}
	r.expect(int(')'))
	return r.setOrigin(&ast.FuncCallExpr{FnName: ast.NewCIStr(name), Args: []ast.ExprNode{expr}}, start)
}

// parseWeightString implements the weightString alternatives.
func (r *rdParser) parseWeightString(start int) ast.ExprNode {
	name := r.cur().lit
	r.advance()
	r.expect(int('('))
	expr := r.parseExpression()
	if r.accept(as) {
		var kind string
		switch r.tok() {
		case charType, character:
			kind = "CHAR"
		case binaryType:
			kind = "BINARY"
		default:
			r.syntaxError()
		}
		r.advance()
		flen := r.parseFieldLen()
		r.expect(int(')'))
		return r.setOrigin(&ast.FuncCallExpr{
			FnName: ast.NewCIStr(name),
			Args: []ast.ExprNode{
				expr,
				ast.NewValueExpr(kind, r.p.charset, r.p.collation),
				ast.NewValueExpr(flen, r.p.charset, r.p.collation),
			},
		}, start)
	}
	r.expect(int(')'))
	return r.setOrigin(&ast.FuncCallExpr{FnName: ast.NewCIStr(name), Args: []ast.ExprNode{expr}}, start)
}

// parseSignedNum implements SignedNum.
func (r *rdParser) parseSignedNum() int64 {
	switch r.tok() {
	case int('-'):
		r.advance()
		t := r.expect(intLit)
		unsigned := getUint64FromNUM(t.item)
		if unsigned > 9223372036854775808 {
			r.actionError(r.actionErrorf("the Signed Value should be at the range of [-9223372036854775808, 9223372036854775807]."))
		} else if unsigned == 9223372036854775808 {
			return -9223372036854775808
		}
		return -int64(unsigned)
	case int('+'):
		r.advance()
	}
	return r.parseInt64Num()
}

// parseInt64Num implements Int64Num.
func (r *rdParser) parseInt64Num() int64 {
	t := r.expect(intLit)
	v, errMsg := getInt64FromNUM(t.item)
	if errMsg != "" {
		r.actionError(r.actionErrorf("%s", errMsg))
	}
	return v
}

// parseSumExprBuiltin implements SumExpr for the current aggregate token.
func (r *rdParser) parseSumExprBuiltin(start int) ast.ExprNode {
	tok := r.tok()
	name := r.cur().lit
	r.advance()
	r.expect(int('('))
	distinct := false
	var args []ast.ExprNode
	var orderBy *ast.OrderByClause

	// BuggyDefaultFalseDistinctOpt / DistinctKwd / "ALL" prefixes.
	parseBuggyDistinct := func() bool {
		d := false
		if r.tok() == distinctKwdTok() || r.tok() == distinctRow {
			d = true
			r.advance()
			r.accept(all) // DistinctKwd "ALL"
		} else if r.accept(all) {
			d = false
		}
		return d
	}

	switch tok {
	case avg, builtinMax, builtinMin, builtinSum, builtinStddevPop,
		builtinStddevSamp, builtinVarPop, builtinVarSamp:
		distinct = parseBuggyDistinct()
		args = []ast.ExprNode{r.parseExpression()}
	case builtinApproxCountDistinct, builtinApproxPercentile:
		// builtinApproxCountDistinct '(' ExpressionList ')'
		args = r.parseExpressionList()
		r.expect(int(')'))
		f := name
		if tok == builtinApproxCountDistinct {
			return r.setOrigin(&ast.AggregateFuncExpr{F: f, Args: args, Distinct: false}, start)
		}
		return r.setOrigin(&ast.AggregateFuncExpr{F: f, Args: args}, start)
	case builtinBitAnd, builtinBitOr, builtinBitXor:
		r.accept(all)
		args = []ast.ExprNode{r.parseExpression()}
	case builtinCount:
		if r.tok() == distinctKwdTok() || r.tok() == distinctRow {
			// builtinCount '(' DistinctKwd ExpressionList ')'
			r.advance()
			args = r.parseExpressionList()
			r.expect(int(')'))
			return r.setOrigin(&ast.AggregateFuncExpr{F: name, Args: args, Distinct: true}, start)
		}
		if r.tok() == int('*') {
			r.advance()
			args = []ast.ExprNode{ast.NewValueExpr(1, r.p.charset, r.p.collation)}
		} else {
			r.accept(all)
			args = []ast.ExprNode{r.parseExpression()}
		}
	case builtinGroupConcat:
		// builtinGroupConcat '(' BuggyDefaultFalseDistinctOpt ExpressionList
		//     OrderByOptional OptGConcatSeparator ')' OptWindowingClause
		distinct = parseBuggyDistinct()
		args = r.parseExpressionList()
		if r.tok() == order {
			orderBy = r.parseOrderBy()
		}
		var sep ast.ExprNode = ast.NewValueExpr(",", "", "")
		if r.accept(separator) {
			sep = ast.NewValueExpr(r.expect(stringLit).lit, "", "")
		}
		args = append(args, sep)
	case jsonArrayagg:
		r.accept(all)
		args = []ast.ExprNode{r.parseExpression()}
	case jsonObjectAgg:
		r.accept(all)
		a1 := r.parseExpression()
		r.expect(int(','))
		r.accept(all)
		a2 := r.parseExpression()
		args = []ast.ExprNode{a1, a2}
	}
	r.expect(int(')'))

	// Stddev/variance aliases normalize the reported function name.
	switch tok {
	case builtinStddevPop:
		name = ast.AggFuncStddevPop
	case builtinVarPop:
		name = ast.AggFuncVarPop
	}

	if spec := r.parseOptWindowingClause(); spec != nil {
		w := &ast.WindowFuncExpr{Name: name, Args: args, Distinct: distinct, Spec: *spec}
		return r.setOrigin(w, start)
	}
	agg := &ast.AggregateFuncExpr{F: name, Args: args, Distinct: distinct}
	if orderBy != nil {
		agg.Order = orderBy
	}
	return r.setOrigin(agg, start)
}

// distinctKwdTok returns the DISTINCT token id (helper because the
// constant is named `distinct` which reads ambiguously at call sites).
func distinctKwdTok() int { return distinct }

// parseOptWindowingClause implements OptWindowingClause.
func (r *rdParser) parseOptWindowingClause() *ast.WindowSpec {
	if r.tok() != over {
		return nil
	}
	spec := r.parseWindowingClause()
	return &spec
}

// parseWindowingClause implements WindowingClause: "OVER" WindowNameOrSpec.
func (r *rdParser) parseWindowingClause() ast.WindowSpec {
	r.expect(over)
	if r.tok() == int('(') {
		return r.parseWindowSpec()
	}
	return ast.WindowSpec{Name: ast.NewCIStr(r.parseIdentifier()), OnlyAlias: true}
}

// parseWindowSpec implements WindowSpec: '(' WindowSpecDetails ')'.
func (r *rdParser) parseWindowSpec() ast.WindowSpec {
	r.expect(int('('))
	spec := ast.WindowSpec{}
	if isIdentifierTok(r.tok()) {
		spec.Ref = ast.NewCIStr(r.parseIdentifier())
	}
	if r.tok() == partition {
		r.advance()
		r.expect(by)
		spec.PartitionBy = &ast.PartitionByClause{Items: r.parseByList()}
	}
	if r.tok() == order {
		r.advance()
		r.expect(by)
		spec.OrderBy = &ast.OrderByClause{Items: r.parseByList()}
	}
	if frame := r.parseOptWindowFrameClause(); frame != nil {
		spec.Frame = frame
	}
	r.expect(int(')'))
	return spec
}

// parseOptWindowFrameClause implements OptWindowFrameClause.
func (r *rdParser) parseOptWindowFrameClause() *ast.FrameClause {
	var tp ast.FrameType
	switch r.tok() {
	case rows:
		tp = ast.FrameType(ast.Rows)
	case rangeKwd:
		tp = ast.FrameType(ast.Ranges)
	case groups:
		tp = ast.FrameType(ast.Groups)
	default:
		return nil
	}
	r.advance()
	var extent ast.FrameExtent
	if r.accept(between) {
		// WindowFrameBetween
		extent.Start = r.parseWindowFrameBound()
		r.expect(and)
		extent.End = r.parseWindowFrameBound()
	} else {
		extent.Start = r.parseWindowFrameBound()
		extent.End = ast.FrameBound{Type: ast.CurrentRow}
	}
	return &ast.FrameClause{Type: tp, Extent: extent}
}

// parseWindowFrameBound implements WindowFrameStart/WindowFrameBound.
func (r *rdParser) parseWindowFrameBound() ast.FrameBound {
	switch r.tok() {
	case unbounded:
		r.advance()
		switch r.tok() {
		case preceding:
			r.advance()
			return ast.FrameBound{Type: ast.Preceding, UnBounded: true}
		case following:
			r.advance()
			return ast.FrameBound{Type: ast.Following, UnBounded: true}
		}
		r.syntaxError()
	case current:
		r.advance()
		r.expect(row)
		return ast.FrameBound{Type: ast.CurrentRow}
	case interval:
		// "INTERVAL" Expression TimeUnit PRECEDING/FOLLOWING
		r.advance()
		expr := r.parseExpression()
		unit := r.parseTimeUnit()
		switch r.tok() {
		case preceding:
			r.advance()
			return ast.FrameBound{Type: ast.Preceding, Expr: expr, Unit: unit}
		case following:
			r.advance()
			return ast.FrameBound{Type: ast.Following, Expr: expr, Unit: unit}
		}
		r.syntaxError()
	case paramMarker:
		offset := r.cur().offset
		r.advance()
		expr := ast.NewParamMarkerExpr(offset)
		switch r.tok() {
		case preceding:
			r.advance()
			return ast.FrameBound{Type: ast.Preceding, Expr: expr}
		case following:
			r.advance()
			return ast.FrameBound{Type: ast.Following, Expr: expr}
		}
		r.syntaxError()
	case intLit, floatLit, decLit:
		// NumLiteral PRECEDING/FOLLOWING
		expr := ast.NewValueExpr(r.cur().item, r.p.charset, r.p.collation)
		r.advance()
		switch r.tok() {
		case preceding:
			r.advance()
			return ast.FrameBound{Type: ast.Preceding, Expr: expr}
		case following:
			r.advance()
			return ast.FrameBound{Type: ast.Following, Expr: expr}
		}
		r.syntaxError()
	}
	r.syntaxError()
	return ast.FrameBound{}
}

// parseByList implements ByList.
func (r *rdParser) parseByList() []*ast.ByItem {
	items := []*ast.ByItem{r.parseByItem()}
	for r.accept(int(',')) {
		items = append(items, r.parseByItem())
	}
	return items
}

// parseByItem implements ByItem, converting integer literals into
// PositionExpr the way the grammar action does.
func (r *rdParser) parseByItem() *ast.ByItem {
	expr := r.parseExpression()
	if valueExpr, ok := expr.(ast.ValueExpr); ok {
		if position, isPosition := valueExpr.GetValue().(int64); isPosition {
			expr = &ast.PositionExpr{N: int(position)}
		}
	}
	switch r.tok() {
	case asc:
		r.advance()
		return &ast.ByItem{Expr: expr, Desc: false}
	case desc:
		r.advance()
		return &ast.ByItem{Expr: expr, Desc: true}
	}
	return &ast.ByItem{Expr: expr, NullOrder: true}
}

// parseOrderBy implements OrderBy: "ORDER" "BY" ByList.
func (r *rdParser) parseOrderBy() *ast.OrderByClause {
	r.expect(order)
	r.expect(by)
	return &ast.OrderByClause{Items: r.parseByList()}
}

// parseWindowFuncCall implements WindowFuncCall.
func (r *rdParser) parseWindowFuncCall(start int) ast.ExprNode {
	tok := r.tok()
	name := r.cur().lit
	r.advance()
	r.expect(int('('))
	w := &ast.WindowFuncExpr{Name: name}
	switch tok {
	case rowNumber, rank, denseRank, cumeDist, percentRank:
		r.expect(int(')'))
	case ntile:
		w.Args = []ast.ExprNode{r.parseSimpleExpr()}
		r.expect(int(')'))
	case lead, lag:
		w.Args = []ast.ExprNode{r.parseExpression()}
		if r.accept(int(',')) {
			// OptLeadLagInfo: NumLiteral or paramMarker, then OptLLDefault
			switch r.tok() {
			case paramMarker:
				w.Args = append(w.Args, ast.NewParamMarkerExpr(r.cur().offset))
				r.advance()
			case intLit, floatLit, decLit:
				w.Args = append(w.Args, ast.NewValueExpr(r.cur().item, r.p.charset, r.p.collation))
				r.advance()
			default:
				r.syntaxError()
			}
			if r.accept(int(',')) {
				w.Args = append(w.Args, r.parseExpression())
			}
		}
		r.expect(int(')'))
		w.IgnoreNull = r.parseOptNullTreatment()
	case firstValue, lastValue:
		w.Args = []ast.ExprNode{r.parseExpression()}
		r.expect(int(')'))
		w.IgnoreNull = r.parseOptNullTreatment()
	case nthValue:
		a1 := r.parseExpression()
		r.expect(int(','))
		a2 := r.parseSimpleExpr()
		w.Args = []ast.ExprNode{a1, a2}
		r.expect(int(')'))
		w.FromLast = r.parseOptFromFirstLast()
		w.IgnoreNull = r.parseOptNullTreatment()
	}
	w.Spec = r.parseWindowingClause()
	return r.setOrigin(w, start)
}

// parseOptNullTreatment implements OptNullTreatment.
func (r *rdParser) parseOptNullTreatment() bool {
	switch r.tok() {
	case respect:
		r.advance()
		r.expect(nulls)
		return false
	case ignore:
		if r.la(1) == nulls {
			r.advance()
			r.advance()
			return true
		}
	}
	return false
}

// parseOptFromFirstLast implements OptFromFirstLast.
func (r *rdParser) parseOptFromFirstLast() bool {
	if r.tok() == from && (r.la(1) == first || r.la(1) == last) {
		fromLast := r.la(1) == last
		r.advance()
		r.advance()
		return fromLast
	}
	return false
}

// parseCastType implements CastType.
func (r *rdParser) parseCastType() *types.FieldType {
	switch r.tok() {
	case binaryType:
		// "BINARY" OptFieldLen
		r.advance()
		flen := r.parseOptFieldLen()
		tp := types.NewFieldType(mysql.TypeVarString)
		tp.SetFlen(flen)
		if tp.GetFlen() != types.UnspecifiedLength {
			tp.SetType(mysql.TypeString)
		}
		tp.SetCharset(charset.CharsetBin)
		tp.SetCollate(charset.CollationBin)
		tp.AddFlag(mysql.BinaryFlag)
		return tp
	case ncharType:
		// "NCHAR" OptFieldLen — the national char cast, which is CHAR in
		// the national character set (the connection charset here).
		r.advance()
		flen := r.parseOptFieldLen()
		tp := types.NewFieldType(mysql.TypeVarString)
		tp.SetFlen(flen)
		tp.SetCharset(r.p.charset)
		tp.SetCollate(r.p.collation)
		return tp
	case charType, character:
		// Char OptFieldLen OptBinary
		r.advance()
		flen := r.parseOptFieldLen()
		opt := r.parseOptBinary()
		tp := types.NewFieldType(mysql.TypeVarString)
		tp.SetFlen(flen)
		tp.SetCharset(opt.Charset)
		if opt.IsBinary {
			tp.AddFlag(mysql.BinaryFlag)
			tp.SetCharset(charset.CharsetBin)
			tp.SetCollate(charset.CollationBin)
		} else if tp.GetCharset() != "" {
			co, err := charset.GetDefaultCollation(tp.GetCharset())
			if err != nil {
				r.actionError(r.actionErrorf("Get collation error for charset: %s", tp.GetCharset()))
			}
			tp.SetCollate(co)
			r.p.explicitCharset = true
		} else {
			tp.SetCharset(r.p.charset)
			tp.SetCollate(r.p.collation)
		}
		return tp
	case dateType:
		r.advance()
		tp := types.NewFieldType(mysql.TypeDate)
		tp.SetCharset(charset.CharsetBin)
		tp.SetCollate(charset.CollationBin)
		tp.AddFlag(mysql.BinaryFlag)
		return tp
	case yearType:
		r.advance()
		tp := types.NewFieldType(mysql.TypeYear)
		tp.SetCharset(charset.CharsetBin)
		tp.SetCollate(charset.CollationBin)
		tp.AddFlag(mysql.BinaryFlag)
		return tp
	case datetimeType:
		// "DATETIME" OptFieldLen
		r.advance()
		dec := r.parseOptFieldLen()
		tp := types.NewFieldType(mysql.TypeDatetime)
		flen, _ := mysql.GetDefaultFieldLengthAndDecimalForCast(mysql.TypeDatetime)
		tp.SetFlen(flen)
		tp.SetDecimal(dec)
		if tp.GetDecimal() > 0 {
			tp.SetFlen(tp.GetFlen() + 1 + tp.GetDecimal())
		}
		tp.SetCharset(charset.CharsetBin)
		tp.SetCollate(charset.CollationBin)
		tp.AddFlag(mysql.BinaryFlag)
		return tp
	case decimalType:
		// "DECIMAL" FloatOpt
		r.advance()
		fopt := r.parseFloatOpt()
		tp := types.NewFieldType(mysql.TypeNewDecimal)
		tp.SetFlen(fopt.Flen)
		tp.SetDecimal(fopt.Decimal)
		tp.SetCharset(charset.CharsetBin)
		tp.SetCollate(charset.CollationBin)
		tp.AddFlag(mysql.BinaryFlag)
		return tp
	case timeType:
		r.advance()
		dec := r.parseOptFieldLen()
		tp := types.NewFieldType(mysql.TypeDuration)
		flen, _ := mysql.GetDefaultFieldLengthAndDecimalForCast(mysql.TypeDuration)
		tp.SetFlen(flen)
		tp.SetDecimal(dec)
		if tp.GetDecimal() > 0 {
			tp.SetFlen(tp.GetFlen() + 1 + tp.GetDecimal())
		}
		tp.SetCharset(charset.CharsetBin)
		tp.SetCollate(charset.CollationBin)
		tp.AddFlag(mysql.BinaryFlag)
		return tp
	case signed:
		r.advance()
		r.acceptInteger()
		tp := types.NewFieldType(mysql.TypeLonglong)
		tp.SetCharset(charset.CharsetBin)
		tp.SetCollate(charset.CollationBin)
		tp.AddFlag(mysql.BinaryFlag)
		return tp
	case unsigned:
		r.advance()
		r.acceptInteger()
		tp := types.NewFieldType(mysql.TypeLonglong)
		tp.AddFlag(mysql.UnsignedFlag | mysql.BinaryFlag)
		tp.SetCharset(charset.CharsetBin)
		tp.SetCollate(charset.CollationBin)
		return tp
	case jsonType:
		r.advance()
		tp := types.NewFieldType(mysql.TypeJSON)
		tp.AddFlag(mysql.BinaryFlag | mysql.ParseToJSONFlag)
		tp.SetCharset(mysql.DefaultCharset)
		tp.SetCollate(mysql.DefaultCollationName)
		return tp
	case doubleType:
		r.advance()
		tp := types.NewFieldType(mysql.TypeDouble)
		flen, decimal := mysql.GetDefaultFieldLengthAndDecimalForCast(mysql.TypeDouble)
		tp.SetFlen(flen)
		tp.SetDecimal(decimal)
		tp.AddFlag(mysql.BinaryFlag)
		tp.SetCharset(charset.CharsetBin)
		tp.SetCollate(charset.CollationBin)
		return tp
	case floatType:
		// "FLOAT" FloatOpt
		r.advance()
		fopt := r.parseFloatOpt()
		tp := types.NewFieldType(mysql.TypeFloat)
		if fopt.Flen >= 54 {
			r.sc.AppendError(ErrTooBigPrecision.GenWithStackByArgs(fopt.Flen, "CAST", 53))
		} else if fopt.Flen >= 25 {
			tp = types.NewFieldType(mysql.TypeDouble)
		}
		flen, decimal := mysql.GetDefaultFieldLengthAndDecimalForCast(tp.GetType())
		tp.SetFlen(flen)
		tp.SetDecimal(decimal)
		tp.AddFlag(mysql.BinaryFlag)
		tp.SetCharset(charset.CharsetBin)
		tp.SetCollate(charset.CollationBin)
		return tp
	case realType:
		r.advance()
		var tp *types.FieldType
		if r.sc.GetSQLMode().HasRealAsFloatMode() {
			tp = types.NewFieldType(mysql.TypeFloat)
		} else {
			tp = types.NewFieldType(mysql.TypeDouble)
		}
		flen, decimal := mysql.GetDefaultFieldLengthAndDecimalForCast(tp.GetType())
		tp.SetFlen(flen)
		tp.SetDecimal(decimal)
		tp.AddFlag(mysql.BinaryFlag)
		tp.SetCharset(charset.CharsetBin)
		tp.SetCollate(charset.CollationBin)
		return tp
	case geometryType, point, linestringType, polygonType, multipointType,
		multilinestringType, multipolygonType, geometryCollectionType:
		// CastType: the spatial types — castable since MySQL 8.0.24
		// (MySQL 26.7 §14.10); postdates the goyacc grammar.
		b := spatialTypeByte(r.tok())
		r.advance()
		tp := types.NewFieldType(b)
		tp.SetCharset(charset.CharsetBin)
		tp.SetCollate(charset.CollationBin)
		tp.AddFlag(mysql.BinaryFlag)
		return tp
	case vectorType:
		// "VECTOR" OptVectorElementType OptFieldLen
		r.advance()
		elemTp := mysql.TypeFloat
		if r.tok() == int('<') {
			r.advance()
			switch r.tok() {
			case floatType:
				elemTp = mysql.TypeFloat
			case doubleType:
				elemTp = mysql.TypeDouble
			default:
				r.syntaxError()
			}
			r.advance()
			r.expect(int('>'))
		}
		flen := r.parseOptFieldLen()
		if elemTp != mysql.TypeFloat {
			r.sc.AppendError(r.actionErrorf("Only VECTOR is supported for now"))
		}
		tp := types.NewFieldType(mysql.TypeTiDBVectorFloat32)
		tp.SetFlen(flen)
		tp.SetDecimal(0)
		tp.SetCharset(charset.CharsetBin)
		tp.SetCollate(charset.CollationBin)
		return tp
	}
	r.syntaxError()
	return nil
}

// acceptInteger consumes OptInteger.
func (r *rdParser) acceptInteger() {
	if r.tok() == integerType || r.tok() == intType {
		r.advance()
	}
}

// parseFloatOpt implements FloatOpt.
func (r *rdParser) parseFloatOpt() *ast.FloatOpt {
	if r.tok() != int('(') {
		return &ast.FloatOpt{Flen: types.UnspecifiedLength, Decimal: types.UnspecifiedLength}
	}
	r.expect(int('('))
	flen := int(r.parseLengthNum())
	if r.accept(int(',')) {
		dec := int(r.parseLengthNum())
		r.expect(int(')'))
		return &ast.FloatOpt{Flen: flen, Decimal: dec}
	}
	r.expect(int(')'))
	return &ast.FloatOpt{Flen: flen, Decimal: types.UnspecifiedLength}
}

// parseOptBinary implements OptBinary.
func (r *rdParser) parseOptBinary() *ast.OptBinary {
	switch r.tok() {
	case binaryType:
		r.advance()
		opt := &ast.OptBinary{IsBinary: true}
		if r.isCharsetKw() {
			r.parseCharsetKw()
			opt.Charset = r.parseCharsetName()
		}
		return opt
	default:
		if r.isCharsetKw() {
			r.parseCharsetKw()
			cs := r.parseCharsetName()
			isBin := r.accept(binaryType)
			return &ast.OptBinary{IsBinary: isBin, Charset: cs}
		}
		return &ast.OptBinary{}
	}
}

// isCharsetKw reports whether CharsetKw starts here:
// "CHARACTER" "SET" | "CHARSET" | "CHAR" "SET".
func (r *rdParser) isCharsetKw() bool {
	switch r.tok() {
	case charsetKwd:
		return true
	case character, charType:
		return r.la(1) == set
	}
	return false
}

// parseCharsetKw consumes CharsetKw.
func (r *rdParser) parseCharsetKw() {
	switch r.tok() {
	case charsetKwd:
		r.advance()
	case character, charType:
		r.advance()
		r.expect(set)
	default:
		r.syntaxError()
	}
}

// setNodeTextSpan reproduces the '(' Expression ')' action:
// parser.setNodeText(expr, src[startOffset:endOffset]).
func (r *rdParser) setNodeTextSpan(n ast.Node, start, end int) {
	r.p.setNodeText(n, r.src[start:end])
}

// endOffsetAt reproduces parser.endOffset: trims whitespace backwards
// from the given offset.
func (r *rdParser) endOffsetAt(offset int) int {
	for offset > 0 && unicode.IsSpace(rune(r.src[offset-1])) {
		offset--
	}
	return offset
}

// parseJSONValueExpr implements the extended JSON_VALUE form
// (MySQL 26.7 §14.17.3):
//
//	"JSON_VALUE" '(' Expression ',' Expression ["RETURNING" CastType]
//	    [JSONOnResponse "ON" "EMPTY"] [JSONOnResponse "ON" "ERROR"] ')'
//
// At least one of the optional clauses must be present; the plain
// two-argument call stays a generic FuncCallExpr (the caller speculates,
// so failing here selects that route).
func (r *rdParser) parseJSONValueExpr() ast.ExprNode {
	r.expect(identifier) // JSON_VALUE
	r.expect(int('('))
	x := &ast.JSONValueExpr{Doc: r.parseExpression()}
	r.expect(int(','))
	x.Path = r.parseExpression()
	if r.accept(returning) {
		x.Returning = r.parseCastType()
	}
	r.parseJSONOnResponses(&x.OnEmpty, &x.OnError)
	if x.Returning == nil && x.OnEmpty == nil && x.OnError == nil {
		r.syntaxError()
	}
	r.expect(int(')'))
	return x
}

// parseJSONOnResponses parses the [on_empty] [on_error] clauses shared by
// JSON_VALUE and JSON_TABLE columns:
// ("NULL" | "ERROR" | "DEFAULT" SimpleExpr) "ON" ("EMPTY" | "ERROR").
func (r *rdParser) parseJSONOnResponses(onEmpty, onError **ast.JSONOnResponse) {
	for r.tok() == null || r.tok() == errorKwd || r.tok() == defaultKwd {
		resp := &ast.JSONOnResponse{}
		switch r.tok() {
		case null:
			r.advance()
			resp.Tp = ast.JSONOnResponseNull
		case errorKwd:
			r.advance()
			resp.Tp = ast.JSONOnResponseError
		case defaultKwd:
			r.advance()
			resp.Tp = ast.JSONOnResponseDefault
			resp.Value = r.parseSimpleExpr()
		}
		r.expect(on)
		if r.accept(empty) {
			*onEmpty = resp
		} else {
			r.expect(errorKwd)
			*onError = resp
		}
	}
}

// parseJSONTableExpr implements the JSON_TABLE table function
// (MySQL 26.7 §14.17.6):
//
//	"JSON_TABLE" '(' Expression ',' Expression "COLUMNS"
//	    '(' JSONTableColumn (',' JSONTableColumn)* ')' ')'
func (r *rdParser) parseJSONTableExpr() *ast.JSONTableExpr {
	r.expect(identifier) // JSON_TABLE
	r.expect(int('('))
	x := &ast.JSONTableExpr{Doc: r.parseExpression()}
	r.expect(int(','))
	x.Path = r.parseExpression()
	r.expect(columns)
	r.expect(int('('))
	x.Columns = []*ast.JSONTableColumn{r.parseJSONTableColumn()}
	for r.accept(int(',')) {
		x.Columns = append(x.Columns, r.parseJSONTableColumn())
	}
	r.expect(int(')'))
	r.expect(int(')'))
	return x
}

// parseJSONTableColumn implements one COLUMNS entry of JSON_TABLE:
//
//	Identifier "FOR" "ORDINALITY"
//	| Identifier Type ["EXISTS"] "PATH" Expression [on_empty] [on_error]
//	| "NESTED" ["PATH"] Expression "COLUMNS" '(' ... ')'
func (r *rdParser) parseJSONTableColumn() *ast.JSONTableColumn {
	if r.tok() == nested {
		r.advance()
		r.accept(path)
		col := &ast.JSONTableColumn{Nested: true, NestedPath: r.parseExpression()}
		r.expect(columns)
		r.expect(int('('))
		col.NestedColumns = []*ast.JSONTableColumn{r.parseJSONTableColumn()}
		for r.accept(int(',')) {
			col.NestedColumns = append(col.NestedColumns, r.parseJSONTableColumn())
		}
		r.expect(int(')'))
		return col
	}
	name := ast.NewCIStr(r.parseIdentifier())
	if r.tok() == forKwd {
		r.advance()
		r.expect(ordinality)
		return &ast.JSONTableColumn{Name: name, ForOrdinality: true}
	}
	col := &ast.JSONTableColumn{Name: name, Tp: r.parseType()}
	col.Exists = r.accept(exists)
	r.expect(path)
	col.Path = r.parseExpression()
	r.parseJSONOnResponses(&col.OnEmpty, &col.OnError)
	return col
}
