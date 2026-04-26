// Copyright 2021 PingCAP, Inc.
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

package terror

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/pingcap/errors"

	"reflect"
)

func TestErrCode(t *testing.T) {
	if !reflect.DeepEqual(ErrCode(1), CodeMissConnectionID) {
		t.Fatalf("got %v, want %v", CodeMissConnectionID, ErrCode(1))
	}
	if !reflect.DeepEqual(ErrCode(2), CodeResultUndetermined) {
		t.Fatalf("got %v, want %v", CodeResultUndetermined, ErrCode(2))
	}
}

func TestTError(t *testing.T) {
	if len(ClassParser.String()) == 0 {
		t.Fatal("expected non-empty")
	}
	if len(ClassOptimizer.String()) == 0 {
		t.Fatal("expected non-empty")
	}
	if len(ClassKV.String()) == 0 {
		t.Fatal("expected non-empty")
	}
	if len(ClassServer.String()) == 0 {
		t.Fatal("expected non-empty")
	}

	parserErr := ClassParser.New(ErrCode(100), "error 100")
	if len(parserErr.Error()) == 0 {
		t.Fatal("expected non-empty")
	}
	if !(ClassParser.EqualClass(parserErr)) {
		t.Fatal("expected true")
	}
	if ClassParser.NotEqualClass(parserErr) {
		t.Fatal("expected false")
	}

	if ClassOptimizer.EqualClass(parserErr) {
		t.Fatal("expected false")
	}
	optimizerErr := ClassOptimizer.New(ErrCode(2), "abc")
	if ClassOptimizer.EqualClass(errors.New("abc")) {
		t.Fatal("expected false")
	}
	if ClassOptimizer.EqualClass(nil) {
		t.Fatal("expected false")
	}
	if !(optimizerErr.Equal(optimizerErr.GenWithStack("def"))) {
		t.Fatal("expected true")
	}
	if optimizerErr.Equal(nil) {
		t.Fatal("expected false")
	}
	if optimizerErr.Equal(errors.New("abc")) {
		t.Fatal("expected false")
	}

	// Test case for FastGen.
	if !(optimizerErr.Equal(optimizerErr.FastGen("def"))) {
		t.Fatal("expected true")
	}
	if !(optimizerErr.Equal(optimizerErr.FastGen("def: %s", "def"))) {
		t.Fatal("expected true")
	}
	kvErr := ClassKV.New(1062, "key already exist")
	e := kvErr.FastGen("Duplicate entry '%d' for key 'PRIMARY'", 1)
	if !reflect.DeepEqual("[kv:1062]Duplicate entry '1' for key 'PRIMARY'", e.Error()) {
		t.Fatalf("got %v, want %v", e.Error(), "[kv:1062]Duplicate entry '1' for key 'PRIMARY'")
	}
	sqlErr := ToSQLError(errors.Cause(e).(*Error))
	if !reflect.DeepEqual("Duplicate entry '1' for key 'PRIMARY'", sqlErr.Message) {
		t.Fatalf("got %v, want %v", sqlErr.Message, "Duplicate entry '1' for key 'PRIMARY'")
	}
	if !reflect.DeepEqual(uint16(1062), sqlErr.Code) {
		t.Fatalf("got %v, want %v", sqlErr.Code, uint16(1062))
	}

	err := errors.Trace(ErrCritical.GenWithStackByArgs("test"))
	if !(ErrCritical.Equal(err)) {
		t.Fatal("expected true")
	}

	err = errors.Trace(ErrCritical)
	if !(ErrCritical.Equal(err)) {
		t.Fatal("expected true")
	}
}

func TestJson(t *testing.T) {
	prevTErr := errors.Normalize("json test", errors.MySQLErrorCode(int(CodeExecResultIsEmpty)))
	buf, err := json.Marshal(prevTErr)
	if err != nil {
		t.Fatal(err)
	}
	var curTErr errors.Error
	err = json.Unmarshal(buf, &curTErr)
	if err != nil {
		t.Fatal(err)
	}
	isEqual := prevTErr.Equal(&curTErr)
	if !(isEqual) {
		t.Fatal("expected true")
	}
}

