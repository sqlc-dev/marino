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

	"github.com/sqlc-dev/marino/format"
)

// The MySQL replication statements (MySQL 26.7 §15.4): PURGE BINARY
// LOGS and RESET BINARY LOGS AND GTIDS on the source side, and CHANGE
// REPLICATION FILTER, RESET REPLICA, START/STOP REPLICA, and START/STOP
// GROUP_REPLICATION on the replica side. (CHANGE REPLICATION SOURCE TO
// and its ReplicationSourceOption node live in misc.go.)

var (
	_ StmtNode = &PurgeBinaryLogsStmt{}
	_ StmtNode = &ResetBinaryLogsAndGtidsStmt{}
	_ StmtNode = &ChangeReplicationFilterStmt{}
	_ StmtNode = &ResetReplicaStmt{}
	_ StmtNode = &StartReplicaStmt{}
	_ StmtNode = &StopReplicaStmt{}
	_ StmtNode = &StartGroupReplicationStmt{}
	_ StmtNode = &StopGroupReplicationStmt{}

	_ SensitiveStmtNode = &StartReplicaStmt{}
	_ SensitiveStmtNode = &StartGroupReplicationStmt{}
)

// PurgeBinaryLogsStmt is a PURGE BINARY LOGS statement:
// PURGE BINARY LOGS {TO 'log_name' | BEFORE datetime_expr}.
// Before is nil for the TO form.
type PurgeBinaryLogsStmt struct {
	stmtNode

	To     string
	Before ExprNode
}

// Restore implements Node interface.
func (n *PurgeBinaryLogsStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("PURGE BINARY LOGS ")
	if n.Before != nil {
		ctx.WriteKeyWord("BEFORE ")
		if err := n.Before.Restore(ctx); err != nil {
			return annotate(err, "An error occurred while restore PurgeBinaryLogsStmt.Before")
		}
		return nil
	}
	ctx.WriteKeyWord("TO ")
	ctx.WriteString(n.To)
	return nil
}

// Accept implements Node Accept interface.
func (n *PurgeBinaryLogsStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*PurgeBinaryLogsStmt)
	if n.Before != nil {
		node, ok := n.Before.Accept(v)
		if !ok {
			return n, false
		}
		n.Before = node.(ExprNode)
	}
	return v.Leave(n)
}

// ResetBinaryLogsAndGtidsStmt is a RESET BINARY LOGS AND GTIDS
// statement (which replaced RESET MASTER; the removed spelling is not
// parsed). To is the TO binary_log_file_index_number clause; zero when
// absent (MySQL requires the index to be positive).
type ResetBinaryLogsAndGtidsStmt struct {
	stmtNode

	To int64
}

// Restore implements Node interface.
func (n *ResetBinaryLogsAndGtidsStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("RESET BINARY LOGS AND GTIDS")
	if n.To != 0 {
		ctx.WriteKeyWord(" TO ")
		ctx.WritePlainf("%d", n.To)
	}
	return nil
}

// Accept implements Node Accept interface.
func (n *ResetBinaryLogsAndGtidsStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*ResetBinaryLogsAndGtidsStmt)
	return v.Leave(n)
}

// ReplicationFilterType names the filter of one CHANGE REPLICATION
// FILTER specification.
type ReplicationFilterType int

const (
	// ReplicationFilterDoDB is REPLICATE_DO_DB.
	ReplicationFilterDoDB ReplicationFilterType = iota
	// ReplicationFilterIgnoreDB is REPLICATE_IGNORE_DB.
	ReplicationFilterIgnoreDB
	// ReplicationFilterDoTable is REPLICATE_DO_TABLE.
	ReplicationFilterDoTable
	// ReplicationFilterIgnoreTable is REPLICATE_IGNORE_TABLE.
	ReplicationFilterIgnoreTable
	// ReplicationFilterWildDoTable is REPLICATE_WILD_DO_TABLE.
	ReplicationFilterWildDoTable
	// ReplicationFilterWildIgnoreTable is REPLICATE_WILD_IGNORE_TABLE.
	ReplicationFilterWildIgnoreTable
	// ReplicationFilterRewriteDB is REPLICATE_REWRITE_DB.
	ReplicationFilterRewriteDB
)

