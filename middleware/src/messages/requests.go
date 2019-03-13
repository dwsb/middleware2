package messages

type Error struct {
	Error string
}

type UserAuthRequest struct {
	Login    string
	Password string
}

type ServiceRequest struct {
	Token string
}

type IsLoggedAuthRequest struct {
	Token string
}