var predefinedErr = ClassExecutor.New(ErrCode(123), "predefiend error")

func example() error {
	err := call()
	return errors.Trace(err)
}

func call() error {
	return predefinedErr.GenWithStack("error message:%s", "abc")
}

func TestErrorEqual(t *testing.T) {
	e1 := errors.New("test error")
	if e1 == nil {
		t.Fatal("expected non-nil")
	}

	e2 := errors.Trace(e1)
	if e2 == nil {
		t.Fatal("expected non-nil")
	}

	e3 := errors.Trace(e2)
	if e3 == nil {
		t.Fatal("expected non-nil")
	}

	if !reflect.DeepEqual(e1, errors.Cause(e2)) {
		t.Fatalf("got %v, want %v", errors.Cause(e2), e1)
	}
	if !reflect.DeepEqual(e1, errors.Cause(e3)) {
		t.Fatalf("got %v, want %v", errors.Cause(e3), e1)
	}
	if !reflect.DeepEqual(errors.Cause(e3), errors.Cause(e2)) {
		t.Fatalf("got %v, want %v", errors.Cause(e2), errors.Cause(e3))
	}

	e4 := errors.New("test error")
	if reflect.DeepEqual(e1, errors.Cause(e4)) {
		t.Fatalf("expected values to differ, both are %v", errors.Cause(e4))
	}

	e5 := errors.Errorf("test error")
	if reflect.DeepEqual(e1, errors.Cause(e5)) {
		t.Fatalf("expected values to differ, both are %v", errors.Cause(e5))
	}

	if !(ErrorEqual(e1, e2)) {
		t.Fatal("expected true")
	}
	if !(ErrorEqual(e1, e3)) {
		t.Fatal("expected true")
	}
	if !(ErrorEqual(e1, e4)) {
		t.Fatal("expected true")
	}
	if !(ErrorEqual(e1, e5)) {
		t.Fatal("expected true")
	}

	var e6 error

	if !(ErrorEqual(nil, nil)) {
		t.Fatal("expected true")
	}
	if !(ErrorNotEqual(e1, e6)) {
		t.Fatal("expected true")
	}
	code1 := ErrCode(9001)
	code2 := ErrCode(9002)
	te1 := ClassParser.Synthesize(code1, "abc")
	te3 := ClassKV.New(code1, "abc")
	te4 := ClassKV.New(code2, "abc")
	if ErrorEqual(te1, te3) {
		t.Fatal("expected false")
	}
	if ErrorEqual(te3, te4) {
		t.Fatal("expected false")
	}
}

func TestLog(t *testing.T) {
	err := fmt.Errorf("xxx")
	Log(err)
}

func TestTraceAndLocation(t *testing.T) {
	err := example()
	stack := errors.ErrorStack(err)
	lines := strings.Split(stack, "\n")
	goroot := strings.ReplaceAll(runtime.GOROOT(), string(os.PathSeparator), "/")
	var sysStack = 0
	for _, line := range lines {
		// When you run test case in the bazel. you will find the difference stack. It looks like this:
		//
		// ```go
		// testing.tRunner
		//   GOROOT/src/testing/testing.go:1576
		// runtime.goexit
		//   src/runtime/asm_arm64.s:1172
		// ```
		//
		// but run with ```go test```. It looks like this:
		//
		// ```go
		// testing.tRunner
		//   /Users/pingcap/.gvm/gos/go1.20.1/src/testing/testing.go:1576
		// runtime.goexit
		//	 /Users/pingcap/.gvm/gos/go1.20.1/src/runtime/asm_arm64.s:1172
		// ```
		//
		// So we have to deal with these boundary conditions.
		if strings.Contains(line, goroot) || strings.Contains(line, "src/runtime") {
			sysStack++
		}
	}
	if !reflect.DeepEqual(9, len(lines)-(2*sysStack)) {
		t.Fatalf("%s: got %v, want %v", fmt.Sprintf("stack =\n%s", stack), len(lines)-(2*sysStack), 9)
	}
	var containTerr bool
	for _, v := range lines {
		if strings.Contains(v, "terror_test.go") {
			containTerr = true
			break
		}
	}
	if !(containTerr) {
		t.Fatal("expected true")
	}
}
