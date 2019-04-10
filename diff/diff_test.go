package diff

import (
	"testing"
)

func Test_NoDiff_TextDiff(t *testing.T) {
	expected := `hello
	test text`
	actual := `hello
	test text`

	diff := TextDiff(expected, actual)
	if diff != "" {
		t.Errorf("Strings were not equal\n%s", diff)
	}
}

func Test_Diff_TextDiff(t *testing.T) {
	expected := `hello
	test text`
	actual := `hello
	test`

	diff := TextDiff(expected, actual)

	expectedDiff := `
--- expected
+++ actual
@@ -1,2 +1,2 @@
 hello
-	test text+	test`

	if diff != expectedDiff {
		t.Errorf("Diff was not equal to expected diff\nExpected:\n%s\n\nActual:\n%s", expectedDiff, diff)
	}
}

func Test_NoDiff_ContextDiff(t *testing.T) {
	expected := `hello
	test text`
	actual := `hello
	test text`

	diff := ContextDiff(expected, actual)
	if diff != "" {
		t.Errorf("Strings were not equal\n%s", diff)
	}
}

func Test_Diff_ContextDiff(t *testing.T) {
	expected := `hello
	test text`
	actual := `hello
	test`

	diff := ContextDiff(expected, actual)

	expectedDiff := `
*** expected
--- actual
***************
*** 1,2 ****
  hello
! 	test text--- 1,2 ----
  hello
! 	test`

	if diff != expectedDiff {
		t.Errorf("Diff was not equal to expected diff\nExpected:\n%s\n\nActual:\n%s", expectedDiff, diff)
	}
}
