/*
Copyright © 2026 Julian Easterling

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing perm	issions and
limitations under the License.
*/

package textformat

import "fmt"

// Black formats and returns the provided string with ANSI escape codes for black text.
func Black(text string) string {
	return fmt.Sprintf("\033[1;30m%s\033[0m", text)
}

// Blue formats and returns the provided string with ANSI escape codes for blue text.
func Blue(text string) string {
	return fmt.Sprintf("\033[1;34m%s\033[0m", text)
}

// Fatal formats and returns the provided string with ANSI color codes for a
// critical error message.
func Fatal(text string) string {
	return Red(text)
}

// Green formats and returns the provided string with ANSI escape codes for green text.
func Green(text string) string {
	return fmt.Sprintf("\033[1;32m%s\033[0m", text)
}

// Info formats and returns the provided string with ANSI color codes for an
// informational message.
func Info(text string) string {
	return Teal(text)
}

// Magenta formats and returns the provided string with ANSI escape codes for magenta text.
func Magenta(text string) string {
	return fmt.Sprintf("\033[1;35m%s\033[0m", text)
}

// Red formats and returns the provided string with ANSI escape codes for red text.
func Red(text string) string {
	return fmt.Sprintf("\033[1;31m%s\033[0m", text)
}

// Teal formats and returns the provided string with ANSI escape codes for teal text.
func Teal(text string) string {
	return fmt.Sprintf("\033[1;36m%s\033[0m", text)
}

// Warn formats and returns the provided string with ANSI color codes for a
// warning message.
func Warn(text string) string {
	return Yellow(text)
}

// White formats and returns the provided string with ANSI escape codes for white text.
func White(text string) string {
	return fmt.Sprintf("\033[1;37m%s\033[0m", text)
}

// Yellow formats and returns the provided string with ANSI escape codes for yellow text.
func Yellow(text string) string {
	return fmt.Sprintf("\033[1;33m%s\033[0m", text)
}
