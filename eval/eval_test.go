package eval

import (
	"testing"

	"goto/lexer"
	"goto/object"
	"goto/parser"
)

func evalInput(inp string) object.Object {
	l := lexer.New(inp)
	p := parser.New(l)

	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, exp int64) bool {
	intobj, ok := obj.(*object.Integer)

	if !ok {
		t.Errorf("Expected object type to be integer. got=%T", obj)
		return false
	}

	if intobj.Value != exp {
		t.Errorf("Expected %d but got %d instead", exp, intobj.Value)
		return false
	}

	return true
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input string
		exp   int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 % 2 + 10", 11},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		out := evalInput(tt.input)
		testIntegerObject(t, out, tt.exp)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, exp bool) bool {
	boolobj, ok := obj.(*object.Boolean)

	if !ok {
		t.Errorf("Expected object type to be boolean. got=%T", obj)
		return false
	}

	if boolobj.Value != exp {
		t.Errorf("Expected %v but got %v instead", exp, boolobj.Value)
		return false
	}

	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input string
		exp   bool
	}{
		{"true", true},
		{"false", false},
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{`"Hello" != "World"`, true},
		{`"Hello" == "World"`, false},
	}

	for _, tt := range tests {
		out := evalInput(tt.input)
		testBooleanObject(t, out, tt.exp)
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("Expected object to be NULL. got=%T", obj)
		return false
	}
	return true
}

func testStringObject(t *testing.T, obj object.Object, exp string) bool {
	strObj, ok := obj.(*object.String)

	if !ok {
		t.Errorf("Expected object type to be String. got=%T", obj)
		return false
	}

	if strObj.Value != exp {
		t.Errorf("Expected %v but got %v instead", exp, strObj.Value)
		return false
	}

	return true
}

func TestEvalStringExpression(t *testing.T) {
	tests := []struct {
		input string
		exp   string
	}{
		{`"Hello" + " " + "World"`, "Hello World"},
		{`var a = "Hello"; a + " World"`, "Hello World"},
	}

	for _, tt := range tests {
		out := evalInput(tt.input)
		testStringObject(t, out, tt.exp)
	}
}
func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input string
		exp   interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 < 2) { if ( 1 < 2 ) { 10; } else { 11;} }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 > 2) { 10 } else if ( 3 > 4 ) { 20 } else { 30 }", 30},
		{"if (1 > 2) { 10 } else if ( 3 < 4 ) { 20 } else { 30 }", 20},
	}

	for _, tt := range tests {
		out := evalInput(tt.input)
		intg, ok := tt.exp.(int)

		if ok {
			testIntegerObject(t, out, int64(intg))
		} else {
			testNullObject(t, out)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input string
		exp   int64
	}{
		// {"func asd() { return 10; }; asd();", 10},
		{
			` func asd () { 
				if (10>1) {
					if (10>1) {
						return 10;
					}
					return 1;
				}
			}
			asd();`,
			10,
		},
	}
	for _, tt := range tests {
		out := evalInput(tt.input)
		testIntegerObject(t, out, tt.exp)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input string
		exp   string
	}{
		{
			"5 + true;",
			"Type Mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"Type Mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"Unknown Operator: -BOOLEAN",
		},
		{
			"true + false;",
			"Unknown Operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"Unknown Operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"Unknown Operator: BOOLEAN + BOOLEAN",
		},
		{
			`if (10>1) {
						if (10>1) {
							return true + false;
							}
						return 1;
					}`,
			"Unknown Operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"Identifier not found: foobar",
		},
		{
			`"Hello" - "World"`,
			"Unknown Operator: STRING - STRING",
		},
	}

	for _, tt := range tests {
		out := evalInput(tt.input)
		errObj, ok := out.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T", out)
			continue
		}
		if errObj.Message != tt.exp {
			t.Errorf("wrong error message. expected=%q, got=%q", tt.exp, errObj.Message)
		}
	}
}

func TestAssigmentStatements(t *testing.T) {
	tests := []struct {
		input string
		exp   int64
	}{
		{"var a = 5; a;", 5},
		{"var a = 5 * 5; a;", 25},
		{"var a = 5; var b = a; b;", 5},
		{"var a = 5; var b = a; var c = a + b + 5; c;", 15},
		{"var a,b = 4,5; a+b;", 9},
		{"var a,b = 4,5; a = 6; a+b;", 11},
		{"var a,b = 4,5; a,b = 5,6; a+b;", 11},
	}

	for _, tt := range tests {
		out := evalInput(tt.input)
		testIntegerObject(t, out, tt.exp)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "func asd (x) { x + 2; }; asd;"
	out := evalInput(input)

	fn, ok := out.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T", out)
	}

	if len(fn.ParameterList.Identifiers) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fn.ParameterList.Identifiers)
	}
	if fn.ParameterList.Identifiers[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.ParameterList.Identifiers[0].String())
	}

	if fn.FuncBody.String() != "{ (x + 2) }" {
		t.Fatalf("body is not (x + 2). got=%q", fn.FuncBody.String())
	}
}

func TestFunctionCall(t *testing.T) {
	tests := []struct {
		input string
		exp   int64
	}{
		{
			"func identity(x) { x; }; identity(5);",
			5,
		},
		{
			"func identity() { return 5; }; identity();",
			5,
		},
		{
			"func double(x) { x * 2; }; double(5);",
			10,
		},
		{
			"func add(x, y) { x + y; }; add(5, 5);",
			10,
		},
		{
			"func add(x, y) { x + y; }; add(5 + 5, add(5, 5));",
			20,
		},
	}
	for _, tt := range tests {
		out := evalInput(tt.input)
		testIntegerObject(t, out, tt.exp)
	}
}

func TestForStatement(t *testing.T) {
	tests := []struct {
		input string
		exp   int64
	}{
		{`for var a = 0; a<10; a=a+1 {
				a = a+2;
				if a>5 {
					break;
				}
			}
			a;`, 8,
		},
		{`for var a = 0; a<10; a=a+1 {
			if a>5 {
				continue;
			}
			a = a+5;
		}
		a`, 10,
		},
	}

	for _, tt := range tests {
		out := evalInput(tt.input)
		testIntegerObject(t, out, tt.exp)
	}
}
