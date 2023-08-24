package monitor_test

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	logic "github.com/Alexey492/synchronizer/logic"

	"github.com/stretchr/testify/assert"
)

func TestMonitor(t *testing.T) {
	srcPath, _ := os.MkdirTemp("", "source")
	defer os.RemoveAll(srcPath)

	dstPath, _ := os.MkdirTemp("", "destination")
	defer os.RemoveAll(dstPath)

	file1Path := filepath.Join(srcPath, "file1.txt")
	file2Path := filepath.Join(srcPath, "file2.txt")
	file3Path := filepath.Join(srcPath, "file3.txt")
	file4Path := filepath.Join(srcPath, "file4.txt")

	createFile(file1Path, "content1")
	createFile(file2Path, "content2")
	createFile(file3Path, "content3")
	createFile(file4Path, "content4")

	logger := log.New(io.Discard, "", log.LstdFlags)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		time.Sleep(10 * time.Second)
		cancel()
	}()

	logic.Monitor(ctx, logger, srcPath, dstPath)

	assert.FileExists(t, filepath.Join(dstPath, "file1.txt"))
	assert.FileExists(t, filepath.Join(dstPath, "file2.txt"))
	assert.FileExists(t, filepath.Join(dstPath, "file3.txt"))
	assert.FileExists(t, filepath.Join(dstPath, "file4.txt"))

	content, _ := os.ReadFile(filepath.Join(dstPath, "file1.txt"))
	assert.Equal(t, "content1", string(content))

	content, _ = os.ReadFile(filepath.Join(dstPath, "file2.txt"))
	assert.Equal(t, "content2", string(content))

	content, _ = os.ReadFile(filepath.Join(dstPath, "file3.txt"))
	assert.Equal(t, "content3", string(content))

	content, _ = os.ReadFile(filepath.Join(dstPath, "file4.txt"))
	assert.Equal(t, "content4", string(content))
}

func createFile(filePath string, content string) {
	f, _ := os.Create(filePath)
	defer f.Close()

	f.WriteString(content)
}
