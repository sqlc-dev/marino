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

// Account management: GRANT (privileges, PROXY, roles) and REVOKE
// (privileges, roles) statements with their helper productions.

import (
	"strings"

	"github.com/sqlc-dev/marino/ast"
	"github.com/sqlc-dev/marino/auth"
	"github.com/sqlc-dev/marino/mysql"
)

func init() {
	rdRegister(grant, (*rdParser).parseGrantStmtFamily)
	rdRegister(revoke, (*rdParser).parseRevokeStmtFamily)
}

// parseGrantStmtFamily implements:
//
//	GrantStmt: "GRANT" RoleOrPrivElemList "ON" ObjectType PrivLevel "TO"
//	           UserSpecList RequireClauseOpt WithGrantOptionOpt
//	GrantProxyStmt: "GRANT" "PROXY" "ON" Username "TO" UsernameList
//	                WithGrantOptionOpt
//	GrantRoleStmt: "GRANT" RoleOrPrivElemList "TO" UsernameList
//
// The three share the leading "GRANT". PROXY cannot start a
// RoleOrPrivElem unless composed into a role name with a following
// '@'/singleAtIdentifier, so "GRANT PROXY ON" is uniquely
// GrantProxyStmt; otherwise GrantStmt and GrantRoleStmt share
// RoleOrPrivElemList and are told apart by "ON" vs "TO", exactly how
// the LR automaton resolves them.
func (r *rdParser) parseGrantStmtFamily() ast.StmtNode {
	r.expect(grant)
	if r.tok() == proxy && r.la(1) == on {
		// GrantProxyStmt: "GRANT" "PROXY" "ON" Username "TO" UsernameList WithGrantOptionOpt
		r.advance()
		r.advance()
		localUser := r.parseGrantUsername()
		r.expect(to)
		externalUsers := r.parseGrantUsernameList()
		withGrant := r.parseWithGrantOptionOpt()
		return &ast.GrantProxyStmt{
			LocalUser:     localUser,
			ExternalUsers: externalUsers,
			WithGrant:     withGrant,
		}
	}
	elems := r.parseRoleOrPrivElemList()
	switch r.tok() {
	case on:
		// GrantStmt
		r.advance()
		objectType := r.parseObjectType()
		level := r.parsePrivLevel()
		r.expect(to)
		users := r.parseGrantUserSpecList()
		tlsOptions := r.parseGrantRequireClauseOpt()
		withGrant := r.parseWithGrantOptionOpt()
		p, err := convertToPriv(elems)
		if err != nil {
			r.actionError(err)
		}
		return &ast.GrantStmt{
			Privs:                 p,
			ObjectType:            objectType,
			Level:                 level,
			Users:                 users,
			AuthTokenOrTLSOptions: tlsOptions,
			WithGrant:             withGrant,
			AsUser:                r.parseGrantAsOpt(),
		}
	case to:
		// GrantRoleStmt
		r.advance()
		users := r.parseGrantUsernameList()
		roles, err := convertToRole(elems)
		if err != nil {
			r.actionError(err)
		}
		stmt := &ast.GrantRoleStmt{
			Roles: roles,
			Users: users,
		}
		if r.tok() == with && r.la(1) == admin {
			// "WITH" "ADMIN" "OPTION"
			r.advance()
			r.advance()
			r.expect(option)
			stmt.AdminOption = true
		}
		return stmt
	default:
		r.syntaxError()
		return nil
	}
}

// parseGrantAsOpt implements the AS user [WITH ROLE ...] clause of
// GrantStmt (nil for the empty alternative):
//
//	empty | "AS" Username ["WITH" "ROLE" ("DEFAULT" | "NONE" | "ALL"
//	    | "ALL" "EXCEPT" RolenameList | RolenameList)]
func (r *rdParser) parseGrantAsOpt() *ast.GrantAs {
	if r.tok() != as {
		return nil
	}
	r.advance()
	grantAs := &ast.GrantAs{User: r.parseGrantUsername()}
	if r.tok() == with && r.la(1) == role {
		r.advance()
		r.advance()
		grantAs.WithRole = true
		switch r.tok() {
		case defaultKwd:
			r.advance()
			grantAs.RoleOpt = ast.SetRoleDefault
		case none:
			r.advance()
			grantAs.RoleOpt = ast.SetRoleNone
		case all:
			r.advance()
			if r.accept(except) {
				grantAs.RoleOpt = ast.SetRoleAllExcept
				grantAs.Roles = r.parseRolenameList()
			} else {
				grantAs.RoleOpt = ast.SetRoleAll
			}
		default:
			grantAs.RoleOpt = ast.SetRoleRegular
			grantAs.Roles = r.parseRolenameList()
		}
	}
	return grantAs
}

