package cmsintegrationcommon

type Config struct {
	// general
	Secret string
	Debug  bool
	// CMS
	CMSBaseURL     string
	CMSToken       string
	CMSIntegration string
	// FME common
	FMEMock      bool
	FMEResultURL string
	// FME v3
	FMEURLV3 string
	// FME v2
	FMEBaseURL          string
	FMEBaseURLV2        string
	FMEToken            string
	FMESkipQualityCheck bool
}
