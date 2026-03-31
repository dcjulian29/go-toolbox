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
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureFileExist_CreatesFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "testfile.txt")
	content := []byte("hello world")

	if err := EnsureFileExist(path, content); err != nil {
		t.Fatalf("EnsureFileExist() returned error: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read created file: %v", err)
	}

	if !bytes.Equal(got, content) {
		t.Errorf("file content = %q, want %q", got, content)
	}
}

func TestEnsureFileExist_CreatesParentDirectories(t *testing.T) {
	path := filepath.Join(t.TempDir(), "a", "b", "c", "testfile.txt")
	content := []byte("nested content")

	if err := EnsureFileExist(path, content); err != nil {
		t.Fatalf("EnsureFileExist() returned error: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read created file: %v", err)
	}

	if !bytes.Equal(got, content) {
		t.Errorf("file content = %q, want %q", got, content)
	}
}

func TestEnsureFileExist_OverwritesExistingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "testfile.txt")

	if err := os.WriteFile(path, []byte("original content that is longer"), FileModeReadable); err != nil {
		t.Fatalf("setup: failed to create initial file: %v", err)
	}

	newContent := []byte("replaced")

	if err := EnsureFileExist(path, newContent); err != nil {
		t.Fatalf("EnsureFileExist() returned error: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read overwritten file: %v", err)
	}

	if !bytes.Equal(got, newContent) {
		t.Errorf("file content = %q, want %q (should truncate, not append)", got, newContent)
	}
}

func TestEnsureFileExist_EmptyContent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "empty.txt")

	if err := EnsureFileExist(path, []byte{}); err != nil {
		t.Fatalf("EnsureFileExist() returned error: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("failed to stat created file: %v", err)
	}

	if info.Size() != 0 {
		t.Errorf("file size = %d, want 0", info.Size())
	}
}

func TestEnsureFileExist_NilContent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nil.txt")

	if err := EnsureFileExist(path, nil); err != nil {
		t.Fatalf("EnsureFileExist() returned error: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("failed to stat created file: %v", err)
	}

	if info.Size() != 0 {
		t.Errorf("file size = %d, want 0", info.Size())
	}
}

func TestEnsureFileExist_BinaryContent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "binary.dat")
	content := []byte{0x00, 0x01, 0xFF, 0xFE, 0x89, 0x50, 0x4E, 0x47}

	if err := EnsureFileExist(path, content); err != nil {
		t.Fatalf("EnsureFileExist() returned error: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read created file: %v", err)
	}

	if !bytes.Equal(got, content) {
		t.Errorf("file content = %v, want %v", got, content)
	}
}
