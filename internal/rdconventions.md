# RD parser porting conventions

(Working notes for the goyacc→recursive-descent rewrite; deleted at the
end of the migration. See PLAN.md for the architecture.)

Every parse function ports specific productions of `parser/parser.y`
(grammar + action) into hand-written Go in package `parser`. The goyacc
parser is still in-tree: anything the RD parser cannot handle falls back
to it, and under `go test` every input the RD parser does handle is
re-parsed by goyacc and compared field-for-field (rd_differential.go).
A divergence panics the test with a dump diff labeled `rd >` / `yacc>`.

## Hard rules

1. Transcribe the action logic EXACTLY — same AST fields, same values,
   same order-sensitive side effects. Do not "improve" anything.
2. Every nontrivial parse function carries a comment naming the
   production(s) it implements.
3. `r.unsupported("thing")` for grammar you are not porting (falls back
   to goyacc; never guess).
4. Never edit parser.y, parser.go, lexer.go, misc.go, yy_parser.go, or
   any _test.go file.
5. Build + run the relevant tests after every family; the FULL parser
   suite (`go test ./parser/ -count=1`) must be green before you stop.

## The rdParser API (rd_parser.go, parse_expr.go, parse_func.go, ...)

Cursor: `r.tok()` current token id (0=EOF), `r.la(k)` lookahead id,
`r.cur()` current *rdToken (`.lit` source spelling, `.item` lexer value,
`.offset` byte offset), `r.advance()`, `r.accept(tok) bool`,
`r.expect(tok) rdToken`, `r.syntaxError()` (panics; during migration =
fallback to goyacc).

Speculation: `r.try(func(){...}) bool` — rewinds on rdSyntaxError inside.

Errors from grammar actions:
- `yylex.AppendError(err); return 1` → `r.actionError(err)` (aborts).
- `yylex.AppendError(err)` without return 1 → `r.sc.AppendError(err)`.
- warnings-as-errors pattern (`parser.lastErrorAsWarn()`) →
  `r.sc.AppendError(w); r.sc.lastErrorAsWarn()`.

Token constants come from the generated parser.go (`selectKwd`, `intLit`,
`identifier`, single chars as `int('(')`...). The mapping literal→name is
greppable in parser.y's %token block. NOTE: the constants `any`, `recover`,
`dump`, `format` shadow Go builtins/imports in package parser.

Common helpers (already written — do NOT redefine):
`parseIdentifier`, `parseStringName`, `parseTableName`,
`parseTableNameOptWild`, `parseColumnName(List)`, `parseExpression`,
`parseExpressionList(Opt)`, `parseSimpleExpr`, `parseSubSelect`,
`parseSelectCore(With)`, `parseTableRefs`, `parseCharsetName`,
`parseCollationName`, `parseFieldLen`, `parseOptFieldLen`,
`parseLengthNum`, `parseInt64Num`, `parseSignedNum`, `parseCastType`,
`parseFloatOpt`, `parseOptBinary`, `parseCharsetKw`/`isCharsetKw`,
`parseTimeUnit`, `parseOrderBy`, `parseByList`, `parseLimitClause`,
`parseSelectStmtLimit`, `parseTextString`, `parseSignedLiteral`,
`parseLiteralExpr`, `parseAssignmentList`, `parseIfExists`/`parseIfNotExists`
(if present), `isIdentifierTok`, `getUint64FromNUM`, `getInt64FromNUM`.

## Identifiers vs keywords

`Identifier: identifier | UnReservedKeyword | NotKeywordToken |
TiDBKeyword` — use `parseIdentifier()` / `isIdentifierTok(tok)`. A
reserved keyword is NEVER an identifier. When a production offers both a
keyword-specific alternative and an Identifier path (e.g. an unreserved
keyword starting a clause), guard with lookahead the way the sibling code
does and note the pseudo-token/conflict comment from parser.y.

## Offsets and node text

- Expression-producing (`%type <expr>`) productions stamp
  `r.setOrigin(node, startOffset)` where startOffset is the offset of the
  production's FIRST token (grab `start := r.cur().offset` before
  consuming). Item-typed productions do NOT stamp.
- `parser.startOffset(&yyS[...])` in actions = that symbol's first-token
  offset. `parser.endOffset(&yyS[...])` = `r.endOffsetAt(offset)` (trims
  whitespace backwards). `parser.yylval.offset` at reduce time = the
  offset of the CURRENT (unconsumed lookahead) token = `r.cur().offset`
  after parsing the production body.
- `parser.setNodeText(n, text)` → `r.p.setNodeText(n, text)`;
  `parser.src` → `r.src`.

## Values

- `$1` of a token = `tok.lit` (source spelling; keyword tokens carry the
  original case). `$1` of intLit/decLit/floatLit/hexLit/bitLit =
  `tok.item` (already converted by the lexer).
- `NUM` in the grammar is intLit; `getUint64FromNUM(tok.item)`.
- `ast.NewCIStr($1)` etc. transcribe directly.

## Statement wiring

parseStatement (rd_parser.go) dispatches on the leading token. Add your
statement's case there (keep the switch tidy; one line per family
calling a parse function in your file).
