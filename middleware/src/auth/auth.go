package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"middleware2/middleware/src/models"
	"middleware2/middleware/src/rabbitmq/consumer"
	"middleware2/middleware/src/utils"

	uuid "github.com/satori/go.uuid"
)

const CONTENT_TYPE = "application/json"

type Auth struct {
	logins   []*models.User
	isLogged map[string]string
}

func New(logins []*models.User) *Auth {
	auth := new(Auth)
	auth.logins = logins
	auth.isLogged = make(map[string]string)

	return auth
}

func (a *Auth) Login(request models.LoginRequest, res *models.LoginResponse) error {
	result := a.validateLogin(request.Login, request.Password)

	if result == -1 {
		return errors.New("Not found")
	}

	res.Token = generateToken()
	a.logins[result].Token = res.Token
	a.isLogged[res.Token] = "ok"

	return nil
}

func (a *Auth) IsLogged(request models.IsLoggedRequest, res *models.IsLoggedResponse) error {
	result := a.isLogged[utils.FormatString(request.Token)] != ""

	res.Result = result
	return nil
}

func (a *Auth) LoginRabbitMQ() error {
	consumer := new(consumer.Consumer)
	consumer.Connect("login")
	answer, err := consumer.Consume(CONTENT_TYPE)
	if err != nil {
		return err
	}

	fmt.Println(string(answer))

	var request models.LoginRequest
	err = json.Unmarshal(answer, &request)
	if err != nil {
		bytes, _ := json.Marshal(models.LoginResponse{Token: ""})
		consumer.Reply(bytes, CONTENT_TYPE)
		return nil
	}

	fmt.Println(request.Login)

	var response models.LoginResponse
	err = a.Login(request, &response)

	fmt.Println(response)
	bytes, _ := json.Marshal(response)
	consumer.Reply(bytes, CONTENT_TYPE)

	return err
}

func (a *Auth) IsLoggedRabbitMQ() error {
	fmt.Println("entrou")
	consumer := consumer.Consumer{}
	consumer.Connect("IsLogged")

	answer, err := consumer.Consume(CONTENT_TYPE)
	if err != nil {
		return err
	}

	fmt.Println(string(answer))

	var request models.IsLoggedRequest
	err = json.Unmarshal(answer, &request)
	if err != nil {
		bytes, _ := json.Marshal(models.IsLoggedResponse{Result: false})
		consumer.Reply(bytes, CONTENT_TYPE)
		return err
	}

	fmt.Println(string(request.Token))
	var response models.IsLoggedResponse
	a.IsLogged(request, &response)

	bytes, _ := json.Marshal(response)
	return consumer.Reply(bytes, CONTENT_TYPE)
}

func (a *Auth) validateLogin(login, password string) int {
	for i, user := range a.logins {
		if user.Login == utils.FormatString(login) && user.Password == utils.FormatString(password) {
			return i
		}
	}

	return -1
}

func generateToken() string {
	return uuid.NewV4().String()
}
