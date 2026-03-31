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
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ExternalProgramEnvCapture runs the specified external program with the given
// parameters and additional environment variables, capturing its standard output
// and returning the result as a trimmed string. If the command fails, stderr
// output is included in the error. Standard error output is discarded on success.
func ExternalProgramEnvCapture(program string, env []string, params ...string) (string, error) {
	var out, errBuf bytes.Buffer

	cmd := exec.Command(program, params...) // #nosec G204
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	cmd.Env = append(os.Environ(), env...)

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%w\n%s", err, errBuf.String())
	}

	return strings.TrimSpace(out.String()), nil
}
