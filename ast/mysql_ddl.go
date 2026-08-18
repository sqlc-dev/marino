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
	"io"

	"github.com/sqlc-dev/marino/auth"
	"github.com/sqlc-dev/marino/format"
	"github.com/sqlc-dev/marino/types"
)

// The MySQL data definition statements that postdate the goyacc grammar
// (MySQL 26.7 §15.1): events, triggers, stored functions (and the ALTER
// forms of stored routines), ALTER VIEW, servers, tablespaces, logfile
// groups, spatial reference systems, libraries, and JSON duality views.
// (The masking policy statements live in ddl.go with the other masking
// nodes, and AlterInstanceStmt in misc.go.)

var (
	_ ExprNode = &JSONDualityObjectExpr{}

	_ DDLNode = &CreateEventStmt{}
	_ DDLNode = &AlterEventStmt{}
	_ DDLNode = &DropEventStmt{}
	_ DDLNode = &CreateTriggerStmt{}
	_ DDLNode = &DropTriggerStmt{}
	_ DDLNode = &CreateFunctionStmt{}
	_ DDLNode = &AlterFunctionStmt{}
	_ DDLNode = &DropFunctionStmt{}
	_ DDLNode = &AlterProcedureStmt{}
	_ DDLNode = &AlterViewStmt{}
	_ DDLNode = &CreateServerStmt{}
	_ DDLNode = &AlterServerStmt{}
	_ DDLNode = &DropServerStmt{}
	_ DDLNode = &CreateTablespaceStmt{}
	_ DDLNode = &AlterTablespaceStmt{}
	_ DDLNode = &DropTablespaceStmt{}
	_ DDLNode = &CreateLogfileGroupStmt{}
	_ DDLNode = &AlterLogfileGroupStmt{}
	_ DDLNode = &DropLogfileGroupStmt{}
	_ DDLNode = &CreateSpatialReferenceSystemStmt{}
	_ DDLNode = &DropSpatialReferenceSystemStmt{}
	_ DDLNode = &CreateLibraryStmt{}
	_ DDLNode = &AlterLibraryStmt{}
	_ DDLNode = &DropLibraryStmt{}
	_ DDLNode = &CreateJSONDualityViewStmt{}
	_ DDLNode = &AlterJSONDualityViewStmt{}

	_ StmtNode = &ReturnStmt{}
)

// restoreDefiner writes a DEFINER = user clause followed by one space.
func restoreDefiner(ctx *format.RestoreCtx, definer *auth.UserIdentity, parent string) error {
	ctx.WriteKeyWord("DEFINER ")
	ctx.WritePlain("= ")
	if err := definer.Restore(ctx); err != nil {
		return annotatef(err, "An error occurred while restore %s.Definer", parent)
	}
	ctx.WritePlain(" ")
	return nil
}

/*
 * Events
 */

// EventSchedule is the ON SCHEDULE clause of CREATE/ALTER EVENT:
//
//	AT timestamp
//	| EVERY interval [STARTS timestamp] [ENDS timestamp]
//
// At is set for the AT form; otherwise Every holds the interval quantity
// and Unit its time unit, with optional Starts and Ends timestamps.
type EventSchedule struct {
	At     ExprNode
	Every  ExprNode
	Unit   TimeUnitType
	Starts ExprNode
	Ends   ExprNode
}

// Restore implements Node interface.
func (n *EventSchedule) Restore(ctx *format.RestoreCtx) error {
	if n.At != nil {
		ctx.WriteKeyWord("AT ")
		if err := n.At.Restore(ctx); err != nil {
			return annotate(err, "An error occurred while restore EventSchedule.At")
		}
		return nil
	}
	ctx.WriteKeyWord("EVERY ")
	if err := n.Every.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore EventSchedule.Every")
	}
	ctx.WritePlain(" ")
	ctx.WriteKeyWord(n.Unit.String())
	if n.Starts != nil {
		ctx.WriteKeyWord(" STARTS ")
		if err := n.Starts.Restore(ctx); err != nil {
			return annotate(err, "An error occurred while restore EventSchedule.Starts")
		}
	}
	if n.Ends != nil {
		ctx.WriteKeyWord(" ENDS ")
		if err := n.Ends.Restore(ctx); err != nil {
			return annotate(err, "An error occurred while restore EventSchedule.Ends")
		}
	}
	return nil
}

func (n *EventSchedule) accept(v Visitor) bool {
	for _, e := range []*ExprNode{&n.At, &n.Every, &n.Starts, &n.Ends} {
		if *e == nil {
			continue
		}
		node, ok := (*e).Accept(v)
		if !ok {
			return false
		}
		*e = node.(ExprNode)
	}
	return true
}

// EventCompletion is the ON COMPLETION clause of CREATE/ALTER EVENT.
type EventCompletion int

const (
	// EventCompletionDefault omits the clause.
	EventCompletionDefault EventCompletion = iota
	// EventCompletionPreserve is ON COMPLETION PRESERVE.
	EventCompletionPreserve
	// EventCompletionNotPreserve is ON COMPLETION NOT PRESERVE.
	EventCompletionNotPreserve
)

// EventState is the ENABLE/DISABLE clause of CREATE/ALTER EVENT.
type EventState int

const (
	// EventStateDefault omits the clause.
	EventStateDefault EventState = iota
	// EventStateEnable is ENABLE.
	EventStateEnable
	// EventStateDisable is DISABLE.
	EventStateDisable
	// EventStateDisableOnReplica is DISABLE ON REPLICA.
	EventStateDisableOnReplica
)

func restoreEventState(ctx *format.RestoreCtx, state EventState) {
	switch state {
	case EventStateEnable:
		ctx.WriteKeyWord(" ENABLE")
	case EventStateDisable:
		ctx.WriteKeyWord(" DISABLE")
	case EventStateDisableOnReplica:
		ctx.WriteKeyWord(" DISABLE ON REPLICA")
	}
}

func restoreEventCompletion(ctx *format.RestoreCtx, completion EventCompletion) {
	switch completion {
	case EventCompletionPreserve:
		ctx.WriteKeyWord(" ON COMPLETION PRESERVE")
	case EventCompletionNotPreserve:
		ctx.WriteKeyWord(" ON COMPLETION NOT PRESERVE")
	}
}

// CreateEventStmt is a CREATE EVENT statement:
//
//	CREATE [DEFINER = user] EVENT [IF NOT EXISTS] event_name
//	    ON SCHEDULE schedule [ON COMPLETION [NOT] PRESERVE]
//	    [ENABLE | DISABLE | DISABLE ON REPLICA] [COMMENT 'string']
//	    DO event_body
type CreateEventStmt struct {
	ddlNode

	IfNotExists bool
	Definer     *auth.UserIdentity // nil when absent
	EventName   *TableName
	Schedule    *EventSchedule
	Completion  EventCompletion
	State       EventState
	Comment     *string
	Body        StmtNode
}

