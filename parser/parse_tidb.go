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

// TiDB-specific management statements: SplitRegionStmt,
// DistributeTableStmt, the CANCEL family (CancelDistributionJobStmt,
// CancelImportStmt, the "CANCEL BR JOB" and "CANCEL TRAFFIC JOBS"
// alternatives of BRIEStmt/TrafficStmt), CalibrateResourceStmt,
// RecommendIndexStmt, the QUERY WATCH statements (AddQueryWatchStmt,
// DropQueryWatchStmt), TrafficStmt, PlanReplayerStmt, RefreshStatsStmt,
// LoadStatsStmt, and LockStatsStmt/UnlockStatsStmt (wired from the
// lock/unlock dispatchers in parse_misc.go).

import (
	"strings"
	"time"

	"github.com/sqlc-dev/marino/ast"
	"github.com/sqlc-dev/marino/duration"
)

func init() {
	rdRegister(split, (*rdParser).parseSplitRegionStmt)
	rdRegister(distribute, (*rdParser).parseDistributeTableStmt)
	rdRegister(cancel, (*rdParser).parseCancelStmtFamily)
	rdRegister(calibrate, (*rdParser).parseCalibrateResourceStmt)
	rdRegister(recommend, (*rdParser).parseRecommendIndexStmt)
	rdRegister(query, (*rdParser).parseQueryWatchStmt)
	rdRegister(traffic, (*rdParser).parseTrafficStmt)
	rdRegister(plan, (*rdParser).parsePlanReplayerStmt)
	rdRegister(refresh, (*rdParser).parseRefreshStatsStmt)
}

// parseSplitRegionStmt implements SplitRegionStmt and SplitSyntaxOption
// (SplitOption/SplitOptionBetween live in parse_column.go).
func (r *rdParser) parseSplitRegionStmt() ast.StmtNode {
	r.expect(split)
	syntaxOpt := &ast.SplitSyntaxOption{}
	switch r.tok() {
	case region:
		// "REGION" "FOR" ["PARTITION"]
		r.advance()
		r.expect(forKwd)
		syntaxOpt.HasRegionFor = true
		if r.accept(partition) {
			syntaxOpt.HasPartition = true
		}
	case partition:
		r.advance()
		syntaxOpt.HasPartition = true
	}
	r.expect(tableKwd)
	table := r.parseTableName()
	partitionNames := r.parsePartitionNameListOpt()
	if r.accept(index) {
		// ... "INDEX" Identifier SplitOption
		indexName := r.parseIdentifier()
		return &ast.SplitRegionStmt{
			SplitSyntaxOpt: syntaxOpt,
			Table:          table,
			PartitionNames: partitionNames,
			IndexName:      ast.NewCIStr(indexName),
			SplitOpt:       r.parseSplitOption(),
		}
	}
	return &ast.SplitRegionStmt{
		SplitSyntaxOpt: syntaxOpt,
		Table:          table,
		PartitionNames: partitionNames,
		SplitOpt:       r.parseSplitOption(),
	}
}

// parseDistributeTableStmt implements DistributeTableStmt.
func (r *rdParser) parseDistributeTableStmt() ast.StmtNode {
	r.expect(distribute)
	r.expect(tableKwd)
	stmt := &ast.DistributeTableStmt{
		Table:          r.parseTableName(),
		PartitionNames: r.parsePartitionNameListOpt(),
	}
	r.expect(rule)
	r.parseEqOpt()
	stmt.Rule = r.expect(stringLit).lit
	r.expect(engine)
	r.parseEqOpt()
	stmt.Engine = r.expect(stringLit).lit
	if r.accept(timeout) {
		r.parseEqOpt()
		stmt.Timeout = r.expect(stringLit).lit
	}
	return stmt
}

