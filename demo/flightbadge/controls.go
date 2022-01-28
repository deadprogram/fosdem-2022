package main

import (
	"machine"
	"time"

	"tinygo.org/x/drivers/shifter"
)

var (
	buttons        shifter.Device
	stickX, stickY machine.ADC

	xPos                                    uint16
	yPos                                    uint16
	b1push, b2push, b3push, b4push, joypush bool
	leftX, leftY, rightX, rightY            int
	droneconnected                          bool
)

func initControls() {
	machine.InitADC()

	buttons = shifter.New(shifter.EIGHT_BITS, machine.BUTTON_LATCH, machine.BUTTON_CLK, machine.BUTTON_OUT)
	buttons.Configure()

	// joystick
	stickX = machine.ADC{machine.D2}
	stickX.Configure(machine.ADCConfig{})

	stickY = machine.ADC{machine.D3}
	stickY.Configure(machine.ADCConfig{})
}

func readControls() {
	for {
		stickmode := "right"
		b1push = false
		b2push = false
		b3push = false
		b4push = false

		pressed, _ := buttons.Read8Input()
		if pressed&machine.BUTTON_LEFT_MASK > 0 {
			b3push = true
			stickmode = "left"
		}
		if pressed&machine.BUTTON_UP_MASK > 0 {
			b1push = true
		}
		if pressed&machine.BUTTON_DOWN_MASK > 0 {
			b2push = true
		}
		if pressed&machine.BUTTON_RIGHT_MASK > 0 {
			b4push = true
		}

		// read control stick
		xPos = stickX.Get() >> 6
		yPos = stickY.Get() >> 6
		if stickmode == "right" {
			// set left to center position
			leftX = center
			leftY = center

			// set right x,y to stick values
			rightX = int(xPos)
			rightY = int(yPos)
		} else {
			// set left x,y to stick values
			leftX = int(xPos)
			leftY = int(yPos)

			// set right to center position
			rightX = center
			rightY = center
		}

		time.Sleep(time.Millisecond * 50)
	}
}

type pair struct {
	x int
	y int
}

func getLeftStick() pair {
	s := pair{x: 0, y: 0}
	s.x = leftX
	s.y = leftY
	return s
}

func getRightStick() pair {
	s := pair{x: 0, y: 0}
	s.x = rightX
	s.y = rightY
	return s
}
