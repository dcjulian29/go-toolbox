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
	"testing"
	"time"
)

func TestExternalProgramContext_Success(t *testing.T) {
	ctx := context.Background()
	err := ExternalProgramContext(ctx, "go", "version")
	if err != nil {
		t.Errorf("ExternalProgramContext(go version) returned error: %v", err)
	}
}

func TestExternalProgramContext_ProgramNotFound(t *testing.T) {
	ctx := context.Background()
	err := ExternalProgramContext(ctx, "program-that-does-not-exist-xyz")
	if err == nil {
		t.Fatal("ExternalProgramContext should return error for nonexistent program")
	}

	var execErr *exec.Error
	if !errors.As(err, &execErr) {
		t.Errorf("expected *exec.Error, got %T: %v", err, err)
	}
}

func TestExternalProgramContext_NonZeroExit(t *testing.T) {
	ctx := context.Background()
	err := ExternalProgramContext(ctx, "go", "nonexistent-subcommand")
	if err == nil {
		t.Fatal("ExternalProgramContext should return error for non-zero exit")
	}

	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		t.Errorf("expected *exec.ExitError, got %T: %v", err, err)
	}
}

func TestExternalProgramContext_NonZeroExitCode(t *testing.T) {
	ctx := context.Background()
	var err error

	if runtime.GOOS == "windows" {
		err = ExternalProgramContext(ctx, "cmd", "/c", "exit", "42")
	} else {
		err = ExternalProgramContext(ctx, "sh", "-c", "exit 42")
	}

	if err == nil {
		t.Fatal("ExternalProgramContext should return error for non-zero exit code")
	}

	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("expected *exec.ExitError, got %T: %v", err, err)
	}

	if exitErr.ExitCode() != 42 {
		t.Errorf("expected exit code 42, got %d", exitErr.ExitCode())
	}
}

func longRunningCommand(ctx context.Context) error {
	// Use ping with a high count as a cross-platform long-running command.
	// ping does not require console input, so it works in test environments
	// where stdin may be redirected.
	if runtime.GOOS == "windows" {
		return ExternalProgramContext(ctx, "ping", "-n", "100", "127.0.0.1")
	}

	return ExternalProgramContext(ctx, "ping", "-c", "100", "127.0.0.1")
}

func TestExternalProgramContext_CancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	err := longRunningCommand(ctx)
	if err == nil {
		t.Fatal("ExternalProgramContext should return error for cancelled context")
	}
}

func TestExternalProgramContext_DeadlineExceeded(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := longRunningCommand(ctx)
	if err == nil {
		t.Fatal("ExternalProgramContext should return error when deadline exceeded")
	}

	if ctx.Err() != context.DeadlineExceeded {
		t.Errorf("expected context.DeadlineExceeded, got %v", ctx.Err())
	}
}

func TestExternalProgramContext_NoParams(t *testing.T) {
	ctx := context.Background()
	err := ExternalProgramContext(ctx, "go")
	if err == nil {
		t.Log("go with no args exited 0 (unexpected but not a test failure)")
	}
}

func TestExternalProgramContext_MultipleParams(t *testing.T) {
	ctx := context.Background()
	err := ExternalProgramContext(ctx, "go", "env", "GOPATH")
	if err != nil {
		t.Errorf("ExternalProgramContext(go env GOPATH) returned error: %v", err)
	}
}

func TestExternalProgramContext_EmptyProgramName(t *testing.T) {
	ctx := context.Background()
	err := ExternalProgramContext(ctx, "")
	if err == nil {
		t.Error("ExternalProgramContext should return error for empty program name")
	}
}

func TestExternalProgramContext_PlatformEchoSuccess(t *testing.T) {
	ctx := context.Background()
	var err error

	if runtime.GOOS == "windows" {
		err = ExternalProgramContext(ctx, "cmd", "/c", "echo", "hello")
	} else {
		err = ExternalProgramContext(ctx, "echo", "hello")
	}

	if err != nil {
		t.Errorf("echo command should succeed: %v", err)
	}
}
