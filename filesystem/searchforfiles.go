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

package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
)

// SearchForFiles searches 'path' for files whose name matches the given
// glob pattern. Directories are excluded from results. Returns an error
// if the path is inaccessible or nothing is found.
func SearchForFiles(path, pattern string) ([]string, error) {
	if _, err := os.Stat(path); err != nil {
		return []string{}, fmt.Errorf("failed to access %q: %w", path, err)
	}

	matches, err := filepath.Glob(filepath.Join(path, pattern))
	if err != nil {
		return []string{}, err
	}

	files := make([]string, 0, len(matches)) //nolint:revive

	for _, match := range matches {
		info, err := os.Stat(match)
		if err != nil {
			continue
		}

		if !info.IsDir() {
			files = append(files, match)
		}
	}

	if len(files) == 0 { //nolint:revive
		return []string{}, fmt.Errorf("no file matching %q found in %q", pattern, path)
	}

	return files, nil
}
