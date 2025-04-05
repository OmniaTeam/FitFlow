package logger

import (
	"io"
	"log"
	"log/slog"
	"os"
)

func GetSlogFileConsoleJsonHandler() (closeLogFile func() error) {

	err := os.MkdirAll("logs", os.ModePerm)
	if err != nil {
		log.Fatalf("failed to create log directory: %v", err)
	}

	logFile, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}

	logWriter := io.MultiWriter(os.Stdout, logFile)

	var loggerOptions slog.HandlerOptions

	handler := slog.Handler(slog.NewJSONHandler(logWriter, &loggerOptions))
	handler = NewHandlerMiddleware(handler)
	slog.SetDefault(slog.New(handler))
	return
}
