# marino development guide

marino is a hand-written recursive-descent parser for the MySQL dialect
(TiDB-flavored). It was ported from a goyacc grammar; PLAN.md records the
architecture of that rewrite and the decisions that survive it.

## Layout

- `parser/rd_parser.go` — parser core: token window over the streaming
  lexer, dispatch (`rdRegister`, one owner per leading token), statement
  loop with the historical statement-text bookkeeping.
- `parser/parse_*.go` — one file per statement family. Every nontrivial
  parse function carries a comment naming the grammar production(s) it
  implements (the grammar reference is the parser.y of the last goyacc
  commit; see PLAN.md).
- `parser/parse_expr.go` / `parse_func.go` — the expression precedence
  ladder and atoms.
- `parser/token_kinds.go`, `hint_token_kinds.go` — token constants and
  lexer value types (snapshotted from the generated parser at removal;
  hand-maintained now). NOTE: `any`, `recover`, `dump`, `format` shadow
  Go builtins/imports inside package parser.
- `parser/keyword_classes.go` + `keywords.go` + `misc.go` tokenMap — the
  keyword tables, kept in exact agreement by `TestKeywordConsistent`.
- `parser/parse_hint.go` — the optimizer-hint sub-parser.
- `parser/rd_errors.go` — error construction; messages reproduce the
  goyacc format byte-for-byte, validated by `TestRDErrorFidelity`
  against `parser/testdata/errors.json` (goldens recorded from the
  goyacc parser before its removal).

## Rules

- The `ast` package is frozen, and the public `parser` API is stable.
- Error messages are part of the contract: `line N column M near "..."`
  built from the offending token's recorded position; action errors use
  `r.actionErrorf` (reduce-time lookahead position). Failures during
  speculation report the farthest position reached (`farthestFail`).
- One deliberate deviation from goyacc, documented in PLAN.md:
  `OriginTextPosition` is the deterministic production-start offset.
- Always run parser tests with a timeout: they are fast, and a hang
  means an infinite loop in the parser:

```sh
go test ./... -count=1 -timeout 120s
```
