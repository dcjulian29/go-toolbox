/*
Copyright © 2026 Julian Easterling

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package configuration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type testConfig struct {
	Name  string   `yaml:"name"`
	Count int      `yaml:"count"`
	Items []string `yaml:"items"`
}

// setHome points the user's home directory at a fresh temp dir for the duration
// of the test, covering both the Unix (HOME) and Windows (USERPROFILE) sources
// that os.UserHomeDir consults.
func setHome(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Setenv("USERPROFILE", dir)

	return dir
}

func writeConfig(t *testing.T, home, name, content string) string {
	t.Helper()

	path := filepath.Join(home, ".config", name)

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	return path
}

func TestPath_UnderConfigDir(t *testing.T) {
	home := setHome(t)

	got, err := New[testConfig]("mytool.yml").Path()
	if err != nil {
		t.Fatalf("Path() returned error: %v", err)
	}

	want := filepath.Join(home, ".config", "mytool.yml")
	if got != want {
		t.Errorf("Path() = %q, want %q", got, want)
	}
}

func TestLoad_MissingFileReturnsZeroValue(t *testing.T) {
	setHome(t)

	cfg, err := New[testConfig]("mytool.yml").Load()
	if err != nil {
		t.Fatalf("Load() returned error for missing file: %v", err)
	}

	if cfg.Name != "" || cfg.Count != 0 || len(cfg.Items) != 0 {
		t.Errorf("Load() of missing file = %+v, want zero value", cfg)
	}
}

func TestSaveThenLoad_RoundTrips(t *testing.T) {
	setHome(t)

	want := testConfig{Name: "acme", Count: 3, Items: []string{"a", "b"}}

	if err := New[testConfig]("mytool.yml").Save(&want); err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}

	// A fresh instance forces a read from disk rather than returning a cache.
	got, err := New[testConfig]("mytool.yml").Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if got.Name != want.Name || got.Count != want.Count || strings.Join(got.Items, ",") != "a,b" {
		t.Errorf("round trip = %+v, want %+v", got, want)
	}
}

func TestSave_WritesYAMLToPath(t *testing.T) {
	home := setHome(t)

	if err := New[testConfig]("mytool.yml").Save(&testConfig{Name: "acme"}); err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(home, ".config", "mytool.yml"))
	if err != nil {
		t.Fatalf("reading saved file: %v", err)
	}

	if !strings.Contains(string(data), "name: acme") {
		t.Errorf("saved file = %q, want it to contain %q", data, "name: acme")
	}
}

func TestSave_OverwritesExistingFile(t *testing.T) {
	setHome(t)

	f := New[testConfig]("mytool.yml")

	if err := f.Save(&testConfig{Name: "first"}); err != nil {
		t.Fatalf("first Save() returned error: %v", err)
	}

	if err := f.Save(&testConfig{Name: "second"}); err != nil {
		t.Fatalf("second Save() returned error: %v", err)
	}

	got, err := New[testConfig]("mytool.yml").Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if got.Name != "second" {
		t.Errorf("after overwrite, Name = %q, want %q", got.Name, "second")
	}
}

func TestLoad_CachesFirstResult(t *testing.T) {
	home := setHome(t)
	writeConfig(t, home, "mytool.yml", "name: original\n")

	f := New[testConfig]("mytool.yml")

	first, err := f.Load()
	if err != nil {
		t.Fatalf("first Load() returned error: %v", err)
	}

	if first.Name != "original" {
		t.Fatalf("first Load() Name = %q, want %q", first.Name, "original")
	}

	// Changing the file after the first load must not affect the cached value.
	writeConfig(t, home, "mytool.yml", "name: changed\n")

	second, err := f.Load()
	if err != nil {
		t.Fatalf("second Load() returned error: %v", err)
	}

	if second.Name != "original" {
		t.Errorf("second Load() Name = %q, want cached %q", second.Name, "original")
	}
}

func TestSave_UpdatesCache(t *testing.T) {
	setHome(t)

	f := New[testConfig]("mytool.yml")

	// Prime the cache from a missing file (once fires, zero value cached).
	if _, err := f.Load(); err != nil {
		t.Fatalf("priming Load() returned error: %v", err)
	}

	if err := f.Save(&testConfig{Name: "saved"}); err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}

	got, err := f.Load()
	if err != nil {
		t.Fatalf("Load() after Save() returned error: %v", err)
	}

	if got.Name != "saved" {
		t.Errorf("Load() after Save() Name = %q, want %q", got.Name, "saved")
	}
}

func TestSave_NilReturnsError(t *testing.T) {
	setHome(t)

	if err := New[testConfig]("mytool.yml").Save(nil); err == nil {
		t.Error("Save(nil) returned nil error, want error")
	}
}

func TestLoad_InvalidYAMLReturnsError(t *testing.T) {
	home := setHome(t)
	writeConfig(t, home, "mytool.yml", "name: [unterminated\n")

	if _, err := New[testConfig]("mytool.yml").Load(); err == nil {
		t.Error("Load() of invalid YAML returned nil error, want error")
	}
}
