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

package hyperv

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dcjulian29/go-toolbox/execute"
	"github.com/dcjulian29/go-toolbox/textformat"
)

// OpenConsole opens the Hyper-V Virtual Machine Connection console for the
// named virtual machine.
func OpenConsole(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("virtual machine name must not be empty")
	}

	script := fmt.Sprintf(`vmconnect.exe localhost "%s"`, textformat.EscapeForPowerShell(name))

	return execute.RunPowerShell(script)
}
