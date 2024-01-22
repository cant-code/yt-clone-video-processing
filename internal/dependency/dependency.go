package dependency

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-stomp/stomp/v3"
	"yt-clone-video-processing/internal/configurations"
)

type Dependency struct {
	Configs  configurations.Config
	S3Client *s3.Client
	MQConn   *stomp.Conn
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

	mq, err := ConnectToMQ(conf.MQ)
	if err != nil {
		return nil, err
	}

	return &Dependency{
		Configs:  conf,
		S3Client: client,
		MQConn:   mq,
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

func ConnectToMQ(config configurations.MQConfig) (*stomp.Conn, error) {
	dial, err := stomp.Dial("tcp",
		GenerateAddress(config),
		stomp.ConnOpt.Login(config.User, config.Password))

	if err != nil {
		return nil, err
	}

	return dial, nil
}

func GenerateAddress(conf configurations.MQConfig) string {
	return fmt.Sprintf("%s:%s", conf.Host, conf.Port)
}
