package configurations

import "github.com/spf13/viper"

type MQConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

type SubscriberConfig struct {
	TranscodingQueue string
	ManagementQueue  string
}

type Buckets struct {
	RawVideos        string
	TranscodedVideos string
}

type AwsConfig struct {
	BaseUrl string
	Region  string
	Buckets Buckets
}

type Config struct {
	MQ   MQConfig
	Jobs SubscriberConfig
	Aws  AwsConfig
}

func LoadConfig() (Config, error) {
	conf := viper.New()

	conf.SetConfigName("config")
	conf.SetConfigType("yml")
	conf.AddConfigPath(".")

	if err := conf.ReadInConfig(); err != nil {
		return Config{}, err
	}

	var config Config
	err := conf.Unmarshal(&config)

	if err != nil {
		return Config{}, err
	}

	return config, nil
}