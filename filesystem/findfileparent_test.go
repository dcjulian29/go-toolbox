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

func TestFindFileParent_FindsInParent(t *testing.T) {
	root := t.TempDir()
	child := filepath.Join(root, "child")

	if err := os.Mkdir(child, 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	target := filepath.Join(root, "target.txt")
	if err := os.WriteFile(target, []byte("found"), FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	orig, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(orig) })

	if err := os.Chdir(child); err != nil {
		t.Fatalf("setup: failed to chdir: %v", err)
	}

	result, err := FindFileParent("target.txt")
	if err != nil {
		t.Fatalf("FindFileParent() returned error: %v", err)
	}

	if result != target {
		t.Errorf("FindFileParent() = %q, want %q", result, target)
	}
}

func TestFindFileParent_FindsInGrandparent(t *testing.T) {
	root := t.TempDir()
	deep := filepath.Join(root, "a", "b")

	if err := os.MkdirAll(deep, 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	target := filepath.Join(root, "config.yml")
	if err := os.WriteFile(target, []byte("found"), FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	orig, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(orig) })

	if err := os.Chdir(deep); err != nil {
		t.Fatalf("setup: failed to chdir: %v", err)
	}

	result, err := FindFileParent("config.yml")
	if err != nil {
		t.Fatalf("FindFileParent() returned error: %v", err)
	}

	if result != target {
		t.Errorf("FindFileParent() = %q, want %q", result, target)
	}
}

func TestFindFileParent_SkipsCurrentDirectory(t *testing.T) {
	root := t.TempDir()
	child := filepath.Join(root, "child")

	if err := os.Mkdir(child, 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	// File only in the current directory, NOT in any parent
	if err := os.WriteFile(filepath.Join(child, "local.txt"), []byte("local"), FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	orig, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(orig) })

	if err := os.Chdir(child); err != nil {
		t.Fatalf("setup: failed to chdir: %v", err)
	}

	_, err := FindFileParent("local.txt")
	if err == nil {
		t.Error("FindFileParent() should not find files in the current directory itself")
	}
}

func TestFindFileParent_NotFound(t *testing.T) {
	root := t.TempDir()

	orig, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(orig) })

	if err := os.Chdir(root); err != nil {
		t.Fatalf("setup: failed to chdir: %v", err)
	}

	_, err := FindFileParent("nonexistent.txt")
	if err == nil {
		t.Error("FindFileParent() should return error when file is not found in any parent")
	}
}

func TestFindFileParent_ReturnsClosestParent(t *testing.T) {
	root := t.TempDir()
	middle := filepath.Join(root, "middle")
	deep := filepath.Join(middle, "deep")

	if err := os.MkdirAll(deep, 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	// Same filename in both root and middle
	for _, dir := range []string{root, middle} {
		if err := os.WriteFile(filepath.Join(dir, "Makefile"), []byte(dir), FileModeReadable); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}

	orig, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(orig) })

	if err := os.Chdir(deep); err != nil {
		t.Fatalf("setup: failed to chdir: %v", err)
	}

	result, err := FindFileParent("Makefile")
	if err != nil {
		t.Fatalf("FindFileParent() returned error: %v", err)
	}

	expected := filepath.Join(middle, "Makefile")
	if result != expected {
		t.Errorf("FindFileParent() = %q, want %q (should return closest parent match)", result, expected)
	}
}

func TestFindFileParent_ExcludesDirectories(t *testing.T) {
	root := t.TempDir()
	child := filepath.Join(root, "child")

	if err := os.Mkdir(child, 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	// Create a directory with the target name in the parent
	if err := os.Mkdir(filepath.Join(root, "target"), 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	orig, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(orig) })

	if err := os.Chdir(child); err != nil {
		t.Fatalf("setup: failed to chdir: %v", err)
	}

	_, err := FindFileParent("target")
	if err == nil {
		t.Error("FindFileParent() should not match directories")
	}
}

func TestFindFileParent_SkipsDirectoryFindsFileAbove(t *testing.T) {
	root := t.TempDir()
	middle := filepath.Join(root, "middle")
	deep := filepath.Join(middle, "deep")

	if err := os.MkdirAll(deep, 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	// Directory named "build" in middle (should be skipped)
	if err := os.Mkdir(filepath.Join(middle, "build"), 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	// File named "build" in root (should be found)
	expected := filepath.Join(root, "build")
	if err := os.WriteFile(expected, []byte("file"), FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	orig, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(orig) })

	if err := os.Chdir(deep); err != nil {
		t.Fatalf("setup: failed to chdir: %v", err)
	}

	result, err := FindFileParent("build")
	if err != nil {
		t.Fatalf("FindFileParent() returned error: %v", err)
	}

	if result != expected {
		t.Errorf("FindFileParent() = %q, want %q (should skip directory and find file above)", result, expected)
	}
}

func TestFindFileParent_ReturnsAbsolutePath(t *testing.T) {
	root := t.TempDir()
	child := filepath.Join(root, "child")

	if err := os.Mkdir(child, 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := os.WriteFile(filepath.Join(root, "abs.txt"), []byte("content"), FileModeReadable); err != nil {
		t.Fatalf("setup: %v", err)
	}

	orig, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(orig) })

	if err := os.Chdir(child); err != nil {
		t.Fatalf("setup: failed to chdir: %v", err)
	}

	result, err := FindFileParent("abs.txt")
	if err != nil {
		t.Fatalf("FindFileParent() returned error: %v", err)
	}

	if !filepath.IsAbs(result) {
		t.Errorf("FindFileParent() returned relative path %q, want absolute", result)
	}
}
