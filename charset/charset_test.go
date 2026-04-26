// Copyright 2015 PingCAP, Inc.
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

package charset

import (
	"math/rand"
	"testing"

	"reflect"
)

func testValidCharset(t *testing.T, charset string, collation string, expect bool) {
	b := ValidCharsetAndCollation(charset, collation)
	if !reflect.DeepEqual(expect, b) {
		t.Fatalf("got %v, want %v", b, expect)
	}
}

func TestValidCharset(t *testing.T) {
	tests := []struct {
		cs   string
		co   string
		succ bool
	}{
		{"utf8", "utf8_general_ci", true},
		{"", "utf8_general_ci", true},
		{"utf8mb4", "utf8mb4_bin", true},
		{"latin1", "latin1_bin", true},
		{"utf8", "utf8_invalid_ci", false},
		{"utf16", "utf16_bin", false},
		{"gb2312", "gb2312_chinese_ci", false},
		{"UTF8", "UTF8_BIN", true},
		{"UTF8", "utf8_bin", true},
		{"UTF8MB4", "utf8mb4_bin", true},
		{"UTF8MB4", "UTF8MB4_bin", true},
		{"UTF8MB4", "UTF8MB4_general_ci", true},
		{"Utf8", "uTf8_bIN", true},
		{"utf8mb3", "", true},
		{"utf8mb3", "utf8mb3_bin", true},
		{"utf8mb3", "utf8mb3_general_ci", true},
		{"utf8mb3", "utf8mb3_unicode_ci", true},
	}
	for _, tt := range tests {
		testValidCharset(t, tt.cs, tt.co, tt.succ)
	}
}

func testGetDefaultCollation(t *testing.T, charset string, expectCollation string, succ bool) {
	b, err := GetDefaultCollation(charset)
	if !succ {
		if err == nil {
			t.Fatal("expected error")
		}
		return
	}
	if !reflect.DeepEqual(expectCollation, b) {
		t.Fatalf("got %v, want %v", b, expectCollation)
	}
}

func TestGetDefaultCollation(t *testing.T) {
	tests := []struct {
		cs   string
		co   string
		succ bool
	}{
		{"utf8", "utf8_bin", true},
		{"UTF8", "utf8_bin", true},
		{"utf8mb4", "utf8mb4_bin", true},
		{"ascii", "ascii_bin", true},
		{"binary", "binary", true},
		{"latin1", "latin1_bin", true},
		{"invalid_cs", "", false},
		{"", "utf8_bin", false},
	}
	for _, tt := range tests {
		testGetDefaultCollation(t, tt.cs, tt.co, tt.succ)
	}

	// Test the consistency of collations table and charset desc table
	charsetNum := 0
	for _, collate := range collations {
		if collate.IsDefault {
			if desc, ok := CharacterSetInfos[collate.CharsetName]; ok {
				if !reflect.DeepEqual(desc.DefaultCollation, collate.Name) {
					t.Fatalf("got %v, want %v", collate.Name, desc.DefaultCollation)
				}
				charsetNum++
			}
		}
	}
	if !reflect.DeepEqual(len(CharacterSetInfos), charsetNum) {
		t.Fatalf("got %v, want %v", charsetNum, len(CharacterSetInfos))
	}
}

func TestGetCharsetDesc(t *testing.T) {
	tests := []struct {
		cs     string
		result string
		succ   bool
	}{
		{"utf8", "utf8", true},
		{"UTF8", "utf8", true},
		{"utf8mb4", "utf8mb4", true},
		{"ascii", "ascii", true},
		{"binary", "binary", true},
		{"latin1", "latin1", true},
		{"invalid_cs", "", false},
		{"", "utf8_bin", false},
	}
	for _, tt := range tests {
		desc, err := GetCharsetInfo(tt.cs)
		if !tt.succ {
			if err == nil {
				t.Fatal("expected error")
			}
		} else {
			if !reflect.DeepEqual(tt.result, desc.Name) {
				t.Fatalf("got %v, want %v", desc.Name, tt.result)
			}
		}
	}
}

func TestGetCollationByName(t *testing.T) {
	for _, collation := range collations {
		coll, err := GetCollationByName(collation.Name)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(collation, coll) {
			t.Fatalf("got %v, want %v", coll, collation)
		}
	}

	_, err := GetCollationByName("non_exist")
	if err == nil || err.Error() != "[ddl:1273]Unknown collation: 'non_exist'" {
		t.Fatalf("expected error %q, got %v", "[ddl:1273]Unknown collation: 'non_exist'", err)
	}
}

func TestValidCustomCharset(t *testing.T) {
	AddCharset(&Charset{"custom", "custom_collation", make(map[string]*Collation), "Custom", 4})
	defer RemoveCharset("custom")
	AddCollation(&Collation{99999, "custom", "custom_collation", true, 8, PadNone})

	tests := []struct {
		cs   string
		co   string
		succ bool
		l    int
		p    string
	}{
		{"custom", "custom_collation", true, 8, PadNone},
		{"utf8", "utf8_invalid_ci", false, 1, PadNone},
	}
	for _, tt := range tests {
		testValidCharset(t, tt.cs, tt.co, tt.succ)
	}
}

func TestUTF8MB3(t *testing.T) {
	colname, err := GetDefaultCollationLegacy("utf8mb3")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(colname, "utf8_bin") {
		t.Fatalf("got %v, want %v", "utf8_bin", colname)
	}

	csinfo, err := GetCharsetInfo("utf8mb3")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(csinfo.Name, "utf8") {
		t.Fatalf("got %v, want %v", "utf8", csinfo.Name)
	}

	tests := []struct {
		cs    string
		alias string
	}{
		{"utf8mb3_bin", "utf8_bin"},
		{"utf8mb3_general_ci", "utf8_general_ci"},
		{"utf8mb3_unicode_ci", "utf8_unicode_ci"},
	}
	for _, tt := range tests {
		col, err := GetCollationByName(tt.cs)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(col.Name, tt.alias) {
			t.Fatalf("got %v, want %v", tt.alias, col.Name)
		}
	}
}

func BenchmarkGetCharsetDesc(b *testing.B) {
	b.ResetTimer()
	charsets := []string{CharsetUTF8, CharsetUTF8MB4, CharsetASCII, CharsetLatin1, CharsetBin}
	index := rand.Intn(len(charsets))
	cs := charsets[index]

	for i := 0; i < b.N; i++ {
		GetCharsetInfo(cs)
	}
}
