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
)

// AppendFile appends the provided content to the file at the specified path.
// The file must already exist.
func AppendFile(path string, content []byte) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, ModeOwnerReadWrite)
	if err != nil {
		return err
	}

	_, err = file.Write(content)

	if closeErr := file.Close(); err == nil {
		err = closeErr
	}

	return err
}