// parseCancelStmtFamily dispatches the statements led by "CANCEL":
//
//	CancelDistributionJobStmt: "CANCEL" "DISTRIBUTION" "JOB" Int64Num
//	CancelImportStmt: "CANCEL" "IMPORT" "JOB" Int64Num
//	BRIEStmt: "CANCEL" "BR" "JOB" Int64Num
//	TrafficStmt: "CANCEL" "TRAFFIC" "JOBS"
func (r *rdParser) parseCancelStmtFamily() ast.StmtNode {
	r.expect(cancel)
	switch r.tok() {
	case distribution:
		r.advance()
		r.expect(job)
		return &ast.CancelDistributionJobStmt{
			JobID: r.parseInt64Num(),
		}
	case importKwd:
		r.advance()
		r.expect(job)
		return &ast.ImportIntoActionStmt{
			Tp:    ast.ImportIntoCancel,
			JobID: r.parseInt64Num(),
		}
	case br:
		r.advance()
		r.expect(job)
		stmt := &ast.BRIEStmt{}
		stmt.Kind = ast.BRIEKindCancelJob
		stmt.JobID = r.parseInt64Num()
		return stmt
	case traffic:
		r.advance()
		r.expect(jobs)
		return &ast.TrafficStmt{
			OpType: ast.TrafficOpCancel,
		}
	}
	r.syntaxError()
	return nil
}

// parseCalibrateResourceStmt implements CalibrateResourceStmt,
// CalibrateOption and CalibrateResourceWorkloadOption.
func (r *rdParser) parseCalibrateResourceStmt() ast.StmtNode {
	r.expect(calibrate)
	r.expect(resource)
	switch r.tok() {
	case workload:
		r.advance()
		var tp ast.CalibrateResourceType
		switch r.tok() {
		case tpcc:
			tp = ast.TPCC
		case oltpReadWrite:
			tp = ast.OLTPREADWRITE
		case oltpReadOnly:
			tp = ast.OLTPREADONLY
		case oltpWriteOnly:
			tp = ast.OLTPWRITEONLY
		case tpch10:
			tp = ast.TPCH10
		default:
			r.syntaxError()
		}
		r.advance()
		return &ast.CalibrateResourceStmt{Tp: tp}
	case startTime, endTime, timeDuration:
		return &ast.CalibrateResourceStmt{
			DynamicCalibrateResourceOptionList: r.parseDynamicCalibrateOptionList(),
		}
	}
	return &ast.CalibrateResourceStmt{}
}

// parseDynamicCalibrateOptionList implements DynamicCalibrateOptionList
// (options separated by nothing or by a comma), including its duplicate
// checks.
func (r *rdParser) parseDynamicCalibrateOptionList() []*ast.DynamicCalibrateResourceOption {
	list := []*ast.DynamicCalibrateResourceOption{r.parseDynamicCalibrateResourceOption()}
	for {
		if !r.accept(int(',')) {
			switch r.tok() {
			case startTime, endTime, timeDuration:
			default:
				return list
			}
		}
		opt := r.parseDynamicCalibrateResourceOption()
		if list[0].Tp == opt.Tp ||
			(len(list) > 1 && list[1].Tp == opt.Tp) {
			r.actionError(r.sc.Errorf("Dupliated options specified"))
		}
		list = append(list, opt)
	}
}

// parseDynamicCalibrateResourceOption implements
// DynamicCalibrateResourceOption.
func (r *rdParser) parseDynamicCalibrateResourceOption() *ast.DynamicCalibrateResourceOption {
	switch r.tok() {
	case startTime:
		// "START_TIME" EqOpt Expression
		r.advance()
		r.parseEqOpt()
		return &ast.DynamicCalibrateResourceOption{Tp: ast.CalibrateStartTime, Ts: r.parseExpression()}
	case endTime:
		// "END_TIME" EqOpt Expression
		r.advance()
		r.parseEqOpt()
		return &ast.DynamicCalibrateResourceOption{Tp: ast.CalibrateEndTime, Ts: r.parseExpression()}
	case timeDuration:
		r.advance()
		r.parseEqOpt()
		if r.accept(interval) {
			// "DURATION" EqOpt "INTERVAL" Expression TimeUnit
			ts := r.parseExpression()
			return &ast.DynamicCalibrateResourceOption{Tp: ast.CalibrateDuration, Ts: ts, Unit: r.parseTimeUnit()}
		}
		// "DURATION" EqOpt stringLit
		dur := r.expect(stringLit).lit
		if _, err := duration.ParseDuration(dur); err != nil {
			r.actionError(r.sc.Errorf("The DURATION option is not a valid duration: %s", err.Error()))
		}
		return &ast.DynamicCalibrateResourceOption{Tp: ast.CalibrateDuration, StrValue: dur}
	}
	r.syntaxError()
	return nil
}

