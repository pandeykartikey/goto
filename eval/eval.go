package eval

import (
	"goto/ast"
	"goto/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func nativeBoolToBooleanObject(input bool) object.Object {
	if input {
		return TRUE
	}
	return FALSE
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt)
	}

	return result
}

// NULL is false and all other values are true
func evalNotOperator(obj object.Object) object.Object {
	switch obj {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

// - Operator can only apply on integer value
func evalNegateOperator(obj object.Object) object.Object {
	intobj, ok := obj.(*object.Integer)
	if !ok {
		// raise error
		return NULL
	}

	intobj.Value = -intobj.Value

	return intobj
}

func evalPrefixExpression(op string, right object.Object) object.Object {
	switch op {
	case "!":
		return evalNotOperator(right)
	case "-":
		return evalNegateOperator(right)
	default:
		return NULL
	}
}

func evalInfixIntegerExpression(op string, left *object.Integer, right *object.Integer) object.Object {
	leftVal, rightVal := left.Value, right.Value

	switch op {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	default:
		return NULL
	}
}

func evalInfixBooleanExpression(op string, left *object.Boolean, right *object.Boolean) object.Object {
	leftVal, rightVal := left.Value, right.Value

	switch op {
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return NULL
	}
}

func evalInfixExpression(op string, left object.Object, right object.Object) object.Object {
	if left.Type() != right.Type() {
		// raise error
		return NULL
	}

	switch left.(type) {
	case *object.Integer:
		return evalInfixIntegerExpression(op, left.(*object.Integer), right.(*object.Integer))
	case *object.Boolean:
		return evalInfixBooleanExpression(op, left.(*object.Boolean), right.(*object.Boolean))
	default:
		return NULL
	}
}

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		right := Eval(node.Right)
		left := Eval(node.Left)
		return evalInfixExpression(node.Operator, left, right)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	}

	return nil
}
