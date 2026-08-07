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

// The BRIE (backup/restore/import/export) statement family: the
// BRIEStmt alternatives led by BACKUP, RESTORE, STOP, PAUSE, RESUME and
// PURGE, plus BRIETables, DBNameList, BRIEOptions and the option-name /
// value helper productions. The "SHOW ..." alternatives are led by the
// show token (parse_show.go's family) and the "CANCEL BR JOB"
// alternative by the cancel token (parse_tidb.go's family).

import (
	"strings"

	"github.com/sqlc-dev/marino/ast"
)

func init() {
	rdRegister(backup, (*rdParser).parseBackupStmt)
	rdRegister(restore, (*rdParser).parseRestoreStmt)
	rdRegister(stop, (*rdParser).parseStopBackupStmt)
	rdRegister(pause, (*rdParser).parsePauseBackupStmt)
	rdRegister(resume, (*rdParser).parseResumeBackupStmt)
	rdRegister(purge, (*rdParser).parsePurgeBackupStmt)
}

// parseBackupStmt implements the BRIEStmt alternatives
//
//	"BACKUP" BRIETables "TO" stringLit BRIEOptions
//	"BACKUP" "LOGS" "TO" stringLit BRIEOptions
func (r *rdParser) parseBackupStmt() ast.StmtNode {
	r.expect(backup)
	if r.accept(logs) {
		stmt := &ast.BRIEStmt{}
		stmt.Kind = ast.BRIEKindStreamStart
		r.expect(to)
		stmt.Storage = r.expect(stringLit).lit
		stmt.Options = r.parseBRIEOptions()
		return stmt
	}
	stmt := r.parseBRIETables()
	stmt.Kind = ast.BRIEKindBackup
	r.expect(to)
	stmt.Storage = r.expect(stringLit).lit
	stmt.Options = r.parseBRIEOptions()
	return stmt
}

// parseRestoreStmt implements the BRIEStmt alternatives
//
//	"RESTORE" BRIETables "FROM" stringLit BRIEOptions
//	"RESTORE" "POINT" "FROM" stringLit BRIEOptions
func (r *rdParser) parseRestoreStmt() ast.StmtNode {
	r.expect(restore)
	if r.accept(point) {
		stmt := &ast.BRIEStmt{}
		stmt.Kind = ast.BRIEKindRestorePIT
		r.expect(from)
		stmt.Storage = r.expect(stringLit).lit
		stmt.Options = r.parseBRIEOptions()
		return stmt
	}
	stmt := r.parseBRIETables()
	stmt.Kind = ast.BRIEKindRestore
	r.expect(from)
	stmt.Storage = r.expect(stringLit).lit
	stmt.Options = r.parseBRIEOptions()
	return stmt
}

// parseStopBackupStmt implements BRIEStmt: "STOP" "BACKUP" "LOGS".
func (r *rdParser) parseStopBackupStmt() ast.StmtNode {
	r.expect(stop)
	r.expect(backup)
	r.expect(logs)
	stmt := &ast.BRIEStmt{}
	stmt.Kind = ast.BRIEKindStreamStop
	return stmt
}

// parsePauseBackupStmt implements BRIEStmt:
// "PAUSE" "BACKUP" "LOGS" BRIEOptions.
func (r *rdParser) parsePauseBackupStmt() ast.StmtNode {
	r.expect(pause)
	r.expect(backup)
	r.expect(logs)
	stmt := &ast.BRIEStmt{}
	stmt.Kind = ast.BRIEKindStreamPause
	stmt.Options = r.parseBRIEOptions()
	return stmt
}

// parseResumeBackupStmt implements BRIEStmt: "RESUME" "BACKUP" "LOGS".
func (r *rdParser) parseResumeBackupStmt() ast.StmtNode {
	r.expect(resume)
	r.expect(backup)
	r.expect(logs)
	stmt := &ast.BRIEStmt{}
	stmt.Kind = ast.BRIEKindStreamResume
	return stmt
}

// parsePurgeBackupStmt implements BRIEStmt:
// "PURGE" "BACKUP" "LOGS" "FROM" stringLit BRIEOptions.
func (r *rdParser) parsePurgeBackupStmt() ast.StmtNode {
	r.expect(purge)
	r.expect(backup)
	r.expect(logs)
	r.expect(from)
	stmt := &ast.BRIEStmt{}
	stmt.Kind = ast.BRIEKindStreamPurge
	stmt.Storage = r.expect(stringLit).lit
	stmt.Options = r.parseBRIEOptions()
	return stmt
}

