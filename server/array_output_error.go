package server

type ArrayOutputError struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func NewArrayOutputError(err error) ArrayOutputError {
	code := 500
	if e, ok := err.(HttpError); ok {
		code = e.Code()
	}

	return ArrayOutputError{
		Code:  code,
		Error: err.Error(),
	}
}
