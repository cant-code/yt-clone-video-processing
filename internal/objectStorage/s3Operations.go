package objectStorage

import (
	"context"
	"fmt"
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

func PutObject(key string, dependency dependency.Dependency) (string, error) {
	file, err := os.Open(fmt.Sprintf("./files/%s", key))
	if err != nil {
		return "", err
	}

	_, err = dependency.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &dependency.Configs.Aws.Buckets.TranscodedVideos,
		Key:    &key,
		Body:   file,
	})
	if err != nil {
		return "", err
	}

	defer func(file *os.File) {
		var fileName = file.Name()
		err := file.Close()
		if err != nil {
			log.Panicln(err)
		}

		err = os.Remove(fileName)
		if err != nil {
			log.Panicln(err)
		}
	}(file)

	return key, nil
}