// Restore implements Node interface.
func (n *CreateEventStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("CREATE ")
	if n.Definer != nil {
		if err := restoreDefiner(ctx, n.Definer, "CreateEventStmt"); err != nil {
			return err
		}
	}
	ctx.WriteKeyWord("EVENT ")
	if n.IfNotExists {
		ctx.WriteKeyWord("IF NOT EXISTS ")
	}
	if err := n.EventName.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore CreateEventStmt.EventName")
	}
	ctx.WriteKeyWord(" ON SCHEDULE ")
	if err := n.Schedule.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore CreateEventStmt.Schedule")
	}
	restoreEventCompletion(ctx, n.Completion)
	restoreEventState(ctx, n.State)
	if n.Comment != nil {
		ctx.WriteKeyWord(" COMMENT ")
		ctx.WriteString(*n.Comment)
	}
	ctx.WriteKeyWord(" DO ")
	if err := n.Body.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore CreateEventStmt.Body")
	}
	return nil
}

// Accept implements Node Accept interface.
func (n *CreateEventStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*CreateEventStmt)
	node, ok := n.EventName.Accept(v)
	if !ok {
		return n, false
	}
	n.EventName = node.(*TableName)
	if !n.Schedule.accept(v) {
		return n, false
	}
	if n.Body != nil {
		node, ok := n.Body.Accept(v)
		if !ok {
			return n, false
		}
		n.Body = node.(StmtNode)
	}
	return v.Leave(n)
}

// AlterEventStmt is an ALTER EVENT statement:
//
//	ALTER [DEFINER = user] EVENT event_name
//	    [ON SCHEDULE schedule] [ON COMPLETION [NOT] PRESERVE]
//	    [RENAME TO new_event_name]
//	    [ENABLE | DISABLE | DISABLE ON REPLICA] [COMMENT 'string']
//	    [DO event_body]
type AlterEventStmt struct {
	ddlNode

	Definer    *auth.UserIdentity // nil when absent
	EventName  *TableName
	Schedule   *EventSchedule // nil when absent
	Completion EventCompletion
	RenameTo   *TableName
	State      EventState
	Comment    *string
	Body       StmtNode // nil when absent
}

// Restore implements Node interface.
func (n *AlterEventStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("ALTER ")
	if n.Definer != nil {
		if err := restoreDefiner(ctx, n.Definer, "AlterEventStmt"); err != nil {
			return err
		}
	}
	ctx.WriteKeyWord("EVENT ")
	if err := n.EventName.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore AlterEventStmt.EventName")
	}
	if n.Schedule != nil {
		ctx.WriteKeyWord(" ON SCHEDULE ")
		if err := n.Schedule.Restore(ctx); err != nil {
			return annotate(err, "An error occurred while restore AlterEventStmt.Schedule")
		}
	}
	restoreEventCompletion(ctx, n.Completion)
	if n.RenameTo != nil {
		ctx.WriteKeyWord(" RENAME TO ")
		if err := n.RenameTo.Restore(ctx); err != nil {
			return annotate(err, "An error occurred while restore AlterEventStmt.RenameTo")
		}
	}
	restoreEventState(ctx, n.State)
	if n.Comment != nil {
		ctx.WriteKeyWord(" COMMENT ")
		ctx.WriteString(*n.Comment)
	}
	if n.Body != nil {
		ctx.WriteKeyWord(" DO ")
		if err := n.Body.Restore(ctx); err != nil {
			return annotate(err, "An error occurred while restore AlterEventStmt.Body")
		}
	}
	return nil
}

// Accept implements Node Accept interface.
func (n *AlterEventStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*AlterEventStmt)
	node, ok := n.EventName.Accept(v)
	if !ok {
		return n, false
	}
	n.EventName = node.(*TableName)
	if n.Schedule != nil && !n.Schedule.accept(v) {
		return n, false
	}
	if n.RenameTo != nil {
		node, ok := n.RenameTo.Accept(v)
		if !ok {
			return n, false
		}
		n.RenameTo = node.(*TableName)
	}
	if n.Body != nil {
		node, ok := n.Body.Accept(v)
		if !ok {
			return n, false
		}
		n.Body = node.(StmtNode)
	}
	return v.Leave(n)
}

// DropEventStmt is a DROP EVENT statement:
// DROP EVENT [IF EXISTS] event_name.
type DropEventStmt struct {
	ddlNode

	IfExists  bool
	EventName *TableName
}

// Restore implements Node interface.
func (n *DropEventStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("DROP EVENT ")
	if n.IfExists {
		ctx.WriteKeyWord("IF EXISTS ")
	}
	if err := n.EventName.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore DropEventStmt.EventName")
	}
	return nil
}

// Accept implements Node Accept interface.
func (n *DropEventStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*DropEventStmt)
	node, ok := n.EventName.Accept(v)
	if !ok {
		return n, false
	}
	n.EventName = node.(*TableName)
	return v.Leave(n)
}

/*
 * Triggers
 */

// TriggerTime is the BEFORE/AFTER modifier of CREATE TRIGGER.
type TriggerTime int

const (
	// TriggerTimeBefore is BEFORE.
	TriggerTimeBefore TriggerTime = iota
	// TriggerTimeAfter is AFTER.
	TriggerTimeAfter
)

// String implements fmt.Stringer interface.
func (n TriggerTime) String() string {
	switch n {
	case TriggerTimeBefore:
		return "BEFORE"
	case TriggerTimeAfter:
		return "AFTER"
	}
	return ""
}

// TriggerEvent is the INSERT/UPDATE/DELETE event of CREATE TRIGGER.
type TriggerEvent int

const (
	// TriggerEventInsert is INSERT.
	TriggerEventInsert TriggerEvent = iota
	// TriggerEventUpdate is UPDATE.
	TriggerEventUpdate
	// TriggerEventDelete is DELETE.
	TriggerEventDelete
)

// String implements fmt.Stringer interface.
func (n TriggerEvent) String() string {
	switch n {
	case TriggerEventInsert:
		return "INSERT"
	case TriggerEventUpdate:
		return "UPDATE"
	case TriggerEventDelete:
		return "DELETE"
	}
	return ""
}

// TriggerOrder is the FOLLOWS/PRECEDES clause of CREATE TRIGGER.
type TriggerOrder struct {
	Precedes     bool // FOLLOWS when false
	OtherTrigger CIStr
}

