package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"middleware2/middleware/src/auth"
	c "middleware2/middleware/src/client"
	"middleware2/middleware/src/library"
	"middleware2/middleware/src/utils"
)

var client c.Client
var user library.User

func main() {
	protocol := os.Args[1]
	authPort := os.Args[2]
	servicePort := os.Args[3]

	go auth.StartServer(protocol, authPort)
	go library.StartServer(protocol, servicePort, authPort)

	client = c.Client{
		Protocol:    protocol,
		AuthPort:    authPort,
		ServicePort: servicePort,
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Escolha uma das opções abaixo: ")
		fmt.Println("1 - Login")
		fmt.Println("2 - Listar livros")

		option, _ := reader.ReadString('\n')
		option = utils.FormatString(option)

		switch option {
		case "1":
			login()
		case "2":
			books()
		default:
			fmt.Printf("%s é uma opção inválida.\n", option)
		}
	}
}

func login() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Login: ")
	login, _ := reader.ReadString('\n')
	login = strings.Replace(login, "\n", "", -1)

	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')
	password = strings.Replace(password, "\n", "", -1)

	user.Login = login
	user.Password = password

	response := client.Login(login, password)
	if response.Error != "" {
		fmt.Println(response.Error)
		return
	}

	fmt.Println("User logged in successfuly")
	user.Token = response.Token
}

func books() {
	booksResponse, err := client.Books(user.Token)
	if err != nil {
		fmt.Println(err)
		return
	}

	if booksResponse.Error != "" {
		fmt.Println(booksResponse.Error)
		return
	}

	fmt.Println(booksResponse.Books)
}
