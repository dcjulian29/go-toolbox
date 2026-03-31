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

func TestSearchForFile_ReturnsLastMatch(t *testing.T) {
	dir := t.TempDir()

	for _, name := range []string{"a.txt", "b.txt", "c.txt"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("content"), FileModeReadable); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}

	result, err := SearchForFile(dir, "*.txt")
	if err != nil {
		t.Fatalf("SearchForFile() returned error: %v", err)
	}

	if filepath.Base(result) != "c.txt" {
		t.Errorf("SearchForFile() = %q, want path ending in 'c.txt' (last lexicographically)", result)
	}
}

func TestSearchForFile_SingleMatch(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "only.txt"), []byte("content"), FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	result, err := SearchForFile(dir, "only.txt")
	if err != nil {
		t.Fatalf("SearchForFile() returned error: %v", err)
	}

	if filepath.Base(result) != "only.txt" {
		t.Errorf("SearchForFile() = %q, want path ending in 'only.txt'", result)
	}
}

func TestSearchForFile_NoMatch(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "file.log"), []byte("content"), FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	_, err := SearchForFile(dir, "*.txt")
	if err == nil {
		t.Error("SearchForFile() should return error when no files match")
	}
}

func TestSearchForFile_BadPattern(t *testing.T) {
	dir := t.TempDir()

	_, err := SearchForFile(dir, "[invalid")
	if err == nil {
		t.Error("SearchForFile() should return error for malformed glob pattern")
	}
}

func TestSearchForFile_NonExistentDirectory(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nonexistent")

	_, err := SearchForFile(path, "*.txt")
	if err == nil {
		t.Error("SearchForFile() should return error for non-existent directory")
	}
}

func TestSearchForFile_ExcludesDirectories(t *testing.T) {
	dir := t.TempDir()

	// Create files and a directory that all match *.log
	if err := os.WriteFile(filepath.Join(dir, "a.log"), []byte("content"), FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	// This directory sorts last but should be skipped
	if err := os.Mkdir(filepath.Join(dir, "z.log"), 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	result, err := SearchForFile(dir, "*.log")
	if err != nil {
		t.Fatalf("SearchForFile() returned error: %v", err)
	}

	if filepath.Base(result) != "a.log" {
		t.Errorf("SearchForFile() = %q, want path ending in 'a.log' (directory 'z.log' should be excluded)", result)
	}
}

func TestSearchForFile_SkipsDirectoryReturnsNextFile(t *testing.T) {
	dir := t.TempDir()

	for _, name := range []string{"a.txt", "b.txt"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("content"), FileModeReadable); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}

	// Directory that sorts between b.txt and z
	if err := os.Mkdir(filepath.Join(dir, "c.txt"), 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	result, err := SearchForFile(dir, "*.txt")
	if err != nil {
		t.Fatalf("SearchForFile() returned error: %v", err)
	}

	if filepath.Base(result) != "b.txt" {
		t.Errorf("SearchForFile() = %q, want path ending in 'b.txt' (directory 'c.txt' should be skipped)", result)
	}
}

func TestSearchForFile_OnlyDirectoriesMatch(t *testing.T) {
	dir := t.TempDir()

	if err := os.Mkdir(filepath.Join(dir, "subdir.txt"), 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	_, err := SearchForFile(dir, "*.txt")
	if err == nil {
		t.Error("SearchForFile() should return error when only directories match the pattern")
	}
}

func TestSearchForFile_EmptyDirectory(t *testing.T) {
	dir := t.TempDir()

	_, err := SearchForFile(dir, "*")
	if err == nil {
		t.Error("SearchForFile() should return error for empty directory")
	}
}

func TestSearchForFile_SymlinkToFileIncluded(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink creation may require elevated privileges on Windows")
	}

	dir := t.TempDir()

	realFile := filepath.Join(dir, "a.txt")
	if err := os.WriteFile(realFile, []byte("content"), FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	// Symlink that sorts after the real file
	if err := os.Symlink(realFile, filepath.Join(dir, "z.txt")); err != nil {
		t.Fatalf("setup: failed to create symlink: %v", err)
	}

	result, err := SearchForFile(dir, "*.txt")
	if err != nil {
		t.Fatalf("SearchForFile() returned error: %v", err)
	}

	if filepath.Base(result) != "z.txt" {
		t.Errorf("SearchForFile() = %q, want path ending in 'z.txt' (symlink should be included as last match)", result)
	}
}

func TestSearchForFile_NonRecursive(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "subdir")

	if err := os.Mkdir(sub, 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := os.WriteFile(filepath.Join(sub, "nested.txt"), []byte("content"), FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	_, err := SearchForFile(dir, "nested.txt")
	if err == nil {
		t.Error("SearchForFile() should not find files in subdirectories")
	}
}
