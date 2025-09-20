package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	busadapter "github.com/tendant/simple-process/pkg/adapters/bus"
	metadataadapter "github.com/tendant/simple-process/pkg/adapters/metadata"
	"github.com/tendant/simple-process/pkg/adapters/storage"
	"github.com/tendant/simple-process/pkg/contracts"
	"github.com/tendant/simple-process/pkg/runner"
	"github.com/tendant/simple-process/pkg/uow"
	"github.com/tendant/simple-process/uows/go/hash"
)

func main() {
	attrs, artifacts, err := runDemo()
	if err != nil {
		panic(err)
	}

	fmt.Printf("attributes stored: %+v\n", attrs)
	fmt.Printf("artifacts stored: %+v\n", artifacts)
}

func runDemo() (map[string]interface{}, []contracts.Artifact, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	storageAdapter := storage.NewInMemoryStorage()
	bus := busadapter.NewMemoryBus(8)
	metadata := metadataadapter.NewMemoryMetadata()

	registry := map[string]uow.UoW{
		"hash": &hash.HashUoW{Storage: storageAdapter},
	}

	done := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		syncRunner := runner.NewSyncRunner()

		for {
			select {
			case <-ctx.Done():
				done <- ctx.Err()
				return
			case job, ok := <-bus.Subscribe():
				if !ok {
					return
				}

				handler, found := registry[job.UoW]
				if !found {
					done <- fmt.Errorf("no handler registered for %s", job.UoW)
					return
				}

				result, err := syncRunner.Run(ctx, handler, job)
				if err != nil {
					done <- err
					return
				}
				if result == nil {
					done <- errors.New("runner returned nil result")
					return
				}

				if err := metadata.UpdateFileAttributes(ctx, result.FileID, result.AttributesPatch); err != nil {
					done <- err
					return
				}
				for _, artifact := range result.Artifacts {
					if err := metadata.CreateArtifact(ctx, result.FileID, artifact); err != nil {
						done <- err
						return
					}
				}

				done <- nil
				return
			}
		}
	}()

	if err := storageAdapter.Put(ctx, "async.txt", strings.NewReader("hello async world")); err != nil {
		bus.Close()
		wg.Wait()
		return nil, nil, err
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
		bus.Close()
		wg.Wait()
		return nil, nil, err
	}

	var resultErr error
	select {
	case resultErr = <-done:
	case <-ctx.Done():
		resultErr = ctx.Err()
	}

	bus.Close()
	wg.Wait()

	if resultErr != nil {
		return nil, nil, resultErr
	}

	attrsSnapshot, artifactsSnapshot := metadata.Snapshot()
	return attrsSnapshot[job.File.ID], artifactsSnapshot[job.File.ID], nil
}
