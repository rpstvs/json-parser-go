package lexer

import (
	"testing"

	"github.com/rpstvs/json-parser-go/internal/Token"
)

func TestNextToken(t *testing.T) {
	input := `{
	"name" : "Stuart",
}
`

	tests := []Token.Token{
		{Type: Token.LeftBrace, Literal: "{", Line: 0},
		{Type: Token.Whitespace, Literal: "\n\t", Line: 1},
		{Type: Token.String, Literal: "name", Line: 1, Prefix: `"`, Suffix: `"`},
		//{Type: Token.Whitespace, Literal: " ", Line: 2},
		{Type: Token.Colon, Literal: ":", Line: 1},
		//{Type: Token.Whitespace, Literal: " ", Line: 2},
		{Type: Token.String, Literal: "Stuart", Line: 1, Prefix: `"`, Suffix: `"`},
		{Type: Token.Comma, Literal: ",", Line: 1},
		//{Type: Token.Whitespace, Literal: "\n\t", Line: 2},
		{Type: Token.RightBrace, Literal: "}", Line: 2},
	}

	l := New(input)

	assertLexerMatches(t, l, tests)
}

func assertLexerMatches(t *testing.T, l *Lexer, tests []Token.Token) {
	for i, expectedToken := range tests {
		actualToken := l.NextToken()

		if actualToken.Type != expectedToken.Type {
			t.Fatalf("tests[%d] - tokenType wrong. Expected %s, Got %s", i, expectedToken.Type, actualToken.Type)
		}

		if actualToken.Literal != expectedToken.Literal {
			t.Fatalf("tests[%d] - tokenType wrong. Expected %s, Got %s", i, expectedToken.Literal, actualToken.Literal)
		}

		if actualToken.Line != expectedToken.Line {
			t.Fatalf("tests[%d] - tokenType wrong. Expected %d, Got %d", i, expectedToken.Line, actualToken.Line)
		}

		if actualToken.Suffix != expectedToken.Suffix {
			t.Fatalf("tests[%d] - tokenType wrong. Expected %s, Got %s", i, expectedToken.Suffix, actualToken.Suffix)
		}
		if actualToken.Prefix != expectedToken.Prefix {
			t.Fatalf("tests[%d] - tokenType wrong. Expected %s, Got %s", i, expectedToken.Prefix, actualToken.Prefix)
		}

	}
}
