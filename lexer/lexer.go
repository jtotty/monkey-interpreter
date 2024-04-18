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

// Set the next character and advance the position in the input string
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

	lexer.skipWhitespace()

	switch lexer.ch {
	// Operators
	case '=':
		if lexer.peekChar() == '=' {
			tok = lexer.newTwoCharToken(token.EQ)
		} else {
			tok = lexer.newToken(token.ASSIGN)
		}
	case '+':
		tok = lexer.newToken(token.PLUS)
	case '-':
		tok = lexer.newToken(token.MINUS)
	case '!':
		if lexer.peekChar() == '=' {
			tok = lexer.newTwoCharToken(token.NOT_EQ)
		} else {
			tok = lexer.newToken(token.BANG)
		}
	case '*':
		tok = lexer.newToken(token.ASTERISK)
	case '/':
		tok = lexer.newToken(token.SLASH)
	case '<':
		tok = lexer.newToken(token.LT)
	case '>':
		tok = lexer.newToken(token.GT)

	// Delimiters
	case ',':
		tok = lexer.newToken(token.COMMA)
	case ';':
		tok = lexer.newToken(token.SEMICOLON)

	// Braces
	case '(':
		tok = lexer.newToken(token.LPAREN)
	case ')':
		tok = lexer.newToken(token.RPAREN)
	case '{':
		tok = lexer.newToken(token.LBRACE)
	case '}':
		tok = lexer.newToken(token.RBRACE)

	// NUL or EOF
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF

	// Identifiers, literals, or illegal
	default:
		if isLetter(lexer.ch) {
			tok.Literal = lexer.readLiteral(isLetter)
			tok.Type = token.LookupIndentifier(tok.Literal)
			return tok
		} else if isDigit(lexer.ch) {
			tok.Literal = lexer.readLiteral(isDigit)
			tok.Type = token.INT
			return tok
		} else {
			tok = lexer.newToken(token.ILLEGAL)
		}
	}

	lexer.readChar()
	return tok
}

func (lexer *Lexer) skipWhitespace() {
	for lexer.ch == ' ' || lexer.ch == '\t' || lexer.ch == '\n' || lexer.ch == '\r' {
		lexer.readChar()
	}
}

func (lexer *Lexer) newToken(tokenType token.TokenType) token.Token {
	return token.Token{Type: tokenType, Literal: string(lexer.ch)}
}

func (lexer *Lexer) newTwoCharToken(tokenType token.TokenType) token.Token {
	ch := lexer.ch
	lexer.readChar()
	literal := string(ch) + string(lexer.ch)
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
func (lexer *Lexer) readLiteral(litmusFunc litmus) string {
	startPostion := lexer.position
	for litmusFunc(lexer.ch) {
		lexer.readChar()
	}
	return lexer.input[startPostion:lexer.position]
}

// Peek ahead in the input without progressing the position
func (lexer *Lexer) peekChar() byte {
	if lexer.readPosition >= len(lexer.input) {
		return 0
	} else {
		return lexer.input[lexer.readPosition]
	}
}
