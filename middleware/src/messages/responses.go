package messages

type UserAuthResponse struct {
	Token string
	Error string
}

type IsLoggedAuthResponse struct {
	Result bool
}
