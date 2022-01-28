// This demo creates an MQTT connection that publishes a message every second
// to an MQTT broker.
//
package main

import (
	"fmt"
	"machine"
	"math/rand"
	"time"

	"tinygo.org/x/drivers/net/mqtt"
	"tinygo.org/x/drivers/wifinina"

	"tinygo.org/x/drivers/ili9341"
	"tinygo.org/x/tinyterm"
)

// MQTT broker to use. Replace with your own info.
const server = "tcp://test.mosquitto.org:1883"

var (
	display  *ili9341.Device
	terminal *tinyterm.Terminal
	adaptor  *wifinina.Device

	cl      mqtt.Client
	topicTx = "tinygo/tx"
	topicRx = "tinygo/rx"
)

func main() {
	initDisplay()

	time.Sleep(3000 * time.Millisecond)
	rand.Seed(time.Now().UnixNano())

	initAdaptor()

	connectToAP()
	connectToMQTT()

	select {}
}

func initAdaptor() {
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
	terminalPrintln("Connecting to " + ssid)
	err := adaptor.ConnectToAccessPoint(ssid, pass, 10*time.Second)
	if err != nil {
		failMessage(err.Error())
	}

	terminalPrintln("Connected.")

	time.Sleep(2 * time.Second)
	ip, _, _, err := adaptor.GetIP()
	for ; err != nil; ip, _, _, err = adaptor.GetIP() {
		failMessage(err.Error())
	}
	terminalPrintln(ip.String())
}

func connectToMQTT() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(server).SetClientID("tinygo-client-" + randomString(10))

	terminalPrintln("Connecting to MQTT broker at")
	terminalPrintln(server)

	cl = mqtt.NewClient(opts)
	if token := cl.Connect(); token.Wait() && token.Error() != nil {
		failMessage(token.Error().Error())
	}

	// subscribe
	token := cl.Subscribe(topicRx, 0, subHandler)
	token.Wait()
	if token.Error() != nil {
		failMessage(token.Error().Error())
	}

	go publishing()
}

func publishing() {
	for i := 0; ; i++ {
		terminalPrintln("Publishing MQTT message...")
		terminalPrintln("")

		data := []byte(fmt.Sprintf(`{"e":[{"n":"hello %d","v":101}]}`, i))
		token := cl.Publish(topicRx, 0, false, data)
		token.Wait()
		if token.Error() != nil {
			terminalPrintln(token.Error().Error())
		}

		time.Sleep(5 * time.Second)
	}
}

func subHandler(client mqtt.Client, msg mqtt.Message) {
	terminalPrintln("Message received on " + msg.Topic())
	terminalPrintln(string(msg.Payload()))
	terminalPrintln("")
}

// Returns an int >= min, < max
func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

// Generate a random string of A-Z chars with len = l
func randomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(randomInt(65, 90))
	}
	return string(bytes)
}
