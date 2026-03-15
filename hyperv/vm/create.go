//go:build windows

package hvvm

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
	"fmt"

	"github.com/dcjulian29/go-toolbox/execute"
	"github.com/dcjulian29/go-toolbox/textformat"
)

// Create creates a new Hyper-V virtual machine with the provided config.
func Create(cfg Config) error {
	if cfg.Generation == 0 { //nolint
		cfg.Generation = 2 //nolint
	}

	script := fmt.Sprintf(
		`New-VM -Generation %d -Name "%s" -MemoryStartupBytes %d `+
			`-SwitchName "%s"-VHDPath "%s" -ErrorAction Stop`,
		cfg.Generation,
		textformat.EscapeForPowershell(cfg.Name),
		cfg.MemoryBytes,
		textformat.EscapeForPowershell(cfg.VirtualSwitch),
		textformat.EscapeForPowershell(cfg.VHDXPath),
	)

	if err := execute.RunPowershell(script); err != nil {
		return fmt.Errorf("creating VM %q: %w", cfg.Name, err)
	}

	return nil
}
