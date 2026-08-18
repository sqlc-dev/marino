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

// CREATE TABLE: CreateTableStmt with TableElementListOpt, the table
// option catalogue (TableOption/PartDefOption), the partition grammar
// (PartitionOpt, PartitionMethod, PartitionDefinition, ...), the split
// index options, and the CREATE ... SELECT tail.

import (
	"fmt"
	"strings"

	"github.com/sqlc-dev/marino/ast"
	"github.com/sqlc-dev/marino/duration"
	"github.com/sqlc-dev/marino/mysql"
)

func init() {
	rdRegister(create, (*rdParser).parseCreateStmtFamily)
}

// parseCreateStmtFamily dispatches the CREATE statement family on the
// token(s) after "CREATE". The "CREATE OR REPLACE" prefix (OrReplace)
// belongs to CreateViewStmt, CreatePolicyStmt, and
// CreateMaskingPolicyStmt, told apart by the token after "REPLACE".
func (r *rdParser) parseCreateStmtFamily() ast.StmtNode {
	switch r.la(1) {
	case tableKwd, temporary:
		// "CREATE" OptTemporary "TABLE" ...
		return r.parseCreateTableStmt()
	case global:
		if r.la(2) == temporary {
			// "CREATE" "GLOBAL" "TEMPORARY" "TABLE" ...
			return r.parseCreateTableStmt()
		}
		// CreateBindingStmt: "CREATE" GlobalScope "BINDING" ...
		return r.parseCreateBindingStmt()
	case database:
		return r.parseCreateDatabaseStmt()
	case spatial:
		if r.la(2) != index {
			// "CREATE" "SPATIAL" "REFERENCE" "SYSTEM" ...
			return r.parseCreateSpatialReferenceSystemStmt()
		}
		return r.parseCreateIndexStmt()
	case index, unique, fulltext, vectorType, columnar:
		// IndexKeyTypeOpt "INDEX"
		return r.parseCreateIndexStmt()
	case view, algorithm, sql:
		// ViewAlgorithm/ViewSQLSecurity ... "VIEW"
		return r.parseCreateViewStmt()
	case definer:
		// "CREATE" "DEFINER" eq Username prefixes views, events,
		// triggers, and stored routines; the token after the clause
		// decides.
		switch r.peekPastDefiner() {
		case event:
			return r.parseCreateEventStmt()
		case trigger:
			return r.parseCreateTriggerStmt()
		case procedure:
			return r.parseCreateProcedureStmt()
		case function:
			return r.parseCreateFunctionFamily()
		default:
			return r.parseCreateViewStmt()
		}
	case or:
		// "CREATE" "OR" "REPLACE" ...
		switch r.la(3) {
		case placement:
			return r.parseCreatePolicyStmt()
		case masking:
			return r.parseCreateMaskingPolicyStmt()
		case spatial:
			return r.parseCreateSpatialReferenceSystemStmt()
		case library:
			return r.parseCreateLibraryStmt()
		case jsonType:
			return r.parseCreateJSONDualityViewStmt()
		default:
			return r.parseCreateViewStmt()
		}
	case user:
		return r.parseCreateUserStmt()
	case role:
		return r.parseCreateRoleStmt()
	case sequence:
		return r.parseCreateSequenceStmt()
	case statistics:
		return r.parseCreateStatisticsStmt()
	case placement:
		return r.parseCreatePolicyStmt()
	case masking:
		return r.parseCreateMaskingPolicyStmt()
	case resource:
		return r.parseCreateResourceGroupStmt()
	case binding, session:
		// CreateBindingStmt: "CREATE" GlobalScope "BINDING" ...
		return r.parseCreateBindingStmt()
	case procedure:
		return r.parseCreateProcedureStmt()
	case function, aggregate:
		return r.parseCreateFunctionFamily()
	case event:
		return r.parseCreateEventStmt()
	case trigger:
		return r.parseCreateTriggerStmt()
	case server:
		return r.parseCreateServerStmt()
	case tablespace, undo:
		return r.parseCreateTablespaceStmt()
	case logfile:
		return r.parseCreateLogfileGroupStmt()
	case library:
		return r.parseCreateLibraryStmt()
	case jsonType:
		return r.parseCreateJSONDualityViewStmt()
	default:
		// No production continues here; the automaton shifts CREATE and
		// errors at the lookahead, so advance before reporting.
		r.advance()
		r.unsupported(fmt.Sprintf("CREATE %s", r.cur().lit))
	}
	return nil
}

// appendWarnf reproduces the `yylex.AppendError(yylex.Errorf(...));
// parser.lastErrorAsWarn()` warning pattern of the option actions. A
// goyacc reduce always has its one-token lookahead lexed before the
// action runs, so the current token is force-lexed first to give Errorf
// the identical scanner position.
func (r *rdParser) appendWarnf(format string, a ...interface{}) {
	r.sc.AppendError(r.actionErrorf(format, a...))
	r.sc.lastErrorAsWarn()
}

// parseCreateTableStmt implements CreateTableStmt (both the
// TableElementListOpt form and the LikeTableWithOrWithoutParen form),
// plus OptTemporary.
func (r *rdParser) parseCreateTableStmt() ast.StmtNode {
	r.expect(create)
	// OptTemporary
	temporaryKeyword := ast.TemporaryNone
	switch r.tok() {
	case temporary:
		r.advance()
		temporaryKeyword = ast.TemporaryLocal
	case global:
		r.advance()
		r.expect(temporary)
		temporaryKeyword = ast.TemporaryGlobal
	}
	r.expect(tableKwd)
	ifNotExists := r.parseIfNotExists()
	table := r.parseTableName()

	if r.tok() == like || (r.tok() == int('(') && r.la(1) == like) {
		// "CREATE" OptTemporary "TABLE" IfNotExists TableName
		// LikeTableWithOrWithoutParen OnCommitOpt
		tmp := &ast.CreateTableStmt{
			Table:            table,
			ReferTable:       r.parseLikeTableWithOrWithoutParen(),
			IfNotExists:      ifNotExists,
			TemporaryKeyword: temporaryKeyword,
		}
		onCommit := r.parseOnCommitOpt()
		if (onCommit != nil && tmp.TemporaryKeyword != ast.TemporaryGlobal) || (tmp.TemporaryKeyword == ast.TemporaryGlobal && onCommit == nil) {
			r.sc.AppendError(r.actionErrorf("GLOBAL TEMPORARY and ON COMMIT DELETE ROWS must appear together"))
		} else {
			if tmp.TemporaryKeyword == ast.TemporaryGlobal {
				tmp.OnCommitDelete = *onCommit
			}
		}
		return tmp
	}

	// "CREATE" OptTemporary "TABLE" IfNotExists TableName TableElementListOpt
	// CreateTableOptionListOpt PartitionOpt SplitIndexListOpt DuplicateOpt
	// AsOpt CreateTableSelectOpt OnCommitOpt
	stmt := r.parseTableElementListOpt()
	stmt.Table = table
	stmt.IfNotExists = ifNotExists
	stmt.TemporaryKeyword = temporaryKeyword
	stmt.Options = r.parseCreateTableOptionListOpt()
	if partition := r.parsePartitionOpt(); partition != nil {
		stmt.Partition = partition
	}
	if splitIndex := r.parseSplitIndexListOpt(); splitIndex != nil {
		stmt.SplitIndex = splitIndex
	}
	stmt.OnDuplicate = r.parseDuplicateOpt()
	// AsOpt
	r.accept(as)
	stmt.Select = r.parseCreateTableSelectOpt().Select
	onCommit := r.parseOnCommitOpt()
	if (onCommit != nil && stmt.TemporaryKeyword != ast.TemporaryGlobal) || (stmt.TemporaryKeyword == ast.TemporaryGlobal && onCommit == nil) {
		r.sc.AppendError(r.actionErrorf("GLOBAL TEMPORARY and ON COMMIT DELETE ROWS must appear together"))
	} else {
		if stmt.TemporaryKeyword == ast.TemporaryGlobal {
			stmt.OnCommitDelete = *onCommit
		}
	}
	return stmt
}

