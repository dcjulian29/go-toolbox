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

package docker

import (
	"github.com/dcjulian29/go-toolbox/filesystem"
)

// NormalizeArguments normalizes arguments to replace Windows backslash
// path separators with forward slashes for use inside the Linux container.
// Deprecated: Use github.com/dcjulian29/go-toolbox/filesystem
// EnsureUnixPathArguments function instead.
func NormalizeArguments() []string {
	return filesystem.EnsureUnixPathArguments()
}
