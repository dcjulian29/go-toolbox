//go:build windows

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

package hypervdisk

import (
	"strings"
	"testing"
)

// --- replaceTemplateMarker tests ---

func Test_replaceTemplateMarker_ReplacesWhenMarkerPresentAndValueProvided(t *testing.T) {
	content := `<ComputerName>{{COMPUTERNAME}}</ComputerName>`
	got, err := replaceTemplateMarker(content, "{{COMPUTERNAME}}", "MYPC", "computer name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(got, "MYPC") {
		t.Errorf("expected replacement value in output, got: %s", got)
	}

	if strings.Contains(got, "{{COMPUTERNAME}}") {
		t.Error("marker should have been replaced")
	}
}

func Test_replaceTemplateMarker_ReturnsErrorWhenMarkerPresentAndValueEmpty(t *testing.T) {
	content := `<ComputerName>{{COMPUTERNAME}}</ComputerName>`
	_, err := replaceTemplateMarker(content, "{{COMPUTERNAME}}", "", "computer name")
	if err == nil {
		t.Fatal("expected error for empty value with marker present")
	}

	if !strings.Contains(err.Error(), "computer name") {
		t.Errorf("error should mention field name, got: %v", err)
	}

	if !strings.Contains(err.Error(), "{{COMPUTERNAME}}") {
		t.Errorf("error should mention marker, got: %v", err)
	}
}

func Test_replaceTemplateMarker_ReturnsErrorWhenMarkerPresentAndValueWhitespace(t *testing.T) {
	content := `<User>{{USER}}</User>`
	_, err := replaceTemplateMarker(content, "{{USER}}", "   ", "user name")
	if err == nil {
		t.Fatal("expected error for whitespace-only value")
	}
}

func Test_replaceTemplateMarker_NoOpWhenMarkerAbsent(t *testing.T) {
	content := `<ComputerName>HARDCODED</ComputerName>`
	got, err := replaceTemplateMarker(content, "{{COMPUTERNAME}}", "", "computer name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != content {
		t.Errorf("content should be unchanged, got: %s", got)
	}
}

func Test_replaceTemplateMarker_ReplacesMultipleOccurrences(t *testing.T) {
	content := `{{USER}} is the admin. Contact {{USER}} for help.`
	got, err := replaceTemplateMarker(content, "{{USER}}", "admin", "user name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(got, "{{USER}}") {
		t.Error("all occurrences of the marker should have been replaced")
	}

	if strings.Count(got, "admin") < 2 {
		t.Errorf("expected at least 2 replacements, got: %s", got)
	}
}

func Test_replaceTemplateMarker_ReturnsEmptyStringOnError(t *testing.T) {
	content := `<User>{{USER}}</User>`
	got, err := replaceTemplateMarker(content, "{{USER}}", "", "user name")
	if err == nil {
		t.Fatal("expected error")
	}

	if got != "" {
		t.Errorf("expected empty string on error, got: %s", got)
	}
}

// --- templateReplacements tests ---

