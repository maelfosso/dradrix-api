package awss3

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"stockinos.com/api/requests"
)

type S3Client struct {
	S3Client        s3.Client
	S3PresignClient s3.PresignClient
	Bucket          string
	// HttpRequester
	requests.HttpRequester
}

func NewS3Client() *S3Client {
	awsConfig, err := config.LoadDefaultConfig(
		context.TODO(),
	)
	if err != nil {
		panic(err)
	}

	s3Client := s3.NewFromConfig(awsConfig)
	s3PresignClient := s3.NewPresignClient(s3Client)
	return &S3Client{
		S3Client:        *s3Client,
		S3PresignClient: *s3PresignClient,
		Bucket:          os.Getenv("AWS_S3_BUCKET"),
		HttpRequester:   requests.HttpRequester{},
	}
}

// GetObject makes a presigned request that can be used to get an object from a bucket.
// The presigned request is valid for the specified number of seconds.
// func (s3 *S3Client) GetObject(
// 	bucketName string, objectKey string, lifetimeSecs int64) (*v4.PresignedHTTPRequest, error) {
// 	request, err := s3.PresignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
// 		Bucket: aws.String(bucketName),
// 		Key:    aws.String(objectKey),
// 	}, func(opts *s3.PresignOptions) {
// 		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
// 	})
// 	if err != nil {
// 		log.Printf("Couldn't get a presigned request to get %v:%v. Here's why: %v\n",
// 			bucketName, objectKey, err)
// 	}
// 	return request, err
// }

// snippet-end:[gov2.s3.PresignGetObject]

// snippet-start:[gov2.s3.PresignPubObject]

// PutObject makes a presigned request that can be used to put an object in a bucket.
// The presigned request is valid for the specified number of seconds.
func (s3Client *S3Client) PutObject(
	objectKey string,
	lifetimeSecs int64,
) (*v4.PresignedHTTPRequest, error) {
	request, err := s3Client.S3PresignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s3Client.Bucket),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
	})
	if err != nil {
		log.Printf("Couldn't get a presigned request to put %v:%v. Here's why: %v\n",
			s3Client.Bucket, objectKey, err)
	}
	return request, err
}

// snippet-end:[gov2.s3.PresignPubObject]

// snippet-start:[gov2.s3.PresignDeleteObject]

// DeleteObject makes a presigned request that can be used to delete an object from a bucket.
// func (presigner Presigner) DeleteObject(bucketName string, objectKey string) (*v4.PresignedHTTPRequest, error) {
// 	request, err := presigner.PresignClient.PresignDeleteObject(context.TODO(), &s3.DeleteObjectInput{
// 		Bucket: aws.String(bucketName),
// 		Key:    aws.String(objectKey),
// 	})
// 	if err != nil {
// 		log.Printf("Couldn't get a presigned request to delete object %v. Here's why: %v\n", objectKey, err)
// 	}
// 	return request, err
// }

func (client *S3Client) UploadFile(uploadKey string, fileToUpload *os.File) error {
	defer fileToUpload.Close()
	fileToUpload.Seek(0, io.SeekStart)

	presignedPutRequest, err := client.PutObject(uploadKey, 60)
	if err != nil {
		log.Println("Error when PutObject: ", err)
		return fmt.Errorf("ERR_S3_UPLF_01")
	}
	log.Printf("Got a presigned %v request to URL:\n\t%v\n", presignedPutRequest.Method,
		presignedPutRequest.URL)

	log.Println("Using net/http to send the request")
	info, err := fileToUpload.Stat()
	if err != nil {
		log.Println("Error on [Stat()]: ", err)
		return fmt.Errorf("ERR_S3_UPLF_03")
	}
	// log.Println("File to upload stat", info.Size(), fileToUpload.)

	putResponse, err := client.Put(presignedPutRequest.URL, info.Size(), fileToUpload)
	if err != nil {
		log.Println("Error on [client.Put]: ", err)
		return fmt.Errorf("ERR_S3_UPLF_04")
	}
	log.Printf("%v object %v with presigned URL returned %v.", presignedPutRequest.Method,
		uploadKey, putResponse.StatusCode)
	// body, err := ioutil.ReadAll(putResponse.Body)
	// if err != nil {
	// 	log.Println("Error reading response body: ", err)
	// 	return err
	// }
	// log.Println("Put Response Body: ", string(body))

	log.Println(strings.Repeat("-", 88))

	return nil
}
