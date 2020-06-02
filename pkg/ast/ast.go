package ast

import (
	"bytes"
	"strings"

	"github.com/adrian83/monkey/pkg/token"
)

type BodyHolder interface {
	BodyStatements() []Statement
}

type Node interface {
	NodeToken() token.Token
	String() string
}

type Statement interface {
	Node
}

type Expression interface {
	Node
}

type Program struct {
	Statements []Statement
}

func (p *Program) BodyStatements() []Statement {
	return p.Statements
}

func (p *Program) NodeToken() token.Token {
	if len(p.Statements) > 0 {
		return p.Statements[0].NodeToken()
	}

	return token.Token{Type: token.EOF, Literal: ""}
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// expression
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) NodeToken() token.Token {
	return i.Token
}

func (i *Identifier) String() string {
	return i.Value
}

// expression
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) NodeToken() token.Token {
	return il.Token
}

func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

// expression
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) NodeToken() token.Token {
	return sl.Token
}

func (sl *StringLiteral) String() string {
	return sl.Token.Literal
}

// expression
type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) NodeToken() token.Token {
	return pe.Token
}

func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

// expression
type InfixExpression struct {
	Token    token.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) NodeToken() token.Token {
	return ie.Token
}

func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}

// expression
type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (b *BooleanLiteral) NodeToken() token.Token {
	return b.Token
}

func (b *BooleanLiteral) String() string {
	return b.Token.Literal
}

// expression
type IfExpression struct {
	Token       token.Token // The 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) NodeToken() token.Token {
	return ie.Token
}

func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

// expression
type FunctionLiteral struct {
	Token      token.Token // The 'fn' token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) NodeToken() token.Token {
	return fl.Token
}

func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.Token.Literal)
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())

	return out.String()
}

// expression
type ArrayLiteral struct {
	Token    token.Token // the '[' token
	Elements []Expression
}

func (al *ArrayLiteral) NodeToken() token.Token {
	return al.Token
}

func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type IndexExpression struct {
	Token token.Token // The [ token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode() {}

func (ie *IndexExpression) NodeToken() token.Token {
	return ie.Token
}

func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}

// expression
type CallExpression struct {
	Token     token.Token // The '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) NodeToken() token.Token {
	return ce.Token
}

func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")

	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

// ------ STATEMENTS ------

// statement
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) NodeToken() token.Token {
	return es.Token
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// statement
type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) NodeToken() token.Token {
	return ls.Token
}

func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.NodeToken().Literal + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")
	return out.String()
}

// statement
type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) BodyStatements() []Statement {
	return bs.Statements
}

func (bs *BlockStatement) NodeToken() token.Token {
	return bs.Token
}

func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// statement
type ReturnStatement struct {
	Token       token.Token // the 'return' token
	ReturnValue Expression
}

func (rs *ReturnStatement) NodeToken() token.Token {
	return rs.Token
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.NodeToken().Literal + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}
