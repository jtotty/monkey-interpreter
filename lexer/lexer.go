package lexer

import "github.com/jtotty/monkey-interpreter/token"

type Lexer struct {
	input        string
	position     int  // current pos in input (current char)
	readPosition int  // current reading pos in input (after current char)
	ch           byte // current char
}

func New(input string) *Lexer {
	lexer := &Lexer{input: input}
	lexer.readChar()
	return lexer
}

// Give the next character and advance the position in the input string
func (lexer *Lexer) readChar() {
	if lexer.readPosition >= len(lexer.input) {
		lexer.ch = 0 // ASCII code for "NUL", i.e. EOF
	} else {
		lexer.ch = lexer.input[lexer.readPosition]
	}

	lexer.position = lexer.readPosition
	lexer.readPosition += 1
}

func (lexer *Lexer) NextToken() token.Token {
	var tok token.Token

	switch lexer.ch {
	// Operators
	case '=':
		tok = newToken(token.ASSIGN, lexer.ch)
	case '+':
		tok = newToken(token.PLUS, lexer.ch)

	// Delimiters
	case ';':
		tok = newToken(token.SEMICOLON, lexer.ch)
	case '(':
		tok = newToken(token.LPAREN, lexer.ch)
	case ')':
		tok = newToken(token.RPAREN, lexer.ch)
	case ',':
		tok = newToken(token.COMMA, lexer.ch)
	case '{':
		tok = newToken(token.LBRACE, lexer.ch)
	case '}':
		tok = newToken(token.RBRACE, lexer.ch)

	// NUL or EOF
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF

	// Identifiers, literals, or illegal
	default:
		if isLetter(lexer.ch) {
			tok.Literal = lexer.readIndentifier()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, lexer.ch)
		}
	}

	lexer.readChar()
	return tok
}

// Reads in an indentifier and advances the lexer's position
// until it encounters a non-letter-character
func (lexer *Lexer) readIndentifier() string {
	position := lexer.position
	for isLetter(lexer.ch) {
		lexer.readChar()
	}
	return lexer.input[position:lexer.position]
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func isLetter(ch byte) bool {
	// Underscore allows snake case variable naming
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}
