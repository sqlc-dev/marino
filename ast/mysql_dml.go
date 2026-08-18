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
	"github.com/sqlc-dev/marino/format"
)

// The MySQL data manipulation statements that postdate the goyacc
// grammar (MySQL 26.7 §15.2): HANDLER, IMPORT TABLE, and LOAD XML.

var (
	_ StmtNode = &HandlerOpenStmt{}
	_ StmtNode = &HandlerReadStmt{}
	_ StmtNode = &HandlerCloseStmt{}
	_ StmtNode = &ImportTableStmt{}
	_ StmtNode = &LoadXMLStmt{}
)

// HandlerOpenStmt is a HANDLER ... OPEN statement:
// HANDLER tbl_name OPEN [[AS] alias].
type HandlerOpenStmt struct {
	stmtNode

	Table *TableName
	Alias CIStr // empty when absent; restored with AS
}

// Restore implements Node interface.
func (n *HandlerOpenStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("HANDLER ")
	if err := n.Table.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore HandlerOpenStmt.Table")
	}
	ctx.WriteKeyWord(" OPEN")
	if n.Alias.O != "" {
		ctx.WriteKeyWord(" AS ")
		ctx.WriteName(n.Alias.O)
	}
	return nil
}

// Accept implements Node Accept interface.
func (n *HandlerOpenStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*HandlerOpenStmt)
	node, ok := n.Table.Accept(v)
	if !ok {
		return n, false
	}
	n.Table = node.(*TableName)
	return v.Leave(n)
}

// HandlerCompareOp is the comparison operator of the indexed-read form
// of HANDLER ... READ. HandlerOpNone selects the direction form.
type HandlerCompareOp int

const (
	// HandlerOpNone selects the FIRST/NEXT/PREV/LAST direction form.
	HandlerOpNone HandlerCompareOp = iota
	// HandlerOpEQ is =.
	HandlerOpEQ
	// HandlerOpLE is <=.
	HandlerOpLE
	// HandlerOpGE is >=.
	HandlerOpGE
	// HandlerOpLT is <.
	HandlerOpLT
	// HandlerOpGT is >.
	HandlerOpGT
)

// String implements fmt.Stringer interface.
func (n HandlerCompareOp) String() string {
	switch n {
	case HandlerOpEQ:
		return "="
	case HandlerOpLE:
		return "<="
	case HandlerOpGE:
		return ">="
	case HandlerOpLT:
		return "<"
	case HandlerOpGT:
		return ">"
	}
	return ""
}

// HandlerReadDirection is the cursor direction of HANDLER ... READ.
type HandlerReadDirection int

const (
	// HandlerReadFirst is FIRST.
	HandlerReadFirst HandlerReadDirection = iota
	// HandlerReadNext is NEXT.
	HandlerReadNext
	// HandlerReadPrev is PREV.
	HandlerReadPrev
	// HandlerReadLast is LAST.
	HandlerReadLast
)

// String implements fmt.Stringer interface.
func (n HandlerReadDirection) String() string {
	switch n {
	case HandlerReadFirst:
		return "FIRST"
	case HandlerReadNext:
		return "NEXT"
	case HandlerReadPrev:
		return "PREV"
	case HandlerReadLast:
		return "LAST"
	}
	return ""
}

// HandlerReadStmt is a HANDLER ... READ statement:
//
//	HANDLER tbl_name READ index_name { = | <= | >= | < | > } (value, ...)
//	    [WHERE where_condition] [LIMIT ...]
//	HANDLER tbl_name READ index_name { FIRST | NEXT | PREV | LAST }
//	    [WHERE where_condition] [LIMIT ...]
//	HANDLER tbl_name READ { FIRST | NEXT }
//	    [WHERE where_condition] [LIMIT ...]
//
// IndexName is empty for the table-scan form. Op is HandlerOpNone for
// the direction forms; otherwise Values holds the comparison values and
// Direction is unused.
type HandlerReadStmt struct {
	stmtNode

	Table     *TableName
	IndexName CIStr
	Op        HandlerCompareOp
	Values    []ExprNode
	Direction HandlerReadDirection
	Where     ExprNode
	Limit     *Limit
}

