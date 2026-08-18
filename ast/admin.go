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

package ast

import (
	"strings"

	"github.com/sqlc-dev/marino/auth"
	"github.com/sqlc-dev/marino/format"
)

// The MySQL database administration statements (MySQL 26.7 §15.7):
// table maintenance (CHECK/CHECKSUM/REPAIR TABLE), components, plugins,
// and loadable functions, CLONE, CACHE INDEX, LOAD INDEX INTO CACHE,
// and RESET PERSIST.

var (
	_ Node = &CacheTableIndex{}

	_ StmtNode = &CheckTableStmt{}
	_ StmtNode = &ChecksumTableStmt{}
	_ StmtNode = &RepairTablesStmt{}
	_ StmtNode = &CreateLoadableFunctionStmt{}
	_ StmtNode = &InstallComponentStmt{}
	_ StmtNode = &UninstallComponentStmt{}
	_ StmtNode = &InstallPluginStmt{}
	_ StmtNode = &UninstallPluginStmt{}
	_ StmtNode = &CloneStmt{}
	_ StmtNode = &CacheIndexStmt{}
	_ StmtNode = &LoadIndexStmt{}
	_ StmtNode = &ResetPersistStmt{}

	_ SensitiveStmtNode = &CloneStmt{}
)

// CheckTableOption is one check option of a CHECK TABLE statement.
type CheckTableOption int

const (
	// CheckTableForUpgrade is FOR UPGRADE.
	CheckTableForUpgrade CheckTableOption = iota
	// CheckTableQuick is QUICK.
	CheckTableQuick
	// CheckTableFast is FAST.
	CheckTableFast
	// CheckTableMedium is MEDIUM.
	CheckTableMedium
	// CheckTableExtended is EXTENDED.
	CheckTableExtended
	// CheckTableChanged is CHANGED.
	CheckTableChanged
)

// String implements fmt.Stringer interface.
func (n CheckTableOption) String() string {
	switch n {
	case CheckTableForUpgrade:
		return "FOR UPGRADE"
	case CheckTableQuick:
		return "QUICK"
	case CheckTableFast:
		return "FAST"
	case CheckTableMedium:
		return "MEDIUM"
	case CheckTableExtended:
		return "EXTENDED"
	case CheckTableChanged:
		return "CHANGED"
	}
	return ""
}

// CheckTableStmt is a CHECK TABLE statement:
// CHECK TABLE tbl_name [, tbl_name] ... [option] ...
// Options keeps the source order (MySQL accepts them in any order and
// repeated).
type CheckTableStmt struct {
	stmtNode

	Tables  []*TableName
	Options []CheckTableOption
}

// Restore implements Node interface.
func (n *CheckTableStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("CHECK TABLE ")
	for i, table := range n.Tables {
		if i != 0 {
			ctx.WritePlain(", ")
		}
		if err := table.Restore(ctx); err != nil {
			return annotatef(err, "An error occurred while restore CheckTableStmt.Tables[%d]", i)
		}
	}
	for _, opt := range n.Options {
		ctx.WritePlain(" ")
		ctx.WriteKeyWord(opt.String())
	}
	return nil
}

// Accept implements Node Accept interface.
func (n *CheckTableStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*CheckTableStmt)
	for i, table := range n.Tables {
		node, ok := table.Accept(v)
		if !ok {
			return n, false
		}
		n.Tables[i] = node.(*TableName)
	}
	return v.Leave(n)
}

// ChecksumType is the QUICK/EXTENDED modifier of a CHECKSUM TABLE
// statement.
type ChecksumType int

const (
	// ChecksumTypeDefault omits the modifier.
	ChecksumTypeDefault ChecksumType = iota
	// ChecksumTypeQuick is QUICK.
	ChecksumTypeQuick
	// ChecksumTypeExtended is EXTENDED.
	ChecksumTypeExtended
)

