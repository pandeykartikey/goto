package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

type CommonPrefixTokenPair struct {
	NextCharacter         byte
	SingleCharacterType   TokenType
	MultipleCharacterType TokenType
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
	COMMA = ","

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	VAR    = "VAR"
	FUNC   = "FUNC"
	TRUE   = "TRUE"
	FALSE  = "FALSE"
	IF     = "IF"
	ELSE   = "ELSE"
	RETURN = "RETURN"
)

var Keywords = map[string]TokenType{
	"var":    VAR,
	"func":   FUNC,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

var SingleCharacterToken = map[byte]TokenType{
	'+': PLUS,
	'-': MINUS,
	'*': MULTIPLY,
	'/': DIVIDE,
	';': SEMI,
	':': COLON,
	',': COMMA,
	'(': LPAREN,
	')': RPAREN,
	'{': LBRACE,
	'}': RBRACE,
}

var CommonPrefixToken = map[byte]CommonPrefixTokenPair{
	'=': {NextCharacter: '=', MultipleCharacterType: EQ, SingleCharacterType: ASSIGN},
	'!': {NextCharacter: '=', MultipleCharacterType: NOT_EQ, SingleCharacterType: NOT},
	'<': {NextCharacter: '=', MultipleCharacterType: LT_EQ, SingleCharacterType: LT},
	'>': {NextCharacter: '=', MultipleCharacterType: GT_EQ, SingleCharacterType: GT},
}

func LookupGroup(s string, m map[string]TokenType, def TokenType) TokenType { // def default token type
	if tok, ok := m[s]; ok {
		return tok
	}

	return def
}
