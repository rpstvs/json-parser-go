package parser

import (
	"fmt"
	"strconv"
	"strings"

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

func (p *Parser) ParseJSON() ast.RootNode {
	var rootNode ast.RootNode
	if p.currentTokenTypeIs(Token.LeftBracket) {
		rootNode.Type = ast.ArrayRoot
	}
	val := p.ParseValue()
	rootNode.RootValue = &val
	return rootNode
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
		return p.parseJSONLiteral()
	}
}

func (p *Parser) parseJSONObject() ast.Value {
	obj := ast.Object{
		Type: "Object",
	}
	objState := ast.ObjStart

	for !p.currentTokenTypeIs(Token.EOF) {
		switch objState {
		case ast.ObjStart:
			if p.currentTokenTypeIs(Token.LeftBrace) {
				objState = ast.ObjOpen
				obj.Start = p.currentToken.Start
				p.nextToken()
			} else {
				p.parseError(fmt.Sprintf("Error parsing JSON Object Expected '{' token, got %s", p.currentToken.Literal))
			}
			return nil
		case ast.ObjOpen:
			if p.currentTokenTypeIs(Token.RightBrace) {
				p.nextToken()
				obj.End = p.currentToken.End
				return obj
			}
			prop := p.parseProperty()
			obj.Children = append(obj.Children, prop)
			objState = ast.ObjProperty
		case ast.ObjProperty:
			if p.currentTokenTypeIs(Token.RightBrace) {
				p.nextToken()
				obj.End = p.currentToken.End
				return obj
			} else if p.currentTokenTypeIs(Token.Comma) {
				objState = ast.ObjComma
				p.nextToken()
			} else {
				p.parseError(fmt.Sprintf("Error parsing JSON Object Expected '}' token, got %s", p.currentToken.Literal))
			}
			return nil
		case ast.ObjComma:
			prop := p.parseProperty()
			if prop.Value != nil {
				obj.Children = append(obj.Children, prop)
				objState = ast.ObjProperty
			}
		}
	}
	obj.End = p.currentToken.Start
	return obj
}

func (p *Parser) parseJSONArray() ast.Value {
	array := ast.Array{
		Type: "Array",
	}
	arrayState := ast.ArrayStart

	switch arrayState {
	case ast.ArrayStart:
		if p.currentTokenTypeIs(Token.LeftBracket) {
			array.Start = p.currentToken.Start
			arrayState = ast.ArrayOpen
		}
		return nil
	case ast.ArrayOpen:
		if p.currentTokenTypeIs(Token.RightBracket) {

			array.End = p.currentToken.End
			p.nextToken()
			return array
		}
		val := p.ParseValue()
		array.Children = append(array.Children, val)
		arrayState = ast.ArrayValue
		if p.peekTokenTypeIs(Token.RightBrace) {
			p.nextToken()
		}
	case ast.ArrayValue:
		if p.currentTokenTypeIs(Token.RightBracket) {
			array.End = p.currentToken.End
			p.nextToken()
			return array
		} else if p.currentTokenTypeIs(Token.Comma) {
			arrayState = ast.ArrayComma
		} else {
			p.parseError(fmt.Sprintf(
				"Error parsing property. Expected RightBrace or Comma token, got: %s",
				p.currentToken.Literal,
			))
		}
	case ast.ArrayComma:
		val := p.ParseValue()
		array.Children = append(array.Children, val)
		arrayState = ast.ArrayComma
	}
	array.End = p.currentToken.Start
	return array
}

func (p *Parser) parseJSONLiteral() ast.Literal {
	val := ast.Literal{
		Type: "Literal",
	}

	defer p.nextToken()

	switch p.currentToken.Type {
	case Token.String:
		val.Value = p.parseString()
		return val
	case Token.Number:
		value, _ := strconv.Atoi(p.currentToken.Literal)
		val.Value = value
		return val
	case Token.False:
		val.Value = false
		return val
	case Token.True:
		val.Value = true
		return val
	default:
		val.Value = "null"
		return val
	}
}

func (p *Parser) parseProperty() ast.Property {
	prop := ast.Property{
		Type: "Property",
	}
	propertyState := ast.PropertyStart

	for !p.currentTokenTypeIs(Token.EOF) {
		switch propertyState {
		case ast.PropertyStart:
			if p.currentTokenTypeIs(Token.String) {
				key := ast.Identifier{
					Type:  "Identifier",
					Value: p.parseString(),
				}
				prop.Key = key
				propertyState = ast.PropertyKey
				p.nextToken()
			} else {
				p.parseError(fmt.Sprintf(
					"Error parsing property start. Expected String token, got: %s",
					p.currentToken.Literal,
				))
			}

		case ast.PropertyKey:
			if p.currentTokenTypeIs(Token.Colon) {
				propertyState = ast.PropertyColon
				p.nextToken()
			} else {
				p.parseError(fmt.Sprintf(
					"Error parsing property. Expected Colon token, got: %s",
					p.currentToken.Literal,
				))
			}
		case ast.PropertyColon:
			val := p.ParseValue()
			prop.Value = val
			p.nextToken()
		}

	}
	return prop
}

func (p *Parser) parseString() string {
	return p.currentToken.Literal
}

func (p *Parser) expectPeekType(t Token.Type) bool {
	if p.peekTokenTypeIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekTokenTypeIs(t Token.Type) bool {
	return p.peekToken.Type == t
}

func (p *Parser) peekError(t Token.Type) {
	msg := fmt.Sprintf(
		"Line: %d: Expected next token to be %s, got: %s instead",
		p.currentToken.Line,
		t,
		p.peekToken.Type,
	)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseError(msg string) {
	p.errors = append(p.errors, msg)
}

func (p *Parser) Errors() string {
	return strings.Join(p.errors, ",")
}
