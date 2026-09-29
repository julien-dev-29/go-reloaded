package main

import (
	"strconv"
)

func hex2Dec(s string) (string, bool) {
	v, err := strconv.ParseInt(s, 16, 64)
	if err != nil {
		return "", false
	}
	return strconv.FormatInt(v, 10), true
}

func bin2Dec(s string) (string, bool) {
	v, err := strconv.ParseInt(s, 2, 64)
	if err != nil {
		return "", false
	}
	return strconv.FormatInt(v, 10), true
}
