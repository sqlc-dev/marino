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

// The column type system: the Type production and its families
// (NumericType, StringType, DateAndTimeType) with their option
// productions (FieldOpts, FloatOpt, OptCharsetWithOptBinary, ...).

import (
	"strings"

	"github.com/sqlc-dev/marino/ast"
	"github.com/sqlc-dev/marino/charset"
	"github.com/sqlc-dev/marino/mysql"
	"github.com/sqlc-dev/marino/types"
)

// parseType implements Type: NumericType | StringType | DateAndTimeType.
func (r *rdParser) parseType() *types.FieldType {
	switch r.tok() {
	case tinyIntType, smallIntType, mediumIntType, middleIntType, intType,
		int1Type, int2Type, int3Type, int4Type, int8Type, integerType, bigIntType:
		// NumericType: IntegerType OptFieldLen FieldOpts
		b := integerTypeByte(r.tok())
		r.advance()
		flen := r.parseOptFieldLen()
		opts := r.parseFieldOpts()
		// TODO: check flen 0
		tp := types.NewFieldType(b)
		tp.SetFlen(flen)
		applyTypeOpts(tp, opts)
		return tp
	case boolType, booleanType:
		// NumericType: BooleanType FieldOpts
		r.advance()
		opts := r.parseFieldOpts()
		// TODO: check flen 0
		tp := types.NewFieldType(mysql.TypeTiny)
		tp.SetFlen(1)
		applyTypeOpts(tp, opts)
		return tp
	case decimalType, numericType, fixed:
		// NumericType: FixedPointType FloatOpt FieldOpts
		r.advance()
		fopt := r.parseFloatOpt()
		opts := r.parseFieldOpts()
		tp := types.NewFieldType(mysql.TypeNewDecimal)
		tp.SetFlen(fopt.Flen)
		tp.SetDecimal(fopt.Decimal)
		applyTypeOpts(tp, opts)
		return tp
	case floatType, realType, doubleType, float4Type, float8Type:
		// NumericType: FloatingPointType FloatOpt FieldOpts
		return r.parseFloatingPointType()
	case bitType:
		// NumericType: BitValueType OptFieldLen
		r.advance()
		flen := r.parseOptFieldLen()
		tp := types.NewFieldType(mysql.TypeBit)
		tp.SetFlen(flen)
		if tp.GetFlen() == types.UnspecifiedLength {
			tp.SetFlen(1)
		}
		return tp
	default:
		return r.parseStringOrDateAndTimeType()
	}
}

// integerTypeByte implements IntegerType.
func integerTypeByte(tok int) byte {
	switch tok {
	case tinyIntType, int1Type:
		return mysql.TypeTiny
	case smallIntType, int2Type:
		return mysql.TypeShort
	case mediumIntType, middleIntType, int3Type:
		return mysql.TypeInt24
	case intType, int4Type, integerType:
		return mysql.TypeLong
	case int8Type, bigIntType:
		return mysql.TypeLonglong
	}
	return 0
}

// parseFloatingPointType implements the
// `NumericType: FloatingPointType FloatOpt FieldOpts` alternative,
// including the FloatingPointType production itself.
func (r *rdParser) parseFloatingPointType() *types.FieldType {
	var b byte
	switch r.tok() {
	case floatType, float4Type:
		b = mysql.TypeFloat
		r.advance()
	case realType:
		// FloatingPointType: "REAL"
		if r.sc.GetSQLMode().HasRealAsFloatMode() {
			b = mysql.TypeFloat
		} else {
			b = mysql.TypeDouble
		}
		r.advance()
	case doubleType:
		// FloatingPointType: "DOUBLE" | "DOUBLE" "PRECISION"
		b = mysql.TypeDouble
		r.advance()
		r.accept(precisionType)
	case float8Type:
		b = mysql.TypeDouble
		r.advance()
	}
	fopt := r.parseFloatOpt()
	opts := r.parseFieldOpts()
	tp := types.NewFieldType(b)
	// check for a double(10) for syntax error
	if tp.GetType() == mysql.TypeDouble && r.p.strictDoubleFieldType {
		if fopt.Flen != types.UnspecifiedLength && fopt.Decimal == types.UnspecifiedLength {
			r.actionError(ErrSyntax)
		}
	}
	tp.SetFlen(fopt.Flen)
	if tp.GetType() == mysql.TypeFloat && fopt.Decimal == types.UnspecifiedLength && tp.GetFlen() <= mysql.MaxDoublePrecisionLength {
		if tp.GetFlen() > mysql.MaxFloatPrecisionLength {
			tp.SetType(mysql.TypeDouble)
		}
		tp.SetFlen(types.UnspecifiedLength)
	}
	tp.SetDecimal(fopt.Decimal)
	applyTypeOpts(tp, opts)
	return tp
}

