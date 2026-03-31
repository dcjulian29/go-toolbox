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

func TestEnsureDirectoryExist_CreatesNewDirectory(t *testing.T) {
	dirPath := filepath.Join(t.TempDir(), "newdir")

	err := EnsureDirectoryExist(dirPath)
	if err != nil {
		t.Fatalf("EnsureDirectoryExist returned unexpected error: %v", err)
	}

	if !DirectoryExist(dirPath) {
		t.Errorf("directory was not created at %q", dirPath)
	}
}

func TestEnsureDirectoryExist_CreatesNestedDirectories(t *testing.T) {
	dirPath := filepath.Join(t.TempDir(), "a", "b", "c", "d")

	err := EnsureDirectoryExist(dirPath)
	if err != nil {
		t.Fatalf("EnsureDirectoryExist returned unexpected error: %v", err)
	}

	if !DirectoryExist(dirPath) {
		t.Errorf("nested directories were not created at %q", dirPath)
	}
}

func TestEnsureDirectoryExist_SucceedsWhenAlreadyExists(t *testing.T) {
	dirPath := t.TempDir()

	err := EnsureDirectoryExist(dirPath)
	if err != nil {
		t.Fatalf("EnsureDirectoryExist should succeed for existing directory, got: %v", err)
	}

	if !DirectoryExist(dirPath) {
		t.Error("existing directory should still exist")
	}
}

func TestEnsureDirectoryExist_IsIdempotent(t *testing.T) {
	dirPath := filepath.Join(t.TempDir(), "idempotent")

	for i := range 3 {
		if err := EnsureDirectoryExist(dirPath); err != nil {
			t.Fatalf("call %d: EnsureDirectoryExist returned unexpected error: %v", i+1, err)
		}
	}

	if !DirectoryExist(dirPath) {
		t.Errorf("directory should exist at %q after multiple calls", dirPath)
	}
}

func TestEnsureDirectoryExist_PreservesExistingContents(t *testing.T) {
	dirPath := filepath.Join(t.TempDir(), "preserve")

	if err := os.Mkdir(dirPath, 0o755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}

	testFile := filepath.Join(dirPath, "existing.txt")

	if err := os.WriteFile(testFile, []byte("keep me"), FileModeReadable); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	if err := EnsureDirectoryExist(dirPath); err != nil {
		t.Fatalf("EnsureDirectoryExist returned unexpected error: %v", err)
	}

	if !FileExist(testFile) {
		t.Error("existing file inside directory was destroyed")
	}

	data, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	if string(data) != "keep me" {
		t.Errorf("file content = %q; expected %q", string(data), "keep me")
	}
}

func TestEnsureDirectoryExist_ReturnsErrorWhenPathIsFile(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "file.txt")

	if err := os.WriteFile(filePath, []byte("I'm a file"), FileModeReadable); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	err := EnsureDirectoryExist(filePath)
	if err == nil {
		t.Fatal("EnsureDirectoryExist should return an error when path is a file")
	}
}

func TestEnsureDirectoryExist_EmptyPath(t *testing.T) {
	// os.Stat("") returns an error, and os.IsNotExist may be true
	// for empty path, so os.MkdirAll("", ...) may be called.
	// os.MkdirAll("") returns nil without creating anything.
	// This test documents the current behavior.

	err := EnsureDirectoryExist("")

	if err != nil {
		t.Logf("EnsureDirectoryExist returned error for empty path: %v", err)
	} else {
		t.Log("EnsureDirectoryExist returned nil for empty path (os.MkdirAll accepts empty string)")
	}
}

func TestEnsureDirectoryExist_ReturnsErrorWhenParentNotWritable(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping permission test on Windows (permission model differs)")
	}

	if os.Getuid() == 0 {
		t.Skip("skipping permission test when running as root")
	}

	dir := t.TempDir()
	locked := filepath.Join(dir, "locked")

	if err := os.Mkdir(locked, 0o755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}

	if err := os.Chmod(locked, 0o555); err != nil {
		t.Fatalf("failed to change permissions: %v", err)
	}

	t.Cleanup(func() {
		os.Chmod(locked, 0o755) //nolint:errcheck
	})

	target := filepath.Join(locked, "subdir")

	err := EnsureDirectoryExist(target)
	if err == nil {
		t.Fatal("EnsureDirectoryExist should have returned an error for non-writable parent")
	}
}

func TestEnsureDirectoryExist_CreatedDirectoryPassesChecks(t *testing.T) {
	dirPath := filepath.Join(t.TempDir(), "verify")

	if err := EnsureDirectoryExist(dirPath); err != nil {
		t.Fatalf("EnsureDirectoryExist returned unexpected error: %v", err)
	}

	if !DirectoryExist(dirPath) {
		t.Error("DirectoryExist returned false for a directory just created by EnsureDirectoryExist")
	}

	if FileExist(dirPath) {
		t.Error("FileExist returned true for a directory (should only be true for files)")
	}
}

func TestEnsureDirectoryExist_ReturnsErrorOnStatPermissionError(t *testing.T) {
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

	innerDir := filepath.Join(restricted, "inner")

	if err := os.Mkdir(innerDir, 0o755); err != nil {
		t.Fatalf("failed to create inner directory: %v", err)
	}

	if err := os.Chmod(restricted, 0o000); err != nil {
		t.Fatalf("failed to change permissions: %v", err)
	}

	t.Cleanup(func() {
		os.Chmod(restricted, 0o755) //nolint:errcheck
	})

	err := EnsureDirectoryExist(innerDir)
	if err == nil {
		t.Fatal("EnsureDirectoryExist should return an error when stat fails with permission denied")
	}
}
