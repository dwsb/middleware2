package auth

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"

	uuid "github.com/satori/go.uuid"

	"middleware2/middleware/src/messages"
	"middleware2/middleware/src/utils"
)

var isLogged = make(map[string]string, 4)

// StartServer ... starts the auth server
func StartServer(protocol, port string) {

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
		fmt.Printf("Auth Serving %s\n", conn.LocalAddr().String())

		action, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			return
		}

		action = strings.Replace(action, "\n", "", 1)
		conn.Write(utils.EncodeString("ok")) // send ok connection

		bytesRequest, err := bufio.NewReader(conn).ReadBytes('\n')
		if err != nil {
			return
		}

		switch action {
		case "login":
			var msgRequest messages.UserAuthRequest
			err = json.Unmarshal(bytesRequest, &msgRequest)
			if err != nil {
				bytes, _ := json.Marshal(&messages.UserAuthResponse{Error: InvalidRequestError{}})
				conn.Write(utils.Encode(bytes))
				continue
			}

			processLogin(conn, msgRequest)
		case "isLogged":
			var msgRequest messages.IsLoggedAuthRequest
			err = json.Unmarshal(bytesRequest, &msgRequest)
			if err != nil {
				// adicionar erro ao messages.IsLoggedAuthResponse
				bytes, _ := json.Marshal(messages.IsLoggedAuthResponse{Result: false})
				conn.Write(utils.Encode(bytes))
				continue
			}

			processIsLogged(conn, msgRequest)
		}
	}
}

func processIsLogged(conn net.Conn, request messages.IsLoggedAuthRequest) {
	result := isLogged[utils.FormatString(request.Token)] != ""

	bytes, _ := json.Marshal(messages.IsLoggedAuthResponse{Result: result})
	conn.Write(utils.Encode(bytes))
}

func processLogin(conn net.Conn, request messages.UserAuthRequest) {
	result := validateLogin(request.Login, request.Password)

	var token string
	var err error

	if result != -1 {
		token = generateToken()
		logins[result].Token = token
		isLogged[token] = "ok"
		err = nil
	} else {
		err = NotFoundError{}
	}

	bytes, _ := json.Marshal(messages.UserAuthResponse{Token: token, Error: err})
	conn.Write(utils.Encode(bytes))
}

func validateLogin(login, password string) int {
	for i, user := range logins {
		if user.Login == utils.FormatString(login) && user.Password == utils.FormatString(password) {
			return i
		}
	}

	return -1
}

func generateToken() string {
	id, _ := uuid.NewV4()
	return id.String()
}
