package library

import (
	"fmt"
)

type Library struct {
	Protocol string
	AuthPort string
}

type NotLoggedError struct {
}

func (e NotLoggedError) Error() string {
	return "User not logged in"
}

type User struct {
	Login    string
	Password string
	Token    string
}

type Book struct {
	Name        string
	Description string
	PublishDate string
	Author      string
	Categories  []string
}

type ServiceResponse struct {
	Books []*Book
	Error error
}

const template string = `Name: %s
Description: %s
PublishDate: %s
Author: %s
Categories: %s

`

func (b *Book) String() string {
	return fmt.Sprintf(template, b.Name, b.Description, b.PublishDate, b.Author, b.Categories)
}

func Books() []*Book {
	return []*Book{
		&Book{
			Name:        "Medicina Interna de Harrison - 2 Volumes",
			Description: "Apresentando os extraordinários avanços ocorridos em todas as áreas da medicina, esta nova edição do Harrison foi amplamente revisada para oferecer uma atualização completa sobre a patogênese das doenças, ensaios clínicos, técnicas de diagnóstico, diretrizes clínicas baseadas em evidências, tratamentos já estabelecidos e métodos recentemente aprovados",
			PublishDate: "21 nov 2016",
			Author:      " Dennis L. Kasper, Stephen L. Hauser, J. Larry Jameson, Anthony S. Fauci, Dan L. Longo, Joseph Loscalzo",
			Categories: []string{
				"Medicina",
				"Especialidades",
			},
		},
		&Book{
			Name:        "Netter Atlas de Anatomia Humana 7ª edição",
			Description: "É um dos nomes mais fortes mundialmente na área de Anatomia, reconhecido pela didática e clareza de suas ilustrações. Figuras modernas, que, em um volume, apresenta todo o corpo humano em descrições detalhadas e clinicamente relevantes",
			PublishDate: "8 dez 2018",
			Author:      "Frank H. Netter",
			Categories: []string{
				"Medicina",
				"Anatomia",
			},
		},
		&Book{
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
}
