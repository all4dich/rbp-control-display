# rbp-control-display

## Comments for main.go

The following comments provide an explanation of the code in `main.go`. These comments are written in GitHub Markdown format using Emacs org-mode syntax.

## Overview

The `main.go` file is a Go program designed to control an LCD display using GPIO pins. It initializes the LCD, displays the current date and time, and updates the display every second. The program uses the `periph.io` library for GPIO control and `argparse` for command-line argument parsing.

## Code Comments

```go
// Import necessary packages
import (
	"fmt"       // For formatted I/O
	"log"       // For logging errors
	"os"        // For command-line arguments
	"time"      // For time-related functions

	"github.com/akamensky/argparse" // For parsing command-line arguments
	"periph.io/x/conn/v3/gpio"     // For GPIO control
	"periph.io/x/conn/v3/gpio/gpioreg" // For GPIO pin registration
	"periph.io/x/host/v3"          // For initializing the host system
)
```

+ The program imports standard Go libraries and external libraries for GPIO control and argument parsing.

```go
// Define GPIO pin constants for the LCD
const (
	rsPin = "GPIO27" // Register Select Pin
	ePin  = "GPIO17" // Enable Pin
	d4Pin = "GPIO25" // Data 4 Pin
	d5Pin = "GPIO24" // Data 5 Pin
	d6Pin = "GPIO23" // Data 6 Pin
	d7Pin = "GPIO22" // Data 7 Pin
)
```

+ These constants define the GPIO pins used to control the LCD. Update these values based on your hardware configuration.

```go
// LCD struct represents an LCD object with GPIO pins for control and data
type LCD struct {
	rs gpio.PinIO // Register Select pin
	e  gpio.PinIO // Enable pin
	d4 gpio.PinIO // Data 4 pin
	d5 gpio.PinIO // Data 5 pin
	d6 gpio.PinIO // Data 6 pin
	d7 gpio.PinIO // Data 7 pin
}
```

+ The `LCD` struct encapsulates the GPIO pins required for controlling the LCD.

```go
// Declare global variables for timing parameters
var a, b int
```

+ The variables `a` and `b` are used to control the timing of the enable pin pulse. They can be customized via command-line arguments.

```go
// NewLCD creates a new LCD object with the specified pins
func NewLCD(rs, e, d4, d5, d6, d7 gpio.PinIO) *LCD {
	return &LCD{rs: rs, e: e, d4: d4, d5: d5, d6: d6, d7: d7}
}
```

+ The `NewLCD` function initializes an `LCD` object with the specified GPIO pins.

```go
// Init initializes the LCD in 4-bit mode
func (lcd *LCD) Init() {
	// Set all pins as outputs and initialize them to low
	lcd.rs.Out(gpio.Low)
	lcd.e.Out(gpio.Low)
	lcd.d4.Out(gpio.Low)
	lcd.d5.Out(gpio.Low)
	lcd.d6.Out(gpio.Low)
	lcd.d7.Out(gpio.Low)

	// Perform the LCD initialization sequence
	time.Sleep(100 * time.Millisecond)
	lcd.write4Bits(0x03)
	time.Sleep(5 * time.Millisecond)
	lcd.write4Bits(0x03)
	time.Sleep(1 * time.Millisecond)
	lcd.write4Bits(0x03)
	time.Sleep(1 * time.Millisecond)
	lcd.write4Bits(0x02)

	// Configure the LCD: 4-bit mode, 2 lines, 5x8 font
	lcd.sendCommand(0x28)
	// Turn on the display, disable cursor and blinking
	lcd.sendCommand(0x0C)
	// Clear the display
	lcd.sendCommand(0x01)
	// Set entry mode
	lcd.sendCommand(0x06)
	time.Sleep(2 * time.Millisecond)
}
```

+ The `Init` function configures the LCD for 4-bit mode and sets its initial state.

```go
// write4Bits sends 4 bits of data to the LCD
func (lcd *LCD) write4Bits(value byte) {
	// Set the data pins based on the 4 bits of the value
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

	// Pulse the enable pin to latch the data
	lcd.pulseEnable()
}
```

+ The `write4Bits` function sends 4 bits of data to the LCD and latches it using the enable pin.

```go
// control_lcd handles the main logic for controlling the LCD
func control_lcd() {
	// Parse command-line arguments for timing parameters
	parser := argparse.NewParser("LCD Display Checker with Date/Time", "LCD Display Checker with Date/Time")
	aPtr := parser.Int("a", "first", &argparse.Options{Required: false, Help: "first number", Default: 1})
	bPtr := parser.Int("b", "second", &argparse.Options{Required: false, Help: "second number", Default: 1})
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Println(err)
	}
	a = *aPtr
	b = *bPtr

	// Initialize the LCD
	rs := gpioreg.ByName(rsPin)
	e := gpioreg.ByName(ePin)
	d4 := gpioreg.ByName(d4Pin)
	d5 := gpioreg.ByName(d5Pin)
	d6 := gpioreg.ByName(d6Pin)
	d7 := gpioreg.ByName(d7Pin)

	if rs == nil || e == nil || d4 == nil || d5 == nil || d6 == nil || d7 == nil {
		log.Fatal("Failed to find one or more pins")
	}

	lcd := NewLCD(rs, e, d4, d5, d6, d7)
	lcd.Init()

	// Display the current date and time
	for {
		now := time.Now()
		curr_date := now.Format("2006-01-02")
		curr_time := now.Format("15:04:05")

		lcd.clear()
		lcd.writeString(curr_date)
		lcd.setCursor(1, 0)
		lcd.writeString(curr_time)

		time.Sleep(1 * time.Second)
	}
}
```

+ The `control_lcd` function initializes the LCD and continuously updates it with the current date and time.

## Notes

+ Ensure that the GPIO pins are correctly connected to the LCD before running the program.
+ The `periph.io` library must be installed and configured for GPIO control.

