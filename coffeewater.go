package main

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/asjoyner/rangesensor"
	"github.com/golang/glog"
)

// Watcher keeps track of the distance seen by a Sensor
type Watcher struct {
	sensor  *rangesensor.Sensor
	history []*rangesensor.Measurement
	mu      sync.RWMutex // protects history
}

// NewWatcher initializes and returns a Watcher
func NewWatcher(s *rangesensor.Sensor) *Watcher {
	w := &Watcher{sensor: s, history: make([]*rangesensor.Measurement, 10)}
	go w.followPin()
	return w
}

// History returns the last 10 measurements
func (w *Watcher) History() []*rangesensor.Measurement {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.history
}

func (w *Watcher) followPin() {
	for {
		res, _ := w.sensor.MeasureDistance()
		w.mu.Lock()
		w.history = append(w.history[1:], res)
		w.mu.Unlock()
		time.Sleep(time.Duration(100) * time.Millisecond)
	}
}

// avgDistance returns the averaged measured distance over the last 1 second
func (w *Watcher) avgDistance() (float32, error) {
	var set []float32
	for i := 0; i < 10; i++ {
		time.Sleep(time.Duration(100) * time.Millisecond)
		distance, err := w.sensor.MeasureDistance()
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
		glog.Infof("Distance: unknown")
		return 0.0, errors.New("distance unknown")
	}
	var avg float32
	for _, cm := range set {
		avg += cm
	}
	avg = avg / float32(numSamples)
	glog.Infof("Distance: %5.2f cm (%d samples)", avg, numSamples)
	return avg, nil
}

func main() {
	s, err := rangesensor.New("13", "16")
	if err != nil {
		fmt.Println("could not configure pin: ", err)
		os.Exit(1)
	}

	w := NewWatcher(s)
	for {
		last10 := w.History()
		latest := last10[len(last10)-1]
		if latest != nil && latest.Trustworthy() && latest.InCentimeters() > 10 {
			fmt.Println("MAIN SCREEN^WPUMP TURN ON")
		}

		for _, m := range last10 {
			if m != nil {
				fmt.Printf("%+v ", *m)
			}
		}
		fmt.Println()
		time.Sleep(time.Second)
	}
}
