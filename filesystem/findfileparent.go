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
)

func FindFileParent(filename string) (string, error) {
	pwd, _ := os.Getwd()
	absStart, err := filepath.Abs(pwd)
	if err != nil {
		return EmptyString, fmt.Errorf("failed to resolve current directory: %w", err)
	}

	current := filepath.Dir(absStart)

	for {
		candidate := filepath.Join(current, filename)

		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}

		parent := filepath.Dir(current)

		if parent == current {
			return EmptyString, errors.New("file not found: no more parent directories to search")
		}

		current = parent
	}
}
