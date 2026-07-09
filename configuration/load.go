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
	"fmt"
	"os"

	"github.com/dcjulian29/go-toolbox/filesystem"
	"gopkg.in/yaml.v3"
)

// Load returns the configuration, reading and parsing the backing file on the
// first call and returning the cached copy on every call thereafter. A missing
// file is not an error: a zero-valued T is returned. Because the result is
// cached on first use, later external edits to the file are not observed until
// the process restarts (or [File.Save] replaces the cache).
func (f *File[T]) Load() (T, error) {
	f.once.Do(func() {
		cfg, err := f.loadFromDisk()

		f.mutex.Lock()
		f.instance = cfg
		f.loadErr = err
		f.mutex.Unlock()
	})

	f.mutex.RLock()
	defer f.mutex.RUnlock()

	return *f.instance, f.loadErr
}

func (f *File[T]) loadFromDisk() (*T, error) {
	cfg := new(T)

	path, err := f.Path()
	if err != nil {
		return cfg, err
	}

	if !filesystem.FileExist(path) {
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("could not read configuration '%s': %w", f.name, err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return cfg, fmt.Errorf("could not parse configuration '%s': %w", f.name, err)
	}

	return cfg, nil
}
