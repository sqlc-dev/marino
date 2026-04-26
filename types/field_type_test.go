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

package types_test

import (
	"fmt"
	"testing"

	"github.com/sqlc-dev/marino/ast"
	"github.com/sqlc-dev/marino/charset"
	"github.com/sqlc-dev/marino/mysql"
	"github.com/sqlc-dev/marino/parser"

	// import parser_driver
	_ "github.com/sqlc-dev/marino/test_driver"
	. "github.com/sqlc-dev/marino/types"

	"reflect"
)

func TestFieldType(t *testing.T) {
	ft := NewFieldType(mysql.TypeDuration)
	if !reflect.DeepEqual(UnspecifiedLength, ft.GetFlen()) {
		t.Fatalf("got %v, want %v", ft.GetFlen(), UnspecifiedLength)
	}
	if !reflect.DeepEqual(UnspecifiedLength, ft.GetDecimal()) {
		t.Fatalf("got %v, want %v", ft.GetDecimal(), UnspecifiedLength)
	}
	ft.SetDecimal(5)
	if !reflect.DeepEqual("time(5)", ft.String()) {
		t.Fatalf("got %v, want %v", ft.String(), "time(5)")
	}
	if HasCharset(ft) {
		t.Fatal("expected false")
	}

	ft = NewFieldType(mysql.TypeLong)
	ft.SetFlen(5)
	ft.SetFlag(mysql.UnsignedFlag | mysql.ZerofillFlag)
	if !reflect.DeepEqual("int(5) UNSIGNED ZEROFILL", ft.String()) {
		t.Fatalf("got %v, want %v", ft.String(), "int(5) UNSIGNED ZEROFILL")
	}
	if !reflect.DeepEqual("int(5) unsigned", ft.InfoSchemaStr()) {
		t.Fatalf("got %v, want %v", ft.InfoSchemaStr(), "int(5) unsigned")
	}
	if HasCharset(ft) {
		t.Fatal("expected false")
	}

	ft = NewFieldType(mysql.TypeFloat)
	ft.SetFlen(12)   // Default
	ft.SetDecimal(3) // Not Default
	if !reflect.DeepEqual("float(12,3)", ft.String()) {
		t.Fatalf("got %v, want %v", ft.String(), "float(12,3)")
	}
	ft = NewFieldType(mysql.TypeFloat)
	ft.SetFlen(12)    // Default
	ft.SetDecimal(-1) // Default
	if !reflect.DeepEqual("float", ft.String()) {
		t.Fatalf("got %v, want %v", ft.String(), "float")
	}
	ft = NewFieldType(mysql.TypeFloat)
	ft.SetFlen(5)     // Not Default
	ft.SetDecimal(-1) // Default
	if !reflect.DeepEqual("float", ft.String()) {
		t.Fatalf("got %v, want %v", ft.String(), "float")
	}
	ft = NewFieldType(mysql.TypeFloat)
	ft.SetFlen(7)    // Not Default
	ft.SetDecimal(3) // Not Default
	if !reflect.DeepEqual("float(7,3)", ft.String()) {
		t.Fatalf("got %v, want %v", ft.String(), "float(7,3)")
	}
	if HasCharset(ft) {
		t.Fatal("expected false")
	}

	ft = NewFieldType(mysql.TypeDouble)
	ft.SetFlen(22)   // Default
	ft.SetDecimal(3) // Not Default
	if !reflect.DeepEqual("double(22,3)", ft.String()) {
		t.Fatalf("got %v, want %v", ft.String(), "double(22,3)")
	}
	ft = NewFieldType(mysql.TypeDouble)
	ft.SetFlen(22)    // Default
	ft.SetDecimal(-1) // Default
	if !reflect.DeepEqual("double", ft.String()) {
		t.Fatalf("got %v, want %v", ft.String(), "double")
	}
	ft = NewFieldType(mysql.TypeDouble)
	ft.SetFlen(5)     // Not Default
	ft.SetDecimal(-1) // Default
	if !reflect.DeepEqual("double", ft.String()) {
		t.Fatalf("got %v, want %v", ft.String(), "double")
	}
	ft = NewFieldType(mysql.TypeDouble)
	ft.SetFlen(7)    // Not Default
	ft.SetDecimal(3) // Not Default
	if !reflect.DeepEqual("double(7,3)", ft.String()) {
		t.Fatalf("got %v, want %v", ft.String(), "double(7,3)")
	}
	if HasCharset(ft) {
		t.Fatal("expected false")
	}

	ft = NewFieldType(mysql.TypeBlob)
	ft.SetFlen(10)
	ft.SetCharset("UTF8")
	ft.SetCollate("UTF8_UNICODE_GI")
	if !reflect.DeepEqual("text CHARACTER SET UTF8 COLLATE UTF8_UNICODE_GI", ft.String()) {
		t.Fatalf("got %v, want %v", ft.String(), "text CHARACTER SET UTF8 COLLATE UTF8_UNICODE_GI")
	}
	if !(HasCharset(ft)) {
		t.Fatal("expected true")
	}

	ft = NewFieldType(mysql.TypeVarchar)
	ft.SetFlen(10)
	ft.AddFlag(mysql.BinaryFlag)
	if !reflect.DeepEqual("varchar(10) BINARY", ft.String()) {
		t.Fatalf("got %v, want %v", ft.String(), "varchar(10) BINARY")
	}
	if HasCharset(ft) {
		t.Fatal("expected false")
	}

	ft = NewFieldType(mysql.TypeString)
	ft.SetCharset(charset.CharsetBin)
	ft.AddFlag(mysql.BinaryFlag)
	if !reflect.DeepEqual("binary(1)", ft.String()) {
		t.Fatalf("got %v, want %v", ft.String(), "binary(1)")
	}
	if HasCharset(ft) {
		t.Fatal("expected false")
	}

	ft = NewFieldType(mysql.TypeEnum)
	ft.SetElems([]string{"a", "b"})
	if !reflect.DeepEqual("enum('a','b')", ft.String()) {
		t.Fatalf("got %v, want %v", ft.String(), "enum('a','b')")
	}
	if !(HasCharset(ft)) {
		t.Fatal("expected true")
	}

	ft = NewFieldType(mysql.TypeEnum)
	ft.SetElems([]string{"'a'", "'b'"})
	if !reflect.DeepEqual("enum('''a''','''b''')", ft.String()) {
		t.Fatalf("got %v, want %v", ft.String(), "enum('''a''','''b''')")
	}
	if !(HasCharset(ft)) {
		t.Fatal("expected true")
	}

	ft = NewFieldType(mysql.TypeEnum)
	ft.SetElems([]string{"a\nb", "a\tb", "a\rb"})
	if !reflect.DeepEqual("enum('a\\nb','a\tb','a\\rb')", ft.String()) {
		t.Fatalf("got %v, want %v", ft.String(), "enum('a\\nb','a\tb','a\\rb')")
	}
	if !(HasCharset(ft)) {
		t.Fatal("expected true")
	}

	ft = NewFieldType(mysql.TypeEnum)
	ft.SetElems([]string{"a\nb", "a'\t\r\nb", "a\rb"})
	if !reflect.DeepEqual("enum('a\\nb','a''	\\r\\nb','a\\rb')", ft.String()) {
		t.Fatalf("got %v, want %v", ft.String(), "enum('a\\nb','a''	\\r\\nb','a\\rb')")
	}
	if !(HasCharset(ft)) {
		t.Fatal("expected true")
	}

	ft = NewFieldType(mysql.TypeSet)
	ft.SetElems([]string{"a", "b"})
	if !reflect.DeepEqual("set('a','b')", ft.String()) {
		t.Fatalf("got %v, want %v", ft.String(), "set('a','b')")
	}
	if !(HasCharset(ft)) {
		t.Fatal("expected true")
	}

	ft = NewFieldType(mysql.TypeSet)
	ft.SetElems([]string{"'a'", "'b'"})
	if !reflect.DeepEqual("set('''a''','''b''')", ft.String()) {
		t.Fatalf("got %v, want %v", ft.String(), "set('''a''','''b''')")
	}
	if !(HasCharset(ft)) {
		t.Fatal("expected true")
	}

	ft = NewFieldType(mysql.TypeSet)
	ft.SetElems([]string{"a\nb", "a'\t\r\nb", "a\rb"})
	if !reflect.DeepEqual("set('a\\nb','a''	\\r\\nb','a\\rb')", ft.String()) {
		t.Fatalf("got %v, want %v", ft.String(), "set('a\\nb','a''	\\r\\nb','a\\rb')")
	}
	if !(HasCharset(ft)) {
		t.Fatal("expected true")
	}

	ft = NewFieldType(mysql.TypeSet)
	ft.SetElems([]string{"a'\nb", "a'b\tc"})
	if !reflect.DeepEqual("set('a''\\nb','a''b	c')", ft.String()) {
		t.Fatalf("got %v, want %v", ft.String(), "set('a''\\nb','a''b	c')")
	}
	if !(HasCharset(ft)) {
		t.Fatal("expected true")
	}

	ft = NewFieldType(mysql.TypeTimestamp)
	ft.SetFlen(8)
	ft.SetDecimal(2)
	if !reflect.DeepEqual("timestamp(2)", ft.String()) {
		t.Fatalf("got %v, want %v", ft.String(), "timestamp(2)")
	}
	if HasCharset(ft) {
		t.Fatal("expected false")
	}
	ft = NewFieldType(mysql.TypeTimestamp)
	ft.SetFlen(8)
	ft.SetDecimal(0)
	if !reflect.DeepEqual("timestamp", ft.String()) {
		t.Fatalf("got %v, want %v", ft.String(), "timestamp")
	}
	if HasCharset(ft) {
		t.Fatal("expected false")
	}

	ft = NewFieldType(mysql.TypeDatetime)
	ft.SetFlen(8)
	ft.SetDecimal(2)
	if !reflect.DeepEqual("datetime(2)", ft.String()) {
		t.Fatalf("got %v, want %v", ft.String(), "datetime(2)")
	}
	if HasCharset(ft) {
		t.Fatal("expected false")
	}
	ft = NewFieldType(mysql.TypeDatetime)
	ft.SetFlen(8)
	ft.SetDecimal(0)
	if !reflect.DeepEqual("datetime", ft.String()) {
		t.Fatalf("got %v, want %v", ft.String(), "datetime")
	}
	if HasCharset(ft) {
		t.Fatal("expected false")
	}

	ft = NewFieldType(mysql.TypeDate)
	ft.SetFlen(8)
	ft.SetDecimal(2)
	if !reflect.DeepEqual("date", ft.String()) {
		t.Fatalf("got %v, want %v", ft.String(), "date")
	}
	if HasCharset(ft) {
		t.Fatal("expected false")
	}
	ft = NewFieldType(mysql.TypeDate)
	ft.SetFlen(8)
	ft.SetDecimal(0)
	if !reflect.DeepEqual("date", ft.String()) {
		t.Fatalf("got %v, want %v", ft.String(), "date")
	}
	if HasCharset(ft) {
		t.Fatal("expected false")
	}

	ft = NewFieldType(mysql.TypeYear)
	ft.SetFlen(4)
	ft.SetDecimal(0)
	if !reflect.DeepEqual("year(4)", ft.String()) {
		t.Fatalf("got %v, want %v", ft.String(), "year(4)")
	}
	if HasCharset(ft) {
		t.Fatal("expected false")
	}
	ft = NewFieldType(mysql.TypeYear)
	ft.SetFlen(2)
	ft.SetDecimal(2)
	if !reflect.DeepEqual("year(2)", ft.String()) {
		t.Fatalf("got %v, want %v", ft.String(), "year(2)")
	}
	if HasCharset(ft) {
		t.Fatal("expected false")
	}

	ft = NewFieldType(mysql.TypeVarchar)
	ft.SetFlen(0)
	ft.SetDecimal(0)
	if !reflect.DeepEqual("varchar(0)", ft.String()) {
		t.Fatalf("got %v, want %v", ft.String(), "varchar(0)")
	}
	if !(HasCharset(ft)) {
		t.Fatal("expected true")
	}

	ft = NewFieldType(mysql.TypeString)
	ft.SetFlen(0)
	ft.SetDecimal(0)
	if !reflect.DeepEqual("char(0)", ft.String()) {
		t.Fatalf("got %v, want %v", ft.String(), "char(0)")
	}
	if !(HasCharset(ft)) {
		t.Fatal("expected true")
	}
}

