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
package io

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
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

func EnsureDirectoryExist(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return err
		}
	}

	return nil
}

func RemoveDirectory(dirPath string) error {
	if DirectoryExists(dirPath) {
		files, err := filepath.Glob(filepath.Join(dirPath, "*"))
		if err != nil {
			return err
		}

		for _, file := range files {
			err := os.RemoveAll(file)
			if err != nil {
				return err
			}
		}

		return os.Remove(dirPath)
	}

	return nil
}

func RemoveFile(filePath string) error {
	if FileExists(filePath) {
		return os.Remove(filePath)
	}

	return nil
}

func FileHash(filePath string) (string, error) {
	hash := sha256.New()

	sourceFile, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	defer sourceFile.Close() //nolint:errcheck

	_, err = io.Copy(hash, sourceFile)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func ScanDirectory(dir_path string, ignore []string) ([]string, []string, error) {
	folders := []string{}
	files := []string{}

	err := filepath.Walk(dir_path, func(path string, f os.FileInfo, err error) error {
		_continue := false

		for _, i := range ignore {
			if strings.Contains(path, i) {
				_continue = true
			}
		}

		if !_continue {
			s, err := os.Stat(path)
			if err != nil {
				return err
			}

			f_mode := s.Mode()

			if f_mode.IsDir() {
				folders = append(folders, path)
			} else if f_mode.IsRegular() {
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
