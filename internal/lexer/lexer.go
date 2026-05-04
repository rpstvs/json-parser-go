package lexer

import (
	"github.com/rpstvs/json-parser-go/internal/Token"
)

type Lexer struct {
	Input        []rune
	char         rune //current char under examination
	position     int  //current position in input (points to current char)
	readPosition int  //current reading position in input (after current char)
	line         int  //line number for error reporting
}

func New(input string) *Lexer {
	l := &Lexer{
		Input: []rune(input),
	}
	l.ReadChar()

	return l
}

func (l *Lexer) ReadChar() {
	if l.readPosition >= len(l.Input) {
		l.char = 0
	} else {
		l.char = l.Input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) NextToken() Token.Token {
	var t Token.Token

	l.SkipWhitespace()

	switch l.char {
	case '{':
		t = NewToken(Token.LeftBrace, l.line, l.position, l.position+1, l.char)
	case '}':
		t = NewToken(Token.RightBrace, l.line, l.position, l.position+1, l.char)
	case '[':
		t = NewToken(Token.LeftBracket, l.line, l.position, l.position+1, l.char)
	case ']':
		t = NewToken(Token.RightBracket, l.line, l.position, l.position+1, l.char)
	case ':':
		t = NewToken(Token.Colon, l.line, l.position, l.position+1, l.char)
	case ',':
		t = NewToken(Token.Comma, l.line, l.position, l.position+1, l.char)
	case '"':
		t.Type = Token.String
		t.Literal = l.readString()
		t.Line = l.line
		t.Start = l.position
		t.End = l.position + 1
	case 0:
		t.Literal = ""
		t.Type = Token.EOF
		t.Line = l.line
	default:
		if isLetter(l.char) {
			t.Start = l.position
			ident := l.readIdentifier()
			t.Literal = ident
			t.Line = l.line
			t.End = l.position

			tokenType, err := Token.LookupIdentifier(ident)
			if err != nil {
				t.Type = Token.Illegal
				return t
			}
			t.Type = tokenType
			t.End = l.position
			return t
		} else if isNumber(l.char) {
			t.Start = l.position
			t.Literal = l.readNumber()
			t.Line = l.line
			t.End = l.position
			return t
		}
		t = NewToken(Token.Illegal, l.line, 1, 2, l.char)
	}
	l.ReadChar()
	return t
}

func (l *Lexer) SkipWhitespace() {
	for l.char == ' ' || l.char == '\t' || l.char == '\n' || l.char == '\r' {
		if l.char == '\n' {
			l.line++
		}
		l.ReadChar()
	}
}

func NewToken(tokenType Token.Type, line, start, end int, char ...rune) Token.Token {
	return Token.Token{
		Type:    tokenType,
		Line:    line,
		Literal: string(char),
		Start:   start,
		End:     end,
	}
}

func (l *Lexer) readString() string {
	position := l.position + 1

	for {
		prevChar := l.char
		l.ReadChar()
		if l.char == '"' && prevChar != '\\' || l.char == 0 {
			break
		}
	}
	return string(l.Input[position:l.position])
}

func (l *Lexer) readNumber() string {
	position := l.position

	for isNumber(l.char) {
		l.ReadChar()
	}
	return string(l.Input[position:l.position])
}

func isNumber(input rune) bool {
	return input >= '0' && input <= '9' || input == '-' || input == '.'
}

func isLetter(input rune) bool {
	return input >= 'a' && input <= 'z' || input >= 'A' && input <= 'Z'
}

func (l *Lexer) readIdentifier() string {
	position := l.position

	for isNumber(l.char) {
		l.ReadChar()
	}
	return string(l.Input[position:l.position])
}
