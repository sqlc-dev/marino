// Copyright 2017 PingCAP, Inc.
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

package ast_test

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/sqlc-dev/marino/ast"
	. "github.com/sqlc-dev/marino/format"
	"github.com/sqlc-dev/marino/mysql"
	"github.com/sqlc-dev/marino/parser"

	"reflect"
)

func TestCacheable(t *testing.T) {
	// test non-SelectStmt
	var stmt Node = &DeleteStmt{}
	if IsReadOnly(stmt, true) {
		t.Fatal("expected false")
	}

	stmt = &InsertStmt{}
	if IsReadOnly(stmt, true) {
		t.Fatal("expected false")
	}

	stmt = &UpdateStmt{}
	if IsReadOnly(stmt, true) {
		t.Fatal("expected false")
	}

	stmt = &ExplainStmt{}
	if !(IsReadOnly(stmt, true)) {
		t.Fatal("expected true")
	}

	stmt = &ExplainStmt{}
	if !(IsReadOnly(stmt, true)) {
		t.Fatal("expected true")
	}

	stmt = &DoStmt{}
	if !(IsReadOnly(stmt, true)) {
		t.Fatal("expected true")
	}

	stmt = &ExplainStmt{
		Stmt: &InsertStmt{},
	}
	if !(IsReadOnly(stmt, true)) {
		t.Fatal("expected true")
	}

	stmt = &ExplainStmt{
		Analyze: true,
		Stmt:    &InsertStmt{},
	}
	if IsReadOnly(stmt, true) {
		t.Fatal("expected false")
	}

	stmt = &ExplainStmt{
		Stmt: &SelectStmt{},
	}
	if !(IsReadOnly(stmt, true)) {
		t.Fatal("expected true")
	}

	stmt = &ExplainStmt{
		Analyze: true,
		Stmt:    &SelectStmt{},
	}
	if !(IsReadOnly(stmt, true)) {
		t.Fatal("expected true")
	}

	stmt = &ShowStmt{}
	if !(IsReadOnly(stmt, true)) {
		t.Fatal("expected true")
	}

	stmt = &ShowStmt{}
	if !(IsReadOnly(stmt, true)) {
		t.Fatal("expected true")
	}

	stmt = &TraceStmt{
		Stmt: &SelectStmt{},
	}
	if !(IsReadOnly(stmt, true)) {
		t.Fatal("expected true")
	}

	stmt = &TraceStmt{
		Stmt: &DeleteStmt{},
	}
	if IsReadOnly(stmt, true) {
		t.Fatal("expected false")
	}
}

func TestUnionReadOnly(t *testing.T) {
	selectReadOnly := &SelectStmt{}
	selectForUpdate := &SelectStmt{
		LockInfo: &SelectLockInfo{LockType: SelectLockForUpdate},
	}
	selectForUpdateNoWait := &SelectStmt{
		LockInfo: &SelectLockInfo{LockType: SelectLockForUpdateNoWait},
	}

	setOprStmt := &SetOprStmt{
		SelectList: &SetOprSelectList{
			Selects: []Node{selectReadOnly, selectReadOnly},
		},
	}
	if !(IsReadOnly(setOprStmt, true)) {
		t.Fatal("expected true")
	}

	setOprStmt.SelectList.Selects = []Node{selectReadOnly, selectReadOnly, selectReadOnly}
	if !(IsReadOnly(setOprStmt, true)) {
		t.Fatal("expected true")
	}

	setOprStmt.SelectList.Selects = []Node{selectReadOnly, selectForUpdate}
	if IsReadOnly(setOprStmt, true) {
		t.Fatal("expected false")
	}

	setOprStmt.SelectList.Selects = []Node{selectReadOnly, selectForUpdateNoWait}
	if IsReadOnly(setOprStmt, true) {
		t.Fatal("expected false")
	}

	setOprStmt.SelectList.Selects = []Node{selectForUpdate, selectForUpdateNoWait}
	if IsReadOnly(setOprStmt, true) {
		t.Fatal("expected false")
	}

	setOprStmt.SelectList.Selects = []Node{selectReadOnly, selectForUpdate, selectForUpdateNoWait}
	if IsReadOnly(setOprStmt, true) {
		t.Fatal("expected false")
	}
}

// CleanNodeText set the text of node and all child node empty.
// For test only.
func CleanNodeText(node Node) {
	var cleaner nodeTextCleaner
	node.Accept(&cleaner)
}

// nodeTextCleaner clean the text of a node and it's child node.
// For test only.
type nodeTextCleaner struct {
}

