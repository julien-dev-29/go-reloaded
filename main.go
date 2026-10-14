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

	if len(os.Args) > 3 {
		fmt.Println("Too much arguments")
		return
	}

	path := os.Args[1]
	data, err := readFile(path)

	if err != nil {
		fmt.Println(err)
		return
	}

	tokens := tokenize(data)

	fixed, err := fix(tokens)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fixed = separatePunctuation(fixed)

	fixed = fixArticle(fixed)

	fixedStr := assemble(fixed)
	os.WriteFile(os.Args[2], []byte(fixedStr), 0644)
}
