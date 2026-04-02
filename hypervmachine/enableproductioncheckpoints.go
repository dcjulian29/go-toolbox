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

// EnableProductionCheckpoints enables production checkpoints for the VM.
// Production checkpoints use Volume Shadow Copy Service (VSS) inside the
// guest to create application-consistent backups without capturing memory
// state. These are safe for domain controllers and production workloads
// but require guest integration services to be running.
func EnableProductionCheckpoints(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("virtual machine name must not be empty")
	}

	script := fmt.Sprintf(
		`Set-VM -Name "%s" -CheckpointType Production -ErrorAction Stop`,
		textformat.EscapeForPowerShell(name),
	)

	if err := execute.RunPowerShell(script); err != nil {
		return fmt.Errorf("enabling production checkpoints for VM %q: %w", name, err)
	}

	return nil
}
