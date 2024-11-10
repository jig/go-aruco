package aruco

import (
	_ "embed"
	"errors"
	"math"
)

type Task struct {
	MarkersChannel chan Markers
}

type Markers []Marker

type Marker struct {
	ID     int     `json:"id"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Z      float64 `json:"z"`
	RollX  float64 `json:"roll-x"`
	PitchY float64 `json:"pitch-y"`
	YawZ   float64 `json:"yaw-z"`
}

func (markers Markers) Marker(id int) (*Marker, error) {
	for _, marker := range markers {
		if marker.ID == id {
			return &marker, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *Marker) GetDistanceXZ() float64 {
	return math.Sqrt(m.X*m.X + m.Z*m.Z)
}

func (m *Marker) GetDistanceXYZ() float64 {
	return math.Sqrt(m.X*m.X + m.Y*m.Y + m.Z*m.Z)
}

func (m *Marker) GetPoseY() float64 {
	return m.PitchY
}
