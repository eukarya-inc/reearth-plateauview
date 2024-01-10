package main

import (
	"flag"
	"os"

	"github.com/eukarya-inc/reearth-plateauview/worker/preparegspatialjp"
	"github.com/samber/lo"
)

func main() {
	config := lo.Must(NewConfig())

	switch os.Args[1] {
	case "prepare-gspatialjp":
		prepareGspatialjp(config)
	}
}

func prepareGspatialjp(conf *Config) {
	config := preparegspatialjp.Config{
		CMSURL:   conf.CMS_URL,
		CMSToken: conf.CMS_Token,
	}

	flag := flag.NewFlagSet("prepare-gspatialjp", flag.ExitOnError)
	flag.StringVar(&config.ProjectID, "project", "", "CMS project id")
	flag.StringVar(&config.CityItemID, "city", "", "CMS city item id")
	flag.BoolVar(&config.WetRun, "wet", false, "wet run")

	if err := flag.Parse(os.Args[2:]); err != nil {
		panic(err)
	}

	if err := preparegspatialjp.Command(&config); err != nil {
		panic(err)
	}
}
