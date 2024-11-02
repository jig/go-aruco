package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/jig/go-aruco"
)

func main() {
	pythonAruco, err := aruco.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		if _, err := pythonAruco.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		pythonAruco.Dispatch(func(markers aruco.Markers) {
			for _, marker := range markers {
				log.Printf("Marker %d:   Z=%.1fcm  X=%.1fcm  pose=%.0f°\n", marker.ID, marker.Z*100, marker.X*100, marker.PitchY)
			}
		})
		wg.Done()
	}()
	wg.Wait()
}
