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

package ast

import (
	"encoding/json"
	"testing"

	"reflect"
)

func TestT(t *testing.T) {
	abc := NewCIStr("aBC")
	if !reflect.DeepEqual("aBC", abc.O) {
		t.Fatalf("got %v, want %v", abc.O, "aBC")
	}
	if !reflect.DeepEqual("abc", abc.L) {
		t.Fatalf("got %v, want %v", abc.L, "abc")
	}
	if !reflect.DeepEqual("aBC", abc.String()) {
		t.Fatalf("got %v, want %v", abc.String(), "aBC")
	}
}

func TestUnmarshalCIStr(t *testing.T) {
	var ci CIStr

	// Test unmarshal CIStr from a single string.
	str := "aaBB"
	buf, err := json.Marshal(str)
	if err != nil {
		t.Fatal(err)
	}
	if ci.UnmarshalJSON(buf) != nil {
		t.Fatal(ci.UnmarshalJSON(buf))
	}
	if !reflect.DeepEqual(str, ci.O) {
		t.Fatalf("got %v, want %v", ci.O, str)
	}
	if !reflect.DeepEqual("aabb", ci.L) {
		t.Fatalf("got %v, want %v", ci.L, "aabb")
	}

	buf, err = json.Marshal(ci)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(`{"O":"aaBB","L":"aabb"}`, string(buf)) {
		t.Fatalf("got %v, want %v", string(buf), `{"O":"aaBB","L":"aabb"}`)
	}
	if ci.UnmarshalJSON(buf) != nil {
		t.Fatal(ci.UnmarshalJSON(buf))
	}
	if !reflect.DeepEqual(str, ci.O) {
		t.Fatalf("got %v, want %v", ci.O, str)
	}
	if !reflect.DeepEqual("aabb", ci.L) {
		t.Fatalf("got %v, want %v", ci.L, "aabb")
	}
}
