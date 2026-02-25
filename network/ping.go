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
	"net"
	"os/exec"
	"runtime"
	"strings"
)

func Ping(address string) bool {
	addr, err := net.ResolveIPAddr("ip", address)
	if err != nil {
		return false
	}

	ip := net.ParseIP(addr.String())
	if ip == nil {
		return false
	}

	var output []byte

	if runtime.GOOS == "windows" {
		output, _ = exec.Command("ping", "-w", "1000", "-n", "1", ip.String()).CombinedOutput() // #nosec G204
	} else {
		output, _ = exec.Command("ping", "-c", "1", ip.String()).CombinedOutput() // #nosec G204
	}

	if strings.Contains(string(output[:]), "TTL") {
		return true
	}

	return false
}
