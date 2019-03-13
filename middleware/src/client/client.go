package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"middleware2/middleware/src/library"
	"middleware2/middleware/src/messages"
	"middleware2/middleware/src/utils"
	"net"
)

type Client struct {
	Protocol    string
	ServicePort string
	AuthPort    string

	AuthConnection    net.Conn
	ServiceConnection net.Conn
}

func (c *Client) openServiceConnection(protocol, port string) {
	if c.ServiceConnection != nil {
		return
	}

	conn, err := utils.OpenConnection(protocol, port)
	if err != nil {
		fmt.Println(err)
	}

	c.ServiceConnection = conn
}

func (c *Client) openAuthConnection(protocol, port string) {
	if c.AuthConnection != nil {
		return
	}

	conn, err := utils.OpenConnection(protocol, port)
	if err != nil {
		fmt.Println(err)
	}

	c.AuthConnection = conn
}

func (c *Client) closeAuthConnection() {
	c.AuthConnection.Close()
	c.AuthConnection = nil
}

func (c *Client) closeServiceConnection() {
	c.ServiceConnection.Close()
	c.ServiceConnection = nil
}

func (c *Client) Login(login, password string) (messages.UserAuthResponse, error) {
	c.openAuthConnection(c.Protocol, c.AuthPort)
	defer c.closeAuthConnection()

	request := messages.UserAuthRequest{
		Login:    login,
		Password: password,
	}

	bytes, _ := json.Marshal(request)

	c.AuthConnection.Write(utils.EncodeString("login"))
	bufio.NewReader(c.AuthConnection).ReadBytes('\n') // wait ok connection from authServer
	c.AuthConnection.Write(utils.Encode(bytes))

	message, err := bufio.NewReader(c.AuthConnection).ReadBytes('\n')
	if err != nil {
		return messages.UserAuthResponse{}, err
	}

	var response messages.UserAuthResponse
	err = json.Unmarshal(message, &response)

	return response, err
}

func (c *Client) Books(token string) (library.ServiceResponse, error) {
	c.openServiceConnection(c.Protocol, c.ServicePort)
	defer c.closeServiceConnection()

	request := messages.ServiceRequest{
		Token: token,
	}

	bytes, _ := json.Marshal(&request)
	c.ServiceConnection.Write(utils.EncodeString("list"))
	bufio.NewReader(c.ServiceConnection).ReadBytes('\n') // wait ok connection from library server
	c.ServiceConnection.Write(utils.Encode(bytes))

	message, err := bufio.NewReader(c.ServiceConnection).ReadBytes('\n')
	if err != nil {
		return library.ServiceResponse{}, err
	}

	var response library.ServiceResponse
	err = json.Unmarshal(message, &response)

	return response, err
}
