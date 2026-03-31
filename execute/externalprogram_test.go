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
	"testing"
)

func TestExternalProgram_Success(t *testing.T) {
	// "go version" is available in any Go test environment and exits 0.
	err := ExternalProgram("go", "version")
	if err != nil {
		t.Errorf("ExternalProgram(go version) returned error: %v", err)
	}
}

func TestExternalProgram_ProgramNotFound(t *testing.T) {
	err := ExternalProgram("program-that-does-not-exist-xyz")
	if err == nil {
		t.Fatal("ExternalProgram should return error for nonexistent program")
	}

	var execErr *exec.Error
	if !errors.As(err, &execErr) {
		t.Errorf("expected *exec.Error, got %T: %v", err, err)
	}
}

func TestExternalProgram_NonZeroExit(t *testing.T) {
	// "go nonexistent-subcommand" exits with a non-zero status.
	err := ExternalProgram("go", "nonexistent-subcommand")
	if err == nil {
		t.Fatal("ExternalProgram should return error for non-zero exit")
	}

	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		t.Errorf("expected *exec.ExitError, got %T: %v", err, err)
	}
}

func TestExternalProgram_NonZeroExitCode(t *testing.T) {
	var err error

	if runtime.GOOS == "windows" {
		err = ExternalProgram("cmd", "/c", "exit", "42")
	} else {
		err = ExternalProgram("sh", "-c", "exit 42")
	}

	if err == nil {
		t.Fatal("ExternalProgram should return error for non-zero exit code")
	}

	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("expected *exec.ExitError, got %T: %v", err, err)
	}

	if exitErr.ExitCode() != 42 {
		t.Errorf("expected exit code 42, got %d", exitErr.ExitCode())
	}
}

func TestExternalProgram_NoParams(t *testing.T) {
	// "go" with no arguments exits non-zero but should not panic or error unexpectedly.
	err := ExternalProgram("go")
	if err == nil {
		t.Log("go with no args exited 0 (unexpected but not a test failure)")
	}
}

func TestExternalProgram_MultipleParams(t *testing.T) {
	err := ExternalProgram("go", "env", "GOPATH")
	if err != nil {
		t.Errorf("ExternalProgram(go env GOPATH) returned error: %v", err)
	}
}

func TestExternalProgram_EmptyProgramName(t *testing.T) {
	err := ExternalProgram("")
	if err == nil {
		t.Error("ExternalProgram should return error for empty program name")
	}
}

func TestExternalProgram_PlatformEchoSuccess(t *testing.T) {
	var err error

	if runtime.GOOS == "windows" {
		err = ExternalProgram("cmd", "/c", "echo", "hello")
	} else {
		err = ExternalProgram("echo", "hello")
	}

	if err != nil {
		t.Errorf("echo command should succeed: %v", err)
	}
}
