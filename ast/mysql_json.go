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

// The MySQL JSON table function and the JSON_VALUE extraction clauses
// (MySQL 26.7 §14.17.6, §14.17.3): JSON_TABLE(expr, path COLUMNS (...))
// as a table factor, and JSON_VALUE(doc, path RETURNING type ON
// EMPTY/ON ERROR).

import (
	"io"
	"strings"

	"github.com/sqlc-dev/marino/format"
	"github.com/sqlc-dev/marino/types"
)

var (
	_ Node          = &JSONOnResponse{}
	_ Node          = &JSONTableColumn{}
	_ ResultSetNode = &JSONTableExpr{}
	_ ExprNode      = &JSONValueExpr{}
)

// JSONOnResponseType is what a JSON function produces for an ON EMPTY or
// ON ERROR condition.
type JSONOnResponseType int

// JSON on-response alternatives.
const (
	// JSONOnResponseNull produces SQL NULL (NULL ON EMPTY/ERROR).
	JSONOnResponseNull JSONOnResponseType = iota
	// JSONOnResponseError raises an error (ERROR ON EMPTY/ERROR).
	JSONOnResponseError
	// JSONOnResponseDefault produces the given value (DEFAULT v ON ...).
	JSONOnResponseDefault
)

// JSONOnResponse is one {NULL | ERROR | DEFAULT value} ON {EMPTY|ERROR}
// clause of JSON_VALUE or a JSON_TABLE column. Whether it answers ON
// EMPTY or ON ERROR is determined by the field holding it.
type JSONOnResponse struct {
	node

	Tp JSONOnResponseType
	// Value is the DEFAULT value; nil for NULL/ERROR.
	Value ExprNode
}

// restoreWithCondition writes the clause followed by " ON <cond>".
func (n *JSONOnResponse) restoreWithCondition(ctx *format.RestoreCtx, cond string) error {
	switch n.Tp {
	case JSONOnResponseNull:
		ctx.WriteKeyWord("NULL")
	case JSONOnResponseError:
		ctx.WriteKeyWord("ERROR")
	case JSONOnResponseDefault:
		ctx.WriteKeyWord("DEFAULT ")
		if err := n.Value.Restore(ctx); err != nil {
			return annotate(err, "An error occurred while restore JSONOnResponse.Value")
		}
	}
	ctx.WriteKeyWord(" ON ")
	ctx.WriteKeyWord(cond)
	return nil
}

// Restore implements Node interface. The ON condition spelling is owned
// by the enclosing node, which restores through restoreWithCondition;
// Restore alone writes the response value only.
func (n *JSONOnResponse) Restore(ctx *format.RestoreCtx) error {
	switch n.Tp {
	case JSONOnResponseNull:
		ctx.WriteKeyWord("NULL")
	case JSONOnResponseError:
		ctx.WriteKeyWord("ERROR")
	case JSONOnResponseDefault:
		ctx.WriteKeyWord("DEFAULT ")
		if err := n.Value.Restore(ctx); err != nil {
			return annotate(err, "An error occurred while restore JSONOnResponse.Value")
		}
	}
	return nil
}

// Accept implements Node Accept interface.
func (n *JSONOnResponse) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*JSONOnResponse)
	if n.Value != nil {
		node, ok := n.Value.Accept(v)
		if !ok {
			return n, false
		}
		n.Value = node.(ExprNode)
	}
	return v.Leave(n)
}

// JSONTableColumn is one column specification of JSON_TABLE's COLUMNS
// clause:
//
//	name FOR ORDINALITY
//	| name type [EXISTS] PATH 'path' [on_empty] [on_error]
//	| NESTED [PATH] 'path' COLUMNS (column[, column]...)
type JSONTableColumn struct {
	node

	// ForOrdinality selects the counter form; only Name is set.
	ForOrdinality bool
	Name          CIStr
	Tp            *types.FieldType
	// Exists selects the EXISTS PATH form.
	Exists bool
	// Path is the column's JSON path (a string literal expression); nil
	// for the NESTED form, which uses NestedPath.
	Path    ExprNode
	OnEmpty *JSONOnResponse
	OnError *JSONOnResponse

	// Nested selects the NESTED PATH form: NestedPath and NestedColumns
	// are set, everything else is zero.
	Nested        bool
	NestedPath    ExprNode
	NestedColumns []*JSONTableColumn
}

// Restore implements Node interface.
func (n *JSONTableColumn) Restore(ctx *format.RestoreCtx) error {
	if n.Nested {
		ctx.WriteKeyWord("NESTED PATH ")
		if err := n.NestedPath.Restore(ctx); err != nil {
			return annotate(err, "An error occurred while restore JSONTableColumn.NestedPath")
		}
		ctx.WriteKeyWord(" COLUMNS ")
		ctx.WritePlain("(")
		for i, col := range n.NestedColumns {
			if i != 0 {
				ctx.WritePlain(", ")
			}
			if err := col.Restore(ctx); err != nil {
				return annotatef(err, "An error occurred while restore JSONTableColumn.NestedColumns[%d]", i)
			}
		}
		ctx.WritePlain(")")
		return nil
	}
	ctx.WriteName(n.Name.O)
	if n.ForOrdinality {
		ctx.WriteKeyWord(" FOR ORDINALITY")
		return nil
	}
	ctx.WritePlain(" ")
	if err := n.Tp.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore JSONTableColumn.Tp")
	}
	if n.Exists {
		ctx.WriteKeyWord(" EXISTS")
	}
	ctx.WriteKeyWord(" PATH ")
	if err := n.Path.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore JSONTableColumn.Path")
	}
	if n.OnEmpty != nil {
		ctx.WritePlain(" ")
		if err := n.OnEmpty.restoreWithCondition(ctx, "EMPTY"); err != nil {
			return err
		}
	}
	if n.OnError != nil {
		ctx.WritePlain(" ")
		if err := n.OnError.restoreWithCondition(ctx, "ERROR"); err != nil {
			return err
		}
	}
	return nil
}

