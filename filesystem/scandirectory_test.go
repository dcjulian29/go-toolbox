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
	"testing"
)

func TestScanDirectory_ReturnsFilesAndFolders(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "subdir")

	if err := os.Mkdir(sub, 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	for _, name := range []string{"a.txt", "b.txt"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("content"), FileModeReadable); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}

	if err := os.WriteFile(filepath.Join(sub, "c.txt"), []byte("content"), FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	folders, files, err := ScanDirectory(dir, []string{})
	if err != nil {
		t.Fatalf("ScanDirectory() returned error: %v", err)
	}

	// Root dir + subdir
	if len(folders) != 2 {
		t.Errorf("ScanDirectory() returned %d folders, want 2", len(folders))
	}

	// a.txt + b.txt + c.txt
	if len(files) != 3 {
		t.Errorf("ScanDirectory() returned %d files, want 3", len(files))
	}
}

func TestScanDirectory_EmptyDirectory(t *testing.T) {
	dir := t.TempDir()

	folders, files, err := ScanDirectory(dir, []string{})
	if err != nil {
		t.Fatalf("ScanDirectory() returned error: %v", err)
	}

	// Root directory itself is included
	if len(folders) != 1 {
		t.Errorf("ScanDirectory() returned %d folders, want 1 (root only)", len(folders))
	}

	if len(files) != 0 {
		t.Errorf("ScanDirectory() returned %d files, want 0", len(files))
	}
}

func TestScanDirectory_RootDirectoryIncluded(t *testing.T) {
	dir := t.TempDir()

	folders, _, err := ScanDirectory(dir, []string{})
	if err != nil {
		t.Fatalf("ScanDirectory() returned error: %v", err)
	}

	if len(folders) == 0 {
		t.Fatal("ScanDirectory() returned no folders, expected root directory")
	}

	if folders[0] != dir {
		t.Errorf("ScanDirectory() first folder = %q, want root %q", folders[0], dir)
	}
}

func TestScanDirectory_NestedDirectories(t *testing.T) {
	dir := t.TempDir()
	deep := filepath.Join(dir, "a", "b", "c")

	if err := os.MkdirAll(deep, 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := os.WriteFile(filepath.Join(deep, "deep.txt"), []byte("content"), FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	folders, files, err := ScanDirectory(dir, []string{})
	if err != nil {
		t.Fatalf("ScanDirectory() returned error: %v", err)
	}

	// dir, a, b, c
	if len(folders) != 4 {
		t.Errorf("ScanDirectory() returned %d folders, want 4", len(folders))
	}

	if len(files) != 1 {
		t.Errorf("ScanDirectory() returned %d files, want 1", len(files))
	}
}

func TestScanDirectory_IgnoreExcludesFiles(t *testing.T) {
	dir := t.TempDir()

	for _, name := range []string{"keep.txt", "ignore_me.txt", "also_keep.txt"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("content"), FileModeReadable); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}

	_, files, err := ScanDirectory(dir, []string{"ignore_me"})
	if err != nil {
		t.Fatalf("ScanDirectory() returned error: %v", err)
	}

	if len(files) != 2 {
		t.Errorf("ScanDirectory() returned %d files, want 2", len(files))
	}

	for _, f := range files {
		if filepath.Base(f) == "ignore_me.txt" {
			t.Error("ScanDirectory() should have excluded 'ignore_me.txt'")
		}
	}
}

func TestScanDirectory_IgnoreSkipsDirectoryDescent(t *testing.T) {
	dir := t.TempDir()
	ignored := filepath.Join(dir, "node_modules")
	nested := filepath.Join(ignored, "pkg")

	if err := os.MkdirAll(nested, 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := os.WriteFile(filepath.Join(nested, "hidden.txt"), []byte("content"), FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := os.WriteFile(filepath.Join(dir, "visible.txt"), []byte("content"), FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	folders, files, err := ScanDirectory(dir, []string{"node_modules"})
	if err != nil {
		t.Fatalf("ScanDirectory() returned error: %v", err)
	}

	// Only root dir — node_modules and its children are skipped
	if len(folders) != 1 {
		t.Errorf("ScanDirectory() returned %d folders, want 1 (ignored dir should not be descended into)", len(folders))
	}

	// Only visible.txt — hidden.txt inside ignored dir is not returned
	if len(files) != 1 {
		t.Errorf("ScanDirectory() returned %d files, want 1", len(files))
	}

	if len(files) > 0 && filepath.Base(files[0]) != "visible.txt" {
		t.Errorf("ScanDirectory() file = %q, want 'visible.txt'", files[0])
	}
}

func TestScanDirectory_MultipleIgnorePatterns(t *testing.T) {
	dir := t.TempDir()

	for _, name := range []string{".git", "vendor", "src"} {
		if err := os.Mkdir(filepath.Join(dir, name), 0755); err != nil {
			t.Fatalf("setup: %v", err)
		}

		if err := os.WriteFile(filepath.Join(dir, name, "file.txt"), []byte("content"), FileModeReadable); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}

	folders, files, err := ScanDirectory(dir, []string{".git", "vendor"})
	if err != nil {
		t.Fatalf("ScanDirectory() returned error: %v", err)
	}

	// root + src
	if len(folders) != 2 {
		t.Errorf("ScanDirectory() returned %d folders, want 2", len(folders))
	}

	// Only src/file.txt
	if len(files) != 1 {
		t.Errorf("ScanDirectory() returned %d files, want 1", len(files))
	}
}

func TestScanDirectory_NilIgnoreList(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "file.txt"), []byte("content"), FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	_, files, err := ScanDirectory(dir, nil)
	if err != nil {
		t.Fatalf("ScanDirectory() returned error: %v", err)
	}

	if len(files) != 1 {
		t.Errorf("ScanDirectory() returned %d files, want 1", len(files))
	}
}

func TestScanDirectory_SymlinksExcluded(t *testing.T) {
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

	_, files, err := ScanDirectory(dir, []string{})
	if err != nil {
		t.Fatalf("ScanDirectory() returned error: %v", err)
	}

	// Only real.txt — symlink is not a regular file
	if len(files) != 1 {
		t.Errorf("ScanDirectory() returned %d files, want 1 (symlinks should be excluded)", len(files))
	}

	if len(files) > 0 && filepath.Base(files[0]) != "real.txt" {
		t.Errorf("ScanDirectory() file = %q, want 'real.txt'", files[0])
	}
}

func TestScanDirectory_IgnoreSubstringMatching(t *testing.T) {
	dir := t.TempDir()

	for _, name := range []string{"test_one.txt", "testing.txt", "latest.txt", "keep.txt"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("content"), FileModeReadable); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}

	_, files, err := ScanDirectory(dir, []string{"test"})
	if err != nil {
		t.Fatalf("ScanDirectory() returned error: %v", err)
	}

	// Only keep.txt — all others contain "test" as a substring
	if len(files) != 1 {
		t.Errorf("ScanDirectory() returned %d files, want 1 (substring matching should exclude 'test_one.txt', 'testing.txt', 'latest.txt')", len(files))
	}
}

func TestScanDirectory_ContinuesPastUnreadableDirectories(t *testing.T) {
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

	_, files, err := ScanDirectory(dir, []string{})
	if err != nil {
		t.Fatalf("ScanDirectory() returned error: %v (should continue past unreadable directories)", err)
	}

	if len(files) != 1 {
		t.Errorf("ScanDirectory() returned %d files, want 1", len(files))
	}
}
