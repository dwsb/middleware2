package auth

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"

	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"

	"middleware2/middleware/src/messages"
	"middleware2/middleware/src/utils"
)

var isLogged = make(map[string]string, 4)

// StartServer ... starts the auth server
func StartServer(protocol, port string) {
	switch strings.ToLower(protocol) {
	case "tcp":
		TCPServer(protocol, port)
	case "udp":
		UDPServer(protocol, port)
	case "rabbitmq":
		RabbitMQServer()
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

func RabbitMQServer() {
	ch, err := utils.ConnectRabbitMQ()
	if err != nil {
		fmt.Println(err)
		return
	}

	q, err := utils.DeclareQueue("auth_service", ch)
	if err != nil {
		fmt.Println(err)
		return
	}

	msgs, err := utils.ConsumeQueue(q.Name, ch)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		go consumeMessages(msgs)
	}

}

func consumeMessages(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		switch string(msg.Body)[:5] {
		case "login":
			fmt.Println("fazer login")
		case "verif":
			fmt.Println("verificar login")
		}
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
				bytes, _ := json.Marshal(&messages.UserAuthResponse{Error: InvalidRequestError{}.Error()})
				conn.Write(utils.Encode(bytes))
				continue
			}

			bytes := processLogin(msgRequest)
			conn.Write(utils.Encode(bytes))

		case "isLogged":
			var msgRequest messages.IsLoggedAuthRequest

			err = json.Unmarshal(bytesRequest, &msgRequest)
			if err != nil {
				// adicionar erro ao messages.IsLoggedAuthResponse
				bytes, _ := json.Marshal(messages.IsLoggedAuthResponse{Result: false})
				conn.Write(utils.Encode(bytes))
				continue
			}

			bytes := processIsLogged(msgRequest)
			conn.Write(utils.Encode(bytes))
		}
	}
}

func handdleUDPConnection(conn net.PacketConn) {
	defer conn.Close()

	for {
		fmt.Printf("Auth Serving %s\n", conn.LocalAddr().String())

		bytesAction := make([]byte, bufferSize)
		nBytes, addr, err := conn.ReadFrom(bytesAction)
		if err != nil {
			fmt.Println(err)
			return
		}

		action := string(bytesAction[:nBytes])
		action = utils.FormatString(action)

		conn.WriteTo(utils.EncodeString("ok"), addr) // send ok connection

		bytesRequest := make([]byte, bufferSize)
		nBytes, addr, err = conn.ReadFrom(bytesRequest)
		if err != nil {
			return
		}

		switch action {
		case "login":
			var msgRequest messages.UserAuthRequest

			err = json.Unmarshal(bytesRequest[:nBytes], &msgRequest)
			if err != nil {
				bytes, _ := json.Marshal(&messages.UserAuthResponse{Error: InvalidRequestError{}.Error()})
				conn.WriteTo(utils.Encode(bytes), addr)
				continue
			}

			bytes := processLogin(msgRequest)
			conn.WriteTo(utils.Encode(bytes), addr)
		case "isLogged":
			var msgRequest messages.IsLoggedAuthRequest

			err = json.Unmarshal(bytesRequest[:nBytes], &msgRequest)
			if err != nil {
				bytes, _ := json.Marshal(messages.IsLoggedAuthResponse{Result: false})
				conn.WriteTo(utils.Encode(bytes), addr)
				continue
			}

			bytes := processIsLogged(msgRequest)
			conn.WriteTo(utils.Encode(bytes), addr)
		}
	}
}

func processIsLogged(request messages.IsLoggedAuthRequest) []byte {
	result := isLogged[utils.FormatString(request.Token)] != ""

	bytes, _ := json.Marshal(messages.IsLoggedAuthResponse{Result: result})
	return bytes
}

func processLogin(request messages.UserAuthRequest) []byte {
	result := validateLogin(request.Login, request.Password)

	var token string
	var err string

	if result != -1 {
		token = generateToken()

		logins[result].Token = token
		isLogged[token] = "ok"
	} else {
		err = NotFoundError{}.Error()
	}

	bytes, _ := json.Marshal(messages.UserAuthResponse{Token: token, Error: err})
	return bytes
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
	return uuid.NewV4().String()
}
