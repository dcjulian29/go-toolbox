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

func TestXMLEscape_Ampersand(t *testing.T) {
	got := XMLEscape("fish & chips")
	want := "fish &amp; chips"
	if got != want {
		t.Errorf("XMLEscape() = %q, want %q", got, want)
	}
}

func TestXMLEscape_LessThan(t *testing.T) {
	got := XMLEscape("a < b")
	want := "a &lt; b"
	if got != want {
		t.Errorf("XMLEscape() = %q, want %q", got, want)
	}
}

func TestXMLEscape_GreaterThan(t *testing.T) {
	got := XMLEscape("a > b")
	want := "a &gt; b"
	if got != want {
		t.Errorf("XMLEscape() = %q, want %q", got, want)
	}
}

func TestXMLEscape_DoubleQuote(t *testing.T) {
	got := XMLEscape("say \"hello\"")
	want := "say &quot;hello&quot;"
	if got != want {
		t.Errorf("XMLEscape() = %q, want %q", got, want)
	}
}

func TestXMLEscape_Apostrophe(t *testing.T) {
	got := XMLEscape("it's")
	want := "it&apos;s"
	if got != want {
		t.Errorf("XMLEscape() = %q, want %q", got, want)
	}
}

func TestXMLEscape_AllFiveEntities(t *testing.T) {
	got := XMLEscape("&<>\"'")
	want := "&amp;&lt;&gt;&quot;&apos;"
	if got != want {
		t.Errorf("XMLEscape() = %q, want %q", got, want)
	}
}

func TestXMLEscape_NoSpecialChars(t *testing.T) {
	input := "hello world 123"
	got := XMLEscape(input)
	if got != input {
		t.Errorf("XMLEscape() = %q, want %q (no escaping needed)", got, input)
	}
}

func TestXMLEscape_EmptyString(t *testing.T) {
	got := XMLEscape("")
	if got != "" {
		t.Errorf("XMLEscape() = %q, want empty string", got)
	}
}

func TestXMLEscape_PreservesExistingEntity(t *testing.T) {
	// An already-escaped entity reference should not be double-escaped.
	got := XMLEscape("&lt;")
	want := "&lt;"
	if got != want {
		t.Errorf("XMLEscape() = %q, want %q (existing entity should be preserved)", got, want)
	}
}

func TestXMLEscape_PreservesAllExistingEntities(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"amp", "&amp;"},
		{"lt", "&lt;"},
		{"gt", "&gt;"},
		{"quot", "&quot;"},
		{"apos", "&apos;"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := XMLEscape(tt.input)
			if got != tt.input {
				t.Errorf("XMLEscape(%q) = %q, want %q (existing entity should be preserved)", tt.input, got, tt.input)
			}
		})
	}
}

func TestXMLEscape_Idempotent(t *testing.T) {
	input := XMLEscape("a & b")
	got := XMLEscape(input)
	if got != input {
		t.Errorf("XMLEscape should be idempotent, got %q, want %q", got, input)
	}
}

func TestXMLEscape_IdempotentAllEntities(t *testing.T) {
	original := "<div class=\"main\">it's a & b</div>"
	first := XMLEscape(original)
	second := XMLEscape(first)
	third := XMLEscape(second)

	if second != first {
		t.Errorf("second pass changed output:\n  first:  %q\n  second: %q", first, second)
	}

	if third != first {
		t.Errorf("third pass changed output:\n  first: %q\n  third: %q", first, third)
	}
}

func TestXMLEscape_BareAmpersandEscaped(t *testing.T) {
	// An ampersand NOT followed by a known entity suffix should be escaped.
	got := XMLEscape("&foo;")
	want := "&amp;foo;"
	if got != want {
		t.Errorf("XMLEscape() = %q, want %q (bare ampersand should be escaped)", got, want)
	}
}

func TestXMLEscape_MixedBareAndEntityAmpersands(t *testing.T) {
	got := XMLEscape("a & b &amp; c & d")
	want := "a &amp; b &amp; c &amp; d"
	if got != want {
		t.Errorf("XMLEscape() = %q, want %q", got, want)
	}
}

func TestXMLEscape_XMLTag(t *testing.T) {
	got := XMLEscape("<div class=\"main\">")
	want := "&lt;div class=&quot;main&quot;&gt;"
	if got != want {
		t.Errorf("XMLEscape() = %q, want %q", got, want)
	}
}

func TestXMLEscape_PreservesWhitespace(t *testing.T) {
	input := "line1\nline2\ttab"
	got := XMLEscape(input)
	if got != input {
		t.Errorf("XMLEscape() = %q, want %q (whitespace should be preserved)", got, input)
	}
}

func TestXMLEscape_Unicode(t *testing.T) {
	input := "caf\u00e9 & cr\u00e8me"
	got := XMLEscape(input)
	want := "caf\u00e9 &amp; cr\u00e8me"
	if got != want {
		t.Errorf("XMLEscape() = %q, want %q", got, want)
	}
}

func TestXMLEscape_OutputLongerThanInput(t *testing.T) {
	input := "&"
	got := XMLEscape(input)
	if len(got) <= len(input) {
		t.Errorf("escaped output %q should be longer than input %q", got, input)
	}
}

func TestXMLEscape_ConsecutiveAmpersands(t *testing.T) {
	got := XMLEscape("&&&&")
	want := "&amp;&amp;&amp;&amp;"
	if got != want {
		t.Errorf("XMLEscape() = %q, want %q", got, want)
	}
}
