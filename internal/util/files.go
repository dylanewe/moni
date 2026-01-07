package util

import (
	"os"
	"path/filepath"
	"slices"
)

func ReadFilesFromFolder(folderPath string, extensions []string) ([]string, error) {
	var contents []string

	err := filepath.WalkDir(folderPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if len(extensions) > 0 {
			ext := filepath.Ext(path)
			if !slices.Contains(extensions, ext) {
				return nil
			}
		}

		contents = append(contents, string(d.Name()))
		return nil
	})

	if err != nil {
		return nil, err
	}

	return contents, nil
}
