package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	
	if len(os.Args) == 1 {
		fmt.Println("Must provide the path of file")
		return
	}

	if len(os.Args) > 2 {
		fmt.Println("Too much arguments")
		return
	}

	path := os.Args[1]
	data, err := readFile(path)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(data))
	normalizedText := normalizeMarkers(string(data))
	formatText := convertMarkers(strings.Fields(normalizedText))
	formatText = applyCase(formatText)
	formatText = separatePunctuation(formatText)
	fmt.Println(strings.Join(formatText, "|"))
}
