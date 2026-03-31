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
	"testing"
)

// saveAndRestoreArgs saves os.Args and registers a cleanup function
// that restores the original value.
func saveAndRestoreArgs(t *testing.T) {
	t.Helper()

	original := make([]string, len(os.Args))
	copy(original, os.Args)

	t.Cleanup(func() {
		os.Args = original
	})
}

func TestEnsureUnixPathArguments_ExcludesProgramName(t *testing.T) {
	saveAndRestoreArgs(t)
	os.Args = []string{`C:\program.exe`, `C:\arg.txt`}

	result := EnsureUnixPathArguments()

	if len(result) != 1 {
		t.Fatalf("got %d results; want 1 (program name excluded)", len(result))
	}

	if result[0] != "/C/arg.txt" {
		t.Errorf("result[0] = %q; want %q", result[0], "/C/arg.txt")
	}
}

func TestEnsureUnixPathArguments_NoArguments(t *testing.T) {
	saveAndRestoreArgs(t)
	os.Args = []string{"program"}

	result := EnsureUnixPathArguments()

	if len(result) != 0 {
		t.Errorf("got %d results; want 0", len(result))
	}
}

func TestEnsureUnixPathArguments_ConvertsArguments(t *testing.T) {
	saveAndRestoreArgs(t)
	os.Args = []string{"program", `C:\Users\test`, `/unix/path`}

	result := EnsureUnixPathArguments()

	if len(result) != 2 {
		t.Fatalf("got %d results; want 2", len(result))
	}

	if result[0] != "/C/Users/test" {
		t.Errorf("result[0] = %q; want %q", result[0], "/C/Users/test")
	}

	if result[1] != "/unix/path" {
		t.Errorf("result[1] = %q; want %q", result[1], "/unix/path")
	}
}
