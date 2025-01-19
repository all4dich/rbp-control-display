package main

import (
	"fmt"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

// LCD pin connections (adjust these to your actual wiring)
const (
	RS = 17 // Register Select
	E  = 27 // Enable
	D4 = 22 // Data 4
	D5 = 23 // Data 5
	D6 = 24 // Data 6
	D7 = 25 // Data 7
)

func main() {
	// Open and defer closing of GPIO
	err := rpio.Open()
	if err != nil {
		panic(fmt.Sprintf("failed to open rpio: %v", err))
	}
	defer rpio.Close()

	// Set pin modes as output
	rpio.PinMode(RS, rpio.Output)
	rpio.PinMode(E, rpio.Output)
	rpio.PinMode(D4, rpio.Output)
	rpio.PinMode(D5, rpio.Output)
	rpio.PinMode(D6, rpio.Output)
	rpio.PinMode(D7, rpio.Output)

	// Initialize LCD
	lcdInit()

	// Display some text
	lcdWriteString("Hello, World!")

	// Move cursor to second line (assuming 16x2 display)
	lcdCmd(0xC0) // Set DDRAM address to 0x40 (second line)
	lcdWriteString("from Go and RPi!")

	for { // Keep the display on
		time.Sleep(time.Second)
	}
}

func lcdPulseEnable() {
	rpio.WritePin(E, rpio.High)
	time.Sleep(time.Microsecond)
	rpio.WritePin(E, rpio.Low)
	time.Sleep(time.Microsecond * 50) // Allow some settling time
}

func lcdWrite4Bits(data byte) {
	rpio.WritePin(D4, data&0x01)
	rpio.WritePin(D5, (data>>1)&0x01)
	rpio.WritePin(D6, (data>>2)&0x01)
	rpio.WritePin(D7, (data>>3)&0x01)

	lcdPulseEnable()
}

func lcdCmd(cmd byte) {
	rpio.WritePin(RS, rpio.Low)
	lcdWrite4Bits(cmd >> 4)
	lcdWrite4Bits(cmd)
}

func lcdChar(char byte) {
	rpio.WritePin(RS, rpio.High)
	lcdWrite4Bits(char >> 4)
	lcdWrite4Bits(char)
}

func lcdWriteString(str string) {
	for _, char := range str {
		lcdChar(byte(char))
	}
}

func lcdInit() {
	// Initialization sequence for 4-bit mode
	time.Sleep(time.Millisecond * 20) // Wait for LCD to power up

	// These steps are crucial for 4-bit mode initialization.
	// The datasheet specifies sending specific nibbles, not bytes.

	lcdWrite4Bits(0x03)
	time.Sleep(time.Millisecond * 5)
	lcdWrite4Bits(0x03)
	time.Sleep(time.Millisecond * 1) // >100us
	lcdWrite4Bits(0x03)
	time.Sleep(time.Microsecond * 100)

	lcdWrite4Bits(0x02)                // 4-bit mode now set
	time.Sleep(time.Microsecond * 100) // >40us

	lcdCmd(0x28) // 2 lines, 5x8 font
	lcdCmd(0x0C) // Display on, cursor off
	lcdCmd(0x06) // Increment cursor
	lcdCmd(0x01) // Clear display
	time.Sleep(time.Millisecond * 2)
}
