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

package hypervhost

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dcjulian29/go-toolbox/execute"
	"github.com/dcjulian29/go-toolbox/textformat"
)

// VMStoragePath returns the configured VM storage path from the Hyper-V host.
func VMStoragePath() (string, error) {
	script := `(Get-VMHost -ErrorAction Stop).VirtualHardDiskPath`
	path, err := execute.RunPowerShellCapture(script)
	if err != nil {
		return textformat.EmptyString, fmt.Errorf("retrieving default hard disk path: %w", err)
	}

	path = strings.TrimSpace(path)

	if path == "" {
		return "", errors.New("Hyper-V host returned an empty storage path")
	}

	return path, nil
}
