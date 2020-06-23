package eval

import (
	"fmt"
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

func isError(obj object.Object) bool {
	if _, ok := obj.(*object.Error); ok {
		return true
	}

	return false
}

func errorMessageToObject(msg string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(msg, a...)}
}

// to evaluate a block of statements, nested to marked true when inside a nested block
func evalStatements(stmts []ast.Statement, env *object.Environment, nested bool) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt, env)

		switch result.(type) {
		case *object.ReturnValue:
			if !nested {
				return result.(*object.ReturnValue).Value
			}
			return result
		case *object.Error:
			return result
		}
	}

	return result
}

func evalIdentifier(id *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(id.Value)

	if !ok {
		return errorMessageToObject("Identifier not found: %s", id.Value)
	}
	return val
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
		return errorMessageToObject("Unknown Operator: -%s", obj.Type())
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
		return errorMessageToObject("Unknown Operator: %s %s", op, right.Type())
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
		return errorMessageToObject("Unknown Operator: %s %s %s", left.Type(), op, right.Type())
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
		return errorMessageToObject("Unknown Operator: %s %s %s", left.Type(), op, right.Type())
	}
}

func evalInfixExpression(op string, left object.Object, right object.Object) object.Object {
	if left.Type() != right.Type() {
		return errorMessageToObject("Type Mismatch: %s %s %s", left.Type(), op, right.Type())
	}

	switch left.(type) {
	case *object.Integer:
		return evalInfixIntegerExpression(op, left.(*object.Integer), right.(*object.Integer))
	case *object.Boolean:
		return evalInfixBooleanExpression(op, left.(*object.Boolean), right.(*object.Boolean))
	default:
		return errorMessageToObject("Unknown Type %s %s", left.Type(), right.Type())
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

func evalIfStatement(ifStmt *ast.IfStatement, env *object.Environment) object.Object {
	cond := Eval(ifStmt.Condition, env)

	if isTrue(cond) {
		return Eval(ifStmt.Consequence, env)
	} else if ifStmt.Alternative != nil {
		return Eval(ifStmt.Alternative, env)
	} else if ifStmt.FollowIf != nil {
		return Eval(ifStmt.FollowIf, env)
	}

	return NULL
}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements, env, false)
	case *ast.IfStatement:
		return evalIfStatement(node, env)
	case *ast.BlockStatement:
		return evalStatements(node.Statements, env, true)
	case *ast.ReturnStatement:
		returnVal := Eval(node.ReturnValue, env)
		if isError(returnVal) {
			return returnVal
		}
		return &object.ReturnValue{Value: returnVal}
	case *ast.VarStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	}

	return nil
}
