package main

import (
	"fmt"
	"image"
	"log"
	"time"

	"github.com/markoxley/alfred/nlp"
	"github.com/markoxley/alfred/ohbot"
	"github.com/markoxley/alfred/state"
	"github.com/markoxley/alfred/vision"
)

var (
	visionExit chan bool
	visionCmd  chan int
	visionOut  chan image.Point
	ohbotState *state.State
	epoch      time.Time

	moveState int
	moveDelay time.Time
)

func main() {
	epoch = time.Now()
	err := startup()
	if err != nil {
		log.Fatalf("Unable to initialise Ohbot: %s", err.Error())
	}
	visionOut = make(chan image.Point)
	visionCmd = make(chan int)
	visionExit = make(chan bool)
	go vision.Run(visionExit, visionCmd, visionOut)

	for {
		if loop() {
			break
		}
	}
	visionExit <- true
	shutdown()
}

func startup() error {
	moveDelay = time.Now()
	err := ohbot.Init("")
	if err != nil {
		return fmt.Errorf("unable to initialise robotics: %s", err.Error())
	}
	log.Print("Robotics initialised")

	err = nlp.Init()
	if err != nil {
		return fmt.Errorf("unable to initialise NLP: %s", err.Error())
	}
	log.Print("NLP initialised")
	err = vision.Init()
	if err != nil {
		return fmt.Errorf("unable to initialise vision: %s", err.Error())
	}
	log.Print("Vision initialised")
	ohbotState = state.Init()
	return nil
}

func loop() bool {
	select {
	case <-visionOut:
		if time.Now().After(epoch) {
			ohbot.Say("Hello. I can see you", nil)
			epoch = time.Now().Add(time.Minute)
		}
	default:
	}

	if time.Now().After(moveDelay) {
		switch moveState {
		case 0:
			ohbot.Move(ohbot.HeadTurn, 1, 1)
		case 1:
			ohbot.Move(ohbot.HeadNod, 1, 1)
		case 2:
			ohbot.Move(ohbot.HeadTurn, 9, 1)
		case 3:
			ohbot.Move(ohbot.HeadNod, 9, 1)
		}
		moveState++
		if moveState == 4 {
			moveState = 0
		}
		moveDelay = time.Now().Add(time.Second)
	}

	// reader := bufio.NewReader(os.Stdin)
	// fmt.Print("-> ")
	// text, _ := reader.ReadString('\n')
	// // convert CRLF to LF
	// text = strings.Replace(text, "\n", "", -1)
	// cls, scr := nlp.GetClass(text)
	// fmt.Printf("Class: %s\tConfidence: %f\n", cls, scr)
	// d, _ := prose.NewDocument(text)
	// fmt.Println(d.Tokens())
	// fmt.Println(d.Entities())
	return false
}

func shutdown() {
	ohbot.Close()
}
