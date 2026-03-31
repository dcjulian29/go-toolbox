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
	"io/fs"
	"path/filepath"
	"strings"
)

// ScanDirectory recursively traverses the specified directory and returns
// separate lists of directories and regular files found within it. Paths
// containing any of the ignore strings are excluded, and ignored directories
// are not descended into.
func ScanDirectory(path string, ignore []string) ([]string, []string, error) {
	folders := []string{}
	files := []string{}

	err := filepath.WalkDir(path, func(entry string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		for _, i := range ignore {
			if strings.Contains(entry, i) {
				if d.IsDir() {
					return filepath.SkipDir
				}

				return nil
			}
		}

		if d.IsDir() {
			folders = append(folders, entry)
		} else if d.Type().IsRegular() {
			files = append(files, entry)
		}

		return nil
	})

	if err != nil {
		return []string{}, []string{}, err
	}

	return folders, files, nil
}
