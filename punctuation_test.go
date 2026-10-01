package main

import (
	"reflect"
	"testing"
)

func Test_separatePunctuation(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{
			"leading punctuation",
			[]string{"there", ",and", "then"},
			[]string{"there", ",", "and", "then"},
		},
		{
			"punct group",
			[]string{"are", "...", "kinda"},
			[]string{"are", "...", "kinda"},
		},
		{
			"question prefix",
			[]string{"boring", ",what", "do"},
			[]string{"boring", ",", "what", "do"},
		},
		{
			"exclamation suffix untouched",
			[]string{"BAMM", "!!"},
			[]string{"BAMM", "!!"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := separatePunctuation(tt.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("separatePunctuation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isPunctChar(t *testing.T) {
	tests := []struct {
		name string
		in   byte
		want bool
	}{
		{".", '.', true},
		{",", ',', true},
		{"?", '?', true},
		{"!", '!', true},
		{"a", 'a', false},
		{"¨", '¨', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isPunctChar(tt.in)

			if got != tt.want {
				t.Errorf("got %t want %t", got, tt.want)
			}
		})
	}
}
