// Package test contains various small helper functions that are useful when
// writing tests.
package test

import "strings"

// ErrorContains checks if the error message in got contains the text in
// expected.
//
// This is safe when got is nil. Use an empty string for expected if you want to
// test that err is nil.
func ErrorContains(got error, expected string) bool {
	if got == nil {
		return expected == ""
	}
	if expected == "" {
		return false
	}
	return strings.Contains(got.Error(), expected)
}
