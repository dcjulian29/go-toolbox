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
	"path/filepath"
	"strings"
	"testing"
)

func Test_FindLatestISO_EmptyDirectoryPath(t *testing.T) {
	_, err := FindLatestISO("", "Win11")
	if err == nil {
		t.Fatal("expected error for empty directory path")
	}

	if !strings.Contains(err.Error(), "directory path") {
		t.Errorf("error should mention directory path, got: %v", err)
	}
}

func Test_FindLatestISO_EmptyPattern(t *testing.T) {
	_, err := FindLatestISO(`C:\isos`, "")
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}

	if !strings.Contains(err.Error(), "pattern") {
		t.Errorf("error should mention pattern, got: %v", err)
	}
}

func Test_FindLatestISO_NoMatches(t *testing.T) {
	dir := t.TempDir()

	_, err := FindLatestISO(dir, "Win11")
	if err == nil {
		t.Fatal("expected error when no files match")
	}

	if !strings.Contains(err.Error(), "no ISO") {
		t.Errorf("error should mention no ISO, got: %v", err)
	}
}

func Test_FindLatestISO_IgnoresNonISOFiles(t *testing.T) {
	dir := t.TempDir()
	createTempFiles(t, dir, []string{
		"Win11-2025-01.vhdx",
		"Win11-2025-06.txt",
	})

	_, err := FindLatestISO(dir, "Win11")
	if err == nil {
		t.Fatal("expected error when no ISO files match")
	}
}

func Test_FindLatestISO_SingleMatch(t *testing.T) {
	dir := t.TempDir()
	createTempFiles(t, dir, []string{"Win11-2025-01.iso"})

	got, err := FindLatestISO(dir, "Win11")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := filepath.Join(dir, "Win11-2025-01.iso")
	if got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

func Test_FindLatestISO_ReturnsAlphabeticallyLast(t *testing.T) {
	dir := t.TempDir()
	createTempFiles(t, dir, []string{
		"Win11-2025-01.iso",
		"Win11-2025-06.iso",
		"Win11-2025-03.iso",
	})

	got, err := FindLatestISO(dir, "Win11")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := filepath.Join(dir, "Win11-2025-06.iso")
	if got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

func Test_FindLatestISO_PrefixFiltersCorrectly(t *testing.T) {
	dir := t.TempDir()
	createTempFiles(t, dir, []string{
		"Win11-2025-01.iso",
		"Win10-2025-06.iso",
		"Win11-2025-03.iso",
	})

	got, err := FindLatestISO(dir, "Win11")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := filepath.Join(dir, "Win11-2025-03.iso")
	if got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}