// parseRevokeStmtFamily implements:
//
//	RevokeStmt: "REVOKE" RoleOrPrivElemList "ON" ObjectType PrivLevel
//	            "FROM" UserSpecList
//	RevokeRoleStmt: "REVOKE" RoleOrPrivElemList "FROM" UsernameList
//
// disambiguated by "ON" vs "FROM" after the shared list.
func (r *rdParser) parseRevokeStmtFamily() ast.StmtNode {
	r.expect(revoke)
	ifExists := r.parseIfExists()
	if r.tok() == proxy && r.la(1) == on {
		// RevokeProxyStmt: "REVOKE" "PROXY" "ON" Username "FROM" UsernameList
		r.advance()
		r.advance()
		localUser := r.parseGrantUsername()
		r.expect(from)
		return &ast.RevokeProxyStmt{
			LocalUser:     localUser,
			ExternalUsers: r.parseGrantUsernameList(),
		}
	}
	elems := r.parseRoleOrPrivElemList()
	switch r.tok() {
	case on:
		// RevokeStmt
		r.advance()
		objectType := r.parseObjectType()
		level := r.parsePrivLevel()
		r.expect(from)
		users := r.parseGrantUserSpecList()
		p, err := convertToPriv(elems)
		if err != nil {
			r.actionError(err)
		}
		return &ast.RevokeStmt{
			Privs:             p,
			ObjectType:        objectType,
			Level:             level,
			Users:             users,
			IfExists:          ifExists,
			IgnoreUnknownUser: r.parseIgnoreUnknownUserOpt(),
		}
	case from:
		// RevokeRoleStmt
		r.advance()
		usernames := r.parseGrantUsernameList()
		// MySQL has special syntax for REVOKE ALL [PRIVILEGES], GRANT OPTION
		// which uses the RevokeRoleStmt syntax but is of type RevokeStmt.
		// It is documented at https://dev.mysql.com/doc/refman/5.7/en/revoke.html
		// as the "second syntax" for REVOKE. It is only valid if *both*
		// ALL PRIVILEGES + GRANT OPTION are specified in that order.
		if isRevokeAllGrant(elems) {
			var users []*ast.UserSpec
			for _, u := range usernames {
				users = append(users, &ast.UserSpec{
					User: u,
				})
			}
			return &ast.RevokeStmt{
				Privs:             []*ast.PrivElem{{Priv: mysql.AllPriv}, {Priv: mysql.GrantPriv}},
				ObjectType:        ast.ObjectTypeNone,
				Level:             &ast.GrantLevel{Level: ast.GrantLevelGlobal},
				Users:             users,
				IfExists:          ifExists,
				IgnoreUnknownUser: r.parseIgnoreUnknownUserOpt(),
			}
		}
		roles, err := convertToRole(elems)
		if err != nil {
			r.actionError(err)
		}
		return &ast.RevokeRoleStmt{
			Roles:             roles,
			Users:             usernames,
			IfExists:          ifExists,
			IgnoreUnknownUser: r.parseIgnoreUnknownUserOpt(),
		}
	default:
		r.syntaxError()
		return nil
	}
}

// parseIgnoreUnknownUserOpt implements the IGNORE UNKNOWN USER clause of
// the REVOKE statements: empty | "IGNORE" "UNKNOWN" "USER".
func (r *rdParser) parseIgnoreUnknownUserOpt() bool {
	if r.tok() != ignore {
		return false
	}
	r.advance()
	r.expect(unknown)
	r.expect(user)
	return true
}

// parseRoleOrPrivElemList implements RoleOrPrivElemList:
// RoleOrPrivElem | RoleOrPrivElemList ',' RoleOrPrivElem.
func (r *rdParser) parseRoleOrPrivElemList() []*ast.RoleOrPriv {
	list := []*ast.RoleOrPriv{r.parseRoleOrPrivElem()}
	for r.accept(int(',')) {
		list = append(list, r.parseRoleOrPrivElem())
	}
	return list
}

