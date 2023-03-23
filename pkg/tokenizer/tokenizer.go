package tokenizer

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Tokenizer struct{}

type Token struct {
	Id    int64
	Value string
}

var tokenDictionary map[int64]Token = map[int64]Token{}
var tokenLookupDictionary map[string]Token = map[string]Token{}
var maxTokenId int64 = 128
var maxNgramId int64 = 10240000

type NGramCollection struct {
	WeightedNgrams map[int64]int64
	NgramLookup    map[string]int64
}

func (t Tokenizer) Ngramize(s string, min, max int) NGramCollection {
	ngrams := NGramCollection{}
	tokens := t.Tokenize(s)

	for i := 0; i < len(tokens); i++ {
		for j := i; j <= len(tokens); j++ {
			if min <= j-i && j-i <= max {
				terms := tokens[i:j]
				ngram := Ngram{
					Tokens: t.tokens2Ids(terms),
				}
				ngramVal := ngram.String()

				if ngramId, ok := ngrams.NgramLookup[ngramVal]; ok {
					ngrams.WeightedNgrams[ngramId] = ngrams.WeightedNgrams[ngramId] + 1
				} else {
					ngramId := maxNgramId
					maxNgramId++
					ngrams.NgramLookup[ngramVal] = ngramId
					ngrams.WeightedNgrams[ngramId] = 1
				}
			}
		}
	}
	return ngrams
}

func (t Tokenizer) Tokenize(s string) []Token {
	tokens := []Token{}

	fixUtf := func(r rune) rune {
		if r == utf8.RuneError {
			return -1
		}
		return r
	}

	fixPunct := func(r rune) bool {
		if unicode.IsPunct(r) {
			return true
		}
		return false
	}

	items := regexp.MustCompile("[\\s]+").Split(s, -1)
	for _, item := range items {
		item = strings.Map(fixUtf, item)
		item = strings.TrimSpace(strings.TrimFunc(item, fixPunct))

		if len(item) > 0 {
			token := t.createToken(item)
			tokens = append(tokens, token)
		}
	}

	return tokens
}

func (t Tokenizer) createToken(val string) Token {
	if token, ok := tokenLookupDictionary[val]; ok {
		return token
	}
	id := maxTokenId
	maxTokenId++
	token := Token{
		Id:    id,
		Value: val,
	}
	tokenDictionary[id] = token
	tokenLookupDictionary[val] = token
	return token
}

func (t Tokenizer) terms2Tokens(terms []string) []Token {
	tokens := []Token{}
	for _, term := range terms {
		if token, ok := tokenLookupDictionary[term]; ok {
			tokens = append(tokens, token)
		}
	}
	return tokens
}

func (t Tokenizer) tokens2Ids(tokens []Token) []int64 {
	ids := []int64{}
	for _, token := range tokens {
		ids = append(ids, token.Id)
	}
	return ids
}
