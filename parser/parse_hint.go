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

// This file is a hand-written recursive-descent port of the goyacc
// optimizer-hint grammar (hintparser.y). It is a second, tiny parser next
// to the statement RD parser in rd_parser.go: it lexes with the exact
// hintScanner the goyacc hint parser drives, one lookahead token at a
// time, so error positions, lexer-generated errors and warnings are
// byte-identical to yyhintParse.
//
// Error model: hintparser.y declares no `error` productions, so goyacc's
// error recovery never finds an error shift — the first syntax error
// appends `Errorf("")` ("Optimizer hint syntax error at ...") and aborts
// the whole parse with a nil result. Unrecognized-but-well-formed hints
// are ordinary productions that emit a warning and produce a nil hint,
// which the list actions skip. Both behaviors are reproduced here.

import (
	"math"
	"strconv"

	"github.com/sqlc-dev/marino/ast"
	"github.com/sqlc-dev/marino/internal/panics"
	"github.com/sqlc-dev/marino/mysql"
)

// rdHintAbort unwinds rdParseHint the way goyacc's `goto ret1` does:
// discard the result, keep whatever errors are on the scanner.
type rdHintAbort struct{}

// rdHintParser is the recursive-descent driver. It holds exactly one
// lookahead token (tok plus its semantic values), mirroring the
// yyhintParse loop which always fetches one lookahead before acting, so
// the scanner position at any error matches goyacc's.
type rdHintParser struct {
	lexer hintScanner

	// current lookahead token: id, ident value and number value
	// (yyhintSymType.ident / .number of the last Lex call).
	tok    int
	ident  string
	number uint64
}

// rdParseHint parses an optimizer hint comment `/*+ ... */`. It is the
// recursive-descent equivalent of hintParser.parse and produces the same
// hints and the same warning/error lists.
func rdParseHint(input string, sqlMode mysql.SQLMode, initPos Pos) ([]*ast.TableOptimizerHint, []error) {
	hp := &rdHintParser{}
	hp.lexer.reset(input[3:])
	hp.lexer.SetSQLMode(sqlMode)
	hp.lexer.r.updatePos(Pos{
		Line:   initPos.Line,
		Col:    initPos.Col + 3, // skipped the initial '/*+'
		Offset: 0,
	})
	hp.lexer.inBangComment = true // skip the final '*/' (we need the '*/' for reporting warnings)

	result := hp.run()

	warns, errs := hp.lexer.Errors()
	if len(errs) == 0 {
		errs = warns
	}
	return result, errs
}

// run implements Start: OptimizerHintList, aborting with a nil result on
// the first syntax error exactly like yyhintParse.
func (hp *rdHintParser) run() (result []*ast.TableOptimizerHint) {
	defer panics.Handle(func(e interface{}) {
		if _, ok := e.(rdHintAbort); ok {
			result = nil
			return
		}
		panic(e)
	})
	hp.next()
	result = hp.parseOptimizerHintList()
	if hp.tok != 0 {
		hp.syntaxError()
	}
	return result
}

// next fetches the next lookahead token from the hint scanner.
func (hp *rdHintParser) next() {
	var lval yyhintSymType
	tok := hp.lexer.Lex(&lval)
	if tok < 0 {
		tok = 0
	}
	hp.tok, hp.ident, hp.number = tok, lval.ident, lval.number
}

// advance consumes the current token, returning its semantic values, and
// fetches the next lookahead (matching goyacc, which always holds one
// lookahead when a reduce action runs).
func (hp *rdHintParser) advance() (ident string, number uint64) {
	ident, number = hp.ident, hp.number
	hp.next()
	return ident, number
}

// expect consumes the current token, which must match tok.
func (hp *rdHintParser) expect(tok int) (ident string, number uint64) {
	if hp.tok != tok {
		hp.syntaxError()
	}
	return hp.advance()
}

// syntaxError reports a syntax error at the current lookahead token and
// aborts the parse. This mirrors goyacc's brand-new-error handling:
// `yylex.AppendError(yylex.Errorf(""))` followed by a failed recovery
// (the hint grammar has no error productions) which returns 1.
func (hp *rdHintParser) syntaxError() {
	hp.lexer.AppendError(hp.lexer.Errorf(""))
	panic(rdHintAbort{})
}

