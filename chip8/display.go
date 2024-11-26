package chip8

import (
	"fmt"
	"strings"
)

const width = 64
const height = 32

var CharSprites = [16][5]uint8{
	{ // 0
		0xF0,
		0x90,
		0x90,
		0x90,
		0xF0,
	},
	{ // 1
		0x20,
		0x60,
		0x20,
		0x20,
		0x70,
	},
	{ // 2
		0xF0,
		0x10,
		0xF0,
		0x80,
		0xF0,
	},
	{ // 3
		0xF0,
		0x10,
		0xF0,
		0x10,
		0xF0,
	},
	{ // 4
		0x90,
		0x90,
		0xF0,
		0x10,
		0x10,
	},
	{ // 5
		0xF0,
		0x80,
		0xF0,
		0x10,
		0xF0,
	},
	{ // 6
		0xF0,
		0x80,
		0xF0,
		0x90,
		0xF0,
	},
	{ // 7
		0xF0,
		0x10,
		0x20,
		0x40,
		0x40,
	},
	{ // 8
		0xF0,
		0x90,
		0xF0,
		0x90,
		0xF0,
	},
	{ // 9
		0xF0,
		0x90,
		0xF0,
		0x10,
		0xF0,
	},
	{ // A
		0xF0,
		0x90,
		0xF0,
		0x90,
		0x90,
	},
	{ // B
		0xE0,
		0x90,
		0xE0,
		0x90,
		0xE0,
	},
	{ // C
		0xF0,
		0x80,
		0x80,
		0x80,
		0xF0,
	},
	{ // D
		0xE0,
		0x90,
		0x90,
		0x90,
		0xE0,
	},
	{ // E
		0xF0,
		0x80,
		0xF0,
		0x80,
		0xF0,
	},
	{ // F
		0xF0,
		0x80,
		0xF0,
		0x80,
		0x80,
	},
}

type Display struct {
	// Structure is [y][x] - makes row-by-row looping easier
	framebuffer [height][width]bool
}

func (display *Display) Set(x uint, y uint, val bool) error {
	if x >= width || y >= height {
		return fmt.Errorf("pixel coordinate out of range: x: %v, y: %v", x, y)
	}

	display.framebuffer[y][x] = val

	return nil
}

func (display *Display) Get(x uint, y uint) (bool, error) {
	if x >= width || y >= height {
		return false, fmt.Errorf("pixel coordinate out of range: x: %v, y: %v", x, y)
	}

	return display.framebuffer[y][x], nil
}

func (display *Display) PrintFrame() string {
	var sb strings.Builder

	sb.WriteRune('\n')

	for y, row := range display.framebuffer {
		for x := range row {
			if display.framebuffer[y][x] {
				sb.WriteString("██")
			} else {
				sb.WriteString("░░")
			}
		}

		sb.WriteRune('\n')
	}

	return sb.String()
}

func NewDisplay() *Display {
	return new(Display)
}
