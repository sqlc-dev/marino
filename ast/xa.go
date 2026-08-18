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

// The MySQL XA transaction statements (MySQL 26.7 §15.3.8) and the
// backup lock statements LOCK INSTANCE FOR BACKUP / UNLOCK INSTANCE
// (§15.3.5).

var (
	_ Node = &XID{}

	_ StmtNode = &XAStmt{}
	_ StmtNode = &LockInstanceStmt{}
	_ StmtNode = &UnlockInstanceStmt{}
)

// XID is the transaction identifier of an XA statement:
// xid: gtrid [, bqual [, formatID]].
// Gtrid and Bqual are string literals; HasBqual distinguishes an absent
// bqual from an empty one, and HasFormatID requires HasBqual.
type XID struct {
	node

	Gtrid       string
	Bqual       string
	HasBqual    bool
	FormatID    uint64
	HasFormatID bool
}

// Restore implements Node interface.
func (n *XID) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteString(n.Gtrid)
	if n.HasBqual {
		ctx.WritePlain(", ")
		ctx.WriteString(n.Bqual)
		if n.HasFormatID {
			ctx.WritePlainf(", %d", n.FormatID)
		}
	}
	return nil
}

// Accept implements Node Accept interface.
func (n *XID) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*XID)
	return v.Leave(n)
}

// XAOp is the operation of an XAStmt.
type XAOp int

const (
	// XAOpStart is XA START (the XA BEGIN spelling parses to it too).
	XAOpStart XAOp = iota
	// XAOpEnd is XA END.
	XAOpEnd
	// XAOpPrepare is XA PREPARE.
	XAOpPrepare
	// XAOpCommit is XA COMMIT.
	XAOpCommit
	// XAOpRollback is XA ROLLBACK.
	XAOpRollback
	// XAOpRecover is XA RECOVER.
	XAOpRecover
)

// XAStartOption is the JOIN/RESUME modifier of XA START.
type XAStartOption int

const (
	// XAStartNone omits the modifier.
	XAStartNone XAStartOption = iota
	// XAStartJoin is XA START ... JOIN.
	XAStartJoin
	// XAStartResume is XA START ... RESUME.
	XAStartResume
)

// XAStmt is an XA transaction statement:
//
//	XA {START|BEGIN} xid [JOIN|RESUME]
//	XA END xid [SUSPEND [FOR MIGRATE]]
//	XA PREPARE xid
//	XA COMMIT xid [ONE PHASE]
//	XA ROLLBACK xid
//	XA RECOVER [CONVERT XID]
//
// Xid is nil for XA RECOVER only. ForMigrate requires Suspend; both
// belong to XA END. OnePhase belongs to XA COMMIT and ConvertXID to
// XA RECOVER.
type XAStmt struct {
	stmtNode

	Op          XAOp
	Xid         *XID
	StartOption XAStartOption
	Suspend     bool
	ForMigrate  bool
	OnePhase    bool
	ConvertXID  bool
}

// Restore implements Node interface.
func (n *XAStmt) Restore(ctx *format.RestoreCtx) error {
	switch n.Op {
	case XAOpStart:
		ctx.WriteKeyWord("XA START ")
	case XAOpEnd:
		ctx.WriteKeyWord("XA END ")
	case XAOpPrepare:
		ctx.WriteKeyWord("XA PREPARE ")
	case XAOpCommit:
		ctx.WriteKeyWord("XA COMMIT ")
	case XAOpRollback:
		ctx.WriteKeyWord("XA ROLLBACK ")
	case XAOpRecover:
		ctx.WriteKeyWord("XA RECOVER")
		if n.ConvertXID {
			ctx.WriteKeyWord(" CONVERT XID")
		}
		return nil
	}
	if err := n.Xid.Restore(ctx); err != nil {
		return annotate(err, "An error occurred while restore XAStmt.Xid")
	}
	switch n.StartOption {
	case XAStartJoin:
		ctx.WriteKeyWord(" JOIN")
	case XAStartResume:
		ctx.WriteKeyWord(" RESUME")
	}
	if n.Suspend {
		ctx.WriteKeyWord(" SUSPEND")
		if n.ForMigrate {
			ctx.WriteKeyWord(" FOR MIGRATE")
		}
	}
	if n.OnePhase {
		ctx.WriteKeyWord(" ONE PHASE")
	}
	return nil
}

// Accept implements Node Accept interface.
func (n *XAStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*XAStmt)
	if n.Xid != nil {
		node, ok := n.Xid.Accept(v)
		if !ok {
			return n, false
		}
		n.Xid = node.(*XID)
	}
	return v.Leave(n)
}

// LockInstanceStmt is a LOCK INSTANCE FOR BACKUP statement.
type LockInstanceStmt struct {
	stmtNode
}

// Restore implements Node interface.
func (n *LockInstanceStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("LOCK INSTANCE FOR BACKUP")
	return nil
}

// Accept implements Node Accept interface.
func (n *LockInstanceStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*LockInstanceStmt)
	return v.Leave(n)
}

// UnlockInstanceStmt is an UNLOCK INSTANCE statement.
type UnlockInstanceStmt struct {
	stmtNode
}

// Restore implements Node interface.
func (n *UnlockInstanceStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("UNLOCK INSTANCE")
	return nil
}

// Accept implements Node Accept interface.
func (n *UnlockInstanceStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*UnlockInstanceStmt)
	return v.Leave(n)
}
