package ohbot

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/huin/goserial"
)

type MotorName uint8

const (
	HeadNod MotorName = iota
	HeadTurn
	EyeTurn
	LidBlink
	TopLip
	BottomLip
	EyeTilt
	MouthOpen
)

var (
	moveMutex sync.Mutex
	blinkChan chan bool
	blinkWg   sync.WaitGroup
)

// Version returns the version of the module
//  @return string
func Version() string {
	return version
}

// Init initialises Ohbot on the specified port.
// If an empty string is passed, the function will
// attempt to find the port that is attached to Ohbot
//  @param portName is the name of the port to be used. If unknown, pass an empty string
//  @return error
func Init(portName string) error {

	// Load the motor definitions
	if err := loadMotorDefs(); err != nil {
		return err
	}

	// Find the port for Ohbot
	ports, err := listSerialPorts()
	if err != nil {
		return err
	}
	if portName == "" {
		for _, p := range ports {
			if CheckPort(p) {
				port = p
				connected = true
				break
			}
		}
		if port == "" {
			return errors.New("unable to find Ohbot port")
		}
	} else {
		if !CheckPort(portName) {
			return fmt.Errorf("unable to find Ohbot on port %s", portName)
		}
		connected = true
		port = portName
	}

	c := &goserial.Config{
		Name: port,
		Baud: 19200,
	}
	ser, err = goserial.OpenPort(c)
	if err != nil {
		return err
	}

	// Load the speech database
	loadSpeechDatabase()
	blinkChan = make(chan bool)
	autoBlink(&blinkWg, blinkChan)
	Reset()
	return nil
}

// CheckPort checks the specified port for the existance of Ohbot
//  @param p is the port
//  @return bool
func CheckPort(p string) bool {
	c := &goserial.Config{
		Name: p,
		Baud: 19200,
	}
	ser, err := goserial.OpenPort(c)
	if err != nil {
		return false
	}
	defer ser.Close()

	msg := "v\n"
	_, err = ser.Write([]byte(msg))
	if err != nil {
		return false
	}
	buf := make([]byte, 256)
	n, err := ser.Read(buf)
	if err != nil {
		return false
	}
	line := string(buf[:n])
	return strings.Contains(line, "v1")
}

// PlaySoundFile makes Ohbot play a sound file
//  @param fp the filepath of the sound file
//  @return error
func PlaySoundFile(fp string) error {
	cmd := exec.Command("aplay", fp)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
	//return cmd.Wait()
}

// Move moves the specified servo
//  @param mn is the name of the servo
//  @param val values for the turn. First value is position, second value is speed
func Move(mn MotorName, val ...float64) {
	moveMutex.Lock()
	defer moveMutex.Unlock()
	pos := float64(5)
	spd := float64(5)
	if len(val) > 0 {
		pos = limit(val[0])
	}
	if len(val) > 1 {
		spd = limit(val[1], 1, 10)
	}
	if pos > 9 && mn == BottomLip {
		topLipFree = true
	}

	if pos <= 5 && mn == BottomLip {
		topLipFree = false
	}

	if pos < 5 && mn == BottomLip {
		pos = 5 - ((5 - pos) / 2)
	}

	m := motors[int(mn)]
	if m.rev {
		pos = 10 - pos
	}

	m.attach()
	absPos := m.absPos(pos)
	spd = float64(250/10) * spd
	msg := fmt.Sprintf("m0%v,%v,%v\n", mn, absPos, spd)
	serWrite(msg)

	m.pos = pos
}

// Attach attaches the specified server
//  @param mn is the name of the servo
func Attach(mn MotorName) {
	motors[int(mn)].attach()
}

// Detach detaches the servo
//  @param mn is the name of the servo to detach
func Detach(mn MotorName) {
	motors[int(mn)].detach()
}

