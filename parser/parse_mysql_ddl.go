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

// The MySQL data definition statements that postdate the goyacc grammar
// (MySQL 26.7 §15.1): events, triggers, stored functions (and the ALTER
// forms of stored routines), ALTER VIEW, servers, tablespaces, logfile
// groups, spatial reference systems, libraries, and JSON duality views.
// The productions are written from the MySQL 26.7 reference manual in
// the same style as parser.y. The CREATE/ALTER/DROP statement heads
// dispatch here from parse_create_table.go, parse_alter.go, and
// parse_drop.go.

import (
	"strings"

	"github.com/sqlc-dev/marino/ast"
	"github.com/sqlc-dev/marino/auth"
)

// parseDefinerOpt implements DefinerOpt: empty | "DEFINER" eq Username.
// It returns nil for the empty alternative.
func (r *rdParser) parseDefinerOpt() *auth.UserIdentity {
	if r.tok() != definer {
		return nil
	}
	r.advance()
	r.expect(eq)
	return r.parseUsername()
}

// peekPastDefiner reports the token following the "DEFINER" eq Username
// clause that begins one token past the cursor (which sits on CREATE or
// ALTER), consuming nothing. It returns -1 when the clause does not
// parse; the caller's default alternative then reports the error.
func (r *rdParser) peekPastDefiner() int {
	tok := -1
	m := r.mark()
	r.try(func() {
		r.advance() // CREATE or ALTER
		r.advance() // DEFINER
		r.expect(eq)
		r.parseUsername()
		tok = r.tok()
	})
	r.rewind(m)
	return tok
}

// expectIdentLit consumes the current token, which must occupy
// identifier position and match name case-insensitively. It is used for
// the words of these productions that are not keywords (non-reserved in
// MySQL).
func (r *rdParser) expectIdentLit(name string) {
	if !isIdentifierTok(r.tok()) || !strings.EqualFold(r.cur().lit, name) {
		r.syntaxError()
	}
	r.advance()
}

/*
 * Events
 */

// parseEventSchedule implements the schedule of ON SCHEDULE:
//
//	"AT" Expression
//	| "EVERY" Expression TimeUnit ["STARTS" Expression]
//	    ["ENDS" Expression]
func (r *rdParser) parseEventSchedule() *ast.EventSchedule {
	if r.accept(at) {
		return &ast.EventSchedule{At: r.parseExpression()}
	}
	r.expect(every)
	sched := &ast.EventSchedule{Every: r.parseExpression(), Unit: r.parseTimeUnit()}
	if r.accept(starts) {
		sched.Starts = r.parseExpression()
	}
	if r.accept(ends) {
		sched.Ends = r.parseExpression()
	}
	return sched
}

// parseEventCompletion implements the ON COMPLETION clause, with ON
// already consumed: "COMPLETION" ["NOT"] "PRESERVE".
func (r *rdParser) parseEventCompletion() ast.EventCompletion {
	r.expect(completion)
	completionV := ast.EventCompletionPreserve
	if r.accept(not) {
		completionV = ast.EventCompletionNotPreserve
	}
	r.expect(preserve)
	return completionV
}

// parseEventStateOpt implements the event state clause:
// empty | "ENABLE" | "DISABLE" ["ON" "REPLICA"].
func (r *rdParser) parseEventStateOpt() ast.EventState {
	switch r.tok() {
	case enable:
		r.advance()
		return ast.EventStateEnable
	case disable:
		r.advance()
		if r.accept(on) {
			r.expect(replica)
			return ast.EventStateDisableOnReplica
		}
		return ast.EventStateDisable
	}
	return ast.EventStateDefault
}

// parseEventBody parses an event, trigger, or function body statement
// and records its source text the way CreateProcedureStmt does.
func (r *rdParser) parseRoutineBody() ast.StmtNode {
	bodyStart := r.cur().offset
	body := r.parseProcedureProcStmt()
	r.p.setNodeText(body, strings.TrimSpace(r.src[bodyStart:r.cur().offset]))
	return body
}