// abort unwinds like a `return 1` from a grammar action: no additional
// error is appended.
func (hp *rdHintParser) abort() {
	panic(rdHintAbort{})
}

func (hp *rdHintParser) warnUnsupportedHint(name string) {
	warn := ErrWarnOptimizerHintUnsupportedHint.FastGenByArgs(name)
	hp.lexer.warns = append(hp.lexer.warns, warn)
}

// isIdentifier reports whether the current token matches the Identifier
// production: hintIdentifier or any hint keyword except "PARTITION".
func (hp *rdHintParser) isIdentifier() bool {
	switch hp.tok {
	case hintIdentifier,
		// MySQL 8.0 hint names
		hintJoinFixedOrder, hintJoinOrder, hintJoinPrefix, hintJoinSuffix,
		hintBKA, hintNoBKA, hintBNL, hintNoBNL, hintHashJoin,
		hintHashJoinBuild, hintHashJoinProbe, hintNoHashJoin, hintMerge,
		hintNoMerge, hintIndexMerge, hintNoIndexMerge, hintMRR, hintNoMRR,
		hintNoICP, hintNoRangeOptimization, hintSkipScan, hintNoSkipScan,
		hintSemijoin, hintNoSemijoin, hintMaxExecutionTime, hintSetVar,
		hintResourceGroup, hintQBName, hintHypoIndex,
		// TiDB hint names
		hintAggToCop, hintLimitToCop, hintIgnorePlanCache, hintWriteSlowLog,
		hintHashAgg, hintMpp1PhaseAgg, hintMpp2PhaseAgg, hintIgnoreIndex,
		hintInlHashJoin, hintIndexHashJoin, hintNoIndexHashJoin, hintInlJoin,
		hintIndexJoin, hintNoIndexJoin, hintInlMergeJoin, hintIndexMergeJoin,
		hintNoIndexMergeJoin, hintMemoryQuota, hintNoSwapJoinInputs,
		hintQueryType, hintReadConsistentReplica, hintReadFromStorage,
		hintSMJoin, hintNoSMJoin, hintBCJoin, hintShuffleJoin, hintStreamAgg,
		hintSwapJoinInputs, hintUseIndexMerge, hintUseIndex, hintOrderIndex,
		hintNoOrderIndex, hintIndexLookUpPushDown, hintNoIndexLookUpPushDown,
		hintUsePlanCache, hintUseToja, hintTimeRange, hintUseCascades,
		hintNthPlan, hintForceIndex, hintStraightJoin, hintLeading,
		hintSemiJoinRewrite, hintNoDecorrelate,
		// other keywords
		hintOLAP, hintOLTP, hintTiKV, hintTiFlash, hintFalse, hintTrue,
		hintMB, hintGB, hintDupsWeedOut, hintFirstMatch, hintLooseScan,
		hintMaterialization:
		return true
	}
	return false
}

// parseIdentifier implements Identifier (all keyword alternatives keep
// the source spelling, like the token's ident value).
func (hp *rdHintParser) parseIdentifier() string {
	if !hp.isIdentifier() {
		hp.syntaxError()
	}
	id, _ := hp.advance()
	return id
}