// Say speaks the selected text
//  @param text is the text to speak
//  @param sc optional additional configuration
func Say(text string, sc *SpeechConfig) {
	startSpeaking()
	text = strings.TrimSpace(text)
	if text == "" {
		return
	}
	text = strings.ReplaceAll(text, "picoh", "peek oh")
	text = strings.ReplaceAll(text, "Picoh", "peek oh")

	if sc == nil {
		sc = NewSpeechConfig()
	}
	generateSpeechFile(text)

	f, err := os.Open(phonemesFile)
	if err != nil {
		return
	}
	defer f.Close()
	phonemes := make([]string, 0)
	times := make([]float64, 0)

	b := bufio.NewReader(f)
	for {
		line, err := b.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return
		}
		vals := strings.Split(string(line), " ")
		if len(vals) < 3 {
			continue
		}

		pt, err := strconv.ParseFloat(vals[0], 64)
		if err != nil {
			continue
		}
		phonemes = append(phonemes, vals[2])
		times = append(times, pt)
	}

	var wg *sync.WaitGroup
	if sc.UntilDone {
		wg = &sync.WaitGroup{}
	}
	if sc.LipSync {
		if sc.SoundDelay > 0 {
			if wg != nil {
				wg.Add(1)
			}
			go func() {
				time.Sleep(time.Second * time.Duration(sc.SoundDelay))
				playSpeech()
				endSpeaking()
				if wg != nil {
					wg.Done()
				}
			}()
		} else {

			if wg != nil {
				wg.Add(2)
			}
			go func() {
				playSpeech()
				endSpeaking()
				if wg != nil {
					wg.Done()
				}
			}()
			//time.Sleep(time.Second * time.Duration(sc.SoundDelay))
			go func() {

				moveSpeech(phonemes, times)
				if wg != nil {
					wg.Done()
				}
			}()
		}
	} else {
		if wg != nil {
			wg.Add(1)
		}
		go func() {
			playSpeech()
			endSpeaking()
			if wg != nil {
				wg.Done()
			}
		}()
	}
	if wg != nil {
		wg.Wait()
	}
}

// Reset sets all the servos to default position and closes the connection
func Reset() {
	for _, m := range motors {
		if m.idx == int(TopLip) || m.idx == int(BottomLip) {
			continue
		}
		Move(MotorName(m.idx), 5)
		//time.Sleep(time.Millisecond * 100)
	}
	Move(MotorName(TopLip), 5)
	Move(MotorName(BottomLip), 5)
	time.Sleep(time.Millisecond * 250)
	DetachAll()
}

// Close detaches all the servos
func Close() {
	stopBlink(&blinkWg, blinkChan)
	DetachAll()
}

// DetachAll releases all the servos
func DetachAll() {
	for _, m := range motors {
		m.detach()
	}
}

// Wait pause execution for s seconds
func Wait(s float64) {
	time.Sleep(time.Millisecond * time.Duration(s*1000))
}

// BotVersion gets the OhBot Arduino code version
func BotVersion() ([]byte, error) {
	msg := "v\n"
	_, err := ser.Write([]byte(msg))
	if err != nil {
		return nil, err
	}
	buf := make([]byte, 256)
	n, err := ser.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

// ReadSensor reads the value at the indicated sensor
func ReadSensor(idx int) (float64, error) {
	msg := fmt.Sprintf("i0%d\n", idx)
	_, err := ser.Write([]byte(msg))
	if err != nil {
		return -1, err
	}
	buf := make([]byte, 256)
	n, err := ser.Read(buf)
	if err != nil {
		return -1, err
	}
	line := string(buf[:n])
	lines := strings.Split(line, ",")
	if len(lines) > 1 {
		indexIn := lines[0]
		indexIn = string(indexIn[1])
		intdex, err := strconv.Atoi(indexIn)
		if err != nil {
			return -1, err
		}
		newVal, err := strconv.Atoi(lines[1])
		if err != nil {
			return -1, err
		}
		sensors[intdex] = limit(float64(newVal*10) / 1024)
	}
	return sensors[idx], nil
}

// SensorValue returns the value read by the sensor indicated by the index
func SensorValue(idx int) float64 {
	return sensors[idx]
}

func Position(mn MotorName) float64 {
	m := motors[int(mn)]
	return m.Position()
}
