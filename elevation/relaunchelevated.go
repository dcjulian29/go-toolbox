//go:build windows

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

package elevation

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

// RelaunchElevated re-executes the current process with administrative privileges
// which triggers the Windows UAC prompt.
func RelaunchElevated() error {
	verb, err := syscall.UTF16PtrFromString("runas")
	if err != nil {
		return err
	}

	exe, err := os.Executable()
	if err != nil {
		return err
	}

	exePtr, err := syscall.UTF16PtrFromString(exe)
	if err != nil {
		return err
	}

	var sb strings.Builder

	for i, a := range os.Args[1:] { //nolint:revive
		if i > 0 { //nolint:revive
			sb.WriteString(" ")
		}

		sb.WriteString(a)
	}

	var argsPtr *uint16
	args := sb.String()

	if args != "" {
		argsPtr, err = syscall.UTF16PtrFromString(args)
		if err != nil {
			return err
		}
	}

	shell32 := syscall.NewLazyDLL("shell32.dll")
	shellExec := shell32.NewProc("ShellExecuteW")

	r, _, _ := shellExec.Call(
		0, //nolint:revive
		uintptr(unsafe.Pointer(verb)),
		uintptr(unsafe.Pointer(exePtr)),
		uintptr(unsafe.Pointer(argsPtr)),
		0, //nolint:revive
		syscall.SW_SHOWNORMAL,
	)

	if r <= 32 { //nolint:revive
		return shellExecuteError(r)
	}

	return nil
}

// nolint
func shellExecuteError(code uintptr) error {
	switch code {
	case 0:
		return errors.New("ShellExecuteW: out of memory")
	case 2:
		return errors.New("ShellExecuteW: file not found")
	case 3:
		return errors.New("ShellExecuteW: path not found")
	case 5:
		return errors.New("ShellExecuteW: access denied")
	case 8:
		return errors.New("ShellExecuteW: out of memory")
	case 26:
		return errors.New("ShellExecuteW: sharing violation")
	case 27:
		return errors.New("ShellExecuteW: invalid file association")
	case 31:
		return errors.New("ShellExecuteW: no application associated with file type")
	case 32:
		return errors.New("ShellExecuteW: DDE transaction failed")
	default:
		return fmt.Errorf("ShellExecuteW failed with code %d", code)
	}
}