func TestHasCharsetFromStmt(t *testing.T) {
	template := "CREATE TABLE t(a %s)"

	types := []struct {
		strType    string
		hasCharset bool
	}{
		{"int", false},
		{"real", false},
		{"float", false},
		{"bit", false},
		{"bool", false},
		{"char(1)", true},
		{"national char(1)", true},
		{"binary", false},
		{"varchar(1)", true},
		{"national varchar(1)", true},
		{"varbinary(1)", false},
		{"year", false},
		{"date", false},
		{"time", false},
		{"datetime", false},
		{"timestamp", false},
		{"blob", false},
		{"tinyblob", false},
		{"mediumblob", false},
		{"longblob", false},
		{"bit", false},
		{"text", true},
		{"tinytext", true},
		{"mediumtext", true},
		{"longtext", true},
		{"json", false},
		{"enum('1')", true},
		{"set('1')", true},
	}

	p := parser.New()
	for _, typ := range types {
		sql := fmt.Sprintf(template, typ.strType)
		stmt, err := p.ParseOneStmt(sql, "", "")
		if err != nil {
			t.Fatal(err)
		}

		col := stmt.(*ast.CreateTableStmt).Cols[0]
		if !reflect.DeepEqual(typ.hasCharset, HasCharset(col.Tp)) {
			t.Fatalf("got %v, want %v", HasCharset(col.Tp), typ.hasCharset)
		}
	}
}

