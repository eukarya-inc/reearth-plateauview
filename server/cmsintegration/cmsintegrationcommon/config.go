package cmsintegrationcommon

type Config struct {
	FMEMock             bool
	FMEBaseURL          string
	FMEToken            string
	FMEResultURL        string
	FMESkipQualityCheck bool
	CMSBaseURL          string
	CMSToken            string
	CMSIntegration      string
	Secret              string
	Debug               bool
}
