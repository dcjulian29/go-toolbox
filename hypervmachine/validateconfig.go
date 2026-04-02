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

	"github.com/dcjulian29/go-toolbox/textformat"
)

// ValidateConfig checks that all required fields in the Config are set and
// that their values are within acceptable ranges.
func ValidateConfig(cfg *Config) error {
	if cfg == nil {
		return errors.New("config must not be nil")
	}

	if strings.TrimSpace(cfg.Name) == textformat.EmptyString {
		return errors.New("virtual machine name must not be empty")
	}

	if strings.TrimSpace(cfg.VHDXPath) == textformat.EmptyString {
		return errors.New("virtual disk path must not be empty")
	}

	if strings.TrimSpace(cfg.VirtualSwitch) == textformat.EmptyString {
		return errors.New("virtual switch must not be empty")
	}

	if cfg.Generation != GenerationV1 && cfg.Generation != GenerationV2 {
		return fmt.Errorf("unsupported VM generation: %d", cfg.Generation)
	}

	if cfg.MemoryBytes <= 0 { //nolint:revive
		return errors.New("memory bytes must be greater than zero")
	}

	if cfg.MaximumMemoryBytes < cfg.MemoryBytes {
		return errors.New("maximum memory bytes must be greater than or equal to memory bytes")
	}

	if cfg.ProcessorCount <= 0 { //nolint:revive
		return errors.New("virtual machine must have one or more processors")
	}

	return nil
}
