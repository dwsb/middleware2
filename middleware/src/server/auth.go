package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"middleware2/middleware/src/messages"
	"middleware2/middleware/src/utils"
)

const ADDRESS = "localhost:1234"

func main() {
	protocol := os.Args[1]
	if strings.ToLower(protocol) != "tcp" && strings.ToLower(protocol) != "udp" {
		fmt.Println("Protocolo inválido.")
		return
	}

	listener, err := net.Listen(protocol, ADDRESS)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to open connection")
			continue
		}

		go handdleConnection(conn)
	}
}

func handdleConnection(conn net.Conn) {
	for {
		fmt.Printf("Serving %s\n", conn.LocalAddr().String())

		bytesRequest, err := bufio.NewReader(conn).ReadBytes('\n')
		if err != nil {
			return
		}

		var msgRequest messages.UserAuthRequest
		err = json.Unmarshal(bytesRequest, &msgRequest)
		if err != nil {
			bytes, _ := json.Marshal(&messages.Error{Error: "Invalid request."})
			conn.Write(utils.Encode(bytes))
			continue
		}

		processRequest(conn, msgRequest)
	}
}

func processRequest(conn net.Conn, request messages.UserAuthRequest) {
	result := true

	// login and generate token for the user
	fmt.Printf("Login: %s\nPassword: %s\n", request.Login, request.Password)

	bytes, _ := json.Marshal(&messages.UserAuthResponse{Result: result})
	conn.Write(utils.Encode(bytes))
}
