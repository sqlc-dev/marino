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

// AnalyzeTableStmt and its helper productions. The AnalyzeOption(List)
// helpers already live in parse_alter.go (ALTER TABLE ... ANALYZE
// PARTITION shares them).

import (
	"github.com/sqlc-dev/marino/ast"
)

func init() {
	rdRegister(analyze, (*rdParser).parseAnalyzeTableStmt)
}

// parseAnalyzeTableStmt implements AnalyzeTableStmt (every alternative).
func (r *rdParser) parseAnalyzeTableStmt() ast.StmtNode {
	r.expect(analyze)
	noWriteToBinLog := r.parseNoWriteToBinLogAliasOpt()
	if r.accept(incremental) {
		// "ANALYZE" NoWriteToBinLogAliasOpt "INCREMENTAL" "TABLE" TableName
		// ["PARTITION" PartitionNameList] "INDEX" IndexNameList
		// AnalyzeOptionListOpt
		r.expect(tableKwd)
		table := r.parseTableName()
		var partitionNames []ast.CIStr
		if r.accept(partition) {
			partitionNames = r.parsePartitionNameList()
		}
		r.expect(index)
		return &ast.AnalyzeTableStmt{
			TableNames:      []*ast.TableName{table},
			NoWriteToBinLog: noWriteToBinLog,
			PartitionNames:  partitionNames,
			IndexNames:      r.parseIndexNameList(),
			IndexFlag:       true,
			Incremental:     true,
			AnalyzeOpts:     r.parseAnalyzeOptionListOpt(),
		}
	}
	r.expect(tableKwd)
	table := r.parseTableName()
	if r.tok() == int(',') {
		// "ANALYZE" NoWriteToBinLogAliasOpt "TABLE" TableNameList
		// AllColumnsOrPredicateColumnsOpt AnalyzeOptionListOpt
		tables := []*ast.TableName{table}
		for r.accept(int(',')) {
			tables = append(tables, r.parseTableName())
		}
		return &ast.AnalyzeTableStmt{
			TableNames:      tables,
			NoWriteToBinLog: noWriteToBinLog,
			ColumnChoice:    r.parseAllColumnsOrPredicateColumnsOpt(),
			AnalyzeOpts:     r.parseAnalyzeOptionListOpt(),
		}
	}
	switch r.tok() {
	case index:
		// "ANALYZE" NoWriteToBinLogAliasOpt "TABLE" TableName "INDEX"
		// IndexNameList AnalyzeOptionListOpt
		r.advance()
		return &ast.AnalyzeTableStmt{
			TableNames:      []*ast.TableName{table},
			NoWriteToBinLog: noWriteToBinLog,
			IndexNames:      r.parseIndexNameList(),
			IndexFlag:       true,
			AnalyzeOpts:     r.parseAnalyzeOptionListOpt(),
		}
	case partition:
		r.advance()
		partitionNames := r.parsePartitionNameList()
		switch r.tok() {
		case index:
			// ... "PARTITION" PartitionNameList "INDEX" IndexNameList
			// AnalyzeOptionListOpt
			r.advance()
			return &ast.AnalyzeTableStmt{
				TableNames:      []*ast.TableName{table},
				NoWriteToBinLog: noWriteToBinLog,
				PartitionNames:  partitionNames,
				IndexNames:      r.parseIndexNameList(),
				IndexFlag:       true,
				AnalyzeOpts:     r.parseAnalyzeOptionListOpt(),
			}
		case columns:
			// ... "PARTITION" PartitionNameList "COLUMNS" IdentList
			// AnalyzeOptionListOpt
			r.advance()
			return &ast.AnalyzeTableStmt{
				TableNames:      []*ast.TableName{table},
				NoWriteToBinLog: noWriteToBinLog,
				PartitionNames:  partitionNames,
				ColumnNames:     r.parseIdentList(),
				ColumnChoice:    ast.ColumnList,
				AnalyzeOpts:     r.parseAnalyzeOptionListOpt(),
			}
		}
		// ... "PARTITION" PartitionNameList AllColumnsOrPredicateColumnsOpt
		// AnalyzeOptionListOpt
		return &ast.AnalyzeTableStmt{
			TableNames:      []*ast.TableName{table},
			NoWriteToBinLog: noWriteToBinLog,
			PartitionNames:  partitionNames,
			ColumnChoice:    r.parseAllColumnsOrPredicateColumnsOpt(),
			AnalyzeOpts:     r.parseAnalyzeOptionListOpt(),
		}
	case update:
		// ... "UPDATE" "HISTOGRAM" "ON" IdentList
		//     ("USING" "DATA" stringLit
		//      | AnalyzeOptionListOpt [("AUTO"|"MANUAL") "UPDATE"])
		r.advance()
		r.expect(histogram)
		r.expect(on)
		stmt := &ast.AnalyzeTableStmt{
			TableNames:         []*ast.TableName{table},
			NoWriteToBinLog:    noWriteToBinLog,
			ColumnNames:        r.parseIdentList(),
			HistogramOperation: ast.HistogramOperationUpdate,
		}
		if r.tok() == using {
			r.advance()
			r.expect(data)
			stmt.HistogramData = r.expect(stringLit).lit
			return stmt
		}
		stmt.AnalyzeOpts = r.parseAnalyzeOptionListOpt()
		switch r.tok() {
		case auto:
			r.advance()
			r.expect(update)
			stmt.HistogramUpdate = ast.HistogramUpdateAuto
		case manual:
			r.advance()
			r.expect(update)
			stmt.HistogramUpdate = ast.HistogramUpdateManual
		}
		return stmt
	case drop:
		// ... "DROP" "HISTOGRAM" "ON" IdentList
		r.advance()
		r.expect(histogram)
		r.expect(on)
		return &ast.AnalyzeTableStmt{
			TableNames:         []*ast.TableName{table},
			NoWriteToBinLog:    noWriteToBinLog,
			ColumnNames:        r.parseIdentList(),
			HistogramOperation: ast.HistogramOperationDrop,
		}
	case columns:
		// ... "COLUMNS" IdentList AnalyzeOptionListOpt
		r.advance()
		return &ast.AnalyzeTableStmt{
			TableNames:      []*ast.TableName{table},
			NoWriteToBinLog: noWriteToBinLog,
			ColumnNames:     r.parseIdentList(),
			ColumnChoice:    ast.ColumnList,
			AnalyzeOpts:     r.parseAnalyzeOptionListOpt(),
		}
	}
	// "ANALYZE" NoWriteToBinLogAliasOpt "TABLE" TableNameList (of one)
	// AllColumnsOrPredicateColumnsOpt AnalyzeOptionListOpt
	return &ast.AnalyzeTableStmt{
		TableNames:      []*ast.TableName{table},
		NoWriteToBinLog: noWriteToBinLog,
		ColumnChoice:    r.parseAllColumnsOrPredicateColumnsOpt(),
		AnalyzeOpts:     r.parseAnalyzeOptionListOpt(),
	}
}

// parseAllColumnsOrPredicateColumnsOpt implements
// AllColumnsOrPredicateColumnsOpt.
func (r *rdParser) parseAllColumnsOrPredicateColumnsOpt() ast.ColumnChoice {
	switch r.tok() {
	case all:
		// "ALL" "COLUMNS"
		r.advance()
		r.expect(columns)
		return ast.AllColumns
	case predicate:
		// "PREDICATE" "COLUMNS"
		r.advance()
		r.expect(columns)
		return ast.PredicateColumns
	}
	return ast.DefaultChoice
}

// parseIdentList implements IdentList.
func (r *rdParser) parseIdentList() []ast.CIStr {
	list := []ast.CIStr{ast.NewCIStr(r.parseIdentifier())}
	for r.accept(int(',')) {
		list = append(list, ast.NewCIStr(r.parseIdentifier()))
	}
	return list
}
