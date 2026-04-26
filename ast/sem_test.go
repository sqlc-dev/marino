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

package ast

import (
	"testing"

	"reflect"
)

func TestShowCommand(t *testing.T) {
	for i := 1; i < showTpCount; i++ {
		stmt := &ShowStmt{
			Tp: ShowStmtType(i),
		}

		if reflect.DeepEqual(stmt.SEMCommand(), UnknownCommand) {
			t.Fatalf("expected values to differ, both are %v", UnknownCommand)
		}
	}
}

func TestAdminCommand(t *testing.T) {
	for i := 1; i < int(adminTpCount); i++ {
		stmt := &AdminStmt{
			Tp: AdminStmtType(i),
		}

		if reflect.DeepEqual(stmt.SEMCommand(), UnknownCommand) {
			t.Fatalf("expected values to differ, both are %v", UnknownCommand)
		}
	}
}

func TestBRIECommand(t *testing.T) {
	for i := range brieKindCount {
		stmt := &BRIEStmt{
			Kind: i,
		}

		if reflect.DeepEqual(stmt.SEMCommand(), UnknownCommand) {
			t.Fatalf("expected values to differ, both are %v", UnknownCommand)
		}
	}
}
