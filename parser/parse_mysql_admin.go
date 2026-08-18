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

// The MySQL database administration statements (MySQL 26.7 §15.7):
// table maintenance (CHECK/CHECKSUM/REPAIR TABLE), INSTALL/UNINSTALL
// COMPONENT and PLUGIN, CREATE FUNCTION for loadable functions, CLONE,
// CACHE INDEX, LOAD INDEX INTO CACHE, and RESET PERSIST. These postdate
// the goyacc grammar; the productions below are written from the MySQL
// 26.7 reference manual in the same style as parser.y. (The MySQL-syntax
// resource group options live with the other resource group productions
// in parse_create_misc.go.)

import (
	"github.com/sqlc-dev/marino/ast"
)

func init() {
	rdRegister(check, (*rdParser).parseCheckTableStmt)
	rdRegister(checksum, (*rdParser).parseChecksumTableStmt)
	rdRegister(repair, (*rdParser).parseRepairTablesStmt)
	rdRegister(install, (*rdParser).parseInstallStmt)
	rdRegister(uninstall, (*rdParser).parseUninstallStmt)
	rdRegister(clone, (*rdParser).parseCloneStmt)
	rdRegister(cache, (*rdParser).parseCacheIndexStmt)
	rdRegister(reset, (*rdParser).parseResetStmt)
}

// parseCheckTableStmt implements CheckTableStmt:
//
//	"CHECK" "TABLE" TableNameList CheckTableOptionListOpt
//	CheckTableOption: "FOR" "UPGRADE" | "QUICK" | "FAST" | "MEDIUM"
//	    | "EXTENDED" | "CHANGED"
func (r *rdParser) parseCheckTableStmt() ast.StmtNode {
	r.expect(check)
	r.expect(tableKwd)
	stmt := &ast.CheckTableStmt{Tables: r.parseTableNameList()}
	for {
		switch r.tok() {
		case forKwd:
			r.advance()
			r.expect(upgrade)
			stmt.Options = append(stmt.Options, ast.CheckTableForUpgrade)
		case quick:
			r.advance()
			stmt.Options = append(stmt.Options, ast.CheckTableQuick)
		case fast:
			r.advance()
			stmt.Options = append(stmt.Options, ast.CheckTableFast)
		case medium:
			r.advance()
			stmt.Options = append(stmt.Options, ast.CheckTableMedium)
		case extended:
			r.advance()
			stmt.Options = append(stmt.Options, ast.CheckTableExtended)
		case changed:
			r.advance()
			stmt.Options = append(stmt.Options, ast.CheckTableChanged)
		default:
			return stmt
		}
	}
}

// parseChecksumTableStmt implements ChecksumTableStmt:
// "CHECKSUM" "TABLE" TableNameList ["QUICK" | "EXTENDED"].
func (r *rdParser) parseChecksumTableStmt() ast.StmtNode {
	r.expect(checksum)
	r.expect(tableKwd)
	stmt := &ast.ChecksumTableStmt{Tables: r.parseTableNameList()}
	switch r.tok() {
	case quick:
		r.advance()
		stmt.Type = ast.ChecksumTypeQuick
	case extended:
		r.advance()
		stmt.Type = ast.ChecksumTypeExtended
	}
	return stmt
}

// parseRepairTablesStmt implements RepairTablesStmt:
//
//	"REPAIR" NoWriteToBinLogAliasOpt "TABLE" TableNameList
//	    RepairTableOptionListOpt
//	RepairTableOption: "QUICK" | "EXTENDED" | "USE_FRM"
//
// The options may come in any order and repeat, as in MySQL.
func (r *rdParser) parseRepairTablesStmt() ast.StmtNode {
	r.expect(repair)
	noWrite := false
	if r.tok() == noWriteToBinLog || r.tok() == local {
		noWrite = true
		r.advance()
	}
	r.expect(tableKwd)
	stmt := &ast.RepairTablesStmt{
		NoWriteToBinLog: noWrite,
		Tables:          r.parseTableNameList(),
	}
	for {
		switch r.tok() {
		case quick:
			r.advance()
			stmt.Quick = true
		case extended:
			r.advance()
			stmt.Extended = true
		case useFrm:
			r.advance()
			stmt.UseFrm = true
		default:
			return stmt
		}
	}
}

