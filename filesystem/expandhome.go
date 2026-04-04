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
	"os"
	"path/filepath"
)

// ExpandHome prefixes the provided path to the user's home directory as defined
// by the operating system.
func ExpandHome(path string) string {
	if len(path) >= 2 {
		if path[:2] == "~\\" || path[:2] == "~/" {
			home, err := os.UserHomeDir()
			if err != nil {
				return path
			}
			return filepath.Join(home, path[2:])
		}
	}

	return path
}
