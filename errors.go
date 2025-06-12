package main

import (
	"errors"
	"fmt"
)

// Error types for better error classification and handling
type ErrorType string

const (
	// FileError represents file system related errors
	FileError ErrorType = "FILE_ERROR"

	// ConfigError represents configuration related errors
	ConfigError ErrorType = "CONFIG_ERROR"

	// ValidationError represents validation errors
	ValidationError ErrorType = "VALIDATION_ERROR"

	// RenderError represents rendering related errors
	RenderError ErrorType = "RENDER_ERROR"

	// FontError represents font loading errors
	FontError ErrorType = "FONT_ERROR"

	// ImageError represents image processing errors
	ImageError ErrorType = "IMAGE_ERROR"

	// TemplateError represents template processing errors
	TemplateError ErrorType = "TEMPLATE_ERROR"
)

// AppError represents a structured application error
type AppError struct {
	Type    ErrorType
	Message string
	Cause   error
	Context map[string]interface{}
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap implements the error unwrapping interface
func (e *AppError) Unwrap() error {
	return e.Cause
}

// WithContext adds context information to the error
func (e *AppError) WithContext(key string, value interface{}) *AppError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// NewAppError creates a new structured application error
func NewAppError(errorType ErrorType, message string, cause error) *AppError {
	return &AppError{
		Type:    errorType,
		Message: message,
		Cause:   cause,
	}
}

// Common error factory functions for consistent error creation

// NewFileError creates a file-related error
func NewFileError(operation string, path string, cause error) *AppError {
	return NewAppError(FileError, fmt.Sprintf("failed to %s file %s", operation, path), cause).
		WithContext("operation", operation).
		WithContext("path", path)
}

// NewConfigError creates a configuration-related error
func NewConfigError(message string, cause error) *AppError {
	return NewAppError(ConfigError, message, cause)
}

// NewValidationError creates a validation error
func NewValidationError(message string) *AppError {
	return NewAppError(ValidationError, message, nil)
}

// NewRenderError creates a rendering-related error
func NewRenderError(component string, cause error) *AppError {
	return NewAppError(RenderError, fmt.Sprintf("failed to render %s", component), cause).
		WithContext("component", component)
}

// NewFontError creates a font-related error
func NewFontError(operation string, fontPath string, cause error) *AppError {
	return NewAppError(FontError, fmt.Sprintf("failed to %s font %s", operation, fontPath), cause).
		WithContext("operation", operation).
		WithContext("fontPath", fontPath)
}

// NewImageError creates an image processing error
func NewImageError(operation string, imagePath string, cause error) *AppError {
	return NewAppError(ImageError, fmt.Sprintf("failed to %s image %s", operation, imagePath), cause).
		WithContext("operation", operation).
		WithContext("imagePath", imagePath)
}

// NewTemplateError creates a template processing error
func NewTemplateError(operation string, cause error) *AppError {
	return NewAppError(TemplateError, fmt.Sprintf("failed to %s template", operation), cause).
		WithContext("operation", operation)
}

// Error checking helper functions

// IsFileError checks if an error is a file-related error
func IsFileError(err error) bool {
	if err == nil {
		return false
	}
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Type == FileError
	}
	return false
}

// IsConfigError checks if an error is a configuration error
func IsConfigError(err error) bool {
	if err == nil {
		return false
	}
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Type == ConfigError
	}
	return false
}

// IsValidationError checks if an error is a validation error
func IsValidationError(err error) bool {
	if err == nil {
		return false
	}
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Type == ValidationError
	}
	return false
}

// IsRenderError checks if an error is a rendering error
func IsRenderError(err error) bool {
	if err == nil {
		return false
	}
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Type == RenderError
	}
	return false
}
