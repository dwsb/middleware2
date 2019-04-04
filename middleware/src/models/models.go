package models

import "fmt"

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

const (
	template string = `Name: %s
Description: %s
PublishDate: %s
Author: %s
Categories: %s
`
	bufferSize = 4096
)

func (b *Book) String() string {
	return fmt.Sprintf(template, b.Name, b.Description, b.PublishDate, b.Author, b.Categories)
}

type LoginRequest struct {
	Login    string
	Password string
}

type IsLoggedRequest struct {
	Token string
}

type ListRequest struct {
	Token string
}

type LoginResponse struct {
	Token string
}

type IsLoggedResponse struct {
	Result bool
}

type ListResponse struct {
	Books []*Book
}
