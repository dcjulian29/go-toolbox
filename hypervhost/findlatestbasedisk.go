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
	"path/filepath"
	"slices"
	"strings"

	"github.com/dcjulian29/go-toolbox/textformat"
)

// FindLatestBaseDisk searches directoryPath for a VHDX file whose name starts with
// the given prefix and returns the full path of the alphabetically last match
// (which is usually the newest dated image).
func FindLatestBaseDisk(directoryPath, pattern string) (string, error) {
	if strings.TrimSpace(directoryPath) == "" {
		return textformat.EmptyString, errors.New("directory path must not be empty")
	}

	if strings.TrimSpace(pattern) == "" {
		return textformat.EmptyString, errors.New("pattern must not be empty")
	}

	matches, err := filepath.Glob(filepath.Join(directoryPath, pattern+"*.vhdx"))
	if err != nil {
		return textformat.EmptyString, err
	}

	if len(matches) == 0 { //nolint:revive
		return textformat.EmptyString, fmt.Errorf("no base disk matching %q found in %q", pattern, directoryPath)
	}

	return slices.Max(matches), nil
}
