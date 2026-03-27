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

// Package filesystem contains robust utility functions for interacting with
// the file system, such as checking file existence, copying files, computing
// file hashes, and finding files.
package filesystem

import "os"

const (
	// FileModeExecutable represents the standard octal file permission mode
	// for executable files and directories (typically 0755).
	FileModeExecutable os.FileMode = 0o755 // rwxr-xr-x — directories, scripts

	// FileModeReadable represents the standard octal file permission mode
	// for readable files (typically 0644).
	FileModeReadable os.FileMode = 0o644 // rw-r--r-- — regular files
)
