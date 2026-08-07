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

// The smaller CREATE families: CreateDatabaseStmt, CreateSequenceStmt,
// CreateUserStmt/CreateRoleStmt, CreateStatisticsStmt, CreatePolicyStmt,
// CreateResourceGroupStmt, and CreateMaskingPolicyStmt. The option-list
// productions these share with their ALTER counterparts
// (DatabaseOption, SequenceOption, ConnectionOptions,
// PasswordOrLockOptions, PlacementOptionList, ResourceGroupOptionList,
// the masking-policy clauses, ...) live in parse_alter.go.

import (
	"strings"
	"time"

	"github.com/sqlc-dev/marino/ast"
	"github.com/sqlc-dev/marino/auth"
)

// parseCreateDatabaseStmt implements:
//
//	CreateDatabaseStmt: "CREATE" DatabaseSym IfNotExists DBName
//	    DatabaseOptionListOpt
func (r *rdParser) parseCreateDatabaseStmt() ast.StmtNode {
	r.expect(create)
	r.expect(database)
	ifNotExists := r.parseIfNotExists()
	name := r.parseIdentifier()
	return &ast.CreateDatabaseStmt{
		IfNotExists: ifNotExists,
		Name:        ast.NewCIStr(name),
		Options:     r.parseDatabaseOptionListOpt(),
	}
}

// parseDatabaseOptionListOpt implements DatabaseOptionListOpt.
func (r *rdParser) parseDatabaseOptionListOpt() []*ast.DatabaseOption {
	if !r.isDatabaseOptionStart() {
		return []*ast.DatabaseOption{}
	}
	return r.parseDatabaseOptionList()
}

// isDatabaseOptionStart reports whether the current token can begin a
// DatabaseOption.
func (r *rdParser) isDatabaseOptionStart() bool {
	switch r.tok() {
	case defaultKwd, charsetKwd, collate, encryption, placement:
		return true
	case character, charType:
		// CharsetKw: "CHARACTER" "SET" | "CHAR" "SET"
		return r.la(1) == set
	case set:
		// "SET" "TIFLASH" "REPLICA" ...
		return r.la(1) == tiFlash
	}
	return false
}

// parseDatabaseOption implements DatabaseOption.
func (r *rdParser) parseDatabaseOption() *ast.DatabaseOption {
	switch r.tok() {
	case set:
		// "SET" "TIFLASH" "REPLICA" LengthNum LocationLabelList
		r.advance()
		r.expect(tiFlash)
		r.expect(replica)
		count := r.parseLengthNum()
		tiflashReplicaSpec := &ast.TiFlashReplicaSpec{
			Count:  count,
			Labels: r.parseLocationLabelList(),
		}
		return &ast.DatabaseOption{
			Tp:             ast.DatabaseSetTiFlashReplica,
			TiFlashReplica: tiflashReplicaSpec,
		}
	case placement:
		// PlacementPolicyOption (without DefaultKwdOpt)
		placementOptions := r.parsePlacementPolicyOption()
		return &ast.DatabaseOption{
			// offset trick, enums are identical but of different type
			Tp:        ast.DatabaseOptionType(placementOptions.Tp),
			Value:     placementOptions.StrValue,
			UintValue: placementOptions.UintValue,
		}
	}
	// DefaultKwdOpt
	r.accept(defaultKwd)
	switch {
	case r.tok() == collate:
		// DefaultKwdOpt "COLLATE" EqOpt CollationName
		r.advance()
		r.parseEqOpt()
		return &ast.DatabaseOption{Tp: ast.DatabaseOptionCollate, Value: r.parseCollationName()}
	case r.tok() == encryption:
		// DefaultKwdOpt "ENCRYPTION" EqOpt EncryptionOpt
		r.advance()
		r.parseEqOpt()
		return &ast.DatabaseOption{Tp: ast.DatabaseOptionEncryption, Value: r.parseEncryptionOpt()}
	case r.tok() == placement:
		// DefaultKwdOpt PlacementPolicyOption
		placementOptions := r.parsePlacementPolicyOption()
		return &ast.DatabaseOption{
			// offset trick, enums are identical but of different type
			Tp:        ast.DatabaseOptionType(placementOptions.Tp),
			Value:     placementOptions.StrValue,
			UintValue: placementOptions.UintValue,
		}
	case r.isCharsetKw():
		// DefaultKwdOpt CharsetKw EqOpt CharsetName
		r.parseCharsetKw()
		r.parseEqOpt()
		return &ast.DatabaseOption{Tp: ast.DatabaseOptionCharset, Value: r.parseCharsetName()}
	}
	r.syntaxError()
	return nil
}

// parseLocationLabelList implements LocationLabelList:
// empty | "LOCATION" "LABELS" StringList.
func (r *rdParser) parseLocationLabelList() []string {
	if r.tok() != location {
		return []string{}
	}
	r.advance()
	r.expect(labels)
	// StringList: stringLit | StringList ',' stringLit
	list := []string{r.expect(stringLit).lit}
	for r.accept(int(',')) {
		list = append(list, r.expect(stringLit).lit)
	}
	return list
}

