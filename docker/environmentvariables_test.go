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

package docker

import (
	"testing"
)

func TestEnvironmentVariablesWithPrefix_MatchingVars(t *testing.T) {
	t.Setenv("MYAPP_DB_HOST", "localhost")
	t.Setenv("MYAPP_DB_PORT", "5432")
	t.Setenv("OTHER_VAR", "ignored")

	result := EnvironmentVariablesWithPrefix("MYAPP")

	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d: %v", len(result), result)
	}

	if result["MYAPP_DB_HOST"] != "localhost" {
		t.Errorf("MYAPP_DB_HOST = %q, want %q", result["MYAPP_DB_HOST"], "localhost")
	}

	if result["MYAPP_DB_PORT"] != "5432" {
		t.Errorf("MYAPP_DB_PORT = %q, want %q", result["MYAPP_DB_PORT"], "5432")
	}
}

func TestEnvironmentVariablesWithPrefix_PrefixWithTrailingUnderscore(t *testing.T) {
	t.Setenv("MYAPP_KEY", "value")

	result := EnvironmentVariablesWithPrefix("MYAPP_")

	if result["MYAPP_KEY"] != "value" {
		t.Errorf("MYAPP_KEY = %q, want %q", result["MYAPP_KEY"], "value")
	}
}

func TestEnvironmentVariablesWithPrefix_NoMatch(t *testing.T) {
	t.Setenv("UNRELATED_VAR", "value")

	result := EnvironmentVariablesWithPrefix("MYAPP")

	if len(result) != 0 {
		t.Errorf("expected empty map, got %v", result)
	}
}

func TestEnvironmentVariablesWithPrefix_EmptyPrefix(t *testing.T) {
	t.Setenv("ANYTHING", "value")

	result := EnvironmentVariablesWithPrefix("")

	if len(result) != 0 {
		t.Errorf("empty prefix should return empty map, got %v", result)
	}
}

func TestEnvironmentVariablesWithPrefix_PrefixAmbiguityResolved(t *testing.T) {
	t.Setenv("APP_REAL", "yes")
	t.Setenv("APPLESAUCE", "no")

	result := EnvironmentVariablesWithPrefix("APP")

	if _, ok := result["APPLESAUCE"]; ok {
		t.Error("APPLESAUCE should not match prefix 'APP' (underscore enforced)")
	}

	if result["APP_REAL"] != "yes" {
		t.Errorf("APP_REAL = %q, want %q", result["APP_REAL"], "yes")
	}
}

func TestEnvironmentVariablesWithPrefix_ValueContainsEquals(t *testing.T) {
	t.Setenv("MYAPP_CONN", "host=localhost;port=5432")

	result := EnvironmentVariablesWithPrefix("MYAPP")

	if result["MYAPP_CONN"] != "host=localhost;port=5432" {
		t.Errorf("MYAPP_CONN = %q, want %q", result["MYAPP_CONN"], "host=localhost;port=5432")
	}
}

func TestEnvironmentVariablesWithStrippedPrefix_KeysStripped(t *testing.T) {
	t.Setenv("MYAPP_DB_HOST", "localhost")
	t.Setenv("MYAPP_DB_PORT", "5432")

	result := EnvironmentVariablesWithStrippedPrefix("MYAPP")

	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d: %v", len(result), result)
	}

	if result["DB_HOST"] != "localhost" {
		t.Errorf("DB_HOST = %q, want %q", result["DB_HOST"], "localhost")
	}

	if result["DB_PORT"] != "5432" {
		t.Errorf("DB_PORT = %q, want %q", result["DB_PORT"], "5432")
	}
}

func TestEnvironmentVariablesWithStrippedPrefix_OriginalKeyAbsent(t *testing.T) {
	t.Setenv("MYAPP_KEY", "value")

	result := EnvironmentVariablesWithStrippedPrefix("MYAPP")

	if _, ok := result["MYAPP_KEY"]; ok {
		t.Error("original prefixed key should not be present after stripping")
	}

	if result["KEY"] != "value" {
		t.Errorf("KEY = %q, want %q", result["KEY"], "value")
	}
}

func TestEnvironmentVariablesWithStrippedPrefix_EmptyPrefix(t *testing.T) {
	result := EnvironmentVariablesWithStrippedPrefix("")

	if len(result) != 0 {
		t.Errorf("empty prefix should return empty map, got %v", result)
	}
}

func TestEnvironmentVariablesWithStrippedPrefix_NoMatch(t *testing.T) {
	t.Setenv("UNRELATED_VAR", "value")

	result := EnvironmentVariablesWithStrippedPrefix("MYAPP")

	if len(result) != 0 {
		t.Errorf("expected empty map, got %v", result)
	}
}

func TestEnvironmentVariablesWithStrippedPrefix_TrailingUnderscore(t *testing.T) {
	t.Setenv("MYAPP_KEY", "value")

	result := EnvironmentVariablesWithStrippedPrefix("MYAPP_")

	if result["KEY"] != "value" {
		t.Errorf("KEY = %q, want %q", result["KEY"], "value")
	}
}
