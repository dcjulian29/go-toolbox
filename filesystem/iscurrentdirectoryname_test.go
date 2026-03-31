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

func TestIsCurrentDirectoryName_Match(t *testing.T) {
	dir := t.TempDir()
	name := filepath.Base(dir)

	original, err := os.Getwd()
	if err != nil {
		t.Fatalf("setup: failed to get working directory: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(original) })

	if err := os.Chdir(dir); err != nil {
		t.Fatalf("setup: failed to change directory: %v", err)
	}

	if !IsCurrentDirectoryName(name) {
		t.Errorf("IsCurrentDirectoryName(%q) = false, want true", name)
	}
}

func TestIsCurrentDirectoryName_NoMatch(t *testing.T) {
	dir := t.TempDir()

	original, err := os.Getwd()
	if err != nil {
		t.Fatalf("setup: failed to get working directory: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(original) })

	if err := os.Chdir(dir); err != nil {
		t.Fatalf("setup: failed to change directory: %v", err)
	}

	if IsCurrentDirectoryName("this-should-not-match") {
		t.Error("IsCurrentDirectoryName() = true for non-matching name, want false")
	}
}

func TestIsCurrentDirectoryName_EmptyString(t *testing.T) {
	if IsCurrentDirectoryName("") {
		t.Error("IsCurrentDirectoryName(\"\") = true, want false")
	}
}

func TestIsCurrentDirectoryName_CaseSensitivity(t *testing.T) {
	dir := t.TempDir()
	name := filepath.Base(dir)

	original, err := os.Getwd()
	if err != nil {
		t.Fatalf("setup: failed to get working directory: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(original) })

	if err := os.Chdir(dir); err != nil {
		t.Fatalf("setup: failed to change directory: %v", err)
	}

	mixedCase := swapCase(name)
	if mixedCase == name {
		t.Skip("directory name has no alphabetic characters to test case sensitivity")
	}

	result := IsCurrentDirectoryName(mixedCase)

	switch runtime.GOOS {
	case "windows", "darwin":
		if !result {
			t.Errorf("IsCurrentDirectoryName(%q) = false on %s, want true (case-insensitive)", mixedCase, runtime.GOOS)
		}
	default:
		if result {
			t.Errorf("IsCurrentDirectoryName(%q) = true on %s, want false (case-sensitive)", mixedCase, runtime.GOOS)
		}
	}
}

func TestIsCurrentDirectoryName_FullPathDoesNotMatch(t *testing.T) {
	dir := t.TempDir()

	original, err := os.Getwd()
	if err != nil {
		t.Fatalf("setup: failed to get working directory: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(original) })

	if err := os.Chdir(dir); err != nil {
		t.Fatalf("setup: failed to change directory: %v", err)
	}

	if IsCurrentDirectoryName(dir) {
		t.Errorf("IsCurrentDirectoryName(%q) = true for full path, want false", dir)
	}
}

// swapCase inverts the case of each ASCII letter in the string.
func swapCase(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			b[i] = c + 32
		} else if c >= 'a' && c <= 'z' {
			b[i] = c - 32
		}
	}

	return string(b)
}
