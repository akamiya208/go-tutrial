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
