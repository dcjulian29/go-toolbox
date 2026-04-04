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

func TestExpandHome(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("could not determine home directory for test setup: %v", err)
	}

	tests := []struct {
		name   string
		input  string
		want   string
		onlyOS string // if set, skip this case on any other GOOS
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
			name:  "unix tilde-slash with no suffix path component",
			input: "~/",
			want:  filepath.Join(home, ""),
		},

		// ── Windows-style prefix — only valid on Windows ──────────────────
		// On Linux/macOS the backslash is a legal filename character, not a
		// separator, so filepath.Join treats the entire suffix as one component.
		// These cases are therefore only meaningful on Windows.
		{
			name:   "windows tilde-backslash prefix expands to home",
			input:  `~\documents`,
			want:   filepath.Join(home, "documents"),
			onlyOS: "windows",
		},
		{
			name:   "windows tilde-backslash prefix with nested path",
			input:  `~\a\b\c`,
			want:   filepath.Join(home, "a", "b", "c"),
			onlyOS: "windows",
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
			if tc.onlyOS != "" && runtime.GOOS != tc.onlyOS {
				t.Skipf("skipping %q: backslash is not a path separator on %s", tc.name, runtime.GOOS)
			}

			got := ExpandHome(tc.input)

			if got != tc.want {
				t.Errorf("ExpandHome(%q)\n  got:  %q\n  want: %q", tc.input, got, tc.want)
			}
		})
	}
}

// TestExpandHome_BackslashOnLinux documents the exact behavior of ExpandHome
// when given a Windows-style tilde path on a non-Windows platform.
//
// Because backslash is a valid filename character on Linux/macOS, the suffix
// after "~\" is treated as a single path component rather than a hierarchy of
// directories. This is a known limitation of the cross-platform implementation
// and this test locks in that documented behavior.
func TestExpandHome_BackslashOnLinux(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("not applicable on Windows — backslash is a real separator there")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("could not determine home directory: %v", err)
	}

	tests := []struct {
		name  string
		input string
		// On Linux the suffix after ~\ is ONE component (backslash is not a separator).
		want string
	}{
		{
			name:  "tilde-backslash suffix is one component on Linux",
			input: `~\documents`,
			// filepath.Join(home, `documents`) — only one level deep because
			// the backslash prefix is consumed by path[:2] but the rest is a
			// plain string with no OS path separator.
			want: filepath.Join(home, `documents`),
		},
		{
			name:  "tilde-backslash nested path is one opaque component on Linux",
			input: `~\a\b\c`,
			// filepath.Join(home, `a\b\c`) — the backslashes are literal
			// characters inside the single filename "a\b\c".
			want: filepath.Join(home, `a\b\c`),
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
func TestExpandHome_HomeEnvUnset(t *testing.T) {
	origHome := os.Getenv("HOME")
	origUserProf := os.Getenv("USERPROFILE")

	t.Cleanup(func() {
		_ = os.Setenv("HOME", origHome)
		_ = os.Setenv("USERPROFILE", origUserProf)
	})

	_ = os.Unsetenv("HOME")
	_ = os.Unsetenv("USERPROFILE")

	// If os.UserHomeDir() still resolves (e.g. via /etc/passwd on Linux),
	// the error path cannot be triggered — skip gracefully.
	if _, err := os.UserHomeDir(); err == nil {
		t.Skip("os.UserHomeDir() succeeded without HOME/USERPROFILE; skipping unset test")
	}

	input := "~/documents"
	got := ExpandHome(input)

	if got != input {
		t.Errorf("ExpandHome(%q) with no home dir\n  got:  %q\n  want: %q (unchanged)", input, got, input)
	}
}
