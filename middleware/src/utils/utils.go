package utils

func Encode(bytes []byte) []byte {
	return []byte(string(bytes) + "\n")
}