// parseRoleOrPrivElem implements:
//
//	RoleOrPrivElem: PrivElem | RolenameWithoutIdent | ExtendedPriv
//	              | "LOAD" "FROM" "S3" | "SELECT" "INTO" "S3"
//	RolenameWithoutIdent: stringLit | RolenameComposed
//	RolenameComposed: StringName '@' StringName
//	                | StringName singleAtIdentifier
//	ExtendedPriv: identifier | ExtendedPriv identifier
//
// A role name spelled with an identifier or unreserved keyword only
// matches here when composed with '@'/singleAtIdentifier (the LALR
// automaton reduces toward StringName only under that lookahead; a
// PrivType keyword or bare identifier reduces to PrivElem/ExtendedPriv
// otherwise, and a plain identifier role resolves later via
// RoleOrPriv.ToRole on its Symbols).
func (r *rdParser) parseRoleOrPrivElem() *ast.RoleOrPriv {
	switch {
	case r.tok() == load:
		// "LOAD" "FROM" "S3"
		r.advance()
		r.expect(from)
		r.expect(s3)
		return &ast.RoleOrPriv{
			Symbols: "LOAD FROM S3",
		}
	case r.tok() == selectKwd && r.la(1) == into:
		// "SELECT" "INTO" "S3"
		r.advance()
		r.advance()
		r.expect(s3)
		return &ast.RoleOrPriv{
			Symbols: "SELECT INTO S3",
		}
	}
	if (r.tok() == stringLit || isIdentifierTok(r.tok())) &&
		(r.la(1) == int('@') || r.la(1) == singleAtIdentifier) {
		// RolenameComposed
		username := r.parseStringName()
		var hostname string
		if r.tok() == singleAtIdentifier {
			hostname = strings.ToLower(strings.TrimPrefix(r.cur().lit, "@"))
			r.advance()
		} else {
			r.expect(int('@'))
			hostname = strings.ToLower(r.parseStringName())
		}
		return &ast.RoleOrPriv{
			Node: &auth.RoleIdentity{Username: username, Hostname: hostname},
		}
	}
	switch r.tok() {
	case stringLit:
		// RolenameWithoutIdent: stringLit
		username := r.cur().lit
		r.advance()
		return &ast.RoleOrPriv{
			Node: &auth.RoleIdentity{Username: username, Hostname: "%"},
		}
	case identifier:
		// ExtendedPriv
		symbols := []string{r.cur().lit}
		r.advance()
		for r.tok() == identifier {
			symbols = append(symbols, r.cur().lit)
			r.advance()
		}
		return &ast.RoleOrPriv{
			Symbols: strings.Join(symbols, " "),
		}
	default:
		return &ast.RoleOrPriv{
			Node: r.parsePrivElem(),
		}
	}
}

// parsePrivElem implements PrivElem:
// PrivType | PrivType '(' ColumnNameList ')'.
func (r *rdParser) parsePrivElem() *ast.PrivElem {
	priv := r.parsePrivType()
	if r.accept(int('(')) {
		cols := r.parseColumnNameList()
		r.expect(int(')'))
		return &ast.PrivElem{
			Priv: priv,
			Cols: cols,
		}
	}
	return &ast.PrivElem{
		Priv: priv,
	}
}

