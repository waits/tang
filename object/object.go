package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/waits/tang/ast"
)

const (
	INTEGER      = "INTEGER"
	BOOLEAN      = "BOOLEAN"
	STRING       = "STRING"
	NULL         = "NULL"
	RETURN_VALUE = "RETURN_VALUE"
	ERROR        = "ERROR"
	FUNCTION     = "FUNCTION"
	BUILTIN      = "BUILTIN"
	LIST         = "LIST"
	TUPLE        = "TUPLE"
	MAP          = "MAP"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING }
func (s *String) Inspect() string  { return s.Value }

type Null struct{}

func (n *Null) Type() ObjectType { return NULL }
func (n *Null) Inspect() string  { return "null" }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Env
}

func (f *Function) Type() ObjectType { return FUNCTION }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

type BuiltinFunction func(args ...Object) Object
type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN }
func (b *Builtin) Inspect() string  { return "builtin function" }

type List struct {
	Elements []Object
}

func (ao *List) Type() ObjectType { return LIST }
func (ao *List) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type Tuple struct {
	Elements []Object
}

func (tp *Tuple) Type() ObjectType { return TUPLE }
func (tp *Tuple) Inspect() string {
	var out bytes.Buffer

	elements := make([]string, len(tp.Elements))
	for i, el := range tp.Elements {
		elements[i] = el.Inspect()
	}

	out.WriteString("(")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString(")")

	return out.String()
}

type MapKey struct {
	Type  ObjectType
	Value uint64
}

func (b *Boolean) MapKey() MapKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return MapKey{Type: b.Type(), Value: value}
}

func (i *Integer) MapKey() MapKey {
	return MapKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (s *String) MapKey() MapKey {
	h := fnv.New64a()

	h.Write([]byte(s.Value))

	return MapKey{Type: s.Type(), Value: h.Sum64()}
}

type MapPair struct {
	Key   Object
	Value Object
}

type Map struct {
	Pairs map[MapKey]MapPair
}

func (h *Map) Type() ObjectType { return MAP }
func (h *Map) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s",
			pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

type Mapable interface {
	MapKey() MapKey
}
