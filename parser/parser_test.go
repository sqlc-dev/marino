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

package parser_test

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/sqlc-dev/marino/ast"
	"github.com/sqlc-dev/marino/charset"
	. "github.com/sqlc-dev/marino/format"
	"github.com/sqlc-dev/marino/mysql"
	"github.com/sqlc-dev/marino/opcode"
	"github.com/sqlc-dev/marino/types"
	"github.com/sqlc-dev/marino/parser"
	"github.com/sqlc-dev/marino/terror"
)

func TestSimple(t *testing.T) {
	p := parser.New()

	reservedKws := []string{
		"add", "all", "alter", "analyze", "and", "as", "asc", "between", "bigint",
		"binary", "blob", "both", "by", "call", "cascade", "case", "change", "character", "check", "collate",
		"column", "constraint", "convert", "create", "cross", "current_date", "current_time",
		"current_timestamp", "current_user", "database", "databases", "day_hour", "day_microsecond",
		"day_minute", "day_second", "decimal", "default", "delete", "desc", "describe",
		"distinct", "distinctRow", "div", "double", "drop", "dual", "else", "enclosed", "escaped",
		"exists", "explain", "false", "float", "fetch", "for", "force", "foreign", "from",
		"fulltext", "grant", "group", "having", "hour_microsecond", "hour_minute",
		"hour_second", "if", "ignore", "in", "index", "infile", "inner", "insert", "int", "into", "integer",
		"interval", "is", "join", "key", "keys", "kill", "leading", "left", "like", "ilike", "limit", "lines", "load",
		"localtime", "localtimestamp", "lock", "longblob", "longtext", "mediumblob", "maxvalue", "mediumint", "mediumtext",
		"minute_microsecond", "minute_second", "mod", "not", "no_write_to_binlog", "null", "numeric",
		"on", "option", "optionally", "or", "order", "outer", "partition", "precision", "primary", "procedure", "range", "read", "real", "recursive",
		"references", "regexp", "rename", "repeat", "replace", "revoke", "restrict", "right", "rlike",
		"schema", "schemas", "second_microsecond", "select", "set", "show", "smallint",
		"starting", "table", "terminated", "then", "tinyblob", "tinyint", "tinytext", "to",
		"trailing", "true", "union", "unique", "unlock", "unsigned",
		"update", "use", "using", "utc_date", "values", "varbinary", "varchar",
		"when", "where", "write", "xor", "year_month", "zerofill",
		"generated", "virtual", "stored", "usage",
		"delayed", "high_priority", "low_priority",
		"cumeDist", "denseRank", "firstValue", "lag", "lastValue", "lead", "nthValue", "ntile",
		"over", "percentRank", "rank", "row", "rows", "rowNumber", "window", "linear",
		"match", "until", "placement", "tablesample", "failedLoginAttempts", "passwordLockTime",
		"cube", "external", "qualify",
		// TODO: support the following keywords
		// "with",
	}
	for _, kw := range reservedKws {
		src := fmt.Sprintf("SELECT * FROM db.%s;", kw)
		_, err := p.ParseOneStmt(src, "", "")

		if err != nil {
			t.Fatalf("%s: %v", fmt.Sprintf("source %s", src), err)
		}

		src = fmt.Sprintf("SELECT * FROM %s.desc", kw)
		_, err = p.ParseOneStmt(src, "", "")
		if err != nil {
			t.Fatalf("%s: %v", fmt.Sprintf("source %s", src), err)
		}

		src = fmt.Sprintf("SELECT t.%s FROM t", kw)
		_, err = p.ParseOneStmt(src, "", "")
		if err != nil {
			t.Fatalf("%s: %v", fmt.Sprintf("source %s", src), err)
		}
	}

	// Testcase for unreserved keywords
	unreservedKws := []string{
		"add_columnar_replica_on_demand", "auto_increment", "after", "begin", "bit", "bool", "boolean", "charset", "columns", "commit",
		"date", "datediff", "datetime", "deallocate", "do", "from_days", "end", "engine", "engines", "execute", "extended", "first", "file", "full",
		"local", "names", "offset", "password", "prepare", "quick", "rollback", "savepoint", "session", "signed",
		"start", "global", "tables", "tablespace", "target", "text", "time", "timestamp", "tidb", "transaction", "truncate", "unknown",
		"value", "warnings", "year", "now", "substr", "subpartition", "subpartitions", "substring", "mode", "any", "some", "user", "identified",
		"collation", "comment", "avg_row_length", "checksum", "compression", "connection", "key_block_size",
		"max_rows", "min_rows", "national", "quarter", "escape", "grants", "status", "fields", "triggers", "language",
		"delay_key_write", "isolation", "partitions", "repeatable", "committed", "uncommitted", "only", "serializable", "level",
		"curtime", "variables", "dayname", "version", "btree", "hash", "row_format", "dynamic", "fixed", "compressed",
		"compact", "redundant", "1 sql_no_cache", "1 sql_cache", "action", "round",
		"enable", "disable", "reverse", "space", "privileges", "get_lock", "release_lock", "sleep", "no", "greatest", "least",
		"binlog", "hex", "unhex", "function", "indexes", "from_unixtime", "processlist", "events", "less", "than", "timediff",
		"ln", "log", "log2", "log10", "timestampdiff", "pi", "proxy", "quote", "none", "super", "shared", "exclusive",
		"always", "stats", "stats_meta", "stats_histogram", "stats_buckets", "stats_healthy", "tidb_version", "replication", "slave", "client",
		"max_connections_per_hour", "max_queries_per_hour", "max_updates_per_hour", "max_user_connections", "event", "reload", "routine", "temporary",
		"following", "preceding", "unbounded", "respect", "nulls", "current", "last", "against", "expansion",
		"chain", "error", "general", "nvarchar", "pack_keys", "p", "shard_row_id_bits", "pre_split_regions",
		"constraints", "role", "replicas", "policy", "s3", "strict", "running", "stop", "preserve", "placement", "attributes", "attribute", "resource",
		"burstable", "calibrate", "masking", "rollup", "manual", "parallel", "channel",
	}
	for _, kw := range unreservedKws {
		src := fmt.Sprintf("SELECT %s FROM tbl;", kw)
		_, err := p.ParseOneStmt(src, "", "")
		if err != nil {
			t.Fatalf("%s: %v", fmt.Sprintf("source %s", src), err)
		}
	}

	// Testcase for prepared statement
	src := "SELECT id+?, id+? from t;"
	_, err := p.ParseOneStmt(src, "", "")
	if err != nil {
		t.Fatal(err)
	}

	// Testcase for -- Comment and unary -- operator
	src = "CREATE TABLE foo (a SMALLINT UNSIGNED, b INT UNSIGNED); -- foo\nSelect --1 from foo;"
	stmts, _, err := p.Parse(src, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if got := len(stmts); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}

	// Testcase for /*! xx */
	// See http://dev.mysql.com/doc/refman/5.7/en/comments.html
	// Fix: https://github.com/pingcap/tidb/issues/971
	src = "/*!40101 SET character_set_client = utf8 */;"
	stmts, _, err = p.Parse(src, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if got := len(stmts); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	stmt := stmts[0]
	_, ok := stmt.(*ast.SetStmt)
	if !(ok) {
		t.Fatal("expected true")
	}

	// for issue #2017
	src = "insert into blobtable (a) values ('/*! truncated */');"
	stmt, err = p.ParseOneStmt(src, "", "")
	if err != nil {
		t.Fatal(err)
	}
	is, ok := stmt.(*ast.InsertStmt)
	if !(ok) {
		t.Fatal("expected true")
	}
	if got := len(is.Lists); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	if got := len(is.Lists[0]); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	if !reflect.DeepEqual("/*! truncated */", is.Lists[0][0].(ast.ValueExpr).GetDatumString()) {
		t.Fatalf("got %v, want %v", is.Lists[0][0].(ast.ValueExpr).GetDatumString(), "/*! truncated */")
	}

	// Testcase for CONVERT(expr,type)
	src = "SELECT CONVERT('111', SIGNED);"
	st, err := p.ParseOneStmt(src, "", "")
	if err != nil {
		t.Fatal(err)
	}
	ss, ok := st.(*ast.SelectStmt)
	if !(ok) {
		t.Fatal("expected true")
	}
	if got := len(ss.Fields.Fields); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	cv, ok := ss.Fields.Fields[0].Expr.(*ast.FuncCastExpr)
	if !(ok) {
		t.Fatal("expected true")
	}
	if !reflect.DeepEqual(ast.CastConvertFunction, cv.FunctionType) {
		t.Fatalf("got %v, want %v", cv.FunctionType, ast.CastConvertFunction)
	}

	// for query start with comment
	srcs := []string{
		"/* some comments */ SELECT CONVERT('111', SIGNED) ;",
		"/* some comments */ /*comment*/ SELECT CONVERT('111', SIGNED) ;",
		"SELECT /*comment*/ CONVERT('111', SIGNED) ;",
		"SELECT CONVERT('111', /*comment*/ SIGNED) ;",
		"SELECT CONVERT('111', SIGNED) /*comment*/;",
	}
	for _, src := range srcs {
		st, err = p.ParseOneStmt(src, "", "")
		if err != nil {
			t.Fatal(err)
		}
		_, ok = st.(*ast.SelectStmt)
		if !(ok) {
			t.Fatal("expected true")
		}
	}

	// for issue #961
	src = "create table t (c int key);"
	st, err = p.ParseOneStmt(src, "", "")
	if err != nil {
		t.Fatal(err)
	}
	cs, ok := st.(*ast.CreateTableStmt)
	if !(ok) {
		t.Fatal("expected true")
	}
	if got := len(cs.Cols); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	if got := len(cs.Cols[0].Options); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	if !reflect.DeepEqual(ast.ColumnOptionPrimaryKey, cs.Cols[0].Options[0].Tp) {
		t.Fatalf("got %v, want %v", cs.Cols[0].Options[0].Tp, ast.ColumnOptionPrimaryKey)
	}

	// for issue #4497
	src = "create table t1(a NVARCHAR(100));"
	_, err = p.ParseOneStmt(src, "", "")
	if err != nil {
		t.Fatal(err)
	}

	// for issue 2803
	src = "use quote;"
	_, err = p.ParseOneStmt(src, "", "")
	if err != nil {
		t.Fatal(err)
	}

	// issue #4354
	src = "select b'';"
	_, err = p.ParseOneStmt(src, "", "")
	if err != nil {
		t.Fatal(err)
	}

	src = "select B'';"
	_, err = p.ParseOneStmt(src, "", "")
	if err != nil {
		t.Fatal(err)
	}

	// src = "select 0b'';"
	// _, err = p.ParseOneStmt(src, "", "")
	// if err == nil { t.Fatal("expected error") }

	// for #4909, support numericType `signed` filedOpt.
	src = "CREATE TABLE t(_sms smallint signed, _smu smallint unsigned);"
	_, err = p.ParseOneStmt(src, "", "")
	if err != nil {
		t.Fatal(err)
	}

	// for #7371, support NATIONAL CHARACTER
	// reference link: https://dev.mysql.com/doc/refman/5.7/en/charset-national.html
	src = "CREATE TABLE t(c1 NATIONAL CHARACTER(10));"
	_, err = p.ParseOneStmt(src, "", "")
	if err != nil {
		t.Fatal(err)
	}

	src = `CREATE TABLE t(a tinyint signed,
		b smallint signed,
		c mediumint signed,
		d int signed,
		e int1 signed,
		f int2 signed,
		g int3 signed,
		h int4 signed,
		i int8 signed,
		j integer signed,
		k bigint signed,
		l bool signed,
		m boolean signed
		);`

	st, err = p.ParseOneStmt(src, "", "")
	if err != nil {
		t.Fatal(err)
	}
	ct, ok := st.(*ast.CreateTableStmt)
	if !(ok) {
		t.Fatal("expected true")
	}
	for _, col := range ct.Cols {
		if !reflect.DeepEqual(uint(0), col.Tp.GetFlag()&mysql.UnsignedFlag) {
			t.Fatalf("got %v, want %v", col.Tp.GetFlag()&mysql.UnsignedFlag, uint(0))
		}
	}

	// for issue #4006
	src = `insert into tb(v) (select v from tb);`
	_, err = p.ParseOneStmt(src, "", "")
	if err != nil {
		t.Fatal(err)
	}

	// for issue #34642
	src = `SELECT a as c having c = a;`
	_, err = p.ParseOneStmt(src, "", "")
	if err != nil {
		t.Fatal(err)
	}

	// for issue #9823
	src = "SELECT 9223372036854775807;"
	st, err = p.ParseOneStmt(src, "", "")
	if err != nil {
		t.Fatal(err)
	}
	sel, ok := st.(*ast.SelectStmt)
	if !(ok) {
		t.Fatal("expected true")
	}
	expr := sel.Fields.Fields[0]
	vExpr := expr.Expr.(*ast.ValueExprBase)
	if !reflect.DeepEqual(ast.KindInt64, vExpr.Kind()) {
		t.Fatalf("got %v, want %v", vExpr.Kind(), ast.KindInt64)
	}
	src = "SELECT 9223372036854775808;"
	st, err = p.ParseOneStmt(src, "", "")
	if err != nil {
		t.Fatal(err)
	}
	sel, ok = st.(*ast.SelectStmt)
	if !(ok) {
		t.Fatal("expected true")
	}
	expr = sel.Fields.Fields[0]
	vExpr = expr.Expr.(*ast.ValueExprBase)
	if !reflect.DeepEqual(ast.KindUint64, vExpr.Kind()) {
		t.Fatalf("got %v, want %v", vExpr.Kind(), ast.KindUint64)
	}

	src = `select 99e+r10 from t1;`
	st, err = p.ParseOneStmt(src, "", "")
	if err != nil {
		t.Fatal(err)
	}
	sel, ok = st.(*ast.SelectStmt)
	if !(ok) {
		t.Fatal("expected true")
	}
	bExpr, ok := sel.Fields.Fields[0].Expr.(*ast.BinaryOperationExpr)
	if !(ok) {
		t.Fatal("expected true")
	}
	if !reflect.DeepEqual(opcode.Plus, bExpr.Op) {
		t.Fatalf("got %v, want %v", bExpr.Op, opcode.Plus)
	}
	if !reflect.DeepEqual("99e", bExpr.L.(*ast.ColumnNameExpr).Name.Name.O) {
		t.Fatalf("got %v, want %v", bExpr.L.(*ast.ColumnNameExpr).Name.Name.O, "99e")
	}
	if !reflect.DeepEqual("r10", bExpr.R.(*ast.ColumnNameExpr).Name.Name.O) {
		t.Fatalf("got %v, want %v", bExpr.R.(*ast.ColumnNameExpr).Name.Name.O, "r10")
	}

	src = `select t./*123*/*,@c3:=0 from t order by t.c1;`
	st, err = p.ParseOneStmt(src, "", "")
	if err != nil {
		t.Fatal(err)
	}
	sel, ok = st.(*ast.SelectStmt)
	if !(ok) {
		t.Fatal("expected true")
	}
	if !reflect.DeepEqual("t", sel.Fields.Fields[0].WildCard.Table.O) {
		t.Fatalf("got %v, want %v", sel.Fields.Fields[0].WildCard.Table.O, "t")
	}
	varExpr, ok := sel.Fields.Fields[1].Expr.(*ast.VariableExpr)
	if !(ok) {
		t.Fatal("expected true")
	}
	if !reflect.DeepEqual("c3", varExpr.Name) {
		t.Fatalf("got %v, want %v", varExpr.Name, "c3")
	}

	src = `select t.1e from test.t;`
	st, err = p.ParseOneStmt(src, "", "")
	if err != nil {
		t.Fatal(err)
	}
	sel, ok = st.(*ast.SelectStmt)
	if !(ok) {
		t.Fatal("expected true")
	}
	colExpr, ok := sel.Fields.Fields[0].Expr.(*ast.ColumnNameExpr)
	if !(ok) {
		t.Fatal("expected true")
	}
	if !reflect.DeepEqual("t", colExpr.Name.Table.O) {
		t.Fatalf("got %v, want %v", colExpr.Name.Table.O, "t")
	}
	if !reflect.DeepEqual("1e", colExpr.Name.Name.O) {
		t.Fatalf("got %v, want %v", colExpr.Name.Name.O, "1e")
	}
	tName := sel.From.TableRefs.Left.(*ast.TableSource).Source.(*ast.TableName)
	if !reflect.DeepEqual("test", tName.Schema.O) {
		t.Fatalf("got %v, want %v", tName.Schema.O, "test")
	}
	if !reflect.DeepEqual("t", tName.Name.O) {
		t.Fatalf("got %v, want %v", tName.Name.O, "t")
	}

	src = "select t. `a` > 10 from t;"
	st, err = p.ParseOneStmt(src, "", "")
	if err != nil {
		t.Fatal(err)
	}
	bExpr, ok = st.(*ast.SelectStmt).Fields.Fields[0].Expr.(*ast.BinaryOperationExpr)
	if !(ok) {
		t.Fatal("expected true")
	}
	if !reflect.DeepEqual(opcode.GT, bExpr.Op) {
		t.Fatalf("got %v, want %v", bExpr.Op, opcode.GT)
	}
	if !reflect.DeepEqual("a", bExpr.L.(*ast.ColumnNameExpr).Name.Name.O) {
		t.Fatalf("got %v, want %v", bExpr.L.(*ast.ColumnNameExpr).Name.Name.O, "a")
	}
	if !reflect.DeepEqual("t", bExpr.L.(*ast.ColumnNameExpr).Name.Table.O) {
		t.Fatalf("got %v, want %v", bExpr.L.(*ast.ColumnNameExpr).Name.Table.O, "t")
	}
	if !reflect.DeepEqual(int64(10), bExpr.R.(ast.ValueExpr).GetValue().(int64)) {
		t.Fatalf("got %v, want %v", bExpr.R.(ast.ValueExpr).GetValue().(int64), int64(10))
	}

	p.SetSQLMode(mysql.ModeANSIQuotes)
	src = `select t."dot"=10 from t;`
	st, err = p.ParseOneStmt(src, "", "")
	if err != nil {
		t.Fatal(err)
	}
	bExpr, ok = st.(*ast.SelectStmt).Fields.Fields[0].Expr.(*ast.BinaryOperationExpr)
	if !(ok) {
		t.Fatal("expected true")
	}
	if !reflect.DeepEqual(opcode.EQ, bExpr.Op) {
		t.Fatalf("got %v, want %v", bExpr.Op, opcode.EQ)
	}
	if !reflect.DeepEqual("dot", bExpr.L.(*ast.ColumnNameExpr).Name.Name.O) {
		t.Fatalf("got %v, want %v", bExpr.L.(*ast.ColumnNameExpr).Name.Name.O, "dot")
	}
	if !reflect.DeepEqual("t", bExpr.L.(*ast.ColumnNameExpr).Name.Table.O) {
		t.Fatalf("got %v, want %v", bExpr.L.(*ast.ColumnNameExpr).Name.Table.O, "t")
	}
	if !reflect.DeepEqual(int64(10), bExpr.R.(ast.ValueExpr).GetValue().(int64)) {
		t.Fatalf("got %v, want %v", bExpr.R.(ast.ValueExpr).GetValue().(int64), int64(10))
	}
}

func TestSpecialComments(t *testing.T) {
	p := parser.New()

	// 1. Make sure /*! ... */ respects the same SQL mode.
	_, err := p.ParseOneStmt(`SELECT /*! '\' */;`, "", "")
	if err == nil {
		t.Fatal("expected error")
	}

	p.SetSQLMode(mysql.ModeNoBackslashEscapes)
	st, err := p.ParseOneStmt(`SELECT /*! '\' */;`, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := st.(*ast.SelectStmt); !ok {
		t.Fatalf("expected type %T, got %T", &ast.SelectStmt{}, st)
	}

	// 2. Make sure multiple statements inside /*! ... */ will not crash
	// (this is issue #330)
	stmts, _, err := p.Parse("/*! SET x = 1; SELECT 2 */", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if got := len(stmts); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if _, ok := stmts[0].(*ast.SetStmt); !ok {
		t.Fatalf("expected type %T, got %T", &ast.SetStmt{}, stmts[0])
	}
	if !reflect.DeepEqual("/*! SET x = 1;", stmts[0].Text()) {
		t.Fatalf("got %v, want %v", stmts[0].Text(), "/*! SET x = 1;")
	}
	if _, ok := stmts[1].(*ast.SelectStmt); !ok {
		t.Fatalf("expected type %T, got %T", &ast.SelectStmt{}, stmts[1])
	}
	if !reflect.DeepEqual(" SELECT 2 */", stmts[1].Text()) {
		t.Fatalf("got %v, want %v", stmts[1].Text(), " SELECT 2 */")
	}
	// ^ not sure if correct approach; having multiple statements in MySQL is a syntax error.

	// 3. Make sure invalid text won't cause infinite loop
	// (this is issue #336)
	st, err = p.ParseOneStmt("SELECT /*+ 😅 */ SLEEP(1);", "", "")
	if err != nil {
		t.Fatal(err)
	}
	sel, ok := st.(*ast.SelectStmt)
	if !(ok) {
		t.Fatal("expected true")
	}
	if got := len(sel.TableHints); got != 0 {
		t.Fatalf("expected length %d, got %d", 0, got)
	}
}

type testCase struct {
	src     string
	ok      bool
	restore string
}

type testErrMsgCase struct {
	src string
	err error
}

func RunTest(t *testing.T, table []testCase, enableWindowFunc bool, MariaDB bool) {
	p := parser.New()
	p.EnableWindowFunc(enableWindowFunc)
	p.SetMariaDB(MariaDB)
	for _, tbl := range table {
		_, _, err := p.Parse(tbl.src, "", "")
		if !tbl.ok {
			if err == nil {
				t.Fatalf("source %v, error %v", tbl.src, err)
			}
			continue
		}
		if err != nil {
			t.Fatalf("%s: %v", fmt.Sprintf("source:\n%v\nerror:\n%v", tbl.src, err), err)
		}
		// restore correctness test
		if tbl.ok {
			RunRestoreTest(t, tbl.src, tbl.restore, enableWindowFunc, MariaDB)
		}
	}
}

func RunRestoreTest(t *testing.T, sourceSQLs, expectSQLs string, enableWindowFunc bool, MariaDB bool) {
	var sb strings.Builder
	p := parser.New()
	p.EnableWindowFunc(enableWindowFunc)
	p.SetMariaDB(MariaDB)
	comment := fmt.Sprintf("source %v", sourceSQLs)
	stmts, _, err := p.Parse(sourceSQLs, "", "")
	if err != nil {
		t.Fatalf("%s: %v", fmt.Sprintf("source %v", sourceSQLs), err)
	}
	restoreSQLs := ""
	for _, stmt := range stmts {
		sb.Reset()
		err = stmt.Restore(NewRestoreCtx(DefaultRestoreFlags, &sb))
		if err != nil {
			t.Fatalf("%v: %v", comment, err)
		}
		restoreSQL := sb.String()
		comment = fmt.Sprintf("source %v; restore %v", sourceSQLs, restoreSQL)
		restoreStmt, err := p.ParseOneStmt(restoreSQL, "", "")
		if err != nil {
			t.Fatalf("%v: %v", comment, err)
		}
		CleanNodeText(stmt)
		CleanNodeText(restoreStmt)
		if !reflect.DeepEqual(stmt, restoreStmt) {
			t.Fatalf("%v: got %v, want %v", comment, restoreStmt, stmt)
		}
		if restoreSQLs != "" {
			restoreSQLs += "; "
		}
		restoreSQLs += restoreSQL
	}
	if !reflect.DeepEqual(expectSQLs, restoreSQLs) {
		t.Fatalf("%s: got %v, want %v", fmt.Sprintf("restore %v; expect %v", restoreSQLs, expectSQLs), restoreSQLs, expectSQLs)
	}
}

func RunErrMsgTest(t *testing.T, table []testErrMsgCase) {
	p := parser.New()
	for _, tbl := range table {
		_, _, err := p.Parse(tbl.src, "", "")
		comment := fmt.Sprintf("source %v", tbl.src)
		if tbl.err != nil {
			if !(terror.ErrorEqual(err, tbl.err)) {
				t.Fatal(comment)
			}
		} else {
			if err != nil {
				t.Fatalf("%v: %v", comment, err)
			}
		}
	}
}

func TestSetVariable(t *testing.T) {
	table := []struct {
		Input      string
		Name       string
		IsGlobal   bool
		IsInstance bool
		IsSystem   bool
	}{

		// Set system variable xx.xx, although xx.xx isn't a system variable, the parser should accept it.
		{"set xx.xx = 666", "xx.xx", false, false, true},
		// Set session system variable xx.xx
		{"set session xx.xx = 666", "xx.xx", false, false, true},
		{"set local xx.xx = 666", "xx.xx", false, false, true},
		{"set global xx.xx = 666", "xx.xx", true, false, true},
		{"set instance xx.xx = 666", "xx.xx", false, true, true},

		{"set @@xx.xx = 666", "xx.xx", false, false, true},
		{"set @@session.xx.xx = 666", "xx.xx", false, false, true},
		{"set @@local.xx.xx = 666", "xx.xx", false, false, true},
		{"set @@global.xx.xx = 666", "xx.xx", true, false, true},
		{"set @@instance.xx.xx = 666", "xx.xx", false, true, true},

		// Set user defined variable xx.xx
		{"set @xx.xx = 666", "xx.xx", false, false, false},
	}

	p := parser.New()
	for _, tbl := range table {
		stmt, err := p.ParseOneStmt(tbl.Input, "", "")
		if err != nil {
			t.Fatal(err)
		}

		setStmt, ok := stmt.(*ast.SetStmt)
		if !(ok) {
			t.Fatal("expected true")
		}
		if got := len(setStmt.Variables); got != 1 {
			t.Fatalf("expected length %d, got %d", 1, got)
		}

		v := setStmt.Variables[0]
		if !reflect.DeepEqual(tbl.Name, v.Name) {
			t.Fatalf("got %v, want %v", v.Name, tbl.Name)
		}
		if !reflect.DeepEqual(tbl.IsGlobal, v.IsGlobal) {
			t.Fatalf("got %v, want %v", v.IsGlobal, tbl.IsGlobal)
		}
		if !reflect.DeepEqual(tbl.IsInstance, v.IsInstance) {
			t.Fatalf("got %v, want %v", v.IsInstance, tbl.IsInstance)
		}
		if !reflect.DeepEqual(tbl.IsSystem, v.IsSystem) {
			t.Fatalf("got %v, want %v", v.IsSystem, tbl.IsSystem)
		}
	}

	_, err := p.ParseOneStmt("set xx.xx.xx = 666", "", "")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestFlushTable(t *testing.T) {
	p := parser.New()
	stmt, _, err := p.Parse("flush local tables tbl1,tbl2 with read lock", "", "")
	if err != nil {
		t.Fatal(err)
	}
	flushTable := stmt[0].(*ast.FlushStmt)
	if !reflect.DeepEqual(ast.FlushTables, flushTable.Tp) {
		t.Fatalf("got %v, want %v", flushTable.Tp, ast.FlushTables)
	}
	if !reflect.DeepEqual("tbl1", flushTable.Tables[0].Name.L) {
		t.Fatalf("got %v, want %v", flushTable.Tables[0].Name.L, "tbl1")
	}
	if !reflect.DeepEqual("tbl2", flushTable.Tables[1].Name.L) {
		t.Fatalf("got %v, want %v", flushTable.Tables[1].Name.L, "tbl2")
	}
	if !(flushTable.NoWriteToBinLog) {
		t.Fatal("expected true")
	}
	if !(flushTable.ReadLock) {
		t.Fatal("expected true")
	}
}

func TestFlushPrivileges(t *testing.T) {
	p := parser.New()
	stmt, _, err := p.Parse("flush privileges", "", "")
	if err != nil {
		t.Fatal(err)
	}
	flushPrivilege := stmt[0].(*ast.FlushStmt)
	if !reflect.DeepEqual(ast.FlushPrivileges, flushPrivilege.Tp) {
		t.Fatalf("got %v, want %v", flushPrivilege.Tp, ast.FlushPrivileges)
	}
}

func TestIdentifier(t *testing.T) {
	table := []testCase{
		// for quote identifier
		{"select `a`, `a.b`, `a b` from t", true, "SELECT `a`,`a.b`,`a b` FROM `t`"},
		// for unquoted identifier
		{"create table MergeContextTest$Simple (value integer not null, primary key (value))", true, "CREATE TABLE `MergeContextTest$Simple` (`value` INT NOT NULL,PRIMARY KEY(`value`))"},
		// for as
		{"select 1 as a, 1 as `a`, 1 as \"a\", 1 as 'a'", true, "SELECT 1 AS `a`,1 AS `a`,1 AS `a`,1 AS `a`"},
		{`select 1 as a, 1 as "a", 1 as 'a'`, true, "SELECT 1 AS `a`,1 AS `a`,1 AS `a`"},
		{`select 1 a, 1 "a", 1 'a'`, true, "SELECT 1 AS `a`,1 AS `a`,1 AS `a`"},
		{`select * from t as "a"`, false, ""},
		{`select * from t a`, true, "SELECT * FROM `t` AS `a`"},
		// reserved keyword can't be used as identifier directly, but A.B pattern is an exception
		{`select * from ROW`, false, ""},
		{`select COUNT from DESC`, false, ""},
		{`select COUNT from SELECT.DESC`, true, "SELECT `COUNT` FROM `SELECT`.`DESC`"},
		{"use `select`", true, "USE `select`"},
		{"use `sel``ect`", true, "USE `sel``ect`"}, //nolint: misspell
		{"use select", false, "USE `select`"},
		{`select * from t as a`, true, "SELECT * FROM `t` AS `a`"},
		{"select 1 full, 1 row, 1 abs", false, ""},
		{"select 1 full, 1 `row`, 1 abs", true, "SELECT 1 AS `full`,1 AS `row`,1 AS `abs`"},
		{"select * from t full, t1 row, t2 abs", false, ""},
		{"select * from t full, t1 `row`, t2 abs", true, "SELECT * FROM ((`t` AS `full`) JOIN `t1` AS `row`) JOIN `t2` AS `abs`"},
		// for issue 1878, identifiers may begin with digit.
		{"create database 123test", true, "CREATE DATABASE `123test`"},
		{"create database 123", false, "CREATE DATABASE `123`"},
		{"create database `123`", true, "CREATE DATABASE `123`"},
		{"create database `12``3`", true, "CREATE DATABASE `12``3`"},
		{"create table `123` (123a1 int)", true, "CREATE TABLE `123` (`123a1` INT)"},
		{"create table 123 (123a1 int)", false, ""},
		{fmt.Sprintf("select * from t%cble", 0), false, ""},
		// for issue 3954, should NOT be recognized as identifiers.
		{`select .78+123`, true, "SELECT 0.78+123"},
		{`select .78+.21`, true, "SELECT 0.78+0.21"},
		{`select .78-123`, true, "SELECT 0.78-123"},
		{`select .78-.21`, true, "SELECT 0.78-0.21"},
		{`select .78--123`, true, "SELECT 0.78--123"},
		{`select .78*123`, true, "SELECT 0.78*123"},
		{`select .78*.21`, true, "SELECT 0.78*0.21"},
		{`select .78/123`, true, "SELECT 0.78/123"},
		{`select .78/.21`, true, "SELECT 0.78/0.21"},
		{`select .78,123`, true, "SELECT 0.78,123"},
		{`select .78,.21`, true, "SELECT 0.78,0.21"},
		{`select .78 , 123`, true, "SELECT 0.78,123"},
		{`select .78.123`, false, ""},
		{`select .78#123`, true, "SELECT 0.78"},
		{`insert float_test values(.67, 'string');`, true, "INSERT INTO `float_test` VALUES (0.67,_UTF8MB4'string')"},
		{`select .78'123'`, true, "SELECT 0.78 AS `123`"},
		{"select .78`123`", true, "SELECT 0.78 AS `123`"},
		{`select .78"123"`, true, "SELECT 0.78 AS `123`"},
		{"select 111 as \xd6\xf7", true, "SELECT 111 AS `??`"},
	}
	RunTest(t, table, false, false)
}

func TestBuiltinFuncAsIdentifier(t *testing.T) {
	whitespaceFuncs := []struct {
		funcName string
		args     string
	}{
		{"BIT_AND", "`c1`"},
		{"BIT_OR", "`c1`"},
		{"BIT_XOR", "`c1`"},
		{"CAST", "1 AS FLOAT"},
		{"COUNT", "1"},
		{"CURDATE", ""},
		{"CURTIME", ""},
		{"DATE_ADD", "_UTF8MB4'2011-11-11 10:10:10', INTERVAL 10 SECOND"},
		{"DATE_SUB", "_UTF8MB4'2011-11-11 10:10:10', INTERVAL 10 SECOND"},
		{"EXTRACT", "SECOND FROM _UTF8MB4'2011-11-11 10:10:10'"},
		{"GROUP_CONCAT", "`c2`, `c1` SEPARATOR ','"},
		{"MAX", "`c1`"},
		{"MID", "_UTF8MB4'Sakila', -5, 3"},
		{"MIN", "`c1`"},
		{"NOW", ""},
		{"POSITION", "_UTF8MB4'bar' IN _UTF8MB4'foobarbar'"},
		{"STDDEV_POP", "`c1`"},
		{"STDDEV_SAMP", "`c1`"},
		{"SUBSTR", "_UTF8MB4'Quadratically', 5"},
		{"SUBSTRING", "_UTF8MB4'Quadratically', 5"},
		{"SUM", "`c1`"},
		{"SYSDATE", ""},
		{"TRIM", "_UTF8MB4' foo '"},
		{"VAR_POP", "`c1`"},
		{"VAR_SAMP", "`c1`"},
	}

	testcases := make([]testCase, 0, 3*len(whitespaceFuncs))
	runTests := func(ignoreSpace bool) {
		p := parser.New()
		if ignoreSpace {
			p.SetSQLMode(mysql.ModeIgnoreSpace)
		}
		for _, c := range testcases {
			_, _, err := p.Parse(c.src, "", "")
			if !c.ok {
				if err == nil {
					t.Fatalf("source %v", c.src)
				}
				continue
			}
			if err != nil {
				t.Fatalf("%s: %v", fmt.Sprintf("source %v", c.src), err)
			}
			if c.ok && !ignoreSpace {
				RunRestoreTest(t, c.src, c.restore, false, false)
			}
		}
	}

	for _, function := range whitespaceFuncs {
		// `x` is recognized as a function name for `x()`.
		testcases = append(testcases, testCase{fmt.Sprintf("select %s(%s)", function.funcName, function.args), true, fmt.Sprintf("SELECT %s(%s)", function.funcName, function.args)})

		// In MySQL, `select x ()` is recognized as a stored function.
		// In TiDB, most of these functions are recognized as identifiers while some are builtin functions (such as COUNT, CURDATE)
		// because the later ones are not added to the token map. We'd better not to modify it since it breaks compatibility.
		// For example, `select CURDATE ()` reports an error in MySQL but it works well for TiDB.

		// `x` is recognized as an identifier for `x ()`.
		testcases = append(testcases, testCase{fmt.Sprintf("create table %s (a int)", function.funcName), true, fmt.Sprintf("CREATE TABLE `%s` (`a` INT)", function.funcName)})

		// `x` is recognized as a function name for `x()`.
		testcases = append(testcases, testCase{fmt.Sprintf("create table %s(a int)", function.funcName), false, ""})
	}
	runTests(false)

	testcases = make([]testCase, 0, 4*len(whitespaceFuncs))
	for _, function := range whitespaceFuncs {
		testcases = append(testcases, testCase{fmt.Sprintf("select %s(%s)", function.funcName, function.args), true, fmt.Sprintf("SELECT %s(%s)", function.funcName, function.args)})
		testcases = append(testcases, testCase{fmt.Sprintf("select %s (%s)", function.funcName, function.args), true, fmt.Sprintf("SELECT %s(%s)", function.funcName, function.args)})
		testcases = append(testcases, testCase{fmt.Sprintf("create table %s (a int)", function.funcName), false, ""})
		testcases = append(testcases, testCase{fmt.Sprintf("create table %s(a int)", function.funcName), false, ""})
	}
	runTests(true)

	normalFuncs := []struct {
		funcName string
		args     string
	}{
		{"ADDDATE", "_UTF8MB4'2011-11-11 10:10:10', INTERVAL 10 SECOND"},
		{"SESSION_USER", ""},
		{"SUBDATE", "_UTF8MB4'2011-11-11 10:10:10', INTERVAL 10 SECOND"},
		{"SYSTEM_USER", ""},
	}

	testcases = make([]testCase, 0, 4*len(normalFuncs))
	for _, function := range normalFuncs {
		// `x` is recognized as a function name for `select x()`.
		testcases = append(testcases, testCase{fmt.Sprintf("select %s(%s)", function.funcName, function.args), true, fmt.Sprintf("SELECT %s(%s)", function.funcName, function.args)})

		// `x` is recognized as a function name for `select x ()`.
		testcases = append(testcases, testCase{fmt.Sprintf("select %s (%s)", function.funcName, function.args), true, fmt.Sprintf("SELECT %s(%s)", function.funcName, function.args)})

		// `x` is recognized as an identifier for `create table x ()`.
		testcases = append(testcases, testCase{fmt.Sprintf("create table %s (a int)", function.funcName), true, fmt.Sprintf("CREATE TABLE `%s` (`a` INT)", function.funcName)})

		// `x` is recognized as an identifier for `create table x()`.
		testcases = append(testcases, testCase{fmt.Sprintf("create table %s(a int)", function.funcName), true, fmt.Sprintf("CREATE TABLE `%s` (`a` INT)", function.funcName)})
	}
	runTests(false)
	runTests(true)
}

func TestHintError(t *testing.T) {
	p := parser.New()
	stmt, warns, err := p.Parse("select /*+ tidb_unknown(T1,t2) */ c1, c2 from t1, t2 where t1.c1 = t2.c1", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if got := len(warns); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	if !reflect.DeepEqual(`[parser:8061]Optimizer hint tidb_unknown is not supported by TiDB and is ignored`, warns[0].Error()) {
		t.Fatalf("got %v, want %v", warns[0].Error(), `[parser:8061]Optimizer hint tidb_unknown is not supported by TiDB and is ignored`)
	}
	if got := len(stmt[0].(*ast.SelectStmt).TableHints); got != 0 {
		t.Fatalf("expected length %d, got %d", 0, got)
	}
	stmt, warns, err = p.Parse("select /*+ TIDB_INLJ(t1, T2) tidb_unknown(T1,t2, 1) */ c1, c2 from t1, t2 where t1.c1 = t2.c1", "", "")
	if got := len(stmt[0].(*ast.SelectStmt).TableHints); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	if err != nil {
		t.Fatal(err)
	}
	if got := len(warns); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	if !reflect.DeepEqual(`[parser:8061]Optimizer hint tidb_unknown is not supported by TiDB and is ignored`, warns[0].Error()) {
		t.Fatalf("got %v, want %v", warns[0].Error(), `[parser:8061]Optimizer hint tidb_unknown is not supported by TiDB and is ignored`)
	}
	_, _, err = p.Parse("select c1, c2 from /*+ tidb_unknow(T1,t2) */ t1, t2 where t1.c1 = t2.c1", "", "")
	if err != nil {
		t.Fatal(err)
	} // Hints are ignored after the "FROM" keyword!
	_, _, err = p.Parse("select1 /*+ TIDB_INLJ(t1, T2) */ c1, c2 from t1, t2 where t1.c1 = t2.c1", "", "")
	if err == nil || err.Error() != "line 1 column 7 near \"select1 /*+ TIDB_INLJ(t1, T2) */ c1, c2 from t1, t2 where t1.c1 = t2.c1\" " {
		t.Fatalf("expected error %q, got %v", "line 1 column 7 near \"select1 /*+ TIDB_INLJ(t1, T2) */ c1, c2 from t1, t2 where t1.c1 = t2.c1\" ", err)
	}
	_, _, err = p.Parse("select /*+ TIDB_INLJ(t1, T2) */ c1, c2 fromt t1, t2 where t1.c1 = t2.c1", "", "")
	if err == nil || err.Error() != "line 1 column 47 near \"t1, t2 where t1.c1 = t2.c1\" " {
		t.Fatalf("expected error %q, got %v", "line 1 column 47 near \"t1, t2 where t1.c1 = t2.c1\" ", err)
	}
	_, _, err = p.Parse("SELECT 1 FROM DUAL WHERE 1 IN (SELECT /*+ DEBUG_HINT3 */ 1)", "", "")
	if err != nil {
		t.Fatal(err)
	}
	stmt, _, err = p.Parse("insert into t select /*+ memory_quota(1 MB) */ * from t;", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if got := len(stmt[0].(*ast.InsertStmt).TableHints); got != 0 {
		t.Fatalf("expected length %d, got %d", 0, got)
	}
	if got := len(stmt[0].(*ast.InsertStmt).Select.(*ast.SelectStmt).TableHints); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	stmt, _, err = p.Parse("insert /*+ memory_quota(1 MB) */ into t select * from t;", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if got := len(stmt[0].(*ast.InsertStmt).TableHints); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}

	_, warns, err = p.Parse("SELECT id FROM tbl WHERE id = 0 FOR UPDATE /*+ xyz */", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if got := len(warns); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	if !regexp.MustCompile(`near '/\*\+' at line 1$`).MatchString(warns[0].Error()) {
		t.Fatalf("expected %q to match %q", warns[0].Error(), `near '/\*\+' at line 1$`)
	}

	_, warns, err = p.Parse("create global binding for select /*+ max_execution_time(1) */ 1 using select /*+ max_execution_time(1) */ 1;\n", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if got := len(warns); got != 0 {
		t.Fatalf("expected length %d, got %d", 0, got)
	}
}

func TestErrorMsg(t *testing.T) {
	p := parser.New()
	_, _, err := p.Parse("select1 1", "", "")
	if err == nil || err.Error() != "line 1 column 7 near \"select1 1\" " {
		t.Fatalf("expected error %q, got %v", "line 1 column 7 near \"select1 1\" ", err)
	}
	_, _, err = p.Parse("select 1 from1 dual", "", "")
	if err == nil || err.Error() != "line 1 column 19 near \"dual\" " {
		t.Fatalf("expected error %q, got %v", "line 1 column 19 near \"dual\" ", err)
	}
	_, _, err = p.Parse("select * from t1 join t2 from t1.a = t2.a;", "", "")
	if err == nil || err.Error() != "line 1 column 29 near \"from t1.a = t2.a;\" " {
		t.Fatalf("expected error %q, got %v", "line 1 column 29 near \"from t1.a = t2.a;\" ", err)
	}
	_, _, err = p.Parse("select * from t1 join t2 one t1.a = t2.a;", "", "")
	if err == nil || err.Error() != "line 1 column 31 near \"t1.a = t2.a;\" " {
		t.Fatalf("expected error %q, got %v", "line 1 column 31 near \"t1.a = t2.a;\" ", err)
	}
	_, _, err = p.Parse("select * from t1 join t2 on t1.a >>> t2.a;", "", "")
	if err == nil || err.Error() != "line 1 column 36 near \"> t2.a;\" " {
		t.Fatalf("expected error %q, got %v", "line 1 column 36 near \"> t2.a;\" ", err)
	}

	_, _, err = p.Parse("create table t(f_year year(5))ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;", "", "")
	if err == nil || err.Error() != "[parser:1818]Supports only YEAR or YEAR(4) column" {
		t.Fatalf("expected error %q, got %v", "[parser:1818]Supports only YEAR or YEAR(4) column", err)
	}

	_, _, err = p.Parse("create table ``.t (id int);", "", "")
	if err == nil || err.Error() != "[parser:1102]Incorrect database name ''" {
		t.Fatalf("expected error %q, got %v", "[parser:1102]Incorrect database name ''", err)
	}

	_, _, err = p.Parse("create table ` `.t (id int);", "", "")
	if err == nil || err.Error() != "[parser:1102]Incorrect database name ' '" {
		t.Fatalf("expected error %q, got %v", "[parser:1102]Incorrect database name ' '", err)
	}

	_, _, err = p.Parse("select ifnull(a,0) & ifnull(a,0) like '55' ESCAPE '\\\\a' from t;", "", "")
	if err == nil || err.Error() != "[parser:1210]Incorrect arguments to ESCAPE" {
		t.Fatalf("expected error %q, got %v", "[parser:1210]Incorrect arguments to ESCAPE", err)
	}

	_, _, err = p.Parse("load data infile 'aaa' into table aaa FIELDS  Enclosed by '\\\\b';", "", "")
	if err == nil || err.Error() != "[parser:1083]Field separator argument is not what is expected; check the manual" {
		t.Fatalf("expected error %q, got %v", "[parser:1083]Field separator argument is not what is expected; check the manual", err)
	}

	_, _, err = p.Parse("load data infile 'aaa' into table aaa FIELDS  Escaped by '\\\\b';", "", "")
	if err == nil || err.Error() != "[parser:1083]Field separator argument is not what is expected; check the manual" {
		t.Fatalf("expected error %q, got %v", "[parser:1083]Field separator argument is not what is expected; check the manual", err)
	}

	_, _, err = p.Parse("load data infile 'aaa' into table aaa FIELDS  Enclosed by '\\\\b' Escaped by '\\\\b' ;", "", "")
	if err == nil || err.Error() != "[parser:1083]Field separator argument is not what is expected; check the manual" {
		t.Fatalf("expected error %q, got %v", "[parser:1083]Field separator argument is not what is expected; check the manual", err)
	}

	_, _, err = p.Parse("ALTER DATABASE `` CHARACTER SET = ''", "", "")
	if err == nil || err.Error() != "[parser:1115]Unknown character set: ''" {
		t.Fatalf("expected error %q, got %v", "[parser:1115]Unknown character set: ''", err)
	}

	_, _, err = p.Parse("ALTER DATABASE t CHARACTER SET = ''", "", "")
	if err == nil || err.Error() != "[parser:1115]Unknown character set: ''" {
		t.Fatalf("expected error %q, got %v", "[parser:1115]Unknown character set: ''", err)
	}

	_, _, err = p.Parse("ALTER SCHEMA t CHARACTER SET = 'SOME_INVALID_CHARSET'", "", "")
	if err == nil || err.Error() != "[parser:1115]Unknown character set: 'SOME_INVALID_CHARSET'" {
		t.Fatalf("expected error %q, got %v", "[parser:1115]Unknown character set: 'SOME_INVALID_CHARSET'", err)
	}

	_, _, err = p.Parse("ALTER DATABASE t COLLATE = ''", "", "")
	if err == nil || err.Error() != "[ddl:1273]Unknown collation: ''" {
		t.Fatalf("expected error %q, got %v", "[ddl:1273]Unknown collation: ''", err)
	}

	_, _, err = p.Parse("ALTER SCHEMA t COLLATE = 'SOME_INVALID_COLLATION'", "", "")
	if err == nil || err.Error() != "[ddl:1273]Unknown collation: 'SOME_INVALID_COLLATION'" {
		t.Fatalf("expected error %q, got %v", "[ddl:1273]Unknown collation: 'SOME_INVALID_COLLATION'", err)
	}

	_, _, err = p.Parse("ALTER DATABASE CHARSET = 'utf8mb4' COLLATE = 'utf8_bin'", "", "")
	if err == nil || err.Error() != "line 1 column 24 near \"= 'utf8mb4' COLLATE = 'utf8_bin'\" " {
		t.Fatalf("expected error %q, got %v", "line 1 column 24 near \"= 'utf8mb4' COLLATE = 'utf8_bin'\" ", err)
	}

	_, _, err = p.Parse("ALTER DATABASE t ENCRYPTION = ''", "", "")
	if err == nil || err.Error() != "[parser:1525]Incorrect argument (should be Y or N) value: ''" {
		t.Fatalf("expected error %q, got %v", "[parser:1525]Incorrect argument (should be Y or N) value: ''", err)
	}

	_, _, err = p.Parse("ALTER DATABASE", "", "")
	if err == nil || err.Error() != "line 1 column 14 near \"\" " {
		t.Fatalf("expected error %q, got %v", "line 1 column 14 near \"\" ", err)
	}

	_, _, err = p.Parse("ALTER SCHEMA `ANY_DB_NAME`", "", "")
	if err == nil || err.Error() != "line 1 column 26 near \"\" " {
		t.Fatalf("expected error %q, got %v", "line 1 column 26 near \"\" ", err)
	}

	_, _, err = p.Parse("alter table t partition by range FIELDS(a)", "", "")
	if err == nil || err.Error() != "[ddl:1492]For RANGE partitions each partition must be defined" {
		t.Fatalf("expected error %q, got %v", "[ddl:1492]For RANGE partitions each partition must be defined", err)
	}

	_, _, err = p.Parse("alter table t partition by list FIELDS(a)", "", "")
	if err == nil || err.Error() != "[ddl:1492]For LIST partitions each partition must be defined" {
		t.Fatalf("expected error %q, got %v", "[ddl:1492]For LIST partitions each partition must be defined", err)
	}

	_, _, err = p.Parse("alter table t partition by list FIELDS(a)", "", "")
	if err == nil || err.Error() != "[ddl:1492]For LIST partitions each partition must be defined" {
		t.Fatalf("expected error %q, got %v", "[ddl:1492]For LIST partitions each partition must be defined", err)
	}

	_, _, err = p.Parse("alter table t partition by list FIELDS(a,b,c)", "", "")
	if err == nil || err.Error() != "[ddl:1492]For LIST partitions each partition must be defined" {
		t.Fatalf("expected error %q, got %v", "[ddl:1492]For LIST partitions each partition must be defined", err)
	}

	_, _, err = p.Parse("alter table t lock = first", "", "")
	if err == nil || err.Error() != "[parser:1801]Unknown LOCK type 'first'" {
		t.Fatalf("expected error %q, got %v", "[parser:1801]Unknown LOCK type 'first'", err)
	}

	_, _, err = p.Parse("alter table t lock = start", "", "")
	if err == nil || err.Error() != "[parser:1801]Unknown LOCK type 'start'" {
		t.Fatalf("expected error %q, got %v", "[parser:1801]Unknown LOCK type 'start'", err)
	}

	_, _, err = p.Parse("alter table t lock = commit", "", "")
	if err == nil || err.Error() != "[parser:1801]Unknown LOCK type 'commit'" {
		t.Fatalf("expected error %q, got %v", "[parser:1801]Unknown LOCK type 'commit'", err)
	}

	_, _, err = p.Parse("alter table t lock = binlog", "", "")
	if err == nil || err.Error() != "[parser:1801]Unknown LOCK type 'binlog'" {
		t.Fatalf("expected error %q, got %v", "[parser:1801]Unknown LOCK type 'binlog'", err)
	}

	_, _, err = p.Parse("alter table t lock = randomStr123", "", "")
	if err == nil || err.Error() != "[parser:1801]Unknown LOCK type 'randomStr123'" {
		t.Fatalf("expected error %q, got %v", "[parser:1801]Unknown LOCK type 'randomStr123'", err)
	}

	_, _, err = p.Parse("create table t (a longtext unicode)", "", "")
	if err == nil || err.Error() != "[parser:1115]Unknown character set: 'ucs2'" {
		t.Fatalf("expected error %q, got %v", "[parser:1115]Unknown character set: 'ucs2'", err)
	}

	_, _, err = p.Parse("create table t (a long byte, b text unicode)", "", "")
	if err == nil || err.Error() != "[parser:1115]Unknown character set: 'ucs2'" {
		t.Fatalf("expected error %q, got %v", "[parser:1115]Unknown character set: 'ucs2'", err)
	}

	_, _, err = p.Parse("create table t (a long ascii, b long unicode)", "", "")
	if err == nil || err.Error() != "[parser:1115]Unknown character set: 'ucs2'" {
		t.Fatalf("expected error %q, got %v", "[parser:1115]Unknown character set: 'ucs2'", err)
	}

	_, _, err = p.Parse("create table t (a text unicode, b mediumtext ascii, c int)", "", "")
	if err == nil || err.Error() != "[parser:1115]Unknown character set: 'ucs2'" {
		t.Fatalf("expected error %q, got %v", "[parser:1115]Unknown character set: 'ucs2'", err)
	}

	_, _, err = p.Parse("select 1 collate some_unknown_collation", "", "")
	if err == nil || err.Error() != "[ddl:1273]Unknown collation: 'some_unknown_collation'" {
		t.Fatalf("expected error %q, got %v", "[ddl:1273]Unknown collation: 'some_unknown_collation'", err)
	}
}

func TestOptimizerHints(t *testing.T) {
	p := parser.New()
	// Test USE_INDEX
	stmt, _, err := p.Parse("select /*+ USE_INDEX(T1,T2), use_index(t3,t4) */ c1, c2 from t1, t2 where t1.c1 = t2.c1", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt := stmt[0].(*ast.SelectStmt)

	hints := selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("use_index", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "use_index")
	}
	if got := len(hints[0].Tables); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	if !reflect.DeepEqual("t1", hints[0].Tables[0].TableName.L) {
		t.Fatalf("got %v, want %v", hints[0].Tables[0].TableName.L, "t1")
	}
	if got := len(hints[0].Indexes); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	if !reflect.DeepEqual("t2", hints[0].Indexes[0].L) {
		t.Fatalf("got %v, want %v", hints[0].Indexes[0].L, "t2")
	}

	if !reflect.DeepEqual("use_index", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "use_index")
	}
	if got := len(hints[1].Tables); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	if !reflect.DeepEqual("t3", hints[1].Tables[0].TableName.L) {
		t.Fatalf("got %v, want %v", hints[1].Tables[0].TableName.L, "t3")
	}
	if got := len(hints[1].Indexes); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	if !reflect.DeepEqual("t4", hints[1].Indexes[0].L) {
		t.Fatalf("got %v, want %v", hints[1].Indexes[0].L, "t4")
	}

	// Test FORCE_INDEX
	stmt, _, err = p.Parse("select /*+ FORCE_INDEX(T1,T2), force_index(t3,t4) RESOURCE_GROUP(rg1)*/ c1, c2 from t1, t2 where t1.c1 = t2.c1", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 3 {
		t.Fatalf("expected length %d, got %d", 3, got)
	}
	if !reflect.DeepEqual("force_index", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "force_index")
	}
	if got := len(hints[0].Tables); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	if !reflect.DeepEqual("t1", hints[0].Tables[0].TableName.L) {
		t.Fatalf("got %v, want %v", hints[0].Tables[0].TableName.L, "t1")
	}
	if got := len(hints[0].Indexes); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	if !reflect.DeepEqual("t2", hints[0].Indexes[0].L) {
		t.Fatalf("got %v, want %v", hints[0].Indexes[0].L, "t2")
	}

	if !reflect.DeepEqual("force_index", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "force_index")
	}
	if got := len(hints[1].Tables); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	if !reflect.DeepEqual("t3", hints[1].Tables[0].TableName.L) {
		t.Fatalf("got %v, want %v", hints[1].Tables[0].TableName.L, "t3")
	}
	if got := len(hints[1].Indexes); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	if !reflect.DeepEqual("t4", hints[1].Indexes[0].L) {
		t.Fatalf("got %v, want %v", hints[1].Indexes[0].L, "t4")
	}

	if !reflect.DeepEqual("resource_group", hints[2].HintName.L) {
		t.Fatalf("got %v, want %v", hints[2].HintName.L, "resource_group")
	}
	if !reflect.DeepEqual(hints[2].HintData, "rg1") {
		t.Fatalf("got %v, want %v", "rg1", hints[2].HintData)
	}

	// Test IGNORE_INDEX
	stmt, _, err = p.Parse("select /*+ IGNORE_INDEX(T1,T2), ignore_index(t3,t4) */ c1, c2 from t1, t2 where t1.c1 = t2.c1", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("ignore_index", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "ignore_index")
	}
	if got := len(hints[0].Tables); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	if !reflect.DeepEqual("t1", hints[0].Tables[0].TableName.L) {
		t.Fatalf("got %v, want %v", hints[0].Tables[0].TableName.L, "t1")
	}
	if got := len(hints[0].Indexes); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	if !reflect.DeepEqual("t2", hints[0].Indexes[0].L) {
		t.Fatalf("got %v, want %v", hints[0].Indexes[0].L, "t2")
	}

	if !reflect.DeepEqual("ignore_index", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "ignore_index")
	}
	if got := len(hints[1].Tables); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	if !reflect.DeepEqual("t3", hints[1].Tables[0].TableName.L) {
		t.Fatalf("got %v, want %v", hints[1].Tables[0].TableName.L, "t3")
	}
	if got := len(hints[1].Indexes); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	if !reflect.DeepEqual("t4", hints[1].Indexes[0].L) {
		t.Fatalf("got %v, want %v", hints[1].Indexes[0].L, "t4")
	}

	// Test ORDER_INDEX
	stmt, _, err = p.Parse("select /*+ ORDER_INDEX(T1,T2), order_index(t3,t4) */ c1, c2 from t1, t2 where t1.c1 = t2.c1", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("order_index", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "order_index")
	}
	if got := len(hints[0].Tables); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	if !reflect.DeepEqual("t1", hints[0].Tables[0].TableName.L) {
		t.Fatalf("got %v, want %v", hints[0].Tables[0].TableName.L, "t1")
	}
	if got := len(hints[0].Indexes); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	if !reflect.DeepEqual("t2", hints[0].Indexes[0].L) {
		t.Fatalf("got %v, want %v", hints[0].Indexes[0].L, "t2")
	}

	if !reflect.DeepEqual("order_index", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "order_index")
	}
	if got := len(hints[1].Tables); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	if !reflect.DeepEqual("t3", hints[1].Tables[0].TableName.L) {
		t.Fatalf("got %v, want %v", hints[1].Tables[0].TableName.L, "t3")
	}
	if got := len(hints[1].Indexes); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	if !reflect.DeepEqual("t4", hints[1].Indexes[0].L) {
		t.Fatalf("got %v, want %v", hints[1].Indexes[0].L, "t4")
	}

	// Test NO_ORDER_INDEX
	stmt, _, err = p.Parse("select /*+ NO_ORDER_INDEX(T1,T2), no_order_index(t3,t4) */ c1, c2 from t1, t2 where t1.c1 = t2.c1", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("no_order_index", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "no_order_index")
	}
	if got := len(hints[0].Tables); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	if !reflect.DeepEqual("t1", hints[0].Tables[0].TableName.L) {
		t.Fatalf("got %v, want %v", hints[0].Tables[0].TableName.L, "t1")
	}
	if got := len(hints[0].Indexes); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	if !reflect.DeepEqual("t2", hints[0].Indexes[0].L) {
		t.Fatalf("got %v, want %v", hints[0].Indexes[0].L, "t2")
	}

	if !reflect.DeepEqual("no_order_index", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "no_order_index")
	}
	if got := len(hints[1].Tables); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	if !reflect.DeepEqual("t3", hints[1].Tables[0].TableName.L) {
		t.Fatalf("got %v, want %v", hints[1].Tables[0].TableName.L, "t3")
	}
	if got := len(hints[1].Indexes); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	if !reflect.DeepEqual("t4", hints[1].Indexes[0].L) {
		t.Fatalf("got %v, want %v", hints[1].Indexes[0].L, "t4")
	}

	// Test INDEX_LOOKUP_PUSHDOWN
	stmt, _, err = p.Parse("select /*+ INDEX_LOOKUP_PUSHDOWN(T1,T2), index_lookup_pushdown(t3,t4) */ c1, c2 from t1, t2 where t1.c1 = t2.c1", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("index_lookup_pushdown", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "index_lookup_pushdown")
	}
	if got := len(hints[0].Tables); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	if !reflect.DeepEqual("t1", hints[0].Tables[0].TableName.L) {
		t.Fatalf("got %v, want %v", hints[0].Tables[0].TableName.L, "t1")
	}
	if got := len(hints[0].Indexes); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	if !reflect.DeepEqual("t2", hints[0].Indexes[0].L) {
		t.Fatalf("got %v, want %v", hints[0].Indexes[0].L, "t2")
	}

	if !reflect.DeepEqual("index_lookup_pushdown", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "index_lookup_pushdown")
	}
	if got := len(hints[1].Tables); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	if !reflect.DeepEqual("t3", hints[1].Tables[0].TableName.L) {
		t.Fatalf("got %v, want %v", hints[1].Tables[0].TableName.L, "t3")
	}
	if got := len(hints[1].Indexes); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	if !reflect.DeepEqual("t4", hints[1].Indexes[0].L) {
		t.Fatalf("got %v, want %v", hints[1].Indexes[0].L, "t4")
	}

	// Test TIDB_SMJ
	stmt, _, err = p.Parse("select /*+ TIDB_SMJ(T1,t2), tidb_smj(T3,t4) */ c1, c2 from t1, t2 where t1.c1 = t2.c1", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("tidb_smj", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "tidb_smj")
	}
	if got := len(hints[0].Tables); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("t1", hints[0].Tables[0].TableName.L) {
		t.Fatalf("got %v, want %v", hints[0].Tables[0].TableName.L, "t1")
	}
	if !reflect.DeepEqual("t2", hints[0].Tables[1].TableName.L) {
		t.Fatalf("got %v, want %v", hints[0].Tables[1].TableName.L, "t2")
	}

	if !reflect.DeepEqual("tidb_smj", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "tidb_smj")
	}
	if got := len(hints[1].Tables); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("t3", hints[1].Tables[0].TableName.L) {
		t.Fatalf("got %v, want %v", hints[1].Tables[0].TableName.L, "t3")
	}
	if !reflect.DeepEqual("t4", hints[1].Tables[1].TableName.L) {
		t.Fatalf("got %v, want %v", hints[1].Tables[1].TableName.L, "t4")
	}

	// Test MERGE_JOIN
	stmt, _, err = p.Parse("select /*+ MERGE_JOIN(t1, T2), merge_join(t3, t4) */ c1, c2 from t1, t2 where t1.c1 = t2.c1", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("merge_join", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "merge_join")
	}
	if got := len(hints[0].Tables); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("t1", hints[0].Tables[0].TableName.L) {
		t.Fatalf("got %v, want %v", hints[0].Tables[0].TableName.L, "t1")
	}
	if !reflect.DeepEqual("t2", hints[0].Tables[1].TableName.L) {
		t.Fatalf("got %v, want %v", hints[0].Tables[1].TableName.L, "t2")
	}

	if !reflect.DeepEqual("merge_join", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "merge_join")
	}
	if got := len(hints[1].Tables); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("t3", hints[1].Tables[0].TableName.L) {
		t.Fatalf("got %v, want %v", hints[1].Tables[0].TableName.L, "t3")
	}
	if !reflect.DeepEqual("t4", hints[1].Tables[1].TableName.L) {
		t.Fatalf("got %v, want %v", hints[1].Tables[1].TableName.L, "t4")
	}

	// TEST BROADCAST_JOIN
	stmt, _, err = p.Parse("select /*+ BROADCAST_JOIN(t1, T2), broadcast_join(t3, t4) */ c1, c2 from t1, t2 where t1.c1 = t2.c1", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("broadcast_join", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "broadcast_join")
	}
	if got := len(hints[0].Tables); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("t1", hints[0].Tables[0].TableName.L) {
		t.Fatalf("got %v, want %v", hints[0].Tables[0].TableName.L, "t1")
	}
	if !reflect.DeepEqual("t2", hints[0].Tables[1].TableName.L) {
		t.Fatalf("got %v, want %v", hints[0].Tables[1].TableName.L, "t2")
	}

	if !reflect.DeepEqual("broadcast_join", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "broadcast_join")
	}
	if got := len(hints[1].Tables); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("t3", hints[1].Tables[0].TableName.L) {
		t.Fatalf("got %v, want %v", hints[1].Tables[0].TableName.L, "t3")
	}
	if !reflect.DeepEqual("t4", hints[1].Tables[1].TableName.L) {
		t.Fatalf("got %v, want %v", hints[1].Tables[1].TableName.L, "t4")
	}

	// Test TIDB_INLJ
	stmt, _, err = p.Parse("select /*+ TIDB_INLJ(t1, T2), tidb_inlj(t3, t4) */ c1, c2 from t1, t2 where t1.c1 = t2.c1", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("tidb_inlj", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "tidb_inlj")
	}
	if got := len(hints[0].Tables); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("t1", hints[0].Tables[0].TableName.L) {
		t.Fatalf("got %v, want %v", hints[0].Tables[0].TableName.L, "t1")
	}
	if !reflect.DeepEqual("t2", hints[0].Tables[1].TableName.L) {
		t.Fatalf("got %v, want %v", hints[0].Tables[1].TableName.L, "t2")
	}

	if !reflect.DeepEqual("tidb_inlj", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "tidb_inlj")
	}
	if got := len(hints[1].Tables); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("t3", hints[1].Tables[0].TableName.L) {
		t.Fatalf("got %v, want %v", hints[1].Tables[0].TableName.L, "t3")
	}
	if !reflect.DeepEqual("t4", hints[1].Tables[1].TableName.L) {
		t.Fatalf("got %v, want %v", hints[1].Tables[1].TableName.L, "t4")
	}

	// Test INL_JOIN
	stmt, _, err = p.Parse("select /*+ INL_JOIN(t1, T2), inl_join(t3, t4) */ c1, c2 from t1, t2 where t1.c1 = t2.c1", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("inl_join", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "inl_join")
	}
	if got := len(hints[0].Tables); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("t1", hints[0].Tables[0].TableName.L) {
		t.Fatalf("got %v, want %v", hints[0].Tables[0].TableName.L, "t1")
	}
	if !reflect.DeepEqual("t2", hints[0].Tables[1].TableName.L) {
		t.Fatalf("got %v, want %v", hints[0].Tables[1].TableName.L, "t2")
	}

	if !reflect.DeepEqual("inl_join", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "inl_join")
	}
	if got := len(hints[1].Tables); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("t3", hints[1].Tables[0].TableName.L) {
		t.Fatalf("got %v, want %v", hints[1].Tables[0].TableName.L, "t3")
	}
	if !reflect.DeepEqual("t4", hints[1].Tables[1].TableName.L) {
		t.Fatalf("got %v, want %v", hints[1].Tables[1].TableName.L, "t4")
	}

	// Test INL_HASH_JOIN
	stmt, _, err = p.Parse("select /*+ INL_HASH_JOIN(t1, T2), inl_hash_join(t3, t4) */ c1, c2 from t1, t2 where t1.c1 = t2.c1", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("inl_hash_join", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "inl_hash_join")
	}
	if got := len(hints[0].Tables); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("t1", hints[0].Tables[0].TableName.L) {
		t.Fatalf("got %v, want %v", hints[0].Tables[0].TableName.L, "t1")
	}
	if !reflect.DeepEqual("t2", hints[0].Tables[1].TableName.L) {
		t.Fatalf("got %v, want %v", hints[0].Tables[1].TableName.L, "t2")
	}

	if !reflect.DeepEqual("inl_hash_join", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "inl_hash_join")
	}
	if got := len(hints[1].Tables); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("t3", hints[1].Tables[0].TableName.L) {
		t.Fatalf("got %v, want %v", hints[1].Tables[0].TableName.L, "t3")
	}
	if !reflect.DeepEqual("t4", hints[1].Tables[1].TableName.L) {
		t.Fatalf("got %v, want %v", hints[1].Tables[1].TableName.L, "t4")
	}

	// Test INL_MERGE_JOIN
	stmt, _, err = p.Parse("select /*+ INL_MERGE_JOIN(t1, T2), inl_merge_join(t3, t4) */ c1, c2 from t1, t2 where t1.c1 = t2.c1", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("inl_merge_join", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "inl_merge_join")
	}
	if got := len(hints[0].Tables); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("t1", hints[0].Tables[0].TableName.L) {
		t.Fatalf("got %v, want %v", hints[0].Tables[0].TableName.L, "t1")
	}
	if !reflect.DeepEqual("t2", hints[0].Tables[1].TableName.L) {
		t.Fatalf("got %v, want %v", hints[0].Tables[1].TableName.L, "t2")
	}

	if !reflect.DeepEqual("inl_merge_join", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "inl_merge_join")
	}
	if got := len(hints[1].Tables); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("t3", hints[1].Tables[0].TableName.L) {
		t.Fatalf("got %v, want %v", hints[1].Tables[0].TableName.L, "t3")
	}
	if !reflect.DeepEqual("t4", hints[1].Tables[1].TableName.L) {
		t.Fatalf("got %v, want %v", hints[1].Tables[1].TableName.L, "t4")
	}

	// Test TIDB_HJ
	stmt, _, err = p.Parse("select /*+ TIDB_HJ(t1, T2), tidb_hj(t3, t4) */ c1, c2 from t1, t2 where t1.c1 = t2.c1", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("tidb_hj", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "tidb_hj")
	}
	if got := len(hints[0].Tables); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("t1", hints[0].Tables[0].TableName.L) {
		t.Fatalf("got %v, want %v", hints[0].Tables[0].TableName.L, "t1")
	}
	if !reflect.DeepEqual("t2", hints[0].Tables[1].TableName.L) {
		t.Fatalf("got %v, want %v", hints[0].Tables[1].TableName.L, "t2")
	}

	if !reflect.DeepEqual("tidb_hj", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "tidb_hj")
	}
	if got := len(hints[1].Tables); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("t3", hints[1].Tables[0].TableName.L) {
		t.Fatalf("got %v, want %v", hints[1].Tables[0].TableName.L, "t3")
	}
	if !reflect.DeepEqual("t4", hints[1].Tables[1].TableName.L) {
		t.Fatalf("got %v, want %v", hints[1].Tables[1].TableName.L, "t4")
	}

	// Test HASH_JOIN
	stmt, _, err = p.Parse("select /*+ HASH_JOIN(t1, T2), hash_join(t3, t4) */ c1, c2 from t1, t2 where t1.c1 = t2.c1", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("hash_join", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "hash_join")
	}
	if got := len(hints[0].Tables); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("t1", hints[0].Tables[0].TableName.L) {
		t.Fatalf("got %v, want %v", hints[0].Tables[0].TableName.L, "t1")
	}
	if !reflect.DeepEqual("t2", hints[0].Tables[1].TableName.L) {
		t.Fatalf("got %v, want %v", hints[0].Tables[1].TableName.L, "t2")
	}

	if !reflect.DeepEqual("hash_join", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "hash_join")
	}
	if got := len(hints[1].Tables); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("t3", hints[1].Tables[0].TableName.L) {
		t.Fatalf("got %v, want %v", hints[1].Tables[0].TableName.L, "t3")
	}
	if !reflect.DeepEqual("t4", hints[1].Tables[1].TableName.L) {
		t.Fatalf("got %v, want %v", hints[1].Tables[1].TableName.L, "t4")
	}

	// Test HASH_JOIN_BUILD and HASH_JOIN_PROBE
	stmt, _, err = p.Parse("select /*+ hash_join_build(t1), hash_join_probe(t4) */ c1, c2 from t1, t2 where t1.c1 = t2.c1", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("hash_join_build", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "hash_join_build")
	}
	if got := len(hints[0].Tables); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	if !reflect.DeepEqual("t1", hints[0].Tables[0].TableName.L) {
		t.Fatalf("got %v, want %v", hints[0].Tables[0].TableName.L, "t1")
	}

	if !reflect.DeepEqual("hash_join_probe", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "hash_join_probe")
	}
	if got := len(hints[1].Tables); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	if !reflect.DeepEqual("t4", hints[1].Tables[0].TableName.L) {
		t.Fatalf("got %v, want %v", hints[1].Tables[0].TableName.L, "t4")
	}

	// Test HASH_JOIN with SWAP_JOIN_INPUTS/NO_SWAP_JOIN_INPUTS
	// t1 for build, t4 for probe
	stmt, _, err = p.Parse("select /*+ HASH_JOIN(t1, T2), hash_join(t3, t4), SWAP_JOIN_INPUTS(t1), NO_SWAP_JOIN_INPUTS(t4)  */ c1, c2 from t1, t2 where t1.c1 = t2.c1", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 4 {
		t.Fatalf("expected length %d, got %d", 4, got)
	}
	if !reflect.DeepEqual("hash_join", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "hash_join")
	}
	if got := len(hints[0].Tables); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("t1", hints[0].Tables[0].TableName.L) {
		t.Fatalf("got %v, want %v", hints[0].Tables[0].TableName.L, "t1")
	}
	if !reflect.DeepEqual("t2", hints[0].Tables[1].TableName.L) {
		t.Fatalf("got %v, want %v", hints[0].Tables[1].TableName.L, "t2")
	}

	if !reflect.DeepEqual("hash_join", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "hash_join")
	}
	if got := len(hints[1].Tables); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("t3", hints[1].Tables[0].TableName.L) {
		t.Fatalf("got %v, want %v", hints[1].Tables[0].TableName.L, "t3")
	}
	if !reflect.DeepEqual("t4", hints[1].Tables[1].TableName.L) {
		t.Fatalf("got %v, want %v", hints[1].Tables[1].TableName.L, "t4")
	}

	if !reflect.DeepEqual("swap_join_inputs", hints[2].HintName.L) {
		t.Fatalf("got %v, want %v", hints[2].HintName.L, "swap_join_inputs")
	}
	if got := len(hints[2].Tables); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	if !reflect.DeepEqual("t1", hints[2].Tables[0].TableName.L) {
		t.Fatalf("got %v, want %v", hints[2].Tables[0].TableName.L, "t1")
	}

	if !reflect.DeepEqual("no_swap_join_inputs", hints[3].HintName.L) {
		t.Fatalf("got %v, want %v", hints[3].HintName.L, "no_swap_join_inputs")
	}
	if got := len(hints[3].Tables); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	if !reflect.DeepEqual("t4", hints[3].Tables[0].TableName.L) {
		t.Fatalf("got %v, want %v", hints[3].Tables[0].TableName.L, "t4")
	}

	// Test MAX_EXECUTION_TIME
	queries := []string{
		"SELECT /*+ MAX_EXECUTION_TIME(1000) */ * FROM t1 INNER JOIN t2 where t1.c1 = t2.c1",
		"SELECT /*+ MAX_EXECUTION_TIME(1000) */ 1",
		"SELECT /*+ MAX_EXECUTION_TIME(1000) */ SLEEP(20)",
		"SELECT /*+ MAX_EXECUTION_TIME(1000) */ 1 FROM DUAL",
	}
	for i, query := range queries {
		stmt, _, err = p.Parse(query, "", "")
		if err != nil {
			t.Fatal(err)
		}
		selectStmt = stmt[0].(*ast.SelectStmt)
		hints = selectStmt.TableHints
		if got := len(hints); got != 1 {
			t.Fatalf("expected length %d, got %d", 1, got)
		}
		if !reflect.DeepEqual("max_execution_time", hints[0].HintName.L) {
			t.Fatalf("case %d: got %v, want %v", i, hints[0].HintName.L, "max_execution_time")
		}
		if !reflect.DeepEqual(uint64(1000), hints[0].HintData.(uint64)) {
			t.Fatalf("got %v, want %v", hints[0].HintData.(uint64), uint64(1000))
		}
	}

	// Test NTH_PLAN
	queries = []string{
		"SELECT /*+ NTH_PLAN(10) */ * FROM t1 INNER JOIN t2 where t1.c1 = t2.c1",
		"SELECT /*+ NTH_PLAN(10) */ 1",
		"SELECT /*+ NTH_PLAN(10) */ SLEEP(20)",
		"SELECT /*+ NTH_PLAN(10) */ 1 FROM DUAL",
	}
	for i, query := range queries {
		stmt, _, err = p.Parse(query, "", "")
		if err != nil {
			t.Fatal(err)
		}
		selectStmt = stmt[0].(*ast.SelectStmt)
		hints = selectStmt.TableHints
		if got := len(hints); got != 1 {
			t.Fatalf("expected length %d, got %d", 1, got)
		}
		if !reflect.DeepEqual("nth_plan", hints[0].HintName.L) {
			t.Fatalf("case %d: got %v, want %v", i, hints[0].HintName.L, "nth_plan")
		}
		if !reflect.DeepEqual(int64(10), hints[0].HintData.(int64)) {
			t.Fatalf("got %v, want %v", hints[0].HintData.(int64), int64(10))
		}
	}

	// Test USE_INDEX_MERGE
	stmt, _, err = p.Parse("select /*+ USE_INDEX_MERGE(t1, c1), use_index_merge(t2, c1), use_index_merge(t3, c1, primary, c2) */ c1, c2 from t1, t2, t3 where t1.c1 = t2.c1 and t3.c2 = t1.c2", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 3 {
		t.Fatalf("expected length %d, got %d", 3, got)
	}
	if !reflect.DeepEqual("use_index_merge", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "use_index_merge")
	}
	if got := len(hints[0].Tables); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	if !reflect.DeepEqual("t1", hints[0].Tables[0].TableName.L) {
		t.Fatalf("got %v, want %v", hints[0].Tables[0].TableName.L, "t1")
	}
	if got := len(hints[0].Indexes); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	if !reflect.DeepEqual("c1", hints[0].Indexes[0].L) {
		t.Fatalf("got %v, want %v", hints[0].Indexes[0].L, "c1")
	}

	if !reflect.DeepEqual("use_index_merge", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "use_index_merge")
	}
	if got := len(hints[1].Tables); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	if !reflect.DeepEqual("t2", hints[1].Tables[0].TableName.L) {
		t.Fatalf("got %v, want %v", hints[1].Tables[0].TableName.L, "t2")
	}
	if got := len(hints[1].Indexes); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	if !reflect.DeepEqual("c1", hints[1].Indexes[0].L) {
		t.Fatalf("got %v, want %v", hints[1].Indexes[0].L, "c1")
	}

	if !reflect.DeepEqual("use_index_merge", hints[2].HintName.L) {
		t.Fatalf("got %v, want %v", hints[2].HintName.L, "use_index_merge")
	}
	if got := len(hints[2].Tables); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	if !reflect.DeepEqual("t3", hints[2].Tables[0].TableName.L) {
		t.Fatalf("got %v, want %v", hints[2].Tables[0].TableName.L, "t3")
	}
	if got := len(hints[2].Indexes); got != 3 {
		t.Fatalf("expected length %d, got %d", 3, got)
	}
	if !reflect.DeepEqual("c1", hints[2].Indexes[0].L) {
		t.Fatalf("got %v, want %v", hints[2].Indexes[0].L, "c1")
	}
	if !reflect.DeepEqual("primary", hints[2].Indexes[1].L) {
		t.Fatalf("got %v, want %v", hints[2].Indexes[1].L, "primary")
	}
	if !reflect.DeepEqual("c2", hints[2].Indexes[2].L) {
		t.Fatalf("got %v, want %v", hints[2].Indexes[2].L, "c2")
	}

	// Test READ_FROM_STORAGE
	stmt, _, err = p.Parse("select /*+ READ_FROM_STORAGE(tiflash[t1, t2], tikv[t3]) */ c1, c2 from t1, t2, t1 t3 where t1.c1 = t2.c1 and t2.c1 = t3.c1", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("read_from_storage", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "read_from_storage")
	}
	if !reflect.DeepEqual("tiflash", hints[0].HintData.(ast.CIStr).L) {
		t.Fatalf("got %v, want %v", hints[0].HintData.(ast.CIStr).L, "tiflash")
	}
	if got := len(hints[0].Tables); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("t1", hints[0].Tables[0].TableName.L) {
		t.Fatalf("got %v, want %v", hints[0].Tables[0].TableName.L, "t1")
	}
	if !reflect.DeepEqual("t2", hints[0].Tables[1].TableName.L) {
		t.Fatalf("got %v, want %v", hints[0].Tables[1].TableName.L, "t2")
	}
	if !reflect.DeepEqual("read_from_storage", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "read_from_storage")
	}
	if !reflect.DeepEqual("tikv", hints[1].HintData.(ast.CIStr).L) {
		t.Fatalf("got %v, want %v", hints[1].HintData.(ast.CIStr).L, "tikv")
	}
	if got := len(hints[1].Tables); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	if !reflect.DeepEqual("t3", hints[1].Tables[0].TableName.L) {
		t.Fatalf("got %v, want %v", hints[1].Tables[0].TableName.L, "t3")
	}

	// Test USE_TOJA
	stmt, _, err = p.Parse("select /*+ USE_TOJA(true), use_toja(false) */ c1, c2 from t1, t2 where t1.c1 = t2.c1", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("use_toja", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "use_toja")
	}
	if !(hints[0].HintData.(bool)) {
		t.Fatal("expected true")
	}

	if !reflect.DeepEqual("use_toja", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "use_toja")
	}
	if hints[1].HintData.(bool) {
		t.Fatal("expected false")
	}

	// Test IGNORE_PLAN_CACHE
	stmt, _, err = p.Parse("select /*+ IGNORE_PLAN_CACHE(), ignore_plan_cache() */ c1, c2 from t1, t2 where t1.c1 = t2.c1", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)
	hints = selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("ignore_plan_cache", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "ignore_plan_cache")
	}
	if !reflect.DeepEqual("ignore_plan_cache", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "ignore_plan_cache")
	}

	stmt, _, err = p.Parse("delete /*+ IGNORE_PLAN_CACHE(), ignore_plan_cache() */ from t where a = 1", "", "")
	if err != nil {
		t.Fatal(err)
	}
	deleteStmt := stmt[0].(*ast.DeleteStmt)
	hints = deleteStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("ignore_plan_cache", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "ignore_plan_cache")
	}
	if !reflect.DeepEqual("ignore_plan_cache", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "ignore_plan_cache")
	}

	stmt, _, err = p.Parse("update /*+  IGNORE_PLAN_CACHE(), ignore_plan_cache() */ t set a = 1 where a = 10", "", "")
	if err != nil {
		t.Fatal(err)
	}
	updateStmt := stmt[0].(*ast.UpdateStmt)
	hints = updateStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("ignore_plan_cache", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "ignore_plan_cache")
	}
	if !reflect.DeepEqual("ignore_plan_cache", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "ignore_plan_cache")
	}

	// Test WRITE_SLOW_LOG
	stmt, _, err = p.Parse("select /*+ WRITE_SLOW_LOG(), write_slow_log() */ c1, c2 from t1, t2 where t1.c1 = t2.c1", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)
	hints = selectStmt.TableHints
	if got := len(hints); got != 0 {
		t.Fatalf("expected length %d, got %d", 0, got)
	}

	stmt, _, err = p.Parse("select /*+ WRITE_SLOW_LOG, write_slow_log*/ c1, c2 from t1, t2 where t1.c1 = t2.c1", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)
	hints = selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("write_slow_log", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "write_slow_log")
	}
	if !reflect.DeepEqual("write_slow_log", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "write_slow_log")
	}

	// Test USE_CASCADES
	stmt, _, err = p.Parse("select /*+ USE_CASCADES(true), use_cascades(false) */ c1, c2 from t1, t2 where t1.c1 = t2.c1", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("use_cascades", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "use_cascades")
	}
	if !(hints[0].HintData.(bool)) {
		t.Fatal("expected true")
	}

	if !reflect.DeepEqual("use_cascades", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "use_cascades")
	}
	if hints[1].HintData.(bool) {
		t.Fatal("expected false")
	}

	// Test USE_PLAN_CACHE
	stmt, _, err = p.Parse("select /*+ USE_PLAN_CACHE(), use_plan_cache() */ c1, c2 from t1, t2 where t1.c1 = t2.c1", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("use_plan_cache", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "use_plan_cache")
	}
	if !reflect.DeepEqual("use_plan_cache", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "use_plan_cache")
	}

	// Test QUERY_TYPE
	stmt, _, err = p.Parse("select /*+ QUERY_TYPE(OLAP), query_type(OLTP) */ c1, c2 from t1, t2 where t1.c1 = t2.c1", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("query_type", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "query_type")
	}
	if !reflect.DeepEqual("olap", hints[0].HintData.(ast.CIStr).L) {
		t.Fatalf("got %v, want %v", hints[0].HintData.(ast.CIStr).L, "olap")
	}
	if !reflect.DeepEqual("query_type", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "query_type")
	}
	if !reflect.DeepEqual("oltp", hints[1].HintData.(ast.CIStr).L) {
		t.Fatalf("got %v, want %v", hints[1].HintData.(ast.CIStr).L, "oltp")
	}

	// Test MEMORY_QUOTA
	stmt, _, err = p.Parse("select /*+ MEMORY_QUOTA(1 MB), memory_quota(1 GB) */ c1, c2 from t1, t2 where t1.c1 = t2.c1", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("memory_quota", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "memory_quota")
	}
	if !reflect.DeepEqual(int64(1024*1024), hints[0].HintData.(int64)) {
		t.Fatalf("got %v, want %v", hints[0].HintData.(int64), int64(1024*1024))
	}
	if !reflect.DeepEqual("memory_quota", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "memory_quota")
	}
	if !reflect.DeepEqual(int64(1024*1024*1024), hints[1].HintData.(int64)) {
		t.Fatalf("got %v, want %v", hints[1].HintData.(int64), int64(1024*1024*1024))
	}

	_, _, err = p.Parse("select /*+ MEMORY_QUOTA(18446744073709551612 MB), memory_quota(8689934592 GB) */ 1", "", "")
	if err != nil {
		t.Fatal(err)
	}

	// Test HASH_AGG
	stmt, _, err = p.Parse("select /*+ HASH_AGG(), hash_agg() */ c1, c2 from t1, t2 where t1.c1 = t2.c1", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("hash_agg", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "hash_agg")
	}
	if !reflect.DeepEqual("hash_agg", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "hash_agg")
	}

	// Test MPPAgg
	stmt, _, err = p.Parse("select /*+ MPP_1PHASE_AGG(), mpp_1phase_agg() */ c1, c2 from t1, t2 where t1.c1 = t2.c1", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("mpp_1phase_agg", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "mpp_1phase_agg")
	}
	if !reflect.DeepEqual("mpp_1phase_agg", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "mpp_1phase_agg")
	}

	stmt, _, err = p.Parse("select /*+ MPP_2PHASE_AGG(), mpp_2phase_agg() */ c1, c2 from t1, t2 where t1.c1 = t2.c1", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("mpp_2phase_agg", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "mpp_2phase_agg")
	}
	if !reflect.DeepEqual("mpp_2phase_agg", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "mpp_2phase_agg")
	}

	// Test ShuffleJoin
	stmt, _, err = p.Parse("select /*+ SHUFFLE_JOIN(t1, t2), shuffle_join(t1, t2) */ * from t1, t2 where t1.c1 = t2.c1", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("shuffle_join", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "shuffle_join")
	}
	if !reflect.DeepEqual("shuffle_join", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "shuffle_join")
	}

	// Test STREAM_AGG
	stmt, _, err = p.Parse("select /*+ STREAM_AGG(), stream_agg() */ c1, c2 from t1, t2 where t1.c1 = t2.c1", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("stream_agg", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "stream_agg")
	}
	if !reflect.DeepEqual("stream_agg", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "stream_agg")
	}

	// Test AGG_TO_COP
	stmt, _, err = p.Parse("select /*+ AGG_TO_COP(), agg_to_cop() */ c1, c2 from t1, t2 where t1.c1 = t2.c1", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("agg_to_cop", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "agg_to_cop")
	}
	if !reflect.DeepEqual("agg_to_cop", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "agg_to_cop")
	}

	// Test NO_INDEX_MERGE
	stmt, _, err = p.Parse("select /*+ NO_INDEX_MERGE(), no_index_merge() */ c1, c2 from t1, t2 where t1.c1 = t2.c1", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("no_index_merge", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "no_index_merge")
	}
	if !reflect.DeepEqual("no_index_merge", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "no_index_merge")
	}

	// Test READ_CONSISTENT_REPLICA
	stmt, _, err = p.Parse("select /*+ READ_CONSISTENT_REPLICA(), read_consistent_replica() */ c1, c2 from t1, t2 where t1.c1 = t2.c1", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("read_consistent_replica", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "read_consistent_replica")
	}
	if !reflect.DeepEqual("read_consistent_replica", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "read_consistent_replica")
	}

	// Test LIMIT_TO_COP
	stmt, _, err = p.Parse("select /*+ LIMIT_TO_COP(), limit_to_cop() */ c1, c2 from t1, t2 where t1.c1 = t2.c1", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("limit_to_cop", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "limit_to_cop")
	}
	if !reflect.DeepEqual("limit_to_cop", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "limit_to_cop")
	}

	// Test CTE MERGE
	stmt, _, err = p.Parse("with cte(x) as (select * from t1) select /*+ MERGE(), merge() */ * from cte;", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("merge", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "merge")
	}
	if !reflect.DeepEqual("merge", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "merge")
	}

	// Test STRAIGHT_JOIN
	stmt, _, err = p.Parse("select /*+ STRAIGHT_JOIN(), straight_join() */ c1, c2 from t1, t2 where t1.c1 = t2.c1", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("straight_join", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "straight_join")
	}
	if !reflect.DeepEqual("straight_join", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "straight_join")
	}

	// Test LEADING
	stmt, _, err = p.Parse("select /*+ LEADING(T1), LEADING(t2, t3), LEADING(T4, t5, t6) */ c1, c2 from t1, t2 where t1.c1 = t2.c1", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 3 {
		t.Fatalf("expected length %d, got %d", 3, got)
	}
	if !reflect.DeepEqual("leading", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "leading")
	}
	leadingList1, ok := hints[0].HintData.(*ast.LeadingList)
	if !(ok) {
		t.Fatal("expected true")
	}
	if got := len(leadingList1.Items); got != 1 {
		t.Fatalf("expected length %d, got %d", 1, got)
	}
	hintTable1, ok := leadingList1.Items[0].(*ast.HintTable)
	if !(ok) {
		t.Fatal("expected true")
	}
	if !reflect.DeepEqual("t1", hintTable1.TableName.L) {
		t.Fatalf("got %v, want %v", hintTable1.TableName.L, "t1")
	}

	if !reflect.DeepEqual("leading", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "leading")
	}
	leadingList2, ok := hints[1].HintData.(*ast.LeadingList)
	if !(ok) {
		t.Fatal("expected true")
	}
	if got := len(leadingList2.Items); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	hintTable2, ok := leadingList2.Items[0].(*ast.HintTable)
	if !(ok) {
		t.Fatal("expected true")
	}
	if !reflect.DeepEqual("t2", hintTable2.TableName.L) {
		t.Fatalf("got %v, want %v", hintTable2.TableName.L, "t2")
	}
	hintTable3, ok := leadingList2.Items[1].(*ast.HintTable)
	if !(ok) {
		t.Fatal("expected true")
	}
	if !reflect.DeepEqual("t3", hintTable3.TableName.L) {
		t.Fatalf("got %v, want %v", hintTable3.TableName.L, "t3")
	}

	if !reflect.DeepEqual("leading", hints[2].HintName.L) {
		t.Fatalf("got %v, want %v", hints[2].HintName.L, "leading")
	}
	leadingList3, ok := hints[2].HintData.(*ast.LeadingList)
	if !(ok) {
		t.Fatal("expected true")
	}
	if got := len(leadingList3.Items); got != 3 {
		t.Fatalf("expected length %d, got %d", 3, got)
	}
	hintTable4, ok := leadingList3.Items[0].(*ast.HintTable)
	if !(ok) {
		t.Fatal("expected true")
	}
	if !reflect.DeepEqual("t4", hintTable4.TableName.L) {
		t.Fatalf("got %v, want %v", hintTable4.TableName.L, "t4")
	}
	hintTable5, ok := leadingList3.Items[1].(*ast.HintTable)
	if !(ok) {
		t.Fatal("expected true")
	}
	if !reflect.DeepEqual("t5", hintTable5.TableName.L) {
		t.Fatalf("got %v, want %v", hintTable5.TableName.L, "t5")
	}
	hintTable6, ok := leadingList3.Items[2].(*ast.HintTable)
	if !(ok) {
		t.Fatal("expected true")
	}
	if !reflect.DeepEqual("t6", hintTable6.TableName.L) {
		t.Fatalf("got %v, want %v", hintTable6.TableName.L, "t6")
	}

	// Test NO_HASH_JOIN
	stmt, _, err = p.Parse("select /*+ NO_HASH_JOIN(t1, t2), NO_HASH_JOIN(t3) */ * from t1, t2, t3", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("no_hash_join", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "no_hash_join")
	}
	if !reflect.DeepEqual(hints[0].Tables[0].TableName.L, "t1") {
		t.Fatalf("got %v, want %v", "t1", hints[0].Tables[0].TableName.L)
	}
	if !reflect.DeepEqual(hints[0].Tables[1].TableName.L, "t2") {
		t.Fatalf("got %v, want %v", "t2", hints[0].Tables[1].TableName.L)
	}

	if !reflect.DeepEqual("no_hash_join", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "no_hash_join")
	}
	if !reflect.DeepEqual(hints[1].Tables[0].TableName.L, "t3") {
		t.Fatalf("got %v, want %v", "t3", hints[1].Tables[0].TableName.L)
	}

	// Test NO_MERGE_JOIN
	stmt, _, err = p.Parse("select /*+ NO_MERGE_JOIN(t1), NO_MERGE_JOIN(t3) */ * from t1, t2, t3", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("no_merge_join", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "no_merge_join")
	}
	if !reflect.DeepEqual(hints[0].Tables[0].TableName.L, "t1") {
		t.Fatalf("got %v, want %v", "t1", hints[0].Tables[0].TableName.L)
	}

	if !reflect.DeepEqual("no_merge_join", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "no_merge_join")
	}
	if !reflect.DeepEqual(hints[1].Tables[0].TableName.L, "t3") {
		t.Fatalf("got %v, want %v", "t3", hints[1].Tables[0].TableName.L)
	}

	// Test INDEX_JOIN
	stmt, _, err = p.Parse("select /*+ INDEX_JOIN(t1), INDEX_JOIN(t3) */ * from t1, t2, t3", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("index_join", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "index_join")
	}
	if !reflect.DeepEqual(hints[0].Tables[0].TableName.L, "t1") {
		t.Fatalf("got %v, want %v", "t1", hints[0].Tables[0].TableName.L)
	}

	if !reflect.DeepEqual("index_join", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "index_join")
	}
	if !reflect.DeepEqual(hints[1].Tables[0].TableName.L, "t3") {
		t.Fatalf("got %v, want %v", "t3", hints[1].Tables[0].TableName.L)
	}

	// Test NO_INDEX_JOIN
	stmt, _, err = p.Parse("select /*+ NO_INDEX_JOIN(t1), NO_INDEX_JOIN(t3) */ * from t1, t2, t3", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("no_index_join", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "no_index_join")
	}
	if !reflect.DeepEqual(hints[0].Tables[0].TableName.L, "t1") {
		t.Fatalf("got %v, want %v", "t1", hints[0].Tables[0].TableName.L)
	}

	if !reflect.DeepEqual("no_index_join", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "no_index_join")
	}
	if !reflect.DeepEqual(hints[1].Tables[0].TableName.L, "t3") {
		t.Fatalf("got %v, want %v", "t3", hints[1].Tables[0].TableName.L)
	}

	// Test INDEX_HASH_JOIN
	stmt, _, err = p.Parse("select /*+ INDEX_HASH_JOIN(t1), INDEX_HASH_JOIN(t3) */ * from t1, t2, t3", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("index_hash_join", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "index_hash_join")
	}
	if !reflect.DeepEqual(hints[0].Tables[0].TableName.L, "t1") {
		t.Fatalf("got %v, want %v", "t1", hints[0].Tables[0].TableName.L)
	}

	if !reflect.DeepEqual("index_hash_join", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "index_hash_join")
	}
	if !reflect.DeepEqual(hints[1].Tables[0].TableName.L, "t3") {
		t.Fatalf("got %v, want %v", "t3", hints[1].Tables[0].TableName.L)
	}

	// Test NO_INDEX_HASH_JOIN
	stmt, _, err = p.Parse("select /*+ NO_INDEX_HASH_JOIN(t1), NO_INDEX_HASH_JOIN(t3) */ * from t1, t2, t3", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("no_index_hash_join", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "no_index_hash_join")
	}
	if !reflect.DeepEqual(hints[0].Tables[0].TableName.L, "t1") {
		t.Fatalf("got %v, want %v", "t1", hints[0].Tables[0].TableName.L)
	}

	if !reflect.DeepEqual("no_index_hash_join", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "no_index_hash_join")
	}
	if !reflect.DeepEqual(hints[1].Tables[0].TableName.L, "t3") {
		t.Fatalf("got %v, want %v", "t3", hints[1].Tables[0].TableName.L)
	}

	// Test INDEX_MERGE_JOIN
	stmt, _, err = p.Parse("select /*+ INDEX_MERGE_JOIN(t1), INDEX_MERGE_JOIN(t3) */ * from t1, t2, t3", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("index_merge_join", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "index_merge_join")
	}
	if !reflect.DeepEqual(hints[0].Tables[0].TableName.L, "t1") {
		t.Fatalf("got %v, want %v", "t1", hints[0].Tables[0].TableName.L)
	}

	if !reflect.DeepEqual("index_merge_join", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "index_merge_join")
	}
	if !reflect.DeepEqual(hints[1].Tables[0].TableName.L, "t3") {
		t.Fatalf("got %v, want %v", "t3", hints[1].Tables[0].TableName.L)
	}

	// Test NO_INDEX_MERGE_JOIN
	stmt, _, err = p.Parse("select /*+ NO_INDEX_MERGE_JOIN(t1), NO_INDEX_MERGE_JOIN(t3) */ * from t1, t2, t3", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("no_index_merge_join", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "no_index_merge_join")
	}
	if !reflect.DeepEqual(hints[0].Tables[0].TableName.L, "t1") {
		t.Fatalf("got %v, want %v", "t1", hints[0].Tables[0].TableName.L)
	}

	if !reflect.DeepEqual("no_index_merge_join", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "no_index_merge_join")
	}
	if !reflect.DeepEqual(hints[1].Tables[0].TableName.L, "t3") {
		t.Fatalf("got %v, want %v", "t3", hints[1].Tables[0].TableName.L)
	}

	// Test HYPO_INDEX
	stmt, _, err = p.Parse("select /*+ HYPO_INDEX(t1, a), HYPO_INDEX(t3, a, b, c) */ * from t1, t2, t3", "", "")
	if err != nil {
		t.Fatal(err)
	}
	selectStmt = stmt[0].(*ast.SelectStmt)

	hints = selectStmt.TableHints
	if got := len(hints); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual("hypo_index", hints[0].HintName.L) {
		t.Fatalf("got %v, want %v", hints[0].HintName.L, "hypo_index")
	}
	if !reflect.DeepEqual(hints[0].Tables[0].TableName.L, "t1") {
		t.Fatalf("got %v, want %v", "t1", hints[0].Tables[0].TableName.L)
	}

	if !reflect.DeepEqual("hypo_index", hints[1].HintName.L) {
		t.Fatalf("got %v, want %v", hints[1].HintName.L, "hypo_index")
	}
	if !reflect.DeepEqual(hints[1].Tables[0].TableName.L, "t3") {
		t.Fatalf("got %v, want %v", "t3", hints[1].Tables[0].TableName.L)
	}
}

