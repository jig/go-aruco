package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"

	aruco "github.com/jig/go-aruco"
)

func main() {
	cmd := exec.Command("python", "/home/pi/git/src/github.com/jig/aruco/aruco.py")
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
		markers := aruco.Markers{}
		if err := json.Unmarshal([]byte(line), &markers); err != nil {
			log.Fatal(err)
		}

		if m7, err := markers.Marker(7); err != nil {
			log.Printf("Marker 7 not visible\n")
		} else {
			fmt.Printf("                                                                                                                    \r")
			// log.Printf("Marker 7: Distance %.0fmm\tPose angle %.0f°; \tcentered at %.2f\t%.2f\r", m7.Distance(), m7.VerticalPoseAngle(), m7.CenterX(), m7.CenterY())
			// fmt.Printf("Marker 7: Distance %4.0fmm\tPose angle %2.0f°; \tView angle %2.0f° (vertical: %2.0f°)\r", m7.Distance(), m7.VerticalPoseAngle()*57.2958, m7.ViewAngleX()*57.2958, m7.ViewAngleY()*57.2958)
			curvature := aruco.InverseChord(m7.Distance(), m7.VerticalPoseAngle()*2)
			if curvature == 0 {
				fmt.Printf("Marker 7: Straight line\r")
			} else {
				fmt.Printf("Marker 7: Radius: %.0f mm\r", 1/curvature)
			}
		}

	}
	if err := cmd.Wait(); err != nil {
		log.Fatalf("Python subroutine failed: %s", err)
	}
}
