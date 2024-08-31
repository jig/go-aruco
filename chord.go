package aruco

import "math"

func InverseChord(dist float64, angle float64) (curvature float64) {
	return 2 * math.Sin(angle/2) / dist
}
