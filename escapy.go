package termy

import (
	"os"
	"strconv"
)

// TODO: Should we use bytes/runes instead of strings?
const (
	esc = "\x1b"
	csi = esc + "["
)

func write(s string) {
	// TODO: Should we use bytes/runes instead of strings?
	// TODO: Should we check for errors?
	os.Stdout.Write([]byte(s))
}

func Home() {
	write(csi + "H")
}

func ClearToEOL() {
	write(csi + "K")
}

func ClearToBOL() {
	write(csi + "1K")
}

func ClearToEOS() {
	write(csi + "J")
}

func ClearScreen() {
	Home()
	ClearToEOS()
}

func SaveCurPos() {
	write(esc + "7")
}

func RestoreCurPos() {
	write(esc + "8")
}

func CurToCol(col int) {
	write(csi + strconv.Itoa(col) + "G")
}

func CurToRow(row int) {
	write(csi + strconv.Itoa(row) + "dd")
}

func Up() {
	write(csi + "A")
}

func Down() {
	write(csi + "B")
}

func Right() {
	write(csi + "C")
}

func Left() {
	write(csi + "C")
}

func MoveUp(lines int) {
	for i := 0; i < lines; i++ {
		Up()
	}
}

func MoveDown(lines int) {
	for i := 0; i < lines; i++ {
		Down()
	}
}

func MoveRight(lines int) {
	for i := 0; i < lines; i++ {
		Right()
	}
}

func MoveLeft(lines int) {
	for i := 0; i < lines; i++ {
		Left()
	}
}

func HideCur() {
	write(csi + "?25l")
}

func ShowCur() {
	write(csi + "?25h")
}

func EnterCaMode() {
	write(csi + "?1049h")
}

func ExitCaMode() {
	write(csi + "?1049l")
}
