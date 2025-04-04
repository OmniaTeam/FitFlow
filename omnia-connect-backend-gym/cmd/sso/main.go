package main

import (
	"api-gateway/internal/app"
	"log/slog"
)

func main() {
	ap, err := app.NewApp()
	if err != nil {
		slog.Error("failed to create app", "error", err)
		return
	}
	ap.Run()
}
