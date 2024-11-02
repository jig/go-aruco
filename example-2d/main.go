package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jig/go-aruco"
	exec "github.com/jig/go-exec"
)

const markerID = 18

func main() {
	pythonAruco, err := exec.NewCmd("python", "-c", markersPythonCode)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	go func() {
		exitCode, err := pythonAruco.Run()
		if err != nil {
			fmt.Println(err)
		}
		os.Exit(exitCode)
	}()
	for data := range pythonAruco.Output() {
		if data.EOF {
			break
		} else if data.IsStderr {
			fmt.Println("ERR:", data.Value)
		} else {
			markers := aruco.Markers{}
			if err := json.Unmarshal([]byte(data.Value), &markers); err != nil {
				log.Println(err)
				continue
			}
			if marker, err := markers.Marker(markerID); err != nil {
				log.Printf("Marker %d not visible\n", markerID)
			} else {
				log.Printf("Marker %d:   Z=%.1fcm  X=%.1fcm  pose=%.0f°\n", markerID, marker.Z*100, marker.X*100, marker.PitchY)
			}
		}
	}
}

//go:embed markers.py
var markersPythonCode string
