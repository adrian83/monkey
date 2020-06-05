package object

import (
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/adrian83/monkey/pkg/ast"
)

const (
	TypeInteger  = "INTEGER"
	TypeString   = "STRING"
	TypeBoolean  = "BOOLEAN"
	TypeArray    = "ARRAY"
	TypeNull     = "NULL"
	TypeError    = "ERROR"
	TypeFunction = "FUNCTION"
	TypeBuiltin  = "BUILTIN"
	TypeHash     = "HASH"

	ReturnVal = "RETURN_VALUE"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

type TypedObject struct {
	objType ObjectType
}

func (o *TypedObject) Type() ObjectType {
	return o.objType
}

func NewInteger(val int64) *Integer {
	return &Integer{
		TypedObject: &TypedObject{objType: TypeInteger},
		Value:       val,
	}
}

type Integer struct {
	*TypedObject
	Value int64
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func NewString(val string) *String {
	return &String{
		TypedObject: &TypedObject{objType: TypeString},
		Value:       val,
	}
}

type String struct {
	*TypedObject
	Value string
}

func (s *String) Inspect() string {
	return s.Value
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

func NewBoolean(val bool) *Boolean {
	return &Boolean{
		TypedObject: &TypedObject{objType: TypeBoolean},
		Value:       val,
	}
}

type Boolean struct {
	*TypedObject
	Value bool
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *Boolean) HashKey() HashKey {
	var value uint64 = 0
	if b.Value {
		value = 1
	}

	return HashKey{Type: b.Type(), Value: value}
}

type Null struct{}

func (n *Null) Type() ObjectType {
	return TypeNull
}

func (n *Null) Inspect() string {
	return "null"
}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType {
	return ReturnVal
}

func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType {
	return TypeError
}

func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType {
	return TypeFunction
}

func (f *Function) Inspect() string {
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	return fmt.Sprintf("fn(%v) {\n%v\n}", strings.Join(params, ", "), f.Body.String())
}

type Array struct {
	Elements []Object
}

func (ao *Array) Type() ObjectType {
	return TypeArray
}

func (ao *Array) Inspect() string {
	elements := []string{}
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}

	return fmt.Sprintf("[%v]", strings.Join(elements, ", "))
}

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType {
	return TypeBuiltin
}

func (b *Builtin) Inspect() string {
	return "builtin function"
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType {
	return TypeHash
}

func (h *Hash) Inspect() string {
	pairs := []string{}
	for _, pair := range h.Pairs {
		p := fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect())
		pairs = append(pairs, p)
	}

	return fmt.Sprintf("{%v}", strings.Join(pairs, ", "))
}

type Hashable interface {
	HashKey() HashKey
}
