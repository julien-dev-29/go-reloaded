package main

import (
	"fmt"
	"math"
	"strconv"
)

func ConvertHex2Dec(dec int, count int, result float64) float64 {
	if dec > 10 {
		return ConvertHex2Dec(dec/10, count+1, float64(dec%10)*math.Pow(16, float64(count)))
	}
	return (result + float64(dec%10)*math.Pow(16, float64(count)))
}

func Hex2Dec(data []byte) ([]byte, error) {
	fmt.Println(data)
	decimal, err := strconv.ParseInt(string(data), 16, 64)
	if err != nil {
		return nil, err
	}

	decimalStr := strconv.FormatInt(decimal, 10)

	data = []byte(decimalStr)
	fmt.Println(data)
	return data, nil
}

func Bin2Dec(bin int, count int, result float64) float64 {
	if bin > 0 {
		return Bin2Dec(bin/10, count+1, result+float64(bin%10)*math.Pow(2, float64(count)))
	}
	return result
}
