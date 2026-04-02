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

// AttachDVD attaches an ISO to the VM's DVD drive.
func AttachDVD(name, path string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("name of virtual machine must not be empty")
	}

	if strings.TrimSpace(path) == "" {
		return errors.New("path to ISO file must not be empty")
	}

	script := fmt.Sprintf(
		`Add-VMDvdDrive -VMName "%s" -Path "%s" -ErrorAction Stop`,
		textformat.EscapeForPowerShell(name),
		textformat.EscapeForPowerShell(path),
	)

	if err := execute.RunPowerShell(script); err != nil {
		return fmt.Errorf("attaching DVD to VM %q: %w", name, err)
	}

	return nil
}
