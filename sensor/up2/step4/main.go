package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/upboard/up2"
)

// Message displays a message on the attached RGB LCD.
func Message(msg string) {
	fmt.Println(msg)

	lcd.Clear()
	lcd.Home()
	lcd.Write(msg)
}

var board *up2.Adaptor

var gp *i2c.GrovePiDriver
var lcd *i2c.GroveLcdDriver

var button *gpio.ButtonDriver
var led *gpio.LedDriver

var dial *aio.GroveRotaryDriver

func main() {
	board = up2.NewAdaptor()

	gp = i2c.NewGrovePiDriver(board)
	lcd = i2c.NewGroveLcdDriver(board)

	button = gpio.NewButtonDriver(gp, "D8")
	led = gpio.NewLedDriver(gp, "D7")

	dial = aio.NewGroveRotaryDriver(gp, "A1", 500*time.Millisecond)

	work := func() {
		Message("ready")

		button.On(gpio.ButtonPush, func(data interface{}) {
			Message("On")

			led.On()
			lcd.SetRGB(0, 255, 0)
		})

		button.On(gpio.ButtonRelease, func(data interface{}) {
			Message("Off")

			led.Off()
			lcd.SetRGB(0, 0, 255)
		})

		dial.On(aio.Data, func(data interface{}) {
			msg := fmt.Sprint("Dial: ", data)
			Message(msg)
		})
	}

	robot := gobot.NewRobot("sensors",
		[]gobot.Connection{board},
		[]gobot.Device{gp, button, led, lcd, dial},
		work,
	)

	robot.Start()
}