// String implements fmt.Stringer interface.
func (n ReplicationFilterType) String() string {
	switch n {
	case ReplicationFilterDoDB:
		return "REPLICATE_DO_DB"
	case ReplicationFilterIgnoreDB:
		return "REPLICATE_IGNORE_DB"
	case ReplicationFilterDoTable:
		return "REPLICATE_DO_TABLE"
	case ReplicationFilterIgnoreTable:
		return "REPLICATE_IGNORE_TABLE"
	case ReplicationFilterWildDoTable:
		return "REPLICATE_WILD_DO_TABLE"
	case ReplicationFilterWildIgnoreTable:
		return "REPLICATE_WILD_IGNORE_TABLE"
	case ReplicationFilterRewriteDB:
		return "REPLICATE_REWRITE_DB"
	}
	return ""
}

// ReplicationRewriteDB is one (from_db, to_db) pair of a
// REPLICATE_REWRITE_DB filter.
type ReplicationRewriteDB struct {
	FromDB CIStr
	ToDB   CIStr
}

// ReplicationFilter is one filter specification of a CHANGE REPLICATION
// FILTER statement. Exactly the value field matching Tp is set — DBNames
// for the _DB filters, Tables for the _TABLE filters, Patterns for the
// _WILD_ filters, and Rewrites for REPLICATE_REWRITE_DB — and an empty
// list (an empty value clears the filter) restores as ().
type ReplicationFilter struct {
	Tp       ReplicationFilterType
	DBNames  []CIStr
	Tables   []*TableName
	Patterns []string
	Rewrites []*ReplicationRewriteDB
}

// Restore implements Node interface.
func (n *ReplicationFilter) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord(n.Tp.String())
	ctx.WritePlain(" = (")
	switch n.Tp {
	case ReplicationFilterDoDB, ReplicationFilterIgnoreDB:
		for i, db := range n.DBNames {
			if i != 0 {
				ctx.WritePlain(", ")
			}
			ctx.WriteName(db.O)
		}
	case ReplicationFilterDoTable, ReplicationFilterIgnoreTable:
		for i, table := range n.Tables {
			if i != 0 {
				ctx.WritePlain(", ")
			}
			if err := table.Restore(ctx); err != nil {
				return annotatef(err, "An error occurred while restore ReplicationFilter.Tables[%d]", i)
			}
		}
	case ReplicationFilterWildDoTable, ReplicationFilterWildIgnoreTable:
		for i, pattern := range n.Patterns {
			if i != 0 {
				ctx.WritePlain(", ")
			}
			ctx.WriteString(pattern)
		}
	case ReplicationFilterRewriteDB:
		for i, rewrite := range n.Rewrites {
			if i != 0 {
				ctx.WritePlain(", ")
			}
			ctx.WritePlain("(")
			ctx.WriteName(rewrite.FromDB.O)
			ctx.WritePlain(", ")
			ctx.WriteName(rewrite.ToDB.O)
			ctx.WritePlain(")")
		}
	}
	ctx.WritePlain(")")
	return nil
}

// ChangeReplicationFilterStmt is a CHANGE REPLICATION FILTER statement:
// CHANGE REPLICATION FILTER filter[, filter] ... [FOR CHANNEL channel].
type ChangeReplicationFilterStmt struct {
	stmtNode

	Filters []*ReplicationFilter
	Channel string // FOR CHANNEL clause; empty when absent
}

