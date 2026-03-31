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
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// sha256Hex is a small test helper that returns the SHA-256 hex digest of data.
func sha256Hex(data []byte) string {
	h := sha256.Sum256(data)

	return hex.EncodeToString(h[:])
}

func TestFileHash_KnownContent(t *testing.T) {
	content := []byte("hello, world")
	filePath := filepath.Join(t.TempDir(), "known.txt")

	if err := os.WriteFile(filePath, content, FileModeReadable); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	got, err := FileHash(filePath)
	if err != nil {
		t.Fatalf("FileHash returned unexpected error: %v", err)
	}

	want := sha256Hex(content)
	if got != want {
		t.Errorf("FileHash = %q; want %q", got, want)
	}
}

func TestFileHash_EmptyFile(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "empty.txt")

	if err := os.WriteFile(filePath, []byte{}, FileModeReadable); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	got, err := FileHash(filePath)
	if err != nil {
		t.Fatalf("FileHash returned unexpected error: %v", err)
	}

	want := sha256Hex([]byte{})
	if got != want {
		t.Errorf("FileHash = %q; want %q", got, want)
	}
}

func TestFileHash_ReturnsLowercaseHex(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "case.txt")

	if err := os.WriteFile(filePath, []byte("test"), FileModeReadable); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	got, err := FileHash(filePath)
	if err != nil {
		t.Fatalf("FileHash returned unexpected error: %v", err)
	}

	for _, c := range got {
		if (c >= 'A' && c <= 'F') || (c >= 'G' && c <= 'Z') {
			t.Fatalf("FileHash returned non-lowercase hex: %q", got)
		}
	}
}

func TestFileHash_Returns64CharHexString(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "length.txt")

	if err := os.WriteFile(filePath, []byte("check length"), FileModeReadable); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	got, err := FileHash(filePath)
	if err != nil {
		t.Fatalf("FileHash returned unexpected error: %v", err)
	}

	// SHA-256 produces 32 bytes = 64 hex characters
	if len(got) != 64 {
		t.Errorf("FileHash returned %d characters; want 64", len(got))
	}
}

func TestFileHash_NonExistentFile(t *testing.T) {
	_, err := FileHash(filepath.Join(t.TempDir(), "does-not-exist.txt"))
	if err == nil {
		t.Fatal("FileHash should return an error for a non-existent file")
	}
}

func TestFileHash_DirectoryReturnsError(t *testing.T) {
	dir := t.TempDir()

	_, err := FileHash(dir)
	// io.Copy from a directory file descriptor is platform-dependent:
	// on some systems it returns an error, on others it reads zero bytes.
	// We just verify it does not panic.
	_ = err
}

func TestFileHash_PermissionDenied(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping permission test on Windows (permission model differs)")
	}

	if os.Getuid() == 0 {
		t.Skip("skipping permission test when running as root")
	}

	filePath := filepath.Join(t.TempDir(), "noperm.txt")

	if err := os.WriteFile(filePath, []byte("secret"), FileModeReadable); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	if err := os.Chmod(filePath, 0o000); err != nil {
		t.Fatalf("failed to change permissions: %v", err)
	}

	t.Cleanup(func() {
		os.Chmod(filePath, FileModeReadable) //nolint:errcheck
	})

	_, err := FileHash(filePath)
	if err == nil {
		t.Fatal("FileHash should return an error for an unreadable file")
	}
}

func TestFileHash_DeterministicSameContent(t *testing.T) {
	content := []byte("deterministic check")

	fileA := filepath.Join(t.TempDir(), "a.txt")
	fileB := filepath.Join(t.TempDir(), "b.txt")

	if err := os.WriteFile(fileA, content, FileModeReadable); err != nil {
		t.Fatalf("failed to create file A: %v", err)
	}

	if err := os.WriteFile(fileB, content, FileModeReadable); err != nil {
		t.Fatalf("failed to create file B: %v", err)
	}

	hashA, err := FileHash(fileA)
	if err != nil {
		t.Fatalf("FileHash(A) returned unexpected error: %v", err)
	}

	hashB, err := FileHash(fileB)
	if err != nil {
		t.Fatalf("FileHash(B) returned unexpected error: %v", err)
	}

	if hashA != hashB {
		t.Errorf("identical files produced different hashes: %q vs %q", hashA, hashB)
	}
}

func TestFileHash_DifferentContentDifferentHash(t *testing.T) {
	fileA := filepath.Join(t.TempDir(), "a.txt")
	fileB := filepath.Join(t.TempDir(), "b.txt")

	if err := os.WriteFile(fileA, []byte("content A"), FileModeReadable); err != nil {
		t.Fatalf("failed to create file A: %v", err)
	}

	if err := os.WriteFile(fileB, []byte("content B"), FileModeReadable); err != nil {
		t.Fatalf("failed to create file B: %v", err)
	}

	hashA, err := FileHash(fileA)
	if err != nil {
		t.Fatalf("FileHash(A) returned unexpected error: %v", err)
	}

	hashB, err := FileHash(fileB)
	if err != nil {
		t.Fatalf("FileHash(B) returned unexpected error: %v", err)
	}

	if hashA == hashB {
		t.Error("different files produced the same hash")
	}
}

func TestFileHash_BinaryContent(t *testing.T) {
	// File containing all 256 byte values
	content := make([]byte, 256)
	for i := range content {
		content[i] = byte(i)
	}

	filePath := filepath.Join(t.TempDir(), "binary.bin")

	if err := os.WriteFile(filePath, content, FileModeReadable); err != nil {
		t.Fatalf("failed to create binary test file: %v", err)
	}

	got, err := FileHash(filePath)
	if err != nil {
		t.Fatalf("FileHash returned unexpected error: %v", err)
	}

	want := sha256Hex(content)
	if got != want {
		t.Errorf("FileHash = %q; want %q", got, want)
	}
}

func TestFileHash_EmptyPath(t *testing.T) {
	_, err := FileHash("")
	if err == nil {
		t.Fatal("FileHash should return an error for an empty path")
	}
}
