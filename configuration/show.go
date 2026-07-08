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

package configuration

import (
	"gopkg.in/yaml.v3"
)

// Show returns the current configuration rendered as YAML. It loads the
// configuration (see [File.Load], so the value is the cached one after the
// first call) and marshals it, giving tools a ready-to-print representation of
// the effective config file contents.
func (f *File[T]) Show() (string, error) {
	cfg, err := f.Load()
	if err != nil {
		return "", err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
