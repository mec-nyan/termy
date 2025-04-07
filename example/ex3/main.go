package main

import (
	"fmt"
	"os"

	"github.com/mec-nyan/termy"
	"github.com/mec-nyan/termy/printer"
)

func main() {
	screen, err := termy.NewDisplay()
	if err != nil {
		panic(err)
	}

	screen.UnCookIt()
	screen.NoEcho()
	defer screen.Restore()

	screen.EnterAltBuf()
	defer screen.ExitAltBuf()

	lines, _, _ := screen.Size()
	screen.SetFg(8).Send()
	for i := 1; i <= lines; i++ {
		screen.MoveTo(1, i)
		screen.Print(fmt.Sprintf("%4d│", i))
	}

	printer := printer.New()

	screen.MoveTo(7, 1)

	printer.SetFg(8).Print("Fu fu fu some text...")

	screen.MoveTo(7, 3)

	printer.SetFg(63).Print("Hello you there!")

	screen.MoveTo(7, 5)

	printer.SetFg(212).CurlyUnderline().Print("I'm curly!")

	screen.MoveTo(7, 8)

	printer.Normal().SetFg(8).Italics(true).Print("Press any key ")
	getChar()
}

func getChar() byte {
	buf := make([]byte, 1)
	os.Stdin.Read(buf)
	return buf[0]
}
