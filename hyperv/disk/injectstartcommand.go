//go:build windows

package hyperv_disk

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
	"fmt"
	"path/filepath"

	"github.com/dcjulian29/go-toolbox/filesystem"
	"github.com/dcjulian29/go-toolbox/textformat"
)

// InjectStartCommand injects the startup command executed just after the first login.
func InjectStartCommand(cfg *InjectConfig) error {
	path := filepath.Join(cfg.MountedDrive, "Windows", "Setup", "Scripts", "SetupComplete.cmd")

	content := "%WINDIR%\\System32\\WindowsPowerShell\\v1.0\\powershell.exe "
	content = content + "-NoProfile -NonInteractive -ExecutionPolicy Bypass -NoLogo -Command "
	content = content + "%WINDIR%\\Setup\\Scripts\\" + filepath.Base(cfg.StartScript)

	fmt.Println(textformat.Info("[inject] start command"))

	return filesystem.EnsureFileExist(path, []byte(content))
}
