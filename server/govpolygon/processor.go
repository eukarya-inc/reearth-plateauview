package govpolygon

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	geojson "github.com/paulmach/go.geojson"
	"github.com/reearth/reearthx/log"
	"github.com/rubenv/topojson"
)

type Processor struct {
	dirpath string
	key1    string
	key2    string
}

func NewProcessor(dirpath, key1, key2 string) *Processor {
	return &Processor{dirpath: dirpath, key1: key1, key2: key2}
}

func (p *Processor) ComputeGeoJSON(ctx context.Context, values []string, citycodem map[string]string) (*geojson.FeatureCollection, error) {
	features, err := loadFeatures(context.Background(), p.dirpath)
	if err != nil {
		return nil, err
	}

	if len(features) == 0 {
		return nil, fmt.Errorf("no features found")
	}

	return computeGeojsonFeatures(features, p.key1, p.key2, values, citycodem), nil
}

func computeGeojsonFeatures(features []*geojson.Feature, key1, key2 string, values []string, citycodem map[string]string) *geojson.FeatureCollection {
	valueSet := map[string]struct{}{}
	for _, v := range values {
		valueSet[v] = struct{}{}
	}

	result := geojson.NewFeatureCollection()
	for _, f := range features {
		v1, ok := f.Properties[key1].(string)
		if !ok {
			continue
		}

		v2, ok := f.Properties[key2].(string)
		if !ok {
			continue
		}

		value := v1 + v2
		if _, ok := valueSet[value]; ok {
			properties := map[string]any{
				"pref": v1,
				"city": v2,
				"code": citycodem[value],
			}
			if citycodem != nil {
				if code := citycodem[value]; code != "" {
					properties["code"] = code
					f.ID = code
				}
			}
			f.Properties = properties
			result.AddFeature(f)
		}
	}

	return result
}

func loadFeatures(ctx context.Context, dirpath string) ([]*geojson.Feature, error) {
	files, err := os.ReadDir(dirpath)
	if err != nil {
		return nil, err
	}

	var features []*geojson.Feature
	for _, f := range files {
		name := f.Name()
		if f.IsDir() || filepath.Ext(name) != ".topojson" {
			continue
		}

		p := filepath.Join(dirpath, name)
		file, err := os.ReadFile(p)
		if err != nil {
			log.Debugfc(ctx, "govpolygon: error reading file (%s): %s", name, err)
			continue
		}

		topology, err := topojson.UnmarshalTopology(file)
		if err != nil {
			log.Debugfc(ctx, "govpolygon: error unmarshaling topojson (%s): %s", name, err)
			continue
		}

		f := topology.ToGeoJSON()
		if f == nil {
			log.Debugfc(ctx, "govpolygon: error converting topojson to geojson (%s)", name)
			continue
		}

		features = append(features, f.Features...)
	}

	return features, nil
}
