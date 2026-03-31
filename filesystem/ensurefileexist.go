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

// EnsureFileExist creates a file at the specified path and writes the provided content to it.
// If the file already exists, it is truncated and overwritten.
func EnsureFileExist(path string, content []byte) error {
	if err := EnsureDirectoryExist(filepath.Dir(path)); err != nil {
		return err
	}

	return os.WriteFile(path, content, FileModeReadable)
}
