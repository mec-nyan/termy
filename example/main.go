package main

import (
	"fmt"
	"os"

	"github.com/mec-nyan/termy"
)

func main() {

	// Handle terminal colour and style.
	display, err := termy.NewDisplay()
	if err != nil {
		panic(err)
	}
	defer display.Restore()

	display.UnCookIt()
	display.NoEcho()

	// Get terminal size.
	rows, cols, _ := display.Size()

	// Or through termy 😜
	rows, cols, _ = display.Size()

	// Save cursor position:
	display.SaveCurPos()
	display.HideCur()

	display.Italics(true).
		Blink(true).
		SetFgRGB(55, 255, 10).
		SetBgHex("#607080").
		Send()

	fmt.Printf("Hello from Termy!!")
	// You can use the pkg global funcs.
	termy.CurToCol(1)
	termy.MoveDown(2)
	fmt.Printf("(%d x %d)", rows, cols)

	display.UseDefault().
		Normal().
		Dim(true).
		Italics(true).
		Send()

	// Or the methods.
	display.CurToCol(1)
	display.MoveDown(4)

	fmt.Printf("Press any key to continue...")

	buffer := make([]byte, 1)
	os.Stdin.Read(buffer)

	display.RestoreCurPos()
	display.ClearToEOS()
	display.ShowCur()
}