func TestParserErrMsg(t *testing.T) {
	commentMsgCases := []testErrMsgCase{
		{"delete from t where a = 7 or 1=1/*' and b = 'p'", errors.New("near '/*' and b = 'p'' at line 1")},
		{"delete from t where a = 7 or\n 1=1/*' and b = 'p'", errors.New("near '/*' and b = 'p'' at line 2")},
		{"select 1/*", errors.New("near '/*' at line 1")},
		{"select 1/* comment */", nil},
	}
	funcCallMsgCases := []testErrMsgCase{
		{"select a.b()", nil},
		{"SELECT foo.bar('baz');", nil},
	}
	RunErrMsgTest(t, commentMsgCases)
	RunErrMsgTest(t, funcCallMsgCases)
}

type subqueryChecker struct {
	text string
	t    *testing.T
}

// Enter implements ast.Visitor interface.
func (sc *subqueryChecker) Enter(inNode ast.Node) (outNode ast.Node, skipChildren bool) {
	if expr, ok := inNode.(*ast.SubqueryExpr); ok {
		if !reflect.DeepEqual(sc.text, expr.Query.Text()) {
			sc.t.Fatalf("got %v, want %v", expr.Query.Text(), sc.text)
		}
		return inNode, true
	}
	return inNode, false
}

