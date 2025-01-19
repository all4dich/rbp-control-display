package main

import (
	"fmt"
	"log"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
)

const (
	RS = 25
	E  = 24
	D4 = 23
	D5 = 17
	D6 = 27
	D7 = 22
)

type LCD struct {
	rs, e, d4, d5, d6, d7 gpio.PinOut
}

func NewLCD() (*LCD, error) {
	if _, err := host.Init(); err != nil {
		return nil, err
	}

	lcd := &LCD{
		rs: gpioreg.ByName(fmt.Sprintf("GPIO%d", RS)),
		e:  gpioreg.ByName(fmt.Sprintf("GPIO%d", E)),
		d4: gpioreg.ByName(fmt.Sprintf("GPIO%d", D4)),
		d5: gpioreg.ByName(fmt.Sprintf("GPIO%d", D5)),
		d6: gpioreg.ByName(fmt.Sprintf("GPIO%d", D6)),
		d7: gpioreg.ByName(fmt.Sprintf("GPIO%d", D7)),
	}

	return lcd, nil
}

func (lcd *LCD) Initialize() {
	// Initialization sequence
	lcd.write4Bits(0x03)
	time.Sleep(5 * time.Millisecond)
	lcd.write4Bits(0x03)
	time.Sleep(100 * time.Microsecond)
	lcd.write4Bits(0x03)
	lcd.write4Bits(0x02)

	// Function set
	lcd.command(0x28) // 4-bit mode, 2 lines, 5x8 font
	lcd.command(0x0C) // Display on, cursor off, blink off
	lcd.command(0x06) // Increment cursor, no shift
	lcd.command(0x01) // Clear display
}

func (lcd *LCD) command(cmd byte) {
	lcd.rs.Out(gpio.Low)
	lcd.write4Bits(cmd >> 4)
	lcd.write4Bits(cmd & 0x0F)
}

func getLevel(data byte) gpio.Level {
	if data != 0 {
		return gpio.High
	} else {
		return gpio.Low
	}
}

func (lcd *LCD) write4Bits(data byte) {
	d4 := data & 0x01
	d5 := (data >> 1) & 0x01
	d6 := (data >> 2) & 0x01
	d7 := (data >> 3) & 0x01

	d4_level := getLevel(d4)
	d5_level := getLevel(d5)
	d6_level := getLevel(d6)
	d7_level := getLevel(d7)

	lcd.d4.Out(d4_level)
	lcd.d5.Out(d5_level)
	lcd.d6.Out(d6_level)
	lcd.d7.Out(d7_level)
	lcd.pulseEnable()
}

func (lcd *LCD) pulseEnable() {
	lcd.e.Out(gpio.High)
	time.Sleep(1 * time.Microsecond)
	lcd.e.Out(gpio.Low)
	time.Sleep(50 * time.Microsecond)
}

func (lcd *LCD) WriteString(s string) {
	for _, c := range s {
		lcd.rs.Out(gpio.High)
		lcd.write4Bits(byte(c) >> 4)
		lcd.write4Bits(byte(c) & 0x0F)
	}
}

func main() {
	lcd, err := NewLCD()
	if err != nil {
		log.Fatal(err)
	}

	lcd.Initialize()
	lcd.WriteString("Hello, World!")
}
