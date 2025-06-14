package main

import (
	"strings"
)

const (
	TypeConfigExtension = ".yaml"
)

// loadTypeConfigSettings loads a type-specific configuration file as ConfigSettings.
// Returns nil if the file doesn't exist (this is not an error condition).
// Returns an error only if the file exists but cannot be parsed.
func loadTypeConfigSettings(configDir, typeName string) (*ConfigSettings, error) {
	// Validate input parameters
	if err := validateTypeConfigInput(configDir, typeName); err != nil {
		return nil, err
	}

	// Return nil for empty type names (not an error)
	if typeName == "" {
		return nil, nil
	}

	typeConfigPath, err := buildTypeConfigPath(configDir, typeName)
	if err != nil {
		return nil, err
	}

	// Check if file exists
	fileOps := NewFileOperations()
	exists, err := fileOps.ExistsAndIsFile(typeConfigPath)
	if err != nil {
		return nil, NewFileError("stat", typeConfigPath, err)
	}
	if !exists {
		return nil, nil // File doesn't exist, not an error
	}

	// Load and parse the configuration
	settings, err := loadAndParseTypeConfigSettings(typeConfigPath)
	if err != nil {
		return nil, err
	}

	return settings, nil
}

// validateTypeConfigInput validates the input parameters for type config loading.
func validateTypeConfigInput(configDir, typeName string) error {
	validator := NewInputValidator()
	pathValidator := NewPathValidator()

	if err := validator.NotEmpty(configDir, "config directory"); err != nil {
		return err
	}

	if err := pathValidator.ValidateSecurePath(typeName); err != nil {
		return err
	}

	return nil
}

// buildTypeConfigPath constructs the file path for a type-specific configuration file.
func buildTypeConfigPath(configDir, typeName string) (string, error) {
	pathValidator := NewPathValidator()
	// Convert type name to lowercase for file name consistency
	// This allows "Book" type to load "book.yaml" file
	fileName := strings.ToLower(typeName) + TypeConfigExtension

	return pathValidator.SecureJoin(configDir, fileName)
}

// loadAndParseTypeConfigSettings reads and parses a type configuration file as ConfigSettings.
func loadAndParseTypeConfigSettings(configPath string) (*ConfigSettings, error) {
	var settings ConfigSettings
	fileOps := NewFileOperations()

	err := fileOps.ReadYAMLFile(configPath, &settings)
	if err != nil {
		return nil, err
	}

	return &settings, nil
}