// parsePrivType implements PrivType.
func (r *rdParser) parsePrivType() mysql.PrivilegeType {
	switch r.tok() {
	case all:
		// "ALL" | "ALL" "PRIVILEGES"
		r.advance()
		r.accept(privileges)
		return mysql.AllPriv
	case alter:
		// "ALTER" | "ALTER" "ROUTINE"
		r.advance()
		if r.accept(routine) {
			return mysql.AlterRoutinePriv
		}
		return mysql.AlterPriv
	case create:
		// "CREATE" ["USER"|"TABLESPACE"|"TEMPORARY" "TABLES"|"VIEW"|"ROLE"|"ROUTINE"]
		r.advance()
		switch r.tok() {
		case user:
			r.advance()
			return mysql.CreateUserPriv
		case tablespace:
			r.advance()
			return mysql.CreateTablespacePriv
		case temporary:
			r.advance()
			r.expect(tables)
			return mysql.CreateTMPTablePriv
		case view:
			r.advance()
			return mysql.CreateViewPriv
		case role:
			r.advance()
			return mysql.CreateRolePriv
		case routine:
			r.advance()
			return mysql.CreateRoutinePriv
		}
		return mysql.CreatePriv
	case trigger:
		r.advance()
		return mysql.TriggerPriv
	case deleteKwd:
		r.advance()
		return mysql.DeletePriv
	case drop:
		// "DROP" | "DROP" "ROLE"
		r.advance()
		if r.accept(role) {
			return mysql.DropRolePriv
		}
		return mysql.DropPriv
	case process:
		r.advance()
		return mysql.ProcessPriv
	case execute:
		r.advance()
		return mysql.ExecutePriv
	case index:
		r.advance()
		return mysql.IndexPriv
	case insert:
		r.advance()
		return mysql.InsertPriv
	case selectKwd:
		r.advance()
		return mysql.SelectPriv
	case super:
		r.advance()
		return mysql.SuperPriv
	case show:
		// "SHOW" "DATABASES" | "SHOW" "VIEW"
		r.advance()
		switch r.tok() {
		case databases:
			r.advance()
			return mysql.ShowDBPriv
		case view:
			r.advance()
			return mysql.ShowViewPriv
		default:
			r.syntaxError()
		}
	case update:
		r.advance()
		return mysql.UpdatePriv
	case grant:
		// "GRANT" "OPTION"
		r.advance()
		r.expect(option)
		return mysql.GrantPriv
	case references:
		r.advance()
		return mysql.ReferencesPriv
	case replication:
		// "REPLICATION" "SLAVE" | "REPLICATION" "CLIENT"
		r.advance()
		switch r.tok() {
		case slave:
			r.advance()
			return mysql.ReplicationSlavePriv
		case client:
			r.advance()
			return mysql.ReplicationClientPriv
		default:
			r.syntaxError()
		}
	case binlog:
		// "BINLOG" "MONITOR" (MariaDB only)
		r.advance()
		r.expect(monitor)
		if r.p.enableMariaDB {
			return mysql.ReplicationClientPriv
		}
		r.actionError(ErrSyntax)
	case usage:
		r.advance()
		return mysql.UsagePriv
	case reload:
		r.advance()
		return mysql.ReloadPriv
	case file:
		r.advance()
		return mysql.FilePriv
	case config:
		r.advance()
		return mysql.ConfigPriv
	case lock:
		// "LOCK" "TABLES"
		r.advance()
		r.expect(tables)
		return mysql.LockTablesPriv
	case event:
		r.advance()
		return mysql.EventPriv
	case shutdown:
		r.advance()
		return mysql.ShutdownPriv
	}
	r.syntaxError()
	return 0
}

// parseObjectType implements ObjectType:
// empty %prec lowerThanFunction | "TABLE" | "FUNCTION" | "PROCEDURE".
// The %prec makes a following FUNCTION always shift into the object
// type rather than reduce the empty alternative, so consuming these
// keywords unconditionally matches the LR resolution.
func (r *rdParser) parseObjectType() ast.ObjectTypeType {
	switch r.tok() {
	case tableKwd:
		r.advance()
		return ast.ObjectTypeTable
	case function:
		r.advance()
		return ast.ObjectTypeFunction
	case procedure:
		r.advance()
		return ast.ObjectTypeProcedure
	}
	return ast.ObjectTypeNone
}

// parsePrivLevel implements PrivLevel:
// '*' | '*' '.' '*' | Identifier '.' '*' | Identifier '.' Identifier | Identifier.
func (r *rdParser) parsePrivLevel() *ast.GrantLevel {
	if r.accept(int('*')) {
		if r.accept(int('.')) {
			r.expect(int('*'))
			return &ast.GrantLevel{
				Level: ast.GrantLevelGlobal,
			}
		}
		// Note: a bare '*' is DB level in the grammar.
		return &ast.GrantLevel{
			Level: ast.GrantLevelDB,
		}
	}
	name := r.parseIdentifier()
	if r.accept(int('.')) {
		if r.accept(int('*')) {
			return &ast.GrantLevel{
				Level:  ast.GrantLevelDB,
				DBName: name,
			}
		}
		tableName := r.parseIdentifier()
		return &ast.GrantLevel{
			Level:     ast.GrantLevelTable,
			DBName:    name,
			TableName: tableName,
		}
	}
	return &ast.GrantLevel{
		Level:     ast.GrantLevelTable,
		TableName: name,
	}
}

