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

package hypervmachine_test

import (
	"strings"
	"testing"

	"github.com/dcjulian29/go-toolbox/hypervmachine"
)

func validConfig() *hypervmachine.Config {
	return &hypervmachine.Config{
		Name:               "TestVM",
		VHDXPath:           `C:\VMs\TestVM\disk.vhdx`,
		VirtualSwitch:      "Default Switch",
		MemoryBytes:        1073741824, // 1 GB
		MaximumMemoryBytes: 2147483648, // 2 GB
		ProcessorCount:     2,
		Generation:         hypervmachine.GenerationV2,
		SecureBoot:         true,
	}
}

// --- Valid configs ---

func TestValidateConfig_ValidConfig(t *testing.T) {
	cfg := validConfig()

	if err := hypervmachine.ValidateConfig(cfg); err != nil {
		t.Errorf("unexpected error for valid config: %v", err)
	}
}

func TestValidateConfig_ValidStaticMemory(t *testing.T) {
	cfg := validConfig()
	cfg.MaximumMemoryBytes = cfg.MemoryBytes

	if err := hypervmachine.ValidateConfig(cfg); err != nil {
		t.Errorf("unexpected error when MaximumMemoryBytes == MemoryBytes: %v", err)
	}
}

func TestValidateConfig_ValidGenerationV1(t *testing.T) {
	cfg := validConfig()
	cfg.Generation = hypervmachine.GenerationV1

	if err := hypervmachine.ValidateConfig(cfg); err != nil {
		t.Errorf("unexpected error for GenerationV1: %v", err)
	}
}

func TestValidateConfig_ValidGenerationV2(t *testing.T) {
	cfg := validConfig()
	cfg.Generation = hypervmachine.GenerationV2

	if err := hypervmachine.ValidateConfig(cfg); err != nil {
		t.Errorf("unexpected error for GenerationV2: %v", err)
	}
}

// --- Nil config ---

func TestValidateConfig_NilConfig(t *testing.T) {
	err := hypervmachine.ValidateConfig(nil)
	if err == nil {
		t.Fatal("expected error for nil config")
	}

	if !strings.Contains(err.Error(), "nil") {
		t.Errorf("error should mention nil: %v", err)
	}
}

// --- Name validation ---

func TestValidateConfig_EmptyName(t *testing.T) {
	cfg := validConfig()
	cfg.Name = ""

	err := hypervmachine.ValidateConfig(cfg)
	if err == nil {
		t.Fatal("expected error for empty name")
	}

	if !strings.Contains(err.Error(), "name") {
		t.Errorf("error should mention name: %v", err)
	}
}

func TestValidateConfig_WhitespaceOnlyName(t *testing.T) {
	cfg := validConfig()
	cfg.Name = "   \t\n"

	err := hypervmachine.ValidateConfig(cfg)
	if err == nil {
		t.Fatal("expected error for whitespace-only name")
	}
}

// --- VHDXPath validation ---

func TestValidateConfig_EmptyVHDXPath(t *testing.T) {
	cfg := validConfig()
	cfg.VHDXPath = ""

	err := hypervmachine.ValidateConfig(cfg)
	if err == nil {
		t.Fatal("expected error for empty VHDXPath")
	}

	if !strings.Contains(err.Error(), "disk") {
		t.Errorf("error should mention disk: %v", err)
	}
}

func TestValidateConfig_WhitespaceOnlyVHDXPath(t *testing.T) {
	cfg := validConfig()
	cfg.VHDXPath = "   "

	err := hypervmachine.ValidateConfig(cfg)
	if err == nil {
		t.Fatal("expected error for whitespace-only VHDXPath")
	}
}

// --- VirtualSwitch validation ---

func TestValidateConfig_EmptyVirtualSwitch(t *testing.T) {
	cfg := validConfig()
	cfg.VirtualSwitch = ""

	err := hypervmachine.ValidateConfig(cfg)
	if err == nil {
		t.Fatal("expected error for empty VirtualSwitch")
	}

	if !strings.Contains(err.Error(), "switch") {
		t.Errorf("error should mention switch: %v", err)
	}
}

