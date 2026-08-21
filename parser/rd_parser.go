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

// This file is the core of the hand-written recursive-descent parser
// that replaced the goyacc-generated one (see PLAN.md).
//
// Tokens are lexed on demand into a small sliding window rather than
// eagerly into a slice: parsing must stay O(1) in memory the way the
// goyacc parser's single-lookahead loop is (the existing test suite
// asserts allocation bounds on large statements). Speculative parses
// pin the window with mark/rewind.
//
// Grammar attribution comments in parse functions name the parser.y
// rule(s) they implement.

import (
	"fmt"
	"os"
	"slices"
	"sync"

	"github.com/sqlc-dev/marino/ast"
	"github.com/sqlc-dev/marino/internal/panics"
)

// rdToken is one lexed token. lit and item mirror the ident/item fields
// the Scanner stores into yySymType.
type rdToken struct {
	tok    int
	lit    string
	item   interface{} // not `any`: the token constant `any` shadows the builtin
	offset int
	// endOffset/endLine/endCol are the reader position after scanning this
	// token, used to reproduce goyacc's error positions.
	endOffset int
	endLine   int
	endCol    int
	// hintPos is set for hintComment tokens.
	hintPos Pos
}

// rdParser drives a recursive-descent parse over a lazily-lexed token
// stream. The Scanner it lexes with is a configuration clone of the
// Parser's lexer, so a fallback to goyacc starts from pristine state.
type rdParser struct {
	p   *Parser
	sc  *Scanner
	src string

	// win holds tokens [base, base+len(win)) of the input; i is the
	// absolute index of the current token. done is set once the EOF
	// token (tok 0) has been lexed; the window then ends with it.
	win   []rdToken
	base  int
	i     int
	done  bool
	marks []int // absolute indices pinned by speculative parses

	// c points at the current token's window slot, kept in sync by
	// advance and rewind (the only movers of i) so that cur and tok are
	// call-free field loads. Slot pointers stay valid across window
	// growth (the old backing array keeps the same values) and are
	// recomputed after compaction, which rewrites slots in place.
	c *rdToken

	// stmtStart mirrors Scanner.stmtStartPos bookkeeping for stmtText().
	stmtStart int
	result    []ast.StmtNode

	// farthestFail remembers the deepest syntax failure across
	// speculative attempts: the LALR automaton is a single forward pass,
	// so its reported error is always at the farthest token reached,
	// which a backtracking parser must reconstruct.
	farthestFail *rdSyntaxError
}

// rdUnsupported unwinds the RD parse when it reaches grammar that is not
// implemented yet; the caller falls back to goyacc.
type rdUnsupported struct {
	construct string
	offset    int
	err       error
}

// rdSyntaxError unwinds the RD parse on invalid input, carrying the
// goyacc-identical error built at the offending token. While goyacc is
// in-tree (and outside errorMode) it is treated like rdUnsupported so
// that reported messages keep coming from goyacc.
type rdSyntaxError struct {
	offset int
	err    error
}

// rdLexError unwinds the RD parse when the scanner reports an error while
// lexing on demand; goyacc redoes the parse and reports it identically.
type rdLexError struct{}

// unsupported aborts the RD parse for a construct that is not implemented.
func (r *rdParser) unsupported(construct string) {
	panic(rdUnsupported{construct: construct, offset: r.cur().offset, err: r.buildSyntaxError(r.cur())})
}

// syntaxError aborts the RD parse at the current token.
func (r *rdParser) syntaxError() {
	e := rdSyntaxError{offset: r.cur().offset, err: r.buildSyntaxError(r.cur())}
	if r.farthestFail == nil || e.offset > r.farthestFail.offset {
		r.farthestFail = &e
	}
	panic(e)
}

// newRDScanner readies the Parser's reusable parse Scanner over sql,
// carrying the same configuration as the parser's lexer (which ParseSQL
// has already reset and applied params to).
func (parser *Parser) newRDScanner(sql string) *Scanner {
	s := &parser.rdScan
	s.reset(sql)
	l := &parser.lexer
	s.client = l.client
	s.connection = l.connection
	s.sqlMode = l.sqlMode
	s.supportWindowFunc = l.supportWindowFunc
	s.skipPositionRecording = l.skipPositionRecording
	s.keepHint = l.keepHint
	return s
}