// parseWithGrantOptionOpt implements WithGrantOptionOpt:
// empty | "WITH" "GRANT" "OPTION" | "WITH" MAX_*-limit NUM (ignored).
func (r *rdParser) parseWithGrantOptionOpt() bool {
	if r.tok() != with {
		return false
	}
	r.advance()
	switch r.tok() {
	case grant:
		r.advance()
		r.expect(option)
		return true
	case maxQueriesPerHour, maxUpdatesPerHour, maxConnectionsPerHour, maxUserConnections:
		r.advance()
		r.expect(intLit)
		return false
	default:
		r.syntaxError()
		return false
	}
}

// parseGrantUsername implements Username:
//
//	StringName | StringName '@' StringName
//	| StringName singleAtIdentifier | "CURRENT_USER" OptionalBraces
func (r *rdParser) parseGrantUsername() *auth.UserIdentity {
	if r.tok() == currentUser {
		// "CURRENT_USER" OptionalBraces
		r.advance()
		if r.accept(int('(')) {
			r.expect(int(')'))
		}
		return &auth.UserIdentity{CurrentUser: true}
	}
	username := r.parseStringName()
	switch r.tok() {
	case int('@'):
		r.advance()
		return &auth.UserIdentity{Username: username, Hostname: strings.ToLower(r.parseStringName())}
	case singleAtIdentifier:
		hostname := strings.ToLower(strings.TrimPrefix(r.cur().lit, "@"))
		r.advance()
		return &auth.UserIdentity{Username: username, Hostname: hostname}
	}
	return &auth.UserIdentity{Username: username, Hostname: "%"}
}

// parseGrantUsernameList implements UsernameList:
// Username | UsernameList ',' Username.
func (r *rdParser) parseGrantUsernameList() []*auth.UserIdentity {
	list := []*auth.UserIdentity{r.parseGrantUsername()}
	for r.accept(int(',')) {
		list = append(list, r.parseGrantUsername())
	}
	return list
}

// parseGrantUserSpec implements UserSpec: Username AuthOption.
func (r *rdParser) parseGrantUserSpec() *ast.UserSpec {
	userSpec := &ast.UserSpec{
		User: r.parseGrantUsername(),
	}
	if r.tok() == discard {
		// "DISCARD" "OLD" "PASSWORD" replaces the IDENTIFIED clause.
		r.advance()
		r.expect(old)
		r.expect(password)
		userSpec.AuthOpt = &ast.AuthOption{DiscardOldPassword: true}
		return userSpec
	}
	if opt := r.parseGrantAuthOption(); opt != nil {
		userSpec.AuthOpt = opt
		// Multi-factor authentication: "AND" IDENTIFIED ... up to twice.
		for r.tok() == and && r.la(1) == identified {
			r.advance()
			userSpec.MoreAuthOpts = append(userSpec.MoreAuthOpts, r.parseGrantAuthOption())
		}
	}
	return userSpec
}

// parseGrantUserSpecList implements UserSpecList:
// UserSpec | UserSpecList ',' UserSpec.
func (r *rdParser) parseGrantUserSpecList() []*ast.UserSpec {
	list := []*ast.UserSpec{r.parseGrantUserSpec()}
	for r.accept(int(',')) {
		list = append(list, r.parseGrantUserSpec())
	}
	return list
}

