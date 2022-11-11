package main

import (
	"fmt"
	"os"

	"github.com/eukarya-inc/reearth-plateauview/server/cmsintegration"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/reearth/reearthx/log"
)

const configPrefix = "REEARTH_PLATEAUVIEW_"

type Config struct {
	Port               uint   `default:"8080" envconfig:"PORT"`
	Host               string `default:"http://localhost:8080"`
	CMS_Webhook_Secret string
	FME_Token          string
	CMS_BaseURL        string
	CMS_Token          string
	Secret             string
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

func (c *Config) CMSIntegration() cmsintegration.Config {
	return cmsintegration.Config{
		FMEBaseURL:       c.FME_Token,
		FMEToken:         c.FME_Token,
		FMEResultURL:     c.Host,
		CMSBaseURL:       c.CMS_BaseURL,
		CMSToken:         c.CMS_Token,
		CMSWebhookSecret: c.CMS_Webhook_Secret,
		Secret:           c.Secret,
	}
}
