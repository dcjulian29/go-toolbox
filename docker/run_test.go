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

package docker

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dcjulian29/go-toolbox/filesystem"
	"github.com/dcjulian29/go-toolbox/textformat"
)

func TestInteractiveArguments_Detached(t *testing.T) {
	opts := ContainerOptions{Interactive: false}

	args := interactiveArguments(opts)

	if len(args) != 1 || args[0] != "--detach" {
		t.Errorf("got %v, want [--detach]", args)
	}
}

func TestInteractiveArguments_InteractiveWithTty(t *testing.T) {
	opts := ContainerOptions{Interactive: true}

	args := interactiveArguments(opts)

	if len(args) != 2 || args[0] != "--interactive" || args[1] != "--tty" {
		t.Errorf("got %v, want [--interactive --tty]", args)
	}
}

func TestInteractiveArguments_InteractiveWithoutTty(t *testing.T) {
	opts := ContainerOptions{Interactive: true, NoTty: true}

	args := interactiveArguments(opts)

	if len(args) != 1 || args[0] != "--interactive" {
		t.Errorf("got %v, want [--interactive]", args)
	}
}

func TestContainerArguments_AllFields(t *testing.T) {
	opts := ContainerOptions{
		HostName:         "myhost",
		Name:             "mycontainer",
		ReadOnly:         true,
		User:             "1000:1000",
		WorkingDirectory: "/app",
	}

	args := containerArguments(opts)

	expected := []string{
		"--rm",
		"--hostname", "myhost",
		"--name", "mycontainer",
		"--read-only",
		"--user", "1000:1000",
		"--workdir", "/app",
	}

	if len(args) != len(expected) {
		t.Fatalf("got %d args %v, want %d args %v", len(args), args, len(expected), expected)
	}

	for i, arg := range expected {
		if args[i] != arg {
			t.Errorf("args[%d] = %q, want %q", i, args[i], arg)
		}
	}
}

func TestContainerArguments_KeepDisablesRm(t *testing.T) {
	opts := ContainerOptions{Keep: true}

	args := containerArguments(opts)

	for _, arg := range args {
		if arg == "--rm" {
			t.Error("--rm should not be present when Keep is true")
		}
	}
}

func TestContainerArguments_EmptyOptions(t *testing.T) {
	opts := ContainerOptions{}

	args := containerArguments(opts)

	if len(args) != 1 || args[0] != "--rm" {
		t.Errorf("got %v, want [--rm]", args)
	}
}

func TestEnvironmentArguments_MultipleVars(t *testing.T) {
	opts := ContainerOptions{
		EnvironmentVariables: map[string]string{
			"DB_HOST": "localhost",
			"DB_PORT": "5432",
		},
	}

	args := environmentArguments(opts)

	if len(args) != 4 {
		t.Fatalf("expected 4 args, got %d: %v", len(args), args)
	}

	pairs := make(map[string]bool)

	for i := 0; i < len(args); i += 2 {
		if args[i] != "-e" {
			t.Errorf("expected '-e' at index %d, got %q", i, args[i])
		}

		pairs[args[i+1]] = true
	}

	if !pairs["DB_HOST=localhost"] {
		t.Error("missing DB_HOST=localhost")
	}

	if !pairs["DB_PORT=5432"] {
		t.Error("missing DB_PORT=5432")
	}
}

func TestEnvironmentArguments_Empty(t *testing.T) {
	opts := ContainerOptions{}

	args := environmentArguments(opts)

	if len(args) != 0 {
		t.Errorf("expected empty slice, got %v", args)
	}
}

func TestPortArguments_MultiplePorts(t *testing.T) {
	opts := ContainerOptions{
		Ports: []string{"8080:80", "443:443"},
	}

	args := portArguments(opts)

	expected := []string{"-p", "8080:80", "-p", "443:443"}

	if len(args) != len(expected) {
		t.Fatalf("got %v, want %v", args, expected)
	}

	for i, arg := range expected {
		if args[i] != arg {
			t.Errorf("args[%d] = %q, want %q", i, args[i], arg)
		}
	}
}

func TestPortArguments_Empty(t *testing.T) {
	opts := ContainerOptions{}

	args := portArguments(opts)

	if len(args) != 0 {
		t.Errorf("expected empty slice, got %v", args)
	}
}

func TestVolumeArguments_MultipleVolumes(t *testing.T) {
	opts := ContainerOptions{
		Volumes: []string{"/host:/container", "/data:/data:ro"},
	}

	args := volumeArguments(opts)

	expected := []string{"--volume", "/host:/container", "--volume", "/data:/data:ro"}

	if len(args) != len(expected) {
		t.Fatalf("got %v, want %v", args, expected)
	}

	for i, arg := range expected {
		if args[i] != arg {
			t.Errorf("args[%d] = %q, want %q", i, args[i], arg)
		}
	}
}

func TestVolumeArguments_Empty(t *testing.T) {
	opts := ContainerOptions{}

	args := volumeArguments(opts)

	if len(args) != 0 {
		t.Errorf("expected empty slice, got %v", args)
	}
}

