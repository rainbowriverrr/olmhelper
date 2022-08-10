package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var startButton *widget.Button

type boss struct {
	timeBetween    time.Duration
	attackSequence []string
}

var olm = boss{
	timeBetween: time.Millisecond * 600,
	attackSequence: []string{
		"nauto",
		"null",
		"sauto",
		"crystal",
		"nauto",
		"null",
		"sauto",
		"lightning",
		"nauto",
		"null",
		"sauto",
		"portals",
	},
}

//cycle() cycles through olm.attackSequence and displays the current attack, then sleeps for olm.timeBetween indefinitely.
//the cycle restarts at the beginning when a restart button is pressed.
func cycle(stop chan string, restart chan string, text *widget.Label) {

	attackIndex := 0
	n := 4
	lastTime := time.Now().Add(-olm.timeBetween)

	for {
		//if channel is empty, change text to the next attack in the sequence and sleep for timeBetween
		select {
		case <-stop:
			text.SetText("Stopped")
			return
		case <-restart:
			attackIndex = 0
			n = 4
			text.SetText(fmt.Sprintf("%s:%d", olm.attackSequence[attackIndex], n))
			lastTime = time.Now()
			n--

		default:
			//checks if 600 milliseconds has passed
			if time.Since(lastTime) > olm.timeBetween {
				text.SetText(fmt.Sprintf("%s:%d", olm.attackSequence[attackIndex], n))
				lastTime = time.Now()
				n--
				//checks if n is 0 if so, reset to 4
				if n == 0 {
					n = 4
					attackIndex++
					//checks if attackIndex is equal to the length of the attackSequence if so, reset to 0
					if attackIndex == len(olm.attackSequence) {
						attackIndex = 0
					}
				}
			}
		}
		//sleeps for a few seconds for system to breathe
		time.Sleep(time.Millisecond * 10)
	}
}

func main() {

	//creates a restart channel
	restart := make(chan string)
	//creates a stop channel
	stop := make(chan string)
	//creates a text label to display the current attack
	main := widget.NewLabel("Nothing Yet")

	a := app.New()
	w := a.NewWindow("Olm Helper")

	restartButton := widget.NewButton("Restart", func() {
		restart <- "restart"
	})
	stopButton := widget.NewButton("Stop", func() {
		stop <- "stop"
		w.SetContent(container.NewGridWithRows(2,
			container.NewCenter(main),
			startButton,
		))
	})
	startButton = widget.NewButton("Start", func() {
		//replaces startButton with restartButton
		w.SetContent(container.NewGridWithRows(3,
			container.NewCenter(main),
			restartButton,
			stopButton))
		//starts the cycle function
		go cycle(stop, restart, main)
	})
	w.SetContent(container.NewGridWithRows(2,
		container.NewCenter(main),
		startButton,
	))
	w.Resize(fyne.NewSize(200, 200))
	w.ShowAndRun()
}
