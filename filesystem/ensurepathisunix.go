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
	"regexp"
	"strings"
)

// driveRegex matches a word boundary (\b), a single letter, a colon,
// and a slash or backslash. This ensures it matches " C:\" or "=D:/"
// but ignores "http://" or "localhost:8080/"
var driveRegex = regexp.MustCompile(`\b([a-zA-Z]):[\\/]`)

// EnsurePathIsUnix normalizes the provided path to replace Windows backslash
// path separators with forward slashes for use inside the Linux container.
func EnsurePathIsUnix(path string) string {
	if !driveRegex.MatchString(path) {
		return path
	}

	path = strings.ReplaceAll(path, "\\", "/")

	path = driveRegex.ReplaceAllStringFunc(path, func(match string) string {
		driveLetter := match[0:1] //nolint:revive

		return "/" + driveLetter + "/"
	})

	return path
}
