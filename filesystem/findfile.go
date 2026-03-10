package filesystem

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
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// FindFile searches the specified directory and its children for files that
// match the provided naming pattern or criteria.
func FindFile(filename string) (string, error) {
	pwd, _ := os.Getwd()
	absStart, err := filepath.Abs(pwd)
	if err != nil {
		return EmptyString, fmt.Errorf("failed to resolve current directory: %w", err)
	}

	if found, err := searchChildren(filename, absStart); err == nil {
		return found, nil
	}

	return EmptyString, fmt.Errorf("failed to find '%s'", filename)
}

func searchChildren(filename, dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	for _, entry := range entries {
		fullPath := filepath.Join(dir, entry.Name())

		if !entry.IsDir() {
			if entry.Name() == filename {
				return fullPath, nil
			}
		} else {
			if found, err := searchChildren(filename, fullPath); err == nil {
				return found, nil
			}
		}
	}

	return EmptyString, errors.New("file not found in children")
}
