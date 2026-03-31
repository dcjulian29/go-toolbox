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

package filesystem

import (
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"testing"
)

func TestSearchForFiles_MatchesFiles(t *testing.T) {
	dir := t.TempDir()

	for _, name := range []string{"a.txt", "b.txt", "c.log"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("content"), FileModeReadable); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}

	results, err := SearchForFiles(dir, "*.txt")
	if err != nil {
		t.Fatalf("SearchForFiles() returned error: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("SearchForFiles() returned %d results, want 2", len(results))
	}
}

func TestSearchForFiles_SingleMatch(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "only.txt"), []byte("content"), FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	results, err := SearchForFiles(dir, "only.txt")
	if err != nil {
		t.Fatalf("SearchForFiles() returned error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("SearchForFiles() returned %d results, want 1", len(results))
	}

	if filepath.Base(results[0]) != "only.txt" {
		t.Errorf("SearchForFiles() result = %q, want path ending in 'only.txt'", results[0])
	}
}

func TestSearchForFiles_NoMatch(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "file.log"), []byte("content"), FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	_, err := SearchForFiles(dir, "*.txt")
	if err == nil {
		t.Error("SearchForFiles() should return error when no files match")
	}
}

func TestSearchForFiles_BadPattern(t *testing.T) {
	dir := t.TempDir()

	_, err := SearchForFiles(dir, "[invalid")
	if err == nil {
		t.Error("SearchForFiles() should return error for malformed glob pattern")
	}
}

func TestSearchForFiles_NonExistentDirectory(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nonexistent")

	_, err := SearchForFiles(path, "*.txt")
	if err == nil {
		t.Error("SearchForFiles() should return error for non-existent directory")
	}
}

func TestSearchForFiles_ExcludesDirectories(t *testing.T) {
	dir := t.TempDir()

	// Create a file and a directory that both match the pattern
	if err := os.WriteFile(filepath.Join(dir, "match.log"), []byte("content"), FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := os.Mkdir(filepath.Join(dir, "also.log"), 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	results, err := SearchForFiles(dir, "*.log")
	if err != nil {
		t.Fatalf("SearchForFiles() returned error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("SearchForFiles() returned %d results, want 1 (directory should be excluded)", len(results))
	}

	if filepath.Base(results[0]) != "match.log" {
		t.Errorf("SearchForFiles() result = %q, want path ending in 'match.log'", results[0])
	}
}

func TestSearchForFiles_OnlyDirectoriesMatch(t *testing.T) {
	dir := t.TempDir()

	// Only a directory matches the pattern
	if err := os.Mkdir(filepath.Join(dir, "subdir.txt"), 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	_, err := SearchForFiles(dir, "*.txt")
	if err == nil {
		t.Error("SearchForFiles() should return error when only directories match the pattern")
	}
}

func TestSearchForFiles_ResultsAreSorted(t *testing.T) {
	dir := t.TempDir()

	for _, name := range []string{"c.txt", "a.txt", "b.txt"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("content"), FileModeReadable); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}

	results, err := SearchForFiles(dir, "*.txt")
	if err != nil {
		t.Fatalf("SearchForFiles() returned error: %v", err)
	}

	if !sort.StringsAreSorted(results) { //nolint:revive
		t.Errorf("SearchForFiles() results are not sorted: %v", results)
	}
}

func TestSearchForFiles_EmptyDirectory(t *testing.T) {
	dir := t.TempDir()

	_, err := SearchForFiles(dir, "*")
	if err == nil {
		t.Error("SearchForFiles() should return error for empty directory")
	}
}

func TestSearchForFiles_SymlinkToFileIncluded(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink creation may require elevated privileges on Windows")
	}

	dir := t.TempDir()

	realFile := filepath.Join(dir, "real.txt")
	if err := os.WriteFile(realFile, []byte("content"), FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := os.Symlink(realFile, filepath.Join(dir, "link.txt")); err != nil {
		t.Fatalf("setup: failed to create symlink: %v", err)
	}

	results, err := SearchForFiles(dir, "*.txt")
	if err != nil {
		t.Fatalf("SearchForFiles() returned error: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("SearchForFiles() returned %d results, want 2 (symlink to file should be included)", len(results))
	}
}

func TestSearchForFiles_NonRecursive(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "subdir")

	if err := os.Mkdir(sub, 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	// File only exists in subdirectory, not in the searched path
	if err := os.WriteFile(filepath.Join(sub, "nested.txt"), []byte("content"), FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	_, err := SearchForFiles(dir, "nested.txt")
	if err == nil {
		t.Error("SearchForFiles() should not find files in subdirectories")
	}
}