func TestCommandArguments_CommandOnly(t *testing.T) {
	opts := ContainerOptions{Command: "/bin/sh"}

	args := commandArguments(opts)

	if len(args) != 1 || args[0] != "/bin/sh" {
		t.Errorf("got %v, want [/bin/sh]", args)
	}
}

func TestCommandArguments_AdditionalArgsOnly(t *testing.T) {
	opts := ContainerOptions{
		AdditionalArgs: []string{"--verbose", "--debug"},
	}

	args := commandArguments(opts)

	if len(args) != 2 || args[0] != "--verbose" || args[1] != "--debug" {
		t.Errorf("got %v, want [--verbose --debug]", args)
	}
}

func TestCommandArguments_CommandAndAdditionalArgs(t *testing.T) {
	opts := ContainerOptions{
		Command:        "/bin/sh",
		AdditionalArgs: []string{"-c", "echo hello"},
	}

	args := commandArguments(opts)

	expected := []string{"/bin/sh", "-c", "echo hello"}

	if len(args) != len(expected) {
		t.Fatalf("got %v, want %v", args, expected)
	}

	for i, arg := range expected {
		if args[i] != arg {
			t.Errorf("args[%d] = %q, want %q", i, args[i], arg)
		}
	}
}

func TestCommandArguments_Empty(t *testing.T) {
	opts := ContainerOptions{}

	args := commandArguments(opts)

	if len(args) != 0 {
		t.Errorf("expected empty slice, got %v", args)
	}
}

func TestEntryArguments_EntryPoint(t *testing.T) {
	opts := ContainerOptions{EntryPoint: "/bin/bash"}

	args, vol, err := entryArguments(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(vol) != 0 {
		t.Errorf("expected no volume, got %v", vol)
	}

	if len(args) != 2 || args[0] != "--entrypoint" || args[1] != "/bin/bash" {
		t.Errorf("got %v, want [--entrypoint /bin/bash]", args)
	}
}

func TestEntryArguments_EntryPointTakesPrecedence(t *testing.T) {
	opts := ContainerOptions{
		EntryPoint:  "/bin/bash",
		EntryScript: "script.sh",
	}

	args, vol, err := entryArguments(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(vol) != 0 {
		t.Errorf("EntryScript volume should be ignored when EntryPoint is set, got %v", vol)
	}

	if len(args) != 2 || args[1] != "/bin/bash" {
		t.Errorf("got %v, want [--entrypoint /bin/bash]", args)
	}
}

func TestEntryArguments_Neither(t *testing.T) {
	opts := ContainerOptions{}

	args, vol, err := entryArguments(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if args != nil {
		t.Errorf("expected nil args, got %v", args)
	}

	if vol != nil {
		t.Errorf("expected nil vol, got %v", vol)
	}
}

func TestEntryArguments_EntryScript(t *testing.T) {
	dir := t.TempDir()
	script := filepath.Join(dir, "myscript.sh")

	if err := os.WriteFile(script, []byte("#!/bin/sh"), filesystem.ModeOwnerReadWrite); err != nil {
		t.Fatalf("failed to create temp script: %v", err)
	}

	opts := ContainerOptions{EntryScript: script}

	args, vol, err := entryArguments(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(args) != 2 || args[0] != "--entrypoint" || args[1] != "/bin/myscript.sh" {
		t.Errorf("args = %v, want [--entrypoint /bin/myscript.sh]", args)
	}

	if len(vol) != 1 {
		t.Fatalf("expected 1 volume, got %d: %v", len(vol), vol)
	}

	if vol[0] == textformat.EmptyString {
		t.Error("volume should not be empty")
	}
}

func TestEntryArguments_EntryScriptDirectory(t *testing.T) {
	dir := t.TempDir()

	opts := ContainerOptions{EntryScript: dir}

	_, _, err := entryArguments(opts)
	if err == nil {
		t.Fatal("expected error for directory path")
	}

	if err.Error() != "path is a directory not a file" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestEntryArguments_EntryScriptNotFound(t *testing.T) {
	opts := ContainerOptions{EntryScript: "/nonexistent/script.sh"}

	_, _, err := entryArguments(opts)
	if err == nil {
		t.Fatal("expected error for nonexistent path")
	}
}

func TestRun_EmptyImage(t *testing.T) {
	opts := ContainerOptions{}

	_, err := Run(opts)
	if err == nil {
		t.Fatal("expected error for empty image")
	}

	if err.Error() != "image is required" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRun_DefaultTag(t *testing.T) {
	opts := ContainerOptions{
		Image: "alpine",
	}

	if opts.Tag != textformat.EmptyString {
		t.Fatal("precondition: Tag should be empty before Run")
	}

	// Run will fail because docker is not available, but we can verify
	// the error is not about the image or tag validation
	_, err := Run(opts)
	if err == nil {
		return
	}

	if err.Error() == "image is required" {
		t.Error("should not fail image validation with Image set")
	}
}
