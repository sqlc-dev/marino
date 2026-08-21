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

// Condition handling statements: SIGNAL, RESIGNAL, and GET DIAGNOSTICS.
// These postdate the goyacc grammar; the productions below are written
// from the MySQL 26.7 reference manual (§15.6.7, Condition Handling)
// in the same style as parser.y.

import (
	"strings"

	"github.com/sqlc-dev/marino/ast"
)

func init() {
	rdRegister(signal, (*rdParser).parseSignalStmt)
	rdRegister(resignal, (*rdParser).parseResignalStmt)
	rdRegister(get, (*rdParser).parseGetDiagnosticsStmt)
}

// signalInfoItemNames are the condition information item names a SIGNAL
// or RESIGNAL SET clause may assign. GET DIAGNOSTICS ... CONDITION reads
// the same items plus RETURNED_SQLSTATE (conditionInfoItemNames), and the
// statement form reads statementInfoItemNames. The names are not keywords
// (they are non-reserved in MySQL); they are matched case-insensitively
// in identifier position.
var signalInfoItemNames = map[string]bool{
	"CLASS_ORIGIN":       true,
	"SUBCLASS_ORIGIN":    true,
	"MESSAGE_TEXT":       true,
	"MYSQL_ERRNO":        true,
	"CONSTRAINT_CATALOG": true,
	"CONSTRAINT_SCHEMA":  true,
	"CONSTRAINT_NAME":    true,
	"CATALOG_NAME":       true,
	"SCHEMA_NAME":        true,
	"TABLE_NAME":         true,
	"COLUMN_NAME":        true,
	"CURSOR_NAME":        true,
}

var statementInfoItemNames = map[string]bool{
	"NUMBER":    true,
	"ROW_COUNT": true,
}

var conditionInfoItemNames = func() map[string]bool {
	m := map[string]bool{"RETURNED_SQLSTATE": true}
	for name := range signalInfoItemNames {
		m[name] = true
	}
	return m
}()

// parseSignalStmt implements SignalStmt:
// "SIGNAL" SignalConditionValue SignalSetOpt.
func (r *rdParser) parseSignalStmt() ast.StmtNode {
	r.expect(signal)
	return &ast.SignalStmt{
		Condition: r.parseSignalConditionValue(),
		SetItems:  r.parseSignalSetOpt(),
	}
}

// parseResignalStmt implements ResignalStmt:
// "RESIGNAL" SignalConditionValueOpt SignalSetOpt.
func (r *rdParser) parseResignalStmt() ast.StmtNode {
	r.expect(resignal)
	x := &ast.ResignalStmt{}
	if r.tok() == sqlstate || isIdentifierTok(r.tok()) {
		x.Condition = r.parseSignalConditionValue()
	}
	x.SetItems = r.parseSignalSetOpt()
	return x
}

// parseSignalConditionValue implements SignalConditionValue:
//
//	"SQLSTATE" optValue stringLit | Identifier
//
// The identifier alternative names a condition declared with
// DECLARE ... CONDITION.
func (r *rdParser) parseSignalConditionValue() *ast.SignalConditionValue {
	if r.tok() == sqlstate {
		r.advance()
		// optValue: empty | "VALUE"
		r.accept(value)
		return &ast.SignalConditionValue{SQLState: r.expect(stringLit).lit}
	}
	return &ast.SignalConditionValue{ConditionName: r.parseIdentifier()}
}

// parseSignalSetOpt implements SignalSetOpt:
// empty | "SET" SignalSetItem (',' SignalSetItem)*.
func (r *rdParser) parseSignalSetOpt() []*ast.SignalSetItem {
	if !r.accept(set) {
		return nil
	}
	items := []*ast.SignalSetItem{r.parseSignalSetItem()}
	for r.accept(int(',')) {
		items = append(items, r.parseSignalSetItem())
	}
	return items
}

