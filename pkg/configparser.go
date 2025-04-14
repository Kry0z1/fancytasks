package tasks

import (
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	JWT JWTConfig `yaml:"jwt"`
}

type JWTConfig struct {
	ExpiresDelta int `yaml:"expires_delta"`
}

func (j JWTConfig) GetExpiresDelta() time.Duration {
	return time.Duration(j.ExpiresDelta) * time.Minute
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
