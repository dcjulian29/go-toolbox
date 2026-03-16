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
	"github.com/dcjulian29/go-toolbox/textformat"
)

// InjectUnattendFile injects the unattend XML to the VHDX file.
func InjectUnattendFile(cfg *InjectConfig) error {
	if cfg.UnattendTemplate != "" && filesystem.FileExists(cfg.UnattendTemplate) {
		raw, err := os.ReadFile(cfg.UnattendTemplate)
		if err != nil {
			return fmt.Errorf("error reading unattend template: %w", err)
		}

		content := string(raw)
		content = strings.ReplaceAll(content, "{{COMPUTERNAME}}", textformat.XMLEscape(cfg.ComputerName))
		content = strings.ReplaceAll(content, "{{PASSWORD}}", textformat.XMLEscape(cfg.UserPassword))
		content = strings.ReplaceAll(content, "{{USER}}", textformat.XMLEscape(cfg.UserName))

		dest := filepath.Join(cfg.MountedDrive, "unattend.xml")

		if err := filesystem.EnsureFileExist(dest, []byte(content)); err != nil {
			return fmt.Errorf("error writing unattend.xml: %w", err)
		}

		fmt.Println(textformat.Info("[inject] unattend.xml"))
	}

	return nil
}
