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
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const FileModeExecutable = 0755

func FileExists(path string) bool {
	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func DirectoryExists(path string) bool {
	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}

	return info.IsDir()
}

func EnsureDirectoryExist(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, FileModeExecutable); err != nil {
			return err
		}
	}

	return nil
}

func RemoveDirectory(path string) error {
	if DirectoryExists(path) {
		files, err := filepath.Glob(filepath.Join(path, "*"))
		if err != nil {
			return err
		}

		for _, file := range files {
			err := os.RemoveAll(file)
			if err != nil {
				return err
			}
		}

		return os.Remove(path)
	}

	return nil
}

func RemoveFile(path string) error {
	if FileExists(path) {
		return os.Remove(path)
	}

	return nil
}

func FileHash(path string) (string, error) {
	hash := sha256.New()

	sourceFile, err := os.Open(path)
	if err != nil {
		return "", err
	}

	defer sourceFile.Close() //nolint:errcheck

	_, err = io.Copy(hash, sourceFile)

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func ScanDirectory(path string, ignore []string) ([]string, []string, error) {
	folders := []string{}
	files := []string{}

	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		_continue := false

		for _, i := range ignore {
			if strings.Contains(path, i) {
				_continue = true
			}
		}

		if !_continue {
			m := f.Mode()

			if m.IsDir() {
				folders = append(folders, path)
			} else if m.IsRegular() {
				files = append(files, path)
			}
		}

		return nil
	})

	if err != nil {
		return []string{}, []string{}, err
	}

	return folders, files, nil
}
