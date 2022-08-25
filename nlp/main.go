package nlp

/**
	Test training for nlp process
**/

import (
	"encoding/gob"
	"fmt"
	"os"
	"regexp"
	"sort"

	"github.com/patrikeh/go-deep"
)

//****************************************
// Json classes for loading training data

// intent object from input json
type intent struct {
	Tag       string   `json:"tag"`
	Patterns  []string `json:"patterns"`
	Responses []string `json:"responses"`
	Function  string   `json:"function"`
}

// intents collection
type intents struct {
	Intents []intent `json:"intents"`
}

//****************************************

// wordPattern for learning
type wordPattern struct {
	words []string
	tag   string
}

type persistantModel struct {
	Words    []string
	Classes  []string
	Training deep.Dump
}

var (
	words   []string
	classes []string
	neural  *deep.Neural
)

func saveNetwork(loc string, w, c []string, n *deep.Neural) {
	p := persistantModel{
		Words:    w,
		Classes:  c,
		Training: *n.Dump(),
	}
	writePickle(loc+nlpData, p)
}

func loadNetwork(loc string) ([]string, []string, *deep.Neural, error) {
	f, err := os.Open(loc)
	if err != nil {
		return nil, nil, nil, err
	}
	defer f.Close()
	dec := gob.NewDecoder(f)
	p := persistantModel{}
	err = dec.Decode(&p)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to covert data file: %s", err.Error())
	}
	return p.Words, p.Classes, deep.FromDump(&p.Training), nil
}

// cleanStringSlice sorts an slice of strings and remove duplicates
//  @param s The slice of strings
//  @return []string
func cleanStringSlice(s []string) []string {
	sort.Strings(s)
	var old string
	res := make([]string, len(s))
	index := 0
	for i, w := range s {
		if i > 0 && w == old {
			continue
		}
		res[index] = w
		index++
		old = w
	}
	return res[:index]
}

// removePlaceHolders removes the placeholders from the input text
//  @param txt The text to clean
//  @return ; string
func removePlaceholders(txt string) string {
	var re = regexp.MustCompile(`(?mU)({{.*}})`)
	var substitution = ""
	return re.ReplaceAllString(txt, substitution)
}

func writePickle(fn string, data interface{}) error {
	f, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := gob.NewEncoder(f)
	return enc.Encode(data)
}

// func readPickle(fn string) (*gob.Decoder, error) {
// 	f, err := os.Open(fn)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer f.Close()
// 	dec := gob.NewDecoder(f)
// 	return dec, nil
// }

func GetClass(sentence string) (string, float64) {
	if len(sentence) == 0 {
		return "", 0
	}
	bag, err := createWordBag(sentence, words)
	if err != nil {
		return "", 0
	}
	predictions := neural.Predict(bag)
	class := ""
	maxVal := float64(0)
	for i, v := range predictions {
		if v > maxVal {
			maxVal = v
			class = classes[i]
		}
	}
	return class, maxVal
}
