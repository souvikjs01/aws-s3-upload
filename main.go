package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var uploader *manager.Uploader
var (
	AWSRegion    string
	AWSAccessKey string
	AWSSecretKey string
	BucketName   string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("error loading env file")
	}
	AWSRegion = os.Getenv("AWS_REGION")
	AWSAccessKey = os.Getenv("AWS_ACCESS_KEY")
	AWSSecretKey = os.Getenv("AWS_SECRET_KEY")
	BucketName = os.Getenv("AWS_BUCKET_NAME")

	if AWSRegion == "" || AWSSecretKey == "" || AWSAccessKey == "" || BucketName == "" {
		fmt.Println("emty credentials")
	}

	// Load custom config with static credentials
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(AWSRegion),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(AWSAccessKey, AWSSecretKey, ""),
		),
	)
	if err != nil {
		panic(fmt.Sprintf("failed to load config, %v", err))
	}

	// Create S3 client
	client := s3.NewFromConfig(cfg)

	// Create uploader
	uploader = manager.NewUploader(client)

}

func main() {

	r := gin.Default()
	r.POST("/upload", uploadFile)
	r.Run(":8000")

}

func uploadFile(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var errors []string
	var uploadedURLs []string

	files := form.File["files"]

	for _, file := range files {
		f, err := file.Open()
		if err != nil {
			errors = append(errors, fmt.Sprintf("Error opening file %s: %s", file.Filename, err.Error()))
			continue
		}
		defer f.Close()

		uploadedURL, err := saveFile(f, file)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Error uploading file %s: %s", file.Filename, err.Error()))
		} else {
			uploadedURLs = append(uploadedURLs, uploadedURL)
		}
	}

	if len(errors) > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors})
	} else {
		c.JSON(http.StatusOK, gin.H{"url": uploadedURLs})
	}
}

func saveFile(fileReader io.Reader, fileHeader *multipart.FileHeader) (string, error) {
	_, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(BucketName),
		Key:    aws.String(fileHeader.Filename),
		Body:   fileReader,
	})

	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", BucketName, AWSRegion, fileHeader.Filename)

	return url, nil
}
