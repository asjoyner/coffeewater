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
		distance, err := s.MeasureDistance()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Distance: %5.2f cm", distance.InCentimeters())
			fmt.Printf(" (%d us)\n", distance.InMicroseconds())
		}
		time.Sleep(time.Duration(250) * time.Millisecond)
	}

}