func TestEnumSetFlen(t *testing.T) {
	p := parser.New()
	cases := []struct {
		sql string
		ex  int
	}{
		{"enum('a')", 1},
		{"enum('a', 'b')", 1},
		{"enum('a', 'bb')", 2},
		{"enum('a', 'b', 'c')", 1},
		{"enum('a', 'bb', 'c')", 2},
		{"enum('a', 'bb', 'c')", 2},
		{"enum('')", 0},
		{"enum('a', '')", 1},
		{"set('a')", 1},
		{"set('a', 'b')", 3},
		{"set('a', 'bb')", 4},
		{"set('a', 'b', 'c')", 5},
		{"set('a', 'bb', 'c')", 6},
		{"set('')", 0},
		{"set('a', '')", 2},
	}

	for _, ca := range cases {
		stmt, err := p.ParseOneStmt(fmt.Sprintf("create table t (e %v)", ca.sql), "", "")
		if err != nil {
			t.Fatal(err)
		}
		col := stmt.(*ast.CreateTableStmt).Cols[0]
		if !reflect.DeepEqual(ca.ex, col.Tp.GetFlen()) {
			t.Fatalf("got %v, want %v", col.Tp.GetFlen(), ca.ex)
		}
	}
}

func TestFieldTypeEqual(t *testing.T) {
	// tp not equal
	ft1 := NewFieldType(mysql.TypeDouble)
	ft2 := NewFieldType(mysql.TypeFloat)
	if !reflect.DeepEqual(false, ft1.Equal(ft2)) {
		t.Fatalf("got %v, want %v", ft1.Equal(ft2), false)
	}

	// decimal not equal
	ft2 = NewFieldType(mysql.TypeDouble)
	ft2.SetDecimal(5)
	if !reflect.DeepEqual(false, ft1.Equal(ft2)) {
		t.Fatalf("got %v, want %v", ft1.Equal(ft2), false)
	}

	// flen not equal and decimal not -1
	ft1.SetDecimal(5)
	ft1.SetFlen(22)
	if !reflect.DeepEqual(false, ft1.Equal(ft2)) {
		t.Fatalf("got %v, want %v", ft1.Equal(ft2), false)
	}

	// flen equal
	ft2.SetFlen(22)
	if !reflect.DeepEqual(true, ft1.Equal(ft2)) {
		t.Fatalf("got %v, want %v", ft1.Equal(ft2), true)
	}

	// decimal is -1
	ft1.SetDecimal(-1)
	ft2.SetDecimal(-1)
	ft1.SetFlen(23)
	if !reflect.DeepEqual(true, ft1.Equal(ft2)) {
		t.Fatalf("got %v, want %v", ft1.Equal(ft2), true)
	}
}

