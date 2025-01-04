package dto

type TaskCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
