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

package network

import (
	"testing"
)

func TestPing_Loopback(t *testing.T) {
	if !Ping("127.0.0.1") {
		t.Error("Ping(127.0.0.1) should succeed on loopback")
	}
}

func TestPing_Localhost(t *testing.T) {
	if !Ping("localhost") {
		t.Error("Ping(localhost) should succeed")
	}
}

func TestPing_EmptyString(t *testing.T) {
	if Ping("") {
		t.Error("Ping(\"\") should return false for empty string")
	}
}

func TestPing_UnresolvableHostname(t *testing.T) {
	if Ping("host.invalid.") {
		t.Error("Ping should return false for unresolvable hostname")
	}
}

func TestPing_InvalidIPAddress(t *testing.T) {
	if Ping("999.999.999.999") {
		t.Error("Ping should return false for invalid IP address")
	}
}

func TestPing_NonRoutableAddress(t *testing.T) {
	// 192.0.2.0/24 is TEST-NET-1, reserved for documentation (RFC 5737).
	// Should be unreachable and timeout quickly.
	if Ping("192.0.2.1") {
		t.Error("Ping should return false for non-routable documentation address")
	}
}

func TestPing_GarbageInput(t *testing.T) {
	if Ping("not-a-valid-address-at-all!!!") {
		t.Error("Ping should return false for garbage input")
	}
}

func TestPing_WhitespaceOnly(t *testing.T) {
	if Ping("   ") {
		t.Error("Ping should return false for whitespace-only input")
	}
}

func TestPing_LoopbackIPv6(t *testing.T) {
	if !Ping("::1") {
		t.Skip("IPv6 loopback ping failed — IPv6 may not be available in this environment")
	}
}
