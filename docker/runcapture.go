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
	"strings"

	"github.com/dcjulian29/go-toolbox/filesystem"
)

// RunCapture runs the docker container with the given parameters
// and captures its standard output returning the result as a string.
// If there is output to standard error, return that as an error.
// Deprecated: Use github.com/dcjulian29/go-toolbox/docker
// Run function and ContainerOptions instead.
func RunCapture(image, tag, envPrefix string) (string, error) {
	opts := ContainerOptions{
		Keep:                 false,
		EnvironmentVariables: EnvironmentVariablesWithPrefix(envPrefix),
		AdditionalArgs:       strings.Join(HostAndWorkArguments(), " "),
		Image:                image,
		Tag:                  tag,
		Command:              strings.Join(filesystem.EnsureUnixPathArguments(), " "),
		Capture:              true,
	}

	return Run(opts)
}
