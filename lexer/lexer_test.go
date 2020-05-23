package lexer

import (
	"testing"

	"pyro/token"
)

func TestNextToken(t *testing.T) {
	input := `var a = b + 6`;

	tests := []struct {
		expectedType token.TokenType
		expectedLiteral string
	}{
		{token.VAR, "var"},
		{token.IDENT, "a"},
		{token.ASSIGN, "="},
		{token.IDENT, "b"},
		{token.PLUS, "+"},
		{token.INT, "6"},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",i, tt.expectedLiteral, tok.Literal)
		}

	}

}