// parseLikeTableWithOrWithoutParen implements LikeTableWithOrWithoutParen.
func (r *rdParser) parseLikeTableWithOrWithoutParen() *ast.TableName {
	if r.accept(like) {
		return r.parseTableName()
	}
	r.expect(int('('))
	r.expect(like)
	tn := r.parseTableName()
	r.expect(int(')'))
	return tn
}

// parseOnCommitOpt implements OnCommitOpt (nil when absent, the DELETE
// ROWS flag when present).
func (r *rdParser) parseOnCommitOpt() *bool {
	if r.tok() != on {
		return nil
	}
	r.advance()
	r.expect(commit)
	var v bool
	switch r.tok() {
	case deleteKwd:
		v = true
	case preserve:
		v = false
	default:
		r.syntaxError()
	}
	r.advance()
	r.expect(rows)
	return &v
}

// parseDuplicateOpt implements DuplicateOpt.
func (r *rdParser) parseDuplicateOpt() ast.OnDuplicateKeyHandlingType {
	switch r.tok() {
	case ignore:
		r.advance()
		return ast.OnDuplicateKeyHandlingIgnore
	case replace:
		r.advance()
		return ast.OnDuplicateKeyHandlingReplace
	}
	return ast.OnDuplicateKeyHandlingError
}

// parseCreateTableSelectOpt implements CreateTableSelectOpt.
func (r *rdParser) parseCreateTableSelectOpt() *ast.CreateTableStmt {
	switch r.tok() {
	case selectKwd, tableKwd, values, with, int('('):
		// SetOprStmt | SelectStmt | SelectStmtWithClause | SubSelect
		stmt, sub := r.parseSelectCore()
		if sub != nil {
			var sel ast.ResultSetNode
			switch x := sub.Query.(type) {
			case *ast.SelectStmt:
				x.IsInBraces = true
				sel = x
			case *ast.SetOprStmt:
				x.IsInBraces = true
				sel = x
			}
			return &ast.CreateTableStmt{Select: sel}
		}
		return &ast.CreateTableStmt{Select: stmt.(ast.ResultSetNode)}
	}
	return &ast.CreateTableStmt{}
}

// parseTableElementListOpt implements TableElementListOpt.
func (r *rdParser) parseTableElementListOpt() *ast.CreateTableStmt {
	var columnDefs []*ast.ColumnDef
	var constraints []*ast.Constraint
	if r.tok() != int('(') {
		return &ast.CreateTableStmt{
			Cols:        columnDefs,
			Constraints: constraints,
		}
	}
	r.advance()
	tes := r.parseTableElementList()
	r.expect(int(')'))
	for _, te := range tes {
		switch te := te.(type) {
		case *ast.ColumnDef:
			columnDefs = append(columnDefs, te)
		case *ast.Constraint:
			constraints = append(constraints, te)
		}
	}
	return &ast.CreateTableStmt{
		Cols:        columnDefs,
		Constraints: constraints,
	}
}

// parseTableElementList implements TableElementList.
func (r *rdParser) parseTableElementList() []interface{} {
	tes := []interface{}{}
	if te := r.parseTableElement(); te != nil {
		tes = append(tes, te)
	}
	for r.accept(int(',')) {
		if te := r.parseTableElement(); te != nil {
			tes = append(tes, te)
		}
	}
	return tes
}

// parseTableElement implements TableElement: ColumnDef |
// ConstraintWithColumnarIndex.
func (r *rdParser) parseTableElement() interface{} {
	switch r.tok() {
	case constraint, primary, unique, fulltext, spatial, foreign, check, key, index:
		// Constraint: ConstraintKeywordOpt ConstraintElem
		return r.parseConstraint()
	case vectorType, columnar:
		// VECTOR/COLUMNAR are also identifiers; "INDEX" commits to
		// ConstraintVectorIndex/ConstraintColumnarIndex.
		if r.la(1) == index {
			return r.parseConstraintVectorOrColumnarIndex()
		}
	}
	return r.parseColumnDef()
}

// isTableOptionStart reports whether the current token can begin a
// TableOption.
func (r *rdParser) isTableOptionStart() bool {
	if r.isPartDefOptionStart() {
		return true
	}
	switch r.tok() {
	case defaultKwd, charsetKwd, collate, force, autoIncrement, autoIdCache,
		autoRandomBase, avgRowLength, connection, checksum, tableChecksum,
		password, compression, keyBlockSize, delayKeyWrite, rowFormat,
		statsPersistent, statsAutoRecalc, statsSamplePages, statsBuckets,
		statsTopN, statsSampleRate, statsColChoice, statsColList,
		shardRowIDBits, preSplitRegions, packKeys, secondaryEngine, union,
		encryption, ttl, ttlEnable, ttlJobInterval, autoextendSize, affinity,
		pageChecksum, pageCompressed, pageCompressionLevel, transactional,
		sequence, ietfQuotes:
		return true
	case character, charType:
		// CharsetKw: "CHARACTER" "SET" | "CHAR" "SET"
		return r.la(1) == set
	}
	return false
}

