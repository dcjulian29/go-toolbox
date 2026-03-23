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
	"os"
	"strings"
)

// EnvironmentVariablesWithPrefix scans the host environment for variables
// with the given prefix and returns them as a map compatible with
// ContainerOptions.EnvironmentVariables. The prefix is preserved in each
// returned key (e.g. "APP_DB_HOST" → "APP_DB_HOST").
func EnvironmentVariablesWithPrefix(prefix string) map[string]string {
	return findEnvironmentVariables(prefix, false)
}

// EnvironmentVariablesWithStrippedPrefix scans the host environment for
// variables with the given prefix and returns them as a map compatible with
// ContainerOptions.EnvironmentVariables. The prefix is removed from each
// returned key (e.g. "APP_DB_HOST" → "DB_HOST").
func EnvironmentVariablesWithStrippedPrefix(prefix string) map[string]string {
	return findEnvironmentVariables(prefix, true)
}

func findEnvironmentVariables(prefix string, stripPrefix bool) map[string]string { //nolint
	if !strings.HasSuffix(prefix, "_") {
		prefix += "_"
	}

	env := make(map[string]string)

	if prefix == "_" { //nolint
		return env
	}

	for _, entry := range os.Environ() {
		key, value, found := strings.Cut(entry, "=")
		if !found || !strings.HasPrefix(key, prefix) {
			continue
		}

		if stripPrefix {
			key = strings.TrimPrefix(key, prefix)
		}

		env[key] = value
	}

	return env
}
