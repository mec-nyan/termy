// escapy.go
// These are standalone versions of the cursor manipulation functions.
// They write directly to os.Stdout.
package termy

import (
	"os"
	"strconv"
)

func Home() {
	writeBytes(csi('H')...)
}

func ClearToEOL() {
	writeBytes(csi('K')...)
}

func ClearToBOL() {
	writeBytes(csi('1', 'K')...)
}

func ClearToEOS() {
	writeBytes(csi('J')...)
}

func ClearScreen() {
	Home()
	ClearToEOS()
}

func SaveCurPos() {
	writeBytes(escape('7')...)
}

func RestoreCurPos() {
	writeBytes(escape('8')...)
}

func CurToCol(col int) {
	code := append(intToBytes(col), 'G')
	writeBytes(csi([]byte(code)...)...)
}

func CurToRow(row int) {
	code := append(intToBytes(row), []byte{'d', 'd'}...)
	writeBytes(csi([]byte(code)...)...)
}

func Up() {
	writeBytes(csi('A')...)
}

func Down() {
	writeBytes(csi('B')...)
}

func Right() {
	writeBytes(csi('c')...)
}

func Left() {
	writeBytes(csi('D')...)
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

func MoveTo(row, col int) {
	writeBytes(csi([]byte(strconv.Itoa(row) + ";" + strconv.Itoa(col) + "H")...)...)
}

func HideCur() {
	writeBytes(csi('?', '2', '5', 'l')...)
}

func ShowCur() {
	writeBytes(csi('?', '2', '5', 'h')...)
}

func EnterCaMode() {
	writeBytes(csi('?', '1', '0', '9', 'h')...)
}

func ExitCaMode() {
	writeBytes(csi('?', '1', '0', '4', '9', 'l')...)
}

// Internal.

func write(s string) {
	// TODO: Should we use bytes/runes instead of strings?
	// TODO: Should we check for errors?
	os.Stdout.Write([]byte(s))
}

func writeBytes(b ...byte) {
	os.Stdout.Write(b)
}

func escape(b ...byte) []byte {
	return append([]byte{'\x1b'}, b...)
}

func csi(b ...byte) []byte {
	return append([]byte{'\x1b', '['}, b...)
}

func intToBytes(n int) []byte {
	out := []byte{}

	for n > 0 {
		x := n % 10
		y := byte('0' + x)
		out = append(out, y)
		n /= 10
	}

	reverse(out)

	return out
}

func reverse(runes []byte) {
	var middle int = len(runes) / 2

	for i := 0; i < middle; i++ {
		runes[i], runes[len(runes)-(1+i)] = runes[len(runes)-(1+i)], runes[i]
	}
}