// parseCreateTableOptionListOpt implements CreateTableOptionListOpt and
// TableOptionList (options separated by nothing or by a comma; a comma
// always continues the list).
func (r *rdParser) parseCreateTableOptionListOpt() []*ast.TableOption {
	if !r.isTableOptionStart() {
		return []*ast.TableOption{}
	}
	opts := []*ast.TableOption{r.parseTableOption()}
	for {
		if r.accept(int(',')) {
			opts = append(opts, r.parseTableOption())
			continue
		}
		if !r.isTableOptionStart() {
			return opts
		}
		opts = append(opts, r.parseTableOption())
	}
}

// parseTableOption implements TableOption (delegating the PartDefOption
// alternatives to parsePartDefOption).
func (r *rdParser) parseTableOption() *ast.TableOption {
	switch r.tok() {
	case defaultKwd:
		// DefaultKwdOpt CharsetKw EqOpt CharsetName
		// DefaultKwdOpt "COLLATE" EqOpt CollationName
		r.advance()
		if r.accept(collate) {
			r.parseEqOpt()
			return &ast.TableOption{Tp: ast.TableOptionCollate, StrValue: r.parseCollationName(),
				UintValue: ast.TableOptionCharsetWithoutConvertTo}
		}
		if !r.isCharsetKw() {
			r.syntaxError()
		}
		r.parseCharsetKw()
		r.parseEqOpt()
		return &ast.TableOption{Tp: ast.TableOptionCharset, StrValue: r.parseCharsetName(),
			UintValue: ast.TableOptionCharsetWithoutConvertTo}
	case charsetKwd, character, charType:
		r.parseCharsetKw()
		r.parseEqOpt()
		return &ast.TableOption{Tp: ast.TableOptionCharset, StrValue: r.parseCharsetName(),
			UintValue: ast.TableOptionCharsetWithoutConvertTo}
	case collate:
		r.advance()
		r.parseEqOpt()
		return &ast.TableOption{Tp: ast.TableOptionCollate, StrValue: r.parseCollationName(),
			UintValue: ast.TableOptionCharsetWithoutConvertTo}
	case force:
		// ForceOpt "AUTO_INCREMENT"/"AUTO_RANDOM_BASE" EqOpt LengthNum
		r.advance()
		switch r.tok() {
		case autoIncrement:
			r.advance()
			r.parseEqOpt()
			return &ast.TableOption{Tp: ast.TableOptionAutoIncrement, UintValue: r.parseLengthNum(), BoolValue: true}
		case autoRandomBase:
			r.advance()
			r.parseEqOpt()
			return &ast.TableOption{Tp: ast.TableOptionAutoRandomBase, UintValue: r.parseLengthNum(), BoolValue: true}
		}
		r.syntaxError()
	case autoIncrement:
		r.advance()
		r.parseEqOpt()
		return &ast.TableOption{Tp: ast.TableOptionAutoIncrement, UintValue: r.parseLengthNum(), BoolValue: false}
	case autoRandomBase:
		r.advance()
		r.parseEqOpt()
		return &ast.TableOption{Tp: ast.TableOptionAutoRandomBase, UintValue: r.parseLengthNum(), BoolValue: false}
	case autoIdCache:
		r.advance()
		r.parseEqOpt()
		return &ast.TableOption{Tp: ast.TableOptionAutoIdCache, UintValue: r.parseLengthNum()}
	case avgRowLength:
		r.advance()
		r.parseEqOpt()
		return &ast.TableOption{Tp: ast.TableOptionAvgRowLength, UintValue: r.parseLengthNum()}
	case connection:
		r.advance()
		r.parseEqOpt()
		return &ast.TableOption{Tp: ast.TableOptionConnection, StrValue: r.expect(stringLit).lit}
	case checksum:
		r.advance()
		r.parseEqOpt()
		return &ast.TableOption{Tp: ast.TableOptionCheckSum, UintValue: r.parseLengthNum()}
	case tableChecksum:
		r.advance()
		r.parseEqOpt()
		return &ast.TableOption{Tp: ast.TableOptionTableCheckSum, UintValue: r.parseLengthNum()}
	case password:
		r.advance()
		r.parseEqOpt()
		opt := &ast.TableOption{Tp: ast.TableOptionPassword, StrValue: r.expect(stringLit).lit}
		r.appendWarnf("The PASSWORD option is parsed but ignored by all storage engines.")
		return opt
	case compression:
		r.advance()
		r.parseEqOpt()
		return &ast.TableOption{Tp: ast.TableOptionCompression, StrValue: r.expect(stringLit).lit}
	case keyBlockSize:
		r.advance()
		r.parseEqOpt()
		return &ast.TableOption{Tp: ast.TableOptionKeyBlockSize, UintValue: r.parseLengthNum()}
	case delayKeyWrite:
		r.advance()
		r.parseEqOpt()
		return &ast.TableOption{Tp: ast.TableOptionDelayKeyWrite, UintValue: r.parseLengthNum()}
	case rowFormat:
		return &ast.TableOption{Tp: ast.TableOptionRowFormat, UintValue: r.parseRowFormatValue()}
	case statsPersistent:
		// "STATS_PERSISTENT" EqOpt StatsPersistentVal
		r.advance()
		r.parseEqOpt()
		r.parseStatsPersistentVal()
		return &ast.TableOption{Tp: ast.TableOptionStatsPersistent}
	case statsAutoRecalc:
		r.advance()
		r.parseEqOpt()
		if r.accept(defaultKwd) {
			opt := &ast.TableOption{Tp: ast.TableOptionStatsAutoRecalc, Default: true}
			r.appendWarnf("The STATS_AUTO_RECALC is parsed but ignored by all storage engines.")
			return opt
		}
		n := r.parseLengthNum()
		if n != 0 && n != 1 {
			r.actionError(r.actionErrorf("The value of STATS_AUTO_RECALC must be one of [0|1|DEFAULT]."))
		}
		opt := &ast.TableOption{Tp: ast.TableOptionStatsAutoRecalc, UintValue: n}
		r.appendWarnf("The STATS_AUTO_RECALC is parsed but ignored by all storage engines.")
		return opt
	case statsSamplePages:
		// Parse it but will ignore it.
		// In MySQL, STATS_SAMPLE_PAGES=N(Where 0<N<=65535) or STAS_SAMPLE_PAGES=DEFAULT.
		// Cause we don't support it, so we don't check range of the value.
		r.advance()
		r.parseEqOpt()
		if r.accept(defaultKwd) {
			// In MySQL, default value of STATS_SAMPLE_PAGES is 0.
			opt := &ast.TableOption{Tp: ast.TableOptionStatsSamplePages, Default: true}
			r.appendWarnf("The STATS_SAMPLE_PAGES is parsed but ignored by all storage engines.")
			return opt
		}
		opt := &ast.TableOption{Tp: ast.TableOptionStatsSamplePages, UintValue: r.parseLengthNum()}
		r.appendWarnf("The STATS_SAMPLE_PAGES is parsed but ignored by all storage engines.")
		return opt
	case statsBuckets:
		r.advance()
		r.parseEqOpt()
		return &ast.TableOption{Tp: ast.TableOptionStatsBuckets, UintValue: r.parseLengthNum()}
	case statsTopN:
		r.advance()
		r.parseEqOpt()
		return &ast.TableOption{Tp: ast.TableOptionStatsTopN, UintValue: r.parseLengthNum()}
	case statsSampleRate:
		// "STATS_SAMPLE_RATE" EqOpt NumLiteral
		r.advance()
		r.parseEqOpt()
		return &ast.TableOption{Tp: ast.TableOptionStatsSampleRate, Value: ast.NewValueExpr(r.parseNumLiteralValue(), "", "")}
	case statsColChoice:
		r.advance()
		r.parseEqOpt()
		return &ast.TableOption{Tp: ast.TableOptionStatsColsChoice, StrValue: r.expect(stringLit).lit}
	case statsColList:
		r.advance()
		r.parseEqOpt()
		return &ast.TableOption{Tp: ast.TableOptionStatsColList, StrValue: r.expect(stringLit).lit}
	case shardRowIDBits:
		r.advance()
		r.parseEqOpt()
		return &ast.TableOption{Tp: ast.TableOptionShardRowID, UintValue: r.parseLengthNum()}
	case preSplitRegions:
		r.advance()
		r.parseEqOpt()
		return &ast.TableOption{Tp: ast.TableOptionPreSplitRegion, UintValue: r.parseLengthNum()}
	case packKeys:
		// "PACK_KEYS" EqOpt StatsPersistentVal — parse it but will ignore it.
		r.advance()
		r.parseEqOpt()
		r.parseStatsPersistentVal()
		return &ast.TableOption{Tp: ast.TableOptionPackKeys}
	case storage:
		r.advance()
		switch r.tok() {
		case engine:
			// PartDefOption: "STORAGE" "ENGINE" EqOpt StringName
			r.advance()
			r.parseEqOpt()
			return &ast.TableOption{Tp: ast.TableOptionEngine, StrValue: r.parseStringName()}
		case memory:
			// Parse it but will ignore it.
			r.advance()
			opt := &ast.TableOption{Tp: ast.TableOptionStorageMedia, StrValue: "MEMORY"}
			r.appendWarnf("The STORAGE clause is parsed but ignored by all storage engines.")
			return opt
		case disk:
			// Parse it but will ignore it.
			r.advance()
			opt := &ast.TableOption{Tp: ast.TableOptionStorageMedia, StrValue: "DISK"}
			r.appendWarnf("The STORAGE clause is parsed but ignored by all storage engines.")
			return opt
		}
		r.syntaxError()
	case secondaryEngine:
		// Parse it but will ignore it.
		r.advance()
		r.parseEqOpt()
		if r.accept(null) {
			opt := &ast.TableOption{Tp: ast.TableOptionSecondaryEngineNull}
			r.appendWarnf("The SECONDARY_ENGINE clause is parsed but ignored by all storage engines.")
			return opt
		}
		opt := &ast.TableOption{Tp: ast.TableOptionSecondaryEngine, StrValue: r.parseStringName()}
		r.appendWarnf("The SECONDARY_ENGINE clause is parsed but ignored by all storage engines.")
		return opt
	case union:
		// "UNION" EqOpt '(' TableNameListOpt ')' — parse it but will ignore it.
		r.advance()
		r.parseEqOpt()
		r.expect(int('('))
		var tableNames []*ast.TableName
		if r.tok() == int(')') {
			tableNames = []*ast.TableName{}
		} else {
			tableNames = r.parseTableNameList()
		}
		r.expect(int(')'))
		opt := &ast.TableOption{
			Tp:         ast.TableOptionUnion,
			TableNames: tableNames,
		}
		r.appendWarnf("The UNION option is parsed but ignored by all storage engines.")
		return opt
	case encryption:
		// "ENCRYPTION" EqOpt EncryptionOpt — parse it but will ignore it.
		r.advance()
		r.parseEqOpt()
		opt := &ast.TableOption{Tp: ast.TableOptionEncryption, StrValue: r.parseEncryptionOpt()}
		r.appendWarnf("The ENCRYPTION option is parsed but ignored by all storage engines.")
		return opt
	case ttl:
		// "TTL" EqOpt Identifier '+' "INTERVAL" Literal TimeUnit
		r.advance()
		r.parseEqOpt()
		id := r.parseIdentifier()
		r.expect(int('+'))
		r.expect(interval)
		litStart := r.cur().offset
		lit := r.parseLiteralExpr(litStart)
		return &ast.TableOption{
			Tp:            ast.TableOptionTTL,
			ColumnName:    &ast.ColumnName{Name: ast.NewCIStr(id)},
			Value:         ast.NewValueExpr(lit, r.p.charset, r.p.collation),
			TimeUnitValue: &ast.TimeUnitExpr{Unit: r.parseTimeUnit()},
		}
	case ttlEnable:
		r.advance()
		r.parseEqOpt()
		onOrOff := strings.ToLower(r.expect(stringLit).lit)
		if onOrOff == "on" {
			return &ast.TableOption{Tp: ast.TableOptionTTLEnable, BoolValue: true}
		} else if onOrOff == "off" {
			return &ast.TableOption{Tp: ast.TableOptionTTLEnable, BoolValue: false}
		}
		r.actionError(r.actionErrorf("The TTL_ENABLE option has to be set 'ON' or 'OFF'"))
	case ttlJobInterval:
		r.advance()
		r.parseEqOpt()
		s := r.expect(stringLit).lit
		_, err := duration.ParseDuration(s)
		if err != nil {
			r.actionError(r.actionErrorf("The TTL_JOB_INTERVAL option is not a valid duration: %s", err.Error()))
		}
		return &ast.TableOption{Tp: ast.TableOptionTTLJobInterval, StrValue: s}
	case autoextendSize:
		// Parse it but will ignore it.
		r.advance()
		r.parseEqOpt()
		opt := &ast.TableOption{Tp: ast.TableOptionAutoextendSize, StrValue: r.parseStringName()}
		r.appendWarnf("The AUTOEXTEND_SIZE option is parsed but ignored by all storage engines.")
		return opt
	case affinity:
		r.advance()
		r.parseEqOpt()
		return &ast.TableOption{Tp: ast.TableOptionAffinity, StrValue: r.parseStringName()}
	// MariaDB specific options
	// - https://mariadb.com/docs/server/reference/sql-statements/data-definition/create/create-table
	case pageChecksum:
		r.advance()
		r.parseEqOpt()
		opt := &ast.TableOption{Tp: ast.TableOptionPageChecksum, UintValue: r.parseLengthNum()}
		r.appendWarnf("The PAGE_CHECKSUM option is parsed but ignored by all storage engines.")
		return opt
	case pageCompressed:
		r.advance()
		r.parseEqOpt()
		opt := &ast.TableOption{Tp: ast.TableOptionPageCompressed, UintValue: r.parseLengthNum()}
		r.appendWarnf("The PAGE_COMPRESSED option is parsed but ignored by all storage engines.")
		return opt
	case pageCompressionLevel:
		r.advance()
		r.parseEqOpt()
		opt := &ast.TableOption{Tp: ast.TableOptionPageCompressionLevel, UintValue: r.parseLengthNum()}
		r.appendWarnf("The PAGE_COMPRESSION_LEVEL option is parsed but ignored by all storage engines.")
		return opt
	case transactional:
		r.advance()
		r.parseEqOpt()
		opt := &ast.TableOption{Tp: ast.TableOptionTransactional, UintValue: r.parseLengthNum()}
		r.appendWarnf("The TRANSACTIONAL option is parsed but ignored by all storage engines.")
		return opt
	case sequence:
		r.advance()
		r.parseEqOpt()
		opt := &ast.TableOption{Tp: ast.TableOptionSequence, UintValue: r.parseLengthNum()}
		r.appendWarnf("The SEQUENCE option is parsed but ignored by all storage engines. Use CREATE SEQUENCE instead.")
		return opt
	case ietfQuotes:
		r.advance()
		r.parseEqOpt()
		opt := &ast.TableOption{Tp: ast.TableOptionIetfQuotes, StrValue: r.parseStringName()}
		r.appendWarnf("The IETF_QUOTES option is parsed but ignored by all storage engines.")
		return opt
	default:
		return r.parsePartDefOption()
	}
	return nil
}

