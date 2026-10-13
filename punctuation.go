package main

import (
	"strings"
)

var punctuations = ",.?;:!"
var suspension = "..."

func separatePunctuation(tokens []string) []string {
	out := make([]string, 0, len(tokens))
	for _, tok := range tokens {
		if isPunctChar(tok[0]) && len(tok) > 1 && !isGroupPunctChars(tok) {
			out = append(out, string(tok[0]))
			out = append(out, tok[1:])
		} else if isPunctChar(tok[len(tok)-1]) && len(tok) > 1 && !isGroupPunctChars(tok) {
			out = append(out, tok[:len(tok)-1])
			out = append(out, string(tok[len(tok)-1]))
		} else {
			out = append(out, tok)
		}
	}
	return out
}

func isPunctChar(char byte) bool {
	return strings.ContainsRune(punctuations, rune(char))
}

func isGroupPunctChars(tok string) bool {
	if len(tok) == 1 {
		return false
	}
	return strings.ContainsRune(punctuations, rune(tok[0])) && strings.ContainsRune(punctuations, rune(tok[1]))
}
