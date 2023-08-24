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

func BenchmarkCreateOrUpdateFile(b *testing.B) {
	// Создаем временную папку для бенчмарка
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Создаем файлы
	srcFile, err := os.CreateTemp(tempDir, "source")
	if err != nil {
		b.Fatal(err)
	}
	defer srcFile.Close()

	_, err = os.CreateTemp(tempDir, "destination")
	if err != nil {
		b.Fatal(err)
	}

	logger := log.New(io.Discard, "", 0)
	var wg sync.WaitGroup

	for i := 0; i < b.N; i++ {
		// Увеличиваем счетчик горутин до запуска горутины
		wg.Add(1)
		go logic.CreateOrUpdateFile(
			logger,
			logic.FileInfo{
				Path:  filepath.Base(srcFile.Name()),
				IsDir: false,
			},
			filepath.Dir(srcFile.Name()),
			tempDir,
			&wg,
		)
	}
	// Ждем завершения всех горутин
	wg.Wait()

	// Проверяем, что файл скопировался
	_, err = os.Stat(filepath.Join(tempDir, filepath.Base(srcFile.Name())))
	if err != nil {
		b.Fatal(err)
	}
}
