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

// ALTER statements: AlterTableStmt (spec list, ANALYZE and COMPACT forms)
// with the whole AlterTableSpec production (LockClause, AlgorithmClause,
// ColumnPosition, AlterOrderList, ...), AlterDatabaseStmt,
// AlterSequenceStmt, AlterInstanceStmt, AlterRangeStmt, AlterPolicyStmt
// (placement policy), AlterResourceGroupStmt, and AlterUserStmt. The
// option catalogues shared with the CREATE-side ports (DatabaseOption,
// SequenceOption, PlacementOptionList, ResourceGroupOptionList, the user
// account options, masking policy options) live in parse_create_misc.go.

import (
	"fmt"
	"strings"

	"github.com/sqlc-dev/marino/ast"
	"github.com/sqlc-dev/marino/mysql"
	"github.com/sqlc-dev/marino/types"
)

func init() {
	rdRegister(alter, (*rdParser).parseAlterStmtFamily)
}

// parseAlterStmtFamily dispatches the ALTER statement family on the token
// after "ALTER": AlterTableStmt | AlterDatabaseStmt | AlterUserStmt |
// AlterInstanceStmt | AlterRangeStmt | AlterSequenceStmt | AlterPolicyStmt
// | AlterResourceGroupStmt.
func (r *rdParser) parseAlterStmtFamily() ast.StmtNode {
	switch r.la(1) {
	case ignore, tableKwd:
		return r.parseAlterTableStmt()
	case database:
		return r.parseAlterDatabaseStmt()
	case user:
		return r.parseAlterUserStmt()
	case instance:
		return r.parseAlterInstanceStmt()
	case rangeKwd:
		return r.parseAlterRangeStmt()
	case sequence:
		return r.parseAlterSequenceStmt()
	case placement:
		return r.parseAlterPolicyStmt()
	case resource:
		return r.parseAlterResourceGroupStmt()
	default:
		r.unsupported(fmt.Sprintf("ALTER %s", r.at(r.i+1).lit))
		return nil
	}
}

// parseNoWriteToBinLogAliasOpt implements NoWriteToBinLogAliasOpt:
// empty | "NO_WRITE_TO_BINLOG" | "LOCAL".
func (r *rdParser) parseNoWriteToBinLogAliasOpt() bool {
	if r.tok() == noWriteToBinLog || r.tok() == local {
		r.advance()
		return true
	}
	return false
}

// parsePartitionNameList implements PartitionNameList: Identifier |
// PartitionNameList ',' Identifier. The comma continuations are consumed
// greedily, matching the shift preference of ',' over the list-ending
// reductions (%prec lowerThanComma) in every use of the production.
func (r *rdParser) parsePartitionNameList() []ast.CIStr {
	names := []ast.CIStr{ast.NewCIStr(r.parseIdentifier())}
	for r.accept(int(',')) {
		names = append(names, ast.NewCIStr(r.parseIdentifier()))
	}
	return names
}

// parseAllOrPartitionNameList implements AllOrPartitionNameList:
// "ALL" (nil) | PartitionNameList.
func (r *rdParser) parseAllOrPartitionNameList() []ast.CIStr {
	if r.accept(all) {
		return nil
	}
	return r.parsePartitionNameList()
}

/**************************************AlterTableStmt***************************************/

// parseAlterTableStmt implements AlterTableStmt: the
// AlterTableSpecListOpt AlterTableSpecSingleOpt form and the
// ANALYZE PARTITION / COMPACT alternatives.
func (r *rdParser) parseAlterTableStmt() ast.StmtNode {
	r.expect(alter)
	// IgnoreOptional (the flag is unused by the actions).
	r.accept(ignore)
	r.expect(tableKwd)
	table := r.parseTableName()

	switch r.tok() {
	case analyze:
		// "ALTER" IgnoreOptional "TABLE" TableName "ANALYZE" "PARTITION"
		// PartitionNameList ["INDEX" IndexNameList] AnalyzeOptionListOpt
		r.advance()
		r.expect(partition)
		partitionNames := r.parsePartitionNameList()
		if r.accept(index) {
			indexNames := r.parseIndexNameList()
			return &ast.AnalyzeTableStmt{
				TableNames:     []*ast.TableName{table},
				PartitionNames: partitionNames,
				IndexNames:     indexNames,
				IndexFlag:      true,
				AnalyzeOpts:    r.parseAnalyzeOptionListOpt(),
			}
		}
		return &ast.AnalyzeTableStmt{TableNames: []*ast.TableName{table}, PartitionNames: partitionNames, AnalyzeOpts: r.parseAnalyzeOptionListOpt()}
	case compact:
		// "ALTER" IgnoreOptional "TABLE" TableName "COMPACT"
		// [ "PARTITION" PartitionNameList ] [ "TIFLASH" "REPLICA" ]
		r.advance()
		if r.accept(partition) {
			partitionNames := r.parsePartitionNameList()
			if r.accept(tiFlash) {
				r.expect(replica)
				return &ast.CompactTableStmt{
					Table:          table,
					PartitionNames: partitionNames,
					ReplicaKind:    ast.CompactReplicaKindTiFlash,
				}
			}
			return &ast.CompactTableStmt{
				Table:          table,
				PartitionNames: partitionNames,
				ReplicaKind:    ast.CompactReplicaKindAll,
			}
		}
		if r.accept(tiFlash) {
			r.expect(replica)
			return &ast.CompactTableStmt{
				Table:       table,
				ReplicaKind: ast.CompactReplicaKindTiFlash,
			}
		}
		return &ast.CompactTableStmt{
			Table:       table,
			ReplicaKind: ast.CompactReplicaKindAll,
		}
	}

	// AlterTableSpecListOpt AlterTableSpecSingleOpt
	specs := make([]*ast.AlterTableSpec, 0, 1)
	if r.isAlterTableSpecStart() {
		specs = append(specs, r.parseAlterTableSpec())
		for r.accept(int(',')) {
			specs = append(specs, r.parseAlterTableSpec())
		}
	}
	if single := r.parseAlterTableSpecSingleOpt(); single != nil {
		specs = append(specs, single)
	}
	return &ast.AlterTableStmt{
		Table: table,
		Specs: specs,
	}
}

// parseIndexNameList implements IndexNameList (which may be empty):
// empty | Identifier | "PRIMARY" | IndexNameList ',' (Identifier|"PRIMARY").
func (r *rdParser) parseIndexNameList() []ast.CIStr {
	var nameList []ast.CIStr
	if r.tok() != primary && !isIdentifierTok(r.tok()) {
		return nameList
	}
	nameList = append(nameList, r.parseIndexNameListItem())
	for r.accept(int(',')) {
		nameList = append(nameList, r.parseIndexNameListItem())
	}
	return nameList
}

func (r *rdParser) parseIndexNameListItem() ast.CIStr {
	if r.tok() == primary {
		lit := r.cur().lit
		r.advance()
		return ast.NewCIStr(lit)
	}
	return ast.NewCIStr(r.parseIdentifier())
}

// parseAnalyzeOptionListOpt implements AnalyzeOptionListOpt and
// AnalyzeOptionList (options separated by nothing or by a comma).
func (r *rdParser) parseAnalyzeOptionListOpt() []ast.AnalyzeOpt {
	if !r.accept(with) {
		return []ast.AnalyzeOpt{}
	}
	opts := []ast.AnalyzeOpt{r.parseAnalyzeOption()}
	for {
		if r.accept(int(',')) {
			opts = append(opts, r.parseAnalyzeOption())
			continue
		}
		switch r.tok() {
		case intLit, decLit, floatLit:
			opts = append(opts, r.parseAnalyzeOption())
			continue
		}
		return opts
	}
}

