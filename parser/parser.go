package parser

import (
	"github.com/jtotty/monkey-interpreter/ast"
	"github.com/jtotty/monkey-interpreter/lexer"
	"github.com/jtotty/monkey-interpreter/token"
)

type Parser struct {
	lexer *lexer.Lexer

	currToken token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	parser := &Parser{lexer: l}

	// Read two tokens, so that currToken and peekToken are both set
	parser.nextToken()
	parser.nextToken()

	return parser
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	return nil
}
