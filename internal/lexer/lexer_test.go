package lexer

import (
	"fmt"
	"testing"

	"github.com/rpstvs/json-parser-go/internal/Token"
)

func TestNextToken(t *testing.T) {
	input := `// Initial comment
{
	"name" : "Stuart", // test comment
}
// ending comment
`

	tests := []Token.Token{
		{Type: Token.LineComment, Literal: " Initial comment\n", Line: 0, Prefix: "//"},
		{Type: Token.LeftBrace, Literal: "{", Line: 1},
		{Type: Token.Whitespace, Literal: "\n\t", Line: 1},
		{Type: Token.String, Literal: "name", Line: 2, Prefix: `"`, Suffix: `"`},
		{Type: Token.Whitespace, Literal: " ", Line: 2},
		{Type: Token.Colon, Literal: ":", Line: 2},
		{Type: Token.Whitespace, Literal: " ", Line: 2},
		{Type: Token.String, Literal: "Stuart", Line: 2, Prefix: `"`, Suffix: `"`},
		{Type: Token.Comma, Literal: ",", Line: 2},
		{Type: Token.Whitespace, Literal: " ", Line: 2},
		{Type: Token.LineComment, Literal: " test comment\n", Line: 2, Prefix: "//"},
		{Type: Token.RightBrace, Literal: "}", Line: 3},
		{Type: Token.Whitespace, Literal: "\n", Line: 3},
		{Type: Token.LineComment, Literal: " ending comment\n", Line: 4, Prefix: "//"},
		{Type: Token.EOF, Literal: "", Line: 5},
	}

	l := New(input)

	assertLexerMatches(t, l, tests)
}

func assertLexerMatches(t *testing.T, l *Lexer, tests []Token.Token) {
	for i, expectedToken := range tests {
		actualToken := l.NextToken()

		if actualToken.Type != expectedToken.Type {
			t.Fatalf("tests[%d] - tokenType wrong. Expected %s, Got %s", i, formatOutputTokenString(expectedToken), formatOutputTokenString(actualToken))
		}

		if actualToken.Literal != expectedToken.Literal {
			t.Fatalf("tests[%d] - tokenType wrong. Expected %v, Got %v", i, formatOutputTokenString(expectedToken), formatOutputTokenString(actualToken))
		}

		if actualToken.Line != expectedToken.Line {
			t.Fatalf("tests[%d] - tokenType wrong. Expected %s, Got %s", i, formatOutputTokenString(expectedToken), formatOutputTokenString(actualToken))
		}

		if actualToken.Suffix != expectedToken.Suffix {
			t.Fatalf("tests[%d] - tokenType wrong. Expected %s, Got %s", i, formatOutputTokenString(expectedToken), formatOutputTokenString(actualToken))
		}
		if actualToken.Prefix != expectedToken.Prefix {
			t.Fatalf("tests[%d] - tokenType wrong. Expected %s, Got %s", i, formatOutputTokenString(expectedToken), formatOutputTokenString(actualToken))
		}

	}
}

func formatOutputTokenString(t Token.Token) string {
	result := fmt.Sprintf("Type: %q; Literal: %q; Line:%d;", t.Type, t.Literal, t.Line)

	if t.Prefix != "" {
		result += fmt.Sprintf("Prefix: %q;", t.Prefix)
	}

	if t.Suffix != "" {
		result += fmt.Sprintf("Suffix: %q;", t.Suffix)
	}
	return result
}
