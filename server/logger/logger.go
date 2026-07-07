package logger

import (
	"log"
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"
)

func InitLogger(LogPath string) logr.Logger {
	logFile, err := os.OpenFile(LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	stdLogger := log.New(logFile, "", log.LstdFlags|log.Lshortfile)
	baseLogger := stdr.New(stdLogger)
	logger := baseLogger.V(0)
	return logger
}
