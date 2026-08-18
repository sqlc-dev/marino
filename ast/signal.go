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

// SIGNAL, RESIGNAL, and GET DIAGNOSTICS: the condition handling
// statements of MySQL's compound statement syntax.

var (
	_ Node = &SignalConditionValue{}
	_ Node = &SignalSetItem{}
	_ Node = &DiagnosticsItem{}

	_ StmtNode = &SignalStmt{}
	_ StmtNode = &ResignalStmt{}
	_ StmtNode = &GetDiagnosticsStmt{}
)

// SignalConditionValue is the condition_value of a SIGNAL or RESIGNAL
// statement: either SQLSTATE [VALUE] 'xxxxx' or the name of a condition
// declared with DECLARE ... CONDITION. Exactly one of SQLState and
// ConditionName is set.
type SignalConditionValue struct {
	node
	SQLState      string
	ConditionName string
}

// Restore implements Node interface.
func (n *SignalConditionValue) Restore(ctx *format.RestoreCtx) error {
	if n.SQLState != "" {
		ctx.WriteKeyWord("SQLSTATE ")
		ctx.WriteString(n.SQLState)
	} else {
		ctx.WriteName(n.ConditionName)
	}
	return nil
}

// Accept implements Node Accept interface.
func (n *SignalConditionValue) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*SignalConditionValue)
	return v.Leave(n)
}

// SignalSetItem is one signal_information_item of a SIGNAL or RESIGNAL
// SET clause: a condition information item name and its value. Name is
// the canonical uppercase item name (e.g. "MESSAGE_TEXT"); Value is a
// literal or a variable.
type SignalSetItem struct {
	node
	Name  string
	Value ExprNode
}

// Restore implements Node interface.
func (n *SignalSetItem) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord(n.Name)
	ctx.WritePlain("=")
	if err := n.Value.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore SignalSetItem.Value")
	}
	return nil
}

// Accept implements Node Accept interface.
func (n *SignalSetItem) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*SignalSetItem)
	node, ok := n.Value.Accept(v)
	if !ok {
		return n, false
	}
	n.Value = node.(ExprNode)
	return v.Leave(n)
}

// SignalStmt is a SIGNAL statement:
// SIGNAL condition_value [SET signal_information_item, ...].
type SignalStmt struct {
	stmtNode
	Condition *SignalConditionValue
	SetItems  []*SignalSetItem
}

// Restore implements Node interface.
func (n *SignalStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("SIGNAL ")
	if err := n.Condition.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore SignalStmt.Condition")
	}
	return restoreSignalSetItems(ctx, n.SetItems)
}

// Accept implements Node Accept interface.
func (n *SignalStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*SignalStmt)
	node, ok := n.Condition.Accept(v)
	if !ok {
		return n, false
	}
	n.Condition = node.(*SignalConditionValue)
	for i, item := range n.SetItems {
		node, ok := item.Accept(v)
		if !ok {
			return n, false
		}
		n.SetItems[i] = node.(*SignalSetItem)
	}
	return v.Leave(n)
}

// ResignalStmt is a RESIGNAL statement:
// RESIGNAL [condition_value] [SET signal_information_item, ...].
// Condition is nil when no condition value is given.
type ResignalStmt struct {
	stmtNode
	Condition *SignalConditionValue
	SetItems  []*SignalSetItem
}

// Restore implements Node interface.
func (n *ResignalStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("RESIGNAL")
	if n.Condition != nil {
		ctx.WritePlain(" ")
		if err := n.Condition.Restore(ctx); err != nil {
			return annotate(err, "An error occurred while restore ResignalStmt.Condition")
		}
	}
	return restoreSignalSetItems(ctx, n.SetItems)
}

// Accept implements Node Accept interface.
func (n *ResignalStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*ResignalStmt)
	if n.Condition != nil {
		node, ok := n.Condition.Accept(v)
		if !ok {
			return n, false
		}
		n.Condition = node.(*SignalConditionValue)
	}
	for i, item := range n.SetItems {
		node, ok := item.Accept(v)
		if !ok {
			return n, false
		}
		n.SetItems[i] = node.(*SignalSetItem)
	}
	return v.Leave(n)
}