// parseRowFormatValue implements RowFormat.
func (r *rdParser) parseRowFormatValue() uint64 {
	r.expect(rowFormat)
	r.parseEqOpt()
	var v uint64
	switch r.tok() {
	case defaultKwd:
		v = ast.RowFormatDefault
	case dynamic:
		v = ast.RowFormatDynamic
	case fixed:
		v = ast.RowFormatFixed
	case compressed:
		v = ast.RowFormatCompressed
	case redundant:
		v = ast.RowFormatRedundant
	case compact:
		v = ast.RowFormatCompact
	case tokudbDefault:
		v = ast.TokuDBRowFormatDefault
	case tokudbFast:
		v = ast.TokuDBRowFormatFast
	case tokudbSmall:
		v = ast.TokuDBRowFormatSmall
	case tokudbZlib:
		v = ast.TokuDBRowFormatZlib
	case tokudbZstd:
		v = ast.TokuDBRowFormatZstd
	case tokudbQuickLZ:
		v = ast.TokuDBRowFormatQuickLZ
	case tokudbLzma:
		v = ast.TokuDBRowFormatLzma
	case tokudbSnappy:
		v = ast.TokuDBRowFormatSnappy
	case tokudbUncompressed:
		v = ast.TokuDBRowFormatUncompressed
	default:
		r.syntaxError()
	}
	r.advance()
	return v
}

