package util

import (
	"log"
	"time"
)

// TimeTrack functions to measure execution time.
// usage: defer util.TimeTrack(time.Now(), "function")
func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
