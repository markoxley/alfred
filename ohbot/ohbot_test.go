package ohbot

import (
	"testing"
	"time"
)

func motorTest(mn ...MotorName) {

	for i := 5; i < 11; i += 5 {
		for p := 0; p < 2; p++ {
			for _, m := range mn {
				Move(m, float64(p*10), float64(i))
			}
			time.Sleep(time.Millisecond * 500)
		}
	}
}

func endTest() {
	time.Sleep(time.Second)
	Reset()
}

func startTest() {
	Init("")
	Reset()
}
func TestInit(t *testing.T) {
	err := Init("")
	if err != nil {
		t.Fatalf("Fail: %s", err.Error())
	}
}

func TestReset(t *testing.T) {
	err := Init("")
	if err != nil {
		t.Fatalf("Fail: %s", err.Error())
	}
	Reset()
}

func TestClose(t *testing.T) {
	Init("")
	Close()
}
func TestEyeTurn(t *testing.T) {
	startTest()
	motorTest(EyeTurn)
	endTest()
}
func TestEyeTilt(t *testing.T) {
	startTest()
	motorTest(EyeTilt)
	endTest()
}

func TestLidBlink(t *testing.T) {
	startTest()
	motorTest(LidBlink)
	endTest()
}

func TestHeadTurn(t *testing.T) {
	startTest()
	motorTest(HeadTurn)
	endTest()
}

func TestHeadNod(t *testing.T) {
	startTest()
	motorTest(HeadNod)
	endTest()
}

func TestSmile(t *testing.T) {
	startTest()
	Smile()
	endTest()
}

func TestSurprise(t *testing.T) {
	startTest()
	Surprise()
	endTest()
}

func TestFrown(t *testing.T) {
	startTest()
	Frown()
	endTest()
}

func TestSleep(t *testing.T) {
	startTest()
	Sleep()
	time.Sleep(time.Second * 2)
	Wakeup()
	endTest()
}

func TestRotate(t *testing.T) {
	startTest()
	Say("Just going to stretch my neck", nil)
	for i := 0; i < 3; i++ {
		Move(HeadTurn, 0, 2)
		Move(HeadNod, 0, 2)
		Wait(.65)
		Move(HeadNod, 10, 2)
		Wait(.65)
		Move(HeadTurn, 10, 2)
		Wait(.65)
		Move(HeadNod, 0, 2)
		Wait(.65)
	}
	Move(HeadTurn, 5)
	Move(HeadNod, 5)
	Wait(.5)

	endTest()
}

func TestSayDate(t *testing.T) {
	startTest()
	SayDate()
	SayDay()
	SayTime()
	endTest()
}

func TestSpeech(t *testing.T) {
	startTest()
	sc := NewSpeechConfig()
	sc.UntilDone = false
	Say("Hello, my name is Alfred.", sc)
	Move(HeadNod, 0)
	Wait(0.5)
	Move(HeadNod, 5)
	for IsSpeaking() {
		time.Sleep(time.Microsecond)
	}
	Say("I am alive and I am ready to party!", nil)
	endTest()
}

func TestCloseEyes(t *testing.T) {
	startTest()
	//Move(LidBlink, 0)
	Wait(10.0)
	endTest()
}
