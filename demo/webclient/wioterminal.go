package main

import (
	"device/sam"
	"fmt"
	"image/color"
	"machine"
	"runtime/interrupt"
	"time"

	"tinygo.org/x/drivers/ili9341"
	"tinygo.org/x/drivers/net"
	"tinygo.org/x/drivers/rtl8720dn"
	"tinygo.org/x/tinyfont/proggy"
	"tinygo.org/x/tinyterm"
)

func init() {
	machine.SPI3.Configure(machine.SPIConfig{
		SCK:       machine.LCD_SCK_PIN,
		SDO:       machine.LCD_SDO_PIN,
		SDI:       machine.LCD_SDI_PIN,
		Frequency: 40000000,
	})
}

func initAdaptor() (*rtl8720dn.RTL8720DN, error) {
	machine.RTL8720D_CHIP_PU.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.RTL8720D_CHIP_PU.Low()
	time.Sleep(100 * time.Millisecond)
	machine.RTL8720D_CHIP_PU.High()
	time.Sleep(1000 * time.Millisecond)

	uart = UARTx{
		UART: &machine.UART{
			Buffer: machine.NewRingBuffer(),
			Bus:    sam.SERCOM0_USART_INT,
			SERCOM: 0,
		},
	}

	uart.Interrupt = interrupt.New(sam.IRQ_SERCOM0_2, handleInterrupt)
	uart.Configure(machine.UARTConfig{TX: machine.PB24, RX: machine.PC24, BaudRate: 614400})

	rtl := rtl8720dn.New(uart)
	_, err := rtl.Rpc_tcpip_adapter_init()
	if err != nil {
		return nil, err
	}

	net.UseDriver(rtl)
	return rtl, nil
}

var (
	black = color.RGBA{0, 0, 0, 255}
	white = color.RGBA{255, 255, 255, 255}
	red   = color.RGBA{255, 0, 0, 255}
	blue  = color.RGBA{0, 0, 255, 255}
	green = color.RGBA{0, 255, 0, 255}

	font = &proggy.TinySZ8pt7b
)

func initDisplay() {
	display = ili9341.NewSPI(
		machine.SPI3,
		machine.LCD_DC,
		machine.LCD_SS_PIN,
		machine.LCD_RESET,
	)
	display.Configure(ili9341.Config{})

	backlight.Configure(machine.PinConfig{machine.PinOutput})

	display.FillScreen(black)
	backlight.High()

	terminal = tinyterm.NewTerminal(display)
	terminal.Configure(&tinyterm.Config{
		Font:       font,
		FontHeight: 10,
		FontOffset: 6,
	})
}

func terminalPrintln(msg string) {
	fmt.Fprintf(terminal, "%s\r\n", msg)
}

func failMessage(err error) {
	for {
		terminalPrintln(err.Error())
		time.Sleep(time.Second)
	}
}

var (
	uart      UARTx
	backlight = machine.LCD_BACKLIGHT
)

func handleInterrupt(interrupt.Interrupt) {
	// should reset IRQ
	uart.Receive(byte((uart.Bus.DATA.Get() & 0xFF)))
	uart.Bus.INTFLAG.SetBits(sam.SERCOM_USART_INT_INTFLAG_RXC)
}

type UARTx struct {
	*machine.UART
}

func (u UARTx) Read(p []byte) (n int, err error) {
	if u.Buffered() == 0 {
		time.Sleep(1 * time.Millisecond)
		return 0, nil
	}
	return u.UART.Read(p)
}
