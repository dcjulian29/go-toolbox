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

// Confirm writes question to out with a suffix showing the default, "[Y/n]"
// when defaultYes is true and "[y/N]" otherwise, and reads the answer from in.
//
// "y" or "yes" is a yes and "n" or "no" is a no, in any case and ignoring
// surrounding whitespace. An empty line takes the default. Any other answer
// is never guessed at: the question is asked again, up to three times in all,
// and then the answer is no.
//
// No input at all (EOF or a closed stdin) is a no even when defaultYes is
// true, so an unattended run never proceeds without an answer.
//
// Cobra commands should pass cmd.InOrStdin() and cmd.OutOrStdout().
//
//nolint:revive // defaultYes selects the default answer; that is the point of the parameter.
func Confirm(in io.Reader, out io.Writer, question string, defaultYes bool) bool {
	suffix := " [y/N] "
	if defaultYes {
		suffix = " [Y/n] "
	}

	return ask(in, out, question+suffix, "Please answer yes or no.", func(line string) (bool, bool) {
		switch strings.ToLower(line) {
		case "":
			return defaultYes, true
		case "y", "yes":
			return true, true
		case "n", "no":
			return false, true
		default:
			return false, false
		}
	})
}
