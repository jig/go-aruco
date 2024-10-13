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
	markersPositions MarkersPositions
	cmd              *exec.Cmd
	scanner          *bufio.Scanner
	mu               sync.RWMutex
	ready            bool
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
			if err := json.Unmarshal([]byte(line), &task.markersPositions); err != nil {
				log.Println(err)
			}
			task.ready = task.isReady()
		}()
	}
	if err := task.cmd.Wait(); err != nil {
		log.Fatalf("Python subroutine failed: %s", err)
	}
}

func (task *Task) GetMarkersPositions() (MarkersPositions, error) {
	task.mu.RLock()
	defer task.mu.RUnlock()

	return task.markersPositions, nil
}

func (task *Task) GetMarkerPosition(markerID int) (*MarkerPosition, error) {
	task.mu.RLock()
	defer task.mu.RUnlock()

	for _, m := range task.markersPositions {
		if m.ID == markerID {
			return task.markersPositions.MarkerPosition(markerID)
		}
	}
	return nil, errors.New("marker not found")
}

func (task *Task) isReady() bool {
	for _, m := range task.markersPositions {
		if m.ID == -1 {
			return false
		}
	}
	return true
}

func (task *Task) IsReady() bool {
	return task.ready
}
