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
		"case 1": {"let x = 5;", "x", 5},
		"case 2": {"let y = 10;", "y", 10},
		"case 3": {"let foobar = 838383;", "foobar", 838383},
		"case 4": {"let x = 5;", "x", 5},
		"case 5": {"let y = true;", "y", true},
		"case 6": {"let foobar = y;", "foobar", "y"},
	}

	for name, tData := range testData {
		data := tData

		t.Run(name, func(t *testing.T) {
			program := parseProgram(t, data.input)

			assertStatementsCount(t, program, 1)

			letStmt := assertLetStatement(t, program.Statements[0])

			assertIdentifierValue(t, data.name, letStmt.Name.Value)
			assertTokenLiteral(t, token.LET, letStmt.NodeToken().Literal)
			assertLiteral(t, letStmt.Value, data.value)
		})
	}
}

func TestReturnStatements(t *testing.T) {
	testData := map[string]struct {
		input string
		value int64
	}{
		"case 1": {"return 5;", 5},
		"case 2": {"return 10;", 10},
		"case 3": {"return 993322;", 993322},
	}

	for name, tData := range testData {
		data := tData

		t.Run(name, func(t *testing.T) {
			program := parseProgram(t, data.input)

			assertStatementsCount(t, program, 1)

			returnStmt := assertReturnStatement(t, program.Statements[0])

			assertTokenLiteral(t, token.RETURN, returnStmt.NodeToken().Literal)
			assertLiteral(t, returnStmt.ReturnValue, data.value)
		})
	}
}

