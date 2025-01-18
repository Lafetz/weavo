package customlogger

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"
	"testing"
)

func captureOutput(f func()) string {
	r, w, _ := os.Pipe()
	stdout := os.Stdout
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = stdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestNewLoggerDevelopment(t *testing.T) {
	output := captureOutput(func() {
		logger := NewLogger(slog.LevelDebug, "development")
		logger.Debug("debug message")
	})
	fmt.Println(output)

	if !bytes.Contains([]byte(output), []byte("level=DEBUG")) {
		t.Fatalf("expected log message to contain 'level=DEBUG', got %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("source=")) {
		t.Fatalf("expected log message to contain 'source=', got %s", output)
	}
}

func TestNewLoggerProduction(t *testing.T) {
	output := captureOutput(func() {
		logger := NewLogger(slog.LevelInfo, "production")
		logger.Info("info message")
	})

	if !bytes.Contains([]byte(output), []byte(`"level":"INFO"`)) {
		t.Fatalf("expected log message to contain '\"level\":\"INFO\"', got %s", output)
	}
	if !bytes.Contains([]byte(output), []byte(`"source":`)) {
		t.Fatalf("expected log message to contain 'source', got %s", output)
	}
}

func TestNewLoggerDebugLevel(t *testing.T) {
	output := captureOutput(func() {
		logger := NewLogger(slog.LevelDebug, "development")
		logger.Debug("debug message")
	})

	if !bytes.Contains([]byte(output), []byte("level=DEBUG")) {
		t.Fatalf("expected log message to contain 'level=DEBUG', got %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("debug message")) {
		t.Fatalf("expected log message to contain 'debug message', got %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("source")) {
		t.Fatalf("expected log message to contain 'source', got %s", output)
	}
}

func TestNewLoggerInfoLevel(t *testing.T) {
	output := captureOutput(func() {
		logger := NewLogger(slog.LevelInfo, "development")
		logger.Info("info message")
	})

	if !bytes.Contains([]byte(output), []byte("level=INFO")) {
		t.Fatalf("expected log message to contain 'level=INFO', got %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("info message")) {
		t.Fatalf("expected log message to contain 'info message', got %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("source")) {
		t.Fatalf("expected log message to contain 'source', got %s", output)
	}
}

func TestNewLoggerWarnLevel(t *testing.T) {
	output := captureOutput(func() {
		logger := NewLogger(slog.LevelWarn, "development")
		logger.Warn("warn message")
	})

	if !bytes.Contains([]byte(output), []byte("level=WARN")) {
		t.Fatalf("expected log message to contain 'level=WARN', got %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("warn message")) {
		t.Fatalf("expected log message to contain 'warn message', got %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("source")) {
		t.Fatalf("expected log message to contain 'source', got %s", output)
	}
}

func TestNewLoggerErrorLevel(t *testing.T) {
	output := captureOutput(func() {
		logger := NewLogger(slog.LevelError, "development")
		logger.Error("error message")
	})

	if !bytes.Contains([]byte(output), []byte("level=ERROR")) {
		t.Fatalf("expected log message to contain 'level=ERROR', got %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("error message")) {
		t.Fatalf("expected log message to contain 'error message', got %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("source")) {
		t.Fatalf("expected log message to contain 'source', got %s", output)
	}
}
