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
	"context"
	"errors"
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestExternalProgramContextCapture_Success(t *testing.T) {
	ctx := context.Background()
	got, err := ExternalProgramContextCapture(ctx, "go", "version")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.HasPrefix(got, "go version") {
		t.Errorf("expected output starting with 'go version', got %q", got)
	}
}

func TestExternalProgramContextCapture_OutputIsTrimmed(t *testing.T) {
	ctx := context.Background()
	var got string
	var err error

	if runtime.GOOS == "windows" {
		got, err = ExternalProgramContextCapture(ctx, "cmd", "/c", "echo", "hello")
	} else {
		got, err = ExternalProgramContextCapture(ctx, "echo", "hello")
	}

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "hello" {
		t.Errorf("expected trimmed output %q, got %q", "hello", got)
	}
}

func TestExternalProgramContextCapture_EmptyOutput(t *testing.T) {
	ctx := context.Background()
	var got string
	var err error

	if runtime.GOOS == "windows" {
		got, err = ExternalProgramContextCapture(ctx, "cmd", "/c", "echo.")
	} else {
		got, err = ExternalProgramContextCapture(ctx, "sh", "-c", "printf ''")
	}

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "" {
		t.Errorf("expected empty output, got %q", got)
	}
}

func TestExternalProgramContextCapture_ProgramNotFound(t *testing.T) {
	ctx := context.Background()
	got, err := ExternalProgramContextCapture(ctx, "program-that-does-not-exist-xyz")
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

func TestExternalProgramContextCapture_NonZeroExit(t *testing.T) {
	ctx := context.Background()
	var got string
	var err error

	if runtime.GOOS == "windows" {
		got, err = ExternalProgramContextCapture(ctx, "cmd", "/c", "exit", "1")
	} else {
		got, err = ExternalProgramContextCapture(ctx, "sh", "-c", "exit 1")
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

func TestExternalProgramContextCapture_ErrorContainsStderr(t *testing.T) {
	ctx := context.Background()
	var err error

	if runtime.GOOS == "windows" {
		_, err = ExternalProgramContextCapture(ctx, "cmd", "/c", "echo error message>&2 & exit 1")
	} else {
		_, err = ExternalProgramContextCapture(ctx, "sh", "-c", "echo 'error message' >&2; exit 1")
	}

	if err == nil {
		t.Fatal("expected error for non-zero exit")
	}

	if !strings.Contains(err.Error(), "error message") {
		t.Errorf("expected error to contain stderr output, got: %v", err)
	}
}

func TestExternalProgramContextCapture_StderrDiscardedOnSuccess(t *testing.T) {
	ctx := context.Background()
	var got string
	var err error

	// Write to stderr but exit 0.
	if runtime.GOOS == "windows" {
		got, err = ExternalProgramContextCapture(ctx, "cmd", "/c", "echo warning>&2 & echo output")
	} else {
		got, err = ExternalProgramContextCapture(ctx, "sh", "-c", "echo 'warning' >&2; echo 'output'")
	}

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "output" {
		t.Errorf("expected %q, got %q", "output", got)
	}
}

func TestExternalProgramContextCapture_CancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	_, err := longRunningCaptureCommand(ctx)
	if err == nil {
		t.Fatal("expected error for cancelled context")
	}
}

func TestExternalProgramContextCapture_DeadlineExceeded(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := longRunningCaptureCommand(ctx)
	if err == nil {
		t.Fatal("expected error when deadline exceeded")
	}

	if ctx.Err() != context.DeadlineExceeded {
		t.Errorf("expected context.DeadlineExceeded, got %v", ctx.Err())
	}
}

func longRunningCaptureCommand(ctx context.Context) (string, error) {
	if runtime.GOOS == "windows" {
		return ExternalProgramContextCapture(ctx, "ping", "-n", "100", "127.0.0.1")
	}

	return ExternalProgramContextCapture(ctx, "ping", "-c", "100", "127.0.0.1")
}

func TestExternalProgramContextCapture_MultipleParams(t *testing.T) {
	ctx := context.Background()
	got, err := ExternalProgramContextCapture(ctx, "go", "env", "GOPATH")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got == "" {
		t.Error("expected non-empty GOPATH output")
	}
}

func TestExternalProgramContextCapture_EmptyProgramName(t *testing.T) {
	ctx := context.Background()
	_, err := ExternalProgramContextCapture(ctx, "")
	if err == nil {
		t.Error("expected error for empty program name")
	}
}