// canStartHint reports whether the current token can begin a
// TableOptimizerHintOpt or StorageOptimizerHintOpt production.
func (hp *rdHintParser) canStartHint() bool {
	switch hp.tok {
	case hintJoinFixedOrder,
		// JoinOrderOptimizerHintName
		hintJoinOrder, hintJoinPrefix, hintJoinSuffix,
		// UnsupportedTableLevelOptimizerHintName
		hintBKA, hintNoBKA, hintBNL, hintNoBNL, hintNoMerge,
		// SupportedTableLevelOptimizerHintName
		hintSMJoin, hintNoSMJoin, hintBCJoin, hintShuffleJoin, hintInlJoin,
		hintIndexJoin, hintNoIndexJoin, hintMerge, hintInlHashJoin,
		hintIndexHashJoin, hintNoIndexHashJoin, hintSwapJoinInputs,
		hintNoSwapJoinInputs, hintInlMergeJoin, hintIndexMergeJoin,
		hintNoIndexMergeJoin, hintHashJoin, hintNoHashJoin,
		hintHashJoinBuild, hintHashJoinProbe, hintHypoIndex,
		hintLeading,
		// UnsupportedIndexLevelOptimizerHintName
		hintIndexMerge, hintMRR, hintNoMRR, hintNoICP,
		hintNoRangeOptimization, hintSkipScan, hintNoSkipScan,
		// SupportedIndexLevelOptimizerHintName
		hintUseIndex, hintIgnoreIndex, hintUseIndexMerge, hintForceIndex,
		hintOrderIndex, hintNoOrderIndex, hintIndexLookUpPushDown,
		hintNoIndexLookUpPushDown,
		// SubqueryOptimizerHintName
		hintSemijoin, hintNoSemijoin,
		hintMaxExecutionTime, hintNthPlan, hintSetVar, hintResourceGroup,
		hintQBName, hintMemoryQuota, hintTimeRange,
		// BooleanHintName
		hintUseToja, hintUseCascades,
		// NullaryHintName
		hintUsePlanCache, hintHashAgg, hintMpp1PhaseAgg, hintMpp2PhaseAgg,
		hintStreamAgg, hintAggToCop, hintLimitToCop, hintNoIndexMerge,
		hintReadConsistentReplica, hintIgnorePlanCache, hintStraightJoin,
		hintSemiJoinRewrite, hintNoDecorrelate,
		hintWriteSlowLog, hintQueryType, hintIdentifier,
		// StorageOptimizerHintOpt
		hintReadFromStorage:
		return true
	}
	return false
}

// parseOptimizerHintList implements:
//
//	OptimizerHintList:
//	  TableOptimizerHintOpt
//	| OptimizerHintList CommaOpt TableOptimizerHintOpt
//	| StorageOptimizerHintOpt
//	| OptimizerHintList CommaOpt StorageOptimizerHintOpt
//
// A nil TableOptimizerHintOpt (unsupported hint) is skipped.
func (hp *rdHintParser) parseOptimizerHintList() []*ast.TableOptimizerHint {
	hints := hp.parseHintListItem(nil)
	for {
		if hp.tok == ',' { // CommaOpt
			hp.next()
		} else if !hp.canStartHint() {
			return hints
		}
		hints = hp.parseHintListItem(hints)
	}
}

// parseHintListItem parses one TableOptimizerHintOpt or
// StorageOptimizerHintOpt and appends its result per the list actions.
func (hp *rdHintParser) parseHintListItem(hints []*ast.TableOptimizerHint) []*ast.TableOptimizerHint {
	if hp.tok == hintReadFromStorage {
		return append(hints, hp.parseStorageOptimizerHintOpt()...)
	}
	if h := hp.parseTableOptimizerHintOpt(); h != nil {
		return append(hints, h)
	}
	return hints
}

