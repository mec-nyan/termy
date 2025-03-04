package main

import (
	"fmt"
	"os"

	"github.com/mec-nyan/termy"
)

func main() {
	fd := int(os.Stdout.Fd())
	ts := termy.New(fd, false)

	_ = ts.Cbreaky()
	defer ts.Restore()

	// Get terminal size.
	rows, cols, _ := ts.Size()

	// Handle terminal colour and style.
	term := termy.NewTermy(os.Stdout)

	// Save cursor position:
	term.SaveCurPos()
	term.HideCur()
	term.Italics().Blink()
	term.SetFgRGB(55, 255, 10).SetBgHex("#607080")
	term.Send()

	fmt.Printf("Hello from Termy!!")
	// You can use the pkg global funcs.
	termy.CurToCol(1)
	termy.MoveDown(2)
	fmt.Printf("(%d x %d)", rows, cols)
	term.UseDefault()
	term.Normal().Dim().Italics()
	term.Send()
	// Or the methods.
	term.CurToCol(1)
	term.MoveDown(4)
	fmt.Printf("Press any key to continue...")

	buffer := make([]byte, 1)
	os.Stdin.Read(buffer)

	term.RestoreCurPos()
	term.ClearToEOS()
	term.ShowCur()
}
