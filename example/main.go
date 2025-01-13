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

	// Handle terminal colour and style.
	term := termy.NewTermy(os.Stdout)

	// Save cursor position:
	term.SaveCurPos()
	term.HideCur()
	term.Italics().Blink()
	term.SetFgRGB(55, 255, 10).SetBgHex("#607080")
	term.Send()

	fmt.Printf("Hello from Termy!!\n\n\n")
	term.UseDefault()
	term.Normal().Dim().Italics()
	term.Send()
	fmt.Printf("Press any key to continue...")

	buffer := make([]byte, 1)
	os.Stdin.Read(buffer)

	term.RestoreCurPos()
	term.ClearToEOS()
	term.ShowCur()
}
