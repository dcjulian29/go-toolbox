//go:build windows

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

package hypervhost

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dcjulian29/go-toolbox/filesystem"
)

func createTempFiles(t *testing.T, dir string, names []string) {
	t.Helper()

	for _, name := range names {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte{}, filesystem.ModeOwnerReadWrite); err != nil {
			t.Fatalf("failed to create temp file %s: %v", name, err)
		}
	}
}

func Test_FindLatestBaseDisk_EmptyDirectoryPath(t *testing.T) {
	_, err := FindLatestBaseDisk("", "base")
	if err == nil {
		t.Fatal("expected error for empty directory path")
	}

	if !strings.Contains(err.Error(), "directory path") {
		t.Errorf("error should mention directory path, got: %v", err)
	}
}

func Test_FindLatestBaseDisk_EmptyPattern(t *testing.T) {
	_, err := FindLatestBaseDisk(`C:\disks`, "")
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}

	if !strings.Contains(err.Error(), "pattern") {
		t.Errorf("error should mention pattern, got: %v", err)
	}
}

func Test_FindLatestBaseDisk_NoMatches(t *testing.T) {
	dir := t.TempDir()

	_, err := FindLatestBaseDisk(dir, "base")
	if err == nil {
		t.Fatal("expected error when no files match")
	}

	if !strings.Contains(err.Error(), "no base disk") {
		t.Errorf("error should mention no base disk, got: %v", err)
	}
}

func Test_FindLatestBaseDisk_IgnoresNonVHDXFiles(t *testing.T) {
	dir := t.TempDir()
	createTempFiles(t, dir, []string{
		"base-2025-01.iso",
		"base-2025-06.txt",
	})

	_, err := FindLatestBaseDisk(dir, "base")
	if err == nil {
		t.Fatal("expected error when no VHDX files match")
	}
}

func Test_FindLatestBaseDisk_SingleMatch(t *testing.T) {
	dir := t.TempDir()
	createTempFiles(t, dir, []string{"base-2025-01.vhdx"})

	got, err := FindLatestBaseDisk(dir, "base")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := filepath.Join(dir, "base-2025-01.vhdx")
	if got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

func Test_FindLatestBaseDisk_ReturnsAlphabeticallyLast(t *testing.T) {
	dir := t.TempDir()
	createTempFiles(t, dir, []string{
		"base-2025-01.vhdx",
		"base-2025-06.vhdx",
		"base-2025-03.vhdx",
	})

	got, err := FindLatestBaseDisk(dir, "base")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := filepath.Join(dir, "base-2025-06.vhdx")
	if got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

func Test_FindLatestBaseDisk_PrefixFiltersCorrectly(t *testing.T) {
	dir := t.TempDir()
	createTempFiles(t, dir, []string{
		"server-2025-01.vhdx",
		"desktop-2025-06.vhdx",
		"server-2025-03.vhdx",
	})

	got, err := FindLatestBaseDisk(dir, "server")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := filepath.Join(dir, "server-2025-03.vhdx")
	if got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}
