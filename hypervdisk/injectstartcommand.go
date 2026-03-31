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
	"path/filepath"
	"strings"

	"github.com/dcjulian29/go-toolbox/filesystem"
)

const prefix = `%WINDIR%\System32\WindowsPowerShell\v1.0\powershell.exe ` +
	`-NoProfile -NonInteractive -ExecutionPolicy Bypass -NoLogo -Command ` +
	`%WINDIR%\Setup\Scripts\`

// InjectStartCommand writes a SetupComplete.cmd file to the mounted VHDX that
// executes the configured start script via PowerShell after the first login.
func InjectStartCommand(cfg *InjectConfig) error {
	if cfg == nil {
		return errors.New("inject config must not be nil")
	}

	if strings.TrimSpace(cfg.MountedDrive) == "" {
		return errors.New("mounted drive must not be empty")
	}

	if strings.TrimSpace(cfg.StartScript) == "" {
		return errors.New("start script must not be empty")
	}

	path := filepath.Join(cfg.MountedDrive, "Windows", "Setup", "Scripts", "SetupComplete.cmd")

	content := prefix + filepath.Base(cfg.StartScript)

	return filesystem.EnsureFileExist(path, []byte(content))
}
