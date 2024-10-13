package aruco

import (
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"os/exec"
	"sync"
)

type Task struct {
	markers []Marker
	cmd     *exec.Cmd
	scanner *bufio.Scanner
	mu      sync.RWMutex
	ready   bool
}

type Marker struct {
	ID     int     `json:"id"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Z      float64 `json:"z"`
	RollX  float64 `json:"roll-x"`
	PitchY float64 `json:"pitch-y"`
	YawZ   float64 `json:"yaw-z"`
}

func NewTask() *Task {
	cmd := exec.Command("python", "/home/pi/git/src/github.com/jig/go-aruco/markers.py")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	return &Task{
		cmd:     cmd,
		scanner: bufio.NewScanner(stdout),
		ready:   false,
	}
}

func (task *Task) Run() {
	for task.scanner.Scan() {
		func() {
			line := task.scanner.Text()
			task.mu.Lock()
			defer task.mu.Unlock()
			if err := json.Unmarshal([]byte(line), &task.markers); err != nil {
				log.Println(err)
			}
			task.ready = task.isReady()
		}()
	}
	if err := task.cmd.Wait(); err != nil {
		log.Fatalf("Python subroutine failed: %s", err)
	}
}

func (task *Task) Marker(id int) (*Marker, error) {
	for _, marker := range task.markers {
		if marker.ID == id {
			return &marker, nil
		}
	}
	return nil, errors.New("not found")
}

func (task *Task) isReady() bool {
	for _, m := range task.markers {
		if m.ID == -1 {
			return false
		}
	}
	return true
}

func (task *Task) IsReady() bool {
	return task.ready
}