// finishCreateLoadableFunctionStmt implements the rest of
// CreateLoadableFunctionStmt (the CREATE FUNCTION statement for
// loadable functions) after parseCreateFunctionFamily has parsed
// through the function name:
//
//	"CREATE" ["AGGREGATE"] "FUNCTION" IfNotExists Identifier
//	    "RETURNS" ("STRING" | "INTEGER" | "INT" | "REAL" | "DECIMAL")
//	    "SONAME" stringLit
func (r *rdParser) finishCreateLoadableFunctionStmt(aggregateOpt, ifNotExists bool, name ast.CIStr) ast.StmtNode {
	r.expect(returns)
	var returnType ast.LoadableFunctionReturnType
	switch r.tok() {
	case stringKwd:
		returnType = ast.LoadableFunctionReturnString
	case intType, integerType:
		returnType = ast.LoadableFunctionReturnInteger
	case realType:
		returnType = ast.LoadableFunctionReturnReal
	case decimalType:
		returnType = ast.LoadableFunctionReturnDecimal
	default:
		r.syntaxError()
	}
	r.advance()
	r.expect(soname)
	return &ast.CreateLoadableFunctionStmt{
		IfNotExists:  ifNotExists,
		Aggregate:    aggregateOpt,
		FunctionName: name,
		ReturnType:   returnType,
		SoName:       r.expect(stringLit).lit,
	}
}

// parseInstallStmt implements InstallComponentStmt and InstallPluginStmt:
//
//	"INSTALL" "COMPONENT" ComponentNameList ["SET" VariableAssignmentList]
//	"INSTALL" "PLUGIN" Identifier "SONAME" stringLit
func (r *rdParser) parseInstallStmt() ast.StmtNode {
	r.expect(install)
	switch r.tok() {
	case component:
		r.advance()
		stmt := &ast.InstallComponentStmt{Components: r.parseComponentNameList()}
		if r.accept(set) {
			stmt.Variables = r.parseVariableAssignmentList()
		}
		return stmt
	case plugin:
		r.advance()
		name := r.parseIdentifier()
		r.expect(soname)
		return &ast.InstallPluginStmt{
			PluginName: ast.NewCIStr(name),
			SoName:     r.expect(stringLit).lit,
		}
	}
	r.syntaxError()
	return nil
}

// parseUninstallStmt implements UninstallComponentStmt and
// UninstallPluginStmt:
//
//	"UNINSTALL" "COMPONENT" ComponentNameList
//	"UNINSTALL" "PLUGIN" Identifier
func (r *rdParser) parseUninstallStmt() ast.StmtNode {
	r.expect(uninstall)
	switch r.tok() {
	case component:
		r.advance()
		return &ast.UninstallComponentStmt{Components: r.parseComponentNameList()}
	case plugin:
		r.advance()
		return &ast.UninstallPluginStmt{PluginName: ast.NewCIStr(r.parseIdentifier())}
	}
	r.syntaxError()
	return nil
}

// parseComponentNameList implements ComponentNameList:
// stringLit | ComponentNameList ',' stringLit.
func (r *rdParser) parseComponentNameList() []string {
	list := []string{r.expect(stringLit).lit}
	for r.accept(int(',')) {
		list = append(list, r.expect(stringLit).lit)
	}
	return list
}

// parseCloneStmt implements CloneStmt:
//
//	"CLONE" "LOCAL" "DATA" "DIRECTORY" EqOpt stringLit
//	"CLONE" "INSTANCE" "FROM" Username ':' NUM "IDENTIFIED" "BY" stringLit
//	    ["DATA" "DIRECTORY" EqOpt stringLit] ["REQUIRE" ["NO"] "SSL"]
func (r *rdParser) parseCloneStmt() ast.StmtNode {
	r.expect(clone)
	switch r.tok() {
	case local:
		r.advance()
		r.expect(data)
		r.expect(directory)
		r.parseEqOpt()
		return &ast.CloneStmt{
			Local:         true,
			DataDirectory: r.expect(stringLit).lit,
		}
	case instance:
		r.advance()
		r.expect(from)
		stmt := &ast.CloneStmt{User: r.parseUsername()}
		r.expect(int(':'))
		stmt.Port = getUint64FromNUM(r.expect(intLit).item)
		r.expect(identified)
		r.expect(by)
		stmt.Password = r.expect(stringLit).lit
		if r.accept(data) {
			r.expect(directory)
			r.parseEqOpt()
			stmt.DataDirectory = r.expect(stringLit).lit
		}
		if r.accept(require) {
			if r.accept(no) {
				stmt.SSL = ast.CloneSSLRequireNo
			} else {
				stmt.SSL = ast.CloneSSLRequire
			}
			r.expect(ssl)
		}
		return stmt
	}
	r.syntaxError()
	return nil
}

