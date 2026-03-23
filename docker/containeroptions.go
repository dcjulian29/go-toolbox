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

// ContainerOptions holds all configuration options for creating a Docker container.
type ContainerOptions struct {
	// EnvironmentVariables is a map of environment variables to set inside the
	// container, equivalent to docker run -e KEY=VALUE.
	EnvironmentVariables map[string]string

	// Name assigns a name to the container, equivalent to docker run --name.
	// If empty, Docker will generate a random name.
	Name string

	// User sets the username or UID used when running the container, equivalent
	// to docker run --user. Accepts the format <name|uid>[:<group|gid>].
	User string

	// Command overrides the default command (CMD) defined in the image.
	Command string

	// EntryPoint overrides the default entry point (ENTRYPOINT) of the image,
	// equivalent to docker run --entrypoint. Takes precedence over EntryScript
	// if both are set.
	EntryPoint string

	// EntryScript overrides the default entry point of the image by mounting
	// the specified host script file into the container at /bin/<filename> and
	// using it as the entry point. Ignored if EntryPoint is also set.
	EntryScript string

	// HostName sets the hostname of the container, equivalent to
	// docker run --hostname.
	HostName string

	// AdditionalArgs contains any extra raw arguments to append to the
	// docker run command that are not covered by the other fields.
	// Arguments are split on whitespace before being passed to the CLI.
	AdditionalArgs string

	// Image is the Docker image to run. Required.
	Image string

	// Tag is the image tag to use, equivalent to the :<tag> suffix in
	// docker run <image>:<tag>. Defaults to "latest" if empty.
	Tag string

	// WorkingDirectory sets the working directory inside the container,
	// equivalent to docker run --workdir.
	WorkingDirectory string

	// Ports is a list of port mappings to publish from the container to the
	// host, equivalent to docker run -p. Each entry should be in Docker's
	// port mapping format, e.g. "8080:80" or "127.0.0.1:8080:80".
	Ports []string

	// Volumes is a list of volume mounts to bind into the container, equivalent
	// to docker run --volume. Each entry should be in Docker's volume format,
	// e.g. "/host/path:/container/path" or "/host/path:/container/path:ro".
	Volumes []string

	// Capture redirects the container's stdout to the caller as a returned
	// string instead of inheriting the parent process's stdout.
	Capture bool

	// Interactive keeps stdin open even if not attached, equivalent to
	// docker run --interactive. When true the container runs in the foreground.
	// When false the container is started in the background with --detach.
	Interactive bool

	// Keep retains the container and its associated anonymous volumes after it
	// exits. When false (the default), --rm is passed and the container is
	// automatically removed on exit.
	Keep bool

	// NoTty disables pseudo-TTY allocation, equivalent to omitting
	// docker run --tty. Only applies when Interactive is true.
	NoTty bool

	// ReadOnly mounts the container's root filesystem as read-only, equivalent
	// to docker run --read-only.
	ReadOnly bool
}