// parseCreateSequenceStmt implements:
//
//	CreateSequenceStmt: "CREATE" "SEQUENCE" IfNotExists TableName
//	    CreateSequenceOptionListOpt CreateTableOptionListOpt
func (r *rdParser) parseCreateSequenceStmt() ast.StmtNode {
	r.expect(create)
	r.expect(sequence)
	ifNotExists := r.parseIfNotExists()
	name := r.parseTableName()
	return &ast.CreateSequenceStmt{
		IfNotExists: ifNotExists,
		Name:        name,
		SeqOptions:  r.parseCreateSequenceOptionListOpt(),
		TblOptions:  r.parseCreateTableOptionListOpt(),
	}
}

// parseCreateSequenceOptionListOpt implements
// CreateSequenceOptionListOpt and SequenceOptionList.
func (r *rdParser) parseCreateSequenceOptionListOpt() []*ast.SequenceOption {
	opts := []*ast.SequenceOption{}
	for r.isSequenceOptionStart() {
		opts = append(opts, r.parseSequenceOption())
	}
	return opts
}

// isSequenceOptionStart reports whether the current token can begin a
// SequenceOption.
func (r *rdParser) isSequenceOptionStart() bool {
	switch r.tok() {
	case increment, start, minValue, nominvalue, maxValue, nomaxvalue,
		cache, nocache, cycle, nocycle, no:
		return true
	}
	return false
}

// parseSequenceOption implements SequenceOption (shared with the ALTER
// SEQUENCE port in parse_alter.go).
func (r *rdParser) parseSequenceOption() *ast.SequenceOption {
	switch r.tok() {
	case increment:
		// "INCREMENT" EqOpt SignedNum | "INCREMENT" "BY" SignedNum
		r.advance()
		if !r.accept(by) {
			r.parseEqOpt()
		}
		return &ast.SequenceOption{Tp: ast.SequenceOptionIncrementBy, IntValue: r.parseSignedNum()}
	case start:
		// "START" EqOpt SignedNum | "START" "WITH" SignedNum
		r.advance()
		if !r.accept(with) {
			r.parseEqOpt()
		}
		return &ast.SequenceOption{Tp: ast.SequenceStartWith, IntValue: r.parseSignedNum()}
	case minValue:
		r.advance()
		r.parseEqOpt()
		return &ast.SequenceOption{Tp: ast.SequenceMinValue, IntValue: r.parseSignedNum()}
	case nominvalue:
		r.advance()
		return &ast.SequenceOption{Tp: ast.SequenceNoMinValue}
	case maxValue:
		r.advance()
		r.parseEqOpt()
		return &ast.SequenceOption{Tp: ast.SequenceMaxValue, IntValue: r.parseSignedNum()}
	case nomaxvalue:
		r.advance()
		return &ast.SequenceOption{Tp: ast.SequenceNoMaxValue}
	case cache:
		r.advance()
		r.parseEqOpt()
		return &ast.SequenceOption{Tp: ast.SequenceCache, IntValue: r.parseSignedNum()}
	case nocache:
		r.advance()
		return &ast.SequenceOption{Tp: ast.SequenceNoCache}
	case cycle:
		r.advance()
		return &ast.SequenceOption{Tp: ast.SequenceCycle}
	case nocycle:
		r.advance()
		return &ast.SequenceOption{Tp: ast.SequenceNoCycle}
	case no:
		// "NO" ("MINVALUE"|"MAXVALUE"|"CACHE"|"CYCLE")
		r.advance()
		switch r.tok() {
		case minValue:
			r.advance()
			return &ast.SequenceOption{Tp: ast.SequenceNoMinValue}
		case maxValue:
			r.advance()
			return &ast.SequenceOption{Tp: ast.SequenceNoMaxValue}
		case cache:
			r.advance()
			return &ast.SequenceOption{Tp: ast.SequenceNoCache}
		case cycle:
			r.advance()
			return &ast.SequenceOption{Tp: ast.SequenceNoCycle}
		}
	}
	r.syntaxError()
	return nil
}