// parseStatsPersistentVal consumes StatsPersistentVal (value ignored).
func (r *rdParser) parseStatsPersistentVal() {
	if r.accept(defaultKwd) {
		return
	}
	r.parseLengthNum()
}

// parseEncryptionOpt implements EncryptionOpt.
func (r *rdParser) parseEncryptionOpt() string {
	// Parse it but will ignore it
	s := r.expect(stringLit).lit
	switch s {
	case "Y", "y":
		r.appendWarnf("The ENCRYPTION clause is parsed but ignored by all storage engines.")
	case "N", "n":
	default:
		r.actionError(ErrWrongValue.GenWithStackByArgs("argument (should be Y or N)", s))
	}
	return s
}

// isPartDefOptionStart reports whether the current token can begin a
// PartDefOption.
func (r *rdParser) isPartDefOptionStart() bool {
	switch r.tok() {
	case comment, engine, storage, engine_attribute, secondaryEngineAttribute,
		insertMethod, data, index, maxRows, minRows, tablespace, nodegroup,
		placement:
		return true
	}
	return false
}

// parsePartDefOption implements PartDefOption.
func (r *rdParser) parsePartDefOption() *ast.TableOption {
	switch r.tok() {
	case comment:
		// "COMMENT" EqOpt stringLit
		r.advance()
		r.parseEqOpt()
		return &ast.TableOption{Tp: ast.TableOptionComment, StrValue: r.expect(stringLit).lit}
	case engine:
		// "ENGINE" EqOpt StringName
		r.advance()
		r.parseEqOpt()
		return &ast.TableOption{Tp: ast.TableOptionEngine, StrValue: r.parseStringName()}
	case storage:
		// "STORAGE" "ENGINE" EqOpt StringName
		r.advance()
		r.expect(engine)
		r.parseEqOpt()
		return &ast.TableOption{Tp: ast.TableOptionEngine, StrValue: r.parseStringName()}
	case engine_attribute:
		// "ENGINE_ATTRIBUTE" EqOpt StringName
		r.advance()
		r.parseEqOpt()
		return &ast.TableOption{Tp: ast.TableOptionEngineAttribute, StrValue: r.parseStringName()}
	case secondaryEngineAttribute:
		// "SECONDARY_ENGINE_ATTRIBUTE" EqOpt stringLit
		r.advance()
		r.parseEqOpt()
		return &ast.TableOption{
			Tp:       ast.TableOptionSecondaryEngineAttribute,
			StrValue: r.expect(stringLit).lit,
		}
	case insertMethod:
		// "INSERT_METHOD" EqOpt StringName
		r.advance()
		r.parseEqOpt()
		return &ast.TableOption{Tp: ast.TableOptionInsertMethod, StrValue: r.parseStringName()}
	case data:
		// "DATA" "DIRECTORY" EqOpt stringLit
		r.advance()
		r.expect(directory)
		r.parseEqOpt()
		return &ast.TableOption{Tp: ast.TableOptionDataDirectory, StrValue: r.expect(stringLit).lit}
	case index:
		// "INDEX" "DIRECTORY" EqOpt stringLit
		r.advance()
		r.expect(directory)
		r.parseEqOpt()
		return &ast.TableOption{Tp: ast.TableOptionIndexDirectory, StrValue: r.expect(stringLit).lit}
	case maxRows:
		r.advance()
		r.parseEqOpt()
		return &ast.TableOption{Tp: ast.TableOptionMaxRows, UintValue: r.parseLengthNum()}
	case minRows:
		r.advance()
		r.parseEqOpt()
		return &ast.TableOption{Tp: ast.TableOptionMinRows, UintValue: r.parseLengthNum()}
	case tablespace:
		// "TABLESPACE" EqOpt Identifier
		r.advance()
		r.parseEqOpt()
		return &ast.TableOption{Tp: ast.TableOptionTablespace, StrValue: r.parseIdentifier()}
	case nodegroup:
		r.advance()
		r.parseEqOpt()
		return &ast.TableOption{Tp: ast.TableOptionNodegroup, UintValue: r.parseLengthNum()}
	case placement:
		placementOptions := r.parsePlacementPolicyOption()
		return &ast.TableOption{
			// offset trick, enums are identical but of different type
			Tp:        ast.TableOptionType(placementOptions.Tp),
			StrValue:  placementOptions.StrValue,
			UintValue: placementOptions.UintValue,
		}
	}
	r.syntaxError()
	return nil
}