// parseAnalyzeOption implements AnalyzeOption.
func (r *rdParser) parseAnalyzeOption() ast.AnalyzeOpt {
	switch r.tok() {
	case intLit:
		v := r.cur().item
		r.advance()
		switch r.tok() {
		case buckets:
			r.advance()
			return ast.AnalyzeOpt{Type: ast.AnalyzeOptNumBuckets, Value: ast.NewValueExpr(v, "", "")}
		case topn:
			r.advance()
			return ast.AnalyzeOpt{Type: ast.AnalyzeOptNumTopN, Value: ast.NewValueExpr(v, "", "")}
		case cmSketch:
			r.advance()
			switch r.tok() {
			case depth:
				r.advance()
				return ast.AnalyzeOpt{Type: ast.AnalyzeOptCMSketchDepth, Value: ast.NewValueExpr(v, "", "")}
			case width:
				r.advance()
				return ast.AnalyzeOpt{Type: ast.AnalyzeOptCMSketchWidth, Value: ast.NewValueExpr(v, "", "")}
			}
		case samples:
			r.advance()
			return ast.AnalyzeOpt{Type: ast.AnalyzeOptNumSamples, Value: ast.NewValueExpr(v, "", "")}
		case sampleRate:
			// NumLiteral "SAMPLERATE"
			r.advance()
			return ast.AnalyzeOpt{Type: ast.AnalyzeOptSampleRate, Value: ast.NewValueExpr(v, "", "")}
		case ndvRate:
			// NumLiteral "NDVRATE"
			r.advance()
			return ast.AnalyzeOpt{Type: ast.AnalyzeOptNDVRate, Value: ast.NewValueExpr(v, "", "")}
		}
	case decLit, floatLit:
		v := r.cur().item
		r.advance()
		switch r.tok() {
		case sampleRate:
			r.advance()
			return ast.AnalyzeOpt{Type: ast.AnalyzeOptSampleRate, Value: ast.NewValueExpr(v, "", "")}
		case ndvRate:
			r.advance()
			return ast.AnalyzeOpt{Type: ast.AnalyzeOptNDVRate, Value: ast.NewValueExpr(v, "", "")}
		}
	}
	r.syntaxError()
	return ast.AnalyzeOpt{}
}

// isAlterTableSpecStart reports whether the current token can begin an
// AlterTableSpec (but not the AlterTableSpecSingleOpt alternatives, whose
// leading tokens PARTITION/REMOVE/REORGANIZE/SPLIT/MERGE never begin an
// AlterTableSpec).
func (r *rdParser) isAlterTableSpecStart() bool {
	switch r.tok() {
	case set, convert, add, last, enable, disable, drop, modify, change,
		attributes, statsOptions, check, coalesce, first, exchange, truncate,
		optimize, repair, importKwd, discard, rebuild, order, alter, rename,
		lock, read, algorithm, force, with, without, secondaryLoad,
		secondaryUnload, cache, nocache:
		return true
	}
	return r.isTableOptionStart()
}

