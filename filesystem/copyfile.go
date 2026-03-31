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
	"io"
	"os"
)

// CopyFile duplicates a file from the source path to the destination path,
// preserving the original file permissions and contents. If the destination
// file already exists, it is overwritten.
func CopyFile(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("source file %q: %w", src, err)
	}

	if info.IsDir() {
		return fmt.Errorf("source %q is a directory, not a file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}

	defer source.Close()

	destination, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}

	_, err = io.Copy(destination, source)

	if closeErr := destination.Close(); err == nil {
		err = closeErr
	}

	return err
}
