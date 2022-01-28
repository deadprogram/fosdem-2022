package main

import (
	"image/color"
	"machine"
	"strconv"
	"time"

	"tinygo.org/x/drivers/st7735"
	"tinygo.org/x/tinydraw"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
)

var (
	display   st7735.Device
	failure   string
	realblack = color.RGBA{0, 0, 0, 255}
	black     = color.RGBA{127, 127, 127, 255}
)

const (
	noInfo = iota
	failureInfo
	connectingInfo
	flyingInfo
)

func initDisplay() {
	machine.SPI1.Configure(machine.SPIConfig{
		SCK:       machine.SPI1_SCK_PIN,
		SDI:       machine.SPI1_SDI_PIN,
		SDO:       machine.SPI1_SDO_PIN,
		Frequency: 8000000,
	})

	display = st7735.New(machine.SPI1, machine.TFT_RST, machine.TFT_DC, machine.TFT_CS, machine.TFT_LITE)
	display.Configure(st7735.Config{
		Rotation: st7735.ROTATION_90,
	})

	display.FillScreen(realblack)
}

func handleDisplay() {
	currentScreen := noInfo
	for {
		switch {
		case len(failure) > 0:
			if currentScreen != failureInfo {
				display.FillScreen(realblack)
				currentScreen = failureInfo
			}

			tinyfont.WriteLine(&display, &freemono.Regular9pt7b, 10, 20, "ERROR", black)
			tinyfont.WriteLine(&display, &freemono.Regular9pt7b, 10, 40, failure, black)

		case !droneconnected:
			if currentScreen != connectingInfo {
				display.FillScreen(realblack)
				currentScreen = connectingInfo
			}

			tinyfont.WriteLine(&display, &freemono.Bold9pt7b, 10, 20, "Connecting", black)
			tinyfont.WriteLine(&display, &freemono.Bold9pt7b, 10, 40, ssid, black)
			time.Sleep(100 * time.Millisecond)

		default:
			if currentScreen != flyingInfo {
				display.FillScreen(realblack)
				currentScreen = flyingInfo
			}

			display.FillRectangle(40, 5, 80, 40, realblack)

			x := strconv.Itoa(int(xPos))
			y := strconv.Itoa(int(yPos))
			msg := "x: " + x
			tinyfont.WriteLine(&display, &freemono.Bold9pt7b, 10, 20, msg, black)

			msg2 := "y: " + y
			tinyfont.WriteLine(&display, &freemono.Bold9pt7b, 10, 40, msg2, black)

			var radius int16 = 4
			if b1push {
				tinydraw.FilledCircle(&display, 16+32*0, 64-radius-1, radius, black)
			} else {
				tinydraw.FilledCircle(&display, 16+32*0, 64-radius-1, radius, realblack)
				tinydraw.Circle(&display, 16+32*0, 64-radius-1, radius, black)
			}
			if b2push {
				tinydraw.FilledCircle(&display, 16+32*1, 64-radius-1, radius, black)
			} else {
				tinydraw.FilledCircle(&display, 16+32*1, 64-radius-1, radius, realblack)
				tinydraw.Circle(&display, 16+32*1, 64-radius-1, radius, black)
			}
			if b3push {
				tinydraw.FilledCircle(&display, 16+32*2, 64-radius-1, radius, black)
			} else {
				tinydraw.FilledCircle(&display, 16+32*2, 64-radius-1, radius, realblack)
				tinydraw.Circle(&display, 16+32*2, 64-radius-1, radius, black)
			}
			if b4push {
				tinydraw.FilledCircle(&display, 16+32*3, 64-radius-1, radius, black)
			} else {
				tinydraw.FilledCircle(&display, 16+32*3, 64-radius-1, radius, realblack)
				tinydraw.Circle(&display, 16+32*3, 64-radius-1, radius, black)
			}
		}

		time.Sleep(100 * time.Millisecond)
	}
}
