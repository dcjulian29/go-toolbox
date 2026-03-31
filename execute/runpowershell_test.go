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

package execute

import (
	"errors"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// requirePowerShell skips the test if no PowerShell executable is available.
func requirePowerShell(t *testing.T) {
	t.Helper()

	if _, err := findPowerShell(); err != nil {
		t.Skip("PowerShell is not installed or available")
	}
}

func TestFindPowerShell_ReturnsPath(t *testing.T) {
	requirePowerShell(t)

	path, err := findPowerShell()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if path == "" {
		t.Fatal("expected non-empty path")
	}

	if !filepath.IsAbs(path) {
		t.Errorf("expected absolute path, got %q", path)
	}
}

func TestFindPowerShell_PrefersPwsh(t *testing.T) {
	if _, err := exec.LookPath("pwsh"); err != nil {
		t.Skip("pwsh is not available")
	}

	path, err := findPowerShell()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	base := strings.ToLower(filepath.Base(path))
	if !strings.HasPrefix(base, "pwsh") {
		t.Errorf("expected pwsh to be preferred, got %q", path)
	}
}

func TestFindPowerShell_FallbackOnWindows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("fallback to powershell only applies on Windows")
	}

	// On Windows, at least one of pwsh or powershell should be found.
	path, err := findPowerShell()
	if err != nil {
		t.Fatalf("expected PowerShell to be available on Windows: %v", err)
	}

	if path == "" {
		t.Fatal("expected non-empty path")
	}
}

func TestRunPowerShell_Success(t *testing.T) {
	requirePowerShell(t)

	err := RunPowerShell("Write-Output 'hello'")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunPowerShell_NonZeroExit(t *testing.T) {
	requirePowerShell(t)

	err := RunPowerShell("exit 1")
	if err == nil {
		t.Fatal("expected error for non-zero exit")
	}

	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		t.Errorf("expected *exec.ExitError, got %T: %v", err, err)
	}
}

func TestRunPowerShell_NonZeroExitCode(t *testing.T) {
	requirePowerShell(t)

	err := RunPowerShell("exit 42")
	if err == nil {
		t.Fatal("expected error for non-zero exit code")
	}

	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("expected *exec.ExitError, got %T: %v", err, err)
	}

	if exitErr.ExitCode() != 42 {
		t.Errorf("expected exit code 42, got %d", exitErr.ExitCode())
	}
}

func TestRunPowerShell_EmptyCommand(t *testing.T) {
	requirePowerShell(t)

	// Empty command should not panic; PowerShell handles it.
	_ = RunPowerShell("")
}

func TestRunPowerShell_InvalidCommand(t *testing.T) {
	requirePowerShell(t)

	err := RunPowerShell("Invoke-NonExistentCmdlet-XYZ")
	if err == nil {
		t.Error("expected error for invalid cmdlet")
	}
}

func TestRunPowerShell_MultiStatement(t *testing.T) {
	requirePowerShell(t)

	err := RunPowerShell("Write-Output 'a'; Write-Output 'b'")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunPowerShell_WhitespaceCommand(t *testing.T) {
	requirePowerShell(t)

	// Whitespace-only command should not panic.
	_ = RunPowerShell("   ")
}
