package main

func ToUppercase(bytes []byte) []byte {
	for i := range bytes {
		if bytes[i] >= 'a' && bytes[i] <= 'z' {
			bytes[i] = bytes[i] - 32
		}
	}
	return bytes
}
