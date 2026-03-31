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

func TestFindFilesByExtension_MatchesFiles(t *testing.T) {
	dir := t.TempDir()

	for _, name := range []string{"a.txt", "b.txt", "c.log"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("content"), FileModeReadable); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}

	files, err := FindFilesByExtension(dir, ".txt")
	if err != nil {
		t.Fatalf("FindFilesByExtension() returned error: %v", err)
	}

	if len(files) != 2 {
		t.Errorf("FindFilesByExtension() returned %d files, want 2", len(files))
	}
}

func TestFindFilesByExtension_WithoutLeadingDot(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "file.txt"), []byte("content"), FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	files, err := FindFilesByExtension(dir, "txt")
	if err != nil {
		t.Fatalf("FindFilesByExtension() returned error: %v", err)
	}

	if len(files) != 1 {
		t.Errorf("FindFilesByExtension() returned %d files, want 1 (extension without dot should work)", len(files))
	}
}

func TestFindFilesByExtension_WithLeadingDot(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "file.txt"), []byte("content"), FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	files, err := FindFilesByExtension(dir, ".txt")
	if err != nil {
		t.Fatalf("FindFilesByExtension() returned error: %v", err)
	}

	if len(files) != 1 {
		t.Errorf("FindFilesByExtension() returned %d files, want 1", len(files))
	}
}

func TestFindFilesByExtension_Recursive(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "a", "b")

	if err := os.MkdirAll(sub, 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := os.WriteFile(filepath.Join(dir, "root.go"), []byte("content"), FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := os.WriteFile(filepath.Join(sub, "deep.go"), []byte("content"), FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	files, err := FindFilesByExtension(dir, ".go")
	if err != nil {
		t.Fatalf("FindFilesByExtension() returned error: %v", err)
	}

	if len(files) != 2 {
		t.Errorf("FindFilesByExtension() returned %d files, want 2 (should search recursively)", len(files))
	}
}

func TestFindFilesByExtension_NoMatches(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "file.txt"), []byte("content"), FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	files, err := FindFilesByExtension(dir, ".go")
	if err != nil {
		t.Fatalf("FindFilesByExtension() returned error: %v", err)
	}

	if len(files) != 0 {
		t.Errorf("FindFilesByExtension() returned %d files, want 0", len(files))
	}
}

func TestFindFilesByExtension_EmptyDirectory(t *testing.T) {
	dir := t.TempDir()

	files, err := FindFilesByExtension(dir, ".txt")
	if err != nil {
		t.Fatalf("FindFilesByExtension() returned error: %v", err)
	}

	if len(files) != 0 {
		t.Errorf("FindFilesByExtension() returned %d files, want 0", len(files))
	}
}

func TestFindFilesByExtension_ExcludesDirectories(t *testing.T) {
	dir := t.TempDir()

	// Create a directory with a matching extension
	if err := os.Mkdir(filepath.Join(dir, "data.txt"), 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	// Create a regular file with a matching extension
	if err := os.WriteFile(filepath.Join(dir, "real.txt"), []byte("content"), FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	files, err := FindFilesByExtension(dir, ".txt")
	if err != nil {
		t.Fatalf("FindFilesByExtension() returned error: %v", err)
	}

	if len(files) != 1 {
		t.Errorf("FindFilesByExtension() returned %d files, want 1 (should exclude directories)", len(files))
	}

	if len(files) > 0 && filepath.Base(files[0]) != "real.txt" {
		t.Errorf("FindFilesByExtension() file = %q, want 'real.txt'", files[0])
	}
}

func TestFindFilesByExtension_SymlinksExcluded(t *testing.T) {
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

	files, err := FindFilesByExtension(dir, ".txt")
	if err != nil {
		t.Fatalf("FindFilesByExtension() returned error: %v", err)
	}

	// Only real.txt — symlink is not a regular file
	if len(files) != 1 {
		t.Errorf("FindFilesByExtension() returned %d files, want 1 (symlinks should be excluded)", len(files))
	}
}

func TestFindFilesByExtension_ContinuesPastUnreadableDirectories(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("permission-based test not reliable on Windows")
	}

	dir := t.TempDir()
	restricted := filepath.Join(dir, "restricted")
	accessible := filepath.Join(dir, "accessible")

	if err := os.Mkdir(restricted, 0000); err != nil {
		t.Fatalf("setup: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(restricted, 0755) })

	if err := os.Mkdir(accessible, 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := os.WriteFile(filepath.Join(accessible, "found.txt"), []byte("content"), FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	files, err := FindFilesByExtension(dir, ".txt")
	if err != nil {
		t.Fatalf("FindFilesByExtension() returned error: %v (should continue past unreadable directories)", err)
	}

	if len(files) != 1 {
		t.Errorf("FindFilesByExtension() returned %d files, want 1", len(files))
	}
}

func TestFindFilesByExtension_ResultsAreConsistent(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "sub")

	if err := os.Mkdir(sub, 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	for _, name := range []string{"c.txt", "a.txt"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("content"), FileModeReadable); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}

	if err := os.WriteFile(filepath.Join(sub, "b.txt"), []byte("content"), FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	files, err := FindFilesByExtension(dir, ".txt")
	if err != nil {
		t.Fatalf("FindFilesByExtension() returned error: %v", err)
	}

	if len(files) != 3 {
		t.Fatalf("FindFilesByExtension() returned %d files, want 3", len(files))
	}

	// WalkDir returns entries in lexicographic order
	if !sort.StringsAreSorted(files) { //nolint:revive
		t.Errorf("FindFilesByExtension() results are not sorted: %v", files)
	}
}

func TestFindFilesByExtension_ReturnsNonNilSlice(t *testing.T) {
	dir := t.TempDir()

	files, err := FindFilesByExtension(dir, ".xyz")
	if err != nil {
		t.Fatalf("FindFilesByExtension() returned error: %v", err)
	}

	if files == nil {
		t.Error("FindFilesByExtension() returned nil, want non-nil empty slice")
	}
}

func TestFindFilesByExtension_CompoundExtension(t *testing.T) {
	dir := t.TempDir()

	for _, name := range []string{"archive.tar.gz", "readme.gz", "data.tar"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("content"), FileModeReadable); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}

	// filepath.Ext returns only the last extension
	files, err := FindFilesByExtension(dir, ".gz")
	if err != nil {
		t.Fatalf("FindFilesByExtension() returned error: %v", err)
	}

	if len(files) != 2 {
		t.Errorf("FindFilesByExtension() returned %d files, want 2 (should match .tar.gz and .gz)", len(files))
	}
}
