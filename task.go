package aruco

import (
	"bufio"
	"encoding/json"
	"log"
	"os/exec"
	"sync"
)

type TaskPosition struct {
	markersPositions MarkersPositions
	cmd              *exec.Cmd
	scanner          *bufio.Scanner
	mu               sync.RWMutex
}

func NewTaskPosition() *TaskPosition {
	cmd := exec.Command("python", "/home/pi/git/src/github.com/jig/borinot7/aruco/python/aruco-pos.py")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	return &TaskPosition{
		cmd:     cmd,
		scanner: bufio.NewScanner(stdout),
	}
}

func (task *TaskPosition) Run() {
	for task.scanner.Scan() {
		func() {
			line := task.scanner.Text()
			task.mu.Lock()
			defer task.mu.Unlock()
			if err := json.Unmarshal([]byte(line), &task.markersPositions); err != nil {
				log.Println(err)
			}
			// log.Println("...")
		}()
	}
	if err := task.cmd.Wait(); err != nil {
		log.Fatalf("Python subroutine failed: %s", err)
	}
}

func (task *TaskPosition) GetMarkersPositions() (MarkersPositions, error) {
	task.mu.RLock()
	defer task.mu.RUnlock()

	return task.markersPositions, nil
}

func (task *TaskPosition) GetMarkerPosition(markerID int) (*MarkerPosition, error) {
	task.mu.RLock()
	defer task.mu.RUnlock()

	return task.markersPositions.MarkerPosition(markerID)
}
