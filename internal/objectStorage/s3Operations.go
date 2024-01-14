package objectStorage

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	"os"
	"yt-clone-video-processing/internal/dependency"
)

func GetObject(key string, dependency dependency.Dependency) (string, error) {
	object, err := dependency.S3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &dependency.Configs.Aws.Buckets.RawVideos,
		Key:    &key,
	})
	if err != nil {
		return "", err
	}

	temp, err := os.CreateTemp("./files", "testfile")
	if err != nil {
		return "", err
	}

	_, err = temp.ReadFrom(object.Body)
	if err != nil {
		return "", err
	}

	defer func(temp *os.File) {
		err := temp.Close()
		if err != nil {
			log.Panic(err)
		}
	}(temp)

	return temp.Name(), nil
}
