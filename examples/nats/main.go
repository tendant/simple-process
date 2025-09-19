//go:build nats

package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	natsclient "github.com/nats-io/nats.go"
	"github.com/tendant/simple-process/core/adapters/storage"
	"github.com/tendant/simple-process/core/contracts"
	"github.com/tendant/simple-process/core/runner"
	"github.com/tendant/simple-process/core/uow"
	natsbus "github.com/tendant/simple-process/transports/nats"
	"github.com/tendant/simple-process/uows/go/hash"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := natsclient.Connect(natsclient.DefaultURL)
	if err != nil {
		log.Fatalf("connect to nats: %v", err)
	}
	defer conn.Drain()

	subject := "simple-process.jobs"
	bus, err := natsbus.NewBus(conn, subject, "simple-process/examples/nats")
	if err != nil {
		log.Fatalf("new bus: %v", err)
	}

	storageAdapter := storage.NewInMemoryStorage()
	registry := map[string]uow.UoW{
		"hash": &hash.HashUoW{Storage: storageAdapter},
	}

	// Seed storage with test blob.
	if err := storageAdapter.Put(ctx, "nats-demo.txt", strings.NewReader("hello from nats")); err != nil {
		log.Fatalf("seed storage: %v", err)
	}

	// Worker subscription.
	_, err = natsbus.SubscribeWorker(conn, subject, "hash-workers", func(jobCtx context.Context, job contracts.Job) error {
		handler, ok := registry[job.UoW]
		if !ok {
			return fmt.Errorf("no handler for %s", job.UoW)
		}

		result, err := runner.NewSyncRunner().Run(jobCtx, handler, job)
		if err != nil {
			return err
		}

		fmt.Printf("job %s finished with attributes %v\n", job.JobID, result.AttributesPatch)
		return nil
	})
	if err != nil {
		log.Fatalf("subscribe: %v", err)
	}

	job := contracts.Job{
		JobID: "nats-job-1",
		UoW:   "hash",
		File: contracts.File{
			ID:   "nats-file-1",
			Blob: contracts.Blob{Location: "nats-demo.txt"},
		},
	}

	asyncRunner := runner.NewAsyncRunner(bus)
	if _, err := asyncRunner.Run(ctx, registry[job.UoW], job); err != nil {
		log.Fatalf("publish job: %v", err)
	}

	// Give the worker time to process before exiting.
	time.Sleep(500 * time.Millisecond)
	fmt.Println("published job to NATS; check worker output above")
}
