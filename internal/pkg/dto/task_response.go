package dto

import "github.com/akamiya208/go-tutrial/internal/pkg/models"

type TaskResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func ToTaskResponse(task models.Task) TaskResponse {
	return TaskResponse{
		ID:          task.ID,
		Name:        task.Name,
		Description: *task.Description,
	}
}
