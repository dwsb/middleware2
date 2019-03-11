package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"middleware2/middleware/src/messages"
	"middleware2/middleware/src/utils"
	"net"
	"os"
	"strings"
)

func main() {
	protocol := os.Args[1]
	host := os.Args[2]
	port := os.Args[3]

	conn, err := net.Dial(protocol, fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()
	handleConnection(conn)
}

func handleConnection(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Login: ")
	login, _ := reader.ReadString('\n')
	login = strings.Replace(login, "\n", "", -1)

	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')
	password = strings.Replace(password, "\n", "", -1)

	request := messages.UserAuthRequest{
		Login:    login,
		Password: password,
	}

	bytes, _ := json.Marshal(&request)
	conn.Write(utils.Encode(bytes))

	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(message)
}
