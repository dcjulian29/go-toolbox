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
	"os"
	"path/filepath"
	"strings"

	"github.com/dcjulian29/go-toolbox/filesystem"
)

// InjectStartScript injects the startup script to the VHDX file
func InjectStartScript(cfg *InjectConfig) error {
	if cfg.StartScript != "" && filesystem.FileExists(cfg.StartScript) {
		raw, err := os.ReadFile(cfg.StartScript)
		if err != nil {
			return fmt.Errorf("error reading start script: %w", err)
		}

		content := strings.ReplaceAll(string(raw), "{{INSTALLPACKAGE}}", cfg.InstallPackage)

		dest := filepath.Join(cfg.MountedDrive, "Windows", "Setup", "Scripts", filepath.Base(cfg.StartScript))

		if err := filesystem.EnsureDirectoryExist(filepath.Dir(dest)); err != nil {
			return fmt.Errorf("creating Scripts dir: %w", err)
		}

		if err := filesystem.EnsureFileExist(dest, []byte(content)); err != nil {
			return fmt.Errorf("writing startup script: %w", err)
		}
	}

	return nil
}