// parseAlterTableSpec implements AlterTableSpec (all alternatives except
// the AlterTableSpecSingleOpt ones).
func (r *rdParser) parseAlterTableSpec() *ast.AlterTableSpec {
	switch r.tok() {
	case set:
		// "SET" ["HYPO"] "TIFLASH" "REPLICA" LengthNum LocationLabelList
		r.advance()
		hypoFlag := r.accept(hypo)
		r.expect(tiFlash)
		r.expect(replica)
		tiflashReplicaSpec := &ast.TiFlashReplicaSpec{
			Count:  r.parseLengthNum(),
			Labels: r.parseLocationLabelList(),
			Hypo:   hypoFlag,
		}
		return &ast.AlterTableSpec{
			Tp:             ast.AlterTableSetTiFlashReplica,
			TiFlashReplica: tiflashReplicaSpec,
		}
	case convert:
		// "CONVERT" "TO" CharsetKw (CharsetName|"DEFAULT") OptCollate
		r.advance()
		r.expect(to)
		r.parseCharsetKw()
		var op *ast.AlterTableSpec
		if r.accept(defaultKwd) {
			op = &ast.AlterTableSpec{
				Tp: ast.AlterTableOption,
				Options: []*ast.TableOption{{Tp: ast.TableOptionCharset, Default: true,
					UintValue: ast.TableOptionCharsetWithConvertTo}},
			}
		} else {
			op = &ast.AlterTableSpec{
				Tp: ast.AlterTableOption,
				Options: []*ast.TableOption{{Tp: ast.TableOptionCharset, StrValue: r.parseCharsetName(),
					UintValue: ast.TableOptionCharsetWithConvertTo}},
			}
		}
		// OptCollate
		if r.accept(collate) {
			op.Options = append(op.Options, &ast.TableOption{Tp: ast.TableOptionCollate, StrValue: r.parseCollationName()})
		}
		return op
	case add:
		return r.parseAlterTableSpecAdd()
	case last:
		// "LAST" "PARTITION" "LESS" "THAN" '(' BitExpr ')' NoWriteToBinLogAliasOpt
		start := r.cur().offset
		r.advance()
		r.expect(partition)
		r.expect(less)
		r.expect(than)
		r.expect(int('('))
		expr := r.parseBitExpr(0)
		r.expect(int(')'))
		noWriteToBinlog := r.parseNoWriteToBinLogAliasOpt()
		if noWriteToBinlog {
			r.appendWarnf("The NO_WRITE_TO_BINLOG option is parsed but ignored for now.")
		}
		partitionMethod := ast.PartitionMethod{Expr: expr}
		r.p.setNodeText(&partitionMethod, r.src[start:r.cur().offset])
		return &ast.AlterTableSpec{
			NoWriteToBinlog: noWriteToBinlog,
			Tp:              ast.AlterTableAddLastPartition,
			Partition:       &ast.PartitionOptions{PartitionMethod: partitionMethod},
		}
	case first:
		// "FIRST" "PARTITION" "LESS" "THAN" '(' BitExpr ')' IfExists
		start := r.cur().offset
		r.advance()
		r.expect(partition)
		r.expect(less)
		r.expect(than)
		r.expect(int('('))
		expr := r.parseBitExpr(0)
		r.expect(int(')'))
		ifExists := r.parseIfExists()
		partitionMethod := ast.PartitionMethod{Expr: expr}
		r.p.setNodeText(&partitionMethod, r.src[start:r.cur().offset])
		return &ast.AlterTableSpec{
			IfExists:  ifExists,
			Tp:        ast.AlterTableDropFirstPartition,
			Partition: &ast.PartitionOptions{PartitionMethod: partitionMethod},
		}
	case enable:
		r.advance()
		if r.tok() == masking {
			// "ENABLE" "MASKING" "POLICY" PolicyName
			r.advance()
			r.expect(policy)
			return &ast.AlterTableSpec{
				Tp:                ast.AlterTableEnableMaskingPolicy,
				MaskingPolicyName: ast.NewCIStr(r.parseIdentifier()),
			}
		}
		// "ENABLE" "KEYS"
		r.expect(keys)
		return &ast.AlterTableSpec{
			Tp: ast.AlterTableEnableKeys,
		}
	case disable:
		r.advance()
		if r.tok() == masking {
			// "DISABLE" "MASKING" "POLICY" PolicyName
			r.advance()
			r.expect(policy)
			return &ast.AlterTableSpec{
				Tp:                ast.AlterTableDisableMaskingPolicy,
				MaskingPolicyName: ast.NewCIStr(r.parseIdentifier()),
			}
		}
		// "DISABLE" "KEYS"
		r.expect(keys)
		return &ast.AlterTableSpec{
			Tp: ast.AlterTableDisableKeys,
		}
	case drop:
		return r.parseAlterTableSpecDrop()
	case modify:
		return r.parseAlterTableSpecModify()
	case change:
		// "CHANGE" ColumnKeywordOpt IfExists ColumnName ColumnDef ColumnPosition
		r.advance()
		r.accept(column)
		ifExists := r.parseIfExists()
		oldColumnName := r.parseColumnName()
		colDef := r.parseColumnDef()
		return &ast.AlterTableSpec{
			IfExists:      ifExists,
			Tp:            ast.AlterTableChangeColumn,
			OldColumnName: oldColumnName,
			NewColumns:    []*ast.ColumnDef{colDef},
			Position:      r.parseColumnPosition(),
		}
	case attributes:
		return &ast.AlterTableSpec{
			Tp:             ast.AlterTableAttributes,
			AttributesSpec: r.parseAttributesOpt(),
		}
	case statsOptions:
		return &ast.AlterTableSpec{
			Tp:               ast.AlterTableStatsOptions,
			StatsOptionsSpec: r.parseStatsOptionsOpt(),
		}
	case check:
		// "CHECK" "PARTITION" AllOrPartitionNameList
		r.advance()
		r.expect(partition)
		names := r.parseAllOrPartitionNameList()
		r.appendWarnf("The CHECK PARTITIONING clause is parsed but not implement yet.")
		ret := &ast.AlterTableSpec{
			Tp: ast.AlterTableCheckPartitions,
		}
		if names == nil {
			ret.OnAllPartitions = true
		} else {
			ret.PartitionNames = names
		}
		return ret
	case coalesce:
		// "COALESCE" "PARTITION" NoWriteToBinLogAliasOpt NUM
		r.advance()
		r.expect(partition)
		noWriteToBinlog := r.parseNoWriteToBinLogAliasOpt()
		num := getUint64FromNUM(r.expect(intLit).item)
		if noWriteToBinlog {
			r.appendWarnf("The NO_WRITE_TO_BINLOG option is parsed but ignored for now.")
		}
		return &ast.AlterTableSpec{
			Tp:              ast.AlterTableCoalescePartitions,
			NoWriteToBinlog: noWriteToBinlog,
			Num:             num,
		}
	case exchange:
		// "EXCHANGE" "PARTITION" Identifier "WITH" "TABLE" TableName WithValidationOpt
		r.advance()
		r.expect(partition)
		name := r.parseIdentifier()
		r.expect(with)
		r.expect(tableKwd)
		newTable := r.parseTableName()
		return &ast.AlterTableSpec{
			Tp:             ast.AlterTableExchangePartition,
			PartitionNames: []ast.CIStr{ast.NewCIStr(name)},
			NewTable:       newTable,
			WithValidation: r.parseWithValidationOpt(),
		}
	case truncate:
		// "TRUNCATE" "PARTITION" AllOrPartitionNameList
		r.advance()
		r.expect(partition)
		names := r.parseAllOrPartitionNameList()
		ret := &ast.AlterTableSpec{
			Tp: ast.AlterTableTruncatePartition,
		}
		if names == nil {
			ret.OnAllPartitions = true
		} else {
			ret.PartitionNames = names
		}
		return ret
	case optimize, repair, rebuild:
		// ("OPTIMIZE"|"REPAIR"|"REBUILD") "PARTITION" NoWriteToBinLogAliasOpt
		// AllOrPartitionNameList
		var tp ast.AlterTableType
		switch r.tok() {
		case optimize:
			tp = ast.AlterTableOptimizePartition
		case repair:
			tp = ast.AlterTableRepairPartition
		case rebuild:
			tp = ast.AlterTableRebuildPartition
		}
		r.advance()
		r.expect(partition)
		ret := &ast.AlterTableSpec{
			NoWriteToBinlog: r.parseNoWriteToBinLogAliasOpt(),
			Tp:              tp,
		}
		names := r.parseAllOrPartitionNameList()
		if names == nil {
			ret.OnAllPartitions = true
		} else {
			ret.PartitionNames = names
		}
		return ret
	case importKwd:
		r.advance()
		if r.accept(partition) {
			// "IMPORT" "PARTITION" AllOrPartitionNameList "TABLESPACE"
			names := r.parseAllOrPartitionNameList()
			r.expect(tablespace)
			ret := &ast.AlterTableSpec{
				Tp: ast.AlterTableImportPartitionTablespace,
			}
			if names == nil {
				ret.OnAllPartitions = true
			} else {
				ret.PartitionNames = names
			}
			r.appendWarnf("The IMPORT PARTITION TABLESPACE clause is parsed but ignored by all storage engines.")
			return ret
		}
		// "IMPORT" "TABLESPACE"
		r.expect(tablespace)
		ret := &ast.AlterTableSpec{
			Tp: ast.AlterTableImportTablespace,
		}
		r.appendWarnf("The IMPORT TABLESPACE clause is parsed but ignored by all storage engines.")
		return ret
	case discard:
		r.advance()
		if r.accept(partition) {
			// "DISCARD" "PARTITION" AllOrPartitionNameList "TABLESPACE"
			names := r.parseAllOrPartitionNameList()
			r.expect(tablespace)
			ret := &ast.AlterTableSpec{
				Tp: ast.AlterTableDiscardPartitionTablespace,
			}
			if names == nil {
				ret.OnAllPartitions = true
			} else {
				ret.PartitionNames = names
			}
			r.appendWarnf("The DISCARD PARTITION TABLESPACE clause is parsed but ignored by all storage engines.")
			return ret
		}
		// "DISCARD" "TABLESPACE"
		r.expect(tablespace)
		ret := &ast.AlterTableSpec{
			Tp: ast.AlterTableDiscardTablespace,
		}
		r.appendWarnf("The DISCARD TABLESPACE clause is parsed but ignored by all storage engines.")
		return ret
	case order:
		// "ORDER" "BY" AlterOrderList: the ','-continuations of
		// AlterOrderList shift over ending the spec (',' is declared after
		// %prec lowerThenOrder), so commas are consumed greedily here.
		r.advance()
		r.expect(by)
		items := []*ast.AlterOrderItem{r.parseAlterOrderItem()}
		for r.accept(int(',')) {
			items = append(items, r.parseAlterOrderItem())
		}
		return &ast.AlterTableSpec{
			Tp:          ast.AlterTableOrderByColumns,
			OrderByList: items,
		}
	case alter:
		return r.parseAlterTableSpecAlter()
	case rename:
		r.advance()
		switch r.tok() {
		case column:
			// "RENAME" "COLUMN" Identifier "TO" Identifier
			r.advance()
			oldColName := &ast.ColumnName{Name: ast.NewCIStr(r.parseIdentifier())}
			r.expect(to)
			newColName := &ast.ColumnName{Name: ast.NewCIStr(r.parseIdentifier())}
			return &ast.AlterTableSpec{
				Tp:            ast.AlterTableRenameColumn,
				OldColumnName: oldColName,
				NewColumnName: newColName,
			}
		case key, index:
			// "RENAME" KeyOrIndex Identifier "TO" Identifier
			r.advance()
			from := r.parseIdentifier()
			r.expect(to)
			return &ast.AlterTableSpec{
				Tp:      ast.AlterTableRenameIndex,
				FromKey: ast.NewCIStr(from),
				ToKey:   ast.NewCIStr(r.parseIdentifier()),
			}
		case to, as, eq:
			// "RENAME" ("TO"|"AS"|eq) TableName
			r.advance()
		}
		// "RENAME" EqOpt TableName (with the empty EqOpt)
		return &ast.AlterTableSpec{
			Tp:       ast.AlterTableRenameTable,
			NewTable: r.parseTableName(),
		}
	case lock:
		return &ast.AlterTableSpec{
			Tp:       ast.AlterTableLock,
			LockType: r.parseLockClause(),
		}
	case read:
		// Writeable: "READ" "WRITE" | "READ" "ONLY"
		r.advance()
		var writeable bool
		switch r.tok() {
		case write:
			writeable = true
		case only:
			writeable = false
		default:
			r.syntaxError()
		}
		r.advance()
		return &ast.AlterTableSpec{
			Tp:        ast.AlterTableWriteable,
			Writeable: writeable,
		}
	case algorithm:
		// Parse it and ignore it. Just for compatibility.
		return &ast.AlterTableSpec{
			Tp:        ast.AlterTableAlgorithm,
			Algorithm: r.parseAlgorithmClause(),
		}
	case force:
		if r.la(1) == autoIncrement || r.la(1) == autoRandomBase {
			// TableOptionList beginning with ForceOpt "AUTO_INCREMENT"/...
			break
		}
		// "FORCE": parse it and ignore it. Just for compatibility.
		r.advance()
		return &ast.AlterTableSpec{
			Tp: ast.AlterTableForce,
		}
	case with:
		// "WITH" "VALIDATION": parse it and ignore it. Just for compatibility.
		r.advance()
		r.expect(validation)
		return &ast.AlterTableSpec{
			Tp: ast.AlterTableWithValidation,
		}
	case without:
		// "WITHOUT" "VALIDATION": parse it and ignore it. Just for compatibility.
		r.advance()
		r.expect(validation)
		return &ast.AlterTableSpec{
			Tp: ast.AlterTableWithoutValidation,
		}
	case secondaryLoad:
		// Parse it and ignore it. Just for compatibility.
		r.advance()
		ret := &ast.AlterTableSpec{
			Tp: ast.AlterTableSecondaryLoad,
		}
		r.appendWarnf("The SECONDARY_LOAD clause is parsed but not implement yet.")
		return ret
	case secondaryUnload:
		// Parse it and ignore it. Just for compatibility.
		r.advance()
		ret := &ast.AlterTableSpec{
			Tp: ast.AlterTableSecondaryUnload,
		}
		r.appendWarnf("The SECONDARY_UNLOAD VALIDATION clause is parsed but not implement yet.")
		return ret
	case cache:
		r.advance()
		return &ast.AlterTableSpec{
			Tp: ast.AlterTableCache,
		}
	case nocache:
		r.advance()
		return &ast.AlterTableSpec{
			Tp: ast.AlterTableNoCache,
		}
	}
	// AlterTableSpec: TableOptionList %prec higherThanComma. The rule's
	// precedence is above ',', so a comma ends the option list (it
	// separates specs) rather than continuing it as in CREATE TABLE.
	if !r.isTableOptionStart() {
		r.syntaxError()
	}
	opts := []*ast.TableOption{r.parseTableOption()}
	for r.isTableOptionStart() {
		opts = append(opts, r.parseTableOption())
	}
	return &ast.AlterTableSpec{
		Tp:      ast.AlterTableOption,
		Options: opts,
	}
}

