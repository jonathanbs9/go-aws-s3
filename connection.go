package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

var (
	AccessKeyID     string
	SecretAccessKey string
	MyRegion        string
)

func ConnectAWS() *session.Session {
	AccessKeyID = GetEnvWithKey("AWS_ACCESS_KEY_ID")
	SecretAccessKey = GetEnvWithKey("AWS_SECRET_ACCESS_KEY")
	MyRegion = GetEnvWithKey("AWS_REGION")

	fmt.Println(AccessKeyID)
	fmt.Println(SecretAccessKey)
	fmt.Println(MyRegion)

	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(MyRegion),
			Credentials: credentials.NewStaticCredentials(
				AccessKeyID,
				SecretAccessKey,
				"", // A token will be created when the session it's used.
			),
		})
	if err != nil {
		log.Println(err.Error())
	}
	return sess
}
