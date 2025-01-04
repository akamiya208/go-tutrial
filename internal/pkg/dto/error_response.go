package dto

type ErrorResponse struct {
	Status string `json:"status"`
	Detail string `json:"detail"`
}
