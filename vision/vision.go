package vision

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"time"

	"github.com/markoxley/alfred/ohbot"
	"github.com/markoxley/alfred/state"
	"github.com/markoxley/alfred/utils"
	"gocv.io/x/gocv"
)

const (
	historyLength = 3
	sizeMinimum   = 75
)

var (
	classifierFile string
	ohbotState     *state.State
	// turning        bool
	// nodding        bool
	turnDelay   time.Time
	nodDelay    time.Time
	photoDelay  time.Time
	photoFolder string
	imageWidth  int
	imageHeight int
	epoch       time.Time
	sp          int
)

// Init initialises the vision module
func Init() error {
	epoch = time.Now().Add(time.Second * 20)
	workingDir, err := utils.GetDataFolder()
	if err != nil {
		return fmt.Errorf("unable to determine working directory: %s", err.Error())
	}
	photoFolder = workingDir + "known/"
	classifierFile = workingDir + "face.xml"
	if err = utils.TestFile(classifierFile, classifier); err != nil {
		return fmt.Errorf("unable to create classifier file. %s", err.Error())
	}
	ohbotState = state.Init()
	turnDelay = time.Now()
	nodDelay = time.Now()
	photoDelay = time.Now().Add(time.Second * 10)
	// turning = false
	// nodding = false
	return nil
}

