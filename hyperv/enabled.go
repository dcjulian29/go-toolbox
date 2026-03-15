//go:build windows

package hyperv

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
)

// Enabled returns an error if the Hyper-V role is not available.
func Enabled() error {
	out, err := execute.RunPowershellCapture(
		`(Get-WindowsOptionalFeature -Online -FeatureName Microsoft-Hyper-V).State`,
	)
	if err != nil {
		return fmt.Errorf("could not query Hyper-V feature state: %w", err)
	}

	if out != "Enabled" {
		return fmt.Errorf("the Hyper-V feature is not enabled on this host (state: %s)", out)
	}

	return nil
}
