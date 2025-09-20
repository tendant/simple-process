//go:build nats && integration

package nats

import (
	"context"
	"testing"
	"time"

	natsclient "github.com/nats-io/nats.go"
	"github.com/tendant/simple-process/pkg/contracts"
)

func TestBusPublishAndSubscribe(t *testing.T) {
	conn, err := natsclient.Connect(natsclient.DefaultURL)
	if err != nil {
		t.Skipf("skipping: unable to connect to NATS (%v)", err)
	}
	t.Cleanup(func() { conn.Drain() })

	bus, err := NewBus(conn, "test.jobs", "simple-process/test")
	if err != nil {
		t.Fatalf("NewBus error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	recv := make(chan contracts.Job, 1)
	sub, err := SubscribeWorker(conn, bus.Subject(), "test-queue", func(_ context.Context, job contracts.Job) error {
		recv <- job
		return nil
	})
	if err != nil {
		t.Fatalf("SubscribeWorker error: %v", err)
	}
	t.Cleanup(func() { sub.Drain() })

	job := contracts.Job{JobID: "integration", UoW: "hash"}
	if err := bus.Publish(ctx, job); err != nil {
		t.Fatalf("Publish error: %v", err)
	}

	select {
	case got := <-recv:
		if got.JobID != job.JobID {
			t.Fatalf("unexpected job: %#v", got)
		}
	case <-ctx.Done():
		t.Fatalf("timed out waiting for job: %v", ctx.Err())
	}
}
