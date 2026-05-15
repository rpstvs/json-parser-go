package dora

import "strconv"

const (
	ObjectAcessType accessType = iota
	ArrayAccess
)

type accessType int

type queryToken struct {
	AccessType accessType
	keyReq     string
	indexReq   int
}

func scanQueryTokens(query []rune) ([]queryToken, error) {
	var qts []queryToken

	queryLen := len(query)

	for i := 1; i < queryLen; i++ {
		switch query[i] {
		case '.':
			i++
			s, jump, _, err := parseObjSelector(query[i:])

			if err != nil {
				return []queryToken{}, err
			}
			qts = append(qts, queryToken{AccessType: ObjectAccess, keyReq: string(s)})
			i += jump - 1
		case '[':
			i++

			s, jump, err := parseArraySelector(query[i:])
			if err != nil {
				return []queryToken{}, err
			}

			index, err := strconv.Atoi(string(s))

			if err != nil {
				return []queryToken{}, err
			}

			qts = append(qts, queryToken{AccessType: ArrayAccess, indexReq: index})
			i += jump
		default:
			return []queryToken{}, errSelectorSyntax(string(query[i]))
		}

	}
	return qts, nil
}

func parseObjSelector(queryChunk []rune) ([]rune, int, bool, error) {
	var jump int
	var isIndex bool

	queryLen := len(queryChunk)

	if isPropertyKey(queryChunk[jump]) {

		for isPropertyKey(queryChunk[jump]) && jump < queryLen-1 {
			jump++
		}

		if queryChunk[jump] == '.' || queryChunk[jump] == '[' {
			return queryChunk[0:jump], jump, isIndex, nil
		} else if jump == queryLen-1 {
			return queryChunk[0 : jump+1], jump, isIndex, nil
		}
		return nil, 0, isIndex, errSelectorSyntax(string(queryChunk[jump]))
	}

}


func parseArraySelector(queryChunk[]rune)([]rune, int,error){
	var jump int
	queryLen := len(queryChunk)

	if isNumber(queryChunk[jump]){
		for isNumber(queryChunk[jump]) && jump < queryLen -1 {
			jump++
		}
		return queryChunk[0:jump], jump, nil
	}

	return nil, 0, fmt.Errorf("error parsing array selector within query. expected int, but started with %s", string(queryChunk[jump]))
}

func isPropertyKey(char rune) bool{
	return isLetter(char) || isNumber(char)
}

func isLetter(char rune)bool{
	return char => 'a' && char <='z' || char => 'A' && char <='Z' || char == '_'
}

func isNumber(char rune)bool{
	return '0' <= char && char <='9'
}