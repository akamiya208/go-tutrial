package main

import (
	"log/slog"

	"github.com/akamiya208/go-tutrial/internal/pkg/models"
	"github.com/akamiya208/go-tutrial/internal/pkg/mysql"
)

func main() {
	mysqlClient, err := mysql.NewMySQLClient()
	if err != nil {
		slog.Error("failed to create mysql client")
		panic(err)
	}

	// Migrate the schema
	mysqlClient.DB().AutoMigrate(&models.Task{})
}