// parseCreateEventStmt implements CreateEventStmt:
//
//	"CREATE" DefinerOpt "EVENT" IfNotExists TableName
//	    "ON" "SCHEDULE" EventSchedule ["ON" "COMPLETION" ["NOT"]
//	    "PRESERVE"] EventStateOpt ["COMMENT" stringLit]
//	    "DO" ProcedureProcStmt
func (r *rdParser) parseCreateEventStmt() ast.StmtNode {
	r.expect(create)
	stmt := &ast.CreateEventStmt{Definer: r.parseDefinerOpt()}
	r.expect(event)
	stmt.IfNotExists = r.parseIfNotExists()
	stmt.EventName = r.parseTableName()
	r.expect(on)
	r.expect(schedule)
	stmt.Schedule = r.parseEventSchedule()
	if r.accept(on) {
		stmt.Completion = r.parseEventCompletion()
	}
	stmt.State = r.parseEventStateOpt()
	if r.accept(comment) {
		s := r.expect(stringLit).lit
		stmt.Comment = &s
	}
	r.expect(do)
	stmt.Body = r.parseRoutineBody()
	return stmt
}

// parseAlterEventStmt implements AlterEventStmt:
//
//	"ALTER" DefinerOpt "EVENT" TableName ["ON" "SCHEDULE" EventSchedule]
//	    ["ON" "COMPLETION" ["NOT"] "PRESERVE"] ["RENAME" "TO" TableName]
//	    EventStateOpt ["COMMENT" stringLit] ["DO" ProcedureProcStmt]
func (r *rdParser) parseAlterEventStmt() ast.StmtNode {
	r.expect(alter)
	stmt := &ast.AlterEventStmt{Definer: r.parseDefinerOpt()}
	r.expect(event)
	stmt.EventName = r.parseTableName()
	if r.tok() == on && r.la(1) == schedule {
		r.advance()
		r.advance()
		stmt.Schedule = r.parseEventSchedule()
	}
	if r.tok() == on && r.la(1) == completion {
		r.advance()
		stmt.Completion = r.parseEventCompletion()
	}
	if r.accept(rename) {
		r.expect(to)
		stmt.RenameTo = r.parseTableName()
	}
	stmt.State = r.parseEventStateOpt()
	if r.accept(comment) {
		s := r.expect(stringLit).lit
		stmt.Comment = &s
	}
	if r.accept(do) {
		stmt.Body = r.parseRoutineBody()
	}
	return stmt
}

/*
 * Triggers
 */

// parseCreateTriggerStmt implements CreateTriggerStmt:
//
//	"CREATE" DefinerOpt "TRIGGER" IfNotExists TableName
//	    ("BEFORE" | "AFTER") ("INSERT" | "UPDATE" | "DELETE")
//	    "ON" TableName "FOR" "EACH" "ROW"
//	    [("FOLLOWS" | "PRECEDES") Identifier] ProcedureProcStmt
func (r *rdParser) parseCreateTriggerStmt() ast.StmtNode {
	r.expect(create)
	stmt := &ast.CreateTriggerStmt{Definer: r.parseDefinerOpt()}
	r.expect(trigger)
	stmt.IfNotExists = r.parseIfNotExists()
	stmt.TriggerName = r.parseTableName()
	switch r.tok() {
	case before:
		stmt.TriggerTime = ast.TriggerTimeBefore
	case after:
		stmt.TriggerTime = ast.TriggerTimeAfter
	default:
		r.syntaxError()
	}
	r.advance()
	switch r.tok() {
	case insert:
		stmt.TriggerEvent = ast.TriggerEventInsert
	case update:
		stmt.TriggerEvent = ast.TriggerEventUpdate
	case deleteKwd:
		stmt.TriggerEvent = ast.TriggerEventDelete
	default:
		r.syntaxError()
	}
	r.advance()
	r.expect(on)
	stmt.Table = r.parseTableName()
	r.expect(forKwd)
	r.expect(each)
	r.expect(row)
	switch r.tok() {
	case follows:
		r.advance()
		stmt.Order = &ast.TriggerOrder{OtherTrigger: ast.NewCIStr(r.parseIdentifier())}
	case precedes:
		r.advance()
		stmt.Order = &ast.TriggerOrder{Precedes: true, OtherTrigger: ast.NewCIStr(r.parseIdentifier())}
	}
	stmt.Body = r.parseRoutineBody()
	return stmt
}

