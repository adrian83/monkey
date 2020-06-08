package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/adrian83/monkey/pkg/ast"
	"github.com/adrian83/monkey/pkg/lexer"
	"github.com/adrian83/monkey/pkg/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX
)

var precedences = map[token.TokenType]int{
	token.OperatorEqual:            EQUALS,
	token.OperatorNotEqual:         EQUALS,
	token.OperatorLowerThan:        LESSGREATER,
	token.OperatorGreaterThan:      LESSGREATER,
	token.OperatorPlus:             SUM,
	token.OperatorMinus:            SUM,
	token.OperatorSlash:            PRODUCT,
	token.OperatorAsterisk:         PRODUCT,
	token.DelimiterLeftParenthesis: CALL,
	token.DelimiterLeftBracket:     INDEX,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	tokens chan token.Token
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		tokens: l.Tokens(),
		errors: []string{},
	}

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.Ident, p.parseIdentifier)
	p.registerPrefix(token.TypeInteger, p.parseIntegerLiteral)
	p.registerPrefix(token.OperatorBang, p.parsePrefixExpression)
	p.registerPrefix(token.OperatorMinus, p.parsePrefixExpression)
	p.registerPrefix(token.KeywordTrue, p.parseBooleanLiteral)
	p.registerPrefix(token.KeywordFalse, p.parseBooleanLiteral)
	p.registerPrefix(token.DelimiterLeftParenthesis, p.parseGroupedExpression)
	p.registerPrefix(token.KeywordIf, p.parseIfExpression)
	p.registerPrefix(token.KeywordFunction, p.parseFunctionLiteral)
	p.registerPrefix(token.TypeString, p.parseStringLiteral)
	p.registerPrefix(token.DelimiterLeftBracket, p.parseArrayLiteral)
	p.registerPrefix(token.DelimiterLeftBrace, p.parseHashLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.OperatorPlus, p.parseInfixExpression)
	p.registerInfix(token.OperatorMinus, p.parseInfixExpression)
	p.registerInfix(token.OperatorSlash, p.parseInfixExpression)
	p.registerInfix(token.OperatorAsterisk, p.parseInfixExpression)
	p.registerInfix(token.OperatorEqual, p.parseInfixExpression)
	p.registerInfix(token.OperatorNotEqual, p.parseInfixExpression)
	p.registerInfix(token.OperatorLowerThan, p.parseInfixExpression)
	p.registerInfix(token.OperatorGreaterThan, p.parseInfixExpression)
	p.registerInfix(token.DelimiterLeftParenthesis, p.parseCallExpression)
	p.registerInfix(token.DelimiterLeftBracket, p.parseIndexExpression)

	return p
}

func (p *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{Token: p.curToken}
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(token.DelimiterRightBrace) {
		p.nextToken()
		key := p.parseExpression(LOWEST)

		if !p.expectPeek(token.DelimiterColon) {
			return nil
		}
		p.nextToken()
		value := p.parseExpression(LOWEST)

		hash.Pairs[key] = value

		if !p.peekTokenIs(token.DelimiterRightBrace) && !p.expectPeek(token.DelimiterComma) {
			return nil
		}
	}

	if !p.expectPeek(token.DelimiterRightBrace) {
		return nil
	}

	return hash
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.DelimiterRightBracket) {
		return nil
	}

	return exp
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	return &ast.CallExpression{
		Token:     p.curToken,
		Function:  function,
		Arguments: p.parseExpressionList(token.DelimiterRightParenthesis),
	}
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	return &ast.ArrayLiteral{
		Token:    p.curToken,
		Elements: p.parseExpressionList(token.DelimiterRightBracket),
	}
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.DelimiterComma) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenIs(token.DelimiterRightParenthesis) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.DelimiterComma) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.DelimiterRightParenthesis) {
		return nil
	}

	return args
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{Token: p.curToken, Value: p.curTokenIs(token.KeywordTrue)}
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(token.DelimiterLeftParenthesis) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.DelimiterLeftBrace) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {

	if p.peekTokenIs(token.DelimiterRightParenthesis) {
		p.nextToken()
		return nil
	}

	p.nextToken()

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	identifiers := []*ast.Identifier{ident}

	for p.peekTokenIs(token.DelimiterComma) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.DelimiterRightParenthesis) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIfExpression() ast.Expression {

	if !p.expectPeek(token.DelimiterLeftParenthesis) {
		return nil
	}

	p.nextToken()
	condition := p.parseExpression(LOWEST)

	if !p.expectPeek(token.DelimiterRightParenthesis) {
		return nil
	}

	if !p.expectPeek(token.DelimiterLeftBrace) {
		return nil
	}

	consequence := p.parseBlockStatement()

	var alternative *ast.BlockStatement
	if p.peekTokenIs(token.KeywordElse) {
		p.nextToken()

		if !p.expectPeek(token.DelimiterLeftBrace) {
			return nil
		}

		alternative = p.parseBlockStatement()
	}

	return &ast.IfExpression{
		Token:       p.curToken,
		Condition:   condition,
		Consequence: consequence,
		Alternative: alternative,
	}
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	stmts := make([]ast.Statement, 0)

	p.nextToken()

	for !p.curTokenIs(token.DelimiterRightBrace) && !p.curTokenIs(token.Eof) {
		stmt := p.parseStatement()
		if stmt != nil {
			stmts = append(stmts, stmt)
		}
		p.nextToken()
	}

	return &ast.BlockStatement{
		Token:      p.curToken,
		Statements: stmts,
	}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.DelimiterRightParenthesis) {
		return nil
	}

	return exp
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = <-p.tokens
}

func (p *Parser) ParseProgram() (*ast.Program, error) {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.Eof {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	if len(p.errors) > 0 {
		errMsg := "error while parsing input: " + strings.Join(p.errors, ", ")
		return program, errors.New(errMsg)
	}

	return program, nil
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.KeywordLet:
		return p.parseLetStatement()
	case token.KeywordReturn:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.DelimiterSemicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.Ident) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.OperatorAssign) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.DelimiterSemicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}

	p.peekError(t)
	return false
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.DelimiterSemicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.DelimiterSemicolon) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}
