package backends

type Backend interface {
	GetUserID(token string) (string, error)
}
