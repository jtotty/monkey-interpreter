package lexer

import (
	"github.com/jtotty/monkey-interpreter/token"
)

type Lexer struct {
	input        string
	position     int  // current pos in input (current char)
	readPosition int  // current reading pos in input (after current char)
	ch           byte // current char
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	// Operators
	case '=':
		if l.peekChar() == '=' {
			tok = l.newTwoCharToken(token.EQ)
		} else {
			tok = l.newToken(token.ASSIGN)
		}
	case '+':
		tok = l.newToken(token.PLUS)
	case '-':
		tok = l.newToken(token.MINUS)
	case '!':
		if l.peekChar() == '=' {
			tok = l.newTwoCharToken(token.NOT_EQ)
		} else {
			tok = l.newToken(token.BANG)
		}
	case '*':
		tok = l.newToken(token.ASTERISK)
	case '/':
		tok = l.newToken(token.SLASH)
	case '<':
		tok = l.newToken(token.LT)
	case '>':
		tok = l.newToken(token.GT)

	// Delimiters
	case ',':
		tok = l.newToken(token.COMMA)
	case ';':
		tok = l.newToken(token.SEMICOLON)

	// Braces
	case '(':
		tok = l.newToken(token.LPAREN)
	case ')':
		tok = l.newToken(token.RPAREN)
	case '{':
		tok = l.newToken(token.LBRACE)
	case '}':
		tok = l.newToken(token.RBRACE)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()

	// NUL or EOF
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF

	// Identifiers, literals, or illegal
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readLiteral(isLetter)
			tok.Type = token.LookupIndentifier(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readLiteral(isDigit)
			tok.Type = token.INT
			return tok
		} else {
			tok = l.newToken(token.ILLEGAL)
		}
	}

	l.readChar()
	return tok
}

// First call sets the starting position (0) and readPosition (1).
// Subsequent calls set the next character and advance the position in the input string.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII code for "NUL", i.e. EOF
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition += 1
}

// Peek ahead in the input without progressing the position
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

// Monkey Programming is whitespace agnostic
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) newToken(tokenType token.TokenType) token.Token {
	return token.Token{Type: tokenType, Literal: string(l.ch)}
}

func (l *Lexer) newTwoCharToken(tokenType token.TokenType) token.Token {
	ch := l.ch
	l.readChar()
	literal := string(ch) + string(l.ch)
	return token.Token{Type: tokenType, Literal: literal}
}

type litmus func(ch byte) bool

func isLetter(ch byte) bool {
	// Underscore allows snake case variable naming
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// Takes a litmus function used on the current char and
// advances the lexer's position until the litmus returns false.
// E.g. test if a char is a number or if it's a letter
func (l *Lexer) readLiteral(litmusFnc litmus) string {
	startPostion := l.position
	for litmusFnc(l.ch) {
		l.readChar()
	}
	return l.input[startPostion:l.position]
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}

	return l.input[position:l.position]
}
