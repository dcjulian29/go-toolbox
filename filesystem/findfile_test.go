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

func TestFindFile_InCurrentDirectory(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "target.txt"), []byte("found"), FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	original, err := os.Getwd()
	if err != nil {
		t.Fatalf("setup: failed to get working directory: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(original) })

	if err = os.Chdir(dir); err != nil {
		t.Fatalf("setup: failed to change directory: %v", err)
	}

	result, err := FindFile("target.txt")
	if err != nil {
		t.Fatalf("FindFile() returned error: %v", err)
	}

	if filepath.Base(result) != "target.txt" {
		t.Errorf("FindFile() = %q, want a path ending in 'target.txt'", result)
	}
}

func TestFindFile_InSubdirectory(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "subdir")

	if err := os.Mkdir(sub, 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := os.WriteFile(filepath.Join(sub, "nested.txt"), []byte("found"), FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	original, err := os.Getwd()
	if err != nil {
		t.Fatalf("setup: failed to get working directory: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(original) })

	if err = os.Chdir(dir); err != nil {
		t.Fatalf("setup: failed to change directory: %v", err)
	}

	result, err := FindFile("nested.txt")
	if err != nil {
		t.Fatalf("FindFile() returned error: %v", err)
	}

	if filepath.Base(result) != "nested.txt" {
		t.Errorf("FindFile() = %q, want a path ending in 'nested.txt'", result)
	}
}

func TestFindFile_DeeplyNested(t *testing.T) {
	dir := t.TempDir()
	deep := filepath.Join(dir, "a", "b", "c")

	if err := os.MkdirAll(deep, 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := os.WriteFile(filepath.Join(deep, "deep.txt"), []byte("found"), FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	original, err := os.Getwd()
	if err != nil {
		t.Fatalf("setup: failed to get working directory: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(original) })

	if err = os.Chdir(dir); err != nil {
		t.Fatalf("setup: failed to change directory: %v", err)
	}

	result, err := FindFile("deep.txt")
	if err != nil {
		t.Fatalf("FindFile() returned error: %v", err)
	}

	expected := filepath.Join(deep, "deep.txt")
	if result != expected {
		t.Errorf("FindFile() = %q, want %q", result, expected)
	}
}

func TestFindFile_NotFound(t *testing.T) {
	dir := t.TempDir()

	original, err := os.Getwd()
	if err != nil {
		t.Fatalf("setup: failed to get working directory: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(original) })

	if err = os.Chdir(dir); err != nil {
		t.Fatalf("setup: failed to change directory: %v", err)
	}

	_, err = FindFile("nonexistent.txt")
	if err == nil {
		t.Error("FindFile() should return error for non-existent file")
	}
}

func TestFindFile_ReturnsAbsolutePath(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "abs.txt"), []byte("found"), FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	original, err := os.Getwd()
	if err != nil {
		t.Fatalf("setup: failed to get working directory: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(original) })

	if err = os.Chdir(dir); err != nil {
		t.Fatalf("setup: failed to change directory: %v", err)
	}

	result, err := FindFile("abs.txt")
	if err != nil {
		t.Fatalf("FindFile() returned error: %v", err)
	}

	if !filepath.IsAbs(result) {
		t.Errorf("FindFile() returned relative path %q, want absolute", result)
	}
}

func TestFindFile_SymlinkLoop(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink creation may require elevated privileges on Windows")
	}

	dir := t.TempDir()
	sub := filepath.Join(dir, "subdir")

	if err := os.Mkdir(sub, 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	// Create a symlink inside subdir that points back to the root dir
	if err := os.Symlink(dir, filepath.Join(sub, "loop")); err != nil {
		t.Fatalf("setup: failed to create symlink: %v", err)
	}

	if err := os.WriteFile(filepath.Join(dir, "findme.txt"), []byte("found"), FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	original, err := os.Getwd()
	if err != nil {
		t.Fatalf("setup: failed to get working directory: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(original) })

	if err = os.Chdir(dir); err != nil {
		t.Fatalf("setup: failed to change directory: %v", err)
	}

	// Should find the file without hanging in an infinite loop
	result, err := FindFile("findme.txt")
	if err != nil {
		t.Fatalf("FindFile() returned error: %v", err)
	}

	if filepath.Base(result) != "findme.txt" {
		t.Errorf("FindFile() = %q, want a path ending in 'findme.txt'", result)
	}
}

func TestFindFile_SymlinkToDirectory(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink creation may require elevated privileges on Windows")
	}

	dir := t.TempDir()
	realSub := filepath.Join(dir, "realdir")
	linkSub := filepath.Join(dir, "linkdir")

	if err := os.Mkdir(realSub, 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := os.WriteFile(filepath.Join(realSub, "linked.txt"), []byte("found"), FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := os.Symlink(realSub, linkSub); err != nil {
		t.Fatalf("setup: failed to create symlink: %v", err)
	}

	original, err := os.Getwd()
	if err != nil {
		t.Fatalf("setup: failed to get working directory: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(original) })

	if err = os.Chdir(dir); err != nil {
		t.Fatalf("setup: failed to change directory: %v", err)
	}

	result, err := FindFile("linked.txt")
	if err != nil {
		t.Fatalf("FindFile() returned error: %v", err)
	}

	if filepath.Base(result) != "linked.txt" {
		t.Errorf("FindFile() = %q, want a path ending in 'linked.txt'", result)
	}
}

func TestFindFile_SymlinkToFile(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink creation may require elevated privileges on Windows")
	}

	dir := t.TempDir()

	realFile := filepath.Join(dir, "real.txt")
	if err := os.WriteFile(realFile, []byte("found"), FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := os.Symlink(realFile, filepath.Join(dir, "link.txt")); err != nil {
		t.Fatalf("setup: failed to create symlink: %v", err)
	}

	original, err := os.Getwd()
	if err != nil {
		t.Fatalf("setup: failed to get working directory: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(original) })

	if err = os.Chdir(dir); err != nil {
		t.Fatalf("setup: failed to change directory: %v", err)
	}

	result, err := FindFile("link.txt")
	if err != nil {
		t.Fatalf("FindFile() returned error: %v", err)
	}

	if filepath.Base(result) != "link.txt" {
		t.Errorf("FindFile() = %q, want a path ending in 'link.txt'", result)
	}
}

func TestFindFile_EmptyDirectory(t *testing.T) {
	dir := t.TempDir()

	original, err := os.Getwd()
	if err != nil {
		t.Fatalf("setup: failed to get working directory: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(original) })

	if err = os.Chdir(dir); err != nil {
		t.Fatalf("setup: failed to change directory: %v", err)
	}

	_, err = FindFile("anything.txt")
	if err == nil {
		t.Error("FindFile() should return error when searching empty directory")
	}
}
