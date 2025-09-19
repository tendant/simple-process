package main

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	busadapter "github.com/tendant/simple-process/core/adapters/bus"
	metadataadapter "github.com/tendant/simple-process/core/adapters/metadata"
	"github.com/tendant/simple-process/core/adapters/storage"
	"github.com/tendant/simple-process/core/contracts"
	"github.com/tendant/simple-process/core/runner"
	"github.com/tendant/simple-process/core/uow"
	"github.com/tendant/simple-process/uows/go/hash"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	storageAdapter := storage.NewInMemoryStorage()
	bus := busadapter.NewMemoryBus(8)
	metadata := metadataadapter.NewMemoryMetadata()

	registry := map[string]uow.UoW{
		"hash": &hash.HashUoW{Storage: storageAdapter},
	}

	processed := make(chan struct{}, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		syncRunner := runner.NewSyncRunner()
		for {
			select {
			case <-ctx.Done():
				return
			case job, ok := <-bus.Subscribe():
				if !ok {
					return
				}

				handler, found := registry[job.UoW]
				if !found {
					fmt.Printf("no handler registered for %s\n", job.UoW)
					continue
				}

				result, err := syncRunner.Run(ctx, handler, job)
				if err != nil {
					fmt.Printf("worker error: %v\n", err)
					continue
				}
				if result == nil {
					continue
				}

				if err := metadata.UpdateFileAttributes(ctx, result.FileID, result.AttributesPatch); err != nil {
					fmt.Printf("metadata error: %v\n", err)
				}
				for _, artifact := range result.Artifacts {
					if err := metadata.CreateArtifact(ctx, result.FileID, artifact); err != nil {
						fmt.Printf("artifact error: %v\n", err)
					}
				}

				select {
				case processed <- struct{}{}:
				default:
				}
			}
		}
	}()

	if err := storageAdapter.Put(ctx, "async.txt", strings.NewReader("hello async world")); err != nil {
		panic(err)
	}

	job := contracts.Job{
		JobID: "async-job-1",
		UoW:   "hash",
		File: contracts.File{
			ID:   "async-file-1",
			Blob: contracts.Blob{Location: "async.txt"},
		},
		Return:  contracts.Return{Type: "metadata"},
		IdemKey: "async-file-1-hash",
	}

	asyncRunner := runner.NewAsyncRunner(bus)
	if _, err := asyncRunner.Run(ctx, registry[job.UoW], job); err != nil {
		panic(err)
	}

	select {
	case <-processed:
		fmt.Println("job processed asynchronously")
	case <-time.After(1 * time.Second):
		fmt.Println("processing timed out")
	}

	bus.Close()
	cancel()
	wg.Wait()

	attrs, artifacts := metadata.Snapshot()
	fmt.Printf("attributes stored: %+v\n", attrs[job.File.ID])
	fmt.Printf("artifacts stored: %+v\n", artifacts[job.File.ID])
}