// parseCreateUserStmt implements:
//
//	CreateUserStmt: "CREATE" "USER" IfNotExists UserSpecList
//	    RequireClauseOpt ConnectionOptions PasswordOrLockOptions
//	    CommentOrAttributeOption ResourceGroupNameOption
func (r *rdParser) parseCreateUserStmt() ast.StmtNode {
	r.expect(create)
	r.expect(user)
	ifNotExists := r.parseIfNotExists()
	specs := r.parseGrantUserSpecList()
	tlsOptions := r.parseGrantRequireClauseOpt()
	resourceOptions := r.parseUserConnectionOptions()
	pwdLockOptions := r.parseUserPasswordOrLockOptions()
	ret := &ast.CreateUserStmt{
		IsCreateRole:          false,
		IfNotExists:           ifNotExists,
		Specs:                 specs,
		AuthTokenOrTLSOptions: tlsOptions,
		ResourceOptions:       resourceOptions,
		PasswordOrLockOptions: pwdLockOptions,
	}
	if opt := r.parseUserCommentOrAttributeOption(); opt != nil {
		ret.CommentOrAttributeOption = opt
	}
	if opt := r.parseUserResourceGroupNameOption(); opt != nil {
		ret.ResourceGroupNameOption = opt
	}
	return ret
}

// parseUserConnectionOptions implements ConnectionOptions and
// ConnectionOptionList (empty, or "WITH" followed by one or more
// options; anything but MAX_USER_CONNECTIONS draws a warning).
func (r *rdParser) parseUserConnectionOptions() []*ast.ResourceOption {
	if r.tok() != with {
		return []*ast.ResourceOption{}
	}
	r.advance()
	l := []*ast.ResourceOption{r.parseUserConnectionOption()}
	for {
		switch r.tok() {
		case maxQueriesPerHour, maxUpdatesPerHour, maxConnectionsPerHour, maxUserConnections:
			l = append(l, r.parseUserConnectionOption())
		default:
			needWarning := false
			for _, option := range l {
				switch option.Type {
				case ast.MaxUserConnections:
				// do nothing.
				default:
					needWarning = true
				}
			}
			if needWarning {
				r.appendWarnf("TiDB does not support WITH ConnectionOptions but MAX_USER_CONNECTIONS now, they would be parsed but ignored.")
			}
			return l
		}
	}
}

// parseUserConnectionOption implements ConnectionOption.
func (r *rdParser) parseUserConnectionOption() *ast.ResourceOption {
	var typ int
	switch r.tok() {
	case maxQueriesPerHour:
		typ = ast.MaxQueriesPerHour
	case maxUpdatesPerHour:
		typ = ast.MaxUpdatesPerHour
	case maxConnectionsPerHour:
		typ = ast.MaxConnectionsPerHour
	case maxUserConnections:
		typ = ast.MaxUserConnections
	default:
		r.syntaxError()
	}
	r.advance()
	return &ast.ResourceOption{
		Type:  typ,
		Count: r.parseInt64Num(),
	}
}

// parseUserPasswordOrLockOptions implements PasswordOrLockOptions and
// PasswordOrLockOptionList.
func (r *rdParser) parseUserPasswordOrLockOptions() []*ast.PasswordOrLockOption {
	l := []*ast.PasswordOrLockOption{}
	for {
		switch r.tok() {
		case account, password, failedLoginAttempts, passwordLockTime:
			l = append(l, r.parseUserPasswordOrLockOption())
		default:
			return l
		}
	}
}

// parseUserPasswordOrLockOption implements PasswordOrLockOption.
func (r *rdParser) parseUserPasswordOrLockOption() *ast.PasswordOrLockOption {
	switch r.tok() {
	case account:
		r.advance()
		switch r.tok() {
		case unlock:
			r.advance()
			return &ast.PasswordOrLockOption{Type: ast.Unlock}
		case lock:
			r.advance()
			return &ast.PasswordOrLockOption{Type: ast.Lock}
		}
		r.syntaxError()
	case failedLoginAttempts:
		// "FAILED_LOGIN_ATTEMPTS" Int64Num
		r.advance()
		return &ast.PasswordOrLockOption{Type: ast.FailedLoginAttempts, Count: r.parseInt64Num()}
	case passwordLockTime:
		// "PASSWORD_LOCK_TIME" Int64Num | "PASSWORD_LOCK_TIME" "UNBOUNDED"
		r.advance()
		if r.accept(unbounded) {
			return &ast.PasswordOrLockOption{Type: ast.PasswordLockTimeUnbounded}
		}
		return &ast.PasswordOrLockOption{Type: ast.PasswordLockTime, Count: r.parseInt64Num()}
	case password:
		r.advance()
		switch r.tok() {
		case history:
			// "PASSWORD" "HISTORY" ("DEFAULT" | NUM)
			r.advance()
			if r.accept(defaultKwd) {
				return &ast.PasswordOrLockOption{Type: ast.PasswordHistoryDefault}
			}
			t := r.expect(intLit)
			return &ast.PasswordOrLockOption{Type: ast.PasswordHistory, Count: t.item.(int64)}
		case reuse:
			// "PASSWORD" "REUSE" "INTERVAL" ("DEFAULT" | NUM "DAY")
			r.advance()
			r.expect(interval)
			if r.accept(defaultKwd) {
				return &ast.PasswordOrLockOption{Type: ast.PasswordReuseDefault}
			}
			t := r.expect(intLit)
			opt := &ast.PasswordOrLockOption{Type: ast.PasswordReuseInterval, Count: t.item.(int64)}
			r.expect(day)
			return opt
		case expire:
			// "PASSWORD" "EXPIRE" ["INTERVAL" Int64Num "DAY" | "NEVER" | "DEFAULT"]
			r.advance()
			switch r.tok() {
			case interval:
				r.advance()
				n := r.parseInt64Num()
				r.expect(day)
				return &ast.PasswordOrLockOption{Type: ast.PasswordExpireInterval, Count: n}
			case never:
				r.advance()
				return &ast.PasswordOrLockOption{Type: ast.PasswordExpireNever}
			case defaultKwd:
				r.advance()
				return &ast.PasswordOrLockOption{Type: ast.PasswordExpireDefault}
			}
			return &ast.PasswordOrLockOption{Type: ast.PasswordExpire}
		case require:
			// "PASSWORD" "REQUIRE" "CURRENT" "DEFAULT"
			r.advance()
			r.expect(current)
			r.expect(defaultKwd)
			return &ast.PasswordOrLockOption{Type: ast.PasswordRequireCurrentDefault}
		}
	}
	r.syntaxError()
	return nil
}

