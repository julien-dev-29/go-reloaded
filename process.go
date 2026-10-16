package main

import "fmt"

func process(data []byte) (string, error) {
	tokens := tokenize(data)

	fixed, err := fix(tokens)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	fixed = separatePunctuation(fixed)

	fixed = fixArticle(fixed)

	return assemble(fixed), nil
}
