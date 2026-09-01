package main

import (
	"slices"
	"testing"
)

func TestToUppercase(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  []byte
	}{
		{"Test with lowercase letters", []byte("yolo les kikis"), []byte("YOLO LES KIKIS")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToUppercase(tt.input)
			if !slices.Equal(got, tt.want) {
				t.Errorf("findWords(%q) = %q, want = %q", tt.input, got, tt.want)
			}
		})
	}
}
