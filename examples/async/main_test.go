package main

import (
	"testing"
)

func TestRunDemoProducesHash(t *testing.T) {
	attrs, artifacts, err := runDemo()
	if err != nil {
		t.Fatalf("runDemo returned error: %v", err)
	}

	if attrs == nil {
		t.Fatalf("expected attributes, got nil")
	}

	hashValue, ok := attrs["sha256"].(string)
	if !ok {
		t.Fatalf("expected sha256 attribute, got %#v", attrs["sha256"])
	}

	const expected = "23d3590d64af323ca8ddbfd54ee96263f8d7fd42fc0db36617cdccd5d1b1482e"
	if hashValue != expected {
		t.Fatalf("unexpected hash: %s", hashValue)
	}

	if len(artifacts) != 0 {
		t.Fatalf("expected no artifacts, got %d", len(artifacts))
	}
}
