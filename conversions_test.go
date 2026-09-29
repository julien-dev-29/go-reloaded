package main

import (
	"testing"
)

func Test_hex2Dec(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
		ok   bool
	}{
		{"Uppercase", "1E", "30", true},
		{"Lowercase", "53", "83", true},
		{"Invalid", "XYZ", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := hex2Dec(tt.in)
			if got != tt.want {
				t.Errorf("Hex2dec() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.ok {
				t.Errorf("Hex2dec() got1 = %v, want %v", got1, tt.ok)
			}
		})
	}
}

func Test_bin2Dec(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
		ok   bool
	}{
		{"simple", "10", "2", true},
		{"long", "1010", "10", true},
		{"zero", "0", "0", true},
		{"invalid", "12", "", false},
		{"invalidDigit", "10AB", "", false},
		{"empty", "", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := bin2Dec(tt.in)
			if got != tt.want {
				t.Errorf("Bin2Dec() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.ok {
				t.Errorf("Bin2Dec() got1 = %v, want %v", got1, tt.ok)
			}
		})
	}
}