// CreateTriggerStmt is a CREATE TRIGGER statement:
//
//	CREATE [DEFINER = user] TRIGGER [IF NOT EXISTS] trigger_name
//	    {BEFORE | AFTER} {INSERT | UPDATE | DELETE} ON tbl_name
//	    FOR EACH ROW [{FOLLOWS | PRECEDES} other_trigger_name]
//	    trigger_body
type CreateTriggerStmt struct {
	ddlNode

	IfNotExists  bool
	Definer      *auth.UserIdentity // nil when absent
	TriggerName  *TableName
	TriggerTime  TriggerTime
	TriggerEvent TriggerEvent
	Table        *TableName
	Order        *TriggerOrder // nil when absent
	Body         StmtNode
}

// Restore implements Node interface.
func (n *CreateTriggerStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("CREATE ")
	if n.Definer != nil {
		if err := restoreDefiner(ctx, n.Definer, "CreateTriggerStmt"); err != nil {
			return err
		}
	}
	ctx.WriteKeyWord("TRIGGER ")
	if n.IfNotExists {
		ctx.WriteKeyWord("IF NOT EXISTS ")
	}
	if err := n.TriggerName.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore CreateTriggerStmt.TriggerName")
	}
	ctx.WritePlain(" ")
	ctx.WriteKeyWord(n.TriggerTime.String())
	ctx.WritePlain(" ")
	ctx.WriteKeyWord(n.TriggerEvent.String())
	ctx.WriteKeyWord(" ON ")
	if err := n.Table.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore CreateTriggerStmt.Table")
	}
	ctx.WriteKeyWord(" FOR EACH ROW")
	if n.Order != nil {
		if n.Order.Precedes {
			ctx.WriteKeyWord(" PRECEDES ")
		} else {
			ctx.WriteKeyWord(" FOLLOWS ")
		}
		ctx.WriteName(n.Order.OtherTrigger.O)
	}
	ctx.WritePlain(" ")
	if err := n.Body.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore CreateTriggerStmt.Body")
	}
	return nil
}

// Accept implements Node Accept interface.
func (n *CreateTriggerStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*CreateTriggerStmt)
	node, ok := n.TriggerName.Accept(v)
	if !ok {
		return n, false
	}
	n.TriggerName = node.(*TableName)
	node, ok = n.Table.Accept(v)
	if !ok {
		return n, false
	}
	n.Table = node.(*TableName)
	if n.Body != nil {
		node, ok := n.Body.Accept(v)
		if !ok {
			return n, false
		}
		n.Body = node.(StmtNode)
	}
	return v.Leave(n)
}

// DropTriggerStmt is a DROP TRIGGER statement:
// DROP TRIGGER [IF EXISTS] [schema_name.]trigger_name.
type DropTriggerStmt struct {
	ddlNode

	IfExists    bool
	TriggerName *TableName
}

// Restore implements Node interface.
func (n *DropTriggerStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("DROP TRIGGER ")
	if n.IfExists {
		ctx.WriteKeyWord("IF EXISTS ")
	}
	if err := n.TriggerName.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore DropTriggerStmt.TriggerName")
	}
	return nil
}

// Accept implements Node Accept interface.
func (n *DropTriggerStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*DropTriggerStmt)
	node, ok := n.TriggerName.Accept(v)
	if !ok {
		return n, false
	}
	n.TriggerName = node.(*TableName)
	return v.Leave(n)
}

/*
 * Stored routines
 */

// RoutineDeterminism is the [NOT] DETERMINISTIC characteristic of a
// stored routine.
type RoutineDeterminism int

const (
	// RoutineDeterminismDefault omits the characteristic.
	RoutineDeterminismDefault RoutineDeterminism = iota
	// RoutineDeterministic is DETERMINISTIC.
	RoutineDeterministic
	// RoutineNotDeterministic is NOT DETERMINISTIC.
	RoutineNotDeterministic
)

// RoutineDataAccess is the SQL data access characteristic of a stored
// routine.
type RoutineDataAccess int

const (
	// RoutineDataAccessDefault omits the characteristic.
	RoutineDataAccessDefault RoutineDataAccess = iota
	// RoutineContainsSQL is CONTAINS SQL.
	RoutineContainsSQL
	// RoutineNoSQL is NO SQL.
	RoutineNoSQL
	// RoutineReadsSQLData is READS SQL DATA.
	RoutineReadsSQLData
	// RoutineModifiesSQLData is MODIFIES SQL DATA.
	RoutineModifiesSQLData
)

// RoutineCharacteristics are the characteristics of a stored routine.
// They restore in the canonical order COMMENT, LANGUAGE, [NOT]
// DETERMINISTIC, data access, SQL SECURITY regardless of source order.
type RoutineCharacteristics struct {
	Comment     *string
	Language    CIStr // empty when absent (SQL, JAVASCRIPT, ...)
	Determinism RoutineDeterminism
	DataAccess  RoutineDataAccess
	Security    *ViewSecurity // nil when absent
}

// Restore implements Node interface. Every present characteristic is
// written with one leading space.
func (n *RoutineCharacteristics) Restore(ctx *format.RestoreCtx) error {
	if n.Comment != nil {
		ctx.WriteKeyWord(" COMMENT ")
		ctx.WriteString(*n.Comment)
	}
	if n.Language.O != "" {
		ctx.WriteKeyWord(" LANGUAGE ")
		ctx.WriteKeyWord(n.Language.O)
	}
	switch n.Determinism {
	case RoutineDeterministic:
		ctx.WriteKeyWord(" DETERMINISTIC")
	case RoutineNotDeterministic:
		ctx.WriteKeyWord(" NOT DETERMINISTIC")
	}
	switch n.DataAccess {
	case RoutineContainsSQL:
		ctx.WriteKeyWord(" CONTAINS SQL")
	case RoutineNoSQL:
		ctx.WriteKeyWord(" NO SQL")
	case RoutineReadsSQLData:
		ctx.WriteKeyWord(" READS SQL DATA")
	case RoutineModifiesSQLData:
		ctx.WriteKeyWord(" MODIFIES SQL DATA")
	}
	if n.Security != nil {
		ctx.WriteKeyWord(" SQL SECURITY ")
		if *n.Security == SecurityInvoker {
			ctx.WriteKeyWord("INVOKER")
		} else {
			ctx.WriteKeyWord("DEFINER")
		}
	}
	return nil
}

// FunctionParam is one parameter of a stored function; unlike procedure
// parameters (StoreParameter) it has no IN/OUT/INOUT mode.
type FunctionParam struct {
	ParamName string
	ParamType *types.FieldType
}

