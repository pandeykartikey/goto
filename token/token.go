package token

type tokenType string

type Token struct {
	Type 	TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF = "EOF"
	
	// Identifiers
	IDENT = "IDENT"

	// Literals
	INT = "INT"
	
	// Operators
	ASSIGN = "="
	PLUS = "+"
	MINUS = "-"
	MULTIPLY = "*"
	DIVIDE = "/"
	
	// Delimiters


	// Keywords
	VAR = "VAR" 
)