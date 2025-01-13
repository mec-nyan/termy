package termy

import (
	"io"
	"strconv"

	"github.com/mec-nyan/termy/colours"
	"github.com/mec-nyan/termy/styles"
)

type Termy struct {
	colours.Colour
	styles.Style
	tty io.Writer
}

func NewTermy(w io.Writer) *Termy {
	return &Termy{
		Colour: colours.Colour{},
		Style:  styles.Style{},
		tty:    w,
	}
}

func (t *Termy) Code() string {
	colourCode := t.Colour.Code()
	styleCode := t.Style.Code()

	if len(colourCode) == 0 {
		return styleCode
	}

	if len(styleCode) == 0 {
		return colourCode
	}

	return styleCode + ";" + colourCode
}

func (t *Termy) Escaped() string {
	code := t.Code()
	if len(code) > 0 {
		return "\x1b[" + code + "m"
	}

	return ""
}

func (t *Termy) Send() {
	t.tty.Write([]byte(t.Escaped()))
}

func (t *Termy) write(s string) {
	t.tty.Write([]byte(s))
}

func (t *Termy) Home() {
	t.write(csi + "H")
}

func (t *Termy) ClearToEOL() {
	t.write(csi + "K")
}

func (t *Termy) ClearToBOL() {
	t.write(csi + "1K")
}

func (t *Termy) ClearToEOS() {
	t.write(csi + "J")
}

func (t *Termy) ClearScreen() {
	Home()
	ClearToEOS()
}

func (t *Termy) SaveCurPos() {
	t.write(esc + "7")
}

func (t *Termy) RestoreCurPos() {
	t.write(esc + "8")
}

func (t *Termy) CurToCol(col int) {
	t.write(csi + strconv.Itoa(col) + "G")
}

func (t *Termy) CurToRow(row int) {
	t.write(csi + strconv.Itoa(row) + "dd")
}

func (t *Termy) Up() {
	t.write(csi + "A")
}

func (t *Termy) Down() {
	t.write(csi + "B")
}

func (t *Termy) Right() {
	t.write(csi + "C")
}

func (t *Termy) Left() {
	t.write(csi + "C")
}

func (t *Termy) MoveUp(lines int) {
	for i := 0; i < lines; i++ {
		Up()
	}
}

func (t *Termy) MoveDown(lines int) {
	for i := 0; i < lines; i++ {
		Down()
	}
}

func (t *Termy) MoveRight(lines int) {
	for i := 0; i < lines; i++ {
		Right()
	}
}

func (t *Termy) MoveLeft(lines int) {
	for i := 0; i < lines; i++ {
		Left()
	}
}

func (t *Termy) HideCur() {
	t.write(csi + "?25l")
}

func (t *Termy) ShowCur() {
	t.write(csi + "?25h")
}

func (t *Termy) EnterCaMode() {
	t.write(csi + "?1049h")
}

func (t *Termy) ExitCaMode() {
	t.write(csi + "?1049l")
}