// ChecksumTableStmt is a CHECKSUM TABLE statement:
// CHECKSUM TABLE tbl_name [, tbl_name] ... [QUICK | EXTENDED].
type ChecksumTableStmt struct {
	stmtNode

	Tables []*TableName
	Type   ChecksumType
}

// Restore implements Node interface.
func (n *ChecksumTableStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("CHECKSUM TABLE ")
	for i, table := range n.Tables {
		if i != 0 {
			ctx.WritePlain(", ")
		}
		if err := table.Restore(ctx); err != nil {
			return annotatef(err, "An error occurred while restore ChecksumTableStmt.Tables[%d]", i)
		}
	}
	switch n.Type {
	case ChecksumTypeQuick:
		ctx.WriteKeyWord(" QUICK")
	case ChecksumTypeExtended:
		ctx.WriteKeyWord(" EXTENDED")
	}
	return nil
}

// Accept implements Node Accept interface.
func (n *ChecksumTableStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*ChecksumTableStmt)
	for i, table := range n.Tables {
		node, ok := table.Accept(v)
		if !ok {
			return n, false
		}
		n.Tables[i] = node.(*TableName)
	}
	return v.Leave(n)
}

// RepairTablesStmt is a REPAIR TABLE statement:
// REPAIR [NO_WRITE_TO_BINLOG | LOCAL] TABLE tbl_name [, tbl_name] ...
// [QUICK] [EXTENDED] [USE_FRM].
// (RepairTableStmt is the unrelated TiDB ADMIN REPAIR TABLE statement.)
type RepairTablesStmt struct {
	stmtNode

	NoWriteToBinLog bool
	Tables          []*TableName
	Quick           bool
	Extended        bool
	UseFrm          bool
}

// Restore implements Node interface.
func (n *RepairTablesStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("REPAIR ")
	if n.NoWriteToBinLog {
		ctx.WriteKeyWord("NO_WRITE_TO_BINLOG ")
	}
	ctx.WriteKeyWord("TABLE ")
	for i, table := range n.Tables {
		if i != 0 {
			ctx.WritePlain(", ")
		}
		if err := table.Restore(ctx); err != nil {
			return annotatef(err, "An error occurred while restore RepairTablesStmt.Tables[%d]", i)
		}
	}
	if n.Quick {
		ctx.WriteKeyWord(" QUICK")
	}
	if n.Extended {
		ctx.WriteKeyWord(" EXTENDED")
	}
	if n.UseFrm {
		ctx.WriteKeyWord(" USE_FRM")
	}
	return nil
}

// Accept implements Node Accept interface.
func (n *RepairTablesStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*RepairTablesStmt)
	for i, table := range n.Tables {
		node, ok := table.Accept(v)
		if !ok {
			return n, false
		}
		n.Tables[i] = node.(*TableName)
	}
	return v.Leave(n)
}

// LoadableFunctionReturnType is the RETURNS type of a loadable function.
type LoadableFunctionReturnType int

const (
	// LoadableFunctionReturnString is RETURNS STRING.
	LoadableFunctionReturnString LoadableFunctionReturnType = iota
	// LoadableFunctionReturnInteger is RETURNS INTEGER (or INT).
	LoadableFunctionReturnInteger
	// LoadableFunctionReturnReal is RETURNS REAL.
	LoadableFunctionReturnReal
	// LoadableFunctionReturnDecimal is RETURNS DECIMAL.
	LoadableFunctionReturnDecimal
)

// String implements fmt.Stringer interface.
func (n LoadableFunctionReturnType) String() string {
	switch n {
	case LoadableFunctionReturnString:
		return "STRING"
	case LoadableFunctionReturnInteger:
		return "INTEGER"
	case LoadableFunctionReturnReal:
		return "REAL"
	case LoadableFunctionReturnDecimal:
		return "DECIMAL"
	}
	return ""
}