// lexOne appends the next token to the window. All scanner state feedback
// (identifierDot, lastKeyword hint gating, SQL-mode token rewrites)
// happens inside Scanner.Lex exactly as in a goyacc parse, because tokens
// are produced in the same order.
func (r *rdParser) lexOne() {
	var v yySymType
	tok := r.sc.Lex(&v)
	// Extend the window and fill the slot in place rather than building
	// an rdToken and copying it in: the struct is large and this is the
	// hottest allocation-free path in the parser. Recycled slots hold
	// stale tokens, so every field is (re)assigned.
	if len(r.win) == cap(r.win) {
		r.win = append(r.win, rdToken{})
	} else {
		r.win = r.win[:len(r.win)+1]
	}
	t := &r.win[len(r.win)-1]
	t.tok, t.lit, t.item, t.offset = tok, v.ident, v.item, v.offset
	p := r.sc.r.pos()
	t.endOffset, t.endLine, t.endCol = p.Offset, p.Line, p.Col
	if tok == hintComment {
		t.hintPos = r.sc.lastHintPos
	} else {
		t.hintPos = Pos{}
	}
	if len(r.sc.errs) > 0 {
		// A lexing problem already recorded its error(s). The parse is
		// abandoned, so the token left in the window is harmless.
		panic(rdLexError{})
	}
	if tok == invalid {
		// The parser side of an invalid token is a plain syntax error at
		// its position.
		panic(rdSyntaxError{offset: t.offset, err: r.buildSyntaxError(t)})
	}
	if tok == 0 {
		r.done = true
	}
}

// at returns the token at absolute index abs, lexing forward as needed.
// Past EOF it returns the EOF token. The in-window fast path is kept
// small enough to inline into tok/la/cur, which are the parser's hottest
// calls.
func (r *rdParser) at(abs int) *rdToken {
	idx := abs - r.base
	if idx >= len(r.win) {
		idx = r.fill(abs)
	}
	return &r.win[idx]
}

// fill lexes until the window covers abs (or EOF) and returns the window
// index to read.
func (r *rdParser) fill(abs int) int {
	for abs-r.base >= len(r.win) {
		if r.done {
			return len(r.win) - 1
		}
		r.lexOne()
	}
	return abs - r.base
}

func (r *rdParser) cur() *rdToken { return r.c }

// tok returns the current token id (0 at EOF).
func (r *rdParser) tok() int { return r.c.tok }

// la returns the token id k positions ahead (la(0) == tok()).
func (r *rdParser) la(k int) int { return r.at(r.i + k).tok }

// advance moves past the current token and opportunistically drops window
// tokens that no active mark or the cursor can reach again.
func (r *rdParser) advance() {
	if r.c.tok != 0 {
		r.i++
	}
	low := r.i
	if len(r.marks) > 0 && r.marks[0] < low {
		low = r.marks[0]
	}
	if low-r.base >= 64 {
		n := copy(r.win, r.win[low-r.base:])
		r.win = r.win[:n]
		r.base = low
	}
	r.c = r.at(r.i)
}

// mark pins the current position for a speculative parse. Every mark is
// released by exactly one rewind (go back) or unmark (keep position).
func (r *rdParser) mark() int {
	r.marks = append(r.marks, r.i)
	return r.i
}

func (r *rdParser) unmark() {
	r.marks = r.marks[:len(r.marks)-1]
}

func (r *rdParser) rewind(m int) {
	r.i = m
	r.c = r.at(m)
	r.unmark()
}

// accept consumes the current token if it matches.
func (r *rdParser) accept(tok int) bool {
	if r.tok() == tok {
		r.advance()
		return true
	}
	return false
}

// expect consumes the current token, which must match.
func (r *rdParser) expect(tok int) rdToken {
	if r.tok() != tok {
		r.syntaxError()
	}
	t := *r.cur()
	r.advance()
	return t
}

// parseRD attempts to parse sql with the recursive-descent parser.
// handled=false means the caller must run the goyacc parser instead;
// nothing observable has happened in that case.
func (parser *Parser) parseRD(sql string) (stmts []ast.StmtNode, warns []error, err error) {
	r := &rdParser{
		p:   parser,
		sc:  parser.newRDScanner(sql),
		src: sql,
		win: parser.rdWin[:0],
	}
	defer func() { parser.rdWin = r.win[:0] }()
	defer panics.Handle(func(e interface{}) {
		// A parse failure is final: the syntax error is appended to the
		// scanner's error list the way yyParse appended Errorf(""), and
		// the first recorded error wins, with warnings returned
		// alongside. Speculative parsing means the automaton's
		// farthest-reached failure is the one to report.
		switch v := e.(type) {
		case rdUnsupported:
			err := v.err
			if r.farthestFail != nil && r.farthestFail.offset > v.offset {
				err = r.farthestFail.err
			}
			r.sc.AppendError(err)
		case rdSyntaxError:
			err := v.err
			if r.farthestFail != nil && r.farthestFail.offset > v.offset {
				err = r.farthestFail.err
			}
			r.sc.AppendError(err)
		case rdLexError, rdActionAbort:
			// Errors already recorded on the scanner.
		default:
			panic(e)
		}
		lexWarns, lexErrs := r.sc.Errors()
		if len(lexWarns) > 0 {
			warns = slices.Clone(lexWarns)
		} else {
			warns = nil
		}
		stmts, err = nil, lexErrs[0]
	})
	r.c = r.at(r.i)
	r.parseStatementList()

	lexWarns, lexErrs := r.sc.Errors()
	if len(lexErrs) > 0 {
		if len(lexWarns) > 0 {
			warns = slices.Clone(lexWarns)
		}
		return nil, warns, lexErrs[0]
	}
	if len(lexWarns) > 0 {
		warns = slices.Clone(lexWarns)
	}
	for _, stmt := range r.result {
		ast.SetFlag(stmt)
	}
	return r.result, warns, nil
}