// parsePlacementPolicyOption implements PlacementPolicyOption.
func (r *rdParser) parsePlacementPolicyOption() *ast.PlacementOption {
	r.expect(placement)
	r.expect(policy)
	if r.tok() == set {
		// "PLACEMENT" "POLICY" "SET" "DEFAULT"
		r.advance()
		t := r.expect(defaultKwd)
		return &ast.PlacementOption{Tp: ast.PlacementOptionPolicy, StrValue: t.lit}
	}
	r.parseEqOpt()
	switch r.tok() {
	case stringLit:
		s := r.cur().lit
		r.advance()
		return &ast.PlacementOption{Tp: ast.PlacementOptionPolicy, StrValue: s}
	case defaultKwd:
		s := r.cur().lit
		r.advance()
		return &ast.PlacementOption{Tp: ast.PlacementOptionPolicy, StrValue: s}
	}
	// PolicyName: Identifier
	return &ast.PlacementOption{Tp: ast.PlacementOptionPolicy, StrValue: r.parseIdentifier()}
}

// parsePartitionOpt implements PartitionOpt.
func (r *rdParser) parsePartitionOpt() *ast.PartitionOptions {
	if r.tok() != partition {
		return nil
	}
	// "PARTITION" "BY" PartitionMethod PartitionNumOpt SubPartitionOpt
	// PartitionDefinitionListOpt UpdateIndexesOpt
	r.advance()
	r.expect(by)
	method := r.parsePartitionMethod()
	method.Num = r.parsePartitionNumOpt()
	sub := r.parseSubPartitionOpt()
	defs := r.parsePartitionDefinitionListOpt()
	updateIndexes := r.parseUpdateIndexesOpt()
	opt := &ast.PartitionOptions{
		PartitionMethod: *method,
		Sub:             sub,
		Definitions:     defs,
		UpdateIndexes:   updateIndexes,
	}
	if err := opt.Validate(); err != nil {
		// yacc runs Validate at the PartitionOpt reduce, which only
		// happens when the lookahead can follow the partition clause;
		// otherwise the syntax error at the lookahead comes first.
		if !partitionOptFollow(r.tok()) {
			r.syntaxError()
		}
		r.actionError(err)
	}
	return opt
}

// partitionOptFollow approximates FOLLOW(PartitionOpt) across its two
// contexts (CREATE TABLE tail and ALTER TABLE spec list).
func partitionOptFollow(tok int) bool {
	switch tok {
	case ignore, replace, as, selectKwd, with, int('('), tableKwd, values,
		on, int(';'), 0, int(','):
		return true
	}
	return false
}

// parsePartitionMethod implements PartitionMethod.
func (r *rdParser) parsePartitionMethod() *ast.PartitionMethod {
	switch r.tok() {
	case linear, key, hash:
		return r.parseSubPartitionMethod()
	case rangeKwd:
		r.advance()
		if r.tok() == fields || r.tok() == columns {
			// "RANGE" FieldsOrColumns '(' ColumnNameList ')' PartitionIntervalOpt
			r.advance()
			r.expect(int('('))
			cols := r.parseColumnNameList()
			r.expect(int(')'))
			return &ast.PartitionMethod{
				Tp:          ast.PartitionTypeRange,
				ColumnNames: cols,
				Interval:    r.parsePartitionIntervalOpt(),
			}
		}
		// "RANGE" '(' BitExpr ')' PartitionIntervalOpt
		r.expect(int('('))
		expr := r.parseBitExpr(0)
		r.expect(int(')'))
		return &ast.PartitionMethod{
			Tp:       ast.PartitionTypeRange,
			Expr:     expr,
			Interval: r.parsePartitionIntervalOpt(),
		}
	case list:
		r.advance()
		if r.tok() == fields || r.tok() == columns {
			// "LIST" FieldsOrColumns '(' ColumnNameList ')'
			r.advance()
			r.expect(int('('))
			cols := r.parseColumnNameList()
			r.expect(int(')'))
			return &ast.PartitionMethod{
				Tp:          ast.PartitionTypeList,
				ColumnNames: cols,
			}
		}
		// "LIST" '(' BitExpr ')'
		r.expect(int('('))
		expr := r.parseBitExpr(0)
		r.expect(int(')'))
		return &ast.PartitionMethod{
			Tp:   ast.PartitionTypeList,
			Expr: expr,
		}
	case systemTime:
		r.advance()
		switch r.tok() {
		case interval:
			// "SYSTEM_TIME" "INTERVAL" Expression TimeUnit
			r.advance()
			expr := r.parseExpression()
			return &ast.PartitionMethod{
				Tp:   ast.PartitionTypeSystemTime,
				Expr: expr,
				Unit: r.parseTimeUnit(),
			}
		case limit:
			// "SYSTEM_TIME" "LIMIT" LengthNum
			r.advance()
			return &ast.PartitionMethod{
				Tp:    ast.PartitionTypeSystemTime,
				Limit: r.parseLengthNum(),
			}
		}
		return &ast.PartitionMethod{
			Tp: ast.PartitionTypeSystemTime,
		}
	}
	r.syntaxError()
	return nil
}

// parseSubPartitionMethod implements SubPartitionMethod (with LinearOpt
// and PartitionKeyAlgorithmOpt).
func (r *rdParser) parseSubPartitionMethod() *ast.PartitionMethod {
	// LinearOpt
	linearOpt := ""
	if r.tok() == linear {
		linearOpt = r.cur().lit
		r.advance()
	}
	switch r.tok() {
	case key:
		// LinearOpt "KEY" PartitionKeyAlgorithmOpt '(' ColumnNameListOpt ')'
		r.advance()
		keyAlgorithm := r.parsePartitionKeyAlgorithmOpt()
		r.expect(int('('))
		var cols []*ast.ColumnName
		if r.tok() == int(')') {
			cols = []*ast.ColumnName{}
		} else {
			cols = r.parseColumnNameList()
		}
		r.expect(int(')'))
		return &ast.PartitionMethod{
			Tp:           ast.PartitionTypeKey,
			Linear:       len(linearOpt) != 0,
			ColumnNames:  cols,
			KeyAlgorithm: keyAlgorithm,
		}
	case hash:
		// LinearOpt "HASH" '(' BitExpr ')'
		r.advance()
		r.expect(int('('))
		expr := r.parseBitExpr(0)
		r.expect(int(')'))
		return &ast.PartitionMethod{
			Tp:     ast.PartitionTypeHash,
			Linear: len(linearOpt) != 0,
			Expr:   expr,
		}
	}
	r.syntaxError()
	return nil
}