/*
 * Stored routines
 */

// parseRoutineCharacteristics implements the characteristic list of the
// stored routine statements:
//
//	("COMMENT" stringLit | "LANGUAGE" (Identifier | "SQL")
//	| ["NOT"] "DETERMINISTIC" | "CONTAINS" "SQL" | "NO" "SQL"
//	| "READS" "SQL" "DATA" | "MODIFIES" "SQL" "DATA"
//	| "SQL" "SECURITY" ("DEFINER" | "INVOKER"))*
//
// MySQL rejects [NOT] DETERMINISTIC in the ALTER statements at a later
// stage; the list is shared here and leaves that to semantics.
func (r *rdParser) parseRoutineCharacteristics() ast.RoutineCharacteristics {
	var c ast.RoutineCharacteristics
	for {
		switch r.tok() {
		case comment:
			r.advance()
			s := r.expect(stringLit).lit
			c.Comment = &s
		case language:
			r.advance()
			if r.tok() == sql {
				r.advance()
				c.Language = ast.NewCIStr("SQL")
			} else {
				c.Language = ast.NewCIStr(strings.ToUpper(r.parseIdentifier()))
			}
		case deterministic:
			r.advance()
			c.Determinism = ast.RoutineDeterministic
		case not:
			if r.la(1) != deterministic {
				return c
			}
			r.advance()
			r.advance()
			c.Determinism = ast.RoutineNotDeterministic
		case contains:
			r.advance()
			r.expect(sql)
			c.DataAccess = ast.RoutineContainsSQL
		case no:
			if r.la(1) != sql {
				return c
			}
			r.advance()
			r.advance()
			c.DataAccess = ast.RoutineNoSQL
		case reads:
			r.advance()
			r.expect(sql)
			r.expect(data)
			c.DataAccess = ast.RoutineReadsSQLData
		case modifies:
			r.advance()
			r.expect(sql)
			r.expect(data)
			c.DataAccess = ast.RoutineModifiesSQLData
		case sql:
			if r.la(1) != security {
				return c
			}
			r.advance()
			r.advance()
			var s ast.ViewSecurity
			switch r.tok() {
			case definer:
				s = ast.SecurityDefiner
			case invoker:
				s = ast.SecurityInvoker
			default:
				r.syntaxError()
			}
			r.advance()
			c.Security = &s
		default:
			return c
		}
	}
}

// parseCreateFunctionFamily dispatches CREATE FUNCTION between the
// stored function form (CreateFunctionStmt, recognized by its mandatory
// parameter list) and the loadable function form
// (CreateLoadableFunctionStmt, whose name is followed directly by
// RETURNS):
//
//	"CREATE" DefinerOpt ["AGGREGATE"] "FUNCTION" IfNotExists ...
func (r *rdParser) parseCreateFunctionFamily() ast.StmtNode {
	r.expect(create)
	definer := r.parseDefinerOpt()
	aggregateOpt := definer == nil && r.accept(aggregate)
	r.expect(function)
	ifNotExists := r.parseIfNotExists()
	name := r.parseTableName()
	if !aggregateOpt && r.tok() == int('(') {
		return r.finishCreateStoredFunctionStmt(definer, ifNotExists, name)
	}
	if definer != nil || name.Schema.O != "" {
		// The loadable form takes neither a definer nor a qualified name.
		r.syntaxError()
	}
	return r.finishCreateLoadableFunctionStmt(aggregateOpt, ifNotExists, name.Name)
}

// finishCreateStoredFunctionStmt implements the rest of
// CreateFunctionStmt after the function name:
//
//	'(' [Identifier Type (',' Identifier Type)*] ')' "RETURNS" Type
//	    RoutineCharacteristics ProcedureProcStmt
func (r *rdParser) finishCreateStoredFunctionStmt(definer *auth.UserIdentity, ifNotExists bool, name *ast.TableName) ast.StmtNode {
	stmt := &ast.CreateFunctionStmt{
		IfNotExists:  ifNotExists,
		Definer:      definer,
		FunctionName: name,
	}
	r.expect(int('('))
	if r.tok() != int(')') {
		for {
			param := &ast.FunctionParam{ParamName: r.parseIdentifier()}
			param.ParamType = r.parseType()
			stmt.Params = append(stmt.Params, param)
			if !r.accept(int(',')) {
				break
			}
		}
	}
	r.expect(int(')'))
	r.expect(returns)
	stmt.ReturnType = r.parseType()
	stmt.Characteristics = r.parseRoutineCharacteristics()
	stmt.Imports = r.parseRoutineImportsOpt()
	if r.tok() == as {
		// The AS 'code' body of a LANGUAGE JAVASCRIPT routine.
		r.advance()
		s := r.expect(stringLit).lit
		stmt.CodeBody = &s
		return stmt
	}
	stmt.Body = r.parseRoutineBody()
	return stmt
}

