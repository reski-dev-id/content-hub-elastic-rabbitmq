package config

import "github.com/spf13/viper"

type Config struct {
	AppPort     string
	DBUrl       string
	RabbitMQUrl string
	ElasticUrl  string
}

func Load() *Config {
	viper.SetConfigFile(".env")

	_ = viper.ReadInConfig()

	viper.AutomaticEnv()

	return &Config{
		AppPort:     viper.GetString("APP_PORT"),
		DBUrl:       viper.GetString("DB_URL"),
		RabbitMQUrl: viper.GetString("RABBITMQ_URL"),
		ElasticUrl:  viper.GetString("ELASTIC_URL"),
	}
}