// parseStatementList implements:
//
//	Start: StatementList
//	StatementList: Statement | StatementList ';' Statement
//
// including the stmtText() bookkeeping the goyacc actions perform: a
// statement's text runs from the end of the previous non-empty statement
// through its terminating ';' (with the historical single-newline trims),
// and empty statements do not advance the start position.
func (r *rdParser) parseStatementList() {
	for {
		stmt := r.parseStatement()
		switch r.tok() {
		case ';':
			if stmt != nil {
				r.finishStmt(stmt, r.cur().offset+1)
			}
			r.advance()
			if r.tok() == 0 {
				return
			}
		case 0:
			if stmt != nil {
				r.finishStmt(stmt, r.cur().offset)
			}
			return
		default:
			r.syntaxError()
		}
	}
}

// finishStmt reproduces Scanner.stmtText() plus the StatementList action:
// set the statement's text and append it to the result.
func (r *rdParser) finishStmt(stmt ast.StmtNode, endPos int) {
	if r.src[endPos-1] == '\n' {
		endPos--
	}
	if r.src[r.stmtStart] == '\n' {
		r.stmtStart++
	}
	text := r.src[r.stmtStart:endPos]
	r.stmtStart = endPos
	r.p.setNodeText(stmt, text)
	r.result = append(r.result, stmt)
}

// parseStatement implements the Statement production. A nil return is an
// EmptyStmt. Statement families not yet ported call r.unsupported, which
// falls back to goyacc.
func (r *rdParser) parseStatement() ast.StmtNode {
	switch r.tok() {
	case ';', 0:
		// EmptyStmt: /* empty */
		return nil
	case do:
		// DoStmt: "DO" ExpressionList
		r.advance()
		return &ast.DoStmt{Exprs: r.parseExpressionList()}
	case selectKwd, tableKwd, values, int('('):
		// Statement: SelectStmt | SetOprStmt | SubSelect
		if r.tok() == int('(') && !r.subSelectFollows() {
			r.syntaxError()
		}
		return r.finishSelectFamily(nil)
	case with:
		// A WITH clause prefixes SELECT, UPDATE, and DELETE statements.
		withClause := r.parseWithClause()
		switch r.tok() {
		case update:
			return r.parseUpdateStmt(withClause)
		case deleteKwd:
			return r.parseDeleteStmt(withClause)
		default:
			return r.finishSelectFamily(withClause)
		}
	case insert:
		return r.parseInsertIntoStmt()
	case replace:
		return r.parseReplaceIntoStmt()
	case update:
		return r.parseUpdateStmt(nil)
	case deleteKwd:
		return r.parseDeleteStmt(nil)
	case load:
		switch r.la(1) {
		case data:
			return r.parseLoadDataStmt()
		case stats:
			return r.parseLoadStatsStmt()
		case index:
			return r.parseLoadIndexStmt()
		case xml:
			return r.parseLoadXMLStmt()
		}
		r.unsupported("LOAD statement")
		return nil
	case importKwd:
		switch r.la(1) {
		case into:
			return r.parseImportIntoStmt()
		case tableKwd:
			return r.parseImportTableStmt()
		}
		r.unsupported("IMPORT statement")
		return nil
	case batch:
		return r.parseNonTransactionalDMLStmt()
	default:
		if f, ok := rdStmtDispatch[r.tok()]; ok {
			return f(r)
		}
		r.unsupported(fmt.Sprintf("statement starting with token %d (%q)", r.tok(), r.cur().lit))
		return nil
	}
}

// rdStmtDispatch maps a statement's leading token to its parse function.
// Families register in their own file's init so files stay independent;
// each leading token has exactly one owner.
var rdStmtDispatch = map[int]func(*rdParser) ast.StmtNode{}

func rdRegister(tok int, f func(*rdParser) ast.StmtNode) {
	if _, dup := rdStmtDispatch[tok]; dup {
		panic(fmt.Sprintf("rdRegister: token %d already registered", tok))
	}
	rdStmtDispatch[tok] = f
}

// finishSelectFamily parses the select statement family after an
// optional pre-parsed WITH clause, applying the Statement: SubSelect
// action to a bare parenthesized query.
func (r *rdParser) finishSelectFamily(withClause *ast.WithClause) ast.StmtNode {
	stmt, sub := r.parseSelectCoreWith(withClause)
	if sub != nil {
		switch x := sub.Query.(type) {
		case *ast.SelectStmt:
			x.IsInBraces = true
			return x
		case *ast.SetOprStmt:
			x.IsInBraces = true
			return x
		}
	}
	return stmt
}

var (
	rdLogMu   sync.Mutex
	rdLogPath = os.Getenv("MARINO_RD_LOG")
)
