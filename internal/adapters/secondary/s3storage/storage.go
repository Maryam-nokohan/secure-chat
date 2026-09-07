package s3storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
	"github.com/maryam-nokohan/secure-chat/pkg"
)

type Config struct {
	Region          string
	Bucket          string
	AccessKeyID     string
	SecretAccessKey string
	Endpoint        string
	PublicEndpoint  string
	UsePathStyle    bool
}

var (
	safeSegment  = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	safeFilename = regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)
	safePath     = regexp.MustCompile(`^[a-zA-Z0-9_./-]+$`)
)


type Storage struct {
	client  *s3.Client
	presign *s3.PresignClient
	bucket  string
}

func New(ctx context.Context, cfg Config) (ports.FileStorage, error) {
	if cfg.Bucket == "" || cfg.Region == "" {
		return nil, fmt.Errorf("s3 storage misconfigured: bucket and region are required")
	}

	var optFns []func(*awsconfig.LoadOptions) error
	optFns = append(optFns, awsconfig.WithRegion(cfg.Region))
	if cfg.AccessKeyID != "" && cfg.SecretAccessKey != "" {
		optFns = append(optFns, awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		))
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, optFns...)
	if err != nil {
		return nil, fmt.Errorf("loading AWS config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
		}
		o.UsePathStyle = cfg.UsePathStyle
	})

	presignEndpoint := cfg.Endpoint
	if cfg.PublicEndpoint != "" {
		presignEndpoint = cfg.PublicEndpoint
	}
	presignClient := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if presignEndpoint != "" {
			o.BaseEndpoint = aws.String(presignEndpoint)
		}
		o.UsePathStyle = cfg.UsePathStyle
	})

	pkg.LogInfo("Initializing S3 storage (bucket=" + cfg.Bucket + ", region=" + cfg.Region + ")")

	return &Storage{client: client, presign: s3.NewPresignClient(presignClient), bucket: cfg.Bucket}, nil
}

func (s *Storage) key(folder, filename string) (string, error) {
	if folder == "" || !safeSegment.MatchString(folder) {
		return "", fmt.Errorf("invalid storage folder")
	}
	if filename == "" || !safeFilename.MatchString(filename) || strings.Contains(filename, "..") {
		return "", fmt.Errorf("invalid filename")
	}
	return folder + "/" + filename, nil
}

func (s *Storage) validatePath(path string) error {
	if path == "" || !safePath.MatchString(path) || strings.Contains(path, "..") {
		return fmt.Errorf("invalid storage path")
	}
	return nil
}

func (s *Storage) Upload(ctx context.Context, folder, filename string, content []byte, contentType string) (string, error) {
	key, err := s.key(folder, filename)
	if err != nil {
		return "", err
	}
	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(s.bucket),
		Key:           aws.String(key),
		Body:          bytes.NewReader(content),
		ContentLength: aws.Int64(int64(len(content))),
		ContentType:   aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("s3 put object: %w", err)
	}
	return key, nil
}

func (s *Storage) Download(ctx context.Context, path string) ([]byte, string, error) {
	if err := s.validatePath(path); err != nil {
		return nil, "", err
	}
	out, err := s.client.GetObject(ctx, &s3.GetObjectInput{Bucket: aws.String(s.bucket), Key: aws.String(path)})
	if err != nil {
		return nil, "", fmt.Errorf("s3 get object: %w", err)
	}
	defer out.Body.Close()
	data, err := io.ReadAll(out.Body)
	if err != nil {
		return nil, "", fmt.Errorf("reading s3 object body: %w", err)
	}
	ct := "application/octet-stream"
	if out.ContentType != nil && *out.ContentType != "" {
		ct = *out.ContentType
	}
	return data, ct, nil
}

func (s *Storage) Delete(ctx context.Context, path string) error {
	if err := s.validatePath(path); err != nil {
		return err
	}
	if _, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{Bucket: aws.String(s.bucket), Key: aws.String(path)}); err != nil {
		return fmt.Errorf("s3 delete object: %w", err)
	}
	return nil
}

func (s *Storage) PresignedURL(ctx context.Context, path string, expiry time.Duration) (string, error) {
	if err := s.validatePath(path); err != nil {
		return "", err
	}
	req, err := s.presign.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket), Key: aws.String(path),
	}, s3.WithPresignExpires(expiry))
	if err != nil {
		return "", fmt.Errorf("presign get object: %w", err)
	}
	return req.URL, nil
}
