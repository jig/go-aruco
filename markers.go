package aruco

import (
	"errors"
	"math"
)

type Markers []Marker

type Marker struct {
	Index       int         `json:"id"`
	ImageWidth  int         `json:"width"`
	ImageHeight int         `json:"height"`
	Corners     [][]Corners `json:"corners"`
}

type Corners []float64

func (markers *Markers) Marker(id int) (*Marker, error) {
	for _, marker := range *markers {
		if marker.Index == id {
			return &marker, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *Marker) Position() (dist float64, angle float64, targetAngle float64, err error) {
	return 0, m.CenterX(), 0, nil
}

func (c *Marker) Size() float64 {
	if len(c.Corners[0]) != 4 {
		return -1
	}
	corners := c.Corners[0]
	size1 := corners[1][0] - corners[0][0]
	size2 := corners[2][1] - corners[1][1]
	size3 := corners[3][1] - corners[0][1]
	size4 := corners[2][0] - corners[3][0]
	size := size1
	if size2 > size {
		size = size2
	}
	if size3 > size {
		size = size3
	}
	if size4 > size {
		size = size4
	}
	return size
}

func (c *Marker) Distance() float64 {
	// for a 185 mm marker
	return 2910.0 * (48.0 / c.Size())

	// in marker units
	// return 15.73 * (48.0 / c.Size())
}

func (c *Marker) CenterX() float64 {
	centerX := 0.0
	for i, xy := range c.Corners[0] {
		centerX += xy[0]
		if i > 4 {
			return -1
		}
	}
	return centerX/4 - float64(c.ImageWidth)/2
}

func (c *Marker) ViewAngleX() float64 {
	return c.CenterX() / 464 * 25 / 57.2958
}

func (c *Marker) ViewAngleY() float64 {
	return c.CenterY() / 464 * 25 / 57.2958
}

func (c *Marker) CenterY() float64 {
	centerY := 0.0
	for i, xy := range c.Corners[0] {
		centerY += xy[1]
		if i > 4 {
			return -1
		}
	}
	return -(centerY/4 - float64(c.ImageHeight)/2)
}

func (c *Marker) VerticalPoseRatio() float64 {
	if len(c.Corners[0]) != 4 {
		return -1
	}
	corners := c.Corners[0]
	h1 := corners[3][1] - corners[0][1]
	h2 := corners[2][1] - corners[1][1]
	w1 := corners[1][0] - corners[0][0]
	w2 := corners[2][0] - corners[3][0]
	h := h1 + h2
	w := w1 + w2
	return w / h
}

func (c *Marker) VerticalPoseAngle() float64 {
	ratio := c.VerticalPoseRatio()
	if ratio > 1 {
		ratio = 1
	}
	return math.Acos(ratio)
}
