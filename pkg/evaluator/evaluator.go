package evaluator

import (
	"github.com/adrian83/monkey/pkg/ast"
	"github.com/adrian83/monkey/pkg/object"
)

var (
	objTrue  = &object.Boolean{Value: true}
	objFalse = &object.Boolean{Value: false}
	objNull  = &object.Null{}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalStatements(node.Statements)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)

	}

	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement)
	}

	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return objTrue
	}
	return objFalse
}
