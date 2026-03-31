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

	"github.com/dcjulian29/go-toolbox/textformat"
)

// SearchForFile searches 'path' for files whose name matches the given
// glob pattern and returns the last match in lexicographic order.
// Directories are excluded from results. Returns an error if the path
// is inaccessible or nothing is found.
func SearchForFile(path, pattern string) (string, error) {
	if _, err := os.Stat(path); err != nil {
		return textformat.EmptyString, fmt.Errorf("failed to access %q: %w", path, err)
	}

	matches, err := filepath.Glob(filepath.Join(path, pattern))
	if err != nil {
		return textformat.EmptyString, err
	}

	var last string

	for i := len(matches) - 1; i >= 0; i-- { //nolint:revive
		info, err := os.Stat(matches[i])
		if err != nil {
			continue
		}

		if !info.IsDir() {
			last = matches[i]

			break
		}
	}

	if last == "" {
		return textformat.EmptyString, fmt.Errorf("no file matching %q found in %q", pattern, path)
	}

	return last, nil
}