// parseTableOptimizerHintOpt implements every TableOptimizerHintOpt
// alternative except READ_FROM_STORAGE (see
// parseStorageOptimizerHintOpt). Unsupported hints warn and return nil.
func (hp *rdHintParser) parseTableOptimizerHintOpt() *ast.TableOptimizerHint {
	switch hp.tok {
	case hintJoinFixedOrder:
		// "JOIN_FIXED_ORDER" '(' QueryBlockOpt ')'
		name, _ := hp.advance()
		hp.expect('(')
		hp.parseQueryBlockOpt()
		hp.expect(')')
		hp.warnUnsupportedHint(name)
		return nil

	case hintJoinOrder, hintJoinPrefix, hintJoinSuffix:
		// JoinOrderOptimizerHintName '(' HintTableList ')'
		name, _ := hp.advance()
		hp.expect('(')
		hp.parseHintTableList()
		hp.expect(')')
		hp.warnUnsupportedHint(name)
		return nil

	case hintBKA, hintNoBKA, hintBNL, hintNoBNL, hintNoMerge:
		// UnsupportedTableLevelOptimizerHintName '(' HintTableListOpt ')'
		name, _ := hp.advance()
		hp.expect('(')
		hp.parseHintTableListOpt()
		hp.expect(')')
		hp.warnUnsupportedHint(name)
		return nil

	case hintSMJoin, hintNoSMJoin, hintBCJoin, hintShuffleJoin, hintInlJoin,
		hintIndexJoin, hintNoIndexJoin, hintMerge, hintInlHashJoin,
		hintIndexHashJoin, hintNoIndexHashJoin, hintSwapJoinInputs,
		hintNoSwapJoinInputs, hintInlMergeJoin, hintIndexMergeJoin,
		hintNoIndexMergeJoin, hintHashJoin, hintNoHashJoin,
		hintHashJoinBuild, hintHashJoinProbe, hintHypoIndex:
		// SupportedTableLevelOptimizerHintName '(' HintTableListOpt ')'
		name, _ := hp.advance()
		hp.expect('(')
		h := hp.parseHintTableListOpt()
		hp.expect(')')
		h.HintName = ast.NewCIStr(name)
		return h

	case hintLeading:
		// "LEADING" '(' QueryBlockOpt LeadingTableList ')'
		name, _ := hp.advance()
		hp.expect('(')
		qb := hp.parseQueryBlockOpt()
		list := hp.parseLeadingTableList()
		hp.expect(')')
		h := &ast.TableOptimizerHint{
			HintName: ast.NewCIStr(name),
			QBName:   ast.NewCIStr(qb),
			HintData: list,
		}
		// For LEADING hints we need to maintain two views of the tables:
		// h.HintData:
		//   - Stores the structured AST node (LeadingList).
		//   - Preserves the nesting and order information of LEADING(...),
		//
		// h.Tables:
		//   - Stores a flat slice of all HintTable elements inside the LeadingList.
		//   - Only used for initialization.
		if leadingList, ok := h.HintData.(*ast.LeadingList); ok {
			// be compatible with the prior flatten writing style
			h.Tables = ast.FlattenLeadingList(leadingList)
		}
		return h

	case hintIndexMerge, hintMRR, hintNoMRR, hintNoICP,
		hintNoRangeOptimization, hintSkipScan, hintNoSkipScan:
		// UnsupportedIndexLevelOptimizerHintName '(' HintIndexList ')'
		name, _ := hp.advance()
		hp.expect('(')
		hp.parseHintIndexList()
		hp.expect(')')
		hp.warnUnsupportedHint(name)
		return nil

	case hintUseIndex, hintIgnoreIndex, hintUseIndexMerge, hintForceIndex,
		hintOrderIndex, hintNoOrderIndex, hintIndexLookUpPushDown,
		hintNoIndexLookUpPushDown:
		// SupportedIndexLevelOptimizerHintName '(' HintIndexList ')'
		name, _ := hp.advance()
		hp.expect('(')
		h := hp.parseHintIndexList()
		hp.expect(')')
		h.HintName = ast.NewCIStr(name)
		return h

	case hintSemijoin, hintNoSemijoin:
		// SubqueryOptimizerHintName '(' QueryBlockOpt SubqueryStrategiesOpt ')'
		name, _ := hp.advance()
		hp.expect('(')
		hp.parseQueryBlockOpt()
		hp.parseSubqueryStrategiesOpt()
		hp.expect(')')
		hp.warnUnsupportedHint(name)
		return nil

	case hintMaxExecutionTime:
		// "MAX_EXECUTION_TIME" '(' QueryBlockOpt hintIntLit ')'
		name, _ := hp.advance()
		hp.expect('(')
		qb := hp.parseQueryBlockOpt()
		_, num := hp.expect(hintIntLit)
		hp.expect(')')
		return &ast.TableOptimizerHint{
			HintName: ast.NewCIStr(name),
			QBName:   ast.NewCIStr(qb),
			HintData: num,
		}

	case hintNthPlan:
		// "NTH_PLAN" '(' QueryBlockOpt hintIntLit ')'
		name, _ := hp.advance()
		hp.expect('(')
		qb := hp.parseQueryBlockOpt()
		_, num := hp.expect(hintIntLit)
		hp.expect(')')
		return &ast.TableOptimizerHint{
			HintName: ast.NewCIStr(name),
			QBName:   ast.NewCIStr(qb),
			HintData: int64(num),
		}

	case hintSetVar:
		// "SET_VAR" '(' Identifier '=' Value ')'
		name, _ := hp.advance()
		hp.expect('(')
		varName := hp.parseIdentifier()
		hp.expect('=')
		value := hp.parseValue()
		hp.expect(')')
		return &ast.TableOptimizerHint{
			HintName: ast.NewCIStr(name),
			HintData: ast.HintSetVar{
				VarName: varName,
				Value:   value,
			},
		}

	case hintResourceGroup:
		// "RESOURCE_GROUP" '(' Identifier ')'
		name, _ := hp.advance()
		hp.expect('(')
		group := hp.parseIdentifier()
		hp.expect(')')
		return &ast.TableOptimizerHint{
			HintName: ast.NewCIStr(name),
			HintData: group,
		}

	case hintQBName:
		// "QB_NAME" '(' Identifier ')'
		// | "QB_NAME" '(' Identifier ',' ViewNameList ')'
		name, _ := hp.advance()
		hp.expect('(')
		qb := hp.parseIdentifier()
		if hp.tok == ',' {
			hp.next()
			views := hp.parseViewNameList()
			hp.expect(')')
			return &ast.TableOptimizerHint{
				HintName: ast.NewCIStr(name),
				QBName:   ast.NewCIStr(qb),
				Tables:   views.Tables,
			}
		}
		hp.expect(')')
		return &ast.TableOptimizerHint{
			HintName: ast.NewCIStr(name),
			QBName:   ast.NewCIStr(qb),
		}

	case hintMemoryQuota:
		// "MEMORY_QUOTA" '(' QueryBlockOpt hintIntLit UnitOfBytes ')'
		name, _ := hp.advance()
		hp.expect('(')
		qb := hp.parseQueryBlockOpt()
		_, num := hp.expect(hintIntLit)
		unit := hp.parseUnitOfBytes()
		hp.expect(')')
		maxValue := uint64(math.MaxInt64) / unit
		if num <= maxValue {
			return &ast.TableOptimizerHint{
				HintName: ast.NewCIStr(name),
				HintData: int64(num * unit),
				QBName:   ast.NewCIStr(qb),
			}
		}
		hp.lexer.AppendError(ErrWarnMemoryQuotaOverflow.GenWithStackByArgs(math.MaxInt))
		hp.lexer.lastErrorAsWarn()
		return nil

	case hintTimeRange:
		// "TIME_RANGE" '(' hintStringLit CommaOpt hintStringLit ')'
		name, _ := hp.advance()
		hp.expect('(')
		from, _ := hp.expect(hintStringLit)
		if hp.tok == ',' { // CommaOpt
			hp.next()
		}
		to, _ := hp.expect(hintStringLit)
		hp.expect(')')
		return &ast.TableOptimizerHint{
			HintName: ast.NewCIStr(name),
			HintData: ast.HintTimeRange{
				From: from,
				To:   to,
			},
		}

	case hintUseToja, hintUseCascades:
		// BooleanHintName '(' QueryBlockOpt HintTrueOrFalse ')'
		name, _ := hp.advance()
		hp.expect('(')
		qb := hp.parseQueryBlockOpt()
		h := hp.parseHintTrueOrFalse()
		hp.expect(')')
		h.HintName = ast.NewCIStr(name)
		h.QBName = ast.NewCIStr(qb)
		return h

	case hintUsePlanCache, hintHashAgg, hintMpp1PhaseAgg, hintMpp2PhaseAgg,
		hintStreamAgg, hintAggToCop, hintLimitToCop, hintNoIndexMerge,
		hintReadConsistentReplica, hintIgnorePlanCache, hintStraightJoin,
		hintSemiJoinRewrite, hintNoDecorrelate:
		// NullaryHintName '(' QueryBlockOpt ')'
		name, _ := hp.advance()
		hp.expect('(')
		qb := hp.parseQueryBlockOpt()
		hp.expect(')')
		return &ast.TableOptimizerHint{
			HintName: ast.NewCIStr(name),
			QBName:   ast.NewCIStr(qb),
		}

	case hintWriteSlowLog:
		// "WRITE_SLOW_LOG"
		name, _ := hp.advance()
		return &ast.TableOptimizerHint{
			HintName: ast.NewCIStr(name),
		}

	case hintQueryType:
		// "QUERY_TYPE" '(' QueryBlockOpt HintQueryType ')'
		name, _ := hp.advance()
		hp.expect('(')
		qb := hp.parseQueryBlockOpt()
		var qt string
		switch hp.tok {
		case hintOLAP, hintOLTP: // HintQueryType: "OLAP" | "OLTP"
			qt, _ = hp.advance()
		default:
			hp.syntaxError()
		}
		hp.expect(')')
		return &ast.TableOptimizerHint{
			HintName: ast.NewCIStr(name),
			QBName:   ast.NewCIStr(qb),
			HintData: ast.NewCIStr(qt),
		}

	case hintIdentifier:
		return hp.parsePseudoHint()

	default:
		hp.syntaxError()
		return nil
	}
}

