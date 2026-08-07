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

// Syntax-error construction for the RD parser, reproducing the goyacc
// path byte-for-byte. The generated parser discards goyacc's verbose
// message ("ignore goyacc error message") and appends Scanner.Errorf(""):
//
//	line %d column %d near "%s"%s %s
//
// where line/column are the reader position after scanning the offending
// lookahead token, the near-text is the rest of the input from that
// token's start (capped at 2048 bytes), the first trailing %s is the
// message (empty for plain syntax errors, immediately after the closing
// quote), and the second is the total-length suffix. There are no error-recovery
// productions in parser.y, so every parse reports at most one syntax
// error and aborts.

import (
	"fmt"
	"strconv"
)

// buildSyntaxError formats the goyacc-identical syntax error for the
// offending token t.
func (r *rdParser) buildSyntaxError(t *rdToken) error {
	val := r.src[t.offset:]
	var lenStr string
	if len(val) > 2048 {
		lenStr = "(total length " + strconv.Itoa(len(val)) + ")"
		val = val[:2048]
	}
	return fmt.Errorf("line %d column %d near \"%s\"%s %s",
		t.endLine, t.endCol, val, "", lenStr)
}

// actionErrorf reproduces yylex.Errorf inside a grammar action: by the
// time an action runs, goyacc has always lexed its one-token lookahead,
// so the reported position is the lookahead token's — robust here even
// if the RD parser peeked further ahead.
func (r *rdParser) actionErrorf(format string, a ...interface{}) error {
	t := r.cur()
	str := fmt.Sprintf(format, a...)
	val := r.src[t.offset:]
	var lenStr string
	if len(val) > 2048 {
		lenStr = "(total length " + strconv.Itoa(len(val)) + ")"
		val = val[:2048]
	}
	return fmt.Errorf("line %d column %d near \"%s\"%s %s",
		t.endLine, t.endCol, val, str, lenStr)
}
