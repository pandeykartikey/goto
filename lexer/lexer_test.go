package lexer

import (
	"testing"

	"goto/token"
)

func TestNextToken(t *testing.T) {
	input := `var a = b + 6;
if 4>=10:
	return a-b;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.VAR, "var"},
		{token.IDENT, "a"},
		{token.ASSIGN, "="},
		{token.IDENT, "b"},
		{token.PLUS, "+"},
		{token.INT, "6"},
		{token.SEMI, ";"},
		{token.IF, "if"},
		{token.INT, "4"},
		{token.GT_EQ, ">="},
		{token.INT, "10"},
		{token.COLON, ":"},
		{token.RETURN, "return"},
		{token.IDENT, "a"},
		{token.MINUS, "-"},
		{token.IDENT, "b"},
		{token.SEMI, ";"},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}

	}

}