// parseCacheIndexStmt implements CacheIndexStmt:
// "CACHE" "INDEX" CacheTableIndexList "IN" KeyCacheName.
func (r *rdParser) parseCacheIndexStmt() ast.StmtNode {
	r.expect(cache)
	r.expect(index)
	stmt := &ast.CacheIndexStmt{
		TableIndexes: r.parseCacheTableIndexList(false),
	}
	r.expect(in)
	// KeyCacheName: Identifier | "DEFAULT" (the default key cache).
	if r.tok() == defaultKwd {
		stmt.KeyCacheName = ast.NewCIStr(r.cur().lit)
		r.advance()
	} else {
		stmt.KeyCacheName = ast.NewCIStr(r.parseIdentifier())
	}
	return stmt
}

// parseLoadIndexStmt implements LoadIndexStmt:
// "LOAD" "INDEX" "INTO" "CACHE" CacheTableIndexList.
func (r *rdParser) parseLoadIndexStmt() ast.StmtNode {
	r.expect(load)
	r.expect(index)
	r.expect(into)
	r.expect(cache)
	return &ast.LoadIndexStmt{TableIndexes: r.parseCacheTableIndexList(true)}
}

// parseCacheTableIndexList implements CacheTableIndexList:
// CacheTableIndex | CacheTableIndexList ',' CacheTableIndex.
func (r *rdParser) parseCacheTableIndexList(allowIgnoreLeaves bool) []*ast.CacheTableIndex {
	list := []*ast.CacheTableIndex{r.parseCacheTableIndex(allowIgnoreLeaves)}
	for r.accept(int(',')) {
		list = append(list, r.parseCacheTableIndex(allowIgnoreLeaves))
	}
	return list
}

// parseCacheTableIndex implements CacheTableIndex:
//
//	TableName ["PARTITION" '(' ("ALL" | PartitionNameList) ')']
//	    [("INDEX" | "KEY") '(' IndexNameList ')'] ["IGNORE" "LEAVES"]
//
// IGNORE LEAVES belongs to LOAD INDEX INTO CACHE only. An index name may
// also be "PRIMARY", naming the primary key.
func (r *rdParser) parseCacheTableIndex(allowIgnoreLeaves bool) *ast.CacheTableIndex {
	ti := &ast.CacheTableIndex{Table: r.parseTableName()}
	if r.tok() == partition {
		r.advance()
		r.expect(int('('))
		if r.tok() == all {
			r.advance()
			ti.AllPartitions = true
		} else {
			ti.Partitions = []ast.CIStr{ast.NewCIStr(r.parseIdentifier())}
			for r.accept(int(',')) {
				ti.Partitions = append(ti.Partitions, ast.NewCIStr(r.parseIdentifier()))
			}
		}
		r.expect(int(')'))
	}
	if r.tok() == index || r.tok() == key {
		r.advance()
		r.expect(int('('))
		ti.Indexes = []ast.CIStr{r.parseCacheIndexName()}
		for r.accept(int(',')) {
			ti.Indexes = append(ti.Indexes, r.parseCacheIndexName())
		}
		r.expect(int(')'))
	}
	if allowIgnoreLeaves && r.tok() == ignore {
		r.advance()
		r.expect(leaves)
		ti.IgnoreLeaves = true
	}
	return ti
}

// parseCacheIndexName implements the key_usage_element production:
// Identifier | "PRIMARY".
func (r *rdParser) parseCacheIndexName() ast.CIStr {
	if r.tok() == primary {
		name := ast.NewCIStr(r.cur().lit)
		r.advance()
		return name
	}
	return ast.NewCIStr(r.parseIdentifier())
}

// parseResetStmt implements the RESET statement family:
// ResetPersistStmt ("RESET" "PERSIST" [["IF" "EXISTS"] Identifier])
// here, with ResetReplicaStmt and ResetBinaryLogsAndGtidsStmt in
// parse_replication.go.
func (r *rdParser) parseResetStmt() ast.StmtNode {
	switch r.la(1) {
	case replica:
		return r.parseResetReplicaStmt()
	case binaryType:
		return r.parseResetBinaryLogsAndGtidsStmt()
	case persist:
	default:
		// The other RESET forms do not parse; fail at RESET the way an
		// unregistered statement head does.
		r.syntaxError()
	}
	r.expect(reset)
	r.expect(persist)
	stmt := &ast.ResetPersistStmt{}
	switch {
	case r.accept(ifKwd):
		r.expect(exists)
		stmt.IfExists = true
		stmt.Variable = r.parseIdentifier()
	case isIdentifierTok(r.tok()):
		stmt.Variable = r.parseIdentifier()
	}
	return stmt
}
