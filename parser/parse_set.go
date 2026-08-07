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

// SET statements: SetStmt (variable assignments, SET PASSWORD, SET
// TRANSACTION, SET CONFIG, SET SESSION_STATES, SET RESOURCE GROUP),
// SetRoleStmt, SetDefaultRoleStmt, and the Username/Rolename helper
// productions shared with the account-management statements.

import (
	"strings"

	"github.com/sqlc-dev/marino/ast"
	"github.com/sqlc-dev/marino/auth"
)

func init() {
	rdRegister(set, (*rdParser).parseSetStmtFamily)
}

// parseSetStmtFamily dispatches the SET-leading statement families:
//
//	SetStmt | SetRoleStmt | SetDefaultRoleStmt | SetBindingStmt
//
// distinguished by lookahead the way the LALR tables do (an unreserved
// keyword after SET followed by "="/":=" is a plain variable name).
func (r *rdParser) parseSetStmtFamily() ast.StmtNode {
	switch r.la(1) {
	case role:
		// SetRoleStmt: "SET" "ROLE" SetRoleOpt — unless `role` is a
		// variable name (SET role = ...).
		if r.la(2) != eq && r.la(2) != assignmentEq && r.la(2) != int('.') {
			return r.parseSetRoleStmt()
		}
	case defaultKwd:
		// "DEFAULT" is reserved: only SetDefaultRoleStmt can follow.
		return r.parseSetDefaultRoleStmt()
	case binding:
		// SetBindingStmt: "SET" "BINDING" BindingStatusType "FOR" ... —
		// unless `binding` is a variable name.
		if r.la(2) != eq && r.la(2) != assignmentEq && r.la(2) != int('.') {
			return r.parseSetBindingStmt()
		}
	}
	return r.parseSetStmt()
}

// parseSetStmt implements SetStmt.
func (r *rdParser) parseSetStmt() ast.StmtNode {
	r.expect(set)
	switch r.tok() {
	case password:
		// "SET" "PASSWORD" EqOrAssignmentEq PasswordOpt
		// "SET" "PASSWORD" "FOR" Username EqOrAssignmentEq PasswordOpt
		// The LALR machine shifts toward these on "="/":="/FOR (shift
		// beats reducing PASSWORD to Identifier).
		if r.la(1) == forKwd || r.la(1) == eq || r.la(1) == assignmentEq {
			return r.parseSetPwdStmt()
		}
	case global, session:
		// "SET" "GLOBAL"/"SESSION" "TRANSACTION" TransactionChars —
		// TransactionChar starts with ISOLATION or READ; anything else
		// makes `transaction` a variable name.
		if r.la(1) == transaction && (r.la(2) == isolation || r.la(2) == read) {
			isGlobal := r.tok() == global
			r.advance()
			r.advance()
			vars := r.parseTransactionChars()
			if isGlobal {
				for _, v := range vars {
					v.IsGlobal = true
				}
			}
			return &ast.SetStmt{Variables: vars}
		}
	case transaction:
		// "SET" "TRANSACTION" TransactionChars
		if r.la(1) == isolation || r.la(1) == read {
			r.advance()
			assigns := r.parseTransactionChars()
			for i := 0; i < len(assigns); i++ {
				if assigns[i].Name == "tx_isolation" {
					// A special session variable that make setting tx_isolation take effect one time.
					assigns[i].Name = "tx_isolation_one_shot"
				}
			}
			return &ast.SetStmt{Variables: assigns}
		}
	case config:
		// "SET" "CONFIG" (Identifier|stringLit) ConfigItemName EqOrAssignmentEq SetExpr
		if isIdentifierTok(r.la(1)) || r.la(1) == stringLit {
			return r.parseSetConfigStmt()
		}
	case sessionStates:
		// "SET" "SESSION_STATES" stringLit
		if r.la(1) == stringLit {
			r.advance()
			t := r.expect(stringLit)
			return &ast.SetSessionStatesStmt{SessionStates: t.lit}
		}
	case resource:
		// "SET" "RESOURCE" "GROUP" ResourceGroupName
		if r.la(1) == group {
			r.advance()
			r.advance()
			return &ast.SetResourceGroupStmt{Name: ast.NewCIStr(r.parseResourceGroupName())}
		}
	}
	// "SET" VariableAssignmentList
	return &ast.SetStmt{Variables: r.parseVariableAssignmentList()}
}

