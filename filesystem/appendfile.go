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
	"os"
)

// AppendFile appends to a file at the specified path and writes the provided content to it.
func AppendFile(path string, content []byte) error {
	if !FileExists(path) {
		return fmt.Errorf("'%s' does not exist", path)
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644) //nolint
	if err != nil {
		return err
	}

	defer file.Close()

	if _, err = file.Write(content); err != nil {
		return err
	}

	return nil
}
