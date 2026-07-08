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

// Package configuration provides a generic, singleton-cached loader and saver
// for a YAML configuration file backed by a caller-supplied struct type.
//
// A caller creates one [File] for its own configuration type, typically as a
// package-level variable so the cached instance behaves as a process-wide
// singleton:
//
//	type Config struct {
//		Editor string `yaml:"editor"`
//	}
//
//	var cfg = configuration.New[Config]("mytool.yml")
//
//	c, _ := cfg.Load()  // reads ~/.config/mytool.yml once, then caches
//	c.Editor = "vim"
//	_ = cfg.Save(&c)    // writes ~/.config/mytool.yml and refreshes the cache
package configuration

import (
	"os"
	"path/filepath"
	"sync"
)

// File is a singleton-cached handle to a YAML configuration file whose contents
// are described by the caller-supplied type T. A File is safe for concurrent
// use. Create one per configuration file with [New].
type File[T any] struct {
	name     string
	instance *T
	loadErr  error
	once     sync.Once
	mutex    sync.RWMutex
}

// New returns a [File] that reads and writes the named file (for example
// "mytool.yml") in the current user's ".config" directory.
func New[T any](name string) *File[T] {
	return &File[T]{name: name}
}

// Path returns the absolute path of the backing configuration file,
// <user home>/.config/<name>. An error is returned only if the user's home
// directory cannot be determined.
func (f *File[T]) Path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".config", f.name), nil
}
