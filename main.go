package main

import (
	"fmt"
	"os"
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
	// tokenizeText := tokenize(data)
	// formatText := convertMarkers(tokenizeText)
	// formatText, err = applyCase(formatText)
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }
	// formatText = separatePunctuation(formatText)
	// fmt.Println(strings.Join(formatText, "|"))
}