// parseAlterTableSpecAdd implements the "ADD" alternatives of AlterTableSpec.
func (r *rdParser) parseAlterTableSpecAdd() *ast.AlterTableSpec {
	r.expect(add)
	switch r.tok() {
	case constraint, primary, unique, fulltext, foreign, check, key, index:
		// "ADD" ConstraintWithColumnarIndex (the Constraint alternative)
		return &ast.AlterTableSpec{
			Tp:         ast.AlterTableAddConstraint,
			Constraint: r.parseConstraint(),
		}
	case vectorType, columnar:
		// ConstraintVectorIndex/ConstraintColumnarIndex when followed by
		// "INDEX"; otherwise VECTOR/COLUMNAR is a column name.
		if r.la(1) == index {
			return &ast.AlterTableSpec{
				Tp:         ast.AlterTableAddConstraint,
				Constraint: r.parseConstraintVectorOrColumnarIndex(),
			}
		}
	case partition:
		// "ADD" "PARTITION" IfNotExists NoWriteToBinLogAliasOpt
		// (PartitionDefinitionListOpt | "PARTITIONS" NUM)
		r.advance()
		ifNotExists := r.parseIfNotExists()
		noWriteToBinlog := r.parseNoWriteToBinLogAliasOpt()
		if r.accept(partitions) {
			num := getUint64FromNUM(r.expect(intLit).item)
			if noWriteToBinlog {
				r.appendWarnf("The NO_WRITE_TO_BINLOG option is parsed but ignored for now.")
			}
			return &ast.AlterTableSpec{
				IfNotExists:     ifNotExists,
				NoWriteToBinlog: noWriteToBinlog,
				Tp:              ast.AlterTableAddPartitions,
				Num:             num,
			}
		}
		defs := r.parsePartitionDefinitionListOpt()
		if noWriteToBinlog {
			r.appendWarnf("The NO_WRITE_TO_BINLOG option is parsed but ignored for now.")
		}
		return &ast.AlterTableSpec{
			IfNotExists:     ifNotExists,
			NoWriteToBinlog: noWriteToBinlog,
			Tp:              ast.AlterTableAddPartitions,
			PartDefinitions: defs,
		}
	case statsExtended:
		// "ADD" "STATS_EXTENDED" IfNotExists Identifier StatsType '(' ColumnNameList ')'
		r.advance()
		ifNotExists := r.parseIfNotExists()
		name := r.parseIdentifier()
		statsType := r.parseStatsType()
		r.expect(int('('))
		cols := r.parseColumnNameList()
		r.expect(int(')'))
		statsSpec := &ast.StatisticsSpec{
			StatsName: name,
			StatsType: statsType,
			Columns:   cols,
		}
		return &ast.AlterTableSpec{
			Tp:          ast.AlterTableAddStatistics,
			IfNotExists: ifNotExists,
			Statistics:  statsSpec,
		}
	case masking:
		// The MASKING token shifts over reducing it to an identifier, so
		// these alternatives win over the plain-ColumnDef path unless only
		// that path can proceed (a qualified column name).
		if r.la(1) == int('.') {
			break
		}
		r.advance()
		switch r.tok() {
		case policy:
			// "ADD" "MASKING" "POLICY" PolicyName "ON" '(' Identifier ')'
			// "AS" Expression MaskingPolicyRestrictOnOpt MaskingPolicyStateOpt
			r.advance()
			name := r.parseIdentifier()
			r.expect(on)
			r.expect(int('('))
			col := r.parseIdentifier()
			r.expect(int(')'))
			r.expect(as)
			expr := r.parseExpression()
			restrictOps := r.parseMaskingPolicyRestrictOnOpt()
			state := r.parseMaskingPolicyStateOpt()
			return &ast.AlterTableSpec{
				Tp:                       ast.AlterTableAddMaskingPolicy,
				MaskingPolicyName:        ast.NewCIStr(name),
				MaskingPolicyColumn:      &ast.ColumnName{Name: ast.NewCIStr(col)},
				MaskingPolicyExpr:        expr,
				MaskingPolicyRestrictOps: restrictOps,
				MaskingPolicyState:       *state,
			}
		case serial:
			// "ADD" "MASKING" "SERIAL" ColumnOptionListOpt ColumnPosition
			// Keep behavior consistent with `ColumnDef: ColumnName "SERIAL" ...`.
			r.advance()
			tp := types.NewFieldType(mysql.TypeLonglong)
			options := []*ast.ColumnOption{{Tp: ast.ColumnOptionNotNull}, {Tp: ast.ColumnOptionAutoIncrement}, {Tp: ast.ColumnOptionUniqKey}}
			options = append(options, r.parseColumnOptionListOpt().Options...)
			tp.AddFlag(mysql.UnsignedFlag)
			colDef := &ast.ColumnDef{
				Name:    &ast.ColumnName{Name: ast.NewCIStr("masking")},
				Tp:      tp,
				Options: options,
			}
			if err := colDef.Validate(); err != nil {
				r.actionError(err)
			}
			return &ast.AlterTableSpec{
				IfNotExists: false,
				Tp:          ast.AlterTableAddColumns,
				NewColumns:  []*ast.ColumnDef{colDef},
				Position:    r.parseColumnPosition(),
			}
		}
		// "ADD" "MASKING" Type ColumnOptionListOpt ColumnPosition
		tp := r.parseType()
		colDef := &ast.ColumnDef{
			Name:    &ast.ColumnName{Name: ast.NewCIStr("masking")},
			Tp:      tp,
			Options: r.parseColumnOptionListOpt().Options,
		}
		if err := colDef.Validate(); err != nil {
			r.actionError(err)
		}
		return &ast.AlterTableSpec{
			IfNotExists: false,
			Tp:          ast.AlterTableAddColumns,
			NewColumns:  []*ast.ColumnDef{colDef},
			Position:    r.parseColumnPosition(),
		}
	}
	// "ADD" ColumnKeywordOpt IfNotExists
	// (ColumnDef ColumnPosition | '(' TableElementList ')')
	r.accept(column)
	ifNotExists := r.parseIfNotExists()
	if r.accept(int('(')) {
		tes := r.parseTableElementList()
		r.expect(int(')'))
		var columnDefs []*ast.ColumnDef
		var constraints []*ast.Constraint
		for _, te := range tes {
			switch te := te.(type) {
			case *ast.ColumnDef:
				columnDefs = append(columnDefs, te)
			case *ast.Constraint:
				constraints = append(constraints, te)
			}
		}
		return &ast.AlterTableSpec{
			IfNotExists:    ifNotExists,
			Tp:             ast.AlterTableAddColumns,
			NewColumns:     columnDefs,
			NewConstraints: constraints,
		}
	}
	colDef := r.parseColumnDef()
	return &ast.AlterTableSpec{
		IfNotExists: ifNotExists,
		Tp:          ast.AlterTableAddColumns,
		NewColumns:  []*ast.ColumnDef{colDef},
		Position:    r.parseColumnPosition(),
	}
}

