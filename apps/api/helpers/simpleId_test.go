package helpers

import (
	"regexp"
	"testing"
)

func TestSimpleIDGenerateFormat(t *testing.T) {
	id := (SimpleID("")).Generate().String()

	if len(id) != 23 {
		t.Fatalf("expected generated id length 23, got %d (%s)", len(id), id)
	}

	pattern := regexp.MustCompile(`^[A-Z0-9]{5}(-[A-Z0-9]{5}){3}$`)
	if !pattern.MatchString(id) {
		t.Fatalf("generated id does not match expected format: %s", id)
	}
}

func TestSimpleIDFromStringAndString(t *testing.T) {
	raw := "ABCDE-12345-FGHIJ-67890"
	id := (SimpleID("")).FromString(raw)

	if id.String() != raw {
		t.Fatalf("expected %s, got %s", raw, id.String())
	}

	if !id.IsValid() {
		t.Fatal("expected IsValid to return true")
	}
}
