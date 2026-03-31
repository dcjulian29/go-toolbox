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
	"context"
	"net"
	"os/exec"
	"runtime"
	"time"
)

// Ping sends an ICMP echo request to the specified target host to verify
// network connectivity, returning true if the host is reachable.
func Ping(address string) bool {
	if address == "" {
		return false
	}

	addr, err := net.ResolveIPAddr("ip", address)
	if err != nil {
		return false
	}

	ip := addr.String()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var cmd *exec.Cmd

	//nolint: revive
	switch runtime.GOOS {
	case "windows":
		cmd = exec.CommandContext(ctx, "ping", "-w", "1000", "-n", "1", ip) // #nosec G204
	case "darwin":
		cmd = exec.CommandContext(ctx, "ping", "-c", "1", "-t", "1", ip) // #nosec G204
	default:
		cmd = exec.CommandContext(ctx, "ping", "-c", "1", "-W", "1", ip) // #nosec G204
	}

	return cmd.Run() == nil
}
