package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/akamensky/argparse"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
)

const (
	rsPin = "GPIO27" // Register Select Pin
	ePin  = "GPIO17" // Enable Pin
	d4Pin = "GPIO25" // Data 4 Pin
	d5Pin = "GPIO24" // Data 5 Pin
	d6Pin = "GPIO23" // Data 6 Pin
	d7Pin = "GPIO22" // Data 7 Pin
)

// LCD represents an LCD object
type LCD struct {
	rs gpio.PinIO
	e  gpio.PinIO
	d4 gpio.PinIO
	d5 gpio.PinIO
	d6 gpio.PinIO
	d7 gpio.PinIO
}

// Declare a, b as integer with default value 1
var a, b int

// NewLCD creates a new LCD object with the specified pins
func NewLCD(rs, e, d4, d5, d6, d7 gpio.PinIO) *LCD {
	return &LCD{rs: rs, e: e, d4: d4, d5: d5, d6: d6, d7: d7}
}

// Init initializes the LCD
func (lcd *LCD) Init() {
	// Set pins as outputs
	lcd.rs.Out(gpio.Low)
	lcd.e.Out(gpio.Low)
	lcd.d4.Out(gpio.Low)
	lcd.d5.Out(gpio.Low)
	lcd.d6.Out(gpio.Low)
	lcd.d7.Out(gpio.Low)

	time.Sleep(100 * time.Millisecond)
	lcd.write4Bits(0x03)
	time.Sleep(5 * time.Millisecond)
	lcd.write4Bits(0x03)
	time.Sleep(1 * time.Millisecond)
	lcd.write4Bits(0x03)
	time.Sleep(1 * time.Millisecond)
	lcd.write4Bits(0x02)

	// Function Set: 4-bit mode, 2 lines, 5x8 font
	lcd.sendCommand(0x28)
	// Display ON/OFF Control: Display on, Cursor off, Blink off
	lcd.sendCommand(0x0C)
	// Display Clear
	lcd.sendCommand(0x01)
	// Entry Mode Set
	lcd.sendCommand(0x06)
	time.Sleep(2 * time.Millisecond)
}

func (lcd *LCD) write4Bits(value byte) {
	// Set data pins based on the 4 bits
	if (value & 0x01) != 0 {
		lcd.d4.Out(gpio.High)
	} else {
		lcd.d4.Out(gpio.Low)
	}
	if (value & 0x02) != 0 {
		lcd.d5.Out(gpio.High)
	} else {
		lcd.d5.Out(gpio.Low)
	}
	if (value & 0x04) != 0 {
		lcd.d6.Out(gpio.High)
	} else {
		lcd.d6.Out(gpio.Low)
	}
	if (value & 0x08) != 0 {
		lcd.d7.Out(gpio.High)
	} else {
		lcd.d7.Out(gpio.Low)
	}

	// Pulse the enable pin
	lcd.pulseEnable()

}
func (lcd *LCD) pulseEnable() {
	lcd.e.Out(gpio.High)
	time.Sleep(time.Duration(a) * time.Microsecond)
	lcd.e.Out(gpio.Low)
	time.Sleep(time.Duration(b) * time.Microsecond)
}

// Send a command to the LCD
func (lcd *LCD) sendCommand(command byte) {
	lcd.rs.Out(gpio.Low) // RS low for command
	lcd.write4Bits(command >> 4)
	lcd.write4Bits(command & 0x0F)
}

// Send a data byte (character) to the LCD
func (lcd *LCD) sendData(data byte) {
	lcd.rs.Out(gpio.High) // RS high for data
	lcd.write4Bits(data >> 4)
	lcd.write4Bits(data & 0x0F)
}

// Clear the LCD screen
func (lcd *LCD) clear() {
	lcd.sendCommand(0x01)
	time.Sleep(2 * time.Millisecond)
}

// Write a string to the LCD
func (lcd *LCD) writeString(str string) {
	for _, char := range str {
		lcd.sendData(byte(char))
	}
}

// Set cursor position to specified row and column (0 indexed)
func (lcd *LCD) setCursor(row, col int) {
	var addr byte
	if row == 0 {
		addr = byte(0x00 + col)
	} else {
		addr = byte(0x40 + col)
	}

	lcd.sendCommand(0x80 | addr)
}

func main() {
	// Initialize host driver
	parser := argparse.NewParser("LCD Display Checker with Date/Time", "LCD Display Checker with Date/Time")
	aPtr := parser.Int("a", "first", &argparse.Options{Required: false, Help: "first number", Default: 1})
	bPtr := parser.Int("b", "second", &argparse.Options{Required: false, Help: "first number", Default: 1})
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Println(err)
	}
	a = *aPtr
	b = *bPtr
	fmt.Println(a)
	fmt.Println(b)
	// Get pins by their names
	rs := gpioreg.ByName(rsPin)
	e := gpioreg.ByName(ePin)
	d4 := gpioreg.ByName(d4Pin)
	d5 := gpioreg.ByName(d5Pin)
	d6 := gpioreg.ByName(d6Pin)
	d7 := gpioreg.ByName(d7Pin)

	if rs == nil || e == nil || d4 == nil || d5 == nil || d6 == nil || d7 == nil {
		log.Fatal("Failed to find one or more pins")
	}

	// Create LCD object
	lcd := NewLCD(rs, e, d4, d5, d6, d7)

	// Initialize the LCD
	lcd.Init()
	fmt.Println("LCD initialized.")

	// Create while loop
	for {
		// Get current date and time
		now := time.Now()
		curr_date := now.Format("2006-01-02")
		curr_time := now.Format("15:04:05")

		lcd.clear()
		lcd.writeString(curr_date)
		lcd.setCursor(1, 0)
		lcd.writeString(curr_time)

		time.Sleep(1 * time.Second)
	}
	//time.Sleep(10 * time.Second)
	//lcd.clear()
	//lcd.writeString("Done")
	//time.Sleep(2 * time.Second)

}
