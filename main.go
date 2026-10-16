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

	fixedStr, err := process(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	os.WriteFile(os.Args[2], []byte(fixedStr), 0644)
}