// parseRecommendIndexStmt implements RecommendIndexStmt.
func (r *rdParser) parseRecommendIndexStmt() ast.StmtNode {
	r.expect(recommend)
	r.expect(index)
	switch r.tok() {
	case run:
		r.advance()
		if r.accept(forKwd) {
			// "RECOMMEND" "INDEX" "RUN" "FOR" stringLit
			// RecommendIndexOptionListOpt
			sql := r.expect(stringLit).lit
			return &ast.RecommendIndexStmt{
				Action:  "run",
				SQL:     sql,
				Options: r.parseRecommendIndexOptionListOpt(),
			}
		}
		// "RECOMMEND" "INDEX" "RUN" RecommendIndexOptionListOpt
		return &ast.RecommendIndexStmt{
			Action:  "run",
			Options: r.parseRecommendIndexOptionListOpt(),
		}
	case show:
		// "RECOMMEND" "INDEX" "SHOW" "OPTION"
		r.advance()
		r.expect(option)
		return &ast.RecommendIndexStmt{
			Action: "show",
		}
	case apply:
		// "RECOMMEND" "INDEX" "APPLY" NUM
		r.advance()
		return &ast.RecommendIndexStmt{
			Action: "apply",
			ID:     r.expect(intLit).item.(int64),
		}
	case ignore:
		// "RECOMMEND" "INDEX" "IGNORE" NUM
		r.advance()
		return &ast.RecommendIndexStmt{
			Action: "ignore",
			ID:     r.expect(intLit).item.(int64),
		}
	case set:
		// "RECOMMEND" "INDEX" "SET" RecommendIndexOptionList
		r.advance()
		return &ast.RecommendIndexStmt{
			Action:  "set",
			Options: r.parseRecommendIndexOptionList(),
		}
	}
	r.syntaxError()
	return nil
}

// parseRecommendIndexOptionListOpt implements RecommendIndexOptionListOpt.
func (r *rdParser) parseRecommendIndexOptionListOpt() []ast.RecommendIndexOption {
	if !r.accept(with) {
		return []ast.RecommendIndexOption{}
	}
	return r.parseRecommendIndexOptionList()
}

// parseRecommendIndexOptionList implements RecommendIndexOptionList.
func (r *rdParser) parseRecommendIndexOptionList() []ast.RecommendIndexOption {
	list := []ast.RecommendIndexOption{r.parseRecommendIndexOption()}
	for r.accept(int(',')) {
		list = append(list, r.parseRecommendIndexOption())
	}
	return list
}

// parseRecommendIndexOption implements RecommendIndexOption:
// Identifier "=" Literal.
func (r *rdParser) parseRecommendIndexOption() ast.RecommendIndexOption {
	name := r.parseIdentifier()
	r.expect(eq)
	lit := r.parseLiteralExpr(r.cur().offset)
	return ast.RecommendIndexOption{
		Option: name,
		Value:  ast.NewValueExpr(lit, r.p.charset, r.p.collation),
	}
}

