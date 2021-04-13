package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
)

// GetEnvWithKey : get env value
func GetEnvWithKey(key string) string {
	return os.Getenv(key)
}

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func main() {
	LoadEnv()

	awsRegion := GetEnvWithKey("AWS_REGION")
	awsAccessKeyID := GetEnvWithKey("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := GetEnvWithKey("AWS_SECRET_ACCESS_KEY")
	awsToken := GetEnvWithKey("AWS_TOKEN")
	myBucket := GetEnvWithKey("BUCKET_NAME")

	fmt.Println("My Region is: ", awsRegion)
	fmt.Println("My access Key is: ", awsAccessKeyID)
	fmt.Println("My Secret Access Key is: ", awsSecretAccessKey)
	fmt.Println("My Token is: ", awsToken)
	fmt.Println("My Bucket is :", myBucket)

	creds := credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, awsToken)

	_, err := creds.Get()
	if err != nil {
		fmt.Println("Error al obtener credenciales | ", err.Error())
	}

	// Set up Configuration and S3 Instance
	cfg := aws.NewConfig().WithRegion(awsRegion).WithCredentials(creds)

	// Create S3 instance
	svc := s3.New(session.New(), cfg)

	// Open the file
	file, err := os.Open("gopher_container.jpg")
	if err != nil {
		log.Println("Error al procesar imagen")
	}
	defer file.Close()

	// Size of the file
	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()

	// Create a buffer of the size file
	buffer := make([]byte, size)

	file.Read(buffer)

	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)

	// Make the set, the path and put together the s3 request
	path := "/uploads/" + file.Name()
	params := &s3.PutObjectInput{
		Bucket:        aws.String("test.solicitudes"),
		Key:           aws.String(path),
		Body:          fileBytes,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fileType),
		ACL:           aws.String("public-read"),
	}

	params2, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String("test.solicitudes"),
		Key:    aws.String(path), //aws.String(awsSecretAccessKey),
	})
	log.Println(params2.Body)

	/*
		urlStr, err := params2.Presign(15 * time.Minute)
		if err != nil {
			log.Println("Failed to sign request", err)
		}
		log.Println(urlStr.)
		log.Println(params2.HTTPRequest.URL.MarshalBinary())
	*/

	// Submit to AWS
	_, err = svc.PutObject(params)
	if err != nil {
		fmt.Println("Error al enviar a AWS")
	}
	filePath := "https://" + myBucket + "." + "s3-" + awsRegion + ".amazonaws.com" + path
	fmt.Printf("File URL: \n  %s", awsutil.StringValue(filePath))

}
