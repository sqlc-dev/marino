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

// CREATE INDEX: CreateIndexStmt with IndexKeyTypeOpt, IndexTypeOpt, the
// IndexPartSpecification/IndexOption helpers shared with the CREATE
// TABLE constraint grammar (parse_column.go), and
// IndexLockAndAlgorithmOpt (parse_drop.go).

import (
	"github.com/sqlc-dev/marino/ast"
)

// parseCreateIndexStmt implements:
//
//	CreateIndexStmt: "CREATE" IndexKeyTypeOpt "INDEX" IfNotExists
//	    Identifier IndexTypeOpt "ON" TableName
//	    '(' IndexPartSpecificationList ')' IndexOptionList
//	    IndexLockAndAlgorithmOpt
func (r *rdParser) parseCreateIndexStmt() ast.StmtNode {
	r.expect(create)
	// IndexKeyTypeOpt
	keyType := ast.IndexKeyTypeNone
	switch r.tok() {
	case unique:
		keyType = ast.IndexKeyTypeUnique
		r.advance()
	case spatial:
		keyType = ast.IndexKeyTypeSpatial
		r.advance()
	case fulltext:
		keyType = ast.IndexKeyTypeFulltext
		r.advance()
	case vectorType:
		keyType = ast.IndexKeyTypeVector
		r.advance()
	case columnar:
		keyType = ast.IndexKeyTypeColumnar
		r.advance()
	}
	r.expect(index)
	ifNotExists := r.parseIfNotExists()
	indexName := r.parseIdentifier()
	// IndexTypeOpt: empty | IndexType ("USING"/"TYPE" IndexTypeName)
	var indexType interface{}
	if r.tok() == using || r.tok() == tp {
		r.advance()
		indexType = r.parseIndexTypeName()
	}
	r.expect(on)
	table := r.parseTableName()
	r.expect(int('('))
	specs := r.parseIndexPartSpecificationList()
	r.expect(int(')'))
	optionList := r.parseIndexOptionList()
	lockAlgOpt := r.parseDropIndexLockAndAlgorithm()

	var indexOption *ast.IndexOption
	if optionList != nil {
		indexOption = optionList
		if indexOption.Tp == ast.IndexTypeInvalid {
			if indexType != nil {
				indexOption.Tp = indexType.(ast.IndexType)
			}
		}
	} else {
		indexOption = &ast.IndexOption{}
		if indexType != nil {
			indexOption.Tp = indexType.(ast.IndexType)
		}
	}
	var indexLockAndAlgorithm *ast.IndexLockAndAlgorithm
	if lockAlgOpt != nil {
		indexLockAndAlgorithm = lockAlgOpt
		if indexLockAndAlgorithm.LockTp == ast.LockTypeDefault && indexLockAndAlgorithm.AlgorithmTp == ast.AlgorithmTypeDefault {
			indexLockAndAlgorithm = nil
		}
	}
	return &ast.CreateIndexStmt{
		IfNotExists:             ifNotExists,
		IndexName:               indexName,
		Table:                   table,
		IndexPartSpecifications: specs,
		IndexOption:             indexOption,
		KeyType:                 keyType,
		LockAlg:                 indexLockAndAlgorithm,
	}
}
