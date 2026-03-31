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
	"testing"
)

// --- Absolute Windows paths ---

func TestEnsurePathIsUnix_AbsolutePathWithBackslashes(t *testing.T) {
	result := EnsurePathIsUnix(`C:\Users\julian\file.txt`)

	if result != "/C/Users/julian/file.txt" {
		t.Errorf("got %q; want %q", result, "/C/Users/julian/file.txt")
	}
}

func TestEnsurePathIsUnix_AbsolutePathWithForwardSlashes(t *testing.T) {
	result := EnsurePathIsUnix("D:/Projects/go/main.go")

	if result != "/D/Projects/go/main.go" {
		t.Errorf("got %q; want %q", result, "/D/Projects/go/main.go")
	}
}

func TestEnsurePathIsUnix_AbsolutePathMixedSeparators(t *testing.T) {
	result := EnsurePathIsUnix(`C:\Users/julian\file.txt`)

	if result != "/C/Users/julian/file.txt" {
		t.Errorf("got %q; want %q", result, "/C/Users/julian/file.txt")
	}
}

func TestEnsurePathIsUnix_LowercaseDriveLetter(t *testing.T) {
	result := EnsurePathIsUnix(`c:\users\julian`)

	if result != "/c/users/julian" {
		t.Errorf("got %q; want %q", result, "/c/users/julian")
	}
}

