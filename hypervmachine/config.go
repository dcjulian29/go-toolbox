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

// Config holds the parameters used when creating a Hyper-V VM.
type Config struct {
	// Name is the display name of the VM in Hyper-V.
	Name string
	// VHDXPath is the full path to the VM's virtual hard disk.
	VHDXPath string
	// VirtualSwitch is the name of the Hyper-V virtual switch to connect.
	VirtualSwitch string
	// MemoryBytes is the startup memory size in bytes. When dynamic memory is
	// enabled, this value is also used as the minimum memory allocation.
	MemoryBytes int64
	// MaximumMemoryBytes is the upper bound for dynamic memory allocation.
	// Must be greater than or equal to MemoryBytes. When set equal to
	// MemoryBytes, the VM uses static (fixed) memory.
	MaximumMemoryBytes int64
	// ProcessorCount is the number of virtual CPUs assigned to the VM.
	ProcessorCount int
	// Generation is the VM generation. Use GenerationV1 or GenerationV2.
	Generation Generation
	// SecureBoot enables Secure Boot (only relevant for GenerationV2).
	SecureBoot bool
}
