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
	expiresDelta int `yaml:"expires_delta"`
}

func (j JWTConfig) GetExpiresDelta() time.Duration {
	return time.Minute * time.Duration(j.expiresDelta)
}

type RedisConfig struct {
	getTimeLimit int `yaml:"get_time_limit"`
}

func (r RedisConfig) GetTimeLimit() time.Duration {
	return time.Millisecond * time.Duration(r.getTimeLimit)
}

var Cfg Config

func init() {
	file, err := os.Open("../config.yaml")
	if err != nil {
		log.Fatalf("Failed to find config file: %s", err.Error())
	}

	err = yaml.NewDecoder(file).Decode(&Cfg)
	if err != nil {
		log.Fatalf("Failed to parse config file: %s", err.Error())
	}
}
