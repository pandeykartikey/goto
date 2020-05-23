package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers
	IDENT = "IDENT"

	// Literals
	INT = "INT"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	MULTIPLY = "*"
	DIVIDE   = "/"
	EQ       = "=="
	NOT      = "!"
	NOT_EQ   = "!="
	LT       = "<"
	LT_EQ    = "<="
	GT       = ">"
	GT_EQ    = ">="

	// Delimiters
	SEMI  = ";"
	COLON = ":"

	// Keywords
	VAR    = "VAR"
	FUNC   = "FUNC"
	TRUE   = "TRUE"
	FALSE  = "FALSE"
	IF     = "IF"
	ELSE   = "ELSE"
	RETURN = "RETURN"
)

var keywords = map[string]TokenType{
	"var":    VAR,
	"func":   FUNC,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENT
}
