package geospatialjpv3

type Config struct {
	CkanBase                  string
	CkanOrg                   string
	CkanToken                 string
	CkanPrivate               bool
	CMSBase                   string
	CMSToken                  string
	CMSIntegration            string
	DisablePublication        bool
	DisableCatalogCheck       bool
	EnablePulicationOnWebhook bool
	PublicationToken          string
	JobName                   string
}
