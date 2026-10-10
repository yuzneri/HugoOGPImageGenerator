package main

import (
	"os"
	"path/filepath"
	"testing"
)

// unsetBuildDraftsEnv removes HUGO_BUILDDRAFTS for the duration of the test.
func unsetBuildDraftsEnv(t *testing.T) {
	t.Helper()
	if value, ok := os.LookupEnv(BuildDraftsEnv); ok {
		t.Cleanup(func() { os.Setenv(BuildDraftsEnv, value) })
	}
	os.Unsetenv(BuildDraftsEnv)
}

func writeProjectFiles(t *testing.T, root string, files map[string]string) {
	t.Helper()
	for path, body := range files {
		full := filepath.Join(root, path)
		if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(body), 0644); err != nil {
			t.Fatal(err)
		}
	}
}

func TestResolveBuildDrafts_ConfigFiles(t *testing.T) {
	unsetBuildDraftsEnv(t)

	tests := []struct {
		name       string
		files      map[string]string
		want       bool
		wantSource string
	}{
		{"no config", nil, false, ""},
		{"toml", map[string]string{"hugo.toml": "buildDrafts = true\n"}, true, "hugo.toml"},
		{"yaml", map[string]string{"hugo.yaml": "buildDrafts: true\n"}, true, "hugo.yaml"},
		{"json", map[string]string{"hugo.json": `{"buildDrafts": true}`}, true, "hugo.json"},
		{"legacy config.toml", map[string]string{"config.toml": "buildDrafts = true\n"}, true, "config.toml"},
		{"case-insensitive key", map[string]string{"hugo.toml": "builddrafts = true\n"}, true, "hugo.toml"},
		{"explicit false", map[string]string{"hugo.toml": "buildDrafts = false\n"}, false, "hugo.toml"},
		{"key not set", map[string]string{"hugo.toml": "title = 'x'\n"}, false, ""},
		{"hugo.toml wins over hugo.yaml", map[string]string{"hugo.toml": "buildDrafts = false\n", "hugo.yaml": "buildDrafts: true\n"}, false, "hugo.toml"},
		{"hugo.* wins over config.*", map[string]string{"hugo.yaml": "buildDrafts: false\n", "config.toml": "buildDrafts = true\n"}, false, "hugo.yaml"},
		{"config directory", map[string]string{"config/_default/hugo.toml": "buildDrafts = true\n"}, true, "config/_default/hugo.toml"},
		{"root file wins over config directory", map[string]string{"hugo.toml": "buildDrafts = false\n", "config/_default/hugo.toml": "buildDrafts = true\n"}, false, "hugo.toml"},
		{"config directory when root file lacks the key", map[string]string{"hugo.toml": "title = 'x'\n", "config/_default/hugo.toml": "buildDrafts = true\n"}, true, "config/_default/hugo.toml"},
		{"table key is not top-level", map[string]string{"hugo.toml": "[params]\nbuildDrafts = true\n"}, false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			writeProjectFiles(t, root, tt.files)

			got, source, err := resolveBuildDrafts(root)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("buildDrafts = %v, want %v", got, tt.want)
			}
			wantSource := ""
			if tt.wantSource != "" {
				wantSource = filepath.Join(root, tt.wantSource)
			}
			if source != wantSource {
				t.Errorf("source = %q, want %q", source, wantSource)
			}
		})
	}
}

func TestResolveBuildDrafts_EnvOverridesConfig(t *testing.T) {
	root := t.TempDir()
	writeProjectFiles(t, root, map[string]string{"hugo.toml": "buildDrafts = true\n"})

	t.Setenv(BuildDraftsEnv, "false")
	got, source, err := resolveBuildDrafts(root)
	if err != nil || got || source != BuildDraftsEnv {
		t.Errorf("got (%v, %q, %v), want (false, %q, nil)", got, source, err, BuildDraftsEnv)
	}

	t.Setenv(BuildDraftsEnv, "true")
	got, _, err = resolveBuildDrafts(t.TempDir())
	if err != nil || !got {
		t.Errorf("got (%v, %v), want (true, nil)", got, err)
	}
}

func TestResolveBuildDrafts_Errors(t *testing.T) {
	t.Run("invalid env value", func(t *testing.T) {
		t.Setenv(BuildDraftsEnv, "yes please")
		if _, _, err := resolveBuildDrafts(t.TempDir()); err == nil {
			t.Error("expected error")
		}
	})

	files := map[string]map[string]string{
		"non-boolean value": {"hugo.toml": "buildDrafts = 'true'\n"},
		"broken toml":       {"hugo.toml": "buildDrafts = \n"},
	}
	for name, f := range files {
		t.Run(name, func(t *testing.T) {
			unsetBuildDraftsEnv(t)
			root := t.TempDir()
			writeProjectFiles(t, root, f)
			if _, _, err := resolveBuildDrafts(root); err == nil {
				t.Error("expected error")
			}
		})
	}
}
