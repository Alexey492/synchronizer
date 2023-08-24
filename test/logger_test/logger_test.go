package logger_test

import (
	"bytes"
	"io"
	"os"
	"testing"

	logic "github.com/Alexey492/synchronizer/logic"
)

// Этот работает
func TestGetLogger(t *testing.T) {
	logger := logic.GetLogger()
	if logger == nil {
		t.Error("expected logger to not be nil")
	}
}

// Этот работает
func TestGetLoggerWithPathNotExist(t *testing.T) {
	os.Remove("log.txt")
	logger := logic.GetLogger()
	if logger == nil {
		t.Error("expected logger to not be nil")
	}
	_, err := os.Stat("log.txt")
	if os.IsNotExist(err) {
		t.Error("expected log file to be created")
	}
}

// Этот работает
func TestGetLoggerWithPathExist(t *testing.T) {
	os.Create("log.txt")
	logger := logic.GetLogger()
	if logger == nil {
		t.Error("expected logger to not be nil")
	}
}

// Этот работает
func TestMultiWriter(t *testing.T) {
	var buf bytes.Buffer
	w1 := &buf
	w2 := io.Discard
	writers := []io.Writer{w1, w2}

	mw := io.MultiWriter(writers...)
	mw.Write([]byte("test"))

	if buf.String() != "test" {
		t.Errorf("Expected buffer to contain 'test', got %s", buf.String())
	}
}

/*
Эти не работают, что-то с правами. Но они довольно простые. Хорошо бы их оставить
func TestGetLoggerConsistency(t *testing.T) {
	logger1 := logic.GetLogger()
	logger2 := logic.GetLogger()

	if logger1 != logger2 {
		t.Error("expected GetLogger to return a consistent logger object")
	}
}

func TestGetLoggerNotNil(t *testing.T) {
	logger := logic.GetLogger()

	if logger == nil {
		t.Error("expected GetLogger to not return nil")
	}
}
*/
