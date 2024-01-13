package config

import "github.com/spf13/viper"

type MQConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

type SubscriberConfig struct {
	Queue string
}

type Config struct {
	MQ   MQConfig
	Jobs SubscriberConfig
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