// Restore implements Node interface.
func (n *FunctionParam) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteName(n.ParamName)
	ctx.WritePlain(" ")
	if err := n.ParamType.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore FunctionParam.ParamType")
	}
	return nil
}

// CreateFunctionStmt is a CREATE FUNCTION statement for stored
// functions (CreateLoadableFunctionStmt is the loadable function form):
//
//	CREATE [DEFINER = user] FUNCTION [IF NOT EXISTS] sp_name
//	    ([func_parameter[, ...]]) RETURNS type [characteristic ...]
//	    routine_body
type CreateFunctionStmt struct {
	ddlNode

	IfNotExists     bool
	Definer         *auth.UserIdentity // nil when absent
	FunctionName    *TableName
	Params          []*FunctionParam
	ReturnType      *types.FieldType
	Characteristics RoutineCharacteristics
	Body            StmtNode
}

// Restore implements Node interface.
func (n *CreateFunctionStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("CREATE ")
	if n.Definer != nil {
		if err := restoreDefiner(ctx, n.Definer, "CreateFunctionStmt"); err != nil {
			return err
		}
	}
	ctx.WriteKeyWord("FUNCTION ")
	if n.IfNotExists {
		ctx.WriteKeyWord("IF NOT EXISTS ")
	}
	if err := n.FunctionName.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore CreateFunctionStmt.FunctionName")
	}
	ctx.WritePlain("(")
	for i, param := range n.Params {
		if i != 0 {
			ctx.WritePlain(", ")
		}
		if err := param.Restore(ctx); err != nil {
			return annotatef(err, "An error occurred while restore CreateFunctionStmt.Params[%d]", i)
		}
	}
	ctx.WritePlain(")")
	ctx.WriteKeyWord(" RETURNS ")
	if err := n.ReturnType.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore CreateFunctionStmt.ReturnType")
	}
	if err := n.Characteristics.Restore(ctx); err != nil {
		return err
	}
	ctx.WritePlain(" ")
	if err := n.Body.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore CreateFunctionStmt.Body")
	}
	return nil
}

// Accept implements Node Accept interface.
func (n *CreateFunctionStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*CreateFunctionStmt)
	node, ok := n.FunctionName.Accept(v)
	if !ok {
		return n, false
	}
	n.FunctionName = node.(*TableName)
	if n.Body != nil {
		node, ok := n.Body.Accept(v)
		if !ok {
			return n, false
		}
		n.Body = node.(StmtNode)
	}
	return v.Leave(n)
}

// AlterFunctionStmt is an ALTER FUNCTION statement:
// ALTER FUNCTION func_name [characteristic ...].
type AlterFunctionStmt struct {
	ddlNode

	FunctionName    *TableName
	Characteristics RoutineCharacteristics
}

// Restore implements Node interface.
func (n *AlterFunctionStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("ALTER FUNCTION ")
	if err := n.FunctionName.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore AlterFunctionStmt.FunctionName")
	}
	return n.Characteristics.Restore(ctx)
}

// Accept implements Node Accept interface.
func (n *AlterFunctionStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*AlterFunctionStmt)
	node, ok := n.FunctionName.Accept(v)
	if !ok {
		return n, false
	}
	n.FunctionName = node.(*TableName)
	return v.Leave(n)
}

// DropFunctionStmt is a DROP FUNCTION statement (for a stored or
// loadable function): DROP FUNCTION [IF EXISTS] func_name.
type DropFunctionStmt struct {
	ddlNode

	IfExists     bool
	FunctionName *TableName
}

// Restore implements Node interface.
func (n *DropFunctionStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("DROP FUNCTION ")
	if n.IfExists {
		ctx.WriteKeyWord("IF EXISTS ")
	}
	if err := n.FunctionName.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore DropFunctionStmt.FunctionName")
	}
	return nil
}

// Accept implements Node Accept interface.
func (n *DropFunctionStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*DropFunctionStmt)
	node, ok := n.FunctionName.Accept(v)
	if !ok {
		return n, false
	}
	n.FunctionName = node.(*TableName)
	return v.Leave(n)
}

// AlterProcedureStmt is an ALTER PROCEDURE statement:
// ALTER PROCEDURE proc_name [characteristic ...].
type AlterProcedureStmt struct {
	ddlNode

	ProcedureName   *TableName
	Characteristics RoutineCharacteristics
}

// Restore implements Node interface.
func (n *AlterProcedureStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("ALTER PROCEDURE ")
	if err := n.ProcedureName.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore AlterProcedureStmt.ProcedureName")
	}
	return n.Characteristics.Restore(ctx)
}

// Accept implements Node Accept interface.
func (n *AlterProcedureStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*AlterProcedureStmt)
	node, ok := n.ProcedureName.Accept(v)
	if !ok {
		return n, false
	}
	n.ProcedureName = node.(*TableName)
	return v.Leave(n)
}

// ReturnStmt is a RETURN statement in a stored function body:
// RETURN expr.
type ReturnStmt struct {
	stmtNode

	Expr ExprNode
}

// Restore implements Node interface.
func (n *ReturnStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("RETURN ")
	if err := n.Expr.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore ReturnStmt.Expr")
	}
	return nil
}

// Accept implements Node Accept interface.
func (n *ReturnStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*ReturnStmt)
	node, ok := n.Expr.Accept(v)
	if !ok {
		return n, false
	}
	n.Expr = node.(ExprNode)
	return v.Leave(n)
}

/*
 * ALTER VIEW
 */

// AlterViewStmt is an ALTER VIEW statement:
//
//	ALTER [ALGORITHM = {UNDEFINED | MERGE | TEMPTABLE}]
//	    [DEFINER = user] [SQL SECURITY {DEFINER | INVOKER}]
//	    VIEW view_name [(column_list)] AS select_statement
//	    [WITH [CASCADED | LOCAL] CHECK OPTION]
//
// The fields mirror CreateViewStmt: Definer defaults to CURRENT_USER
// and CheckOption to CASCADED when their clauses are absent.
type AlterViewStmt struct {
	ddlNode

	ViewName    *TableName
	Cols        []CIStr
	Select      StmtNode
	Algorithm   ViewAlgorithm
	Definer     *auth.UserIdentity
	Security    ViewSecurity
	CheckOption ViewCheckOption
}

