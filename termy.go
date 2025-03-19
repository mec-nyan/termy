package termy

import (
	"io"
	"os"
	"strconv"

	"github.com/mec-nyan/termy/colours"
	"github.com/mec-nyan/termy/styles"
)

const (
	_csi = "\x1b["
	_esc = "\x1b"
)

// Termy struct handles in-band colour and style commands for its tty.
// tty would normally be os.Stdout.
type Termy struct {
	colours.Colour
	styles.Style
	tty io.Writer
	// Experimental! 
	// Access terminal settings/features through Termy
	// TODO: If this works out fine, maybe rename "Termy" to something like "Screeny"
	// and TermSettings to "Termy" 🤔
	TermSettings
}

// NewTermy sets up a new Termy struct to handle in-band signalling to the selected io.Writer.
func NewTermy(w io.Writer) *Termy {
	stdout, ok := w.(*os.File)
	if !ok {
		stdout = os.Stdout
	}

	return &Termy{
		Colour: colours.Colour{},
		Style:  styles.Style{},
		tty:    w,
		// The defaults are OK for the other fields.
		TermSettings: TermSettings{
			fd: int(stdout.Fd()),
		},
	}
}

// Code generates the code for the currently selected colours and/or style.
// It doesn't prepend the CSI.
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

// Escaped converts the colour and style sequence in an in-band command.
// prepending the CSI and appending a terminator string.
func (t *Termy) Escaped() string {
	code := t.Code()
	if len(code) > 0 {
		return "\x1b[" + code + "m"
	}

	return ""
}

// Send actually sends the in-band signal to the terminal/selected writer.
func (t *Termy) Send() {
	t.tty.Write([]byte(t.Escaped()))
}

// Cursor manipulation:
//
// Home moves the cursor to the top left corner of the terminal.
func (t *Termy) Home() {
	t.write(_csi + "H")
}

// Clear to end of line.
func (t *Termy) ClearToEOL() {
	t.write(_csi + "K")
}

// Clear to the beginning of line.
func (t *Termy) ClearToBOL() {
	t.write(_csi + "1K")
}

// Clear to end of screen.
func (t *Termy) ClearToEOS() {
	t.write(_csi + "J")
}

// Clear the screen and move the cursor to the upper left corner.
func (t *Termy) ClearScreen() {
	Home()
	ClearToEOS()
}

// Save the current cursor position.
func (t *Termy) SaveCurPos() {
	t.write(_esc + "7")
}

// Restore the cursor position to a previously saved one.
func (t *Termy) RestoreCurPos() {
	t.write(_esc + "8")
}

// Move cursor to column "col".
func (t *Termy) CurToCol(col int) {
	t.write(_csi + strconv.Itoa(col) + "G")
}

// Move cursor to row "row".
func (t *Termy) CurToRow(row int) {
	t.write(_csi + strconv.Itoa(row) + "dd")
}

// Move cursor up one row.
func (t *Termy) Up() {
	t.write(_csi + "A")
}

// Move cursor down one row.
func (t *Termy) Down() {
	t.write(_csi + "B")
}

// Move cursor one column to the right.
func (t *Termy) Right() {
	t.write(_csi + "C")
}

// Move cursor one column to the left.
func (t *Termy) Left() {
	t.write(_csi + "D")
}

// Move cursor "lines" rows up.
func (t *Termy) MoveUp(lines int) {
	for i := 0; i < lines; i++ {
		Up()
	}
}

// Move cursor to line "y" col "x"
func (t *Termy) MoveTo(x, y int) {
	t.write(_csi + strconv.Itoa(y) + ";" + strconv.Itoa(x) + "H")
}

// Move cursor "lines" rows down.
func (t *Termy) MoveDown(lines int) {
	for i := 0; i < lines; i++ {
		Down()
	}
}

// Move cursor "cols" columns to the right.
func (t *Termy) MoveRight(cols int) {
	for i := 0; i < cols; i++ {
		Right()
	}
}

// Move cursor "cols" columns to the left.
func (t *Termy) MoveLeft(cols int) {
	for i := 0; i < cols; i++ {
		Left()
	}
}

// Make cursor invisible.
func (t *Termy) HideCur() {
	t.write(_csi + "?25l")
}

// Make cursor visible.
func (t *Termy) ShowCur() {
	t.write(_csi + "?25h")
}

// Enter alt buffer mode.
func (t *Termy) EnterCaMode() {
	t.write(_csi + "?1049h")
}

// Exit alt buffer mode.
func (t *Termy) ExitCaMode() {
	t.write(_csi + "?1049l")
}

// Internal.

func (t *Termy) write(s string) {
	t.tty.Write([]byte(s))
}