// parseStatsType implements StatsType.
func (r *rdParser) parseStatsType() uint8 {
	var tp uint8
	switch r.tok() {
	case cardinality:
		tp = ast.StatsTypeCardinality
	case dependency:
		tp = ast.StatsTypeDependency
	case correlation:
		tp = ast.StatsTypeCorrelation
	default:
		r.syntaxError()
	}
	r.advance()
	return tp
}

// parseAlterTableSpecDrop implements the "DROP" alternatives of
// AlterTableSpec.
func (r *rdParser) parseAlterTableSpecDrop() *ast.AlterTableSpec {
	r.expect(drop)
	switch r.tok() {
	case primary:
		// "DROP" "PRIMARY" "KEY"
		r.advance()
		r.expect(key)
		return &ast.AlterTableSpec{Tp: ast.AlterTableDropPrimaryKey}
	case partition:
		// "DROP" "PARTITION" IfExists PartitionNameList %prec lowerThanComma
		r.advance()
		ifExists := r.parseIfExists()
		return &ast.AlterTableSpec{
			IfExists:       ifExists,
			Tp:             ast.AlterTableDropPartition,
			PartitionNames: r.parsePartitionNameList(),
		}
	case statsExtended:
		// "DROP" "STATS_EXTENDED" IfExists Identifier
		r.advance()
		ifExists := r.parseIfExists()
		statsSpec := &ast.StatisticsSpec{
			StatsName: r.parseIdentifier(),
		}
		return &ast.AlterTableSpec{
			Tp:         ast.AlterTableDropStatistics,
			IfExists:   ifExists,
			Statistics: statsSpec,
		}
	case key, index:
		// "DROP" KeyOrIndex IfExists Identifier
		r.advance()
		ifExists := r.parseIfExists()
		return &ast.AlterTableSpec{
			IfExists: ifExists,
			Tp:       ast.AlterTableDropIndex,
			Name:     r.parseIdentifier(),
		}
	case foreign:
		// "DROP" "FOREIGN" "KEY" Symbol
		r.advance()
		r.expect(key)
		return &ast.AlterTableSpec{
			Tp:   ast.AlterTableDropForeignKey,
			Name: r.parseIdentifier(),
		}
	case check, constraint:
		// "DROP" CheckConstraintKeyword Identifier
		// Parse it and ignore it. Just for compatibility.
		r.advance()
		c := &ast.Constraint{
			Name: r.parseIdentifier(),
		}
		return &ast.AlterTableSpec{
			Tp:         ast.AlterTableDropCheck,
			Constraint: c,
		}
	case masking:
		// As in ADD, the MASKING token shifts unless only the qualified
		// column-name path can proceed.
		if r.la(1) == int('.') {
			break
		}
		r.advance()
		if r.tok() == policy {
			// "DROP" "MASKING" "POLICY" PolicyName
			r.advance()
			return &ast.AlterTableSpec{
				Tp:                ast.AlterTableDropMaskingPolicy,
				MaskingPolicyName: ast.NewCIStr(r.parseIdentifier()),
			}
		}
		// "DROP" "MASKING" RestrictOrCascadeOpt
		r.parseRestrictOrCascadeOpt()
		return &ast.AlterTableSpec{
			IfExists:      false,
			Tp:            ast.AlterTableDropColumn,
			OldColumnName: &ast.ColumnName{Name: ast.NewCIStr("masking")},
		}
	}
	// "DROP" ColumnKeywordOpt IfExists ColumnName RestrictOrCascadeOpt
	r.accept(column)
	ifExists := r.parseIfExists()
	col := r.parseColumnName()
	r.parseRestrictOrCascadeOpt()
	return &ast.AlterTableSpec{
		IfExists:      ifExists,
		Tp:            ast.AlterTableDropColumn,
		OldColumnName: col,
	}
}

