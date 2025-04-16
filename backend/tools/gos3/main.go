package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/arian-nj/master-card/back/internal/dbconf"
	"github.com/arian-nj/master-card/back/internal/version"
	"github.com/arian-nj/master-card/back/sqldb"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	s3Cl, err := NewS3Client()
	if err != nil {
		return err
	}

	// release_mode := os.Getenv("RELEASE_MODE")
	// if release_mode == "" {
	// 	log.Fatal("RELEASE_MODE is empty")
	// }
	release_mode := "release"

	queries, poll, err := dbconf.SetupDB()
	if err != nil {
		return err
	}
	defer poll.Close()

	pvRow, err := queries.GetVersion(context.Background(), release_mode)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	if pvRow.ID == 0 { // if no verson exist
		v, err := version.NewVersion("0.2.0")
		if err != nil {
			return err
		}
		pvRow, err = queries.InsertVersion(context.Background(), sqldb.InsertVersionParams{
			Rmode:         release_mode,
			VersionNumber: v.String(),
		})
		if err != nil {
			return err
		}
	} else { // if a version already exist
		v, err := version.NewVersion(pvRow.VersionNumber)
		if err != nil {
			return err
		}
		v.Patch++
		pvRow, err = queries.UpdateVersion(context.Background(), sqldb.UpdateVersionParams{
			Rmode:         release_mode,
			VersionNumber: v.String(),
		})
		if err != nil {
			return err
		}
	}

	patchesPath := "../../../files/patches/"
	dirEntries, err := os.ReadDir(patchesPath)
	if err != nil {
		return err
	}

	if len(dirEntries) == 0 {
		return fmt.Errorf("dirEntities len is 0")
	}

	entery := dirEntries[0]
	log.Println(entery.Name())
	fb, err := os.ReadFile(patchesPath + entery.Name())
	if err != nil {
		return err
	}

	_, err = s3Cl.UploadUsingS3(bytes.NewReader(fb), FilePathFromVersion(pvRow.VersionNumber, version.ReleasModes(release_mode)))
	if err != nil {
		return err
	}

	result, err := s3Cl.ListFiles(s3Cl.BucketName, FolderPathFromVersion(version.ReleasModes(release_mode)))
	if err != nil {
		return err
	}

	time_now := time.Now().UTC()
	for _, item := range result {
		if time_now.Sub(item.LastModified.UTC()).Hours() > 1 {
			err := s3Cl.DeleteFile(s3Cl.BucketName, *item.Key)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func FolderPathFromVersion(mode version.ReleasModes) string {
	return fmt.Sprintf("patches/%s/", mode)
}

func FilePathFromVersion(v string, mode version.ReleasModes) string {
	log.Println(v)
	return fmt.Sprintf("%sGameContentV_%s.pck", FolderPathFromVersion(mode), v)
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
