package colours

import "testing"

func TestColours(t *testing.T) {
	// Given
	cases := []struct {
		name   string
		want   string
		action func(*Colour) string
	}{
		{
			name: "Emtpy fg",
			want: "",
			action: func(c *Colour) string {
				return c.Fg()
			},
		},
		{
			name: "Emtpy bg",
			want: "",
			action: func(c *Colour) string {
				return c.Bg()
			},
		},
		{
			name: "Emtpy seq",
			want: "",
			action: func(c *Colour) string {
				return c.Code()
			},
		},
		{
			name: "Use default fg",
			want: "39",
			action: func(c *Colour) string {
				c.UseDefaultFg()
				return c.Code()
			},
		},
		{
			name: "Use default bg",
			want: "49",
			action: func(c *Colour) string {
				c.UseDefaultBg()
				return c.Code()
			},
		},
		{
			name: "Use default fg and bg",
			want: "39;49",
			action: func(c *Colour) string {
				c.UseDefault()
				return c.Code()
			},
		},
		{
			name: "Reset fg",
			want: "49",
			action: func(c *Colour) string {
				c.UseDefault()
				c.ResetFg()
				return c.Code()
			},
		},
		{
			name: "Reset bg",
			want: "39",
			action: func(c *Colour) string {
				c.UseDefault()
				c.ResetBg()
				return c.Code()
			},
		},
		{
			name: "Reset fg and bg",
			want: "",
			action: func(c *Colour) string {
				c.UseDefault()
				c.Reset()
				return c.Code()
			},
		},
		// SetFg and SetBg.
		{
			name: "Set fg (int)",
			want: "38:5:1",
			action: func(c *Colour) string {
				c.SetFg(Red)
				return c.Code()
			},
		},
		{
			name: "Set fg (int)(invalid)",
			want: "39",
			action: func(c *Colour) string {
				c.SetFg(256)
				return c.Code()
			},
		},
		{
			name: "Set bg (int)",
			want: "48:5:1",
			action: func(c *Colour) string {
				c.SetBg(Red)
				return c.Code()
			},
		},
		{
			name: "Set fg (int)(invalid)",
			want: "49",
			action: func(c *Colour) string {
				c.SetBg(-1)
				return c.Code()
			},
		},
		// SetFgRGB and SetBgRGB.
		{
			name: "Set fg (rgb)",
			want: "38:2:253:254:255",
			action: func(c *Colour) string {
				c.SetFgRGB(253, 254, 255)
				return c.Code()
			},
		},
		{
			name: "Set fg (rgb)(invalid)",
			want: "39",
			action: func(c *Colour) string {
				c.SetFgRGB(256, 0, 0)
				return c.Code()
			},
		},
		{
			name: "Set bg (rgb)",
			want: "48:2:127:128:129",
			action: func(c *Colour) string {
				c.SetBgRGB(127, 128, 129)
				return c.Code()
			},
		},
		{
			name: "Set fg (int)(invalid)",
			want: "49",
			action: func(c *Colour) string {
				c.SetBgRGB(-1, 0, 0)
				return c.Code()
			},
		},
		// SetFgHex and SetBgHex.
		{
			name: "Set fg (hex)",
			want: "38:2:255:254:253",
			action: func(c *Colour) string {
				c.SetFgHex("#FFFEFD")
				return c.Code()
			},
		},
		{
			name: "Set fg (hex)(no #)",
			want: "38:2:255:254:253",
			action: func(c *Colour) string {
				c.SetFgHex("FFFEFD")
				return c.Code()
			},
		},
		{
			name: "Set fg (hex)(invalid)",
			want: "39",
			action: func(c *Colour) string {
				c.SetFgHex("#-xFFEF")
				return c.Code()
			},
		},
		{
			name: "Set bg (hex)",
			want: "48:2:127:128:129",
			action: func(c *Colour) string {
				c.SetBgHex("#7F8081")
				return c.Code()
			},
		},
		{
			name: "Set bg (hex)(no #)",
			want: "48:2:127:128:129",
			action: func(c *Colour) string {
				c.SetBgHex("7F8081")
				return c.Code()
			},
		},
		{
			name: "Set fg (hex)(invalid)",
			want: "49",
			action: func(c *Colour) string {
				c.SetBgHex("-FFFFFF")
				return c.Code()
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			colour := Colour{}
			got := c.action(&colour)
			if got != c.want {
				t.Errorf("err: want (%s), got (%s)", c.want, got)
			}
		})
	}
}
