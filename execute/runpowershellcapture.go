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
	"errors"
	"strings"
)

// RunPowerShellCapture executes a PowerShell command and captures its standard
// output, returning the result as a trimmed string. If the command fails,
// stderr output is included in the error. Standard error output is discarded
// on success.
func RunPowerShellCapture(command string) (string, error) {
	if strings.TrimSpace(command) == "" {
		return "", errors.New("command must not be empty")
	}

	shell, err := findPowerShell()
	if err != nil {
		return "", err //nolint:revive
	}

	return ExternalProgramCapture(shell, "-NoProfile", "-NonInteractive",
		"-ExecutionPolicy", "Bypass", "-Command", command)
}
