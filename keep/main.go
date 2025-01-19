package main

import (
	"fmt"
	"github.com/brian-armstrong/gpio"
	"github.com/mskrha/rpi-lcd"
	"time"
)

func main() {
	// Initialize GPIO
	err := gpio.Open()
	if err != nil {
		panic(err)
	}
	defer gpio.Close()

	// Initialize LCD
	lcd, err := rpi_lcd.New(
		gpio.NewPin(27), // RS
		gpio.NewPin(17), // E
		gpio.NewPin(25), // D4
		gpio.NewPin(24), // D5
		gpio.NewPin(23), // D6
		gpio.NewPin(22), // D7
	)
	if err != nil {
		panic(err)
	}
	defer lcd.Close()

	// Write a message
	lcd.Clear()
	lcd.WriteString("Hello, World!")
	lcd.Clear()

	// Display current date and time
	now := time.Now()
	dateStr := now.Format("%D")
	timeStr := now.Format("%T")
	lcd.SetCursor(0, 0)
	lcd.WriteString("Date: " + dateStr)
	lcd.SetCursor(1, 0)
	lcd.WriteString("Time: " + timeStr)
}
