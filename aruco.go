package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"math"
	"os/exec"
)

type Markers []Marker

type Marker struct {
	Index       int         `json:"id"`
	ImageWidth  int         `json:"width"`
	ImageHeight int         `json:"height"`
	Corners     [][]Corners `json:"corners"`
}

type Corners []float64

func main() {
	cmd := exec.Command("python", "../../jig/aruco/aruco.py")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		// log.Println(line)
		markers := Markers{}
		if err := json.Unmarshal([]byte(line), &markers); err != nil {
			log.Fatal(err)
		}

		if m7, err := markers.Marker(7); err != nil {
			log.Printf("Marker 7 not visible\n")
		} else {
			log.Printf("Marker 7: Distance %.0fmm\tPose angle %.0f°; \tcentered at %.2f\t%.2f\n", m7.Distance(), m7.VerticalPoseAngle(), m7.CenterX(), m7.CenterY())
		}

	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}

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
	return 1000.0 * (48.0 / c.Size())
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

func (c *Marker) CenterY() float64 {
	centerY := 0.0
	for i, xy := range c.Corners[0] {
		centerY += xy[1]
		if i > 4 {
			return -1
		}
	}
	return centerY/4 - float64(c.ImageHeight)/2
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
	return math.Acos(ratio) * 57.2958
}
