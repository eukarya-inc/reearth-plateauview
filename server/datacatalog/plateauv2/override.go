package plateauv2

import "github.com/eukarya-inc/reearth-plateauview/server/datacatalog/datacatalogutil"

type Override struct {
	Name   string
	Type   string
	TypeEn string
	Area   string
	Layers []string
}

func (o Override) Merge(p Override) Override {
	if o.Name == "" {
		o.Name = p.Name
	}
	if o.Type == "" {
		o.Type = p.Type
	}
	if o.TypeEn == "" {
		o.TypeEn = p.TypeEn
	}
	if o.Area == "" {
		o.Area = p.Area
	}
	if len(o.Layers) == 0 {
		o.Layers = p.Layers
	}
	return o
}

func (o Override) LayersIfSupported(ty string) []string {
	if datacatalogutil.IsLayerSupported(ty) {
		return o.Layers
	}
	return nil
}
