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

package execute

import (
	"context"
	"os"
	"os/exec"
)

// ExternalProgramContext runs the specified external program with the given parameters,
// binding its standard input, output, and error streams directly to the host OS standard
// streams. If the context is cancelled or its deadline expires, the process is killed.
func ExternalProgramContext(ctx context.Context, program string, params ...string) error {
	cmd := exec.CommandContext(ctx, program, params...) // #nosec G204
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	return cmd.Run()
}
