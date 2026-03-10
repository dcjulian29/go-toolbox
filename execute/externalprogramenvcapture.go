package execute

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

import (
	"os"
	"os/exec"
)

// ExternalProgramEnvCapture runs the specified external program with the given parameters
// and additional environment variables, capturing and returning its combined output as a string.
func ExternalProgramEnvCapture(program string, env []string, params ...string) (string, error) {
	cmd := exec.Command(program, params...)
	cmd.Stdin = os.Stdin
	cmd.Env = append(os.Environ(), env...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(output[:]), nil
}
