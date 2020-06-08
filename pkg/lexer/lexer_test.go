package lexer

import (
	"testing"

	"github.com/adrian83/monkey/pkg/token"
)

func TestNextToken(t *testing.T) {
	input := `let five = 5;
	let ten = 10;
	let add = fn(x, y) { x + y; };
	let result = add(five, ten);
	!-/*5;
	5 < 10 > 5;
	
	if (5 < 10) {
		return true;
	} else {
		return false;
	}
	
	10 == 10;
	10 != 9;
	"foobar"
	"foo bar"
	[1, 2];
	{"foo": "bar"}
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.KeywordLet, "let"},
		{token.Ident, "five"},
		{token.OperatorAssign, "="},
		{token.TypeInteger, "5"},
		{token.DelimiterSemicolon, ";"},
		{token.KeywordLet, "let"},
		{token.Ident, "ten"},
		{token.OperatorAssign, "="},
		{token.TypeInteger, "10"},
		{token.DelimiterSemicolon, ";"},
		{token.KeywordLet, "let"},
		{token.Ident, "add"},
		{token.OperatorAssign, "="},
		{token.KeywordFunction, "fn"},
		{token.DelimiterLeftParenthesis, "("},
		{token.Ident, "x"},
		{token.DelimiterComma, ","},
		{token.Ident, "y"},
		{token.DelimiterRightParenthesis, ")"},
		{token.DelimiterLeftBrace, "{"},
		{token.Ident, "x"},
		{token.OperatorPlus, "+"},
		{token.Ident, "y"},
		{token.DelimiterSemicolon, ";"},
		{token.DelimiterRightBrace, "}"},
		{token.DelimiterSemicolon, ";"},
		{token.KeywordLet, "let"},
		{token.Ident, "result"},
		{token.OperatorAssign, "="},
		{token.Ident, "add"},
		{token.DelimiterLeftParenthesis, "("},
		{token.Ident, "five"},
		{token.DelimiterComma, ","},
		{token.Ident, "ten"},
		{token.DelimiterRightParenthesis, ")"},
		{token.DelimiterSemicolon, ";"},
		{token.OperatorBang, "!"},
		{token.OperatorMinus, "-"},
		{token.OperatorSlash, "/"},
		{token.OperatorAsterisk, "*"},
		{token.TypeInteger, "5"},
		{token.DelimiterSemicolon, ";"},
		{token.TypeInteger, "5"},
		{token.OperatorLowerThan, "<"},
		{token.TypeInteger, "10"},
		{token.OperatorGreaterThan, ">"},
		{token.TypeInteger, "5"},
		{token.DelimiterSemicolon, ";"},
		{token.KeywordIf, "if"},
		{token.DelimiterLeftParenthesis, "("},
		{token.TypeInteger, "5"},
		{token.OperatorLowerThan, "<"},
		{token.TypeInteger, "10"},
		{token.DelimiterRightParenthesis, ")"},
		{token.DelimiterLeftBrace, "{"},
		{token.KeywordReturn, "return"},
		{token.KeywordTrue, "true"},
		{token.DelimiterSemicolon, ";"},
		{token.DelimiterRightBrace, "}"},
		{token.KeywordElse, "else"},
		{token.DelimiterLeftBrace, "{"},
		{token.KeywordReturn, "return"},
		{token.KeywordFalse, "false"},
		{token.DelimiterSemicolon, ";"},
		{token.DelimiterRightBrace, "}"},
		{token.TypeInteger, "10"},
		{token.OperatorEqual, "=="},
		{token.TypeInteger, "10"},
		{token.DelimiterSemicolon, ";"},
		{token.TypeInteger, "10"},
		{token.OperatorNotEqual, "!="},
		{token.TypeInteger, "9"},
		{token.DelimiterSemicolon, ";"},
		{token.TypeString, "foobar"},
		{token.TypeString, "foo bar"},
		{token.DelimiterLeftBracket, "["},
		{token.TypeInteger, "1"},
		{token.DelimiterComma, ","},
		{token.TypeInteger, "2"},
		{token.DelimiterRightBracket, "]"},
		{token.DelimiterSemicolon, ";"},
		{token.DelimiterLeftBrace, "{"},
		{token.TypeString, "foo"},
		{token.DelimiterColon, ":"},
		{token.TypeString, "bar"},
		{token.DelimiterRightBrace, "}"},
		{token.Eof, ""}}

	l := New(input)

	i := 0
	for tok := range l.Tokens() {

		tt := tests[i]

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q, token: %v", i, tt.expectedType, tok.Type, tok)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q, token: %v", i, tt.expectedLiteral, tok.Literal, tok)
		}

		i++
	}
}
