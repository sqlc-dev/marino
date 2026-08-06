// Copyright 2026 The sqlc Authors.
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

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/sqlc-dev/marino/mysql"
)

// rdHintCorpus is every hint-looking string from hintparser_test.go's
// tables plus constructed edge cases: all productions, all list/comma
// variants, and erroneous inputs (the yacc hint grammar has no error
// recovery, so any syntax error must abort with a nil result and the
// identical error list).
var rdHintCorpus = []string{
	// From TestParseHint in hintparser_test.go.
	"",
	"MEMORY_QUOTA(8 MB) MEMORY_QUOTA(6 GB)",
	"QB_NAME(qb1) QB_NAME(`qb2`), QB_NAME(TRUE) QB_NAME(\"ANSI quoted\") QB_NAME(_utf8), QB_NAME(0b10) QB_NAME(0x1a)",
	"QB_NAME(1)",
	"QB_NAME('string literal')",
	"QB_NAME(many identifiers)",
	"QB_NAME(@qb1)",
	"QB_NAME(b'10')",
	"QB_NAME(x'1a')",
	"JOIN_FIXED_ORDER() BKA()",
	"HASH_JOIN() TIDB_HJ(@qb1) INL_JOIN(x, `y y`.z) MERGE_JOIN(w@`First QB`)",
	"USE_INDEX_MERGE(@qb1 tbl1 x, y, z) IGNORE_INDEX(tbl2@qb2) USE_INDEX(tbl3 PRIMARY) FORCE_INDEX(tbl4@qb3 c1) INDEX_LOOKUP_PUSHDOWN(tbl5@qb6 c3)",
	"USE_INDEX(@qb1 tbl1 partition(p0) x) USE_INDEX_MERGE(@qb2 tbl2@qb2 partition(p0, p1) x, y, z)",
	`SET_VAR(sbs = 16M) SET_VAR(fkc=OFF) SET_VAR(os="mcb=off") set_var(abc=1) set_var(os2='mcb2=off')`,
	"USE_TOJA(TRUE) IGNORE_PLAN_CACHE() USE_CASCADES(TRUE) QUERY_TYPE(@qb1 OLAP) QUERY_TYPE(OLTP) NO_INDEX_MERGE() RESOURCE_GROUP(rg1)",
	"READ_FROM_STORAGE(@foo TIKV[a, b], TIFLASH[c, d]) HASH_AGG() SEMI_JOIN_REWRITE() READ_FROM_STORAGE(TIKV[e])",
	"WRITE_SLOW_LOG, WRITE_SLOW_LOG",
	"WRITE_SLOW_LOG()",
	"unknown_hint()",
	"set_var(timestamp = 1.5)",
	"set_var(timestamp = _utf8mb4'1234')",
	"set_var(timestamp = 9999999999999999999999999999999999999)",
	"time_range('2020-02-20 12:12:12',456)",
	"time_range(456,'2020-02-20 12:12:12')",
	"TIME_RANGE('2020-02-20 12:12:12','2020-02-20 13:12:12')",
	"LEADING(a,(b,(c,d)))",
	"LEADING(a,b,c)",
	"LEADING((a,b),(c,d))",
	"LEADING(x,(y,z),w)",

	// JOIN_FIXED_ORDER / join-order hints.
	"JOIN_FIXED_ORDER()",
	"JOIN_FIXED_ORDER(@qb)",
	"JOIN_FIXED_ORDER(t1)",
	"JOIN_ORDER(t1, t2)",
	"JOIN_ORDER(@qb t1, d.t2@q2)",
	"JOIN_ORDER()",
	"JOIN_PREFIX(t1)",
	"JOIN_SUFFIX(t1 , t2)",

	// Table-level hints and HintTableListOpt shapes.
	"BKA()",
	"BKA(t1)",
	"NO_BKA(@qb t1)",
	"BNL(t1, t2)",
	"NO_BNL(@qb)",
	"NO_MERGE(t)",
	"MERGE()",
	"MERGE_JOIN(t1)",
	"NO_MERGE_JOIN(t1)",
	"BROADCAST_JOIN(t1, t2)",
	"SHUFFLE_JOIN(t1)",
	"INL_JOIN(db.t@qb)",
	"INDEX_JOIN(t1)",
	"NO_INDEX_JOIN(t1)",
	"INL_HASH_JOIN(t1)",
	"INDEX_HASH_JOIN(t1)",
	"NO_INDEX_HASH_JOIN(t1)",
	"SWAP_JOIN_INPUTS(t1)",
	"NO_SWAP_JOIN_INPUTS(t1)",
	"INL_MERGE_JOIN(t1)",
	"INDEX_MERGE_JOIN(t1)",
	"NO_INDEX_MERGE_JOIN(t1)",
	"HASH_JOIN_BUILD(t1)",
	"HASH_JOIN_PROBE(t1)",
	"HYPO_INDEX(t1)",
	"TIDB_SMJ(t1)",
	"TIDB_INLJ(t1)",
	"HASH_JOIN(@qb1)",
	"HASH_JOIN(@qb1 t1)",
	"HASH_JOIN(@qb1, t1)",
	"HASH_JOIN(t1@qb1, t2@qb2)",
	"HASH_JOIN(BKA)",
	"HASH_JOIN(TRUE.FALSE)",
	"HASH_JOIN(t1,)",
	"HASH_JOIN(t1 t2)",
	"MERGE_JOIN(t1@qb1 partition(p0))",
	"MERGE_JOIN(d.t1@qb1 partition(p0 p1, p2))",
	"MERGE_JOIN(t1 partition())",
	"MERGE_JOIN(t1 partition(p0,))",
	"MERGE_JOIN(PARTITION)",
	"MERGE_JOIN(t1.PARTITION)",

	// LEADING.
	"LEADING(@qb a, b)",
	"LEADING((a))",
	"LEADING(t1@qb1)",
	"LEADING(d.t)",
	"LEADING(d.t@qb PARTITION(p0))",
	"LEADING(((a,b),c),d)",
	"LEADING()",
	"LEADING(a b)",
	"LEADING(a,(b,c)",
	"LEADING(a,())",
	"leading(a)",

	// Index-level hints.
	"INDEX_MERGE(t1 x)",
	"MRR(t1 x)",
	"NO_MRR(t1 x)",
	"NO_ICP(t1 x)",
	"NO_RANGE_OPTIMIZATION(t1 x)",
	"SKIP_SCAN(t1 x)",
	"NO_SKIP_SCAN(t1 x)",
	"ORDER_INDEX(t1 x)",
	"NO_ORDER_INDEX(t1 x)",
	"NO_INDEX_LOOKUP_PUSHDOWN(t1 x)",
	"USE_INDEX()",
	"USE_INDEX(@qb)",
	"USE_INDEX(t1)",
	"USE_INDEX(t1,)",
	"USE_INDEX(t1, x)",
	"USE_INDEX(t1 x, y)",
	"USE_INDEX(t1 a b)",
	"USE_INDEX(t1 x,)",
	"USE_INDEX(t1 PARTITION(p0) idx1)",
	"USE_INDEX(t1 PARTITION(p0 p1) idx1, idx2)",
	"USE_INDEX(d.t1@qb x)",
	"USE_INDEX(t1 MB, GB)",
	"use_index(t x)",

	// Subquery hints.
	"SEMIJOIN()",
	"SEMIJOIN(FIRSTMATCH, LOOSESCAN)",
	"SEMIJOIN(@qb DUPSWEEDOUT)",
	"NO_SEMIJOIN(MATERIALIZATION)",
	"SEMIJOIN(MATERIALIZATION FIRSTMATCH)",
	"SEMIJOIN(x)",
	"SEMIJOIN(FIRSTMATCH,)",

	// MAX_EXECUTION_TIME / NTH_PLAN.
	"MAX_EXECUTION_TIME(1000)",
	"MAX_EXECUTION_TIME(@qb 100)",
	"MAX_EXECUTION_TIME(x)",
	"MAX_EXECUTION_TIME()",
	"MAX_EXECUTION_TIME(18446744073709551615)",
	"MAX_EXECUTION_TIME(99999999999999999999999999)",
	"NTH_PLAN(2)",
	"NTH_PLAN(@qb 3)",
	"NTH_PLAN(-1)",

	// SET_VAR and Value.
	"SET_VAR(x=y)",
	"SET_VAR(x = 1)",
	"SET_VAR(x = +5)",
	"SET_VAR(x = -5)",
	"SET_VAR(x = -9223372036854775807)",
	"SET_VAR(x = -9223372036854775808)",
	"SET_VAR(x = -9223372036854775809)",
	"SET_VAR(x = -18446744073709551615)",
	"SET_VAR(x = -9223372036854775809 y)",
	"SET_VAR(x = -9223372036854775808 y)",
	"SET_VAR(x = +)",
	"SET_VAR(x = -)",
	"SET_VAR(x =)",
	"SET_VAR(x)",
	"SET_VAR(PARTITION=1)",
	"SET_VAR(TIKV=OLAP)",
	"SET_VAR(x=`quoted`)",
	"SET_VAR(a=b) SET_VAR(",

	// RESOURCE_GROUP / QB_NAME.
	"RESOURCE_GROUP(TRUE)",
	"RESOURCE_GROUP()",
	"RESOURCE_GROUP(rg1, rg2)",
	"QB_NAME(a,)",
	"QB_NAME(a, .v1)",
	"QB_NAME(a, v1.v2@qb.v3)",
	"QB_NAME(a, @qb)",
	"QB_NAME(a, v1@qb1)",
	"QB_NAME(a, v1,v2)",

	// MEMORY_QUOTA.
	"MEMORY_QUOTA(8 MB)",
	"MEMORY_QUOTA(@qb 8 GB)",
	"MEMORY_QUOTA(8589934592 GB)",
	"MEMORY_QUOTA(8589934592 GB) HASH_AGG()",
	"HASH_AGG() MEMORY_QUOTA(8589934592 GB), STREAM_AGG()",
	"MEMORY_QUOTA(9223372036854775807 MB) BKA(t)",
	"MEMORY_QUOTA(8589934592 GB",
	"MEMORY_QUOTA(8589934592 GB x)",
	"MEMORY_QUOTA(1)",
	"MEMORY_QUOTA(1 KB)",
	"MEMORY_QUOTA(MB MB)",

	// TIME_RANGE.
	"TIME_RANGE('a','b')",
	"TIME_RANGE('a' 'b')",
	"TIME_RANGE('a')",
	"TIME_RANGE(a,'b')",
	"TIME_RANGE('a','b','c')",

	// Boolean/nullary/query-type hints.
	"USE_TOJA(TRUE)",
	"USE_TOJA(@qb FALSE)",
	"USE_CASCADES(FALSE)",
	"USE_TOJA(1)",
	"USE_TOJA()",
	"USE_PLAN_CACHE()",
	"MPP_1PHASE_AGG()",
	"MPP_2PHASE_AGG()",
	"STREAM_AGG(@qb)",
	"AGG_TO_COP()",
	"LIMIT_TO_COP()",
	"READ_CONSISTENT_REPLICA()",
	"STRAIGHT_JOIN()",
	"NO_DECORRELATE()",
	"NO_INDEX_MERGE(t1 x)",
	"WRITE_SLOW_LOG",
	"WRITE_SLOW_LOG WRITE_SLOW_LOG",
	"Write_Slow_Log",
	"QUERY_TYPE(OLAP)",
	"QUERY_TYPE(olap)",
	"QUERY_TYPE(x)",
	"QUERY_TYPE()",

	// READ_FROM_STORAGE.
	"READ_FROM_STORAGE(TIKV[t1])",
	"READ_FROM_STORAGE(tikv[t1])",
	"READ_FROM_STORAGE(TIKV[t1 partition(p0)])",
	"READ_FROM_STORAGE(TIKV[@qb t1], TIFLASH[t2@qb2])",
	"READ_FROM_STORAGE(TIKV[d.t1, t2])",
	"READ_FROM_STORAGE(TIKV)",
	"READ_FROM_STORAGE(OLAP[t])",
	"READ_FROM_STORAGE(TIKV[t],)",
	"READ_FROM_STORAGE(TIKV[])",
	"READ_FROM_STORAGE(TIKV[t] TIFLASH[u])",
	"READ_FROM_STORAGE()",

	// Pseudo (unknown identifier) hints.
	"unknown_hint(123)",
	"unknown_hint(@qb 1)",
	"unknown_hint(@qb x)",
	"unknown_hint(p1, p2)",
	"unknown_hint(p1 p2)",
	"unknown_hint(p1, p2, 10)",
	"unknown_hint(p1 10)",
	"unknown_hint(p1 10 p2)",
	"unknown_hint(x=1)",
	"unknown_hint(x='str')",
	"unknown_hint(x=-3)",
	"unknown_hint(x=-9223372036854775809)",
	"unknown_hint(x, 1, y)",
	"unknown_hint(1,2)",
	"unknown_hint(x=)",
	"unknown_hint(TRUE, FALSE)",
	"unknown_hint(a,)",
	"unknown_hint",
	"unknown_hint(",
	"unknown_hint MERGE()",

	// Separators, stray tokens, misc errors.
	"HASH_AGG() HASH_AGG()",
	"HASH_AGG(), HASH_AGG()",
	"HASH_AGG(),, HASH_AGG()",
	",HASH_AGG()",
	"HASH_AGG(),",
	"HASH_AGG() )",
	"HASH_AGG() BKA(t) unknown_hint(1) USE_INDEX(t x)",
	"BKA(t) %",
	"5",
	"PARTITION",
	"TRUE",
	"@qb",
	"(",
	"=",
	"USE_INDEX(t1",
	"QB_NAME(",
	"   ",
	"\n\nQB_NAME(\nqb1\n)",
	"QB_NAME(\n",
}

