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

package parser

// Keyword bookkeeping consistency. Before the goyacc removal this test
// cross-checked parser.y's token sections; the sources of truth are now
// the hand-maintained tables themselves: the lexer's tokenMap (misc.go),
// the identifier-position keyword classes (keyword_classes.go), and the
// exported Keywords catalogue (keywords.go). This test keeps the three
// in exact agreement so an edit to one surfaces immediately.

import (
	"testing"
)

func TestKeywordConsistent(t *testing.T) {
	// Aliases map distinct spellings onto the same token.
	for k, v := range aliases {
		if k == v {
			t.Fatalf("alias %q maps to itself", k)
		}
		if tokenMap[k] != tokenMap[v] {
			t.Fatalf("alias %q (%d) and %q (%d) have different tokens", k, tokenMap[k], v, tokenMap[v])
		}
	}

	toSet := func(names []string) map[string]bool {
		s := make(map[string]bool, len(names))
		for _, n := range names {
			if s[n] {
				t.Fatalf("duplicate keyword-class entry %q", n)
			}
			s[n] = true
		}
		return s
	}
	unreserved := toSet(unReservedKeywordNames)
	notKeyword := toSet(notKeywordTokenNames)
	tidb := toSet(tiDBKeywordNames)

	// The classes are disjoint.
	for n := range unreserved {
		if notKeyword[n] || tidb[n] {
			t.Fatalf("%q appears in more than one keyword class", n)
		}
	}
	for n := range notKeyword {
		if tidb[n] {
			t.Fatalf("%q appears in more than one keyword class", n)
		}
	}

	// Partition the Keywords catalogue by section.
	sections := map[string]map[string]bool{}
	for _, kw := range Keywords {
		m := sections[kw.Section]
		if m == nil {
			m = map[string]bool{}
			sections[kw.Section] = m
		}
		if m[kw.Word] {
			t.Fatalf("duplicate Keywords entry %q", kw.Word)
		}
		m[kw.Word] = true
		if kw.Reserved != (kw.Section == "reserved") {
			t.Fatalf("Keywords entry %q: Reserved=%v inconsistent with Section=%q", kw.Word, kw.Reserved, kw.Section)
		}
	}
	reserved := sections["reserved"]

	// Reserved keywords can never appear in identifier position.
	for n := range reserved {
		if unreserved[n] || notKeyword[n] || tidb[n] {
			t.Fatalf("reserved keyword %q also appears in an identifier keyword class", n)
		}
	}

	// The TiDB section matches its class exactly.
	for n := range tidb {
		if !sections["tidb"][n] {
			t.Fatalf("TiDBKeyword %q missing from Keywords", n)
		}
	}
	for n := range sections["tidb"] {
		if !tidb[n] {
			t.Fatalf("Keywords tidb entry %q missing from tiDBKeywordNames", n)
		}
	}

	// The unreserved section matches its class except MONITOR, a
	// MariaDB-gated keyword the historical generator omitted.
	for n := range sections["unreserved"] {
		if !unreserved[n] {
			t.Fatalf("Keywords unreserved entry %q missing from unReservedKeywordNames", n)
		}
	}
	for n := range unreserved {
		if !sections["unreserved"][n] && n != "MONITOR" {
			t.Fatalf("unreserved keyword %q missing from Keywords", n)
		}
	}

	// Every keyword the lexer knows is classified somewhere, and every
	// classified keyword is known to a lexer table (tokenMap, or
	// windowFuncTokenMap for the window-function keywords).
	classified := map[string]bool{}
	for _, set := range []map[string]bool{reserved, unreserved, notKeyword, tidb} {
		for n := range set {
			classified[n] = true
		}
	}
	for n := range tokenMap {
		if !classified[n] && aliases[n] == "" {
			t.Fatalf("tokenMap keyword %q is not classified anywhere", n)
		}
	}
	for n := range classified {
		if _, inTok := tokenMap[n]; !inTok {
			if _, inWin := windowFuncTokenMap[n]; !inWin {
				t.Fatalf("classified keyword %q is in neither tokenMap nor windowFuncTokenMap", n)
			}
		}
	}
}
