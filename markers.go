package aruco

import "errors"

type MarkersPositions []MarkerPosition

type MarkerPosition struct {
	ID     int     `json:"id"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Z      float64 `json:"z"`
	RollX  float64 `json:"roll-x"`
	PitchY float64 `json:"pitch-y"`
	YawZ   float64 `json:"yaw-z"`
}

func (markersPositions *MarkersPositions) MarkerPosition(id int) (*MarkerPosition, error) {
	for _, marker := range *markersPositions {
		if marker.ID == id {
			return &marker, nil
		}
	}
	return nil, errors.New("not found")
}