// Restore implements Node interface.
func (n *HandlerReadStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("HANDLER ")
	if err := n.Table.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore HandlerReadStmt.Table")
	}
	ctx.WriteKeyWord(" READ ")
	if n.IndexName.O != "" {
		ctx.WriteName(n.IndexName.O)
		ctx.WritePlain(" ")
	}
	if n.Op != HandlerOpNone {
		ctx.WritePlain(n.Op.String())
		ctx.WritePlain(" (")
		for i, value := range n.Values {
			if i != 0 {
				ctx.WritePlain(", ")
			}
			if err := value.Restore(ctx); err != nil {
				return annotatef(err, "An error occurred while restore HandlerReadStmt.Values[%d]", i)
			}
		}
		ctx.WritePlain(")")
	} else {
		ctx.WriteKeyWord(n.Direction.String())
	}
	if n.Where != nil {
		ctx.WriteKeyWord(" WHERE ")
		if err := n.Where.Restore(ctx); err != nil {
			return annotate(err, "An error occurred while restore HandlerReadStmt.Where")
		}
	}
	if n.Limit != nil {
		ctx.WritePlain(" ")
		if err := n.Limit.Restore(ctx); err != nil {
			return annotate(err, "An error occurred while restore HandlerReadStmt.Limit")
		}
	}
	return nil
}

// Accept implements Node Accept interface.
func (n *HandlerReadStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*HandlerReadStmt)
	node, ok := n.Table.Accept(v)
	if !ok {
		return n, false
	}
	n.Table = node.(*TableName)
	for i, value := range n.Values {
		node, ok := value.Accept(v)
		if !ok {
			return n, false
		}
		n.Values[i] = node.(ExprNode)
	}
	if n.Where != nil {
		node, ok := n.Where.Accept(v)
		if !ok {
			return n, false
		}
		n.Where = node.(ExprNode)
	}
	if n.Limit != nil {
		node, ok := n.Limit.Accept(v)
		if !ok {
			return n, false
		}
		n.Limit = node.(*Limit)
	}
	return v.Leave(n)
}

// HandlerCloseStmt is a HANDLER ... CLOSE statement:
// HANDLER tbl_name CLOSE.
type HandlerCloseStmt struct {
	stmtNode

	Table *TableName
}

// Restore implements Node interface.
func (n *HandlerCloseStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("HANDLER ")
	if err := n.Table.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore HandlerCloseStmt.Table")
	}
	ctx.WriteKeyWord(" CLOSE")
	return nil
}

// Accept implements Node Accept interface.
func (n *HandlerCloseStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*HandlerCloseStmt)
	node, ok := n.Table.Accept(v)
	if !ok {
		return n, false
	}
	n.Table = node.(*TableName)
	return v.Leave(n)
}

// ImportTableStmt is an IMPORT TABLE statement:
// IMPORT TABLE FROM sdi_file [, sdi_file] ...
// (ImportIntoStmt is the unrelated TiDB IMPORT INTO statement.)
type ImportTableStmt struct {
	stmtNode

	SdiFiles []string
}

// Restore implements Node interface.
func (n *ImportTableStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("IMPORT TABLE FROM ")
	for i, file := range n.SdiFiles {
		if i != 0 {
			ctx.WritePlain(", ")
		}
		ctx.WriteString(file)
	}
	return nil
}

// Accept implements Node Accept interface.
func (n *ImportTableStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*ImportTableStmt)
	return v.Leave(n)
}

