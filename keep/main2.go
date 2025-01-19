package main

import (
	"fmt"
	"log"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/conn/v3/lcd"
	"periph.io/x/host/v3"
	"time"
)

func main() {
	// Initialize periph.io
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Define GPIO pins
	pinRS := gpioreg.ByName("GPIO27")
	pinE := gpioreg.ByName("GPIO17")
	pinsData := []gpio.PinOut{
		gpioreg.ByName("GPIO25"),
		gpioreg.ByName("GPIO24"),
		gpioreg.ByName("GPIO23"),
		gpioreg.ByName("GPIO22"),
	}

	// Initialize the LCD
	l, err := lcd.NewHD44780(pinRS, pinE, pinsData, &lcd.HD44780Opts{
		NumRows: 2,
		NumCols: 16,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Clear the display
	if err := l.Clear(); err != nil {
		log.Fatal(err)
	}

	// Write a message
	if err := l.Write([]byte("Hello, World!")); err != nil {
		log.Fatal(err)
	}
	time.Sleep(2 * time.Second)
	if err := l.Clear(); err != nil {
		log.Fatal(err)
	}

	// Get current date and time
	now := time.Now()
	dateStr := now.Format("01/02/06")
	timeStr := now.Format("15:04:05")

	// Display date and time
	if err := l.SetCursor(0, 0); err != nil {
		log.Fatal(err)
	}
	if err := l.Write([]byte("Date: " + dateStr)); err != nil {
		log.Fatal(err)
	}
	if err := l.SetCursor(1, 0); err != nil {
		log.Fatal(err)
	}
	if err := l.Write([]byte("Time: " + timeStr)); err != nil {
		log.Fatal(err)
	}

	// Close the LCD
	if err := l.Halt(); err != nil {
		log.Fatal(err)
	}
}
