package lib

import (
	"testing"
	"time"
)

func TestContains(t *testing.T) {
	values := []string{"a", "b", "c"}
	if !contains(values, "b") {
		t.Fatal("expected contains to find existing value")
	}
	if contains(values, "z") {
		t.Fatal("did not expect contains to find missing value")
	}
}

func TestFormatRemainingDuration(t *testing.T) {
	cases := []struct {
		input    time.Duration
		expected string
	}{
		{input: 0, expected: "0s"},
		{input: 5 * time.Second, expected: "5s"},
		{input: 65 * time.Second, expected: "1m 5s"},
		{input: 3661 * time.Second, expected: "1h 1m 1s"},
	}

	for _, tc := range cases {
		got := formatRemainingDuration(tc.input)
		if got != tc.expected {
			t.Fatalf("duration %v expected %s, got %s", tc.input, tc.expected, got)
		}
	}
}
