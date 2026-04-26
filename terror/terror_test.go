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
	"errors"
	"fmt"
	"reflect"
	"testing"
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
	sqlErr := ToSQLError(cause(e).(*Error))
	if !reflect.DeepEqual("Duplicate entry '1' for key 'PRIMARY'", sqlErr.Message) {
		t.Fatalf("got %v, want %v", sqlErr.Message, "Duplicate entry '1' for key 'PRIMARY'")
	}
	if !reflect.DeepEqual(uint16(1062), sqlErr.Code) {
		t.Fatalf("got %v, want %v", sqlErr.Code, uint16(1062))
	}

	err := ErrCritical.GenWithStackByArgs("test")
	if !(ErrCritical.Equal(err)) {
		t.Fatal("expected true")
	}

	if !(ErrCritical.Equal(ErrCritical)) {
		t.Fatal("expected true")
	}
}

func TestErrorEqual(t *testing.T) {
	e1 := errors.New("test error")
	if e1 == nil {
		t.Fatal("expected non-nil")
	}

	if !reflect.DeepEqual(e1, cause(e1)) {
		t.Fatalf("got %v, want %v", cause(e1), e1)
	}

	e4 := errors.New("test error")
	if !(ErrorEqual(e1, e4)) {
		t.Fatal("expected true")
	}
	e5 := fmt.Errorf("test error")
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
