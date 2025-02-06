package backend

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"rip/lib/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var bucketName string = "code-inspector"

func InitializeMinIO() (*minio.Client, error) {
	endpoint, accessKey, secretKey, useSSL, err := config.FromEnvMinIO()
	if err != nil {
		return nil, fmt.Errorf("failed to get MinIO configuration from environment: %v", err)
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MinIO client: %v", err)
	}

	log.Println("MinIO client initialized successfully")

	found, err := client.BucketExists(context.Background(), bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check if bucket exists: %v", err)
	}

	if !found {
		err = client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %v", err)
		}
		log.Printf("Bucket '%s' created successfully\n", bucketName)
	} else {
		log.Printf("Bucket '%s' already exists\n", bucketName)
	}

	return client, nil
}

func (app *App) uploadImageToMinIO(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("[err] failed to open image file: %w", err)
	}
	defer src.Close()

	_, err = app.minioClient.PutObject(context.Background(), "code-inspector", file.Filename, src, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})
	if err != nil {
		return "", fmt.Errorf("[err] failed to upload image to MinIO: %w", err)
	}

	imageURL := fmt.Sprintf("%s/%s/%s", app.minioClient.EndpointURL(), "code-inspector", file.Filename)
	return imageURL, nil
}

func (app *App) deleteImageFromMinIO(imgLink string) error {
	if imgLink == "" {
		return nil
	}

	objectName := extractObjectNameFromURL(imgLink)
	err := app.minioClient.RemoveObject(context.Background(), "code-inspector", objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("[err] failed to delete image from MinIO: %w", err)
	}

	return nil
}