// parsePartitionKeyAlgorithmOpt implements PartitionKeyAlgorithmOpt.
func (r *rdParser) parsePartitionKeyAlgorithmOpt() *ast.PartitionKeyAlgorithm {
	if r.tok() != algorithm {
		return nil
	}
	// "ALGORITHM" eq NUM
	r.advance()
	r.expect(eq)
	tp := getUint64FromNUM(r.expect(intLit).item)
	if tp != 1 && tp != 2 {
		r.actionError(ErrSyntax)
	}
	return &ast.PartitionKeyAlgorithm{
		Type: tp,
	}
}

// parsePartitionIntervalOpt implements PartitionIntervalOpt (with
// IntervalExpr, FirstAndLastPartOpt, NullPartOpt, and MaxValPartOpt).
func (r *rdParser) parsePartitionIntervalOpt() *ast.PartitionInterval {
	if r.tok() != interval {
		return nil
	}
	// "INTERVAL" '(' IntervalExpr ')' FirstAndLastPartOpt NullPartOpt MaxValPartOpt
	start := r.cur().offset
	r.advance()
	r.expect(int('('))
	// IntervalExpr: BitExpr | BitExpr TimeUnit
	expr := r.parseBitExpr(0)
	timeUnit := ast.TimeUnitInvalid
	if _, isUnit := timeUnitFor(r.tok()); isUnit {
		timeUnit = r.parseTimeUnit()
	}
	r.expect(int(')'))
	// FirstAndLastPartOpt: First/LastRangeEnd defaults to nil
	firstAndLast := ast.PartitionInterval{}
	if r.tok() == first {
		// "FIRST" "PARTITION" "LESS" "THAN" '(' BitExpr ')'
		// "LAST" "PARTITION" "LESS" "THAN" '(' BitExpr ')'
		r.advance()
		r.expect(partition)
		r.expect(less)
		r.expect(than)
		r.expect(int('('))
		firstExpr := r.parseBitExpr(0)
		r.expect(int(')'))
		r.expect(last)
		r.expect(partition)
		r.expect(less)
		r.expect(than)
		r.expect(int('('))
		lastExpr := r.parseBitExpr(0)
		r.expect(int(')'))
		firstAndLast = ast.PartitionInterval{
			FirstRangeEnd: &firstExpr,
			LastRangeEnd:  &lastExpr,
		}
	}
	// NullPartOpt
	nullPart := false
	if r.tok() == null {
		r.advance()
		r.expect(partition)
		nullPart = true
	}
	// MaxValPartOpt
	maxValPart := false
	if r.tok() == maxValue {
		r.advance()
		r.expect(partition)
		maxValPart = true
	}
	// The action's node text runs from the "INTERVAL" token
	// (parser.yyVAL.offset) to the unconsumed lookahead
	// (parser.yylval.offset).
	end := r.cur().offset
	partitionInterval := &ast.PartitionInterval{
		IntervalExpr:  ast.PartitionIntervalExpr{Expr: expr, TimeUnit: timeUnit},
		FirstRangeEnd: firstAndLast.FirstRangeEnd,
		LastRangeEnd:  firstAndLast.LastRangeEnd,
		NullPart:      nullPart,
		MaxValPart:    maxValPart,
	}
	r.p.setNodeText(partitionInterval, r.src[start:end])
	return partitionInterval
}

// parsePartitionNumOpt implements PartitionNumOpt.
func (r *rdParser) parsePartitionNumOpt() uint64 {
	if r.tok() != partitions {
		return uint64(0)
	}
	r.advance()
	res := r.parseLengthNum()
	if res == 0 {
		r.actionError(ast.ErrNoParts.GenWithStackByArgs("partitions"))
	}
	return res
}

// parseSubPartitionOpt implements SubPartitionOpt.
func (r *rdParser) parseSubPartitionOpt() *ast.PartitionMethod {
	if r.tok() != subpartition {
		return nil
	}
	// "SUBPARTITION" "BY" SubPartitionMethod SubPartitionNumOpt
	r.advance()
	r.expect(by)
	method := r.parseSubPartitionMethod()
	method.Num = r.parseSubPartitionNumOpt()
	return method
}

// parseSubPartitionNumOpt implements SubPartitionNumOpt.
func (r *rdParser) parseSubPartitionNumOpt() uint64 {
	if r.tok() != subpartitions {
		return uint64(0)
	}
	r.advance()
	res := r.parseLengthNum()
	if res == 0 {
		r.actionError(ast.ErrNoParts.GenWithStackByArgs("subpartitions"))
	}
	return res
}

// parsePartitionDefinitionListOpt implements PartitionDefinitionListOpt
// and PartitionDefinitionList.
func (r *rdParser) parsePartitionDefinitionListOpt() []*ast.PartitionDefinition {
	if r.tok() != int('(') {
		return nil
	}
	r.advance()
	defs := []*ast.PartitionDefinition{r.parsePartitionDefinition()}
	for r.accept(int(',')) {
		defs = append(defs, r.parsePartitionDefinition())
	}
	r.expect(int(')'))
	return defs
}

// parsePartitionDefinition implements PartitionDefinition.
func (r *rdParser) parsePartitionDefinition() *ast.PartitionDefinition {
	// "PARTITION" Identifier PartDefValuesOpt PartDefOptionList SubPartDefinitionListOpt
	r.expect(partition)
	name := r.parseIdentifier()
	clause := r.parsePartDefValuesOpt()
	options := r.parsePartDefOptionList()
	return &ast.PartitionDefinition{
		Name:    ast.NewCIStr(name),
		Clause:  clause,
		Options: options,
		Sub:     r.parseSubPartDefinitionListOpt(),
	}
}

// parsePartDefOptionList implements PartDefOptionList.
func (r *rdParser) parsePartDefOptionList() []*ast.TableOption {
	list := make([]*ast.TableOption, 0)
	for r.isPartDefOptionStart() {
		list = append(list, r.parsePartDefOption())
	}
	return list
}

