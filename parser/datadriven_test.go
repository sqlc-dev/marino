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
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/sqlc-dev/marino/format"
	"github.com/sqlc-dev/marino/mysql"
	"github.com/sqlc-dev/marino/parser"
)

// TestParserData runs the file-driven parser test suite under
// testdata/parser. Each subdirectory is one test group holding two files:
//
//	input.sql  — the source SQL cases, separated by lines reading "-- case".
//	             An optional first line "-- flags: a, b" configures the
//	             parser for the whole group (window_func, mariadb,
//	             real_as_float).
//	output.sql — one segment per input case, separated the same way. A
//	             segment is either the expected Restore() output of the
//	             parsed statements (joined by "; " for multi-statement
//	             sources), or "-- error: <message>" with the exact error
//	             the parse must fail with.
//
// Every passing case is additionally round-tripped: the restored SQL is
// re-parsed and the two ASTs must be deep-equal (after CleanNodeText).
//
// Run with -update to regenerate every output.sql from current parser
// behavior:
//
//	go test ./parser -run TestParserData -update
var updateParserData = flag.Bool("update", false, "regenerate testdata/parser output.sql goldens")

const (
	parserDataDir = "testdata/parser"
	caseDelimiter = "\n-- case\n"
	errorMarker   = "-- error: "
	flagsMarker   = "-- flags:"
)

func TestParserData(t *testing.T) {
	entries, err := os.ReadDir(parserDataDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		t.Run(e.Name(), func(t *testing.T) {
			runParserDataGroup(t, filepath.Join(parserDataDir, e.Name()))
		})
	}
}

type parserDataFlags struct {
	windowFunc  bool
	mariaDB     bool
	realAsFloat bool
}

func runParserDataGroup(t *testing.T, dir string) {
	inputPath := filepath.Join(dir, "input.sql")
	outputPath := filepath.Join(dir, "output.sql")

	inBuf, err := os.ReadFile(inputPath)
	if err != nil {
		t.Fatal(err)
	}
	flags, srcs := decodeParserDataInput(t, string(inBuf))

	var wants []string
	if !*updateParserData {
		outBuf, err := os.ReadFile(outputPath)
		if err != nil {
			t.Fatal(err)
		}
		wants = decodeParserDataSegments(string(outBuf))
		if len(wants) != len(srcs) {
			t.Fatalf("%s: %d input cases but %d output segments", dir, len(srcs), len(wants))
		}
	}

	// One parser is reused across the whole group, mirroring the historical
	// RunTest helper; restore round-trips use fresh parsers below.
	p := newParserDataParser(flags)
	var gots []string
	for i, src := range srcs {
		comment := fmt.Sprintf("%s case %d: source %v", dir, i, src)
		_, _, err := p.Parse(src, "", "")

		var got string
		if err != nil {
			got = errorMarker + err.Error()
		} else {
			got = restoreRoundTrip(t, flags, src, comment)
		}
		gots = append(gots, got)

		if *updateParserData {
			continue
		}
		want := wants[i]
		if errMsg, isErr := strings.CutPrefix(want, errorMarker); isErr {
			if err == nil {
				t.Fatalf("%v: expected error %q, but parse succeeded", comment, errMsg)
			}
			if err.Error() != errMsg {
				t.Fatalf("%v: error mismatch\ngot:  %v\nwant: %v", comment, err.Error(), errMsg)
			}
			continue
		}
		if err != nil {
			t.Fatalf("%v: %v", comment, err)
		}
		if got != want {
			t.Fatalf("%v: restore mismatch\ngot:  %v\nwant: %v", comment, got, want)
		}
	}

	if *updateParserData {
		if err := os.WriteFile(outputPath, []byte(strings.Join(gots, caseDelimiter)+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func newParserDataParser(flags parserDataFlags) *parser.Parser {
	p := parser.New()
	p.EnableWindowFunc(flags.windowFunc)
	p.SetMariaDB(flags.mariaDB)
	if flags.realAsFloat {
		p.SetSQLMode(mysql.ModeRealAsFloat)
	}
	return p
}

// restoreRoundTrip re-parses src with a fresh parser, restores every
// statement, verifies each restored statement parses back to a deep-equal
// AST, and returns the restored statements joined by "; ".
func restoreRoundTrip(t *testing.T, flags parserDataFlags, src, comment string) string {
	p := newParserDataParser(flags)
	stmts, _, err := p.Parse(src, "", "")
	if err != nil {
		t.Fatalf("%v: %v", comment, err)
	}
	var sb strings.Builder
	restoreSQLs := ""
	for _, stmt := range stmts {
		sb.Reset()
		if err := stmt.Restore(format.NewRestoreCtx(format.DefaultRestoreFlags, &sb)); err != nil {
			t.Fatalf("%v: %v", comment, err)
		}
		restoreSQL := sb.String()
		restoreStmt, err := p.ParseOneStmt(restoreSQL, "", "")
		if err != nil {
			t.Fatalf("%v: restore %v: %v", comment, restoreSQL, err)
		}
		CleanNodeText(stmt)
		CleanNodeText(restoreStmt)
		if !reflect.DeepEqual(stmt, restoreStmt) {
			t.Fatalf("%v: restore %v: got %v, want %v", comment, restoreSQL, restoreStmt, stmt)
		}
		if restoreSQLs != "" {
			restoreSQLs += "; "
		}
		restoreSQLs += restoreSQL
	}
	return restoreSQLs
}

func decodeParserDataInput(t *testing.T, content string) (parserDataFlags, []string) {
	var flags parserDataFlags
	if strings.HasPrefix(content, flagsMarker) {
		nl := strings.IndexByte(content, '\n')
		if nl < 0 {
			t.Fatalf("input.sql is only a flags header")
		}
		for _, f := range strings.Split(strings.TrimPrefix(content[:nl], flagsMarker), ",") {
			switch name := strings.TrimSpace(f); name {
			case "window_func":
				flags.windowFunc = true
			case "mariadb":
				flags.mariaDB = true
			case "real_as_float":
				flags.realAsFloat = true
			default:
				t.Fatalf("unknown flag %q in input.sql header", name)
			}
		}
		content = content[nl+1:]
	}
	return flags, decodeParserDataSegments(content)
}

func decodeParserDataSegments(content string) []string {
	content = strings.TrimSuffix(content, "\n")
	return strings.Split(content, caseDelimiter)
}
