package filesystem

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

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"

	"github.com/dcjulian29/go-toolbox/textformat"
)

// FileHash computes and returns the cryptographic hash (e.g., SHA256) of
// the specified file's contents.
func FileHash(path string) (string, error) {
	hash := sha256.New()

	sourceFile, err := os.Open(path)
	if err != nil {
		return textformat.EmptyString, err
	}

	defer sourceFile.Close() //nolint:errcheck

	_, err = io.Copy(hash, sourceFile)

	if err != nil {
		return textformat.EmptyString, err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