// CreateLoadableFunctionStmt is a CREATE FUNCTION statement for loadable
// functions:
// CREATE [AGGREGATE] FUNCTION [IF NOT EXISTS] function_name
// RETURNS {STRING | INTEGER | REAL | DECIMAL} SONAME shared_library_name.
type CreateLoadableFunctionStmt struct {
	stmtNode

	IfNotExists  bool
	Aggregate    bool
	FunctionName CIStr
	ReturnType   LoadableFunctionReturnType
	SoName       string
}

// Restore implements Node interface.
func (n *CreateLoadableFunctionStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("CREATE ")
	if n.Aggregate {
		ctx.WriteKeyWord("AGGREGATE ")
	}
	ctx.WriteKeyWord("FUNCTION ")
	if n.IfNotExists {
		ctx.WriteKeyWord("IF NOT EXISTS ")
	}
	ctx.WriteName(n.FunctionName.O)
	ctx.WriteKeyWord(" RETURNS ")
	ctx.WriteKeyWord(n.ReturnType.String())
	ctx.WriteKeyWord(" SONAME ")
	ctx.WriteString(n.SoName)
	return nil
}

// Accept implements Node Accept interface.
func (n *CreateLoadableFunctionStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*CreateLoadableFunctionStmt)
	return v.Leave(n)
}

// InstallComponentStmt is an INSTALL COMPONENT statement:
// INSTALL COMPONENT component_name [, component_name] ...
// [SET variable = expr [, variable = expr] ...].
type InstallComponentStmt struct {
	stmtNode

	Components []string
	Variables  []*VariableAssignment
}

// Restore implements Node interface.
func (n *InstallComponentStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("INSTALL COMPONENT ")
	for i, component := range n.Components {
		if i != 0 {
			ctx.WritePlain(", ")
		}
		ctx.WriteString(component)
	}
	for i, v := range n.Variables {
		if i == 0 {
			ctx.WriteKeyWord(" SET ")
		} else {
			ctx.WritePlain(", ")
		}
		if err := v.Restore(ctx); err != nil {
			return annotatef(err, "An error occurred while restore InstallComponentStmt.Variables[%d]", i)
		}
	}
	return nil
}

// Accept implements Node Accept interface.
func (n *InstallComponentStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*InstallComponentStmt)
	for i, v2 := range n.Variables {
		node, ok := v2.Accept(v)
		if !ok {
			return n, false
		}
		n.Variables[i] = node.(*VariableAssignment)
	}
	return v.Leave(n)
}

// UninstallComponentStmt is an UNINSTALL COMPONENT statement:
// UNINSTALL COMPONENT component_name [, component_name] ...
type UninstallComponentStmt struct {
	stmtNode

	Components []string
}

// Restore implements Node interface.
func (n *UninstallComponentStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("UNINSTALL COMPONENT ")
	for i, component := range n.Components {
		if i != 0 {
			ctx.WritePlain(", ")
		}
		ctx.WriteString(component)
	}
	return nil
}

// Accept implements Node Accept interface.
func (n *UninstallComponentStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*UninstallComponentStmt)
	return v.Leave(n)
}

// InstallPluginStmt is an INSTALL PLUGIN statement:
// INSTALL PLUGIN plugin_name SONAME shared_library_name.
type InstallPluginStmt struct {
	stmtNode

	PluginName CIStr
	SoName     string
}

// Restore implements Node interface.
func (n *InstallPluginStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("INSTALL PLUGIN ")
	ctx.WriteName(n.PluginName.O)
	ctx.WriteKeyWord(" SONAME ")
	ctx.WriteString(n.SoName)
	return nil
}

// Accept implements Node Accept interface.
func (n *InstallPluginStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*InstallPluginStmt)
	return v.Leave(n)
}

// UninstallPluginStmt is an UNINSTALL PLUGIN statement:
// UNINSTALL PLUGIN plugin_name.
type UninstallPluginStmt struct {
	stmtNode

	PluginName CIStr
}