// parseQueryWatchStmt implements AddQueryWatchStmt and DropQueryWatchStmt.
func (r *rdParser) parseQueryWatchStmt() ast.StmtNode {
	r.expect(query)
	r.expect(watch)
	switch r.tok() {
	case add:
		// "QUERY" "WATCH" "ADD" QueryWatchOptionList
		r.advance()
		return &ast.AddQueryWatchStmt{
			QueryWatchOptionList: r.parseQueryWatchOptionList(),
		}
	case remove:
		r.advance()
		if r.tok() == intLit {
			// "QUERY" "WATCH" "REMOVE" NUM
			return &ast.DropQueryWatchStmt{
				IntValue: r.expect(intLit).item.(int64),
			}
		}
		r.expect(resource)
		r.expect(group)
		if r.tok() == singleAtIdentifier {
			// "QUERY" "WATCH" "REMOVE" "RESOURCE" "GROUP" UserVariable
			return &ast.DropQueryWatchStmt{
				GroupNameExpr: r.parseUserVariable(),
			}
		}
		// "QUERY" "WATCH" "REMOVE" "RESOURCE" "GROUP" ResourceGroupName
		return &ast.DropQueryWatchStmt{
			GroupNameStr: ast.NewCIStr(r.parseResourceGroupName()),
		}
	}
	r.syntaxError()
	return nil
}

// parseQueryWatchOptionList implements QueryWatchOptionList (options
// separated by nothing or by a comma), including its duplicate checks.
func (r *rdParser) parseQueryWatchOptionList() []*ast.QueryWatchOption {
	list := []*ast.QueryWatchOption{r.parseQueryWatchOption()}
	for {
		if !r.accept(int(',')) {
			switch r.tok() {
			case resource, action, sql, plan:
			default:
				return list
			}
		}
		opt := r.parseQueryWatchOption()
		if !ast.CheckQueryWatchAppend(list, opt) {
			r.actionError(r.sc.Errorf("Dupliated options specified"))
		}
		list = append(list, opt)
	}
}

// parseQueryWatchOption implements QueryWatchOption and
// QueryWatchTextOption.
func (r *rdParser) parseQueryWatchOption() *ast.QueryWatchOption {
	switch r.tok() {
	case resource:
		// "RESOURCE" "GROUP" (ResourceGroupName | UserVariable)
		r.advance()
		r.expect(group)
		if r.tok() == singleAtIdentifier {
			return &ast.QueryWatchOption{
				Tp: ast.QueryWatchResourceGroup,
				ResourceGroupOption: &ast.QueryWatchResourceGroupOption{
					GroupNameExpr: r.parseUserVariable(),
				},
			}
		}
		return &ast.QueryWatchOption{
			Tp: ast.QueryWatchResourceGroup,
			ResourceGroupOption: &ast.QueryWatchResourceGroupOption{
				GroupNameStr: ast.NewCIStr(r.parseResourceGroupName()),
			},
		}
	case action:
		// "ACTION" EqOpt ResourceGroupRunawayActionOption
		r.advance()
		r.parseEqOpt()
		return &ast.QueryWatchOption{
			Tp:           ast.QueryWatchAction,
			ActionOption: r.parseResourceGroupRunawayActionOption(),
		}
	case sql:
		r.advance()
		if r.accept(digest) {
			// "SQL" "DIGEST" SimpleExpr
			return &ast.QueryWatchOption{
				Tp: ast.QueryWatchType,
				TextOption: &ast.QueryWatchTextOption{
					Type:        ast.WatchSimilar,
					PatternExpr: r.parseSimpleExpr(),
				},
			}
		}
		// "SQL" "TEXT" ResourceGroupRunawayWatchOption "TO" SimpleExpr
		r.expect(textType)
		watchType := r.parseResourceGroupRunawayWatchOption()
		r.expect(to)
		return &ast.QueryWatchOption{
			Tp: ast.QueryWatchType,
			TextOption: &ast.QueryWatchTextOption{
				Type:          watchType,
				PatternExpr:   r.parseSimpleExpr(),
				TypeSpecified: true,
			},
		}
	case plan:
		// "PLAN" "DIGEST" SimpleExpr
		r.advance()
		r.expect(digest)
		return &ast.QueryWatchOption{
			Tp: ast.QueryWatchType,
			TextOption: &ast.QueryWatchTextOption{
				Type:        ast.WatchPlan,
				PatternExpr: r.parseSimpleExpr(),
			},
		}
	}
	r.syntaxError()
	return nil
}

