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
	"runtime"
	"strings"
	"testing"
)

func TestHostContainerVolume_ReturnsCurrentDirectory(t *testing.T) {
	volume, work, err := HostContainerVolume()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if volume == "" {
		t.Error("volume should not be empty")
	}

	if work == "" {
		t.Error("work should not be empty")
	}
}

func TestHostContainerVolume_WorkContainsNoBackslashes(t *testing.T) {
	_, work, err := HostContainerVolume()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(work, "\\") {
		t.Errorf("work path should not contain backslashes: %q", work)
	}
}

func TestHostContainerVolume_WorkContainsNoColons(t *testing.T) {
	_, work, err := HostContainerVolume()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(work, ":") {
		t.Errorf("work path should not contain colons: %q", work)
	}
}

func TestHostContainerVolume_WorkStartsWithForwardSlash(t *testing.T) {
	_, work, err := HostContainerVolume()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.HasPrefix(work, "/") {
		t.Errorf("work path should start with /: %q", work)
	}
}

func TestHostContainerVolume_VolumeContainsSeparator(t *testing.T) {
	volume, _, err := HostContainerVolume()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(volume, ":") {
		t.Errorf("volume should contain a colon separator: %q", volume)
	}
}

func TestHostContainerVolume_LinuxVolume(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("test validates Linux/macOS behavior")
	}

	volume, _, err := HostContainerVolume()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if volume != "/:/" {
		t.Errorf("volume = %q, want %q", volume, "/:/")
	}
}

func TestHostContainerVolume_LinuxWorkMatchesPwd(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("test validates Linux/macOS behavior")
	}

	pwd, _ := os.Getwd()

	_, work, err := HostContainerVolume()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if work != pwd {
		t.Errorf("work = %q, want %q", work, pwd)
	}
}

func TestWindowsContainerVolume_StandardPath(t *testing.T) {
	volume, work, err := windowsContainerVolume(`C:\Users\dev\project`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if volume != `C:\:/C` {
		t.Errorf("volume = %q, want %q", volume, `C:\:/C`)
	}

	if work != "/C/Users/dev/project" {
		t.Errorf("work = %q, want %q", work, "/C/Users/dev/project")
	}
}

func TestWindowsContainerVolume_DifferentDriveLetter(t *testing.T) {
	volume, work, err := windowsContainerVolume(`D:\Data\repo`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if volume != `D:\:/D` {
		t.Errorf("volume = %q, want %q", volume, `D:\:/D`)
	}

	if work != "/D/Data/repo" {
		t.Errorf("work = %q, want %q", work, "/D/Data/repo")
	}
}

func TestWindowsContainerVolume_ForwardSlashPath(t *testing.T) {
	volume, work, err := windowsContainerVolume("C:/Users/dev/project")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if volume != `C:\:/C` {
		t.Errorf("volume = %q, want %q", volume, `C:\:/C`)
	}

	if work != "/C/Users/dev/project" {
		t.Errorf("work = %q, want %q", work, "/C/Users/dev/project")
	}
}

func TestWindowsContainerVolume_DriveRootOnly(t *testing.T) {
	volume, work, err := windowsContainerVolume(`C:\`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if volume != `C:\:/C` {
		t.Errorf("volume = %q, want %q", volume, `C:\:/C`)
	}

	if work != "/C/" {
		t.Errorf("work = %q, want %q", work, "/C/")
	}
}

func TestWindowsContainerVolume_NestedPath(t *testing.T) {
	_, work, err := windowsContainerVolume(`C:\a\b\c\d\e`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if work != "/C/a/b/c/d/e" {
		t.Errorf("work = %q, want %q", work, "/C/a/b/c/d/e")
	}
}

func TestWindowsContainerVolume_InvalidPathTooShort(t *testing.T) {
	_, _, err := windowsContainerVolume("C")
	if err == nil {
		t.Fatal("expected error for path too short")
	}
}

func TestWindowsContainerVolume_InvalidPathNoColon(t *testing.T) {
	_, _, err := windowsContainerVolume("/home/user/project")
	if err == nil {
		t.Fatal("expected error for non-Windows path")
	}
}

func TestWindowsContainerVolume_InvalidPathNoSeparator(t *testing.T) {
	_, _, err := windowsContainerVolume("C:project")
	if err == nil {
		t.Fatal("expected error for missing separator after drive letter")
	}
}
