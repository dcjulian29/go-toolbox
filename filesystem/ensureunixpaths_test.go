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
	"testing"
)

func TestEnsureUnixPaths_ConvertsBackslashes(t *testing.T) {
	input := []string{`C:\Users\test\file.txt`, `D:\data\output`}

	result := EnsureUnixPaths(input)

	if len(result) != 2 {
		t.Fatalf("got %d results; want 2", len(result))
	}

	if result[0] != "/C/Users/test/file.txt" {
		t.Errorf("result[0] = %q; want %q", result[0], "/C/Users/test/file.txt")
	}

	if result[1] != "/D/data/output" {
		t.Errorf("result[1] = %q; want %q", result[1], "/D/data/output")
	}
}

func TestEnsureUnixPaths_EmptySlice(t *testing.T) {
	result := EnsureUnixPaths([]string{})

	if len(result) != 0 {
		t.Errorf("got %d results; want 0", len(result))
	}
}

func TestEnsureUnixPaths_UnixPathsUnchanged(t *testing.T) {
	input := []string{"/usr/local/bin", "./relative/path"}

	result := EnsureUnixPaths(input)

	if result[0] != "/usr/local/bin" {
		t.Errorf("result[0] = %q; want %q", result[0], "/usr/local/bin")
	}

	if result[1] != "./relative/path" {
		t.Errorf("result[1] = %q; want %q", result[1], "./relative/path")
	}
}

func TestEnsureUnixPaths_MixedSeparators(t *testing.T) {
	input := []string{`C:\mixed/path\to/file`}

	result := EnsureUnixPaths(input)

	if result[0] != "/C/mixed/path/to/file" {
		t.Errorf("result[0] = %q; want %q", result[0], "/C/mixed/path/to/file")
	}
}

func TestEnsureUnixPaths_NonPathStrings(t *testing.T) {
	input := []string{"--verbose", "-n", "42", "plain-arg"}

	result := EnsureUnixPaths(input)

	for i, want := range input {
		if result[i] != want {
			t.Errorf("result[%d] = %q; want %q", i, result[i], want)
		}
	}
}

func TestEnsureUnixPaths_EmptyStringElement(t *testing.T) {
	input := []string{""}

	result := EnsureUnixPaths(input)

	if len(result) != 1 {
		t.Fatalf("got %d results; want 1", len(result))
	}

	if result[0] != "" {
		t.Errorf("result[0] = %q; want empty string", result[0])
	}
}

func TestEnsureUnixPaths_SingleElement(t *testing.T) {
	input := []string{`path\to\file`}

	result := EnsureUnixPaths(input)

	if result[0] != "path/to/file" {
		t.Errorf("result[0] = %q; want %q", result[0], "path/to/file")
	}
}

func TestEnsureUnixPaths_DoesNotMutateInput(t *testing.T) {
	input := []string{`C:\original\path`, `D:\another\path`}
	originalFirst := input[0]
	originalSecond := input[1]

	_ = EnsureUnixPaths(input)

	if input[0] != originalFirst {
		t.Errorf("input[0] was mutated: got %q; want %q", input[0], originalFirst)
	}

	if input[1] != originalSecond {
		t.Errorf("input[1] was mutated: got %q; want %q", input[1], originalSecond)
	}
}

func TestEnsureUnixPaths_NilSlice(t *testing.T) {
	result := EnsureUnixPaths(nil)

	if len(result) != 0 {
		t.Errorf("got %d results; want 0", len(result))
	}
}