// Restore implements Node interface.
func (n *ChangeReplicationFilterStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("CHANGE REPLICATION FILTER ")
	for i, f := range n.Filters {
		if i != 0 {
			ctx.WritePlain(", ")
		}
		if err := f.Restore(ctx); err != nil {
			return annotatef(err, "An error occurred while restore ChangeReplicationFilterStmt.Filters[%d]", i)
		}
	}
	if n.Channel != "" {
		ctx.WriteKeyWord(" FOR CHANNEL ")
		ctx.WriteString(n.Channel)
	}
	return nil
}

// Accept implements Node Accept interface.
func (n *ChangeReplicationFilterStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*ChangeReplicationFilterStmt)
	for _, f := range n.Filters {
		for i, table := range f.Tables {
			node, ok := table.Accept(v)
			if !ok {
				return n, false
			}
			f.Tables[i] = node.(*TableName)
		}
	}
	return v.Leave(n)
}

// ResetReplicaStmt is a RESET REPLICA statement:
// RESET REPLICA [ALL] [FOR CHANNEL channel].
type ResetReplicaStmt struct {
	stmtNode

	All     bool
	Channel string // FOR CHANNEL clause; empty when absent
}

// Restore implements Node interface.
func (n *ResetReplicaStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("RESET REPLICA")
	if n.All {
		ctx.WriteKeyWord(" ALL")
	}
	if n.Channel != "" {
		ctx.WriteKeyWord(" FOR CHANNEL ")
		ctx.WriteString(n.Channel)
	}
	return nil
}

// Accept implements Node Accept interface.
func (n *ResetReplicaStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*ResetReplicaStmt)
	return v.Leave(n)
}

// ReplicaThreadType is one thread type of a START/STOP REPLICA
// statement.
type ReplicaThreadType int

const (
	// ReplicaThreadIO is IO_THREAD.
	ReplicaThreadIO ReplicaThreadType = iota
	// ReplicaThreadSQL is SQL_THREAD.
	ReplicaThreadSQL
)

// String implements fmt.Stringer interface.
func (n ReplicaThreadType) String() string {
	switch n {
	case ReplicaThreadIO:
		return "IO_THREAD"
	case ReplicaThreadSQL:
		return "SQL_THREAD"
	}
	return ""
}

func restoreReplicaThreadTypes(ctx *format.RestoreCtx, types []ReplicaThreadType) {
	for i, t := range types {
		if i != 0 {
			ctx.WritePlain(",")
		}
		ctx.WritePlain(" ")
		ctx.WriteKeyWord(t.String())
	}
}

// StartReplicaStmt is a START REPLICA statement:
//
//	START REPLICA [thread_types] [UNTIL until_option]
//	    [connection_options] [FOR CHANNEL channel]
//
// Until and ConnectionOptions reuse the generic name/value option node
// of CHANGE REPLICATION SOURCE TO; the bare UNTIL SQL_AFTER_MTS_GAPS
// option has a nil Value.
type StartReplicaStmt struct {
	stmtNode

	ThreadTypes       []ReplicaThreadType
	Until             []*ReplicationSourceOption
	ConnectionOptions []*ReplicationSourceOption
	Channel           string // FOR CHANNEL clause; empty when absent
}

// Restore implements Node interface.
func (n *StartReplicaStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("START REPLICA")
	restoreReplicaThreadTypes(ctx, n.ThreadTypes)
	for i, opt := range n.Until {
		if i == 0 {
			ctx.WriteKeyWord(" UNTIL ")
		} else {
			ctx.WritePlain(", ")
		}
		if err := opt.Restore(ctx); err != nil {
			return annotatef(err, "An error occurred while restore StartReplicaStmt.Until[%d]", i)
		}
	}
	for i, opt := range n.ConnectionOptions {
		ctx.WritePlain(" ")
		if err := opt.Restore(ctx); err != nil {
			return annotatef(err, "An error occurred while restore StartReplicaStmt.ConnectionOptions[%d]", i)
		}
	}
	if n.Channel != "" {
		ctx.WriteKeyWord(" FOR CHANNEL ")
		ctx.WriteString(n.Channel)
	}
	return nil
}

