package main

import (
	"github.com/pimvanhespen/go-pi-lcd1602"
	"github.com/pimvanhespen/go-pi-lcd1602/synchronized"
	"time"
)

func main() {
	// Initialize the LCD with default pin configuration
	lcd := lcd1602.New(27, 17, []int{25, 24, 23, 22}, 16)
	syncLcd := synchronized.NewSynchronizedLCD(lcd)

	// Initialize the LCD
	syncLcd.Initialize()
	defer syncLcd.Close()

	// Write messages to the LCD
	syncLcd.WriteLines("Hello, Raspberry", "Pi 5 with Go!")

	// Keep the message displayed for 10 seconds
	time.Sleep(10 * time.Second)
}