// parseFieldOpts implements FieldOpts (a possibly-empty sequence of
// FieldOpt: "UNSIGNED" | "SIGNED" | "ZEROFILL").
func (r *rdParser) parseFieldOpts() []*ast.TypeOpt {
	opts := []*ast.TypeOpt{}
	for {
		switch r.tok() {
		case unsigned:
			r.advance()
			opts = append(opts, &ast.TypeOpt{IsUnsigned: true})
		case signed:
			r.advance()
			opts = append(opts, &ast.TypeOpt{IsUnsigned: false})
		case zerofill:
			r.advance()
			opts = append(opts, &ast.TypeOpt{IsZerofill: true, IsUnsigned: true})
		default:
			return opts
		}
	}
}

// applyTypeOpts reproduces the FieldOpts flag loop shared by the
// NumericType actions.
func applyTypeOpts(tp *types.FieldType, opts []*ast.TypeOpt) {
	for _, o := range opts {
		if o.IsUnsigned {
			tp.AddFlag(mysql.UnsignedFlag)
		}
		if o.IsZerofill {
			tp.AddFlag(mysql.ZerofillFlag)
		}
	}
}

// parseStringOrDateAndTimeType implements StringType and DateAndTimeType.
func (r *rdParser) parseStringOrDateAndTimeType() *types.FieldType {
	switch r.tok() {
	case charType, character:
		// Varchar: "CHARACTER" "VARYING" | "CHAR" "VARYING"
		if r.la(1) == varying {
			r.advance()
			r.advance()
			return r.parseVarcharTail()
		}
		// StringType: Char FieldLen OptBinary | Char OptBinary
		r.advance()
		return r.parseCharTail()
	case ncharType:
		// NVarchar: "NCHAR" "VARCHAR" | "NCHAR" "VARCHARACTER" | "NCHAR" "VARYING"
		if r.la(1) == varcharType || r.la(1) == varcharacter || r.la(1) == varying {
			r.advance()
			r.advance()
			return r.parseVarcharTail()
		}
		// NChar: "NCHAR"
		r.advance()
		return r.parseCharTail()
	case national:
		switch r.la(1) {
		case character, charType:
			if r.la(2) == varying {
				// NVarchar: "NATIONAL" "CHARACTER" "VARYING" | "NATIONAL" "CHAR" "VARYING"
				r.advance()
				r.advance()
				r.advance()
				return r.parseVarcharTail()
			}
			// NChar: "NATIONAL" "CHARACTER" | "NATIONAL" "CHAR"
			r.advance()
			r.advance()
			return r.parseCharTail()
		case varcharType, varcharacter:
			// NVarchar: "NATIONAL" "VARCHAR" | "NATIONAL" "VARCHARACTER"
			r.advance()
			r.advance()
			return r.parseVarcharTail()
		}
		r.syntaxError()
	case varcharType, varcharacter:
		// Varchar: "VARCHAR" | "VARCHARACTER"
		r.advance()
		return r.parseVarcharTail()
	case nvarcharType:
		// NVarchar: "NVARCHAR"
		r.advance()
		return r.parseVarcharTail()
	case binaryType:
		// StringType: "BINARY" OptFieldLen
		r.advance()
		flen := r.parseOptFieldLen()
		tp := types.NewFieldType(mysql.TypeString)
		tp.SetFlen(flen)
		tp.SetCharset(charset.CharsetBin)
		tp.SetCollate(charset.CharsetBin)
		tp.AddFlag(mysql.BinaryFlag)
		return tp
	case varbinaryType:
		// StringType: "VARBINARY" FieldLen
		r.advance()
		flen := r.parseFieldLen()
		tp := types.NewFieldType(mysql.TypeVarchar)
		tp.SetFlen(flen)
		tp.SetCharset(charset.CharsetBin)
		tp.SetCollate(charset.CharsetBin)
		tp.AddFlag(mysql.BinaryFlag)
		return tp
	case tinyblobType, blobType, mediumblobType, longblobType:
		// StringType: BlobType
		tp := r.parseBlobType()
		tp.SetCharset(charset.CharsetBin)
		tp.SetCollate(charset.CharsetBin)
		tp.AddFlag(mysql.BinaryFlag)
		return tp
	case tinytextType, textType, mediumtextType, longtextType:
		// StringType: TextType OptCharsetWithOptBinary
		tp := r.parseTextType()
		opt := r.parseOptCharsetWithOptBinary()
		tp.SetCharset(opt.Charset)
		if opt.Charset == charset.CharsetBin {
			tp.AddFlag(mysql.BinaryFlag)
			tp.SetCollate(charset.CollationBin)
		}
		if opt.IsBinary {
			tp.AddFlag(mysql.BinaryFlag)
		}
		return tp
	case enum:
		// StringType: "ENUM" '(' TextStringList ')' OptCharsetWithOptBinary
		r.advance()
		r.expect(int('('))
		elems := r.parseTextStringList()
		r.expect(int(')'))
		opt := r.parseOptCharsetWithOptBinary()
		tp := types.NewFieldType(mysql.TypeEnum)
		tp.SetElems(make([]string, len(elems)))
		fieldLen := -1 // enum_flen = max(ele_flen)
		for i, e := range elems {
			trimmed := strings.TrimRight(e.Value, " ")
			tp.SetElemWithIsBinaryLit(i, trimmed, e.IsBinaryLiteral)
			if len(trimmed) > fieldLen {
				fieldLen = len(trimmed)
			}
		}
		tp.SetFlen(fieldLen)
		tp.SetCharset(opt.Charset)
		if opt.IsBinary {
			tp.AddFlag(mysql.BinaryFlag)
		}
		return tp
	case set:
		// StringType: "SET" '(' TextStringList ')' OptCharsetWithOptBinary
		r.advance()
		r.expect(int('('))
		elems := r.parseTextStringList()
		r.expect(int(')'))
		opt := r.parseOptCharsetWithOptBinary()
		tp := types.NewFieldType(mysql.TypeSet)
		tp.SetElems(make([]string, len(elems)))
		fieldLen := len(elems) - 1 // set_flen = sum(ele_flen) + number_of_ele - 1
		for i, e := range elems {
			trimmed := strings.TrimRight(e.Value, " ")
			tp.SetElemWithIsBinaryLit(i, trimmed, e.IsBinaryLiteral)
			fieldLen += len(trimmed)
		}
		tp.SetFlen(fieldLen)
		tp.SetCharset(opt.Charset)
		if opt.IsBinary {
			tp.AddFlag(mysql.BinaryFlag)
		}
		return tp
	case jsonType:
		// StringType: "JSON"
		r.advance()
		tp := types.NewFieldType(mysql.TypeJSON)
		tp.SetDecimal(0)
		tp.SetCharset(charset.CharsetBin)
		tp.SetCollate(charset.CollationBin)
		return tp
	case geometryType, point, linestringType, polygonType, multipointType,
		multilinestringType, multipolygonType, geometryCollectionType:
		// SpatialType: "GEOMETRY" | "POINT" | "LINESTRING" | "POLYGON"
		// | "MULTIPOINT" | "MULTILINESTRING" | "MULTIPOLYGON"
		// | "GEOMETRYCOLLECTION" — MySQL 26.7 §13.4.1. The spatial types
		// postdate the goyacc grammar; GEOMCOLLECTION lexes as its
		// GEOMETRYCOLLECTION synonym. Spatial values are stored in a
		// binary format, like JSON.
		b := spatialTypeByte(r.tok())
		r.advance()
		tp := types.NewFieldType(b)
		tp.SetCharset(charset.CharsetBin)
		tp.SetCollate(charset.CollationBin)
		return tp
	case long:
		return r.parseLongType()
	case vectorType:
		// StringType: "VECTOR" OptVectorElementType OptFieldLen
		r.advance()
		elementType := r.parseOptVectorElementType()
		flen := r.parseOptFieldLen()
		if elementType.Tp != mysql.TypeFloat {
			r.sc.AppendError(r.actionErrorf("Only VECTOR is supported for now"))
		}
		tp := types.NewFieldType(mysql.TypeTiDBVectorFloat32)
		tp.SetFlen(flen)
		tp.SetDecimal(0)
		tp.SetCharset(charset.CharsetBin)
		tp.SetCollate(charset.CollationBin)
		return tp
	case dateType:
		// DateAndTimeType: "DATE"
		r.advance()
		return types.NewFieldType(mysql.TypeDate)
	case datetimeType:
		// DateAndTimeType: "DATETIME" OptFieldLen
		r.advance()
		dec := r.parseOptFieldLen()
		tp := types.NewFieldType(mysql.TypeDatetime)
		tp.SetFlen(mysql.MaxDatetimeWidthNoFsp)
		tp.SetDecimal(dec)
		if tp.GetDecimal() > 0 {
			tp.SetFlen(tp.GetFlen() + 1 + tp.GetDecimal())
		}
		return tp
	case timestampType:
		// DateAndTimeType: "TIMESTAMP" OptFieldLen
		r.advance()
		dec := r.parseOptFieldLen()
		tp := types.NewFieldType(mysql.TypeTimestamp)
		tp.SetFlen(mysql.MaxDatetimeWidthNoFsp)
		tp.SetDecimal(dec)
		if tp.GetDecimal() > 0 {
			tp.SetFlen(tp.GetFlen() + 1 + tp.GetDecimal())
		}
		return tp
	case timeType:
		// DateAndTimeType: "TIME" OptFieldLen
		r.advance()
		dec := r.parseOptFieldLen()
		tp := types.NewFieldType(mysql.TypeDuration)
		tp.SetFlen(mysql.MaxDurationWidthNoFsp)
		tp.SetDecimal(dec)
		if tp.GetDecimal() > 0 {
			tp.SetFlen(tp.GetFlen() + 1 + tp.GetDecimal())
		}
		return tp
	case yearType, sqlTsiYear:
		// DateAndTimeType: Year OptFieldLen FieldOpts
		r.advance()
		flen := r.parseOptFieldLen()
		r.parseFieldOpts()
		tp := types.NewFieldType(mysql.TypeYear)
		tp.SetFlen(flen)
		if tp.GetFlen() != types.UnspecifiedLength && tp.GetFlen() != 4 {
			// yylex.AppendError(...); return -1 — an aborting action.
			r.actionError(ErrInvalidYearColumnLength.GenWithStackByArgs())
		}
		return tp
	}
	r.syntaxError()
	return nil
}