// parsePartDefValuesOpt implements PartDefValuesOpt.
func (r *rdParser) parsePartDefValuesOpt() ast.PartitionDefinitionClause {
	switch r.tok() {
	case values:
		r.advance()
		if r.accept(less) {
			r.expect(than)
			if r.tok() == maxValue {
				// "VALUES" "LESS" "THAN" "MAXVALUE"
				r.advance()
				return &ast.PartitionDefinitionClauseLessThan{
					Exprs: []ast.ExprNode{&ast.MaxValueExpr{}},
				}
			}
			// "VALUES" "LESS" "THAN" '(' MaxValueOrExpressionList ')'
			r.expect(int('('))
			exprs := r.parseMaxValueOrExpressionList()
			r.expect(int(')'))
			return &ast.PartitionDefinitionClauseLessThan{
				Exprs: exprs,
			}
		}
		// "VALUES" "IN" '(' DefaultOrExpressionList ')'
		r.expect(in)
		r.expect(int('('))
		exprs := r.parseDefaultOrExpressionList()
		r.expect(int(')'))
		values := make([][]ast.ExprNode, 0, len(exprs))
		for _, expr := range exprs {
			if row, ok := expr.(*ast.RowExpr); ok {
				values = append(values, row.Values)
			} else {
				values = append(values, []ast.ExprNode{expr})
			}
		}
		return &ast.PartitionDefinitionClauseIn{Values: values}
	case defaultKwd:
		r.advance()
		return &ast.PartitionDefinitionClauseIn{
			Values: [][]ast.ExprNode{{&ast.DefaultExpr{}}},
		}
	case history:
		r.advance()
		return &ast.PartitionDefinitionClauseHistory{Current: false}
	case current:
		r.advance()
		return &ast.PartitionDefinitionClauseHistory{Current: true}
	}
	return &ast.PartitionDefinitionClauseNone{}
}

// parseMaxValueOrExpressionList implements MaxValueOrExpressionList.
func (r *rdParser) parseMaxValueOrExpressionList() []ast.ExprNode {
	exprs := []ast.ExprNode{r.parseMaxValueOrExpression()}
	for r.accept(int(',')) {
		exprs = append(exprs, r.parseMaxValueOrExpression())
	}
	return exprs
}

// parseMaxValueOrExpression implements MaxValueOrExpression.
func (r *rdParser) parseMaxValueOrExpression() ast.ExprNode {
	if r.tok() == maxValue {
		start := r.cur().offset
		r.advance()
		return r.setOrigin(&ast.MaxValueExpr{}, start)
	}
	return r.parseBitExpr(0)
}

// parseDefaultOrExpressionList implements DefaultOrExpressionList.
func (r *rdParser) parseDefaultOrExpressionList() []ast.ExprNode {
	exprs := []ast.ExprNode{r.parseDefaultOrExpression()}
	for r.accept(int(',')) {
		exprs = append(exprs, r.parseDefaultOrExpression())
	}
	return exprs
}

// parseDefaultOrExpression implements DefaultOrExpression.
func (r *rdParser) parseDefaultOrExpression() ast.ExprNode {
	if r.tok() == defaultKwd {
		start := r.cur().offset
		r.advance()
		return r.setOrigin(&ast.DefaultExpr{}, start)
	}
	return r.parseBitExpr(0)
}

// parseSubPartDefinitionListOpt implements SubPartDefinitionListOpt and
// SubPartDefinitionList.
func (r *rdParser) parseSubPartDefinitionListOpt() []*ast.SubPartitionDefinition {
	list := make([]*ast.SubPartitionDefinition, 0)
	if r.tok() != int('(') {
		return list
	}
	r.advance()
	list = append(list, r.parseSubPartDefinition())
	for r.accept(int(',')) {
		list = append(list, r.parseSubPartDefinition())
	}
	r.expect(int(')'))
	return list
}

// parseSubPartDefinition implements SubPartDefinition.
func (r *rdParser) parseSubPartDefinition() *ast.SubPartitionDefinition {
	// "SUBPARTITION" Identifier PartDefOptionList
	r.expect(subpartition)
	name := r.parseIdentifier()
	return &ast.SubPartitionDefinition{
		Name:    ast.NewCIStr(name),
		Options: r.parsePartDefOptionList(),
	}
}

// parseUpdateIndexesOpt implements UpdateIndexesOpt and UpdateIndexesList.
func (r *rdParser) parseUpdateIndexesOpt() []*ast.Constraint {
	if r.tok() != update {
		return nil
	}
	// "UPDATE" "INDEXES" '(' UpdateIndexesList ')'
	r.advance()
	r.expect(indexes)
	r.expect(int('('))
	list := []*ast.Constraint{r.parseUpdateIndexElem()}
	for r.accept(int(',')) {
		list = append(list, r.parseUpdateIndexElem())
	}
	r.expect(int(')'))
	return list
}

// parseUpdateIndexElem implements UpdateIndexElem: Identifier GlobalOrLocal.
func (r *rdParser) parseUpdateIndexElem() *ast.Constraint {
	name := r.parseIdentifier()
	var g bool
	switch r.tok() {
	case local:
		g = false
	case global:
		g = true
	default:
		r.syntaxError()
	}
	r.advance()
	opt := &ast.IndexOption{Global: g}
	return &ast.Constraint{
		Name:   name,
		Option: opt,
	}
}

// parseSplitIndexListOpt implements SplitIndexListOpt and SplitIndexList.
func (r *rdParser) parseSplitIndexListOpt() []*ast.SplitIndexOption {
	if r.tok() != split {
		return nil
	}
	list := []*ast.SplitIndexOption{r.parseSplitIndexOption()}
	for r.tok() == split {
		list = append(list, r.parseSplitIndexOption())
	}
	return list
}

// parseSplitIndexOption implements SplitIndexOption.
func (r *rdParser) parseSplitIndexOption() *ast.SplitIndexOption {
	r.expect(split)
	switch r.tok() {
	case primary:
		// "SPLIT" "PRIMARY" "KEY" SplitOptionBetween
		r.advance()
		r.expect(key)
		return &ast.SplitIndexOption{
			PrimaryKey: true,
			IndexName:  ast.NewCIStr(mysql.PrimaryKeyName),
			SplitOpt:   r.parseSplitOptionBetween(),
		}
	case index:
		// "SPLIT" "INDEX" Identifier SplitOptionBetween
		r.advance()
		name := r.parseIdentifier()
		return &ast.SplitIndexOption{
			IndexName: ast.NewCIStr(name),
			SplitOpt:  r.parseSplitOptionBetween(),
		}
	}
	// "SPLIT" SplitOptionBetween
	return &ast.SplitIndexOption{
		TableLevel: true,
		SplitOpt:   r.parseSplitOptionBetween(),
	}
}
