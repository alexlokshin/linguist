package tokenizer

import (
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Tokenizer struct {
	Ngrams NGramCollection
}

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
	NgramIdLookup  map[int64]string
}

func NewTokenizer() Tokenizer {
	return Tokenizer{
		Ngrams: NGramCollection{
			WeightedNgrams: map[int64]int64{},
			NgramLookup:    map[string]int64{},
			NgramIdLookup:  map[int64]string{},
		},
	}
}

func (t Tokenizer) Ngramize(s string, min, max int) {
	tokens := t.Tokenize(s)

	for i := 0; i < len(tokens); i++ {
		for j := i; j <= len(tokens); j++ {
			if min <= j-i && j-i <= max {
				terms := tokens[i:j]
				ngram := Ngram{
					Tokens: t.tokens2Ids(terms),
				}
				ngramVal := ngram.String()

				if ngramId, ok := t.Ngrams.NgramLookup[ngramVal]; ok {
					t.Ngrams.WeightedNgrams[ngramId] = t.Ngrams.WeightedNgrams[ngramId] + 1
				} else {
					ngramId := maxNgramId
					maxNgramId++
					t.Ngrams.NgramLookup[ngramVal] = ngramId
					t.Ngrams.NgramIdLookup[ngramId] = ngramVal
					t.Ngrams.WeightedNgrams[ngramId] = 1
				}
			}
		}
	}
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
		return unicode.IsPunct(r)
	}

	items := regexp.MustCompile(`[\s]+`).Split(s, -1)
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

func (t Tokenizer) SortNgrams() []int64 {
	keys := make([]int64, 0, len(t.Ngrams.WeightedNgrams))
	for key := range t.Ngrams.WeightedNgrams {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool { return t.Ngrams.WeightedNgrams[keys[i]] > t.Ngrams.WeightedNgrams[keys[j]] })
	return keys
}
