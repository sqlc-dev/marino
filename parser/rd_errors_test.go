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

// Error-fidelity corpus: every invalid input harvested from the suite's
// fallback log, with the goyacc parser's exact error and warning strings
// recorded while it was still in-tree. The RD parser's error mode must
// reproduce them byte-for-byte. Regenerate (requires the goyacc parser):
//
//	rm -f /tmp/rd.log
//	MARINO_RD_LOG=/tmp/rd.log go test ./parser/ -count=1
//	MARINO_GEN_ERRCORPUS=/tmp/rd.log go test ./parser/ -run TestGenerateErrorCorpus -count=1

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"
)

// Fields are []byte so invalid-UTF-8 inputs survive the JSON round trip
// (base64), which string fields silently corrupt to U+FFFD.
type errCorpusEntry struct {
	SQL   []byte   `json:"sql"`
	Err   []byte   `json:"err"`
	Warns [][]byte `json:"warns,omitempty"`
}

const errCorpusPath = "testdata/errors.json"

func TestGenerateErrorCorpus(t *testing.T) {
	logPath := os.Getenv("MARINO_GEN_ERRCORPUS")
	if logPath == "" {
		t.Skip("set MARINO_GEN_ERRCORPUS to the fallback log to regenerate")
	}
	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	var sqls []string
	for _, rec := range strings.Split(string(data), "\x00\n") {
		parts := strings.SplitN(rec, "\x00", 2)
		if len(parts) != 2 || seen[parts[1]] {
			continue
		}
		seen[parts[1]] = true
		sqls = append(sqls, parts[1])
	}
	sort.Strings(sqls)
	var entries []errCorpusEntry
	for _, sql := range sqls {
		p := New()
		_, warns, perr := p.yaccParseSQL(sql, CharsetConnection(""), CollationConnection(""))
		if perr == nil {
			// Accepted under the default configuration (the log records
			// no config); the differential harness owns accepted inputs.
			continue
		}
		e := errCorpusEntry{SQL: []byte(sql), Err: []byte(perr.Error())}
		for _, w := range warns {
			e.Warns = append(e.Warns, []byte(w.Error()))
		}
		entries = append(entries, e)
	}
	if err := os.MkdirAll("testdata", 0o755); err != nil {
		t.Fatal(err)
	}
	f, err := os.Create(errCorpusPath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for _, e := range entries {
		b, _ := json.Marshal(e)
		w.Write(b)
		w.WriteByte('\n')
	}
	w.Flush()
	t.Logf("wrote %d entries", len(entries))
}

// TestRDErrorFidelity checks the RD parser's error mode against the
// recorded corpus.
func TestRDErrorFidelity(t *testing.T) {
	data, err := os.ReadFile(errCorpusPath)
	if err != nil {
		t.Skip("no error corpus present")
	}
	mismatches := 0
	sc := bufio.NewScanner(strings.NewReader(string(data)))
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	for sc.Scan() {
		var e errCorpusEntry
		if err := json.Unmarshal(sc.Bytes(), &e); err != nil {
			t.Fatal(err)
		}
		sql := string(e.SQL)
		p := New()
		resetParams(p)
		p.lexer.reset(sql)
		CharsetConnection("").ApplyOn(p)
		CollationConnection("").ApplyOn(p)
		p.src = sql
		_, warns, perr, handled := p.parseRDMode(sql, true)
		report := func(format string, args ...interface{}) {
			mismatches++
			if mismatches <= 20 {
				t.Errorf("SQL: %s\n  %s", sql, fmt.Sprintf(format, args...))
			}
		}
		if !handled {
			report("error mode did not handle the input")
			continue
		}
		if perr == nil {
			report("rd accepted; yacc error: %s", e.Err)
			continue
		}
		if perr.Error() != string(e.Err) {
			report("error mismatch\n  rd  : %s\n  yacc: %s", perr.Error(), e.Err)
			continue
		}
		var gotWarns []string
		for _, w := range warns {
			gotWarns = append(gotWarns, w.Error())
		}
		if len(gotWarns) != len(e.Warns) {
			report("warn count mismatch rd=%v yacc=%v", gotWarns, e.Warns)
			continue
		}
		for i := range gotWarns {
			if gotWarns[i] != string(e.Warns[i]) {
				report("warn mismatch\n  rd  : %s\n  yacc: %s", gotWarns[i], e.Warns[i])
				break
			}
		}
	}
	if mismatches > 20 {
		t.Errorf("... and %d more mismatches", mismatches-20)
	}
	if mismatches > 0 {
		t.Logf("total mismatches: %d", mismatches)
	}
}
