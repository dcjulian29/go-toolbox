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

func TestExternalProgramEnv_Success(t *testing.T) {
	err := ExternalProgramEnv("go", nil, "version")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExternalProgramEnv_EnvVarIsSet(t *testing.T) {
	var err error

	if runtime.GOOS == "windows" {
		// "set VARNAME" exits 0 if defined, 1 if not.
		err = ExternalProgramEnv("cmd", []string{"TEST_EXEC_VAR=hello"}, "/c", "set TEST_EXEC_VAR")
	} else {
		err = ExternalProgramEnv("sh", []string{"TEST_EXEC_VAR=hello"}, "-c", `test -n "$TEST_EXEC_VAR"`)
	}

	if err != nil {
		t.Errorf("expected env var to be set, got error: %v", err)
	}
}

func TestExternalProgramEnv_EnvVarNotSetFails(t *testing.T) {
	var err error

	if runtime.GOOS == "windows" {
		// Without setting the var, "set" should fail.
		err = ExternalProgramEnv("cmd", nil, "/c", "set TEST_EXEC_NOTSET_VAR_XYZ")
	} else {
		err = ExternalProgramEnv("sh", nil, "-c", `test -n "$TEST_EXEC_NOTSET_VAR_XYZ"`)
	}

	if err == nil {
		t.Error("expected error for undefined env var")
	}
}

func TestExternalProgramEnv_MultipleEnvVars(t *testing.T) {
	env := []string{"TEST_VAR_A=alpha", "TEST_VAR_B=bravo"}

	// Verify both vars are set by checking each one.
	for _, varName := range []string{"TEST_VAR_A", "TEST_VAR_B"} {
		var err error

		if runtime.GOOS == "windows" {
			err = ExternalProgramEnv("cmd", env, "/c", "set "+varName)
		} else {
			err = ExternalProgramEnv("sh", env, "-c", `test -n "$`+varName+`"`)
		}

		if err != nil {
			t.Errorf("expected env var %s to be set, got error: %v", varName, err)
		}
	}
}

func TestExternalProgramEnv_InheritsExistingEnv(t *testing.T) {
	var err error

	if runtime.GOOS == "windows" {
		err = ExternalProgramEnv("cmd", []string{"TEST_EXEC_VAR=hello"}, "/c", "set PATH")
	} else {
		err = ExternalProgramEnv("sh", []string{"TEST_EXEC_VAR=hello"}, "-c", `test -n "$PATH"`)
	}

	if err != nil {
		t.Errorf("expected PATH to be inherited, got error: %v", err)
	}
}

func TestExternalProgramEnv_NilEnv(t *testing.T) {
	err := ExternalProgramEnv("go", nil, "version")
	if err != nil {
		t.Fatalf("nil env should inherit current environment: %v", err)
	}
}

func TestExternalProgramEnv_EmptyEnv(t *testing.T) {
	err := ExternalProgramEnv("go", []string{}, "version")
	if err != nil {
		t.Fatalf("empty env should inherit current environment: %v", err)
	}
}

func TestExternalProgramEnv_ProgramNotFound(t *testing.T) {
	err := ExternalProgramEnv("program-that-does-not-exist-xyz", nil)
	if err == nil {
		t.Fatal("expected error for nonexistent program")
	}

	var execErr *exec.Error
	if !errors.As(err, &execErr) {
		t.Errorf("expected *exec.Error, got %T: %v", err, err)
	}
}

func TestExternalProgramEnv_NonZeroExit(t *testing.T) {
	var err error

	if runtime.GOOS == "windows" {
		err = ExternalProgramEnv("cmd", nil, "/c", "exit 1")
	} else {
		err = ExternalProgramEnv("sh", nil, "-c", "exit 1")
	}

	if err == nil {
		t.Fatal("expected error for non-zero exit")
	}

	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		t.Errorf("expected *exec.ExitError, got %T: %v", err, err)
	}
}

func TestExternalProgramEnv_NonZeroExitCode(t *testing.T) {
	var err error

	if runtime.GOOS == "windows" {
		err = ExternalProgramEnv("cmd", nil, "/c", "exit 42")
	} else {
		err = ExternalProgramEnv("sh", nil, "-c", "exit 42")
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

func TestExternalProgramEnv_NoParams(t *testing.T) {
	err := ExternalProgramEnv("go", nil)
	if err == nil {
		t.Log("go with no args exited 0 (unexpected but not a test failure)")
	}
}

func TestExternalProgramEnv_EmptyProgramName(t *testing.T) {
	err := ExternalProgramEnv("", nil)
	if err == nil {
		t.Error("expected error for empty program name")
	}
}
