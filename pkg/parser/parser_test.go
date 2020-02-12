package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/adrian83/monkey/pkg/ast"
	"github.com/adrian83/monkey/pkg/lexer"
	"github.com/adrian83/monkey/pkg/token"
)

func TestLetStatements(t *testing.T) {
	testData := map[string]struct {
		input string
		name  string
		value int64
	}{
		"first":  {"let x = 5;", "x", 5},
		"second": {"let y = 10;", "y", 10},
		"third":  {"let foobar = 838383;", "foobar", 838383},
	}

	for name, tData := range testData {
		data := tData

		t.Run(name, func(t *testing.T) {
			l := lexer.New(data.input)
			p := New(l)

			program := p.ParseProgram()
			assertNoParsingErrors(t, p.Errors())
			assertNotEmptyProgram(t, program)

			stmt := program.Statements[0]
			letStmt := assertLetStatement(t, stmt)
			assertIdentifierValue(t, data.name, letStmt.Name.Value)
			assertTokenLiteral(t, strings.ToLower(token.LET), stmt.TokenLiteral())
		})
	}
}

func TestReturnStatements(t *testing.T) {
	testData := map[string]struct {
		input string
		value int64
	}{
		"return 5":      {"return 5;", 5},
		"return 10":     {"return 10;", 10},
		"return 993322": {"return 993322;", 993322},
	}

	for name, tData := range testData {
		data := tData

		t.Run(name, func(t *testing.T) {
			l := lexer.New(data.input)
			p := New(l)

			program := p.ParseProgram()
			assertNoParsingErrors(t, p.Errors())
			assertNotEmptyProgram(t, program)

			returnStmt := assertReturnStatement(t, program.Statements[0])
			assertTokenLiteral(t, strings.ToLower(token.RETURN), returnStmt.TokenLiteral())
		})
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	assertNoParsingErrors(t, p.Errors())
	assertNotEmptyProgram(t, program)

	stmt := assertExpressionStatement(t, program.Statements[0])
	ident := assertIdentifier(t, stmt.Expression)
	assertIdentifierValue(t, "foobar", ident.Value)
	assertTokenLiteral(t, "foobar", ident.TokenLiteral())
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	assertNoParsingErrors(t, p.Errors())
	assertNotEmptyProgram(t, program)

	stmt := assertExpressionStatement(t, program.Statements[0])
	intLit := assertIntegerLiteral(t, stmt.Expression)
	if intLit.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d", 5, intLit.Value)
	}

	assertTokenLiteral(t, "5", intLit.TokenLiteral())
}

func TestParsingPrefixExpressions(t *testing.T) {
	testData := map[string]struct {
		input    string
		operator string
		value    int64
	}{
		"bang":  {"!5;", "!", 5},
		"minus": {"-15;", "-", 15},
	}

	for name, tData := range testData {
		data := tData

		t.Run(name, func(t *testing.T) {

			l := lexer.New(data.input)
			p := New(l)
			program := p.ParseProgram()

			assertNoParsingErrors(t, p.Errors())
			assertNotEmptyProgram(t, program)

			stmt := assertExpressionStatement(t, program.Statements[0])

			exp, ok := stmt.Expression.(*ast.PrefixExpression)
			if !ok {
				t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
			}
			if exp.Operator != data.operator {
				t.Fatalf("exp.Operator is not '%s'. got=%s", data.operator, exp.Operator)
			}
			if !testIntegerLiteral(t, exp.Right, data.value) {
				return
			}
		})
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		assertNoParsingErrors(t, p.Errors())

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		exprStmt := assertExpressionStatement(t, program.Statements[0])

		exp, ok := exprStmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp is not ast.InfixExpression. got=%T", exprStmt.Expression)
		}

		if !testIntegerLiteral(t, exp.Left, tt.leftValue) {
			return
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}

		if !testIntegerLiteral(t, exp.Right, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		assertNoParsingErrors(t, p.Errors())

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func assertTokenLiteral(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Errorf("invalid literal, expected: %v, actual: %v", expected, actual)
	}
}

func assertIdentifierValue(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Errorf("invalid identifier value, expected: %v, actual: %v", expected, actual)
	}
}

func assertNoParsingErrors(t *testing.T, errors []string) {
	if len(errors) > 0 {
		t.Errorf("unexpected parser errors: %v", errors)
	}
}

func assertExpressionStatement(t *testing.T, stmt ast.Statement) *ast.ExpressionStatement {
	expStmt, ok := stmt.(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("invalid type of statement: %v, expected: *ast.ExpressionStatement, actual: %T", stmt, stmt)
	}

	return expStmt
}

func assertLetStatement(t *testing.T, stmt ast.Statement) *ast.LetStatement {
	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("invalid type of statement: %v, expected: *ast.LetStatement, actual: %T", stmt, stmt)
	}

	return letStmt
}

func assertIdentifier(t *testing.T, exp ast.Expression) *ast.Identifier {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("invalid type of expression: %v, expected: *ast.Identifier, actual: %T", ident, ident)
	}

	return ident
}

func assertIntegerLiteral(t *testing.T, exp ast.Expression) *ast.IntegerLiteral {
	intLit, ok := exp.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("invalid type of expression: %v, expected: *ast.IntegerLiteral, actual: %T", intLit, intLit)
	}

	return intLit
}

func assertReturnStatement(t *testing.T, stmt ast.Statement) *ast.ReturnStatement {
	retStmt, ok := stmt.(*ast.ReturnStatement)
	if !ok {
		t.Errorf("invalid type of statement: %v, expected: *ast.ReturnStatement, actual: %T", stmt, stmt)
	}

	return retStmt
}

func assertNotEmptyProgram(t *testing.T, program *ast.Program) {
	if len(program.Statements) == 0 {
		t.Error("program doesn't contain any statements")
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value,
			integ.TokenLiteral())
		return false
	}

	return true
}