// Restore implements Node interface.
func (n *AlterViewStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("ALTER ")
	ctx.WriteKeyWord("ALGORITHM ")
	ctx.WritePlainf("= %s ", n.Algorithm.String())
	if err := restoreDefiner(ctx, n.Definer, "AlterViewStmt"); err != nil {
		return err
	}
	ctx.WriteKeyWord("SQL SECURITY ")
	ctx.WriteKeyWord(n.Security.String())
	ctx.WriteKeyWord(" VIEW ")
	if err := n.ViewName.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore AlterViewStmt.ViewName")
	}
	for i, col := range n.Cols {
		if i == 0 {
			ctx.WritePlain(" (")
		} else {
			ctx.WritePlain(",")
		}
		ctx.WriteName(col.O)
		if i == len(n.Cols)-1 {
			ctx.WritePlain(")")
		}
	}
	ctx.WriteKeyWord(" AS ")
	if err := n.Select.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore AlterViewStmt.Select")
	}
	if n.CheckOption != CheckOptionCascaded {
		ctx.WriteKeyWord(" WITH ")
		ctx.WriteKeyWord(n.CheckOption.String())
		ctx.WriteKeyWord(" CHECK OPTION")
	}
	return nil
}

// Accept implements Node Accept interface.
func (n *AlterViewStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*AlterViewStmt)
	node, ok := n.ViewName.Accept(v)
	if !ok {
		return n, false
	}
	n.ViewName = node.(*TableName)
	selnode, ok := n.Select.Accept(v)
	if !ok {
		return n, false
	}
	n.Select = selnode.(StmtNode)
	return v.Leave(n)
}

/*
 * Servers
 */

// ServerOption is one option of CREATE/ALTER SERVER:
// {HOST | DATABASE | USER | PASSWORD | SOCKET | OWNER} 'string'
// | PORT port_num. Names are stored uppercase.
type ServerOption struct {
	Name  string
	Value ValueExpr
}

// Restore implements Node interface.
func (n *ServerOption) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord(n.Name)
	ctx.WritePlain(" ")
	if err := n.Value.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore ServerOption.Value")
	}
	return nil
}

// CreateServerStmt is a CREATE SERVER statement:
//
//	CREATE SERVER server_name FOREIGN DATA WRAPPER wrapper_name
//	    OPTIONS (option [, option] ...)
type CreateServerStmt struct {
	ddlNode

	ServerName CIStr
	Wrapper    CIStr
	Options    []*ServerOption
}

// Restore implements Node interface.
func (n *CreateServerStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("CREATE SERVER ")
	ctx.WriteName(n.ServerName.O)
	ctx.WriteKeyWord(" FOREIGN DATA WRAPPER ")
	ctx.WriteName(n.Wrapper.O)
	ctx.WriteKeyWord(" OPTIONS ")
	ctx.WritePlain("(")
	for i, opt := range n.Options {
		if i != 0 {
			ctx.WritePlain(", ")
		}
		if err := opt.Restore(ctx); err != nil {
			return annotatef(err, "An error occurred while restore CreateServerStmt.Options[%d]", i)
		}
	}
	ctx.WritePlain(")")
	return nil
}

// Accept implements Node Accept interface.
func (n *CreateServerStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*CreateServerStmt)
	return v.Leave(n)
}

// AlterServerStmt is an ALTER SERVER statement:
// ALTER SERVER server_name OPTIONS (option [, option] ...).
type AlterServerStmt struct {
	ddlNode

	ServerName CIStr
	Options    []*ServerOption
}

// Restore implements Node interface.
func (n *AlterServerStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("ALTER SERVER ")
	ctx.WriteName(n.ServerName.O)
	ctx.WriteKeyWord(" OPTIONS ")
	ctx.WritePlain("(")
	for i, opt := range n.Options {
		if i != 0 {
			ctx.WritePlain(", ")
		}
		if err := opt.Restore(ctx); err != nil {
			return annotatef(err, "An error occurred while restore AlterServerStmt.Options[%d]", i)
		}
	}
	ctx.WritePlain(")")
	return nil
}

// Accept implements Node Accept interface.
func (n *AlterServerStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*AlterServerStmt)
	return v.Leave(n)
}

// DropServerStmt is a DROP SERVER statement:
// DROP SERVER [IF EXISTS] server_name.
type DropServerStmt struct {
	ddlNode

	IfExists   bool
	ServerName CIStr
}

// Restore implements Node interface.
func (n *DropServerStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("DROP SERVER ")
	if n.IfExists {
		ctx.WriteKeyWord("IF EXISTS ")
	}
	ctx.WriteName(n.ServerName.O)
	return nil
}

// Accept implements Node Accept interface.
func (n *DropServerStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*DropServerStmt)
	return v.Leave(n)
}

/*
 * Tablespaces and logfile groups
 */

// TablespaceOptionType names one option of the CREATE/ALTER/DROP
// TABLESPACE and CREATE/ALTER/DROP LOGFILE GROUP statements.
type TablespaceOptionType int

const (
	// TablespaceOptionInitialSize is INITIAL_SIZE.
	TablespaceOptionInitialSize TablespaceOptionType = iota
	// TablespaceOptionAutoextendSize is AUTOEXTEND_SIZE.
	TablespaceOptionAutoextendSize
	// TablespaceOptionMaxSize is MAX_SIZE.
	TablespaceOptionMaxSize
	// TablespaceOptionExtentSize is EXTENT_SIZE.
	TablespaceOptionExtentSize
	// TablespaceOptionNodegroup is NODEGROUP.
	TablespaceOptionNodegroup
	// TablespaceOptionFileBlockSize is FILE_BLOCK_SIZE.
	TablespaceOptionFileBlockSize
	// TablespaceOptionUndoBufferSize is UNDO_BUFFER_SIZE.
	TablespaceOptionUndoBufferSize
	// TablespaceOptionRedoBufferSize is REDO_BUFFER_SIZE.
	TablespaceOptionRedoBufferSize
	// TablespaceOptionWait is WAIT.
	TablespaceOptionWait
	// TablespaceOptionEncryption is ENCRYPTION.
	TablespaceOptionEncryption
	// TablespaceOptionComment is COMMENT.
	TablespaceOptionComment
	// TablespaceOptionEngine is ENGINE.
	TablespaceOptionEngine
	// TablespaceOptionEngineAttribute is ENGINE_ATTRIBUTE.
	TablespaceOptionEngineAttribute
)

// String implements fmt.Stringer interface.
func (n TablespaceOptionType) String() string {
	switch n {
	case TablespaceOptionInitialSize:
		return "INITIAL_SIZE"
	case TablespaceOptionAutoextendSize:
		return "AUTOEXTEND_SIZE"
	case TablespaceOptionMaxSize:
		return "MAX_SIZE"
	case TablespaceOptionExtentSize:
		return "EXTENT_SIZE"
	case TablespaceOptionNodegroup:
		return "NODEGROUP"
	case TablespaceOptionFileBlockSize:
		return "FILE_BLOCK_SIZE"
	case TablespaceOptionUndoBufferSize:
		return "UNDO_BUFFER_SIZE"
	case TablespaceOptionRedoBufferSize:
		return "REDO_BUFFER_SIZE"
	case TablespaceOptionWait:
		return "WAIT"
	case TablespaceOptionEncryption:
		return "ENCRYPTION"
	case TablespaceOptionComment:
		return "COMMENT"
	case TablespaceOptionEngine:
		return "ENGINE"
	case TablespaceOptionEngineAttribute:
		return "ENGINE_ATTRIBUTE"
	}
	return ""
}

