package lexer

import (
	"github.com/adrian83/monkey/pkg/token"
)

const (
	lowerCaseA = 'a'
	upperCaseA = 'A'
	lowerCaseZ = 'z'
	upperCaseZ = 'Z'
)

func New(input string) *Lexer {
	l := &Lexer{
		input: input,
	}
	return l
}

type Lexer struct {
	input    string
	position int
}

func (l *Lexer) Tokens() chan token.Token {
	result := make(chan token.Token, 10)

	go func() {
		for {
			t := l.nextToken()
			result <- t
			if t.Type == token.Eof {
				close(result)
				break
			}
		}
	}()

	return result
}

func (l *Lexer) nextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	ch := l.peekChar()

	switch ch {
	case '"':
		l.readChar()
		tok = newToken(token.TypeString, l.readString())
	case '+':
		tok = newToken(token.OperatorPlus, string(ch))
	case '-':
		tok = newToken(token.OperatorMinus, string(ch))
	case '/':
		tok = newToken(token.OperatorSlash, string(ch))
	case '*':
		tok = newToken(token.OperatorAsterisk, string(ch))
	case '<':
		tok = newToken(token.OperatorLowerThan, string(ch))
	case '>':
		tok = newToken(token.OperatorGreaterThan, string(ch))
	case ';':
		tok = newToken(token.DelimiterSemicolon, string(ch))
	case '(':
		tok = newToken(token.DelimiterLeftParenthesis, string(ch))
	case ')':
		tok = newToken(token.DelimiterRightParenthesis, string(ch))
	case ',':
		tok = newToken(token.DelimiterComma, string(ch))
	case '{':
		tok = newToken(token.DelimiterLeftBrace, string(ch))
	case '}':
		tok = newToken(token.DelimiterRightBrace, string(ch))
	case '[':
		tok = newToken(token.DelimiterLeftBracket, string(ch))
	case ']':
		tok = newToken(token.DelimiterRightBracket, string(ch))
	case ':':
		tok = newToken(token.DelimiterColon, string(ch))
	case 0:
		return newToken(token.Eof, "")
	case '=':
		chars := l.peekTwoChars()
		if chars == token.OperatorEqual {
			tok = token.Token{Type: token.OperatorEqual, Literal: chars}
			l.readChar()
		} else {
			tok = newToken(token.OperatorAssign, string(ch))
		}

	case '!':
		chars := l.peekTwoChars()
		if chars == token.OperatorNotEqual {
			tok = token.Token{Type: token.OperatorNotEqual, Literal: chars}
			l.readChar()
		} else {
			tok = newToken(token.OperatorBang, string(ch))
		}

	default:
		if isLetter(ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(ch) {
			tok.Type = token.TypeInteger
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.Illegal, string(ch))
		}
	}

	l.readChar()

	return tok
}

func newToken(tokenType token.TokenType, literal string) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: literal,
	}
}

func (l *Lexer) readString() string {

	position := l.position
	for {
		ch := l.peekChar()
		if ch != '"' && ch != 0 {
			l.readChar()
		} else {
			break
		}
	}

	result := l.input[position:l.position]
	return result
}

func (l *Lexer) readIdentifier() string {
	position := l.position

	for {
		ch := l.peekChar()
		if isLetter(ch) {
			l.readChar()
			ch = l.peekChar()
		} else {
			break
		}
	}

	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return lowerCaseA <= ch && ch <= lowerCaseZ || upperCaseA <= ch && ch <= upperCaseZ || ch == '_'
}

func (l *Lexer) skipWhitespace() {
	for {
		ch := l.peekChar()
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			l.readChar()
		} else {
			break
		}
	}
}

func (l *Lexer) readNumber() string {
	position := l.position

	for {
		ch := l.peekChar()
		if isDigit(ch) {
			ch = l.readChar()
		} else {
			break
		}

	}

	return l.input[position:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) peekChar() byte {
	if l.position >= len(l.input) {
		return 0
	}

	return l.input[l.position]
}

func (l *Lexer) peekTwoChars() string {
	if l.position+2 >= len(l.input) {
		return ""
	}

	return string(l.input[l.position : l.position+2])
}

func (l *Lexer) readChar() byte {
	var ch byte = 0
	if l.position < len(l.input) {
		ch = l.input[l.position]
	}
	l.position++
	return ch
}
