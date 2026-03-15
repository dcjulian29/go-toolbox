//go:build windows

package host

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
	"path/filepath"

	"github.com/dcjulian29/go-toolbox/filesystem"
)

// FindLatestBaseDisk searches directoryPath for a VHDX whose name matches
// the pattern and returns the full path of the alphabetically last match
// (which is usually the newest dated image).
func FindLatestBaseDisk(directoryPath, pattern string) (string, error) {
	matches, err := filepath.Glob(filepath.Join(directoryPath, pattern))
	if err != nil {
		return "", err //nolint
	}

	if len(matches) == 0 { //nolint
		return "", fmt.Errorf("no base disk matching '%q' found in '%s'", pattern, directoryPath) //nolint
	}

	latest := matches[0] //nolint

	for _, m := range matches[1:] { //nolint
		if m > latest {
			latest = m
		}
	}

	if !filesystem.FileExists(latest) {
		return "", fmt.Errorf("base disk %q not found", latest) //nolint
	}

	return latest, nil
}
