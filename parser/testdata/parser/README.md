# File-driven parser tests

Each subdirectory is one test group, run as `TestParserData/<name>` (see
`parser/datadriven_test.go`). A group holds two files whose cases are
separated by lines reading exactly `-- case`:

- `input.sql` — the source SQL for each case. An optional first line
  `-- flags: a, b` configures the parser for the whole group; recognized
  flags are `window_func`, `mariadb`, and `real_as_float`.
- `output.sql` — one segment per input case. A segment is either the
  expected `Restore()` output of the parsed statements (joined by `; `
  for multi-statement sources), or `-- error: <message>` with the exact
  error the parse must fail with. Error messages are part of the parser's
  compatibility contract (see CLAUDE.md), so they are matched verbatim.

Every passing case is also round-tripped: the restored SQL is re-parsed
and the two ASTs must be deep-equal.

To regenerate all `output.sql` goldens from current parser behavior:

```sh
go test ./parser -run TestParserData -update
```

Review the diff before committing — the goldens are the contract.
