package tasks

import (
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	JWT   JWTConfig   `yaml:"jwt"`
	Redis RedisConfig `yaml:"redis"`
}

type JWTConfig struct {
	ExpiresDelta int `yaml:"expires_delta"`
}

func (j JWTConfig) GetExpiresDelta() time.Duration {
	return time.Duration(j.ExpiresDelta) * time.Minute
}

type RedisConfig struct {
	GetTimeLimit int `yaml:"get_time_limit"`
}

func (r RedisConfig) GetTimeLimitInMilliseconds() time.Duration {
	return time.Duration(r.GetTimeLimit) * time.Millisecond
}

var Cfg Config

func init() {
	file, err := os.Open("config.yaml")
	if err != nil {
		log.Fatalf("Failed to find config file: %s", err.Error())
	}

	err = yaml.NewDecoder(file).Decode(&Cfg)
	if err != nil {
		log.Fatalf("Failed to parse config file: %s", err.Error())
	}
}