// parseUserCommentOrAttributeOption implements
// CommentOrAttributeOption (nil for the empty alternative).
func (r *rdParser) parseUserCommentOrAttributeOption() *ast.CommentOrAttributeOption {
	switch r.tok() {
	case comment:
		r.advance()
		return &ast.CommentOrAttributeOption{Type: ast.UserCommentType, Value: r.expect(stringLit).lit}
	case attribute:
		r.advance()
		return &ast.CommentOrAttributeOption{Type: ast.UserAttributeType, Value: r.expect(stringLit).lit}
	}
	return nil
}

// parseUserResourceGroupNameOption implements ResourceGroupNameOption
// (nil for the empty alternative).
func (r *rdParser) parseUserResourceGroupNameOption() *ast.ResourceGroupNameOption {
	if r.tok() != resource {
		return nil
	}
	r.advance()
	r.expect(group)
	return &ast.ResourceGroupNameOption{Value: r.parseResourceGroupName()}
}

// parseCreateRoleStmt implements:
//
//	CreateRoleStmt: "CREATE" "ROLE" IfNotExists RoleSpecList
func (r *rdParser) parseCreateRoleStmt() ast.StmtNode {
	r.expect(create)
	r.expect(role)
	ifNotExists := r.parseIfNotExists()
	// RoleSpecList: RoleSpec | RoleSpecList ',' RoleSpec
	specs := []*ast.UserSpec{r.parseRoleSpec()}
	for r.accept(int(',')) {
		specs = append(specs, r.parseRoleSpec())
	}
	return &ast.CreateUserStmt{
		IsCreateRole: true,
		IfNotExists:  ifNotExists,
		Specs:        specs,
	}
}

// parseRoleSpec implements RoleSpec: Rolename.
func (r *rdParser) parseRoleSpec() *ast.UserSpec {
	role := r.parseRolename()
	return &ast.UserSpec{
		User: &auth.UserIdentity{
			Username: role.Username,
			Hostname: role.Hostname,
		},
		IsRole: true,
	}
}

// parseCreateStatisticsStmt implements:
//
//	CreateStatisticsStmt: "CREATE" "STATISTICS" IfNotExists Identifier
//	    '(' StatsType ')' "ON" TableName '(' ColumnNameList ')'
func (r *rdParser) parseCreateStatisticsStmt() ast.StmtNode {
	r.expect(create)
	r.expect(statistics)
	ifNotExists := r.parseIfNotExists()
	statsName := r.parseIdentifier()
	r.expect(int('('))
	// StatsType: "CARDINALITY" | "DEPENDENCY" | "CORRELATION"
	var statsType uint8
	switch r.tok() {
	case cardinality:
		statsType = ast.StatsTypeCardinality
	case dependency:
		statsType = ast.StatsTypeDependency
	case correlation:
		statsType = ast.StatsTypeCorrelation
	default:
		r.syntaxError()
	}
	r.advance()
	r.expect(int(')'))
	r.expect(on)
	table := r.parseTableName()
	r.expect(int('('))
	cols := r.parseColumnNameList()
	r.expect(int(')'))
	return &ast.CreateStatisticsStmt{
		IfNotExists: ifNotExists,
		StatsName:   statsName,
		StatsType:   statsType,
		Table:       table,
		Columns:     cols,
	}
}

// parseCreatePolicyStmt implements:
//
//	CreatePolicyStmt: "CREATE" OrReplace "PLACEMENT" "POLICY"
//	    IfNotExists PolicyName PlacementOptionList
func (r *rdParser) parseCreatePolicyStmt() ast.StmtNode {
	r.expect(create)
	orReplace := false
	if r.tok() == or {
		r.advance()
		r.expect(replace)
		orReplace = true
	}
	r.expect(placement)
	r.expect(policy)
	ifNotExists := r.parseIfNotExists()
	policyName := r.parseIdentifier()
	return &ast.CreatePlacementPolicyStmt{
		OrReplace:        orReplace,
		IfNotExists:      ifNotExists,
		PolicyName:       ast.NewCIStr(policyName),
		PlacementOptions: r.parsePlacementOptionList(),
	}
}

