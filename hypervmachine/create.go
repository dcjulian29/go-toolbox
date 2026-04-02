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
	"fmt"

	"github.com/dcjulian29/go-toolbox/execute"
	"github.com/dcjulian29/go-toolbox/textformat"
)

// Create creates a new Hyper-V virtual machine with the provided config.
// If post-creation configuration fails, the VM may exist in a partially
// configured state and should be removed manually.
func Create(cfg *Config) error {
	if err := ValidateConfig(cfg); err != nil {
		return err
	}

	script := fmt.Sprintf(
		`New-VM -Generation %d -Name "%s" -MemoryStartupBytes %d `+
			`-SwitchName "%s" -VHDPath "%s" -ErrorAction Stop`,
		cfg.Generation,
		textformat.EscapeForPowerShell(cfg.Name),
		cfg.MemoryBytes,
		textformat.EscapeForPowerShell(cfg.VirtualSwitch),
		textformat.EscapeForPowerShell(cfg.VHDXPath),
	)

	if err := execute.RunPowerShell(script); err != nil {
		return fmt.Errorf("creating VM %q: %w", cfg.Name, err)
	}

	if cfg.MaximumMemoryBytes != cfg.MemoryBytes {
		if err := SetDynamicMemory(cfg.Name, cfg.MemoryBytes, cfg.MaximumMemoryBytes); err != nil {
			return fmt.Errorf("configuring dynamic memory for VM %q: %w", cfg.Name, err)
		}
	}

	if cfg.ProcessorCount > 0 { //nolint:revive
		if err := SetProcessorCount(cfg.Name, cfg.ProcessorCount); err != nil {
			return fmt.Errorf("setting processor count for VM %q: %w", cfg.Name, err)
		}
	}

	if cfg.Generation == GenerationV2 && !cfg.SecureBoot {
		if err := DisableSecureBoot(cfg.Name); err != nil {
			return fmt.Errorf("disabling Secure Boot for VM %q: %w", cfg.Name, err)
		}
	}

	return nil
}
