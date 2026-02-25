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
package color

import "fmt"

func Black(colorString string) string {
	c := Color("\033[1;30m%s\033[0m")

	return c(colorString)
}

func Color(colorString string) func(...any) string {
	sprint := func(args ...any) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}

	return sprint
}

func Fatal(colorString string) string {
	return Red(colorString)
}

func Green(colorString string) string {
	c := Color("\033[1;32m%s\033[0m")

	return c(colorString)
}

func Info(colorString string) string {
	return Teal(colorString)
}

func Magenta(colorString string) string {
	c := Color("\033[1;35m%s\033[0m")

	return c(colorString)
}

func Purple(colorString string) string {
	c := Color("\033[1;34m%s\033[0m")

	return c(colorString)
}

func Red(colorString string) string {
	c := Color("\033[1;31m%s\033[0m")

	return c(colorString)
}

func Teal(colorString string) string {
	c := Color("\033[1;36m%s\033[0m")

	return c(colorString)
}

func Warn(colorString string) string {
	return Yellow(colorString)
}

func White(colorString string) string {
	c := Color("\033[1;37m%s\033[0m")

	return c(colorString)
}

func Yellow(colorString string) string {
	c := Color("\033[1;33m%s\033[0m")

	return c(colorString)
}
