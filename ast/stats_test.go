// Copyright 2025 PingCAP, Inc.
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
	"strings"
	"testing"

	"github.com/sqlc-dev/marino/ast"
	"github.com/sqlc-dev/marino/format"
	"github.com/sqlc-dev/marino/parser"

	"reflect"
)

func TestRefreshStatsStmt(t *testing.T) {
	tests := []struct {
		sql     string
		want    string
		modeSet bool
		mode    ast.RefreshStatsMode
	}{
		{
			sql:  "REFRESH STATS *.*",
			want: "REFRESH STATS *.*",
		},
		{
			sql:  "refresh stats *.*",
			want: "REFRESH STATS *.*",
		},
		{
			sql:  "REFRESH STATS db1.*",
			want: "REFRESH STATS `db1`.*",
		},
		{
			sql:  "REFRESH STATS db1.t1",
			want: "REFRESH STATS `db1`.`t1`",
		},
		{
			sql:  "REFRESH STATS table1",
			want: "REFRESH STATS `table1`",
		},
		{
			sql:  "REFRESH STATS table1, table2",
			want: "REFRESH STATS `table1`, `table2`",
		},
		{
			sql:  "REFRESH STATS *.*, db1.*, db2.t1, table1, table2",
			want: "REFRESH STATS *.*, `db1`.*, `db2`.`t1`, `table1`, `table2`",
		},
		{
			sql:     "REFRESH STATS table1 full",
			want:    "REFRESH STATS `table1` FULL",
			modeSet: true,
			mode:    ast.RefreshStatsModeFull,
		},
		{
			sql:  "REFRESH STATS table1 cluster",
			want: "REFRESH STATS `table1` CLUSTER",
		},
		{
			sql:     "REFRESH STATS db1.* lite cluster",
			want:    "REFRESH STATS `db1`.* LITE CLUSTER",
			modeSet: true,
			mode:    ast.RefreshStatsModeLite,
		},
	}

	p := parser.New()
	for _, test := range tests {
		stmt, err := p.ParseOneStmt(test.sql, "", "")
		if err != nil {
			t.Fatal(err)
		}
		rs := stmt.(*ast.RefreshStatsStmt)
		if test.modeSet {
			if rs.RefreshMode == nil {
				t.Fatal("expected non-nil")
			}
			if !reflect.DeepEqual(test.mode, *rs.RefreshMode) {
				t.Fatalf("got %v, want %v", *rs.RefreshMode, test.mode)
			}
		} else {
			if rs.RefreshMode != nil {
				t.Fatalf("expected nil, got %v", rs.RefreshMode)
			}
		}
		var sb strings.Builder
		err = stmt.Restore(format.NewRestoreCtx(format.DefaultRestoreFlags, &sb))
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(test.want, sb.String()) {
			t.Fatalf("got %v, want %v", sb.String(), test.want)
		}
	}
}