// isDirectPlacementOptionStart reports whether the current token can
// begin a DirectPlacementOption.
func (r *rdParser) isDirectPlacementOptionStart() bool {
	switch r.tok() {
	case primaryRegion, regions, followers, voters, learners, schedule,
		constraints, leaderConstraints, followerConstraints,
		voterConstraints, learnerConstraints, survivalPreferences:
		return true
	}
	return false
}

// parsePlacementOptionList implements PlacementOptionList (options
// separated by nothing or by a comma; a comma always continues the
// list).
func (r *rdParser) parsePlacementOptionList() []*ast.PlacementOption {
	opts := []*ast.PlacementOption{r.parseDirectPlacementOption()}
	for {
		if r.accept(int(',')) {
			opts = append(opts, r.parseDirectPlacementOption())
			continue
		}
		if !r.isDirectPlacementOptionStart() {
			return opts
		}
		opts = append(opts, r.parseDirectPlacementOption())
	}
}

// parseDirectPlacementOption implements DirectPlacementOption.
func (r *rdParser) parseDirectPlacementOption() *ast.PlacementOption {
	switch r.tok() {
	case primaryRegion:
		r.advance()
		r.parseEqOpt()
		return &ast.PlacementOption{Tp: ast.PlacementOptionPrimaryRegion, StrValue: r.expect(stringLit).lit}
	case regions:
		r.advance()
		r.parseEqOpt()
		return &ast.PlacementOption{Tp: ast.PlacementOptionRegions, StrValue: r.expect(stringLit).lit}
	case followers:
		r.advance()
		r.parseEqOpt()
		cnt := r.parseLengthNum()
		if cnt == 0 {
			r.actionError(r.sc.Errorf("FOLLOWERS must be positive"))
		}
		return &ast.PlacementOption{Tp: ast.PlacementOptionFollowerCount, UintValue: cnt}
	case voters:
		r.advance()
		r.parseEqOpt()
		return &ast.PlacementOption{Tp: ast.PlacementOptionVoterCount, UintValue: r.parseLengthNum()}
	case learners:
		r.advance()
		r.parseEqOpt()
		return &ast.PlacementOption{Tp: ast.PlacementOptionLearnerCount, UintValue: r.parseLengthNum()}
	case schedule:
		r.advance()
		r.parseEqOpt()
		return &ast.PlacementOption{Tp: ast.PlacementOptionSchedule, StrValue: r.expect(stringLit).lit}
	case constraints:
		r.advance()
		r.parseEqOpt()
		return &ast.PlacementOption{Tp: ast.PlacementOptionConstraints, StrValue: r.expect(stringLit).lit}
	case leaderConstraints:
		r.advance()
		r.parseEqOpt()
		return &ast.PlacementOption{Tp: ast.PlacementOptionLeaderConstraints, StrValue: r.expect(stringLit).lit}
	case followerConstraints:
		r.advance()
		r.parseEqOpt()
		return &ast.PlacementOption{Tp: ast.PlacementOptionFollowerConstraints, StrValue: r.expect(stringLit).lit}
	case voterConstraints:
		r.advance()
		r.parseEqOpt()
		return &ast.PlacementOption{Tp: ast.PlacementOptionVoterConstraints, StrValue: r.expect(stringLit).lit}
	case learnerConstraints:
		r.advance()
		r.parseEqOpt()
		return &ast.PlacementOption{Tp: ast.PlacementOptionLearnerConstraints, StrValue: r.expect(stringLit).lit}
	case survivalPreferences:
		r.advance()
		r.parseEqOpt()
		return &ast.PlacementOption{Tp: ast.PlacementOptionSurvivalPreferences, StrValue: r.expect(stringLit).lit}
	}
	r.syntaxError()
	return nil
}

// parseCreateResourceGroupStmt implements:
//
//	CreateResourceGroupStmt: "CREATE" "RESOURCE" "GROUP" IfNotExists
//	    ResourceGroupName ResourceGroupOptionList
func (r *rdParser) parseCreateResourceGroupStmt() ast.StmtNode {
	r.expect(create)
	r.expect(resource)
	r.expect(group)
	ifNotExists := r.parseIfNotExists()
	name := r.parseResourceGroupName()
	return &ast.CreateResourceGroupStmt{
		IfNotExists:             ifNotExists,
		ResourceGroupName:       ast.NewCIStr(name),
		ResourceGroupOptionList: r.parseResourceGroupOptionList(),
	}
}

