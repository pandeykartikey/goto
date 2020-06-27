package eval

import (
	"fmt"

	"goto/ast"
	"goto/object"
)

var (
	NULL        = &object.Null{}
	DEFAULT_INT = &object.Integer{Value: 0}
	TRUE        = &object.Boolean{Value: true}
	FALSE       = &object.Boolean{Value: false}
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
	case "%":
		return &object.Integer{Value: leftVal % rightVal}
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

func evalInfixStringExpression(op string, left *object.String, right *object.String) object.Object {
	leftVal, rightVal := left.Value, right.Value

	switch op {
	case "+":
		return &object.String{Value: leftVal + rightVal}
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
	case *object.String:
		return evalInfixStringExpression(op, left.(*object.String), right.(*object.String))
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

func evalExpressionList(exprList *ast.ExpressionList, env *object.Environment) object.Object {
	var objList []*object.Object
	for _, expr := range exprList.Expressions {
		obj := Eval(*expr, env)
		if isError(obj) {
			return obj
		}
		objList = append(objList, &obj)
	}
	return &object.List{Value: objList}
}

func evalAssignStatement(assignStmt *ast.AssignStatement, env *object.Environment) object.Object {

	var (
		valueList *object.List
		ok        bool
	)

	if assignStmt.ValueList != nil {
		evaluatedList := evalExpressionList(assignStmt.ValueList, env)
		if isError(evaluatedList) {
			return evaluatedList
		}
		valueList, ok = evaluatedList.(*object.List)
		if !ok {
			return nil
		}
	}

	for idx, ident := range assignStmt.NameList.Identifiers {
		switch assignStmt.TokenLiteral() {
		case "var":
			if valueList != nil {
				if _, ok = env.Create(ident.Value, *valueList.Value[idx]); !ok {
					return errorMessageToObject("An identifier already exists with that name")
				}
			} else {
				if _, ok = env.Create(ident.Value, DEFAULT_INT); !ok {
					return errorMessageToObject("An identifier already exists with that name")
				}
			}
		case "=":
			if _, ok = env.Update(ident.Value, *valueList.Value[idx]); !ok {
				return errorMessageToObject("An identifier does not exists with that name")
			}
		default:
			return errorMessageToObject("Unexpected Error encountered")
		}
	}

	return nil
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

func evalFuncStatement(funcStmt *ast.FuncStatement, env *object.Environment) object.Object {
	funcObj := &object.Function{
		ParameterList: funcStmt.ParameterList,
		FuncBody:      funcStmt.FuncBody,
	}

	if _, ok := env.Create(funcStmt.Name.Value, funcObj); !ok {
		return errorMessageToObject("A function already exists with that name")
	}

	return nil
}

func addArgumentsToEnvironment(fn *object.Function, objList *object.List, env *object.Environment) *object.Environment {
	extendedEnv := object.ExtendEnv(env)

	for idx, param := range fn.ParameterList.Identifiers {
		extendedEnv.Create(param.Value, *objList.Value[idx])
	}

	return extendedEnv
}

func evalCallExpression(name string, obj object.Object, env *object.Environment) object.Object {
	args, ok := obj.(*object.List)
	if !ok {
		return errorMessageToObject("Unknown Operator: %s", obj.Type())
	}
	fn, ok := env.Get(name)
	if !ok {
		return errorMessageToObject("Function not found: %s", name)
	}
	fnObj, ok := fn.(*object.Function)
	if !ok {
		return errorMessageToObject("Function not found: %s", name)
	}
	extendedEnv := addArgumentsToEnvironment(fnObj, args, env)
	return evalStatements(fnObj.FuncBody.Statements, extendedEnv, false)
}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements, env, false)
	case *ast.FuncStatement:
		return evalFuncStatement(node, env)
	case *ast.CallExpression:
		args := evalExpressionList(node.ArgumentList, env)
		if isError(args) {
			return args
		}
		return evalCallExpression(node.FunctionName.Value, args, env)
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
	case *ast.AssignStatement:
		return evalAssignStatement(node, env)
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
	case *ast.String:
		return &object.String{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	}

	return nil
}
