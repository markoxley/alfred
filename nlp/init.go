package nlp

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/markoxley/alfred/utils"
)

const (
	nlpData     = "nlp.mdl"
	intentsData = "intents.json"
)

// init automatically initialises the module
func Init() error {
	rand.Seed(time.Now().UnixNano())
	workingDir, err := utils.GetDataFolder()
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("unable to get data folder. %s", err.Error())
	}
	for {
		loaded, err := loadModel(workingDir + nlpData)
		if err != nil {
			return fmt.Errorf("unable to load model data: %s", err.Error())
		}
		if !loaded {
			log.Println("Training NLP model")
			err = createModel(workingDir, workingDir+intentsData)
			if err != nil {
				return fmt.Errorf("unable to create model data: %s", err.Error())
			}
			continue
		}
		log.Println("NLP model loaded")
		break
	}
	return nil
}

// loadModel loads the trained model
func loadModel(fp string) (bool, error) {
	wds, clss, n, err := loadNetwork(fp)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, fmt.Errorf("unable to parse model: %s", err.Error())
	}
	words = wds
	classes = clss
	neural = n
	return true, nil
}

// createModel creates a new model
func createModel(wd string, fp string) error {
	err := utils.TestFile(fp, nlpOriginalData)
	if err != nil {
		return fmt.Errorf("unable to create initial intents file: %s", err.Error())
	}
	return Train(wd, fp)
}
