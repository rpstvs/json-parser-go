package parser

import (
	"github.com/rpstvs/json-parser-go/internal/Token"
	"github.com/rpstvs/json-parser-go/internal/ast"
	"github.com/rpstvs/json-parser-go/internal/lexer"
)

type Parser struct {
	lexer        *lexer.Lexer
	errors       []string
	currentToken Token.Token
	peekToken    Token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer: l,
	}

	p.nextToken()
	p.nextToken()

	return p

}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) currentTokenTypeIs(t Token.Type) bool {
	return p.currentToken.Type == t
}

func (p *Parser) ParseValue() ast.Value {
	switch p.currentToken.Type {
	case Token.LeftBrace:
		return p.parseJSONObject()
	case Token.LeftBracket:
		return p.parseJSONArray()
	default:
		return p.parseLiteral()
	}
}

func (p *Parser) parseJSONObject() ast.Value {

}