// parsePseudoHint implements the unsupported pseudo hints:
//
//	hintIdentifier '(' QueryBlockOpt hintIntLit ')'
//	hintIdentifier '(' PartitionList ')'
//	hintIdentifier '(' PartitionList CommaOpt hintIntLit ')'
//	hintIdentifier '(' Identifier '=' Value ')'
//
// All alternatives warn and produce a nil hint.
func (hp *rdHintParser) parsePseudoHint() *ast.TableOptimizerHint {
	name, _ := hp.advance()
	hp.expect('(')
	switch {
	case hp.tok == hintSingleAtIdentifier:
		// QueryBlockOpt hintIntLit
		hp.next()
		hp.expect(hintIntLit)
	case hp.tok == hintIntLit:
		// empty QueryBlockOpt, hintIntLit
		hp.next()
	case hp.isIdentifier():
		hp.parseIdentifier()
		if hp.tok == '=' {
			// Identifier '=' Value
			hp.next()
			hp.parseValue()
		} else {
			// PartitionList [CommaOpt hintIntLit]
		list:
			for {
				switch {
				case hp.tok == ',': // CommaOpt
					hp.next()
					if hp.tok == hintIntLit {
						hp.next()
						break list
					}
					hp.parseIdentifier()
				case hp.isIdentifier():
					hp.parseIdentifier()
				case hp.tok == hintIntLit:
					hp.next()
					break list
				default:
					break list
				}
			}
		}
	default:
		hp.syntaxError()
	}
	hp.expect(')')
	hp.warnUnsupportedHint(name)
	return nil
}

