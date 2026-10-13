package main

import (
	"reflect"
	"testing"
)

func Test_readFile(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    []byte
		wantErr bool
	}{
		{"test file", "test.txt", []byte("yolo les kikis!"), false},
		{"test file do not exists", "test2.txt", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readFile(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("readFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readFile() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
