package docker

import (
	"fmt"

	"github.com/dcjulian29/go-toolbox/execute"
)

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

// RunCapture runs the docker container with the given parameters
// and captures its standard output returning the result as a string.
// If there is output to standard error, return that as an error.
func RunCapture(image, tag, envPrefix string) (string, error) {
	args := []string{
		"run",
		"--rm",
	}

	args = append(args, EnvironmentVariables(envPrefix)...)
	args = append(args, HostAndWorkArguments()...)
	args = append(args, fmt.Sprintf("%s:%s", image, tag))
	args = append(args, NormalizeArguments()...)

	return execute.ExternalProgramCapture("docker", args...)
}
