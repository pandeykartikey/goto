package token

type TokenType string

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
	SEMI = ";"

	// Keywords
	VAR = "VAR" 
)

var keywords = map[string]TokenType {
	"var" : VAR,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENT
}