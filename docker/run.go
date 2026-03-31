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
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dcjulian29/go-toolbox/execute"
	"github.com/dcjulian29/go-toolbox/filesystem"
	"github.com/dcjulian29/go-toolbox/textformat"
)

// Run builds and executes a docker run command from the provided options.
// When opts.Capture is true the container's stdout is returned as a trimmed
// string; otherwise stdout and stderr are streamed to the terminal and an
// empty string is returned.
func Run(opts ContainerOptions) (string, error) {
	if opts.Image == textformat.EmptyString {
		return textformat.EmptyString, errors.New("image is required")
	}

	if opts.Tag == textformat.EmptyString {
		opts.Tag = "latest"
	}

	entryArgs, entryVol, err := entryArguments(opts)
	if err != nil {
		return textformat.EmptyString, err
	}

	opts.Volumes = append(opts.Volumes, entryVol...)

	args := []string{"run"}
	args = append(args, interactiveArguments(opts)...)
	args = append(args, containerArguments(opts)...)
	args = append(args, entryArgs...)
	args = append(args, environmentArguments(opts)...)
	args = append(args, volumeArguments(opts)...)
	args = append(args, portArguments(opts)...)
	args = append(args, fmt.Sprintf("%s:%s", opts.Image, opts.Tag))
	args = append(args, commandArguments(opts)...)

	if opts.Capture {
		return execute.ExternalProgramCapture("docker", args...)
	}

	return textformat.EmptyString, execute.ExternalProgram("docker", args...)
}

func commandArguments(opts ContainerOptions) []string {
	var args []string

	if opts.Command != textformat.EmptyString {
		args = append(args, opts.Command)
	}

	if len(opts.AdditionalArgs) > 0 { //nolint:revive
		args = append(args, opts.AdditionalArgs...)
	}

	return args
}

func containerArguments(opts ContainerOptions) []string {
	var args []string

	if !opts.Keep {
		args = append(args, "--rm")
	}

	if opts.HostName != textformat.EmptyString {
		args = append(args, "--hostname", opts.HostName)
	}

	if opts.Name != textformat.EmptyString {
		args = append(args, "--name", opts.Name)
	}

	if opts.ReadOnly {
		args = append(args, "--read-only")
	}

	if opts.User != textformat.EmptyString {
		args = append(args, "--user", opts.User)
	}

	if opts.WorkingDirectory != textformat.EmptyString {
		args = append(args, "--workdir", opts.WorkingDirectory)
	}

	return args
}

func entryArguments(opts ContainerOptions) (args []string, entryVolume []string, err error) {
	if opts.EntryPoint != textformat.EmptyString {
		return []string{"--entrypoint", opts.EntryPoint}, nil, nil
	}

	if opts.EntryScript == textformat.EmptyString {
		return nil, nil, nil
	}

	abs, err := filepath.Abs(opts.EntryScript)
	if err != nil {
		return nil, nil, fmt.Errorf("resolving absolute path of EntryScript: %w", err)
	}

	info, err := os.Stat(abs)

	if err != nil {
		return nil, nil, err
	}

	if info.IsDir() {
		return nil, nil, errors.New("path is a directory not a file")
	}

	filename := filepath.Base(opts.EntryScript)
	mountPoint := filesystem.EnsurePathIsUnix(abs)
	volume := mountPoint + ":/bin/" + filename

	return []string{"--entrypoint", "/bin/" + filename}, []string{volume}, nil
}

func environmentArguments(opts ContainerOptions) []string {
	args := make([]string, 0, len(opts.EnvironmentVariables)*2) //nolint

	for k, v := range opts.EnvironmentVariables {
		args = append(args, "-e", fmt.Sprintf("%s=%s", k, v))
	}

	return args
}

func interactiveArguments(opts ContainerOptions) []string {
	if !opts.Interactive {
		return []string{"--detach"}
	}

	if opts.NoTty {
		return []string{"--interactive"}
	}

	return []string{"--interactive", "--tty"}
}

func portArguments(opts ContainerOptions) []string {
	args := make([]string, 0, len(opts.Ports)*2) //nolint

	for _, port := range opts.Ports {
		args = append(args, "-p", port)
	}

	return args
}

func volumeArguments(opts ContainerOptions) []string {
	args := make([]string, 0, len(opts.Volumes)*2) //nolint

	for _, vol := range opts.Volumes {
		args = append(args, "--volume", vol)
	}

	return args
}
