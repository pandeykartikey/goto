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

// to evaluate a block of statements, nested to marked true when inside a nested block
func evalStatements(stmts []ast.Statement, nested bool) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt)

		if retVal, ok := result.(*object.ReturnValue); ok {
			if !nested {
				return retVal.Value
			}
			return retVal
		}
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

func isTrue(obj object.Object) bool { // TODO: merge with eval not operator
	switch obj {
	case TRUE:
		return true
	case FALSE:
		return false
	case NULL:
		return false
	default:
		return true
	}
}

func evalIfStatement(ifStmt *ast.IfStatement) object.Object {
	cond := Eval(ifStmt.Condition)

	if isTrue(cond) {
		return Eval(ifStmt.Consequence)
	} else if ifStmt.Alternative != nil {
		return Eval(ifStmt.Alternative)
	} else if ifStmt.FollowIf != nil {
		return Eval(ifStmt.FollowIf)
	}

	return NULL
}

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements, false)
	case *ast.IfStatement:
		return evalIfStatement(node)
	case *ast.BlockStatement:
		return evalStatements(node.Statements, true)
	case *ast.ReturnStatement:
		returnValue := Eval(node.ReturnValue)
		return &object.ReturnValue{Value: returnValue}
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
