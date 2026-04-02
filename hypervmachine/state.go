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

// State returns the current state of a virtual machine ("Running", "Off", "Saved", ...).
func State(name string) (string, error) {
	if strings.TrimSpace(name) == "" {
		return textformat.EmptyString, errors.New("virtual machine name must not be empty")
	}

	if !Exist(name) {
		return textformat.EmptyString, errors.New("virtual machine does not exist")
	}

	script := fmt.Sprintf(`(Get-VM -VMName '%s').State`, textformat.EscapeForPowerShell(name))

	return execute.RunPowerShellCapture(script)
}
