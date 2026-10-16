package main

import (
	"reflect"
	"testing"
)

func Test_assemble(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want string
	}{
		{"basic", []string{"yolo", "les", "kikis"}, "yolo les kikis"},
		{"basic", []string{"yolo", "les", "kikis", "!"}, "yolo les kikis!"},
		{"group punct", []string{"Punctuation", "tests", "are", "...", "kinda", "boring", ",", "what", "do", "you", "think", "?"}, "Punctuation tests are... kinda boring, what do you think?"},
		{"quote punct", []string{"As", "Elton", "John", "said", ":", "'", "I", "am", "the", "most", "well-known", "homosexual", "in", "the", "world", "'"}, "As Elton John said: 'I am the most well-known homosexual in the world'"},
		{"quote punct", []string{"I", "am", "exactly", "how", "they", "describe", "me", ":", "'", "awesome", "'"}, "I am exactly how they describe me: 'awesome'"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := assemble(tt.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
