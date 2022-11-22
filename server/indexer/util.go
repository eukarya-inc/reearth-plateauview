package indexer

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"

	"github.com/qmuntal/gltf"
	b3dms "github.com/reearth/go3dtiles/b3dm"
	tiles "github.com/reearth/go3dtiles/tileset"
	"gonum.org/v1/gonum/mat"
)

type CesiumRTC struct {
	Center []float64
}

func getRtcTransform(ft *b3dms.B3dmFeatureTable, gltf *gltf.Document) *mat.Dense {
	rtcCenter := ft.RtcCenter
	if len(rtcCenter) == 0 {
		var temp CesiumRTC
		json.Unmarshal(gltf.Extensions["CESIUM_RTC"].(json.RawMessage), &temp)
		rtcCenter = temp.Center
	}
	rtcTransform := eyeMat(4)
	if len(rtcCenter) > 0 {
		rtcTransform = mat4FromCartesian(cartesianFromSlice(rtcCenter))
	}
	return rtcTransform
}

// Creates a rotation matrix around the x-axis.
func getYUpToZUp() *mat.Dense {
	sinAngle, cosAngle := math.Sincos(math.Pi / 2)
	d := []float64{
		1.0,
		0.0,
		0.0,
		0.0,
		0.0,
		cosAngle,
		sinAngle,
		0.0,
		0.0,
		-sinAngle,
		cosAngle,
		0.0,
		0.0,
		0.0,
		0.0,
		1.0,
	}

	return mat.NewDense(4, 4, d)
}

func getZUpTransform(ts *tiles.Tileset) *mat.Dense {
	// discuss if we need gltfAxisUpAxis
	upAxis := AXIS_Y
	transform := eyeMat(4)
	if upAxis == AXIS_Y {
		transform = getYUpToZUp()
	}
	return transform
}

var floatType = reflect.TypeOf(float64(0))

func getFloat(unk interface{}) (float64, error) {
	v := reflect.ValueOf(unk)
	v = reflect.Indirect(v)
	if !v.Type().ConvertibleTo(floatType) {
		return 0, fmt.Errorf("cannot convert %v to float64", v.Type())
	}
	fv := v.Convert(floatType)
	return fv.Float(), nil
}

var intType = reflect.TypeOf(int(0))

func getInt(unk interface{}) (int64, error) {
	v := reflect.ValueOf(unk)
	v = reflect.Indirect(v)
	if !v.Type().ConvertibleTo(intType) {
		return 0, fmt.Errorf("cannot convert %v to int", v.Type())
	}
	iv := v.Convert(intType)
	return iv.Int(), nil
}

func Map(elem []interface{}, f func(interface{}) (float64, error)) ([]float64, error) {
	result := make([]float64, len(elem))
	var err error
	for i, v := range elem {
		result[i], err = f(v)
		if err != nil {
			return nil, fmt.Errorf("failed to apply function the function: %w", err)
		}
	}
	return result, err
}

func minMaxOfSlice(array []float64) (float64, float64) {
	var max float64 = array[0]
	var min float64 = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
