// closevalve simply sets a GPIO pin to low and exits.
package main

import (
	"flag"
	"fmt"
	"os"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
)

var (
	valvePin = flag.String("valvePin", "5", "GPIO pin connected the water solneoid valve.")
)

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
}