// parseRoutineImportsOpt implements the USING clause of LANGUAGE
// JAVASCRIPT routines:
// empty | "USING" '(' TableName ["AS" Identifier] (',' ...)* ')'.
func (r *rdParser) parseRoutineImportsOpt() []*ast.RoutineImport {
	if r.tok() != using {
		return nil
	}
	r.advance()
	r.expect(int('('))
	var imports []*ast.RoutineImport
	for {
		imp := &ast.RoutineImport{Library: r.parseTableName()}
		if r.accept(as) {
			imp.Alias = ast.NewCIStr(r.parseIdentifier())
		}
		imports = append(imports, imp)
		if !r.accept(int(',')) {
			break
		}
	}
	r.expect(int(')'))
	return imports
}

// parseAlterFunctionStmt implements AlterFunctionStmt:
// "ALTER" "FUNCTION" TableName RoutineCharacteristics.
func (r *rdParser) parseAlterFunctionStmt() ast.StmtNode {
	r.expect(alter)
	r.expect(function)
	return &ast.AlterFunctionStmt{
		FunctionName:    r.parseTableName(),
		Characteristics: r.parseRoutineCharacteristics(),
	}
}

// parseAlterProcedureStmt implements AlterProcedureStmt:
// "ALTER" "PROCEDURE" TableName RoutineCharacteristics.
func (r *rdParser) parseAlterProcedureStmt() ast.StmtNode {
	r.expect(alter)
	r.expect(procedure)
	return &ast.AlterProcedureStmt{
		ProcedureName:   r.parseTableName(),
		Characteristics: r.parseRoutineCharacteristics(),
	}
}

/*
 * ALTER VIEW
 */

// parseAlterViewStmt implements AlterViewStmt:
//
//	"ALTER" ViewAlgorithm ViewDefiner ViewSQLSecurity "VIEW" ViewName
//	    ViewFieldList "AS" CreateViewSelectOpt ViewCheckOption
//
// mirroring parseCreateViewStmt, including its recording of the select
// statement's source text.
func (r *rdParser) parseAlterViewStmt() ast.StmtNode {
	r.expect(alter)
	algorithmV := ast.AlgorithmUndefined
	if r.tok() == algorithm {
		r.advance()
		r.expect(eq)
		switch r.tok() {
		case undefined:
			algorithmV = ast.AlgorithmUndefined
		case merge:
			algorithmV = ast.AlgorithmMerge
		case temptable:
			algorithmV = ast.AlgorithmTemptable
		default:
			r.syntaxError()
		}
		r.advance()
	}
	definerV := &auth.UserIdentity{CurrentUser: true}
	if d := r.parseDefinerOpt(); d != nil {
		definerV = d
	}
	securityV := ast.SecurityDefiner
	if r.tok() == sql {
		r.advance()
		r.expect(security)
		switch r.tok() {
		case definer:
			securityV = ast.SecurityDefiner
		case invoker:
			securityV = ast.SecurityInvoker
		default:
			r.syntaxError()
		}
		r.advance()
	}
	r.expect(view)
	viewName := r.parseTableName()
	var cols []ast.CIStr
	if r.tok() == int('(') {
		r.advance()
		cols = []ast.CIStr{ast.NewCIStr(r.parseIdentifier())}
		for r.accept(int(',')) {
			cols = append(cols, ast.NewCIStr(r.parseIdentifier()))
		}
		r.expect(int(')'))
	}
	r.expect(as)
	startOffset := r.cur().offset
	selStmt := r.parseCreateViewSelectOpt()
	endOffset := r.cur().offset
	x := &ast.AlterViewStmt{
		ViewName:  viewName,
		Select:    selStmt,
		Algorithm: algorithmV,
		Definer:   definerV,
		Security:  securityV,
	}
	if cols != nil {
		x.Cols = cols
	}
	// A bare WITH CHECK OPTION means CASCADED, as in MySQL.
	if r.tok() == with {
		r.advance()
		switch r.tok() {
		case cascaded:
			x.CheckOption = ast.CheckOptionCascaded
			r.advance()
		case local:
			x.CheckOption = ast.CheckOptionLocal
			r.advance()
		default:
			x.CheckOption = ast.CheckOptionCascaded
		}
		r.expect(check)
		r.expect(option)
	} else {
		x.CheckOption = ast.CheckOptionCascaded
	}
	r.p.setNodeText(selStmt, strings.TrimSpace(r.src[startOffset:endOffset]))
	return x
}

