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

// ExternalProgramCapture runs the specified external program with the given parameters
// and captures its combined standard output and standard error, returning the result as a string.
func ExternalProgramCapture(program string, params ...string) (string, error) {
	cmd := exec.Command(program, params...)
	cmd.Stdin = os.Stdin
	out, err := cmd.CombinedOutput()

	if err != nil {
		return "", err
	}

	return string(out), nil
}
