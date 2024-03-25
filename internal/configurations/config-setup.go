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

type DBConfig struct {
	Url      string
	Username string
	Password string
}

type Config struct {
	MQ   MQConfig
	Jobs SubscriberConfig
	Aws  AwsConfig
	DB   DBConfig
}

func LoadConfig() (Config, error) {
	conf := viper.New()

	conf.SetConfigName("config")
	conf.SetConfigType("yml")
	conf.AddConfigPath("./config")
	conf.AutomaticEnv()

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
