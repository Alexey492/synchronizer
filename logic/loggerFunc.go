package logic

import (
	"io"
	"log"
	"os"
)

// MultiWriter вызывает io.MultiWriter с переданными в него писателями.
func MultiWriter(writers ...io.Writer) io.Writer {
	return io.MultiWriter(writers...)
}

// GetLogger возвращает экземпляр логгера с настройками.
// Логи записываются и в stdout, и в файл по пути "log.txt"
func GetLogger() *log.Logger {
	logFile, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	return log.New(MultiWriter(os.Stdout, logFile), "", log.Ldate|log.Ltime)
}
