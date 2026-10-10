package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

// BuildDraftsEnv is the environment variable Hugo reads for the buildDrafts setting.
const BuildDraftsEnv = "HUGO_BUILDDRAFTS"

// hugoConfigNames lists project configuration files in Hugo's lookup order.
var hugoConfigNames = []string{
	"hugo.toml", "hugo.yaml", "hugo.json",
	"config.toml", "config.yaml", "config.json",
}

// resolveBuildDrafts reports whether Hugo publishes drafts for the project, and where that was decided.
// Like Hugo, the HUGO_BUILDDRAFTS environment variable takes precedence over the configuration file.
// The project configuration file is read first, then config/_default/.
// Per-environment directories (config/production etc.) and the --config flag are not followed.
func resolveBuildDrafts(projectRoot string) (bool, string, error) {
	if value, ok := os.LookupEnv(BuildDraftsEnv); ok {
		drafts, err := strconv.ParseBool(value)
		if err != nil {
			return false, "", NewValidationError(fmt.Sprintf("invalid %s value %q: %v", BuildDraftsEnv, value, err))
		}
		return drafts, BuildDraftsEnv, nil
	}

	for _, dir := range []string{projectRoot, filepath.Join(projectRoot, "config", "_default")} {
		path, ok := findHugoConfig(dir)
		if !ok {
			continue
		}
		drafts, found, err := readBuildDrafts(path)
		if err != nil {
			return false, "", err
		}
		if found {
			return drafts, path, nil
		}
	}

	return false, "", nil
}

// findHugoConfig returns the first Hugo configuration file in dir.
func findHugoConfig(dir string) (string, bool) {
	for _, name := range hugoConfigNames {
		path := filepath.Join(dir, name)
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			return path, true
		}
	}
	return "", false
}

// readBuildDrafts reads the top-level buildDrafts key from a Hugo configuration file.
// Hugo configuration keys are case-insensitive.
func readBuildDrafts(path string) (bool, bool, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return false, false, NewFileError("read", path, err)
	}

	settings := map[string]interface{}{}
	switch filepath.Ext(path) {
	case ".toml":
		err = toml.Unmarshal(content, &settings)
	case ".yaml":
		err = yaml.Unmarshal(content, &settings)
	case ".json":
		err = json.Unmarshal(content, &settings)
	}
	if err != nil {
		return false, false, NewConfigError(fmt.Sprintf("failed to parse Hugo configuration %s", path), err)
	}

	for key, value := range settings {
		if !strings.EqualFold(key, "buildDrafts") {
			continue
		}
		drafts, ok := value.(bool)
		if !ok {
			return false, false, NewValidationError(fmt.Sprintf("buildDrafts in %s must be a boolean, got %v", path, value))
		}
		return drafts, true, nil
	}

	return false, false, nil
}