// TablespaceOption is one option of a tablespace or logfile group
// statement. The size options, NODEGROUP, and FILE_BLOCK_SIZE use
// UintValue; ENCRYPTION, COMMENT, and ENGINE_ATTRIBUTE use StrValue
// restored as a string literal; ENGINE uses StrValue restored as a
// plain word; WAIT uses neither.
type TablespaceOption struct {
	Tp        TablespaceOptionType
	StrValue  string
	UintValue uint64
}

// Restore implements Node interface.
func (n *TablespaceOption) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord(n.Tp.String())
	switch n.Tp {
	case TablespaceOptionWait:
	case TablespaceOptionEngine:
		ctx.WritePlain(" = ")
		ctx.WriteKeyWord(n.StrValue)
	case TablespaceOptionEncryption, TablespaceOptionComment, TablespaceOptionEngineAttribute:
		ctx.WritePlain(" = ")
		ctx.WriteString(n.StrValue)
	default:
		ctx.WritePlainf(" = %d", n.UintValue)
	}
	return nil
}

func restoreTablespaceOptions(ctx *format.RestoreCtx, options []*TablespaceOption, parent string) error {
	for i, opt := range options {
		ctx.WritePlain(" ")
		if err := opt.Restore(ctx); err != nil {
			return annotatef(err, "An error occurred while restore %s.Options[%d]", parent, i)
		}
	}
	return nil
}

// CreateTablespaceStmt is a CREATE TABLESPACE statement:
//
//	CREATE [UNDO] TABLESPACE tablespace_name [ADD DATAFILE 'file_name']
//	    [USE LOGFILE GROUP logfile_group] [option ...]
type CreateTablespaceStmt struct {
	ddlNode

	Undo            bool
	Name            CIStr
	Datafile        string // ADD DATAFILE clause; empty when absent
	UseLogfileGroup CIStr  // USE LOGFILE GROUP clause; empty when absent
	Options         []*TablespaceOption
}

// Restore implements Node interface.
func (n *CreateTablespaceStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("CREATE ")
	if n.Undo {
		ctx.WriteKeyWord("UNDO ")
	}
	ctx.WriteKeyWord("TABLESPACE ")
	ctx.WriteName(n.Name.O)
	if n.Datafile != "" {
		ctx.WriteKeyWord(" ADD DATAFILE ")
		ctx.WriteString(n.Datafile)
	}
	if n.UseLogfileGroup.O != "" {
		ctx.WriteKeyWord(" USE LOGFILE GROUP ")
		ctx.WriteName(n.UseLogfileGroup.O)
	}
	return restoreTablespaceOptions(ctx, n.Options, "CreateTablespaceStmt")
}

// Accept implements Node Accept interface.
func (n *CreateTablespaceStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*CreateTablespaceStmt)
	return v.Leave(n)
}

// AlterTablespaceOp selects the form of an ALTER TABLESPACE statement.
type AlterTablespaceOp int

const (
	// AlterTablespaceOptionsOnly changes options alone.
	AlterTablespaceOptionsOnly AlterTablespaceOp = iota
	// AlterTablespaceAddDatafile is ADD DATAFILE.
	AlterTablespaceAddDatafile
	// AlterTablespaceDropDatafile is DROP DATAFILE.
	AlterTablespaceDropDatafile
	// AlterTablespaceRenameTo is RENAME TO.
	AlterTablespaceRenameTo
	// AlterTablespaceSetActive is SET ACTIVE.
	AlterTablespaceSetActive
	// AlterTablespaceSetInactive is SET INACTIVE.
	AlterTablespaceSetInactive
)

// AlterTablespaceStmt is an ALTER [UNDO] TABLESPACE statement:
//
//	ALTER [UNDO] TABLESPACE tablespace_name
//	    {{ADD | DROP} DATAFILE 'file_name' | RENAME TO new_name
//	     | SET {ACTIVE | INACTIVE}} [option ...]
type AlterTablespaceStmt struct {
	ddlNode

	Undo     bool
	Name     CIStr
	Op       AlterTablespaceOp
	Datafile string
	NewName  CIStr
	Options  []*TablespaceOption
}

// Restore implements Node interface.
func (n *AlterTablespaceStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("ALTER ")
	if n.Undo {
		ctx.WriteKeyWord("UNDO ")
	}
	ctx.WriteKeyWord("TABLESPACE ")
	ctx.WriteName(n.Name.O)
	switch n.Op {
	case AlterTablespaceAddDatafile:
		ctx.WriteKeyWord(" ADD DATAFILE ")
		ctx.WriteString(n.Datafile)
	case AlterTablespaceDropDatafile:
		ctx.WriteKeyWord(" DROP DATAFILE ")
		ctx.WriteString(n.Datafile)
	case AlterTablespaceRenameTo:
		ctx.WriteKeyWord(" RENAME TO ")
		ctx.WriteName(n.NewName.O)
	case AlterTablespaceSetActive:
		ctx.WriteKeyWord(" SET ACTIVE")
	case AlterTablespaceSetInactive:
		ctx.WriteKeyWord(" SET INACTIVE")
	}
	return restoreTablespaceOptions(ctx, n.Options, "AlterTablespaceStmt")
}

// Accept implements Node Accept interface.
func (n *AlterTablespaceStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*AlterTablespaceStmt)
	return v.Leave(n)
}

// DropTablespaceStmt is a DROP [UNDO] TABLESPACE statement:
// DROP [UNDO] TABLESPACE tablespace_name [ENGINE [=] engine_name].
type DropTablespaceStmt struct {
	ddlNode

	Undo    bool
	Name    CIStr
	Options []*TablespaceOption
}

// Restore implements Node interface.
func (n *DropTablespaceStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("DROP ")
	if n.Undo {
		ctx.WriteKeyWord("UNDO ")
	}
	ctx.WriteKeyWord("TABLESPACE ")
	ctx.WriteName(n.Name.O)
	return restoreTablespaceOptions(ctx, n.Options, "DropTablespaceStmt")
}

// Accept implements Node Accept interface.
func (n *DropTablespaceStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*DropTablespaceStmt)
	return v.Leave(n)
}

