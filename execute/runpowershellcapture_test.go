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
	"strings"
	"testing"
)

func skipIfPowerShellNotAvailable(t *testing.T) {
	t.Helper()

	if _, err := findPowerShell(); err != nil {
		t.Skipf("PowerShell not available: %v", err)
	}
}

func TestRunPowerShellCapture_EmptyCommand(t *testing.T) {
	_, err := RunPowerShellCapture("")
	if err == nil {
		t.Fatal("expected error for empty command")
	}

	if !strings.Contains(err.Error(), "command must not be empty") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestRunPowerShellCapture_WhitespaceCommand(t *testing.T) {
	_, err := RunPowerShellCapture("   \t  ")
	if err == nil {
		t.Fatal("expected error for whitespace-only command")
	}

	if !strings.Contains(err.Error(), "command must not be empty") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestRunPowerShellCapture_SimpleOutput(t *testing.T) {
	skipIfPowerShellNotAvailable(t)

	out, err := RunPowerShellCapture("Write-Output 'hello world'")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out != "hello world" {
		t.Errorf("got %q, want %q", out, "hello world")
	}
}

func TestRunPowerShellCapture_OutputIsTrimmed(t *testing.T) {
	skipIfPowerShellNotAvailable(t)

	out, err := RunPowerShellCapture("Write-Output '  padded  '")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out != "padded" {
		t.Errorf("output not trimmed: got %q", out)
	}
}

func TestRunPowerShellCapture_MultilineOutput(t *testing.T) {
	skipIfPowerShellNotAvailable(t)

	out, err := RunPowerShellCapture("Write-Output 'line1'; Write-Output 'line2'")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(out, "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d: %q", len(lines), out)
	}

	if strings.TrimSpace(lines[0]) != "line1" {
		t.Errorf("line 1: got %q, want %q", lines[0], "line1")
	}

	if strings.TrimSpace(lines[1]) != "line2" {
		t.Errorf("line 2: got %q, want %q", lines[1], "line2")
	}
}

func TestRunPowerShellCapture_EmptyOutput(t *testing.T) {
	skipIfPowerShellNotAvailable(t)

	out, err := RunPowerShellCapture("Write-Output ''")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out != "" {
		t.Errorf("expected empty string, got %q", out)
	}
}

func TestRunPowerShellCapture_NonZeroExitReturnsError(t *testing.T) {
	skipIfPowerShellNotAvailable(t)

	_, err := RunPowerShellCapture("exit 1")
	if err == nil {
		t.Fatal("expected error for non-zero exit")
	}
}

func TestRunPowerShellCapture_StderrIncludedInError(t *testing.T) {
	skipIfPowerShellNotAvailable(t)

	_, err := RunPowerShellCapture("Write-Error 'something broke'; exit 1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "something broke") {
		t.Errorf("error should contain stderr text, got: %v", err)
	}
}

func TestRunPowerShellCapture_StderrDiscardedOnSuccess(t *testing.T) {
	skipIfPowerShellNotAvailable(t)

	out, err := RunPowerShellCapture(
		"Write-Error 'ignore this' -ErrorAction Continue; Write-Output 'success'",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out != "success" {
		t.Errorf("got %q, want %q", out, "success")
	}
}