// Restore implements Node interface.
func (n *UninstallPluginStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("UNINSTALL PLUGIN ")
	ctx.WriteName(n.PluginName.O)
	return nil
}

// Accept implements Node Accept interface.
func (n *UninstallPluginStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*UninstallPluginStmt)
	return v.Leave(n)
}

// CloneSSLType is the [REQUIRE [NO] SSL] clause of CLONE INSTANCE.
type CloneSSLType int

const (
	// CloneSSLDefault omits the clause.
	CloneSSLDefault CloneSSLType = iota
	// CloneSSLRequire is REQUIRE SSL.
	CloneSSLRequire
	// CloneSSLRequireNo is REQUIRE NO SSL.
	CloneSSLRequireNo
)

// CloneStmt is a CLONE statement:
//
//	CLONE LOCAL DATA DIRECTORY [=] 'clone_dir'
//	CLONE INSTANCE FROM 'user'@'host':port IDENTIFIED BY 'password'
//	    [DATA DIRECTORY [=] 'clone_dir'] [REQUIRE [NO] SSL]
//
// Local selects the first form; the remote form sets User, Port, and
// Password. DataDirectory is empty when the optional clause is absent
// from the remote form.
type CloneStmt struct {
	stmtNode

	Local         bool
	DataDirectory string
	User          *auth.UserIdentity
	Port          uint64
	Password      string
	SSL           CloneSSLType
}

// Restore implements Node interface.
func (n *CloneStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("CLONE ")
	if n.Local {
		ctx.WriteKeyWord("LOCAL DATA DIRECTORY ")
		ctx.WritePlain("= ")
		ctx.WriteString(n.DataDirectory)
		return nil
	}
	ctx.WriteKeyWord("INSTANCE FROM ")
	if err := n.User.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore CloneStmt.User")
	}
	ctx.WritePlainf(":%d", n.Port)
	ctx.WriteKeyWord(" IDENTIFIED BY ")
	ctx.WriteString(n.Password)
	if n.DataDirectory != "" {
		ctx.WriteKeyWord(" DATA DIRECTORY ")
		ctx.WritePlain("= ")
		ctx.WriteString(n.DataDirectory)
	}
	switch n.SSL {
	case CloneSSLRequire:
		ctx.WriteKeyWord(" REQUIRE SSL")
	case CloneSSLRequireNo:
		ctx.WriteKeyWord(" REQUIRE NO SSL")
	}
	return nil
}

// SecureText implements SensitiveStatement interface.
func (n *CloneStmt) SecureText() string {
	masked := *n
	if !masked.Local {
		masked.Password = "xxxxxx"
	}
	var sb strings.Builder
	_ = masked.Restore(format.NewRestoreCtx(format.DefaultRestoreFlags, &sb))
	return sb.String()
}

// Accept implements Node Accept interface.
func (n *CloneStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*CloneStmt)
	return v.Leave(n)
}

// CacheTableIndex is one table specification of a CACHE INDEX or LOAD
// INDEX INTO CACHE statement:
// tbl_name [PARTITION (partition_list)] [{INDEX | KEY} (index_name, ...)]
// [IGNORE LEAVES]. AllPartitions is PARTITION (ALL); IgnoreLeaves is only
// produced by LOAD INDEX INTO CACHE.
type CacheTableIndex struct {
	node

	Table         *TableName
	Partitions    []CIStr
	AllPartitions bool
	Indexes       []CIStr
	IgnoreLeaves  bool
}

