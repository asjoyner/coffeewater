package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/asjoyner/rangesensor"
	"github.com/golang/glog"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
)

var (
	valvePin   = flag.String("valvePin", "5", "GPIO pin connected the water solneoid valve.")
	echoPin    = flag.String("echoPin", "13", "GPIO pin connected the HC-SR04 echo pin.")
	triggerPin = flag.String("triggerPin", "16", "GPIO pin connected the HC-SR04 trigger pin.")
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

// CloseOnSigTerm sets the valve GPIO pin low when SIGTERM is received, prints
// a notification, then exits.
func CloseOnSigTerm(valve gpio.PinIO) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nReceived SIGTERM, closing valve and exiting.")
		valve.Out(gpio.Low)
		os.Exit(0)
	}()
}

func main() {
	// Initialize the water valve
	if _, err := host.Init(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	valve := gpioreg.ByName(*valvePin)
	if valve == nil {
		fmt.Fprintln(os.Stderr, "no GPIO valve pin named: ", *valvePin)
		os.Exit(2)
	}
	if err := valve.Out(gpio.Low); err != nil {
		fmt.Fprintln(os.Stderr, "could not configure valve output pin: ", *valvePin)
		os.Exit(3)
	}
	CloseOnSigTerm(valve)

	// Initialize the HC-SR04 sensor
	s, err := rangesensor.New(*echoPin, *triggerPin)
	if err != nil {
		fmt.Println("could not configure rangesensor: ", err)
		os.Exit(4)
	}

	// Poll the sensor in a loop
	w := NewWatcher(s)
	for {
		last10 := w.History()
		latest := last10[len(last10)-1]
		if latest != nil && latest.Trustworthy() && latest.InCentimeters() > 10 {
			valve.Out(gpio.High) // Open the water valve
		} else {
			valve.Out(gpio.Low) // Close the water valve
		}

		for _, m := range last10 {
			if m != nil {
				fmt.Printf("%5.2f ", m.InCentimeters())
			}
		}
		fmt.Println()
		time.Sleep(time.Second)
	}
}
