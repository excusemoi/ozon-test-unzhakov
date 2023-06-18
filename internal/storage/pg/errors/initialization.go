package errors

type InitializationErr struct {
	Err error
}

func (e InitializationErr) Error() string {
	return e.Err.Error()
}

type IncorrectInitializationErr struct{}

func (e IncorrectInitializationErr) Error() string {
	return "linkStoragePg: incorrect initialization"
}

type InitializedIncorrectlyErr struct{}

func (e InitializedIncorrectlyErr) Error() string {
	return "linkStoragePg: initialized incorrectly"
}
