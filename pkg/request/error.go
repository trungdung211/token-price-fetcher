package request

type AppError struct {
	StatusCode int
	Code       string
	Msg        string
}

func NewError(statusCode int, code string, msg string) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Code:       code,
		Msg:        msg,
	}
}

func (e *AppError) Error() string {
	return e.Msg
}
