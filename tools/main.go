package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

type Application struct {
	client *s3.Client
	*S3Config
}

type S3Config struct {
	AccessKey  string
	BucketName string
	SecretKey  string
	Endpoint   string
}

func main() {
	app := Application{}
	app.SetCredentials()

	fileNamePtr := flag.String("file", "", "file name to upload to s3")
	flag.Parse()
	fmt.Printf("file name is %q\n", *fileNamePtr)

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

func (app *Application) SetCredentials() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app.S3Config = &S3Config{
		AccessKey:  os.Getenv("LIARA_ACCESS_KEY"),
		BucketName: os.Getenv("LIARA_BUCKET_NAME"),
		SecretKey:  os.Getenv("LIARA_SECRET_KEY"),
		Endpoint:   os.Getenv("LIARA_ENDPOINT"),
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
	if err != nil {
		panic(err)
	}

	// Initialize S3 client with custom configuration
	cfg.Credentials = aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
		return aws.Credentials{
			AccessKeyID:     app.AccessKey,
			SecretAccessKey: app.SecretKey,
		}, nil
	})
	cfg.BaseEndpoint = aws.String(app.Endpoint)

	app.client = s3.NewFromConfig(cfg)
}
