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

package prompt

import (
	"bytes"
	"strings"
	"testing"
)

func TestConfirmText_Answers(t *testing.T) {
	const question = `Type "host1" to remove it:`

	tests := []struct {
		name     string
		input    string
		expected string
		want     bool
		asked    int
	}{
		{"exact match", "host1\n", "host1", true, 1},
		{"case and whitespace ignored", "  HOST1 \r\n", "host1", true, 1},
		{"expected whitespace ignored", "host1\n", " host1 ", true, 1},
		{"last line without newline", "host1", "host1", true, 1},
		{"EOF is no", "", "host1", false, 1},
		{"empty line asks again", "\nhost1\n", "host1", true, 2},
		{"mismatch asks again", "host2\nhost1\n", "host1", true, 2},
		{"mismatch then EOF", "host2\n", "host1", false, 2},
		{"retry limit is no", "a\nb\nc\nhost1\n", "host1", false, 3},
		{"empty expected is never matched", "\n\n\n", "", false, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer

			got := ConfirmText(strings.NewReader(tt.input), &out, question, tt.expected)
			if got != tt.want {
				t.Errorf("ConfirmText(%q, %q) = %v, want %v", tt.input, tt.expected, got, tt.want)
			}

			if asked := strings.Count(out.String(), question); asked != tt.asked {
				t.Errorf("question asked %d times, want %d; output %q", asked, tt.asked, out.String())
			}
		})
	}
}

func TestConfirmText_Output(t *testing.T) {
	var out bytes.Buffer

	ConfirmText(strings.NewReader("host2\nhost1\n"), &out, "Type host1:", "host1")

	want := "Type host1: The text did not match.\nType host1: "
	if out.String() != want {
		t.Errorf("wrote %q, want %q", out.String(), want)
	}
}
