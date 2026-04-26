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

	"fmt"
	"reflect"
)

func TestLateralParsing(t *testing.T) {
	p := parser.New()

	testCases := []struct {
		name         string
		sql          string
		expectError  bool
		checkLateral bool // whether to verify Lateral flag is set
		columnNames  []string
	}{
		{
			name:         "LATERAL with comma syntax",
			sql:          "SELECT * FROM t1, LATERAL (SELECT t1.a) AS dt",
			expectError:  false,
			checkLateral: true,
		},
		{
			name:         "LATERAL with LEFT JOIN",
			sql:          "SELECT * FROM t1 LEFT JOIN LATERAL (SELECT t1.b) AS dt ON true",
			expectError:  false,
			checkLateral: true,
		},
		{
			name:         "LATERAL with CROSS JOIN",
			sql:          "SELECT * FROM t1 CROSS JOIN LATERAL (SELECT t1.c) AS dt",
			expectError:  false,
			checkLateral: true,
		},
		{
			name:         "LATERAL with RIGHT JOIN",
			sql:          "SELECT * FROM t1 RIGHT JOIN LATERAL (SELECT t1.d) AS dt ON true",
			expectError:  false, // Parser allows it, planner will reject
			checkLateral: true,
		},
		{
			name:         "LATERAL with INNER JOIN",
			sql:          "SELECT * FROM t1 JOIN LATERAL (SELECT t1.e) AS dt ON true",
			expectError:  false,
			checkLateral: true,
		},
		{
			name:         "LATERAL with complex subquery",
			sql:          "SELECT * FROM t1, LATERAL (SELECT t1.a, COUNT(*) FROM t2 WHERE t2.x = t1.x GROUP BY t1.a) AS dt",
			expectError:  false,
			checkLateral: true,
		},
		{
			name:         "LATERAL with nested subquery",
			sql:          "SELECT * FROM t1, LATERAL (SELECT * FROM (SELECT t1.a) AS inner_dt) AS dt",
			expectError:  false,
			checkLateral: true,
		},
		{
			name:         "Multiple LATERAL joins",
			sql:          "SELECT * FROM t1, LATERAL (SELECT t1.a) AS dt1, LATERAL (SELECT t1.b) AS dt2",
			expectError:  false,
			checkLateral: true,
		},
		{
			name:        "Non-LATERAL derived table",
			sql:         "SELECT * FROM t1, (SELECT a FROM t2) AS dt",
			expectError: false,
			// Lateral flag should be false for non-LATERAL
		},
		{
			name:         "LATERAL with WHERE clause",
			sql:          "SELECT * FROM t1, LATERAL (SELECT * FROM t2 WHERE t2.x = t1.x) AS dt WHERE dt.y > 10",
			expectError:  false,
			checkLateral: true,
		},
		{
			name:         "LATERAL with column list",
			sql:          "SELECT * FROM t1, LATERAL (SELECT t1.a, t1.b) AS dt(c1, c2)",
			expectError:  false,
			checkLateral: true,
			columnNames:  []string{"c1", "c2"},
		},
		{
			name:         "LATERAL with column list no AS",
			sql:          "SELECT * FROM t1, LATERAL (SELECT t1.a) dt(col1)",
			expectError:  false,
			checkLateral: true,
			columnNames:  []string{"col1"},
		},
		{
			name:         "LATERAL with column list and JOIN",
			sql:          "SELECT * FROM t1 LEFT JOIN LATERAL (SELECT t1.a, t1.b, t1.c) AS dt(x, y, z) ON true",
			expectError:  false,
			checkLateral: true,
			columnNames:  []string{"x", "y", "z"},
		},
		{
			name:        "LATERAL without alias is rejected",
			sql:         "SELECT * FROM t1, LATERAL (SELECT t1.a)",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stmt, err := p.ParseOneStmt(tc.sql, "", "")

			if tc.expectError {
				if err == nil {
					t.Fatalf("Expected parsing to fail for: %s", tc.sql)
				}
				return
			}

			if err != nil {
				t.Fatalf("%s: %v", fmt.Sprintf("Failed to parse: %s", tc.sql), err)
			}
			if stmt == nil {
				t.Fatal("expected non-nil")
			}

			// Test round-trip: parse -> restore -> parse again
			var sb strings.Builder
			restoreCtx := format.NewRestoreCtx(format.RestoreStringSingleQuotes, &sb)
			err = stmt.Restore(restoreCtx)
			if err != nil {
				t.Fatalf("%v: %v", "Failed to restore statement", err)
			}

			restored := sb.String()
			if tc.checkLateral {
				// Verify LATERAL keyword is preserved in restoration
				if !strings.Contains(restored, "LATERAL") {
					t.Fatalf("%s: expected %q to contain %q", fmt.Sprintf("LATERAL keyword missing in restored SQL: %s", restored), restored, "LATERAL")
				}
			}

			// Parse the restored SQL to ensure it's valid (round-trip test)
			stmt2, err := p.ParseOneStmt(restored, "", "")
			if err != nil {
				t.Fatalf("%s: %v", fmt.Sprintf("Failed to parse restored SQL: %s", restored), err)
			}
			if stmt2 == nil {
				t.Fatal("expected non-nil")
			}

			// Verify AST flags on both original and round-tripped statements.
			for _, stmtToCheck := range []struct {
				label string
				node  ast.StmtNode
			}{
				{"original", stmt},
				{"round-trip", stmt2},
			} {
				selectStmt, ok := stmtToCheck.node.(*ast.SelectStmt)
				if !(ok) {
					t.Fatalf("[%s] Statement should be SelectStmt", stmtToCheck.label)
				}
				if selectStmt.From == nil {
					t.Fatalf("[%s] FROM clause should not be nil", stmtToCheck.label)
				}

				if tc.checkLateral {
					lateralTS := findLateralTableSource(selectStmt.From.TableRefs)
					if lateralTS == nil {
						t.Fatalf("[%s] LATERAL TableSource not found for: %s", stmtToCheck.label, tc.sql)
					}

					if len(tc.columnNames) > 0 {
						if got := len(lateralTS.ColumnNames); got != len(tc.columnNames) {
							t.Fatalf("%s: expected length %d, got %d", fmt.Sprintf("[%s] column name count mismatch", stmtToCheck.label), len(tc.columnNames), got)
						}
						for i, expected := range tc.columnNames {
							if !reflect.DeepEqual(expected, lateralTS.ColumnNames[i].L) {
								t.Fatalf("%s: got %v, want %v", fmt.Sprintf("[%s] column name mismatch at index %d", stmtToCheck.label, i), lateralTS.ColumnNames[i].L, expected)
							}
						}
					}
				} else {
					lateralTS := findLateralTableSource(selectStmt.From.TableRefs)
					if lateralTS != nil {
						t.Fatalf("[%s] Lateral should be false for non-LATERAL query: %s", stmtToCheck.label, tc.sql)
					}
				}
			}
		})
	}
}

// findLateralTableSource recursively searches for the first LATERAL TableSource in a ResultSetNode.
func findLateralTableSource(node ast.ResultSetNode) *ast.TableSource {
	if node == nil {
		return nil
	}

	switch n := node.(type) {
	case *ast.TableSource:
		if n.Lateral {
			return n
		}
		return findLateralTableSource(n.Source)
	case *ast.Join:
		if ts := findLateralTableSource(n.Left); ts != nil {
			return ts
		}
		return findLateralTableSource(n.Right)
	}

	return nil
}