// parseAlterTableSpecModify implements the "MODIFY" alternatives of
// AlterTableSpec.
func (r *rdParser) parseAlterTableSpecModify() *ast.AlterTableSpec {
	r.expect(modify)
	if r.tok() == masking && r.la(1) != int('.') {
		r.advance()
		switch r.tok() {
		case policy:
			// "MODIFY" "MASKING" "POLICY" PolicyName "SET" ...
			r.advance()
			name := r.parseIdentifier()
			r.expect(set)
			if r.accept(restrict) {
				// ... "SET" "RESTRICT" "ON" ('(' MaskingPolicyRestrictOperationList ')' | "NONE")
				r.expect(on)
				if r.accept(none) {
					return &ast.AlterTableSpec{
						Tp:                       ast.AlterTableModifyMaskingPolicyRestrictOn,
						MaskingPolicyName:        ast.NewCIStr(name),
						MaskingPolicyRestrictOps: ast.MaskingPolicyRestrictOpNone,
					}
				}
				r.expect(int('('))
				// MaskingPolicyRestrictOperationList
				ops := r.parseMaskingPolicyRestrictOperation()
				for r.accept(int(',')) {
					ops |= r.parseMaskingPolicyRestrictOperation()
				}
				r.expect(int(')'))
				return &ast.AlterTableSpec{
					Tp:                       ast.AlterTableModifyMaskingPolicyRestrictOn,
					MaskingPolicyName:        ast.NewCIStr(name),
					MaskingPolicyRestrictOps: ops,
				}
			}
			// ... "SET" Identifier "=" Expression
			opt := r.parseIdentifier()
			r.expect(eq)
			expr := r.parseExpression()
			if !strings.EqualFold(opt, "expression") {
				r.actionError(r.actionErrorf("unsupported masking policy modify option: %s", opt))
			}
			return &ast.AlterTableSpec{
				Tp:                ast.AlterTableModifyMaskingPolicyExpression,
				MaskingPolicyName: ast.NewCIStr(name),
				MaskingPolicyExpr: expr,
			}
		case serial:
			// "MODIFY" "MASKING" "SERIAL" ColumnOptionListOpt ColumnPosition
			// Keep behavior consistent with `ColumnDef: ColumnName "SERIAL" ...`.
			r.advance()
			tp := types.NewFieldType(mysql.TypeLonglong)
			options := []*ast.ColumnOption{{Tp: ast.ColumnOptionNotNull}, {Tp: ast.ColumnOptionAutoIncrement}, {Tp: ast.ColumnOptionUniqKey}}
			options = append(options, r.parseColumnOptionListOpt().Options...)
			tp.AddFlag(mysql.UnsignedFlag)
			colDef := &ast.ColumnDef{
				Name:    &ast.ColumnName{Name: ast.NewCIStr("masking")},
				Tp:      tp,
				Options: options,
			}
			if err := colDef.Validate(); err != nil {
				r.actionError(err)
			}
			return &ast.AlterTableSpec{
				IfExists:   false,
				Tp:         ast.AlterTableModifyColumn,
				NewColumns: []*ast.ColumnDef{colDef},
				Position:   r.parseColumnPosition(),
			}
		}
		// "MODIFY" "MASKING" Type ColumnOptionListOpt ColumnPosition
		tp := r.parseType()
		colDef := &ast.ColumnDef{
			Name:    &ast.ColumnName{Name: ast.NewCIStr("masking")},
			Tp:      tp,
			Options: r.parseColumnOptionListOpt().Options,
		}
		if err := colDef.Validate(); err != nil {
			r.actionError(err)
		}
		return &ast.AlterTableSpec{
			IfExists:   false,
			Tp:         ast.AlterTableModifyColumn,
			NewColumns: []*ast.ColumnDef{colDef},
			Position:   r.parseColumnPosition(),
		}
	}
	// "MODIFY" ColumnKeywordOpt IfExists ColumnDef ColumnPosition
	r.accept(column)
	ifExists := r.parseIfExists()
	colDef := r.parseColumnDef()
	return &ast.AlterTableSpec{
		IfExists:   ifExists,
		Tp:         ast.AlterTableModifyColumn,
		NewColumns: []*ast.ColumnDef{colDef},
		Position:   r.parseColumnPosition(),
	}
}

// parseAlterTableSpecAlter implements the "ALTER"-led alternatives of
// AlterTableSpec: ALTER [COLUMN] col SET DEFAULT/DROP DEFAULT, ALTER
// CHECK/CONSTRAINT, and ALTER INDEX visibility.
func (r *rdParser) parseAlterTableSpecAlter() *ast.AlterTableSpec {
	r.expect(alter)
	switch r.tok() {
	case check, constraint:
		// "ALTER" CheckConstraintKeyword Identifier EnforcedOrNot
		r.advance()
		name := r.parseIdentifier()
		c := &ast.Constraint{
			Name:     name,
			Enforced: r.parseEnforcedOrNot(),
		}
		return &ast.AlterTableSpec{
			Tp:         ast.AlterTableAlterCheck,
			Constraint: c,
		}
	case index:
		// "ALTER" "INDEX" Identifier IndexInvisible
		r.advance()
		name := r.parseIdentifier()
		return &ast.AlterTableSpec{
			Tp:         ast.AlterTableIndexInvisible,
			IndexName:  ast.NewCIStr(name),
			Visibility: r.parseIndexInvisible(),
		}
	}
	// "ALTER" ColumnKeywordOpt ColumnName
	// ("SET" "DEFAULT" (SignedLiteral | '(' Expression ')') | "DROP" "DEFAULT")
	r.accept(column)
	col := r.parseColumnName()
	switch r.tok() {
	case set:
		r.advance()
		r.expect(defaultKwd)
		var expr ast.ExprNode
		if r.accept(int('(')) {
			expr = r.parseExpression()
			r.expect(int(')'))
		} else {
			expr = r.parseSignedLiteral()
		}
		option := &ast.ColumnOption{Expr: expr}
		colDef := &ast.ColumnDef{
			Name:    col,
			Options: []*ast.ColumnOption{option},
		}
		return &ast.AlterTableSpec{
			Tp:         ast.AlterTableAlterColumn,
			NewColumns: []*ast.ColumnDef{colDef},
		}
	case drop:
		r.advance()
		r.expect(defaultKwd)
		colDef := &ast.ColumnDef{
			Name: col,
		}
		return &ast.AlterTableSpec{
			Tp:         ast.AlterTableAlterColumn,
			NewColumns: []*ast.ColumnDef{colDef},
		}
	}
	r.syntaxError()
	return nil
}

// parseEnforcedOrNot implements EnforcedOrNot: "ENFORCED" | NotSym "ENFORCED".
func (r *rdParser) parseEnforcedOrNot() bool {
	if r.tok() == not || r.tok() == not2 {
		r.advance()
		r.expect(enforced)
		return false
	}
	r.expect(enforced)
	return true
}

// parseIndexInvisible implements IndexInvisible: "VISIBLE" | "INVISIBLE".
func (r *rdParser) parseIndexInvisible() ast.IndexVisibility {
	switch r.tok() {
	case visible:
		r.advance()
		return ast.IndexVisibilityVisible
	case invisible:
		r.advance()
		return ast.IndexVisibilityInvisible
	}
	r.syntaxError()
	return ast.IndexVisibilityVisible
}

// parseWithValidationOpt implements WithValidationOpt and WithValidation.
func (r *rdParser) parseWithValidationOpt() bool {
	switch r.tok() {
	case with:
		r.advance()
		r.expect(validation)
		return true
	case without:
		r.advance()
		r.expect(validation)
		return false
	}
	return true
}

// parseAttributesOpt implements AttributesOpt.
func (r *rdParser) parseAttributesOpt() *ast.AttributesSpec {
	r.expect(attributes)
	r.parseEqOpt()
	if r.accept(defaultKwd) {
		return &ast.AttributesSpec{Default: true}
	}
	return &ast.AttributesSpec{Default: false, Attributes: r.expect(stringLit).lit}
}

// parseStatsOptionsOpt implements StatsOptionsOpt.
func (r *rdParser) parseStatsOptionsOpt() *ast.StatsOptionsSpec {
	r.expect(statsOptions)
	r.parseEqOpt()
	if r.accept(defaultKwd) {
		return &ast.StatsOptionsSpec{Default: true}
	}
	return &ast.StatsOptionsSpec{Default: false, StatsOptions: r.expect(stringLit).lit}
}

