package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"slices"
)

type Marker struct {
	start int
	end   int
	name  string
	value string
}

func main() {
	fsys := os.DirFS(".")
	data, err := fs.ReadFile(fsys, "test.txt")
	if err != nil {
		log.Fatal(err)
	}
	findMarkers(data)
}

// Find the "("
func findMarkers(data []byte) {
	//markers := make([]Marker, 2000)
	for i := range data {
		if data[i] == '(' {
			var marker Marker
			j := i
			for data[j] != ')' {
				j++
			}
			name := data[i+1 : j]
			var value byte
			if isAValue(name) {
				name, value = splitValue(name)

			}
			marker.name = string(name)
			marker.value = string(value)
			fmt.Println("Nouveau marker")
			fmt.Print(marker.start)
			fmt.Print(marker.end)
			fmt.Print(marker.name)
			fmt.Print(marker.value)
			fmt.Println()
		}
	}
}

// Est ce qu'il y a une value, en gros est ce qu'il ya une virgule
func splitValue(name []byte) ([]byte, byte) {
	var k int
	var j int
	if isAValue(name) {
		for i := range name {
			if name[i] == ',' {
				k = i
				j = i
				for name[j] <= '0' || name[j] >= '9' {
					j++
				}
			}
		}
		return name[0:k], name[j]
	} else {
		return name, name[j]
	}
}

func isAValue(name []byte) bool {
	return slices.Contains(name, ',')
}
