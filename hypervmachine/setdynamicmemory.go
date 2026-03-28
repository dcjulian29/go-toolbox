//go:build windows

package hypervmachine

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

// SetDynamicMemory enables dynamic memory with the given startup/min/max values.
func SetDynamicMemory(name string, startBytes, minBytes, maxBytes int64) error {
	script := fmt.Sprintf(
		`Set-VMMemory -VMName "%s" -DynamicMemoryEnabled $true `+
			`-StartupBytes %d -MinimumBytes %d -MaximumBytes %d -ErrorAction Stop`,
		textformat.EscapeForPowerShell(name), startBytes, minBytes, maxBytes,
	)

	return execute.RunPowershell(script)
}
