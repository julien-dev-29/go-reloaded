package main

import (
	"strconv"
	"strings"
	"unicode"
)

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func fix(tokens []string) ([]string, error) {
	out := make([]string, 0, len(tokens))

	for i := 0; i <= len(tokens)-1; i++ {
		switch tokens[i] {

		case "(hex)":
			if val, ok := hex2Dec(out[len(out)-1]); ok {
				out[len(out)-1] = val
			} else {
				out = append(out, tokens[i])
			}

		case "(bin)":
			if val, ok := bin2Dec(out[len(out)-1]); ok {
				out[len(out)-1] = val
			} else {
				out = append(out, tokens[i])
			}

		case "(low)":
			out[len(out)-1] = strings.ToLower(out[len(out)-1])

		case "(up)":
			out[len(out)-1] = strings.ToUpper(out[len(out)-1])

		case "(cap)":
			out[len(out)-1] = capitalize(out[len(out)-1])

		case "(low,":
			if i < len(tokens) {
				num, ctrl := getNum(tokens, i)
				switch ctrl {
				case 1:
					continue
				}
				for j := range num {
					out[len(out)-1-j] = strings.ToLower(out[len(out)-1-j])
				}
				i++
			}

		case "(up,":
			if i < len(tokens) {
				num, ctrl := getNum(tokens, i)
				switch ctrl {
				case 1:
					continue
				}
				for j := range num {
					out[len(out)-1-j] = strings.ToUpper(out[len(out)-1-j])
				}
				i++
			}
		case "(cap,":
			if i < len(tokens) {
				num, ctrl := getNum(tokens, i)
				switch ctrl {
				case 1:
					continue
				}
				for j := range num {
					out[len(out)-1-j] = capitalize(out[len(out)-1-j])
				}
				i++
			}
		default:
			out = append(out, tokens[i])
		}
	}
	return out, nil
}

func getNum(tokens []string, i int) (int, int) {
	numStr := tokens[i+1][:len(tokens[i+1])-1]
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, 1
	}
	return num, 0
}

func fixArticle(tokens []string) []string {
	out := append([]string{}, tokens...)
	for i, tok := range out {
		if isStartWithVowelOrH(tok) && i-1 > 0 {
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
