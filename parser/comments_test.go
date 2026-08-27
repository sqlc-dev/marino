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

import (
	"reflect"
	"testing"
)

// TestParserComments pins the comment-recording contract: every skipped
// comment is recorded exactly once with offsets that slice its text back
// out of the source, and constructs the lexer turns into SQL (executable
// comments, optimizer hints) are not comments.
func TestParserComments(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want []string // want[i] == sql[Begin:End] of comment i
	}{
		{
			name: "no comments",
			sql:  "SELECT 1",
			want: nil,
		},
		{
			name: "dash line comment above statement",
			sql:  "-- name: GetOne :one\nSELECT 1",
			want: []string{"-- name: GetOne :one"},
		},
		{
			name: "hash line comment",
			sql:  "# leading\nSELECT 1",
			want: []string{"# leading"},
		},
		{
			name: "hash comment at EOF without newline",
			sql:  "SELECT 1; # trailing",
			want: []string{"# trailing"},
		},
		{
			name: "block comment inside statement",
			sql:  "SELECT /* inline note */ 1",
			want: []string{"/* inline note */"},
		},
		{
			name: "multi-line block comment",
			sql:  "SELECT 1 /* one\n   two */ + 2",
			want: []string{"/* one\n   two */"},
		},
		{
			name: "star-heavy block comment",
			sql:  "SELECT /** starry **/ 1",
			want: []string{"/** starry **/"},
		},
		{
			name: "comments between statements",
			sql:  "SELECT 1; -- after one\n-- before two\nSELECT 2;",
			want: []string{"-- after one", "-- before two"},
		},
		{
			name: "every syntax in one file",
			sql:  "-- dash\n# hash\nSELECT /* block */ 1;",
			want: []string{"-- dash", "# hash", "/* block */"},
		},
		{
			// `AS` makes Lex look ahead for `OF` with a saved-and-restored
			// reader, scanning the comment twice; it must be recorded once.
			name: "comment re-scanned by lookahead",
			sql:  "SELECT 1 AS /* dup */ x",
			want: []string{"/* dup */"},
		},
		{
			name: "double dash without space is an operator",
			sql:  "SELECT 1 --2",
			want: nil,
		},
		{
			name: "executable comment is SQL",
			sql:  "SELECT /*!80000 1 */",
			want: nil,
		},
		{
			name: "optimizer hint is a token",
			sql:  "SELECT /*+ MAX_EXECUTION_TIME(1000) */ 1",
			want: nil,
		},
		{
			name: "misplaced optimizer hint is a comment",
			sql:  "SELECT 1 + /*+ MAX_EXECUTION_TIME(1000) */ 2",
			want: []string{"/*+ MAX_EXECUTION_TIME(1000) */"},
		},
	}
	p := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, _, err := p.Parse(tt.sql, "", ""); err != nil {
				t.Fatalf("Parse(%q): %v", tt.sql, err)
			}
			var got []string
			last := -1
			for _, c := range p.Comments() {
				if c.Begin <= last {
					t.Fatalf("comments out of order: Begin %d after %d", c.Begin, last)
				}
				last = c.Begin
				got = append(got, tt.sql[c.Begin:c.End])
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse(%q) comments = %q, want %q", tt.sql, got, tt.want)
			}
		})
	}
}

// TestParserCommentsReset checks that a parse clears the comments of the
// one before it.
func TestParserCommentsReset(t *testing.T) {
	p := New()
	if _, _, err := p.Parse("SELECT 1 -- one", "", ""); err != nil {
		t.Fatal(err)
	}
	if got := len(p.Comments()); got != 1 {
		t.Fatalf("first parse recorded %d comments, want 1", got)
	}
	if _, _, err := p.Parse("SELECT 2", "", ""); err != nil {
		t.Fatal(err)
	}
	if got := len(p.Comments()); got != 0 {
		t.Fatalf("second parse kept %d stale comments, want 0", got)
	}
}
