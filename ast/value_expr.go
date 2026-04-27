// Copyright 2019 PingCAP, Inc.
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

//go:build !codes

package ast

import (
	"fmt"
	"io"
	"strconv"

	"github.com/sqlc-dev/marino/charset"
	"github.com/sqlc-dev/marino/format"
	"github.com/sqlc-dev/marino/mysql"
)

var (
	_ ParamMarkerExpr = &ParamMarkerExprBase{}
	_ ValueExpr       = &ValueExprBase{}
)

// ValueExprBase is the simple value expression.
type ValueExprBase struct {
	exprNode
	Datum
	projectionOffset int
}

// Restore implements Node interface.
func (n *ValueExprBase) Restore(ctx *format.RestoreCtx) error {
	switch n.Kind() {
	case KindNull:
		ctx.WriteKeyWord("NULL")
	case KindInt64:
		if n.Type.GetFlag()&mysql.IsBooleanFlag != 0 {
			if n.GetInt64() > 0 {
				ctx.WriteKeyWord("TRUE")
			} else {
				ctx.WriteKeyWord("FALSE")
			}
		} else {
			ctx.WritePlain(strconv.FormatInt(n.GetInt64(), 10))
		}
	case KindUint64:
		ctx.WritePlain(strconv.FormatUint(n.GetUint64(), 10))
	case KindFloat32:
		ctx.WritePlain(strconv.FormatFloat(n.GetFloat64(), 'e', -1, 32))
	case KindFloat64:
		ctx.WritePlain(strconv.FormatFloat(n.GetFloat64(), 'e', -1, 64))
	case KindString:
		if n.Type.GetCharset() != "" &&
			!ctx.Flags.HasStringWithoutCharset() &&
			(!ctx.Flags.HasStringWithoutDefaultCharset() || n.Type.GetCharset() != mysql.DefaultCharset) {
			ctx.WritePlain("_")
			ctx.WriteKeyWord(n.Type.GetCharset())
		}
		ctx.WriteString(n.GetString())
	case KindBytes:
		ctx.WriteString(n.GetString())
	case KindMysqlDecimal:
		ctx.WritePlain(n.GetMysqlDecimal().String())
	case KindBinaryLiteral:
		if n.Type.GetCharset() != "" && n.Type.GetCharset() != mysql.DefaultCharset &&
			!ctx.Flags.HasStringWithoutCharset() &&
			n.Type.GetCharset() != charset.CharsetBin {
			ctx.WritePlain("_")
			ctx.WriteKeyWord(n.Type.GetCharset() + " ")
		}
		if n.Type.GetFlag()&mysql.UnsignedFlag != 0 {
			ctx.WritePlainf("x'%x'", n.GetBytes())
		} else {
			ctx.WritePlain(n.GetBinaryLiteral().ToBitLiteralString(true))
		}
	case KindMysqlDuration, KindMysqlEnum,
		KindMysqlBit, KindMysqlSet, KindMysqlTime,
		KindInterface, KindMinNotNull, KindMaxValue,
		KindRaw, KindMysqlJSON:
		// TODO implement Restore function
		return fmt.Errorf("not implemented")
	default:
		return fmt.Errorf("can't format to string")
	}
	return nil
}

// GetDatumString implements the ValueExpr interface.
func (n *ValueExprBase) GetDatumString() string {
	return n.GetString()
}

// Format the ExprNode into a Writer.
func (n *ValueExprBase) Format(w io.Writer) {
	var s string
	switch n.Kind() {
	case KindNull:
		s = "NULL"
	case KindInt64:
		if n.Type.GetFlag()&mysql.IsBooleanFlag != 0 {
			if n.GetInt64() > 0 {
				s = "TRUE"
			} else {
				s = "FALSE"
			}
		} else {
			s = strconv.FormatInt(n.GetInt64(), 10)
		}
	case KindUint64:
		s = strconv.FormatUint(n.GetUint64(), 10)
	case KindFloat32:
		s = strconv.FormatFloat(n.GetFloat64(), 'e', -1, 32)
	case KindFloat64:
		s = strconv.FormatFloat(n.GetFloat64(), 'e', -1, 64)
	case KindString, KindBytes:
		s = strconv.Quote(n.GetString())
	case KindMysqlDecimal:
		s = n.GetMysqlDecimal().String()
	case KindBinaryLiteral:
		if n.Type.GetFlag()&mysql.UnsignedFlag != 0 {
			s = fmt.Sprintf("x'%x'", n.GetBytes())
		} else {
			s = n.GetBinaryLiteral().ToBitLiteralString(true)
		}
	default:
		panic("Can't format to string")
	}
	_, _ = fmt.Fprint(w, s)
}

// NewValueExpr creates a ValueExpr with value, and sets default field type.
func NewValueExpr(value any, charset string, collate string) ValueExpr {
	if ve, ok := value.(*ValueExprBase); ok {
		return ve
	}
	ve := &ValueExprBase{}
	ve.SetValue(value)
	DefaultTypeForValue(value, &ve.Type, charset, collate)
	ve.projectionOffset = -1
	return ve
}

// SetProjectionOffset sets ValueExprBase.projectionOffset for logical plan builder.
func (n *ValueExprBase) SetProjectionOffset(offset int) {
	n.projectionOffset = offset
}

// GetProjectionOffset returns ValueExprBase.projectionOffset.
func (n *ValueExprBase) GetProjectionOffset() int {
	return n.projectionOffset
}

// Accept implements Node interface.
func (n *ValueExprBase) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*ValueExprBase)
	return v.Leave(n)
}

// ParamMarkerExprBase expression holds a place for another expression.
// Used in parsing prepare statement.
type ParamMarkerExprBase struct {
	ValueExprBase
	Offset    int
	Order     int
	InExecute bool
}

// Restore implements Node interface.
func (n *ParamMarkerExprBase) Restore(ctx *format.RestoreCtx) error {
	ctx.WritePlain("?")
	return nil
}

// NewParamMarkerExpr creates a ParamMarkerExpr.
func NewParamMarkerExpr(offset int) ParamMarkerExpr {
	return &ParamMarkerExprBase{
		Offset: offset,
	}
}

// Format the ExprNode into a Writer.
func (n *ParamMarkerExprBase) Format(w io.Writer) {
	panic("Not implemented")
}

// Accept implements Node Accept interface.
func (n *ParamMarkerExprBase) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	n = newNode.(*ParamMarkerExprBase)
	return v.Leave(n)
}

// SetOrder implements the ParamMarkerExpr interface.
func (n *ParamMarkerExprBase) SetOrder(order int) {
	n.Order = order
}

// NewDecimal creates a *MyDecimal value.
func NewDecimal(s string) (*MyDecimal, error) {
	dec := new(MyDecimal)
	err := dec.FromString([]byte(s))
	return dec, err
}
