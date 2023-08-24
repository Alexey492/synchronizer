package logic

import (
	"context"
	"log"
	"path/filepath"
	"sync"
	"time"
)

func Monitor(ctx context.Context, logger *log.Logger, srcDir, dstDir string) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			srcFiles, err := ListFiles(srcDir)
			if err != nil {
				logger.Printf("Failed to list source files: %v", err)
				continue
			}

			dstFiles, err := ListFiles(dstDir)
			if err != nil {
				logger.Printf("Failed to list destination files: %v", err)
				continue
			}

			createFiles := make([]FileInfo, 0)
			updateFiles := make([]FileInfo, 0)
			deleteFiles := make([]FileInfo, 0)

			for _, srcFile := range srcFiles {
				found := false
				for _, dstFile := range dstFiles {
					if srcFile.Path == dstFile.Path {
						found = true
						if srcFile.ModTime.After(dstFile.ModTime) || srcFile.Size != dstFile.Size {
							updateFiles = append(updateFiles, srcFile)
						}
						break
					}
				}
				if !found {
					createFiles = append(createFiles, srcFile)
				}
			}

			for _, dstFile := range dstFiles {
				found := false
				for _, srcFile := range srcFiles {
					if srcFile.Path == dstFile.Path {
						found = true
						break
					}
				}
				if !found {
					deleteFiles = append(deleteFiles, dstFile)
				}
			}

			if len(createFiles) > 0 || len(updateFiles) > 0 || len(deleteFiles) > 0 {
				logger.Printf("Changes detected: create=%d, update=%d, delete=%d", len(createFiles), len(updateFiles), len(deleteFiles))
			}

			var wg sync.WaitGroup

			for _, file := range deleteFiles {
				wg.Add(1)
				go DeleteFile(logger, filepath.Join(dstDir, file.Path), &wg)
			}

			for _, file := range createFiles {
				wg.Add(1)
				go CreateOrUpdateFile(logger, file, srcDir, dstDir, &wg)
			}

			for _, file := range updateFiles {
				wg.Add(1)
				go CreateOrUpdateFile(logger, file, srcDir, dstDir, &wg)
			}

			wg.Wait()

			logger.Printf("Syncing done")
		}
	}
}
