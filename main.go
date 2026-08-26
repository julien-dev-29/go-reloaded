package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
)

type Marker struct {
	start int
	end   int
	name  string
	value string
}

func main() {
	fsys := os.DirFS(".")
	data, err := fs.ReadFile(fsys, "test.txt")
	if err != nil {
		log.Fatal(err)
	}
	markers := findMarkers(data)
	for i := range markers {
		fmt.Println(markers[i].name)
	}
}

// Find the "("
func findMarkers(data []byte) []Marker {
	count := countMarkers(data)
	markers := make([]Marker, count)
	count = 0
	for i := range data {
		if data[i] == '(' {
			count++
			var marker Marker
			j := i
			for data[j] != ')' {
				j++
			}
			name := string(data[i+1 : j])
			var value string
			if isAValue(name) {
				name, value = splitValue(name)

			}
			marker.start = i
			marker.end = j
			marker.name = string(name)
			marker.value = string(value)
			markers[count-1] = marker
		}
	}
	return markers
}

// Est ce qu'il y a une value, en gros est ce qu'il y a une virgule
func splitValue(name string) (string, string) {
	var k int
	var j int
	for i := range name {
		if name[i] == ',' {
			k = i
			j = i
			for name[j] <= '0' || name[j] >= '9' {
				j++
			}
		}
	}
	return name[0:k], string(name[j])
}

func isAValue(name string) bool {
	for _, c := range name {
		if c == ',' {
			return true
		}
	}
	return false
}

func countMarkers(data []byte) int {
	count := 0
	for i := range data {
		if data[i] == '(' {
			count++
		}
	}
	return count
}
