package main

import (
	"log"
	"time"

	aruco "github.com/jig/go-aruco"
)

const markerID = 7

func main() {
	log.Println("Launching Python script...")
	task := aruco.NewTask()

	log.Println("Starting to run the task...")
	go task.Run()

	log.Println("Waiting for samples...")
	for !task.IsReady() {
		time.Sleep(250 * time.Millisecond)
	}
	for {
		time.Sleep(250 * time.Millisecond)
		if marker, err := task.GetMarkerPosition(markerID); err != nil {
			log.Printf("Marker %d not visible\n", markerID)
		} else {
			log.Printf("Marker %d:   Z=%.1fcm  X=%.1fcm  pose=%.0f°\n", markerID, marker.Z*100, marker.X*100, marker.PitchY)
		}
	}
}