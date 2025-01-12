package dto

type TaskUpdateRequest struct {
	Name        string `json:"name" example:"taskUpdateName"`
	Description string `json:"description" example:"taskUpdateDescription"`
}
