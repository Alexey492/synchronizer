package copyFiles_test

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"testing"

	logic "github.com/Alexey492/synchronizer/logic"
)

func TestCreateOrUpdateFileDir(t *testing.T) {
	testDir := "./testDir"

	// Удаляем тестовую директорию, чтобы не было проблем с её созданием
	os.RemoveAll(testDir)

	logger := log.New(io.Discard, "", log.LstdFlags)
	var wg sync.WaitGroup

	wg.Add(1)

	file := logic.FileInfo{Path: "testDir", IsDir: true}
	logic.CreateOrUpdateFile(logger, file, ".", testDir, &wg)

	wg.Wait()

	// Проверяем, что директория успешно создалась
	if _, err := os.Stat(filepath.Join(testDir, "testDir")); os.IsNotExist(err) {
		t.Error("Failed to create directory")
	}

	// Удаляем тестовую директорию
	os.RemoveAll(testDir)
}

func TestCreateOrUpdateFileCopy(t *testing.T) {
	testDir := "./testDir"

	// Удаляем тестовую директорию, чтобы не было проблем с её созданием
	os.RemoveAll(testDir)

	logger := log.New(io.Discard, "", log.LstdFlags)
	var wg sync.WaitGroup

	wg.Add(1)

	srcFile, err := os.Create(filepath.Join(".", "test.txt"))
	if err != nil {
		t.Error(err)
	}
	defer srcFile.Close()
	srcFile.WriteString("test file")

	file := logic.FileInfo{Path: "test.txt", IsDir: false}

	// Создаем тестовую директорию перед копированием файлов
	if err := os.Mkdir(testDir, 0744); err != nil {
		t.Fatal(err)
	}

	logic.CreateOrUpdateFile(logger, file, ".", testDir, &wg)

	wg.Wait()

	// Проверяем, что файл успешно скопировался
	dstFile, err := os.Open(filepath.Join(testDir, "test.txt"))
	if err != nil {
		t.Error(err)
	}
	defer dstFile.Close()
	contents, err := io.ReadAll(dstFile)
	if err != nil {
		t.Error(err)
	}

	if string(contents) != "test file" {
		t.Error("Failed to copy file")
	}

	// Удаляем тестовую директорию
	os.RemoveAll(testDir)
	os.Remove(filepath.Join(".", "test.txt"))
}

func TestCreateOrUpdateFileSrcChk(t *testing.T) {
	logger := log.New(io.Discard, "", log.LstdFlags)
	var wg sync.WaitGroup

	// Создаем тестовый файл
	srcFile, err := os.Create(filepath.Join(".", "test.txt"))
	if err != nil {
		t.Error(err)
	}
	defer srcFile.Close()
	srcFile.WriteString("test file")

	// Проверка сущестования файла
	file := logic.FileInfo{Path: "test.txt", IsDir: false}
	wg.Add(1)
	go func() {
		exist, err := CheckSourceFile(file, ".", &wg)
		if err != nil {
			t.Error(err)
		}
		if !exist {
			t.Error("Source file doesn't exist")
		}
	}()

	// Удаление тестового файла
	file = logic.FileInfo{Path: "test.txt", IsDir: false}
	wg.Add(1)
	logic.CreateOrUpdateFile(logger, file, ".", ".", &wg)

	wg.Wait()

	// Закрываем и удаляем тестовый файл
	srcFilePath := filepath.Join(".", "test.txt")
	err = srcFile.Close()
	if err != nil {
		t.Error("Error while closing srcFile", err)
	}

	err = os.Remove(srcFilePath)
	if err != nil {
		t.Error("Failed to delete file", err)
	}
}

func CheckSourceFile(file logic.FileInfo, basePath string, wg *sync.WaitGroup) (bool, error) {
	defer wg.Done()

	fileAbsPath := filepath.Join(basePath, file.Path)

	if _, err := os.Stat(fileAbsPath); os.IsNotExist(err) {
		// Файл не найден
		return false, nil
	} else if err != nil {
		// Произошла другая ошибка
		return false, err
	}

	// Файл найден
	return true, nil
}
