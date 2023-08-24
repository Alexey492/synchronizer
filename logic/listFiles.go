package logic

import (
	"os"
	"path/filepath"
)

func ListFiles(dir string) ([]FileInfo, error) {
	var files []FileInfo

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			relPath, err := filepath.Rel(dir, path)
			if err != nil {
				return err
			}

			files = append(files, FileInfo{
				Path:    relPath,
				Size:    info.Size(),
				ModTime: info.ModTime(),
				IsDir:   false, // файл
			})
		} else {
			relPath, err := filepath.Rel(dir, path)
			if err != nil {
				return err
			}

			files = append(files, FileInfo{
				Path:    relPath,
				Size:    info.Size(),
				ModTime: info.ModTime(),
				IsDir:   true, // директория
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}