func restoreSignalSetItems(ctx *format.RestoreCtx, items []*SignalSetItem) error {
	for i, item := range items {
		if i == 0 {
			ctx.WriteKeyWord(" SET ")
		} else {
			ctx.WritePlain(", ")
		}
		if err := item.Restore(ctx); err != nil {
			return annotate(err, "An error occurred while restore SignalSetItem")
		}
	}
	return nil
}

// DiagnosticsArea selects which diagnostics area a GET DIAGNOSTICS
// statement reads.
type DiagnosticsArea int

const (
	// DiagnosticsAreaDefault omits the area keyword (the current area).
	DiagnosticsAreaDefault DiagnosticsArea = iota
	// DiagnosticsAreaCurrent is GET CURRENT DIAGNOSTICS.
	DiagnosticsAreaCurrent
	// DiagnosticsAreaStacked is GET STACKED DIAGNOSTICS.
	DiagnosticsAreaStacked
)

// DiagnosticsItem is one target = item_name assignment of a GET
// DIAGNOSTICS statement. Target is a user variable or a stored program
// local variable; Name is the canonical uppercase statement or condition
// information item name (e.g. "ROW_COUNT", "RETURNED_SQLSTATE").
type DiagnosticsItem struct {
	node
	Target ExprNode
	Name   string
}

// Restore implements Node interface.
func (n *DiagnosticsItem) Restore(ctx *format.RestoreCtx) error {
	if err := n.Target.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore DiagnosticsItem.Target")
	}
	ctx.WritePlain("=")
	ctx.WriteKeyWord(n.Name)
	return nil
}

// Accept implements Node Accept interface.
func (n *DiagnosticsItem) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*DiagnosticsItem)
	node, ok := n.Target.Accept(v)
	if !ok {
		return n, false
	}
	n.Target = node.(ExprNode)
	return v.Leave(n)
}

// GetDiagnosticsStmt is a GET DIAGNOSTICS statement:
//
//	GET [CURRENT | STACKED] DIAGNOSTICS
//	    { statement_information_item [, ...]
//	    | CONDITION condition_number condition_information_item [, ...] }
//
// ConditionNumber is non-nil for the CONDITION form.
type GetDiagnosticsStmt struct {
	stmtNode
	Area            DiagnosticsArea
	ConditionNumber ExprNode
	Items           []*DiagnosticsItem
}

// Restore implements Node interface.
func (n *GetDiagnosticsStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("GET ")
	switch n.Area {
	case DiagnosticsAreaCurrent:
		ctx.WriteKeyWord("CURRENT ")
	case DiagnosticsAreaStacked:
		ctx.WriteKeyWord("STACKED ")
	}
	ctx.WriteKeyWord("DIAGNOSTICS ")
	if n.ConditionNumber != nil {
		ctx.WriteKeyWord("CONDITION ")
		if err := n.ConditionNumber.Restore(ctx); err != nil {
			return annotate(err, "An error occurred while restore GetDiagnosticsStmt.ConditionNumber")
		}
		ctx.WritePlain(" ")
	}
	for i, item := range n.Items {
		if i != 0 {
			ctx.WritePlain(", ")
		}
		if err := item.Restore(ctx); err != nil {
			return annotate(err, "An error occurred while restore DiagnosticsItem")
		}
	}
	return nil
}

// Accept implements Node Accept interface.
func (n *GetDiagnosticsStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*GetDiagnosticsStmt)
	if n.ConditionNumber != nil {
		node, ok := n.ConditionNumber.Accept(v)
		if !ok {
			return n, false
		}
		n.ConditionNumber = node.(ExprNode)
	}
	for i, item := range n.Items {
		node, ok := item.Accept(v)
		if !ok {
			return n, false
		}
		n.Items[i] = node.(*DiagnosticsItem)
	}
	return v.Leave(n)
}
