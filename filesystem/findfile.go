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

// FindFile searches the current working directory and its children for a file
// that matches the provided filename.
func FindFile(filename string) (string, error) {
	absStart, err := filepath.Abs(".")
	if err != nil {
		return textformat.EmptyString, fmt.Errorf("failed to resolve current directory: %w", err)
	}

	visited := make(map[string]bool)

	found, err := searchChildren(filename, absStart, visited)

	if err != nil {
		return textformat.EmptyString, fmt.Errorf("failed to find '%s': %w", filename, err)
	}

	return found, nil
}

func searchChildren(filename, dir string, visited map[string]bool) (string, error) {
	realDir, err := filepath.EvalSymlinks(dir)
	if err != nil {
		return textformat.EmptyString, fmt.Errorf("failed to resolve path %s: %w", dir, err)
	}

	if visited[realDir] {
		return textformat.EmptyString, errors.New("directory already visited")
	}

	visited[realDir] = true

	entries, err := os.ReadDir(dir)
	if err != nil {
		return textformat.EmptyString, fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	for _, entry := range entries {
		fullPath := filepath.Join(dir, entry.Name())

		isDir := entry.IsDir()

		if entry.Type()&os.ModeSymlink != 0 { //nolint
			info, err := os.Stat(fullPath)
			if err != nil {
				continue
			}

			isDir = info.IsDir()
		}

		if isDir {
			if found, err := searchChildren(filename, fullPath, visited); err == nil {
				return found, nil
			}
		} else if entry.Name() == filename {
			return fullPath, nil
		}
	}

	return textformat.EmptyString, errors.New("file not found in children")
}
