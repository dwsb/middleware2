package server

import (
	"net"
	"net/rpc"
)

func Start(serverName string, serverClass interface{}, address string) error {
	server := rpc.NewServer()
	server.RegisterName(serverName, serverClass)

	listen, err := net.Listen("tcp", address)

	if err != nil {
		return err
	}

	go server.Accept(listen)
	return nil
}