// parseGrantAuthOption implements AuthOption (nil for the empty
// alternative):
//
//	empty | "IDENTIFIED" "BY" (AuthString | "RANDOM" "PASSWORD")
//	| "IDENTIFIED" "WITH" AuthPlugin
//	  ["BY" (AuthString | "RANDOM" "PASSWORD") | "AS" HashString]
//	| "IDENTIFIED" "BY" "PASSWORD" HashString
//
// with AuthString: stringLit and AuthPlugin: StringName, followed by the
// optional dual-password clauses ["REPLACE" stringLit] ["RETAIN"
// "CURRENT" "PASSWORD"].
func (r *rdParser) parseGrantAuthOption() *ast.AuthOption {
	if r.tok() != identified {
		return nil
	}
	r.advance()
	var opt *ast.AuthOption
	switch r.tok() {
	case by:
		r.advance()
		switch {
		case r.tok() == random:
			// "BY" "RANDOM" "PASSWORD"
			r.advance()
			r.expect(password)
			opt = &ast.AuthOption{ByRandomPassword: true}
		case r.accept(password):
			// "IDENTIFIED" "BY" "PASSWORD" HashString
			opt = &ast.AuthOption{
				AuthPlugin:   mysql.AuthNativePassword,
				HashString:   r.parseGrantHashString(),
				ByHashString: true,
			}
		default:
			opt = &ast.AuthOption{
				AuthString:   r.expect(stringLit).lit,
				ByAuthString: true,
			}
		}
	case with:
		r.advance()
		plugin := r.parseStringName()
		switch r.tok() {
		case by:
			r.advance()
			if r.tok() == random {
				r.advance()
				r.expect(password)
				opt = &ast.AuthOption{AuthPlugin: plugin, ByRandomPassword: true}
			} else {
				opt = &ast.AuthOption{
					AuthPlugin:   plugin,
					AuthString:   r.expect(stringLit).lit,
					ByAuthString: true,
				}
			}
		case as:
			r.advance()
			opt = &ast.AuthOption{
				AuthPlugin:   plugin,
				HashString:   r.parseGrantHashString(),
				ByHashString: true,
			}
		default:
			opt = &ast.AuthOption{
				AuthPlugin: plugin,
			}
		}
	default:
		r.syntaxError()
		return nil
	}
	if r.tok() == replace && r.la(1) == stringLit {
		r.advance()
		s := r.expect(stringLit).lit
		opt.ReplacePassword = &s
	}
	if r.tok() == retain {
		r.advance()
		r.expect(current)
		r.expect(password)
		opt.RetainCurrentPassword = true
	}
	return opt
}

// parseGrantHashString implements HashString: stringLit | hexLit.
func (r *rdParser) parseGrantHashString() string {
	if r.tok() == stringLit {
		lit := r.cur().lit
		r.advance()
		return lit
	}
	t := r.expect(hexLit)
	return t.item.(ast.BinaryLiteral).ToString()
}

// parseGrantRequireClauseOpt implements RequireClauseOpt and
// RequireClause:
//
//	RequireClauseOpt: empty | RequireClause
//	RequireClause: "REQUIRE" "NONE" | "REQUIRE" "SSL" | "REQUIRE" "X509"
//	             | "REQUIRE" RequireList
//	RequireList: RequireListElement | RequireList "AND" RequireListElement
//	           | RequireList RequireListElement
func (r *rdParser) parseGrantRequireClauseOpt() []*ast.AuthTokenOrTLSOption {
	if r.tok() != require {
		return []*ast.AuthTokenOrTLSOption{}
	}
	r.advance()
	switch r.tok() {
	case none:
		r.advance()
		return []*ast.AuthTokenOrTLSOption{{Type: ast.TlsNone}}
	case ssl:
		r.advance()
		return []*ast.AuthTokenOrTLSOption{{Type: ast.Ssl}}
	case x509:
		r.advance()
		return []*ast.AuthTokenOrTLSOption{{Type: ast.X509}}
	}
	l := []*ast.AuthTokenOrTLSOption{r.parseGrantRequireListElement()}
	for {
		switch r.tok() {
		case and:
			r.advance()
			l = append(l, r.parseGrantRequireListElement())
		case issuer, subject, cipher, san, tokenIssuer:
			l = append(l, r.parseGrantRequireListElement())
		default:
			return l
		}
	}
}

// parseGrantRequireListElement implements RequireListElement:
// "ISSUER"|"SUBJECT"|"CIPHER"|"SAN"|"TOKEN_ISSUER" stringLit.
func (r *rdParser) parseGrantRequireListElement() *ast.AuthTokenOrTLSOption {
	var typ ast.AuthTokenOrTLSOptionType
	switch r.tok() {
	case issuer:
		typ = ast.Issuer
	case subject:
		typ = ast.Subject
	case cipher:
		typ = ast.Cipher
	case san:
		typ = ast.SAN
	case tokenIssuer:
		typ = ast.TokenIssuer
	default:
		r.syntaxError()
	}
	r.advance()
	return &ast.AuthTokenOrTLSOption{
		Type:  typ,
		Value: r.expect(stringLit).lit,
	}
}
