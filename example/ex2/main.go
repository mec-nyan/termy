package main

import (
	"fmt"
	"os"

	"github.com/mec-nyan/termy"
)

func main() {
	// Handle terminal colour and style.
	screen := termy.NewTermy(os.Stdout)
	_ = screen.Cbreaky()

	// Or through termy 😜
	rows, cols, _ := screen.Size()

	defer screen.Restore()


	// Save cursor position:
	screen.SaveCurPos()
	screen.HideCur()
	screen.Italics().Blink()
	screen.SetFgRGB(20, 20, 20).SetBgHex("#FF8040")
	screen.Send()

	fmt.Printf("New Termy now integrates term settings!!")
	// You can use the pkg global funcs.
	termy.CurToCol(1)
	termy.MoveDown(2)
	fmt.Printf("(%d x %d)", rows, cols)
	screen.UseDefault()
	screen.Normal().Dim().Italics()
	screen.Send()
	// Or the methods.
	screen.CurToCol(1)
	screen.MoveDown(4)
	fmt.Printf("Press any key to continue...")

	buffer := make([]byte, 1)
	os.Stdin.Read(buffer)

	screen.RestoreCurPos()
	screen.ClearToEOS()
	screen.ShowCur()
}
