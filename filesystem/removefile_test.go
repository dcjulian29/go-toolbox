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
	"testing"
)

func TestRemoveFile_ExistingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "testfile.txt")

	if err := os.WriteFile(path, []byte("hello"), FileModeReadable); err != nil {
		t.Fatalf("setup: failed to create test file: %v", err)
	}

	err := RemoveFile(path)
	if err != nil {
		t.Errorf("RemoveFile() returned error: %v", err)
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("file still exists after RemoveFile()")
	}
}

func TestRemoveFile_NonExistentFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "doesnotexist.txt")

	err := RemoveFile(path)
	if err != nil {
		t.Errorf("RemoveFile() on non-existent file returned error: %v", err)
	}
}

func TestRemoveFile_NonExistentParentDirectory(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nodir", "nofile.txt")

	err := RemoveFile(path)
	if err != nil {
		t.Errorf("RemoveFile() with non-existent parent directory returned error: %v", err)
	}
}

func TestRemoveFile_NonEmptyDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "subdir")

	if err := os.Mkdir(dir, 0755); err != nil {
		t.Fatalf("setup: failed to create directory: %v", err)
	}

	if err := os.WriteFile(filepath.Join(dir, "file.txt"), []byte("hello"), FileModeReadable); err != nil {
		t.Fatalf("setup: failed to create file in directory: %v", err)
	}

	err := RemoveFile(dir)
	if err == nil {
		t.Error("RemoveFile() on non-empty directory should return error")
	}
}

func TestRemoveFile_EmptyDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "emptydir")

	if err := os.Mkdir(dir, 0755); err != nil {
		t.Fatalf("setup: failed to create directory: %v", err)
	}

	err := RemoveFile(dir)
	if err != nil {
		t.Errorf("RemoveFile() on empty directory returned error: %v", err)
	}
}

func TestRemoveFile_Idempotent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "testfile.txt")

	if err := os.WriteFile(path, []byte("hello"), FileModeReadable); err != nil {
		t.Fatalf("setup: failed to create test file: %v", err)
	}

	if err := RemoveFile(path); err != nil {
		t.Fatalf("first RemoveFile() returned error: %v", err)
	}

	if err := RemoveFile(path); err != nil {
		t.Errorf("second RemoveFile() returned error: %v", err)
	}
}
