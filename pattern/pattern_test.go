package pattern

import (
	"testing"
)

func TestPatternCreation(t *testing.T) {
	type test struct {
		input       string
		expr        string
		expectMatch bool
		ignoreCase  bool
		name        string
	}

	var tests = []test{
		{"foo", "foo", true, false, "exact match"},
		{"foo", "bo", false, false, "invalid wildcard"},
		{"foo", ".*", false, false, "regex without regex delimiters"},
		{"foo", "/.*/", true, false, "regex with regex delimiters"},
		{"foo", "Foo", false, false, "wrong case"},
		{"foo", "Foo", true, true, "wrong case but with case insensitive matching"},
		{"foo", "", false, false, "empty expression should fail"},
		{"foo", "", false, true, "empty expression should fail also for case insensitive match"},
		{"foo", "/bar|.oo/", true, false, "or-expression"},
		{"Foo", "/bar|.O+/", true, true, "or-expression with case insensitive matching"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			match, err := NewMatcher(tc.expr, tc.ignoreCase)
			if err != nil {
				t.Fatal("failed to parse expression")
			}
			matched := match(tc.input)
			if !matched && tc.expectMatch == true {
				t.Error("failed to match when expected")
			} else if matched && tc.expectMatch == false {
				t.Error("matched when it shouldn't")
			}

		})
	}
}
