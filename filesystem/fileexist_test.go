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

func createTestFile(t *testing.T, path string, content []byte) {
	t.Helper()

	if err := os.WriteFile(path, content, FileModeReadable); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
}

func TestFileExist_ReturnsTrueForExistingFile(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "exists.txt")
	createTestFile(t, filePath, []byte("hello"))

	if !FileExist(filePath) {
		t.Errorf("FileExist(%q) = false; expected true", filePath)
	}
}

func TestFileExist_ReturnsTrueForEmptyFile(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "empty.txt")
	createTestFile(t, filePath, []byte(""))

	if !FileExist(filePath) {
		t.Errorf("FileExist(%q) = false; expected true", filePath)
	}
}

func TestFileExist_ReturnsTrueForFileInNestedDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "a", "b", "c")

	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("failed to create nested directory: %v", err)
	}

	filePath := filepath.Join(dir, "deep.txt")
	createTestFile(t, filePath, []byte("deep"))

	if !FileExist(filePath) {
		t.Errorf("FileExist(%q) = false; expected true", filePath)
	}
}

func TestFileExist_ReturnsFalseForNonExistentPath(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "nope.txt")

	if FileExist(filePath) {
		t.Errorf("FileExist(%q) = true; expected false", filePath)
	}
}

func TestFileExist_ReturnsFalseForDirectory(t *testing.T) {
	dir := t.TempDir()

	if FileExist(dir) {
		t.Errorf("FileExist(%q) = true; expected false (path is a directory)", dir)
	}
}

func TestFileExist_ReturnsFalseForEmptyStringPath(t *testing.T) {
	if FileExist("") {
		t.Error("FileExist(\"\") = true; expected false")
	}
}

func TestFileExist_ReturnsTrueForSymlinkToFile(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping symlink test on Windows (requires elevated privileges)")
	}

	dir := t.TempDir()
	realFile := filepath.Join(dir, "real.txt")
	linkPath := filepath.Join(dir, "link.txt")
	createTestFile(t, realFile, []byte("target"))

	if err := os.Symlink(realFile, linkPath); err != nil {
		t.Skipf("symlinks not supported: %v", err)
	}

	if !FileExist(linkPath) {
		t.Errorf("FileExist(%q) = false; expected true (symlink points to a file)", linkPath)
	}
}

func TestFileExist_ReturnsFalseForSymlinkToDirectory(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping symlink test on Windows (requires elevated privileges)")
	}

	dir := t.TempDir()
	realDir := filepath.Join(dir, "realdir")
	linkPath := filepath.Join(dir, "link")

	if err := os.Mkdir(realDir, 0o755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}

	if err := os.Symlink(realDir, linkPath); err != nil {
		t.Skipf("symlinks not supported: %v", err)
	}

	if FileExist(linkPath) {
		t.Errorf("FileExist(%q) = true; expected false (symlink points to a directory)", linkPath)
	}
}

func TestFileExist_ReturnsFalseForBrokenSymlink(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping symlink test on Windows (requires elevated privileges)")
	}

	dir := t.TempDir()
	target := filepath.Join(dir, "deleted.txt")
	linkPath := filepath.Join(dir, "broken")

	createTestFile(t, target, []byte("temporary"))

	if err := os.Symlink(target, linkPath); err != nil {
		t.Skipf("symlinks not supported: %v", err)
	}

	// Remove the target so the symlink is now broken
	os.Remove(target)

	if FileExist(linkPath) {
		t.Errorf("FileExist(%q) = true; expected false (broken symlink)", linkPath)
	}
}

func TestFileExist_DoesNotPanicOnPermissionError(t *testing.T) {
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

	innerFile := filepath.Join(restricted, "secret.txt")

	if err := os.WriteFile(innerFile, []byte("secret"), FileModeReadable); err != nil {
		t.Fatalf("failed to create inner file: %v", err)
	}

	// Remove all permissions from the parent directory
	if err := os.Chmod(restricted, 0o000); err != nil {
		t.Fatalf("failed to change permissions: %v", err)
	}

	t.Cleanup(func() {
		os.Chmod(restricted, 0o755) //nolint:errcheck
	})

	// Should return false without panicking
	if FileExist(innerFile) {
		t.Errorf("FileExist(%q) = true; expected false (parent has no permissions)", innerFile)
	}
}
