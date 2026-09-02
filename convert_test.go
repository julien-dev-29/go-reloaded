package main

import (
	"fmt"
	"slices"
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
		input []byte
		want  []byte
	}{
		{"Test converting 42", []byte("42"), []byte("66")},
		{"Test converting 500", []byte("500"), []byte("1280")},
		{"Test converting 600", []byte("600"), []byte("1536")},
		{"Test converting 800", []byte("800"), []byte("2048")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := Hex2Dec(tt.input); !slices.Equal(got, tt.want) {
				if err != nil {
					fmt.Println("Error:", err)
				}
				t.Errorf("Hex2Dec() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func TestConvertHex2Dec(t *testing.T) {
	type args struct {
		dec    int
		count  int
		result float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertHex2Dec(tt.args.dec, tt.args.count, tt.args.result); got != tt.want {
				t.Errorf("ConvertHex2Dec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBin2Dec(t *testing.T) {
	tests := []struct {
		name  string
		input int
		want  float64
	}{
		{"Test function with 101", 101, 5},
		{"Test function with 111101", 111101, 61},
		{"Test function with 10011", 10011, 19},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Bin2Dec(tt.input, 0, 0)
			if got != tt.want {
				t.Errorf("Bin2Dec() = %v, want %v", got, tt.want)
			}
		})
	}
}
