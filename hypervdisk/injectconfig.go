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

package hypervdisk

// InjectConfig holds all the values needed for injecting files to a VHDX file.
type InjectConfig struct {
	// ComputerName is the hostname to set in the injected configuration.
	ComputerName string

	// InstallPackage is the path to a package to install during installation.
	InstallPackage string

	// MountedDrive is the drive letter of the mounted VHDX (e.g. "E:").
	MountedDrive string

	// StartScript is the path to a script to run on first boot.
	StartScript string

	// UnattendTemplate is the path to the unattend.xml template file.
	UnattendTemplate string

	// UserName is the local account name to create.
	UserName string

	// UserPassword is the password for the local account.
	UserPassword string
}
