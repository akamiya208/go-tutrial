package dto

import "github.com/akamiya208/go-tutrial/internal/pkg/models"

type TaskResponse struct {
	ID          uint   `json:"id" example:"1"`
	Name        string `json:"name" example:"taskName"`
	Description string `json:"description" example:"taskDescription"`
}

func ToTaskResponse(task models.Task) TaskResponse {
	return TaskResponse{
		ID:          task.ID,
		Name:        task.Name,
		Description: *task.Description,
	}
}
