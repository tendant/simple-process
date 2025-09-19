package main

import (
	"context"
	"fmt"

	"github.com/tendant/simple-process/core/adapters/storage"
	"github.com/tendant/simple-process/core/contracts"
	"github.com/tendant/simple-process/core/runner"
	"github.com/tendant/simple-process/uows/go/hash"
)

func main() {
	// Create a new in-memory storage adapter for the example.
	storage := storage.NewInMemoryStorage()

	// Create a new synchronous runner.
	runner := runner.NewSyncRunner()

	// Create a new hash UoW.
	hashUoW := &hash.HashUoW{Storage: storage}

	// Create a sample job.
	job := contracts.Job{
		JobID: "example-job-1",
		UoW:   "hash",
		File: contracts.File{
			ID:       "example-file-1",
			Blob:     contracts.Blob{Location: "example.txt"},
		},
	}

	// Upload a sample file to the in-memory storage.
	storage.Put(context.Background(), "example.txt", "hello world")

	// Run the UoW.
	result, err := runner.Run(context.Background(), hashUoW, job)
	if err != nil {
		panic(err)
	}

	// Print the result.
	fmt.Printf("Result: %+v\n", result)
}

