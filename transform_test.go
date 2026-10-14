package main

import (
	"slices"
	"testing"
)

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

func Test_fix(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{
			"capitalize",
			[]string{"yolo", "lala", "popo", "(cap)"},
			[]string{"yolo", "lala", "Popo"},
		}, {
			"capitalize count 3",
			[]string{"yolo", "lala", "popo", "(cap,", "3)"},
			[]string{"Yolo", "Lala", "Popo"},
		},
		{
			"up simple",
			[]string{"Ready", "for", "the", "show", "(up)", "!"},
			[]string{"Ready", "for", "the", "SHOW", "!"},
		},
		{
			"up count=3",
			[]string{"Ready", "for", "the", "show", "(up,", "3)", "!"},
			[]string{"Ready", "FOR", "THE", "SHOW", "!"},
		},
		{
			"low simple",
			[]string{"Ready", "for", "the", "SHOW", "(low)", "!"},
			[]string{"Ready", "for", "the", "show", "!"},
		},
		{
			"low count=2",
			[]string{"Ready", "for", "THE", "SHOW", "(low,", "2)", "!"},
			[]string{"Ready", "for", "the", "show", "!"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := fix(tt.in); !slices.Equal(got, tt.want) {
				t.Errorf("want %v, got %v", tt.want, got)
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

func Test_isStartWithVowelOrH(t *testing.T) {
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
		{"A", "A", true},
		{"E", "E", true},
		{"I", "I", true},
		{"O", "O", true},
		{"U", "U", true},
		{"Y", "Y", true},
		{"h", "h", true},
		{"z", "z", false},
		{"w", "w", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isStartWithVowelOrH(tt.in)

			if got != tt.want {
				t.Errorf("got %t want %t", got, tt.want)
			}
		})
	}
}