// parseSetPwdStmt implements the SET PASSWORD alternatives of SetStmt
// (SET already consumed).
func (r *rdParser) parseSetPwdStmt() ast.StmtNode {
	r.expect(password)
	if r.accept(forKwd) {
		user := r.parseUsername()
		r.parseEqOrAssignmentEq()
		return &ast.SetPwdStmt{User: user, Password: r.parsePasswordOpt()}
	}
	r.parseEqOrAssignmentEq()
	return &ast.SetPwdStmt{Password: r.parsePasswordOpt()}
}

// parseSetConfigStmt implements the SET CONFIG alternatives of SetStmt
// (SET already consumed).
func (r *rdParser) parseSetConfigStmt() ast.StmtNode {
	r.expect(config)
	st := &ast.SetConfigStmt{}
	if r.tok() == stringLit {
		// "SET" "CONFIG" stringLit ConfigItemName EqOrAssignmentEq SetExpr
		st.Instance = r.cur().lit
		r.advance()
	} else {
		// "SET" "CONFIG" Identifier ConfigItemName EqOrAssignmentEq SetExpr
		st.Type = strings.ToLower(r.parseIdentifier())
	}
	st.Name = r.parseConfigItemName()
	r.parseEqOrAssignmentEq()
	st.Value = r.parseSetExpr()
	return st
}

// parseConfigItemName implements ConfigItemName:
// Identifier | Identifier '.' ConfigItemName | Identifier '-' ConfigItemName.
func (r *rdParser) parseConfigItemName() string {
	name := r.parseIdentifier()
	for r.tok() == int('.') || r.tok() == int('-') {
		sep := "."
		if r.tok() == int('-') {
			sep = "-"
		}
		r.advance()
		name += sep + r.parseIdentifier()
	}
	return name
}

// parseResourceGroupName implements ResourceGroupName: Identifier | "DEFAULT".
func (r *rdParser) parseResourceGroupName() string {
	if r.tok() == defaultKwd {
		lit := r.cur().lit
		r.advance()
		return lit
	}
	return r.parseIdentifier()
}

// parseEqOrAssignmentEq implements EqOrAssignmentEq: eq | assignmentEq.
func (r *rdParser) parseEqOrAssignmentEq() {
	if !r.accept(eq) && !r.accept(assignmentEq) {
		r.syntaxError()
	}
}

// parseTransactionChars implements TransactionChars.
func (r *rdParser) parseTransactionChars() []*ast.VariableAssignment {
	vars := r.parseTransactionChar()
	for r.accept(int(',')) {
		vars = append(vars, r.parseTransactionChar()...)
	}
	return vars
}

// parseTransactionChar implements TransactionChar.
func (r *rdParser) parseTransactionChar() []*ast.VariableAssignment {
	switch r.tok() {
	case isolation:
		// "ISOLATION" "LEVEL" IsolationLevel
		r.advance()
		r.expect(level)
		isoLevel := r.parseIsolationLevel()
		varAssigns := []*ast.VariableAssignment{}
		expr := ast.NewValueExpr(isoLevel, r.p.charset, r.p.collation)
		varAssigns = append(varAssigns, &ast.VariableAssignment{Name: "tx_isolation", Value: expr, IsSystem: true})
		return varAssigns
	case read:
		r.advance()
		if r.accept(write) {
			// "READ" "WRITE"
			varAssigns := []*ast.VariableAssignment{}
			expr := ast.NewValueExpr("0", r.p.charset, r.p.collation)
			varAssigns = append(varAssigns, &ast.VariableAssignment{Name: "tx_read_only", Value: expr, IsSystem: true})
			return varAssigns
		}
		r.expect(only)
		if r.tok() == asof {
			// "READ" "ONLY" AsOfClause
			varAssigns := []*ast.VariableAssignment{}
			asof := r.parseAsOfClause()
			if asof != nil {
				varAssigns = append(varAssigns, &ast.VariableAssignment{Name: "tx_read_ts", Value: asof.TsExpr, IsSystem: true})
			}
			return varAssigns
		}
		// "READ" "ONLY"
		varAssigns := []*ast.VariableAssignment{}
		expr := ast.NewValueExpr("1", r.p.charset, r.p.collation)
		varAssigns = append(varAssigns, &ast.VariableAssignment{Name: "tx_read_only", Value: expr, IsSystem: true})
		return varAssigns
	}
	r.syntaxError()
	return nil
}