func TestValidateConfig_WhitespaceOnlyVirtualSwitch(t *testing.T) {
	cfg := validConfig()
	cfg.VirtualSwitch = "\t"

	err := hypervmachine.ValidateConfig(cfg)
	if err == nil {
		t.Fatal("expected error for whitespace-only VirtualSwitch")
	}
}

// --- Generation validation ---

func TestValidateConfig_InvalidGenerationZero(t *testing.T) {
	cfg := validConfig()
	cfg.Generation = 0

	err := hypervmachine.ValidateConfig(cfg)
	if err == nil {
		t.Fatal("expected error for generation 0")
	}

	if !strings.Contains(err.Error(), "generation") {
		t.Errorf("error should mention generation: %v", err)
	}
}

func TestValidateConfig_InvalidGenerationThree(t *testing.T) {
	cfg := validConfig()
	cfg.Generation = 3

	err := hypervmachine.ValidateConfig(cfg)
	if err == nil {
		t.Fatal("expected error for generation 3")
	}
}

func TestValidateConfig_InvalidGenerationNegative(t *testing.T) {
	cfg := validConfig()
	cfg.Generation = -1

	err := hypervmachine.ValidateConfig(cfg)
	if err == nil {
		t.Fatal("expected error for negative generation")
	}
}

// --- MemoryBytes validation ---

func TestValidateConfig_ZeroMemoryBytes(t *testing.T) {
	cfg := validConfig()
	cfg.MemoryBytes = 0

	err := hypervmachine.ValidateConfig(cfg)
	if err == nil {
		t.Fatal("expected error for zero MemoryBytes")
	}

	if !strings.Contains(err.Error(), "memory") {
		t.Errorf("error should mention memory: %v", err)
	}
}

func TestValidateConfig_NegativeMemoryBytes(t *testing.T) {
	cfg := validConfig()
	cfg.MemoryBytes = -1048576

	err := hypervmachine.ValidateConfig(cfg)
	if err == nil {
		t.Fatal("expected error for negative MemoryBytes")
	}
}

// --- MaximumMemoryBytes validation ---

func TestValidateConfig_MaxMemoryLessThanMemory(t *testing.T) {
	cfg := validConfig()
	cfg.MemoryBytes = 2147483648
	cfg.MaximumMemoryBytes = 1073741824

	err := hypervmachine.ValidateConfig(cfg)
	if err == nil {
		t.Fatal("expected error when MaximumMemoryBytes < MemoryBytes")
	}

	if !strings.Contains(err.Error(), "maximum memory") {
		t.Errorf("error should mention maximum memory: %v", err)
	}
}

func TestValidateConfig_MaxMemoryEqualToMemory(t *testing.T) {
	cfg := validConfig()
	cfg.MaximumMemoryBytes = cfg.MemoryBytes

	if err := hypervmachine.ValidateConfig(cfg); err != nil {
		t.Errorf("MaximumMemoryBytes == MemoryBytes should be valid (static memory): %v", err)
	}
}

// --- ProcessorCount validation ---

func TestValidateConfig_ZeroProcessorCount(t *testing.T) {
	cfg := validConfig()
	cfg.ProcessorCount = 0

	err := hypervmachine.ValidateConfig(cfg)
	if err == nil {
		t.Fatal("expected error for zero ProcessorCount")
	}

	if !strings.Contains(err.Error(), "processor") {
		t.Errorf("error should mention processor: %v", err)
	}
}

func TestValidateConfig_NegativeProcessorCount(t *testing.T) {
	cfg := validConfig()
	cfg.ProcessorCount = -1

	err := hypervmachine.ValidateConfig(cfg)
	if err == nil {
		t.Fatal("expected error for negative ProcessorCount")
	}
}
