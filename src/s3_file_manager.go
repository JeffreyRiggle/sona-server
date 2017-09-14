package main

import (
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// S3FileManager provides the ability to store attachments in AWS S3
type S3FileManager struct {
	Region string // The region to store objects in.
	Bucket string // The bucket to store objects in.
}

// Initialize will setup S3. This will make sure S3 can be connected to.
// This will also create the configured bucket if it does not already exist.
func (manager S3FileManager) Initialize() {
	svc := s3.New(CreateSession(manager.Region))

	_, err := svc.GetBucketLocation(&s3.GetBucketLocationInput{
		Bucket: aws.String(manager.Bucket),
	})

	if err == nil {
		logManager.LogPrintf("Bucket %v exists starting s3 file manager\n", manager.Bucket)
		return
	}

	logManager.LogPrintf("Got an error. Attempting to create bucket. %v", err)
	svc2 := s3.New(CreateSession(manager.Region))
	result, err2 := svc2.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(manager.Bucket),
	})

	if err2 != nil {
		if aerr, ok := err2.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeBucketAlreadyExists:
				logManager.LogPrintln(s3.ErrCodeBucketAlreadyExists, aerr.Error())
			case s3.ErrCodeBucketAlreadyOwnedByYou:
				logManager.LogPrintln(s3.ErrCodeBucketAlreadyOwnedByYou, aerr.Error())
			default:
				logManager.LogPrintf("Unable to create bucket. %v", aerr.Error())
				panic(err2)
			}
		} else {
			logManager.LogPrintf("Unable to create bucket %v. Error %v\n", manager.Bucket, err2)
			panic(err2)
		}
	}

	logManager.LogPrintf("Created bucket %v, with result %v\n", manager.Bucket, result)
}

// SaveFile will attempt to save an attachment to the configured s3 bucket.
// If the attempt fails a false will be returned.
func (manager S3FileManager) SaveFile(incident string, fileName string, file multipart.File) (string, bool) {
	uploader := s3manager.NewUploader(CreateSession(manager.Region))

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(manager.Bucket),
		Key:    aws.String(incident + "/" + fileName),
		Body:   file,
	})

	if err != nil {
		return "", false
	}

	return result.Location, true
}

// LoadFile will attempt to load an attachment out of an s3 bucket.
// If the file cannot be returned a false will be returned.
func (manager S3FileManager) LoadFile(incident string, fileName string) (io.ReadSeeker, os.FileInfo, bool, func()) {
	downloader := s3manager.NewDownloader(CreateSession(manager.Region))

	f, err := ioutil.TempFile("", fileName)

	if err != nil {
		return nil, nil, false, nil
	}

	_, err2 := downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(manager.Bucket),
		Key:    aws.String(incident + "/" + fileName),
	})

	if err2 != nil {
		f.Close()
		defer os.Remove(f.Name())
		return nil, nil, false, nil
	}

	d, err3 := f.Stat()
	if err3 != nil {
		f.Close()
		defer os.Remove(f.Name())
		return nil, nil, false, nil
	}

	callback := func() {
		if f != nil {
			f.Close()
			os.Remove(f.Name())
		}
	}

	return f, d, true, callback
}

// CreateSession will create a session.Session for an AWS region.
func CreateSession(region string) *session.Session {
	return session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))
}
