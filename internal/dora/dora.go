package dora

import (
	"strconv"

	"github.com/rpstvs/json-parser-go/internal/ast"
	"github.com/rpstvs/json-parser-go/internal/lexer"
	"github.com/rpstvs/json-parser-go/internal/parser"
)

type Client struct {
	input       []byte
	tree        *ast.RootNode
	query       []byte
	parsedQuery []queryToken
	result      string
}

func NewFromString(json string) (*Client, error) {
	l := lexer.New(json)
	p := parser.New(l)

	tree := p.ParseJSON()

	return &Client{
		tree:  &tree,
		input: l.Input,
	}, nil
}

func NewFromBytes(bytes []byte) (*Client, error) {
	return NewFromString(string(bytes))
}

func (c *Client) GetString(query string) (string, error) {
	result, err := c.get(query)

	if err != nil {
		return "", err
	}

	return result, nil

}

func (c *Client) GetBool(query string) (bool, error) {
	result, err := c.get(query)

	if err != nil {
		return false, err
	}

	s, err := strconv.ParseBool(result)

	if err != nil {
		return false, err
	}

	return s, nil

}

func (c *Client) getFloat(query string) (float64, error) {
	result, err := c.get(query)

	if err != nil {
		return 0.00, err
	}

	s, err := strconv.ParseFloat(result, 64)

	if err != nil {
		return 0.00, err
	}

	return s, nil

}
