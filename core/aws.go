package core

import (
	"bytes"
	"mime/multipart"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func NewAWSSession(accessKeyID, secretAccessKey, region string) (*session.Session, error) {
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(region),
			Credentials: credentials.NewStaticCredentials(
				accessKeyID,
				secretAccessKey,
				"",
			),
		})
	return sess, err
}

func UploadImage(sess *session.Session, file multipart.File, header *multipart.FileHeader, bucket string, region string) (string, error) {
	uploader := s3manager.NewUploader(sess)

	filename := header.Filename
	//TODO: Handle up
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		ACL:    aws.String("public-read"),
		Key:    aws.String(filename),
		Body:   file,
	})
	if err != nil {
		return "", err
	}
	return "https://" + bucket + "." + "s3-" + region + ".amazonaws.com/" + filename, err
}

func PutImage(sess *session.Session, file multipart.File, handler *multipart.FileHeader, fileName string) error {

	size := handler.Size
	buffer := make([]byte, size)
	_, err := file.Read(buffer)
	if err != nil {
		return err
	}

	_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String("did2"),
		Key:                  aws.String(fileName),
		ACL:                  aws.String("public-read"),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(int64(size)),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
		StorageClass:         aws.String("INTELLIGENT_TIERING"),
	})

	return err
}
