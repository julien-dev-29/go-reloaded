package main

import (
	"testing"
)

func TestConvertDecToHex(t *testing.T) {
	tests := []struct {
		name  string
		input int
		want  float64
	}{
		{"Test converting 42", 42, 66},
		{"Test converting 96", 96, 150},
		{"Test converting 500", 500, 1280},
		{"Test converting 600", 600, 1536},
		{"Test converting 800", 800, 2048},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertHex2Dec(tt.input, 0, 0)
			if got != tt.want {
				t.Errorf("ConvertDecToHex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHex2Dec(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int64
	}{
		{"Test converting 42", "42", 66},
		{"Test converting 500", "500", 1280},
		{"Test converting 600", "600", 1536},
		{"Test converting 800", "800", 2048},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Hex2Dec(tt.input); got != tt.want {
				t.Errorf("Hex2Dec() = %v, want %v", got, tt.want)
			}
		})
	}
}
