package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/alexlokshin/linguist/pkg/tokenizer"
)

func main() {
	// app := fiber.New()

	// app.Get("/", func(c *fiber.Ctx) error {
	// 	return c.SendString("Hello, World!")
	// })

	// log.Fatal(app.Listen(":3000"))

	tk := tokenizer.NewTokenizer()

	csvFile, _ := os.Open("./pkg/tokenizer/testdata/holmes.txt")
	reader := csv.NewReader(bufio.NewReader(csvFile))
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		tk.Ngramize(strings.ToLower(line[1]), 2, 6)
	}

	// f, err := os.OpenFile("stopwords.txt",
	// 	os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	log.Println(err)
	// }
	// defer f.Close()

	ngramIds := tk.SortNgrams()
	for _, ngramId := range ngramIds {
		fmt.Printf("%d\t%s\n", tk.Ngrams.WeightedNgrams[ngramId], tk.Ngrams.NgramIdLookup[ngramId])
		// stopWord := false
		// prompt := &survey.Confirm{
		// 	Message: "Stopword?",
		// }
		// survey.AskOne(prompt, &stopWord)
		// if stopWord {
		// 	if _, err := f.WriteString(tk.Ngrams.NgramIdLookup[ngramId] + "\n"); err != nil {
		// 		log.Println(err)
		// 	}
		// }
	}
}
