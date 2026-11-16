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
	"io"
	"strings"
)

// ConfirmText writes question to out and reads a line from in, returning true
// only when the line matches expected, ignoring case and surrounding
// whitespace. It is meant for destructive actions, where the user types the
// name of the thing being removed:
//
//	prompt.ConfirmText(in, out, `Type "host1" to remove it:`, "host1")
//
// There is no default. An empty line or a mismatch asks again, up to three
// times in all, and then the answer is no; so is no input at all (EOF or a
// closed stdin). An empty expected can never be matched, so it is always no.
//
// Cobra commands should pass cmd.InOrStdin() and cmd.OutOrStdout().
func ConfirmText(in io.Reader, out io.Writer, question, expected string) bool {
	want := strings.TrimSpace(expected)

	return ask(in, out, question+" ", "The text did not match.", func(line string) (bool, bool) {
		// An empty line is unsettled rather than a no, so a stray Enter
		// re-asks instead of abandoning the prompt.
		matched := line != "" && strings.EqualFold(line, want)

		return matched, matched
	})
}
