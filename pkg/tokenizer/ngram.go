package tokenizer

import "strings"

type Ngram struct {
	Id     int64
	Tokens []int64
}

func (ng Ngram) String() string {
	sb := strings.Builder{}
	for i, token := range ng.Tokens {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(tokenDictionary[token].Value)
	}
	return sb.String()
}
