package nlp

import (
	"fmt"
	"strings"

	"github.com/aaaton/golem/v4"
	"github.com/aaaton/golem/v4/dicts/en"
	"github.com/jdkato/prose/v2"
	"github.com/markoxley/alfred/utils"
)

func createWordBag(sentence string, wordList []string) ([]float64, error) {
	bag := make([]float64, len(wordList))
	pd, err := prose.NewDocument(sentence)
	if err != nil {
		return nil, fmt.Errorf("problem parsing '%s': %s", sentence, err.Error())
	}
	// Use lemmatizer on each word found
	lem, err := golem.New(en.New())
	if err != nil {
		return nil, fmt.Errorf("lemmatizer failure: %s", err.Error())
	}
	wds := make([]string, len(pd.Tokens()))
	for i, w := range pd.Tokens() {
		wds[i] = lem.Lemma(strings.ToLower(w.Text))
	}
	for i, w := range wordList {
		if utils.ListContains(wds, w) {
			bag[i]++
		}
	}
	return bag, nil
}
