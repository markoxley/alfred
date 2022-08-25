package ohbot

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"
)

type phrase struct {
	set      string
	variable string
	text     string
}

type SpeechConfig struct {
	UntilDone  bool
	LipSync    bool
	HDMIAudo   bool
	SoundDelay float64
}

var (
	phraseList         []*phrase
	phenomeTop         map[string]float64
	phenomeBottom      map[string]float64
	speechDatabaseFile string
	speechMutex        sync.Mutex
	speaking           int
	textToSpeak        []string
)

func NewSpeechConfig() *SpeechConfig {
	return &SpeechConfig{
		UntilDone:  true,
		LipSync:    true,
		HDMIAudo:   false,
		SoundDelay: 0,
	}
}
func newPhrase(set, variable, text string) *phrase {
	return &phrase{
		set:      set,
		variable: variable,
		text:     text,
	}
}

func IsSpeaking() bool {
	speechMutex.Lock()
	defer speechMutex.Unlock()
	return speaking > 0
}

func startSpeaking() {
	speechMutex.Lock()
	defer speechMutex.Unlock()
	speaking++
}

func endSpeaking() {
	speechMutex.Lock()
	defer speechMutex.Unlock()
	speaking--
}

func loadSpeechDatabase() {
	cf, err := os.Open(speechDatabaseFile)
	if err != nil {
		log.Fatalf("Unable to open speech database: %s", err.Error())
	}
	defer cf.Close()
	r := csv.NewReader(cf)
	for {
		r, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		phraseList = append(phraseList, newPhrase(r[0], r[1], r[2]))
	}
}

func generateSpeechFile(text string) error {
	// Note: This does not support gTTs
	re := regexp.MustCompile(`(?m)[^ .a-zA-Z0-9?\']`)
	st := re.ReplaceAllString(text, "")
	var bc string
	var args []string
	bc = "/usr/bin/festival"
	args = []string{
		"-b",
		fmt.Sprintf("(set! mytext (Utterance Text \"%s\"))", st),
		"(utt.synth mytext)",
		fmt.Sprintf("(utt.save.wave mytext \"%s\")", speechAudioFile),
		fmt.Sprintf("(utt.save.segs mytext \"%s\")", phonemesFile),
	}

	cmd := exec.Command(bc, args...)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func playSpeech() {
	PlaySoundFile(speechAudioFile)
}

func moveSpeech(ph []string, tm []float64) {
	startTime := time.Now()
	timeNow := float64(0)
	totalTime := tm[len(tm)-1]
	currentX := -1
	for timeNow < totalTime {
		timeNow = float64(time.Since(startTime).Seconds())
		for i, t := range tm {
			if timeNow > t && i > currentX {
				currentX = i
				phenome := strings.TrimSpace(ph[i])

				posTop, posBottom := getPhenome(phenome)
				Move(TopLip, posTop, 10)
				Move(BottomLip, posBottom, 10)

			}
		}
	}
	Move(TopLip, 5)
	Move(BottomLip, 5)
}

func getPhenome(ph string) (float64, float64) {
	t := float64(5)
	b := float64(5)
	if v, ok := phenomeTop[ph]; ok {
		t = v
	}
	if v, ok := phenomeBottom[ph]; ok {
		b = v
	}
	return t, b
}
