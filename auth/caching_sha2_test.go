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

package auth

import (
	"encoding/hex"
	"testing"

	"github.com/sqlc-dev/marino/mysql"

	"reflect"
)

var foobarPwdSHA2Hash, _ = hex.DecodeString("24412430303524031A69251C34295C4B35167C7F1E5A7B63091349503974624D34504B5A424679354856336868686F52485A736E4A733368786E427575516C73446469496537")

func TestCheckShaPasswordGood(t *testing.T) {
	pwd := "foobar"
	r, err := CheckHashingPassword(foobarPwdSHA2Hash, pwd, mysql.AuthCachingSha2Password)
	if err != nil {
		t.Fatal(err)
	}
	if !(r) {
		t.Fatal("expected true")
	}
}

func TestCheckShaPasswordBad(t *testing.T) {
	pwd := "not_foobar"
	pwhash, _ := hex.DecodeString("24412430303524031A69251C34295C4B35167C7F1E5A7B63091349503974624D34504B5A424679354856336868686F52485A736E4A733368786E427575516C73446469496537")
	r, err := CheckHashingPassword(pwhash, pwd, mysql.AuthCachingSha2Password)
	if err != nil {
		t.Fatal(err)
	}
	if r {
		t.Fatal("expected false")
	}
}

func TestCheckShaPasswordShort(t *testing.T) {
	pwd := "not_foobar"
	pwhash, _ := hex.DecodeString("aaaaaaaa")
	_, err := CheckHashingPassword(pwhash, pwd, mysql.AuthCachingSha2Password)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCheckShaPasswordDigestTypeIncompatible(t *testing.T) {
	pwd := "not_foobar"
	pwhash, _ := hex.DecodeString("24422430303524031A69251C34295C4B35167C7F1E5A7B63091349503974624D34504B5A424679354856336868686F52485A736E4A733368786E427575516C73446469496537")
	_, err := CheckHashingPassword(pwhash, pwd, mysql.AuthCachingSha2Password)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCheckShaPasswordIterationsInvalid(t *testing.T) {
	pwd := "not_foobar"
	pwhash, _ := hex.DecodeString("24412430304724031A69251C34295C4B35167C7F1E5A7B63091349503974624D34504B5A424679354856336868686F52485A736E4A733368786E427575516C73446469496537")
	_, err := CheckHashingPassword(pwhash, pwd, mysql.AuthCachingSha2Password)
	if err == nil {
		t.Fatal("expected error")
	}
}

// The output from NewHashPassword is not stable as the hash is based on the generated salt.
// This is why CheckHashingPassword is used here.
func TestNewSha2Password(t *testing.T) {
	pwd := "testpwd"
	pwhash := NewHashPassword(pwd, mysql.AuthCachingSha2Password)
	r, err := CheckHashingPassword([]byte(pwhash), pwd, mysql.AuthCachingSha2Password)
	if err != nil {
		t.Fatal(err)
	}
	if !(r) {
		t.Fatal("expected true")
	}

	for r := range pwhash {
		if !(pwhash[r] < uint8(128)) {
			t.Fatalf("expected %v < %v", pwhash[r], uint8(128))
		}
		if reflect.DeepEqual(pwhash[r], 0) {
			t.Fatalf("expected values to differ, both are %v", 0)
		} // NUL
		if reflect.DeepEqual(pwhash[r], 36) {
			t.Fatalf("expected values to differ, both are %v", 36)
		} // '$'
	}
}

func BenchmarkShaPassword(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m, err := CheckHashingPassword(foobarPwdSHA2Hash, "foobar", mysql.AuthCachingSha2Password)
		if err != nil {
			b.Fatalf("expected nil, got %v", err)
		}
		if !(m) {
			b.Fatal("expected true")
		}
	}
}
