package main

import (
	"bufio"
	"fmt"
	"strings"
	"time"

	"tinygo.org/x/drivers/ili9341"

	"tinygo.org/x/drivers/net/http"
	"tinygo.org/x/drivers/rtl8720dn"
	"tinygo.org/x/tinyterm"
)

var (
	display  *ili9341.Device
	terminal *tinyterm.Terminal
	adaptor  *rtl8720dn.RTL8720DN

	buf [0x400]byte
	err error
)

// change this to the URL that you want your HTTP request to go to
const url = "http://tinygo.org/"

func main() {
	initDisplay()

	terminalPrintln("Starting RTL8720DN adapter")

	adaptor, err = initAdaptor()
	if err != nil {
		failMessage(err)
	}
	http.SetBuf(buf[:])

	terminalPrintln("Connecting to access point")
	err = adaptor.ConnectToAccessPoint(ssid, pass, 10*time.Second)
	if err != nil {
		failMessage(err)
	}
	terminalPrintln("Connected")
	terminalPrintln("")

	displayIPAddress()

	cnt := 1
	for {
		makeHTTPRequest()

		fmt.Fprintf(terminal, "-------- %d --------\r\n", cnt)
		time.Sleep(10 * time.Second)
		cnt++
	}
}

func displayIPAddress() {
	ip, subnet, gateway, err := adaptor.GetIP()
	if err != nil {
		failMessage(err)
	}

	fmt.Fprintf(terminal, "IP Address : %s\r\n", ip)
	fmt.Fprintf(terminal, "Mask       : %s\r\n", subnet)
	fmt.Fprintf(terminal, "Gateway    : %s\r\n", gateway)
}

func makeHTTPRequest() {
	resp, err := http.Get(url)
	if err != nil {
		failMessage(err)
	}

	fmt.Fprintf(terminal, "%s %s\r\n", resp.Proto, resp.Status)
	for k, v := range resp.Header {
		fmt.Fprintf(terminal, "%s: %s\r\n", k, strings.Join(v, " "))
	}
	terminalPrintln("")

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		terminalPrintln(scanner.Text())
	}
	resp.Body.Close()
}