// SecureText implements SensitiveStatement interface.
func (n *StartReplicaStmt) SecureText() string {
	opts := make([]*ReplicationSourceOption, 0, len(n.ConnectionOptions))
	for _, opt := range n.ConnectionOptions {
		if strings.Contains(opt.Name, "PASSWORD") {
			opt = &ReplicationSourceOption{Name: opt.Name, Value: NewValueExpr("xxxxxx", "", "")}
		}
		opts = append(opts, opt)
	}
	masked := &StartReplicaStmt{
		ThreadTypes:       n.ThreadTypes,
		Until:             n.Until,
		ConnectionOptions: opts,
		Channel:           n.Channel,
	}
	var sb strings.Builder
	_ = masked.Restore(format.NewRestoreCtx(format.DefaultRestoreFlags, &sb))
	return sb.String()
}

// Accept implements Node Accept interface.
func (n *StartReplicaStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*StartReplicaStmt)
	return v.Leave(n)
}

// StopReplicaStmt is a STOP REPLICA statement:
// STOP REPLICA [thread_types] [FOR CHANNEL channel].
type StopReplicaStmt struct {
	stmtNode

	ThreadTypes []ReplicaThreadType
	Channel     string // FOR CHANNEL clause; empty when absent
}

// Restore implements Node interface.
func (n *StopReplicaStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("STOP REPLICA")
	restoreReplicaThreadTypes(ctx, n.ThreadTypes)
	if n.Channel != "" {
		ctx.WriteKeyWord(" FOR CHANNEL ")
		ctx.WriteString(n.Channel)
	}
	return nil
}

// Accept implements Node Accept interface.
func (n *StopReplicaStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*StopReplicaStmt)
	return v.Leave(n)
}

// StartGroupReplicationStmt is a START GROUP_REPLICATION statement:
// START GROUP_REPLICATION [USER='u' [, PASSWORD='p'] [, DEFAULT_AUTH='a']],
// with the credential options reusing the generic name/value option
// node of CHANGE REPLICATION SOURCE TO.
type StartGroupReplicationStmt struct {
	stmtNode

	Options []*ReplicationSourceOption
}

// Restore implements Node interface.
func (n *StartGroupReplicationStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("START GROUP_REPLICATION")
	for i, opt := range n.Options {
		if i == 0 {
			ctx.WritePlain(" ")
		} else {
			ctx.WritePlain(", ")
		}
		if err := opt.Restore(ctx); err != nil {
			return annotatef(err, "An error occurred while restore StartGroupReplicationStmt.Options[%d]", i)
		}
	}
	return nil
}

// SecureText implements SensitiveStatement interface.
func (n *StartGroupReplicationStmt) SecureText() string {
	opts := make([]*ReplicationSourceOption, 0, len(n.Options))
	for _, opt := range n.Options {
		if strings.Contains(opt.Name, "PASSWORD") {
			opt = &ReplicationSourceOption{Name: opt.Name, Value: NewValueExpr("xxxxxx", "", "")}
		}
		opts = append(opts, opt)
	}
	masked := &StartGroupReplicationStmt{Options: opts}
	var sb strings.Builder
	_ = masked.Restore(format.NewRestoreCtx(format.DefaultRestoreFlags, &sb))
	return sb.String()
}

// Accept implements Node Accept interface.
func (n *StartGroupReplicationStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*StartGroupReplicationStmt)
	return v.Leave(n)
}

// StopGroupReplicationStmt is a STOP GROUP_REPLICATION statement.
type StopGroupReplicationStmt struct {
	stmtNode
}

// Restore implements Node interface.
func (n *StopGroupReplicationStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("STOP GROUP_REPLICATION")
	return nil
}

// Accept implements Node Accept interface.
func (n *StopGroupReplicationStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*StopGroupReplicationStmt)
	return v.Leave(n)
}
