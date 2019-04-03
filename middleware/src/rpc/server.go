package main

import (
	"net"
	"net/rpc"
)

func Start(serverName string, serverClass *interface{}, protocol string, address string) error {
	server := rpc.NewServer()
	server.RegisterName(serverName, serverClass)

	listen, err := net.Listen(protocol, address)

	if err != nil {
		return err
	}

	go server.Accept(listen)
	return nil
}
