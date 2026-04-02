//go:build windows

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

package hypervmachine

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dcjulian29/go-toolbox/execute"
	"github.com/dcjulian29/go-toolbox/textformat"
)

// SetProcessorCount sets the number of virtual processors on an existing VM.
func SetProcessorCount(name string, count int) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("virtual machine name must not be empty")
	}

	if !Exist(name) {
		return errors.New("virtual machine does not exist")
	}

	if count < 1 { //nolint:revive
		return errors.New("processor count must be greater than zero")
	}

	script := fmt.Sprintf(
		`Set-VMProcessor -VMName "%s" -Count %d -ErrorAction Stop`,
		textformat.EscapeForPowerShell(name), count,
	)

	if err := execute.RunPowerShell(script); err != nil {
		return fmt.Errorf("setting processor count for VM %q: %w", name, err)
	}

	return nil
}
