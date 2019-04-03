package main

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"
	"strings"

	"middleware2/middleware/src/library"
	"middleware2/middleware/src/utils"
)

var client c.Client
var user library.User

const PROTOCOL = "tcp"
const HOST = "localhost:"

func main() {
	authPort := os.Args[1]
	libraryPort := os.Args[2]

	auth := new(Auth{})
	library := new(Library{})

	server.Start("Auth", auth, PROTOCOL, HOST+authPort)
	server.Start("Library", library, PROTOCOL, HOST+libraryPort)

	authClient, authErr := rpc.Dial(PROTOCOL, HOST+authPort)
	libraryClient, libraryErr := rpc.Dial(PROTOCOL, HOST+libraryPort)

	if authErr != nil {
		fmt.Println(authErr)
		return
	}

	if libraryErr != nil {
		fmt.Println(libraryErr)
		return
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

	var loginResponse LoginResponse
	loginErr := authClient.Call("Auth.Login", LoginRequest{Login: login, Password: password}, &loginResponse)

	if loginErr != nil {
		fmt.Println(loginErr)
		return
	}

	fmt.Println("User logged in successfuly")
	user.Token = loginResponse.Token
}

func books() {
	var listResponse ListResponse
	libraryErr := libraryClient.Call("Library.List", ListRequest{Token: user.token}, &listResponse)

	if libraryErr != nil {
		fmt.Println(libraryErr)
		return
	}

	fmt.Println(listResponse.Books)
}
