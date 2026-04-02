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

package hypervmachine

// SecureBootTemplate identifies the Secure Boot certificate template
// applied to a Generation 2 VM.
type SecureBootTemplate string

const (
	// MicrosoftUEFICertificateAuthority uses the Microsoft UEFI CA,
	// suitable for Linux distributions and other non-Windows OSes.
	MicrosoftUEFICertificateAuthority SecureBootTemplate = "MicrosoftUEFICertificateAuthority"

	// MicrosoftWindows uses the Microsoft Windows certificate template.
	MicrosoftWindows SecureBootTemplate = "MicrosoftWindows"

	// OpenSourceShieldedVM uses the open-source shielded VM template,
	// enabling Secure Boot for Linux VMs running on a guarded Hyper-V
	// host with Host Guardian Service (HGS) attestation.
	OpenSourceShieldedVM SecureBootTemplate = "OpenSourceShieldedVM"

	// LinuxVM is a convenience alias for MicrosoftUEFICertificateAuthority.
	LinuxVM SecureBootTemplate = MicrosoftUEFICertificateAuthority

	// WindowsVM is a convenience alias for MicrosoftWindows.
	WindowsVM SecureBootTemplate = MicrosoftWindows
)