// parseResourceGroupRunawayWatchOption implements
// ResourceGroupRunawayWatchOption: "EXACT" | "SIMILAR" | "PLAN".
func (r *rdParser) parseResourceGroupRunawayWatchOption() ast.RunawayWatchType {
	switch r.tok() {
	case exact:
		r.advance()
		return ast.WatchExact
	case similar:
		r.advance()
		return ast.WatchSimilar
	case plan:
		r.advance()
		return ast.WatchPlan
	}
	r.syntaxError()
	return 0
}

// parseTrafficStmt implements the TrafficStmt alternatives led by
// "TRAFFIC" (the "SHOW"/"CANCEL" ones are led by other families'
// tokens), plus TrafficCaptureOptList and TrafficReplayOptList.
func (r *rdParser) parseTrafficStmt() ast.StmtNode {
	r.expect(traffic)
	switch r.tok() {
	case capture:
		// "TRAFFIC" "CAPTURE" "TO" stringLit TrafficCaptureOptList
		r.advance()
		r.expect(to)
		x := &ast.TrafficStmt{
			OpType: ast.TrafficOpCapture,
			Dir:    r.expect(stringLit).lit,
		}
		options := []*ast.TrafficOption{r.parseTrafficCaptureOpt()}
		for {
			switch r.tok() {
			case timeDuration, encryptionMethod, compress:
				options = append(options, r.parseTrafficCaptureOpt())
				continue
			}
			break
		}
		x.Options = options
		return x
	case replay:
		// "TRAFFIC" "REPLAY" "FROM" stringLit TrafficReplayOptList
		r.advance()
		r.expect(from)
		x := &ast.TrafficStmt{
			OpType: ast.TrafficOpReplay,
			Dir:    r.expect(stringLit).lit,
		}
		options := []*ast.TrafficOption{r.parseTrafficReplayOpt()}
		for {
			switch r.tok() {
			case user, password, speed, readOnly:
				options = append(options, r.parseTrafficReplayOpt())
				continue
			}
			break
		}
		x.Options = options
		return x
	}
	r.syntaxError()
	return nil
}

// parseTrafficCaptureOpt implements TrafficCaptureOpt.
func (r *rdParser) parseTrafficCaptureOpt() *ast.TrafficOption {
	switch r.tok() {
	case timeDuration:
		// "DURATION" EqOpt stringLit
		r.advance()
		r.parseEqOpt()
		dur := r.expect(stringLit).lit
		if _, err := time.ParseDuration(dur); err != nil {
			r.actionError(r.sc.Errorf("The DURATION option is not a valid duration: %s", err.Error()))
		}
		return &ast.TrafficOption{OptionType: ast.TrafficOptionDuration, StrValue: dur}
	case encryptionMethod:
		// "ENCRYPTION_METHOD" EqOpt stringLit
		r.advance()
		r.parseEqOpt()
		return &ast.TrafficOption{OptionType: ast.TrafficOptionEncryptionMethod, StrValue: r.expect(stringLit).lit}
	case compress:
		// "COMPRESS" EqOpt Boolean
		r.advance()
		r.parseEqOpt()
		return &ast.TrafficOption{OptionType: ast.TrafficOptionCompress, BoolValue: r.parseBoolean()}
	}
	r.syntaxError()
	return nil
}

// parseTrafficReplayOpt implements TrafficReplayOpt.
func (r *rdParser) parseTrafficReplayOpt() *ast.TrafficOption {
	switch r.tok() {
	case user:
		// "USER" EqOpt stringLit
		r.advance()
		r.parseEqOpt()
		return &ast.TrafficOption{OptionType: ast.TrafficOptionUsername, StrValue: r.expect(stringLit).lit}
	case password:
		// "PASSWORD" EqOpt stringLit
		r.advance()
		r.parseEqOpt()
		return &ast.TrafficOption{OptionType: ast.TrafficOptionPassword, StrValue: r.expect(stringLit).lit}
	case speed:
		// "SPEED" EqOpt NumLiteral
		r.advance()
		r.parseEqOpt()
		return &ast.TrafficOption{OptionType: ast.TrafficOptionSpeed, FloatValue: ast.NewValueExpr(r.parseNumLiteralValue(), "", "")}
	case readOnly:
		// "READ_ONLY" EqOpt Boolean
		r.advance()
		r.parseEqOpt()
		return &ast.TrafficOption{OptionType: ast.TrafficOptionReadOnly, BoolValue: r.parseBoolean()}
	}
	r.syntaxError()
	return nil
}

