package termy

import (
	"fmt"
	"strconv"

	"github.com/mec-nyan/termy/byteme"
	"github.com/mec-nyan/termy/printer"
	"github.com/mec-nyan/termy/term"
)

const (
	// CSI (Control sequence introducer) is sent before most control codes.
	_csi = "\x1b["
	// ESC is used before some control codes that don't begin with the CSI.
	_esc = "\x1b"
)

const (
	// These flags lets us keep track of some internal state.
	altBuf uint = 1 << iota
	altCharSet
)

// Display takes care of handling your terminal and setting things up for your application.
// It allows you to set the terminal to the desired state (i.e. enable/disable line buffering)
// and controls in-band signalling for it.
// It also connects your program with the three file streams STDOUT, STDIN and STDERR.
type Display struct {
	term.Settings
	printer.Printer
	flags uint
}

// NewDisplay initialise a new Display structure with the default settings.
// For now I'm keeping that name "NewDisplay" instead of just "New" since it's
// more explicit about what it really does.
func NewDisplay() (*Display, error) {
	printer := printer.New()

	settings := term.New(int(printer.Stdout.Fd()))

	err := settings.Init()
	if err != nil {
		return nil, err
	}

	return &Display{
		Settings: *settings,
		Printer:  printer,
	}, nil
}

///////////////////////////////
// ** Cursor manipulation ** //
///////////////////////////////

// Home moves the cursor to the top left corner of the terminal.
func (d *Display) Home() {
	d.write(_csi + "H")
}

// Clear to end of line.
func (d *Display) ClearToEOL() {
	d.write(_csi + "K")
}

// Clear to the beginning of line.
func (d *Display) ClearToBOL() {
	d.write(_csi + "1K")
}

// Clear to end of screen.
func (d *Display) ClearToEOS() {
	d.write(_csi + "J")
}

// Clear the screen and move the cursor to the upper left corner.
func (d *Display) ClearScreen() {
	Home()
	ClearToEOS()
}

// Save the current cursor position.
func (d *Display) SaveCurPos() {
	d.write(_esc + "7")
}

// Restore the cursor position to a previously saved one.
func (d *Display) RestoreCurPos() {
	d.write(_esc + "8")
}

// Move cursor to column "col".
func (d *Display) CurToCol(col int) {
	d.write(_csi + strconv.Itoa(col) + "G")
}

// Move cursor to row "row".
func (d *Display) CurToRow(row int) {
	d.write(_csi + strconv.Itoa(row) + "dd")
}

// Move cursor up one row.
func (d *Display) Up() {
	d.write(_csi + "A")
}

// Move cursor down one row.
func (d *Display) Down() {
	d.write(_csi + "B")
}

// Move cursor one column to the right.
func (d *Display) Right() {
	d.write(_csi + "C")
}

// Move cursor one column to the left.
func (d *Display) Left() {
	d.write(_csi + "D")
}

// Move cursor "lines" rows up.
func (d *Display) MoveUp(lines int) {
	for i := 0; i < lines; i++ {
		Up()
	}
}

// Move cursor "lines" rows down.
func (d *Display) MoveDown(lines int) {
	for i := 0; i < lines; i++ {
		Down()
	}
}

// Move cursor "cols" columns to the right.
func (d *Display) MoveRight(cols int) {
	for i := 0; i < cols; i++ {
		Right()
	}
}

// Move cursor "cols" columns to the left.
func (d *Display) MoveLeft(cols int) {
	for i := 0; i < cols; i++ {
		Left()
	}
}

// Move cursor to line "y" col "x"
func (d *Display) MoveTo(x, y int) {
	d.write(_csi + strconv.Itoa(y) + ";" + strconv.Itoa(x) + "H")
}

// Make cursor invisible.
func (d *Display) HideCur() {
	d.write(_csi + "?25l")
}

// Make cursor visible.
func (d *Display) ShowCur() {
	d.write(_csi + "?25h")
}

// Enter alt buffer mode.
func (d *Display) EnterAltBuf() {
	if d.inAltBuf() {
		return
	}
	d.write(_csi + "?1049h")
	// Set flag to save state.
	d.flags |= altBuf
}

// Exit alt buffer mode.
func (d *Display) ExitAltBuf() {
	if d.inAltBuf() {
		d.write(_csi + "?1049l")
		// Clear flag.
		d.flags &^= altBuf
	}
}

// Enter alternate character set mode.
func (d *Display) EnterACS() {
	if d.inAltCharSet() {
		return
	}
	d.write(_esc + "(0")
	d.flags |= altCharSet
}

