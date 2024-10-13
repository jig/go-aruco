package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"

	aruco "github.com/jig/go-aruco"
)

const markerID = 7

func main() {
	cmd := exec.Command("python", "/home/pi/git/src/github.com/jig/go-aruco/markers.py")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("starting")
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		// log.Println(line)
		markersPositions := aruco.MarkersPositions{}
		if err := json.Unmarshal([]byte(line), &markersPositions); err != nil {
			log.Fatal(err)
		}

		if marker, err := markersPositions.MarkerPosition(markerID); err != nil {
			log.Printf("Marker %d not visible\n", markerID)
		} else {
			log.Printf("Marker %d:   Z=%.1f  X=%.1f  pose=%.0f°\n", markerID, marker.Z*100, marker.X*100, marker.PitchY)
		}

	}
	if err := cmd.Wait(); err != nil {
		log.Fatalf("Python subroutine failed: %s", err)
	}
}
