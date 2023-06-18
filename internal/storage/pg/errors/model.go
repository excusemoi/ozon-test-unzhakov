package errors

type ModelErr struct {
	Err error
}

func (e ModelErr) Error() string {
	return e.Err.Error()
}

type ModelNilErr struct {
}

func (e ModelNilErr) Error() string {
	return "pg: Model(nil)"
}
