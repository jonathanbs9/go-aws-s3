package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
)

//var Session *session.Session = ConnectAWS()

//GetEnvWithKey : get env value
func GetEnvWithKey(key string) string {
	return os.Getenv(key)
}

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

/*func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("sess", Session)
		return next(c)
	}
}*/

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
	file, err := os.Open("gopher2_test.jpg")
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
	}

	params2, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String("test.solicitudes"),
		Key:    aws.String(path), //aws.String(awsSecretAccessKey),
	})

	urlStr, err := params2.Presign(15 * time.Minute)
	if err != nil {
		log.Println("Failed to sign request", err)
	}
	log.Println(urlStr)

	// Submit to AWS
	_, err = svc.PutObject(params)
	if err != nil {
		fmt.Println("Error al enviar a AWS")
	}
	filePath := "https://" + myBucket + "." + "s3-" + awsRegion + ".amazonaws.com" + path
	fmt.Printf("response \n  %s", awsutil.StringValue(filePath))

	//session := ConnectAWS()

	//e := echo.New()

	//e.Use(ServerHeader)

	/*e.POST("/uploadCert", func(c echo.Context, session *session.Session) error {
		sess, _ := session.Get("session", c)

		sess.Options = &sessions.Options{
			AwsRegion:    GetEnvWithKey("AWS_REGION"),
			AccesKey:     GetEnvWithKey("AWS_ACCESS_KEY_ID"),
			SecretAccess: GetEnvWithKey("AWS_SECRET_ACCESS_KEY"),
		}
		sess.Save(c.Request(), c.Response())
		return c.NoContent(http.StatusOK)
	})*/

	//e.Logger.Fatal(e.Start(":5000"))

}