// LoadXMLStmt is a LOAD XML statement:
//
//	LOAD XML [LOW_PRIORITY | CONCURRENT] [LOCAL] INFILE 'file_name'
//	    [REPLACE | IGNORE] INTO TABLE [db_name.]tbl_name
//	    [CHARACTER SET charset_name]
//	    [ROWS IDENTIFIED BY '<tagname>'] [IGNORE number {LINES | ROWS}]
//	    [(field_name_or_user_var [, field_name_or_user_var] ...)]
//	    [SET (col_name={expr | DEFAULT}) [, col_name={expr | DEFAULT}] ...]
//
// The IGNORE number LINES spelling parses like ROWS, so Restore()
// canonicalizes it to ROWS.
type LoadXMLStmt struct {
	stmtNode

	LowPriority        bool
	Concurrent         bool
	FileLocRef         FileLocRefTp
	Path               string
	OnDuplicate        OnDuplicateKeyHandlingType
	Table              *TableName
	Charset            *string
	RowsIdentifiedBy   string // the '<tagname>' literal; empty when absent
	IgnoreRows         *uint64
	ColumnsAndUserVars []*ColumnNameOrUserVar
	ColumnAssignments  []*Assignment
}

// Restore implements Node interface.
func (n *LoadXMLStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("LOAD XML ")
	if n.LowPriority {
		ctx.WriteKeyWord("LOW_PRIORITY ")
	} else if n.Concurrent {
		ctx.WriteKeyWord("CONCURRENT ")
	}
	if n.FileLocRef == FileLocClient {
		ctx.WriteKeyWord("LOCAL ")
	}
	ctx.WriteKeyWord("INFILE ")
	ctx.WriteString(n.Path)
	if n.OnDuplicate == OnDuplicateKeyHandlingReplace {
		ctx.WriteKeyWord(" REPLACE")
	} else if n.OnDuplicate == OnDuplicateKeyHandlingIgnore {
		ctx.WriteKeyWord(" IGNORE")
	}
	ctx.WriteKeyWord(" INTO TABLE ")
	if err := n.Table.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore LoadXMLStmt.Table")
	}
	if n.Charset != nil {
		ctx.WriteKeyWord(" CHARACTER SET ")
		ctx.WritePlain(*n.Charset)
	}
	if n.RowsIdentifiedBy != "" {
		ctx.WriteKeyWord(" ROWS IDENTIFIED BY ")
		ctx.WriteString(n.RowsIdentifiedBy)
	}
	if n.IgnoreRows != nil {
		ctx.WriteKeyWord(" IGNORE ")
		ctx.WritePlainf("%d", *n.IgnoreRows)
		ctx.WriteKeyWord(" ROWS")
	}
	if len(n.ColumnsAndUserVars) != 0 {
		ctx.WritePlain(" (")
		for i, c := range n.ColumnsAndUserVars {
			if i != 0 {
				ctx.WritePlain(",")
			}
			if err := c.Restore(ctx); err != nil {
				return annotatef(err, "An error occurred while restore LoadXMLStmt.ColumnsAndUserVars[%d]", i)
			}
		}
		ctx.WritePlain(")")
	}
	if len(n.ColumnAssignments) != 0 {
		ctx.WriteKeyWord(" SET")
		for i, assign := range n.ColumnAssignments {
			if i != 0 {
				ctx.WritePlain(",")
			}
			ctx.WritePlain(" ")
			if err := assign.Restore(ctx); err != nil {
				return annotatef(err, "An error occurred while restore LoadXMLStmt.ColumnAssignments[%d]", i)
			}
		}
	}
	return nil
}

// Accept implements Node Accept interface.
func (n *LoadXMLStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*LoadXMLStmt)
	if n.Table != nil {
		node, ok := n.Table.Accept(v)
		if !ok {
			return n, false
		}
		n.Table = node.(*TableName)
	}
	for i, cuVars := range n.ColumnsAndUserVars {
		node, ok := cuVars.Accept(v)
		if !ok {
			return n, false
		}
		n.ColumnsAndUserVars[i] = node.(*ColumnNameOrUserVar)
	}
	for i, assignment := range n.ColumnAssignments {
		node, ok := assignment.Accept(v)
		if !ok {
			return n, false
		}
		n.ColumnAssignments[i] = node.(*Assignment)
	}
	return v.Leave(n)
}
