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

func TestDirectoryExist_ReturnsTrueForExistingDirectory(t *testing.T) {
	dir := t.TempDir()

	if !DirectoryExist(dir) {
		t.Errorf("DirectoryExist(%q) = false; expected true", dir)
	}
}

func TestDirectoryExist_ReturnsTrueForNestedDirectory(t *testing.T) {
	parent := t.TempDir()
	child := filepath.Join(parent, "subdir", "deep")

	if err := os.MkdirAll(child, 0o755); err != nil {
		t.Fatalf("failed to create nested directory: %v", err)
	}

	if !DirectoryExist(child) {
		t.Errorf("DirectoryExist(%q) = false; expected true", child)
	}
}

func TestDirectoryExist_ReturnsFalseForNonExistentPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nonexistent")

	if DirectoryExist(path) {
		t.Errorf("DirectoryExist(%q) = true; expected false", path)
	}
}

func TestDirectoryExist_ReturnsFalseForRegularFile(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "file.txt")

	if err := os.WriteFile(filePath, []byte("hello"), FileModeReadable); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	if DirectoryExist(filePath) {
		t.Errorf("DirectoryExist(%q) = true; expected false (it is a file, not a directory)", filePath)
	}
}

func TestDirectoryExist_ReturnsFalseForEmptyStringPath(t *testing.T) {
	if DirectoryExist("") {
		t.Error("DirectoryExist(\"\") = true; expected false")
	}
}

func TestDirectoryExist_ReturnsTrueForCurrentDirectoryDot(t *testing.T) {
	if !DirectoryExist(".") {
		t.Error("DirectoryExist(\".\") = false; expected true")
	}
}

func TestDirectoryExist_ReturnsTrueForRootDirectory(t *testing.T) {
	root := "/"
	if runtime.GOOS == "windows" {
		root = `C:\`
	}

	if !DirectoryExist(root) {
		t.Errorf("DirectoryExist(%q) = false; expected true", root)
	}
}

func TestDirectoryExist_ReturnsFalseForSymlinkToFile(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping symlink test on Windows (requires elevated privileges)")
	}

	dir := t.TempDir()
	filePath := filepath.Join(dir, "realfile.txt")
	linkPath := filepath.Join(dir, "link")

	if err := os.WriteFile(filePath, []byte("data"), FileModeReadable); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	if err := os.Symlink(filePath, linkPath); err != nil {
		t.Skipf("symlinks not supported: %v", err)
	}

	if DirectoryExist(linkPath) {
		t.Errorf("DirectoryExist(%q) = true; expected false (symlink points to a file)", linkPath)
	}
}

func TestDirectoryExist_ReturnsTrueForSymlinkToDirectory(t *testing.T) {
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

	if !DirectoryExist(linkPath) {
		t.Errorf("DirectoryExist(%q) = false; expected true (symlink points to a directory)", linkPath)
	}
}

func TestDirectoryExist_DoesNotPanicOnPermissionError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping permission test on Windows (chmod does not restrict directory traversal)")
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

	// Remove all permissions from the parent to make inner inaccessible
	if err := os.Chmod(restricted, 0o000); err != nil {
		t.Fatalf("failed to change permissions: %v", err)
	}

	// Restore permissions on cleanup so t.TempDir() can remove it
	t.Cleanup(func() {
		os.Chmod(restricted, 0o755) //nolint:errcheck
	})

	// This should return false without panicking
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("DirectoryExist panicked on permission error: %v", r)
		}
	}()

	result := DirectoryExist(inner)

	if result {
		t.Errorf("DirectoryExist(%q) = true; expected false (parent has no permissions)", inner)
	}
}
