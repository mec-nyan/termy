package main

import (
	"fmt"
	"os"

	"github.com/mec-nyan/termy"
)

func main() {
	ts := termy.NewTerm(int(os.Stdout.Fd()))

	_ = ts.UnCookIt()
	defer ts.Restore()

	// Get terminal size.
	rows, cols, _ := ts.Size()

	// Handle terminal colour and style.
	display, err := termy.NewDisplay(os.Stdout)
	if err != nil {
		panic(err)
	}

	// Or through termy 😜
	rows, cols, _ = display.Size()

	// Save cursor position:
	display.SaveCurPos()
	display.HideCur()
	display.Italics().Blink()
	display.SetFgRGB(55, 255, 10).SetBgHex("#607080")
	display.Send()

	fmt.Printf("Hello from Termy!!")
	// You can use the pkg global funcs.
	termy.CurToCol(1)
	termy.MoveDown(2)
	fmt.Printf("(%d x %d)", rows, cols)
	display.UseDefault()
	display.Normal().Dim().Italics()
	display.Send()
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
