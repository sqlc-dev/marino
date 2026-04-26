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

package util

import (
	"fmt"
	"testing"

	"reflect"
)

func TestUnescapeChar(t *testing.T) {
	tests := []struct {
		input byte
		want  []byte
	}{
		// Standard single-byte escapes
		{'n', []byte{'\n'}},
		{'0', []byte{0}},
		{'b', []byte{8}},
		{'Z', []byte{26}},
		{'r', []byte{'\r'}},
		{'t', []byte{'\t'}},

		// Preserve both backslash and character
		{'%', []byte{'\\', '%'}},
		{'_', []byte{'\\', '_'}},

		// Self-escaping characters (backslash removed)
		{'\\', []byte{'\\'}},
		{'\'', []byte{'\''}},
		{'"', []byte{'"'}},

		// Any other character just returns itself (backslash removed)
		{'a', []byte{'a'}},
		{'z', []byte{'z'}},
		{'1', []byte{'1'}},
		{' ', []byte{' '}},
	}
	for _, tt := range tests {
		got := UnescapeChar(tt.input)
		if !reflect.DeepEqual(tt.want, got) {
			t.Fatalf("%s: got %v, want %v", fmt.Sprintf("UnescapeChar(%q)", tt.input), got, tt.want)
		}
	}
}