// parseBRIETables implements BRIETables (and DBNameList).
func (r *rdParser) parseBRIETables() *ast.BRIEStmt {
	switch r.tok() {
	case database:
		// DatabaseSym ('*' | DBNameList)
		r.advance()
		if r.accept(int('*')) {
			return &ast.BRIEStmt{}
		}
		schemas := []string{r.parseIdentifier()}
		for r.accept(int(',')) {
			schemas = append(schemas, r.parseIdentifier())
		}
		return &ast.BRIEStmt{Schemas: schemas}
	case tableKwd:
		// "TABLE" TableNameList
		r.advance()
		return &ast.BRIEStmt{Tables: r.parseTableNameList()}
	}
	r.syntaxError()
	return nil
}

// parseBRIEOptions implements BRIEOptions: a possibly-empty run of
// BRIEOption with no separators.
func (r *rdParser) parseBRIEOptions() []*ast.BRIEOption {
	options := []*ast.BRIEOption{}
	for {
		switch r.tok() {
		case concurrency, resume, checksumConcurrency, compressionLevel,
			sendCredentialsToTiKV, online, checkpoint, skipSchemaFiles,
			strictFormat, csvNotNull, csvBackslashEscape,
			csvTrimLastSeparators, waitTiflashReady, withSysTable,
			ignoreStats, loadStats, tikvImporter, csvSeparator, csvDelimiter,
			csvNull, compressionType, encryptionMethod, encryptionKeyFile,
			backend, onDuplicate, on, snapshot, lastBackup, rateLimit,
			csvHeader, checksum, analyze, fullBackupStorage, restoredTS,
			startTS, untilTS, gcTTL:
			options = append(options, r.parseBRIEOption())
		default:
			return options
		}
	}
}

