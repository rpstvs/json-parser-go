package Token

type Type string

type Token struct {
	Type    Type
	Literal string
	Line    int
	Start   int
	End     int
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

	True  Type = "TRUE"
	False Type = "FALSE"
	Null  Type = "NULL"
)