/*
 * Servers
 */

// parseServerOptions implements the OPTIONS clause of CREATE/ALTER
// SERVER: "OPTIONS" '(' ServerOption (',' ServerOption)* ')'.
func (r *rdParser) parseServerOptions() []*ast.ServerOption {
	r.expect(options)
	r.expect(int('('))
	opts := []*ast.ServerOption{r.parseServerOption()}
	for r.accept(int(',')) {
		opts = append(opts, r.parseServerOption())
	}
	r.expect(int(')'))
	return opts
}

// parseServerOption implements ServerOption:
//
//	("HOST" | "DATABASE" | "USER" | "PASSWORD" | "SOCKET" | "OWNER")
//	    stringLit
//	| "PORT" NUM
//
// HOST, SOCKET, OWNER, and PORT are not keywords and are matched in
// identifier position.
func (r *rdParser) parseServerOption() *ast.ServerOption {
	var name string
	switch {
	case r.tok() == user || r.tok() == password || r.tok() == database:
		name = strings.ToUpper(r.cur().lit)
	case isIdentifierTok(r.tok()):
		name = strings.ToUpper(r.cur().lit)
		switch name {
		case "HOST", "SOCKET", "OWNER", "PORT":
		default:
			r.syntaxError()
		}
	default:
		r.syntaxError()
	}
	r.advance()
	if name == "PORT" {
		return &ast.ServerOption{Name: name, Value: ast.NewValueExpr(r.expect(intLit).item, "", "")}
	}
	return &ast.ServerOption{Name: name, Value: ast.NewValueExpr(r.expect(stringLit).lit, "", "")}
}

// parseCreateServerStmt implements CreateServerStmt:
//
//	"CREATE" "SERVER" Identifier "FOREIGN" "DATA" "WRAPPER" Identifier
//	    ServerOptions
func (r *rdParser) parseCreateServerStmt() ast.StmtNode {
	r.expect(create)
	r.expect(server)
	stmt := &ast.CreateServerStmt{ServerName: ast.NewCIStr(r.parseIdentifier())}
	r.expect(foreign)
	r.expect(data)
	r.expect(wrapper)
	stmt.Wrapper = ast.NewCIStr(r.parseIdentifier())
	stmt.Options = r.parseServerOptions()
	return stmt
}

// parseAlterServerStmt implements AlterServerStmt:
// "ALTER" "SERVER" Identifier ServerOptions.
func (r *rdParser) parseAlterServerStmt() ast.StmtNode {
	r.expect(alter)
	r.expect(server)
	return &ast.AlterServerStmt{
		ServerName: ast.NewCIStr(r.parseIdentifier()),
		Options:    r.parseServerOptions(),
	}
}

/*
 * Tablespaces and logfile groups
 */

