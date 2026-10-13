package main

import "strings"

func tokenize(data []byte) []string {
	return strings.Fields(string(data))
}
