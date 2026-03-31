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

func TestExternalProgramEnvCapture_Success(t *testing.T) {
	got, err := ExternalProgramEnvCapture("go", nil, "version")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.HasPrefix(got, "go version") {
		t.Errorf("expected output starting with 'go version', got %q", got)
	}
}

func TestExternalProgramEnvCapture_OutputIsTrimmed(t *testing.T) {
	var got string
	var err error

	if runtime.GOOS == "windows" {
		got, err = ExternalProgramEnvCapture("cmd", nil, "/c", "echo", "hello")
	} else {
		got, err = ExternalProgramEnvCapture("echo", nil, "hello")
	}

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "hello" {
		t.Errorf("expected trimmed output %q, got %q", "hello", got)
	}
}

func TestExternalProgramEnvCapture_EmptyOutput(t *testing.T) {
	var got string
	var err error

	if runtime.GOOS == "windows" {
		got, err = ExternalProgramEnvCapture("cmd", nil, "/c", "echo.")
	} else {
		got, err = ExternalProgramEnvCapture("sh", nil, "-c", "printf ''")
	}

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "" {
		t.Errorf("expected empty output, got %q", got)
	}
}

func TestExternalProgramEnvCapture_EnvVarCaptured(t *testing.T) {
	var got string
	var err error

	env := []string{"TEST_CAPTURE_VAR=captured_value"}

	if runtime.GOOS == "windows" {
		got, err = ExternalProgramEnvCapture("cmd", env, "/c", "echo", "%TEST_CAPTURE_VAR%")
	} else {
		got, err = ExternalProgramEnvCapture("sh", env, "-c", `echo "$TEST_CAPTURE_VAR"`)
	}

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "captured_value" {
		t.Errorf("expected %q, got %q", "captured_value", got)
	}
}

func TestExternalProgramEnvCapture_MultipleEnvVars(t *testing.T) {
	var got string
	var err error

	env := []string{"VAR_A=alpha", "VAR_B=bravo"}

	if runtime.GOOS == "windows" {
		got, err = ExternalProgramEnvCapture("cmd", env, "/c", "echo", "%VAR_A%-%VAR_B%")
	} else {
		got, err = ExternalProgramEnvCapture("sh", env, "-c", `echo "$VAR_A-$VAR_B"`)
	}

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "alpha-bravo" {
		t.Errorf("expected %q, got %q", "alpha-bravo", got)
	}
}

func TestExternalProgramEnvCapture_InheritsExistingEnv(t *testing.T) {
	// GOPATH should be available even with custom env vars.
	got, err := ExternalProgramEnvCapture("go", []string{"TEST_EXEC_VAR=hello"}, "env", "GOPATH")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got == "" {
		t.Error("expected GOPATH to be inherited, got empty string")
	}
}

func TestExternalProgramEnvCapture_NilEnv(t *testing.T) {
	got, err := ExternalProgramEnvCapture("go", nil, "version")
	if err != nil {
		t.Fatalf("nil env should inherit current environment: %v", err)
	}

	if !strings.HasPrefix(got, "go version") {
		t.Errorf("expected output starting with 'go version', got %q", got)
	}
}

func TestExternalProgramEnvCapture_EmptyEnv(t *testing.T) {
	got, err := ExternalProgramEnvCapture("go", []string{}, "version")
	if err != nil {
		t.Fatalf("empty env should inherit current environment: %v", err)
	}

	if !strings.HasPrefix(got, "go version") {
		t.Errorf("expected output starting with 'go version', got %q", got)
	}
}

func TestExternalProgramEnvCapture_ProgramNotFound(t *testing.T) {
	got, err := ExternalProgramEnvCapture("program-that-does-not-exist-xyz", nil)
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

func TestExternalProgramEnvCapture_NonZeroExit(t *testing.T) {
	var got string
	var err error

	if runtime.GOOS == "windows" {
		got, err = ExternalProgramEnvCapture("cmd", nil, "/c", "exit", "1")
	} else {
		got, err = ExternalProgramEnvCapture("sh", nil, "-c", "exit 1")
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

func TestExternalProgramEnvCapture_NonZeroExitCode(t *testing.T) {
	var err error

	if runtime.GOOS == "windows" {
		_, err = ExternalProgramEnvCapture("cmd", nil, "/c", "exit", "42")
	} else {
		_, err = ExternalProgramEnvCapture("sh", nil, "-c", "exit 42")
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

func TestExternalProgramEnvCapture_ErrorContainsStderr(t *testing.T) {
	var err error

	if runtime.GOOS == "windows" {
		_, err = ExternalProgramEnvCapture("cmd", nil, "/c", "echo error message>&2 & exit 1")
	} else {
		_, err = ExternalProgramEnvCapture("sh", nil, "-c", "echo 'error message' >&2; exit 1")
	}

	if err == nil {
		t.Fatal("expected error for non-zero exit")
	}

	if !strings.Contains(err.Error(), "error message") {
		t.Errorf("expected error to contain stderr output, got: %v", err)
	}
}

func TestExternalProgramEnvCapture_StderrDiscardedOnSuccess(t *testing.T) {
	var got string
	var err error

	if runtime.GOOS == "windows" {
		got, err = ExternalProgramEnvCapture("cmd", nil, "/c", "echo warning>&2 & echo output")
	} else {
		got, err = ExternalProgramEnvCapture("sh", nil, "-c", "echo 'warning' >&2; echo 'output'")
	}

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "output" {
		t.Errorf("expected %q, got %q", "output", got)
	}
}

func TestExternalProgramEnvCapture_EnvVarWithEnv(t *testing.T) {
	// Verify env vars work together with capture.
	var err error

	env := []string{"TEST_CAPTURE_VAR=hello"}

	if runtime.GOOS == "windows" {
		_, err = ExternalProgramEnvCapture("cmd", env, "/c", "echo", "%TEST_CAPTURE_VAR%")
	} else {
		_, err = ExternalProgramEnvCapture("sh", env, "-c", `echo "$TEST_CAPTURE_VAR"`)
	}

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExternalProgramEnvCapture_MultipleParams(t *testing.T) {
	got, err := ExternalProgramEnvCapture("go", nil, "env", "GOPATH")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got == "" {
		t.Error("expected non-empty GOPATH output")
	}
}

func TestExternalProgramEnvCapture_EmptyProgramName(t *testing.T) {
	_, err := ExternalProgramEnvCapture("", nil)
	if err == nil {
		t.Error("expected error for empty program name")
	}
}
