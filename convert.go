package main

import (
	"fmt"
	"math"
	"strconv"
)

func ConvertHex2Dec(dec int, count int, result float64) float64 {
	if dec > 10 {
		return ConvertHex2Dec(dec/10, count+1, float64((dec%10))*math.Pow(16, float64(count)))
	}
	return (result + float64((dec%10))*math.Pow(16, float64(count)))
}

func Hex2Dec(str string) int64 {
	decimal, err := strconv.ParseInt(str, 16, 64)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	return decimal
}
