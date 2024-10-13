package aruco

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"errors"
	"log"
	"os/exec"
)

type Task struct {
	MarkersChannel chan Markers

	cmd     *exec.Cmd
	scanner *bufio.Scanner
}

//go:embed markers.py
var markersPythonCode string

func NewTask() *Task {
	// cmd := exec.Command("python", "/home/pi/git/src/github.com/jig/go-aruco/markers.py")
	cmd := exec.Command("python", "-c", markersPythonCode)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	return &Task{
		MarkersChannel: make(chan Markers),
		cmd:            cmd,
		scanner:        bufio.NewScanner(stdout),
	}
}

func (task *Task) Run() {
	for task.scanner.Scan() {
		line := task.scanner.Text()
		markers := []Marker{}
		if err := json.Unmarshal([]byte(line), &markers); err != nil {
			log.Println(err)
		}
		task.MarkersChannel <- markers
	}
	if err := task.cmd.Wait(); err != nil {
		log.Fatalf("Python subroutine failed: %s", err)
	}
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
