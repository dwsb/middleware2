package main

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"
	"strings"

	"middleware2/middleware/src/auth"
	"middleware2/middleware/src/library"
	"middleware2/middleware/src/models"
	"middleware2/middleware/src/rpc/server"
	"middleware2/middleware/src/utils"
)

var (
	user          models.User
	authClient    *rpc.Client
	libraryClient *rpc.Client
)

const PROTOCOL = "tcp"
const HOST = "localhost:"

var LOGINS = []*models.User{
	&models.User{
		Login:    "diogo",
		Password: "nogueira",
	},
	&models.User{
		Login:    "marcela",
		Password: "azevedo",
	},
	&models.User{
		Login:    "luiz",
		Password: "reis",
	},
	&models.User{
		Login:    "edjan",
		Password: "michiles",
	},
}

var BOOKS = []*models.Book{
	&models.Book{
		Name:        "Medicina Interna de Harrison - 2 Volumes",
		Description: "Apresentando os extraordinários avanços ocorridos em todas as áreas da medicina, esta nova edição do Harrison foi amplamente revisada para oferecer uma atualização completa sobre a patogênese das doenças, ensaios clínicos, técnicas de diagnóstico, diretrizes clínicas baseadas em evidências, tratamentos já estabelecidos e métodos recentemente aprovados",
		PublishDate: "21 nov 2016",
		Author:      " Dennis L. Kasper, Stephen L. Hauser, J. Larry Jameson, Anthony S. Fauci, Dan L. Longo, Joseph Loscalzo",
		Categories: []string{
			"Medicina",
			"Especialidades",
		},
	},
	&models.Book{
		Name:        "Netter Atlas de Anatomia Humana 7ª edição",
		Description: "É um dos nomes mais fortes mundialmente na área de Anatomia, reconhecido pela didática e clareza de suas ilustrações. Figuras modernas, que, em um volume, apresenta todo o corpo humano em descrições detalhadas e clinicamente relevantes",
		PublishDate: "8 dez 2018",
		Author:      "Frank H. Netter",
		Categories: []string{
			"Medicina",
			"Anatomia",
		},
	},
	&models.Book{
		Name:        "Cultura Inglesa. Go Beyond - Caixa com Worbook: Student's Pack With Worbook",
		Description: "Go Beyond is an exciting 6-level American English course for teenagers learning English. The course covers CEFR levels A1+ through to B2, + all levels being based on mapping of the requirements of the CEFR and international exams.",
		PublishDate: "2 out 2018",
		Author:      "Rebbeca Robb Benne",
		Categories: []string{
			"Inglês e outras línguas",
			"Educação",
			"Didáticos",
		},
	},
}

func main() {
	authPort := os.Args[1]
	libraryPort := os.Args[2]

	auth := auth.New(LOGINS)
	library := library.New(BOOKS, HOST+authPort)

	server.Start("Auth", auth, HOST+authPort)
	server.Start("Library", library, HOST+libraryPort)

	var authErr, libraryErr error
	authClient, authErr = rpc.Dial(PROTOCOL, HOST+authPort)
	libraryClient, libraryErr = rpc.Dial(PROTOCOL, HOST+libraryPort)

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

	var loginResponse models.LoginResponse
	loginErr := authClient.Call("Auth.Login", models.LoginRequest{Login: login, Password: password}, &loginResponse)

	if loginErr != nil {
		fmt.Println(loginErr)
		return
	}

	fmt.Println("User logged in successfuly")
	user.Token = loginResponse.Token
}

func books() {
	var listResponse models.ListResponse
	libraryErr := libraryClient.Call("Library.ListRPC", models.ListRequest{Token: user.Token}, &listResponse)

	if libraryErr != nil {
		fmt.Println(libraryErr)
		return
	}

	fmt.Println(listResponse.Books)
}