// spatialTypeByte maps a spatial type keyword token to its type byte.
func spatialTypeByte(tok int) byte {
	switch tok {
	case geometryType:
		return mysql.TypeGeometry
	case point:
		return mysql.TypePoint
	case linestringType:
		return mysql.TypeLineString
	case polygonType:
		return mysql.TypePolygon
	case multipointType:
		return mysql.TypeMultiPoint
	case multilinestringType:
		return mysql.TypeMultiLineString
	case multipolygonType:
		return mysql.TypeMultiPolygon
	case geometryCollectionType:
		return mysql.TypeGeometryCollection
	}
	return 0
}

// parseCharTail finishes `Char/NChar FieldLen OptBinary` and
// `Char/NChar OptBinary` after the type keyword(s) have been consumed.
func (r *rdParser) parseCharTail() *types.FieldType {
	tp := types.NewFieldType(mysql.TypeString)
	if r.tok() == int('(') {
		tp.SetFlen(r.parseFieldLen())
	}
	opt := r.parseOptBinary()
	tp.SetCharset(opt.Charset)
	if opt.IsBinary {
		tp.AddFlag(mysql.BinaryFlag)
	}
	return tp
}

// parseVarcharTail finishes `Varchar/NVarchar FieldLen OptBinary` after
// the type keyword(s) have been consumed.
func (r *rdParser) parseVarcharTail() *types.FieldType {
	flen := r.parseFieldLen()
	opt := r.parseOptBinary()
	tp := types.NewFieldType(mysql.TypeVarchar)
	tp.SetFlen(flen)
	tp.SetCharset(opt.Charset)
	if opt.IsBinary {
		tp.AddFlag(mysql.BinaryFlag)
	}
	return tp
}

