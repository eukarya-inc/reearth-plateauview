package plateauv2

import (
	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/datacatalogutil"
	"github.com/reearth/reearthx/util"
	"github.com/samber/lo"
)

type Override struct {
	Name         string
	SubName      string
	Type         string
	TypeEn       string
	Type2        string
	Type2En      string
	Area         string
	ItemName     string
	Group        string
	Layers       []string
	Root         bool
	Order        *int
	DatasetOrder *int
}

func (o Override) Merge(p Override) Override {
	if o.Name == "" {
		o.Name = p.Name
	}
	if o.SubName == "" {
		o.SubName = p.SubName
	}
	if o.Type == "" {
		o.Type = p.Type
	}
	if o.TypeEn == "" {
		o.TypeEn = p.TypeEn
	}
	if o.Type2 == "" {
		o.Type2 = p.Type2
	}
	if o.Type2En == "" {
		o.Type2En = p.Type2En
	}
	if o.Area == "" {
		o.Area = p.Area
	}
	if o.ItemName == "" {
		o.ItemName = p.ItemName
	}
	if o.Group == "" {
		o.Group = p.Group
	}
	if len(o.Layers) == 0 {
		o.Layers = p.Layers
	}
	if !o.Root {
		o.Root = p.Root
	}
	if o.Order == nil {
		o.Order = util.CloneRef(p.Order)
	}
	if o.DatasetOrder == nil {
		o.DatasetOrder = util.CloneRef(p.DatasetOrder)
	}
	return o
}

func (o Override) LayersIfSupported(ty string) []string {
	if datacatalogutil.IsLayerSupported(ty) {
		return o.Layers
	}
	return nil
}

func (o Override) Item() ItemOverride {
	return ItemOverride{
		Name:   o.ItemName,
		Layers: o.Layers,
		Order:  lo.FromPtr(o.DatasetOrder),
	}
}

type ItemOverride struct {
	Name   string
	Layers []string
	Order  int
}

func (o ItemOverride) Merge(p ItemOverride) ItemOverride {
	if o.Name == "" {
		o.Name = p.Name
	}
	if len(o.Layers) == 0 {
		o.Layers = p.Layers
	}
	return o
}

func (o ItemOverride) LayersIfSupported(ty string) []string {
	if datacatalogutil.IsLayerSupported(ty) {
		return o.Layers
	}
	return nil
}