// parseBRIEOption implements BRIEOption together with
// BRIEIntegerOptionName, BRIEBooleanOptionName, BRIEStringOptionName and
// BRIEKeywordOptionName.
func (r *rdParser) parseBRIEOption() *ast.BRIEOption {
	switch r.tok() {
	case concurrency, resume, checksumConcurrency, compressionLevel:
		// BRIEIntegerOptionName EqOpt LengthNum
		var tp ast.BRIEOptionType
		switch r.tok() {
		case concurrency:
			tp = ast.BRIEOptionConcurrency
		case resume:
			tp = ast.BRIEOptionResume
		case checksumConcurrency:
			tp = ast.BRIEOptionChecksumConcurrency
		case compressionLevel:
			tp = ast.BRIEOptionCompressionLevel
		}
		r.advance()
		r.parseEqOpt()
		return &ast.BRIEOption{
			Tp:        tp,
			UintValue: r.parseLengthNum(),
		}
	case sendCredentialsToTiKV, online, checkpoint, skipSchemaFiles,
		strictFormat, csvNotNull, csvBackslashEscape, csvTrimLastSeparators,
		waitTiflashReady, withSysTable, ignoreStats, loadStats:
		// BRIEBooleanOptionName EqOpt Boolean
		var tp ast.BRIEOptionType
		switch r.tok() {
		case sendCredentialsToTiKV:
			tp = ast.BRIEOptionSendCreds
		case online:
			tp = ast.BRIEOptionOnline
		case checkpoint:
			tp = ast.BRIEOptionCheckpoint
		case skipSchemaFiles:
			tp = ast.BRIEOptionSkipSchemaFiles
		case strictFormat:
			tp = ast.BRIEOptionStrictFormat
		case csvNotNull:
			tp = ast.BRIEOptionCSVNotNull
		case csvBackslashEscape:
			tp = ast.BRIEOptionCSVBackslashEscape
		case csvTrimLastSeparators:
			tp = ast.BRIEOptionCSVTrimLastSeparators
		case waitTiflashReady:
			tp = ast.BRIEOptionWaitTiflashReady
		case withSysTable:
			tp = ast.BRIEOptionWithSysTable
		case ignoreStats:
			tp = ast.BRIEOptionIgnoreStats
		case loadStats:
			tp = ast.BRIEOptionLoadStats
		}
		r.advance()
		r.parseEqOpt()
		value := uint64(0)
		if r.parseBoolean() {
			value = 1
		}
		return &ast.BRIEOption{
			Tp:        tp,
			UintValue: value,
		}
	case tikvImporter, csvSeparator, csvDelimiter, csvNull, compressionType,
		encryptionMethod, encryptionKeyFile:
		// BRIEStringOptionName EqOpt stringLit
		var tp ast.BRIEOptionType
		switch r.tok() {
		case tikvImporter:
			tp = ast.BRIEOptionTiKVImporter
		case csvSeparator:
			tp = ast.BRIEOptionCSVSeparator
		case csvDelimiter:
			tp = ast.BRIEOptionCSVDelimiter
		case csvNull:
			tp = ast.BRIEOptionCSVNull
		case compressionType:
			tp = ast.BRIEOptionCompression
		case encryptionMethod:
			tp = ast.BRIEOptionEncryptionMethod
		case encryptionKeyFile:
			tp = ast.BRIEOptionEncryptionKeyFile
		}
		r.advance()
		r.parseEqOpt()
		return &ast.BRIEOption{
			Tp:       tp,
			StrValue: r.expect(stringLit).lit,
		}
	case backend, onDuplicate, on:
		// BRIEKeywordOptionName EqOpt StringNameOrBRIEOptionKeyword
		var tp ast.BRIEOptionType
		switch r.tok() {
		case backend:
			tp = ast.BRIEOptionBackend
		case onDuplicate:
			tp = ast.BRIEOptionOnDuplicate
		case on:
			// "ON" "DUPLICATE"
			tp = ast.BRIEOptionOnDuplicate
			r.advance()
			r.expect(duplicate)
			r.parseEqOpt()
			return &ast.BRIEOption{
				Tp:       tp,
				StrValue: strings.ToLower(r.parseStringNameOrBRIEOptionKeyword()),
			}
		}
		r.advance()
		r.parseEqOpt()
		return &ast.BRIEOption{
			Tp:       tp,
			StrValue: strings.ToLower(r.parseStringNameOrBRIEOptionKeyword()),
		}
	case snapshot:
		// "SNAPSHOT" EqOpt (LengthNum TimestampUnit "AGO" | stringLit |
		// LengthNum)
		r.advance()
		r.parseEqOpt()
		if r.tok() == stringLit {
			return &ast.BRIEOption{
				Tp:       ast.BRIEOptionBackupTS,
				StrValue: r.expect(stringLit).lit,
			}
		}
		n := r.parseLengthNum()
		switch r.tok() {
		case microsecond, second, minute, hour, day, week, month, quarter, yearType:
			unitType := r.parseTimestampUnit()
			r.expect(ago)
			unit, err := unitType.Duration()
			if err != nil {
				r.actionError(err)
			}
			// TODO: check overflow?
			return &ast.BRIEOption{
				Tp:        ast.BRIEOptionBackupTimeAgo,
				UintValue: n * uint64(unit),
			}
		}
		return &ast.BRIEOption{
			Tp:        ast.BRIEOptionBackupTSO,
			UintValue: n,
		}
	case lastBackup:
		// "LAST_BACKUP" EqOpt (stringLit | LengthNum)
		r.advance()
		r.parseEqOpt()
		if r.tok() == stringLit {
			return &ast.BRIEOption{
				Tp:       ast.BRIEOptionLastBackupTS,
				StrValue: r.expect(stringLit).lit,
			}
		}
		return &ast.BRIEOption{
			Tp:        ast.BRIEOptionLastBackupTSO,
			UintValue: r.parseLengthNum(),
		}
	case rateLimit:
		// "RATE_LIMIT" EqOpt LengthNum "MB" '/' "SECOND"
		r.advance()
		r.parseEqOpt()
		n := r.parseLengthNum()
		r.expect(mb)
		r.expect(int('/'))
		r.expect(second)
		// TODO: check overflow?
		return &ast.BRIEOption{
			Tp:        ast.BRIEOptionRateLimit,
			UintValue: n * 1048576,
		}
	case csvHeader:
		// "CSV_HEADER" EqOpt (FieldsOrColumns | LengthNum)
		r.advance()
		r.parseEqOpt()
		if r.tok() == fields || r.tok() == columns {
			r.advance()
			return &ast.BRIEOption{
				Tp:        ast.BRIEOptionCSVHeader,
				UintValue: ast.BRIECSVHeaderIsColumns,
			}
		}
		return &ast.BRIEOption{
			Tp:        ast.BRIEOptionCSVHeader,
			UintValue: r.parseLengthNum(),
		}
	case checksum:
		// "CHECKSUM" EqOpt (Boolean | OptionLevel)
		r.advance()
		r.parseEqOpt()
		if r.tok() == off || r.tok() == optional || r.tok() == required {
			return &ast.BRIEOption{
				Tp:        ast.BRIEOptionChecksum,
				UintValue: uint64(r.parseBRIEOptionLevel()),
			}
		}
		value := uint64(0)
		if r.parseBoolean() {
			value = 1
		}
		return &ast.BRIEOption{
			Tp:        ast.BRIEOptionChecksum,
			UintValue: value,
		}
	case analyze:
		// "ANALYZE" EqOpt (Boolean | OptionLevel)
		r.advance()
		r.parseEqOpt()
		if r.tok() == off || r.tok() == optional || r.tok() == required {
			return &ast.BRIEOption{
				Tp:        ast.BRIEOptionAnalyze,
				UintValue: uint64(r.parseBRIEOptionLevel()),
			}
		}
		value := uint64(0)
		if r.parseBoolean() {
			value = 1
		}
		return &ast.BRIEOption{
			Tp:        ast.BRIEOptionAnalyze,
			UintValue: value,
		}
	case fullBackupStorage:
		// "FULL_BACKUP_STORAGE" EqOpt stringLit
		r.advance()
		r.parseEqOpt()
		return &ast.BRIEOption{
			Tp:       ast.BRIEOptionFullBackupStorage,
			StrValue: r.expect(stringLit).lit,
		}
	case restoredTS:
		// "RESTORED_TS" EqOpt stringLit
		r.advance()
		r.parseEqOpt()
		return &ast.BRIEOption{
			Tp:       ast.BRIEOptionRestoredTS,
			StrValue: r.expect(stringLit).lit,
		}
	case startTS:
		// "START_TS" EqOpt stringLit
		r.advance()
		r.parseEqOpt()
		return &ast.BRIEOption{
			Tp:       ast.BRIEOptionStartTS,
			StrValue: r.expect(stringLit).lit,
		}
	case untilTS:
		// "UNTIL_TS" EqOpt stringLit
		r.advance()
		r.parseEqOpt()
		return &ast.BRIEOption{
			Tp:       ast.BRIEOptionUntilTS,
			StrValue: r.expect(stringLit).lit,
		}
	case gcTTL:
		// "GC_TTL" EqOpt stringLit
		r.advance()
		r.parseEqOpt()
		return &ast.BRIEOption{
			Tp:       ast.BRIEOptionGCTTL,
			StrValue: r.expect(stringLit).lit,
		}
	}
	r.syntaxError()
	return nil
}