// parseTablespaceOptions implements the option list shared by the
// tablespace and logfile group statements:
//
//	(("INITIAL_SIZE" | "MAX_SIZE" | "EXTENT_SIZE" | "FILE_BLOCK_SIZE"
//	  | "UNDO_BUFFER_SIZE" | "REDO_BUFFER_SIZE" | "AUTOEXTEND_SIZE"
//	  | "NODEGROUP") EqOpt LengthNum
//	| ("ENCRYPTION" | "COMMENT" | "ENGINE_ATTRIBUTE") EqOpt stringLit
//	| "ENGINE" EqOpt (Identifier | stringLit)
//	| "WAIT")*
//
// The size option names that are not keywords are matched in identifier
// position.
func (r *rdParser) parseTablespaceOptions() []*ast.TablespaceOption {
	var opts []*ast.TablespaceOption
	for {
		switch r.tok() {
		case autoextendSize:
			r.advance()
			r.parseEqOpt()
			opts = append(opts, &ast.TablespaceOption{Tp: ast.TablespaceOptionAutoextendSize, UintValue: getUint64FromNUM(r.expect(intLit).item)})
		case nodegroup:
			r.advance()
			r.parseEqOpt()
			opts = append(opts, &ast.TablespaceOption{Tp: ast.TablespaceOptionNodegroup, UintValue: getUint64FromNUM(r.expect(intLit).item)})
		case wait:
			r.advance()
			opts = append(opts, &ast.TablespaceOption{Tp: ast.TablespaceOptionWait})
		case encryption:
			r.advance()
			r.parseEqOpt()
			opts = append(opts, &ast.TablespaceOption{Tp: ast.TablespaceOptionEncryption, StrValue: r.expect(stringLit).lit})
		case comment:
			r.advance()
			r.parseEqOpt()
			opts = append(opts, &ast.TablespaceOption{Tp: ast.TablespaceOptionComment, StrValue: r.expect(stringLit).lit})
		case engine_attribute:
			r.advance()
			r.parseEqOpt()
			opts = append(opts, &ast.TablespaceOption{Tp: ast.TablespaceOptionEngineAttribute, StrValue: r.expect(stringLit).lit})
		case engine:
			r.advance()
			r.parseEqOpt()
			var name string
			if r.tok() == stringLit {
				name = r.cur().lit
				r.advance()
			} else {
				name = r.parseIdentifier()
			}
			opts = append(opts, &ast.TablespaceOption{Tp: ast.TablespaceOptionEngine, StrValue: name})
		default:
			var tp ast.TablespaceOptionType
			switch {
			case !isIdentifierTok(r.tok()):
				return opts
			case strings.EqualFold(r.cur().lit, "INITIAL_SIZE"):
				tp = ast.TablespaceOptionInitialSize
			case strings.EqualFold(r.cur().lit, "MAX_SIZE"):
				tp = ast.TablespaceOptionMaxSize
			case strings.EqualFold(r.cur().lit, "EXTENT_SIZE"):
				tp = ast.TablespaceOptionExtentSize
			case strings.EqualFold(r.cur().lit, "FILE_BLOCK_SIZE"):
				tp = ast.TablespaceOptionFileBlockSize
			case strings.EqualFold(r.cur().lit, "UNDO_BUFFER_SIZE"):
				tp = ast.TablespaceOptionUndoBufferSize
			case strings.EqualFold(r.cur().lit, "REDO_BUFFER_SIZE"):
				tp = ast.TablespaceOptionRedoBufferSize
			default:
				return opts
			}
			r.advance()
			r.parseEqOpt()
			opts = append(opts, &ast.TablespaceOption{Tp: tp, UintValue: getUint64FromNUM(r.expect(intLit).item)})
		}
	}
}

// parseCreateTablespaceStmt implements CreateTablespaceStmt:
//
//	"CREATE" ["UNDO"] "TABLESPACE" Identifier
//	    ["ADD" "DATAFILE" stringLit] ["USE" "LOGFILE" "GROUP" Identifier]
//	    TablespaceOptions
func (r *rdParser) parseCreateTablespaceStmt() ast.StmtNode {
	r.expect(create)
	stmt := &ast.CreateTablespaceStmt{Undo: r.accept(undo)}
	r.expect(tablespace)
	stmt.Name = ast.NewCIStr(r.parseIdentifier())
	if r.accept(add) {
		r.expect(datafile)
		stmt.Datafile = r.expect(stringLit).lit
	}
	if r.accept(use) {
		r.expect(logfile)
		r.expect(group)
		stmt.UseLogfileGroup = ast.NewCIStr(r.parseIdentifier())
	}
	stmt.Options = r.parseTablespaceOptions()
	return stmt
}

