package main

import (
	"bytes"
	"os"
	"testing"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger()
	if logger == nil {
		t.Error("NewLogger should return a non-nil logger")
	}
}

func TestLogger_Warning(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	logger := NewLogger()
	logger.Warning("This is a warning: %s", "test message")

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	expected := "Warning: This is a warning: test message\n"
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}
}

func TestLogger_Error(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	logger := NewLogger()
	logger.Error("This is an error: %d", 42)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	expected := "Error: This is an error: 42\n"
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}
}

func TestLogger_Info(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	logger := NewLogger()
	logger.Info("This is info: %v", true)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	expected := "Info: This is info: true\n"
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}
}

func TestLogger_Warning_NoArgs(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	logger := NewLogger()
	logger.Warning("Simple warning")

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	expected := "Warning: Simple warning\n"
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}
}

func TestLogger_Error_NoArgs(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	logger := NewLogger()
	logger.Error("Simple error")

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	expected := "Error: Simple error\n"
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}
}

func TestLogger_Info_NoArgs(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	logger := NewLogger()
	logger.Info("Simple info")

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	expected := "Info: Simple info\n"
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}
}

func TestLogger_MultipleArgs(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	logger := NewLogger()
	logger.Warning("Multiple args: %s, %d, %v", "text", 123, false)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	expected := "Warning: Multiple args: text, 123, false\n"
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}
}

func TestDefaultLogger(t *testing.T) {
	if DefaultLogger == nil {
		t.Error("DefaultLogger should not be nil")
	}

	// Test that DefaultLogger is a Logger instance
	if _, ok := interface{}(DefaultLogger).(*Logger); !ok {
		t.Error("DefaultLogger should be a *Logger instance")
	}
}

// Helper function to capture stdout for testing
func captureOutput(f func()) string {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}

func TestLogger_UsingCaptureHelper(t *testing.T) {
	logger := NewLogger()

	output := captureOutput(func() {
		logger.Info("Test with helper")
	})

	expected := "Info: Test with helper\n"
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}
}

// Note: Logger.Fatal() is not tested here because it calls log.Fatalf()
// which terminates the program. In a real scenario, you might want to:
// 1. Refactor Fatal to use dependency injection for testability
// 2. Use a testing framework that can handle program termination
// 3. Test Fatal in integration tests rather than unit tests