func TestCompactStr(t *testing.T) {
	cases := []struct {
		t     byte   // Field Type
		flen  int    // Field Length
		flags uint   // Field Flags, e.g. ZEROFILL
		e1    string // Expected string with TiDBStrictIntegerDisplayWidth disabled
		e2    string // Expected string with TiDBStrictIntegerDisplayWidth enabled
	}{
		// TINYINT(1) is considered a bool by connectors, this should always display
		// the display length.
		{mysql.TypeTiny, 1, 0, `tinyint(1)`, `tinyint(1)`},
		{mysql.TypeTiny, 2, 0, `tinyint(2)`, `tinyint`},

		// If the ZEROFILL flag is set the display length should not be hidden.
		{mysql.TypeLong, 10, 0, `int(10)`, `int`},
		{mysql.TypeLong, 10, mysql.ZerofillFlag, `int(10)`, `int(10)`},
	}
	for _, cc := range cases {
		ft := NewFieldType(cc.t)
		ft.SetFlen(cc.flen)
		ft.SetFlag(cc.flags)

		TiDBStrictIntegerDisplayWidth = false
		if !reflect.DeepEqual(cc.e1, ft.CompactStr()) {
			t.Fatalf("got %v, want %v", ft.CompactStr(), cc.e1)
		}

		TiDBStrictIntegerDisplayWidth = true
		if !reflect.DeepEqual(cc.e2, ft.CompactStr()) {
			t.Fatalf("got %v, want %v", ft.CompactStr(), cc.e2)
		}
	}
}
