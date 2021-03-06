# coffeewater
Tired of filling the water in your single-cup coffee machine?  Have a Benjamin to spare (~$100), and an afternoon to kill?  This is just what you need!

coffeewater is the code behind CoffeBot, a tiny computer and some plumbing parts to automate an annoying and repetitive task.

## Roughly what you need

You need the following basic components:

1. A single-cup coffee machine.
1. A computer with a few GPIO pins, such as a Raspberry Pi.
1. An HC-SR04 compatible ultrasonic range sensor.
1. A solenoid controlled water valve.
1. A few bits of plumbing and wire to connect everything together.

## How it works

CoffeeBot uses an ultrasonic range sensor mounted inside the top of the coffee machine's water reservoir to read the water level.  When the water level falls below a certain point, it triggers a solenoid valve to open and begin filling the reservoir.  To feed water into and out of the value, you must attach the water valve between a nearby sink's plumbing and the reservoir.  CoffeeBot's code includes several safety measures to ensure it won't accidentally flood your kitchen, and to ensure you never have to fill the reservior by hand again!

## Bill of materials

Here's the precise set of materials, with links, that you need to assemble a CoffeeBot of your very own.  Do take note that you may have some of these items already (such as a few resistors, or an old USB charger to power the Pi), or that you may be able to find them cheaper at a local plumbing supply store or big-box hardware store (such as a faucet supply line).  You may also be able to easily substitute some less expensive plastic plumbing parts, such as for the brass Sink Tee and Threaded Coupling I have linked to below.  You could cut costs more aggressively by dropping down to a Pi Zero and replacing the Perma-Proto hat with a spare piece of perfboard.

| Description                                                                 |   Cost |
| --------------------------------------------------------------------------- | ------:|
| [Rapsberry Pi Model 3 A+](https://www.adafruit.com/product/4027)            | $25.00 |
| [A cheap SDCard](https://smile.amazon.com/gp/product/B079GTYCW4)            |  $8.99 |
| [Rapsberry Pi case](https://www.adafruit.com/product/2361)                  |  $5.00 |
| [Perma-Proto Model A+ Hat](https://www.adafruit.com/product/2310)           |  $4.95 |
| [Rapsberry Pi power supply](https://www.adafruit.com/product/1995)          |  $7.50 |
| [HC-SR04 Ultrasonic Distance Sensor](https://www.adafruit.com/product/4007) |  $3.95 |
| [Solenoid Water Valve](https://www.adafruit.com/product/997)                |  $6.95 |
| [Water Valve Power Supply](https://www.adafruit.com/product/798)            |  $8.95 |
| [5.5x2.1mmm Barrel Power Supply Jack](https://www.digikey.com/product-detail/en/tensility-international-corp/54-00133/839-1516-ND/9685442)            |  $0.97 |
| [TIP120 Transistor](https://www.adafruit.com/product/976)                   |  $2.50 |
| [Kickback Protection Diode](https://www.adafruit.com/product/755)           |  $1.50 |
| [Faucet Supply Line](https://www.lowes.com/pd/Homewerks-Worldwide-3-8-in-Compression-12-in-Braided-Stainless-Steel-Faucet-Supply-Line/1000011602)               |  $5.22 |
| [Sink Tee](https://www.lowes.com/pd/B-K-3-8-in-Compression-Tee-Adapter-Fitting/1000505459)                                                                      |  $8.46 |
| [1/2" Threaded Coupling](https://www.lowes.com/pd/B-K-1-2-in-Threaded-Coupling-Fitting/1000505577)                                                              |  $6.88 |
| [1/2" MIP Adapter to 1/4" OD push fit adapter](https://www.lowes.com/pd/SharkBite-1-4-in-Push-to-Connect-x-1-2-in-Mip-dia-Male-Adapter-Push-Fitting/1000192601) |  $3.98 |
| [1/4" OD PEX water tubing and adapters](https://smile.amazon.com/gp/product/B07CRMDDYG)                                                                         | $14.87 |

### Supplemental

If you don't have a pile of resistors and a few optoisolators lying around, you might need to order some.  You should order a bunch of these if this kind of project sounds fun, because you'll inevitably find uses for these cheap components.

| Description        | Cost           |
| ------------- |-------------:|
| [an optoisolator](https://www.digikey.com/product-detail/en/taiwan-semiconductor-corporation/TPC817C-C9G/TPC817CC9G-ND/7359670) |  $0.39 |
| [some resistors](https://smile.amazon.com/Resistor-Assorted-Resistors-Assortment-Experiments/dp/B07L851T3V)                     | $14.99 |
| [some breakaway headers](https://www.adafruit.com/product/392)                                                                  |  $4.95 |
| [some female/female jumper wires](https://www.adafruit.com/product/1950)                                                        |  $1.95 |
| [~10 feet of >=4 conductor wire](https://www.lowes.com/pd/Southwire-18-4-Jacketed-Sprinkler-Wire-By-the-Foot/50142294)          |  $3.30 |

You could do without the headers if you just solder everything together, but it's much easier to work with on the future if you can disconnect each of the pieces.  You could also use roughly any wire with at least 4 conductors to hook up your HC-SR04, such as an old CAT5 patch cord or some alarm or thermostat wire.

## Perma-Proto Board Layout

TODO: Describe how to lay out the parts on the perma-proto board, and what each one does.

## Wiring assembly

TODO: Describe how to wire the 3 parts together.

## Software Setup

TODO: Describe how to setup coffeewater to be run by systemd.