func TestIdentifierExpression(t *testing.T) {
	testData := map[string]struct {
		input   string
		literal string
	}{
		"case 1": {"foo", "foo"},
		"case 2": {"bar;", "bar"},
	}
	for name, tData := range testData {
		data := tData

		t.Run(name, func(t *testing.T) {
			program := parseProgram(t, data.input)

			assertStatementsCount(t, program, 1)

			stmt := toExpressionStatement(t, program.Statements[0])
			ident := assertIdentifier(t, stmt.Expression)

			assertIdentifierValue(t, data.literal, ident.Value)
			assertTokenLiteral(t, data.literal, ident.NodeToken().Literal)
		})
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	testData := map[string]struct {
		input   string
		value   int64
		literal string
	}{
		"case 1": {"5", 5, "5"},
		"case 2": {"6;", 6, "6"},
	}

	for name, tData := range testData {
		data := tData

		t.Run(name, func(t *testing.T) {
			program := parseProgram(t, data.input)

			assertStatementsCount(t, program, 1)

			stmt := toExpressionStatement(t, program.Statements[0])

			assertLiteral(t, stmt.Expression, data.value)
			assertTokenLiteral(t, data.literal, stmt.NodeToken().Literal)
		})
	}
}

func TestStringLiteralExpression(t *testing.T) {
	testData := map[string]struct {
		input string
		value string
	}{
		"case 1": {`"foo"`, "foo"},
		"case 2": {`"bar";`, "bar"},
	}

	for name, tData := range testData {
		data := tData

		t.Run(name, func(t *testing.T) {
			program := parseProgram(t, data.input)

			assertStatementsCount(t, program, 1)

			expStmt := toExpressionStatement(t, program.Statements[0])
			strLit := toStringLiteral(t, expStmt.Expression)

			assertLiteral(t, strLit, data.value)
			assertTokenLiteral(t, data.value, strLit.NodeToken().Literal)
		})
	}
}

func TestBooleanExpression(t *testing.T) {
	testData := map[string]struct {
		input   string
		value   bool
		literal string
	}{
		"case 1": {"true", true, token.TRUE},
		"case 2": {"true;", true, token.TRUE},
		"case 3": {"false", false, token.FALSE},
		"case 4": {"false;", false, token.FALSE},
	}

	for name, tData := range testData {
		data := tData

		t.Run(name, func(t *testing.T) {
			program := parseProgram(t, data.input)

			assertStatementsCount(t, program, 1)

			stmt := toExpressionStatement(t, program.Statements[0])
			boolLit := assertBooleanLiteral(t, stmt.Expression)

			assertLiteral(t, stmt.Expression, data.value)
			assertTokenLiteral(t, data.literal, boolLit.NodeToken().Literal)
		})
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	testData := map[string]struct {
		input    string
		operator string
		value    interface{}
	}{
		"case 1": {"!5;", "!", 5},
		"case 2": {"-15;", "-", 15},
		"case 3": {"!true;", "!", true},
		"case 4": {"!false;", "!", false},
	}

	for name, tData := range testData {
		data := tData

		t.Run(name, func(t *testing.T) {
			program := parseProgram(t, data.input)

			assertStatementsCount(t, program, 1)

			stmt := toExpressionStatement(t, program.Statements[0])
			prfxExp := assertPrefixExpression(t, stmt.Expression)

			assertOperator(t, prfxExp.Operator, data.operator)
			assertLiteral(t, prfxExp.Right, data.value)
		})
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	testData := map[string]struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		"case 1":  {"5 + 5;", 5, "+", 5},
		"case 2":  {"5 - 5;", 5, "-", 5},
		"case 3":  {"5 * 5;", 5, "*", 5},
		"case 4":  {"5 / 5;", 5, "/", 5},
		"case 5":  {"5 > 5;", 5, ">", 5},
		"case 6":  {"5 < 5;", 5, "<", 5},
		"case 7":  {"5 == 5;", 5, "==", 5},
		"case 8":  {"5 != 5;", 5, "!=", 5},
		"case 9":  {"true != false;", true, "!=", false},
		"case 10": {"true == true;", true, "==", true},
		"case 11": {"false == false;", false, "==", false},
	}

	for name, tData := range testData {
		data := tData

		t.Run(name, func(t *testing.T) {
			program := parseProgram(t, data.input)

			assertStatementsCount(t, program, 1)

			exprStmt := toExpressionStatement(t, program.Statements[0])
			infixExp := assertInfixExpression(t, exprStmt.Expression)

			assertOperator(t, data.operator, infixExp.Operator)
			assertLiteral(t, infixExp.Left, data.leftValue)
			assertLiteral(t, infixExp.Right, data.rightValue)
		})
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	testData := map[string]struct {
		input      string
		expected   string
		statements int
	}{
		"case 1":  {"true", "true", 1},
		"case 2":  {"false", "false", 1},
		"case 3":  {"3 > 5 == false", "((3 > 5) == false)", 1},
		"case 4":  {"3 < 5 == true", "((3 < 5) == true)", 1},
		"case 5":  {"-a * b", "((-a) * b)", 1},
		"case 6":  {"!-a", "(!(-a))", 1},
		"case 7":  {"a + b + c", "((a + b) + c)", 1},
		"case 8":  {"a + b - c", "((a + b) - c)", 1},
		"case 9":  {"a * b * c", "((a * b) * c)", 1},
		"case 10": {"a * b / c", "((a * b) / c)", 1},
		"case 11": {"a + b / c", "(a + (b / c))", 1},
		"case 12": {"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)", 1},
		"case 13": {"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)", 2},
		"case 14": {"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))", 1},
		"case 15": {"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))", 1},
		"case 16": {"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))", 1},
		"case 17": {"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)", 1},
		"case 18": {"(5 + 5) * 2", "((5 + 5) * 2)", 1},
		"case 19": {"2 / (5 + 5)", "(2 / (5 + 5))", 1},
		"case 20": {"-(5 + 5)", "(-(5 + 5))", 1},
		"case 21": {"!(true == true)", "(!(true == true))", 1},
		"case 22": {"a + add(b * c) + d", "((a + add((b * c))) + d)", 1},
		"case 23": {"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))", "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))", 1},
		"case 24": {"add(a + b + c * d / f + g)", "add((((a + b) + ((c * d) / f)) + g))", 1},
		"case 25": {"a * [1, 2, 3, 4][b * c] * d", "((a * ([1, 2, 3, 4][(b * c)])) * d)", 1},
		"case 26": {"add(a * b[2], b[1], 2 * [1, 2][1])", "add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))", 1},
	}

	for name, tData := range testData {
		data := tData

		t.Run(name, func(t *testing.T) {
			program := parseProgram(t, data.input)

			assertStatementsCount(t, program, data.statements)

			actual := program.String()
			assert.Equal(t, data.expected, actual)
		})
	}
}

func TestIfExpression(t *testing.T) {
	testData := map[string]struct {
		input        string
		condLeft     string
		condRight    string
		condOperator string
		consequence  string
	}{
		"case 1": {`if (x < y) { x }`, "x", "y", token.LT, "x"},
		"case 2": {`if (a == b) { b }`, "a", "b", token.EQ, "b"},
	}

	for name, tData := range testData {
		data := tData

		t.Run(name, func(t *testing.T) {
			program := parseProgram(t, data.input)

			assertStatementsCount(t, program, 1)

			expStmt := toExpressionStatement(t, program.Statements[0])
			ifExp := assertIfExpression(t, expStmt.Expression)

			infixExp := assertInfixExpression(t, ifExp.Condition)

			assertOperator(t, data.condOperator, infixExp.Operator)
			assertLiteral(t, infixExp.Left, data.condLeft)
			assertLiteral(t, infixExp.Right, data.condRight)

			assertStatementsCount(t, ifExp.Consequence, 1)

			consExpStmt := toExpressionStatement(t, ifExp.Consequence.Statements[0])
			ident := assertIdentifier(t, consExpStmt.Expression)

			assertIdentifierValue(t, data.consequence, ident.Value)

			assert.Nil(t, ifExp.Alternative)
		})
	}
}

func TestIfElseExpression(t *testing.T) {
	testData := map[string]struct {
		input        string
		condLeft     string
		condRight    string
		condOperator string
		consequence  string
		alternative  string
	}{
		"case 1": {`if (x < y) { x } else { y }`, "x", "y", token.LT, "x", "y"},
		"case 2": {`if (a == b) { b } else { a }`, "a", "b", token.EQ, "b", "a"},
	}

	for name, tData := range testData {
		data := tData

		t.Run(name, func(t *testing.T) {
			program := parseProgram(t, data.input)

			assertStatementsCount(t, program, 1)

			expStmt := toExpressionStatement(t, program.Statements[0])
			ifExp := assertIfExpression(t, expStmt.Expression)

			infixExp := assertInfixExpression(t, ifExp.Condition)

			assertOperator(t, data.condOperator, infixExp.Operator)
			assertLiteral(t, infixExp.Left, data.condLeft)
			assertLiteral(t, infixExp.Right, data.condRight)

			assertStatementsCount(t, ifExp.Consequence, 1)

			consExpStmt := toExpressionStatement(t, ifExp.Consequence.Statements[0])
			consIdent := assertIdentifier(t, consExpStmt.Expression)

			assertIdentifierValue(t, data.consequence, consIdent.Value)

			altExpStmt := toExpressionStatement(t, ifExp.Alternative.Statements[0])
			altIdent := assertIdentifier(t, altExpStmt.Expression)

			assertIdentifierValue(t, data.alternative, altIdent.Value)
		})
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	testData := map[string]struct {
		input   string
		params  []string
		bodyExp []string
	}{
		"case 1": {
			`fn(x, y) { x + y; }`,
			[]string{"x", "y"},
			[]string{"x", token.PLUS, "y"},
		},
		"case 2": {
			`fn(a, b) { a * b }`,
			[]string{"a", "b"},
			[]string{"a", token.ASTERISK, "b"},
		},
	}

	for name, tData := range testData {
		data := tData

		t.Run(name, func(t *testing.T) {
			program := parseProgram(t, data.input)

			assertStatementsCount(t, program, 1)

			expStmt := toExpressionStatement(t, program.Statements[0])
			function := assertFunctionLiteral(t, expStmt.Expression)

			assertParameters(t, function.Parameters, data.params)
			assertStatementsCount(t, function.Body, 1)

			bodyExpStmt := toExpressionStatement(t, function.Body.Statements[0])
			infixExp := assertInfixExpression(t, bodyExpStmt.Expression)

			assertOperator(t, data.bodyExp[1], infixExp.Operator)
			assertLiteral(t, infixExp.Left, data.bodyExp[0])
			assertLiteral(t, infixExp.Right, data.bodyExp[2])
		})
	}
}

func TestFunctionParameterParsing(t *testing.T) {
	testData := map[string]struct {
		input          string
		expectedParams []string
	}{
		"case 1": {input: "fn() {};", expectedParams: []string{}},
		"case 2": {input: "fn(x) {};", expectedParams: []string{"x"}},
		"case 3": {input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for name, tData := range testData {
		data := tData

		t.Run(name, func(t *testing.T) {
			program := parseProgram(t, data.input)

			assertStatementsCount(t, program, 1)

			expStmt := toExpressionStatement(t, program.Statements[0])
			function := assertFunctionLiteral(t, expStmt.Expression)

			assertParameters(t, function.Parameters, data.expectedParams)
		})
	}
}

func TestCallExpressionParsing(t *testing.T) {
	type param struct {
		leftVal  int64
		rightVal int64
		operator string
	}

	testData := map[string]struct {
		input    string
		funcName string
		params   []param
	}{
		"case 1": {
			"add(5 - 1, 2 * 3, 4 + 5);",
			"add",
			[]param{
				{5, 1, token.MINUS},
				{2, 3, token.ASTERISK},
				{4, 5, token.PLUS},
			},
		},
		"case 2": {
			"sub(2 + 2, 7 * 7)",
			"sub",
			[]param{
				{2, 2, token.PLUS},
				{7, 7, token.ASTERISK},
			},
		},
	}

	for name, tData := range testData {
		data := tData

		t.Run(name, func(t *testing.T) {
			program := parseProgram(t, data.input)

			assertStatementsCount(t, program, 1)

			expStmt := toExpressionStatement(t, program.Statements[0])
			expCall := assertCallExpression(t, expStmt.Expression)

			altIdent := assertIdentifier(t, expCall.Function)
			assertIdentifierValue(t, data.funcName, altIdent.Value)

			assert.Len(t, expCall.Arguments, len(data.params))
			for i, arg := range expCall.Arguments {
				infixExp := assertInfixExpression(t, arg)

				assertOperator(t, infixExp.Operator, data.params[i].operator)
				assertLiteral(t, infixExp.Left, data.params[i].leftVal)
				assertLiteral(t, infixExp.Right, data.params[i].rightVal)
			}
		})
	}
}

func TestParsingArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	program := parseProgram(t, input)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Expression)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
	}

	assertLiteral(t, array.Elements[0], 1)

	infixExp1 := assertInfixExpression(t, array.Elements[1])
	assertOperator(t, infixExp1.Operator, token.ASTERISK)
	assertLiteral(t, infixExp1.Left, 2)
	assertLiteral(t, infixExp1.Right, 2)

	infixExp2 := assertInfixExpression(t, array.Elements[2])
	assertOperator(t, infixExp2.Operator, token.PLUS)
	assertLiteral(t, infixExp2.Left, 3)
	assertLiteral(t, infixExp2.Right, 3)
}

func TestParsingIndexExpressions(t *testing.T) {
	input := "myArray[1 + 2]"

	program := parseProgram(t, input)

	stmt := toExpressionStatement(t, program.Statements[0])

	indexExp := toIndexExpression(t, stmt.Expression)

	assertLiteral(t, indexExp.Left, "myArray")

	infixExp := assertInfixExpression(t, indexExp.Index)
	assertOperator(t, infixExp.Operator, token.PLUS)
	assertLiteral(t, infixExp.Left, 1)
	assertLiteral(t, infixExp.Right, 2)
}

func TestParsingHashLiteralsStringKeys(t *testing.T) {
	testData := map[string]struct {
		input    string
		expected map[string]int64
	}{
		"case 1": {`{"one": 1, "two": 2, "three": 3}`, map[string]int64{"one": 1, "two": 2, "three": 3}},
		"case 2": {`{}`, nil},
	}

	for name, tData := range testData {
		data := tData

		t.Run(name, func(t *testing.T) {
			program := parseProgram(t, data.input)

			expStmt := toExpressionStatement(t, program.Statements[0])
			hashLit := toHashLiteral(t, expStmt.Expression)

			if len(hashLit.Pairs) != len(data.expected) {
				t.Errorf("hash.Pairs has wrong length. got=%d", len(hashLit.Pairs))
			}

			for key, value := range hashLit.Pairs {

				strLit := toStringLiteral(t, key)
				expectedValue := data.expected[strLit.String()]

				assertLiteral(t, value, expectedValue)
			}
		})
	}
}

// func TestParsingHashLiteralsWithExpressions(t *testing.T) {
// 	input := `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`

// 	program := parseProgram(t, input)

// 	stmt := program.Statements[0].(*ast.ExpressionStatement)
// 	hash, ok := stmt.Expression.(*ast.HashLiteral)
// 	if !ok {
// 		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
// 	}

// 	if len(hash.Pairs) != 3 {
// 		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
// 	}

// 	tests := map[string]func(ast.Expression){
// 		"one": func(e ast.Expression) {
// 			testInfixExpression(t, e, 0, "+", 1)
// 		},
// 		"two": func(e ast.Expression) {
// 			testInfixExpression(t, e, 10, "-", 8)
// 		},
// 		"three": func(e ast.Expression) {
// 			testInfixExpression(t, e, 15, "/", 5)
// 		},
// 	}

// 	for key, value := range hash.Pairs {
// 		literal, ok := key.(*ast.StringLiteral)
// 		if !ok {
// 			t.Errorf("key is not ast.StringLiteral. got=%T", key)
// 			continue
// 		}

// 		testFunc, ok := tests[literal.String()]
// 		if !ok {
// 			t.Errorf("No test function for key %q found", literal.String())
// 			continue
// 		}

// 		testFunc(value)
// 	}
// }

func assertTokenLiteral(t *testing.T, expected, actual string) {
	if strings.ToLower(expected) != strings.ToLower(actual) {
		t.Errorf("invalid literal, expected: %v, actual: %v", expected, actual)
	}
}

func toStringLiteral(t *testing.T, exp ast.Expression) *ast.StringLiteral {
	strLit, ok := exp.(*ast.StringLiteral)
	if !ok {
		t.Errorf("invalid ast.Expression type, expected: *ast.StringLiteral, actual: %T", exp)
	}
	return strLit
}

func toHashLiteral(t *testing.T, exp ast.Expression) *ast.HashLiteral {
	hashLit, ok := exp.(*ast.HashLiteral)
	if !ok {
		t.Errorf("invalid ast.Expression type, expected: *ast.HashLiteral, actual: %T", exp)
	}
	return hashLit
}

func toIndexExpression(t *testing.T, exp ast.Expression) *ast.IndexExpression {
	indexExp, ok := exp.(*ast.IndexExpression)
	if !ok {
		t.Errorf("invalid ast.Expression type, expected: *ast.IndexExpression, actual: %T", exp)
	}
	return indexExp
}

func assertIdentifierValue(t *testing.T, expected, actual string) {
	if strings.ToLower(expected) != strings.ToLower(actual) {
		t.Errorf("invalid identifier value, expected: %v, actual: %v", expected, actual)
	}
}

func toExpressionStatement(t *testing.T, stmt ast.Statement) *ast.ExpressionStatement {
	expStmt, ok := stmt.(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("invalid ast.Statement type, expected: *ast.ExpressionStatement, actual: %T", stmt)
	}
	return expStmt
}

func assertCallExpression(t *testing.T, exp ast.Expression) *ast.CallExpression {
	callExp, ok := exp.(*ast.CallExpression)
	if !ok {
		t.Errorf("invalid ast.Expression type, expected: *ast.CallExpression, actual: %T", exp)
	}
	return callExp
}

func assertLetStatement(t *testing.T, stmt ast.Statement) *ast.LetStatement {
	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("invalid ast.Statement type, expected: *ast.LetStatement, actual: %T", stmt)
	}

	return letStmt
}

func assertIdentifier(t *testing.T, exp ast.Expression) *ast.Identifier {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("invalid ast.Expression type, expected: *ast.Identifier, actual: %T", exp)
	}
	return ident
}

func assertPrefixExpression(t *testing.T, exp ast.Expression) *ast.PrefixExpression {
	prfxExp, ok := exp.(*ast.PrefixExpression)
	if !ok {
		t.Errorf("invalid ast.Expression type, expected: *ast.PrefixExpression, actual: %T", exp)
	}
	return prfxExp
}

func assertFunctionLiteral(t *testing.T, exp ast.Expression) *ast.FunctionLiteral {
	funcLit, ok := exp.(*ast.FunctionLiteral)
	if !ok {
		t.Errorf("invalid ast.Expression type, expected: *ast.FunctionLiteral, actual: %T", exp)
	}
	return funcLit
}

func assertIfExpression(t *testing.T, exp ast.Expression) *ast.IfExpression {
	ifExp, ok := exp.(*ast.IfExpression)
	if !ok {
		t.Errorf("invalid ast.Expression type, expected: *ast.IfExpression, actual: %T", exp)
	}
	return ifExp
}

func assertIntegerLiteral(t *testing.T, exp ast.Expression) *ast.IntegerLiteral {
	intLit, ok := exp.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("invalid ast.Expression type, expected: *ast.IntegerLiteral, actual: %T", exp)
	}
	return intLit
}

