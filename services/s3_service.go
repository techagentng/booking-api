package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

// S3Service handles file uploads to AWS S3
type S3Service struct {
	client     *s3.Client
	bucketName string
	region     string
}

// NewS3Service creates a new S3 service instance
func NewS3Service() (*S3Service, error) {
	bucketName := os.Getenv("AWS_BUCKET")
	region := os.Getenv("S3_REGION")
	if region == "" {
		region = os.Getenv("AWS_REGION")
	}
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	if bucketName == "" || region == "" || accessKey == "" || secretKey == "" {
		return nil, fmt.Errorf("missing AWS configuration")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(cfg)

	return &S3Service{
		client:     client,
		bucketName: bucketName,
		region:     region,
	}, nil
}

// UploadFile uploads a file to S3 and returns the URL
func (s *S3Service) UploadFile(file multipart.File, header *multipart.FileHeader, folder string) (string, error) {
	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%s/%s_%s%s", folder, time.Now().Format("20060102"), uuid.New().String()[:8], ext)

	// Determine content type
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = getContentType(ext)
	}

	// Upload to S3
	_, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(filename),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	// Return the S3 URL
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucketName, s.region, filename)
	return url, nil
}

// DeleteFile deletes a file from S3
func (s *S3Service) DeleteFile(fileURL string) error {
	// Extract key from URL
	key := extractKeyFromURL(fileURL, s.bucketName)
	if key == "" {
		return fmt.Errorf("invalid file URL")
	}

	_, err := s.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	return nil
}

// ValidateFile validates file type and size
func ValidateFile(header *multipart.FileHeader, maxSizeMB int64, allowedTypes []string) error {
	// Check file size
	if header.Size > maxSizeMB*1024*1024 {
		return fmt.Errorf("file size exceeds %dMB limit", maxSizeMB)
	}

	// Check file type
	contentType := header.Header.Get("Content-Type")
	ext := strings.ToLower(filepath.Ext(header.Filename))

	isValidType := false
	for _, t := range allowedTypes {
		if contentType == t || ext == "."+strings.TrimPrefix(t, "image/") || ext == "."+strings.TrimPrefix(t, "application/") {
			isValidType = true
			break
		}
	}

	if !isValidType {
		return fmt.Errorf("invalid file type: %s", contentType)
	}

	return nil
}

// getContentType returns content type based on file extension
func getContentType(ext string) string {
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".webp":
		return "image/webp"
	case ".pdf":
		return "application/pdf"
	default:
		return "application/octet-stream"
	}
}

// extractKeyFromURL extracts the S3 key from a full URL
func extractKeyFromURL(url, bucket string) string {
	// Handle URL format: https://bucket.s3.region.amazonaws.com/key
	prefix := fmt.Sprintf("https://%s.s3.", bucket)
	if strings.HasPrefix(url, prefix) {
		parts := strings.SplitN(url, ".amazonaws.com/", 2)
		if len(parts) == 2 {
			return parts[1]
		}
	}
	return ""
}
