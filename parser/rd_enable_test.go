// Copyright 2026 The sqlc Authors.
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

package parser

import "os"

// The differential harness (rd_differential.go) is on for every test in
// this package — internal and external test files share the test binary —
// so the whole existing suite exercises rd-vs-goyacc equivalence.
// MARINO_NO_DIFF=1 disables it (benchmark comparisons).
func init() {
	rdDifferential = os.Getenv("MARINO_NO_DIFF") == ""
}