func Test_templateReplacements_AllMarkersReplaced(t *testing.T) {
	raw := []byte(`<u>{{COMPUTERNAME}}</u><p>{{PASSWORD}}</p><n>{{USER}}</n>`)
	cfg := &InjectConfig{
		ComputerName: "TESTPC",
		UserPassword: "Secret123",
		UserName:     "admin",
	}

	got, err := templateReplacements(raw, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result := string(got)

	//nolint
	if strings.Contains(result, "{{COMPUTERNAME}}") ||
		strings.Contains(result, "{{PASSWORD}}") ||
		strings.Contains(result, "{{USER}}") {
		t.Errorf("unreplaced markers remain: %s", result)
	}
}

func Test_templateReplacements_NoMarkersInTemplate(t *testing.T) {
	raw := []byte(`<settings><timezone>UTC</timezone></settings>`)
	cfg := &InjectConfig{}

	got, err := templateReplacements(raw, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(got) != string(raw) {
		t.Error("content should be unchanged when no markers present")
	}
}

func Test_templateReplacements_ErrorOnEmptyComputerName(t *testing.T) {
	raw := []byte(`<name>{{COMPUTERNAME}}</name>`)
	cfg := &InjectConfig{
		ComputerName: "",
		UserPassword: "pass",
		UserName:     "user",
	}

	_, err := templateReplacements(raw, cfg)
	if err == nil {
		t.Fatal("expected error for empty computer name")
	}

	if !strings.Contains(err.Error(), "computer name") {
		t.Errorf("error should mention computer name, got: %v", err)
	}
}

func Test_templateReplacements_ErrorOnEmptyPassword(t *testing.T) {
	raw := []byte(`<pass>{{PASSWORD}}</pass>`)
	cfg := &InjectConfig{
		ComputerName: "PC",
		UserPassword: "",
		UserName:     "user",
	}

	_, err := templateReplacements(raw, cfg)
	if err == nil {
		t.Fatal("expected error for empty password")
	}

	if !strings.Contains(err.Error(), "user password") {
		t.Errorf("error should mention user password, got: %v", err)
	}
}

func Test_templateReplacements_ErrorOnEmptyUserName(t *testing.T) {
	raw := []byte(`<user>{{USER}}</user>`)
	cfg := &InjectConfig{
		ComputerName: "PC",
		UserPassword: "pass",
		UserName:     "",
	}

	_, err := templateReplacements(raw, cfg)
	if err == nil {
		t.Fatal("expected error for empty user name")
	}

	if !strings.Contains(err.Error(), "user name") {
		t.Errorf("error should mention user name, got: %v", err)
	}
}

func Test_templateReplacements_PartialMarkers(t *testing.T) {
	raw := []byte(`<name>{{COMPUTERNAME}}</name><tz>UTC</tz>`)
	cfg := &InjectConfig{
		ComputerName: "MYPC",
	}

	got, err := templateReplacements(raw, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(string(got), "{{COMPUTERNAME}}") {
		t.Error("marker should have been replaced")
	}
}

func Test_templateReplacements_ReturnsNilOnError(t *testing.T) {
	raw := []byte(`<name>{{COMPUTERNAME}}</name>`)
	cfg := &InjectConfig{ComputerName: ""}

	got, err := templateReplacements(raw, cfg)
	if err == nil {
		t.Fatal("expected error")
	}

	if got != nil {
		t.Errorf("expected nil on error, got: %v", got)
	}
}

// --- InjectUnattendFile guard tests ---

func Test_InjectUnattendFile_NilConfig(t *testing.T) {
	err := InjectUnattendFile(nil)
	if err == nil {
		t.Fatal("expected error for nil config")
	}

	if !strings.Contains(err.Error(), "nil") {
		t.Errorf("error should mention nil, got: %v", err)
	}
}

func Test_InjectUnattendFile_EmptyTemplate(t *testing.T) {
	cfg := &InjectConfig{UnattendTemplate: ""}

	err := InjectUnattendFile(cfg)
	if err == nil {
		t.Fatal("expected error for empty template")
	}

	if !strings.Contains(err.Error(), "empty") {
		t.Errorf("error should mention empty, got: %v", err)
	}
}

func Test_InjectUnattendFile_WhitespaceTemplate(t *testing.T) {
	cfg := &InjectConfig{UnattendTemplate: "   "}

	err := InjectUnattendFile(cfg)
	if err == nil {
		t.Fatal("expected error for whitespace-only template")
	}
}

func Test_InjectUnattendFile_TemplateMissing(t *testing.T) {
	cfg := &InjectConfig{UnattendTemplate: `C:\nonexistent\template.xml`}

	err := InjectUnattendFile(cfg)
	if err == nil {
		t.Fatal("expected error for missing template file")
	}

	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("error should mention not found, got: %v", err)
	}
}
