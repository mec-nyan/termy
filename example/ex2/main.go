package main

import (
	"fmt"
	"os"

	"github.com/mec-nyan/termy"
)

func main() {
	// Handle terminal colour and style.
	screen, err := termy.NewDisplay(os.Stdout)
	if err != nil {
		panic(err)
	}
	screen.UnCookIt()
	screen.NoEcho()

	// Or through termy 😜
	rows, cols, _ := screen.Size()

	defer screen.Restore()

	// Save cursor position:
	screen.SaveCurPos()
	screen.HideCur()

	screen.Italics(true).
		Blink(true).
		SetFgRGB(20, 20, 20).
		SetBgHex("#FF8040").
		Send()

	fmt.Printf("New Termy now integrates term settings!!")
	// You can use the pkg global funcs...
	termy.CurToCol(1)
	termy.MoveDown(2)

	fmt.Printf("(%d x %d)", rows, cols)

	screen.UseDefault().
		Normal().
		Dim(true).
		Italics(true).
		Send()

	// ... or the methods.
	screen.CurToCol(1)
	screen.MoveDown(4)
	fmt.Printf("Press any key to continue...")

	getc()

	fmt.Print("\nEnter your name: ")
	screen.Echo()

Loop1:
	for {
		b := getc()
		switch b {
		case '\n':
			break Loop1
		default:
		}
	}

	fmt.Print("\nEnter your password: ")
	screen.NoEcho()

Loop2:
	for {
		b := getc()
		switch b {
		case '\n':
			break Loop2
		default:
			fmt.Print("*")
		}
	}

	screen.RestoreCurPos()
	screen.ClearToEOS()
	screen.ShowCur()
}

func getc() byte {
	buffer := make([]byte, 1)
	os.Stdin.Read(buffer)
	return buffer[0]
}
