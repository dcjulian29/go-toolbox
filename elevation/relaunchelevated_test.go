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

package elevation

import (
	"strings"
	"testing"
)

func TestShellExecuteError_KnownCodes(t *testing.T) {
	tests := []struct {
		contains string
		code     uintptr
	}{
		{"out of memory", 0},
		{"file not found", 2},
		{"path not found", 3},
		{"access denied", 5},
		{"out of memory", 8},
		{"sharing violation", 26},
		{"invalid file association", 27},
		{"no application associated", 31},
		{"DDE transaction failed", 32},
	}

	for _, tt := range tests {
		err := shellExecuteError(tt.code)
		if err == nil {
			t.Errorf("code %d: expected error, got nil", tt.code)

			continue
		}

		if !strings.Contains(err.Error(), tt.contains) {
			t.Errorf("code %d: got %q, want substring %q", tt.code, err.Error(), tt.contains)
		}
	}
}

func TestShellExecuteError_UnknownCode(t *testing.T) {
	err := shellExecuteError(17)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "17") {
		t.Errorf("expected error to contain code 17, got %q", err.Error())
	}
}
