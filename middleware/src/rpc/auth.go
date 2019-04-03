package main

import "middleware2/middleware/src/utils"

type Auth struct {
	isLogged map[string]string
	logins   []*User
}

func (t *Auth) Login(request LoginRequest, res *LoginResponse) error {
	result := validateLogin(request.Login, request.Password)

	if result == -1 {
		return NotFoundError{}
	}

	res.Token = generateToken()
	logins[result].Token = res.Token
	isLogged[res.Token] = "ok"

	return nil
}

func (t *Auth) IsLogged(request IsLoggedRequest, res *IsLoggedResponse) error {
	result := isLogged[utils.FormatString(request.Token)] != ""

	res.Result = result
	return nil
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
