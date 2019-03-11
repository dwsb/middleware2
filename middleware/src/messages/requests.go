package messages

type Error struct {
	Error string
}

type UserAuthRequest struct {
	Login    string
	Password string
}
