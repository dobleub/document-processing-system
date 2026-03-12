package helpers

import "testing"

func TestCleanFileName(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{input: "folder/sub/file.txt", expected: "file.txt"},
		{input: `folder\\sub\\file.txt`, expected: "file.txt"},
		{input: "file.txt", expected: "file.txt"},
	}

	for _, tc := range cases {
		got := CleanFileName(tc.input)
		if got != tc.expected {
			t.Fatalf("input %q expected %q, got %q", tc.input, tc.expected, got)
		}
	}
}

func TestStringArrayConversions(t *testing.T) {
	arr := []string{"ACTIVE", "DRAFT"}
	joined := StringArrayToString(arr)

	if joined != "[ACTIVE DRAFT]" {
		t.Fatalf("expected joined string '[ACTIVE DRAFT]', got %q", joined)
	}

	restored := StringToArray(joined)
	if len(restored) != 2 || restored[0] != "ACTIVE" || restored[1] != "DRAFT" {
		t.Fatalf("unexpected restored array: %#v", restored)
	}
}

func TestCleanQueryString(t *testing.T) {
	input := "  SELECT\t*\nFROM   process   WHERE   status = 'DONE'  "
	got := CleanQueryString(input)
	expected := "SELECT *FROM process WHERE status = 'DONE'"

	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}