// Leave implements ast.Visitor interface.
func (sc *subqueryChecker) Leave(inNode ast.Node) (node ast.Node, ok bool) {
	return inNode, true
}

func TestSubquery(t *testing.T) {
	tests := []struct {
		input string
		text  string
	}{
		{"SELECT 1 > (select 1)", "select 1"},
		{"SELECT 1 > (select 1 union select 2)", "select 1 union select 2"},
	}
	p := parser.New()
	for _, tbl := range tests {
		stmt, err := p.ParseOneStmt(tbl.input, "", "")
		if err != nil {
			t.Fatal(err)
		}
		stmt.Accept(&subqueryChecker{
			text: tbl.text,
			t:    t,
		})
	}
}

func checkOrderBy(t *testing.T, s ast.Node, hasOrderBy []bool, i int) int {
	switch x := s.(type) {
	case *ast.SelectStmt:
		if !reflect.DeepEqual(hasOrderBy[i], x.OrderBy != nil) {
			t.Fatalf("got %v, want %v", x.OrderBy != nil, hasOrderBy[i])
		}
		return i + 1
	case *ast.SetOprSelectList:
		for _, sel := range x.Selects {
			i = checkOrderBy(t, sel, hasOrderBy, i)
		}
		return i
	}
	return i
}

func TestUnionOrderBy(t *testing.T) {
	p := parser.New()
	p.EnableWindowFunc(false)

	tests := []struct {
		src        string
		hasOrderBy []bool
	}{
		{"select 2 as a from dual union select 1 as b from dual order by a", []bool{false, false, true}},
		{"select 2 as a from dual union (select 1 as b from dual order by a)", []bool{false, true, false}},
		{"(select 2 as a from dual order by a) union select 1 as b from dual order by a", []bool{true, false, true}},
		{"select 1 a, 2 b from dual order by a", []bool{true}},
		{"select 1 a, 2 b from dual", []bool{false}},
	}

	for _, tbl := range tests {
		stmt, _, err := p.Parse(tbl.src, "", "")
		if err != nil {
			t.Fatal(err)
		}
		us, ok := stmt[0].(*ast.SetOprStmt)
		if ok {
			var i int
			for _, s := range us.SelectList.Selects {
				i = checkOrderBy(t, s, tbl.hasOrderBy, i)
			}
			if !reflect.DeepEqual(tbl.hasOrderBy[i], us.OrderBy != nil) {
				t.Fatalf("got %v, want %v", us.OrderBy != nil, tbl.hasOrderBy[i])
			}
		}
		ss, ok := stmt[0].(*ast.SelectStmt)
		if ok {
			if !reflect.DeepEqual(tbl.hasOrderBy[0], ss.OrderBy != nil) {
				t.Fatalf("got %v, want %v", ss.OrderBy != nil, tbl.hasOrderBy[0])
			}
		}
	}
}

