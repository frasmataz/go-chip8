package chip8

import (
	"fmt"
	"strings"
)

const width = 64
const height = 32

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
