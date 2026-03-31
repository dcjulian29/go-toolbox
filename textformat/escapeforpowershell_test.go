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
	"testing"
)

func TestEscapeForPowerShell_EscapesBacktick(t *testing.T) {
	got := EscapeForPowerShell("hello`world")
	want := "hello``world"
	if got != want {
		t.Errorf("EscapeForPowerShell() = %q, want %q", got, want)
	}
}

func TestEscapeForPowerShell_EscapesDollar(t *testing.T) {
	got := EscapeForPowerShell("hello$world")
	want := "hello`$world"
	if got != want {
		t.Errorf("EscapeForPowerShell() = %q, want %q", got, want)
	}
}

func TestEscapeForPowerShell_EscapesDoubleQuote(t *testing.T) {
	got := EscapeForPowerShell(`hello"world`)
	want := "hello`\"world"
	if got != want {
		t.Errorf("EscapeForPowerShell() = %q, want %q", got, want)
	}
}

func TestEscapeForPowerShell_AllSpecialChars(t *testing.T) {
	got := EscapeForPowerShell("`$\"")
	want := "```$`\""
	if got != want {
		t.Errorf("EscapeForPowerShell() = %q, want %q", got, want)
	}
}

func TestEscapeForPowerShell_NoSpecialChars(t *testing.T) {
	input := "hello world 123"
	got := EscapeForPowerShell(input)
	if got != input {
		t.Errorf("EscapeForPowerShell() = %q, want %q (no escaping needed)", got, input)
	}
}

func TestEscapeForPowerShell_EmptyString(t *testing.T) {
	got := EscapeForPowerShell("")
	if got != "" {
		t.Errorf("EscapeForPowerShell(\"\") = %q, want \"\"", got)
	}
}

func TestEscapeForPowerShell_DollarVariable(t *testing.T) {
	got := EscapeForPowerShell("path is $env:PATH")
	want := "path is `$env:PATH"
	if got != want {
		t.Errorf("EscapeForPowerShell() = %q, want %q", got, want)
	}
}

func TestEscapeForPowerShell_Subexpression(t *testing.T) {
	got := EscapeForPowerShell("result: $(Get-Date)")
	want := "result: `$(Get-Date)"
	if got != want {
		t.Errorf("EscapeForPowerShell() = %q, want %q", got, want)
	}
}

func TestEscapeForPowerShell_ConsecutiveBackticks(t *testing.T) {
	got := EscapeForPowerShell("``")
	want := "````"
	if got != want {
		t.Errorf("EscapeForPowerShell() = %q, want %q", got, want)
	}
}

func TestEscapeForPowerShell_ConsecutiveDollars(t *testing.T) {
	got := EscapeForPowerShell("$$")
	want := "`$`$"
	if got != want {
		t.Errorf("EscapeForPowerShell() = %q, want %q", got, want)
	}
}

func TestEscapeForPowerShell_MixedContent(t *testing.T) {
	got := EscapeForPowerShell("Write-Host \"Hello $name\" `n")
	want := "Write-Host `\"Hello `$name`\" ``n"
	if got != want {
		t.Errorf("EscapeForPowerShell() = %q, want %q", got, want)
	}
}

func TestEscapeForPowerShell_SingleQuoteNotEscaped(t *testing.T) {
	input := "it's a test"
	got := EscapeForPowerShell(input)
	if got != input {
		t.Errorf("EscapeForPowerShell() = %q, want %q (single quotes should not be escaped)", got, input)
	}
}

func TestEscapeForPowerShell_BackslashNotEscaped(t *testing.T) {
	input := `C:\Users\test`
	got := EscapeForPowerShell(input)
	if got != input {
		t.Errorf("EscapeForPowerShell() = %q, want %q (backslashes should not be escaped)", got, input)
	}
}

func TestEscapeForPowerShell_PreservesNewlines(t *testing.T) {
	input := "line1\nline2\r\n"
	got := EscapeForPowerShell(input)
	if got != input {
		t.Errorf("EscapeForPowerShell() = %q, want %q (newlines should be preserved)", got, input)
	}
}
