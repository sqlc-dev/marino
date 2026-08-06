# marino development guide

marino is a hand-written recursive-descent parser for the MySQL dialect
(TiDB-flavored). Read PLAN.md first — it defines the architecture of the
goyacc→recursive-descent rewrite, what is frozen (the `ast` package and
the public `parser` API), and the milestones.

## During the migration (goyacc still in-tree)

The RD parser lives in `parser/parse_*.go` + `parser/rd_*.go`; the goyacc
parser (`parser/parser.y` → generated `parser/parser.go`) is the fallback
for anything unimplemented and the *oracle* for everything implemented:

- `ParseSQL` tries the RD parser first; `r.unsupported(...)` unwinds to a
  goyacc fallback, so the tree stays green at every stage.
- Under `go test`, every input the RD parser handles is re-parsed by
  goyacc and compared field-for-field (`rd_differential.go`, rendered by
  `internal/dump`). A divergence panics the test with a `rd >`/`yacc>`
  diff. The yacc side is always right — fix the transcription.
- Porting conventions (API, offsets, error actions):
  `internal/rdconventions.md`.

### The loop

```sh
rm -f /tmp/rd.log
MARINO_RD_LOG=/tmp/rd.log go test ./parser/ -count=1 -timeout 500s
go run ./cmd/next-test /tmp/rd.log     # ranked remaining fallbacks
# ... port the next family in parser/parse_*.go ...
go test ./... -count=1                 # everything green before commit
```

Always run the parser tests with a timeout: they are fast, and a hang
means an infinite loop in the RD parser.

## Rules

- **Never** edit `parser/parser.y`, the generated `parser/parser.go` /
  `parser/hintparser.go`, or any `_test.go` of the existing suite. The
  suite is the acceptance gate and must pass unmodified.
- The `ast` package is frozen: the RD parser produces byte-identical
  trees (the differential enforces it).
- Exception, documented in PLAN.md and `internal/dump/dump.go`:
  `OriginTextPosition` is stamped as the deterministic production-start
  offset; the goyacc value depended on parser-stack layout and is not
  reproduced.
- Statement families register their leading token via `rdRegister` in
  their own file's `init` (see `rd_parser.go`), one owner per token.
- `MARINO_NO_DIFF=1` disables the in-process differential (benchmarks).
