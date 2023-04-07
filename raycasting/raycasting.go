package raycasting

import (
	"math"

	"golang.org/x/image/math/f64"
)

func PolarCoordinatesToDirectionVector(phi float64, theta float64) f64.Vec3 {
	var cos_phi = math.Cos(phi);
	return f64.Vec3{(cos_phi * math.Sin(theta)), math.Sin(phi), (cos_phi * math.Cos(theta))}
}
