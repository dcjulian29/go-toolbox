package docker

import (
	"fmt"
	"os"
	"strings"
)

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

// HostContainerVolume normalizes the working directory path by replacing Windows
// backslash path separators with forward slashes for use inside the Linux container
// and returns the volume mapping to replicate the host path inside the container
// and the working directory inside the Linux container.
func HostContainerVolume() (string, string) {
	pwd, _ := os.Getwd()
	pwd = strings.ReplaceAll(strings.ReplaceAll(pwd, "\\", "/"), ":", "")

	host := fmt.Sprintf("%s:\\", string(pwd[0]))
	container := fmt.Sprintf("/%s", string(pwd[0]))

	data := pwd[2:]

	work := fmt.Sprintf("%s/%s", container, data)

	return fmt.Sprintf("%s:%s", host, container), work
}
