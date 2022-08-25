package ohbot

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

var (
	monthNames []string = []string{
		"January",
		"February",
		"March",
		"April",
		"May",
		"June",
		"July",
		"August",
		"September",
		"October",
		"November",
		"December",
	}

	exact         []string = []string{"", "exactly", "precisely"}
	approximate   []string = []string{"", "roughly", "about", "approximately"}
	morningText   []string = []string{"A. M.", "in the morning"}
	afternoonText []string = []string{"P. M.", "in the afternoon"}
	eveningText   []string = []string{"P. M.", "in the evening"}
	dateText      []string = []string{
		"Today is %s",
		"The date today is %s",
		"It is %s today",
		"Today's date is %s",
	}
	timeSentenceText []string = []string{
		"It is %s",
		"The time is %s",
		"It is now %s",
	}
	dayText []string = []string{
		"It is %s",
		"Today is %s",
		"The day today is %s",
	}
	timeText []string = []string{
		"o clock",
		"five past",
		"ten past",
		"quarter past",
		"twenty five past",
		"twenty five past",
		"half past",
		"twenty five to",
		"twenty to",
		"quarter to",
		"ten to",
		"five to",
	}
)

func SayDate() {
	dt := time.Now()
	y := dt.Year()
	d := dt.Day()
	m := monthNames[int(dt.Month())-1]
	var ts string
	if rand.Int()%2 == 0 {
		ts = fmt.Sprintf("%s %d %d", m, d, y)
	} else {
		ts = fmt.Sprintf("the %d of %s %d", d, m, y)
	}
	ts = fmt.Sprintf(dateText[rand.Int()%len(dateText)], ts)
	Say(ts, nil)

}

func SayTime() {

	tm := time.Now()
	h := tm.Hour()
	m := tm.Minute()
	appm := int(math.Round(float64(m / 5)))
	if m%5 > 2 {
		appm++
		if appm == len(timeText) {
			appm = 0
		}
	}
	mt := timeText[appm]

	var dayText string
	if h >= 12 {
		if h > 17 {
			dayText = eveningText[rand.Int()%len(eveningText)]
		} else {
			dayText = afternoonText[rand.Int()%len(afternoonText)]
		}
		h -= 12

	} else {
		dayText = morningText[rand.Int()%len(morningText)]
	}
	if appm > 7 {
		h++
	}
	if h == 0 {
		h = 12
	}
	var precision string
	if m%5 == 0 {
		precision = exact[rand.Int()%len(exact)]
	} else {
		precision = approximate[rand.Int()%len(approximate)]
	}
	var text string

	if m == 0 {
		text = fmt.Sprintf("%s %d %s %s", precision, h, mt, dayText)
	} else {
		text = fmt.Sprintf("%s %s %d, %s", precision, mt, h, dayText)
	}

	text = fmt.Sprintf(timeSentenceText[rand.Int()%len(timeSentenceText)], text)
	Say(text, nil)
}

func SayDay() {
	dt := time.Now()
	d := dt.Weekday().String()
	text := fmt.Sprintf(dayText[rand.Int()%len(dayText)], d)
	Say(text, nil)
}
