package preparegspatialjp

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/samber/lo"
)

func TestCommand(t *testing.T) {
	t.Skip()
	lo.Must0(godotenv.Load("../.env"))

	conf := Config{
		CMSURL:   os.Getenv("REEARTH_CMS_URL"),
		CMSToken: os.Getenv("REEARTH_CMS_TOKEN"),
		// ProjectID:  "",
		// CityItemID: "",
	}

	if err := Command(&conf); err != nil {
		t.Fatal(err)
	}
}
