package preparegspatialjp

import (
	"github.com/k0kubun/pp/v3"
)

var ppp *pp.PrettyPrinter

func init() {
	ppp = pp.New()
	ppp.SetColoringEnabled(false)
}