// parsePlanReplayerStmt implements PlanReplayerStmt and
// PlanReplayerDumpOpt.
func (r *rdParser) parsePlanReplayerStmt() ast.StmtNode {
	r.expect(plan)
	r.expect(replayer)
	switch r.tok() {
	case dump:
		r.advance()
		// PlanReplayerDumpOpt: empty | "WITH" "STATS" AsOfClause
		var histInfo *ast.AsOfClause
		if r.accept(with) {
			r.expect(stats)
			histInfo = r.parseAsOfClause()
		}
		r.expect(explain)
		analyzeFlag := r.accept(analyze)
		switch {
		case r.tok() == slow:
			// ... "EXPLAIN" ["ANALYZE"] "SLOW" "QUERY" WhereClauseOptional
			// OrderByOptional SelectStmtLimitOpt
			r.advance()
			r.expect(query)
			x := &ast.PlanReplayerStmt{
				Stmt:    nil,
				Analyze: analyzeFlag,
				Load:    false,
				File:    "",
			}
			if histInfo != nil {
				x.HistoricalStatsInfo = histInfo
			}
			if where := r.parseWhereClauseOptional(); where != nil {
				x.Where = where
			}
			if r.tok() == order {
				x.OrderBy = r.parseOrderBy()
			}
			if r.tok() == limit || r.tok() == fetch {
				x.Limit = r.parseSelectStmtLimit()
			}
			return x
		case r.tok() == stringLit:
			// ... "EXPLAIN" ["ANALYZE"] stringLit
			x := &ast.PlanReplayerStmt{
				Stmt:    nil,
				Analyze: analyzeFlag,
				Load:    false,
				File:    r.expect(stringLit).lit,
			}
			if histInfo != nil {
				x.HistoricalStatsInfo = histInfo
			}
			return x
		case r.tok() == int('(') && r.la(1) == stringLit:
			// ... "EXPLAIN" ["ANALYZE"] '(' StringList ')'
			r.advance()
			stmtList := []string{r.expect(stringLit).lit}
			for r.accept(int(',')) {
				stmtList = append(stmtList, r.expect(stringLit).lit)
			}
			r.expect(int(')'))
			x := &ast.PlanReplayerStmt{
				Stmt:     nil,
				Analyze:  analyzeFlag,
				Load:     false,
				File:     "",
				StmtList: stmtList,
			}
			if histInfo != nil {
				x.HistoricalStatsInfo = histInfo
			}
			return x
		default:
			// ... "EXPLAIN" ["ANALYZE"] ExplainableStmt, recording the
			// statement's source text from its first token to the end of
			// the input (parser.src[startOffset:] in the actions).
			startOffset := r.cur().offset
			stmt := r.parseExplainableStmt()
			x := &ast.PlanReplayerStmt{
				Stmt:    stmt,
				Analyze: analyzeFlag,
				Load:    false,
				File:    "",
				Where:   nil,
				OrderBy: nil,
				Limit:   nil,
			}
			if histInfo != nil {
				x.HistoricalStatsInfo = histInfo
			}
			r.p.setNodeText(x.Stmt, strings.TrimSpace(r.src[startOffset:]))
			return x
		}
	case load:
		// "PLAN" "REPLAYER" "LOAD" stringLit
		r.advance()
		return &ast.PlanReplayerStmt{
			Stmt:    nil,
			Analyze: false,
			Load:    true,
			File:    r.expect(stringLit).lit,
			Where:   nil,
			OrderBy: nil,
			Limit:   nil,
		}
	case capture:
		r.advance()
		if r.accept(remove) {
			// "PLAN" "REPLAYER" "CAPTURE" "REMOVE" stringLit stringLit
			sqlDigest := r.expect(stringLit).lit
			return &ast.PlanReplayerStmt{
				Stmt:       nil,
				Analyze:    false,
				Remove:     true,
				SQLDigest:  sqlDigest,
				PlanDigest: r.expect(stringLit).lit,
				Where:      nil,
				OrderBy:    nil,
				Limit:      nil,
			}
		}
		// "PLAN" "REPLAYER" "CAPTURE" stringLit stringLit
		sqlDigest := r.expect(stringLit).lit
		return &ast.PlanReplayerStmt{
			Stmt:       nil,
			Analyze:    false,
			Capture:    true,
			SQLDigest:  sqlDigest,
			PlanDigest: r.expect(stringLit).lit,
			Where:      nil,
			OrderBy:    nil,
			Limit:      nil,
		}
	}
	r.syntaxError()
	return nil
}

