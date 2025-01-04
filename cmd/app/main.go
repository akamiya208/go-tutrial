package main

import (
	"log/slog"
	"net/http"
	"os"

	handlers "github.com/akamiya208/go-tutrial/internal/app"
	"github.com/akamiya208/go-tutrial/internal/pkg/mysql"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	mysqlClient, err := mysql.NewMySQLClient()
	if err != nil {
		slog.Error("failed to create mysql client")
		panic(err)
	}

	taskHandler := handlers.NewTaskHandler(mysqlClient)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/tasks/{taskId}", taskHandler.HandleGetTask)
	mux.HandleFunc("PATCH /api/v1/tasks/{taskId}", taskHandler.HandleUpdateTask)
	mux.HandleFunc("DELETE /api/v1/tasks/{taskId}", taskHandler.HandleDeleteTask)
	mux.HandleFunc("GET /api/v1/tasks", taskHandler.HandleGetTasks)
	mux.HandleFunc("POST /api/v1/tasks", taskHandler.HandleCreateTask)

	slog.Info("starting server")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		slog.Error("failed to start server", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
