package parser

import (
	"strings"
	"testing"

	"github.com/adrian83/monkey/pkg/ast"
	"github.com/adrian83/monkey/pkg/lexer"
	"github.com/adrian83/monkey/pkg/token"

	"github.com/stretchr/testify/assert"
)

func TestLetStatements(t *testing.T) {
	testData := map[string]struct {
		input string
		name  string
		value interface{}
	}{
		"first":  {"let x = 5;", "x", 5},
		"second": {"let y = 10;", "y", 10},
		"third":  {"let foobar = 838383;", "foobar", 838383},
		"forth":  {"let x = 5;", "x", 5},
		"fifth":  {"let y = true;", "y", true},
		"sixth":  {"let foobar = y;", "foobar", "y"},
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

			assertLiteral(t, letStmt.Value, data.value)
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
func TestBooleanExpression(t *testing.T) {
	testData := map[string]struct {
		input   string
		value   bool
		literal string
	}{
		"only true":            {"true", true, token.TRUE},
		"true with semicolon":  {"true;", true, token.TRUE},
		"only false":           {"false", false, token.FALSE},
		"flase with semicolon": {"false;", false, token.FALSE},
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
			boolLit := assertBooleanLiteral(t, stmt.Expression)
			if boolLit.Value != data.value {
				t.Errorf("literal.Value not %v. got=%v", data.value, boolLit.Value)
			}

			assertTokenLiteral(t, strings.ToLower(data.literal), boolLit.TokenLiteral())
		})
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	testData := map[string]struct {
		input    string
		operator string
		value    int64
	}{
		"bang and int":  {"!5;", "!", 5},
		"minus and int": {"-15;", "-", 15},
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

			assertLiteral(t, exp.Right, data.value)
		})
	}
}

func TestParsingBooleanPrefixExpressions(t *testing.T) {
	testData := map[string]struct {
		input    string
		operator string
		value    bool
	}{
		"bang and bool true":  {"!true;", "!", true},
		"bang and bool false": {"!false;", "!", false},
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

			assertLiteral(t, exp.Right, data.value)
		})
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true != false;", true, "!=", false},
		{"true == true;", true, "==", true},
		{"false == false;", false, "==", false},
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

		assertInfixExpression(t, exprStmt.Expression, tt.leftValue, tt.rightValue, tt.operator)
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
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
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
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

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	assertNoParsingErrors(t, p.Errors())

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T",
			stmt.Expression)
	}

	assertInfixExpression(t, exp.Condition, "x", "y", token.LT)

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	ident := assertIdentifier(t, consequence.Expression)
	assertIdentifierValue(t, "x", ident.Value)

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil. got=%+v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	assertNoParsingErrors(t, p.Errors())

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T",
			stmt.Expression)
	}

	assertInfixExpression(t, exp.Condition, "x", "y", token.LT)

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	ident := assertIdentifier(t, consequence.Expression)
	assertIdentifierValue(t, "x", ident.Value)

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Alternative.Statements[0])
	}

	altIdent := assertIdentifier(t, alternative.Expression)
	assertIdentifierValue(t, "y", altIdent.Value)
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	assertNoParsingErrors(t, p.Errors())

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T",
			stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d\n",
			len(function.Parameters))
	}

	assertLiteral(t, function.Parameters[0], "x")
	assertLiteral(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n", len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T",
			function.Body.Statements[0])
	}

	assertInfixExpression(t, bodyStmt.Expression, "x", "y", token.PLUS)
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		assertNoParsingErrors(t, p.Errors())

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. want %d, got=%d\n",
				len(tt.expectedParams), len(function.Parameters))
		}

		for i, ident := range tt.expectedParams {
			//testLiteralExpression(t, function.Parameters[i], ident)
			assertLiteral(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	assertNoParsingErrors(t, p.Errors())

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
			stmt.Expression)
	}

	altIdent := assertIdentifier(t, exp.Function)
	assertIdentifierValue(t, "add", altIdent.Value)

	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}

	assertLiteral(t, exp.Arguments[0], 1)
	assertInfixExpression(t, exp.Arguments[1], 2, 3, token.ASTERISK)
	assertInfixExpression(t, exp.Arguments[2], 4, 5, token.PLUS)
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

func assertBooleanLiteral(t *testing.T, exp ast.Expression) *ast.BooleanLiteral {
	boolLit, ok := exp.(*ast.BooleanLiteral)
	if !ok {
		t.Errorf("invalid type of expression: %v, expected: *ast.BooleanLiteral, actual: %T", boolLit, boolLit)
	}

	return boolLit
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

func assertInfixExpression(t *testing.T, expr ast.Expression, leftVal, rightVal interface{}, operator string) {

	exp, ok := expr.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("exp is not ast.InfixExpression. got=%T", exp)
	}

	if exp.Operator != operator {
		t.Fatalf("exp.Operator is not '%s'. got=%s", operator, exp.Operator)
	}

	assertLiteral(t, exp.Left, leftVal)
	assertLiteral(t, exp.Right, rightVal)
}

func assertLiteral(t *testing.T, il ast.Expression, value interface{}) {

	switch v := value.(type) {
	case bool:
		boolLit, ok := il.(*ast.BooleanLiteral)
		if !ok {
			t.Errorf("il not *ast.BooleanLiteral. got=%T", il)
		}
		assert.Equal(t, v, boolLit.Value)

	case int:
		intLit, ok := il.(*ast.IntegerLiteral)
		if !ok {
			t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		}
		assert.Equal(t, int64(v), intLit.Value)

	case int64:
		intLit, ok := il.(*ast.IntegerLiteral)
		if !ok {
			t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		}
		assert.Equal(t, v, intLit.Value)
	case string:
		ident, ok := il.(*ast.Identifier)
		if !ok {
			t.Errorf("il not *ast.Identifier. got=%T", il)
		}
		assert.Equal(t, v, ident.Value)
	default:
		t.Errorf("unknown type, got=%T", value)
	}
}
