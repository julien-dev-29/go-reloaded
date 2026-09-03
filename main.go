package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"strconv"
)

type Marker struct {
	startIndex int
	endIndex   int
	option     string
	value      int
}

type Word struct {
	start int
	end   int
}

func main() {
	data := readFile("test.txt")
	fmt.Println(string(data))
	markers := findMarkers(data)
	result := process(data, markers)
	fmt.Println(string(result))
}

func process(data []byte, markers []Marker) []byte {
	var result []byte
	for _, marker := range markers {
		name := marker.option
		switch name {
		case "cap":
			wordsSlice := data[getWordsToFormatIndex(data, marker) : marker.startIndex-1]
			capWordsSlices := Capitalize(wordsSlice)
			fmt.Println(string(capWordsSlices))
		case "up":
			//result = ToUppercase(data[findWordsStart(data, marker) : marker.position-1])
		case "low":
			//result = ToLowercase(data[findWordsStart(data, marker) : marker.position-1])
		case "hex":
			//slices.Grow()()(data[findWordsStart(data, marker) : marker.position-1])
		default:

		}
	}
	return result
}

func readFile(filename string) []byte {
	fsys := os.DirFS(".")
	data, err := fs.ReadFile(fsys, filename)
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
			marker.startIndex = i
			j := i
			for data[j] != ')' {
				j++
			}
			marker.endIndex = j
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

func getWordsToFormatIndex(data []byte, marker Marker) int {
	fmt.Println(marker)
	i := 0
	startIndex := marker.startIndex - 2
	count := 0
	for count != marker.value {
		i++
		if startIndex-i == 0 {
			return 0
		}
		if data[startIndex-i] == ' ' {
			count++
		}
	}
	i--
	return startIndex - i
}
