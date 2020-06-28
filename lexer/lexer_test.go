package lexer

import (
	"testing"

	"goto/token"
)

func TestNextToken(t *testing.T) {
	input := `var a = b % 6;
				if 4>=10 {
					return a-b;
				}
				"foobar"
				"foo bar"
				a,b = 4,5;
				for a=3; a>5 && a<4 ; a=a**1 {
					continue;
					break;
				}
				`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.VAR, "var"},
		{token.IDENT, "a"},
		{token.ASSIGN, "="},
		{token.IDENT, "b"},
		{token.MOD, "%"},
		{token.INT, "6"},
		{token.SEMI, ";"},
		{token.IF, "if"},
		{token.INT, "4"},
		{token.GT_EQ, ">="},
		{token.INT, "10"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IDENT, "a"},
		{token.MINUS, "-"},
		{token.IDENT, "b"},
		{token.SEMI, ";"},
		{token.RBRACE, "}"},
		{token.STRING, "foobar"},
		{token.STRING, "foo bar"},
		{token.IDENT, "a"},
		{token.COMMA, ","},
		{token.IDENT, "b"},
		{token.ASSIGN, "="},
		{token.INT, "4"},
		{token.COMMA, ","},
		{token.INT, "5"},
		{token.SEMI, ";"},
		{token.FOR, "for"},
		{token.IDENT, "a"},
		{token.ASSIGN, "="},
		{token.INT, "3"},
		{token.SEMI, ";"},
		{token.IDENT, "a"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.AND, "&&"},
		{token.IDENT, "a"},
		{token.LT, "<"},
		{token.INT, "4"},
		{token.SEMI, ";"},
		{token.IDENT, "a"},
		{token.ASSIGN, "="},
		{token.IDENT, "a"},
		{token.POW, "**"},
		{token.INT, "1"},
		{token.LBRACE, "{"},
		{token.CONTINUE, "continue"},
		{token.SEMI, ";"},
		{token.BREAK, "break"},
		{token.SEMI, ";"},
		{token.RBRACE, "}"},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - Type wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}

	}

}