// parseStorageOptimizerHintOpt implements:
//
//	StorageOptimizerHintOpt:
//	  "READ_FROM_STORAGE" '(' QueryBlockOpt HintStorageTypeAndTableList ')'
func (hp *rdHintParser) parseStorageOptimizerHintOpt() []*ast.TableOptimizerHint {
	storageName, _ := hp.advance()
	hp.expect('(')
	qb := hp.parseQueryBlockOpt()
	// HintStorageTypeAndTableList:
	//   HintStorageTypeAndTable
	// | HintStorageTypeAndTableList ',' HintStorageTypeAndTable
	hs := []*ast.TableOptimizerHint{hp.parseHintStorageTypeAndTable()}
	for hp.tok == ',' {
		hp.next()
		hs = append(hs, hp.parseHintStorageTypeAndTable())
	}
	hp.expect(')')
	name := ast.NewCIStr(storageName)
	qbName := ast.NewCIStr(qb)
	for _, h := range hs {
		h.HintName = name
		h.QBName = qbName
	}
	return hs
}

// parseHintStorageTypeAndTable implements:
//
//	HintStorageTypeAndTable: HintStorageType '[' HintTableList ']'
func (hp *rdHintParser) parseHintStorageTypeAndTable() *ast.TableOptimizerHint {
	var storageType string
	switch hp.tok {
	case hintTiKV, hintTiFlash: // HintStorageType: "TIKV" | "TIFLASH"
		storageType, _ = hp.advance()
	default:
		hp.syntaxError()
	}
	hp.expect('[')
	h := hp.parseHintTableList()
	hp.expect(']')
	h.HintData = ast.NewCIStr(storageType)
	return h
}

// parseQueryBlockOpt implements QueryBlockOpt: /* empty */ | hintSingleAtIdentifier.
func (hp *rdHintParser) parseQueryBlockOpt() string {
	if hp.tok == hintSingleAtIdentifier {
		qb, _ := hp.advance()
		return qb
	}
	return ""
}

// parseHintTableListOpt implements HintTableListOpt: HintTableList | QueryBlockOpt.
func (hp *rdHintParser) parseHintTableListOpt() *ast.TableOptimizerHint {
	if hp.tok == hintSingleAtIdentifier {
		qb, _ := hp.advance()
		if hp.isIdentifier() {
			return hp.parseHintTableListTail(qb)
		}
		return &ast.TableOptimizerHint{
			QBName: ast.NewCIStr(qb),
		}
	}
	if hp.isIdentifier() {
		return hp.parseHintTableListTail("")
	}
	return &ast.TableOptimizerHint{
		QBName: ast.NewCIStr(""),
	}
}

