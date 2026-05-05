package ast

const (
	ObjectRoot RootNodeType = iota
	ArrayRoot
)

type RootNodeType int

type RootNode struct {
	RootValue *Value
	Type      RootNodeType
}

type Object struct {
	Type     string
	Children []Property
	Start    int
	End      int
}

type Array struct {
	Type     string
	Children []Value
	Start    int
	End      int
}

type Literal struct {
	Type  string
	Value Value
}

type Property struct {
	Type  string
	Key   Identifier
	Value Value
}

type Identifier struct {
	Type  string
	Value string
}

type Value interface{} //or any

type state int

const (
	ObjStart state = iota
	ObjOpen
	ObjProperty
	ObjComma

	PropertyStart
	PropertyKey
	PropertyColon

	ArrayStart
	ArrayOpen
	ArrayValue
	ArrayComma

	StringStart
	StringQuoteOrChar
	Escape

	NumberStart
	NumberMinus
	NumberZero
	NumberDigit
	NumberPoint
	NumberDigitFraction
	NumberExp
	NumberExpDigitOrSign
)
