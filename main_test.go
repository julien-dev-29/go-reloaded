package main

import (
	"io/fs"
	"log"
	"os"
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

func Test_getWordsToFormatIndex(t *testing.T) {
	data := readFile("test.txt")
	markers := findMarkers(data)
	tests := []struct {
		name   string
		data   []byte
		marker Marker
		want   int
	}{
		{"String with test.txt and marker 1", data, markers[0], 0},
		{"String with test.txt and marker 2", data, markers[1], 52},
		{"String with test.txt and marker 3", data, markers[2], 91},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getWordsToFormatIndex(tt.data, tt.marker)
			if got != tt.want {
				t.Errorf("findWords(%q, %q) = %v, want = %v", tt.data, tt.marker, got, tt.want)
			}
		})
	}
}

// func Test_process(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		data    []byte
// 		markers []Marker
// 		want    []byte
// 	}{
// 		{"Test with test.txt and marker 1", readFile("test.txt"), findMarkers(readFile("test.txt")), []byte("It")},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := process(tt.data, tt.markers)
// 			if !slices.Equal(got, tt.want) {
// 				t.Errorf("process() = %v, want %v", string(got), string(tt.want))
// 			}
// 		})
// 	}
// }