// hintDiffModes are the sqlMode/initPos combinations each corpus entry is
// parsed under.
var hintDiffCases = []struct {
	mode mysql.SQLMode
	pos  Pos
}{
	{0, Pos{Line: 1}},
	{mysql.ModeANSIQuotes, Pos{Line: 1}},
	{0, Pos{Line: 3, Col: 11, Offset: 42}},
}

// TestRDParseHintDifferential parses every corpus entry with both the
// goyacc hint parser and rdParseHint and requires identical hints (deep
// equality) and identical error lists (by Error() string).
func TestRDParseHintDifferential(t *testing.T) {
	for _, tc := range hintDiffCases {
		for _, frag := range rdHintCorpus {
			input := "/*+" + frag + "*/"

			yHints, yErrs := newHintParser().parse(input, tc.mode, tc.pos)
			rHints, rErrs := rdParseHint(input, tc.mode, tc.pos)

			id := fmt.Sprintf("input=%q mode=%v pos=%+v", frag, tc.mode, tc.pos)
			if !reflect.DeepEqual(yHints, rHints) {
				t.Fatalf("%s:\nhints diverge\nyacc: %#v\nrd  : %#v", id, yHints, rHints)
			}
			if len(yErrs) != len(rErrs) {
				t.Fatalf("%s:\nerror count diverges\nyacc: %q\nrd  : %q", id, errStrings(yErrs), errStrings(rErrs))
			}
			for i := range yErrs {
				if yErrs[i].Error() != rErrs[i].Error() {
					t.Fatalf("%s:\nerror %d diverges\nyacc: %s\nrd  : %s", id, i, yErrs[i], rErrs[i])
				}
			}
		}
	}
}

