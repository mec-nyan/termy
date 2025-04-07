package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mec-nyan/termy"
)

func main() {
	// Handle terminal colour and style.
	screen, err := termy.NewDisplay()
	if err != nil {
		panic(err)
	}
	defer clean(screen)

	setup(screen)

	showInfo(screen)
	showLove(screen)

	lookupState(screen)

	getName(screen)
	getPassword(screen)

	showCurPos(screen)

	bye(screen)
}

func setup(screen *termy.Display) {
	// Disable line buffering.
	screen.UnCookIt()
	// Do not output user input straight away.
	screen.NoEcho()

	// Save cursor position: We'll clear the screen and leave everything as it was
	// before our app started.
	screen.SaveCurPos()
	screen.HideCur()
}

func showInfo(screen *termy.Display) {
	rows, cols, _ := screen.Size()

	screen.Italics(true).
		Blink(true).
		SetFgRGB(20, 20, 20).
		SetBgHex("#FF8040").
		Send()

	screen.Print("New Termy now integrates term settings!!")
	// You can use the pkg global funcs...
	termy.CurToCol(1)
	termy.MoveDown(2)

	screen.Print(fmt.Sprintf("(%d x %d)", rows, cols))
}

func showLove(screen *termy.Display) {
	screen.UseDefault().
		Normal().
		Dim(true).
		Italics(true).
		Send()

	// ... or the methods.
	screen.CurToCol(1)
	screen.MoveDown(4)
	screen.Print("This is some text for u 🩷")

	getc()
}

func lookupState(screen *termy.Display) {
	var msg string
	if screen.Echoing() {
		msg = "\n\tI'm repeating everything you type!"
	} else {
		msg = "\n\tYour text disappears!"
	}
	screen.Print(msg)

	if screen.Cooked() {
		msg = "\n\tI'm in my normal state"
	} else {
		msg = "\n\tRaw raw raw!!!"
	}
	screen.Print(msg)
}

func getName(screen *termy.Display) {
	screen.Print("\nEnter your name: ")
	screen.Echo()

	for {
		b := getc()
		switch b {
		case '\n':
			return
		default:
		}
	}
}

func getPassword(screen *termy.Display) {
	screen.Print("\nEnter your password: ")
	screen.NoEcho()

Loop:
	for {
		b := getc()
		switch b {
		case '\n':
			break Loop
		default:
			screen.Print("*")
		}
	}
	screen.Print("\n")
}

func showCurPos(screen *termy.Display) {
	x, y := screen.CurPos()

	screen.Print(fmt.Sprintf("\n\nCur pos? x: %d, y: %d", x, y))

	getc()
}

func bye(screen *termy.Display) {
	screen.SetFg(4).Italics(true).Send()
	screen.Print("\n\nLater mate!")
	time.Sleep(1 * time.Second)
}

func getc() byte {
	buffer := make([]byte, 1)
	os.Stdin.Read(buffer)
	return buffer[0]
}

func clean(screen *termy.Display) {
	screen.RestoreCurPos()
	screen.ClearToEOS()
	screen.Restore()
	screen.UseDefault()
	screen.ShowCur()
}
