package main

import (
	"io/fs"
	"log"
	"os"
	"slices"
	"testing"
)

func Test_countMarkers(t *testing.T) {
	fsys := os.DirFS(".")
	data, err := fs.ReadFile(fsys, "test.txt")
	if err != nil {
		log.Fatal(err)
	}
	count := countMarkers(data)
	want := 4
	if count != want {
		t.Errorf("Error")
	}
}

func Test_splitOptionAndValueFromMarker(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantOption string
		wantValue  int
	}{
		{"valid input", "cap, 8", "cap", 8},
		{"another input", "foo, 3", "foo", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotValue := splitOptionAndValueFromMarker(tt.input)
			if tt.wantOption != gotName {
				t.Errorf("splitValue(%q) gotName = %q, want %q", tt.input, gotName, tt.wantOption)
			}
			if tt.wantValue != gotValue {
				t.Errorf("splitValue(%q) gotValue = %q want %q", tt.input, gotValue, tt.wantValue)
			}
		})
	}

}

func Test_isLower(t *testing.T) {
	tests := []struct {
		name  string
		input byte
		want  bool
	}{
		{"lowercase letter", 'c', true},
		{"uppercase letter", 'C', false},
		{"digit", '1', false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isLower(tt.input); got != tt.want {
				t.Errorf("isLower(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func Test_findWords(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  []byte
	}{
		{"String with one word", readFile("test.txt"), []byte("it was the age of foolishness")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findWords(tt.input, findMarkers(readFile("test.txt"))[2])
			if !slices.Equal(got, tt.want) {
				t.Errorf("findWords(%q) = %q, want = %q", tt.input, got, tt.want)
			}
		})
	}
}

func Test_findStartOfWord(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  int
	}{
		{"Trouver l'index de début d'un mot dans une phrase", readFile("test.txt")[100:120], 9},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findStartOfWord(tt.input)
			if got != tt.want {
				t.Errorf("findWords(%q) = %d, want = %d", tt.input, got, tt.want)
			}
		})
	}
}
