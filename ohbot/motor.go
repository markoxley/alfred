package ohbot

import (
	"fmt"
	"strconv"

	"github.com/beevik/etree"
)

// motor is a struct to hold details for a specified servo
type motor struct {
	pos, min, max, rest float64
	rev, attached       bool
	mType               string
	idx                 int
}

var (
	motors            []*motor
	ohbotMotorDefFile string
)

// newMotor creates a new Motor object with default settings
func newMotor() *motor {
	return &motor{
		pos:      11,
		min:      0,
		max:      0,
		rest:     0,
		rev:      false,
		attached: false,
		mType:    "",
		idx:      0,
	}
}

func loadMotorDefs() error {
	tree := etree.NewDocument()
	if err := tree.ReadFromFile(ohbotMotorDefFile); err != nil {
		return fmt.Errorf(fmt.Sprintf("Unable to read motor definitions file: %s", err.Error()))
	}

	root := tree.SelectElement("Motors")
	for _, element := range root.SelectElements("Motor") {
		indexStr := element.SelectAttrValue("Motor", "")
		index, err := strconv.Atoi(indexStr)
		if err != nil {
			return fmt.Errorf(fmt.Sprintf("Error reading motor min def: %s", err.Error()))
		}
		min, err := strconv.ParseFloat(element.SelectAttrValue("Min", ""), 64)
		if err != nil {
			return fmt.Errorf(fmt.Sprintf("Error reading motor min def: %s", err.Error()))
		}
		max, err := strconv.ParseFloat(element.SelectAttrValue("Max", ""), 64)
		if err != nil {
			return fmt.Errorf(fmt.Sprintf("Error reading motor max def: %s", err.Error()))
		}
		rest, err := strconv.ParseFloat(element.SelectAttrValue("RestPosition", ""), 64)
		if err != nil {
			return fmt.Errorf(fmt.Sprintf("Error reading motor rest def: %s", err.Error()))
		}
		rev := element.SelectAttrValue("Reverse", "") == "True"

		motors[index].min = min / 1000 * 180
		motors[index].max = max / 1000 * 180
		motors[index].rest = rest
		motors[index].pos = rest
		motors[index].rev = rev
		motors[index].idx = index
	}
	return nil
}

func (m *motor) attach() {
	if m.attached {
		return
	}
	msg := fmt.Sprintf("a0%d\n", m.idx)
	serWrite(msg)
	m.attached = true
}

func (m *motor) detach() {
	msg := fmt.Sprintf("d0%d\n", m.idx)
	serWrite(msg)
	m.attached = false
}

func (m *motor) absPos(p float64) uint16 {
	mr := m.max - m.min
	sp := (mr / 10) * p
	return uint16(sp + m.min)
}

// Position returns the current position of the motor
func (m *motor) Position() float64 {
	return m.pos
}
