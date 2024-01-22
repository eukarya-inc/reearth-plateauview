package main

import (
	"fmt"
	"os"

	"github.com/eukarya-inc/reearth-plateauview/server/cmsintegration"
	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog"
	"github.com/eukarya-inc/reearth-plateauview/server/opinion"
	"github.com/eukarya-inc/reearth-plateauview/server/plateaucms"
	"github.com/eukarya-inc/reearth-plateauview/server/sdkapi"
	"github.com/eukarya-inc/reearth-plateauview/server/searchindex"
	"github.com/eukarya-inc/reearth-plateauview/server/sidebar"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/reearth/reearthx/log"
)

const configPrefix = "REEARTH_PLATEAUVIEW"

type Config struct {
	Port                              uint   `default:"8080" envconfig:"PORT"`
	Host                              string `default:"http://localhost:8080"`
	Debug                             bool
	Origin                            []string
	Secret                            string
	Delegate_URL                      string
	CMS_Webhook_Secret                string
	CMS_BaseURL                       string
	CMS_Token                         string
	CMS_IntegrationID                 string
	CMS_PlateauProject                string
	CMS_SystemProject                 string
	CMS_TokenProject                  string
	FME_BaseURL                       string
	FME_BaseURL_V2                    string
	FME_URL_V3                        string
	FME_Mock                          bool
	FME_Token                         string
	FME_SkipQualityCheck              bool
	Ckan_BaseURL                      string
	Ckan_Org                          string
	Ckan_Token                        string
	Ckan_Private                      bool
	SDK_Token                         string
	SendGrid_APIKey                   string
	Opinion_From                      string
	Opinion_FromName                  string
	Opinion_To                        string
	Opinion_ToName                    string
	Sidebar_Token                     string
	Share_Disable                     bool
	Geospatialjp_Publication_Disable  bool
	Geospatialjp_CatalocCheck_Disable bool
	Geospatialjp_JobName              string
	DataConv_Disable                  bool
	Indexer_Delegate                  bool
	DataCatalog_DisableCache          bool
	DataCatalog_CacheUpdateKey        string
	DataCatalog_PlaygroundEndpoint    string
	DataCatalog_CacheTTL              int
	DataCatalog_GQL_MaxComplexity     int
	SDKAPI_DisableCache               bool
	SDKAPI_CacheTTL                   int
	GCParcent                         int
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
		Host:                            c.Host,
		FMEMock:                         c.FME_Mock,
		FMEBaseURL:                      c.FME_BaseURL,
		FMEToken:                        c.FME_Token,
		FMEBaseURLV2:                    c.FME_BaseURL_V2,
		FMEURLV3:                        c.FME_URL_V3,
		FMESkipQualityCheck:             c.FME_SkipQualityCheck,
		CMSBaseURL:                      c.CMS_BaseURL,
		CMSToken:                        c.CMS_Token,
		CMSIntegration:                  c.CMS_IntegrationID,
		Secret:                          c.Secret,
		Debug:                           c.Debug,
		CkanBaseURL:                     c.Ckan_BaseURL,
		CkanOrg:                         c.Ckan_Org,
		CkanToken:                       c.Ckan_Token,
		CkanPrivate:                     c.Ckan_Private,
		DisableGeospatialjpPublication:  c.Geospatialjp_Publication_Disable,
		DisableGeospatialjpCatalogCheck: c.Geospatialjp_CatalocCheck_Disable,
		DisableDataConv:                 c.DataConv_Disable,
		APIToken:                        c.Sidebar_Token,
		GeospatialjpJobName:             c.Geospatialjp_JobName,
	}
}

func (c *Config) SearchIndex() searchindex.Config {
	return searchindex.Config{
		CMSBase:           c.CMS_BaseURL,
		CMSToken:          c.CMS_Token,
		CMSStorageProject: c.CMS_SystemProject,
		Delegate:          c.Indexer_Delegate,
		DelegateURL:       c.Delegate_URL,
		Debug:             c.Debug,
		// CMSModel: c.CMS_Model,
		// CMSStorageModel:   c.CMS_IndexerStorageModel,
	}
}

func (c *Config) SDKAPI() sdkapi.Config {
	return sdkapi.Config{
		CMSBaseURL: c.CMS_BaseURL,
		CMSToken:   c.CMS_Token,
		Project:    c.CMS_PlateauProject,
		// Model:      c.CMS_SDKModel,
		Token:        c.SDK_Token,
		DisableCache: c.SDKAPI_DisableCache,
		CacheTTL:     c.SDKAPI_CacheTTL,
	}
}

func (c *Config) Opinion() opinion.Config {
	return opinion.Config{
		SendGridAPIKey: c.SendGrid_APIKey,
		From:           c.Opinion_From,
		FromName:       c.Opinion_FromName,
		To:             c.Opinion_To,
		ToName:         c.Opinion_ToName,
	}
}

func (c *Config) Sidebar() sidebar.Config {
	return sidebar.Config{
		Config:       c.plateauCMS(),
		DisableShare: c.Share_Disable,
	}
}

func (c *Config) DataCatalog() datacatalog.Config {
	return datacatalog.Config{
		Config:               c.plateauCMS(),
		CacheUpdateKey:       c.DataCatalog_CacheUpdateKey,
		PlaygroundEndpoint:   c.DataCatalog_PlaygroundEndpoint,
		GraphqlMaxComplexity: c.DataCatalog_GQL_MaxComplexity,
		DisableCache:         c.DataCatalog_DisableCache,
		CacheTTL:             c.DataCatalog_CacheTTL,
	}
}

func (c *Config) plateauCMS() plateaucms.Config {
	return plateaucms.Config{
		CMSBaseURL:      c.CMS_BaseURL,
		CMSMainToken:    c.CMS_Token,
		CMSTokenProject: c.CMS_TokenProject,
		// compat
		CMSMainProject: c.CMS_SystemProject,
		AdminToken:     c.Sidebar_Token,
	}
}
