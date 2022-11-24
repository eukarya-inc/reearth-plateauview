package indexer

import (
	"math"
)

const (
	_EPSILON1           = 0.1
	_EPSILON12          = 0.000000000001
	_EPSILON14          = 0.00000000000001
)

const (
	_PI     = math.Pi
	_TWO_PI = 2 * math.Pi
	_DEGREES_PER_RADIAN = 180.0 / math.Pi
)

func zeroToTwoPi(angle float64) float64 {
	if angle >= 0 && angle <= _TWO_PI {
		// Early exit if the input is already inside the range. This avoids
		// unnecessary math which could introduce floating point error.
		return angle
	}
	mod := math.Mod(angle, _TWO_PI)
	if math.Abs(mod) < _EPSILON14 && math.Abs(angle) > _EPSILON14 {
		return _TWO_PI
	}
	return mod
}

func negativePiToPi(angle float64) float64 {
	if angle >= -_PI && angle <= _PI {
		// Early exit if the input is already inside the range. This avoids
		// unnecessary math which could introduce floating point error.
		return angle
	}
	return zeroToTwoPi(angle+_PI) - _PI
}

func toDegrees(radians float64) float64 {
	return radians * _DEGREES_PER_RADIAN
}
