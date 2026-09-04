package main

import (
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
	mylog := log.New(os.Stdout, "jurol:", log.LstdFlags)
	mylog.SetFlags(log.Lmicroseconds)
	data := readFile("test.txt")
	for i, b := range data {
		mylog.Println("Byte", i, ":", string(b))
		if b == '(' {
			j := i
			mylog.Println("----------------------->", string(b))
		}
	}
}

func readFile(filename string) []byte {
	fsys := os.DirFS(".")
	data, err := fs.ReadFile(fsys, filename)
	if err != nil {
		log.Fatal(err)
	}
	return data
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
func isMarkerGotValue(name []byte) bool {
	for _, c := range name {
		if c == ',' {
			return true
		}
	}
	return false
}

func isALetter(r rune) bool {
	if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
		return true
	}
	return false
}

func getWordsToFormatIndex(data []byte, marker Marker) int {
	i := 0
	startIndex := marker.startIndex - 2
	count := 0
	for count < marker.value {
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
