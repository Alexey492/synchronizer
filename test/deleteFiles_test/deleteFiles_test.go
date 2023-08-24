package deleteFiles_test

import (
	"log"
	"os"
	"sync"
	"testing"

	logic "github.com/Alexey492/synchronizer/logic"
)

func TestDeleteFile(t *testing.T) {
	logger := log.New(os.Stdout, "TestDeleteFile ", log.LstdFlags)

	//Тест на удаление файла
	wg := sync.WaitGroup{}
	wg.Add(1)
	go logic.DeleteFile(logger, "testfile.txt", &wg)
	wg.Wait()
	if _, err := os.Stat("testfile.txt"); !os.IsNotExist(err) {
		t.Errorf("Error: Expected file to be deleted, but still exists")
	}

	// Проверяем удаление несуществующих файлов
	wg.Add(1)
	go logic.DeleteFile(logger, "non-existent-file.txt", &wg)
	wg.Wait()

	// Тест для удаления запертого файла
	file, err := os.OpenFile("lockedfile.txt", os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if _, err := os.Stat("lockedfile.txt"); os.IsNotExist(err) {
		t.Errorf("Error: Expected file to be created and exist, but it doesn't")
	}
	wg.Add(1)
	go logic.DeleteFile(logger, "lockedfile.txt", &wg)
	wg.Wait()
	if _, err := os.Stat("lockedfile.txt"); os.IsNotExist(err) {
		t.Errorf("Error: Expected file to still exist, but it doesn't")
	}
	file.Close()
}
