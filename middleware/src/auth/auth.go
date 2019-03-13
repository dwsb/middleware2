package auth

import "middleware2/middleware/src/library"

type NotFoundError struct {
}

func (e NotFoundError) Error() string {
	return "User not found"
}

type InvalidRequestError struct {
}

func (e InvalidRequestError) Error() string {
	return "Invalid Request"
}

var logins = []*library.User{
	&library.User{
		Login:    "diogo",
		Password: "nogueira",
	},
	&library.User{
		Login:    "marcela",
		Password: "azevedo",
	},
	&library.User{
		Login:    "luiz",
		Password: "reis",
	},
	&library.User{
		Login:    "edjan",
		Password: "michiles",
	},
}
