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
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dcjulian29/go-toolbox/textformat"
)

// FindFileParent searches for a file matching the specified name starting
// from the parent of the current working directory and moving upwards
// through ancestor directories until the filesystem root is reached.
// Directories are not matched even if they share the same name.
func FindFileParent(filename string) (string, error) {
	absStart, err := filepath.Abs(".")
	if err != nil {
		return textformat.EmptyString, fmt.Errorf("failed to resolve current directory: %w", err)
	}

	current := filepath.Dir(absStart)

	for {
		candidate := filepath.Join(current, filename)

		info, err := os.Stat(candidate)
		if err == nil && !info.IsDir() {
			return candidate, nil
		}

		parent := filepath.Dir(current)

		if parent == current {
			return textformat.EmptyString, errors.New("file not found: no more parent directories to search")
		}

		current = parent
	}
}
