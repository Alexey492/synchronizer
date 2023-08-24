package logger

import "fmt"

type LogWriter struct {
}

func (writer LogWriter) Write(bytes []byte) (n int, err error) {
	// здесь выполняется логирование
	fmt.Println(string(bytes))
	return len(bytes), nil
}
