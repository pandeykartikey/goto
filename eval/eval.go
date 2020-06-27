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
func evalStatements(stmts []ast.Statement, env *object.Environment, insideFunc bool) object.Object {
	var result object.Object

	for _, stmt := range stmts {

		result = evalProgram(stmt, env)

		switch result.(type) {
		case *object.ReturnValue:
			if insideFunc {
				return result.(*object.ReturnValue).Value
			}
			return result
		case *object.Error, *object.LoopControl:
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

func evalNotOperator(obj object.Object) object.Object {
	return nativeBoolToBooleanObject(!isTrue(obj))
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

// NULL is false and all other values are true
func isTrue(obj object.Object) bool {
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

	if exprList == nil {
		return &object.List{Value: objList}
	}

	for _, expr := range exprList.Expressions {
		obj := evalProgram(*expr, env)
		if isError(obj) {
			return obj
		}
		objList = append(objList, &obj)
	}

	return &object.List{Value: objList}
}

func evalAssignment(assignStmt *ast.Assignment, env *object.Environment) object.Object {

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
	cond := evalProgram(ifStmt.Condition, env)

	if isError(cond) {
		return cond
	}

	if isTrue(cond) {
		return evalProgram(ifStmt.Consequence, env)
	} else if ifStmt.Alternative != nil {
		return evalProgram(ifStmt.Alternative, env)
	} else if ifStmt.FollowIf != nil {
		return evalProgram(ifStmt.FollowIf, env)
	}

	return NULL
}

func evalForStatement(forStmt *ast.ForStatement, env *object.Environment) object.Object {
	out := evalAssignment(forStmt.Init, env)
	if isError(out) {
		return out
	}
forLoop:
	for {

		cond := evalProgram(forStmt.Condition, env)
		if isError(cond) {
			return cond
		}
		if !isTrue(cond) {
			break
		}

		out = evalStatements(forStmt.ForBody.Statements, env, false)

		switch out.(type) {
		case *object.LoopControl:
			if out.Inspect() == "break" {
				break forLoop
			}
		case *object.Error, *object.ReturnValue:
			return out
		}

		out = evalProgram(forStmt.Update, env)
		if isError(out) {
			return out
		}
	}

	return nil
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

	if fn.ParameterList == nil {
		return extendedEnv
	}

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

	return evalStatements(fnObj.FuncBody.Statements, extendedEnv, true)
}

func evalProgram(node ast.Node, env *object.Environment) object.Object {
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
	case *ast.ForStatement:
		return evalForStatement(node, env)
	case *ast.IfStatement:
		return evalIfStatement(node, env)
	case *ast.BlockStatement:
		return evalStatements(node.Statements, env, false)
	case *ast.ReturnStatement:
		returnVal := evalProgram(node.ReturnValue, env)
		if isError(returnVal) {
			return returnVal
		}
		return &object.ReturnValue{Value: returnVal}
	case *ast.LoopControlStatement:
		return &object.LoopControl{Value: node.TokenLiteral()}
	case *ast.Assignment:
		return evalAssignment(node, env)
	case *ast.ExpressionStatement:
		return evalProgram(node.Expression, env)
	case *ast.PrefixExpression:
		right := evalProgram(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		right := evalProgram(node.Right, env)
		if isError(right) {
			return right
		}
		left := evalProgram(node.Left, env)
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

func Eval(node ast.Node, env *object.Environment) object.Object {
	out := evalProgram(node, env)
	switch out.(type) {
	case *object.ReturnValue:
		return errorMessageToObject("return used outside function")
	case *object.LoopControl:
		return errorMessageToObject("break or continue used outside for loop")
	default:
		return out
	}
}