func TestPriority(t *testing.T) {
	p := parser.New()
	stmt, _, err := p.Parse("select HIGH_PRIORITY * from t", "", "")
	if err != nil {
		t.Fatal(err)
	}
	sel := stmt[0].(*ast.SelectStmt)
	if !reflect.DeepEqual(mysql.HighPriority, sel.SelectStmtOpts.Priority) {
		t.Fatalf("got %v, want %v", sel.SelectStmtOpts.Priority, mysql.HighPriority)
	}
}

func TestSQLNoCache(t *testing.T) {
	table := []testCase{
		{`select SQL_NO_CACHE * from t`, false, ""},
		{`select SQL_CACHE * from t`, true, "SELECT * FROM `t`"},
		{`select * from t`, true, "SELECT * FROM `t`"},
	}

	p := parser.New()
	for _, tbl := range table {
		stmt, _, err := p.Parse(tbl.src, "", "")
		if err != nil {
			t.Fatal(err)
		}

		sel := stmt[0].(*ast.SelectStmt)
		if !reflect.DeepEqual(tbl.ok, sel.SelectStmtOpts.SQLCache) {
			t.Fatalf("got %v, want %v", sel.SelectStmtOpts.SQLCache, tbl.ok)
		}
	}
}

func TestBinding(t *testing.T) {
	p := parser.New()
	sms, _, err := p.Parse("create global binding for select * from t using select * from t use index(a)", "", "")
	if err != nil {
		t.Fatal(err)
	}
	v, ok := sms[0].(*ast.CreateBindingStmt)
	if !(ok) {
		t.Fatal("expected true")
	}
	if !reflect.DeepEqual("select * from t", v.OriginNode.Text()) {
		t.Fatalf("got %v, want %v", v.OriginNode.Text(), "select * from t")
	}
	if !reflect.DeepEqual("select * from t use index(a)", v.HintedNode.Text()) {
		t.Fatalf("got %v, want %v", v.HintedNode.Text(), "select * from t use index(a)")
	}
	if !(v.GlobalScope) {
		t.Fatal("expected true")
	}
}

