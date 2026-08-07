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

// RENAME statements: RenameTableStmt (TableToTableList) and
// RenameUserStmt (UserToUserList).

import (
	"github.com/sqlc-dev/marino/ast"
)

func init() {
	rdRegister(rename, (*rdParser).parseRenameStmtFamily)
}

// parseRenameStmtFamily implements:
//
//	RenameTableStmt: "RENAME" "TABLE" TableToTableList
//	RenameUserStmt: "RENAME" "USER" UserToUserList
func (r *rdParser) parseRenameStmtFamily() ast.StmtNode {
	r.expect(rename)
	switch r.tok() {
	case tableKwd:
		r.advance()
		// TableToTableList: TableToTable | TableToTableList ',' TableToTable
		list := []*ast.TableToTable{r.parseTableToTable()}
		for r.accept(int(',')) {
			list = append(list, r.parseTableToTable())
		}
		return &ast.RenameTableStmt{
			TableToTables: list,
		}
	case user:
		r.advance()
		// UserToUserList: UserToUser | UserToUserList ',' UserToUser
		list := []*ast.UserToUser{r.parseUserToUser()}
		for r.accept(int(',')) {
			list = append(list, r.parseUserToUser())
		}
		return &ast.RenameUserStmt{
			UserToUsers: list,
		}
	}
	r.syntaxError()
	return nil
}

// parseTableToTable implements TableToTable: TableName "TO" TableName.
func (r *rdParser) parseTableToTable() *ast.TableToTable {
	oldTable := r.parseTableName()
	r.expect(to)
	return &ast.TableToTable{
		OldTable: oldTable,
		NewTable: r.parseTableName(),
	}
}

// parseUserToUser implements UserToUser: Username "TO" Username.
func (r *rdParser) parseUserToUser() *ast.UserToUser {
	oldUser := r.parseUsername()
	r.expect(to)
	return &ast.UserToUser{
		OldUser: oldUser,
		NewUser: r.parseUsername(),
	}
}
