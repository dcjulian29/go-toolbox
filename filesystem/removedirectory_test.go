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

func TestRemoveDirectory_RemovesEmptyDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "empty")

	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}

	err := RemoveDirectory(dir)
	if err != nil {
		t.Fatalf("RemoveDirectory returned unexpected error: %v", err)
	}

	if DirectoryExist(dir) {
		t.Errorf("directory still exists at %q after removal", dir)
	}
}

func TestRemoveDirectory_RemovesDirectoryWithFiles(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "withfiles")

	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}

	for _, name := range []string{"a.txt", "b.txt", "c.txt"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("data"), FileModeReadable); err != nil {
			t.Fatalf("failed to create file %q: %v", name, err)
		}
	}

	err := RemoveDirectory(dir)
	if err != nil {
		t.Fatalf("RemoveDirectory returned unexpected error: %v", err)
	}

	if DirectoryExist(dir) {
		t.Errorf("directory still exists at %q after removal", dir)
	}
}

func TestRemoveDirectory_RemovesNestedDirectories(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested")
	deep := filepath.Join(dir, "a", "b", "c")

	if err := os.MkdirAll(deep, 0o755); err != nil {
		t.Fatalf("failed to create nested directories: %v", err)
	}

	if err := os.WriteFile(filepath.Join(deep, "file.txt"), []byte("deep"), FileModeReadable); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	err := RemoveDirectory(dir)
	if err != nil {
		t.Fatalf("RemoveDirectory returned unexpected error: %v", err)
	}

	if DirectoryExist(dir) {
		t.Errorf("directory still exists at %q after removal", dir)
	}
}

func TestRemoveDirectory_RemovesHiddenFiles(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "withdotfiles")

	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}

	// Create visible and hidden files
	for _, name := range []string{"visible.txt", ".hidden", ".gitignore"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("data"), FileModeReadable); err != nil {
			t.Fatalf("failed to create file %q: %v", name, err)
		}
	}

	err := RemoveDirectory(dir)
	if err != nil {
		t.Fatalf("RemoveDirectory returned unexpected error: %v", err)
	}

	if DirectoryExist(dir) {
		t.Errorf("directory still exists at %q after removal (dotfiles may have been left behind)", dir)
	}
}

func TestRemoveDirectory_ReturnsNilForNonExistentPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "does-not-exist")

	err := RemoveDirectory(path)
	if err != nil {
		t.Fatalf("RemoveDirectory should return nil for non-existent path, got: %v", err)
	}
}

func TestRemoveDirectory_ReturnsErrorForFile(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "file.txt")

	if err := os.WriteFile(filePath, []byte("data"), FileModeReadable); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	err := RemoveDirectory(filePath)
	if err == nil {
		t.Fatal("RemoveDirectory should return an error when path is a file")
	}

	// The file should not have been deleted
	if !FileExist(filePath) {
		t.Error("file was unexpectedly removed")
	}
}

func TestRemoveDirectory_EmptyPath(t *testing.T) {
	// os.Stat("") returns an error; os.IsNotExist behavior may vary.
	// This test documents the current behavior.
	err := RemoveDirectory("")

	if err != nil {
		t.Logf("RemoveDirectory returned error for empty path: %v", err)
	} else {
		t.Log("RemoveDirectory returned nil for empty path")
	}
}

func TestRemoveDirectory_StatPermissionError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping permission test on Windows (permission model differs)")
	}

	if os.Getuid() == 0 {
		t.Skip("skipping permission test when running as root")
	}

	dir := t.TempDir()
	restricted := filepath.Join(dir, "noaccess")

	if err := os.Mkdir(restricted, 0o755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}

	inner := filepath.Join(restricted, "inner")

	if err := os.Mkdir(inner, 0o755); err != nil {
		t.Fatalf("failed to create inner directory: %v", err)
	}

	// Remove all permissions from parent so stat on inner fails with permission denied
	if err := os.Chmod(restricted, 0o000); err != nil {
		t.Fatalf("failed to change permissions: %v", err)
	}

	t.Cleanup(func() {
		os.Chmod(restricted, 0o755) //nolint:errcheck
	})

	err := RemoveDirectory(inner)
	if err == nil {
		t.Fatal("RemoveDirectory should return error when stat fails with permission denied")
	}
}

func TestRemoveDirectory_RemovalPermissionDenied(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping permission test on Windows (permission model differs)")
	}

	if os.Getuid() == 0 {
		t.Skip("skipping permission test when running as root")
	}

	dir := t.TempDir()
	target := filepath.Join(dir, "protected")

	if err := os.Mkdir(target, 0o755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}

	if err := os.WriteFile(filepath.Join(target, "file.txt"), []byte("data"), FileModeReadable); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	// Remove write permission from parent so target cannot be removed
	if err := os.Chmod(dir, 0o555); err != nil {
		t.Fatalf("failed to change permissions: %v", err)
	}

	t.Cleanup(func() {
		os.Chmod(dir, 0o755) //nolint:errcheck
	})

	err := RemoveDirectory(target)
	if err == nil {
		t.Fatal("RemoveDirectory should return error when parent is not writable")
	}
}

func TestRemoveDirectory_IsIdempotent(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "idempotent")

	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}

	for i := range 3 {
		if err := RemoveDirectory(dir); err != nil {
			t.Fatalf("call %d: RemoveDirectory returned unexpected error: %v", i+1, err)
		}
	}
}
