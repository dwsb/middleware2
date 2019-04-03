package main

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
