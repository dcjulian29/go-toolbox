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
	"fmt"
	"os"
)

// EnsureDirectoryExist verifies that a directory exists at the given path,
// creating it and any necessary parent directories if it does not exist.
func EnsureDirectoryExist(path string) error {
	info, err := os.Stat(path)

	if err == nil {
		if !info.IsDir() {
			return fmt.Errorf("'%s' exists but is not a directory", path)
		}

		return nil
	}

	if !os.IsNotExist(err) {
		return err
	}

	return os.MkdirAll(path, FileModeExecutable)
}
