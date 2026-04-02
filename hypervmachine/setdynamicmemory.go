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

// SetDynamicMemory enables dynamic memory with the given range.
// The startup memory is set equal to minBytes.
func SetDynamicMemory(name string, minBytes, maxBytes int64) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("virtual machine name must not be empty")
	}

	if !Exist(name) {
		return errors.New("virtual machine does not exist")
	}

	if minBytes <= 0 { //nolint:revive
		return errors.New("minimum bytes must be greater than zero")
	}

	if maxBytes <= minBytes {
		return errors.New("maximum memory bytes must be greater than minimum memory bytes")
	}

	script := fmt.Sprintf(
		`Set-VMMemory -VMName "%s" -DynamicMemoryEnabled $true `+
			`-StartupBytes %d -MinimumBytes %d -MaximumBytes %d -ErrorAction Stop`,
		textformat.EscapeForPowerShell(name), minBytes, minBytes, maxBytes,
	)

	if err := execute.RunPowerShell(script); err != nil {
		return fmt.Errorf("setting dynamic memory for VM %q: %w", name, err)
	}

	return nil
}
