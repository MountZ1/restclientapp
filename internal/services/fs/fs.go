package services

import (
	"os"
	"path/filepath"
	"restclient/internal/collections"
	"strings"
)

type Create struct {
	Location string
	Type     string
	Name     string
}

func Initialize() (string, []os.DirEntry, error) {
	dataPath := Path()

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

func CreateRequestCollection(request Create) error {
	base := Path()
	if request.Location != "" {
		base = request.Location
	}

	switch request.Type {
	case "collection":
		if err := os.Mkdir(filepath.Join(base, request.Name), 0755); err != nil {
			return err
		}
	case "request":
		f, err := os.Create(filepath.Join(base, "GET-"+request.Name+".json"))
		if err != nil {
			return err
		}
		f.Close()
	}
	return nil
}
