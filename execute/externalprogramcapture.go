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
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/dcjulian29/go-toolbox/textformat"
)

// ExternalProgramCapture runs the docker container with the given parameters
// and captures its standard output returning the result as a string.
// If there is output to standard error, return that as an error.
func ExternalProgramCapture(program string, params ...string) (string, error) {
	var out, errBuf bytes.Buffer

	cmd := exec.Command(program, params...)
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return textformat.EmptyString, fmt.Errorf("%w\n%s", err, errBuf.String())
	}

	return strings.TrimSpace(out.String()), nil
}
