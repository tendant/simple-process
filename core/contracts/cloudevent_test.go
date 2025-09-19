package contracts

import (
	"encoding/json"
	"testing"
)

func TestNewJobCloudEventRoundTrip(t *testing.T) {
	job := Job{JobID: "job-123", UoW: "hash"}

	event, err := NewJobCloudEvent("example/source", job)
	if err != nil {
		t.Fatalf("NewJobCloudEvent returned error: %v", err)
	}

	if event.SpecVersion != "1.0" {
		t.Fatalf("unexpected spec version: %s", event.SpecVersion)
	}
	if event.Type != "simpleprocess.job" {
		t.Fatalf("unexpected type: %s", event.Type)
	}
	if event.Source != "example/source" {
		t.Fatalf("unexpected source: %s", event.Source)
	}
	if event.ID != job.JobID {
		t.Fatalf("expected event ID to match job ID; got %s", event.ID)
	}

	decoded, err := event.DecodeJob()
	if err != nil {
		t.Fatalf("DecodeJob returned error: %v", err)
	}
	if decoded.JobID != job.JobID || decoded.UoW != job.UoW {
		t.Fatalf("decoded job mismatch: %#v", decoded)
	}
}

func TestDecodeJobRejectsWrongContentType(t *testing.T) {
	payload, _ := json.Marshal(Job{JobID: "x"})
	event := CloudEvent{DataContentType: "application/xml", Data: payload}

	if _, err := event.DecodeJob(); err == nil {
		t.Fatalf("expected error for unsupported content type")
	}
}

func TestNewJobCloudEventRequiresJobID(t *testing.T) {
	_, err := NewJobCloudEvent("src", Job{})
	if err == nil {
		t.Fatalf("expected error when job ID missing")
	}
}
