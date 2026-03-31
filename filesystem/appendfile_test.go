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

func TestAppendFile_BasicAppend(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")

	if err := os.WriteFile(path, []byte("hello"), ModeOwnerReadWrite); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := AppendFile(path, []byte(" world")); err != nil {
		t.Fatalf("AppendFile() returned error: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading file: %v", err)
	}

	if string(got) != "hello world" {
		t.Errorf("AppendFile() file content = %q, want %q", got, "hello world")
	}
}

func TestAppendFile_MultipleAppends(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")

	if err := os.WriteFile(path, []byte("a"), ModeOwnerReadWrite); err != nil {
		t.Fatalf("setup: %v", err)
	}

	for _, s := range []string{"b", "c", "d"} {
		if err := AppendFile(path, []byte(s)); err != nil {
			t.Fatalf("AppendFile() returned error: %v", err)
		}
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading file: %v", err)
	}

	if string(got) != "abcd" {
		t.Errorf("AppendFile() file content = %q, want %q", got, "abcd")
	}
}

func TestAppendFile_EmptyContent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")
	original := []byte("original")

	if err := os.WriteFile(path, original, ModeOwnerReadWrite); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := AppendFile(path, []byte{}); err != nil {
		t.Fatalf("AppendFile() returned error: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading file: %v", err)
	}

	if string(got) != string(original) {
		t.Errorf("AppendFile() file content = %q, want %q (empty append should be no-op)", got, original)
	}
}

func TestAppendFile_NilContent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")
	original := []byte("original")

	if err := os.WriteFile(path, original, ModeOwnerReadWrite); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := AppendFile(path, nil); err != nil {
		t.Fatalf("AppendFile() returned error: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading file: %v", err)
	}

	if string(got) != string(original) {
		t.Errorf("AppendFile() file content = %q, want %q (nil append should be no-op)", got, original)
	}
}

func TestAppendFile_BinaryContent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "binary.dat")

	if err := os.WriteFile(path, []byte{}, ModeOwnerReadWrite); err != nil {
		t.Fatalf("setup: %v", err)
	}

	content := make([]byte, 256)
	for i := range content {
		content[i] = byte(i)
	}

	if err := AppendFile(path, content); err != nil {
		t.Fatalf("AppendFile() returned error: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading file: %v", err)
	}

	if !bytes.Equal(got, content) {
		t.Errorf("AppendFile() binary content mismatch: got %d bytes, want %d bytes", len(got), len(content))
	}
}

func TestAppendFile_FileDoesNotExist(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nonexistent.txt")

	err := AppendFile(path, []byte("content"))
	if err == nil {
		t.Fatal("AppendFile() returned nil, want error for nonexistent file")
	}
}

func TestAppendFile_DirectoryDoesNotExist(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nonexistent", "file.txt")

	err := AppendFile(path, []byte("content"))
	if err == nil {
		t.Fatal("AppendFile() returned nil, want error for nonexistent parent directory")
	}
}

func TestAppendFile_PreservesExistingContent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")
	original := "line 1\nline 2\n"

	if err := os.WriteFile(path, []byte(original), ModeOwnerReadWrite); err != nil {
		t.Fatalf("setup: %v", err)
	}

	appended := "line 3\n"

	if err := AppendFile(path, []byte(appended)); err != nil {
		t.Fatalf("AppendFile() returned error: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading file: %v", err)
	}

	want := original + appended
	if string(got) != want {
		t.Errorf("AppendFile() file content = %q, want %q", got, want)
	}
}

func TestAppendFile_AppendToEmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.txt")

	if err := os.WriteFile(path, []byte{}, ModeOwnerReadWrite); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := AppendFile(path, []byte("content")); err != nil {
		t.Fatalf("AppendFile() returned error: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading file: %v", err)
	}

	if string(got) != "content" {
		t.Errorf("AppendFile() file content = %q, want %q", got, "content")
	}
}

func TestAppendFile_LargeAppend(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "large.dat")

	if err := os.WriteFile(path, []byte{}, ModeOwnerReadWrite); err != nil {
		t.Fatalf("setup: %v", err)
	}

	content := bytes.Repeat([]byte("x"), 1024*1024) // 1 MB

	if err := AppendFile(path, content); err != nil {
		t.Fatalf("AppendFile() returned error: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat file: %v", err)
	}

	if info.Size() != int64(len(content)) {
		t.Errorf("AppendFile() file size = %d, want %d", info.Size(), len(content))
	}
}
