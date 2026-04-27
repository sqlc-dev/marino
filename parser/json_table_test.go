// Copyright 2026 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package parser_test

import (
	"strings"
	"testing"

	"github.com/sqlc-dev/marino/ast"
	"github.com/sqlc-dev/marino/format"
	"github.com/sqlc-dev/marino/parser"
)

func TestJSONTableParsing(t *testing.T) {
	p := parser.New()

	cases := []struct {
		name        string
		sql         string
		expectError bool
		// expectedRestore is the expected canonical-restored form. If empty,
		// the test only verifies round-trip parse-restore-parse stability.
		expectedRestore string
	}{
		{
			name: "Basic projection from issue example",
			sql: `SELECT * FROM JSON_TABLE(
				'[{"a":1,"b":"x"},{"a":2,"b":"y"}]',
				'$[*]' COLUMNS (
					a INT          PATH '$.a',
					b VARCHAR(10)  PATH '$.b'
				)
			) AS jt`,
		},
		{
			name: "Joined with NESTED PATH and FOR ORDINALITY",
			sql: `SELECT t.id, jt.ord, jt.zip
				FROM customer t,
				     JSON_TABLE(
				       t.addresses,
				       '$[*]' COLUMNS (
				         ord FOR ORDINALITY,
				         zip VARCHAR(10) PATH '$.zip' DEFAULT '00000' ON EMPTY,
				         NESTED PATH '$.tags[*]' COLUMNS (tag VARCHAR(32) PATH '$')
				       )
				     ) AS jt`,
		},
		{
			name: "Validation example",
			sql:  `SELECT * FROM JSON_TABLE('[{"a":1},{"a":2}]', '$[*]' COLUMNS (a INT PATH '$.a')) AS jt`,
		},
		{
			name: "ON EMPTY and ON ERROR with NULL",
			sql:  `SELECT * FROM JSON_TABLE('[]', '$' COLUMNS (a INT PATH '$.a' NULL ON EMPTY NULL ON ERROR)) AS jt`,
		},
		{
			name: "ON EMPTY and ON ERROR with ERROR",
			sql:  `SELECT * FROM JSON_TABLE('[]', '$' COLUMNS (a INT PATH '$.a' ERROR ON EMPTY ERROR ON ERROR)) AS jt`,
		},
		{
			name: "Only ON ERROR",
			sql:  `SELECT * FROM JSON_TABLE('[]', '$' COLUMNS (a INT PATH '$.a' DEFAULT '0' ON ERROR)) AS jt`,
		},
		{
			name: "EXISTS PATH",
			sql:  `SELECT * FROM JSON_TABLE('{"a":1}', '$' COLUMNS (has_a INT EXISTS PATH '$.a')) AS jt`,
		},
		{
			name: "NESTED without PATH keyword",
			sql:  `SELECT * FROM JSON_TABLE('[]', '$' COLUMNS (a INT PATH '$.a', NESTED '$.children[*]' COLUMNS (b INT PATH '$.b'))) AS jt`,
		},
		{
			name: "JSON_TABLE without alias",
			sql:  `SELECT * FROM JSON_TABLE('[]', '$' COLUMNS (a INT PATH '$.a'))`,
		},
		{
			name: "Deeply nested NESTED PATH",
			sql: `SELECT * FROM JSON_TABLE(
				'[{"x":[{"y":[1,2]}]}]',
				'$[*]' COLUMNS (
					NESTED PATH '$.x[*]' COLUMNS (
						NESTED PATH '$.y[*]' COLUMNS (v INT PATH '$')
					)
				)
			) AS jt`,
		},
		{
			name:        "Missing COLUMNS clause is rejected",
			sql:         `SELECT * FROM JSON_TABLE('[]', '$') AS jt`,
			expectError: true,
		},
		{
			name:        "Missing path argument is rejected",
			sql:         `SELECT * FROM JSON_TABLE('[]') AS jt`,
			expectError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stmt, err := p.ParseOneStmt(tc.sql, "", "")
			if tc.expectError {
				if err == nil {
					t.Fatalf("expected parse error for: %s", tc.sql)
				}
				return
			}
			if err != nil {
				t.Fatalf("failed to parse %q: %v", tc.sql, err)
			}

			selectStmt, ok := stmt.(*ast.SelectStmt)
			if !ok {
				t.Fatalf("expected *ast.SelectStmt, got %T", stmt)
			}
			if selectStmt.From == nil {
				t.Fatalf("expected FROM clause for: %s", tc.sql)
			}
			if findJSONTable(selectStmt.From.TableRefs) == nil {
				t.Fatalf("expected to find a JSONTable node in: %s", tc.sql)
			}

			// Restore round-trip.
			var sb strings.Builder
			if err := stmt.Restore(format.NewRestoreCtx(format.DefaultRestoreFlags, &sb)); err != nil {
				t.Fatalf("failed to restore %q: %v", tc.sql, err)
			}
			restored := sb.String()
			if !strings.Contains(strings.ToUpper(restored), "JSON_TABLE") {
				t.Fatalf("restored SQL missing JSON_TABLE: %s", restored)
			}
			if _, err := p.ParseOneStmt(restored, "", ""); err != nil {
				t.Fatalf("failed to re-parse restored SQL %q: %v", restored, err)
			}
		})
	}
}

