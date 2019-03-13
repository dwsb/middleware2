package utils

import (
	"fmt"
	"net"
	"strings"
)

func Encode(bytes []byte) []byte {
	return []byte(string(bytes) + "\n")
}

func EncodeString(value string) []byte {
	return []byte(value + "\n")
}

func OpenConnection(protocol, port string) (net.Conn, error) {
	return net.Dial(protocol, fmt.Sprintf(":%s", port))
}

func FormatString(s string) string {
	s = strings.Replace(s, "\n", "", 1)
	s = strings.Replace(s, "\r", "", 1)

	return s
}
