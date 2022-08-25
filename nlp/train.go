package nlp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"

	"github.com/aaaton/golem/v4"
	"github.com/aaaton/golem/v4/dicts/en"
	"github.com/jdkato/prose/v2"
	"github.com/markoxley/alfred/utils"
	"github.com/patrikeh/go-deep"
	"github.com/patrikeh/go-deep/training"
)

// Train trains the data from the intents json
func Train(wd string, ij string) error {
	// Load intents json file
	fb, err := ioutil.ReadFile(ij)
	if err != nil {
		return fmt.Errorf("failed to open file: %s", err.Error())
	}

	// Convert the intents into an intents collection
	var js intents
	err = json.Unmarshal(fb, &js)
	if err != nil {
		log.Fatalf("Unable to unmarshal json: %s", err.Error())
	}

	newWords := []string{}
	newClasses := []string{}
	documents := []wordPattern{}
	ignore_chars := []string{"?", "!", "<", ">", ",", "."}

	// Loop through each intent in the collection
	for _, intent := range js.Intents {

		// Loop through each pattern for the current intent
		for _, pattern := range intent.Patterns {

			// Rmove placeholders
			np := removePlaceholders(pattern)
			// Use prose to parse the pattern
			pd, err := prose.NewDocument(np)
			if err != nil {
				return fmt.Errorf("problem parsing '%s': %s", pattern, err.Error())
			}

			// Create a wordlist for the pattern
			wordList := make([]string, 0)

			// Add each token to the word list for the pattern and the word list for the overall network
			for _, token := range pd.Tokens() {

				// if the word is in the ignore list, ignore it
				if utils.ListContains(ignore_chars, token.Text) {
					continue
				}
				wordList = append(wordList, token.Text)
				newWords = append(newWords, token.Text)
			}

			// Add the wordlist with the intent to the document
			documents = append(documents, wordPattern{
				words: wordList,
				tag:   intent.Tag,
			})

			// If we do not already have an intent tag stored in the classes, add it
			if !utils.ListContains(newClasses, intent.Tag) {
				newClasses = append(newClasses, intent.Tag)
			}
		}
	}

	// Use lemmatizer on each word found
	lem, err := golem.New(en.New())
	if err != nil {
		return fmt.Errorf("lemmatizer failure: %s", err.Error())
	}

	for i, w := range newWords {
		newWords[i] = lem.Lemma(strings.ToLower(w))
	}

	// Sort and remove duplicates
	newWords = cleanStringSlice(newWords)
	newClasses = cleanStringSlice(newClasses)

	trainingData := training.Examples(make([]training.Example, 0, len(documents)))
	// gobrainData := [][][]float64{}
	for _, document := range documents {
		bag := make([]float64, len(newWords))
		wordPatterns := document.words
		for i, w := range wordPatterns {
			wordPatterns[i] = lem.Lemma(strings.ToLower(w))
		}
		for i, w := range newWords {
			if utils.ListContains(wordPatterns, w) {
				bag[i]++
			}
		}
		output := make([]float64, len(newClasses))
		output[utils.ListIndex(newClasses, document.tag)] = 1
		trainingData = append(trainingData, training.Example{
			Input:    bag,
			Response: output,
		})
		// gobrainData = append(gobrainData, [][]float64{bag, output})

	}

	rand.Shuffle(len(trainingData), func(i, j int) {
		trainingData[i], trainingData[j] = trainingData[j], trainingData[i]
	})

	// rand.Shuffle(len(gobrainData), func(i, j int) {
	// 	gobrainData[i], gobrainData[j] = gobrainData[j], gobrainData[i]
	// })

	layerSize := ((len(trainingData[0].Input) * 2) / 3) + len(trainingData[0].Response)
	//	 int(len(trainingData[0].Input) - ((len(trainingData[0].Input) - len(trainingData[0].Response)) / 2))
	n := deep.NewNeural(&deep.Config{

		Inputs: len(trainingData[0].Input),
		Layout: []int{layerSize, len(trainingData[0].Response)},
		//Loss:       deep.LossCrossEntropy,
		Activation: deep.ActivationSigmoid,
		Mode:       deep.ModeMultiClass,
		Weight:     deep.NewNormal(1.0, 0.0),
		Bias:       true,
	})
	optimizer := training.NewSGD(0.01, 0.1, 1e-10, true)
	trainer := training.NewTrainer(optimizer, 1000)

	//training, heldout := trainingData.Split(.95)
	// trainer.Train(n, training, heldout, 100000)
	trainer.Train(n, trainingData, trainingData, 10000)

	// log.Println("Gobrain training")
	// ff := &gobrain.FeedForward{
	// 	Regression: true,
	// }
	// ff.Init(len(gobrainData[0][0]), layerSize, len(gobrainData[0][1]))

	// ff.Train(gobrainData, 100000, 0.6, 0.4, true)
	//ff.Test(gobrainData)
	saveNetwork(wd, newWords, newClasses, n)
	words = newWords
	classes = newClasses
	log.Println("Training complete.")
	return nil
}
