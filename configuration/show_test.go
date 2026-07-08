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

package configuration

import (
	"strings"
	"testing"
)

func TestShow_RendersSavedConfigAsYAML(t *testing.T) {
	setHome(t)

	f := New[testConfig]("mytool.yml")

	if err := f.Save(&testConfig{Name: "acme", Count: 2, Items: []string{"a", "b"}}); err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}

	out, err := f.Show()
	if err != nil {
		t.Fatalf("Show() returned error: %v", err)
	}

	for _, want := range []string{"name: acme", "count: 2", "- a", "- b"} {
		if !strings.Contains(out, want) {
			t.Errorf("Show() = %q, want it to contain %q", out, want)
		}
	}
}

func TestShow_MissingFileRendersZeroValue(t *testing.T) {
	setHome(t)

	out, err := New[testConfig]("mytool.yml").Show()
	if err != nil {
		t.Fatalf("Show() returned error: %v", err)
	}

	if !strings.Contains(out, "name:") || !strings.Contains(out, "count: 0") {
		t.Errorf("Show() of missing file = %q, want zero-valued YAML", out)
	}
}
