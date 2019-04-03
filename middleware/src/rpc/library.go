package main

import "net/rpc"

type Library struct {
	books        []*Book
	authProtocol string
	authAddress  string
}

func (t *Library) List(request ListRequest, res *ListResponse) error {
	client, err := rpc.Dial(authProtocol, authAddress)

	if err != nil {
		return err
	}

	var response IsLoggedResponse
	err = client.Call("Auth.IsLogged", IsLoggedRequest{Token: request.Token}, &response)

	if err != nil {
		return err
	}

	if !response.Result {
		return NotLoggedError{}
	}

	res.Books = books
	return nil
}
