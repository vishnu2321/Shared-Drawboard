package main

import (
	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/shared-drawboard/internal/handler"
	"github.com/shared-drawboard/pkg/logger"
)

func main() {
	logger.Info("Server is starting...")

	if err := godotenv.Load(); err != nil {
		logger.Error("Error loading environment variables: %s", err)
		os.Exit(1)
	}

	PORT := os.Getenv("PORT")
	if PORT == "" {
		logger.Error("PORT is not set")
		PORT = ":8080"
	}
	if PORT[0] != ':' {
		PORT = ":" + PORT
	}

	handler, err := handler.New()
	if err != nil {
		logger.Error("Error creating handler: %s", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		logger.Info("Server running on %s", PORT)
		if err := http.ListenAndServe(PORT, handler.Router); err != nil {
			logger.Error("Server stopped with error: %s", err)
		}
	}()

	wg.Wait()

}
