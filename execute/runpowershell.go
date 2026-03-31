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
	"os/exec"
	"runtime"
	"strings"
)

// RunPowerShell executes a PowerShell command and streams stdout/stderr to the
// caller's terminal. It prefers pwsh (PowerShell Core) and falls back to
// powershell on Windows. An error is returned if PowerShell is not available
// or the command exits with a non-zero status.
func RunPowerShell(command string) error {
	if strings.TrimSpace(command) == "" {
		return errors.New("command must not be empty")
	}

	shell, err := findPowerShell()
	if err != nil {
		return err
	}

	return ExternalProgram(shell, "-NoProfile", "-NonInteractive",
		"-ExecutionPolicy", "Bypass", "-Command", command)
}

func findPowerShell() (string, error) {
	if path, err := exec.LookPath("pwsh"); err == nil {
		return path, nil
	}

	if runtime.GOOS == "windows" {
		if path, err := exec.LookPath("powershell"); err == nil {
			return path, nil
		}
	}

	return "", errors.New("PowerShell is not installed or available")
}
