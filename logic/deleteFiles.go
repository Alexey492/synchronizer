package logic

import (
	"log"
	"os"
	"sync"
)

func DeleteFile(logger *log.Logger, path string, wg *sync.WaitGroup) {
	defer wg.Done()

	err := os.Remove(path)
	if err != nil {
		logger.Printf("Failed to delete file %s: %v", path, err)
		return
	}
	logger.Printf("Deleted file %s", path)
}