// Restore implements Node interface.
func (n *CacheTableIndex) Restore(ctx *format.RestoreCtx) error {
	if err := n.Table.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore CacheTableIndex.Table")
	}
	if n.AllPartitions {
		ctx.WriteKeyWord(" PARTITION ")
		ctx.WritePlain("(")
		ctx.WriteKeyWord("ALL")
		ctx.WritePlain(")")
	} else if len(n.Partitions) > 0 {
		ctx.WriteKeyWord(" PARTITION ")
		ctx.WritePlain("(")
		for i, partition := range n.Partitions {
			if i != 0 {
				ctx.WritePlain(", ")
			}
			ctx.WriteName(partition.O)
		}
		ctx.WritePlain(")")
	}
	if len(n.Indexes) > 0 {
		ctx.WriteKeyWord(" INDEX ")
		ctx.WritePlain("(")
		for i, idx := range n.Indexes {
			if i != 0 {
				ctx.WritePlain(", ")
			}
			ctx.WriteName(idx.O)
		}
		ctx.WritePlain(")")
	}
	if n.IgnoreLeaves {
		ctx.WriteKeyWord(" IGNORE LEAVES")
	}
	return nil
}

// Accept implements Node Accept interface.
func (n *CacheTableIndex) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*CacheTableIndex)
	node, ok := n.Table.Accept(v)
	if !ok {
		return n, false
	}
	n.Table = node.(*TableName)
	return v.Leave(n)
}

// CacheIndexStmt is a CACHE INDEX statement:
//
//	CACHE INDEX {tbl_index_list [, tbl_index_list] ...
//	    | tbl_name PARTITION (partition_list)} IN key_cache_name
type CacheIndexStmt struct {
	stmtNode

	TableIndexes []*CacheTableIndex
	KeyCacheName CIStr
}

// Restore implements Node interface.
func (n *CacheIndexStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("CACHE INDEX ")
	for i, ti := range n.TableIndexes {
		if i != 0 {
			ctx.WritePlain(", ")
		}
		if err := ti.Restore(ctx); err != nil {
			return annotatef(err, "An error occurred while restore CacheIndexStmt.TableIndexes[%d]", i)
		}
	}
	ctx.WriteKeyWord(" IN ")
	ctx.WriteName(n.KeyCacheName.O)
	return nil
}

// Accept implements Node Accept interface.
func (n *CacheIndexStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*CacheIndexStmt)
	for i, ti := range n.TableIndexes {
		node, ok := ti.Accept(v)
		if !ok {
			return n, false
		}
		n.TableIndexes[i] = node.(*CacheTableIndex)
	}
	return v.Leave(n)
}

// LoadIndexStmt is a LOAD INDEX INTO CACHE statement:
//
//	LOAD INDEX INTO CACHE
//	    tbl_index_or_partition [, tbl_index_or_partition] ...
type LoadIndexStmt struct {
	stmtNode

	TableIndexes []*CacheTableIndex
}

// Restore implements Node interface.
func (n *LoadIndexStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("LOAD INDEX INTO CACHE ")
	for i, ti := range n.TableIndexes {
		if i != 0 {
			ctx.WritePlain(", ")
		}
		if err := ti.Restore(ctx); err != nil {
			return annotatef(err, "An error occurred while restore LoadIndexStmt.TableIndexes[%d]", i)
		}
	}
	return nil
}

// Accept implements Node Accept interface.
func (n *LoadIndexStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*LoadIndexStmt)
	for i, ti := range n.TableIndexes {
		node, ok := ti.Accept(v)
		if !ok {
			return n, false
		}
		n.TableIndexes[i] = node.(*CacheTableIndex)
	}
	return v.Leave(n)
}

// ResetPersistStmt is a RESET PERSIST statement:
// RESET PERSIST [[IF EXISTS] system_var_name].
// Variable is empty when every persisted variable is removed; IfExists
// requires Variable.
type ResetPersistStmt struct {
	stmtNode

	IfExists bool
	Variable string
}

// Restore implements Node interface.
func (n *ResetPersistStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("RESET PERSIST")
	if n.Variable != "" {
		ctx.WritePlain(" ")
		if n.IfExists {
			ctx.WriteKeyWord("IF EXISTS ")
		}
		ctx.WriteName(n.Variable)
	}
	return nil
}

// Accept implements Node Accept interface.
func (n *ResetPersistStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*ResetPersistStmt)
	return v.Leave(n)
}
