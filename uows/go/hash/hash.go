package hash

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"strings"

	"github.com/tendant/simple-process/pkg/adapters"
	"github.com/tendant/simple-process/pkg/contracts"
)

// HashUoW is a UoW that calculates the SHA256 hash of a file.
type HashUoW struct {
	Storage adapters.Storage
}

// Process executes the hash calculation.
func (u *HashUoW) Process(ctx context.Context, job contracts.Job) (*contracts.Result, error) {
	reader, err := u.Storage.Get(ctx, job.File.Blob.Location)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	h := sha256.New()
	if _, err := io.Copy(h, reader); err != nil {
		return nil, err
	}

	sha256sum := fmt.Sprintf("%x", h.Sum(nil))

	artifactLocation := fmt.Sprintf("artifacts/%s.sha256", job.File.ID)
	if err := u.Storage.Put(ctx, artifactLocation, strings.NewReader(sha256sum)); err != nil {
		return nil, err
	}

	return &contracts.Result{
		JobID:  job.JobID,
		UoW:    job.UoW,
		FileID: job.File.ID,
		AttributesPatch: map[string]interface{}{
			"sha256": sha256sum,
		},
		Artifacts: []contracts.Artifact{
			{
				Kind:     "checksum",
				MIME:     "text/plain",
				Bytes:    int64(len(sha256sum)),
				Location: artifactLocation,
			},
		},
	}, nil
}
