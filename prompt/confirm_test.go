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

func TestConfirm_Answers(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		defaultYes bool
		want       bool
		asked      int
	}{
		{"y", "y\n", false, true, 1},
		{"yes", "yes\n", false, true, 1},
		{"upper case with CRLF", "YES\r\n", false, true, 1},
		{"surrounding whitespace", "  y  \n", false, true, 1},
		{"n with default yes", "n\n", true, false, 1},
		{"no with default yes", "No\n", true, false, 1},
		{"empty line takes default no", "\n", false, false, 1},
		{"empty line takes default yes", "\n", true, true, 1},
		{"blank line takes default yes", "   \r\n", true, true, 1},
		{"EOF is no with default no", "", false, false, 1},
		{"EOF is no even with default yes", "", true, false, 1},
		{"last line without newline", "y", false, true, 1},
		{"unrecognized then yes", "yep\ny\n", false, true, 2},
		{"unrecognized then empty takes default", "maybe\n\n", true, true, 2},
		{"unrecognized then EOF", "yep\n", true, false, 2},
		{"retry limit is no", "a\nb\nc\ny\n", true, false, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer

			got := Confirm(strings.NewReader(tt.input), &out, "Proceed?", tt.defaultYes)
			if got != tt.want {
				t.Errorf("Confirm(%q, defaultYes=%v) = %v, want %v", tt.input, tt.defaultYes, got, tt.want)
			}

			if asked := strings.Count(out.String(), "Proceed?"); asked != tt.asked {
				t.Errorf("question asked %d times, want %d; output %q", asked, tt.asked, out.String())
			}
		})
	}
}

func TestConfirm_SuffixShowsDefault(t *testing.T) {
	tests := []struct {
		want       string
		defaultYes bool
	}{
		{"Proceed? [y/N] ", false},
		{"Proceed? [Y/n] ", true},
	}

	for _, tt := range tests {
		var out bytes.Buffer

		Confirm(strings.NewReader("y\n"), &out, "Proceed?", tt.defaultYes)

		if out.String() != tt.want {
			t.Errorf("defaultYes=%v wrote %q, want %q", tt.defaultYes, out.String(), tt.want)
		}
	}
}

func TestConfirm_RetryExplainsAnswers(t *testing.T) {
	var out bytes.Buffer

	Confirm(strings.NewReader("yep\ny\n"), &out, "Proceed?", false)

	want := "Proceed? [y/N] Please answer yes or no.\nProceed? [y/N] "
	if out.String() != want {
		t.Errorf("wrote %q, want %q", out.String(), want)
	}
}
