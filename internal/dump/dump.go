// Package dump renders AST nodes as deterministic, human-reviewable trees.
//
// The output pins the entire tree — every exported struct field plus the
// base-node state reachable through the ast.Node/ast.ExprNode interfaces
// (Text, OriginalText, OriginTextPosition, GetFlag, GetType) — so that two
// parsers producing "the same" AST can be compared field-for-field. It is
// the golden format of the parser corpus and the comparator of the
// differential harness; it is not a formatter and has no compatibility
// promise beyond being deterministic for equal trees.
package dump

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/sqlc-dev/marino/ast"
	"github.com/sqlc-dev/marino/types"
)

// Statements renders a parse result.
func Statements(stmts []ast.StmtNode) string {
	var d dumper
	d.buf.Grow(1024)
	fmt.Fprintf(&d.buf, "[]ast.StmtNode len=%d\n", len(stmts))
	for i, s := range stmts {
		d.indent = 0
		fmt.Fprintf(&d.buf, "[%d] ", i)
		d.value(reflect.ValueOf(s))
		d.buf.WriteByte('\n')
	}
	return d.buf.String()
}

// Any renders an arbitrary value with the same rules as Statements.
func Any(v any) string {
	var d dumper
	d.value(reflect.ValueOf(v))
	d.buf.WriteByte('\n')
	return d.buf.String()
}

type dumper struct {
	buf    strings.Builder
	indent int
	onPath map[uintptr]bool // cycle guard: pointers on the current path
}

func (d *dumper) nl() {
	d.buf.WriteByte('\n')
	for range d.indent {
		d.buf.WriteString("  ")
	}
}

func (d *dumper) value(v reflect.Value) {
	if !v.IsValid() {
		d.buf.WriteString("nil")
		return
	}
	switch v.Kind() {
	case reflect.Interface:
		if v.IsNil() {
			d.buf.WriteString("nil")
			return
		}
		d.value(v.Elem())
	case reflect.Ptr:
		if v.IsNil() {
			d.buf.WriteString("nil")
			return
		}
		p := v.Pointer()
		if d.onPath[p] {
			d.buf.WriteString("<cycle>")
			return
		}
		if d.onPath == nil {
			d.onPath = make(map[uintptr]bool)
		}
		d.onPath[p] = true
		d.buf.WriteByte('&')
		d.value(v.Elem())
		delete(d.onPath, p)
	case reflect.Struct:
		d.structValue(v)
	case reflect.Slice, reflect.Array:
		if v.Kind() == reflect.Slice && v.IsNil() {
			d.buf.WriteString("nil")
			return
		}
		if v.Type().Elem().Kind() == reflect.Uint8 {
			// []byte: hex, so binary payloads stay printable.
			fmt.Fprintf(&d.buf, "0x%x", v.Bytes())
			return
		}
		fmt.Fprintf(&d.buf, "%s len=%d [", v.Type(), v.Len())
		d.indent++
		for i := range v.Len() {
			d.nl()
			fmt.Fprintf(&d.buf, "%d: ", i)
			d.value(v.Index(i))
		}
		d.indent--
		if v.Len() > 0 {
			d.nl()
		}
		d.buf.WriteByte(']')
	case reflect.Map:
		if v.IsNil() {
			d.buf.WriteString("nil")
			return
		}
		keys := v.MapKeys()
		type kv struct {
			ks string
			k  reflect.Value
		}
		pairs := make([]kv, 0, len(keys))
		for _, k := range keys {
			pairs = append(pairs, kv{fmt.Sprintf("%v", k), k})
		}
		sort.Slice(pairs, func(i, j int) bool { return pairs[i].ks < pairs[j].ks })
		fmt.Fprintf(&d.buf, "%s len=%d {", v.Type(), v.Len())
		d.indent++
		for _, p := range pairs {
			d.nl()
			fmt.Fprintf(&d.buf, "%s: ", p.ks)
			d.value(v.MapIndex(p.k))
		}
		d.indent--
		if len(pairs) > 0 {
			d.nl()
		}
		d.buf.WriteByte('}')
	case reflect.String:
		d.buf.WriteString(strconv.Quote(v.String()))
	case reflect.Float32, reflect.Float64:
		d.buf.WriteString(strconv.FormatFloat(v.Float(), 'g', -1, 64))
	case reflect.Func, reflect.Chan, reflect.UnsafePointer:
		fmt.Fprintf(&d.buf, "<%s>", v.Type())
	default:
		fmt.Fprintf(&d.buf, "%v", v)
	}
}

func (d *dumper) structValue(v reflect.Value) {
	t := v.Type()
	if t == reflect.TypeOf(types.FieldType{}) {
		d.fieldType(v)
		return
	}
	fmt.Fprintf(&d.buf, "%s {", t)
	d.indent++
	d.nodeState(v)
	for i := range t.NumField() {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		d.nl()
		d.buf.WriteString(f.Name)
		d.buf.WriteString(": ")
		d.value(v.Field(i))
	}
	d.indent--
	d.nl()
	d.buf.WriteByte('}')
}

// nodeState prints the base-node state that lives in unexported embedded
// structs and is only reachable through interface methods. Values arrive
// here dereferenced, so re-take the address to find the methods.
func (d *dumper) nodeState(v reflect.Value) {
	if !v.CanAddr() {
		return
	}
	a := v.Addr()
	if !a.CanInterface() {
		return
	}
	iv := a.Interface()
	if n, ok := iv.(ast.Node); ok {
		// OriginalText, not Text(): Text() memoizes a converted copy on
		// the node, and a dump must never mutate the tree it renders
		// (the test suite deep-compares trees including that memo state).
		//
		// OriginTextPosition is deliberately NOT rendered: the goyacc
		// runtime restamps it from stale expression values left in
		// reused parser-stack slots (empty reductions point yyVAL one
		// past the top of stack), so its value depends on LALR stack
		// layout, not on the statement. The RD parser stamps the
		// deterministic production-start offset instead — the value
		// every explicit test assertion and downstream consumer
		// expects — and those assertions remain the gate for it.
		d.nl()
		fmt.Fprintf(&d.buf, ".originalText: %s", strconv.Quote(n.OriginalText()))
	}
	if e, ok := iv.(ast.ExprNode); ok {
		d.nl()
		fmt.Fprintf(&d.buf, ".flag: %d", e.GetFlag())
		d.nl()
		d.buf.WriteString(".type: ")
		d.fieldType(reflect.ValueOf(e.GetType()).Elem())
	}
}

// fieldType renders types.FieldType through its getters; its fields are
// unexported.
func (d *dumper) fieldType(v reflect.Value) {
	if !v.CanAddr() {
		// Copy so the getter methods (pointer receivers) are reachable.
		c := reflect.New(v.Type())
		c.Elem().Set(v)
		v = c.Elem()
	}
	ft := v.Addr().Interface().(*types.FieldType)
	fmt.Fprintf(&d.buf, "FieldType{tp:%d flag:%d flen:%d dec:%d charset:%q collate:%q elems:%q array:%v}",
		ft.GetType(), ft.GetFlag(), ft.GetFlen(), ft.GetDecimal(), ft.GetCharset(), ft.GetCollate(), ft.GetElems(), ft.IsArray())
}
