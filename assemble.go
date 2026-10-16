package main

import "strings"

func assemble(tokens []string) string {
	var b strings.Builder
	open := false
	for i, tok := range tokens {

		if tok == "'" && open == false {
			b.WriteString(" ")
			b.WriteString(tok)
			open = true
		} else if tok == "'" && open == true {
			b.WriteString(tok)
			open = false
		} else if i < len(tokens)-1 && !isPunctChar(tokens[i+1]) && !isGroupPunctChars(tokens[i+1]) {
			b.WriteString(tok)
			b.WriteString(" ")
		} else {
			b.WriteString(tok)
		}
	}

	return b.String()
}
