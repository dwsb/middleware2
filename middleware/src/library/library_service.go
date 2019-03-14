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

	switch strings.ToLower(protocol) {
	case "tcp":
		TCPServer(protocol, port)
	case "udp":
		UDPServer(protocol, port)
	default:
		fmt.Println("Protocolo inválido.")
		return
	}
}

func TCPServer(protocol, port string) {
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

func UDPServer(protocol, port string) {
	pc, err := net.ListenPacket("udp", fmt.Sprintf(":%s", port))
	if err != nil {
		return
	}

	go handdleUDPConnection(pc)
}

func handdleUDPConnection(conn net.PacketConn) {
	defer conn.Close()

	for {
		fmt.Printf("Library Serving %s\n", conn.LocalAddr().String())
		bytesAction := make([]byte, bufferSize)

		nBytes, addr, err := conn.ReadFrom(bytesAction)
		if err != nil {
			return
		}

		_, err = conn.WriteTo(utils.EncodeString("ok"), addr)
		if err != nil {
			return
		}

		action := string(bytesAction[:nBytes])
		action = strings.Replace(action, "\n", "", 1)

		bytesRequest := make([]byte, bufferSize)

		nBytes, addr, err = conn.ReadFrom(bytesRequest)
		if err != nil {
			return
		}

		switch action {
		case "list":
			var msgRequest messages.ServiceRequest

			err = json.Unmarshal(bytesRequest[:nBytes], &msgRequest)
			if err != nil {
				bytes, _ := json.Marshal(&messages.Error{Error: "Invalid request."})
				conn.WriteTo(utils.Encode(bytes), addr)
				continue
			}

			bytes := processList(msgRequest)
			conn.WriteTo(utils.Encode(bytes), addr)
		}
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

			bytes := processList(msgRequest)
			conn.Write(utils.Encode(bytes))
		}
	}
}

func processList(request messages.ServiceRequest) []byte {
	authConn, err := utils.OpenConnection(library.Protocol, library.AuthPort)
	if err != nil {
		return nil
	}

	defer authConn.Close()

	bytes, _ := json.Marshal(&messages.IsLoggedAuthRequest{Token: request.Token})

	authConn.Write(utils.EncodeString("isLogged"))

	bufio.NewReader(authConn).ReadBytes('\n') // wait ok connection from auth server

	authConn.Write(utils.Encode(bytes))

	bytesResponse, err := bufio.NewReader(authConn).ReadBytes('\n')
	if err != nil {
		return nil
	}

	var response messages.IsLoggedAuthResponse
	json.Unmarshal(bytesResponse, &response)

	var books []*Book
	var errorMessage string

	if response.Result {
		books = Books()
	} else {
		errorMessage = NotLoggedError{}.Error()
	}

	bytes, _ = json.Marshal(ServiceResponse{Books: books, Error: errorMessage})
	return bytes
}
