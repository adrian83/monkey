package token

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

	codeKeywordFunction = "fn"
	codeKeywordLet      = "let"
	codeKeywordTrue     = "true"
	codeKeywordFalse    = "false"
	codeKeywordIf       = "if"
	codeKeywordElse     = "else"
	codeKeywordReturn   = "return"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	codeKeywordFunction: KeywordFunction,
	codeKeywordLet:      KeywordLet,
	codeKeywordTrue:     KeywordTrue,
	codeKeywordFalse:    KeywordFalse,
	codeKeywordIf:       KeywordIf,
	codeKeywordElse:     KeywordElse,
	codeKeywordReturn:   KeywordReturn,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return Ident
}
