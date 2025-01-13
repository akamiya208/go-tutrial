package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/akamiya208/go-tutrial/internal/pkg/dto"
	"github.com/akamiya208/go-tutrial/internal/pkg/models"
	"github.com/akamiya208/go-tutrial/internal/pkg/mysql"
)

type TaskHandler struct {
	mysqlClient mysql.IClient
}

func NewTaskHandler(client mysql.IClient) *TaskHandler {
	return &TaskHandler{mysqlClient: client}
}

// @Summary		get a task
// @Description	get task by taskId
// @Tags			tasks
// @Accept			json
// @Produce		json
// @Param			taskId	path		uint	true	"Task ID"
// @Success		200		{object}	dto.TaskResponse
// @Failure		400		{object}	dto.ErrorResponse
// @Failure		404		{object}	dto.ErrorResponse
// @Failure		500		{object}	dto.ErrorResponse
// @Router			/api/v1/tasks/{taskId} [get]
func (h *TaskHandler) HandleGetTask(w http.ResponseWriter, r *http.Request) {
	taskID, err := strconv.Atoi(r.PathValue("taskId"))
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	task, err := h.mysqlClient.GetTask(uint(taskID))
	if err != nil {
		if err.Error() == "record not found" {
			writeErrorResponse(w, http.StatusNotFound, err.Error())
		} else {
			writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	writeJSONResponse(w, http.StatusOK, dto.ToTaskResponse(task))
}

// @Summary		list tasks
// @Description	list tasks by task name
// @Tags			tasks
// @Accept			json
// @Produce		json
// @Param			name	query		string	true	"search by name"
// @Success		200		{object}	[]dto.TaskResponse
// @Failure		400		{object}	dto.ErrorResponse
// @Failure		404		{object}	dto.ErrorResponse
// @Failure		500		{object}	dto.ErrorResponse
// @Router			/api/v1/tasks [get]
func (h *TaskHandler) HandleGetTasks(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		writeErrorResponse(w, http.StatusBadRequest, "name is required")
		return
	}

	tasks, err := h.mysqlClient.GetTasksByName(name)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses := make([]dto.TaskResponse, len(tasks))
	for i, task := range tasks {
		responses[i] = dto.ToTaskResponse(task)
	}
	writeJSONResponse(w, http.StatusOK, responses)
}

// @Summary		create a task
// @Description	create a task
// @Tags			tasks
// @Accept			json
// @Produce		json
// @Param			reqest	body		dto.TaskCreateRequest	true	"Create task"
// @Success		201		{object}	dto.TaskResponse
// @Failure		400		{object}	dto.ErrorResponse
// @Failure		404		{object}	dto.ErrorResponse
// @Failure		500		{object}	dto.ErrorResponse
// @Router			/api/v1/tasks [post]
func (h *TaskHandler) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	var request dto.TaskCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	task := models.Task{Name: request.Name, Description: &request.Description}
	if err := h.mysqlClient.CreateTask(&task); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSONResponse(w, http.StatusCreated, dto.ToTaskResponse(task))
}

// @Summary		update a task
// @Description	update a task
// @Tags			tasks
// @Accept			json
// @Produce		json
// @Param			taskId	path		uint	true	"Task ID"
// @Param			reqest	body		dto.TaskUpdateRequest	true	"Update task"
// @Success		200		{object}	dto.TaskResponse
// @Failure		400		{object}	dto.ErrorResponse
// @Failure		404		{object}	dto.ErrorResponse
// @Failure		500		{object}	dto.ErrorResponse
// @Router			/api/v1/tasks/{taskId} [patch]
func (h *TaskHandler) HandleUpdateTask(w http.ResponseWriter, r *http.Request) {
	taskID, err := strconv.Atoi(r.PathValue("taskId"))
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	task, err := h.mysqlClient.GetTask(uint(taskID))
	if err != nil {
		if err.Error() == "record not found" {
			writeErrorResponse(w, http.StatusNotFound, err.Error())
		} else {
			writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	var request dto.TaskUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	task.Name = request.Name
	task.Description = &request.Description

	if err := h.mysqlClient.UpdateTask(&task); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSONResponse(w, http.StatusOK, dto.ToTaskResponse(task))
}

// @Summary		delete a task
// @Description	delete a task by taskId
// @Tags			tasks
// @Accept			json
// @Produce		json
// @Param			taskId	path	uint	true	"Task ID"
// @Success		204
// @Failure		400	{object}	dto.ErrorResponse
// @Failure		404	{object}	dto.ErrorResponse
// @Failure		500	{object}	dto.ErrorResponse
// @Router			/api/v1/tasks/{taskId} [delete]
func (h *TaskHandler) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	taskID, err := strconv.Atoi(r.PathValue("taskId"))
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	task, err := h.mysqlClient.GetTask(uint(taskID))
	if err != nil {
		if err.Error() == "record not found" {
			writeErrorResponse(w, http.StatusNotFound, err.Error())
		} else {
			writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	if err := h.mysqlClient.DeleteTask(&task); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusNoContent, nil)
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := dto.ErrorResponse{Status: http.StatusText(statusCode), Detail: message}
	json.NewEncoder(w).Encode(response)
}

func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
