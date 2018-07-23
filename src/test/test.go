package main

import (
	"fmt"

	"github.com/fluhus/gostuff/nlp/wordnet"
)

func main() {
	/* questions := make([]string, 0)
	questions = append(questions, "what is go in golang?", "it's not right", "beauty's gold the", "a ; b; v")
	questions = append(questions, "Range over integers in golang")
	for _, q := range questions {
		cleanContent := stopwords.CleanString(q, "en", true)
		println(cleanContent)
	} */

	wn, err := wordnet.Parse("../../wordnet")
	if err != nil {
		fmt.Println(err)
	} else {
		cat := wn.Search("feline")["n"][0]
		dog := wn.Search("cat")["n"][0]
		similarity := wn.WupSimilarity(cat, dog, false)
		fmt.Println(similarity)
	}

}