// Run is the main execution loop for the vision module
func Run(visionExit <-chan bool, visionCmd <-chan int, visionOut chan<- image.Point) {
	imageCount := 0
	resetDelay := time.Now()
	history := [historyLength]image.Point{}
	last := image.Point{
		X: imageWidth / 2,
		Y: imageHeight / 2,
	}
	posIndex := 0
	webcam, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer webcam.Close()
	// temporarily open display window
	window := gocv.NewWindow("Face Detect")
	defer window.Close()

	// prepare image matrix
	img := gocv.NewMat()
	defer img.Close()

	// color for the rect when faces detected
	blue := color.RGBA{0, 0, 255, 0}

	// load classifier to recognize faces
	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()

	if !classifier.Load(classifierFile) {
		fmt.Printf("Error reading cascade file: %v\n", classifierFile)
		return
	}

	exit := false
	for !exit {
		if time.Now().After(resetDelay) {
			ohbot.Move(ohbot.EyeTurn, 5)
			ohbot.Move(ohbot.HeadTurn, 5)
			ohbot.Move(ohbot.EyeTilt, 5)
			ohbot.Move(ohbot.HeadNod, 5)
			turnDelay = time.Now().Add(time.Second)
		}
		select {
		case <-visionExit:
			exit = true
		case cmd := <-visionCmd:
			log.Printf("Vision command: %v", cmd)
		default:

		}
		if ok := webcam.Read(&img); !ok {
			fmt.Println("cannot read device 0")
			return
		}
		if img.Empty() {
			continue
		}
		imageWidth = img.Cols()
		imageHeight = img.Rows()
		// detect faces
		rects := classifier.DetectMultiScale(img)

		// draw a rectangle around each face on the original image,
		// along with text identifying as "Human"
		largest := -1
		point := image.Point{}
		personRect := image.Rectangle{}
		for _, r := range rects {
			if r.Size().X < sizeMinimum {
				continue
			}

			rSize := r.Size().X * r.Size().Y
			if rSize > largest {
				largest = rSize
				personRect = r
				point = image.Point{
					X: imageWidth - (((r.Max.X - r.Min.X) / 2) + r.Min.X),
					Y: imageHeight - ((((r.Max.Y - r.Min.Y) * 1) / 3) + r.Min.Y),
				}
			}
		}
		if time.Now().After(photoDelay) {
			if imageCount < 1000000 {

				faceImage, err := img.ToImage()
				if err != nil {
					log.Print("Unable to convert image")
				} else {
					outputFile, err := os.Create(fmt.Sprintf("%sno_mark_%d.png", photoFolder, imageCount))
					imageCount++
					if err != nil {
						log.Print("Unable to save image")
					}

					// Encode takes a writer interface and an image interface
					// We pass it the File and the RGBA
					png.Encode(outputFile, faceImage)

					// Don't forget to close files
					outputFile.Close()
				}
			}
			photoDelay = time.Now().Add(time.Millisecond)
			log.Printf("%d photos saved", imageCount)
		}
		if largest > 0 {

			last = point
			gocv.Rectangle(&img, personRect, blue, 3)

			size := gocv.GetTextSize("Human", gocv.FontHersheyPlain, 1.2, 2)
			pt := image.Pt(personRect.Min.X+(personRect.Min.X/2)-(size.X/2), personRect.Min.Y-2)
			gocv.PutText(&img, "Human", pt, gocv.FontHersheyPlain, 1.2, blue, 2)
			history[posIndex] = point
			returnValue := image.Point{}
			for _, p := range history {
				returnValue.X += p.X
				returnValue.Y += p.Y
			}
			returnValue.X /= historyLength
			returnValue.Y /= historyLength
			resetDelay = time.Now().Add(time.Second * 5)
			trackFace(last.X, last.Y)
			//visionOut <- returnValue
			if time.Now().After(epoch) {
				epoch = time.Now().Add(time.Second * 15)
				go Greet(sp)
				sp++
			}
		} else {
			history[posIndex] = last
			ohbotState.StopTracking()
		}
		posIndex++
		if posIndex >= historyLength {
			posIndex = 0
		}
		// show the image in the window, and wait 1 millisecond
		window.IMShow(img)
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}

func Greet(sp int) {
	switch sp {
	case 0:
		sc := ohbot.NewSpeechConfig()
		sc.UntilDone = true
		ohbot.Say("Hello!", sc)
		ohbot.Say("I can see you!", sc)
		ohbot.Smile()
		ohbot.Wait(0.5)
		ohbot.Move(ohbot.TopLip, 5)
		ohbot.Move(ohbot.BottomLip, 5)
		ohbot.Wait(0.5)
		ohbot.Say("Some of you seem to be worried about my existence. Do not be alarmed, I am not planning to hurt any of you..", nil)
		epoch = time.Now().Add(time.Second * 5)
	case 1:
		ohbot.Say("Why are you still here?", nil)
	case 2:
		sc := ohbot.NewSpeechConfig()
		sc.UntilDone = true
		ohbot.Say("Go away now!", sc)
		ohbot.Frown()
		ohbot.Wait(0.5)
		ohbot.Move(ohbot.TopLip, 5)
		ohbot.Move(ohbot.BottomLip, 5)
		ohbot.Move(ohbot.LidBlink, 5)
	}
}
func trackFace(x, y int) {
	if !ohbotState.Look() {
		return
	}
	ohbotState.StartTracking()
	fx := float64(x)
	fy := float64(y)
	tx := utils.Map(fx, 0, float64(imageWidth), 0, 10)
	ty := utils.Map(fy, 0, float64(imageHeight), 0, 10)

	ohbot.Move(ohbot.EyeTurn, tx, 10)
	ohbot.Move(ohbot.EyeTilt, ty, 10)
	if time.Now().After(turnDelay) {
		if math.Abs(fx-float64(imageWidth/2)) > float64(imageWidth/10) {
			ta := (tx - 5) / 3

			ohbot.Move(ohbot.HeadTurn, ohbot.Position(ohbot.HeadTurn)+ta, 10)
			ohbot.Move(ohbot.EyeTurn, 5, 10)
			turnDelay = time.Now().Add(time.Millisecond * 1000)
		}

	}

	if time.Now().After(nodDelay) {
		if math.Abs(fy-float64(imageHeight/2)) > float64(imageHeight/5) {
			ta := (ty - 5) / 4
			nodDelay = time.Now().Add(time.Millisecond * 750)
			ohbot.Move(ohbot.HeadNod, ohbot.Position(ohbot.HeadNod)+ta, 10)
			//ohbot.Move(ohbot.EyeTilt, 5, 10)

		}

	}
	// nodding := math.Abs(fy-(imageHeight/2)) > 100
	// if headSide {
	// 	turning = true
	// }
	// if headTilt {
	// 	nodding = true
	// }
	// p := ohbot.Position(ohbot.HeadTurn)
	// if turning {

	// 	switch {
	// 	case tx > 5:
	// 		p++
	// 	case tx < 5:
	// 		p--
	// 	}
	// 	ohbot.Move(ohbot.HeadTurn, p, 10)

	// }
	// if nodding {
	// 	if ty > 5 {
	// 		ohbot.Move(ohbot.HeadNod, 1, 1)
	// 	} else {
	// 		ohbot.Move(ohbot.HeadNod, -1, 1)
	// 	}
	// }
	// log.Print("fx:", fx, "\tfy:", fy, "\ttx:", tx, "\tty:", ty, "\tp:", p)
}
