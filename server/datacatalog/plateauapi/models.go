package plateauapi

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/datacatalogv2/datacatalogutil"
	"github.com/samber/lo"
)

func FindItem(d Dataset, id ID) DatasetItem {
	res, _ := lo.Find(d.GetItems(), func(i DatasetItem) bool {
		return i.GetID() == id
	})
	return res
}

type Areas map[AreaType][]Area

func (a Areas) Append(cat AreaType, as []Area) {
	a[cat] = append(a[cat], as...)
}

func (a Areas) All() []Area {
	entries := lo.Entries(a)
	sort.Slice(entries, func(i, j int) bool {
		return slices.Index(AllAreaType, entries[i].Key) < slices.Index(AllAreaType, entries[j].Key)
	})
	return lo.FlatMap(entries, func(e lo.Entry[AreaType, []Area], _ int) []Area {
		return e.Value
	})
}

func (a Areas) Filter(f func(a Area) bool) []Area {
	return lo.Filter(a.All(), func(a Area, _ int) bool {
		return f(a)
	})
}

func (a Areas) Find(f func(a Area) bool) Area {
	res, _ := lo.Find(a.All(), f)
	return res
}

func (a Areas) FindByCode(code AreaCode) Area {
	return a.Find(func(a Area) bool {
		return a.GetCode() == code
	})
}

func (a Areas) FindByCodeAndType(code AreaCode, ty AreaType) Area {
	areas := a[ty]
	if areas == nil {
		return nil
	}
	res, _ := lo.Find(areas, func(a Area) bool {
		return a.GetCode() == code
	})
	return res
}

func (a *Areas) Area(id ID) Area {
	for _, area := range a.All() {
		if area.GetID() == id {
			return area
		}
	}
	return nil
}

func DatasetTypeCategoryFromDataset(d Dataset) DatasetTypeCategory {
	switch d.(type) {
	case PlateauDataset:
		return DatasetTypeCategoryPlateau
	case *PlateauDataset:
		return DatasetTypeCategoryPlateau
	case RelatedDataset:
		return DatasetTypeCategoryRelated
	case *RelatedDataset:
		return DatasetTypeCategoryRelated
	case GenericDataset:
		return DatasetTypeCategoryGeneric
	case *GenericDataset:
		return DatasetTypeCategoryGeneric
	}
	return ""
}

type Datasets map[DatasetTypeCategory][]Dataset

func (d Datasets) Append(cat DatasetTypeCategory, ds []Dataset) {
	d[cat] = append(d[cat], ds...)
}

func (d Datasets) All() []Dataset {
	entries := lo.Entries(d)
	sort.Slice(entries, func(i, j int) bool {
		return slices.Index(AllDatasetTypeCategory, entries[i].Key) < slices.Index(AllDatasetTypeCategory, entries[j].Key)
	})
	return lo.FlatMap(entries, func(e lo.Entry[DatasetTypeCategory, []Dataset], _ int) []Dataset {
		return e.Value
	})
}

func (d Datasets) Filter(f func(d Dataset) bool) []Dataset {
	return lo.Filter(d.All(), func(d Dataset, _ int) bool {
		return f(d)
	})
}

func (d *Datasets) Dataset(id ID) Dataset {
	for _, ds := range d.All() {
		if ds.GetID() == id {
			return ds
		}
	}
	return nil
}

func (d *Datasets) Item(id ID) DatasetItem {
	for _, ds := range d.All() {
		if item := FindItem(ds, id); item != nil {
			return item
		}
	}
	return nil
}

type DatasetTypes map[DatasetTypeCategory][]DatasetType

func (d DatasetTypes) Append(cat DatasetTypeCategory, ds []DatasetType) {
	d[cat] = append(d[cat], ds...)
}

func (d DatasetTypes) All() []DatasetType {
	entries := lo.Entries(d)
	sort.Slice(entries, func(i, j int) bool {
		return slices.Index(AllDatasetTypeCategory, entries[i].Key) < slices.Index(AllDatasetTypeCategory, entries[j].Key)
	})
	return lo.FlatMap(entries, func(e lo.Entry[DatasetTypeCategory, []DatasetType], _ int) []DatasetType {
		return e.Value
	})
}

func (d DatasetTypes) DatasetTypesByCategories(categories []DatasetTypeCategory) (res []DatasetType) {
	for _, cat := range categories {
		res = append(res, d[cat]...)
	}

	slices.SortStableFunc(res, func(a, b DatasetType) int {
		return strings.Compare(a.GetCode(), b.GetCode())
	})

	return res
}

func (d DatasetTypes) Filter(f func(d DatasetType) bool) []DatasetType {
	return lo.Filter(d.All(), func(d DatasetType, _ int) bool {
		return f(d)
	})
}

func (d *DatasetTypes) DatasetType(id ID) DatasetType {
	for _, ds := range d.All() {
		if ds.GetID() == id {
			return ds
		}
	}
	return nil
}

func (s *PlateauSpec) Minor(name string) *PlateauSpecMinor {
	for _, minor := range s.MinorVersions {
		if minor.Name == name {
			return minor
		}
	}
	return nil
}

func FindSpecMinorByName(specs []PlateauSpec, name string) *PlateauSpecMinor {
	for _, spec := range specs {
		if specMinor := spec.Minor(name); specMinor != nil {
			return specMinor
		}
	}
	return nil
}

func stageFromDataset(ds Dataset) string {
	admin := ds.GetAdmin()
	if admin == nil {
		return ""
	}

	return stageFrom(admin)
}

func stageFromCityGMLDataset(ds *CityGMLDataset) string {
	admin := ds.Admin
	if admin == nil {
		return ""
	}

	return stageFrom(admin)
}

func stageFrom(admin any) string {
	if admin == nil {
		return ""
	}

	m, ok := admin.(map[string]any)
	if !ok || m == nil {
		return ""
	}

	stage, ok := m["stage"]
	if !ok {
		return ""
	}

	s, ok := stage.(string)
	if !ok {
		return ""
	}

	return s
}

func (d PlateauDatasetType) GetYear() int {
	return d.Year
}

func (d PlateauSpec) GetYear() int {
	return d.Year
}

func (d PlateauSpecMinor) GetYear() int {
	return d.Year
}

var _ YearNode = (*PlateauDatasetType)(nil)
var _ YearNode = (*PlateauSpec)(nil)
var _ YearNode = (*PlateauSpecMinor)(nil)

func IsLayerSupported(format DatasetFormat) bool {
	return datacatalogutil.IsLayerSupported(strings.ToLower(string(format)))
}

func CityGMLDatasetIDFrom(areaCode AreaCode) ID {
	return NewID(areaCode.String(), TypeCityGML)
}

type VagueID interface {
	VagueID() string
}

func (a *PlateauDataset) VagueID() string {
	return fmt.Sprintf("d_%s_%s", mostDetailedAreaCodeFrom(a), a.TypeCode)
}

func getVagueID(n any) string {
	if v, ok := n.(VagueID); ok && v != nil {
		return v.VagueID()
	}
	return ""
}
