// This example connects to Access Point and prints some info
package main

import (
	"machine"
	"strconv"
	"time"

	"tinygo.org/x/drivers/wifinina"
)

// these are the default pins for the Arduino Nano33 IoT.
// change these to connect to a different UART or pins for the ESP8266/ESP32
var (

	// these are the default pins for the Arduino Nano33 IoT.
	spi = machine.NINA_SPI

	// this is the ESP chip that has the WIFININA firmware flashed on it
	adaptor *wifinina.Device
)

func main() {
	setup()
	connectToAP()

	for {
		println("----------------------------------------")
		printSSID()
		printRSSI()
		printMac()
		printIPs()
		printTime()
		time.Sleep(1 * time.Second)
	}

}

func setup() {
	// Configure SPI for 8Mhz, Mode 0, MSB First
	spi.Configure(machine.SPIConfig{
		Frequency: 8 * 1e6,
		SDO:       machine.NINA_SDO,
		SDI:       machine.NINA_SDI,
		SCK:       machine.NINA_SCK,
	})

	adaptor = wifinina.New(spi,
		machine.NINA_CS,
		machine.NINA_ACK,
		machine.NINA_GPIO0,
		machine.NINA_RESETN)
	adaptor.Configure()
}

// connect to access point
func connectToAP() {
	time.Sleep(2 * time.Second)
	println("Connecting to " + ssid)
	err := adaptor.ConnectToAccessPoint(ssid, pass, 10*time.Second)
	if err != nil {
		failMessage(err.Error())
	}

	println("Connected.")
}

func printSSID() {
	print("SSID: ")
	ssid, err := adaptor.GetCurrentSSID()
	if err != nil {
		println("Unknown (error: ", err.Error(), ")")
		return
	}
	println(ssid)
}

func printRSSI() {
	print("RSSI: ")
	rssi, err := adaptor.GetCurrentRSSI()
	if err != nil {
		println("Unknown (error: ", err.Error(), ")")
		return
	}
	println(strconv.Itoa(int(rssi)))
}

func printIPs() {
	ip, subnet, gateway, err := adaptor.GetIP()
	if err != nil {
		println("IP: Unknown (error: ", err.Error(), ")")
		return
	}
	println("IP: ", ip.String())
	println("Subnet: ", subnet.String())
	println("Gateway: ", gateway.String())
}

func printTime() {
	print("Time: ")
	t, err := adaptor.GetTime()
	for {
		if err != nil {
			println("Unknown (error: ", err.Error(), ")")
			return
		}
		if t != 0 {
			break
		}
		time.Sleep(time.Second)
		t, err = adaptor.GetTime()
	}
	println(time.Unix(int64(t), 0).String())
}

func printMac() {
	print("MAC: ")
	mac, err := adaptor.GetMACAddress()
	if err != nil {
		println("Unknown (", err.Error(), ")")
	}
	println(mac.String())
}

func failMessage(msg string) {
	for {
		println(msg)
		time.Sleep(1 * time.Second)
	}
}
