package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// b2Client is the global Backblaze B2 client
var b2Client *s3.Client
var b2Bucket string

// initStorage sets up the Backblaze B2 client using credentials from .env
func initStorage() {
	keyID := os.Getenv("B2_KEY_ID")
	appKey := os.Getenv("B2_APP_KEY")
	endpoint := os.Getenv("B2_ENDPOINT")
	b2Bucket = os.Getenv("B2_BUCKET")

	// If no B2 credentials, fall back to local storage
	if keyID == "" || appKey == "" || endpoint == "" || b2Bucket == "" {
		fmt.Println("Storage: no B2 credentials found, using local disk")
		return
	}

	// Create the S3 compatible client pointing at Backblaze B2
	b2Client = s3.New(s3.Options{
		BaseEndpoint: aws.String("https://" + endpoint),
		Region:       "auto",
		Credentials:  credentials.NewStaticCredentialsProvider(keyID, appKey, ""),
	})

	fmt.Println("Storage: connected to Backblaze B2")
}

// uploadFile uploads a file to B2 or saves it locally if B2 is not configured
func uploadFile(fileName string, file io.Reader, contentType string) error {
	if b2Client == nil {
		// Fall back to local disk
		return saveFileLocally(fileName, file)
	}

	_, err := b2Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(b2Bucket),
		Key:         aws.String(fileName),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	return err
}

// deleteFile removes a file from B2 or local disk
func deleteFile(fileName string) error {
	if fileName == "" {
		return nil
	}

	if b2Client == nil {
		return os.Remove(localPath(fileName))
	}

	_, err := b2Client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(b2Bucket),
		Key:    aws.String(fileName),
	})
	return err
}

// getFileURL returns a signed URL for private B2 files valid for 1 hour,
// or a local path if B2 is not configured
func getSignedURL(fileName string) (string, error) {
	if b2Client == nil {
		return "/stream-local/" + fileName, nil
	}

	presignClient := s3.NewPresignClient(b2Client)

	req, err := presignClient.PresignGetObject(context.Background(),
		&s3.GetObjectInput{
			Bucket: aws.String(b2Bucket),
			Key:    aws.String(fileName),
		},
		s3.WithPresignExpires(time.Hour),
	)
	if err != nil {
		return "", err
	}

	return req.URL, nil
}

// serveFile streams a file directly from B2 to the browser,
// or serves it from local disk if B2 is not configured
func serveFile(fileName string, w io.Writer) error {
	if b2Client == nil {
		f, err := os.Open(localPath(fileName))
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(w, f)
		return err
	}

	result, err := b2Client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(b2Bucket),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return err
	}
	defer result.Body.Close()

	_, err = io.Copy(w, result.Body)
	return err
}

// saveFileLocally saves a file to the local uploads directory
func saveFileLocally(fileName string, file io.Reader) error {
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "uploads"
	}

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return err
	}

	dst, err := os.Create(localPath(fileName))
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	return err
}

// localPath returns the full local path for a filename
func localPath(fileName string) string {
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "uploads"
	}
	return uploadDir + "/" + fileName
}

// initStorageFatal is like initStorage but exits if B2 credentials are missing
// in production mode
func mustInitStorage() {
	initStorage()

	if os.Getenv("APP_ENV") == "production" && b2Client == nil {
		log.Fatal("Storage: B2 credentials required in production mode")
	}
}