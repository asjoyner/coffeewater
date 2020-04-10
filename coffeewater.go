package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/asjoyner/rangesensor"
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

// safeToFill sanity checks values, and indicates if filling is appropriate.
//
// top defines the short expected measured distance, in centimeters, from the
// sensor to the literal high water mark
//
// bottom defines the longest expected measured distance, in centimeters, to
// the bottom of the container
//
// target indicates the measured height of the water level, in centimeters from
// the sensor.
//
// fill indicates the measured distance to where coffeebot should start
// refilling the water resivoir.  Together with target, this defines the two
// values used to maintain the desired water level.
//
// values is the ordered history of measurements, in centimeters, beginning
// with the most recent.
//
// The bool return value indicates if it is safe to fill, and the string
// describes the current status of the water level.
func safeToFill(top, bottom float32, values []*rangesensor.Measurement) (bool, float32, string) {
	if len(values) < 10 {
		return false, 0.0, fmt.Sprintf("insufficient history: %v", len(values))
	}
	var summary strings.Builder
	var sum float32
	unknown := 0
	valid := 0
	tooLow := 0
	tooHigh := 0
	for _, v := range values {
		if v == nil || !v.Trustworthy() {
			summary.WriteString("nil, ")
			unknown++
			continue
		}
		vcm := v.InCentimeters()
		if vcm < top {
			tooHigh++
			summary.WriteString(fmt.Sprintf("^%5.2f ", vcm))
			continue
		}
		if vcm > bottom {
			tooLow++
			summary.WriteString(fmt.Sprintf("v%5.2f ", vcm))
			continue
		}
		summary.WriteString(fmt.Sprintf("%5.2f ", vcm))
		sum += vcm
		valid++
	}
	if valid < 7 {
		return false, 0.0, fmt.Sprintf("too many uncertain values: %s", summary.String())
		// (this also avoids a divide-by-zero error below)
	}
	average := sum / float32(valid)
	summary.WriteString(fmt.Sprintf("avg of %d: %5.2f", valid, average))
	if tooHigh > 3 {
		return false, average, fmt.Sprintf("too many high values: %s", summary.String())
	}
	if tooLow > 3 {
		return false, average, fmt.Sprintf("too many low values: %s", summary.String())
	}
	return true, average, summary.String()
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
	var filling bool
	var target float32 = 9.0
	var fill float32 = 14.0
	for {
		time.Sleep(time.Second)

		ok, avg, status := safeToFill(5.0, 20.0, w.History())
		if !ok {
			fmt.Println("unsafe condition: ", status)
			filling = false
			valve.Out(gpio.Low) // Close the water valve
			continue
		}

		// go through the conditions from "top to bottom"
		if avg < target {
			fmt.Println("comfortably above target: ", status)
			filling = false
			valve.Out(gpio.Low) // Close the water valve
			continue
		}
		if filling {
			fmt.Println("filling: ", status)
			continue
		}
		if avg < fill {
			valve.Out(gpio.Low) // Close the water valve
			fmt.Println("just a little low: ", status)
			continue
		}
		if avg > fill {
			fmt.Println("starting fill: ", status)
			filling = true
			valve.Out(gpio.High) // Open the water valve
			continue
		}
		fmt.Println("unexpected condition: ", status)
		valve.Out(gpio.Low) // Close the water valve

	}
}
