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
)

// XMLEscape escapes the five predefined XML entities in s. The function is
// idempotent — already-escaped entities are not double-escaped.
func XMLEscape(s string) string {
	var b strings.Builder
	b.Grow(len(s))

	for i := 0; i < len(s); i++ { //nolint:revive
		switch s[i] {
		case '&':
			if isXMLEntity(s[i:]) {
				end := strings.IndexByte(s[i:], ';')
				b.WriteString(s[i : i+end+1])
				i += end
			} else {
				b.WriteString("&amp;")
			}

		case '<':
			b.WriteString("&lt;")

		case '>':
			b.WriteString("&gt;")

		case '"':
			b.WriteString("&quot;")

		case '\'':
			b.WriteString("&apos;")

		default:
			b.WriteByte(s[i])
		}
	}

	return b.String()
}

func isXMLEntity(s string) bool {
	return strings.HasPrefix(s, "&amp;") ||
		strings.HasPrefix(s, "&lt;") ||
		strings.HasPrefix(s, "&gt;") ||
		strings.HasPrefix(s, "&quot;") ||
		strings.HasPrefix(s, "&apos;")
}