func assertBooleanLiteral(t *testing.T, exp ast.Expression) *ast.BooleanLiteral {
	boolLit, ok := exp.(*ast.BooleanLiteral)
	if !ok {
		t.Errorf("invalid ast.Expression type, expected: *ast.BooleanLiteral, actual: %T", exp)
	}
	return boolLit
}

func assertReturnStatement(t *testing.T, stmt ast.Statement) *ast.ReturnStatement {
	retStmt, ok := stmt.(*ast.ReturnStatement)
	if !ok {
		t.Errorf("invalid ast.Statement type, expected: *ast.ReturnStatement, actual: %T", stmt)
	}
	return retStmt
}

func assertInfixExpression(t *testing.T, exp ast.Expression) *ast.InfixExpression {
	infixExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("invalid ast.Expression type, expected: *ast.InfixExpression, actual: %T", exp)
	}
	return infixExp
}

func assertLiteral(t *testing.T, il ast.Expression, value interface{}) {

	switch expVal := il.(type) {
	case *ast.BooleanLiteral:
		assert.Equal(t, value, expVal.Value)
	case *ast.IntegerLiteral:

		switch numbVal := value.(type) {
		case int:
			assert.Equal(t, int64(numbVal), expVal.Value)
		case int64:
			assert.Equal(t, numbVal, expVal.Value)
		default:
			t.Errorf("unsuported number type: %T", numbVal)
		}

	case *ast.Identifier:
		assert.Equal(t, value, expVal.Value)
	case *ast.StringLiteral:
		assert.Equal(t, value, expVal.Value)
	default:
		t.Errorf("unsuported ast.Expression type: %T", il)
	}
}

func parseProgram(t *testing.T, input string) *ast.Program {
	l := lexer.New(input)
	p := New(l)

	program, err := p.ParseProgram()
	if err != nil {
		t.Errorf("cannot parse input, error: %v", err)
	}

	if len(program.Statements) == 0 {
		t.Error("empty program")
	}

	return program
}

func assertStatementsCount(t *testing.T, bodyHolder ast.BodyHolder, count int) {
	if len(bodyHolder.BodyStatements()) != count {
		t.Fatalf("invalid number of statements in bodyHolder.BodyStatements(), expected: %v, actual, %v", count, len(bodyHolder.BodyStatements()))
	}
}

func assertOperator(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Fatalf("invalid operator, expected: %v, actual: %v\n", expected, actual)
	}
}

func assertParameters(t *testing.T, params []*ast.Identifier, expected []string) {
	assert.Len(t, params, len(expected))

	for i, param := range params {
		assertLiteral(t, param, expected[i])
	}
}
