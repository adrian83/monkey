package token

import "fmt"

type Operator string

const (
	Illegal = "ILLEGAL"
	Eof     = "EOF"

	// Identifiers + literals
	Ident       = "IDENT" // add, foobar, x, y, ...
	TypeInteger = "INT"   // 1343456
	TypeString  = "STRING"

	// Operators
	OperatorAssign      = "="
	OperatorPlus        = "+"
	OperatorMinus       = "-"
	OperatorBang        = "!"
	OperatorAsterisk    = "*"
	OperatorSlash       = "/"
	OperatorEqual       = "=="
	OperatorNotEqual    = "!="
	OperatorLowerThan   = "<"
	OperatorGreaterThan = ">"

	// Delimiters
	DelimiterComma            = ","
	DelimiterSemicolon        = ";"
	DelimiterLeftParenthesis  = "("
	DelimiterRightParenthesis = ")"
	DelimiterLeftBrace        = "{"
	DelimiterRightBrace       = "}"
	DelimiterLeftBracket      = "["
	DelimiterRightBracket     = "]"
	DelimiterColon            = ":"

	// Keywords
	KeywordFunction = "FUNCTION"
	KeywordLet      = "LET"
	KeywordTrue     = "TRUE"
	KeywordFalse    = "FALSE"
	KeywordIf       = "IF"
	KeywordElse     = "ELSE"
	KeywordReturn   = "RETURN"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

func (t *Token) String() string {
	return fmt.Sprintf("Token: {Type: %v, Literal: %v}", t.Type, t.Literal)
}

var keywords = map[string]TokenType{
	"fn":     KeywordFunction,
	"let":    KeywordLet,
	"true":   KeywordTrue,
	"false":  KeywordFalse,
	"if":     KeywordIf,
	"else":   KeywordElse,
	"return": KeywordReturn,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return Ident
}
