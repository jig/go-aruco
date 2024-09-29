package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"

	aruco "github.com/jig/go-aruco"
)

const markerID = 13

func main() {
	cmd := exec.Command("python", "/home/pi/git/src/github.com/jig/aruco/aruco.py")
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
		markers := aruco.Markers{}
		if err := json.Unmarshal([]byte(line), &markers); err != nil {
			log.Fatal(err)
		}

		if m2, err := markers.Marker(markerID); err != nil {
			log.Printf("Marker %d not visible\n", markerID)
		} else {
			fmt.Printf("                                                                              \r")
			fmt.Printf("Marker %d: Distance %.0fmm\tPose angle %.0f°; \tcentered at (%.2f,%.2f)\r",
				markerID, m2.Distance(), m2.VerticalPoseAngle()*57.2958, m2.CenterX(), m2.CenterY())
			// fmt.Printf("Marker %d: Distance %.0fmm\tPose angle %.0f°; \tView angle %2.0f° (vertical: %2.0f°)\r",
			// 	markerID, m2.Distance(), m2.VerticalPoseAngle()*57.2958, m2.ViewAngleX()*57.2958, m2.ViewAngleY()*57.2958)

			// curvature := aruco.InverseChord(m2.Distance(), m2.VerticalPoseAngle()*2)
			// if curvature == 0 {
			// 	fmt.Printf("Marker %d: Straight line\r", markerID)
			// } else {
			// 	fmt.Printf("Marker %d: Radius: %.0f mm\r", markerID, 1/curvature)
			// }
		}

	}
	if err := cmd.Wait(); err != nil {
		log.Fatalf("Python subroutine failed: %s", err)
	}
}
