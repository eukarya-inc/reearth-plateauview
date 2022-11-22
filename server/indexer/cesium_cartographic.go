package indexer

import (
	"fmt"
	"math"
)

var (
	_wgs84OneOverRadii = Cartesian3{
		X: 1.0 / 6378137.0,
		Y: 1.0 / 6378137.0,
		Z: 1.0 / 6356752.3142451793,
	}
	_wgs84OneOverRadiiSquared = Cartesian3{
		X: 1.0 / (6378137.0 * 6378137.0),
		Y: 1.0 / (6378137.0 * 6378137.0),
		Z: 1.0 / (6356752.3142451793 * 6356752.3142451793),
	}
	_wgs84CenterToleranceSquared = _EPSILON1
)

type Cartographic struct {
	Longitude float64
	Latitude  float64
	Height    float64
}

func cartographicFromCartesian3(cs *Cartesian3) (*Cartographic, error) {
	oneOverRadii := _wgs84OneOverRadii
	oneOverRadiiSquared := _wgs84OneOverRadiiSquared
	centerToleranceSquared := _wgs84CenterToleranceSquared

	p, err := scaleToGeodeticSurface(cs, &oneOverRadii, &oneOverRadiiSquared, centerToleranceSquared)

	if err != nil {
		return nil, fmt.Errorf("failed to do scaleToGeoticSurface transformation: %v", err)
	}

	n := multiplyCartesian3Components(p, &oneOverRadiiSquared)
	n = normalize(n)

	h := subtract(cs, p)

	longitude := math.Atan2(n.Y, n.X)
	latitude := math.Asin(n.Z)
	height := sign(dot(h, cs)) * h.magnitude()

	return &Cartographic{
		Longitude: longitude,
		Latitude:  latitude,
		Height:    height,
	}, nil

}
