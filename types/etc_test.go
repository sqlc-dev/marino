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

package types

import (
	"testing"

	"github.com/sqlc-dev/marino/mysql"

	"reflect"
)

func TestStrToType(t *testing.T) {
	for tp, str := range type2Str {
		a := StrToType(str)
		if !reflect.DeepEqual(tp, a) {
			t.Fatalf("got %v, want %v", a, tp)
		}
	}

	tp := StrToType("blob")
	if !reflect.DeepEqual(tp, mysql.TypeBlob) {
		t.Fatalf("got %v, want %v", mysql.TypeBlob, tp)
	}

	tp = StrToType("binary")
	if !reflect.DeepEqual(tp, mysql.TypeString) {
		t.Fatalf("got %v, want %v", mysql.TypeString, tp)
	}
}
