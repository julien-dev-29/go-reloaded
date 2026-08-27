package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
)

type Marker struct {
	position int
	option   string
	value    string
}

func main() {
	data := readFile()
	markers := findMarkers(data)
	for _, marker := range markers {
		name := marker.option
		switch name {
		case "cap":
			fmt.Println(marker.option)
		default:

		}
	}
}

func readFile() []byte {
	fsys := os.DirFS(".")
	data, err := fs.ReadFile(fsys, "test.txt")
	if err != nil {
		log.Fatal(err)
	}
	return data
}

// Find the "("
func findMarkers(data []byte) []Marker {
	count := countMarkers(data)
	markers := make([]Marker, count)
	markerIndex := 0
	for i := range data {
		if data[i] == '(' {
			var marker Marker
			marker.position = i
			j := i
			for data[j] != ')' {
				j++
			}
			bracketIndex := 1
			markerContent := string(data[i+bracketIndex : j])
			if isMarkerGotOptionAndValue(markerContent) {
				marker.option, marker.value = splitOptionAndValueFromMarker(markerContent)
			} else {
				marker.option = markerContent
				marker.value = "1"
			}
			markers[markerIndex] = marker
			count++
		}
	}
	return markers
}

func splitOptionAndValueFromMarker(name string) (string, string) {
	var optionEndIndex int
	var j int
	for i := range name {
		if name[i] == ',' {
			optionEndIndex = i
			j = i
			for name[j] <= '0' || name[j] >= '9' {
				j++
			}
		}
	}
	return name[0:optionEndIndex], string(name[j]) // TODO fix values > 9
}

// Est ce qu'il y a une value, en gros est ce qu'il y a une virgule
func isMarkerGotOptionAndValue(name string) bool {
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

func isLower(r rune) bool {
	if r >= 'a' && r <= 'z' {
		return true
	}
	return false
}

func isALetter(r rune) bool {
	if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
		return true
	}
	return false
}

func capitalize(s string) string {
	rs := []rune(s)

	for i := range rs {
		if i == 0 && isLower(rs[i]) {
			rs[i] = rs[i] - 32
		}
		if i > 0 && isALetter(rs[i]) && !isLower(rs[i]) {
			rs[i] = rs[i] + 32
		}
	}
	return string(rs)
}

// func findWords(s string, number int) []string {

// }
