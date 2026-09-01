package main_test

import (
	main "go-reloaded"
	"slices"
	"testing"
)

func TestToLowercase(t *testing.T) {
	tests := []struct {
		name  string
		bytes []byte
		want  []byte
	}{
		{"Test with uppercase letters", []byte("YOLO LES KIKIS"), []byte("yolo les kikis")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := main.ToLowercase(tt.bytes)
			if !slices.Equal(got, tt.want) {
				t.Errorf("ToLowercase() = %q, want %q", got, tt.want)
			}
		})
	}
}
