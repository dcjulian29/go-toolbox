package docker

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dcjulian29/go-toolbox/execute"
	"github.com/dcjulian29/go-toolbox/filesystem"
	"github.com/dcjulian29/go-toolbox/textformat"
)

// Run builds and executes a `docker run` command based on the
// provided options. It returns the combined output and any error encountered.
func Run(opts ContainerOptions) (string, error) {
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

	return textformat.EmptyString, execute.ExternalProgram("docker", args...) //nolint
}

func commandArguments(opts ContainerOptions) []string {
	var args []string

	if opts.Command != textformat.EmptyString {
		args = append(args, opts.Command)
	}

	if opts.AdditionalArgs != textformat.EmptyString {
		args = append(args, strings.Fields(opts.AdditionalArgs)...)
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
		args = append(args, "--user="+opts.User)
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
		return nil, nil, fmt.Errorf("resolving absolute directory of EntryScript: %w", err)
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
