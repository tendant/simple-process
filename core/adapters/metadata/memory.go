package metadata

import (
	"context"
	"sync"

	"github.com/tendant/simple-process/core/contracts"
)

// MemoryMetadata stores file attributes and artifacts in memory for demos/tests.
type MemoryMetadata struct {
	mu         sync.Mutex
	attributes map[string]map[string]interface{}
	artifacts  map[string][]contracts.Artifact
}

// NewMemoryMetadata returns an initialized in-memory metadata service.
func NewMemoryMetadata() *MemoryMetadata {
	return &MemoryMetadata{
		attributes: make(map[string]map[string]interface{}),
		artifacts:  make(map[string][]contracts.Artifact),
	}
}

// UpdateFileAttributes merges the patch into the stored attribute set.
func (m *MemoryMetadata) UpdateFileAttributes(ctx context.Context, fileID string, attributesPatch map[string]interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	attrs, ok := m.attributes[fileID]
	if !ok {
		attrs = make(map[string]interface{})
		m.attributes[fileID] = attrs
	}
	for k, v := range attributesPatch {
		attrs[k] = v
	}
	return nil
}

// CreateArtifact appends the artifact to the file's artifact list.
func (m *MemoryMetadata) CreateArtifact(ctx context.Context, fileID string, artifact contracts.Artifact) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.artifacts[fileID] = append(m.artifacts[fileID], artifact)
	return nil
}

// Snapshot provides a copy of stored attributes and artifacts for inspection.
func (m *MemoryMetadata) Snapshot() (map[string]map[string]interface{}, map[string][]contracts.Artifact) {
	m.mu.Lock()
	defer m.mu.Unlock()

	attrsCopy := make(map[string]map[string]interface{}, len(m.attributes))
	for id, attrs := range m.attributes {
		clone := make(map[string]interface{}, len(attrs))
		for k, v := range attrs {
			clone[k] = v
		}
		attrsCopy[id] = clone
	}

	artifactsCopy := make(map[string][]contracts.Artifact, len(m.artifacts))
	for id, artifacts := range m.artifacts {
		clone := make([]contracts.Artifact, len(artifacts))
		copy(clone, artifacts)
		artifactsCopy[id] = clone
	}

	return attrsCopy, artifactsCopy
}
