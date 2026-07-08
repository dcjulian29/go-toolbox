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
	"errors"

	"github.com/dcjulian29/go-toolbox/filesystem"
	"gopkg.in/yaml.v3"
)

// Save writes cfg to the backing configuration file, creating the ".config"
// directory if necessary, and refreshes the cached instance so subsequent
// [File.Load] calls return the saved value. It returns an error if cfg is nil
// or the file cannot be marshaled or written.
func (f *File[T]) Save(cfg *T) error {
	if cfg == nil {
		return errors.New("cannot save an uninitialized configuration")
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	path, err := f.Path()
	if err != nil {
		return err
	}

	f.mutex.Lock()
	defer f.mutex.Unlock()

	if err := filesystem.EnsureFileExist(path, data); err != nil {
		return err
	}

	f.instance = cfg

	return nil
}
