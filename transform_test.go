package main

import (
	"slices"
	"testing"
)

func Test_normalizeMarkers(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"bare", "(up)", "(up)"},
		{"spaced inside", "( up )", "(up)"},
		{"number spaced", "(up, 2)", "(up,2)"},
		{"number tight", "(up,2)", "(up,2)"},
		{"space before comma", "(up , 2)", "(up,2)"},
		{"hex", "(hex)", "(hex)"},
		{"bin", "(bin)", "(bin)"},
		{"low", "(low)", "(low)"},
		{"cap", "(cap)", "(cap)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeMarkers(tt.in); got != tt.want {
				t.Errorf("normalizeMarker() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_capitalize(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"normal", "yolo", "Yolo"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := capitalize(tt.in); got != tt.want {
				t.Errorf("isLowercase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertMarkers(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{
			"hex converts previous word",
			[]string{"1E", "(hex)", "files", "were", "added"},
			[]string{"30", "files", "were", "added"},
		},
		{
			"bin converts previous word",
			[]string{"It", "has", "been", "10", "(bin)", "years"},
			[]string{"It", "has", "been", "2", "years"},
		},
		{
			"invalid hex marker kept",
			[]string{"XYZ", "(hex)"},
			[]string{"XYZ", "(hex)"},
		},
		{
			"marker at start kept",
			[]string{"(hex)", "value"},
			[]string{"(hex)", "value"},
		},
		{
			"multiple markers",
			[]string{"42", "(hex)", "and", "10", "(bin)"},
			[]string{"66", "and", "2"},
		},
		{
			"no markers unchanged",
			[]string{"hello", "world"},
			[]string{"hello", "world"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := convertMarkers(tc.in)
			if len(got) != len(tc.want) {
				t.Fatalf("convertMarkers(%v) length=%d, want %d: %v", tc.in, len(got), len(tc.want), got)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Fatalf("convertMarkers(%v)=%v, want %v", tc.in, got, tc.want)
				}
			}
		})
	}
}

func Test_applyCase(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{
			"capitalize",
			[]string{"yolo", "lala", "popo", "(cap,3)"},
			[]string{"Yolo", "Lala", "Popo"},
		},
		{
			"up simple",
			[]string{"Ready", "for", "the", "show", "(up)", "!"},
			[]string{"Ready", "for", "the", "SHOW", "!"},
		},
		{
			"up count=3",
			[]string{"Ready", "for", "the", "show", "(up,3)", "!"},
			[]string{"Ready", "FOR", "THE", "SHOW", "!"},
		},
		{
			"low simple",
			[]string{"Ready", "for", "the", "SHOW", "(low)", "!"},
			[]string{"Ready", "for", "the", "show", "!"},
		},
		{
			"low count=2",
			[]string{"Ready", "for", "THE", "SHOW", "(low,2)", "!"},
			[]string{"Ready", "for", "the", "show", "!"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := applyCase(tt.in); !slices.Equal(got, tt.want) {
				t.Errorf("applyCase(%v)=%v, got %v", tt.in, tt.want, got)
			}
		})
	}
}

func Test_parseMarker(t *testing.T) {
	tests := []struct {
		name  string
		in    string
		want1 string
		want2 int
	}{
		{"Marker with a value", "(up,2)", "up", 2},
		{"Marker without a value", "(up)", "up", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1, got2 := parseMarker(tt.in)
			if got1 != tt.want1 || got2 != tt.want2 {
				t.Errorf("out1 %s want %s, out2 %d want2 %d", got1, tt.want1, got2, tt.want2)
			}
		})
	}
}

func Test_fixArticle(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{"lowercase article with vowel", []string{"bearing", "a", "untold"}, []string{"bearing", "an", "untold"}},
		{"lowercase article without vowel", []string{"bearing", "a", "told"}, []string{"bearing", "a", "told"}},
		{"lowercase article with 'h'", []string{"bearing", "a", "hold"}, []string{"bearing", "an", "hold"}},
		{"uppercase article with vowel", []string{"bearing", "A", "untold"}, []string{"bearing", "An", "untold"}},
		{"uppercase article without vowel", []string{"bearing", "A", "told"}, []string{"bearing", "A", "told"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fixArticle(tt.in)

			if !slices.Equal(got, tt.want) {
				t.Errorf("got %s want %s", got, tt.want)
			}
		})
	}
}

func Test_isStartWithVowel(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want bool
	}{
		{"a", "a", true},
		{"e", "e", true},
		{"i", "i", true},
		{"o", "o", true},
		{"u", "u", true},
		{"y", "y", true},
		{"z", "z", false},
		{"w", "w", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isStartWithVowel(tt.in)

			if got != tt.want {
				t.Errorf("got %t want %t", got, tt.want)
			}
		})
	}
}

func Test_isVowel(t *testing.T) {
	tests := []struct {
		name string
		in   rune
		want bool
	}{
		{"a", 'a', true},
		{"e", 'e', true},
		{"i", 'i', true},
		{"o", 'o', true},
		{"u", 'u', true},
		{"y", 'y', true},
		{"A", 'A', true},
		{"E", 'E', true},
		{"I", 'I', true},
		{"O", 'O', true},
		{"U", 'U', true},
		{"Y", 'Y', true},
		{"b", 'b', false},
		{"Z", 'Z', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isVowel(tt.in)

			if got != tt.want {
				t.Errorf("got %t want %t", got, tt.want)
			}
		})
	}
}

func Test_separatePunctuation(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{"basic", []string{"Punctuation", "tests", "are", "...", "kinda", "boring", ",what do you think", "?"},
			[]string{"Punctuation", "tests", "are...", "kinda", "boring,", "what", "do", "you", "think?"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := separatePunctuation(tt.in)

			if !slices.Equal(got, tt.want) {
				t.Errorf("got %s want %s", got, tt.want)
			}
		})
	}
}

func Test_startWithComma(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want bool
	}{
		{"basic", ",what", true},
		{"invalid", "what", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := startWithComma(tt.in)

			if got != tt.want {
				t.Errorf("got %t want %t", got, tt.want)
			}
		})
	}
}

func Test_addComma(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"basic", "yolo", "yolo,"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := addComma(tt.in)

			if got != tt.want {
				t.Errorf("got %s want %s", got, tt.want)
			}
		})
	}
}

func Test_deleteComma(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"basic", ",yolo", "yolo"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := deleteComma(tt.in)

			if got != tt.want {
				t.Errorf("got %s want %s", got, tt.want)
			}
		})
	}
}

// func Test_assemble(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		in   []string
// 		want string
// 	}{
// 		{"basic", []string{"yolo", "les", "kikis"}, "yolo les kikis"},
// 	}
// }
