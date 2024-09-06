package config

import "github.com/spf13/viper"

type Config struct {
	Logger Logger
	App    App
	API    API
}

type Logger struct {
	Level string
}

type App struct {
	Database Database
}

type Database struct {
	Host           string
	Port           int
	Name           string
	User           string
	Password       string
	MigrationsPath string
}

type API struct {
	HTTP
	GRPC
	WebSocket
}

type HTTP struct {
	Host string
	Port string
}

type GRPC struct {
	Port int
}

type WebSocket struct {
	ReadBufferSize  int
	WriteBufferSize int
	Port            int
}

func Parse() (*Config, error) {
	v := viper.New()

	v.SetConfigType("yml")

	v.SetConfigFile("config.yml")

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config

	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
