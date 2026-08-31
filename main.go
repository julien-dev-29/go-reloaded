package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"strconv"
)

type Marker struct {
	position int
	option   string
	value    int
}

type Word struct {
	start int
	end   int
}

func main() {
	data := readFile()
	markers := findMarkers(data)
	for _, marker := range markers {
		name := marker.option
		switch name {
		case "cap":
			Capitalize(findWords(data, marker))
			fmt.Println(string(data))
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
			if isMarkerGotValue(markerContent) {
				marker.option, marker.value = splitOptionAndValueFromMarker(markerContent)
			} else {
				marker.option = markerContent
				marker.value = 1
			}
			markers[markerIndex] = marker
			markerIndex++
		}
	}
	return markers
}

func splitOptionAndValueFromMarker(name string) (string, int) {
	var optionEndIndex int
	for i := range name {
		if name[i] == ',' {
			optionEndIndex = i
		}
	}
	value, err := strconv.Atoi(string(name[optionEndIndex+2:]))
	if err != nil {
		log.Fatal("Error:", err)
	}
	return name[0:optionEndIndex], value // TODO fix values > 9
}

// Est ce qu'il y a une value, en gros est ce qu'il y a une virgule
func isMarkerGotValue(name string) bool {
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

func isALetter(r rune) bool {
	if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
		return true
	}
	return false
}

func findWords(data []byte, marker Marker) []byte {
	splitFirstSpace := 1
	index := marker.position
	fmt.Println(marker.value)
	for i := 0; i < marker.value; i++ {
		index = findStartOfWord(data[:index-splitFirstSpace])
	}
	return data[index : marker.position-splitFirstSpace]
}

func findStartOfWord(data []byte) int {
	result := 0
	for i := len(data) - 1; i >= 0; i-- {
		if data[i] == ' ' {
			result = i + 1
			break
		}
	}
	return result
}