// parseIsolationLevel implements IsolationLevel.
func (r *rdParser) parseIsolationLevel() string {
	switch r.tok() {
	case repeatable:
		// "REPEATABLE" "READ"
		r.advance()
		r.expect(read)
		return ast.RepeatableRead
	case read:
		r.advance()
		switch r.tok() {
		case committed:
			// "READ" "COMMITTED"
			r.advance()
			return ast.ReadCommitted
		case uncommitted:
			// "READ" "UNCOMMITTED"
			r.advance()
			return ast.ReadUncommitted
		}
	case serializable:
		// "SERIALIZABLE"
		r.advance()
		return ast.Serializable
	}
	r.syntaxError()
	return ""
}

// parseSetExpr implements SetExpr: "ON" | "BINARY" | ExprOrDefault.
func (r *rdParser) parseSetExpr() ast.ExprNode {
	switch r.tok() {
	case on:
		start := r.cur().offset
		r.advance()
		return r.setOrigin(ast.NewValueExpr("ON", r.p.charset, r.p.collation), start)
	case binaryType:
		// The bare "BINARY" value only when nothing that could continue a
		// BINARY expression follows (the reduce lookaheads).
		if next := r.la(1); next == int(',') || next == int(';') || next == 0 {
			start := r.cur().offset
			r.advance()
			return r.setOrigin(ast.NewValueExpr("BINARY", r.p.charset, r.p.collation), start)
		}
	}
	return r.parseExprOrDefault()
}

// parseVariableAssignmentList implements VariableAssignmentList.
func (r *rdParser) parseVariableAssignmentList() []*ast.VariableAssignment {
	list := []*ast.VariableAssignment{r.parseVariableAssignment()}
	for r.accept(int(',')) {
		list = append(list, r.parseVariableAssignment())
	}
	return list
}

// parseVariableName implements VariableName: Identifier | Identifier '.' Identifier.
func (r *rdParser) parseVariableName() string {
	name := r.parseIdentifier()
	if r.accept(int('.')) {
		name = name + "." + r.parseIdentifier()
	}
	return name
}

