package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"path"
	"strings"

	geojson "github.com/paulmach/go.geojson"
	"github.com/samber/lo"
	"github.com/spf13/afero"
	"github.com/vincent-petithory/dataurl"
)

const (
	wallHeight        = 100
	wallImageName     = "yellow_gradient.png"
	billboardImageDir = "billboard_image"
)

var (
	//go:embed yellow_gradient.png
	wallImage        []byte
	wallImageDataURL = dataurl.New(wallImage, http.DetectContentType(wallImage)).String()
)

type Convert struct {
	Input    string   `help:"入力元ディレクトリパス。デフォルトはカレントディレクトリです。"`
	Output   string   `help:"出力先ディレクトリパス。デフォルトはカレントディレクトリです。"`
	InputFS  afero.Fs `opts:"-"`
	OutputFS afero.Fs `opts:"-"`
}

func (c *Convert) Execute() error {
	if c.Input == "" {
		c.InputFS = afero.NewBasePathFs(afero.NewOsFs(), c.Input)
	}
	if c.Output == "" || path.Clean(c.Output) == "." {
		c.OutputFS = afero.NewOsFs()
	} else {
		c.OutputFS = afero.NewBasePathFs(afero.NewOsFs(), c.Output)
	}
	return c.execute()
}

func (c *Convert) execute() error {
	files, err := afero.ReadDir(c.InputFS, "")
	if err != nil {
		return err
	}

	for _, fi := range files {
		t := detectType(fi)
		if t == "" {
			continue
		}

		name := fi.Name()
		fc, err := c.loadGeoJSON(name)
		if err != nil {
			return err
		}

		id := strings.TrimSuffix(name, path.Ext(name))
		var czml any
		switch t {
		case "landmark":
			czml, err = ConvertLandmark(fc, id)
		case "border":
			czml, err = ConvertBorder(fc, id)
		}

		if err != nil {
			return err
		}
		if czml == nil {
			continue
		}

		if err := c.writeCZML(id, czml); err != nil {
			return err
		}
	}

	return nil
}

func (c *Convert) loadGeoJSON(path string) (*geojson.FeatureCollection, error) {
	f, err := c.InputFS.Open(path)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = f.Close()
	}()

	fc := &geojson.FeatureCollection{}
	if err := json.NewDecoder(f).Decode(&fc); err != nil {
		return fc, err
	}
	return fc, nil
}

func (c *Convert) writeCZML(name string, d any) error {
	f, err := c.OutputFS.Create(name + ".czml")
	if err != nil {
		return err
	}

	defer func() {
		_ = f.Close()
	}()

	return json.NewEncoder(f).Encode(d)
}

func detectType(fi fs.FileInfo) string {
	n := fi.Name()
	ext := path.Ext(n)
	if ext != ".geojson" && ext != ".json" {
		return ""
	}
	fn := strings.TrimSuffix(n, ext)
	if strings.HasSuffix(fn, "_landmark") || strings.HasSuffix(fn, "_station") {
		return "landmark"
	} else if strings.HasSuffix(fn, "_border") {
		return "border"
	}
	return ""
}

// ConvertLandmark は国土数値情報を基に作成されたランドマーク・鉄道駅GeoJSONデータをPLATEAU VIEW用のCZMLに変換します。
func ConvertLandmark(fc *geojson.FeatureCollection, id string) (any, error) {
	packets := make([]any, 0, len(fc.Features))
	for i, f := range fc.Features {
		if len(f.Geometry.Point) < 2 {
			continue
		}
		if len(f.Geometry.Point) == 2 {
			f.Geometry.Point = append(f.Geometry.Point, 0)
		}

		name, _ := f.PropertyString("名称")
		if name == "" {
			name, _ = f.PropertyString("駅名")
		}
		if name == "" {
			continue
		}

		image, err := GenerateLandmarkImage(name)
		if err != nil {
			return nil, err
		}

		packets = append(packets, map[string]any{
			"id":          fmt.Sprintf("%s_%d", id, i),
			"name":        name,
			"description": name,
			"billboard": map[string]any{
				"eyeOffset": map[string]any{
					"cartesian": []int{0, 0, 0},
				},
				"horizontalOrigin": "CENTER",
				"image":            dataurl.New(image, http.DetectContentType(image)).String(),
				"pixelOffset": map[string]any{
					"cartesian2": []int{0, 0},
				},
				"scale":          0.5,
				"show":           true,
				"verticalOrigin": "BOTTOM",
				"sizeInMeters":   true,
			},
			"position": map[string]any{
				"cartographicDegrees": f.Geometry.Point,
			},
			"properties": f.Properties,
		})
	}

	return czml(id, packets...), nil
}

// GenerateLandmarkImage はランドマーク用の画像を生成します。
func GenerateLandmarkImage(name string) ([]byte, error) {
	return nil, nil
}

// ConvertBorder は国土数値情報を基に作成された行政界GeoJSONデータをPLATEAU VIEW用のCZMLに変換します。
func ConvertBorder(fc *geojson.FeatureCollection, id string) (any, error) {
	packets := make([]any, 0, len(fc.Features))
	for i, f := range fc.Features {
		if len(f.Geometry.Polygon) == 0 || len(f.Geometry.Polygon[0]) == 0 {
			continue
		}

		positions := lo.FlatMap(f.Geometry.Polygon[0], func(p []float64, _ int) []float64 {
			if len(p) < 2 {
				return nil
			}
			return []float64{p[0], p[1], wallHeight}
		})

		packets = append(packets, map[string]any{
			"id": fmt.Sprintf("%s_%d", id, i+1),
			"wall": map[string]any{
				"material": map[string]any{
					"image": map[string]any{
						"image":       wallImageDataURL,
						"repeat":      true,
						"transparent": true,
					},
				},
				"positions": map[string]any{
					"cartographicDegrees": positions,
				},
			},
			"properties": f.Properties,
		})
	}

	return czml(id, packets...), nil
}

func czml(name string, packets ...any) any {
	return append(
		[]any{
			map[string]any{
				"id":      "document",
				"name":    name,
				"version": "1.0",
			},
		},
		packets...,
	)
}
