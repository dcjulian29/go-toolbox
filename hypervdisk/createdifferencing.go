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

package hypervdisk

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dcjulian29/go-toolbox/execute"
	"github.com/dcjulian29/go-toolbox/textformat"
)

// CreateDifferencing creates a differencing VHDX disk at vhdxPath linked to the
// parent disk at referencePath.
func CreateDifferencing(referencePath, vhdxPath string) error {
	if strings.TrimSpace(referencePath) == "" {
		return errors.New("reference path must not be empty")
	}

	if strings.TrimSpace(vhdxPath) == "" {
		return errors.New("VHDX path must not be empty")
	}

	script := fmt.Sprintf(
		`New-VHD -ParentPath "%s" -Path "%s" -Differencing -ErrorAction Stop`,
		textformat.EscapeForPowerShell(referencePath),
		textformat.EscapeForPowerShell(vhdxPath),
	)

	if err := execute.RunPowerShell(script); err != nil {
		return fmt.Errorf("creating differencing VHDX: %w", err)
	}

	return nil
}
