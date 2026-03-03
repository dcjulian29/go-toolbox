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
)

func FindFilesByExtension(path, extension string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(path, func(f string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(d.Name()) == extension {
			files = append(files, f)
		}

		return nil
	})

	if err != nil {
		return []string{}, err
	}

	return files, nil
}
