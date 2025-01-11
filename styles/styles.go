// Package styles serves to generate the sequence to signal in-band styles to the terminal.
// It doesn't return escape sequences, so you can combine these i.e. with colours.
package styles

import "strconv"

type Attribute int

// Unexported: These consts are used within this package only.
const (
	// TODO: Check if "normal" is really necessary.
	// normal Attribute = iota
	_ Attribute = iota
	bold
	dim
	italics
	underline
	double
	curly
	blink
	blinkFast
	reverse
	hidden
	strikeout
)

// Style holds the selected styles.
// Since all different styles can be combined with each others, we keep them as boolean values
// and update the sequence when we call Code.
type Style struct {
	sequence string

	// TODO: Check if "normal" is really necessary.
	normal    bool
	bold      bool
	dim       bool
	italics   bool
	underline bool // Single underline.
	double    bool // Double underline.
	curly     bool // Curly underline.
	blink     bool
	reverse   bool
	hidden    bool
	strikeout bool
}

// Code gets the attributes sequence.
func (s *Style) Code() string {
	s.update()
	return s.sequence
}

// Reset turns every attribute off.
func (s *Style) Reset() *Style {
	s.sequence = ""
	s.normal = true
	s.bold = false
	s.dim = false
	s.italics = false
	s.underline = false
	s.double = false
	s.curly = false
	s.blink = false
	s.reverse = false
	s.hidden = false
	s.strikeout = false
	return s
}

// Normal is an alias for Reset.
// It turns every attribute off.
func (s *Style) Normal() *Style {
	return s.Reset()
}

// Set bold attribute.
func (s *Style) Bold() *Style {
	return s.setAttr(bold, true)
}

// Remove bold attribute.
func (s *Style) NoBold() *Style {
	return s.setAttr(bold, false)
}

// Set dim attribute.
func (s *Style) Dim() *Style {
	return s.setAttr(dim, true)
}

// Remove dim attribute.
func (s *Style) NoDim() *Style {
	return s.setAttr(dim, false)
}

// Set italics attribute.
func (s *Style) Italics() *Style {
	return s.setAttr(italics, true)
}

// Remove italics attribute.
func (s *Style) NoItalics() *Style {
	return s.setAttr(italics, false)
}

// Set underline attribute.
func (s *Style) Underline() *Style {
	return s.setAttr(underline, true)
}

// Remove underline attribute (including single, double and curly).
func (s *Style) NoUnderline() *Style {
	return s.setAttr(underline, false)
}

// Set single underline attribute (Same as "Underline").
func (s *Style) UnderlineSingle() *Style {
	return s.Underline()
}

// Set double underline attribute.
func (s *Style) UnderlineDouble() *Style {
	return s.setAttr(double, true)
}

// Set curly underline attribute.
func (s *Style) UnderlineCurly() *Style {
	return s.setAttr(curly, true)
}

// Set blink attribute.
func (s *Style) Blink() *Style {
	return s.setAttr(blink, true)
}

// Remove blink attribute.
func (s *Style) NoBlink() *Style {
	return s.setAttr(blink, false)
}

// Set reverse or standout attribute.
func (s *Style) Reverse() *Style {
	return s.setAttr(reverse, true)
}

// Remove reverse or standout attribute.
func (s *Style) NoReverse() *Style {
	return s.setAttr(reverse, false)
}

// Set hidden attribute.
func (s *Style) Hidden() *Style {
	return s.setAttr(hidden, true)
}

// Remove hidden attribute.
func (s *Style) NoHidden() *Style {
	return s.setAttr(hidden, false)
}

// Set strikeout attribute.
func (s *Style) Strikeout() *Style {
	return s.setAttr(strikeout, true)
}

// Remove strikeout attribute.
func (s *Style) NoStrikeout() *Style {
	return s.setAttr(strikeout, false)
}

// Internal.
//
// Should we export setAttr? Would it be useful beyond this package?
func (s *Style) setAttr(attr Attribute, on bool) *Style {
	switch attr {
	// case normal:
	// 	s.normal = on
	case bold:
		s.bold = on
	case dim:
		s.dim = on
	case italics:
		s.italics = on
	// underline is a special case.
	case underline:
		// Set single underline or none.
		s.underline = on
		s.double = false
		s.curly = false
	case double:
		// Set double underline or none.
		s.underline = on
		s.double = on
		s.curly = false
	case curly:
		// Set curly underline or none.
		s.underline = on
		s.double = false
		s.curly = on
	case blink, blinkFast:
		s.blink = on
	case reverse:
		s.reverse = on
	case hidden:
		s.hidden = on
	case strikeout:
		s.strikeout = on
	}
	return s
}

func (s *Style) update() {
	// Attributes can combine with each others.
	// TODO: Check if "normal" is really necessary.
	attrs := []int{}
	ul_type := 0
	if s.bold {
		attrs = append(attrs, 1)
	}
	if s.dim {
		attrs = append(attrs, 2)
	}
	if s.italics {
		attrs = append(attrs, 3)
	}
	if s.underline {
		attrs = append(attrs, 4)
		ul_type = 1
	}
	if s.double {
		ul_type = 2
	}
	if s.curly {
		ul_type = 3
	}
	if s.blink {
		attrs = append(attrs, 5)
	}
	if s.reverse {
		attrs = append(attrs, 7)
	}
	if s.hidden {
		attrs = append(attrs, 8)
	}
	if s.strikeout {
		attrs = append(attrs, 9)
	}

	if len(attrs) == 0 {
		s.sequence = "0"
		return
	}

	s.sequence = ""
	for i, n := range attrs {
		s.sequence += strconv.Itoa(n)
		// Handle special underline case.
		if n == 4 {
			s.sequence += ":" + strconv.Itoa(ul_type)
		}
		// If it's not the last element, add a separator (";").
		if i < len(attrs)-1 {
			s.sequence += ";"
		}
	}
}
