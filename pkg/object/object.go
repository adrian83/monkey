package object

import "fmt"

const (
	typeInteger = "INTEGER"
	typeBoolean = "BOOLEAN"
	typeNull    = "NULL"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) Type() ObjectType {
	return typeInteger
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return typeBoolean
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

type Null struct{}

func (n *Null) Type() ObjectType {
	return typeNull
}

func (n *Null) Inspect() string {
	return "null"
}
