package model

type ErrorResponse struct {
	Version string `json:"version"`
	Message string `json:"Message"`
}

func NewErrorResponse(msg string) ErrorResponse {
	return ErrorResponse{
		Version: "1.0.0",
		Message: msg,
	}
}