// parseHintTableList implements:
//
//	HintTableList: QueryBlockOpt HintTable | HintTableList ',' HintTable
func (hp *rdHintParser) parseHintTableList() *ast.TableOptimizerHint {
	return hp.parseHintTableListTail(hp.parseQueryBlockOpt())
}

// parseHintTableListTail parses HintTableList after its QueryBlockOpt.
func (hp *rdHintParser) parseHintTableListTail(qb string) *ast.TableOptimizerHint {
	h := &ast.TableOptimizerHint{
		Tables: []ast.HintTable{hp.parseHintTable()},
		QBName: ast.NewCIStr(qb),
	}
	for hp.tok == ',' {
		hp.next()
		h.Tables = append(h.Tables, hp.parseHintTable())
	}
	return h
}

// parseHintTable implements:
//
//	HintTable:
//	  Identifier QueryBlockOpt PartitionListOpt
//	| Identifier '.' Identifier QueryBlockOpt PartitionListOpt
func (hp *rdHintParser) parseHintTable() ast.HintTable {
	first := hp.parseIdentifier()
	if hp.tok == '.' {
		hp.next()
		table := hp.parseIdentifier()
		qb := hp.parseQueryBlockOpt()
		return ast.HintTable{
			DBName:        ast.NewCIStr(first),
			TableName:     ast.NewCIStr(table),
			QBName:        ast.NewCIStr(qb),
			PartitionList: hp.parsePartitionListOpt(),
		}
	}
	qb := hp.parseQueryBlockOpt()
	return ast.HintTable{
		TableName:     ast.NewCIStr(first),
		QBName:        ast.NewCIStr(qb),
		PartitionList: hp.parsePartitionListOpt(),
	}
}

// parsePartitionListOpt implements:
//
//	PartitionListOpt: /* empty */ | "PARTITION" '(' PartitionList ')'
//	PartitionList: Identifier | PartitionList CommaOpt Identifier
func (hp *rdHintParser) parsePartitionListOpt() []ast.CIStr {
	if hp.tok != hintPartition {
		return nil
	}
	hp.next()
	hp.expect('(')
	list := []ast.CIStr{ast.NewCIStr(hp.parseIdentifier())}
	for {
		if hp.tok == ',' { // CommaOpt
			hp.next()
			list = append(list, ast.NewCIStr(hp.parseIdentifier()))
		} else if hp.isIdentifier() {
			list = append(list, ast.NewCIStr(hp.parseIdentifier()))
		} else {
			break
		}
	}
	hp.expect(')')
	return list
}

// parseHintIndexList implements:
//
//	HintIndexList: QueryBlockOpt HintTable CommaOpt IndexNameListOpt
//	IndexNameListOpt: /* empty */ | IndexNameList
//	IndexNameList: Identifier | IndexNameList ',' Identifier
func (hp *rdHintParser) parseHintIndexList() *ast.TableOptimizerHint {
	qb := hp.parseQueryBlockOpt()
	table := hp.parseHintTable()
	if hp.tok == ',' { // CommaOpt
		hp.next()
	}
	var h *ast.TableOptimizerHint
	if hp.isIdentifier() {
		h = &ast.TableOptimizerHint{
			Indexes: []ast.CIStr{ast.NewCIStr(hp.parseIdentifier())},
		}
		for hp.tok == ',' {
			hp.next()
			h.Indexes = append(h.Indexes, ast.NewCIStr(hp.parseIdentifier()))
		}
	} else {
		h = &ast.TableOptimizerHint{}
	}
	h.Tables = []ast.HintTable{table}
	h.QBName = ast.NewCIStr(qb)
	return h
}

// parseViewNameList implements:
//
//	ViewNameList: ViewNameList '.' ViewName | ViewName
func (hp *rdHintParser) parseViewNameList() *ast.TableOptimizerHint {
	h := &ast.TableOptimizerHint{
		Tables: []ast.HintTable{hp.parseViewName()},
	}
	for hp.tok == '.' {
		hp.next()
		h.Tables = append(h.Tables, hp.parseViewName())
	}
	return h
}

