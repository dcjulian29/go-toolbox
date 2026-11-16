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

// Package prompt asks the user to confirm an action before it is taken.
package prompt

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// maxAttempts caps how many lines a prompt reads before giving up with a no,
// so a pipe of unrecognized input cannot loop forever.
const maxAttempts = 3

// ask writes text to out and reads lines from in until decide settles one or
// maxAttempts lines have been read; the limit is a no. decide receives the
// line with surrounding whitespace (including the CR of a Windows CRLF)
// trimmed, and returns the answer and whether the line settled it. Before
// each re-ask, retry is written on its own line.
//
// Abandoned input, meaning EOF or a read error with nothing typed, is always
// a no and never reaches decide, so an unattended run cannot fall through to
// a prompt's default.
func ask(in io.Reader, out io.Writer, text, retry string, decide func(line string) (answer, settled bool)) bool {
	// One reader for every attempt: a fresh bufio.Reader per attempt would
	// drop whatever the previous one had already buffered from a pipe.
	reader := bufio.NewReader(in)

	for attempt := range maxAttempts {
		if attempt > 0 {
			_, _ = fmt.Fprintln(out, retry)
		}

		_, _ = fmt.Fprint(out, text)

		line, err := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		// A final line without a newline still counts as typed input; only
		// an empty read at EOF is abandoned.
		if err != nil && line == "" {
			return false
		}

		if answer, settled := decide(line); settled {
			return answer
		}
	}

	return false
}