// Exit alternate character set mode.
func (d *Display) ExitACS() {
	if d.inAltCharSet() {
		d.write(_esc + "(B")
		d.flags &^= altCharSet
	}
}

// Delete character.
func (d *Display) DelChar() {
	d.write(_csi + "P")
}

// Delete line.
func (d *Display) DelLine() {
	d.write(_csi + "M")
}

// Inset line.
func (d *Display) InsLine() {
	d.write(_csi + "L")
}

// Colours.

// Some of the following routines work on a "best effort" basis
// and don't return error.

// Set the foreground colour using the terminal's theme.
// Assume 256 colours as it's available on most terminals.
// See Colour.SetFg()
func (d *Display) SetFg(colour int) *Display {
	d.Colour.SetFg(colour)
	return d
}

// Set the background colour using the terminal's theme.
// Assume 256 colours as it's available on most terminals.
// See Colour.SetBg()
func (d *Display) SetBg(colour int) *Display {
	d.Colour.SetBg(colour)
	return d
}

// Use default foreground colour.
func (d *Display) UseDefaultFg() *Display {
	d.Colour.UseDefaultFg()
	return d
}

// Use default background colour.
func (d *Display) UseDefaultBg() *Display {
	d.Colour.UseDefaultBg()
	return d
}

// Use default foreground and background colours.
func (d *Display) UseDefault() *Display {
	d.Colour.UseDefault()
	return d
}

// Reset foreground colour code.
// NOTE: this just empties the colour code, it does not
// restore the foreground colour!
// For that, use UseDefaultFg instead.
func (d *Display) ResetFg() *Display {
	d.Colour.ResetFg()
	return d
}

// Reset background colour code.
// NOTE: this just empties the colour code, it does not
// restore the background colour!
// For that, use UseDefaultBg instead.
func (d *Display) ResetBg() *Display {
	d.Colour.ResetBg()
	return d
}

// Reset foreground and background colour codes.
// NOTE: this just empties the colour code, it does not
// restore the foreground or background colours!
// For that, use UseDefault instead.
func (d *Display) Reset() *Display {
	d.Colour.Reset()
	return d
}

// SetFgRGB sets the foreground using rgb values.
// Each value should be in the range 0-255, otherwise
// the default fg is used.
func (d *Display) SetFgRGB(r, g, b int) *Display {
	d.Colour.SetFgRGB(r, g, b)
	return d
}

// SetBgRGB sets the background using rgb values.
// Each value should be in the range 0-255, otherwise
// the default bg is used.
func (d *Display) SerBgRGB(r, g, b int) *Display {
	d.Colour.SetBgRGB(r, g, b)
	return d
}

// SetFgHex sets the foreground to the hex colour provided.
// Valid hex codes are string in the "#RRGGBB" or "RRGGBB" format.
// If an invalid string/colour is given, the default fg is set instead.
func (d *Display) SetFgHex(colour string) *Display {
	d.Colour.SetFgHex(colour)
	return d
}

// SetBgHex sets the background to the hex colour provided.
// Valid hex codes are string in the "#RRGGBB" or "RRGGBB" format.
// If an invalid string/colour is given, the default bg is set instead.
func (d *Display) SetBgHex(colour string) *Display {
	d.Colour.SetBgHex(colour)
	return d
}

// Code generates the code for the currently selected colours and/or style.
// It doesn't prepend the CSI.
// NOTE: This function may not need to be exported.
func (d *Display) Code() string {
	colourCode := d.Colour.Code()
	styleCode := d.Style.Code()

	if len(colourCode) == 0 {
		return styleCode
	}

	if len(styleCode) == 0 {
		return colourCode
	}

	return styleCode + ";" + colourCode
}

// Send actually sends the in-band signal to the terminal/selected writer.
func (d *Display) Send() {
	d.Stdout.Write(byteme.UnsafeStrToBytes(d.escaped()))
}

// Text style.

// Reset all style attributes.
func (d *Display) Normal() *Display {
	d.Style.Normal()
	return d
}

// Set/unset bold attribute.
func (d *Display) Bold(on bool) *Display {
	if on {
		d.Style.Bold()
	} else {
		d.Style.NoBold()
	}
	return d
}

// Set/unset dim attribute.
func (d *Display) Dim(on bool) *Display {
	if on {
		d.Style.Dim()
	} else {
		d.Style.NoDim()
	}
	return d
}

// Set/unset italics attribute.
func (d *Display) Italics(on bool) *Display {
	if on {
		d.Style.Italics()
	} else {
		d.Style.NoItalics()
	}
	return d
}

