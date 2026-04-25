package services

import (
	"os"
	"path/filepath"
	"restclient/internal/collections"
	"strings"
)

func Initialize() (string, []os.DirEntry, error) {
	exepath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	exeDir := filepath.Dir(exepath)
	dataPath := filepath.Join(exeDir, "data")

	// files, err := os.ReadDir(dataPath)

	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		err = os.Mkdir(dataPath, 0755)
		if err != nil {
			return "", nil, err
		}
	}

	files, err := os.ReadDir(dataPath)
	if err != nil {
		return "", nil, err
	}

	return dataPath, files, nil
}

func StartColection() string {
	var b strings.Builder
	_, files, err := Initialize()
	if err != nil {
		// panic(err)
		return "something went wrong, please contact customer service"
	}

	for _, f := range files {
		if f.IsDir() {
			b.WriteString(collections.IconFolder + " " + f.Name() + "\n")
		} else if strings.ToLower(filepath.Ext(f.Name())) == ".json" {
			b.WriteString(collections.IconFile + " " + f.Name() + "\n")
		}
	}

	return b.String()
}

func GetCollectionItems() []os.DirEntry {
	_, files, err := Initialize()
	if err != nil {
		return []os.DirEntry{}
	}

	var filtered []os.DirEntry
	for _, f := range files {
		if f.IsDir() || strings.ToLower(filepath.Ext(f.Name())) == ".json" {
			filtered = append(filtered, f)
		}
	}
	return filtered
}

func read() {
}
