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

func TestCopyFile_BasicCopy(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "source.txt")
	dst := filepath.Join(dir, "destination.txt")
	content := []byte("hello world")

	if err := os.WriteFile(src, content, FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := CopyFile(src, dst); err != nil {
		t.Fatalf("CopyFile() returned error: %v", err)
	}

	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("reading destination: %v", err)
	}

	if string(got) != string(content) {
		t.Errorf("CopyFile() destination content = %q, want %q", got, content)
	}
}

func TestCopyFile_PreservesPermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("file permission preservation not reliable on Windows")
	}

	dir := t.TempDir()
	src := filepath.Join(dir, "source.sh")
	dst := filepath.Join(dir, "destination.sh")

	if err := os.WriteFile(src, []byte("#!/bin/sh"), FileModeExecutable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := CopyFile(src, dst); err != nil {
		t.Fatalf("CopyFile() returned error: %v", err)
	}

	srcInfo, err := os.Stat(src)
	if err != nil {
		t.Fatalf("stat source: %v", err)
	}

	dstInfo, err := os.Stat(dst)
	if err != nil {
		t.Fatalf("stat destination: %v", err)
	}

	if srcInfo.Mode().Perm() != dstInfo.Mode().Perm() {
		t.Errorf("CopyFile() destination permissions = %v, want %v", dstInfo.Mode().Perm(), srcInfo.Mode().Perm())
	}
}

func TestCopyFile_OverwritesExistingDestination(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "source.txt")
	dst := filepath.Join(dir, "destination.txt")

	if err := os.WriteFile(src, []byte("new content"), FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := os.WriteFile(dst, []byte("old content that is longer"), FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := CopyFile(src, dst); err != nil {
		t.Fatalf("CopyFile() returned error: %v", err)
	}

	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("reading destination: %v", err)
	}

	if string(got) != "new content" {
		t.Errorf("CopyFile() destination content = %q, want %q (should fully overwrite)", got, "new content")
	}
}

func TestCopyFile_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "empty.txt")
	dst := filepath.Join(dir, "copy.txt")

	if err := os.WriteFile(src, []byte{}, FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := CopyFile(src, dst); err != nil {
		t.Fatalf("CopyFile() returned error: %v", err)
	}

	info, err := os.Stat(dst)
	if err != nil {
		t.Fatalf("stat destination: %v", err)
	}

	if info.Size() != 0 {
		t.Errorf("CopyFile() destination size = %d, want 0", info.Size())
	}
}

func TestCopyFile_BinaryContent(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "binary.dat")
	dst := filepath.Join(dir, "copy.dat")

	content := make([]byte, 256)
	for i := range content {
		content[i] = byte(i)
	}

	if err := os.WriteFile(src, content, FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := CopyFile(src, dst); err != nil {
		t.Fatalf("CopyFile() returned error: %v", err)
	}

	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("reading destination: %v", err)
	}

	if len(got) != len(content) {
		t.Fatalf("CopyFile() destination length = %d, want %d", len(got), len(content))
	}

	for i := range content {
		if got[i] != content[i] {
			t.Errorf("CopyFile() byte mismatch at index %d: got %d, want %d", i, got[i], content[i])

			break
		}
	}
}

func TestCopyFile_SourceDoesNotExist(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "nonexistent.txt")
	dst := filepath.Join(dir, "destination.txt")

	err := CopyFile(src, dst)
	if err == nil {
		t.Fatal("CopyFile() returned nil, want error for nonexistent source")
	}
}

func TestCopyFile_SourceIsDirectory(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "subdir")
	dst := filepath.Join(dir, "destination.txt")

	if err := os.Mkdir(src, 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	err := CopyFile(src, dst)
	if err == nil {
		t.Fatal("CopyFile() returned nil, want error when source is a directory")
	}
}

func TestCopyFile_DestinationDirectoryDoesNotExist(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "source.txt")
	dst := filepath.Join(dir, "nonexistent", "destination.txt")

	if err := os.WriteFile(src, []byte("content"), FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	err := CopyFile(src, dst)
	if err == nil {
		t.Fatal("CopyFile() returned nil, want error when destination directory does not exist")
	}
}

func TestCopyFile_SourceNotModified(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "source.txt")
	dst := filepath.Join(dir, "destination.txt")
	content := []byte("original content")

	if err := os.WriteFile(src, content, FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := CopyFile(src, dst); err != nil {
		t.Fatalf("CopyFile() returned error: %v", err)
	}

	got, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("reading source: %v", err)
	}

	if string(got) != string(content) {
		t.Errorf("CopyFile() modified source: got %q, want %q", got, content)
	}
}

func TestCopyFile_DestinationIsIndependentCopy(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "source.txt")
	dst := filepath.Join(dir, "destination.txt")

	if err := os.WriteFile(src, []byte("original"), FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := CopyFile(src, dst); err != nil {
		t.Fatalf("CopyFile() returned error: %v", err)
	}

	// Modify the source after copying
	if err := os.WriteFile(src, []byte("modified"), FileModeReadable); err != nil {
		t.Fatalf("modifying source: %v", err)
	}

	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("reading destination: %v", err)
	}

	if string(got) != "original" {
		t.Errorf("CopyFile() destination content = %q, want %q (should be independent copy)", got, "original")
	}
}

func TestCopyFile_SameDirectory(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "file.txt")
	dst := filepath.Join(dir, "file_copy.txt")

	if err := os.WriteFile(src, []byte("content"), FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := CopyFile(src, dst); err != nil {
		t.Fatalf("CopyFile() returned error: %v", err)
	}

	// Both files should exist
	for _, path := range []string{src, dst} {
		if _, err := os.Stat(path); err != nil {
			t.Errorf("CopyFile() file %q does not exist after copy", path)
		}
	}
}

func TestCopyFile_DifferentDirectory(t *testing.T) {
	dir := t.TempDir()
	srcDir := filepath.Join(dir, "src")
	dstDir := filepath.Join(dir, "dst")

	if err := os.Mkdir(srcDir, 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := os.Mkdir(dstDir, 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	src := filepath.Join(srcDir, "file.txt")
	dst := filepath.Join(dstDir, "file.txt")

	if err := os.WriteFile(src, []byte("content"), FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := CopyFile(src, dst); err != nil {
		t.Fatalf("CopyFile() returned error: %v", err)
	}

	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("reading destination: %v", err)
	}

	if string(got) != "content" {
		t.Errorf("CopyFile() destination content = %q, want %q", got, "content")
	}
}