// CreateLogfileGroupStmt is a CREATE LOGFILE GROUP statement:
//
//	CREATE LOGFILE GROUP logfile_group ADD UNDOFILE 'undo_file'
//	    [option ...]
type CreateLogfileGroupStmt struct {
	ddlNode

	GroupName CIStr
	Undofile  string
	Options   []*TablespaceOption
}

// Restore implements Node interface.
func (n *CreateLogfileGroupStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("CREATE LOGFILE GROUP ")
	ctx.WriteName(n.GroupName.O)
	ctx.WriteKeyWord(" ADD UNDOFILE ")
	ctx.WriteString(n.Undofile)
	return restoreTablespaceOptions(ctx, n.Options, "CreateLogfileGroupStmt")
}

// Accept implements Node Accept interface.
func (n *CreateLogfileGroupStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*CreateLogfileGroupStmt)
	return v.Leave(n)
}

// AlterLogfileGroupStmt is an ALTER LOGFILE GROUP statement:
//
//	ALTER LOGFILE GROUP logfile_group ADD UNDOFILE 'file_name'
//	    [option ...]
type AlterLogfileGroupStmt struct {
	ddlNode

	GroupName CIStr
	Undofile  string
	Options   []*TablespaceOption
}

// Restore implements Node interface.
func (n *AlterLogfileGroupStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("ALTER LOGFILE GROUP ")
	ctx.WriteName(n.GroupName.O)
	ctx.WriteKeyWord(" ADD UNDOFILE ")
	ctx.WriteString(n.Undofile)
	return restoreTablespaceOptions(ctx, n.Options, "AlterLogfileGroupStmt")
}

// Accept implements Node Accept interface.
func (n *AlterLogfileGroupStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*AlterLogfileGroupStmt)
	return v.Leave(n)
}

// DropLogfileGroupStmt is a DROP LOGFILE GROUP statement:
// DROP LOGFILE GROUP logfile_group ENGINE [=] engine_name.
type DropLogfileGroupStmt struct {
	ddlNode

	GroupName CIStr
	Options   []*TablespaceOption
}

// Restore implements Node interface.
func (n *DropLogfileGroupStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("DROP LOGFILE GROUP ")
	ctx.WriteName(n.GroupName.O)
	return restoreTablespaceOptions(ctx, n.Options, "DropLogfileGroupStmt")
}

// Accept implements Node Accept interface.
func (n *DropLogfileGroupStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*DropLogfileGroupStmt)
	return v.Leave(n)
}

/*
 * Spatial reference systems
 */

// SRSAttributeType names one attribute of CREATE SPATIAL REFERENCE
// SYSTEM.
type SRSAttributeType int

const (
	// SRSAttributeName is NAME 'srs_name'.
	SRSAttributeName SRSAttributeType = iota
	// SRSAttributeDefinition is DEFINITION 'definition'.
	SRSAttributeDefinition
	// SRSAttributeOrganization is ORGANIZATION 'org_name' IDENTIFIED BY
	// srs_id.
	SRSAttributeOrganization
	// SRSAttributeDescription is DESCRIPTION 'description'.
	SRSAttributeDescription
)

// String implements fmt.Stringer interface.
func (n SRSAttributeType) String() string {
	switch n {
	case SRSAttributeName:
		return "NAME"
	case SRSAttributeDefinition:
		return "DEFINITION"
	case SRSAttributeOrganization:
		return "ORGANIZATION"
	case SRSAttributeDescription:
		return "DESCRIPTION"
	}
	return ""
}

// SRSAttribute is one attribute of CREATE SPATIAL REFERENCE SYSTEM,
// kept in source order. OrgID belongs to SRSAttributeOrganization only.
type SRSAttribute struct {
	Tp       SRSAttributeType
	StrValue string
	OrgID    uint64
}

// Restore implements Node interface.
func (n *SRSAttribute) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord(n.Tp.String())
	ctx.WritePlain(" ")
	ctx.WriteString(n.StrValue)
	if n.Tp == SRSAttributeOrganization {
		ctx.WriteKeyWord(" IDENTIFIED BY ")
		ctx.WritePlainf("%d", n.OrgID)
	}
	return nil
}

// CreateSpatialReferenceSystemStmt is a CREATE SPATIAL REFERENCE SYSTEM
// statement:
//
//	CREATE [OR REPLACE] SPATIAL REFERENCE SYSTEM [IF NOT EXISTS] srid
//	    srs_attribute ...
type CreateSpatialReferenceSystemStmt struct {
	ddlNode

	OrReplace   bool
	IfNotExists bool
	Srid        uint64
	Attributes  []*SRSAttribute
}

// Restore implements Node interface.
func (n *CreateSpatialReferenceSystemStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("CREATE ")
	if n.OrReplace {
		ctx.WriteKeyWord("OR REPLACE ")
	}
	ctx.WriteKeyWord("SPATIAL REFERENCE SYSTEM ")
	if n.IfNotExists {
		ctx.WriteKeyWord("IF NOT EXISTS ")
	}
	ctx.WritePlainf("%d", n.Srid)
	for i, attr := range n.Attributes {
		ctx.WritePlain(" ")
		if err := attr.Restore(ctx); err != nil {
			return annotatef(err, "An error occurred while restore CreateSpatialReferenceSystemStmt.Attributes[%d]", i)
		}
	}
	return nil
}

// Accept implements Node Accept interface.
func (n *CreateSpatialReferenceSystemStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*CreateSpatialReferenceSystemStmt)
	return v.Leave(n)
}

// DropSpatialReferenceSystemStmt is a DROP SPATIAL REFERENCE SYSTEM
// statement: DROP SPATIAL REFERENCE SYSTEM [IF EXISTS] srid.
type DropSpatialReferenceSystemStmt struct {
	ddlNode

	IfExists bool
	Srid     uint64
}

// Restore implements Node interface.
func (n *DropSpatialReferenceSystemStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("DROP SPATIAL REFERENCE SYSTEM ")
	if n.IfExists {
		ctx.WriteKeyWord("IF EXISTS ")
	}
	ctx.WritePlainf("%d", n.Srid)
	return nil
}

// Accept implements Node Accept interface.
func (n *DropSpatialReferenceSystemStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*DropSpatialReferenceSystemStmt)
	return v.Leave(n)
}

/*
 * Libraries
 */

// CreateLibraryStmt is a CREATE LIBRARY statement:
//
//	CREATE [OR REPLACE] LIBRARY [IF NOT EXISTS] [schema.]library_name
//	    LANGUAGE language_name AS 'code'
type CreateLibraryStmt struct {
	ddlNode

	OrReplace   bool
	IfNotExists bool
	LibraryName *TableName
	Language    CIStr
	Code        string
}

