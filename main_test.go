package main

import (
	"io/fs"
	"log"
	"os"
	"testing"
)

func TestCountMarkers(t *testing.T) {
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

func TestSplitOptionAndValueFromMarker(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantOption string
		wantValue  string
	}{
		{"valid input", "cap, 8", "cap", "8"},
		{"another input", "foo, 3", "foo", "3"},
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

func TestIsLower(t *testing.T) {
	tests := []struct {
		name  string
		input rune
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

func TestCapitalize(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"String with only lower characters", "yolo", "Yolo"},
		{"String with lower and upper characters", "yoPo66", "Yopo66"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := capitalize(tt.input)
			if got != tt.want {
				t.Errorf("capitalize(%q) = %q, want = %q", tt.input, got, tt.want)
			}
		})
	}

}
