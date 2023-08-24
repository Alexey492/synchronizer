package logic

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

func CreateOrUpdateFile(logger *log.Logger, file FileInfo, srcDir, dstDir string, wg *sync.WaitGroup) {
	defer wg.Done()

	srcPath := filepath.Join(srcDir, file.Path)
	dstPath := filepath.Join(dstDir, file.Path)

	if file.IsDir {
		// Если путь является директорией, то нужно его создать в директории назначения
		err := os.MkdirAll(dstPath, 0755)
		if err != nil {
			logger.Printf("Failed to create destination directory %s: %v", dstPath, err)
			return
		}
		logger.Printf("Created directory %s", dstPath)
		return
	}

	srcFile, err := os.Open(srcPath)
	if err != nil {
		logger.Printf("Failed to open source file %s: %v", srcPath, err)
		return
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dstPath)
	if err != nil {
		logger.Printf("Failed to create destination file %s: %v", dstPath, err)
		return
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		logger.Printf("Failed to copy file %s to %s: %v", srcPath, dstPath, err)
		return
	}

	logger.Printf("Created or updated file %s", dstPath)
}
