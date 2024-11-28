package chip8

import (
	"fmt"
	"strings"
)

const MemorySize = 0x1000

type Memory struct {
	Memory [MemorySize]uint8
}

func NewMemory() *Memory {
	mem := new(Memory)

	// Load default char sprites into 'interpreter area' (0x000 - 0x1FF) of memory
	for ci, char := range CharSprites {
		for bi, _byte := range char {
			mem.Set8(uint16(ci*5+bi), _byte)
		}
	}

	return mem
}

func (mem *Memory) Get8(addr uint16) (uint8, error) {
	if addr > uint16(len(mem.Memory)-1) {
		return 0x0, fmt.Errorf("memory address out of bounds: %v, capacity %v", addr, len(mem.Memory))
	}

	return mem.Memory[addr], nil
}

func (mem *Memory) Get16(addr uint16) (uint16, error) {
	if addr > uint16(len(mem.Memory)-2) {
		return 0x0, fmt.Errorf("memory address out of bounds: %v, capacity %v", addr, len(mem.Memory))
	}

	return uint16(mem.Memory[addr])<<8 | uint16(mem.Memory[addr+1]), nil
}

func (mem *Memory) Set8(addr uint16, val uint8) error {
	if addr > uint16(len(mem.Memory)-1) {
		return fmt.Errorf("memory address out of bounds: %v, capacity %v", addr, len(mem.Memory))
	}

	mem.Memory[addr] = val

	return nil
}

func (mem *Memory) Set16(addr uint16, val uint16) error {
	if addr > uint16(len(mem.Memory)-2) {
		return fmt.Errorf("memory address out of bounds: %v, capacity %v", addr, len(mem.Memory))
	}

	mem.Memory[addr] = uint8(val & 0xFF00 >> 8)
	mem.Memory[addr+1] = uint8(val & 0x00FF)

	return nil
}

func (mem *Memory) GetPrettyMemoryState() string {
	const columns = 16

	// Pre-compute output string size to minimize memory copying while building string
	// 6 chars per row adddress  (⏎xxx:␣) + 3 chars for each byte (xx␣)
	//outputSize := (len(mem.Memory) * 3) + ((len(mem.Memory) / columns) * 6)

	var sb strings.Builder
	sb.Grow(1)

	for addr, val := range mem.Memory {
		if addr%columns == 0 {
			sb.WriteString(fmt.Sprintf("\n%03x: ", addr))
		}

		sb.WriteString(fmt.Sprintf("%02x ", val))
	}

	return sb.String()
}
