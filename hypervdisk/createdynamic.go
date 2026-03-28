//go:build windows

package hypervdisk

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

// CreateDynamic creates a new dynamic VHDX with the given size.
func CreateDynamic(vhdxPath string, sizeBytes int64) error {
	script := fmt.Sprintf(
		`New-VHD -Path "%s" -SizeBytes %d -Dynamic -ErrorAction Stop`,
		textformat.EscapeForPowerShell(vhdxPath), sizeBytes,
	)

	if err := execute.RunPowershell(script); err != nil {
		return fmt.Errorf("creating dynamic VHDX: %w", err)
	}

	return nil
}
