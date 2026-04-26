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

package mysql

import (
	"testing"
)

func TestFlags(t *testing.T) {
	if !(HasNotNullFlag(NotNullFlag)) {
		t.Fatal("expected true")
	}
	if !(HasUniKeyFlag(UniqueKeyFlag)) {
		t.Fatal("expected true")
	}
	if !(HasNotNullFlag(NotNullFlag)) {
		t.Fatal("expected true")
	}
	if !(HasNoDefaultValueFlag(NoDefaultValueFlag)) {
		t.Fatal("expected true")
	}
	if !(HasAutoIncrementFlag(AutoIncrementFlag)) {
		t.Fatal("expected true")
	}
	if !(HasUnsignedFlag(UnsignedFlag)) {
		t.Fatal("expected true")
	}
	if !(HasZerofillFlag(ZerofillFlag)) {
		t.Fatal("expected true")
	}
	if !(HasBinaryFlag(BinaryFlag)) {
		t.Fatal("expected true")
	}
	if !(HasPriKeyFlag(PriKeyFlag)) {
		t.Fatal("expected true")
	}
	if !(HasMultipleKeyFlag(MultipleKeyFlag)) {
		t.Fatal("expected true")
	}
	if !(HasTimestampFlag(TimestampFlag)) {
		t.Fatal("expected true")
	}
	if !(HasOnUpdateNowFlag(OnUpdateNowFlag)) {
		t.Fatal("expected true")
	}
}
