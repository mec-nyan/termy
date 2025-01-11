package styles

import "testing"

func TestCode(t *testing.T) {
	// Given
	cases := []struct {
		name    string
		want    string
		actions func(s *Style)
	}{
		{
			name:    "None",
			want:    "0",
			actions: func(s *Style) {},
		},
		{
			name: "Bold",
			want: "1",
			actions: func(s *Style) {
				s.Bold()
			},
		},
		{
			name: "Dim",
			want: "2",
			actions: func(s *Style) {
				s.Dim()
			},
		},
		{
			name: "Italics",
			want: "3",
			actions: func(s *Style) {
				s.Italics()
			},
		},
		{
			name: "Underline",
			want: "4:1",
			actions: func(s *Style) {
				s.Underline()
			},
		},
		{
			name: "UnderlineDouble",
			want: "4:2",
			actions: func(s *Style) {
				s.UnderlineDouble()
			},
		},
		{
			name: "UnderlineCurly",
			want: "4:3",
			actions: func(s *Style) {
				s.UnderlineCurly()
			},
		},
		{
			name: "Blink",
			want: "5",
			actions: func(s *Style) {
				s.Blink()
			},
		},
		{
			name: "Reverse",
			want: "7",
			actions: func(s *Style) {
				s.Reverse()
			},
		},
		{
			name: "Hidden",
			want: "8",
			actions: func(s *Style) {
				s.Hidden()
			},
		},
		{
			name: "Strikeout",
			want: "9",
			actions: func(s *Style) {
				s.Strikeout()
			},
		},
		{
			name: "All (Single)",
			want: "1;2;3;4:1;5;7;8;9",
			actions: func(s *Style) {
				s.Bold()
				s.Dim()
				s.Italics()
				s.UnderlineSingle()
				s.Blink()
				s.Reverse()
				s.Hidden()
				s.Strikeout()
			},
		},
		{
			name: "All (Double)",
			want: "1;2;3;4:2;5;7;8;9",
			actions: func(s *Style) {
				s.Bold()
				s.Dim()
				s.Italics()
				s.UnderlineDouble()
				s.Blink()
				s.Reverse()
				s.Hidden()
				s.Strikeout()
			},
		},
		{
			name: "All (Curly)",
			want: "1;2;3;4:3;5;7;8;9",
			actions: func(s *Style) {
				s.Bold()
				s.Dim()
				s.Italics()
				s.UnderlineCurly()
				s.Blink()
				s.Reverse()
				s.Hidden()
				s.Strikeout()
			},
		},
		{
			name: "Reset",
			want: "0",
			actions: func(s *Style) {
				s.Bold()
				s.Dim()
				s.Italics()
				s.UnderlineCurly()
				s.Blink()
				s.Reverse()
				s.Hidden()
				s.Strikeout()

				s.Normal()
			},
		},
		{
			name: "Reset (Manual)",
			want: "0",
			actions: func(s *Style) {
				s.Bold()
				s.Dim()
				s.Italics()
				s.Underline()
				s.Blink()
				s.Reverse()
				s.Hidden()
				s.Strikeout()

				s.NoBold()
				s.NoDim()
				s.NoItalics()
				s.NoUnderline()
				s.NoBlink()
				s.NoReverse()
				s.NoHidden()
				s.NoStrikeout()
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := Style{}
			c.actions(&s)

			got := s.Code()
			if got != c.want {
				t.Errorf("want: '%s', got: '%s' :(", c.want, got)
			}
		})
	}
}
