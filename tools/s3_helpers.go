package main

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// UploadUsingS3 uploads a file to the specified S3-compatible storage.
func (app *Application) UploadUsingS3(fileContent *bytes.Reader, fileName string) (string, error) {

	// Specify the destination key in the bucket
	destinationKey := "uploads/" + fileName

	// Use the S3 client to upload the file
	_, err := app.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(app.BucketName),
		Key:    aws.String(destinationKey),
		Body:   fileContent,
	})
	if err != nil {
		return "", err
	}

	// Generate a permanent link
	permanentLink := fmt.Sprintf("%s/%s/%s", app.Endpoint, app.BucketName, destinationKey)

	return permanentLink, nil
}

// GetTemporaryLink generates a temporary signed URL for the given object in the bucket.
func (app *Application) GetTemporaryLink(bucketName, objectKey string, expiration time.Duration) (string, error) {
	// Generate a pre-signed URL for the object
	presignClient := s3.NewPresignClient(app.client)
	presignedReq, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, s3.WithPresignExpires(expiration))
	if err != nil {
		return "", err
	}

	return presignedReq.URL, nil
}

// ListFiles lists all files in the specified bucket and prefix.
func (app *Application) ListFiles(bucketName, prefix string) ([]string, error) {
	// List objects in the bucket
	result, err := app.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return nil, err
	}

	var files []string
	for _, item := range result.Contents {
		files = append(files, aws.ToString(item.Key))
	}

	return files, nil
}

// DeleteFile deletes a file from the specified bucket.
func (app *Application) DeleteFile(bucketName, objectKey string) error {
	// Delete the object
	_, err := app.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	return err
}
