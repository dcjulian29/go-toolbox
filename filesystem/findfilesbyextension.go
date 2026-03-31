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

// FindFilesByExtension recursively searches a directory for all files matching
// the given file extension. The extension may be specified with or without a
// leading dot (e.g., ".txt" or "txt"). Directories are excluded from results.
func FindFilesByExtension(path, extension string) ([]string, error) {
	if !strings.HasPrefix(extension, ".") {
		extension = "." + extension
	}

	files := []string{}

	err := filepath.WalkDir(path, func(entry string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if !d.IsDir() && filepath.Ext(d.Name()) == extension {
			files = append(files, entry)
		}

		return nil
	})

	if err != nil {
		return []string{}, err
	}

	return files, nil
}
