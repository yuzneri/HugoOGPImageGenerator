package main

import (
	"errors"
	"testing"
)

func TestAppError_Error(t *testing.T) {
	tests := []struct {
		name     string
		appError *AppError
		expected string
	}{
		{
			name: "Error with cause",
			appError: &AppError{
				Type:    FileError,
				Message: "failed to read file",
				Cause:   errors.New("no such file"),
			},
			expected: "FILE_ERROR: failed to read file: no such file",
		},
		{
			name: "Error without cause",
			appError: &AppError{
				Type:    ValidationError,
				Message: "invalid input",
				Cause:   nil,
			},
			expected: "VALIDATION_ERROR: invalid input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.appError.Error()
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestAppError_Unwrap(t *testing.T) {
	originalErr := errors.New("original error")
	appErr := &AppError{
		Type:    FileError,
		Message: "wrapped error",
		Cause:   originalErr,
	}

	unwrapped := appErr.Unwrap()
	if unwrapped != originalErr {
		t.Errorf("Expected unwrapped error to be %v, got %v", originalErr, unwrapped)
	}

	// Test with no cause
	appErrNoCause := &AppError{
		Type:    ValidationError,
		Message: "no cause",
		Cause:   nil,
	}

	unwrappedNil := appErrNoCause.Unwrap()
	if unwrappedNil != nil {
		t.Errorf("Expected unwrapped error to be nil, got %v", unwrappedNil)
	}
}

func TestAppError_WithContext(t *testing.T) {
	appErr := &AppError{
		Type:    FileError,
		Message: "test error",
	}

	result := appErr.WithContext("operation", "read").WithContext("path", "/test/path")

	if result.Context["operation"] != "read" {
		t.Errorf("Expected operation context to be 'read', got %v", result.Context["operation"])
	}

	if result.Context["path"] != "/test/path" {
		t.Errorf("Expected path context to be '/test/path', got %v", result.Context["path"])
	}
}

func TestNewAppError(t *testing.T) {
	cause := errors.New("underlying error")
	appErr := NewAppError(ConfigError, "test message", cause)

	if appErr.Type != ConfigError {
		t.Errorf("Expected type %v, got %v", ConfigError, appErr.Type)
	}

	if appErr.Message != "test message" {
		t.Errorf("Expected message 'test message', got %q", appErr.Message)
	}

	if appErr.Cause != cause {
		t.Errorf("Expected cause %v, got %v", cause, appErr.Cause)
	}
}

func TestNewFileError(t *testing.T) {
	cause := errors.New("permission denied")
	fileErr := NewFileError("read", "/test/file.txt", cause)

	if fileErr.Type != FileError {
		t.Errorf("Expected type %v, got %v", FileError, fileErr.Type)
	}

	expectedMessage := "failed to read file /test/file.txt"
	if fileErr.Message != expectedMessage {
		t.Errorf("Expected message %q, got %q", expectedMessage, fileErr.Message)
	}

	if fileErr.Cause != cause {
		t.Errorf("Expected cause %v, got %v", cause, fileErr.Cause)
	}

	// Check context
	if fileErr.Context["operation"] != "read" {
		t.Errorf("Expected operation context to be 'read', got %v", fileErr.Context["operation"])
	}

	if fileErr.Context["path"] != "/test/file.txt" {
		t.Errorf("Expected path context to be '/test/file.txt', got %v", fileErr.Context["path"])
	}
}

func TestNewConfigError(t *testing.T) {
	cause := errors.New("invalid YAML")
	configErr := NewConfigError("failed to parse config", cause)

	if configErr.Type != ConfigError {
		t.Errorf("Expected type %v, got %v", ConfigError, configErr.Type)
	}

	if configErr.Message != "failed to parse config" {
		t.Errorf("Expected message 'failed to parse config', got %q", configErr.Message)
	}

	if configErr.Cause != cause {
		t.Errorf("Expected cause %v, got %v", cause, configErr.Cause)
	}
}

func TestNewValidationError(t *testing.T) {
	validationErr := NewValidationError("input is required")

	if validationErr.Type != ValidationError {
		t.Errorf("Expected type %v, got %v", ValidationError, validationErr.Type)
	}

	if validationErr.Message != "input is required" {
		t.Errorf("Expected message 'input is required', got %q", validationErr.Message)
	}

	if validationErr.Cause != nil {
		t.Errorf("Expected cause to be nil, got %v", validationErr.Cause)
	}
}

func TestNewRenderError(t *testing.T) {
	cause := errors.New("font not found")
	renderErr := NewRenderError("title", cause)

	if renderErr.Type != RenderError {
		t.Errorf("Expected type %v, got %v", RenderError, renderErr.Type)
	}

	expectedMessage := "failed to render title"
	if renderErr.Message != expectedMessage {
		t.Errorf("Expected message %q, got %q", expectedMessage, renderErr.Message)
	}

	if renderErr.Cause != cause {
		t.Errorf("Expected cause %v, got %v", cause, renderErr.Cause)
	}

	if renderErr.Context["component"] != "title" {
		t.Errorf("Expected component context to be 'title', got %v", renderErr.Context["component"])
	}
}

func TestNewFontError(t *testing.T) {
	cause := errors.New("invalid font format")
	fontErr := NewFontError("parse", "test.ttf", cause)

	if fontErr.Type != FontError {
		t.Errorf("Expected type %v, got %v", FontError, fontErr.Type)
	}

	expectedMessage := "failed to parse font test.ttf"
	if fontErr.Message != expectedMessage {
		t.Errorf("Expected message %q, got %q", expectedMessage, fontErr.Message)
	}

	if fontErr.Cause != cause {
		t.Errorf("Expected cause %v, got %v", cause, fontErr.Cause)
	}

	if fontErr.Context["operation"] != "parse" {
		t.Errorf("Expected operation context to be 'parse', got %v", fontErr.Context["operation"])
	}

	if fontErr.Context["fontPath"] != "test.ttf" {
		t.Errorf("Expected fontPath context to be 'test.ttf', got %v", fontErr.Context["fontPath"])
	}
}

func TestNewImageError(t *testing.T) {
	cause := errors.New("corrupted image")
	imageErr := NewImageError("decode", "test.jpg", cause)

	if imageErr.Type != ImageError {
		t.Errorf("Expected type %v, got %v", ImageError, imageErr.Type)
	}

	expectedMessage := "failed to decode image test.jpg"
	if imageErr.Message != expectedMessage {
		t.Errorf("Expected message %q, got %q", expectedMessage, imageErr.Message)
	}

	if imageErr.Cause != cause {
		t.Errorf("Expected cause %v, got %v", cause, imageErr.Cause)
	}

	if imageErr.Context["operation"] != "decode" {
		t.Errorf("Expected operation context to be 'decode', got %v", imageErr.Context["operation"])
	}

	if imageErr.Context["imagePath"] != "test.jpg" {
		t.Errorf("Expected imagePath context to be 'test.jpg', got %v", imageErr.Context["imagePath"])
	}
}

func TestNewTemplateError(t *testing.T) {
	cause := errors.New("syntax error")
	templateErr := NewTemplateError("parse", cause)

	if templateErr.Type != TemplateError {
		t.Errorf("Expected type %v, got %v", TemplateError, templateErr.Type)
	}

	expectedMessage := "failed to parse template"
	if templateErr.Message != expectedMessage {
		t.Errorf("Expected message %q, got %q", expectedMessage, templateErr.Message)
	}

	if templateErr.Cause != cause {
		t.Errorf("Expected cause %v, got %v", cause, templateErr.Cause)
	}

	if templateErr.Context["operation"] != "parse" {
		t.Errorf("Expected operation context to be 'parse', got %v", templateErr.Context["operation"])
	}
}

func TestIsFileError(t *testing.T) {
	fileErr := NewFileError("read", "test.txt", errors.New("test"))
	if !IsFileError(fileErr) {
		t.Error("Expected IsFileError to return true for FileError")
	}

	configErr := NewConfigError("test", nil)
	if IsFileError(configErr) {
		t.Error("Expected IsFileError to return false for ConfigError")
	}

	if IsFileError(nil) {
		t.Error("Expected IsFileError to return false for nil error")
	}

	standardErr := errors.New("standard error")
	if IsFileError(standardErr) {
		t.Error("Expected IsFileError to return false for standard error")
	}
}

func TestIsConfigError(t *testing.T) {
	configErr := NewConfigError("test", nil)
	if !IsConfigError(configErr) {
		t.Error("Expected IsConfigError to return true for ConfigError")
	}

	fileErr := NewFileError("read", "test.txt", errors.New("test"))
	if IsConfigError(fileErr) {
		t.Error("Expected IsConfigError to return false for FileError")
	}

	if IsConfigError(nil) {
		t.Error("Expected IsConfigError to return false for nil error")
	}
}

func TestIsValidationError(t *testing.T) {
	validationErr := NewValidationError("test")
	if !IsValidationError(validationErr) {
		t.Error("Expected IsValidationError to return true for ValidationError")
	}

	fileErr := NewFileError("read", "test.txt", errors.New("test"))
	if IsValidationError(fileErr) {
		t.Error("Expected IsValidationError to return false for FileError")
	}

	if IsValidationError(nil) {
		t.Error("Expected IsValidationError to return false for nil error")
	}
}

func TestIsRenderError(t *testing.T) {
	renderErr := NewRenderError("test", errors.New("test"))
	if !IsRenderError(renderErr) {
		t.Error("Expected IsRenderError to return true for RenderError")
	}

	fileErr := NewFileError("read", "test.txt", errors.New("test"))
	if IsRenderError(fileErr) {
		t.Error("Expected IsRenderError to return false for FileError")
	}

	if IsRenderError(nil) {
		t.Error("Expected IsRenderError to return false for nil error")
	}
}

func TestErrorChaining(t *testing.T) {
	originalErr := errors.New("original error")
	wrappedErr := NewFileError("read", "test.txt", originalErr)

	// Test errors.Is
	if !errors.Is(wrappedErr, originalErr) {
		t.Error("Expected errors.Is to find the original error in the chain")
	}

	// Test errors.As
	var appErr *AppError
	if !errors.As(wrappedErr, &appErr) {
		t.Error("Expected errors.As to find AppError in the chain")
	}

	if appErr.Type != FileError {
		t.Errorf("Expected type %v, got %v", FileError, appErr.Type)
	}
}
