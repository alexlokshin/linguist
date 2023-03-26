package tokenizer

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"testing"
)

func TestTokenize(t *testing.T) {
	tk := NewTokenizer()

	csvFile, _ := os.Open("./testdata/products.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		tk.Ngramize(strings.ToLower(line[1]), 1, 6)
	}

	ngramIds := tk.SortNgrams()
	for _, ngramId := range ngramIds {
		fmt.Printf("%d\t%s\n", tk.Ngrams.WeightedNgrams[ngramId], tk.Ngrams.NgramIdLookup[ngramId])
	}
}
