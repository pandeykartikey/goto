package parser

import (
	"fmt"
	"testing"

	"pyro/ast"
	"pyro/lexer"
)

func parseInput(t *testing.T, input string, n int) *ast.Program {
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("Parse Program returned nil")
	}

	if len(program.Statements) != n {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", n, len(program.Statements))
	}

	return program
}

func assertExpressionStatement(t *testing.T, program *ast.Program) ast.Expression {

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	return stmt.Expression
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))

	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}

func testVarStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "var" {
		t.Errorf("s.TokenLiteral not 'var'. got=%q", s.TokenLiteral())
		return false
	}

	varStmt, ok := s.(*ast.VarStatement)

	if !ok {
		t.Errorf("s not *ast.VarStatement. got=%T", s)
		return false
	}
	if varStmt.Name.Value != name {
		t.Errorf("VarStmt.Name.Value not '%s'. got=%s", name, varStmt.Name.Value)
		return false
	}
	if varStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, varStmt.Name)
		return false
	}
	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)

	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", exp)
		return false
	}
	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}
	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value, ident.TokenLiteral())
		return false
	}

	return true
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
		t.Errorf("integ.TokenLiteral not %d. got=%s", value, integ.TokenLiteral())
		return false
	}
	return true
}

func testBooleanLiteral(t *testing.T, bl ast.Expression, value bool) bool {
	boolean, ok := bl.(*ast.Boolean)

	if !ok {
		t.Errorf("bl not *ast.Boolean. got=%T", bl)
		return false
	}

	if boolean.Value != value {
		t.Errorf("boolean.Value not %v. got=%v", value, boolean.Value)
		return false
	}
	if boolean.TokenLiteral() != fmt.Sprintf("%v", value) {
		t.Errorf("boolean.TokenLiteral not %v. got=%v", value, boolean.TokenLiteral())
		return false
	}
	return true
}

func testIdentifierList(t *testing.T, il ast.Expression, value []string) bool {
	identlist, ok := il.(*ast.IdentifierList)

	if !ok {
		t.Errorf("il is not *ast.IdentifierList. got=%T", il)
		return false
	}

	for idx, ident := range identlist.Identifiers {
		if ident.Value != value[idx] {
			t.Errorf("ident.Value is not %s. got=%s", value[idx], ident.Value)
			return false
		}
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)

	return false
}

func TestVarStatements(t *testing.T) {
	input := `var a = b + 6;
var bar = 5;
var foo = 2;`

	program := parseInput(t, input, 3)

	tests := []struct {
		expectedIdentifier string
	}{
		{"a"},
		{"bar"},
		{"foo"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testVarStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	input := `return 5;
return 10;
return 993322;`

	program := parseInput(t, input, 3)

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.returnStatement. got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q", returnStmt.TokenLiteral())
		}
	}
}

