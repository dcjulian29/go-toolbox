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
	"runtime"
	"strings"
	"testing"
)

func TestExternalProgramCapture_Success(t *testing.T) {
	got, err := ExternalProgramCapture("go", "version")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.HasPrefix(got, "go version") {
		t.Errorf("expected output starting with 'go version', got %q", got)
	}
}

func TestExternalProgramCapture_OutputIsTrimmed(t *testing.T) {
	var got string
	var err error

	if runtime.GOOS == "windows" {
		got, err = ExternalProgramCapture("cmd", "/c", "echo", "hello")
	} else {
		got, err = ExternalProgramCapture("echo", "hello")
	}

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "hello" {
		t.Errorf("expected trimmed output %q, got %q", "hello", got)
	}
}

func TestExternalProgramCapture_EmptyOutput(t *testing.T) {
	var got string
	var err error

	if runtime.GOOS == "windows" {
		got, err = ExternalProgramCapture("cmd", "/c", "echo.")
	} else {
		got, err = ExternalProgramCapture("sh", "-c", "printf ''")
	}

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "" {
		t.Errorf("expected empty output, got %q", got)
	}
}

func TestExternalProgramCapture_ProgramNotFound(t *testing.T) {
	got, err := ExternalProgramCapture("program-that-does-not-exist-xyz")
	if err == nil {
		t.Fatal("expected error for nonexistent program")
	}

	if got != "" {
		t.Errorf("expected empty output on error, got %q", got)
	}

	var execErr *exec.Error
	if !errors.As(err, &execErr) {
		t.Errorf("expected *exec.Error, got %T: %v", err, err)
	}
}

func TestExternalProgramCapture_NonZeroExit(t *testing.T) {
	var got string
	var err error

	if runtime.GOOS == "windows" {
		got, err = ExternalProgramCapture("cmd", "/c", "exit", "1")
	} else {
		got, err = ExternalProgramCapture("sh", "-c", "exit 1")
	}

	if err == nil {
		t.Fatal("expected error for non-zero exit")
	}

	if got != "" {
		t.Errorf("expected empty output on error, got %q", got)
	}

	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		t.Errorf("expected *exec.ExitError, got %T: %v", err, err)
	}
}

func TestExternalProgramCapture_NonZeroExitCode(t *testing.T) {
	var err error

	if runtime.GOOS == "windows" {
		_, err = ExternalProgramCapture("cmd", "/c", "exit", "42")
	} else {
		_, err = ExternalProgramCapture("sh", "-c", "exit 42")
	}

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

func TestExternalProgramCapture_ErrorContainsStderr(t *testing.T) {
	var err error

	if runtime.GOOS == "windows" {
		_, err = ExternalProgramCapture("cmd", "/c", "echo error message>&2 & exit 1")
	} else {
		_, err = ExternalProgramCapture("sh", "-c", "echo 'error message' >&2; exit 1")
	}

	if err == nil {
		t.Fatal("expected error for non-zero exit")
	}

	if !strings.Contains(err.Error(), "error message") {
		t.Errorf("expected error to contain stderr output, got: %v", err)
	}
}

func TestExternalProgramCapture_StderrDiscardedOnSuccess(t *testing.T) {
	var got string
	var err error

	// Write to stderr but exit 0.
	if runtime.GOOS == "windows" {
		got, err = ExternalProgramCapture("cmd", "/c", "echo warning>&2 & echo output")
	} else {
		got, err = ExternalProgramCapture("sh", "-c", "echo 'warning' >&2; echo 'output'")
	}

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "output" {
		t.Errorf("expected %q, got %q", "output", got)
	}
}

func TestExternalProgramCapture_NoParams(t *testing.T) {
	_, err := ExternalProgramCapture("go")
	if err == nil {
		t.Log("go with no args exited 0 (unexpected but not a test failure)")
	}
}

func TestExternalProgramCapture_MultipleParams(t *testing.T) {
	got, err := ExternalProgramCapture("go", "env", "GOPATH")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got == "" {
		t.Error("expected non-empty GOPATH output")
	}
}

func TestExternalProgramCapture_EmptyProgramName(t *testing.T) {
	_, err := ExternalProgramCapture("")
	if err == nil {
		t.Error("expected error for empty program name")
	}
}

func TestExternalProgramCapture_PlatformEchoSuccess(t *testing.T) {
	var got string
	var err error

	if runtime.GOOS == "windows" {
		got, err = ExternalProgramCapture("cmd", "/c", "echo", "hello world")
	} else {
		got, err = ExternalProgramCapture("echo", "hello world")
	}

	if err != nil {
		t.Fatalf("echo command should succeed: %v", err)
	}

	if !strings.Contains(got, "hello world") {
		t.Errorf("expected output to contain 'hello world', got %q", got)
	}
}
