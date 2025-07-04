package main

import (
	"log/slog"

	"github.com/jeffscottbrown/satchel/repository"
	"github.com/jeffscottbrown/satchel/server"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		slog.Info("Could not load .env file, proceeding without it", slog.Any("error", err))
	}
	repository.InitializeDatabase()
	server.Run()
}
