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
	"fmt"
	"path/filepath"
)

// SearchForFile searches 'path' for the last file whose name matches the given
// glob pattern. Returns an error if nothing is found.
func SearchForFile(path, pattern string) (string, error) {
	matches, err := filepath.Glob(filepath.Join(path, pattern))
	if err != nil {
		return "", err
	}

	if len(matches) == 0 { //nolint:revive
		return "", fmt.Errorf("no file matching '%q' found in '%s'", pattern, path)
	}

	return matches[len(matches)-1], nil //nolint:revive
}
