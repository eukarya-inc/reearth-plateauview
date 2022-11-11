package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/reearth/reearthx/log"
)

const configPrefix = "REEARTH_PLATEAUVIEW_"

type Config struct {
	Port           uint   `default:"8080" envconfig:"PORT"`
	Host           string `default:"http://localhost:8080"`
	Webhook_Secret string
	FME_Secret     string
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil && !os.IsNotExist(err) {
		return nil, err
	} else if err == nil {
		log.Infof("config: .env loaded")
	}

	var c Config
	err := envconfig.Process(configPrefix, &c)

	return &c, err
}

func (c *Config) Print() string {
	s := fmt.Sprintf("%+v", c)
	return s
}
