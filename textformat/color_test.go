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

package textformat

import (
	"strings"
	"testing"
)

func TestColorFunctions_CorrectANSICodes(t *testing.T) {
	tests := []struct {
		fn   func(string) string
		name string
		code string
	}{
		{Black, "Black", "\033[1;30m"},
		{Blue, "Blue", "\033[1;34m"},
		{Green, "Green", "\033[1;32m"},
		{Magenta, "Magenta", "\033[1;35m"},
		{Red, "Red", "\033[1;31m"},
		{Teal, "Teal", "\033[1;36m"},
		{White, "White", "\033[1;37m"},
		{Yellow, "Yellow", "\033[1;33m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fn("test")
			want := tt.code + "test" + "\033[0m"
			if got != want {
				t.Errorf("%s(\"test\") = %q, want %q", tt.name, got, want)
			}
		})
	}
}

func TestColorFunctions_ContainText(t *testing.T) {
	tests := []struct {
		fn   func(string) string
		name string
	}{
		{Black, "Black"},
		{Blue, "Blue"},
		{Green, "Green"},
		{Magenta, "Magenta"},
		{Red, "Red"},
		{Teal, "Teal"},
		{White, "White"},
		{Yellow, "Yellow"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text := "hello world"
			got := tt.fn(text)
			if !strings.Contains(got, text) {
				t.Errorf("%s(%q) = %q, does not contain original text", tt.name, text, got)
			}
		})
	}
}

func TestColorFunctions_ResetSuffix(t *testing.T) {
	tests := []struct {
		fn   func(string) string
		name string
	}{
		{Black, "Black"},
		{Blue, "Blue"},
		{Green, "Green"},
		{Magenta, "Magenta"},
		{Red, "Red"},
		{Teal, "Teal"},
		{White, "White"},
		{Yellow, "Yellow"},
	}

	reset := "\033[0m"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fn("test")
			if !strings.HasSuffix(got, reset) {
				t.Errorf("%s(\"test\") = %q, does not end with reset code %q", tt.name, got, reset)
			}
		})
	}
}

func TestColorFunctions_EmptyString(t *testing.T) {
	tests := []struct {
		fn   func(string) string
		name string
	}{
		{Black, "Black"},
		{Blue, "Blue"},
		{Green, "Green"},
		{Magenta, "Magenta"},
		{Red, "Red"},
		{Teal, "Teal"},
		{White, "White"},
		{Yellow, "Yellow"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fn("")
			if !strings.Contains(got, "\033[0m") {
				t.Errorf("%s(\"\") = %q, expected reset code even for empty string", tt.name, got)
			}
		})
	}
}

func TestColorFunctions_SpecialCharacters(t *testing.T) {
	special := "hello\nworld\ttab \"quoted\" 100%"

	tests := []struct {
		fn   func(string) string
		name string
	}{
		{Black, "Black"},
		{Blue, "Blue"},
		{Green, "Green"},
		{Magenta, "Magenta"},
		{Red, "Red"},
		{Teal, "Teal"},
		{White, "White"},
		{Yellow, "Yellow"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fn(special)
			if !strings.Contains(got, special) {
				t.Errorf("%s() did not preserve special characters, got %q", tt.name, got)
			}
		})
	}
}

func TestColorFunctions_UniqueColorCodes(t *testing.T) {
	colors := map[string]func(string) string{
		"Black":   Black,
		"Blue":    Blue,
		"Green":   Green,
		"Magenta": Magenta,
		"Red":     Red,
		"Teal":    Teal,
		"White":   White,
		"Yellow":  Yellow,
	}

	results := make(map[string]string)

	for name, fn := range colors {
		results[name] = fn("test")
	}

	for name1, result1 := range results {
		for name2, result2 := range results {
			if name1 != name2 && result1 == result2 {
				t.Errorf("%s and %s produce identical output: %q", name1, name2, result1)
			}
		}
	}
}

func TestFatal_MatchesRed(t *testing.T) {
	text := "critical error"
	if Fatal(text) != Red(text) {
		t.Errorf("Fatal(%q) = %q, want %q (should match Red)", text, Fatal(text), Red(text))
	}
}

func TestInfo_MatchesTeal(t *testing.T) {
	text := "info message"
	if Info(text) != Teal(text) {
		t.Errorf("Info(%q) = %q, want %q (should match Teal)", text, Info(text), Teal(text))
	}
}

func TestWarn_MatchesYellow(t *testing.T) {
	text := "warning message"
	if Warn(text) != Yellow(text) {
		t.Errorf("Warn(%q) = %q, want %q (should match Yellow)", text, Warn(text), Yellow(text))
	}
}