// parseAlterTablespaceStmt implements AlterTablespaceStmt:
//
//	"ALTER" ["UNDO"] "TABLESPACE" Identifier
//	    [("ADD" | "DROP") "DATAFILE" stringLit | "RENAME" "TO" Identifier
//	     | "SET" ("ACTIVE" | "INACTIVE")] TablespaceOptions
//
// ACTIVE and INACTIVE are not keywords and are matched in identifier
// position.
func (r *rdParser) parseAlterTablespaceStmt() ast.StmtNode {
	r.expect(alter)
	stmt := &ast.AlterTablespaceStmt{Undo: r.accept(undo)}
	r.expect(tablespace)
	stmt.Name = ast.NewCIStr(r.parseIdentifier())
	switch r.tok() {
	case add:
		r.advance()
		r.expect(datafile)
		stmt.Op = ast.AlterTablespaceAddDatafile
		stmt.Datafile = r.expect(stringLit).lit
	case drop:
		r.advance()
		r.expect(datafile)
		stmt.Op = ast.AlterTablespaceDropDatafile
		stmt.Datafile = r.expect(stringLit).lit
	case rename:
		r.advance()
		r.expect(to)
		stmt.Op = ast.AlterTablespaceRenameTo
		stmt.NewName = ast.NewCIStr(r.parseIdentifier())
	case set:
		r.advance()
		if isIdentifierTok(r.tok()) && strings.EqualFold(r.cur().lit, "ACTIVE") {
			stmt.Op = ast.AlterTablespaceSetActive
		} else if isIdentifierTok(r.tok()) && strings.EqualFold(r.cur().lit, "INACTIVE") {
			stmt.Op = ast.AlterTablespaceSetInactive
		} else {
			r.syntaxError()
		}
		r.advance()
	}
	stmt.Options = r.parseTablespaceOptions()
	return stmt
}

// parseCreateLogfileGroupStmt implements CreateLogfileGroupStmt:
//
//	"CREATE" "LOGFILE" "GROUP" Identifier "ADD" "UNDOFILE" stringLit
//	    TablespaceOptions
func (r *rdParser) parseCreateLogfileGroupStmt() ast.StmtNode {
	r.expect(create)
	r.expect(logfile)
	r.expect(group)
	stmt := &ast.CreateLogfileGroupStmt{GroupName: ast.NewCIStr(r.parseIdentifier())}
	r.expect(add)
	r.expect(undofile)
	stmt.Undofile = r.expect(stringLit).lit
	stmt.Options = r.parseTablespaceOptions()
	return stmt
}

// parseAlterLogfileGroupStmt implements AlterLogfileGroupStmt:
//
//	"ALTER" "LOGFILE" "GROUP" Identifier "ADD" "UNDOFILE" stringLit
//	    TablespaceOptions
func (r *rdParser) parseAlterLogfileGroupStmt() ast.StmtNode {
	r.expect(alter)
	r.expect(logfile)
	r.expect(group)
	stmt := &ast.AlterLogfileGroupStmt{GroupName: ast.NewCIStr(r.parseIdentifier())}
	r.expect(add)
	r.expect(undofile)
	stmt.Undofile = r.expect(stringLit).lit
	stmt.Options = r.parseTablespaceOptions()
	return stmt
}

/*
 * Spatial reference systems
 */

