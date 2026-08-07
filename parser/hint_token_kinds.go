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

import (
	"github.com/sqlc-dev/marino/ast"
)

// Token kinds and the lexer value type of the optimizer-hint parser,
// snapshotted verbatim from the goyacc-generated hintparser.go at its
// removal; hand-maintained from then on.

type yyhintSymType struct {
	yys            int
	offset         int
	ident          string
	number         uint64
	hint           *ast.TableOptimizerHint
	hints          []*ast.TableOptimizerHint
	table          ast.HintTable
	modelIdents    []ast.CIStr
	leadingList    *ast.LeadingList
	leadingElement interface{} // Modified: Represents either *ast.HintTable or *ast.LeadingList
}

const (
	yyhintDefault             = 57437
	yyhintEOFCode             = 57344
	yyhintErrCode             = 57345
	hintAggToCop              = 57380
	hintBCJoin                = 57403
	hintBKA                   = 57355
	hintBNL                   = 57357
	hintDupsWeedOut           = 57433
	hintFalse                 = 57429
	hintFirstMatch            = 57434
	hintForceIndex            = 57419
	hintGB                    = 57432
	hintHashAgg               = 57383
	hintHashJoin              = 57359
	hintHashJoinBuild         = 57360
	hintHashJoinProbe         = 57361
	hintHypoIndex             = 57379
	hintIdentifier            = 57347
	hintIgnoreIndex           = 57386
	hintIgnorePlanCache       = 57381
	hintIndexHashJoin         = 57390
	hintIndexJoin             = 57387
	hintIndexLookUpPushDown   = 57411
	hintIndexMerge            = 57365
	hintIndexMergeJoin        = 57394
	hintInlHashJoin           = 57389
	hintInlJoin               = 57392
	hintInlMergeJoin          = 57393
	hintIntLit                = 57346
	hintInvalid               = 57348
	hintJoinFixedOrder        = 57351
	hintJoinOrder             = 57352
	hintJoinPrefix            = 57353
	hintJoinSuffix            = 57354
	hintLeading               = 57421
	hintLimitToCop            = 57418
	hintLooseScan             = 57435
	hintMB                    = 57431
	hintMRR                   = 57367
	hintMaterialization       = 57436
	hintMaxExecutionTime      = 57375
	hintMemoryQuota           = 57396
	hintMerge                 = 57363
	hintMpp1PhaseAgg          = 57384
	hintMpp2PhaseAgg          = 57385
	hintNoBKA                 = 57356
	hintNoBNL                 = 57358
	hintNoDecorrelate         = 57423
	hintNoHashJoin            = 57362
	hintNoICP                 = 57369
	hintNoIndexHashJoin       = 57391
	hintNoIndexJoin           = 57388
	hintNoIndexLookUpPushDown = 57412
	hintNoIndexMerge          = 57366
	hintNoIndexMergeJoin      = 57395
	hintNoMRR                 = 57368
	hintNoMerge               = 57364
	hintNoOrderIndex          = 57410
	hintNoRangeOptimization   = 57370
	hintNoSMJoin              = 57402
	hintNoSemijoin            = 57374
	hintNoSkipScan            = 57372
	hintNoSwapJoinInputs      = 57397
	hintNthPlan               = 57417
	hintOLAP                  = 57424
	hintOLTP                  = 57425
	hintOrderIndex            = 57409
	hintPartition             = 57426
	hintQBName                = 57378
	hintQueryType             = 57398
	hintReadConsistentReplica = 57399
	hintReadFromStorage       = 57400
	hintResourceGroup         = 57377
	hintSMJoin                = 57401
	hintSemiJoinRewrite       = 57422
	hintSemijoin              = 57373
	hintSetVar                = 57376
	hintShuffleJoin           = 57404
	hintSingleAtIdentifier    = 57349
	hintSkipScan              = 57371
	hintStraightJoin          = 57420
	hintStreamAgg             = 57405
	hintStringLit             = 57350
	hintSwapJoinInputs        = 57406
	hintTiFlash               = 57428
	hintTiKV                  = 57427
	hintTimeRange             = 57415
	hintTrue                  = 57430
	hintUseCascades           = 57416
	hintUseIndex              = 57408
	hintUseIndexMerge         = 57407
	hintUsePlanCache          = 57413
	hintUseToja               = 57414
	hintWriteSlowLog          = 57382

	yyhintMaxDepth = 200
	yyhintTabOfs   = -229
)
