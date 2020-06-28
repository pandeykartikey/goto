package token

type Type string

type Token struct {
	Type    Type
	Literal string
}

type CommonPrefixTokenPair struct {
	NextCharacter         byte
	SingleCharacterType   Type
	MultipleCharacterType Type
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers
	IDENT = "IDENT"

	// Literals
	INT    = "INT"
	STRING = "STRING"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	MULTIPLY = "*"
	DIVIDE   = "/"
	MOD      = "%"
	POW      = "**"
	EQ       = "=="
	NOT      = "!"
	NOT_EQ   = "!="
	LT       = "<"
	LT_EQ    = "<="
	GT       = ">"
	GT_EQ    = ">="
	AND      = "&&"
	OR       = "||"

	// Delimiters
	SEMI  = ";"
	COLON = ":"
	COMMA = ","
	QUOTE = "\""

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	VAR      = "VAR"
	FUNC     = "FUNC"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	FOR      = "FOR"
	CONTINUE = "CONTINUE"
	BREAK    = "BREAK"
)

var Keywords = map[string]Type{
	"var":      VAR,
	"func":     FUNC,
	"true":     TRUE,
	"false":    FALSE,
	"if":       IF,
	"else":     ELSE,
	"return":   RETURN,
	"for":      FOR,
	"continue": CONTINUE,
	"break":    BREAK,
}

var SingleCharacterToken = map[byte]Type{
	'+': PLUS,
	'-': MINUS,
	'/': DIVIDE,
	'%': MOD,
	';': SEMI,
	':': COLON,
	',': COMMA,
	'(': LPAREN,
	')': RPAREN,
	'{': LBRACE,
	'}': RBRACE,
	'[': LBRACKET,
	']': RBRACKET,
}

var CommonPrefixToken = map[byte]CommonPrefixTokenPair{
	'=': {NextCharacter: '=', MultipleCharacterType: EQ, SingleCharacterType: ASSIGN},
	'*': {NextCharacter: '*', MultipleCharacterType: POW, SingleCharacterType: MULTIPLY},
	'&': {NextCharacter: '&', MultipleCharacterType: AND, SingleCharacterType: AND},
	'|': {NextCharacter: '|', MultipleCharacterType: OR, SingleCharacterType: OR},
	'!': {NextCharacter: '=', MultipleCharacterType: NOT_EQ, SingleCharacterType: NOT},
	'<': {NextCharacter: '=', MultipleCharacterType: LT_EQ, SingleCharacterType: LT},
	'>': {NextCharacter: '=', MultipleCharacterType: GT_EQ, SingleCharacterType: GT},
}

func LookupGroup(s string, m map[string]Type, def Type) Type { // def default token type
	if tok, ok := m[s]; ok {
		return tok
	}

	return def
}
