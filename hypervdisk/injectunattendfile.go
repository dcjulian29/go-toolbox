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
	"os"
	"path/filepath"
	"strings"

	"github.com/dcjulian29/go-toolbox/filesystem"
	"github.com/dcjulian29/go-toolbox/textformat"
)

// InjectUnattendFile injects the unattend XML to the VHDX file.
func InjectUnattendFile(cfg *InjectConfig) error {
	if cfg == nil {
		return errors.New("inject config must not be nil")
	}

	if strings.TrimSpace(cfg.UnattendTemplate) == "" {
		return errors.New("unattend template must not be empty")
	}

	if !filesystem.FileExist(cfg.UnattendTemplate) {
		return fmt.Errorf("unattend template not found: %s", cfg.UnattendTemplate)
	}

	if strings.TrimSpace(cfg.MountedDrive) == "" {
		return errors.New("mounted drive must not be empty")
	}

	raw, err := os.ReadFile(cfg.UnattendTemplate)
	if err != nil {
		return fmt.Errorf("reading unattend template: %w", err)
	}

	if len(raw) == 0 { //nolint:revive
		return errors.New("unattend template file must not be empty")
	}

	replaced, err := templateReplacements(raw, cfg)
	if err != nil {
		return err
	}

	dest := filepath.Join(cfg.MountedDrive, "unattend.xml")

	if err := filesystem.EnsureFileExist(dest, replaced); err != nil {
		return fmt.Errorf("writing unattend.xml: %w", err)
	}

	return nil
}

func templateReplacements(raw []byte, cfg *InjectConfig) ([]byte, error) {
	var err error

	content := string(raw)

	content, err = replaceTemplateMarker(content, "{{COMPUTERNAME}}", cfg.ComputerName, "computer name")
	if err != nil {
		return nil, err
	}

	content, err = replaceTemplateMarker(content, "{{PASSWORD}}", cfg.UserPassword, "user password")
	if err != nil {
		return nil, err
	}

	content, err = replaceTemplateMarker(content, "{{USER}}", cfg.UserName, "user name")
	if err != nil {
		return nil, err
	}

	return []byte(content), nil
}

func replaceTemplateMarker(content, marker, value, fieldName string) (string, error) {
	if !strings.Contains(content, marker) {
		return content, nil
	}

	if strings.TrimSpace(value) == "" { //nolint:revive
		return "", fmt.Errorf("template contains %s but %s is empty", marker, fieldName)
	}

	return strings.ReplaceAll(content, marker, textformat.XMLEscape(value)), nil
}
