package contracts

import (
	"encoding/json"
	"fmt"
	"time"
)

// CloudEvent represents a minimal CloudEvents v1.0 envelope used by this project.
type CloudEvent struct {
	SpecVersion     string          `json:"specversion"`
	Type            string          `json:"type"`
	Source          string          `json:"source"`
	ID              string          `json:"id"`
	Time            time.Time       `json:"time"`
	DataContentType string          `json:"datacontenttype"`
	Data            json.RawMessage `json:"data"`
}

const (
	cloudEventSpecVersion = "1.0"
	jobEventType          = "simpleprocess.job"
	jobDataContentType    = "application/json"
)

// NewJobCloudEvent wraps a Job in a CloudEvent envelope.
func NewJobCloudEvent(source string, job Job) (CloudEvent, error) {
	if job.JobID == "" {
		return CloudEvent{}, fmt.Errorf("job id is required")
	}

	payload, err := json.Marshal(job)
	if err != nil {
		return CloudEvent{}, fmt.Errorf("marshal job: %w", err)
	}

	if source == "" {
		source = "simple-process"
	}

	return CloudEvent{
		SpecVersion:     cloudEventSpecVersion,
		Type:            jobEventType,
		Source:          source,
		ID:              job.JobID,
		Time:            time.Now().UTC(),
		DataContentType: jobDataContentType,
		Data:            payload,
	}, nil
}

// DecodeJob extracts a Job from the CloudEvent payload.
func (e CloudEvent) DecodeJob() (Job, error) {
	if e.DataContentType != jobDataContentType {
		return Job{}, fmt.Errorf("unexpected data content type: %s", e.DataContentType)
	}

	var job Job
	if err := json.Unmarshal(e.Data, &job); err != nil {
		return Job{}, fmt.Errorf("decode job: %w", err)
	}
	return job, nil
}