// parseSignalSetItem implements SignalSetItem (signal_information_item):
// SignalInfoItemName eq SignalAllowedExpr.
func (r *rdParser) parseSignalSetItem() *ast.SignalSetItem {
	name := r.parseInfoItemName(signalInfoItemNames)
	r.expect(eq)
	return &ast.SignalSetItem{Name: name, Value: r.parseSignalAllowedExpr()}
}

// parseInfoItemName consumes an Identifier and validates it against the
// allowed information item names of the enclosing production, returning
// the canonical uppercase spelling.
func (r *rdParser) parseInfoItemName(allowed map[string]bool) string {
	if !isIdentifierTok(r.tok()) {
		r.syntaxError()
	}
	name := strings.ToUpper(r.cur().lit)
	if !allowed[name] {
		r.syntaxError()
	}
	r.advance()
	return name
}

// parseSignalAllowedExpr implements SignalAllowedExpr (the manual's
// simple_value_specification): Literal | UserVariable | SystemVariable
// | Identifier, the identifier being a stored program variable.
func (r *rdParser) parseSignalAllowedExpr() ast.ExprNode {
	start := r.cur().offset
	switch {
	case r.tok() == singleAtIdentifier:
		return r.parseUserVariable()
	case r.tok() == doubleAtIdentifier:
		return r.parseSystemVariable(start)
	case isIdentifierTok(r.tok()):
		name := &ast.ColumnName{Name: ast.NewCIStr(r.parseIdentifier())}
		return r.setOrigin(&ast.ColumnNameExpr{Name: name}, start)
	}
	return r.parseLiteralExpr(start)
}

// parseGetDiagnosticsStmt implements GetDiagnosticsStmt:
//
//	"GET" OptDiagnosticsArea "DIAGNOSTICS" (
//	    DiagnosticsItem (',' DiagnosticsItem)*
//	    | "CONDITION" SignalAllowedExpr DiagnosticsItem (',' DiagnosticsItem)* )
//
// with OptDiagnosticsArea: empty | "CURRENT" | "STACKED". The statement
// form assigns statementInfoItemNames, the CONDITION form
// conditionInfoItemNames.
func (r *rdParser) parseGetDiagnosticsStmt() ast.StmtNode {
	r.expect(get)
	x := &ast.GetDiagnosticsStmt{}
	switch r.tok() {
	case current:
		r.advance()
		x.Area = ast.DiagnosticsAreaCurrent
	case stacked:
		r.advance()
		x.Area = ast.DiagnosticsAreaStacked
	}
	r.expect(diagnostics)
	names := statementInfoItemNames
	if r.tok() == condition {
		r.advance()
		x.ConditionNumber = r.parseSignalAllowedExpr()
		names = conditionInfoItemNames
	}
	x.Items = []*ast.DiagnosticsItem{r.parseDiagnosticsItem(names)}
	for r.accept(int(',')) {
		x.Items = append(x.Items, r.parseDiagnosticsItem(names))
	}
	return x
}

// parseDiagnosticsItem implements DiagnosticsItem
// (statement_information_item / condition_information_item):
// DiagnosticsTarget eq InfoItemName.
func (r *rdParser) parseDiagnosticsItem(allowed map[string]bool) *ast.DiagnosticsItem {
	item := &ast.DiagnosticsItem{Target: r.parseDiagnosticsTarget()}
	r.expect(eq)
	item.Name = r.parseInfoItemName(allowed)
	return item
}

// parseDiagnosticsTarget implements DiagnosticsTarget (the manual's
// simple_target_specification): UserVariable | Identifier, the
// identifier being a stored program variable.
func (r *rdParser) parseDiagnosticsTarget() ast.ExprNode {
	start := r.cur().offset
	if r.tok() == singleAtIdentifier {
		return r.parseUserVariable()
	}
	name := &ast.ColumnName{Name: ast.NewCIStr(r.parseIdentifier())}
	return r.setOrigin(&ast.ColumnNameExpr{Name: name}, start)
}
