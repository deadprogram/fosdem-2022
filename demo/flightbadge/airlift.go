package main

import (
	"machine"

	"tinygo.org/x/drivers/wifinina"
)

var (
	// interface for the AirLift WiFi Featherwing
	spi = machine.SPI0

	// ESP32 on the Featherwing with the WIFININA firmware flashed on it
	adaptor *wifinina.Device
)

func initAdaptor() *wifinina.Device {
	// Configure SPI for 8Mhz, Mode 0, MSB First
	spi.Configure(machine.SPIConfig{
		Frequency: 8 * 1e6,
		SDO:       machine.SPI0_SDO_PIN,
		SDI:       machine.SPI0_SDI_PIN,
		SCK:       machine.SPI0_SCK_PIN,
	})

	adaptor = wifinina.New(spi,
		machine.D13,
		machine.D11,
		machine.D10,
		machine.D12)
	adaptor.Configure()

	return adaptor
}