// parseBlobType implements BlobType (except the "LONG" "VARBINARY"
// alternative, which parseLongType owns).
func (r *rdParser) parseBlobType() *types.FieldType {
	switch r.tok() {
	case tinyblobType:
		r.advance()
		return types.NewFieldType(mysql.TypeTinyBlob)
	case blobType:
		// BlobType: "BLOB" OptFieldLen
		r.advance()
		tp := types.NewFieldType(mysql.TypeBlob)
		tp.SetFlen(r.parseOptFieldLen())
		return tp
	case mediumblobType:
		r.advance()
		return types.NewFieldType(mysql.TypeMediumBlob)
	case longblobType:
		r.advance()
		return types.NewFieldType(mysql.TypeLongBlob)
	}
	r.syntaxError()
	return nil
}

// parseTextType implements TextType.
func (r *rdParser) parseTextType() *types.FieldType {
	switch r.tok() {
	case tinytextType:
		r.advance()
		return types.NewFieldType(mysql.TypeTinyBlob)
	case textType:
		// TextType: "TEXT" OptFieldLen
		r.advance()
		tp := types.NewFieldType(mysql.TypeBlob)
		tp.SetFlen(r.parseOptFieldLen())
		return tp
	case mediumtextType:
		r.advance()
		return types.NewFieldType(mysql.TypeMediumBlob)
	case longtextType:
		r.advance()
		return types.NewFieldType(mysql.TypeLongBlob)
	}
	r.syntaxError()
	return nil
}

