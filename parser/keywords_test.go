// Copyright 2023 PingCAP, Inc.
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

package parser_test

import (
	"testing"

	"github.com/sqlc-dev/marino/parser"

	"reflect"
)

func TestKeywords(t *testing.T) {
	// Test for the first keyword
	if !reflect.DeepEqual("ADD", parser.Keywords[0].Word) {
		t.Fatalf("got %v, want %v", parser.Keywords[0].Word, "ADD")
	}
	if !reflect.DeepEqual(true, parser.Keywords[0].Reserved) {
		t.Fatalf("got %v, want %v", parser.Keywords[0].Reserved, true)
	}

	// Make sure TiDBKeywords are included.
	found := false
	for _, kw := range parser.Keywords {
		if kw.Word == "ADMIN" {
			found = true
		}
	}
	if !reflect.DeepEqual(found, true) {
		t.Fatalf("%v: got %v, want %v", "TiDBKeyword ADMIN is part of the list", true, found)
	}
}

func TestKeywordsLength(t *testing.T) {
	if !reflect.DeepEqual(679, len(parser.Keywords)) {
		t.Fatalf("got %v, want %v", len(parser.Keywords), 679)
	}

	reservedNr := 0
	for _, kw := range parser.Keywords {
		if kw.Reserved {
			reservedNr += 1
		}
	}
	if !reflect.DeepEqual(233, reservedNr) {
		t.Fatalf("got %v, want %v", reservedNr, 233)
	}
}

func TestKeywordsSorting(t *testing.T) {
	for i, kw := range parser.Keywords {
		if i > 1 && parser.Keywords[i-1].Word > kw.Word && parser.Keywords[i-1].Section == kw.Section {
			t.Errorf("%s should come after %s, please update parser.y and re-generate keywords.go\n",
				parser.Keywords[i-1].Word, kw.Word)
		}
	}
}
