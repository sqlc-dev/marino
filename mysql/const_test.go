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
	"strings"
)

func TestSQLMode(t *testing.T) {
	// ref https://dev.mysql.com/doc/internals/en/query-event.html#q-sql-mode-code,
	hardCode := []struct {
		code  SQLMode
		value int
	}{{
		ModeRealAsFloat, 0x00000001,
	}, {
		ModePipesAsConcat, 0x00000002,
	}, {
		ModeANSIQuotes, 0x00000004,
	}, {
		ModeIgnoreSpace, 0x00000008,
	}, {
		ModeNotUsed, 0x00000010,
	}, {
		ModeOnlyFullGroupBy, 0x00000020,
	}, {
		ModeNoUnsignedSubtraction, 0x00000040,
	}, {
		ModeNoDirInCreate, 0x00000080,
	}, {
		ModePostgreSQL, 0x00000100,
	}, {
		ModeOracle, 0x00000200,
	}, {
		ModeMsSQL, 0x00000400,
	}, {
		ModeDb2, 0x00000800,
	}, {
		ModeMaxdb, 0x00001000,
	}, {
		ModeNoKeyOptions, 0x00002000,
	}, {
		ModeNoTableOptions, 0x00004000,
	}, {
		ModeNoFieldOptions, 0x00008000,
	}, {
		ModeMySQL323, 0x00010000,
	}, {
		ModeMySQL40, 0x00020000,
	}, {
		ModeANSI, 0x00040000,
	}, {
		ModeNoAutoValueOnZero, 0x00080000,
	}, {
		ModeNoBackslashEscapes, 0x00100000,
	}, {
		ModeStrictTransTables, 0x00200000,
	}, {
		ModeStrictAllTables, 0x00400000,
	}, {
		ModeNoZeroInDate, 0x00800000,
	}, {
		ModeNoZeroDate, 0x01000000,
	}, {
		ModeInvalidDates, 0x02000000,
	}, {
		ModeErrorForDivisionByZero, 0x04000000,
	}, {
		ModeTraditional, 0x08000000,
	}, {
		ModeNoAutoCreateUser, 0x10000000,
	}, {
		ModeHighNotPrecedence, 0x20000000,
	}, {
		ModeNoEngineSubstitution, 0x40000000,
	}, {
		ModePadCharToFullLength, 0x80000000,
	}}

	for _, ca := range hardCode {
		if !reflect.DeepEqual(ca.value, int(ca.code)) {
			t.Fatalf("got %v, want %v", int(ca.code), ca.value)
		}
	}
}

func TestVersionSeparator(t *testing.T) {
	// DO NOT change the value of VersionSeparator.
	if !reflect.DeepEqual("-TiDB-", VersionSeparator) {
		t.Fatalf("got %v, want %v", VersionSeparator, "-TiDB-")
	}
}

func TestBuildTiDBXReleaseVersion(t *testing.T) {
	tidbXVersion, err := BuildTiDBXReleaseVersion("v26.3.0")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual("CLOUD.202603.0", tidbXVersion) {
		t.Fatalf("got %v, want %v", tidbXVersion, "CLOUD.202603.0")
	}

	tidbXVersion, err = BuildTiDBXReleaseVersion("v26.3.0-xxx")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual("CLOUD.202603.0-xxx", tidbXVersion) {
		t.Fatalf("got %v, want %v", tidbXVersion, "CLOUD.202603.0-xxx")
	}

	serverVersion, err := BuildTiDBXServerVersion("v26.3.0")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual("8.0.11-TiDB-CLOUD.202603.0", serverVersion) {
		t.Fatalf("got %v, want %v", serverVersion, "8.0.11-TiDB-CLOUD.202603.0")
	}

	serverVersion, err = BuildTiDBXServerVersion("v26.3.0-xxx")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual("8.0.11-TiDB-CLOUD.202603.0-xxx", serverVersion) {
		t.Fatalf("got %v, want %v", serverVersion, "8.0.11-TiDB-CLOUD.202603.0-xxx")
	}

	for _, ver := range []string{"26.1.1", "v26xxxx", "v24.1.1", "v26.0.1", "v26.13.1"} {
		_, err = BuildTiDBXReleaseVersion(ver)
		if err == nil || !strings.Contains(err.Error(), "invalid TiDB release version") {
			t.Fatalf("expected error containing %q, got %v", "invalid TiDB release version", err)
		}
	}
}

func TestNormalizeTiDBReleaseVersionForNextGen(t *testing.T) {
	if !reflect.DeepEqual(tidbXPlaceholderReleaseVersion, NormalizeTiDBReleaseVersionForNextGen(legacyTiDBReleaseVersionPlaceholder)) {
		t.Fatalf("got %v, want %v", NormalizeTiDBReleaseVersionForNextGen(legacyTiDBReleaseVersionPlaceholder), tidbXPlaceholderReleaseVersion)
	}
	if !reflect.DeepEqual("v26.3.0", NormalizeTiDBReleaseVersionForNextGen("v26.3.0")) {
		t.Fatalf("got %v, want %v", NormalizeTiDBReleaseVersionForNextGen("v26.3.0"), "v26.3.0")
	}
}
