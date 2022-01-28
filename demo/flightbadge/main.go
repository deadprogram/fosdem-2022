// TinyGo flight stick using analog joystick and 5 buttons
// Outputs data via serial port in very simple space-delimited format
// End of each line of data has "CR-LF" aka 0x13 0x10
// Each update is sent every 50 ms.
package main

import (
	"time"

	tello "github.com/hybridgroup/tinygo-tello"
)

// Tello drone info here
const (
	ssid    = "TELLO-C48E59"
	pass    = ""
	speed   = 30
	center  = 660
	detente = 300
)

var (
	drone *tello.Tello
)

func main() {
	a := initAdaptor()
	drone = tello.New(a, "8888")

	initDisplay()
	go handleDisplay()

	initControls()
	go readControls()

	connectToAP(droneConnected)
}

func droneConnected() {
	println("Starting drone...")
	if err := drone.Start(); err != nil {
		failMessage(err.Error())
	}

	println("Drone started.")

	time.Sleep(1 * time.Second)

	println("Starting video...")
	if err := drone.StartVideo(); err != nil {
		failMessage(err.Error())
	}
	println("Video started.")

	droneconnected = true
	controlDrone()
}

// connect to drone wifi
func connectToAP(connectHandler func()) {
	var err error
	time.Sleep(2 * time.Second)
	for i := 0; i < 3; i++ {
		println("Connecting to " + ssid)
		err = adaptor.ConnectToAccessPoint(ssid, pass, 10*time.Second)
		if err != nil {
			println(err)
			time.Sleep(1 * time.Second)
			continue
		}

		// success
		println("Connected.")
		time.Sleep(3 * time.Second)
		if connectHandler != nil {
			connectHandler()
		}
		return
	}

	// couldn't connect to AP
	failMessage(err.Error())
}

func controlDrone() {
	startvid := true

	for {
		switch {
		case b1push:
			println("takeoff")
			err := drone.TakeOff()
			if err != nil {
				println(err)
			}

		case b2push:
			println("land")
			err := drone.Land()
			if err != nil {
				println(err)
			}
		}

		rightStick := getRightStick()
		switch {
		case rightStick.y+detente < center:
			drone.Backward(speed)
		case rightStick.y-detente > center:
			drone.Forward(speed)
		default:
			drone.Forward(0)
		}

		switch {
		case rightStick.x-detente > center:
			drone.Right(speed)
		case rightStick.x+detente < center:
			drone.Left(speed)
		default:
			drone.Right(0)
		}

		leftStick := getLeftStick()
		switch {
		case leftStick.y+detente < center:
			drone.Down(speed)
		case leftStick.y-detente > center:
			drone.Up(speed)
		default:
			drone.Up(0)
		}

		switch {
		case leftStick.x-detente > center:
			drone.Clockwise(speed)
		case leftStick.x+detente < center:
			drone.CounterClockwise(speed)
		default:
			drone.Clockwise(0)
		}

		if startvid {
			drone.StartVideo()
			startvid = false
		} else {
			startvid = true
		}

		time.Sleep(50 * time.Millisecond)
	}
}

func failMessage(msg string) {
	failure = msg
	for {
		println(msg)
		time.Sleep(1 * time.Second)
	}
}
