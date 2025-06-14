package main

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// FileOperations provides common file operation utilities
type FileOperations struct{}

// NewFileOperations creates a new FileOperations instance
func NewFileOperations() *FileOperations {
	return &FileOperations{}
}

// ExistsAndIsFile checks if a file exists and is not a directory
func (fo *FileOperations) ExistsAndIsFile(path string) (bool, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return !info.IsDir(), nil
}

// ReadYAMLFile reads and unmarshals a YAML file into the provided target
func (fo *FileOperations) ReadYAMLFile(path string, target interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return NewFileError("read", path, err)
	}

	err = yaml.Unmarshal(data, target)
	if err != nil {
		return NewConfigError("failed to unmarshal YAML file "+filepath.Base(path), err)
	}

	return nil
}

// WriteYAMLFile marshals data to YAML and writes it to a file
func (fo *FileOperations) WriteYAMLFile(path string, data interface{}) error {
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return NewConfigError("failed to marshal YAML data", err)
	}

	err = os.WriteFile(path, yamlData, 0644)
	if err != nil {
		return NewFileError("write", path, err)
	}

	return nil
}

// PathValidator provides secure path validation utilities
type PathValidator struct{}

// NewPathValidator creates a new PathValidator instance
func NewPathValidator() *PathValidator {
	return &PathValidator{}
}

// ValidateSecurePath validates that a path doesn't contain directory traversal attempts
func (pv *PathValidator) ValidateSecurePath(name string) error {
	if strings.Contains(name, "..") || strings.Contains(name, "/") || strings.Contains(name, "\\") {
		return NewValidationError("path contains invalid characters or directory traversal: " + name)
	}
	return nil
}

// SecureJoin safely joins a base directory with a relative path, ensuring the result stays within the base
func (pv *PathValidator) SecureJoin(baseDir, relativePath string) (string, error) {
	// Validate the relative path first
	if err := pv.ValidateSecurePath(relativePath); err != nil {
		return "", err
	}

	path := filepath.Join(baseDir, relativePath)

	// Ensure the resulting path is within the base directory
	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		return "", NewFileError("get absolute base directory", baseDir, err)
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", NewFileError("get absolute path", path, err)
	}

	relPath, err := filepath.Rel(absBaseDir, absPath)
	if err != nil || strings.HasPrefix(relPath, "..") {
		return "", NewValidationError("path escapes base directory: " + relativePath)
	}

	return path, nil
}

// InputValidator provides common input validation utilities
type InputValidator struct{}

// NewInputValidator creates a new InputValidator instance
func NewInputValidator() *InputValidator {
	return &InputValidator{}
}

// NotEmpty validates that a string is not empty
func (iv *InputValidator) NotEmpty(value, fieldName string) error {
	if value == "" {
		return NewValidationError(fieldName + " cannot be empty")
	}
	return nil
}

// NotNil validates that a pointer is not nil
func (iv *InputValidator) NotNil(value interface{}, fieldName string) error {
	if value == nil {
		return NewValidationError(fieldName + " cannot be nil")
	}
	return nil
}

// ValidateFileExtension validates that a filename has the expected extension
func (iv *InputValidator) ValidateFileExtension(filename, expectedExt string) error {
	if !strings.HasSuffix(filename, expectedExt) {
		return NewValidationError("file must have " + expectedExt + " extension: " + filename)
	}
	return nil
}
