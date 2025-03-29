package termy

import (
	"errors"
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

// Display struct handles in-band colour and style commands for its tty.
// tty would normally be os.Stdout.
type Display struct {
	colours.Colour
	styles.Style
	tty io.Writer
	// Experimental!
	// Access terminal settings/features through Termy
	// TODO: If this works out fine, maybe rename "Termy" to something like "Screeny"
	// and Settings to "Termy" 🤔
	Settings
}

// NewDisplay sets up a new Display struct to handle in-band signalling to the selected io.Writer.
func NewDisplay(w io.Writer) (*Display, error) {
	stdout, ok := w.(*os.File)
	if !ok {
		return nil, errors.New("Failed to initialise display.")
	}

	return &Display{
		Colour: colours.Colour{},
		Style:  styles.Style{},
		tty:    w,
		Settings: Settings{
			// The defaults are OK for the other fields.
			fd: int(stdout.Fd()),
		},
	}, nil
}

// Code generates the code for the currently selected colours and/or style.
// It doesn't prepend the CSI.
func (t *Display) Code() string {
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
func (t *Display) Escaped() string {
	code := t.Code()
	if len(code) > 0 {
		return "\x1b[" + code + "m"
	}

	return ""
}

// Send actually sends the in-band signal to the terminal/selected writer.
func (t *Display) Send() {
	t.tty.Write([]byte(t.Escaped()))
}

// Cursor manipulation:
//
// Home moves the cursor to the top left corner of the terminal.
func (t *Display) Home() {
	t.write(_csi + "H")
}

// Clear to end of line.
func (t *Display) ClearToEOL() {
	t.write(_csi + "K")
}

// Clear to the beginning of line.
func (t *Display) ClearToBOL() {
	t.write(_csi + "1K")
}

// Clear to end of screen.
func (t *Display) ClearToEOS() {
	t.write(_csi + "J")
}

// Clear the screen and move the cursor to the upper left corner.
func (t *Display) ClearScreen() {
	Home()
	ClearToEOS()
}

// Save the current cursor position.
func (t *Display) SaveCurPos() {
	t.write(_esc + "7")
}

// Restore the cursor position to a previously saved one.
func (t *Display) RestoreCurPos() {
	t.write(_esc + "8")
}

// Move cursor to column "col".
func (t *Display) CurToCol(col int) {
	t.write(_csi + strconv.Itoa(col) + "G")
}

// Move cursor to row "row".
func (t *Display) CurToRow(row int) {
	t.write(_csi + strconv.Itoa(row) + "dd")
}

// Move cursor up one row.
func (t *Display) Up() {
	t.write(_csi + "A")
}

// Move cursor down one row.
func (t *Display) Down() {
	t.write(_csi + "B")
}

// Move cursor one column to the right.
func (t *Display) Right() {
	t.write(_csi + "C")
}

// Move cursor one column to the left.
func (t *Display) Left() {
	t.write(_csi + "D")
}

// Move cursor "lines" rows up.
func (t *Display) MoveUp(lines int) {
	for i := 0; i < lines; i++ {
		Up()
	}
}

// Move cursor to line "y" col "x"
func (t *Display) MoveTo(x, y int) {
	t.write(_csi + strconv.Itoa(y) + ";" + strconv.Itoa(x) + "H")
}

// Move cursor "lines" rows down.
func (t *Display) MoveDown(lines int) {
	for i := 0; i < lines; i++ {
		Down()
	}
}

// Move cursor "cols" columns to the right.
func (t *Display) MoveRight(cols int) {
	for i := 0; i < cols; i++ {
		Right()
	}
}

// Move cursor "cols" columns to the left.
func (t *Display) MoveLeft(cols int) {
	for i := 0; i < cols; i++ {
		Left()
	}
}

// Make cursor invisible.
func (t *Display) HideCur() {
	t.write(_csi + "?25l")
}

// Make cursor visible.
func (t *Display) ShowCur() {
	t.write(_csi + "?25h")
}

// Enter alt buffer mode.
func (t *Display) EnterCaMode() {
	t.write(_csi + "?1049h")
}

// Exit alt buffer mode.
func (t *Display) ExitCaMode() {
	t.write(_csi + "?1049l")
}

// Internal.

func (t *Display) write(s string) {
	t.tty.Write([]byte(s))
}