// parseAlterOrderItem implements AlterOrderItem: ColumnName OptOrder.
func (r *rdParser) parseAlterOrderItem() *ast.AlterOrderItem {
	col := r.parseColumnName()
	return &ast.AlterOrderItem{Column: col, Desc: r.parseOptOrder()}
}

// parseColumnPosition implements ColumnPosition:
// empty | "FIRST" | "AFTER" ColumnName.
func (r *rdParser) parseColumnPosition() *ast.ColumnPosition {
	switch r.tok() {
	case first:
		r.advance()
		return &ast.ColumnPosition{Tp: ast.ColumnPositionFirst}
	case after:
		r.advance()
		return &ast.ColumnPosition{
			Tp:             ast.ColumnPositionAfter,
			RelativeColumn: r.parseColumnName(),
		}
	}
	return &ast.ColumnPosition{Tp: ast.ColumnPositionNone}
}

// parseLockClause implements LockClause.
func (r *rdParser) parseLockClause() ast.LockType {
	r.expect(lock)
	r.parseEqOpt()
	if r.accept(defaultKwd) {
		return ast.LockTypeDefault
	}
	name := r.parseIdentifier()
	switch strings.ToUpper(name) {
	case "NONE":
		return ast.LockTypeNone
	case "SHARED":
		return ast.LockTypeShared
	case "EXCLUSIVE":
		return ast.LockTypeExclusive
	}
	r.actionError(ErrUnknownAlterLock.GenWithStackByArgs(name))
	return ast.LockTypeDefault
}

// parseAlgorithmClause implements AlgorithmClause.
func (r *rdParser) parseAlgorithmClause() ast.AlgorithmType {
	r.expect(algorithm)
	r.parseEqOpt()
	switch r.tok() {
	case defaultKwd:
		r.advance()
		return ast.AlgorithmTypeDefault
	case copyKwd:
		r.advance()
		return ast.AlgorithmTypeCopy
	case inplace:
		r.advance()
		return ast.AlgorithmTypeInplace
	case instant:
		r.advance()
		return ast.AlgorithmTypeInstant
	case identifier:
		// "ALGORITHM" EqOpt identifier: an unknown algorithm is an error.
		lit := r.cur().lit
		r.advance()
		r.actionError(ErrUnknownAlterAlgorithm.GenWithStackByArgs(lit))
	}
	r.syntaxError()
	return ast.AlgorithmTypeDefault
}

// parseAlterTableSpecSingleOpt implements AlterTableSpecSingleOpt (nil for
// the empty PartitionOpt alternative).
func (r *rdParser) parseAlterTableSpecSingleOpt() *ast.AlterTableSpec {
	switch r.tok() {
	case partition:
		if r.la(1) == by {
			// PartitionOpt (non-empty)
			return &ast.AlterTableSpec{
				Tp:        ast.AlterTablePartition,
				Partition: r.parsePartitionOpt(),
			}
		}
		// "PARTITION" Identifier (AttributesOpt | PartDefOptionList)
		r.advance()
		name := r.parseIdentifier()
		if r.tok() == attributes {
			return &ast.AlterTableSpec{
				Tp:             ast.AlterTablePartitionAttributes,
				PartitionNames: []ast.CIStr{ast.NewCIStr(name)},
				AttributesSpec: r.parseAttributesOpt(),
			}
		}
		return &ast.AlterTableSpec{
			Tp:             ast.AlterTablePartitionOptions,
			PartitionNames: []ast.CIStr{ast.NewCIStr(name)},
			Options:        r.parsePartDefOptionList(),
		}
	case remove:
		r.advance()
		switch r.tok() {
		case partitioning:
			// "REMOVE" "PARTITIONING"
			r.advance()
			return &ast.AlterTableSpec{
				Tp: ast.AlterTableRemovePartitioning,
			}
		case ttl:
			// "REMOVE" "TTL"
			r.advance()
			return &ast.AlterTableSpec{
				Tp: ast.AlterTableRemoveTTL,
			}
		}
		r.syntaxError()
	case reorganize:
		// "REORGANIZE" "PARTITION" NoWriteToBinLogAliasOpt ReorganizePartitionRuleOpt
		r.advance()
		r.expect(partition)
		noWriteToBinlog := r.parseNoWriteToBinLogAliasOpt()
		ret := r.parseReorganizePartitionRuleOpt()
		ret.NoWriteToBinlog = noWriteToBinlog
		return ret
	case split:
		if r.la(1) == maxValue {
			// "SPLIT" "MAXVALUE" "PARTITION" "LESS" "THAN" '(' BitExpr ')'
			start := r.cur().offset
			r.advance()
			r.advance()
			r.expect(partition)
			r.expect(less)
			r.expect(than)
			r.expect(int('('))
			expr := r.parseBitExpr(0)
			r.expect(int(')'))
			partitionMethod := ast.PartitionMethod{Expr: expr}
			r.p.setNodeText(&partitionMethod, r.src[start:r.cur().offset])
			return &ast.AlterTableSpec{
				Tp:        ast.AlterTableReorganizeLastPartition,
				Partition: &ast.PartitionOptions{PartitionMethod: partitionMethod},
			}
		}
		// SplitIndexOption
		return &ast.AlterTableSpec{
			Tp:         ast.AlterTableSplitIndex,
			SplitIndex: r.parseSplitIndexOption(),
		}
	case merge:
		// "MERGE" "FIRST" "PARTITION" "LESS" "THAN" '(' BitExpr ')'
		start := r.cur().offset
		r.advance()
		r.expect(first)
		r.expect(partition)
		r.expect(less)
		r.expect(than)
		r.expect(int('('))
		expr := r.parseBitExpr(0)
		r.expect(int(')'))
		partitionMethod := ast.PartitionMethod{Expr: expr}
		r.p.setNodeText(&partitionMethod, r.src[start:r.cur().offset])
		return &ast.AlterTableSpec{
			Tp:        ast.AlterTableReorganizeFirstPartition,
			Partition: &ast.PartitionOptions{PartitionMethod: partitionMethod},
		}
	}
	return nil
}

// parseReorganizePartitionRuleOpt implements ReorganizePartitionRuleOpt.
func (r *rdParser) parseReorganizePartitionRuleOpt() *ast.AlterTableSpec {
	if !isIdentifierTok(r.tok()) {
		// empty %prec lowerThanRemove
		return &ast.AlterTableSpec{
			Tp:              ast.AlterTableReorganizePartition,
			OnAllPartitions: true,
		}
	}
	// PartitionNameList "INTO" '(' PartitionDefinitionList ')'
	names := r.parsePartitionNameList()
	r.expect(into)
	if r.tok() != int('(') {
		r.syntaxError()
	}
	defs := r.parsePartitionDefinitionListOpt()
	return &ast.AlterTableSpec{
		Tp:              ast.AlterTableReorganizePartition,
		PartitionNames:  names,
		PartDefinitions: defs,
	}
}

/**************************************AlterDatabaseStmt***************************************/

