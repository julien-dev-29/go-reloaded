package main

import (
	"os"
	"reflect"
	"testing"
)

func Test_tokenize(t *testing.T) {
	args, err := os.ReadFile("test.txt")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		args []byte
		want []string
	}{
		{"test1", args, []string{"yolo", "les", "kikis!"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tokenize(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tokenize() = %v, want %v", got, tt.want)
			}
		})
	}
}
