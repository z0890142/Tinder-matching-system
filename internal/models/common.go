package models

type ErrorResponse struct {
	Error ErrorDetails `json:"error" validate:"required"`
}

type ErrorDetails struct {
	Code    string `json:"code" validate:"required"`
	Message string `json:"message" validate:"required"`
}

type CodeMessage struct {
	Code    int32  `validate:"required" example:"0"`
	Message string `validate:"required" example:"success"`
}
