package geospatialjpv3

import "github.com/k0kubun/pp/v3"

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

var ppp *pp.PrettyPrinter

func init() {
	ppp = pp.New()
	ppp.SetColoringEnabled(false)
}