// parseRefreshStatsStmt implements RefreshStatsStmt, RefreshStatsModeOpt
// and RefreshStatsClusterOpt (StatsObjectList lives in parse_misc.go).
func (r *rdParser) parseRefreshStatsStmt() ast.StmtNode {
	r.expect(refresh)
	r.expect(stats)
	stmt := &ast.RefreshStatsStmt{
		RefreshObjects: r.parseStatsObjectList(),
	}
	switch r.tok() {
	case full:
		r.advance()
		mode := ast.RefreshStatsModeFull
		stmt.RefreshMode = &mode
	case lite:
		r.advance()
		mode := ast.RefreshStatsModeLite
		stmt.RefreshMode = &mode
	}
	stmt.IsClusterWide = r.accept(cluster)
	return stmt
}

// parseLoadStatsStmt implements LoadStatsStmt: "LOAD" "STATS" stringLit.
// It is dispatched from the load case of parseStatement.
func (r *rdParser) parseLoadStatsStmt() ast.StmtNode {
	r.expect(load)
	r.expect(stats)
	return &ast.LoadStatsStmt{
		Path: r.expect(stringLit).lit,
	}
}

// parseLockStatsStmt implements LockStatsStmt (dispatched from the lock
// statement family owner).
func (r *rdParser) parseLockStatsStmt() ast.StmtNode {
	r.expect(lock)
	r.expect(stats)
	table := r.parseTableName()
	if r.accept(partition) {
		if r.accept(int('(')) {
			// "LOCK" "STATS" TableName "PARTITION" '(' PartitionNameList ')'
			table.PartitionNames = r.parsePartitionNameList()
			r.expect(int(')'))
		} else {
			// "LOCK" "STATS" TableName "PARTITION" PartitionNameList
			table.PartitionNames = r.parsePartitionNameList()
		}
		return &ast.LockStatsStmt{
			Tables: []*ast.TableName{table},
		}
	}
	// "LOCK" "STATS" TableNameList
	tables := []*ast.TableName{table}
	for r.accept(int(',')) {
		tables = append(tables, r.parseTableName())
	}
	return &ast.LockStatsStmt{
		Tables: tables,
	}
}

// parseUnlockStatsStmt implements UnlockStatsStmt (dispatched from the
// unlock statement family owner).
func (r *rdParser) parseUnlockStatsStmt() ast.StmtNode {
	r.expect(unlock)
	r.expect(stats)
	table := r.parseTableName()
	if r.accept(partition) {
		if r.accept(int('(')) {
			// "UNLOCK" "STATS" TableName "PARTITION" '(' PartitionNameList ')'
			table.PartitionNames = r.parsePartitionNameList()
			r.expect(int(')'))
		} else {
			// "UNLOCK" "STATS" TableName "PARTITION" PartitionNameList
			table.PartitionNames = r.parsePartitionNameList()
		}
		return &ast.UnlockStatsStmt{
			Tables: []*ast.TableName{table},
		}
	}
	// "UNLOCK" "STATS" TableNameList
	tables := []*ast.TableName{table}
	for r.accept(int(',')) {
		tables = append(tables, r.parseTableName())
	}
	return &ast.UnlockStatsStmt{
		Tables: tables,
	}
}