// parseLongType implements the "LONG"-led StringType alternatives:
//
//	StringType: "LONG" Varchar OptCharsetWithOptBinary
//	          | "LONG" OptCharsetWithOptBinary
//	          | BlobType ("LONG" "VARBINARY") — with the StringType: BlobType wrap
func (r *rdParser) parseLongType() *types.FieldType {
	r.expect(long)
	if r.tok() == varbinaryType {
		// BlobType: "LONG" "VARBINARY", then StringType: BlobType.
		r.advance()
		tp := types.NewFieldType(mysql.TypeMediumBlob)
		tp.SetCharset(charset.CharsetBin)
		tp.SetCollate(charset.CharsetBin)
		tp.AddFlag(mysql.BinaryFlag)
		return tp
	}
	if r.tok() == varcharType || r.tok() == varcharacter ||
		((r.tok() == character || r.tok() == charType) && r.la(1) == varying) {
		// Varchar: "CHARACTER" "VARYING" | "CHAR" "VARYING" | "VARCHAR" | "VARCHARACTER"
		if r.tok() == character || r.tok() == charType {
			r.advance()
			r.advance()
		} else {
			r.advance()
		}
	}
	// Both remaining alternatives share the same action.
	opt := r.parseOptCharsetWithOptBinary()
	tp := types.NewFieldType(mysql.TypeMediumBlob)
	tp.SetCharset(opt.Charset)
	if opt.Charset == charset.CharsetBin {
		tp.AddFlag(mysql.BinaryFlag)
		tp.SetCollate(charset.CollationBin)
	}
	if opt.IsBinary {
		tp.AddFlag(mysql.BinaryFlag)
	}
	return tp
}

// parseOptVectorElementType implements OptVectorElementType.
func (r *rdParser) parseOptVectorElementType() *ast.VectorElementType {
	if r.tok() != int('<') {
		return &ast.VectorElementType{Tp: mysql.TypeFloat}
	}
	r.advance()
	var tp byte
	switch r.tok() {
	case floatType:
		tp = mysql.TypeFloat
	case doubleType:
		tp = mysql.TypeDouble
	default:
		r.syntaxError()
	}
	r.advance()
	r.expect(int('>'))
	return &ast.VectorElementType{Tp: tp}
}

// parseOptCharsetWithOptBinary implements OptCharsetWithOptBinary.
func (r *rdParser) parseOptCharsetWithOptBinary() *ast.OptBinary {
	switch r.tok() {
	case ascii:
		r.advance()
		return &ast.OptBinary{
			IsBinary: false,
			Charset:  charset.CharsetLatin1,
		}
	case unicodeSym:
		r.advance()
		cs, err := charset.GetCharsetInfo("ucs2")
		if err != nil {
			r.actionError(ErrUnknownCharacterSet.GenWithStackByArgs("ucs2"))
		}
		return &ast.OptBinary{
			IsBinary: false,
			Charset:  cs.Name,
		}
	case byteType:
		r.advance()
		return &ast.OptBinary{
			IsBinary: false,
			Charset:  charset.CharsetBin,
		}
	default:
		return r.parseOptBinary()
	}
}

// parseTextStringList implements TextStringList.
func (r *rdParser) parseTextStringList() []*ast.TextString {
	list := []*ast.TextString{r.parseTextString()}
	for r.accept(int(',')) {
		list = append(list, r.parseTextString())
	}
	return list
}
