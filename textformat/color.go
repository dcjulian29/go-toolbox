package textformat

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

import "fmt"

// Black formats and returns the provided string with ANSI escape codes for black text
func Black(colorString string) string {
	c := color("\033[1;30m%s\033[0m")

	return c(colorString)
}

// Fatal formats and outputs a critical error message using ANSI color codes,
// typically used to indicate an unrecoverable error before terminating
// the application.
func Fatal(colorString string) string {
	return Red(colorString)
}

// Green formats and returns the provided string with ANSI escape codes for green text
func Green(colorString string) string {
	c := color("\033[1;32m%s\033[0m")

	return c(colorString)
}

// Info formats and outputs an informational message using ANSI color codes,
// typically used for standard program output or logging general execution steps.
func Info(colorString string) string {
	return Teal(colorString)
}

// Magenta formats and returns the provided string with ANSI escape codes for magenta text
func Magenta(colorString string) string {
	c := color("\033[1;35m%s\033[0m")

	return c(colorString)
}

// Purple formats and returns the provided string with ANSI escape codes for purple text
func Purple(colorString string) string {
	c := color("\033[1;34m%s\033[0m")

	return c(colorString)
}

// Red formats and returns the provided string with ANSI escape codes for red text
func Red(colorString string) string {
	c := color("\033[1;31m%s\033[0m")

	return c(colorString)
}

// Teal formats and returns the provided string with ANSI escape codes for teal text
func Teal(colorString string) string {
	c := color("\033[1;36m%s\033[0m")

	return c(colorString)
}

// Warn formats and outputs a warning message using ANSI color codes,
// typically used to highlight non-critical issues or potential problems.
func Warn(colorString string) string {
	return Yellow(colorString)
}

// White formats and returns the provided string with ANSI escape codes for white text
func White(colorString string) string {
	c := color("\033[1;37m%s\033[0m")

	return c(colorString)
}

// Yellow formats and returns the provided string with ANSI escape codes for yellow text
func Yellow(colorString string) string {
	c := color("\033[1;33m%s\033[0m")

	return c(colorString)
}

func color(colorString string) func(...any) string {
	sprint := func(args ...any) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}

	return sprint
}