func TestView(t *testing.T) {
	// Test case for the text of the select statement in create view statement.
	p := parser.New()
	sms, _, err := p.Parse("create view v as select * from t", "", "")
	if err != nil {
		t.Fatal(err)
	}
	v, ok := sms[0].(*ast.CreateViewStmt)
	if !(ok) {
		t.Fatal("expected true")
	}
	if !reflect.DeepEqual(ast.AlgorithmUndefined, v.Algorithm) {
		t.Fatalf("got %v, want %v", v.Algorithm, ast.AlgorithmUndefined)
	}
	if !reflect.DeepEqual("select * from t", v.Select.Text()) {
		t.Fatalf("got %v, want %v", v.Select.Text(), "select * from t")
	}
	if !reflect.DeepEqual(ast.SecurityDefiner, v.Security) {
		t.Fatalf("got %v, want %v", v.Security, ast.SecurityDefiner)
	}
	if !reflect.DeepEqual(ast.CheckOptionCascaded, v.CheckOption) {
		t.Fatalf("got %v, want %v", v.CheckOption, ast.CheckOptionCascaded)
	}

	src := `CREATE OR REPLACE ALGORITHM = UNDEFINED DEFINER = root@localhost
                  SQL SECURITY DEFINER
			      VIEW V(a,b,c) AS select c,d,e from t
                  WITH CASCADED CHECK OPTION;`

	var st ast.StmtNode
	st, err = p.ParseOneStmt(src, "", "")
	if err != nil {
		t.Fatal(err)
	}
	v, ok = st.(*ast.CreateViewStmt)
	if !(ok) {
		t.Fatal("expected true")
	}
	if !(v.OrReplace) {
		t.Fatal("expected true")
	}
	if !reflect.DeepEqual(ast.AlgorithmUndefined, v.Algorithm) {
		t.Fatalf("got %v, want %v", v.Algorithm, ast.AlgorithmUndefined)
	}
	if !reflect.DeepEqual("root", v.Definer.Username) {
		t.Fatalf("got %v, want %v", v.Definer.Username, "root")
	}
	if !reflect.DeepEqual("localhost", v.Definer.Hostname) {
		t.Fatalf("got %v, want %v", v.Definer.Hostname, "localhost")
	}
	if !reflect.DeepEqual(ast.NewCIStr("a"), v.Cols[0]) {
		t.Fatalf("got %v, want %v", v.Cols[0], ast.NewCIStr("a"))
	}
	if !reflect.DeepEqual(ast.NewCIStr("b"), v.Cols[1]) {
		t.Fatalf("got %v, want %v", v.Cols[1], ast.NewCIStr("b"))
	}
	if !reflect.DeepEqual(ast.NewCIStr("c"), v.Cols[2]) {
		t.Fatalf("got %v, want %v", v.Cols[2], ast.NewCIStr("c"))
	}
	if !reflect.DeepEqual("select c,d,e from t", v.Select.Text()) {
		t.Fatalf("got %v, want %v", v.Select.Text(), "select c,d,e from t")
	}
	if !reflect.DeepEqual(ast.SecurityDefiner, v.Security) {
		t.Fatalf("got %v, want %v", v.Security, ast.SecurityDefiner)
	}
	if !reflect.DeepEqual(ast.CheckOptionCascaded, v.CheckOption) {
		t.Fatalf("got %v, want %v", v.CheckOption, ast.CheckOptionCascaded)
	}

	src = `
CREATE VIEW v1 AS SELECT * FROM t;
CREATE VIEW v2 AS SELECT 123123123123123;
`
	nodes, _, err := p.Parse(src, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if got := len(nodes); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	if !reflect.DeepEqual(nodes[0].(*ast.CreateViewStmt).Select.Text(), "SELECT * FROM t") {
		t.Fatalf("got %v, want %v", "SELECT * FROM t", nodes[0].(*ast.CreateViewStmt).Select.Text())
	}
	if !reflect.DeepEqual(nodes[1].(*ast.CreateViewStmt).Select.Text(), "SELECT 123123123123123") {
		t.Fatalf("got %v, want %v", "SELECT 123123123123123", nodes[1].(*ast.CreateViewStmt).Select.Text())
	}
}

func TestTimestampDiffUnit(t *testing.T) {
	// Test case for timestampdiff unit.
	// TimeUnit should be unified to upper case.
	p := parser.New()
	stmt, _, err := p.Parse("SELECT TIMESTAMPDIFF(MONTH,'2003-02-01','2003-05-01'), TIMESTAMPDIFF(month,'2003-02-01','2003-05-01');", "", "")
	if err != nil {
		t.Fatal(err)
	}
	ss := stmt[0].(*ast.SelectStmt)
	fields := ss.Fields.Fields
	if got := len(fields); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}
	expr := fields[0].Expr
	f, ok := expr.(*ast.FuncCallExpr)
	if !(ok) {
		t.Fatal("expected true")
	}
	if !reflect.DeepEqual(ast.TimeUnitMonth, f.Args[0].(*ast.TimeUnitExpr).Unit) {
		t.Fatalf("got %v, want %v", f.Args[0].(*ast.TimeUnitExpr).Unit, ast.TimeUnitMonth)
	}

	expr = fields[1].Expr
	f, ok = expr.(*ast.FuncCallExpr)
	if !(ok) {
		t.Fatal("expected true")
	}
	if !reflect.DeepEqual(ast.TimeUnitMonth, f.Args[0].(*ast.TimeUnitExpr).Unit) {
		t.Fatalf("got %v, want %v", f.Args[0].(*ast.TimeUnitExpr).Unit, ast.TimeUnitMonth)
	}
}