// parseVariableAssignment implements VariableAssignment.
func (r *rdParser) parseVariableAssignment() *ast.VariableAssignment {
	switch r.tok() {
	case global:
		// "GLOBAL" VariableName EqOrAssignmentEq SetExpr
		if isIdentifierTok(r.la(1)) {
			r.advance()
			name := r.parseVariableName()
			r.parseEqOrAssignmentEq()
			return &ast.VariableAssignment{Name: name, Value: r.parseSetExpr(), IsGlobal: true, IsSystem: true}
		}
	case instance:
		// "INSTANCE" VariableName EqOrAssignmentEq SetExpr
		if isIdentifierTok(r.la(1)) {
			r.advance()
			name := r.parseVariableName()
			r.parseEqOrAssignmentEq()
			return &ast.VariableAssignment{Name: name, Value: r.parseSetExpr(), IsInstance: true, IsSystem: true}
		}
	case session, local:
		// "SESSION"/"LOCAL" VariableName EqOrAssignmentEq SetExpr
		if isIdentifierTok(r.la(1)) {
			r.advance()
			name := r.parseVariableName()
			r.parseEqOrAssignmentEq()
			return &ast.VariableAssignment{Name: name, Value: r.parseSetExpr(), IsSystem: true}
		}
	case doubleAtIdentifier:
		// doubleAtIdentifier EqOrAssignmentEq SetExpr
		v := strings.ToLower(r.cur().lit)
		r.advance()
		r.parseEqOrAssignmentEq()
		var isGlobal bool
		var isInstance bool
		if strings.HasPrefix(v, "@@global.") {
			isGlobal = true
			v = strings.TrimPrefix(v, "@@global.")
		} else if strings.HasPrefix(v, "@@instance.") {
			isInstance = true
			v = strings.TrimPrefix(v, "@@instance.")
		} else if strings.HasPrefix(v, "@@session.") {
			v = strings.TrimPrefix(v, "@@session.")
		} else if strings.HasPrefix(v, "@@local.") {
			v = strings.TrimPrefix(v, "@@local.")
		} else if strings.HasPrefix(v, "@@") {
			v = strings.TrimPrefix(v, "@@")
		}
		return &ast.VariableAssignment{Name: v, Value: r.parseSetExpr(), IsGlobal: isGlobal, IsInstance: isInstance, IsSystem: true}
	case singleAtIdentifier:
		// singleAtIdentifier EqOrAssignmentEq Expression
		v := strings.TrimPrefix(r.cur().lit, "@")
		r.advance()
		r.parseEqOrAssignmentEq()
		return &ast.VariableAssignment{Name: v, Value: r.parseExpression()}
	case names:
		// "NAMES" CharsetName ["COLLATE" ("DEFAULT"|StringName)] | "NAMES" "DEFAULT"
		// — unless `names` is a variable name.
		if r.la(1) != eq && r.la(1) != assignmentEq && r.la(1) != int('.') {
			r.advance()
			if r.tok() == defaultKwd {
				r.advance()
				v := &ast.DefaultExpr{}
				return &ast.VariableAssignment{Name: ast.SetNames, Value: v}
			}
			cs := r.parseCharsetName()
			if r.accept(collate) {
				if r.tok() == defaultKwd {
					r.advance()
					return &ast.VariableAssignment{
						Name:  ast.SetNames,
						Value: ast.NewValueExpr(cs, "", ""),
					}
				}
				return &ast.VariableAssignment{
					Name:        ast.SetNames,
					Value:       ast.NewValueExpr(cs, "", ""),
					ExtendValue: ast.NewValueExpr(r.parseStringName(), "", ""),
				}
			}
			return &ast.VariableAssignment{
				Name:  ast.SetNames,
				Value: ast.NewValueExpr(cs, "", ""),
			}
		}
	}
	if r.isCharsetKw() {
		// CharsetKw CharsetNameOrDefault — for the single-token "CHARSET"
		// spelling, `charset` followed by "="/":="/"." is a variable name.
		if r.tok() != charsetKwd || (r.la(1) != eq && r.la(1) != assignmentEq && r.la(1) != int('.')) {
			r.parseCharsetKw()
			return &ast.VariableAssignment{Name: ast.SetCharset, Value: r.parseCharsetNameOrDefault()}
		}
	}
	// VariableName EqOrAssignmentEq SetExpr
	name := r.parseVariableName()
	r.parseEqOrAssignmentEq()
	return &ast.VariableAssignment{Name: name, Value: r.parseSetExpr(), IsSystem: true}
}

// parseCharsetNameOrDefault implements CharsetNameOrDefault (an <expr>
// production: the result is origin-stamped).
func (r *rdParser) parseCharsetNameOrDefault() ast.ExprNode {
	start := r.cur().offset
	if r.tok() == defaultKwd {
		r.advance()
		return r.setOrigin(&ast.DefaultExpr{}, start)
	}
	return r.setOrigin(ast.NewValueExpr(r.parseCharsetName(), "", ""), start)
}

// parseSetRoleStmt implements SetRoleStmt: "SET" "ROLE" SetRoleOpt.
func (r *rdParser) parseSetRoleStmt() ast.StmtNode {
	r.expect(set)
	r.expect(role)
	return r.parseSetRoleOpt()
}

// parseSetRoleOpt implements SetRoleOpt.
func (r *rdParser) parseSetRoleOpt() *ast.SetRoleStmt {
	switch r.tok() {
	case all:
		// "ALL" "EXCEPT" RolenameList
		if r.la(1) == except {
			r.advance()
			r.advance()
			return &ast.SetRoleStmt{SetRoleOpt: ast.SetRoleAllExcept, RoleList: r.parseRolenameList()}
		}
	case defaultKwd:
		// "DEFAULT"
		r.advance()
		return &ast.SetRoleStmt{SetRoleOpt: ast.SetRoleDefault, RoleList: nil}
	}
	return r.parseSetDefaultRoleOpt()
}

