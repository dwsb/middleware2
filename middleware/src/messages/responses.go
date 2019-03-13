package messages

type UserAuthResponse struct {
	Token string
	Error error
}

type IsLoggedAuthResponse struct {
	Result bool
}
