package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

func main() {
	s3Cl, err := NewS3Client()
	if err != nil {
		log.Fatal(err)
	}
	// dbconf.SetupDB()

	patchesPath := "../../../files/patches/"
	dirEntries, err := os.ReadDir(patchesPath)
	if err != nil {
		log.Fatal(err)
	}

	if len(dirEntries) == 0 {
		log.Fatal("dirEntities len is 0")
	}

	entery := dirEntries[0]
	fmt.Println(entery.Name())
	fb, err := os.ReadFile(patchesPath + entery.Name())
	if err != nil {
		log.Fatal(err)
	}
	_, err = s3Cl.UploadUsingS3(bytes.NewReader(fb), FilePathFromVersion("0.2.0"))
	if err != nil {
		log.Fatal(err)
	}

}

func FilePathFromVersion(v string) string {
	return fmt.Sprintf("patches/GameContentV_%s.pck", v)
}

type S3Config struct {
	AccessKey  string
	BucketName string
	SecretKey  string
	Endpoint   string
}

type s3Client struct {
	*s3.Client
	*S3Config
}

func NewS3Client() (*s3Client, error) {
	cl := &s3Client{}
	err := godotenv.Load("../../.env")
	if err != nil {

		return nil, fmt.Errorf("error loading .env file")
	}

	cl.S3Config = &S3Config{
		AccessKey:  os.Getenv("LIARA_ACCESS_KEY"),
		BucketName: os.Getenv("LIARA_BUCKET_NAME"),
		SecretKey:  os.Getenv("LIARA_SECRET_KEY"),
		Endpoint:   os.Getenv("LIARA_ENDPOINT"),
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
	if err != nil {
		return nil, err
	}

	// Initialize S3 client with custom configuration
	cfg.Credentials = aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
		return aws.Credentials{
			AccessKeyID:     cl.AccessKey,
			SecretAccessKey: cl.SecretKey,
		}, nil
	})
	cfg.BaseEndpoint = aws.String(cl.Endpoint)

	cl.Client = s3.NewFromConfig(cfg)
	return cl, nil
}

// // Upload a file
// fileContent := bytes.NewReader([]byte("Hello, World!"))
// fileName := "example.txt"

// permanentLink, err := app.UploadUsingS3(fileContent, fileName)
// if err != nil {
// 	log.Fatalf("Failed to upload file: %v", err)
// }
// fmt.Printf("File uploaded successfully. Permanent link: %s\n", permanentLink)

// // Generate a temporary link (valid for 1 hour)
// temporaryLink, err := app.GetTemporaryLink(app.BucketName, "uploads/example.txt", time.Hour)
// if err != nil {
// 	log.Fatalf("Failed to generate temporary link: %v", err)
// }
// fmt.Printf("Temporary link (valid for 1 hour): %s\n", temporaryLink)

// // List files in the bucket
// files, err := app.ListFiles(app.BucketName, "uploads/")
// if err != nil {
// 	log.Fatalf("Failed to list files: %v", err)
// }
// fmt.Println("Files in the bucket:")
// for _, file := range files {
// 	fmt.Println(file)
// }

// // Delete a file
// err = app.DeleteFile(app.BucketName, "uploads/example.txt")
// if err != nil {
// 	log.Fatalf("Failed to delete file: %v", err)
// }
// fmt.Println("File deleted successfully.")
