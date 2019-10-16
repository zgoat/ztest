// Package diff has some helpers for text diffs.
package diff // import "zgo.at/ztest/diff"

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/sqm/go-difflib/difflib"
)

// Diff returns a unified diff of the two passed arguments, or "" if they are
// the same.
func Diff(expected, actual interface{}) string {
	if reflect.DeepEqual(expected, actual) {
		return ""
	}

	scs := spew.ConfigState{
		Indent:                  "  ",
		DisableMethods:          true,
		DisablePointerAddresses: true,
		DisableCapacities:       true,
		SortKeys:                true,
	}

	udiff := difflib.UnifiedDiff{
		A:        strings.SplitAfter(scs.Sdump(expected), "\n"),
		FromFile: "expected",
		B:        strings.SplitAfter(scs.Sdump(actual), "\n"),
		ToFile:   "actual",
		Context:  2,
	}
	diff, err := difflib.GetUnifiedDiffString(udiff)
	if err != nil {
		panic(fmt.Sprintf("Error producing diff: %s\n", err))
	}

	if diff != "" {
		diff = "\n" + diff
	}
	return diff
}
