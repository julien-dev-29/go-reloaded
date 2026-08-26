package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"testing"
)

func TestCountMarkers(t *testing.T) {
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

func TestSplitValue(t *testing.T) {
	expectName := "yolo"
	expectValue := 8
	name, value := splitValue("yolo, 8")
	if name != expectName || value != fmt.Sprint(expectValue) {
		t.Errorf("Error")
	}
}
