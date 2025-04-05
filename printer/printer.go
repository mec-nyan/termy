/*
Package printe tries to separete the process of setting colours,
styles and printing to the terminal from the display.

TODO: In this stage a lot of the functionality will be duplicated.
Some of it may be useful in the Display object for more simple cases.
*/
package printer

import (
	"io"
	"unsafe"

	"github.com/mec-nyan/termy/colour"
	"github.com/mec-nyan/termy/style"
)

// Printer prints to the writer with selected colour and style.
type Printer struct {
	colour.Colour
	style.Style
	writer io.Writer
}

func New(w io.Writer) *Printer {
	return &Printer{writer: w}
}

// Set the foreground colour using the terminal's theme.
// Assume 256 colours as it's available on most terminals.
// See Colour.SetFg()
func (p *Printer) SetFg(colour int) *Printer {
	p.Colour.SetFg(colour)
	return p
}

// Set the background colour using the terminal's theme.
// Assume 256 colours as it's available on most terminals.
// See Colour.SetBg()
func (p *Printer) SetBg(colour int) *Printer {
	p.Colour.SetBg(colour)
	return p
}

// Use default foreground colour.
func (p *Printer) UseDefaultFg() *Printer {
	p.Colour.UseDefaultFg()
	return p
}

// Use default background colour.
func (p *Printer) UseDefaultBg() *Printer {
	p.Colour.UseDefaultBg()
	return p
}

// Use default foreground and background colours.
func (p *Printer) UseDefault() *Printer {
	p.Colour.UseDefault()
	return p
}

// Reset foreground colour code.
// NOTE: this just empties the colour code, it does not
// restore the foreground colour!
// For that, use UseDefaultFg instead.
func (p *Printer) ResetFg() *Printer {
	p.Colour.ResetFg()
	return p
}

// Reset background colour code.
// NOTE: this just empties the colour code, it does not
// restore the background colour!
// For that, use UseDefaultBg instead.
func (p *Printer) ResetBg() *Printer {
	p.Colour.ResetBg()
	return p
}

// Reset foreground and background colour codes.
// NOTE: this just empties the colour code, it does not
// restore the foreground or background colours!
// For that, use UseDefault instead.
func (p *Printer) Reset() *Printer {
	p.Colour.Reset()
	return p
}

// SetFgRGB sets the foreground using rgb values.
// Each value should be in the range 0-255, otherwise
// the default fg is used.
func (p *Printer) SetFgRGB(r, g, b int) *Printer {
	p.Colour.SetFgRGB(r, g, b)
	return p
}

// SetBgRGB sets the background using rgb values.
// Each value should be in the range 0-255, otherwise
// the default bg is used.
func (p *Printer) SerBgRGB(r, g, b int) *Printer {
	p.Colour.SetBgRGB(r, g, b)
	return p
}

// SetFgHex sets the foreground to the hex colour provided.
// Valid hex codes are string in the "#RRGGBB" or "RRGGBB" format.
// If an invalid string/colour is given, the default fg is set instead.
func (p *Printer) SetFgHex(colour string) *Printer {
	p.Colour.SetFgHex(colour)
	return p
}

// SetBgHex sets the background to the hex colour provided.
// Valid hex codes are string in the "#RRGGBB" or "RRGGBB" format.
// If an invalid string/colour is given, the default bg is set instead.
func (p *Printer) SetBgHex(colour string) *Printer {
	p.Colour.SetBgHex(colour)
	return p
}

// Text style.

// Reset all style attributes.
func (p *Printer) Normal() *Printer {
	p.Style.Normal()
	return p
}

// Set/unset bold attribute.
func (p *Printer) Bold(on bool) *Printer {
	if on {
		p.Style.Bold()
	} else {
		p.Style.NoBold()
	}
	return p
}

// Set/unset dim attribute.
func (p *Printer) Dim(on bool) *Printer {
	if on {
		p.Style.Dim()
	} else {
		p.Style.NoDim()
	}
	return p
}

// Set/unset italics attribute.
func (p *Printer) Italics(on bool) *Printer {
	if on {
		p.Style.Italics()
	} else {
		p.Style.NoItalics()
	}
	return p
}

// Set/unset blink attribute.
func (p *Printer) Blink(on bool) *Printer {
	if on {
		p.Style.Blink()
	} else {
		p.Style.NoBlink()
	}
	return p
}

// Set/unset reverse/standout attribute.
func (p *Printer) Reverse(on bool) *Printer {
	if on {
		p.Style.Reverse()
	} else {
		p.Style.NoReverse()
	}
	return p
}

// Set/unset hidden attribute.
func (p *Printer) Hidden(on bool) *Printer {
	if on {
		p.Style.Hidden()
	} else {
		p.Style.NoHidden()
	}
	return p
}

// Set/unset strikeout attribute.
func (p *Printer) Strikeout(on bool) *Printer {
	if on {
		p.Style.Strikeout()
	} else {
		p.Style.NoStrikeout()
	}
	return p
}

// Set/unset single underline attribute.
func (p *Printer) Underline(on bool) *Printer {
	if on {
		p.Style.Underline()
	} else {
		p.Style.NoUnderline()
	}
	return p
}

// Set/unset single underline attribute.
// This is an alias for Underline.
func (p *Printer) SingleUnderline(on bool) *Printer {
	return p.Underline(on)
}

// Set/unset double underline attribute.
// To turn it off, use d.Underline(false)
func (p *Printer) DoubleUnderline() *Printer {
	p.Style.UnderlineDouble()
	return p
}

// Set/unset curly underline attribute.
// To turn it off, use d.Underline(false)
func (p *Printer) CurlyUnderline() *Printer {
	p.Style.UnderlineCurly()
	return p
}

// Code generates the code for the currently selected colours and/or style.
// It doesn't prepend the CSI.
func (p *Printer) Code() string {
	colourCode := p.Colour.Code()
	styleCode := p.Style.Code()

	if len(colourCode) == 0 {
		return styleCode
	}

	if len(styleCode) == 0 {
		return colourCode
	}

	return styleCode + ";" + colourCode
}

// escaped converts the colour and style sequence in an in-band command.
// prepending the CSI and appending a terminator string.
func (p *Printer) escaped() string {
	code := p.Code()
	if len(code) > 0 {
		return "\x1b[" + code + "m"
	}

	return ""
}

// Send actually sends the in-band signal to the terminal/selected writer.
func (p *Printer) Send() {
	p.writer.Write(unsafeStrToBytes(p.escaped()))
}

// PrintBytes prints out a slice of bytes with the printer style.
func (p *Printer) PrintBytes(b []byte) (int, error) {
	p.Send()
	// Should we clear at the end?
	b = append(b, unsafeStrToBytes("\x1b[0m")...)
	return p.writer.Write(b)
}

// Print prints a utf-8 encoded string.
func (p *Printer) Print(s string) (int, error) {
	// Oh yes, Go is a memory safe language... or is it?
	return p.PrintBytes(unsafeStrToBytes(s))
}

// -------- Internal -------- //

// Memory safe language my a**
// We can use []byte(s) but it seems to be, lets say, not ideal...
func unsafeStrToBytes(s string) []byte {
	if s == "" {
		return nil
	}
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