func TestFlushStatsDeltaScoped(t *testing.T) {
	tests := []struct {
		sql     string
		want    string
		objects int // expected number of FlushObjects
		cluster bool
	}{
		{
			sql:     "FLUSH STATS_DELTA *.*",
			want:    "FLUSH STATS_DELTA *.*",
			objects: 1,
		},
		{
			sql:     "FLUSH STATS_DELTA *.* CLUSTER",
			want:    "FLUSH STATS_DELTA *.* CLUSTER",
			objects: 1,
			cluster: true,
		},
		{
			sql:     "FLUSH STATS_DELTA db1.*",
			want:    "FLUSH STATS_DELTA `db1`.*",
			objects: 1,
		},
		{
			sql:     "FLUSH STATS_DELTA db1.t1",
			want:    "FLUSH STATS_DELTA `db1`.`t1`",
			objects: 1,
		},
		{
			sql:     "FLUSH STATS_DELTA db1.t1 CLUSTER",
			want:    "FLUSH STATS_DELTA `db1`.`t1` CLUSTER",
			objects: 1,
			cluster: true,
		},
		{
			sql:     "FLUSH STATS_DELTA table1",
			want:    "FLUSH STATS_DELTA `table1`",
			objects: 1,
		},
		{
			sql:     "FLUSH STATS_DELTA db1.t1, db2.*, *.*",
			want:    "FLUSH STATS_DELTA `db1`.`t1`, `db2`.*, *.*",
			objects: 3,
		},
		{
			sql:     "FLUSH STATS_DELTA db1.t1, db2.* CLUSTER",
			want:    "FLUSH STATS_DELTA `db1`.`t1`, `db2`.* CLUSTER",
			objects: 2,
			cluster: true,
		},
	}

	p := parser.New()
	for _, test := range tests {
		t.Run(test.sql, func(t *testing.T) {
			stmt, err := p.ParseOneStmt(test.sql, "", "")
			if err != nil {
				t.Fatal(err)
			}
			fs := stmt.(*ast.FlushStmt)
			if !reflect.DeepEqual(ast.FlushStatsDelta, fs.Tp) {
				t.Fatalf("got %v, want %v", fs.Tp, ast.FlushStatsDelta)
			}
			if got := len(fs.FlushObjects); got != test.objects {
				t.Fatalf("expected length %d, got %d", test.objects, got)
			}
			if !reflect.DeepEqual(test.cluster, fs.IsCluster) {
				t.Fatalf("got %v, want %v", fs.IsCluster, test.cluster)
			}
			var sb strings.Builder
			err = stmt.Restore(format.NewRestoreCtx(format.DefaultRestoreFlags, &sb))
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(test.want, sb.String()) {
				t.Fatalf("got %v, want %v", sb.String(), test.want)
			}
		})
	}

	dedupTests := []struct {
		name string
		sql  string
		want string
	}{
		{
			name: "global overrides all",
			sql:  "FLUSH STATS_DELTA table1, db1.t1, *.*, db2.t2",
			want: "FLUSH STATS_DELTA *.*",
		},
		{
			name: "database removes prior tables",
			sql:  "FLUSH STATS_DELTA db1.t1, db2.t1, db1.*, db2.t2",
			want: "FLUSH STATS_DELTA `db2`.`t1`, `db1`.*, `db2`.`t2`",
		},
		{
			name: "table duplicates case insensitive",
			sql:  "FLUSH STATS_DELTA db1.t1, db1.T1, db2.t1",
			want: "FLUSH STATS_DELTA `db1`.`t1`, `db2`.`t1`",
		},
	}
	for _, test := range dedupTests {
		t.Run(test.name, func(t *testing.T) {
			stmt, err := p.ParseOneStmt(test.sql, "", "")
			if err != nil {
				t.Fatal(err)
			}
			fs := stmt.(*ast.FlushStmt)
			fs.DedupFlushObjects()
			var sb strings.Builder
			err = fs.Restore(format.NewRestoreCtx(format.DefaultRestoreFlags, &sb))
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(test.want, sb.String()) {
				t.Fatalf("got %v, want %v", sb.String(), test.want)
			}
		})
	}
}

func TestRefreshStatsStmtDedup(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want string
	}{
		{
			name: "global overrides all",
			sql:  "REFRESH STATS table1, db1.t1, *.*, db2.t2",
			want: "REFRESH STATS *.*",
		},
		{
			name: "database removes prior tables",
			sql:  "REFRESH STATS db1.t1, db2.t1, db1.*, db2.t2",
			want: "REFRESH STATS `db2`.`t1`, `db1`.*, `db2`.`t2`",
		},
		{
			name: "table duplicates case insensitive",
			sql:  "REFRESH STATS db1.t1, db1.T1, db2.t1",
			want: "REFRESH STATS `db1`.`t1`, `db2`.`t1`",
		},
		{
			name: "table duplicates without database",
			sql:  "REFRESH STATS table1, table1, table2",
			want: "REFRESH STATS `table1`, `table2`",
		},
		{
			name: "database duplicates case insensitive",
			sql:  "REFRESH STATS db1.*, DB1.*, db2.t1",
			want: "REFRESH STATS `db1`.*, `db2`.`t1`",
		},
	}

	p := parser.New()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stmt, err := p.ParseOneStmt(test.sql, "", "")
			if err != nil {
				t.Fatal(err)
			}
			rs := stmt.(*ast.RefreshStatsStmt)
			rs.Dedup()
			var sb strings.Builder
			err = rs.Restore(format.NewRestoreCtx(format.DefaultRestoreFlags, &sb))
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(test.want, sb.String()) {
				t.Fatalf("got %v, want %v", sb.String(), test.want)
			}
		})
	}
}