// parseCreateSpatialReferenceSystemStmt implements
// CreateSpatialReferenceSystemStmt:
//
//	"CREATE" OrReplace "SPATIAL" "REFERENCE" "SYSTEM" IfNotExists NUM
//	    SRSAttribute*
//	SRSAttribute: ("NAME" | "DEFINITION" | "DESCRIPTION") stringLit
//	    | "ORGANIZATION" stringLit "IDENTIFIED" "BY" NUM
//
// REFERENCE and the attribute names are not keywords and are matched in
// identifier position.
func (r *rdParser) parseCreateSpatialReferenceSystemStmt() ast.StmtNode {
	r.expect(create)
	stmt := &ast.CreateSpatialReferenceSystemStmt{}
	if r.tok() == or {
		r.advance()
		r.expect(replace)
		stmt.OrReplace = true
	}
	r.expect(spatial)
	r.expectIdentLit("REFERENCE")
	r.expect(system)
	stmt.IfNotExists = r.parseIfNotExists()
	stmt.Srid = getUint64FromNUM(r.expect(intLit).item)
	for isIdentifierTok(r.tok()) {
		attr := &ast.SRSAttribute{}
		switch strings.ToUpper(r.cur().lit) {
		case "NAME":
			attr.Tp = ast.SRSAttributeName
		case "DEFINITION":
			attr.Tp = ast.SRSAttributeDefinition
		case "ORGANIZATION":
			attr.Tp = ast.SRSAttributeOrganization
		case "DESCRIPTION":
			attr.Tp = ast.SRSAttributeDescription
		default:
			r.syntaxError()
		}
		r.advance()
		attr.StrValue = r.expect(stringLit).lit
		if attr.Tp == ast.SRSAttributeOrganization {
			r.expect(identified)
			r.expect(by)
			attr.OrgID = getUint64FromNUM(r.expect(intLit).item)
		}
		stmt.Attributes = append(stmt.Attributes, attr)
	}
	return stmt
}

/*
 * Libraries
 */

// parseCreateLibraryStmt implements CreateLibraryStmt:
//
//	"CREATE" OrReplace "LIBRARY" IfNotExists TableName
//	    "LANGUAGE" Identifier "AS" stringLit
func (r *rdParser) parseCreateLibraryStmt() ast.StmtNode {
	r.expect(create)
	stmt := &ast.CreateLibraryStmt{}
	if r.tok() == or {
		r.advance()
		r.expect(replace)
		stmt.OrReplace = true
	}
	r.expect(library)
	stmt.IfNotExists = r.parseIfNotExists()
	stmt.LibraryName = r.parseTableName()
	r.expect(language)
	stmt.Language = ast.NewCIStr(strings.ToUpper(r.parseIdentifier()))
	r.expect(as)
	stmt.Code = r.expect(stringLit).lit
	return stmt
}

// parseAlterLibraryStmt implements AlterLibraryStmt:
// "ALTER" "LIBRARY" TableName "COMMENT" stringLit.
func (r *rdParser) parseAlterLibraryStmt() ast.StmtNode {
	r.expect(alter)
	r.expect(library)
	stmt := &ast.AlterLibraryStmt{LibraryName: r.parseTableName()}
	r.expect(comment)
	stmt.Comment = r.expect(stringLit).lit
	return stmt
}

/*
 * JSON duality views
 */

// parseCreateJSONDualityViewStmt implements CreateJSONDualityViewStmt:
//
//	"CREATE" OrReplace "JSON" "DUALITY" "VIEW" IfNotExists TableName
//	    "AS" CreateViewSelectOpt
//
// recording the select statement's source text like CreateViewStmt.
func (r *rdParser) parseCreateJSONDualityViewStmt() ast.StmtNode {
	r.expect(create)
	stmt := &ast.CreateJSONDualityViewStmt{}
	if r.tok() == or {
		r.advance()
		r.expect(replace)
		stmt.OrReplace = true
	}
	r.expect(jsonType)
	r.expect(duality)
	r.expect(view)
	stmt.IfNotExists = r.parseIfNotExists()
	stmt.ViewName = r.parseTableName()
	r.expect(as)
	startOffset := r.cur().offset
	stmt.Select = r.parseCreateViewSelectOpt()
	r.p.setNodeText(stmt.Select, strings.TrimSpace(r.src[startOffset:r.cur().offset]))
	return stmt
}

// parseAlterJSONDualityViewStmt implements AlterJSONDualityViewStmt:
// "ALTER" "JSON" "DUALITY" "VIEW" TableName "AS" CreateViewSelectOpt.
func (r *rdParser) parseAlterJSONDualityViewStmt() ast.StmtNode {
	r.expect(alter)
	r.expect(jsonType)
	r.expect(duality)
	r.expect(view)
	stmt := &ast.AlterJSONDualityViewStmt{ViewName: r.parseTableName()}
	r.expect(as)
	startOffset := r.cur().offset
	stmt.Select = r.parseCreateViewSelectOpt()
	r.p.setNodeText(stmt.Select, strings.TrimSpace(r.src[startOffset:r.cur().offset]))
	return stmt
}
