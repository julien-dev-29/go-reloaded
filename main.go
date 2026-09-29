package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {

	if len(os.Args) > 2 {
		fmt.Println("Too much arguments")
		return
	}
	path := os.Args[1]
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Erreur pendant la lecture du fichier")
		return
	}
	fmt.Println(string(data))
	formatText := normalizeMarkers(string(data))
	fmt.Println(formatText)
	words := strings.FieldsSeq(formatText)
	for word := range words{
		fmt.Println(word)
	}
}