// isDirectResourceGroupOptionStart reports whether the current token
// can begin a DirectResourceGroupOption.
func (r *rdParser) isDirectResourceGroupOptionStart() bool {
	switch r.tok() {
	case ruRate, priority, burstable, queryLimit, background:
		return true
	}
	return false
}

// parseResourceGroupOptionList implements ResourceGroupOptionList,
// including its duplicate-option check.
func (r *rdParser) parseResourceGroupOptionList() []*ast.ResourceGroupOption {
	opts := []*ast.ResourceGroupOption{r.parseDirectResourceGroupOption()}
	for {
		if !r.accept(int(',')) && !r.isDirectResourceGroupOptionStart() {
			return opts
		}
		next := r.parseDirectResourceGroupOption()
		if !ast.CheckAppend(opts, next) {
			r.actionError(r.sc.Errorf("Dupliated options specified"))
		}
		opts = append(opts, next)
	}
}

// parseDirectResourceGroupOption implements DirectResourceGroupOption.
func (r *rdParser) parseDirectResourceGroupOption() *ast.ResourceGroupOption {
	switch r.tok() {
	case ruRate:
		// "RU_PER_SEC" EqOpt (LengthNum | "UNLIMITED")
		r.advance()
		r.parseEqOpt()
		if r.accept(unlimited) {
			return &ast.ResourceGroupOption{Tp: ast.ResourceRURate, Burstable: ast.BurstableUnlimited}
		}
		return &ast.ResourceGroupOption{Tp: ast.ResourceRURate, UintValue: r.parseLengthNum()}
	case priority:
		// "PRIORITY" EqOpt ResourceGroupPriorityOption
		r.advance()
		r.parseEqOpt()
		var v uint64
		switch r.tok() {
		case low:
			v = uint64(1)
		case medium:
			v = uint64(8)
		case high:
			v = uint64(16)
		default:
			r.syntaxError()
		}
		r.advance()
		return &ast.ResourceGroupOption{Tp: ast.ResourcePriority, UintValue: v}
	case burstable:
		// "BURSTABLE" | "BURSTABLE" EqOpt ("MODERATED"|"UNLIMITED"|"OFF")
		r.advance()
		switch r.tok() {
		case eq, moderated, unlimited, off:
			r.parseEqOpt()
			switch r.tok() {
			case moderated:
				r.advance()
				return &ast.ResourceGroupOption{Tp: ast.ResourceBurstable, Burstable: ast.BurstableModerated}
			case unlimited:
				r.advance()
				return &ast.ResourceGroupOption{Tp: ast.ResourceBurstable, Burstable: ast.BurstableUnlimited}
			case off:
				r.advance()
				return &ast.ResourceGroupOption{Tp: ast.ResourceBurstable, Burstable: ast.BurstableDisable}
			}
			r.syntaxError()
		}
		return &ast.ResourceGroupOption{Tp: ast.ResourceBurstable, Burstable: ast.BurstableModerated}
	case queryLimit:
		// "QUERY_LIMIT" EqOpt ('(' [ResourceGroupRunawayOptionList] ')' | "NULL")
		r.advance()
		r.parseEqOpt()
		if r.accept(null) {
			return &ast.ResourceGroupOption{Tp: ast.ResourceGroupRunaway, RunawayOptionList: nil}
		}
		r.expect(int('('))
		if r.accept(int(')')) {
			return &ast.ResourceGroupOption{Tp: ast.ResourceGroupRunaway, RunawayOptionList: nil}
		}
		list := r.parseResourceGroupRunawayOptionList()
		r.expect(int(')'))
		return &ast.ResourceGroupOption{Tp: ast.ResourceGroupRunaway, RunawayOptionList: list}
	case background:
		// "BACKGROUND" EqOpt ('(' [ResourceGroupBackgroundOptionList] ')' | "NULL")
		r.advance()
		r.parseEqOpt()
		if r.accept(null) {
			return &ast.ResourceGroupOption{Tp: ast.ResourceGroupBackground, BackgroundOptions: nil}
		}
		r.expect(int('('))
		if r.accept(int(')')) {
			return &ast.ResourceGroupOption{Tp: ast.ResourceGroupBackground, BackgroundOptions: nil}
		}
		list := r.parseResourceGroupBackgroundOptionList()
		r.expect(int(')'))
		return &ast.ResourceGroupOption{Tp: ast.ResourceGroupBackground, BackgroundOptions: list}
	}
	r.syntaxError()
	return nil
}

// parseResourceGroupRunawayOptionList implements
// ResourceGroupRunawayOptionList, including its duplicate check.
func (r *rdParser) parseResourceGroupRunawayOptionList() []*ast.ResourceGroupRunawayOption {
	opts := []*ast.ResourceGroupRunawayOption{r.parseDirectResourceGroupRunawayOption()}
	for {
		if !r.accept(int(',')) {
			switch r.tok() {
			case execElapsed, processedKeys, ru, action, watch:
			default:
				return opts
			}
		}
		next := r.parseDirectResourceGroupRunawayOption()
		if !ast.CheckRunawayAppend(opts, next) {
			r.actionError(r.sc.Errorf("Dupliated runaway options specified"))
		}
		opts = append(opts, next)
	}
}

