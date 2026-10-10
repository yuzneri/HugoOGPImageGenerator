package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// writeContentFiles creates empty files at the given paths relative to root.
func writeContentFiles(t *testing.T, root string, paths ...string) {
	t.Helper()
	for _, p := range paths {
		full := filepath.Join(root, p)
		if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte("---\ntitle: x\n---\n"), 0644); err != nil {
			t.Fatal(err)
		}
	}
}

func TestResolveArticle(t *testing.T) {
	root := t.TempDir()
	writeContentFiles(t, root,
		"posts/bundle/index.md",
		"posts/_index.md",
		"posts/single.md",
		"both/index.md",
		"both/_index.md",
	)
	os.MkdirAll(filepath.Join(root, "empty"), 0755)
	os.WriteFile(filepath.Join(root, "posts/image.png"), nil, 0644)

	tests := []struct {
		name string
		path string
		want Article
	}{
		{"bundle directory", "posts/bundle", Article{"posts/bundle/index.md", "posts/bundle", "posts/bundle"}},
		{"bundle index file", "posts/bundle/index.md", Article{"posts/bundle/index.md", "posts/bundle", "posts/bundle"}},
		{"section directory", "posts", Article{"posts/_index.md", "posts", "posts"}},
		{"section index file", "posts/_index.md", Article{"posts/_index.md", "posts", "posts"}},
		{"single file", "posts/single.md", Article{"posts/single.md", "posts", "posts/single"}},
		{"directory prefers index.md", "both", Article{"both/index.md", "both", "both"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolveArticle(filepath.Join(root, tt.path))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			want := Article{
				File:    filepath.Join(root, tt.want.File),
				Dir:     filepath.Join(root, tt.want.Dir),
				PageDir: filepath.Join(root, tt.want.PageDir),
			}
			if *got != want {
				t.Errorf("got %+v, want %+v", *got, want)
			}
		})
	}

	for _, p := range []string{"missing", "missing.md", "empty", "posts/image.png"} {
		t.Run("error "+p, func(t *testing.T) {
			if _, err := resolveArticle(filepath.Join(root, p)); err == nil {
				t.Errorf("expected error for %s", p)
			}
		})
	}
}

func TestFindArticles(t *testing.T) {
	root := t.TempDir()
	writeContentFiles(t, root,
		"_index.md",
		"about.md",
		"posts/_index.md",
		"posts/single.md",
		"posts/bundle/index.md",
		"posts/bundle/resource.md",
		"posts/bundle/sub/nested.md",
		"posts/nested/deep.md",
	)
	os.WriteFile(filepath.Join(root, "posts/notes.txt"), nil, 0644)

	got, err := findArticles(root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var rel []string
	for _, f := range got {
		r, _ := filepath.Rel(root, f)
		rel = append(rel, filepath.ToSlash(r))
	}

	want := []string{
		"_index.md",
		"about.md",
		"posts/_index.md",
		"posts/bundle/index.md",
		"posts/nested/deep.md",
		"posts/single.md",
	}
	if !reflect.DeepEqual(rel, want) {
		t.Errorf("got %v, want %v", rel, want)
	}
}
