package main

import "middleware2/middleware/src/messages"

type DbAuth struct {
	Users []*messages.UserAuthRequest
}
