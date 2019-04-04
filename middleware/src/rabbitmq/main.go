package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"middleware2/middleware/src/auth"
	"middleware2/middleware/src/library"
	"os"
	"strings"

	"middleware2/middleware/src/models"
	"middleware2/middleware/src/rabbitmq/producer"
	"middleware2/middleware/src/utils"
)

const CONTENT_TYPE = "application/json"

var (
	user   models.User
	client *producer.Producer
)

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
	client = new(producer.Producer)

	authServer := auth.New(LOGINS)
	libraryServer := library.NewMQ(BOOKS)

	go authServer.LoginRabbitMQ()
	go authServer.IsLoggedRabbitMQ()
	go libraryServer.ListRabbitMQ()

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

	request := models.LoginRequest{Login: login, Password: password}
	bytes, _ := json.Marshal(request)
	client.Connect("login")
	defer client.Close()
	fmt.Println("Login connected")
	bytes, err := client.ProduceAndWaitReply(bytes, CONTENT_TYPE)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Get login response")
	var loginResponse models.LoginResponse
	loginErr := json.Unmarshal(bytes, &loginResponse)
	if loginErr != nil {
		fmt.Println(loginErr)
		return
	}

	fmt.Println("User logged in successfuly. Token: " + loginResponse.Token)
	user.Token = loginResponse.Token

}

func books() {
	request := models.ListRequest{Token: user.Token}
	bytes, _ := json.Marshal(request)
	client.Connect("list")
	defer client.Close()

	bytes, err := client.ProduceAndWaitReply(bytes, CONTENT_TYPE)
	if err != nil {
		fmt.Println(err)
		return
	}

	var listResponse models.ListResponse
	err = json.Unmarshal(bytes, &listResponse)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(listResponse.Books)
}