// Set/unset blink attribute.
func (d *Display) Blink(on bool) *Display {
	if on {
		d.Style.Blink()
	} else {
		d.Style.NoBlink()
	}
	return d
}

// Set/unset reverse/standout attribute.
func (d *Display) Reverse(on bool) *Display {
	if on {
		d.Style.Reverse()
	} else {
		d.Style.NoReverse()
	}
	return d
}

// Set/unset hidden attribute.
func (d *Display) Hidden(on bool) *Display {
	if on {
		d.Style.Hidden()
	} else {
		d.Style.NoHidden()
	}
	return d
}

// Set/unset strikeout attribute.
func (d *Display) Strikeout(on bool) *Display {
	if on {
		d.Style.Strikeout()
	} else {
		d.Style.NoStrikeout()
	}
	return d
}

// Set/unset single underline attribute.
func (d *Display) Underline(on bool) *Display {
	if on {
		d.Style.Underline()
	} else {
		d.Style.NoUnderline()
	}
	return d
}

// Set/unset single underline attribute.
// This is an alias for Underline.
func (d *Display) SingleUnderline(on bool) *Display {
	return d.Underline(on)
}

// Set/unset double underline attribute.
// To turn it off, use d.Underline(false)
func (d *Display) DoubleUnderline() *Display {
	d.Style.UnderlineDouble()
	return d
}

// Set/unset curly underline attribute.
// To turn it off, use d.Underline(false)
func (d *Display) CurlyUnderline() *Display {
	d.Style.UnderlineCurly()
	return d
}

// Printing functions that can be accessed directly.

// PrintBytes prints out a slice of bytes.
func (d *Display) PrintBytes(b []byte) (int, error) {
	return d.Stdout.Write(b)
}

// PrintBytesAt prints a slice of byte at (x, y).
func (d *Display) PrintBytesAt(x, y int, b []byte) (int, error) {
	// TODO: Boundary check.
	d.MoveTo(x, y)
	return d.PrintBytes(b)
}

// TODO: PrintNBytes should count COLUMNS and not BYTES nor CHARACTERS.
func (d *Display) PrintNBytes(cols int, b []byte) (int, error) {
	return 0, nil
}

// TODO: PrintNBytes should count COLUMNS and not BYTES nor CHARACTERS.
func (d *Display) PrintNBytesAt(x, y, cols int, b []byte) (int, error) {
	return 0, nil
}

// Print prints a utf-8 encoded string.
func (d *Display) Print(s string) (int, error) {
	// Oh yes, Go is a memory safe language... or is it?
	return d.PrintBytes(byteme.UnsafeStrToBytes(s))

}

// PrintAt prints a utf-8 encoded string at (x, y).
func (d *Display) PrintAt(x, y int, s string) (int, error) {
	return d.PrintBytesAt(x, y, byteme.UnsafeStrToBytes(s))
}

func (d *Display) CurPos() (int, int) {
	// We don't want our terminal to print that to stdout.
	if d.Echoing() {
		d.NoEcho()
		defer d.Echo()
	}
	// Not quite sure if we need this, but it seems like it.
	if d.Cooked() {
		d.UnCookIt()
		defer d.CookIt()
	}

	x := 0
	y := 0
	// TODO: There has to be a better way using std[out|in] directly.
	// It will work for now.
	fmt.Print("\x1b[6n")
	fmt.Scanf("\x1b[%d;%dR", &y, &x)

	return x, y
}

// Internal.

// write is a wrapper for Stdout.Write.
func (d *Display) write(s string) {
	d.Stdout.Write(byteme.UnsafeStrToBytes(s))
}

// escaped converts the colour and style sequence in an in-band command.
// prepending the CSI and appending a terminator string.
func (d *Display) escaped() string {
	code := d.Code()
	if len(code) > 0 {
		return "\x1b[" + code + "m"
	}

	return ""
}

// Flags handling.

// Check if the display is in alternate buffer mode.
// NOTE: This is internal. It will check the state saved by our application.
// It will NOT check your emulator state directly.
// If the terminal was put in alt buf mode by any other way (i.e. you print the codes directly)
// this will not reflect that. (BTW don't do that! Use the methods provided).
func (d *Display) inAltBuf() bool {
	f := d.flags & altBuf
	return f == altBuf
}

// Check if the display is in alternate character set mode.
// NOTE: This is internal. It will check the state saved by our application.
// It will NOT check your emulator state directly.
// If the terminal was put in ACS mode by any other way (i.e. you print the codes directly)
// this will not reflect that. (BTW don't do that! Use the methods provided).
func (d *Display) inAltCharSet() bool {
	f := d.flags & altCharSet
	return f == altCharSet
}
