package dora

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/rpstvs/json-parser-go/internal/ast"
)

var ErrNoDollarSignRoot = errors.New("Incorrect syntax: query must start with $ representing the root object")
var ErrWrongObjectRootSelector = errors.New("Incorrect syntax. Your root JSON type is an object. Therefore, path queries must" +
	"begin by selecting a `key` from your root object. Ex: `$.keyOnRootObject` or `$[\"keyOnRootObject\"]`")
var ErrWrongArrRootSelector = errors.New(
	"Incorrect syntax. Your root JSON type is an array. Therefore, path queries must" +
		"begin by selecting an item by index on the root array. Ex: `$[0]` or `$[1]`",
)

func (c *Client) prepAndExecQuery(query string) error {
	if err := c.prepareQuery(query, c.tree.Type); err != nil {
		return err
	}

	if err := c.executeQuery(); err != nil {
		return err
	}
	return nil
}

func (c *Client) get(query string) (string, error) {
	if err := c.prepAndExecQuery(query); err != nil {
		return "", err
	}

	return c.result, nil
}

func (c *Client) prepareQuery(query string, rootNodeType ast.RootNodeType) error {

	if err := validateQueryRoot(query, c.tree.Type); err != nil {
		return err
	}

	c.SetQuery([]byte(query))

	if err := c.parseQuery(); err != nil {
		return err
	}

	return nil
}

func (c *Client) SetQuery(query []byte) {
	c.query = query
}

func (c *Client) parseQuery() error {
	tokens, err := scanQueryTokens(c.query)

	if err != nil {
		return err
	}
	c.parsedQuery = tokens
	return nil
}

const Object = "object"
const Array = "array"

func (c *Client) executeQuery() error {
	rootVal := *c.tree.RootValue
	obj, _ := rootVal.(ast.Object)
	arr, ok := rootVal.(ast.Array)

	currentType := Object

	if ok {
		currentType = Array
	}

	parsedQueryLen := len(c.parsedQuery)

	for i := 0; i < parsedQueryLen; i++ {
		if i == parsedQueryLen-1 {
			c.setFinalValue(currentType, i, obj, arr)
		}
		if c.parsedQuery[i].AccessType == ObjectAcess {
			if currentType != Object {
				return fmt.Errorf("error")
			}

			var found bool

			for _, v := range obj.Children {
				if v.Key.Value == c.parsedQuery[i].keyReq {
					found = true
					o, astObj := v.Value.(ast.Object)
					a, astArr := v.Value.(ast.Array)

					if astObj {
						obj = o
						currentType = Object
						break
					}

					if astArr {
						arr = a
						currentType = Array
						break
					}
				}
			}
			if !found {
				return fmt.Errorf("Sorry, could not find a key with that value. Key: %s", c.parsedQuery[i].keyReq)
			}
		} else {
			if currentType != Array {
				return fmt.Errorf("error not array")
			}

			qt := c.parsedQuery[i]
			val := arr.Children[qt.indexReq]

			switch v := val.(type) {
			case ast.Object:
				obj = v
				currentType = Object
				break
			case ast.Array:
				arr = v
				currentType = Array
				break
			case ast.Literal:
				if i == parsedQueryLen-1 {
					c.setResultFromLiteral(v.Value)
				} else {
					fmt.Println("error")
				}
			}
		}
	}
	return nil
}

func (c *Client) setFinalValue(currentType string, index int, obj ast.Object, arr ast.Array) {
	if currentType == Object {
		r := c.parsedQuery[index].keyReq

		for _, v := range obj.Children {
			if r == v.Key.Value {
				c.setResultFromValue(v.Value)
				break
			}
		}
		return

	}
	ind := c.parsedQuery[index].indexReq
	c.setResultFromValue(arr.Children[ind])
}

func (c *Client) setResultFromValue(value ast.Value) {
	switch val := value.(type) {
	case ast.Literal:
		c.setResultFromLiteral(val.Value)
	case ast.Object:
		c.result = string(c.input[val.Start:val.End])
	case ast.Array:
		c.result = string(c.input[val.Start:val.End])
	}
}

func (c *Client) setResultFromLiteral(value ast.Value) {
	switch lit := value.(type) {
	case string:
		c.result = lit
	case int:
		c.result = strconv.Itoa(lit)
	case bool:
		c.result = fmt.Sprintf("%v", lit)
	case nil:
		c.result = "null"
	}
}

func validateQueryRoot(query string, rooNodeType ast.RootNodeType) error {
	if query[0] != '$' {
		return ErrNoDollarSignRoot
	}
	validObjQueryRoot := query[1] == '.'

	if rooNodeType == ast.ObjectRoot && !validObjQueryRoot {
		return ErrWrongObjectRootSelector
	}

	validArrQueryRoot := query[1] == '['

	if rooNodeType == ast.ArrayRoot && !validArrQueryRoot {
		return ErrWrongArrRootSelector
	}
	return nil
}

func errSelectorSyntax(operator string) error {
	return fmt.Errorf("Error parsing query, expected either a `.` for selections on an object or a `[` for selections on an array. Got: %s",
		operator)
}