// parseStringNameOrBRIEOptionKeyword implements
// StringNameOrBRIEOptionKeyword: StringName | "IGNORE" | "REPLACE".
func (r *rdParser) parseStringNameOrBRIEOptionKeyword() string {
	if r.tok() == ignore || r.tok() == replace {
		lit := r.cur().lit
		r.advance()
		return lit
	}
	return r.parseStringName()
}

// parseBoolean implements Boolean: NUM | "FALSE" | "TRUE".
func (r *rdParser) parseBoolean() bool {
	switch r.tok() {
	case intLit:
		v := r.cur().item
		r.advance()
		return v.(int64) != 0
	case falseKwd:
		r.advance()
		return false
	case trueKwd:
		r.advance()
		return true
	}
	r.syntaxError()
	return false
}

// parseBRIEOptionLevel implements OptionLevel:
// "OFF" | "OPTIONAL" | "REQUIRED".
func (r *rdParser) parseBRIEOptionLevel() ast.BRIEOptionLevel {
	switch r.tok() {
	case off:
		r.advance()
		return ast.BRIEOptionLevelOff
	case optional:
		r.advance()
		return ast.BRIEOptionLevelOptional
	case required:
		r.advance()
		return ast.BRIEOptionLevelRequired
	}
	r.syntaxError()
	return 0
}
