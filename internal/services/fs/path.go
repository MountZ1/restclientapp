package services

import (
	"os"
	"path/filepath"
)

func Path() string {
	exepath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exepath, _ = filepath.EvalSymlinks(exepath)
	exeDir := filepath.Dir(exepath)
	dataPath := filepath.Join(exeDir, "data")

	if err := os.MkdirAll(dataPath, 0755); err != nil {
		return ""
	}
	return dataPath
}

func LogPath() string {
	exepath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exepath, _ = filepath.EvalSymlinks(exepath)
	exeDir := filepath.Dir(exepath)

	logPath := filepath.Join(exeDir, "log")
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		os.MkdirAll(logPath, 0755)
	}
	return logPath
}
