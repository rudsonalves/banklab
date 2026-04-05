package sharederrors

type AppError struct {
	Code    string
	Message string
	Status  int
}

func (e AppError) Error() string {
	return e.Message
}
