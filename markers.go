package aruco

import (
	"bufio"
	_ "embed"
	"errors"
	"os/exec"
)

type Task struct {
	MarkersChannel chan Markers

	cmd     *exec.Cmd
	scanner *bufio.Scanner
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
