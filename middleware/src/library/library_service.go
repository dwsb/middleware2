package library

import (
	"bufio"
	"encoding/json"
	"fmt"
	"middleware2/middleware/src/messages"
	"middleware2/middleware/src/utils"
	"net"
	"strings"
)

var library Library

func StartServer(protocol, port, authPort string) {
	library = Library{
		Protocol: protocol,
		AuthPort: authPort,
	}

	if strings.ToLower(protocol) != "tcp" && strings.ToLower(protocol) != "udp" {
		fmt.Println("Protocolo inválido.")
		return
	}

	listener, err := net.Listen(protocol, fmt.Sprintf(":%s", port))
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
		fmt.Printf("Library Serving %s\n", conn.LocalAddr().String())

		action, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			return
		}

		action = strings.Replace(action, "\n", "", 1)
		conn.Write(utils.EncodeString("ok")) // send connection ok

		bytesRequest, err := bufio.NewReader(conn).ReadBytes('\n')
		if err != nil {
			return
		}

		switch action {
		case "list":
			var msgRequest messages.ServiceRequest
			err = json.Unmarshal(bytesRequest, &msgRequest)
			if err != nil {
				bytes, _ := json.Marshal(&messages.Error{Error: "Invalid request."})
				conn.Write(utils.Encode(bytes))
				continue
			}

			processList(conn, msgRequest)
		}
	}
}

func processList(conn net.Conn, request messages.ServiceRequest) {
	authConn, err := utils.OpenConnection(library.Protocol, library.AuthPort)
	if err != nil {
		return
	}

	defer authConn.Close()

	bytes, _ := json.Marshal(&messages.IsLoggedAuthRequest{Token: request.Token})

	authConn.Write(utils.EncodeString("isLogged"))
	bufio.NewReader(authConn).ReadBytes('\n') // wait ok connection from auth server
	authConn.Write(utils.Encode(bytes))

	bytesResponse, err := bufio.NewReader(authConn).ReadBytes('\n')
	if err != nil {
		return
	}

	var response messages.IsLoggedAuthResponse
	json.Unmarshal(bytesResponse, &response)

	var books []*Book

	if response.Result {
		books = Books()
	} else {
		err = NotLoggedError{}
	}

	bytes, _ = json.Marshal(&ServiceResponse{Books: books, Error: err})
	conn.Write(utils.Encode(bytes))
}
