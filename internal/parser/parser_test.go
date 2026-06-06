package parser

import (
	"testing"

	"github.com/rpstvs/json-parser-go/internal/ast"
	"github.com/rpstvs/json-parser-go/internal/lexer"
)

func TestJSONObjectChildren(t *testing.T) {
	tests := []struct {
		input       string
		childrenLen int
	}{
		{input: "{\"key0\": \"value0\"}", childrenLen: 1},
		{input: "{\"key0\": \"value0\" }", childrenLen: 1},
		{input: "{\"key1\": \"value1\", \"key2\": \"value2\"}", childrenLen: 2},
		{input: "{\"key3\": [\"value3\", \"value4\"]}", childrenLen: 1},
		{input: "{\"key4\": [\"value5\", {\"key5\": \"value6\"}]}", childrenLen: 1},
		{input: "{\"key5\":\" value7\", \"key6\": \"value7\"}", childrenLen: 2},
		{input: "{\"key5\":\" value7\", \"key6\": \"value7\", \"key7\": \"value8\"}", childrenLen: 3},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)

		program := p.ParseJSON()

		rv := *program.RootValue

		val := rv.(ast.Object)

		checkParserForErrors(t, p)

		if len(val.Children) != test.childrenLen {
			t.Fatalf("The length of the children does not contain 1 statement.Expected: %d; Got: %d", test.childrenLen, len(val.Children))
		}
	}
}

func checkParserForErrors(t *testing.T, p *Parser) {
	errors := p.errors

	if len(errors) == 0 {
		return
	}

	t.Errorf("Parser has %d errors", len(errors))

	for _, val := range p.errors {
		t.Errorf("Parser error: %q", val)
	}
	t.FailNow()
}
