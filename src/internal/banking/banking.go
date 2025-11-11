package banking

type Banking interface {
	StartSession() (string, error)
	EndSession() error
}
