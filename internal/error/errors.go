package error

import "fmt"

type Http struct {
	Detail ErrorDetails `json:"error" validate:"required"`
}

type ErrorDetails struct {
	Code       string `json:"code" validate:"required"`
	Message    string `json:"message" validate:"required"`
	StatusCode int    `json:"-" validate:"required"`
}

func (e Http) Error() string {
	return fmt.Sprintf("Code : %s , Message : %s", e.Detail.Code, e.Detail.Message)
}

func NewHttpError(StatusCode int, code, message string) Http {
	return Http{
		Detail: ErrorDetails{
			StatusCode: StatusCode,
			Code:       code,
			Message:    message,
		},
	}
}
