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

// This file is the core of the hand-written recursive-descent parser that
// replaces the goyacc-generated one (see PLAN.md). While the goyacc parser
// is still in-tree it remains the fallback: any construct the RD parser
// does not implement yet — and, during the migration, any input the RD
// parser would reject — is re-parsed by goyacc, so behavior is unchanged
// until coverage is complete.
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

	// stmtStart mirrors Scanner.stmtStartPos bookkeeping for stmtText().
	stmtStart int
	result    []ast.StmtNode
}

// rdUnsupported unwinds the RD parse when it reaches grammar that is not
// implemented yet; the caller falls back to goyacc.
type rdUnsupported struct {
	construct string
	offset    int
}

// rdSyntaxError unwinds the RD parse on invalid input. While goyacc is
// in-tree it is treated exactly like rdUnsupported so that error messages
// keep coming from goyacc; when goyacc is removed it becomes the real
// error path.
type rdSyntaxError struct {
	offset int
}

// rdLexError unwinds the RD parse when the scanner reports an error while
// lexing on demand; goyacc redoes the parse and reports it identically.
type rdLexError struct{}

// unsupported aborts the RD parse for a construct that is not implemented.
func (r *rdParser) unsupported(construct string) {
	panic(rdUnsupported{construct: construct, offset: r.cur().offset})
}

// syntaxError aborts the RD parse at the current token.
func (r *rdParser) syntaxError() {
	panic(rdSyntaxError{offset: r.cur().offset})
}

// newRDScanner returns a new Scanner over sql carrying the same
// configuration as the parser's lexer (which ParseSQL has already reset
// and applied params to).
func (parser *Parser) newRDScanner(sql string) *Scanner {
	s := NewScanner(sql)
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
	t := rdToken{tok: tok, lit: v.ident, item: v.item, offset: v.offset}
	p := r.sc.r.pos()
	t.endOffset, t.endLine, t.endCol = p.Offset, p.Line, p.Col
	if tok == hintComment {
		t.hintPos = r.sc.lastHintPos
	}
	if tok == invalid || len(r.sc.errs) > 0 {
		// A lexing problem: let goyacc redo the parse and report it.
		panic(rdLexError{})
	}
	r.win = append(r.win, t)
	if tok == 0 {
		r.done = true
	}
}

// at returns the token at absolute index abs, lexing forward as needed.
// Past EOF it returns the EOF token.
func (r *rdParser) at(abs int) *rdToken {
	for abs-r.base >= len(r.win) {
		if r.done {
			return &r.win[len(r.win)-1]
		}
		r.lexOne()
	}
	return &r.win[abs-r.base]
}

func (r *rdParser) cur() *rdToken { return r.at(r.i) }

// tok returns the current token id (0 at EOF).
func (r *rdParser) tok() int { return r.at(r.i).tok }

// la returns the token id k positions ahead (la(0) == tok()).
func (r *rdParser) la(k int) int { return r.at(r.i + k).tok }

// advance moves past the current token and opportunistically drops window
// tokens that no active mark or the cursor can reach again.
func (r *rdParser) advance() {
	if t := r.at(r.i); t.tok != 0 {
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
func (parser *Parser) parseRD(sql string) (stmts []ast.StmtNode, warns []error, err error, handled bool) {
	r := &rdParser{
		p:   parser,
		sc:  parser.newRDScanner(sql),
		src: sql,
	}
	defer panics.Handle(func(e interface{}) {
		switch v := e.(type) {
		case rdUnsupported:
			logRDFallback(sql, fmt.Sprintf("unsupported %s at %d", v.construct, v.offset))
			handled = false
		case rdSyntaxError:
			logRDFallback(sql, fmt.Sprintf("syntax error at %d", v.offset))
			handled = false
		case rdLexError:
			logRDFallback(sql, "lex error")
			handled = false
		case rdActionAbort:
			// The aborting error lives on the discarded scanner clone;
			// goyacc re-parses and reports it identically.
			logRDFallback(sql, "action abort")
			handled = false
		default:
			panic(e)
		}
	})
	r.parseStatementList()

	lexWarns, lexErrs := r.sc.Errors()
	if len(lexErrs) > 0 {
		// Errors raised by warning/validation helpers during the parse:
		// goyacc would report these identically, but routing through it
		// keeps ordering and messages authoritative during the migration.
		logRDFallback(sql, "deferred error: "+lexErrs[0].Error())
		return nil, nil, nil, false
	}
	if len(lexWarns) > 0 {
		warns = slices.Clone(lexWarns)
	}
	for _, stmt := range r.result {
		ast.SetFlag(stmt)
	}
	return r.result, warns, nil, true
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
		if r.la(1) != data {
			r.unsupported("LOAD statement")
		}
		return r.parseLoadDataStmt()
	case importKwd:
		if r.la(1) != into {
			r.unsupported("IMPORT statement")
		}
		return r.parseImportIntoStmt()
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

// logRDFallback records inputs the RD parser could not handle, as the
// migration's todo list (see cmd/next-test).
func logRDFallback(sql, reason string) {
	if rdLogPath == "" {
		return
	}
	rdLogMu.Lock()
	defer rdLogMu.Unlock()
	f, err := os.OpenFile(rdLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer f.Close()
	fmt.Fprintf(f, "%s\x00%s\x00\n", reason, sql)
}
