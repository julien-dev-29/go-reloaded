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

func applyCase(tokens []string) []string {
	out := make([]string, 0, len(tokens))
	for _, tok := range tokens {
		if len(out) == 0 {
			out = append(out, tok)
			continue
		}
		option := ""
		val := 1
		if []byte(tok)[0] == '(' && []byte(tok)[len(tok)-1] == ')' {
			option, val = parseMarker(tok)
		}
		switch option {
		case "low":
			for i := range val {
				out[len(out)-1-int(i)] = strings.ToLower(out[len(out)-1-int(i)])
			}
		case "up":
			for i := range val {
				out[len(out)-1-int(i)] = strings.ToUpper(out[len(out)-1-int(i)])
			}
		case "cap":
			for i := range val {
				out[len(out)-1-i] = capitalize(out[len(out)-1-i])
			}
		default:
			out = append(out, tok)
		}
	}
	return out
}

func parseMarker(s string) (string, int) {
	if !strings.ContainsRune(s, ',') {
		return string([]byte(s)[1 : len(s)-1]), 1
	}
	bytes := []byte(s)[1 : len(s)-1]
	str := string(bytes)
	yolo := strings.Split(str, ",")
	value, err := strconv.ParseInt(yolo[1], 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
	}
	return yolo[0], int(value)
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


