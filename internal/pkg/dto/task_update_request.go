package dto

type TaskUpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
