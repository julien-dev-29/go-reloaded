package main

import "strings"

func assemble(tokens []string) string {
	var b strings.Builder

	for i, tok := range tokens {
		if tok == "'" {
			b.WriteString(tok)
		} else if i < len(tokens)-1 && !isPunctChar(tokens[i+1]) && !isGroupPunctChars(tokens[i+1]) {
			b.WriteString(tok)
			b.WriteString(" ")
		} else {
			b.WriteString(tok)
		}
	}

	return b.String()
}
