//go:build s3

package s3

import (
	"context"
	"errors"
	"io"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/tendant/simple-process/core/adapters"
)

// Config captures the information required to construct an S3/MinIO storage adapter.
type Config struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	UseSSL          bool
	Region          string
	Bucket          string
	Prefix          string
	PresignExpiry   time.Duration
}

// Storage implements adapters.Storage using an S3-compatible backend.
type Storage struct {
	client *minio.Client
	bucket string
	prefix string
	expiry time.Duration
}

// New initialises the storage adapter using the provided configuration.
func New(cfg Config) (*Storage, error) {
	if cfg.Endpoint == "" {
		return nil, errors.New("endpoint is required")
	}
	if cfg.Bucket == "" {
		return nil, errors.New("bucket is required")
	}
	if cfg.PresignExpiry <= 0 {
		cfg.PresignExpiry = 15 * time.Minute
	}

	opts := &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, cfg.SessionToken),
		Secure: cfg.UseSSL,
		Region: cfg.Region,
	}

	client, err := minio.New(cfg.Endpoint, opts)
	if err != nil {
		return nil, err
	}

	return NewWithClient(client, cfg.Bucket, cfg.Prefix, cfg.PresignExpiry), nil
}

// NewWithClient wires an existing MinIO client into the adapter.
func NewWithClient(client *minio.Client, bucket, prefix string, presignExpiry time.Duration) *Storage {
	if presignExpiry <= 0 {
		presignExpiry = 15 * time.Minute
	}

	return &Storage{
		client: client,
		bucket: bucket,
		prefix: strings.Trim(prefix, "/"),
		expiry: presignExpiry,
	}
}

// Get returns a reader for the given blob location.
func (s *Storage) Get(ctx context.Context, location string) (io.ReadCloser, error) {
	key, bucket := s.resolve(location)
	object, err := s.client.GetObject(ctx, bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return object, nil
}

// Put uploads a blob from a reader to the given location.
func (s *Storage) Put(ctx context.Context, location string, reader io.Reader) error {
	key, bucket := s.resolve(location)
	_, err := s.client.PutObject(ctx, bucket, key, reader, -1, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	return err
}

// PresignGet generates a presigned URL for getting a blob.
func (s *Storage) PresignGet(ctx context.Context, location string) (string, error) {
	key, bucket := s.resolve(location)
	url, err := s.client.PresignedGetObject(ctx, bucket, key, s.expiry, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

func (s *Storage) resolve(location string) (key, bucket string) {
	bucket = s.bucket
	key = location

	if strings.HasPrefix(location, "s3://") || strings.HasPrefix(location, "minio://") {
		if u, err := url.Parse(location); err == nil {
			if u.Host != "" {
				bucket = u.Host
			}
			key = strings.TrimPrefix(u.Path, "/")
		}
	}

	if s.prefix != "" {
		key = path.Join(s.prefix, key)
	}

	key = strings.TrimLeft(key, "/")
	return key, bucket
}

var _ adapters.Storage = (*Storage)(nil)
