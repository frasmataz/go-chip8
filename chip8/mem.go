package chip8

import (
	"fmt"
)

type Memory struct {
	Memory [0xFFF]uint8
}

func (mem *Memory) Get(addr uint16) (uint8, error) {
	if addr > uint16(len(mem.Memory)) {
		return 0x0, fmt.Errorf("memory address out of bounds: %v, capacity %v", addr, len(mem.Memory))
	}

	return mem.Memory[addr], nil
}

func (mem *Memory) Set(addr uint16, val uint8) error {
	if addr > uint16(len(mem.Memory)) {
		return fmt.Errorf("memory address out of bounds: %v, capacity %v", addr, len(mem.Memory))
	}

	mem.Memory[addr] = val

	return nil
}

func NewMemory() *Memory {
	return new(Memory)
}
