package services

import (
	"os"
	"path/filepath"
)

type Create struct {
	Location string
	Type     string
	Name     string
}

func CreateRequestCollection(request Create) error {
	path := Path()

	switch request.Type {
	case "collection":
		newPath := filepath.Join(path, request.Location, request.Name)
		err := os.Mkdir(newPath, 0755)
		if err != nil {
			return err
		}
	case "request":
		newPath := filepath.Join(path, request.Location, "GET-"+request.Name+".json")
		f, err := os.Create(newPath)
		if err != nil {
			return err
		}

		f.Close()
	}

	return nil
}