// parseDirectResourceGroupRunawayOption implements
// DirectResourceGroupRunawayOption.
func (r *rdParser) parseDirectResourceGroupRunawayOption() *ast.ResourceGroupRunawayOption {
	switch r.tok() {
	case execElapsed:
		// "EXEC_ELAPSED" EqOpt stringLit
		r.advance()
		r.parseEqOpt()
		s := r.expect(stringLit).lit
		_, err := time.ParseDuration(s)
		if err != nil {
			r.actionError(r.sc.Errorf("The EXEC_ELAPSED option is not a valid duration: %s", err.Error()))
		}
		return &ast.ResourceGroupRunawayOption{
			Tp:         ast.RunawayRule,
			RuleOption: &ast.ResourceGroupRunawayRuleOption{Tp: ast.RunawayRuleExecElapsed, ExecElapsed: s},
		}
	case processedKeys:
		// "PROCESSED_KEYS" EqOpt intLit
		r.advance()
		r.parseEqOpt()
		t := r.expect(intLit)
		return &ast.ResourceGroupRunawayOption{
			Tp:         ast.RunawayRule,
			RuleOption: &ast.ResourceGroupRunawayRuleOption{Tp: ast.RunawayRuleProcessedKeys, ProcessedKeys: t.item.(int64)},
		}
	case ru:
		// "RU" EqOpt intLit
		r.advance()
		r.parseEqOpt()
		t := r.expect(intLit)
		return &ast.ResourceGroupRunawayOption{
			Tp:         ast.RunawayRule,
			RuleOption: &ast.ResourceGroupRunawayRuleOption{Tp: ast.RunawayRuleRequestUnit, RequestUnit: t.item.(int64)},
		}
	case action:
		// "ACTION" EqOpt ResourceGroupRunawayActionOption
		r.advance()
		r.parseEqOpt()
		return &ast.ResourceGroupRunawayOption{
			Tp:           ast.RunawayAction,
			ActionOption: r.parseResourceGroupRunawayActionOption(),
		}
	case watch:
		// "WATCH" EqOpt ResourceGroupRunawayWatchOption WatchDurationOption
		r.advance()
		r.parseEqOpt()
		var watchType ast.RunawayWatchType
		switch r.tok() {
		case exact:
			watchType = ast.WatchExact
		case similar:
			watchType = ast.WatchSimilar
		case plan:
			watchType = ast.WatchPlan
		default:
			r.syntaxError()
		}
		r.advance()
		// WatchDurationOption: empty | "DURATION" EqOpt (stringLit | "UNLIMITED")
		durOpt := ""
		if r.accept(timeDuration) {
			r.parseEqOpt()
			if !r.accept(unlimited) {
				durOpt = r.expect(stringLit).lit
			}
		}
		dur := strings.ToLower(durOpt)
		if dur == "unlimited" {
			dur = ""
		}
		if len(dur) > 0 {
			_, err := time.ParseDuration(dur)
			if err != nil {
				r.actionError(r.sc.Errorf("The WATCH DURATION option is not a valid duration: %s", err.Error()))
			}
		}
		return &ast.ResourceGroupRunawayOption{
			Tp: ast.RunawayWatch,
			WatchOption: &ast.ResourceGroupRunawayWatchOption{
				Type:     watchType,
				Duration: dur,
			},
		}
	}
	r.syntaxError()
	return nil
}

// parseResourceGroupRunawayActionOption implements
// ResourceGroupRunawayActionOption.
func (r *rdParser) parseResourceGroupRunawayActionOption() *ast.ResourceGroupRunawayActionOption {
	switch r.tok() {
	case dryRun:
		r.advance()
		return &ast.ResourceGroupRunawayActionOption{Type: ast.RunawayActionDryRun}
	case cooldown:
		r.advance()
		return &ast.ResourceGroupRunawayActionOption{Type: ast.RunawayActionCooldown}
	case kill:
		r.advance()
		return &ast.ResourceGroupRunawayActionOption{Type: ast.RunawayActionKill}
	case switchGroup:
		// "SWITCH_GROUP" '(' ResourceGroupName ')'
		r.advance()
		r.expect(int('('))
		name := r.parseResourceGroupName()
		r.expect(int(')'))
		return &ast.ResourceGroupRunawayActionOption{
			Type:            ast.RunawayActionSwitchGroup,
			SwitchGroupName: ast.NewCIStr(name),
		}
	}
	r.syntaxError()
	return nil
}