// parseSetDefaultRoleOpt implements SetDefaultRoleOpt.
func (r *rdParser) parseSetDefaultRoleOpt() *ast.SetRoleStmt {
	switch r.tok() {
	case none:
		r.advance()
		return &ast.SetRoleStmt{SetRoleOpt: ast.SetRoleNone, RoleList: nil}
	case all:
		r.advance()
		return &ast.SetRoleStmt{SetRoleOpt: ast.SetRoleAll, RoleList: nil}
	}
	return &ast.SetRoleStmt{SetRoleOpt: ast.SetRoleRegular, RoleList: r.parseRolenameList()}
}

// parseSetDefaultRoleStmt implements SetDefaultRoleStmt:
// "SET" "DEFAULT" "ROLE" SetDefaultRoleOpt "TO" UsernameList.
func (r *rdParser) parseSetDefaultRoleStmt() ast.StmtNode {
	r.expect(set)
	r.expect(defaultKwd)
	r.expect(role)
	tmp := r.parseSetDefaultRoleOpt()
	r.expect(to)
	return &ast.SetDefaultRoleStmt{
		SetRoleOpt: tmp.SetRoleOpt,
		RoleList:   tmp.RoleList,
		UserList:   r.parseUsernameList(),
	}
}

// parseUsername implements Username.
func (r *rdParser) parseUsername() *auth.UserIdentity {
	if r.tok() == currentUser {
		// "CURRENT_USER" OptionalBraces
		r.advance()
		if r.accept(int('(')) {
			r.expect(int(')'))
		}
		return &auth.UserIdentity{CurrentUser: true}
	}
	name := r.parseStringName()
	switch r.tok() {
	case int('@'):
		// StringName '@' StringName
		r.advance()
		return &auth.UserIdentity{Username: name, Hostname: strings.ToLower(r.parseStringName())}
	case singleAtIdentifier:
		// StringName singleAtIdentifier
		host := strings.ToLower(strings.TrimPrefix(r.cur().lit, "@"))
		r.advance()
		return &auth.UserIdentity{Username: name, Hostname: host}
	}
	return &auth.UserIdentity{Username: name, Hostname: "%"}
}

// parseUsernameList implements UsernameList.
func (r *rdParser) parseUsernameList() []*auth.UserIdentity {
	list := []*auth.UserIdentity{r.parseUsername()}
	for r.accept(int(',')) {
		list = append(list, r.parseUsername())
	}
	return list
}

// parsePasswordOpt implements PasswordOpt: stringLit | "PASSWORD" '(' AuthString ')'.
func (r *rdParser) parsePasswordOpt() string {
	if r.tok() == password {
		r.advance()
		r.expect(int('('))
		s := r.parseAuthString()
		r.expect(int(')'))
		return s
	}
	return r.expect(stringLit).lit
}

// parseAuthString implements AuthString: stringLit.
func (r *rdParser) parseAuthString() string {
	return r.expect(stringLit).lit
}

// parseRolename implements Rolename:
//
//	Rolename: RoleNameString | RolenameComposed
//	RoleNameString: stringLit | identifier
//	RolenameComposed: StringName '@' StringName | StringName singleAtIdentifier
//
// A bare role name must be a string literal or a plain identifier token
// (not an unreserved keyword); with a host part any StringName works.
func (r *rdParser) parseRolename() *auth.RoleIdentity {
	bareOK := r.tok() == stringLit || r.tok() == identifier
	name := r.parseStringName()
	switch r.tok() {
	case int('@'):
		r.advance()
		return &auth.RoleIdentity{Username: name, Hostname: strings.ToLower(r.parseStringName())}
	case singleAtIdentifier:
		host := strings.ToLower(strings.TrimPrefix(r.cur().lit, "@"))
		r.advance()
		return &auth.RoleIdentity{Username: name, Hostname: host}
	}
	if !bareOK {
		r.syntaxError()
	}
	return &auth.RoleIdentity{Username: name, Hostname: "%"}
}

// parseRolenameList implements RolenameList.
func (r *rdParser) parseRolenameList() []*auth.RoleIdentity {
	list := []*auth.RoleIdentity{r.parseRolename()}
	for r.accept(int(',')) {
		list = append(list, r.parseRolename())
	}
	return list
}
