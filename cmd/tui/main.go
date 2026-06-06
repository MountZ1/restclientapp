package main

import (
	"os"
	"path/filepath"
	services "restclient/internal/services/fs"
	"restclient/internal/services/logger"
	"restclient/internal/tui"
)

func main() {
	logPath := services.LogPath()
	logger.Init(filepath.Join(logPath, "app.log"))

	err := tui.Run()
	if err != nil {
		logger.Error("Failed to run the app: %+v", err)
		os.Exit(1)
	}
}
