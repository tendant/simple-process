//go:build nats

package nats

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	natsclient "github.com/nats-io/nats.go"
	"github.com/tendant/simple-process/core/contracts"
)

// Bus publishes Jobs to NATS subjects so remote workers can execute them.
type Bus struct {
	conn    *natsclient.Conn
	subject string
	source  string
}

// NewBus wires an existing NATS connection into the adapters.Bus interface.
func NewBus(conn *natsclient.Conn, subject, source string) (*Bus, error) {
	if conn == nil {
		return nil, errors.New("nats connection is required")
	}
	if subject == "" {
		return nil, errors.New("subject is required")
	}
	if source == "" {
		source = "simple-process/nats"
	}
	return &Bus{conn: conn, subject: subject, source: source}, nil
}

// Publish serialises the job to JSON and pushes it onto the configured subject.
func (b *Bus) Publish(ctx context.Context, job contracts.Job) error {
	event, err := contracts.NewJobCloudEvent(b.source, job)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal cloudevent: %w", err)
	}

	// Honour context cancellation by using RequestWithContext semantics.
	// NATS does not natively accept contexts, so we rely on PublishMsgAsync.
	msg := &natsclient.Msg{Subject: b.subject, Data: payload}
	if err := b.conn.PublishMsg(msg); err != nil {
		return err
	}

	// Flush to ensure delivery before returning; respect deadlines if present.
	deadline, ok := ctx.Deadline()
	if ok {
		timeout := time.Until(deadline)
		if timeout <= 0 {
			timeout = time.Millisecond
		}
		return b.conn.FlushTimeout(timeout)
	}

	return b.conn.Flush()
}

// Subject exposes the NATS subject used for publishing jobs.
func (b *Bus) Subject() string {
	return b.subject
}

// WorkerHandler processes a job pulled from NATS.
type WorkerHandler func(context.Context, contracts.Job) error

// SubscribeWorker attaches a queue subscription that hands jobs to the provided handler.
func SubscribeWorker(conn *natsclient.Conn, subject, queue string, handler WorkerHandler) (*natsclient.Subscription, error) {
	if conn == nil {
		return nil, errors.New("nats connection is required")
	}
	if subject == "" {
		return nil, errors.New("subject is required")
	}
	if handler == nil {
		return nil, errors.New("handler is required")
	}
	if queue == "" {
		queue = "simple-process-workers"
	}

	return conn.QueueSubscribe(subject, queue, func(msg *natsclient.Msg) {
		var event contracts.CloudEvent
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			fmt.Printf("nats worker: failed to decode event: %v\n", err)
			return
		}

		job, err := event.DecodeJob()
		if err != nil {
			fmt.Printf("nats worker: failed to extract job: %v\n", err)
			return
		}

		if err := handler(context.Background(), job); err != nil {
			fmt.Printf("nats worker: handler error: %v\n", err)
		}
	})
}
