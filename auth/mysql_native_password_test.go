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

package auth

import (
	"testing"

	"reflect"
)

func TestEncodePassword(t *testing.T) {
	pwd := "123"
	if !reflect.DeepEqual("*23AE809DDACAF96AF0FD78ED04B6A265E05AA257", EncodePassword(pwd)) {
		t.Fatalf("got %v, want %v", EncodePassword(pwd), "*23AE809DDACAF96AF0FD78ED04B6A265E05AA257")
	}
	if !reflect.DeepEqual(EncodePasswordBytes([]byte(pwd)), EncodePassword(pwd)) {
		t.Fatalf("got %v, want %v", EncodePassword(pwd), EncodePasswordBytes([]byte(pwd)))
	}
}

func TestDecodePassword(t *testing.T) {
	x, err := DecodePassword(EncodePassword("123"))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(Sha1Hash(Sha1Hash([]byte("123"))), x) {
		t.Fatalf("got %v, want %v", x, Sha1Hash(Sha1Hash([]byte("123"))))
	}
}

func TestCheckScramble(t *testing.T) {
	pwd := "abc"
	salt := []byte{85, 92, 45, 22, 58, 79, 107, 6, 122, 125, 58, 80, 12, 90, 103, 32, 90, 10, 74, 82}
	auth := []byte{24, 180, 183, 225, 166, 6, 81, 102, 70, 248, 199, 143, 91, 204, 169, 9, 161, 171, 203, 33}
	encodepwd := EncodePassword(pwd)
	hpwd, err := DecodePassword(encodepwd)
	if err != nil {
		t.Fatal(err)
	}

	res := CheckScrambledPassword(salt, hpwd, auth)
	if !(res) {
		t.Fatal("expected true")
	}

	// Do not panic for invalid input.
	res = CheckScrambledPassword(salt, hpwd, []byte("xxyyzz"))
	if res {
		t.Fatal("expected false")
	}
}