// Enter implements Visitor interface.
func (checker *nodeTextCleaner) Enter(in Node) (out Node, skipChildren bool) {
	in.SetText(nil, "")
	in.SetOriginTextPosition(0)
	if v, ok := in.(ValueExpr); ok && v != nil {
		tpFlag := v.GetType().GetFlag()
		if tpFlag&mysql.UnderScoreCharsetFlag != 0 {
			// ignore underscore charset flag to let `'abc' = _utf8'abc'` pass
			tpFlag ^= mysql.UnderScoreCharsetFlag
			v.GetType().SetFlag(tpFlag)
		}
	}

	switch node := in.(type) {
	case *Constraint:
		if node.Option != nil {
			if node.Option.KeyBlockSize == 0x0 && node.Option.Tp == 0 && node.Option.Comment == "" {
				node.Option = nil
			}
		}
	case *FuncCallExpr:
		node.FnName.O = strings.ToLower(node.FnName.O)
		switch node.FnName.L {
		case "convert":
			node.Args[1].(*ValueExprBase).Datum.SetBytes(nil)
		}
	case *AggregateFuncExpr:
		node.F = strings.ToLower(node.F)
	case *FieldList:
		for _, f := range node.Fields {
			f.Offset = 0
		}
	case *AlterTableSpec:
		for _, opt := range node.Options {
			opt.StrValue = strings.ToLower(opt.StrValue)
		}
	case *Join:
		node.ExplicitParens = false
	case *ColumnDef:
		node.Tp.CleanElemIsBinaryLit()
	}
	return in, false
}

// Leave implements Visitor interface.
func (checker *nodeTextCleaner) Leave(in Node) (out Node, ok bool) {
	return in, true
}

type NodeRestoreTestCase struct {
	sourceSQL string
	expectSQL string
}

func runNodeRestoreTest(t *testing.T, nodeTestCases []NodeRestoreTestCase, template string, extractNodeFunc func(node Node) Node) {
	runNodeRestoreTestWithFlags(t, nodeTestCases, template, extractNodeFunc, DefaultRestoreFlags)
}

func runNodeRestoreTestWithFlags(t *testing.T, nodeTestCases []NodeRestoreTestCase, template string, extractNodeFunc func(node Node) Node, flags RestoreFlags) {
	p := parser.New()
	p.EnableWindowFunc(true)
	for _, testCase := range nodeTestCases {
		sourceSQL := fmt.Sprintf(template, testCase.sourceSQL)
		expectSQL := fmt.Sprintf(template, testCase.expectSQL)
		stmt, err := p.ParseOneStmt(sourceSQL, "", "")
		comment := fmt.Sprintf("source %#v", testCase)
		if err != nil {
			t.Fatalf("%v: %v", comment, err)
		}
		var sb strings.Builder
		err = extractNodeFunc(stmt).Restore(NewRestoreCtx(flags, &sb))
		if err != nil {
			t.Fatalf("%v: %v", comment, err)
		}
		restoreSql := fmt.Sprintf(template, sb.String())
		comment = fmt.Sprintf("source %#v; restore %v", testCase, restoreSql)
		if !reflect.DeepEqual(expectSQL, restoreSql) {
			t.Fatalf("%v: got %v, want %v", comment, restoreSql, expectSQL)
		}
		stmt2, err := p.ParseOneStmt(restoreSql, "", "")
		if err != nil {
			t.Fatalf("%v: %v", comment, err)
		}
		CleanNodeText(stmt)
		CleanNodeText(stmt2)
		if !reflect.DeepEqual(stmt, stmt2) {
			t.Fatalf("%v: got %v, want %v", comment, stmt2, stmt)
		}
	}
}

// runNodeRestoreTestWithFlagsStmtChange likes runNodeRestoreTestWithFlags but not check if the ASTs are same.
// Sometimes the AST are different and it's expected.
func runNodeRestoreTestWithFlagsStmtChange(t *testing.T, nodeTestCases []NodeRestoreTestCase, template string, extractNodeFunc func(node Node) Node, flags RestoreFlags) {
	p := parser.New()
	p.EnableWindowFunc(true)
	for _, testCase := range nodeTestCases {
		sourceSQL := fmt.Sprintf(template, testCase.sourceSQL)
		expectSQL := fmt.Sprintf(template, testCase.expectSQL)
		stmt, err := p.ParseOneStmt(sourceSQL, "", "")
		comment := fmt.Sprintf("source %#v", testCase)
		if err != nil {
			t.Fatalf("%v: %v", comment, err)
		}
		var sb strings.Builder
		err = extractNodeFunc(stmt).Restore(NewRestoreCtx(flags, &sb))
		if err != nil {
			t.Fatalf("%v: %v", comment, err)
		}
		restoreSql := fmt.Sprintf(template, sb.String())
		comment = fmt.Sprintf("source %#v; restore %v", testCase, restoreSql)
		if !reflect.DeepEqual(expectSQL, restoreSql) {
			t.Fatalf("%v: got %v, want %v", comment, restoreSql, expectSQL)
		}
	}
}
