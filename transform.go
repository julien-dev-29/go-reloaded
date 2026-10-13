package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var markerRegex = regexp.MustCompile(`\(\s*(hex|bin|up|low|cap)\s*(?:,\s*(\d+)\s*)?\)`)

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func normalizeMarkers(s string) string {
	return markerRegex.ReplaceAllStringFunc(s, func(m string) string {
		subs := markerRegex.FindStringSubmatch(m)
		if subs[2] != "" {
			return "(" + subs[1] + "," + subs[2] + ")"
		}
		return "(" + subs[1] + ")"
	})
}

func convertMarkers(tokens []string) []string {
	out := make([]string, 0, len(tokens))
	for _, tok := range tokens {
		if len(out) == 0 {
			out = append(out, tok)
			continue
		}
		switch tok {
		case "(hex)":
			if val, ok := hex2Dec(out[len(out)-1]); ok {
				out[len(out)-1] = val
			} else {
				out = append(out, tok)
			}
		case "(bin)":
			if val, ok := bin2Dec(out[len(out)-1]); ok {
				out[len(out)-1] = val
			} else {
				out = append(out, tok)
			}
		default:
			out = append(out, tok)
		}
	}
	return out
}

func applyCase(tokens []string) ([]string, error) {
	out := make([]string, 0, len(tokens))
	for i, tok := range tokens {
		if len(out) == 0 {
			out = append(out, tok)
			continue
		}
		var (
			option string
			val    int
		)
		if tok[len(tok)-1] == ')' {
			if tok[0] == '(' {
				option, val = tok[1:len(tok)-1], 1
			} else if i >= 2 {
				convertVal, err := strconv.ParseInt(string(tok[0:len(tok)-1]), 10, 64)
				if err != nil {
					continue
				}
				option, val = tokens[i-2][1:], int(convertVal)
			}
		}
		fmt.Println(option)
		switch option {
		case "low":
			for i := range val {
				out[len(out)-3-int(i)] = strings.ToLower(out[len(out)-3-int(i)])
			}
		case "up":
			for i := range val {
				out[len(out)-3-int(i)] = strings.ToUpper(out[len(out)-3-int(i)])
			}
		case "cap":
			for i := range val {
				out[len(out)-3-i] = capitalize(out[len(out)-3-i])
			}

		default:
			out = append(out, tok)
		}
	}
	return out, nil
}

func fixArticle(tokens []string) []string {
	out := append([]string{}, tokens...)
	for i, tok := range out {
		if isStartWithVowelOrH(tok) {
			switch out[i-1] {
			case "a":
				out[i-1] = "an"
			case "A":
				out[i-1] = "An"
			}
		}
	}
	return out
}

func isStartWithVowelOrH(token string) bool {
	if token == "" {
		return false
	}
	switch token[0] {
	case 'a', 'e', 'i', 'o', 'u', 'y', 'h', 'A', 'E', 'I', 'O', 'U', 'Y', 'H':
		return true
	}
	return false
}
