// Package colours generate the colour codes to be embedded in a colour escape sequence.
package colours

import (
	"fmt"
)

const (
	Black int = iota
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
	BrightBlack
	BrightRed
	BrightGreen
	BrightYellow
	BrightBlue
	BrightMagenta
	BrightCyan
	BrightWhite
)

// Colour holds the sequences for fg and bg using 256 colours.
// It doesn't add the CSI nor the terminator, so you can combine theses codes
// with other (i.e. with bold, underline, etc).
type Colour struct {
	fg, bg string
}

// Reset empty the sequences. It will produce no changes at all.
// If you wish to use the default colours, use UseDefault.
func (c *Colour) Reset() {
	c.ResetFg()
	c.ResetBg()
}

// ResetFg cleans the fg sequence. It doesn't set the default fg.
// If you want to use the default fg use UseDefaultFg instead.
func (c *Colour) ResetFg() {
	c.fg = ""
}

// ResetBg cleans the bg sequence. It doesn't set the default bg.
// If you want to use the default bg use UseDefaultBg instead.
func (c *Colour) ResetBg() {
	c.bg = ""
}

// UseDefault sets the default fg and bg.
func (c *Colour) UseDefault() {
	c.UseDefaultFg()
	c.UseDefaultBg()
}

// UseDefaultFg sets the default fg.
func (c *Colour) UseDefaultFg() {
	c.fg = "39"
}

// UseDefaultBg sets the default bg.
func (c *Colour) UseDefaultBg() {
	c.bg = "49"
}

// SetFg sets the foreground colour in the 0-255 range.
// If you pass an invalid value (n<0 or n>255) it will fallback on the default fg.
func (c *Colour) SetFg(colour int) {
	// Invalid argument? Don't panic! Just use the default colour.
	if !in255range(colour) {
		c.UseDefaultFg()
		return
	}
	c.fg = fmt.Sprintf("38:5:%d", colour)
}

// SetBg sets the background colour in the 0-255 range.
// If you pass an invalid value (n<0 or n>255) it will fallback on the default bg.
func (c *Colour) SetBg(colour int) {
	// Invalid argument? Don't panic! Just use the default colour.
	if !in255range(colour) {
		c.UseDefaultBg()
		return
	}
	c.bg = fmt.Sprintf("48:5:%d", colour)
}

// SetFgRGB sets the foreground using rgb values.
// Each value should be in the range 0-255, otherwise the default fg is used.
func (c *Colour) SetFgRGB(r, g, b int) {
	if !isValidRGB(r, g, b) {
		c.UseDefaultFg()
		return
	}
	c.fg = fmt.Sprintf("38:2:%d:%d:%d", r, g, b)
}

// SetBgRGB sets the background using rgb values.
// Each value should be in the range 0-255, otherwise the default bg is used.
func (c *Colour) SetBgRGB(r, g, b int) {
	if !isValidRGB(r, g, b) {
		c.UseDefaultBg()
		return
	}
	c.bg = fmt.Sprintf("48:2:%d:%d:%d", r, g, b)
}

func isValidRGB(r, g, b int) bool {
	if !in255range(r) || !in255range(g) || !in255range(b) {
		return false
	}
	return true
}

func in255range(n int) bool {
	if n < 0 || n > 255 {
		return false
	}
	return true
}