func TestIdentifierStatement(t *testing.T) {
	input := `foobar;`

	program := parseInput(t, input, 1)
	expstmt := assertExpressionStatement(t, program)

	if testIdentifier(t, expstmt, "foobar") {
		return
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := `5;`

	program := parseInput(t, input, 1)
	expstmt := assertExpressionStatement(t, program)

	if !testLiteralExpression(t, expstmt, 5) {
		return
	}

}

func TestParsingPrefixExpression(t *testing.T) {
	input := []struct {
		input        string
		operator     string
		integerValue interface{}
	}{
		{"!5;", "!", 5},
		{"-1;", "-", 1},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range input {
		program := parseInput(t, tt.input, 1)
		expstmt := assertExpressionStatement(t, program)

		exp, ok := expstmt.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", expstmt)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.integerValue) {
			return
		}
	}
}

func testInfixExpression(t *testing.T, input ast.Expression, leftValue interface{}, operator string, rightValue interface{}) bool {
	exp, ok := input.(*ast.InfixExpression)

	if !ok {
		t.Fatalf("exp is not ast.InfixExpression. got=%T", input)
		return false
	}
	if !testLiteralExpression(t, exp.Left, leftValue) {
		return false
	}
	if exp.Operator != operator {
		t.Fatalf("exp.Operator is not '%s'. got=%s", operator, exp.Operator)
		return false
	}
	if !testLiteralExpression(t, exp.Right, rightValue) {
		return false
	}

	return true

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
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}
	for _, tt := range infixTests {
		program := parseInput(t, tt.input, 1)

		expstmt := assertExpressionStatement(t, program)

		if !testInfixExpression(t, expstmt, tt.leftValue, tt.operator, tt.rightValue) {
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
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
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
		program := parseInput(t, tt.input, 1)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	input := `false;`

	program := parseInput(t, input, 1)
	expstmt := assertExpressionStatement(t, program)

	if testLiteralExpression(t, expstmt, false) {
		return
	}
}

func TestIfStatements(t *testing.T) {
	input := `if a==b {
var a = 6;
} else if b==c {
	var b = 5;
} else {
	var c = 10;
}`

	program := parseInput(t, input, 1)

	stmt, ok := program.Statements[0].(*ast.IfStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.IfStatement. got=%T", program.Statements[0])
	}

	if !testInfixExpression(t, stmt.Condition, "a", "==", "b") {
		return
	}

	if len(stmt.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n", len(stmt.Consequence.Statements))
	}

	if !testVarStatement(t, stmt.Consequence.Statements[0], "a") {
		return
	}

	if stmt.Alternative != nil {
		t.Errorf("stmt.Alternative.Statements was not nil. got=%+v", stmt.Alternative)
	}

	if !testInfixExpression(t, stmt.FollowIf.Condition, "b", "==", "c") {
		return
	}

	if len(stmt.FollowIf.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n", len(stmt.FollowIf.Consequence.Statements))
	}

	if !testVarStatement(t, stmt.FollowIf.Consequence.Statements[0], "b") {
		return
	}

	if stmt.FollowIf.FollowIf != nil {
		return
	}

	if len(stmt.FollowIf.Alternative.Statements) != 1 {
		t.Errorf("Alternative is not 1 statements. got=%d\n", len(stmt.FollowIf.Alternative.Statements))
	}

	if !testVarStatement(t, stmt.FollowIf.Alternative.Statements[0], "c") {
		return
	}
}

func TestFuncStatement(t *testing.T) {
	input := `func abc (x, y) {
		return x;
	}`

	program := parseInput(t, input, 1)

	stmt, ok := program.Statements[0].(*ast.FuncStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.FuncStatement. got=%T", program.Statements[0])
	}

	if testIdentifier(t, stmt.Name, "abc") {
		return
	}

	if testIdentifierList(t, stmt.ParameterList, []string{"x", "y"}) {
		return
	}

	if len(stmt.FuncBody.Statements) != 1 {
		t.Errorf("FuncBody is not 1 statements. got=%d\n", len(stmt.FuncBody.Statements))
	}

	returnStmt, ok := stmt.FuncBody.Statements[0].(*ast.ReturnStatement)
	if !ok {
		t.Errorf("stmt.FuncBody.Statements[0] not *ast.returnStatement. got=%T", stmt)
	}

	if returnStmt.TokenLiteral() != "return" {
		t.Errorf("returnStmt.TokenLiteral not 'return', got %q", returnStmt.TokenLiteral())
	}

	if returnStmt.ReturnValue != nil {
		t.Errorf("returnStmt.ReturnValue is nil")
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2*3, 4+5)"

	program := parseInput(t, input, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)

	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, exp.FunctionName, "add") {
		return
	}

	args := exp.ArgumentList.Expressions

	if len(args) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(args))
	}

	testLiteralExpression(t, *args[0], 1)
	testInfixExpression(t, *args[1], 2, "*", 3)
	testInfixExpression(t, *args[2], 4, "+", 5)

}