// parseAlterDatabaseStmt implements AlterDatabaseStmt:
//
//	"ALTER" DatabaseSym DBName DatabaseOptionList
//	| "ALTER" DatabaseSym DatabaseOptionList
//
// A leading token that can begin a DatabaseOption selects the no-name
// alternative, EXCEPT the unreserved keywords whose token precedence is
// above the empty DefaultKwdOpt's %prec lowerThanCharsetKwd: CHARSET and
// ENCRYPTION are shifted into DBName instead (so "ALTER DATABASE CHARSET
// = x" is a syntax error, exactly as in the LALR tables). PLACEMENT's
// precedence is below lowerThanCharsetKwd, so it selects the option list.
func (r *rdParser) parseAlterDatabaseStmt() ast.StmtNode {
	r.expect(alter)
	r.expect(database)
	if r.isAlterDatabaseOptionListStart() {
		return &ast.AlterDatabaseStmt{
			Name:                 ast.NewCIStr(""),
			AlterDefaultDatabase: true,
			Options:              r.parseDatabaseOptionList(),
		}
	}
	name := r.parseIdentifier()
	return &ast.AlterDatabaseStmt{
		Name:                 ast.NewCIStr(name),
		AlterDefaultDatabase: false,
		Options:              r.parseDatabaseOptionList(),
	}
}

// isAlterDatabaseOptionListStart reports whether the current token begins
// the DatabaseOptionList of the no-name AlterDatabaseStmt alternative (see
// parseAlterDatabaseStmt for why CHARSET/ENCRYPTION are excluded).
func (r *rdParser) isAlterDatabaseOptionListStart() bool {
	switch r.tok() {
	case defaultKwd, collate, set, placement, character, charType:
		return true
	}
	return false
}

// parseDatabaseOptionList implements DatabaseOptionList (one or more
// juxtaposed options; the CREATE-side helpers own DatabaseOption itself).
func (r *rdParser) parseDatabaseOptionList() []*ast.DatabaseOption {
	opts := []*ast.DatabaseOption{r.parseDatabaseOption()}
	for r.isDatabaseOptionStart() {
		opts = append(opts, r.parseDatabaseOption())
	}
	return opts
}

/**************************************AlterUserStmt***************************************/

// parseAlterUserStmt implements AlterUserStmt.
func (r *rdParser) parseAlterUserStmt() ast.StmtNode {
	r.expect(alter)
	r.expect(user)
	ifExists := r.parseIfExists()
	if r.tok() == user && r.la(1) == int('(') {
		// "ALTER" "USER" IfExists "USER" '(' ')' "IDENTIFIED" "BY" AuthString
		r.advance()
		r.expect(int('('))
		r.expect(int(')'))
		r.expect(identified)
		r.expect(by)
		auth := &ast.AuthOption{
			AuthString:   r.parseAuthString(),
			ByAuthString: true,
		}
		return &ast.AlterUserStmt{
			IfExists:    ifExists,
			CurrentAuth: auth,
		}
	}
	// "ALTER" "USER" IfExists UserSpecList RequireClauseOpt ConnectionOptions
	// PasswordOrLockOptions CommentOrAttributeOption ResourceGroupNameOption
	specs := r.parseGrantUserSpecList()
	tlsOptions := r.parseGrantRequireClauseOpt()
	resourceOptions := r.parseUserConnectionOptions()
	passwordOrLockOptions := r.parseUserPasswordOrLockOptions()
	ret := &ast.AlterUserStmt{
		IfExists:              ifExists,
		Specs:                 specs,
		AuthTokenOrTLSOptions: tlsOptions,
		ResourceOptions:       resourceOptions,
		PasswordOrLockOptions: passwordOrLockOptions,
	}
	if opt := r.parseUserCommentOrAttributeOption(); opt != nil {
		ret.CommentOrAttributeOption = opt
	}
	if opt := r.parseUserResourceGroupNameOption(); opt != nil {
		ret.ResourceGroupNameOption = opt
	}
	return ret
}

/**************************************Other ALTER kinds***************************************/

// parseAlterInstanceStmt implements AlterInstanceStmt and InstanceOption.
func (r *rdParser) parseAlterInstanceStmt() ast.StmtNode {
	r.expect(alter)
	r.expect(instance)
	// InstanceOption: "RELOAD" "TLS" ["NO" "ROLLBACK" "ON" "ERROR"]
	r.expect(reload)
	r.expect(tls)
	if r.accept(no) {
		r.expect(rollback)
		r.expect(on)
		r.expect(errorKwd)
		return &ast.AlterInstanceStmt{
			ReloadTLS:         true,
			NoRollbackOnError: true,
		}
	}
	return &ast.AlterInstanceStmt{
		ReloadTLS: true,
	}
}

// parseAlterRangeStmt implements AlterRangeStmt:
// "ALTER" "RANGE" Identifier PlacementPolicyOption.
func (r *rdParser) parseAlterRangeStmt() ast.StmtNode {
	r.expect(alter)
	r.expect(rangeKwd)
	name := r.parseIdentifier()
	option := r.parsePlacementPolicyOption()
	return &ast.AlterRangeStmt{RangeName: ast.NewCIStr(name), PlacementOption: option}
}

// parseAlterPolicyStmt implements AlterPolicyStmt:
// "ALTER" "PLACEMENT" "POLICY" IfExists PolicyName PlacementOptionList.
func (r *rdParser) parseAlterPolicyStmt() ast.StmtNode {
	r.expect(alter)
	r.expect(placement)
	r.expect(policy)
	ifExists := r.parseIfExists()
	name := r.parseIdentifier()
	return &ast.AlterPlacementPolicyStmt{
		IfExists:         ifExists,
		PolicyName:       ast.NewCIStr(name),
		PlacementOptions: r.parsePlacementOptionList(),
	}
}

// parseAlterResourceGroupStmt implements AlterResourceGroupStmt:
// "ALTER" "RESOURCE" "GROUP" IfExists ResourceGroupName ResourceGroupOptionList.
func (r *rdParser) parseAlterResourceGroupStmt() ast.StmtNode {
	r.expect(alter)
	r.expect(resource)
	r.expect(group)
	ifExists := r.parseIfExists()
	name := r.parseResourceGroupName()
	return &ast.AlterResourceGroupStmt{
		IfExists:                ifExists,
		ResourceGroupName:       ast.NewCIStr(name),
		ResourceGroupOptionList: r.parseResourceGroupOptionList(),
	}
}

/**************************************AlterSequenceStmt***************************************/

// parseAlterSequenceStmt implements AlterSequenceStmt:
// "ALTER" "SEQUENCE" IfExists TableName AlterSequenceOptionList.
func (r *rdParser) parseAlterSequenceStmt() ast.StmtNode {
	r.expect(alter)
	r.expect(sequence)
	ifExists := r.parseIfExists()
	name := r.parseTableName()
	opts := []*ast.SequenceOption{r.parseAlterSequenceOption()}
	for r.isAlterSequenceOptionStart() {
		opts = append(opts, r.parseAlterSequenceOption())
	}
	return &ast.AlterSequenceStmt{
		IfExists:   ifExists,
		Name:       name,
		SeqOptions: opts,
	}
}

func (r *rdParser) isAlterSequenceOptionStart() bool {
	return r.tok() == restart || r.isSequenceOptionStart()
}

// parseAlterSequenceOption implements AlterSequenceOption: SequenceOption
// plus the RESTART forms.
func (r *rdParser) parseAlterSequenceOption() *ast.SequenceOption {
	if r.tok() != restart {
		return r.parseSequenceOption()
	}
	r.advance()
	switch r.tok() {
	case with:
		// "RESTART" "WITH" SignedNum
		r.advance()
		return &ast.SequenceOption{Tp: ast.SequenceRestartWith, IntValue: r.parseSignedNum()}
	case eq, intLit, int('+'), int('-'):
		// "RESTART" EqOpt SignedNum
		r.parseEqOpt()
		return &ast.SequenceOption{Tp: ast.SequenceRestartWith, IntValue: r.parseSignedNum()}
	}
	// "RESTART"
	return &ast.SequenceOption{Tp: ast.SequenceRestart}
}