// parseResourceGroupBackgroundOptionList implements
// ResourceGroupBackgroundOptionList, including its duplicate check.
func (r *rdParser) parseResourceGroupBackgroundOptionList() []*ast.ResourceGroupBackgroundOption {
	opts := []*ast.ResourceGroupBackgroundOption{r.parseDirectResourceGroupBackgroundOption()}
	for {
		if !r.accept(int(',')) {
			switch r.tok() {
			case taskTypes, utilizationLimit:
			default:
				return opts
			}
		}
		next := r.parseDirectResourceGroupBackgroundOption()
		if !ast.CheckBackgroundAppend(opts, next) {
			r.actionError(r.sc.Errorf("Dupliated background options specified"))
		}
		opts = append(opts, next)
	}
}

// parseDirectResourceGroupBackgroundOption implements
// DirectResourceGroupBackgroundOption.
func (r *rdParser) parseDirectResourceGroupBackgroundOption() *ast.ResourceGroupBackgroundOption {
	switch r.tok() {
	case taskTypes:
		// "TASK_TYPES" EqOpt stringLit
		r.advance()
		r.parseEqOpt()
		return &ast.ResourceGroupBackgroundOption{Type: ast.BackgroundOptionTaskNames, StrValue: r.expect(stringLit).lit}
	case utilizationLimit:
		// "UTILIZATION_LIMIT" EqOpt LengthNum
		r.advance()
		r.parseEqOpt()
		return &ast.ResourceGroupBackgroundOption{Type: ast.BackgroundUtilizationLimit, UintValue: r.parseLengthNum()}
	}
	r.syntaxError()
	return nil
}

// parseCreateMaskingPolicyStmt implements:
//
//	CreateMaskingPolicyStmt: "CREATE" OrReplace "MASKING" "POLICY"
//	    IfNotExists PolicyName "ON" TableName '(' Identifier ')' "AS"
//	    Expression MaskingPolicyRestrictOnOpt MaskingPolicyStateOpt
func (r *rdParser) parseCreateMaskingPolicyStmt() ast.StmtNode {
	r.expect(create)
	orReplace := false
	if r.tok() == or {
		r.advance()
		r.expect(replace)
		orReplace = true
	}
	r.expect(masking)
	r.expect(policy)
	ifNotExists := r.parseIfNotExists()
	policyName := r.parseIdentifier()
	r.expect(on)
	table := r.parseTableName()
	r.expect(int('('))
	column := r.parseIdentifier()
	r.expect(int(')'))
	r.expect(as)
	expr := r.parseExpression()
	restrictOps := r.parseMaskingPolicyRestrictOnOpt()
	state := r.parseMaskingPolicyStateOpt()
	if orReplace && ifNotExists {
		r.actionError(r.sc.Errorf("'OR REPLACE' and 'IF NOT EXISTS' are mutually exclusive"))
	}
	return &ast.CreateMaskingPolicyStmt{
		OrReplace:          orReplace,
		IfNotExists:        ifNotExists,
		PolicyName:         ast.NewCIStr(policyName),
		Table:              table,
		Column:             &ast.ColumnName{Name: ast.NewCIStr(column)},
		Expr:               expr,
		RestrictOps:        restrictOps,
		MaskingPolicyState: *state,
	}
}

// parseMaskingPolicyRestrictOnOpt implements
// MaskingPolicyRestrictOnOpt and MaskingPolicyRestrictOperationList.
func (r *rdParser) parseMaskingPolicyRestrictOnOpt() ast.MaskingPolicyRestrictOps {
	if r.tok() != restrict {
		return ast.MaskingPolicyRestrictOpNone
	}
	r.advance()
	r.expect(on)
	if r.accept(none) {
		// "RESTRICT" "ON" "NONE"
		return ast.MaskingPolicyRestrictOpNone
	}
	r.expect(int('('))
	ops := r.parseMaskingPolicyRestrictOperation()
	for r.accept(int(',')) {
		ops |= r.parseMaskingPolicyRestrictOperation()
	}
	r.expect(int(')'))
	return ops
}

// parseMaskingPolicyRestrictOperation implements
// MaskingPolicyRestrictOperation: Identifier.
func (r *rdParser) parseMaskingPolicyRestrictOperation() ast.MaskingPolicyRestrictOps {
	name := r.parseIdentifier()
	op, ok := getMaskingPolicyRestrictOp(name)
	if !ok {
		r.actionError(r.sc.Errorf("unsupported masking policy restrict operation: %s", name))
	}
	return op
}

// parseMaskingPolicyStateOpt implements MaskingPolicyStateOpt.
func (r *rdParser) parseMaskingPolicyStateOpt() *ast.MaskingPolicyState {
	switch r.tok() {
	case enable:
		r.advance()
		return &ast.MaskingPolicyState{Enabled: true, Explicit: true}
	case disable:
		r.advance()
		return &ast.MaskingPolicyState{Enabled: false, Explicit: true}
	}
	return &ast.MaskingPolicyState{Enabled: true, Explicit: false}
}