// Restore implements Node interface.
func (n *CreateLibraryStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("CREATE ")
	if n.OrReplace {
		ctx.WriteKeyWord("OR REPLACE ")
	}
	ctx.WriteKeyWord("LIBRARY ")
	if n.IfNotExists {
		ctx.WriteKeyWord("IF NOT EXISTS ")
	}
	if err := n.LibraryName.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore CreateLibraryStmt.LibraryName")
	}
	ctx.WriteKeyWord(" LANGUAGE ")
	ctx.WriteKeyWord(n.Language.O)
	ctx.WriteKeyWord(" AS ")
	ctx.WriteString(n.Code)
	return nil
}

// Accept implements Node Accept interface.
func (n *CreateLibraryStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*CreateLibraryStmt)
	node, ok := n.LibraryName.Accept(v)
	if !ok {
		return n, false
	}
	n.LibraryName = node.(*TableName)
	return v.Leave(n)
}

// AlterLibraryStmt is an ALTER LIBRARY statement:
// ALTER LIBRARY [schema.]library_name COMMENT 'string'.
type AlterLibraryStmt struct {
	ddlNode

	LibraryName *TableName
	Comment     string
}

// Restore implements Node interface.
func (n *AlterLibraryStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("ALTER LIBRARY ")
	if err := n.LibraryName.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore AlterLibraryStmt.LibraryName")
	}
	ctx.WriteKeyWord(" COMMENT ")
	ctx.WriteString(n.Comment)
	return nil
}

// Accept implements Node Accept interface.
func (n *AlterLibraryStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*AlterLibraryStmt)
	node, ok := n.LibraryName.Accept(v)
	if !ok {
		return n, false
	}
	n.LibraryName = node.(*TableName)
	return v.Leave(n)
}

// DropLibraryStmt is a DROP LIBRARY statement:
// DROP LIBRARY [IF EXISTS] [schema.]library_name.
type DropLibraryStmt struct {
	ddlNode

	IfExists    bool
	LibraryName *TableName
}

// Restore implements Node interface.
func (n *DropLibraryStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("DROP LIBRARY ")
	if n.IfExists {
		ctx.WriteKeyWord("IF EXISTS ")
	}
	if err := n.LibraryName.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore DropLibraryStmt.LibraryName")
	}
	return nil
}

// Accept implements Node Accept interface.
func (n *DropLibraryStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*DropLibraryStmt)
	node, ok := n.LibraryName.Accept(v)
	if !ok {
		return n, false
	}
	n.LibraryName = node.(*TableName)
	return v.Leave(n)
}

/*
 * JSON duality views
 */

// CreateJSONDualityViewStmt is a CREATE JSON DUALITY VIEW statement:
//
//	CREATE [OR REPLACE] JSON DUALITY VIEW [IF NOT EXISTS] view_name
//	    AS select_statement
type CreateJSONDualityViewStmt struct {
	ddlNode

	OrReplace   bool
	IfNotExists bool
	ViewName    *TableName
	Select      StmtNode
}

// Restore implements Node interface.
func (n *CreateJSONDualityViewStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("CREATE ")
	if n.OrReplace {
		ctx.WriteKeyWord("OR REPLACE ")
	}
	ctx.WriteKeyWord("JSON DUALITY VIEW ")
	if n.IfNotExists {
		ctx.WriteKeyWord("IF NOT EXISTS ")
	}
	if err := n.ViewName.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore CreateJSONDualityViewStmt.ViewName")
	}
	ctx.WriteKeyWord(" AS ")
	if err := n.Select.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore CreateJSONDualityViewStmt.Select")
	}
	return nil
}

// Accept implements Node Accept interface.
func (n *CreateJSONDualityViewStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*CreateJSONDualityViewStmt)
	node, ok := n.ViewName.Accept(v)
	if !ok {
		return n, false
	}
	n.ViewName = node.(*TableName)
	selnode, ok := n.Select.Accept(v)
	if !ok {
		return n, false
	}
	n.Select = selnode.(StmtNode)
	return v.Leave(n)
}

// AlterJSONDualityViewStmt is an ALTER JSON DUALITY VIEW statement:
// ALTER JSON DUALITY VIEW view_name AS select_statement.
type AlterJSONDualityViewStmt struct {
	ddlNode

	ViewName *TableName
	Select   StmtNode
}

// Restore implements Node interface.
func (n *AlterJSONDualityViewStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("ALTER JSON DUALITY VIEW ")
	if err := n.ViewName.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore AlterJSONDualityViewStmt.ViewName")
	}
	ctx.WriteKeyWord(" AS ")
	if err := n.Select.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore AlterJSONDualityViewStmt.Select")
	}
	return nil
}

// Accept implements Node Accept interface.
func (n *AlterJSONDualityViewStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*AlterJSONDualityViewStmt)
	node, ok := n.ViewName.Accept(v)
	if !ok {
		return n, false
	}
	n.ViewName = node.(*TableName)
	selnode, ok := n.Select.Accept(v)
	if !ok {
		return n, false
	}
	n.Select = selnode.(StmtNode)
	return v.Leave(n)
}

// JSONDualityObjectItem is one 'key' : value pair of a
// JSON_DUALITY_OBJECT expression.
type JSONDualityObjectItem struct {
	Key   string
	Value ExprNode
}

// JSONDualityObjectExpr is a JSON_DUALITY_OBJECT('key' : value, ...)
// expression, the select-list constructor of JSON duality views.
type JSONDualityObjectExpr struct {
	exprNode

	Items []*JSONDualityObjectItem
}

// Restore implements Node interface.
func (n *JSONDualityObjectExpr) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("JSON_DUALITY_OBJECT")
	ctx.WritePlain("(")
	for i, item := range n.Items {
		if i != 0 {
			ctx.WritePlain(", ")
		}
		ctx.WriteString(item.Key)
		ctx.WritePlain(" : ")
		if err := item.Value.Restore(ctx); err != nil {
			return annotatef(err, "An error occurred while restore JSONDualityObjectExpr.Items[%d]", i)
		}
	}
	ctx.WritePlain(")")
	return nil
}

// Format the ExprNode into a Writer.
func (n *JSONDualityObjectExpr) Format(w io.Writer) {
	panic("Not implemented")
}

// Accept implements Node Accept interface.
func (n *JSONDualityObjectExpr) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*JSONDualityObjectExpr)
	for _, item := range n.Items {
		node, ok := item.Value.Accept(v)
		if !ok {
			return n, false
		}
		item.Value = node.(ExprNode)
	}
	return v.Leave(n)
}