func TestFuncCallExprOffset(t *testing.T) {
	// Test case for offset field on func call expr.
	p := parser.New()
	stmt, _, err := p.Parse("SELECT s.a(), b();", "", "")
	if err != nil {
		t.Fatal(err)
	}
	ss := stmt[0].(*ast.SelectStmt)
	fields := ss.Fields.Fields
	if got := len(fields); got != 2 {
		t.Fatalf("expected length %d, got %d", 2, got)
	}

	{
		// s.a()
		expr := fields[0].Expr
		f, ok := expr.(*ast.FuncCallExpr)
		if !(ok) {
			t.Fatal("expected true")
		}
		if !reflect.DeepEqual(7, f.OriginTextPosition()) {
			t.Fatalf("got %v, want %v", f.OriginTextPosition(), 7)
		}
	}

	{
		// b()
		expr := fields[1].Expr
		f, ok := expr.(*ast.FuncCallExpr)
		if !(ok) {
			t.Fatal("expected true")
		}
		if !reflect.DeepEqual(14, f.OriginTextPosition()) {
			t.Fatalf("got %v, want %v", f.OriginTextPosition(), 14)
		}
	}
}

func TestSQLModeANSIQuotes(t *testing.T) {
	p := parser.New()
	p.SetSQLMode(mysql.ModeANSIQuotes)
	tests := []string{
		`CREATE TABLE "table" ("id" int)`,
		`select * from t "tt"`,
	}
	for _, test := range tests {
		_, _, err := p.Parse(test, "", "")
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestDDLStatements(t *testing.T) {
	p := parser.New()
	// Tests that whatever the charset it is define, we always assign utf8 charset and utf8_bin collate.
	createTableStr := `CREATE TABLE t (
		a varchar(64) binary,
		b char(10) charset utf8 collate utf8_general_ci,
		c text charset latin1) ENGINE=innoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin`
	stmts, _, err := p.Parse(createTableStr, "", "")
	if err != nil {
		t.Fatal(err)
	}
	stmt := stmts[0].(*ast.CreateTableStmt)
	if !(mysql.HasBinaryFlag(stmt.Cols[0].Tp.GetFlag())) {
		t.Fatal("expected true")
	}
	for _, colDef := range stmt.Cols[1:] {
		if mysql.HasBinaryFlag(colDef.Tp.GetFlag()) {
			t.Fatal("expected false")
		}
	}
	for _, tblOpt := range stmt.Options {
		switch tblOpt.Tp {
		case ast.TableOptionCharset:
			if !reflect.DeepEqual("utf8", tblOpt.StrValue) {
				t.Fatalf("got %v, want %v", tblOpt.StrValue, "utf8")
			}
		case ast.TableOptionCollate:
			if !reflect.DeepEqual("utf8_bin", tblOpt.StrValue) {
				t.Fatalf("got %v, want %v", tblOpt.StrValue, "utf8_bin")
			}
		}
	}
	createTableStr = `CREATE TABLE t (
		a varbinary(64),
		b binary(10),
		c blob)`
	stmts, _, err = p.Parse(createTableStr, "", "")
	if err != nil {
		t.Fatal(err)
	}
	stmt = stmts[0].(*ast.CreateTableStmt)
	for _, colDef := range stmt.Cols {
		if !reflect.DeepEqual(charset.CharsetBin, colDef.Tp.GetCharset()) {
			t.Fatalf("got %v, want %v", colDef.Tp.GetCharset(), charset.CharsetBin)
		}
		if !reflect.DeepEqual(charset.CollationBin, colDef.Tp.GetCollate()) {
			t.Fatalf("got %v, want %v", colDef.Tp.GetCollate(), charset.CollationBin)
		}
		if !(mysql.HasBinaryFlag(colDef.Tp.GetFlag())) {
			t.Fatal("expected true")
		}
	}
	// Test set collate for all column types
	createTableStr = `CREATE TABLE t (
		c_int int collate utf8_bin,
		c_real real collate utf8_bin,
		c_float float collate utf8_bin,
		c_bool bool collate utf8_bin,
		c_char char collate utf8_bin,
		c_binary binary collate utf8_bin,
		c_varchar varchar(2) collate utf8_bin,
		c_year year collate utf8_bin,
		c_date date collate utf8_bin,
		c_time time collate utf8_bin,
		c_datetime datetime collate utf8_bin,
		c_timestamp timestamp collate utf8_bin,
		c_tinyblob tinyblob collate utf8_bin,
		c_blob blob collate utf8_bin,
		c_mediumblob mediumblob collate utf8_bin,
		c_longblob longblob collate utf8_bin,
		c_bit bit collate utf8_bin,
		c_long_varchar long varchar collate utf8_bin,
		c_tinytext tinytext collate utf8_bin,
		c_text text collate utf8_bin,
		c_mediumtext mediumtext collate utf8_bin,
		c_longtext longtext collate utf8_bin,
		c_decimal decimal collate utf8_bin,
		c_numeric numeric collate utf8_bin,
		c_enum enum('1') collate utf8_bin,
		c_set set('1') collate utf8_bin,
		c_json json collate utf8_bin)`
	_, _, err = p.Parse(createTableStr, "", "")
	if err != nil {
		t.Fatal(err)
	}

	createTableStr = `CREATE TABLE t (c_double double(10))`
	_, _, err = p.Parse(createTableStr, "", "")
	if err == nil || err.Error() != "[parser:1149]You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use" {
		t.Fatalf("expected error %q, got %v", "[parser:1149]You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use", err)
	}
	p.SetStrictDoubleTypeCheck(false)
	_, _, err = p.Parse(createTableStr, "", "")
	if err != nil {
		t.Fatal(err)
	}
	p.SetStrictDoubleTypeCheck(true)

	createTableStr = `CREATE TABLE t (c_double double(10, 2))`
	_, _, err = p.Parse(createTableStr, "", "")
	if err != nil {
		t.Fatal(err)
	}

	createTableStr = `create global temporary table t010(local_01 int, local_03 varchar(20))`
	_, _, err = p.Parse(createTableStr, "", "")
	if err == nil || err.Error() != "line 1 column 70 near \"\"GLOBAL TEMPORARY and ON COMMIT DELETE ROWS must appear together " {
		t.Fatalf("expected error %q, got %v", "line 1 column 70 near \"\"GLOBAL TEMPORARY and ON COMMIT DELETE ROWS must appear together ", err)
	}

	createTableStr = `create global temporary table t010(local_01 int, local_03 varchar(20)) on commit preserve rows`
	_, _, err = p.Parse(createTableStr, "", "")
	if err != nil {
		t.Fatal(err)
	}
}

func TestChangeReplicationSource(t *testing.T) {
	// SOURCE_PASSWORD is masked in the sensitive-statement text.
	p := parser.New()
	stmt, err := p.ParseOneStmt("change replication source to source_user = 'repl', source_password = 'hunter2'", "", "")
	if err != nil {
		t.Fatal(err)
	}
	sensitive, ok := stmt.(ast.SensitiveStmtNode)
	if !ok {
		t.Fatalf("expected ChangeReplicationSourceStmt to implement SensitiveStmtNode, got %T", stmt)
	}
	secure := sensitive.SecureText()
	if strings.Contains(secure, "hunter2") {
		t.Fatalf("SecureText leaked the password: %q", secure)
	}
	want := "CHANGE REPLICATION SOURCE TO SOURCE_USER = 'repl', SOURCE_PASSWORD = 'xxxxxx'"
	if secure != want {
		t.Fatalf("got %q, want %q", secure, want)
	}
}

func TestTableSample(t *testing.T) {
	p := parser.New()
	cases := []string{
		"select * from tbl tablesample (33.3 + 44.4);",
		"select * from tbl tablesample (33.3 + 44.4 percent);",
		"select * from tbl tablesample (33 + 44 rows);",
		"select * from tbl tablesample (33 + 44 rows) repeatable (55 + 66);",
		"select * from tbl tablesample (200);",
		"select * from tbl tablesample (-10);",
		"select * from tbl tablesample (null);",
		"select * from tbl tablesample (33.3 rows);",
		"select * from tbl tablesample (-4 rows);",
		"select * from tbl tablesample (50) repeatable ('ssss');",
		"delete from tbl using tbl2 tablesample(10 rows) repeatable (111) where tbl.id = tbl2.id",
		"update tbl tablesample regions() set id = '1'",
	}
	for _, sql := range cases {
		_, err := p.ParseOneStmt(sql, "", "")
		if err != nil {
			t.Fatalf("%s: %v", fmt.Sprintf("source %v", sql), err)
		}
	}
}

func TestGeneratedColumn(t *testing.T) {
	tests := []struct {
		input string
		ok    bool
		expr  string
	}{
		{"create table t (c int, d int generated always as (c + 1) virtual)", true, "c + 1"},
		{"create table t (c int, d int as (   c + 1   ) virtual)", true, "c + 1"},
		{"create table t (c int, d int as (1 + 1) stored)", true, "1 + 1"},
	}
	p := parser.New()
	for _, tbl := range tests {
		stmtNodes, _, err := p.Parse(tbl.input, "", "")
		if tbl.ok {
			if err != nil {
				t.Fatal(err)
			}
			stmtNode := stmtNodes[0]
			for _, col := range stmtNode.(*ast.CreateTableStmt).Cols {
				for _, opt := range col.Options {
					if opt.Tp == ast.ColumnOptionGenerated {
						if !reflect.DeepEqual(tbl.expr, opt.Expr.Text()) {
							t.Fatalf("got %v, want %v", opt.Expr.Text(), tbl.expr)
						}
					}
				}
			}
		} else {
			if err == nil {
				t.Fatal("expected error")
			}
		}
	}

	_, _, err := p.Parse("create table t1 (a int, b int as (a + 1) default 10);", "", "")
	if !reflect.DeepEqual(err.Error(), "[ddl:1221]Incorrect usage of DEFAULT and generated column") {
		t.Fatalf("got %v, want %v", "[ddl:1221]Incorrect usage of DEFAULT and generated column", err.Error())
	}
	_, _, err = p.Parse("create table t1 (a int, b int as (a + 1) on update now());", "", "")
	if !reflect.DeepEqual(err.Error(), "[ddl:1221]Incorrect usage of ON UPDATE and generated column") {
		t.Fatalf("got %v, want %v", "[ddl:1221]Incorrect usage of ON UPDATE and generated column", err.Error())
	}
	_, _, err = p.Parse("create table t1 (a int, b int as (a + 1) auto_increment);", "", "")
	if !reflect.DeepEqual(err.Error(), "[ddl:1221]Incorrect usage of AUTO_INCREMENT and generated column") {
		t.Fatalf("got %v, want %v", "[ddl:1221]Incorrect usage of AUTO_INCREMENT and generated column", err.Error())
	}
}

func TestSetTransaction(t *testing.T) {
	// Set transaction is equivalent to setting the global or session value of tx_isolation.
	// For example:
	// SET SESSION TRANSACTION ISOLATION LEVEL READ COMMITTED
	// SET SESSION tx_isolation='READ-COMMITTED'
	tests := []struct {
		input    string
		isGlobal bool
		value    string
	}{
		{
			"SET SESSION TRANSACTION ISOLATION LEVEL READ COMMITTED",
			false, "READ-COMMITTED",
		},
		{
			"SET GLOBAL TRANSACTION ISOLATION LEVEL REPEATABLE READ",
			true, "REPEATABLE-READ",
		},
	}
	p := parser.New()
	for _, tbl := range tests {
		stmt1, err := p.ParseOneStmt(tbl.input, "", "")
		if err != nil {
			t.Fatal(err)
		}
		setStmt := stmt1.(*ast.SetStmt)
		vars := setStmt.Variables[0]
		if !reflect.DeepEqual("tx_isolation", vars.Name) {
			t.Fatalf("got %v, want %v", vars.Name, "tx_isolation")
		}
		if !reflect.DeepEqual(tbl.isGlobal, vars.IsGlobal) {
			t.Fatalf("got %v, want %v", vars.IsGlobal, tbl.isGlobal)
		}
		if !reflect.DeepEqual(true, vars.IsSystem) {
			t.Fatalf("got %v, want %v", vars.IsSystem, true)
		}
		if !reflect.DeepEqual(tbl.value, vars.Value.(ast.ValueExpr).GetValue()) {
			t.Fatalf("got %v, want %v", vars.Value.(ast.ValueExpr).GetValue(), tbl.value)
		}
	}
}

func TestSideEffect(t *testing.T) {
	// This test cover a bug that parse an error SQL doesn't leave the parser in a
	// clean state, cause the following SQL parse fail.
	p := parser.New()
	_, err := p.ParseOneStmt("create table t /*!50100 'abc', 'abc' */;", "", "")
	if err == nil {
		t.Fatal("expected error")
	}

	_, err = p.ParseOneStmt("show tables;", "", "")
	if err != nil {
		t.Fatal(err)
	}
}

func TestTablePartition(t *testing.T) {
	// Check comment content.
	p := parser.New()
	stmt, err := p.ParseOneStmt("create table t (id int) partition by range (id) (partition p0 values less than (10) comment 'check')", "", "")
	if err != nil {
		t.Fatal(err)
	}
	createTable := stmt.(*ast.CreateTableStmt)
	comment, ok := createTable.Partition.Definitions[0].Comment()
	if !(ok) {
		t.Fatal("expected true")
	}
	if !reflect.DeepEqual("check", comment) {
		t.Fatalf("got %v, want %v", comment, "check")
	}
}

func TestTablePartitionNameList(t *testing.T) {
	table := []testCase{
		{`select * from t partition (p0,p1)`, true, ""},
	}

	p := parser.New()
	for _, tbl := range table {
		stmt, _, err := p.Parse(tbl.src, "", "")
		if err != nil {
			t.Fatal(err)
		}

		sel := stmt[0].(*ast.SelectStmt)
		source, ok := sel.From.TableRefs.Left.(*ast.TableSource)
		if !(ok) {
			t.Fatal("expected true")
		}
		tableName, ok := source.Source.(*ast.TableName)
		if !(ok) {
			t.Fatal("expected true")
		}
		if got := len(tableName.PartitionNames); got != 2 {
			t.Fatalf("expected length %d, got %d", 2, got)
		}
		if !reflect.DeepEqual(ast.CIStr{O: "p0", L: "p0"}, tableName.PartitionNames[0]) {
			t.Fatalf("got %v, want %v", tableName.PartitionNames[0], ast.CIStr{O: "p0", L: "p0"})
		}
		if !reflect.DeepEqual(ast.CIStr{O: "p1", L: "p1"}, tableName.PartitionNames[1]) {
			t.Fatalf("got %v, want %v", tableName.PartitionNames[1], ast.CIStr{O: "p1", L: "p1"})
		}
	}
}

func TestNotExistsSubquery(t *testing.T) {
	table := []testCase{
		{`select * from t1 where not exists (select * from t2 where t1.a = t2.a)`, true, ""},
	}

	p := parser.New()
	for _, tbl := range table {
		stmt, _, err := p.Parse(tbl.src, "", "")
		if err != nil {
			t.Fatal(err)
		}

		sel := stmt[0].(*ast.SelectStmt)
		exists, ok := sel.Where.(*ast.ExistsSubqueryExpr)
		if !(ok) {
			t.Fatal("expected true")
		}
		if !reflect.DeepEqual(tbl.ok, exists.Not) {
			t.Fatalf("got %v, want %v", exists.Not, tbl.ok)
		}
	}
}

func TestWindowFunctionIdentifier(t *testing.T) {
	//nolint: prealloc
	var table []testCase
	for key := range parser.WindowFuncTokenMapForTest {
		table = append(table, testCase{fmt.Sprintf("select 1 %s", key), false, fmt.Sprintf("SELECT 1 AS `%s`", key)})
	}
	RunTest(t, table, true, false)

	for i := range table {
		table[i].ok = true
	}
	RunTest(t, table, false, false)
}

type windowFrameBoundChecker struct {
	fb     *ast.FrameBound
	exprRc int
	unit   ast.TimeUnitType
	t      *testing.T
}

// Enter implements ast.Visitor interface.
func (wfc *windowFrameBoundChecker) Enter(inNode ast.Node) (outNode ast.Node, skipChildren bool) {
	if _, ok := inNode.(*ast.FrameBound); ok {
		wfc.fb = inNode.(*ast.FrameBound)
		if wfc.fb.Unit != ast.TimeUnitInvalid {
			_, ok := wfc.fb.Expr.(ast.ValueExpr)
			if ok {
				wfc.t.Fatal("expected false")
			}
		}
	}
	return inNode, false
}

// Leave implements ast.Visitor interface.
func (wfc *windowFrameBoundChecker) Leave(inNode ast.Node) (node ast.Node, ok bool) {
	if _, ok := inNode.(*ast.FrameBound); ok {
		wfc.fb = nil
	}
	if wfc.fb != nil {
		if inNode == wfc.fb.Expr {
			wfc.exprRc++
		}
		wfc.unit = wfc.fb.Unit
	}
	return inNode, true
}

// For issue #51
// See https://github.com/pingcap/parser/pull/51 for details
func TestVisitFrameBound(t *testing.T) {
	p := parser.New()
	p.EnableWindowFunc(true)
	table := []struct {
		s      string
		exprRc int
		unit   ast.TimeUnitType
	}{
		{`SELECT AVG(val) OVER (RANGE INTERVAL 1+3 MINUTE_SECOND PRECEDING) FROM t;`, 1, ast.TimeUnitMinuteSecond},
		{`SELECT AVG(val) OVER (RANGE 5 PRECEDING) FROM t;`, 1, ast.TimeUnitInvalid},
		{`SELECT AVG(val) OVER () FROM t;`, 0, ast.TimeUnitInvalid},
	}
	for _, tbl := range table {
		stmt, err := p.ParseOneStmt(tbl.s, "", "")
		if err != nil {
			t.Fatal(err)
		}
		checker := windowFrameBoundChecker{t: t}
		stmt.Accept(&checker)
		if !reflect.DeepEqual(tbl.exprRc, checker.exprRc) {
			t.Fatalf("got %v, want %v", checker.exprRc, tbl.exprRc)
		}
		if !reflect.DeepEqual(tbl.unit, checker.unit) {
			t.Fatalf("got %v, want %v", checker.unit, tbl.unit)
		}
	}
}

func TestFieldText(t *testing.T) {
	p := parser.New()
	stmts, _, err := p.Parse("select a from t", "", "")
	if err != nil {
		t.Fatal(err)
	}
	tmp := stmts[0].(*ast.SelectStmt)
	if !reflect.DeepEqual("a", tmp.Fields.Fields[0].Text()) {
		t.Fatalf("got %v, want %v", tmp.Fields.Fields[0].Text(), "a")
	}

	sqls := []string{
		"trace select a from t",
		"trace format = 'row' select a from t",
		"trace format = 'json' select a from t",
	}
	for _, sql := range sqls {
		stmts, _, err = p.Parse(sql, "", "")
		if err != nil {
			t.Fatal(err)
		}
		traceStmt := stmts[0].(*ast.TraceStmt)
		if !reflect.DeepEqual(sql, traceStmt.Text()) {
			t.Fatalf("got %v, want %v", traceStmt.Text(), sql)
		}
		if !reflect.DeepEqual("select a from t", traceStmt.Stmt.Text()) {
			t.Fatalf("got %v, want %v", traceStmt.Stmt.Text(), "select a from t")
		}
	}
}

// See https://github.com/pingcap/parser/issue/94
func TestQuotedSystemVariables(t *testing.T) {
	p := parser.New()

	st, err := p.ParseOneStmt(
		"select @@Sql_Mode, @@`SQL_MODE`, @@session.`sql_mode`, @@global.`s ql``mode`, @@session.'sql\\nmode', @@local.\"sql\\\"mode\", @@instance.sql_mode;",
		"",
		"",
	)
	if err != nil {
		t.Fatal(err)
	}
	ss := st.(*ast.SelectStmt)
	expected := []*ast.VariableExpr{
		{
			Name:          "sql_mode",
			IsGlobal:      false,
			IsSystem:      true,
			ExplicitScope: false,
		},
		{
			Name:          "sql_mode",
			IsGlobal:      false,
			IsSystem:      true,
			ExplicitScope: false,
		},
		{
			Name:          "sql_mode",
			IsGlobal:      false,
			IsSystem:      true,
			ExplicitScope: true,
		},
		{
			Name:          "s ql`mode",
			IsGlobal:      true,
			IsSystem:      true,
			ExplicitScope: true,
		},
		{
			Name:          "sql\nmode",
			IsGlobal:      false,
			IsSystem:      true,
			ExplicitScope: true,
		},
		{
			Name:          `sql"mode`,
			IsGlobal:      false,
			IsSystem:      true,
			ExplicitScope: true,
		},
		{
			Name:          "sql_mode",
			IsGlobal:      false,
			IsSystem:      true,
			IsInstance:    true,
			ExplicitScope: true,
		},
	}

	if got := len(ss.Fields.Fields); got != len(expected) {
		t.Fatalf("expected length %d, got %d", len(expected), got)
	}
	for i, field := range ss.Fields.Fields {
		ve := field.Expr.(*ast.VariableExpr)
		comment := fmt.Sprintf("field %d, ve = %v", i, ve)
		if !reflect.DeepEqual(expected[i].Name, ve.Name) {
			t.Fatalf("%v: got %v, want %v", comment, ve.Name, expected[i].Name)
		}
		if !reflect.DeepEqual(expected[i].IsGlobal, ve.IsGlobal) {
			t.Fatalf("%v: got %v, want %v", comment, ve.IsGlobal, expected[i].IsGlobal)
		}
		if !reflect.DeepEqual(expected[i].IsInstance, ve.IsInstance) {
			t.Fatalf("%v: got %v, want %v", comment, ve.IsInstance, expected[i].IsInstance)
		}
		if !reflect.DeepEqual(expected[i].IsSystem, ve.IsSystem) {
			t.Fatalf("%v: got %v, want %v", comment, ve.IsSystem, expected[i].IsSystem)
		}
		if !reflect.DeepEqual(expected[i].ExplicitScope, ve.ExplicitScope) {
			t.Fatalf("%v: got %v, want %v", comment, ve.ExplicitScope, expected[i].ExplicitScope)
		}
	}
}

// See https://github.com/pingcap/parser/issue/95
func TestQuotedVariableColumnName(t *testing.T) {
	p := parser.New()

	st, err := p.ParseOneStmt(
		"select @abc, @`abc`, @'aBc', @\"AbC\", @6, @`6`, @'6', @\"6\", @@sql_mode, @@`sql_mode`, @;",
		"",
		"",
	)
	if err != nil {
		t.Fatal(err)
	}
	ss := st.(*ast.SelectStmt)
	expected := []string{
		"@abc",
		"@`abc`",
		"@'aBc'",
		`@"AbC"`,
		"@6",
		"@`6`",
		"@'6'",
		`@"6"`,
		"@@sql_mode",
		"@@`sql_mode`",
		"@",
	}

	if got := len(ss.Fields.Fields); got != len(expected) {
		t.Fatalf("expected length %d, got %d", len(expected), got)
	}
	for i, field := range ss.Fields.Fields {
		if !reflect.DeepEqual(expected[i], field.Text()) {
			t.Fatalf("got %v, want %v", field.Text(), expected[i])
		}
	}
}

func TestCharset(t *testing.T) {
	p := parser.New()

	st, err := p.ParseOneStmt("ALTER SCHEMA GLOBAL DEFAULT CHAR SET utf8mb4", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if st.(*ast.AlterDatabaseStmt) == nil {
		t.Fatal("expected non-nil")
	}
	st, err = p.ParseOneStmt("ALTER DATABASE CHAR SET = utf8mb4", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if st.(*ast.AlterDatabaseStmt) == nil {
		t.Fatal("expected non-nil")
	}
	st, err = p.ParseOneStmt("ALTER DATABASE DEFAULT CHAR SET = utf8mb4", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if st.(*ast.AlterDatabaseStmt) == nil {
		t.Fatal("expected non-nil")
	}
}

func TestUnderscoreCharset(t *testing.T) {
	p := parser.New()
	tests := []struct {
		cs        string
		parseFail bool
		unSupport bool
	}{
		{"utf8", false, false},
		{"gbk", false, true},
		{"ujis", false, true},
		{"gbk1", true, true},
		{"ujisx", true, true},
	}
	for _, tt := range tests {
		sql := fmt.Sprintf("select hex(_%s '3F')", tt.cs)
		_, err := p.ParseOneStmt(sql, "", "")
		if tt.parseFail {
			if err == nil || err.Error() != fmt.Sprintf("line 1 column %d near \"'3F')\" ", len(tt.cs)+17) {
				t.Fatalf("expected error %q, got %v", fmt.Sprintf("line 1 column %d near \"'3F')\" ", len(tt.cs)+17), err)
			}
		} else if tt.unSupport {
			if err == nil || err.Error() != ast.ErrUnknownCharacterSet.GenWithStack("Unsupported character introducer: '%-.64s'", tt.cs).Error() {
				t.Fatalf("expected error %q, got %v", ast.ErrUnknownCharacterSet.GenWithStack("Unsupported character introducer: '%-.64s'", tt.cs).Error(), err)
			}
		} else {
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestFulltextSearch(t *testing.T) {
	p := parser.New()

	st, err := p.ParseOneStmt("SELECT * FROM fulltext_test WHERE MATCH(content) AGAINST('search')", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if st.(*ast.SelectStmt) == nil {
		t.Fatal("expected non-nil")
	}

	st, err = p.ParseOneStmt("SELECT * FROM fulltext_test WHERE MATCH() AGAINST('search')", "", "")
	if err == nil {
		t.Fatal("expected error")
	}
	if st != nil {
		t.Fatalf("expected nil, got %v", st)
	}

	st, err = p.ParseOneStmt("SELECT * FROM fulltext_test WHERE MATCH(content) AGAINST()", "", "")
	if err == nil {
		t.Fatal("expected error")
	}
	if st != nil {
		t.Fatalf("expected nil, got %v", st)
	}

	st, err = p.ParseOneStmt("SELECT * FROM fulltext_test WHERE MATCH(content) AGAINST('search' IN)", "", "")
	if err == nil {
		t.Fatal("expected error")
	}
	if st != nil {
		t.Fatalf("expected nil, got %v", st)
	}

	st, err = p.ParseOneStmt("SELECT * FROM fulltext_test WHERE MATCH(content) AGAINST('search' IN BOOLEAN MODE WITH QUERY EXPANSION)", "", "")
	if err == nil {
		t.Fatal("expected error")
	}
	if st != nil {
		t.Fatalf("expected nil, got %v", st)
	}

	st, err = p.ParseOneStmt("SELECT * FROM fulltext_test WHERE MATCH(title,content) AGAINST('search' IN NATURAL LANGUAGE MODE)", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if st.(*ast.SelectStmt) == nil {
		t.Fatal("expected non-nil")
	}
	writer := bytes.NewBufferString("")
	st.(*ast.SelectStmt).Where.Format(writer)
	if !reflect.DeepEqual("MATCH(title,content) AGAINST(\"search\")", writer.String()) {
		t.Fatalf("got %v, want %v", writer.String(), "MATCH(title,content) AGAINST(\"search\")")
	}

	st, err = p.ParseOneStmt("SELECT * FROM fulltext_test WHERE MATCH(title,content) AGAINST('search' IN BOOLEAN MODE)", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if st.(*ast.SelectStmt) == nil {
		t.Fatal("expected non-nil")
	}
	writer.Reset()
	st.(*ast.SelectStmt).Where.Format(writer)
	if !reflect.DeepEqual("MATCH(title,content) AGAINST(\"search\" IN BOOLEAN MODE)", writer.String()) {
		t.Fatalf("got %v, want %v", writer.String(), "MATCH(title,content) AGAINST(\"search\" IN BOOLEAN MODE)")
	}

	st, err = p.ParseOneStmt("SELECT * FROM fulltext_test WHERE MATCH(title,content) AGAINST('search' WITH QUERY EXPANSION)", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if st.(*ast.SelectStmt) == nil {
		t.Fatal("expected non-nil")
	}
	writer.Reset()
	st.(*ast.SelectStmt).Where.Format(writer)
	if !reflect.DeepEqual("MATCH(title,content) AGAINST(\"search\" WITH QUERY EXPANSION)", writer.String()) {
		t.Fatalf("got %v, want %v", writer.String(), "MATCH(title,content) AGAINST(\"search\" WITH QUERY EXPANSION)")
	}
}

func TestSignedInt64OutOfRange(t *testing.T) {
	p := parser.New()
	cases := []string{
		"recover table by job 18446744073709551612",
		"recover table t 18446744073709551612",
		"admin check index t idx (0, 18446744073709551612)",
		"create user abc@def with max_queries_per_hour 18446744073709551612",
	}

	for _, s := range cases {
		_, err := p.ParseOneStmt(s, "", "")
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "out of range") {
			t.Fatalf("expected %q to contain %q", err.Error(), "out of range")
		}
	}
}

// CleanNodeText set the text of node and all child node empty.
// For test only.
func CleanNodeText(node ast.Node) {
	var cleaner nodeTextCleaner
	node.Accept(&cleaner)
}

// nodeTextCleaner clean the text of a node and it's child node.
// For test only.
type nodeTextCleaner struct {
}

func cleanPartition(n ast.Node) {
	if p, ok := n.(*ast.PartitionOptions); ok && p != nil {
		var tmpCleaner nodeTextCleaner
		if p.Interval != nil {
			p.Interval.SetText(nil, "")
			p.Interval.SetOriginTextPosition(0)
			p.Interval.IntervalExpr.Expr.Accept(&tmpCleaner)
			if p.Interval.FirstRangeEnd != nil {
				(*p.Interval.FirstRangeEnd).Accept(&tmpCleaner)
			}
			if p.Interval.LastRangeEnd != nil {
				(*p.Interval.LastRangeEnd).Accept(&tmpCleaner)
			}
		}
	}
}

// Enter implements Visitor interface.
func (checker *nodeTextCleaner) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	in.SetText(nil, "")
	in.SetOriginTextPosition(0)
	if v, ok := in.(ast.ValueExpr); ok && v != nil {
		tpFlag := v.GetType().GetFlag()
		if tpFlag&mysql.UnderScoreCharsetFlag != 0 {
			// ignore underscore charset flag to let `'abc' = _utf8'abc'` pass
			tpFlag ^= mysql.UnderScoreCharsetFlag
			v.GetType().SetFlag(tpFlag)
		}
	}

	switch node := in.(type) {
	case *ast.PatternLikeOrIlikeExpr:
		if node.Escape == '\\' {
			node.EscapeExplicit = false
		}
	case *ast.CreateTableStmt:
		for _, opt := range node.Options {
			switch opt.Tp {
			case ast.TableOptionCharset:
				opt.StrValue = strings.ToUpper(opt.StrValue)
			case ast.TableOptionCollate:
				opt.StrValue = strings.ToUpper(opt.StrValue)
			}
		}
		for _, col := range node.Cols {
			col.Tp.SetCharset(strings.ToUpper(col.Tp.GetCharset()))
			col.Tp.SetCollate(strings.ToUpper(col.Tp.GetCollate()))

			for i, option := range col.Options {
				if option.Tp == 0 && option.Expr == nil && !option.Stored && option.Refer == nil {
					col.Options = slices.Delete(col.Options, i, i+1)
				}
			}
		}
	case *ast.DeleteStmt:
		for _, tableHint := range node.TableHints {
			tableHint.HintName.O = ""
		}
	case *ast.UpdateStmt:
		for _, tableHint := range node.TableHints {
			tableHint.HintName.O = ""
		}
	case *ast.Constraint:
		if node.Option != nil {
			if node.Option.KeyBlockSize == 0x0 && node.Option.Tp == 0 && node.Option.Comment == "" {
				node.Option = nil
			}
		}
	case *ast.FuncCallExpr:
		node.FnName.O = strings.ToLower(node.FnName.O)
		node.SetOriginTextPosition(0)
	case *ast.AggregateFuncExpr:
		node.F = strings.ToLower(node.F)
	case *ast.SelectField:
		node.Offset = 0
	case *ast.ValueExprBase:
		if node.Kind() == ast.KindMysqlDecimal {
			_ = node.GetMysqlDecimal().FromString(node.GetMysqlDecimal().ToString())
		}
	case *ast.GrantStmt:
		var privs []*ast.PrivElem
		for _, v := range node.Privs {
			if v.Priv != 0 {
				privs = append(privs, v)
			}
		}
		node.Privs = privs
	case *ast.AlterTableStmt:
		var specs []*ast.AlterTableSpec
		for _, v := range node.Specs {
			if v.Tp != 0 && !(v.Tp == ast.AlterTableOption && len(v.Options) == 0) {
				specs = append(specs, v)
			}
		}
		node.Specs = specs
	case *ast.Join:
		node.ExplicitParens = false
	case *ast.ColumnDef:
		node.Tp.CleanElemIsBinaryLit()
	case *ast.PartitionOptions:
		cleanPartition(node)
	case *ast.ProcedureBlock:
		// ProcedureBlock.Accept deliberately does not traverse
		// ProcedureProcStmts; clean them explicitly so restored
		// procedure bodies compare deep-equal.
		var tmpCleaner nodeTextCleaner
		for _, stmt := range node.ProcedureProcStmts {
			stmt.Accept(&tmpCleaner)
		}
	case *ast.ProcedureInfo:
		// The parameter list's recorded source text is not normalized by
		// Restore().
		node.ProcedureParamStr = ""
	case *ast.StoreParameter:
		// Restore() prints the type via CompactStr, which substitutes
		// the default display width for an unspecified one; canonicalize
		// the parsed type the same way so IN a INT compares deep-equal
		// with its restored IN a INT(11) spelling.
		if node.ParamType.GetFlen() == types.UnspecifiedLength {
			flen, _ := mysql.GetDefaultFieldLengthAndDecimal(node.ParamType.GetType())
			node.ParamType.SetFlen(flen)
		}
		if node.ParamType.GetDecimal() == types.UnspecifiedLength {
			_, decimal := mysql.GetDefaultFieldLengthAndDecimal(node.ParamType.GetType())
			node.ParamType.SetDecimal(decimal)
		}
	}
	return in, false
}

// Leave implements Visitor interface.
func (checker *nodeTextCleaner) Leave(in ast.Node) (out ast.Node, ok bool) {
	return in, true
}

func TestStatisticsOps(t *testing.T) {
	p := parser.New()
	sms, _, err := p.Parse("create statistics if not exists stats1 (cardinality) on t(a,b,c)", "", "")
	if err != nil {
		t.Fatal(err)
	}
	v, ok := sms[0].(*ast.CreateStatisticsStmt)
	if !(ok) {
		t.Fatal("expected true")
	}
	if !(v.IfNotExists) {
		t.Fatal("expected true")
	}
	if !reflect.DeepEqual("stats1", v.StatsName) {
		t.Fatalf("got %v, want %v", v.StatsName, "stats1")
	}
	if !reflect.DeepEqual(ast.StatsTypeCardinality, v.StatsType) {
		t.Fatalf("got %v, want %v", v.StatsType, ast.StatsTypeCardinality)
	}
	if !reflect.DeepEqual(ast.CIStr{O: "t", L: "t"}, v.Table.Name) {
		t.Fatalf("got %v, want %v", v.Table.Name, ast.CIStr{O: "t", L: "t"})
	}
	if got := len(v.Columns); got != 3 {
		t.Fatalf("expected length %d, got %d", 3, got)
	}
	if !reflect.DeepEqual(ast.CIStr{O: "a", L: "a"}, v.Columns[0].Name) {
		t.Fatalf("got %v, want %v", v.Columns[0].Name, ast.CIStr{O: "a", L: "a"})
	}
	if !reflect.DeepEqual(ast.CIStr{O: "b", L: "b"}, v.Columns[1].Name) {
		t.Fatalf("got %v, want %v", v.Columns[1].Name, ast.CIStr{O: "b", L: "b"})
	}
	if !reflect.DeepEqual(ast.CIStr{O: "c", L: "c"}, v.Columns[2].Name) {
		t.Fatalf("got %v, want %v", v.Columns[2].Name, ast.CIStr{O: "c", L: "c"})
	}
}

func TestHighNotPrecedenceMode(t *testing.T) {
	p := parser.New()
	var sb strings.Builder

	sms, _, err := p.Parse("SELECT NOT 1 BETWEEN -5 AND 5", "", "")
	if err != nil {
		t.Fatal(err)
	}
	v, ok := sms[0].(*ast.SelectStmt)
	if !(ok) {
		t.Fatal("expected true")
	}
	v1, ok := v.Fields.Fields[0].Expr.(*ast.UnaryOperationExpr)
	if !(ok) {
		t.Fatal("expected true")
	}
	if !reflect.DeepEqual(opcode.Not, v1.Op) {
		t.Fatalf("got %v, want %v", v1.Op, opcode.Not)
	}
	err = sms[0].Restore(NewRestoreCtx(DefaultRestoreFlags, &sb))
	if err != nil {
		t.Fatal(err)
	}
	restoreSQL := sb.String()
	if !reflect.DeepEqual("SELECT NOT 1 BETWEEN -5 AND 5", restoreSQL) {
		t.Fatalf("got %v, want %v", restoreSQL, "SELECT NOT 1 BETWEEN -5 AND 5")
	}
	sb.Reset()

	sms, _, err = p.Parse("SELECT !1 BETWEEN -5 AND 5", "", "")
	if err != nil {
		t.Fatal(err)
	}
	v, ok = sms[0].(*ast.SelectStmt)
	if !(ok) {
		t.Fatal("expected true")
	}
	_, ok = v.Fields.Fields[0].Expr.(*ast.BetweenExpr)
	if !(ok) {
		t.Fatal("expected true")
	}
	err = sms[0].Restore(NewRestoreCtx(DefaultRestoreFlags, &sb))
	if err != nil {
		t.Fatal(err)
	}
	restoreSQL = sb.String()
	if !reflect.DeepEqual("SELECT !1 BETWEEN -5 AND 5", restoreSQL) {
		t.Fatalf("got %v, want %v", restoreSQL, "SELECT !1 BETWEEN -5 AND 5")
	}
	sb.Reset()

	p = parser.New()
	p.SetSQLMode(mysql.ModeHighNotPrecedence)
	sms, _, err = p.Parse("SELECT NOT 1 BETWEEN -5 AND 5", "", "")
	if err != nil {
		t.Fatal(err)
	}
	v, ok = sms[0].(*ast.SelectStmt)
	if !(ok) {
		t.Fatal("expected true")
	}
	_, ok = v.Fields.Fields[0].Expr.(*ast.BetweenExpr)
	if !(ok) {
		t.Fatal("expected true")
	}
	err = sms[0].Restore(NewRestoreCtx(DefaultRestoreFlags, &sb))
	if err != nil {
		t.Fatal(err)
	}
	restoreSQL = sb.String()
	if !reflect.DeepEqual("SELECT !1 BETWEEN -5 AND 5", restoreSQL) {
		t.Fatalf("got %v, want %v", restoreSQL, "SELECT !1 BETWEEN -5 AND 5")
	}
}

func TestWithoutCharsetFlags(t *testing.T) {
	type testCaseWithFlag struct {
		src     string
		ok      bool
		restore string
		flag    RestoreFlags
	}

	flag := RestoreStringSingleQuotes | RestoreSpacesAroundBinaryOperation | RestoreBracketAroundBinaryOperation | RestoreNameBackQuotes
	cases := []testCaseWithFlag{
		{"select 'a'", true, "SELECT 'a'", flag | RestoreStringWithoutCharset},
		{"select _utf8'a'", true, "SELECT 'a'", flag | RestoreStringWithoutCharset},
		{"select _utf8mb4'a'", true, "SELECT 'a'", flag | RestoreStringWithoutCharset},
		{"select _utf8 X'D0B1'", true, "SELECT x'd0b1'", flag | RestoreStringWithoutCharset},

		{"select _utf8mb4'a'", true, "SELECT 'a'", flag | RestoreStringWithoutDefaultCharset},
		{"select _utf8'a'", true, "SELECT _utf8'a'", flag | RestoreStringWithoutDefaultCharset},
		{"select _utf8'a'", true, "SELECT _utf8'a'", flag | RestoreStringWithoutDefaultCharset},
		{"select _utf8 X'D0B1'", true, "SELECT _utf8 x'd0b1'", flag | RestoreStringWithoutDefaultCharset},
	}

	p := parser.New()
	p.EnableWindowFunc(false)
	for _, tbl := range cases {
		stmts, _, err := p.Parse(tbl.src, "", "")
		if !tbl.ok {
			if err == nil {
				t.Fatal("expected error")
			}
			continue
		}
		if err != nil {
			t.Fatal(err)
		}
		// restore correctness test
		var sb strings.Builder
		restoreSQLs := ""
		for _, stmt := range stmts {
			sb.Reset()
			ctx := NewRestoreCtx(tbl.flag, &sb)
			ctx.DefaultDB = "test"
			err = stmt.Restore(ctx)
			if err != nil {
				t.Fatal(err)
			}
			restoreSQL := sb.String()
			if restoreSQLs != "" {
				restoreSQLs += "; "
			}
			restoreSQLs += restoreSQL
		}
		if !reflect.DeepEqual(tbl.restore, restoreSQLs) {
			t.Fatalf("got %v, want %v", restoreSQLs, tbl.restore)
		}
	}
}

func TestRestoreBinOpWithBrackets(t *testing.T) {
	cases := []testCase{
		{"select mod(a+b, 4)+1", true, "SELECT (((`a` + `b`) % 4) + 1)"},
		{"SELECT MOD(10, 2 BETWEEN 0 and 5)", true, "SELECT (10 % (2 BETWEEN 0 AND 5))"}, // issue #59000
		{"select mod( year(a) - abs(weekday(a) + dayofweek(a)), 4) + 1", true, "SELECT (((year(`a`) - abs((weekday(`a`) + dayofweek(`a`)))) % 4) + 1)"},
	}

	p := parser.New()
	p.EnableWindowFunc(false)
	for _, tbl := range cases {
		_, _, err := p.Parse(tbl.src, "", "")
		comment := fmt.Sprintf("source %v", tbl.src)
		if !tbl.ok {
			if err == nil {
				t.Fatal(comment)
			}
			continue
		}
		if err != nil {
			t.Fatalf("%v: %v", comment, err)
		}
		// restore correctness test
		if tbl.ok {
			var sb strings.Builder
			comment := fmt.Sprintf("source %v", tbl.src)
			stmts, _, err := p.Parse(tbl.src, "", "")
			if err != nil {
				t.Fatalf("%v: %v", comment, err)
			}
			restoreSQLs := ""
			for _, stmt := range stmts {
				sb.Reset()
				ctx := NewRestoreCtx(RestoreStringSingleQuotes|RestoreSpacesAroundBinaryOperation|RestoreBracketAroundBinaryOperation|RestoreStringWithoutCharset|RestoreNameBackQuotes, &sb)
				ctx.DefaultDB = "test"
				err = stmt.Restore(ctx)
				if err != nil {
					t.Fatalf("%v: %v", comment, err)
				}
				restoreSQL := sb.String()
				comment = fmt.Sprintf("source %v; restore %v", tbl.src, restoreSQL)
				if restoreSQLs != "" {
					restoreSQLs += "; "
				}
				restoreSQLs += restoreSQL
			}
			comment = fmt.Sprintf("restore %v; expect %v", restoreSQLs, tbl.restore)
			if !reflect.DeepEqual(tbl.restore, restoreSQLs) {
				t.Fatalf("%v: got %v, want %v", comment, restoreSQLs, tbl.restore)
			}
		}
	}
}

// For CTE bindings.
func TestCTEBindings(t *testing.T) {
	table := []testCase{
		{"WITH `cte` AS (SELECT * from t) SELECT `col1`,`col2` FROM `cte`", true, "WITH `cte` AS (SELECT * FROM `test`.`t`) SELECT `col1`,`col2` FROM `cte`"},
		{"WITH `cte` (col1, col2) AS (SELECT * from t UNION ALL SELECT 3,4) SELECT col1, col2 FROM cte;", true, "WITH `cte` (`col1`, `col2`) AS (SELECT * FROM `test`.`t` UNION ALL SELECT 3,4) SELECT `col1`,`col2` FROM `cte`"},
		{"WITH `cte` AS (SELECT * from t), cte2 as (select * from cte) SELECT `col1`,`col2` FROM `cte`", true, "WITH `cte` AS (SELECT * FROM `test`.`t`), `cte2` AS (SELECT * FROM `cte`) SELECT `col1`,`col2` FROM `cte`"},
		{"WITH RECURSIVE cte (n) AS (  SELECT * from t  UNION ALL  SELECT n + 1 FROM cte WHERE n < 5)SELECT * FROM cte;", true, "WITH RECURSIVE `cte` (`n`) AS (SELECT * FROM `test`.`t` UNION ALL SELECT `n` + 1 FROM `cte` WHERE `n` < 5) SELECT * FROM `cte`"},
		{"with cte(a) as (select * from t) update t, cte set t.a=1  where t.a=cte.a;", true, "WITH `cte` (`a`) AS (SELECT * FROM `test`.`t`) UPDATE (`test`.`t`) JOIN `cte` SET `t`.`a`=1 WHERE `t`.`a` = `cte`.`a`"},
		{"with cte(a) as (select * from t) delete t from t, cte where t.a=cte.a;", true, "WITH `cte` (`a`) AS (SELECT * FROM `test`.`t`) DELETE `test`.`t` FROM (`test`.`t`) JOIN `cte` WHERE `t`.`a` = `cte`.`a`"},
		{"WITH cte1 AS (SELECT * from t) SELECT * FROM (WITH cte2 AS (SELECT * from cte1) SELECT * FROM cte2 JOIN cte1) AS dt;", true, "WITH `cte1` AS (SELECT * FROM `test`.`t`) SELECT * FROM (WITH `cte2` AS (SELECT * FROM `cte1`) SELECT * FROM `cte2` JOIN `cte1`) AS `dt`"},
		{"WITH cte AS (SELECT * from t) SELECT /*+ MAX_EXECUTION_TIME(1000) */ * FROM cte;", true, "WITH `cte` AS (SELECT * FROM `test`.`t`) SELECT /*+ MAX_EXECUTION_TIME(1000)*/ * FROM `cte`"},
		{"with cte as (table t) table cte;", true, "WITH `cte` AS (TABLE `test`.`t`) TABLE `cte`"},
		{"with cte as (select * from t) select 1 union with cte as (select * from t) select * from cte;", false, ""},
		{"with cte as (select * from t) (select * from t);", true, "WITH `cte` AS (SELECT * FROM `test`.`t`) (SELECT * FROM `test`.`t`)"},
		{"with cte as (select 1) (select 1 union select * from t)", true, "WITH `cte` AS (SELECT 1) (SELECT 1 UNION SELECT * FROM `test`.`t`)"},
		{"select * from (with cte as (select * from t) select 1 union select * from t) qn", true, "SELECT * FROM (WITH `cte` AS (SELECT * FROM `test`.`t`) SELECT 1 UNION SELECT * FROM `test`.`t`) AS `qn`"},
		{"select * from t where 1 > (with cte as (select * from t) select * from cte)", true, "SELECT * FROM `test`.`t` WHERE 1 > (WITH `cte` AS (SELECT * FROM `test`.`t`) SELECT * FROM `cte`)"},
		{"( with cte(n) as ( select * from t )  select n+1 from cte  union select n+2 from cte) union select 1", true, "(WITH `cte` (`n`) AS (SELECT * FROM `test`.`t`) SELECT `n` + 1 FROM `cte` UNION SELECT `n` + 2 FROM `cte`) UNION SELECT 1"},
		{"( with cte(n) as ( select * from t )  select n+1 from cte) union select * from t", true, "(WITH `cte` (`n`) AS (SELECT * FROM `test`.`t`) SELECT `n` + 1 FROM `cte`) UNION SELECT * FROM `test`.`t`"},
		{"with cte as (select * from t union select * from cte) select * from cte", true, "WITH `cte` AS (SELECT * FROM `test`.`t` UNION SELECT * FROM `test`.`cte`) SELECT * FROM `cte`"},
	}

	p := parser.New()
	p.EnableWindowFunc(false)
	for _, tbl := range table {
		_, _, err := p.Parse(tbl.src, "", "")
		comment := fmt.Sprintf("source %v", tbl.src)
		if !tbl.ok {
			if err == nil {
				t.Fatal(comment)
			}
			continue
		}
		if err != nil {
			t.Fatalf("%v: %v", comment, err)
		}
		// restore correctness test
		if tbl.ok {
			var sb strings.Builder
			comment := fmt.Sprintf("source %v", tbl.src)
			stmts, _, err := p.Parse(tbl.src, "", "")
			if err != nil {
				t.Fatalf("%v: %v", comment, err)
			}
			restoreSQLs := ""
			for _, stmt := range stmts {
				sb.Reset()
				ctx := NewRestoreCtx(RestoreStringSingleQuotes|RestoreSpacesAroundBinaryOperation|RestoreStringWithoutCharset|RestoreNameBackQuotes, &sb)
				ctx.DefaultDB = "test"
				err = stmt.Restore(ctx)
				if err != nil {
					t.Fatalf("%v: %v", comment, err)
				}
				restoreSQL := sb.String()
				comment = fmt.Sprintf("source %v; restore %v", tbl.src, restoreSQL)
				if restoreSQLs != "" {
					restoreSQLs += "; "
				}
				restoreSQLs += restoreSQL
			}
			comment = fmt.Sprintf("restore %v; expect %v", restoreSQLs, tbl.restore)
			if !reflect.DeepEqual(tbl.restore, restoreSQLs) {
				t.Fatalf("%v: got %v, want %v", comment, restoreSQLs, tbl.restore)
			}
		}
	}
}

func TestPlanReplayer(t *testing.T) {
	p := parser.New()
	sms, _, err := p.Parse("PLAN REPLAYER DUMP EXPLAIN SELECT a FROM t", "", "")
	if err != nil {
		t.Fatal(err)
	}
	v, ok := sms[0].(*ast.PlanReplayerStmt)
	if !(ok) {
		t.Fatal("expected true")
	}
	if !reflect.DeepEqual("SELECT a FROM t", v.Stmt.Text()) {
		t.Fatalf("got %v, want %v", v.Stmt.Text(), "SELECT a FROM t")
	}
	if v.Analyze {
		t.Fatal("expected false")
	}

	sms, _, err = p.Parse("PLAN REPLAYER DUMP EXPLAIN ANALYZE SELECT a FROM t", "", "")
	if err != nil {
		t.Fatal(err)
	}
	v, ok = sms[0].(*ast.PlanReplayerStmt)
	if !(ok) {
		t.Fatal("expected true")
	}
	if !reflect.DeepEqual("SELECT a FROM t", v.Stmt.Text()) {
		t.Fatalf("got %v, want %v", v.Stmt.Text(), "SELECT a FROM t")
	}
	if !(v.Analyze) {
		t.Fatal("expected true")
	}

	// Multiple SQL records: EXPLAIN ( "sql1", "sql2", ... )
	sms, _, err = p.Parse("PLAN REPLAYER DUMP EXPLAIN ('SELECT * FROM t1', 'SELECT * FROM t2')", "", "")
	if err != nil {
		t.Fatal(err)
	}
	v, ok = sms[0].(*ast.PlanReplayerStmt)
	if !(ok) {
		t.Fatal("expected true")
	}
	if v.Stmt != nil {
		t.Fatalf("expected nil, got %v", v.Stmt)
	}
	if v.Analyze {
		t.Fatal("expected false")
	}
	if !reflect.DeepEqual([]string{"SELECT * FROM t1", "SELECT * FROM t2"}, v.StmtList) {
		t.Fatalf("got %v, want %v", v.StmtList, []string{"SELECT * FROM t1", "SELECT * FROM t2"})
	}

	sms, _, err = p.Parse("PLAN REPLAYER DUMP EXPLAIN ANALYZE ('SELECT * FROM t1')", "", "")
	if err != nil {
		t.Fatal(err)
	}
	v, ok = sms[0].(*ast.PlanReplayerStmt)
	if !(ok) {
		t.Fatal("expected true")
	}
	if v.Stmt != nil {
		t.Fatalf("expected nil, got %v", v.Stmt)
	}
	if !(v.Analyze) {
		t.Fatal("expected true")
	}
	if !reflect.DeepEqual([]string{"SELECT * FROM t1"}, v.StmtList) {
		t.Fatalf("got %v, want %v", v.StmtList, []string{"SELECT * FROM t1"})
	}
}

func TestTrafficStmt(t *testing.T) {
	table := []testCase{
		{"traffic capture to '/tmp' duration='1s' encryption_method='aes' compress=true", true, "TRAFFIC CAPTURE TO '/tmp' DURATION = '1s' ENCRYPTION_METHOD = 'aes' COMPRESS = TRUE"},
		{"traffic capture to '/tmp' duration '1s' encryption_method 'aes' compress true", true, "TRAFFIC CAPTURE TO '/tmp' DURATION = '1s' ENCRYPTION_METHOD = 'aes' COMPRESS = TRUE"},
		{"traffic capture to '/tmp' encryption_method='aes' duration='1s'", true, "TRAFFIC CAPTURE TO '/tmp' ENCRYPTION_METHOD = 'aes' DURATION = '1s'"},
		{"traffic capture to '/tmp' duration='1m'", true, "TRAFFIC CAPTURE TO '/tmp' DURATION = '1m'"},
		{"traffic capture to '/tmp' duration='1'", false, ""},
		{"traffic capture to '/tmp' duration=1s", false, ""},
		{"traffic capture to '/tmp' compress='true'", false, ""},
		{"traffic capture duration='1m'", false, ""},
		{"traffic capture", false, ""},
		{"traffic replay from '/tmp' user='root' password='123456' speed=1.0 read_only=true", true, "TRAFFIC REPLAY FROM '/tmp' USER = 'root' PASSWORD = '123456' SPEED = 1.0 READONLY = TRUE"},
		{"traffic replay from '/tmp' user 'root' password '123456' speed 1.0 read_only true", true, "TRAFFIC REPLAY FROM '/tmp' USER = 'root' PASSWORD = '123456' SPEED = 1.0 READONLY = TRUE"},
		{"traffic replay from '/tmp' speed 1.0 user='root'", true, "TRAFFIC REPLAY FROM '/tmp' SPEED = 1.0 USER = 'root'"},
		{"traffic replay from '/tmp' speed=1", true, "TRAFFIC REPLAY FROM '/tmp' SPEED = 1"},
		{"traffic replay from '/tmp' speed=0.5", true, "TRAFFIC REPLAY FROM '/tmp' SPEED = 0.5"},
		{"traffic replay from '/tmp' speed=-1", false, ""},
		{"traffic replay speed=1", false, ""},
		{"traffic replay", false, ""},
		{"show traffic jobs", true, "SHOW TRAFFIC JOBS"},
		{"show traffic jobs duration='1m'", false, ""},
		{"show traffic", false, ""},
		{"cancel traffic jobs", true, "CANCEL TRAFFIC JOBS"},
		{"cancel traffic jobs duration='1m'", false, ""},
		{"cancel traffic", false, ""},
		{"traffic test", false, ""},
		{"traffic", false, ""},
	}

	p := parser.New()
	var sb strings.Builder
	for _, tbl := range table {
		stmts, _, err := p.Parse(tbl.src, "", "")
		if !tbl.ok {
			if err == nil {
				t.Fatal(tbl.src)
			}
			continue
		}
		if err != nil {
			t.Fatalf("%v: %v", tbl.src, err)
		}
		if got := len(stmts); got != 1 {
			t.Fatalf("expected length %d, got %d", 1, got)
		}
		v, ok := stmts[0].(*ast.TrafficStmt)
		if !(ok) {
			t.Fatal("expected true")
		}
		switch v.OpType {
		case ast.TrafficOpCapture, ast.TrafficOpReplay:
			if !reflect.DeepEqual("/tmp", v.Dir) {
				t.Fatalf("got %v, want %v", v.Dir, "/tmp")
			}
		}
		sb.Reset()
		ctx := NewRestoreCtx(RestoreStringSingleQuotes|RestoreSpacesAroundBinaryOperation|RestoreStringWithoutCharset|RestoreNameBackQuotes, &sb)
		err = v.Restore(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(tbl.restore, sb.String()) {
			t.Fatalf("got %v, want %v", sb.String(), tbl.restore)
		}
	}
}

func TestGBKEncoding(t *testing.T) {
	p := parser.New()
	gbkEncoding, _ := charset.Lookup("gbk")
	encoder := gbkEncoding.NewEncoder()
	sql, err := encoder.String("create table 测试表 (测试列 varchar(255) default 'GBK测试用例');")
	if err != nil {
		t.Fatal(err)
	}

	stmt, _, err := p.ParseSQL(sql)
	if err != nil {
		t.Fatal(err)
	}
	checker := &gbkEncodingChecker{}
	_, _ = stmt[0].Accept(checker)
	if reflect.DeepEqual("测试表", checker.tblName) {
		t.Fatalf("expected values to differ, both are %v", checker.tblName)
	}
	if reflect.DeepEqual("测试列", checker.colName) {
		t.Fatalf("expected values to differ, both are %v", checker.colName)
	}

	gbkOpt := parser.CharsetClient("gbk")
	stmt, _, err = p.ParseSQL(sql, gbkOpt)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = stmt[0].Accept(checker)
	if !reflect.DeepEqual("测试表", checker.tblName) {
		t.Fatalf("got %v, want %v", checker.tblName, "测试表")
	}
	if !reflect.DeepEqual("测试列", checker.colName) {
		t.Fatalf("got %v, want %v", checker.colName, "测试列")
	}
	if !reflect.DeepEqual("GBK测试用例", checker.expr) {
		t.Fatalf("got %v, want %v", checker.expr, "GBK测试用例")
	}

	_, _, err = p.ParseSQL("select _gbk '\xc6\x5c' from dual;")
	if err == nil {
		t.Fatal("expected error")
	}

	for _, test := range []struct {
		sql string
		err bool
	}{
		{"select '\xc6\x5c' from `\xab\x60`;", false},
		{`prepare p1 from "insert into t values ('中文');";`, false},
		{"select '啊';", false},
		{"create table t1(s set('a一','b二','c三'));", false},
		{"insert into t3 values('一a');", false},
		{"select '\xa5\x5c'", false},
		{"select '''\xa5\x5c'", false},
		{"select ```\xa5\x5c`", false},
		{"select '\x65\x5c'", true},
	} {
		_, _, err = p.ParseSQL(test.sql, gbkOpt)
		if test.err {
			if err == nil {
				t.Fatal(test.sql)
			}
		} else {
			if err != nil {
				t.Fatalf("%v: %v", test.sql, err)
			}
		}
	}
}

func TestGB18030Encoding(t *testing.T) {
	p := parser.New()
	gb18030Encoding, _ := charset.Lookup("gb18030")
	encoder := gb18030Encoding.NewEncoder()
	sql, err := encoder.String("create table 测试表 (测试列 varchar(255) default 'GB18030测试用例');")
	if err != nil {
		t.Fatal(err)
	}

	stmt, _, err := p.ParseSQL(sql)
	if err != nil {
		t.Fatal(err)
	}
	checker := &gbkEncodingChecker{}
	_, _ = stmt[0].Accept(checker)
	if reflect.DeepEqual("测试表", checker.tblName) {
		t.Fatalf("expected values to differ, both are %v", checker.tblName)
	}
	if reflect.DeepEqual("测试列", checker.colName) {
		t.Fatalf("expected values to differ, both are %v", checker.colName)
	}

	gb18030Opt := parser.CharsetClient("gb18030")
	stmt, _, err = p.ParseSQL(sql, gb18030Opt)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = stmt[0].Accept(checker)
	if !reflect.DeepEqual("测试表", checker.tblName) {
		t.Fatalf("got %v, want %v", checker.tblName, "测试表")
	}
	if !reflect.DeepEqual("测试列", checker.colName) {
		t.Fatalf("got %v, want %v", checker.colName, "测试列")
	}
	if !reflect.DeepEqual("GB18030测试用例", checker.expr) {
		t.Fatalf("got %v, want %v", checker.expr, "GB18030测试用例")
	}

	_, _, err = p.ParseSQL("select _gbk '\xc6\x5c' from dual;")
	if err == nil {
		t.Fatal("expected error")
	}

	for _, test := range []struct {
		sql string
		err bool
	}{
		{"select '\xc6\x5c' from `\xab\x60`;", false},
		{`prepare p1 from "insert into t values ('中文');";`, false},
		{"select '啊';", false},
		{"create table t1(s set('a一','b二','c三'));", false},
		{"insert into t3 values('一a');", false},
		{"select '\xa5\x5c'", false},
		{"select '''\xa5\x5c'", false},
		{"select ```\xa5\x5c`", false},
		{"select '\x65\x5c'", true},
	} {
		_, _, err = p.ParseSQL(test.sql, gb18030Opt)
		if test.err {
			if err == nil {
				t.Fatal(test.sql)
			}
		} else {
			if err != nil {
				t.Fatalf("%v: %v", test.sql, err)
			}
		}
	}
}

type gbkEncodingChecker struct {
	tblName string
	colName string
	expr    string
}

func (g *gbkEncodingChecker) Enter(n ast.Node) (node ast.Node, skipChildren bool) {
	if tn, ok := n.(*ast.TableName); ok {
		g.tblName = tn.Name.O
		return n, false
	}
	if cn, ok := n.(*ast.ColumnName); ok {
		g.colName = cn.Name.O
		return n, false
	}
	if c, ok := n.(*ast.ColumnOption); ok {
		if ve, ok := c.Expr.(ast.ValueExpr); ok {
			g.expr = ve.GetString()
			return n, false
		}
	}
	return n, false
}

func (g *gbkEncodingChecker) Leave(n ast.Node) (node ast.Node, ok bool) {
	return n, true
}

func TestInsertStatementMemoryAllocation(t *testing.T) {
	sql := "insert t values (1)" + strings.Repeat(",(1)", 1000)
	var oldStats, newStats runtime.MemStats
	runtime.ReadMemStats(&oldStats)
	_, err := parser.New().ParseOneStmt(sql, "", "")
	if err != nil {
		t.Fatal(err)
	}
	runtime.ReadMemStats(&newStats)
	if !(int(newStats.TotalAlloc-oldStats.TotalAlloc) < 1024*500) {
		t.Fatalf("expected %v < %v", int(newStats.TotalAlloc-oldStats.TotalAlloc), 1024*500)
	}
}

func TestCharsetIntroducer(t *testing.T) {
	p := parser.New()
	defer charset.RemoveCharset("gbk")
	// `_gbk` is treated as a character set.
	_, _, err := p.Parse("select _gbk 'a';", "", "")
	if err == nil || err.Error() != "[ddl:1115]Unsupported character introducer: 'gbk'" {
		t.Fatalf("expected error %q, got %v", "[ddl:1115]Unsupported character introducer: 'gbk'", err)
	}
	_, _, err = p.Parse("select _gbk 0x1234;", "", "")
	if err == nil || err.Error() != "[ddl:1115]Unsupported character introducer: 'gbk'" {
		t.Fatalf("expected error %q, got %v", "[ddl:1115]Unsupported character introducer: 'gbk'", err)
	}
	_, _, err = p.Parse("select _gbk 0b101001;", "", "")
	if err == nil || err.Error() != "[ddl:1115]Unsupported character introducer: 'gbk'" {
		t.Fatalf("expected error %q, got %v", "[ddl:1115]Unsupported character introducer: 'gbk'", err)
	}
}

func TestIssue45898(t *testing.T) {
	p := parser.New()
	p.ParseSQL("a.")
	stmts, _, err := p.ParseSQL("select count(1) from t")
	if err != nil {
		t.Fatal(err)
	}
	var sb strings.Builder
	restoreCtx := NewRestoreCtx(DefaultRestoreFlags, &sb)
	sb.Reset()
	stmts[0].Restore(restoreCtx)
	if !reflect.DeepEqual("SELECT COUNT(1) FROM `t`", sb.String()) {
		t.Fatalf("got %v, want %v", sb.String(), "SELECT COUNT(1) FROM `t`")
	}
}

func TestMultiStmt(t *testing.T) {
	p := parser.New()
	stmts, _, err := p.Parse("SELECT 'foo'; SELECT 'foo;bar','baz'; select 'foo' , 'bar' , 'baz' ;select 1", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(len(stmts), 4) {
		t.Fatalf("got %v, want %v", 4, len(stmts))
	}
	stmt1 := stmts[0].(*ast.SelectStmt)
	stmt2 := stmts[1].(*ast.SelectStmt)
	stmt3 := stmts[2].(*ast.SelectStmt)
	stmt4 := stmts[3].(*ast.SelectStmt)
	if !reflect.DeepEqual("'foo'", stmt1.Fields.Fields[0].Text()) {
		t.Fatalf("got %v, want %v", stmt1.Fields.Fields[0].Text(), "'foo'")
	}
	if !reflect.DeepEqual("'foo;bar'", stmt2.Fields.Fields[0].Text()) {
		t.Fatalf("got %v, want %v", stmt2.Fields.Fields[0].Text(), "'foo;bar'")
	}
	if !reflect.DeepEqual("'baz'", stmt2.Fields.Fields[1].Text()) {
		t.Fatalf("got %v, want %v", stmt2.Fields.Fields[1].Text(), "'baz'")
	}
	if !reflect.DeepEqual("'foo'", stmt3.Fields.Fields[0].Text()) {
		t.Fatalf("got %v, want %v", stmt3.Fields.Fields[0].Text(), "'foo'")
	}
	if !reflect.DeepEqual("'bar'", stmt3.Fields.Fields[1].Text()) {
		t.Fatalf("got %v, want %v", stmt3.Fields.Fields[1].Text(), "'bar'")
	}
	if !reflect.DeepEqual("'baz'", stmt3.Fields.Fields[2].Text()) {
		t.Fatalf("got %v, want %v", stmt3.Fields.Fields[2].Text(), "'baz'")
	}
	if !reflect.DeepEqual("1", stmt4.Fields.Fields[0].Text()) {
		t.Fatalf("got %v, want %v", stmt4.Fields.Fields[0].Text(), "1")
	}
}