// findJSONTable walks a FROM-clause tree looking for a JSONTable wrapped in a TableSource.
func findJSONTable(node ast.ResultSetNode) *ast.JSONTable {
	if node == nil {
		return nil
	}
	switch n := node.(type) {
	case *ast.JSONTable:
		return n
	case *ast.TableSource:
		return findJSONTable(n.Source)
	case *ast.Join:
		if jt := findJSONTable(n.Left); jt != nil {
			return jt
		}
		return findJSONTable(n.Right)
	}
	return nil
}

func TestJSONTableASTContents(t *testing.T) {
	p := parser.New()
	sql := `SELECT * FROM JSON_TABLE(
		'[{"a":1}]',
		'$[*]' COLUMNS (
			ord FOR ORDINALITY,
			a INT PATH '$.a' DEFAULT '99' ON EMPTY ERROR ON ERROR,
			NESTED PATH '$.tags[*]' COLUMNS (tag VARCHAR(32) PATH '$')
		)
	) AS jt`

	stmt, err := p.ParseOneStmt(sql, "", "")
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}
	jt := findJSONTable(stmt.(*ast.SelectStmt).From.TableRefs)
	if jt == nil {
		t.Fatalf("expected JSONTable node")
	}

	if jt.Path != "$[*]" {
		t.Fatalf("Path: got %q, want %q", jt.Path, "$[*]")
	}
	if got := len(jt.Columns); got != 3 {
		t.Fatalf("expected 3 columns, got %d", got)
	}

	// Column 0: ord FOR ORDINALITY
	if jt.Columns[0].Tp != ast.JSONTableColumnForOrdinality {
		t.Fatalf("column 0 Tp: got %d", jt.Columns[0].Tp)
	}
	if jt.Columns[0].Name.L != "ord" {
		t.Fatalf("column 0 Name: got %q", jt.Columns[0].Name.L)
	}

	// Column 1: a INT PATH '$.a' DEFAULT '99' ON EMPTY ERROR ON ERROR
	c := jt.Columns[1]
	if c.Tp != ast.JSONTableColumnPath {
		t.Fatalf("column 1 Tp: got %d", c.Tp)
	}
	if c.Name.L != "a" {
		t.Fatalf("column 1 Name: got %q", c.Name.L)
	}
	if c.Path != "$.a" {
		t.Fatalf("column 1 Path: got %q", c.Path)
	}
	if !c.HasOnEmpty || c.OnEmpty.Tp != ast.JSONTableOnHandlerDefault || c.OnEmpty.DefaultValue != "99" {
		t.Fatalf("column 1 OnEmpty unexpected: hasOnEmpty=%v handler=%+v", c.HasOnEmpty, c.OnEmpty)
	}
	if !c.HasOnError || c.OnError.Tp != ast.JSONTableOnHandlerError {
		t.Fatalf("column 1 OnError unexpected: hasOnError=%v handler=%+v", c.HasOnError, c.OnError)
	}

	// Column 2: NESTED PATH '$.tags[*]' COLUMNS (tag VARCHAR(32) PATH '$')
	nested := jt.Columns[2]
	if nested.Tp != ast.JSONTableColumnNested {
		t.Fatalf("column 2 Tp: got %d", nested.Tp)
	}
	if nested.Path != "$.tags[*]" {
		t.Fatalf("column 2 Path: got %q", nested.Path)
	}
	if got := len(nested.NestedColumns); got != 1 {
		t.Fatalf("expected 1 nested column, got %d", got)
	}
	if nested.NestedColumns[0].Name.L != "tag" || nested.NestedColumns[0].Path != "$" {
		t.Fatalf("nested col 0 unexpected: name=%q path=%q", nested.NestedColumns[0].Name.L, nested.NestedColumns[0].Path)
	}
}
