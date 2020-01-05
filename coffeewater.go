package main

import (
	"fmt"
	"os"
	"time"

	"github.com/asjoyner/rangesensor"
)

func main() {

	s, err := rangesensor.New("22", "23")
	if err != nil {
		fmt.Println("could not configure pin: ", err)
		os.Exit(1)
	}

	for {
		var set []float32
		for i := 0; i < 10; i++ {
			time.Sleep(time.Duration(100) * time.Millisecond)
			distance, err := s.MeasureDistance()
			if err != nil {
				continue // failure to read the sensor
			}
			cm := distance.InCentimeters()
			if cm > 200 {
				continue // spurious result from the sensor
			}
			set = append(set, cm)
		}
		numSamples := len(set)
		if numSamples < 3 {
			fmt.Printf("Distance: unknown\n")
			continue
		}
		var avg float32
		for _, cm := range set {
			avg += cm
		}
		avg = avg / float32(numSamples)
		fmt.Printf("Distance: %5.2f cm (%d samples)\n", avg, numSamples)
	}
}