func TestEnsurePathIsUnix_DriveRootWithBackslash(t *testing.T) {
	result := EnsurePathIsUnix(`C:\`)

	if result != "/C/" {
		t.Errorf("got %q; want %q", result, "/C/")
	}
}

func TestEnsurePathIsUnix_DriveRootWithForwardSlash(t *testing.T) {
	result := EnsurePathIsUnix("C:/")

	if result != "/C/" {
		t.Errorf("got %q; want %q", result, "/C/")
	}
}

func TestEnsurePathIsUnix_DeeplyNestedAbsolutePath(t *testing.T) {
	result := EnsurePathIsUnix(`C:\a\b\c\d\e\f\g\file.txt`)

	if result != "/C/a/b/c/d/e/f/g/file.txt" {
		t.Errorf("got %q; want %q", result, "/C/a/b/c/d/e/f/g/file.txt")
	}
}

func TestEnsurePathIsUnix_AbsolutePathWithSpaces(t *testing.T) {
	result := EnsurePathIsUnix(`C:\Program Files\My App\file.txt`)

	if result != "/C/Program Files/My App/file.txt" {
		t.Errorf("got %q; want %q", result, "/C/Program Files/My App/file.txt")
	}
}

// --- Relative Windows paths ---

func TestEnsurePathIsUnix_DotBackslashRelative(t *testing.T) {
	result := EnsurePathIsUnix(`.\folder\file.txt`)

	if result != "./folder/file.txt" {
		t.Errorf("got %q; want %q", result, "./folder/file.txt")
	}
}

func TestEnsurePathIsUnix_DotDotBackslashRelative(t *testing.T) {
	result := EnsurePathIsUnix(`..\parent\child\file.txt`)

	if result != "../parent/child/file.txt" {
		t.Errorf("got %q; want %q", result, "../parent/child/file.txt")
	}
}

func TestEnsurePathIsUnix_BareRelativeWithoutDotPrefix(t *testing.T) {
	result := EnsurePathIsUnix(`folder\subfolder\file.txt`)

	if result != "folder/subfolder/file.txt" {
		t.Errorf("got %q; want %q", result, "folder/subfolder/file.txt")
	}
}

func TestEnsurePathIsUnix_RelativeSingleBackslashSegment(t *testing.T) {
	result := EnsurePathIsUnix(`folder\file.txt`)

	if result != "folder/file.txt" {
		t.Errorf("got %q; want %q", result, "folder/file.txt")
	}
}

// --- Embedded paths with equals prefix ---

func TestEnsurePathIsUnix_EqualsPrefixAbsolutePath(t *testing.T) {
	result := EnsurePathIsUnix(`--config=C:\Users\julian\config.yaml`)

	if result != "--config=/C/Users/julian/config.yaml" {
		t.Errorf("got %q; want %q", result, "--config=/C/Users/julian/config.yaml")
	}
}

func TestEnsurePathIsUnix_EqualsPrefixRelativePath(t *testing.T) {
	result := EnsurePathIsUnix(`--output=.\build\result.bin`)

	if result != "--output=./build/result.bin" {
		t.Errorf("got %q; want %q", result, "--output=./build/result.bin")
	}
}

func TestEnsurePathIsUnix_EqualsPrefixBareRelative(t *testing.T) {
	result := EnsurePathIsUnix(`--dir=src\main\resources`)

	if result != "--dir=src/main/resources" {
		t.Errorf("got %q; want %q", result, "--dir=src/main/resources")
	}
}

// --- Embedded paths with quotes ---

func TestEnsurePathIsUnix_DoubleQuotedAbsolutePath(t *testing.T) {
	result := EnsurePathIsUnix(`"C:\Users\julian\file.txt"`)

	if result != `"/C/Users/julian/file.txt"` {
		t.Errorf("got %q; want %q", result, `"/C/Users/julian/file.txt"`)
	}
}

func TestEnsurePathIsUnix_SingleQuotedRelativePath(t *testing.T) {
	result := EnsurePathIsUnix(`'.\folder\file.txt'`)

	if result != `'./folder/file.txt'` {
		t.Errorf("got %q; want %q", result, `'./folder/file.txt'`)
	}
}

func TestEnsurePathIsUnix_DoubleQuotedRelativeAfterEquals(t *testing.T) {
	result := EnsurePathIsUnix(`--config=".\config\app.yaml"`)

	if result != `--config="./config/app.yaml"` {
		t.Errorf("got %q; want %q", result, `--config="./config/app.yaml"`)
	}
}

// --- Embedded paths with spaces ---

func TestEnsurePathIsUnix_SpaceSeparatedPaths(t *testing.T) {
	result := EnsurePathIsUnix(`copy C:\source\file.txt C:\dest\file.txt`)

	if result != "copy /C/source/file.txt /C/dest/file.txt" {
		t.Errorf("got %q; want %q", result, "copy /C/source/file.txt /C/dest/file.txt")
	}
}

func TestEnsurePathIsUnix_SpacePrefixedRelativePath(t *testing.T) {
	result := EnsurePathIsUnix(`run .\scripts\build.sh`)

	if result != "run ./scripts/build.sh" {
		t.Errorf("got %q; want %q", result, "run ./scripts/build.sh")
	}
}

// --- Already-Unix paths pass through unchanged ---

func TestEnsurePathIsUnix_UnixAbsolutePathUnchanged(t *testing.T) {
	result := EnsurePathIsUnix("/home/julian/file.txt")

	if result != "/home/julian/file.txt" {
		t.Errorf("got %q; want %q", result, "/home/julian/file.txt")
	}
}

func TestEnsurePathIsUnix_UnixRelativePathUnchanged(t *testing.T) {
	result := EnsurePathIsUnix("./folder/file.txt")

	if result != "./folder/file.txt" {
		t.Errorf("got %q; want %q", result, "./folder/file.txt")
	}
}

func TestEnsurePathIsUnix_UnixDotDotRelativeUnchanged(t *testing.T) {
	result := EnsurePathIsUnix("../parent/file.txt")

	if result != "../parent/file.txt" {
		t.Errorf("got %q; want %q", result, "../parent/file.txt")
	}
}

// --- Strings that should not be modified ---

func TestEnsurePathIsUnix_HttpURLNotModified(t *testing.T) {
	result := EnsurePathIsUnix("http://example.com/path")

	if result != "http://example.com/path" {
		t.Errorf("got %q; want %q", result, "http://example.com/path")
	}
}

func TestEnsurePathIsUnix_HttpsURLNotModified(t *testing.T) {
	result := EnsurePathIsUnix("https://example.com/path/to/resource")

	if result != "https://example.com/path/to/resource" {
		t.Errorf("got %q; want %q", result, "https://example.com/path/to/resource")
	}
}

func TestEnsurePathIsUnix_LocalhostWithPortNotModified(t *testing.T) {
	result := EnsurePathIsUnix("localhost:8080/api")

	if result != "localhost:8080/api" {
		t.Errorf("got %q; want %q", result, "localhost:8080/api")
	}
}

func TestEnsurePathIsUnix_PlainTextNotModified(t *testing.T) {
	result := EnsurePathIsUnix("hello world")

	if result != "hello world" {
		t.Errorf("got %q; want %q", result, "hello world")
	}
}

// --- Edge cases ---

func TestEnsurePathIsUnix_EmptyString(t *testing.T) {
	result := EnsurePathIsUnix("")

	if result != "" {
		t.Errorf("got %q; want %q", result, "")
	}
}

func TestEnsurePathIsUnix_SingleFilenameNoSeparators(t *testing.T) {
	result := EnsurePathIsUnix("file.txt")

	if result != "file.txt" {
		t.Errorf("got %q; want %q", result, "file.txt")
	}
}

func TestEnsurePathIsUnix_SingleBackslashOnly(t *testing.T) {
	result := EnsurePathIsUnix(`\`)

	if result != `\` {
		t.Errorf("got %q; want %q", result, `\`)
	}
}

func TestEnsurePathIsUnix_SingleForwardSlashOnly(t *testing.T) {
	result := EnsurePathIsUnix("/")

	if result != "/" {
		t.Errorf("got %q; want %q", result, "/")
	}
}

func TestEnsurePathIsUnix_TrailingBackslash(t *testing.T) {
	result := EnsurePathIsUnix(`C:\Users\julian\`)

	if result != "/C/Users/julian/" {
		t.Errorf("got %q; want %q", result, "/C/Users/julian/")
	}
}

func TestEnsurePathIsUnix_FileExtensionWithDots(t *testing.T) {
	result := EnsurePathIsUnix(`C:\project\archive.tar.gz`)

	if result != "/C/project/archive.tar.gz" {
		t.Errorf("got %q; want %q", result, "/C/project/archive.tar.gz")
	}
}

func TestEnsurePathIsUnix_UNCStylePath(t *testing.T) {
	result := EnsurePathIsUnix(`\\server\share\folder`)

	if result != "//server/share/folder" {
		t.Errorf("got %q; want %q", result, "//server/share/folder")
	}
}
