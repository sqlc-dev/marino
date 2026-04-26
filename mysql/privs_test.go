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

package mysql

import (
	"testing"

	"reflect"
)

func TestPrivString(t *testing.T) {
	for i := 0; ; i++ {
		p := PrivilegeType(1 << i)
		if p > AllPriv {
			break
		}
		if reflect.DeepEqual("", p.String()) {
			t.Fatalf("expected values to differ, both are %v", p.String())
		}
	}
}

func TestPrivColumn(t *testing.T) {
	for _, p := range AllGlobalPrivs {
		if len(p.ColumnString()) == 0 {
			t.Fatalf("%s", p)
		}
		np, ok := NewPrivFromColumn(p.ColumnString())
		if !(ok) {
			t.Fatalf("%s", p)
		}
		if !reflect.DeepEqual(p, np) {
			t.Fatalf("got %v, want %v", np, p)
		}
	}
	for _, p := range StaticGlobalOnlyPrivs {
		if len(p.ColumnString()) == 0 {
			t.Fatalf("%s", p)
		}
		np, ok := NewPrivFromColumn(p.ColumnString())
		if !(ok) {
			t.Fatalf("%s", p)
		}
		if !reflect.DeepEqual(p, np) {
			t.Fatalf("got %v, want %v", np, p)
		}
	}
	for _, p := range AllDBPrivs {
		if len(p.ColumnString()) == 0 {
			t.Fatalf("%s", p)
		}
		np, ok := NewPrivFromColumn(p.ColumnString())
		if !(ok) {
			t.Fatalf("%s", p)
		}
		if !reflect.DeepEqual(p, np) {
			t.Fatalf("got %v, want %v", np, p)
		}
	}
}

func TestPrivSetString(t *testing.T) {
	for _, p := range AllTablePrivs {
		if len(p.SetString()) == 0 {
			t.Fatalf("%s", p)
		}
		np, ok := NewPrivFromSetEnum(p.SetString())
		if !(ok) {
			t.Fatalf("%s", p)
		}
		if !reflect.DeepEqual(p, np) {
			t.Fatalf("got %v, want %v", np, p)
		}
	}
	for _, p := range AllColumnPrivs {
		if len(p.SetString()) == 0 {
			t.Fatalf("%s", p)
		}
		np, ok := NewPrivFromSetEnum(p.SetString())
		if !(ok) {
			t.Fatalf("%s", p)
		}
		if !reflect.DeepEqual(p, np) {
			t.Fatalf("got %v, want %v", np, p)
		}
	}
}

func TestPrivsHas(t *testing.T) {
	// it is a simple helper, does not handle all&dynamic privs
	privs := Privileges{AllPriv}
	if !(privs.Has(AllPriv)) {
		t.Fatal("expected true")
	}
	if privs.Has(InsertPriv) {
		t.Fatal("expected false")
	}

	// multiple privs
	privs = Privileges{InsertPriv, SelectPriv}
	if !(privs.Has(SelectPriv)) {
		t.Fatal("expected true")
	}
	if !(privs.Has(InsertPriv)) {
		t.Fatal("expected true")
	}
	if privs.Has(DropPriv) {
		t.Fatal("expected false")
	}
}

func TestPrivAllConsistency(t *testing.T) {
	// AllPriv in mysql.user columns.
	for priv := CreatePriv; priv != AllPriv; priv = priv << 1 {
		_, ok := Priv2UserCol[priv]
		if !(ok) {
			t.Fatalf("priv fail %d", priv)
		}
	}

	if !reflect.DeepEqual(len(AllGlobalPrivs)+1, len(Priv2UserCol)) {
		t.Fatalf("got %v, want %v", len(Priv2UserCol), len(AllGlobalPrivs)+1)
	}

	// USAGE privilege doesn't have a column in Priv2UserCol
	// ALL privilege doesn't have a column in Priv2UserCol
	// so it's +2
	if !reflect.DeepEqual(len(Priv2UserCol)+2, len(Priv2Str)) {
		t.Fatalf("got %v, want %v", len(Priv2Str), len(Priv2UserCol)+2)
	}
}
