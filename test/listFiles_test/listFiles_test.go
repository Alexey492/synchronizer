package listFiles_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	logic "github.com/Alexey492/synchronizer/logic"
)

func TestListFiles(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Errorf("Unable to get current directory: %s", err.Error())
	}
	files, err := logic.ListFiles(dir)
	if err != nil {
		t.Errorf("Error while getting files: %s", err.Error())
	}
	if len(files) == 0 {
		t.Errorf("No files found in directory: %s", dir)
	}
}

func TestListFilesNoDir(t *testing.T) {
	var dir string = "fake_directory"
	_, err := logic.ListFiles(dir)
	if err == nil {
		t.Errorf("Expected error while getting files from non-existent directory")
	}
}

func TestListFilesNonExistentDir(t *testing.T) {
	_, err := logic.ListFiles("./non_existent_dir")
	if err == nil {
		t.Error("expected an error but got none")
	}
}

func TestListFilesInvalidRelPath(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Errorf("Unable to get current directory: %s", err.Error())
	}
	var fakePath string = "/fake/path/to/file"
	relPath, err := filepath.Rel(dir, fakePath)
	if err == nil {
		t.Errorf("Expected error while getting relative path for non-existent file")
	}

	fakeInfo := fakeFileInfo{
		name:    "fakeFileInfo",
		size:    0,
		mode:    0,
		modTime: time.Time{},
		isDir:   false,
	}
	file := logic.FileInfo{
		Path:    relPath,
		Size:    fakeInfo.Size(),
		ModTime: fakeInfo.ModTime(),
		IsDir:   fakeInfo.IsDir(),
	}

	if !(file.Path == "" && file.Size == 0 && file.ModTime.IsZero() && !file.IsDir) {
		t.Errorf("Expected empty file info for non-existent file")
	}
}

//Моки для тестов

type fakeFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
}

func (f fakeFileInfo) Name() string {
	return f.name
}

func (f fakeFileInfo) Size() int64 {
	return f.size
}

func (f fakeFileInfo) Mode() os.FileMode {
	return f.mode
}

func (f fakeFileInfo) ModTime() time.Time {
	return f.modTime
}

func (f fakeFileInfo) IsDir() bool {
	return f.isDir
}

func (f fakeFileInfo) Sys() interface{} {
	return nil
}
