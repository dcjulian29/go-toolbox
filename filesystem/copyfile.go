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
	"io"
	"os"
)

// CopyFile duplicates a file from the source path to the destination path,
// preserving the original file permissions and contents.
func CopyFile(src, dst string) error {
	if !FileExists(src) {
		return fmt.Errorf("source file '%s' does not exists", src)
	}

	if FileExists(dst) {
		if err := os.Remove(dst); err != nil {
			return err
		}
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}

	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}

	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return err
	}

	return nil
}
