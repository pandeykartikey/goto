package lexer

import (
	"goto/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readSequence(f func(byte) bool) string {
	position := l.position

	for f(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	if l.ch == 0 {
		tok.Type = token.EOF
		tok.Literal = ""
	} else if toktype, ok := token.SingleCharacterToken[l.ch]; ok {
		tok = newToken(toktype, l.ch)
	} else if commonprefixtok, ok := token.CommonPrefixToken[l.ch]; ok {
		if l.peekChar() == commonprefixtok.NextCharacter {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: commonprefixtok.MultipleCharacterType, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(commonprefixtok.SingleCharacterType, l.ch)
		}
	} else if !isNotQuote(l.ch) {
		l.readChar()
		tok.Literal = l.readSequence(isNotQuote)
		l.readChar()
		tok.Type = token.STRING
		return tok
	} else {
		if isLetter(l.ch) {
			tok.Literal = l.readSequence(isLetter)
			tok.Type = token.LookupGroup(tok.Literal, token.Keywords, token.IDENT)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readSequence(isDigit)
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func newToken(Type token.Type, ch byte) token.Token {
	return token.Token{Type: Type, Literal: string(ch)}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isNotQuote(ch byte) bool {
	return ch != '"'
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}
