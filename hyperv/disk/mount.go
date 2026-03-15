//go:build windows

package hvdisk

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
	"path/filepath"

	"github.com/dcjulian29/go-toolbox/execute"
	"github.com/dcjulian29/go-toolbox/textformat"
)

// Mount mounts the VHDX and returns the drive letter assigned by Windows.
func Mount(path string) (string, error) {
	script := fmt.Sprintf(
		`$v = Mount-VHD -Path "%s" -PassThru -ErrorAction Stop; `+
			`($v | Get-Disk | Get-Partition | Get-Volume).DriveLetter`,
		textformat.EscapeForPowershell(path),
	)

	letter, err := execute.RunPowershellCapture(script)
	if err != nil {
		return "", fmt.Errorf("an error occured when mounting VHDX %s: %w", filepath.Base(path), err) //nolint
	}

	if letter == "" {
		return "", fmt.Errorf("no drive letter assigned after mounting %s", filepath.Base(path)) //nolint
	}

	return letter, nil
}
