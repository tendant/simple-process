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

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/tendant/simple-process/core/adapters"
)

// Config captures the information required to construct an S3-compatible storage adapter.
type Config struct {
	Region         string
	Bucket         string
	Prefix         string
	Endpoint       string
	ForcePathStyle bool
	Credentials    aws.CredentialsProvider
	PresignExpiry  time.Duration
}

// Storage implements adapters.Storage using an AWS S3 compatible backend.
type Storage struct {
	client   *awss3.Client
	uploader *manager.Uploader
	presign  *awss3.PresignClient
	bucket   string
	prefix   string
	expiry   time.Duration
}

// New initialises the storage adapter using the provided configuration.
func New(ctx context.Context, cfg Config) (*Storage, error) {
	if cfg.Region == "" {
		return nil, errors.New("region is required")
	}
	if cfg.Bucket == "" {
		return nil, errors.New("bucket is required")
	}
	if cfg.PresignExpiry <= 0 {
		cfg.PresignExpiry = 15 * time.Minute
	}

	loadOpts := []func(*awsconfig.LoadOptions) error{awsconfig.WithRegion(cfg.Region)}
	if cfg.Credentials != nil {
		loadOpts = append(loadOpts, awsconfig.WithCredentialsProvider(cfg.Credentials))
	}
	if cfg.Endpoint != "" {
		endpoint := cfg.Endpoint
		loadOpts = append(loadOpts, awsconfig.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(func(service, region string, _ ...interface{}) (aws.Endpoint, error) {
			if service == awss3.ServiceID {
				return aws.Endpoint{URL: endpoint, SigningRegion: cfg.Region, HostnameImmutable: true}, nil
			}
			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		})))
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, loadOpts...)
	if err != nil {
		return nil, err
	}

	s3Opts := func(o *awss3.Options) {
		o.UsePathStyle = cfg.ForcePathStyle
	}

	client := awss3.NewFromConfig(awsCfg, s3Opts)
	return NewWithClient(client, cfg.Bucket, cfg.Prefix, cfg.PresignExpiry), nil
}

// NewWithClient wires an existing S3 client into the adapter.
func NewWithClient(client *awss3.Client, bucket, prefix string, presignExpiry time.Duration) *Storage {
	if presignExpiry <= 0 {
		presignExpiry = 15 * time.Minute
	}

	return &Storage{
		client:   client,
		uploader: manager.NewUploader(client),
		presign:  awss3.NewPresignClient(client),
		bucket:   bucket,
		prefix:   strings.Trim(prefix, "/"),
		expiry:   presignExpiry,
	}
}

// Get returns a reader for the given blob location.
func (s *Storage) Get(ctx context.Context, location string) (io.ReadCloser, error) {
	key, bucket := s.resolve(location)
	output, err := s.client.GetObject(ctx, &awss3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	return output.Body, nil
}

// Put uploads a blob from a reader to the given location.
func (s *Storage) Put(ctx context.Context, location string, reader io.Reader) error {
	key, bucket := s.resolve(location)
	_, err := s.uploader.Upload(ctx, &awss3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        reader,
		ContentType: aws.String("application/octet-stream"),
	})
	return err
}

// PresignGet generates a presigned URL for getting a blob.
func (s *Storage) PresignGet(ctx context.Context, location string) (string, error) {
	key, bucket := s.resolve(location)
	result, err := s.presign.PresignGetObject(ctx, &awss3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, awss3.WithPresignExpires(s.expiry))
	if err != nil {
		return "", err
	}
	return result.URL, nil
}

func (s *Storage) resolve(location string) (key, bucket string) {
	bucket = s.bucket
	key = location

	if strings.Contains(location, "://") {
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
