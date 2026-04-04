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

func TestExpandHome(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("could not determine home directory for test setup: %v", err)
	}

	tests := []struct {
		name  string
		input string
		want  string
	}{
		// ── Unix-style prefix ─────────────────────────────────────────────
		{
			name:  "unix tilde-slash prefix expands to home",
			input: "~/documents",
			want:  filepath.Join(home, "documents"),
		},
		{
			name:  "unix tilde-slash prefix with nested path",
			input: "~/a/b/c",
			want:  filepath.Join(home, "a", "b", "c"),
		},
		{
			name:  "unix tilde-slash prefix with no suffix path component",
			input: "~/",
			want:  filepath.Join(home, ""),
		},

		// ── Windows-style prefix ──────────────────────────────────────────
		{
			name:  "windows tilde-backslash prefix expands to home",
			input: `~\documents`,
			want:  filepath.Join(home, "documents"),
		},
		{
			name:  "windows tilde-backslash prefix with nested path",
			input: `~\a\b\c`,
			want:  filepath.Join(home, "a", "b", "c"),
		},

		// ── No expansion expected ─────────────────────────────────────────
		{
			name:  "bare tilde is returned unchanged",
			input: "~",
			want:  "~",
		},
		{
			name:  "absolute path is returned unchanged",
			input: "/usr/local/bin",
			want:  "/usr/local/bin",
		},
		{
			name:  "relative path without tilde is returned unchanged",
			input: "some/relative/path",
			want:  "some/relative/path",
		},
		{
			name:  "empty string is returned unchanged",
			input: "",
			want:  "",
		},
		{
			name:  "tilde in the middle of a path is returned unchanged",
			input: "/home/~user/docs",
			want:  "/home/~user/docs",
		},
		{
			name:  "tilde at end of path is returned unchanged",
			input: "/some/path~",
			want:  "/some/path~",
		},
		{
			name:  "single character string is returned unchanged",
			input: "a",
			want:  "a",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ExpandHome(tc.input)
			if got != tc.want {
				t.Errorf("ExpandHome(%q)\n  got:  %q\n  want: %q", tc.input, got, tc.want)
			}
		})
	}
}

// TestExpandHome_HomeEnvUnset verifies that ExpandHome returns the original
// path unchanged when the operating system cannot resolve the home directory.
// We simulate this by unsetting both HOME and USERPROFILE before the call and
// restoring them afterwards.
func TestExpandHome_HomeEnvUnset(t *testing.T) {
	origHome := os.Getenv("HOME")
	origUserProf := os.Getenv("USERPROFILE")

	t.Cleanup(func() {
		_ = os.Setenv("HOME", origHome)
		_ = os.Setenv("USERPROFILE", origUserProf)
	})

	_ = os.Unsetenv("HOME")
	_ = os.Unsetenv("USERPROFILE")

	// os.UserHomeDir() will now fail on most platforms.
	// If it still succeeds (e.g. via /etc/passwd), skip the test rather than
	// report a false failure.
	if _, err := os.UserHomeDir(); err == nil {
		t.Skip("os.UserHomeDir() succeeded even without HOME/USERPROFILE; skipping unset test")
	}

	input := "~/documents"
	got := ExpandHome(input)

	if got != input {
		t.Errorf("ExpandHome(%q) with no home dir\n  got:  %q\n  want: %q (unchanged)", input, got, input)
	}
}
