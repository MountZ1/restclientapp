package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"

	appgui "restclient/internal/app-gui"
	services "restclient/internal/services/fs"
	"restclient/internal/services/logger"
)

func main() {
	logPath := services.LogPath()
	logger.Init(filepath.Join(logPath, "app.log"))

	if err := run(); err != nil {
		logger.Error("fatal: %v", err)
		fmt.Fprintln(os.Stderr, "restclient: fatal error, see app.log for details:", err)
		os.Exit(1)
	}
}

func run() (err error) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("panic recovered: %v\n%s", r, debug.Stack())
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	appgui.Run()
	return nil
}