// Accept implements Node Accept interface.
func (n *JSONTableColumn) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*JSONTableColumn)
	if n.Path != nil {
		node, ok := n.Path.Accept(v)
		if !ok {
			return n, false
		}
		n.Path = node.(ExprNode)
	}
	if n.OnEmpty != nil {
		node, ok := n.OnEmpty.Accept(v)
		if !ok {
			return n, false
		}
		n.OnEmpty = node.(*JSONOnResponse)
	}
	if n.OnError != nil {
		node, ok := n.OnError.Accept(v)
		if !ok {
			return n, false
		}
		n.OnError = node.(*JSONOnResponse)
	}
	if n.NestedPath != nil {
		node, ok := n.NestedPath.Accept(v)
		if !ok {
			return n, false
		}
		n.NestedPath = node.(ExprNode)
	}
	for i, col := range n.NestedColumns {
		node, ok := col.Accept(v)
		if !ok {
			return n, false
		}
		n.NestedColumns[i] = node.(*JSONTableColumn)
	}
	return v.Leave(n)
}

// JSONTableExpr is the JSON_TABLE table function, used as a table
// factor: JSON_TABLE(expr, path COLUMNS (column[, column]...)).
type JSONTableExpr struct {
	node

	Doc     ExprNode
	Path    ExprNode
	Columns []*JSONTableColumn
}

func (*JSONTableExpr) resultSet() {}

// Restore implements Node interface.
func (n *JSONTableExpr) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("JSON_TABLE")
	ctx.WritePlain("(")
	if err := n.Doc.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore JSONTableExpr.Doc")
	}
	ctx.WritePlain(", ")
	if err := n.Path.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore JSONTableExpr.Path")
	}
	ctx.WriteKeyWord(" COLUMNS ")
	ctx.WritePlain("(")
	for i, col := range n.Columns {
		if i != 0 {
			ctx.WritePlain(", ")
		}
		if err := col.Restore(ctx); err != nil {
			return annotatef(err, "An error occurred while restore JSONTableExpr.Columns[%d]", i)
		}
	}
	ctx.WritePlain(")")
	ctx.WritePlain(")")
	return nil
}

// Accept implements Node Accept interface.
func (n *JSONTableExpr) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*JSONTableExpr)
	node, ok := n.Doc.Accept(v)
	if !ok {
		return n, false
	}
	n.Doc = node.(ExprNode)
	node, ok = n.Path.Accept(v)
	if !ok {
		return n, false
	}
	n.Path = node.(ExprNode)
	for i, col := range n.Columns {
		colNode, ok := col.Accept(v)
		if !ok {
			return n, false
		}
		n.Columns[i] = colNode.(*JSONTableColumn)
	}
	return v.Leave(n)
}

// JSONValueExpr is JSON_VALUE(doc, path [RETURNING type] [on_empty]
// [on_error]). The plain two-argument call without any of the optional
// clauses parses as an ordinary FuncCallExpr, not this node.
type JSONValueExpr struct {
	funcNode

	Doc  ExprNode
	Path ExprNode
	// Returning is the RETURNING cast target; nil when absent.
	Returning *types.FieldType
	OnEmpty   *JSONOnResponse
	OnError   *JSONOnResponse
}

// Restore implements Node interface.
func (n *JSONValueExpr) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("JSON_VALUE")
	ctx.WritePlain("(")
	if err := n.Doc.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore JSONValueExpr.Doc")
	}
	ctx.WritePlain(", ")
	if err := n.Path.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore JSONValueExpr.Path")
	}
	if n.Returning != nil {
		ctx.WriteKeyWord(" RETURNING ")
		n.Returning.RestoreAsCastType(ctx, false)
	}
	if n.OnEmpty != nil {
		ctx.WritePlain(" ")
		if err := n.OnEmpty.restoreWithCondition(ctx, "EMPTY"); err != nil {
			return err
		}
	}
	if n.OnError != nil {
		ctx.WritePlain(" ")
		if err := n.OnError.restoreWithCondition(ctx, "ERROR"); err != nil {
			return err
		}
	}
	ctx.WritePlain(")")
	return nil
}

// Format the ExprNode into a Writer.
func (n *JSONValueExpr) Format(w io.Writer) {
	var sb strings.Builder
	if err := n.Restore(format.NewRestoreCtx(format.DefaultRestoreFlags, &sb)); err != nil {
		return
	}
	io.WriteString(w, sb.String())
}

// Accept implements Node Accept interface.
func (n *JSONValueExpr) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*JSONValueExpr)
	node, ok := n.Doc.Accept(v)
	if !ok {
		return n, false
	}
	n.Doc = node.(ExprNode)
	node, ok = n.Path.Accept(v)
	if !ok {
		return n, false
	}
	n.Path = node.(ExprNode)
	if n.OnEmpty != nil {
		onNode, ok := n.OnEmpty.Accept(v)
		if !ok {
			return n, false
		}
		n.OnEmpty = onNode.(*JSONOnResponse)
	}
	if n.OnError != nil {
		onNode, ok := n.OnError.Accept(v)
		if !ok {
			return n, false
		}
		n.OnError = onNode.(*JSONOnResponse)
	}
	return v.Leave(n)
}
