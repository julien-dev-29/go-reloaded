package main

func ToLowercase(bytes []byte) []byte {
	for i := range bytes {
		if bytes[i] >= 'A' && bytes[i] <= 'Z' {
			bytes[i] = bytes[i] + 32
		}
	}
	return bytes
}