// TestRDParseHintInStatements runs every corpus entry through full
// statement parses so the in-process rd-vs-yacc differential harness
// (rd_differential.go) compares hint ASTs and hint warnings inside real
// statements too.
func TestRDParseHintInStatements(t *testing.T) {
	p := New()
	for _, frag := range rdHintCorpus {
		sql := "select /*+ " + frag + " */ * from t1, t2 where t1.a = t2.a;"
		// Parse errors (e.g. from hint fragments that break the outer
		// comment or statement) are fine: the differential harness panics
		// on any rd/yacc divergence, which is what this test is for.
		_, _, _ = p.Parse(sql, "", "")

		sql = "update /*+ " + frag + " */ t1 set a = 1;"
		_, _, _ = p.Parse(sql, "", "")
	}
}

// TestRDParseHintRandomDifferential compares the two hint parsers on
// pseudo-random token sequences (deterministic seed) drawn from the hint
// grammar's vocabulary, covering token orders the curated corpus misses.
func TestRDParseHintRandomDifferential(t *testing.T) {
	vocab := []string{
		"USE_INDEX", "HASH_JOIN", "unknown_hint", "QB_NAME", "SET_VAR",
		"MEMORY_QUOTA", "LEADING", "READ_FROM_STORAGE", "SEMIJOIN",
		"TIME_RANGE", "MAX_EXECUTION_TIME", "NTH_PLAN", "WRITE_SLOW_LOG",
		"RESOURCE_GROUP", "JOIN_ORDER", "BKA", "MRR", "HASH_AGG",
		"QUERY_TYPE", "USE_TOJA", "JOIN_FIXED_ORDER", "TIDB_INLJ",
		"TIKV", "TIFLASH", "OLAP", "OLTP", "TRUE", "FALSE", "MB", "GB",
		"PARTITION", "FIRSTMATCH", "DUPSWEEDOUT",
		"t1", "d", "x", "y", "`q q`",
		"(", ")", ",", ".", "@qb", "=", "[", "]", "+", "-",
		"1", "1000", "18446744073709551615", "9223372036854775808",
		"9223372036854775809", "'str'", "\"dq\"", "1.5", "\n",
	}
	// Simple deterministic PRNG (xorshift) so failures are reproducible.
	state := uint64(0x9E3779B97F4A7C15)
	next := func(n int) int {
		state ^= state << 13
		state ^= state >> 7
		state ^= state << 17
		return int(state % uint64(n))
	}
	for iter := 0; iter < 30000; iter++ {
		n := 1 + next(9)
		frag := ""
		for i := 0; i < n; i++ {
			if i > 0 {
				frag += " "
			}
			frag += vocab[next(len(vocab))]
		}
		input := "/*+" + frag + "*/"
		mode := mysql.SQLMode(0)
		if next(2) == 1 {
			mode = mysql.ModeANSIQuotes
		}
		pos := Pos{Line: 1}

		yHints, yErrs := newHintParser().parse(input, mode, pos)
		rHints, rErrs := rdParseHint(input, mode, pos)

		id := fmt.Sprintf("iter=%d input=%q mode=%v", iter, frag, mode)
		if !reflect.DeepEqual(yHints, rHints) {
			t.Fatalf("%s:\nhints diverge\nyacc: %#v\nrd  : %#v", id, yHints, rHints)
		}
		if len(yErrs) != len(rErrs) {
			t.Fatalf("%s:\nerror count diverges\nyacc: %q\nrd  : %q", id, errStrings(yErrs), errStrings(rErrs))
		}
		for i := range yErrs {
			if yErrs[i].Error() != rErrs[i].Error() {
				t.Fatalf("%s:\nerror %d diverges\nyacc: %s\nrd  : %s", id, i, yErrs[i], rErrs[i])
			}
		}
	}
}

func errStrings(errs []error) []string {
	out := make([]string, len(errs))
	for i, e := range errs {
		out[i] = e.Error()
	}
	return out
}
