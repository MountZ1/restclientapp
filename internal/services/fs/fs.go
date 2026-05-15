package services

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
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

func RenameRequestOrCollection(path string, name string) error {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	nameOnly := strings.TrimSuffix(base, ext)
	prefix, after, found := strings.Cut(nameOnly, "-")

	var newName string
	if found && after != "" && isHTTPMethod(prefix) {
		newName = prefix + "-" + name + ext
	} else {
		newName = name
	}

	newPath := filepath.Join(filepath.Dir(path), newName)
	return os.Rename(path, newPath)
}

func isHTTPMethod(s string) bool {
	switch s {
	case "GET", "POST", "PUT", "DELETE", "PATCH":
		return true
	}
	return false
}

func DestroyRequestOrCollection(path string) error {
	return os.RemoveAll(path)
}

func DuplicateRequestOrCollection(path string, isDir bool) error {
	if isDir {
		return duplicateCollection(path)
	}
	dir := filepath.Dir(path)
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	nameOnly := strings.TrimSuffix(base, ext)
	newTarget := generateUniqueName(dir, nameOnly, ext)
	return duplicateRequest(path, newTarget)
}

func duplicateRequest(path string, target string) error {
	in, err := os.Open(path)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(target)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func duplicateCollection(src string) error {
	dir := filepath.Dir(src)
	base := filepath.Base(src)
	dst := generateUniqueName(dir, base, "")

	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		return duplicateRequest(path, target)
	})
}

func generateUniqueName(dir, nameOnly, ext string) string {
	indexRegex := regexp.MustCompile(`\(\d+\)$`)

	base := indexRegex.ReplaceAllString(nameOnly, "")

	for i := 1; ; i++ {
		candidate := filepath.Join(dir, fmt.Sprintf("%s(%d)%s", base, i, ext))
		if _, err := os.Stat(candidate); os.IsNotExist(err) {
			return candidate
		}
	}
}
