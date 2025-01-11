// Package colours generate the colour codes to be embedded in a colour escape sequence.
package colours

import (
	"fmt"
	"strconv"
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

// Colour holds the sequences for fg and bg.
// It doesn't add the CSI nor the terminator, so you can combine theses codes
// with other (i.e. with bold, underline, etc).
type Colour struct {
	// fg and bg are private, so we don't overwrite them by mistake.
	fg, bg string
}

// Fg returns the currently set fg sequence.
func (c *Colour) Fg() string {
	return c.fg
}

// Bg returns the currently set bg sequence.
func (c *Colour) Bg() string {
	return c.bg
}

// Code return the sequence for setting the foreground and background for the current Colour's state.
func (c *Colour) Code() string {
	// If both fg and bg are empty, it will return an empty string, which is valid.
	if c.fg == "" {
		return c.bg
	}
	if c.bg == "" {
		return c.fg
	}
	return c.fg + ";" + c.bg
}

// Reset empty the sequences. It will produce no changes beyond cleaning its state.
// If you wish to use the default colours, use UseDefault instead.
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

// SetFgHex sets the foreground to the hex colour provided.
// If an invalid string/colour is given, the default fg is set instead.
func (c *Colour) SetFgHex(colour string) {
	c.SetFgRGB(fromHex(colour))
}

// SetBgHex sets the background to the hex colour provided.
// If an invalid string/colour is given, the default bg is set instead.
func (c *Colour) SetBgHex(colour string) {
	c.SetBgRGB(fromHex(colour))
}

// fromHex is a wrapper around hexToRGB.
// It handles the error so you can use it where you just need the rgb values.
// It will try to parse a string containing an hex colour.
// That string should be in the format "#RRGGBB" or "RRGGBB".
// If it fails, it will return an invalid rgb colour so the defaults are used instead.
func fromHex(hex string) (r, g, b int) {
	r, g, b, err := hexToRGB(hex)
	if err != nil {
		return -1, -1, -1
	}
	return
}

// Internal.

// hexToRGB tries to parse a string containing an hex colour.
func hexToRGB(colour string) (r, g, b int, err error) {
	colour, err = getHexColour(colour)
	if err != nil {
		return
	}

	// If we've passed that check, colour contains a valid hex number,
	// and ParseInt will never fail.

	rgb, err := strconv.ParseInt(colour, 16, 0)
	if err != nil {
		return
	}

	b = int(rgb % 0x100)
	rgb /= 0x100

	g = int(rgb % 0x100)
	rgb /= 0x100

	r = int(rgb)

	return
}

// getHexColour extracts the hex number from a string in the format "#RRGGBB" or "RRGGBB".
// TODO: Is the validation of ParseInt enough? (We wouldn't need isHexNum and isHexDigit).
// It isn't. We need to validate that these is a colour (0xRRGGBB) and not just an hex.
// i.e. ParseInt won't give an error for "F000" (les than six digits) but we should.
func getHexColour(str string) (hex string, err error) {
	// Remove leading '#', if any.
	if str[0] == '#' {
		str = str[1:]
	}
	// We only accept "#RRGGBB" or "RRGGBB" colours.
	if len(str) != 6 {
		return "", fmt.Errorf("'%s' is not a valid hex colour", str)
	}
	// It must be a valid hexadecimal number. Negatives are not allowed.
	if !isHexNum(str) {
		return "", fmt.Errorf("'%s' is not a valid hex number", str)
	}

	return str, nil
}

func isHexNum(str string) bool {
	for _, r := range str {
		if !isHexDigit(r) {
			return false
		}
	}
	return true
}

func isHexDigit(h rune) bool {
	hexDigits := "0123456789ABCDEFabcdef"
	for _, r := range hexDigits {
		if r == h {
			return true
		}
	}
	return false
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
