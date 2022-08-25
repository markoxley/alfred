package ohbot

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/beevik/etree"
	"github.com/markoxley/alfred/utils"
)

const (
	version = "1.0.0"
	dirName = "ohbotData"
)

var (
	sensors []float64
	//shapeList       []float64
	port            string
	writing         bool
	connected       bool
	topLipFree      bool
	ser             io.ReadWriteCloser
	workingDir      string
	pathSep         string
	speechAudioFile string
	settingsFile    string
	phonemesFile    string
)

func init() {
	rand.Seed(time.Now().UnixNano())
	sensors = []float64{0, 0, 0, 0, 0, 0, 0, 0}
	for i := uint8(0); i <= uint8(MouthOpen); i++ {
		motors = append(motors, newMotor())
	}
	writing = false
	connected = false
	topLipFree = false
	workingDir, err := utils.GetDataFolder()
	if err != nil {
		log.Fatalf("Unable to determine working directory: %s", err.Error())
	}
	ohbotMotorDefFile = workingDir + "MotorDefinitionsv21.omd"
	speechAudioFile = workingDir + "ohbotspeech.wav"
	settingsFile = workingDir + "OhbotSettings.xml"
	phonemesFile = workingDir + "phonemes"
	//settingsFile = workingDir + settingsFile
	if err = utils.TestFile(settingsFile, settingsDef); err != nil {
		log.Fatalf("Unable to create default XML file. %s", err.Error())
	}

	if err = loadSettings(); err != nil {
		log.Fatalf("Unable to load settings. %s", err.Error())
	}

	if err = utils.TestFile(speechDatabaseFile, speechDef); err != nil {
		log.Fatalf("Unable to create speech database file. %s", err.Error())
	}

	if err = utils.TestFile(ohbotMotorDefFile, motorDef); err != nil {
		log.Fatalf("Unable to create ohbot motor definition file. %s", err.Error())
	}

	phenomeTop = map[string]float64{
		"p":  5,
		"b":  5,
		"m":  5,
		"ae": 7,
		"ax": 7,
		"ah": 7,
		"aw": 10,
		"aa": 10,
		"ao": 10,
		"ow": 10,
		"ey": 7,
		"eh": 7,
		"uh": 7,
		"ay": 7,
		"h":  7,
		"er": 8,
		"r":  8,
		"l":  8,
		"y":  6,
		"iy": 6,
		"ih": 6,
		"ix": 6,
		"w":  6,
		"uw": 6,
		"oy": 6,
		"s":  5,
		"z":  5,
		"sh": 5,
		"ch": 5,
		"jh": 5,
		"zh": 5,
		"th": 5,
		"dh": 5,
		"d":  5,
		"t":  5,
		"n":  5,
		"k":  5,
		"g":  5,
		"ng": 5,
		"f":  6,
		"v":  6,
	}

	phenomeBottom = map[string]float64{
		"p":  5,
		"b":  5,
		"m":  5,
		"ae": 8,
		"ax": 8,
		"ah": 8,
		"aw": 5,
		"aa": 10,
		"ao": 10,
		"ow": 10,
		"ey": 7,
		"eh": 7,
		"uh": 7,
		"ay": 7,
		"h":  7,
		"er": 8,
		"r":  8,
		"l":  8,
		"y":  6,
		"iy": 6,
		"ih": 6,
		"ix": 6,
		"w":  6,
		"uw": 6,
		"oy": 6,
		"s":  6,
		"z":  6,
		"sh": 6,
		"ch": 6,
		"jh": 6,
		"zh": 6,
		"th": 6,
		"dh": 6,
		"d":  6,
		"t":  6,
		"n":  6,
		"k":  6,
		"g":  6,
		"ng": 6,
		"f":  5,
		"v":  5,
	}
}

func isDigit(s string) bool {
	_, err := strconv.Atoi(s)
	return err != nil
}

func loadSettings() error {
	tree := etree.NewDocument()
	if err := tree.ReadFromFile(settingsFile); err != nil {
		return fmt.Errorf("unable to read settings file: %s", err.Error())
	}

	root := tree.SelectElement("SettingList")
	for _, element := range root.SelectElements("Setting") {
		value := element.SelectAttrValue("Value", "")
		switch element.SelectAttrValue("Name", "") {
		case "SpeechDBFile":
			speechDatabaseFile = workingDir + value
		case "MotorDefFile":
			ohbotMotorDefFile = workingDir + value
		}
	}
	return nil
}

func listSerialPorts() ([]string, error) {

	dir, err := os.ReadDir("/dev")
	if err != nil {
		return nil, err
	}
	res := make([]string, 0, len(dir))
	for _, d := range dir {
		if d.IsDir() {
			continue
		}
		if len(d.Name()) < 6 {
			continue
		}
		if d.Name()[:4] == "ttyA" {
			res = append(res, "/dev/"+d.Name())
		}
	}
	return res, nil
}

func limit(v float64, l ...float64) float64 {
	mn := float64(0)
	mx := float64(9)
	if len(l) > 0 {
		mn = l[0]
	}
	if len(l) > 1 {
		mx = l[1]
	}
	if v > mx {
		return mx
	}
	if v < mn {
		return mn
	}
	return v
}

func serWrite(s string) {
	if !connected {
		return
	}
	writing = true
	ser.Write([]byte(s))
	writing = false
}

func autoBlink(wg *sync.WaitGroup, stop chan bool) {
	// Move(LidBlink, 5, 5)
	wg.Add(1)
	go func() {
		time.Sleep(time.Second * 2)
		brk := false
		for {
			select {
			case <-stop:
				brk = true
			default:
			}
			if brk {
				break
			}
			delay := 5000 + rand.Intn(5000)
			time.Sleep(time.Millisecond * time.Duration(delay))
			Blink()
		}
		wg.Done()
	}()
}

func stopBlink(wg *sync.WaitGroup, stop chan bool) {
	stop <- true
	wg.Wait()
}
