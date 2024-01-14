package dependency

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"yt-clone-video-processing/internal/configurations"
)

type Dependency struct {
	Configs  configurations.Config
	S3Client *s3.Client
}

func GetDependencies() (*Dependency, error) {
	conf, err := configurations.LoadConfig()
	if err != nil {
		return nil, err
	}

	client, err := GetClient(conf.Aws)
	if err != nil {
		return nil, err
	}

	return &Dependency{
		Configs:  conf,
		S3Client: client,
	}, nil
}

func GetClient(awsConfig configurations.AwsConfig) (*s3.Client, error) {
	defaultConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(defaultConfig, func(options *s3.Options) {
		options.BaseEndpoint = aws.String(awsConfig.BaseUrl)
		options.Region = awsConfig.Region
		options.UsePathStyle = true
	})

	return client, nil
}
