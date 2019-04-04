package library

import (
	"encoding/json"
	"errors"
	"middleware2/middleware/src/models"
	"middleware2/middleware/src/rabbitmq/consumer"
	"middleware2/middleware/src/rabbitmq/producer"
	"net/rpc"
)

const CONTENT_TYPE = "application/json"

type Library struct {
	books       []*models.Book
	authAddress string
}

func New(books []*models.Book, authAddress string) *Library {
	library := new(Library)
	library.books = books
	library.authAddress = authAddress

	return library
}

func NewMQ(books []*models.Book) *Library {
	library := new(Library)
	library.books = books

	return library
}

func (l *Library) ListRPC(request models.ListRequest, res *models.ListResponse) error {
	client, err := rpc.Dial("tcp", l.authAddress)

	if err != nil {
		return err
	}

	var response models.IsLoggedResponse
	err = client.Call("Auth.IsLogged", models.IsLoggedRequest{Token: request.Token}, &response)

	if err != nil {
		return err
	}

	if !response.Result {
		return errors.New("Usuario not logged")
	}

	res.Books = l.books
	return nil
}

func (l *Library) ListRabbitMQ() error {
	consumer := new(consumer.Consumer)
	consumer.Connect("list")
	answer, err := consumer.Consume(CONTENT_TYPE)
	if err != nil {
		return err
	}

	var request models.ListRequest
	err = json.Unmarshal(answer, &request)
	if err != nil {
		return err
	}

	producer := new(producer.Producer)
	err = producer.Connect("isLogged")
	if err != nil {
		return err
	}

	bytes, _ := json.Marshal(request)
	reply, err := producer.ProduceAndWaitReply(bytes, "application/json")
	if err != nil {
		return err
	}

	var isLogged models.IsLoggedResponse
	json.Unmarshal(reply, &isLogged)

	if !isLogged.Result {
		return errors.New("não está logado")
	}

	var res *models.ListResponse
	res.Books = l.books

	bytes, _ = json.Marshal(res)
	return consumer.Reply(bytes, CONTENT_TYPE)
}