// parseViewName implements ViewName: Identifier QueryBlockOpt | QueryBlockOpt.
func (hp *rdHintParser) parseViewName() ast.HintTable {
	if hp.isIdentifier() {
		name, _ := hp.advance()
		qb := hp.parseQueryBlockOpt()
		return ast.HintTable{
			TableName: ast.NewCIStr(name),
			QBName:    ast.NewCIStr(qb),
		}
	}
	return ast.HintTable{
		QBName: ast.NewCIStr(hp.parseQueryBlockOpt()),
	}
}

// parseLeadingTableList implements:
//
//	LeadingTableList: LeadingTableElement | LeadingTableList ',' LeadingTableElement
func (hp *rdHintParser) parseLeadingTableList() *ast.LeadingList {
	list := &ast.LeadingList{Items: []interface{}{hp.parseLeadingTableElement()}}
	for hp.tok == ',' {
		hp.next()
		list.Items = append(list.Items, hp.parseLeadingTableElement())
	}
	return list
}

// parseLeadingTableElement implements:
//
//	LeadingTableElement: HintTable | '(' LeadingTableList ')'
func (hp *rdHintParser) parseLeadingTableElement() interface{} {
	if hp.tok == '(' {
		hp.next()
		list := hp.parseLeadingTableList()
		hp.expect(')')
		return list
	}
	tmp := hp.parseHintTable()
	return &tmp
}

// parseSubqueryStrategiesOpt implements:
//
//	SubqueryStrategiesOpt: /* empty */ | SubqueryStrategies
//	SubqueryStrategies: SubqueryStrategy | SubqueryStrategies ',' SubqueryStrategy
//	SubqueryStrategy: "DUPSWEEDOUT" | "FIRSTMATCH" | "LOOSESCAN" | "MATERIALIZATION"
func (hp *rdHintParser) parseSubqueryStrategiesOpt() {
	switch hp.tok {
	case hintDupsWeedOut, hintFirstMatch, hintLooseScan, hintMaterialization:
		hp.next()
	default:
		return
	}
	for hp.tok == ',' {
		hp.next()
		switch hp.tok {
		case hintDupsWeedOut, hintFirstMatch, hintLooseScan, hintMaterialization:
			hp.next()
		default:
			hp.syntaxError()
		}
	}
}

// parseValue implements:
//
//	Value: hintStringLit | Identifier | hintIntLit | '+' hintIntLit | '-' hintIntLit
func (hp *rdHintParser) parseValue() string {
	switch hp.tok {
	case hintStringLit:
		v, _ := hp.advance()
		return v
	case hintIntLit:
		_, num := hp.advance()
		return strconv.FormatUint(num, 10)
	case '+':
		hp.next()
		_, num := hp.expect(hintIntLit)
		return strconv.FormatUint(num, 10)
	case '-':
		hp.next()
		_, num := hp.expect(hintIntLit)
		// goyacc's exact-lookahead LALR tables only reduce this
		// alternative (running its action) when the lookahead is in
		// Follow(Value) = {')'}; on any other token the parser reports a
		// plain syntax error there instead, so the range check below must
		// not run.
		if hp.tok != ')' {
			hp.syntaxError()
		}
		if num > 9223372036854775808 {
			hp.lexer.AppendError(hp.lexer.Errorf("the Signed Value should be at the range of [-9223372036854775808, 9223372036854775807]."))
			hp.abort() // `return 1` in the grammar action
			return ""
		} else if num == 9223372036854775808 {
			signedOne := int64(1)
			return strconv.FormatInt(signedOne<<63, 10)
		}
		return strconv.FormatInt(-int64(num), 10)
	default:
		return hp.parseIdentifier()
	}
}

// parseUnitOfBytes implements UnitOfBytes: "MB" | "GB".
func (hp *rdHintParser) parseUnitOfBytes() uint64 {
	switch hp.tok {
	case hintMB:
		hp.next()
		return 1024 * 1024
	case hintGB:
		hp.next()
		return 1024 * 1024 * 1024
	default:
		hp.syntaxError()
		return 0
	}
}

// parseHintTrueOrFalse implements HintTrueOrFalse: "TRUE" | "FALSE".
func (hp *rdHintParser) parseHintTrueOrFalse() *ast.TableOptimizerHint {
	switch hp.tok {
	case hintTrue:
		hp.next()
		return &ast.TableOptimizerHint{HintData: true}
	case hintFalse:
		hp.next()
		return &ast.TableOptimizerHint{HintData: false}
	default:
		hp.syntaxError()
		return nil
	}
}
