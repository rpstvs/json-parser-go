package Token

import "fmt"

type Type string

type Token struct {
	Type    Type
	Literal string
	Line    int
	Start   int
	End     int
	Prefix  string
	Suffix  string
}

const (
	Illegal Type = "ILLEGAL"

	EOF Type = "EOF"

	String Type = "STRING"
	Number Type = "NUMBER"

	LeftBrace    Type = "{"
	RightBrace   Type = "}"
	LeftBracket  Type = "["
	RightBracket Type = "]"
	Comma        Type = ","
	Colon        Type = ":"

	Whitespace Type = "WHITESPACE"

	LineComment  Type = "LINECOMMENT"
	BlockComment Type = "BLOCKCOMMENT"

	True  Type = "TRUE"
	False Type = "FALSE"
	Null  Type = "NULL"
)

var validJsonIdentifiers = map[string]Type{
	"False": False,
	"True":  True,
	"Null":  Null,
}

func LookupIdentifier(ident string) (Type, error) {
	if val, ok := validJsonIdentifiers[ident]; ok {
		return val, nil
	}
	return Illegal, fmt.Errorf("Expected valid json identifier. Found %s", ident)
}
